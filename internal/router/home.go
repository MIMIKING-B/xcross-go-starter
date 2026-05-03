package router

import (
	"context"
	"xcross-go-starter/internal/consts"
	"xcross-go-starter/internal/controller/home/base"
	"xcross-go-starter/internal/service"
	"xcross-go-starter/utility/simple"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Home 前台页面路由
func Home(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(service.Middleware().HomeAuth)
		// 允许通过根地址访问的路由可以加到这里，访问地址：http://127.0.0.1:8000
		group.Bind(
			base.Site, // 基础
		)
		// 默认访问地址：http://127.0.0.1:8000/home
		group.Group(simple.RouterPrefix(ctx, consts.AppHome), func(group *ghttp.RouterGroup) {
			group.Bind(
				base.Site, // 基础
			)
		})
	})
}
