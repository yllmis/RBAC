package tools

import "fmt"

var (
	OK           = ECode{Code: 0}
	NotLogin     = ECode{Code: 10001, Message: "用户未登录"}
	ParamErr     = ECode{Code: 10002, Message: "参数错误"}
	UserErr      = ECode{Code: 10003, Message: "用户名或密码错误"}
	NoPermission = ECode{Code: 10004, Message: "没有权限访问"}
	NotFound     = ECode{Code: 10005, Message: "资源未找到"}
)

type ECode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (e *ECode) String() string {
	return fmt.Sprintf("code:%d,message:%s", e.Code, e.Message)
}
