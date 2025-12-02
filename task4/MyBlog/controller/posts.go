package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github/task4/MyBlog/dao/mysql"
	"github/task4/MyBlog/logic"
	"github/task4/MyBlog/models"
	"go.uber.org/zap"
	"net/http"
)

func PostsAllHandler(c *gin.Context) {
	//1.获取所有文章列表，包含(id,titile,users_id, create_at)
	data, err := logic.GetPostsList()
	if err != nil {
		zap.L().Error("logic.GetPostsList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)
}

func PostsDetailHandler(c *gin.Context) {
	title := c.Param("title")
	//	fmt.Println("Controlloer title:", title)
	data, err := logic.GetPostsDetailByTitle(title)
	if err != nil {
		zap.L().Error("logic.GetPostsDetail() failed", zap.Error(err))
		ResponseError(c, CodePostNotExist)
		return
	}
	ResponseSuccess(c, data)
}

func CreatePostHandler(c *gin.Context) {

	p := new(models.ParamPost)

	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("CreatePost with invalid param", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})

		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
		return
	}

	userid := c.GetInt64(CtxUserIDKey)

	err := logic.CreatePost(p.Title, p.Content, userid)
	if err != nil {
		zap.L().Error("logic.CreatePost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, "创建帖子成功")
}

func UpdatePostHandler(c *gin.Context) {
	p := new(models.ParamPost)

	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("Update with invalid param", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})

		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
		return
	}

	err := logic.UpdatePost(p.Title, p.Content)
	if err != nil {
		if errors.Is(err, mysql.ErrorPostNotExist) {
			zap.L().Error("logic.UpdatePost() post not found", zap.Error(err))
			ResponseError(c, CodePostNotExist)
			return
		}

		zap.L().Error("logic.UpdatePost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	} else {
		ResponseSuccess(c, "帖子更新成功")
	}

}

func DeletePostHandler(c *gin.Context) {
	p := new(models.ParamDeletePost)

	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("Delete with invalid param", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})

		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
		return
	}

	err := logic.DeletePost(p.Title)
	if err != nil {
		if errors.Is(err, mysql.ErrorPostNotExist) {
			zap.L().Error("logic.DeletePost() post not found", zap.Error(err))
			ResponseError(c, CodePostNotExist)
			return
		}
		zap.L().Error("logic.DeletePost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	} else {
		ResponseSuccess(c, "帖子删除成功")
	}

}
