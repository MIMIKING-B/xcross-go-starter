// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	IMiddleware interface {
		// HomeAuth 前台页面鉴权中间件
		HomeAuth(r *ghttp.Request)
		// CORS allows Cross-origin resource sharing.
		CORS(r *ghttp.Request)
		// DemoLimit 演示系统操作限制
		DemoLimit(r *ghttp.Request)
		// Ctx 初始化请求上下文
		Ctx(r *ghttp.Request)
	}
)

var (
	localMiddleware IMiddleware
)

func Middleware() IMiddleware {
	if localMiddleware == nil {
		panic("implement not found for interface IMiddleware, forgot register?")
	}
	return localMiddleware
}

func RegisterMiddleware(i IMiddleware) {
	localMiddleware = i
}
