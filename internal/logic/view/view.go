package view

import (
	"context"
	"xcross-go-starter/internal/model"
	"xcross-go-starter/internal/service"
	"xcross-go-starter/utility/charset"
	"xcross-go-starter/utility/simple"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type sView struct{}

func init() {
	service.RegisterView(New())
}

func New() *sView {
	return &sView{}
}

// Render 渲染默认模板页面
func (s *sView) Render(ctx context.Context, data ...model.View) {
	s.RenderTpl(ctx, g.Cfg().MustGet(ctx, "viewer.homeLayout").String(), data...)
}

// RenderTpl 渲染指定模板页面
func (s *sView) RenderTpl(ctx context.Context, tpl string, data ...model.View) {
	var (
		viewObj = model.View{}
		request = g.RequestFromCtx(ctx)
	)
	if len(data) > 0 {
		viewObj = data[0]
	}
	if viewObj.Title == "" {
		viewObj.Title = g.Cfg().MustGet(ctx, `setting.title`).String()
	} else {
		viewObj.Title = viewObj.Title + ` - ` + g.Cfg().MustGet(ctx, `setting.title`).String()
	}
	if viewObj.Keywords == "" {
		viewObj.Keywords = g.Cfg().MustGet(ctx, `setting.keywords`).String()
	}
	if viewObj.Description == "" {
		viewObj.Description = g.Cfg().MustGet(ctx, `setting.description`).String()
	}
	if viewObj.IpcCode == "" {
		viewObj.IpcCode = g.Cfg().MustGet(ctx, `setting.icpCode`).String()
	}
	if viewObj.GET == nil {
		viewObj.GET = request.GetQueryMap()
	}
	// 去掉空数据
	viewData := gconv.Map(viewObj)
	for k, v := range viewData {
		if g.IsEmpty(v) {
			delete(viewData, k)
		}
	}
	// 渲染模板
	_ = request.Response.WriteTpl(tpl, viewData)
}

// Error 自定义错误页面
func (s *sView) Error(ctx context.Context, err error) {
	var (
		request = g.RequestFromCtx(ctx)
		code    = gerror.Code(err)
		stack   string
	)
	// 是否输出错误堆栈到页面
	if simple.Debug(ctx) {
		stack = charset.SerializeStack(err)
	}
	request.Response.ClearBuffer()
	_ = request.Response.WriteTplContent(simple.DefaultErrorTplContent(ctx), g.Map{
		"code":    code.Code(),
		"message": code.Message(),
		"stack":   stack,
	})
}
