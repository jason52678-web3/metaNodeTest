package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github/task4/MyBlog/models"
	"go.uber.org/zap"
	"time"
)

var (
	ErrorPostNotExist = errors.New("帖子不存在")
)

func GetPostsList() (postsList []*models.Posts, err error) {
	sqlStr := "select id, title, user_id, created_at from posts"

	postsList = make([]*models.Posts, 0, 2)
	if err := db.Select(&postsList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("文章库里面是空的:", zap.Error(err))
			fmt.Println("文章库是空的")
			err = nil
		} else {
			fmt.Println("文章篇数：", len(postsList))
		}

	}
	return
}

func GetPostsDetailByTitle(title string) (post *models.Posts, err error) {
	sqlStr := "select id, title, user_id, content, created_at from posts where title = ?"
	post = &models.Posts{}

	if err := db.Get(post, sqlStr, title); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("数据库中没有这篇文章:", zap.Error(err))
			return nil, ErrorPostNotExist
		}
	}

	return post, nil
}

func CreatePost(title, content string, userid int64) error {
	sqlStr := "insert into posts(title, user_id, content,created_at,updated_at) values(?, ?, ?,?,?)"
	_, err := db.Exec(sqlStr, title, userid, content, time.Now(), time.Now())
	if err != nil {
		zap.L().Error("发布文章失败:", zap.Error(err))
		return err
	}
	return nil
}

func UpdatePost(title, content string) error {
	sqlStr := "update posts set title=?, content=?, updated_at=? where id=?"

	result, err := db.Exec(sqlStr, title, content, time.Now(), 30)

	if err != nil {
		zap.L().Error("更新文章失败:", zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Error("获取影响更新行数失败:", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		zap.L().Error("没有找到文章:", zap.Error(err))
		return ErrorPostNotExist
	}

	return nil
}

func DeletePost(title string) error {
	sqlStr := "delete from posts where title=?"
	result, err := db.Exec(sqlStr, title)

	if err != nil {
		zap.L().Error("删除文章失败:", zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Error("获取影响更新行数失败:", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		zap.L().Error("没有找到要删除的文章:", zap.Error(err))
		return ErrorPostNotExist
	}

	return nil
}
