package simple

import (
	"context"
	"sync"
)

// EventFunc 事件回调函数类型
// 所有注册的事件处理函数都必须遵循这个签名：接收上下文 + 可变参数
type EventFunc func(ctx context.Context, args ...interface{})

// sEvent 事件总线核心结构体
// 负责管理所有事件分组、回调函数，自带互斥锁保证并发安全
type sEvent struct {
	sync.Mutex                        // 嵌入互斥锁，解决并发读写map问题
	list       map[string][]EventFunc // key：事件分组名，value：该分组下所有回调函数
}

// event 全局单例事件总线实例
// 整个程序共用一个事件总线对象
var event *sEvent

// Event 获取全局单例事件总线实例
// 懒加载初始化：第一次调用时才创建实例，节约资源
func Event() *sEvent {
	if event == nil {
		event = &sEvent{
			list: make(map[string][]EventFunc),
		}
	}
	return event
}

// Register 注册事件回调函数
// group：事件分组名称（如：user.register、order.pay）
// callback：事件触发时要执行的函数
func (e *sEvent) Register(group string, callback EventFunc) {
	e.Lock()         // 加锁：防止并发注册导致map并发异常
	defer e.Unlock() // 函数结束自动解锁
	// 将回调函数追加到对应分组的列表中
	e.list[group] = append(e.list[group], callback)
}

// Call 触发指定分组的所有事件回调
// group：要触发的事件分组
// ctx：上下文，用于传递请求信息、超时控制等
// args：触发事件时传递的参数（可变参数，可传任意类型、任意数量）
func (e *sEvent) Call(group string, ctx context.Context, args ...interface{}) {
	// 先判断该分组是否存在注册事件
	if events, ok := e.list[group]; ok {
		// 遍历执行该分组下所有回调函数
		for _, f := range events {
			f(ctx, args...)
		}
	}
}

// Remove 删除整个事件分组
// 会清空该分组下所有已注册的回调函数
func (e *sEvent) Remove(group string) {
	delete(e.list, group)
}

// Clear 清空所有事件分组和回调函数
// 重置为初始空状态
func (e *sEvent) Clear() {
	e.list = make(map[string][]EventFunc)
}
