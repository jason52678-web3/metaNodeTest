package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github/task4/MyBlog/dao/mysql"
	"github/task4/MyBlog/logic"
	"github/task4/MyBlog/models"
	"go.uber.org/zap"
	"net/http"
)

func SignUpHandler(c *gin.Context) {
	//1.获取参数与参数校验
	p := new(models.ParamSignUp)

	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})

		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
		return
	}

	//手动对请求参数进行详细的业务规则校验
	if p.Password != p.RePassword {
		zap.L().Error("SignUp with invalid param")
		ResponseErrorWithMsg(c, CodeInvalidParam, "两次输入的密码不一致")
		return
	}

	fmt.Println(p)

	//2.业务处理

	if err := logic.SignUp(p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseErrorWithMsg(c, CodeServerBusy, err.Error())
		return
	}

	//3.返回响应
	ResponseSuccess(c, "注册成功")
}

func LoginHandler(c *gin.Context) {
	//1.获取参数及参数验证
	p := new(models.ParamLogin)

	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
		return
	}

	//2.逻辑处理
	token, err := logic.Login(p)
	if err != nil {
		zap.L().Error("Login with invalid param", zap.String("username", p.Username),
			zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	//3.返回相应
	ResponseSuccess(c, token)
}
