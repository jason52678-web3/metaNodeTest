package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/gorm/schema"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 开关：true 使用软删除（deleted_at），false 使用硬删除
const UseSoftDelete = true

var (
	dsnTpl = "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	db     *gorm.DB
)

// 状态常量
const (
	StatusNoComment  = "无评论"
	StatusHasComment = "已评论"
)

// User 用户
type User struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"size:50;not null"`
	PostsCount int64  `gorm:"default:0;not null"` // 文章数
	Posts      []Post `gorm:"foreignKey:UserID"`
}

func (User) TableName() string { return "user" }

// Post 文章
type Post struct {
	ID            uint      `gorm:"primaryKey"`
	Title         string    `gorm:"size:200;not null"`
	Content       string    `gorm:"type:longtext;not null"`
	UserID        uint      `gorm:"not null;index"`
	CommentStatus string    `gorm:"size:20;default:已评论;index"` // 已评论 | 无评论
	Comments      []Comment `gorm:"foreignKey:PostID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (Post) TableName() string { return "posts" }

// Comment 评论
type Comment struct {
	ID        uint           `gorm:"primaryKey"`
	Content   string         `gorm:"type:text;not null"`
	PostID    uint           `gorm:"not null;index"`
	Post      *Post          `gorm:"foreignKey:PostID"`
	DeletedAt gorm.DeletedAt `gorm:"index"` // 软删除字段
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Comment) TableName() string { return "comments" }

// 初始化数据库连接
func initDB() error {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	name := os.Getenv("DB_NAME")

	var err error
	dsn := fmt.Sprintf(dsnTpl, user, pass, host, port, name)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 关闭复数表名
		},
	})
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql db: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移
	if err = db.AutoMigrate(&User{}, &Post{}, &Comment{}); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}

// 清理与初始化测试数据
func resetData() error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM comments").Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM posts").Error; err != nil {
			return err
		}
		return tx.Exec("UPDATE user SET posts_count = 0").Error
	})
}

// 文章创建后：原子自增作者的文章数
func (p *Post) AfterCreate(tx *gorm.DB) error {
	return tx.Exec("UPDATE user SET posts_count = posts_count + 1 WHERE id = ?", p.UserID).Error
}

// 评论创建后：原子更新文章评论状态
func (c *Comment) AfterCreate(tx *gorm.DB) error {
	return updatePostCommentStatus(tx, c.PostID)
}

// 评论删除后：原子更新文章评论状态
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	return updatePostCommentStatus(tx, c.PostID)
}

// 原子更新文章评论状态：根据是否启用软删除选择计数条件
func updatePostCommentStatus(tx *gorm.DB, postID uint) error {
	var countCond string
	if UseSoftDelete {
		countCond = "deleted_at IS NULL" //没有标记为删除的评论
	} else {
		countCond = "1=1"
	}
	return tx.Exec(`
		UPDATE posts 
		SET comment_status = CASE 
			WHEN (SELECT COUNT(1) FROM comments WHERE post_id = ? AND `+countCond+`) = 0 
			THEN ? ELSE ? 
		END 
		WHERE id = ?`,
		postID, StatusNoComment, StatusHasComment, postID).Error
}

// 读取环境变量或默认值
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	if err := initDB(); err != nil {
		log.Fatalf("init db failed: %v", err)
	}
	defer func() {
		if sqlDB, _ := db.DB(); sqlDB != nil {
			_ = sqlDB.Close()
		}
	}()

	if err := resetData(); err != nil {
		log.Fatalf("reset data failed: %v", err)
	}

	// 测试1：创建文章后，作者的 PostsCount 是否 +1
	var u User
	if err := db.Create(&User{Name: "Alice"}).Error; err != nil {
		log.Fatalf("create user failed: %v", err)
	}
	if err := db.First(&u, "name = ?", "Alice").Error; err != nil {
		log.Fatalf("find user failed: %v", err)
	}

	post := Post{Title: "第一篇文章", Content: "内容", UserID: u.ID}
	if err := db.Create(&post).Error; err != nil {
		log.Fatalf("create post failed: %v", err)
	}

	var u2 User
	if err := db.First(&u2, u.ID).Error; err != nil {
		log.Fatalf("find user2 failed: %v", err)
	}
	if u2.PostsCount != 1 {
		log.Fatalf("期望 PostsCount=1，实际=%d", u2.PostsCount)
	}
	fmt.Printf("[Test1] 创建文章后 PostsCount 校验通过：%d\n", u2.PostsCount)

	// 测试2：删除评论后，评论数为 0 时更新文章状态为“无评论”
	// 2.1 新增一条评论
	cmt := Comment{Content: "第一条评论", PostID: post.ID}
	if err := db.Create(&cmt).Error; err != nil {
		log.Fatalf("create comment failed: %v", err)
	}

	// 删除这条评论
	if err := db.Delete(&cmt).Error; err != nil {
		log.Fatalf("delete comment failed: %v", err)
	}

	var p Post
	if err := db.First(&p, post.ID).Error; err != nil {
		log.Fatalf("find post failed: %v", err)
	}
	if p.CommentStatus != StatusNoComment {
		log.Fatalf("期望 CommentStatus=%s，实际=%s", StatusNoComment, p.CommentStatus)
	}
	fmt.Printf("[Test2-1] 删除唯一评论后状态校验通过：%s\n", p.CommentStatus)

	// 2.2 再新增两条评论，再删除一条，确保状态仍为“已评论”
	//db.Create(&Comment{Content: "<UNK>", PostID: post.ID}).Error
	if err := db.Create(&Comment{Content: "第二条评论", PostID: post.ID}).Error; err != nil {
		log.Fatalf("create comment2 failed: %v", err)
	}
	if err := db.Create(&Comment{Content: "第三条评论", PostID: post.ID}).Error; err != nil {
		log.Fatalf("create comment3 failed: %v", err)
	}
	// 删除一条
	if err := db.Delete(&Comment{}, "content = ?", "第二条评论").Error; err != nil {
		log.Fatalf("delete comment2 failed: %v", err)
	}

	var p2 Post
	if err := db.First(&p2, post.ID).Error; err != nil {
		log.Fatalf("find post2 failed: %v", err)
	}
	if p2.CommentStatus != StatusHasComment {
		log.Fatalf("期望 CommentStatus=%s，实际=%s", StatusHasComment, p2.CommentStatus)
	}
	fmt.Printf("[Test2-2] 删除后仍有评论，状态校验通过：%s\n", p2.CommentStatus)

	fmt.Println("所有测试通过")
}
