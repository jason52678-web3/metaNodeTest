package models

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamPost 文章参数
type ParamPost struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type ParamDeletePost struct {
	Title string `json:"title" binding:"required"`
}

type ParamComment struct {
	Content string `json:"content" binding:"required"`
	PostId  int    `json:"post_id" binding:"required"`
}
