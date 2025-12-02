package controller

type ResCode int64

const CtxUserIDKey string = "UserID"

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeInvalidToken
	CodeNeedLogin
	CodePostNotExist
	CodeCommentNotExist
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户已经存在",
	CodeUserNotExist:    "该用户不存在",
	CodeInvalidPassword: "密码错误",
	CodeServerBusy:      "服务繁忙",
	CodeInvalidToken:    "无效的Token",
	CodeNeedLogin:       "需要登录",
	CodePostNotExist:    "要操作的帖子不存在",
	CodeCommentNotExist: "指定的评论帖子不存在",
}

func (res ResCode) Msg() string {
	msg, ok := codeMsgMap[res]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
