// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"github.com/MIMIKING-B/xcross-go-starter/internal/model"
)

type (
	IView interface {
		// Render 渲染默认模板页面
		Render(ctx context.Context, data ...model.View)
		// RenderTpl 渲染指定模板页面
		RenderTpl(ctx context.Context, tpl string, data ...model.View)
		// Error 自定义错误页面
		Error(ctx context.Context, err error)
	}
)

var (
	localView IView
)

func View() IView {
	if localView == nil {
		panic("implement not found for interface IView, forgot register?")
	}
	return localView
}

func RegisterView(i IView) {
	localView = i
}
