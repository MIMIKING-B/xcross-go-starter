package middleware

import (
	"net/http"
	"strings"
	"xcross-go-starter/internal/consts"
	"xcross-go-starter/internal/library/contexts"
	"xcross-go-starter/internal/library/response"
	"xcross-go-starter/internal/model"
	"xcross-go-starter/internal/service"
	"xcross-go-starter/utility/simple"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
	"go.opentelemetry.io/otel/attribute"
)

type sMiddleware struct {
	LoginUrl         string // 登录路由地址
	DemoWhiteList    g.Map  // 演示模式放行的路由白名单
	NotRecordRequest g.Map  // 不记录请求数据的路由（当前请求数据过大时会影响响应效率，可以将路径放到该选项中改善）
}

func init() {
	service.RegisterMiddleware(NewMiddleware())
}

func NewMiddleware() *sMiddleware {
	return &sMiddleware{
		LoginUrl: "/common",
		DemoWhiteList: g.Map{
			"/admin/site/accountLogin": struct{}{}, // 账号登录
		},
		NotRecordRequest: g.Map{
			"/admin/upload/file": struct{}{}, // 上传文件
		},
	}
}

// CORS allows Cross-origin resource sharing.
func (s *sMiddleware) CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

// DemoLimit 演示系统操作限制
func (s *sMiddleware) DemoLimit(r *ghttp.Request) {
	if !simple.IsDemo(r.Context()) {
		r.Middleware.Next()
		return
	}
	if r.Method == http.MethodPost {
		if _, ok := s.DemoWhiteList[r.URL.Path]; ok {
			r.Middleware.Next()
			return
		}
		response.JsonExit(r, gcode.CodeNotSupported.Code(), "演示系统禁止操作！")
		return
	}
	r.Middleware.Next()
}

// Ctx 初始化请求上下文
func (s *sMiddleware) Ctx(r *ghttp.Request) {
	// 国际化
	r.SetCtx(gi18n.WithLanguage(r.Context(), simple.GetHeaderLocale(r.Context())))
	// 链路追踪
	if g.Cfg().MustGet(r.Context(), "jaeger.switch").Bool() {
		ctx, span := gtrace.NewSpan(r.Context(), "middleware.ctx")
		span.SetAttributes(attribute.KeyValue{
			Key:   "traceID",
			Value: attribute.StringValue(gctx.CtxId(ctx)),
		})
		span.End()
		r.SetCtx(ctx)
	}
	data := make(g.Map)
	if _, ok := s.NotRecordRequest[r.URL.Path]; ok {
		data["request.body"] = gjson.New(nil)
	} else {
		data["request.body"] = gjson.New(r.GetBodyString())
	}
	contexts.Init(r, &model.Context{
		Data:   data,
		Module: getModule(r.URL.Path),
	})
	if len(r.Cookie.GetSessionId()) == 0 {
		r.Cookie.SetSessionId(gctx.CtxId(r.Context()))
	}
	r.SetCtx(r.GetNeverDoneCtx())
	r.Middleware.Next()
}

// getModule 获取当前模块名称
func getModule(path string) (module string) {
	slice := strings.Split(path, "/")
	if len(slice) < 2 {
		module = consts.AppDefault
		return
	}
	if slice[1] == "" {
		module = consts.AppDefault
		return
	}
	return slice[1]
}
