package middleware

import (
	"github.com/MIMIKING-B/xcross-go-starter/internal/consts"
	"github.com/MIMIKING-B/xcross-go-starter/internal/library/response"
	"github.com/MIMIKING-B/xcross-go-starter/utility/charset"
	"github.com/MIMIKING-B/xcross-go-starter/utility/simple"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gmeta"
)

// ResponseHandler 全局HTTP响应预处理中间件
// 统一处理接口响应格式、错误码、异常页面、响应类型分发
func (s *sMiddleware) ResponseHandler(r *ghttp.Request) {
	// 执行后续中间件/控制器逻辑
	r.Middleware.Next()
	// 统一接管HTTP错误状态码
	switch r.Response.Status {
	case 403:
		r.Response.Writeln("403 - 网站拒绝显示此网页")
		return
	case 404:
		r.Response.Writeln("404 - 请求资源不存在")
		return
	}
	// 获取当前响应的Content-Type类型
	contentType := getContentType(r)
	// 如果已经存在响应内容（非流类型），则直接返回不做处理
	if contentType != consts.HTTPContentTypeStream && r.Response.BufferLength() > 0 {
		return
	}
	// 根据不同响应类型分发处理
	switch contentType {
	case consts.HTTPContentTypeHtml:
		// HTML页面响应
		s.responseHtml(r)
		return
	case consts.HTTPContentTypeXml:
		// XML格式响应
		s.responseXml(r)
		return
	case consts.HTTPContentTypeStream:
		// 流类型不处理
	case consts.HTTPContentTypeOctetStream:
		// 二进制流不处理
	default:
		// 默认统一返回JSON格式
		responseJson(r)
	}
}

// responseHtml HTML页面错误响应处理
// 当页面请求发生错误时，渲染统一的错误模板页面
func (s *sMiddleware) responseHtml(r *ghttp.Request) {
	code, message, resp := parseResponse(r)
	// 业务正常则不处理
	if code == gcode.CodeOK.Code() {
		return
	}

	// 清空原有响应缓冲区
	r.Response.ClearBuffer()
	// 渲染系统默认错误模板页面
	_ = r.Response.WriteTplContent(simple.DefaultErrorTplContent(r.Context()), g.Map{
		"code":    code,
		"message": message,
		"stack":   resp,
	})
}

// responseXml XML格式响应处理
// 统一封装XML格式返回数据
func (s *sMiddleware) responseXml(r *ghttp.Request) {
	code, message, data := parseResponse(r)
	response.RXml(r, code, message, data)
}

// responseJson JSON格式响应处理
// 项目最常用：统一封装JSON格式返回数据
func responseJson(r *ghttp.Request) {
	code, message, data := parseResponse(r)
	response.RJson(r, code, message, data)
}

// parseResponse 解析请求响应与错误信息
// 统一提取：错误码、错误信息、响应数据
func parseResponse(r *ghttp.Request) (code int, message string, resp interface{}) {
	ctx := r.Context()
	// 获取请求执行过程中的错误信息
	err := r.GetError()
	// 无错误：返回成功+业务响应数据
	if err == nil {
		return gcode.CodeOK.Code(), "操作成功", r.GetHandlerResponse()
	}
	// 调试模式：输出详细错误信息与堆栈
	if simple.Debug(ctx) {
		// 获取最顶层错误信息
		message = gerror.Current(err).Error()
		// HTML页面返回格式化堆栈字符串
		if getContentType(r) == consts.HTTPContentTypeHtml {
			resp = charset.SerializeStack(err)
		} else {
			// 接口返回结构化堆栈
			resp = charset.ParseErrStack(err)
		}
	} else {
		// 生产模式：只返回友好提示信息
		message = consts.ErrorMessage(gerror.Current(err))
	}
	// 获取错误码
	code = gerror.Code(err).Code()
	// 错误日志记录区分
	// CodeNil(-1)：业务可控错误 → 只记录到文件，不输出到控制台
	// 其他错误：系统异常 → 记录错误日志+堆栈
	if code == gcode.CodeNil.Code() {
		g.Log().Stdout(false).Infof(ctx, "exception:%v", err)
	} else {
		g.Log().Errorf(ctx, "exception:%v", err)
	}
	return
}

// getContentType 获取响应Content-Type类型
// 优先级：响应头设置 > 结构体meta注解 > 默认JSON
func getContentType(r *ghttp.Request) (contentType string) {
	// 先从响应头获取Content-Type
	contentType = r.Response.Header().Get("Content-Type")
	if contentType != "" {
		return
	}
	// 从返回结构体的mime元标签获取
	mime := gmeta.Get(r.GetHandlerResponse(), "mime").String()
	if mime == "" {
		// 默认使用JSON格式
		contentType = consts.HTTPContentTypeJson
	} else {
		contentType = mime
	}
	return
}
