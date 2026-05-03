// Package global 全局变量、公共工具、路由元信息管理包，提供HTTP路由元信息解析与缓存能力
package global

import (
	"reflect"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/gtag"
)

// HTTPRouter 扩展GF框架原生路由项，增加接口文档相关元数据字段
type HTTPRouter struct {
	ghttp.RouterItem        // 继承GF框架原生路由项结构体
	Tags             string `json:"tags"         dc:"接口所属的标签，用于接口分类"`
	Summary          string `json:"summary"      dc:"接口/参数概要描述"`
	Description      string `json:"description"  dc:"接口/参数详细描述"`
}

var (
	// httpRoutes 全局HTTP路由元信息缓存，key为GenRouteKey生成的唯一标识
	httpRoutes map[string]*HTTPRouter
	// routeMutex 路由缓存并发安全锁
	routeMutex sync.Mutex
	// shortTypeMapForTag 标签简写映射表，用于兼容GF标签的简写形式
	shortTypeMapForTag = map[string]string{
		gtag.SummaryShort:      gtag.Summary,
		gtag.SummaryShort2:     gtag.Summary,
		gtag.DescriptionShort:  gtag.Description,
		gtag.DescriptionShort2: gtag.Description,
	}
)

// GetRequestRoute 获取当前HTTP请求对应的路由元信息
// 参数 r: 当前请求对象
// 返回值: 对应路由的扩展HTTPRouter对象，不存在则返回nil
func GetRequestRoute(r *ghttp.Request) *HTTPRouter {
	key := GenFilterRequestKey(r)
	routes := LoadHTTPRoutes(r)
	router, ok := routes[key]
	if !ok {
		return nil
	}
	return router
}

// GenFilterRequestKey 根据当前请求对象生成唯一路由key
// 参数 r: 当前请求对象
// 返回值: 格式："请求方法 路由路径"（大写）
func GenFilterRequestKey(r *ghttp.Request) string {
	return GenRouteKey(r.Method, r.Request.URL.Path)
}

// GenFilterRouteKey 根据GF原生路由对象生成唯一路由key
// 参数 r: GF框架路由对象
// 返回值: 格式："请求方法 路由路径"（大写）
func GenFilterRouteKey(r *ghttp.Router) string {
	return GenRouteKey(r.Method, r.Uri)
}

// GenRouteKey 统一生成路由唯一标识key
// 参数 method: HTTP请求方法(GET/POST等)
// 参数 path: 路由地址
// 返回值: 拼接后的唯一key，统一转为大写
func GenRouteKey(method, path string) string {
	return strings.ToUpper(method) + " " + path
}

// LoadHTTPRoutes 加载并缓存服务所有路由元信息
// 首次调用时解析GF服务路由并解析接口元数据，后续直接返回缓存
// 参数 r: 请求对象，用于获取服务路由信息
// 返回值: 全量路由元信息map，key为唯一路由标识
func LoadHTTPRoutes(r *ghttp.Request) map[string]*HTTPRouter {
	if httpRoutes == nil {
		routeMutex.Lock()
		defer routeMutex.Unlock()

		// 双重检查，防止并发重复初始化
		if httpRoutes != nil {
			return httpRoutes
		}

		// 初始化路由缓存map
		httpRoutes = make(map[string]*HTTPRouter, len(r.Server.GetRoutes()))
		// 遍历服务所有路由项，解析并存储元信息
		for _, v := range r.Server.GetRoutes() {
			key := GenFilterRouteKey(v.Handler.Router)
			if _, ok := httpRoutes[key]; !ok {
				router := new(HTTPRouter)
				router.RouterItem = v
				httpRoutes[key] = setRouterMeta(router)
			}
		}
	}
	return httpRoutes
}

// setRouterMeta 为路由对象设置接口元数据（从请求输入结构体解析tags/summary/description）
// 参数 router: 扩展路由对象
// 返回值: 填充完元数据的路由对象
func setRouterMeta(router *HTTPRouter) *HTTPRouter {
	// 非严格路由不解析元数据
	if !router.RouterItem.Handler.Info.IsStrictRoute {
		return router
	}

	// 获取处理器反射对象
	var reflectValue = reflect.ValueOf(router.Handler.Info.Value.Interface())
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}

	// 非函数类型不处理
	if reflectValue.Kind() != reflect.Func {
		return router
	}

	// 校验函数参数/返回值数量（GF标准控制器方法：2入参2出参）
	var reflectType = reflect.TypeOf(router.Handler.Info.Value.Interface())
	if reflectType.NumIn() != 2 || reflectType.NumOut() != 2 {
		return router
	}

	// 构造输入结构体实例，用于获取meta元标签
	var inputObject reflect.Value
	if reflectType.In(1).Kind() == reflect.Ptr {
		inputObject = reflect.New(reflectType.In(1).Elem()).Elem()
	} else {
		inputObject = reflect.New(reflectType.In(1)).Elem()
	}

	// 填充简写标签，解析元数据
	inputMetaMap := fillMapWithShortTags(gmeta.Data(inputObject.Interface()))
	router.Tags = inputMetaMap["tags"]
	router.Summary = inputMetaMap[gtag.Summary]
	router.Description = inputMetaMap[gtag.Description]

	return router
}

// fillMapWithShortTags 填充标签简写映射，将GF简写标签转换为标准标签
// 参数 m: 原始元数据map
// 返回值: 处理后的元数据map
func fillMapWithShortTags(m map[string]string) map[string]string {
	for k, v := range shortTypeMapForTag {
		// 标准标签为空且简写标签存在时，赋值简写内容到标准标签
		if m[v] == "" && m[k] != "" {
			m[v] = m[k]
		}
	}
	return m
}
