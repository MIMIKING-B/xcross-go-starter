package contexts

import (
	"context"
	"time"
)

// detached 自定义结构体，实现了 context.Context 接口
// 作用：包装一个已有的上下文 ctx，**继承它的键值对数据，但切断取消/超时信号**
type detached struct {
	ctx context.Context // 持有原始父上下文，用于读取键值对
}

// Deadline 实现 context.Context 接口
// 返回：无截止时间，永远不会因为时间到期自动取消
// 原因：我们需要一个“永远不会超时/自动取消”的上下文
func (detached) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

// Done 实现 context.Context 接口
// 返回：nil 通道，nil 通道永远不会被关闭，永远不会收到取消信号
// 原因：让这个上下文**完全不响应父上下文的取消操作**
// 标准 ctx.Done() 关闭时，当前 detached 不会有任何反应
func (detached) Done() <-chan struct{} {
	return nil
}

// Err 实现 context.Context 接口
// 返回：nil 错误，永远不会标记为已取消/超时
// 原因：配合 Done()，让外部判断时永远认为上下文是正常未取消状态
func (detached) Err() error {
	return nil
}

// Value 实现 context.Context 接口
// 作用：**从原始父上下文里读取键值数据**（如 traceID、userID、requestID）
// 原因：我们只想要取消隔离，不想要丢失链路追踪、用户信息等重要数据
func (d detached) Value(key interface{}) interface{} {
	return d.ctx.Value(key)
}

// Detach 对外工具方法：将一个上下文“脱离取消链”
// newCtx := contexts.Detach(oldCtx)
// 输入：父上下文（可能会被取消、超时）
// 输出：新上下文 = 继承数据 + 不会被取消/超时
// 用途：HTTP 请求结束后、客户端断开后，仍需要后台继续执行的异步任务
func Detach(ctx context.Context) context.Context {
	return detached{ctx: ctx}
}
