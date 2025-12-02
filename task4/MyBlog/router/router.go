package router

import (
	"github.com/gin-gonic/gin"
	"github/task4/MyBlog/controller"
	"github/task4/MyBlog/logger"
	"github/task4/MyBlog/middlewares"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")
	//注册用户
	v1.POST("/signup", controller.SignUpHandler)

	//登录
	v1.POST("/login", controller.LoginHandler)

	v1.GET("/postsall", controller.PostsAllHandler)
	v1.GET("/postsdetail/:title", controller.PostsDetailHandler)

	v1.Use(middlewares.JWTAuthMiddleware()) //应用JWT认证中间件
	{
		v1.POST("createpost", controller.CreatePostHandler)
		v1.POST("updatepost", controller.UpdatePostHandler)
		v1.POST("deletepost", controller.DeletePostHandler)

		v1.POST("createcomment", controller.CreateCommentHandler)
		v1.GET("getpostcomments/:post_id", controller.GetPostCommentsHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	return r

}
