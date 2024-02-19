package pt

import (
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/event"
	"git.golaxy.org/core/util/container"
	"git.golaxy.org/core/util/generic"
	"git.golaxy.org/core/util/iface"
	"git.golaxy.org/core/util/option"
	"git.golaxy.org/core/util/uid"
)

// Option 所有选项设置器
type Option struct{}

type (
	ComponentCtor = generic.DelegateAction2[ec.Entity, ec.Component] // 组件构造函数
	EntityCtor    = generic.DelegateAction1[ec.Entity]               // 实体构造函数
)

// ConstructEntityOptions 创建实体的所有选项
type ConstructEntityOptions struct {
	ec.EntityOptions
	ComponentCtor ComponentCtor // 组件构造函数
	EntityCtor    EntityCtor    // 实体构造函数
}

// Default 默认值
func (Option) Default() option.Setting[ConstructEntityOptions] {
	return func(o *ConstructEntityOptions) {
		ec.Option{}.Default()(&o.EntityOptions)
		Option{}.ComponentCtor(nil)
		Option{}.EntityCtor(nil)
	}
}

// CompositeFace 扩展者，在扩展实体自身能力时使用
func (Option) CompositeFace(face iface.Face[ec.Entity]) option.Setting[ConstructEntityOptions] {
	return func(o *ConstructEntityOptions) {
		ec.Option{}.CompositeFace(face)(&o.EntityOptions)
	}
}

// Prototype 实体原型名称
func (Option) Prototype(pt string) option.Setting[ConstructEntityOptions] {
	return func(o *ConstructEntityOptions) {
		ec.Option{}.Prototype(pt)(&o.EntityOptions)
	}
}

// Scope 可访问作用域
func (Option) Scope(s ec.Scope) option.Setting[ConstructEntityOptions] {
	return func(o *ConstructEntityOptions) {
		ec.Option{}.Scope(s)(&o.EntityOptions)
	}
}

// PersistId 实体持久化Id
func (Option) PersistId(id uid.Id) option.Setting[ConstructEntityOptions] {
	return func(o *ConstructEntityOptions) {
		ec.Option{}.PersistId(id)(&o.EntityOptions)
	}
}

// AwakeOnFirstAccess 开启组件被首次访问时，检测并调用Awake()
func (Option) AwakeOnFirstAccess(b bool) option.Setting[ConstructEntityOptions] {
	return func(o *ConstructEntityOptions) {
		ec.Option{}.AwakeOnFirstAccess(b)(&o.EntityOptions)
	}
}

// Meta Meta信息
func (Option) Meta(m ec.Meta) option.Setting[ConstructEntityOptions] {
	return func(o *ConstructEntityOptions) {
		ec.Option{}.Meta(m)(&o.EntityOptions)
	}
}

// FaceAnyAllocator 自定义FaceAny内存分配器，用于提高性能，通常传入运行时上下文中的FaceAnyAllocator
func (Option) FaceAnyAllocator(allocator container.Allocator[iface.FaceAny]) option.Setting[ConstructEntityOptions] {
	return func(o *ConstructEntityOptions) {
		ec.Option{}.FaceAnyAllocator(allocator)(&o.EntityOptions)
	}
}

// HookAllocator 自定义Hook内存分配器，用于提高性能，通常传入运行时上下文中的HookAllocator
func (Option) HookAllocator(allocator container.Allocator[event.Hook]) option.Setting[ConstructEntityOptions] {
	return func(o *ConstructEntityOptions) {
		ec.Option{}.HookAllocator(allocator)(&o.EntityOptions)
	}
}

// GCCollector 自定义GC收集器，通常不传或者传入运行时上下文
func (Option) GCCollector(collector container.GCCollector) option.Setting[ConstructEntityOptions] {
	return func(o *ConstructEntityOptions) {
		ec.Option{}.GCCollector(collector)(&o.EntityOptions)
	}
}

// ComponentCtor 组件构造函数
func (Option) ComponentCtor(ctor ComponentCtor) option.Setting[ConstructEntityOptions] {
	return func(o *ConstructEntityOptions) {
		o.ComponentCtor = ctor
	}
}

// EntityCtor 实体构造函数
func (Option) EntityCtor(ctor EntityCtor) option.Setting[ConstructEntityOptions] {
	return func(o *ConstructEntityOptions) {
		o.EntityCtor = ctor
	}
}
