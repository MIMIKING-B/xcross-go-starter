package simple

import (
	"context"
	"github.com/MIMIKING-B/xcross-go-starter/internal/consts"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
)

// RouterPrefix 获取应用路由前缀
func RouterPrefix(ctx context.Context, app string) string {
	return g.Cfg().MustGet(ctx, "router."+app+".prefix", "/"+app+"").String()
}

// DefaultErrorTplContent 获取默认的错误模板内容
func DefaultErrorTplContent(ctx context.Context) string {
	return gfile.GetContents(g.Cfg().MustGet(ctx, "viewer.paths").String() + "/error/default.html")
}

// GetHeaderLocale 获取请求头语言设置
// gf支持格式：en/ja/ru/zh-CN/zh-TW
func GetHeaderLocale(ctx context.Context) (lang string) {
	lang = g.Cfg().MustGet(ctx, "system.i18n.defaultLanguage", consts.SysDefaultLanguage).String()
	// 没有开启国际化，使用默认语言
	if !g.Cfg().MustGet(ctx, "system.i18n.switch", true).Bool() {
		return
	}

	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return
	}
	locale := r.Header.Get("Locale")
	// 简体中文
	if locale == "zh-CN" || locale == "zh-Hans" || locale == "zh" || locale == "ZH" {
		lang = "zh-CN"
		return
	}
	// 繁体
	if locale == "zh-TW" || locale == "zh-Hant" {
		lang = "zh-TW"
		return
	}
	// 英文
	if locale == "en" || locale == "EN" {
		lang = "en"
		return
	}
	// 更多语言
	// ...
	return
}

// SafeGo 安全的调用协程，遇到错误时输出错误日志而不是抛出panic
func SafeGo(ctx context.Context, f func(ctx context.Context), lv ...int) {
	g.Go(ctx, f, func(ctx context.Context, err error) {
		var level = glog.LEVEL_ERRO
		if len(lv) > 0 {
			level = lv[0]
		}
		Logf(level, ctx, "SafeGo exec failed:%+v", err)
	})
}

// Logf 打印对应的错误日志
func Logf(level int, ctx context.Context, format string, v ...interface{}) {
	switch level {
	case glog.LEVEL_DEBU:
		g.Log().Debugf(ctx, format, v...)
	case glog.LEVEL_INFO:
		g.Log().Infof(ctx, format, v...)
	case glog.LEVEL_NOTI:
		g.Log().Noticef(ctx, format, v...)
	case glog.LEVEL_WARN:
		g.Log().Warningf(ctx, format, v...)
	case glog.LEVEL_ERRO:
		g.Log().Errorf(ctx, format, v...)
	case glog.LEVEL_CRIT:
		g.Log().Criticalf(ctx, format, v...)
	case glog.LEVEL_PANI:
		g.Log().Panicf(ctx, format, v...)
	case glog.LEVEL_FATA:
		g.Log().Fatalf(ctx, format, v...)
	default:
		g.Log().Errorf(ctx, format, v...)
	}
}
