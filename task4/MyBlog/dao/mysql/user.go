package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github/task4/MyBlog/models"
)

const secret string = "metaNode Project"

var (
	ErrorUserExist      = errors.New("用户已存在")
	ErrorUserNotExist   = errors.New("用户不存在")
	ErrorInvalidPssword = errors.New("用户或密码错误")
)

// 把每一步数据库操作封装成函数
// 等待logic层根据业务需求调用
func CheckUserExist(username string) error {
	sqlStr := `select count(user_id) from user where username = ? `
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return nil
}

func InsertUser(user *models.User) (err error) {
	//对密码进行加密
	user.Password = encryptPassword(user.Password)

	sqlStr := `insert into user (user_id,username,password) values (?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	if err != nil {
		fmt.Println("插入数据库错误:", err.Error())
	}
	return err

}

func encryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))

}

func Login(user *models.User) (err error) {
	oPassWord := user.Password //用户登录密码
	sqlStr := `select user_id, username, password from user where username = ? `
	err = db.Get(user, sqlStr, user.Username)
	if err != nil {
		return ErrorUserNotExist
	}

	//判断密码是否正确
	password := encryptPassword(oPassWord)
	if password != user.Password {
		return ErrorInvalidPssword
	}
	return nil
}
