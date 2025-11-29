package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(w http.ResponseWriter, data interface{}) {
	httpx.OkJson(w, Response{
		Code: 0,
		Msg:  "成功",
		Data: data,
	})
}

// SuccessMsg 成功响应（自定义消息）
func SuccessMsg(w http.ResponseWriter, msg string, data interface{}) {
	httpx.OkJson(w, Response{
		Code: 0,
		Msg:  msg,
		Data: data,
	})
}

// Error 错误响应
func Error(w http.ResponseWriter, code int, msg string) {
	httpx.OkJson(w, Response{
		Code: code,
		Msg:  msg,
	})
}

// ErrorData 错误响应（带数据）
func ErrorData(w http.ResponseWriter, code int, msg string, data interface{}) {
	httpx.OkJson(w, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
