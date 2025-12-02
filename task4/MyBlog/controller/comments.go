package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github/task4/MyBlog/dao/mysql"
	"github/task4/MyBlog/logic"
	"github/task4/MyBlog/models"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func CreateCommentHandler(c *gin.Context) {
	p := new(models.ParamComment)

	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("CreateComment with invalid param", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})

		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
		return
	}

	userid := c.GetInt64(CtxUserIDKey)

	err := logic.CreateComment(p.PostId, p.Content, userid)
	if err != nil {
		if errors.Is(err, mysql.ErrorCommentsNotExist) {
			zap.L().Error("CreateComment with invalid param", zap.Error(err))
			ResponseError(c, CodeCommentNotExist)
			return
		}
		zap.L().Error("logic.CreateComment() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, "创建评论成功")
}

func GetPostCommentsHandler(c *gin.Context) {
	postid, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		zap.L().Error("logic.GetPostCommentsHandler() failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, err := logic.GetPostCommentsList(postid)
	if err != nil {
		if errors.Is(err, mysql.ErrorCommentsNotExist) {
			zap.L().Error("logic.GetPostCommentsHandler() failed", zap.Error(err))
			ResponseError(c, CodeCommentNotExist)
			return
		}
		zap.L().Error("logic.GetPostCommentsList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)

}
