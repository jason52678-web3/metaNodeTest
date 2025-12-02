package logic

import (
	"errors"
	"github/task4/MyBlog/dao/mysql"
	"github/task4/MyBlog/models"
	"github/task4/MyBlog/pkg/jwt"
	"github/task4/MyBlog/pkg/snowflake"
)

// 处理业务逻辑
func SignUp(p *models.ParamSignUp) (err error) {
	//1.判断用户存不存在
	err = mysql.CheckUserExist(p.Username)
	if err != nil {
		return errors.New("数据库查询错误")
	}

	//.生成UID
	userID := snowflake.GenID()
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	//.保存数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (token string, err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// Login传递的是指针，可以拿到user.UserID
	if err := mysql.Login(user); err != nil {
		return " ", err
	}

	return jwt.GenToken(user.UserID, user.Username)

}
