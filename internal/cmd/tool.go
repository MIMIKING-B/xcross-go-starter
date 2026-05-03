// Package cmd 命令行工具包
// 提供项目内常用命令行功能
package cmd

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gres"
)

// 列出所有打包文件
// go run main.go tools -m=gres -a1=dump
// 查看某个打包文件内容
// go run main.go tools -m=gres -a1=content -a2=resource/template/home/index.html
var (
	// Tools 命令行工具入口
	// 原始功能：提供项目常用工具命令行入口
	Tools = &gcmd.Command{
		Name:        "tools",
		Brief:       "常用工具",
		Description: `项目常用命令行调试工具`,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 获取所有命令行参数
			args := parser.GetOptAll()
			g.Log().Debugf(ctx, "tools args:%v", args)
			// 参数不能为空
			if len(args) == 0 {
				err = gerror.New("tools args cannot be empty.")
				return
			}
			// 获取 -m 参数（执行方法类型）
			// 原始功能：判断要执行哪种工具
			method, ok := args["m"]
			if !ok {
				err = gerror.New("tools method cannot be empty.")
				return
			}
			// 根据 -m 参数分发执行逻辑
			switch method {
			case "gres":
				// 执行 gres 资源打包调试工具
				err = handleGRes(ctx, args)
			default:
				err = gerror.Newf("tools method[%v] does not exist", method)
			}
			// 执行成功提示
			if err == nil {
				g.Log().Info(ctx, "tools exec successful!")
			}
			return
		},
	}
)

// handleGRes
// 处理 gres 资源打包调试命令
// 原始功能：查看 gf pack 打包进程序的静态资源
//
// -------------------------- 重要说明 --------------------------
// gres = GoFrame 资源打包工具
// gf pack = GoFrame 框架的打包命令
// 作用：将 config/html/js/css 等静态文件打包到 go 源码中
// 最终程序变成一个独立二进制文件，不需要附带任何静态文件
// 打包命令示例：gf pack config,static,template packed/data.go
// --------------------------------------------------------------
func handleGRes(ctx context.Context, args map[string]string) (err error) {
	// 获取 -a1 参数（子命令类型）
	a1, ok := args["a1"]
	if !ok {
		err = gerror.New("gres args cannot be empty.")
		return
	}
	switch a1 {
	// -a1=dump
	// 原始功能：打印所有已打包进程序的文件列表
	case "dump":
		gres.Dump()
	// -a1=content -a2=文件路径
	// 原始功能：查看某个打包文件的内容
	case "content":
		// 获取 -a2 参数（要查看的文件路径）
		path, ok := args["a2"]
		if !ok {
			err = gerror.New("缺少查看文件路径参数：`a2`")
			return
		}
		// 判断文件是否被打包进程序
		if !gres.Contains(path) {
			err = gerror.Newf("没有找到资源文件:%v", path)
			return
		}
		// 获取文件内容
		content := string(gres.GetContent(path))
		// 内容不能为空
		if len(content) == 0 {
			err = gerror.Newf("没有找到资源文件内容，请确认传入`a2`参数是一个文件，a2:%v", path)
			return
		}
		// 输出文件内容
		fmt.Println("以下是资源文件内容:")
		fmt.Println(content)
	default:
		err = gerror.Newf("handleGRes a1 is invalid, a1:%v", a1)
	}
	return
}
