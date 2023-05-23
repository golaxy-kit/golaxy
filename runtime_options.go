package golaxy

import (
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/util"
	"time"
)

// RuntimeOptions 创建运行时的所有选项
type RuntimeOptions struct {
	CompositeFace        util.Face[Runtime] // 扩展者，需要扩展运行时自身功能时需要使用
	EnableAutoRun        bool               // 是否开启自动运行
	ProcessQueueCapacity int                // 任务处理流水线大小
	ProcessQueueTimeout  time.Duration      // 当任务处理流水线满时，向其插入代码片段的超时时间，为0表示不等待直接报错
	SyncCallTimeout      time.Duration      // 同步调用超时时间，为0表示不处理超时，此时两个运行时互相同步调用会死锁
	Frame                runtime.Frame      // 帧，设置为nil表示不使用帧更新特性
	GCInterval           time.Duration      // GC间隔时长
	CustomGC             func(rt Runtime)   // 自定义GC
}

// RuntimeOption 创建运行时的选项设置器
type RuntimeOption func(o *RuntimeOptions)

// RuntimeDefault 运行时的默认值
func (WithOption) RuntimeDefault() RuntimeOption {
	return func(o *RuntimeOptions) {
		WithOption{}.RuntimeCompositeFace(util.Face[Runtime]{})(o)
		WithOption{}.RuntimeEnableAutoRun(false)(o)
		WithOption{}.RuntimeProcessQueueCapacity(128)(o)
		WithOption{}.RuntimeProcessQueueTimeout(0)(o)
		WithOption{}.RuntimeSyncCallTimeout(3 * time.Second)(o)
		WithOption{}.RuntimeFrame(nil)(o)
		WithOption{}.RuntimeGCInterval(10 * time.Second)(o)
		WithOption{}.RuntimeCustomGC(nil)(o)
	}
}

// RuntimeCompositeFace 运行时的扩展者，需要扩展运行时自身功能时需要使用
func (WithOption) RuntimeCompositeFace(face util.Face[Runtime]) RuntimeOption {
	return func(o *RuntimeOptions) {
		o.CompositeFace = face
	}
}

// RuntimeEnableAutoRun 运行时是否开启自动运行
func (WithOption) RuntimeEnableAutoRun(b bool) RuntimeOption {
	return func(o *RuntimeOptions) {
		o.EnableAutoRun = b
	}
}

// RuntimeProcessQueueCapacity 运行时的任务处理流水线大小
func (WithOption) RuntimeProcessQueueCapacity(cap int) RuntimeOption {
	return func(o *RuntimeOptions) {
		if cap <= 0 {
			panic("RuntimeProcessQueueCapacity less equal 0 is invalid")
		}
		o.ProcessQueueCapacity = cap
	}
}

// RuntimeProcessQueueTimeout 运行时的当任务处理流水线满时，向其插入代码片段的超时时间，为0表示不等待直接报错
func (WithOption) RuntimeProcessQueueTimeout(dur time.Duration) RuntimeOption {
	return func(o *RuntimeOptions) {
		o.ProcessQueueTimeout = dur
	}
}

// RuntimeSyncCallTimeout 运行时的同步调用超时时间，为0表示不处理超时，此时两个运行时互相同步调用会死锁
func (WithOption) RuntimeSyncCallTimeout(dur time.Duration) RuntimeOption {
	return func(o *RuntimeOptions) {
		o.SyncCallTimeout = dur
	}
}

// RuntimeFrame 运行时的帧，设置为nil表示不使用帧更新特性
func (WithOption) RuntimeFrame(frame runtime.Frame) RuntimeOption {
	return func(o *RuntimeOptions) {
		o.Frame = frame
	}
}

// RuntimeGCInterval 运行时的GC间隔时长
func (WithOption) RuntimeGCInterval(dur time.Duration) RuntimeOption {
	return func(o *RuntimeOptions) {
		if dur <= 0 {
			panic("RuntimeGCInterval less equal 0 is invalid")
		}
		o.GCInterval = dur
	}
}

// RuntimeCustomGC 运行时的自定义GC
func (WithOption) RuntimeCustomGC(fn func(rt Runtime)) RuntimeOption {
	return func(o *RuntimeOptions) {
		o.CustomGC = fn
	}
}
