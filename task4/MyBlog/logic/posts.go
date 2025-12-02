package logic

import (
	"github/task4/MyBlog/dao/mysql"
	"github/task4/MyBlog/models"
)

func GetPostsList() ([]*models.Posts, error) {
	//查找数据库，查找到所有的文章列表
	return mysql.GetPostsList()
}

func GetPostsDetailByTitle(title string) (*models.Posts, error) {
	return mysql.GetPostsDetailByTitle(title)
}

func CreatePost(title, content string, userid int64) error {
	return mysql.CreatePost(title, content, userid)
}

func UpdatePost(title, content string) error {
	return mysql.UpdatePost(title, content)
}

func DeletePost(title string) error {
	return mysql.DeletePost(title)
}
