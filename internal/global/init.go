package global

import (
	"context"
	"fmt"
	"runtime"
	"xcross-go-starter/internal/consts"
	"xcross-go-starter/internal/library/cache"
	"xcross-go-starter/utility/validate"

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
	if err := gtime.SetTimeZone("Asia/Shanghai"); err != nil {
		g.Log().Fatalf(ctx, "时区设置异常 err：%+v", err)
		return
	}
	fmt.Printf("当前运行环境：%v, 运行根路径为：%v \r\n初始化版本：v%v, gf版本：%v \n", runtime.GOOS, gfile.Pwd(), consts.VersionApp, gf.VERSION)
	// 设置缓存适配器
	cache.SetAdapter(ctx)
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
