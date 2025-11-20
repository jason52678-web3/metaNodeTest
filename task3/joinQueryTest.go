package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Id        uint      `gorm:"primaryKey" `
	Username  string    `gorm:"size:50;not null;uniqueIndex"`
	Email     string    `gorm:"size:255;not null;uniqueIndex"`
	Password  string    `gorm:"size:255;not null"`
	Nickname  string    `gorm:"size:50;not null"`
	Avatar    string    `gorm:"size:500;not null"`
	Bio       string    `gorm:"size:500;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	DeletedAt time.Time `gorm:"not null"`

	Posts []Post `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Post struct {
	Id        uint      `gorm:"primaryKey" `
	Title     string    `gorm:"size:200;not null"`
	Content   string    `gorm:"type:longtext;not null"`
	Slug      string    `gorm:"size:255;not null"`
	Status    string    `gorm:"size:20;default:draft"` //draft/published/archived
	ViewCount int       `gorm:"default:0"`
	UserId    uint      `gorm:"not null;index"`
	User      User      `gorm:"foreignKey:UserId"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	DeletedAt time.Time `gorm:"not null"`

	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

type Comment struct {
	Id        uint      `gorm:"primaryKey" `
	Content   string    `gorm:"type:longtext;not null"`
	Author    string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:255;not null"`
	Website   string    `gorm:"size:500;not null"`
	PostId    uint      `gorm:"not null;index"`
	Post      Post      `gorm:"foreignKey:PostId"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	DeletedAt time.Time `gorm:"index"`
}

type PostWithCount struct {
	Post
	CommentCount int `gorm:"column:comment_count"`
}

func GetPostMostComments(db *gorm.DB) (*PostWithCount, error) {
	var result PostWithCount
	//err := db.Model(&Post{}).
	err := db.Table("posts").
		Select("posts.*, COUNT(comments.id) AS comment_count").
		Joins("JOIN comments on posts.id = comments.post_id").
		Group("posts.id").
		Order("comment_count desc").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		fmt.Println("GetPostMostComments error:", err)
		return nil, err
	}
	return &result, nil
}

func GetUserPostWithCommentsByid(ctx context.Context, db *gorm.DB, userId uint) (*User, error) {
	var user User

	err := db.WithContext(ctx).Preload("Posts.Comments").First(&user, userId).Error
	if err != nil {
		fmt.Println("get infor error", err)
		return nil, err
	}
	return &user, nil
}

func main() {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := user + ":" + pass + "@tcp(localhost:3306)/" + dbname + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("mysql open err:", err)
		return
	}

	//err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	//if err != nil {
	//	fmt.Println("mysql auto migrate err:", err)
	//	return
	//}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//1.使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
	var targetUserId uint = 1
	userInfo, err := GetUserPostWithCommentsByid(ctx, db, targetUserId)
	if err != nil {
		fmt.Println("mysql get user post err:", err)
		return
	}

	fmt.Printf("UserID:%d has published %d Posts", targetUserId, len(userInfo.Posts))
	for _, post := range userInfo.Posts {
		fmt.Printf("\rTitle:%s comments numbers:%d \n", post.Id, len(post.Comments))
		for _, comment := range post.Comments {
			fmt.Printf("\rAuthor:%s  Conten:%s Time:%s\n",
				comment.Author,
				comment.Content,
				comment.CreatedAt.Format(time.RFC3339),
				//comment.CreatedAt.Format("2006-01-02 15:04:05"),
			)
		}
	}

	//2.使用Gorm查询评论数量最多的文章信息。
	result, err := GetPostMostComments(db)
	if err != nil {
		fmt.Println("mysql get post err:", err)
		return
	}
	fmt.Printf("Title:%s comments：%d\n", result.Title, result.CommentCount)

}
