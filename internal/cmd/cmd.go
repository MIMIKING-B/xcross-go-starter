package cmd

import (
	"context"

	"github.com/MIMIKING-B/xcross-go-starter/utility/simple"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = &gcmd.Command{
		Description: `默认启动所有服务`,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			return All.Func(ctx, parser)
		},
	}

	Help = &gcmd.Command{
		Name:  "help",
		Brief: "查看帮助",
		Description: `
		命令提示符
		---------------------------------------------------------------------------------
		启动服务
		>> 所有服务  [go run main.go]   热编译  [gf run main.go]
		>> HTTP服务  [go run main.go http]
		>> 查看帮助  [go run main.go help]
		---------------------------------------------------------------------------------
		工具
		>> 打印所有打包的资源文件列表  [go run main.go tools -m=gres -a1=dump]
		>> 打印指定打包的资源文件内容  [go run main.go tools -m=gres -a1=content -a2=resource/template/home/index.html]
		---------------------------------------------------------------------------------
    `,
	}

	All = &gcmd.Command{
		Name:        "all",
		Brief:       "start all server",
		Description: "this is the command entry for starting all server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			g.Log().Debug(ctx, "starting all server")
			// 需要启动的服务
			var allServers = []*gcmd.Command{Http}
			for _, server := range allServers {
				var cmd = server
				simple.SafeGo(ctx, func(ctx context.Context) {
					if err := cmd.Func(ctx, parser); err != nil {
						g.Log().Fatalf(ctx, "%v start fail:%v", cmd.Name, err)
					}
				})
			}
			// 信号监听
			signalListen(ctx, signalHandlerForOverall)
			<-serverCloseSignal
			serverWg.Wait()
			g.Log().Debug(ctx, "all service successfully closed ..")
			return
		},
	}
)

func init() {
	if err := Main.AddCommand(All, Http, Tools, Help); err != nil {
		panic(err)
	}
}
