package mysql

import (
	"database/sql"
	"errors"
	"github/task4/MyBlog/models"
	"go.uber.org/zap"
	"time"
)

var (
	ErrorCommentsNotExist = errors.New("指定的评论帖子不存在")
)

func CreateComment(postid int, content string, userid int64) error {
	postStr := "select * from posts where id = ?"
	post := &models.Posts{}
	err := db.Get(post, postStr, postid)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("想评论的文章不存在:", zap.Error(err))
		}
		return ErrorCommentsNotExist
	}

	sqlStr := "insert into comments(content, user_id, post_id,created_at) values(?, ?, ?,?)"
	_, err = db.Exec(sqlStr, content, userid, postid, time.Now())
	if err != nil {
		zap.L().Error("发布文章评论失败:", zap.Error(err))
		return err
	}
	return nil
}

func GetPostCommentsList(postid int) ([]*models.Comments, error) {
	sqlStr := "select * from comments where post_id = ?"

	commentsList := make([](*models.Comments), 0, 2)
	err := db.Select(&commentsList, sqlStr, postid)
	if err != nil {
		zap.L().Warn("查询错误：", zap.Error(err))
		return nil, err
	}

	if len(commentsList) == 0 {
		zap.L().Warn("要查找文章没有评论：", zap.Error(err))
		return nil, ErrorCommentsNotExist
	}

	return commentsList, nil
}
