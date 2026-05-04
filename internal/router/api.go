package router

import (
	"context"
	"github.com/MIMIKING-B/xcross-go-starter/internal/consts"
	"github.com/MIMIKING-B/xcross-go-starter/utility/simple"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Api 前台路由
func Api(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group(simple.RouterPrefix(ctx, consts.AppApi), func(group *ghttp.RouterGroup) {
		group.Bind(
		// TODO 可以添加自定义的路由
		)
	})
}
