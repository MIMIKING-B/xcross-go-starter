package global

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/MIMIKING-B/xcross-go-starter/internal/consts"
	"github.com/MIMIKING-B/xcross-go-starter/internal/library/cache"
	"github.com/MIMIKING-B/xcross-go-starter/internal/service"
	"github.com/MIMIKING-B/xcross-go-starter/utility/simple"
	"github.com/MIMIKING-B/xcross-go-starter/utility/validate"

	"github.com/gogf/gf/contrib/trace/jaeger/v2"
	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gmode"
)

// Init 项目初始化
func Init(ctx context.Context) {
	// 设置gf运行模式
	SetGFMode(ctx)
	// 默认上海时区
	var err error
	if runtime.GOOS == "windows" {
		// Windows 用固定时区（UTC+8）
		cst := time.FixedZone("CST", 8*3600)
		err := gtime.SetTimeZone(cst.String())
		if err != nil {
			return
		}
		time.Local = cst
	} else {
		// Linux/macOS 用标准时区
		err = gtime.SetTimeZone("Asia/Shanghai")
	}
	if err != nil {
		g.Log().Fatalf(ctx, "时区设置异常 err：%+v", err)
		return
	}
	fmt.Printf("当前运行环境：%v, 运行根路径为：%v \r\n初始化版本：v%v, gf版本：%v \n", runtime.GOOS, gfile.Pwd(), consts.VersionApp, gf.VERSION)
	// 初始化链路追踪
	InitTrace(ctx)
	// 设置缓存适配器
	cache.SetAdapter(ctx)
	// 初始化功能库配置
	service.SysConfig().InitConfig(ctx)
}

// InitTrace 初始化链路追踪
func InitTrace(ctx context.Context) {
	if !g.Cfg().MustGet(ctx, "jaeger.switch").Bool() {
		return
	}
	tp, err := jaeger.Init(simple.AppName(ctx), g.Cfg().MustGet(ctx, "jaeger.endpoint").String())
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	simple.Event().Register(consts.EventServerClose, func(ctx context.Context, args ...interface{}) {
		_ = tp.Shutdown(ctx)
		g.Log().Debug(ctx, "jaeger closed ..")
	})
}

// SetGFMode 设置gf运行模式
func SetGFMode(ctx context.Context) {
	mode := g.Cfg().MustGet(ctx, "system.mode").String()
	if len(mode) == 0 {
		mode = gmode.NOT_SET
	}
	var modes = []string{gmode.DEVELOP, gmode.TESTING, gmode.STAGING, gmode.PRODUCT}
	// 如果是有效的运行模式，就进行设置
	if validate.InSlice(modes, mode) {
		gmode.Set(mode)
	}
}
