package base

import (
	"context"
	"github.com/MIMIKING-B/xcross-go-starter/api/home/base"
	"github.com/MIMIKING-B/xcross-go-starter/internal/consts"
	"github.com/MIMIKING-B/xcross-go-starter/internal/model"
	"github.com/MIMIKING-B/xcross-go-starter/internal/service"
	"github.com/MIMIKING-B/xcross-go-starter/utility/simple"

	"github.com/gogf/gf/v2/frame/g"
)

// Site 基础
var Site = cSite{}

type cSite struct{}

func (a *cSite) Index(ctx context.Context, _ *base.SiteIndexReq) (res *base.SiteIndexRes, err error) {
	service.View().Render(ctx, model.View{Data: g.Map{
		"name":    simple.AppName(ctx),
		"version": consts.VersionApp,
	}})

	// err = gerror.New("这是一个测试错误")
	// return

	// err = gerror.NewCode(gcode.New(10000, "这是一个测试自定义错误码错误", nil))
	// return

	// service.View().Error(ctx, gerror.New("这是一个允许被自定义格式的错误，默认和通用错误格式一致，你可以修改它"))
	// return
	return
}
