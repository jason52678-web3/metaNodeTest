package logic

import (
	"github/task4/MyBlog/dao/mysql"
	"github/task4/MyBlog/models"
)

func CreateComment(postid int, content string, userid int64) error {
	return mysql.CreateComment(postid, content, userid)
}

func GetPostCommentsList(postid int) ([]*models.Comments, error) {
	return mysql.GetPostCommentsList(postid)
}
