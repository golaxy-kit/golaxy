package runtime

import (
	"kit.golaxy.org/golaxy/internal"
	"kit.golaxy.org/golaxy/service"
)

// Deprecated: UnsafeContext 访问运行时上下文内部方法
func UnsafeContext(ctx Context) _UnsafeContext {
	return _UnsafeContext{
		Context: ctx,
	}
}

type _UnsafeContext struct {
	Context
}

// Init 初始化
func (uc _UnsafeContext) Init(serviceCtx service.Context, opts ContextOptions) {
	uc.Context.init(serviceCtx, opts)
}

// GetOptions 获取运行时上下文所有选项
func (uc _UnsafeContext) GetOptions() *ContextOptions {
	return uc.getOptions()
}

// SetFrame 设置帧
func (uc _UnsafeContext) SetFrame(frame Frame) {
	uc.setFrame(frame)
}

// SetCallee 设置调用接受者
func (uc _UnsafeContext) SetCallee(callee Callee) {
	uc.setCallee(callee)
}

// GetServiceCtx 获取服务上下文
func (uc _UnsafeContext) GetServiceCtx() service.Context {
	return uc.getServiceCtx()
}

// GC GC
func (uc _UnsafeContext) GC() {
	uc.gc()
}

// MarkRunning 标记运行时已经开始运行
func (uc _UnsafeContext) MarkRunning(v bool) bool {
	return internal.UnsafeRunningState(uc.Context).MarkRunning(v)
}

// MarkPaired 标记与运行时已经配对
func (uc _UnsafeContext) MarkPaired(v bool) bool {
	return internal.UnsafeContext(uc.Context).MarkPaired(v)
}
