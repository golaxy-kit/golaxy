package golaxy

import (
	"fmt"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/event"
	"kit.golaxy.org/golaxy/internal"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/util/iface"
	"kit.golaxy.org/golaxy/util/uid"
)

// NewRuntime 创建运行时
func NewRuntime(ctx runtime.Context, options ...RuntimeOption) Runtime {
	opts := RuntimeOptions{}
	Option{}.Runtime.Default()(&opts)

	for i := range options {
		options[i](&opts)
	}

	return UnsafeNewRuntime(ctx, opts)
}

// Deprecated: UnsafeNewRuntime 内部创建运行时
func UnsafeNewRuntime(ctx runtime.Context, options RuntimeOptions) Runtime {
	if !options.CompositeFace.IsNil() {
		options.CompositeFace.Iface.init(ctx, &options)
		return options.CompositeFace.Iface
	}

	runtime := &RuntimeBehavior{}
	runtime.init(ctx, &options)

	return runtime.opts.CompositeFace.Iface
}

// Runtime 运行时接口
type Runtime interface {
	_Runtime
	internal.ContextResolver
	internal.Running

	// GetContext 获取运行时上下文
	GetContext() runtime.Context
}

type _Runtime interface {
	init(ctx runtime.Context, opts *RuntimeOptions)
	getOptions() *RuntimeOptions
}

// RuntimeBehavior 运行时行为，在需要扩展运行时能力时，匿名嵌入至运行时结构体中
type RuntimeBehavior struct {
	opts            RuntimeOptions
	ctx             runtime.Context
	hooksMap        map[uid.Id][3]event.Hook
	processQueue    chan func()
	eventUpdate     event.Event
	eventLateUpdate event.Event
}

// GetContext 获取运行时上下文
func (_runtime *RuntimeBehavior) GetContext() runtime.Context {
	return _runtime.ctx
}

// ResolveContext 解析上下文
func (_runtime *RuntimeBehavior) ResolveContext() iface.Cache {
	return iface.Iface2Cache[runtime.Context](_runtime.ctx)
}

func (_runtime *RuntimeBehavior) init(ctx runtime.Context, opts *RuntimeOptions) {
	if ctx == nil {
		panic(fmt.Errorf("%w: %w: ctx is nil", ErrRuntime, ErrArgs))
	}

	if opts == nil {
		panic(fmt.Errorf("%w: %w: opts is nil", ErrRuntime, ErrArgs))
	}

	if !runtime.UnsafeContext(ctx).MarkPaired(true) {
		panic(fmt.Errorf("%w: ctx already paired", ErrRuntime))
	}

	_runtime.opts = *opts

	if _runtime.opts.CompositeFace.IsNil() {
		_runtime.opts.CompositeFace = iface.NewFace[Runtime](_runtime)
	}

	_runtime.ctx = ctx
	_runtime.hooksMap = make(map[uid.Id][3]event.Hook)

	_runtime.eventUpdate.Init(ctx.GetAutoRecover(), ctx.GetReportError(), event.EventRecursion_Disallow, ctx.GetHookAllocator(), nil)
	_runtime.eventLateUpdate.Init(ctx.GetAutoRecover(), ctx.GetReportError(), event.EventRecursion_Disallow, ctx.GetHookAllocator(), nil)

	_runtime.changeRunningState(runtime.RunningState_Birth)

	if _runtime.opts.AutoRun {
		_runtime.opts.CompositeFace.Iface.Run()
	}
}

func (_runtime *RuntimeBehavior) getOptions() *RuntimeOptions {
	return &_runtime.opts
}

// OnEntityMgrAddEntity 事件处理器：实体管理器添加实体
func (_runtime *RuntimeBehavior) OnEntityMgrAddEntity(entityMgr runtime.IEntityMgr, entity ec.Entity) {
	_runtime.connectEntity(entity)
	_runtime.initEntity(entity)
}

// OnEntityMgrRemoveEntity 事件处理器：实体管理器删除实体
func (_runtime *RuntimeBehavior) OnEntityMgrRemoveEntity(entityMgr runtime.IEntityMgr, entity ec.Entity) {
	_runtime.disconnectEntity(entity)
	_runtime.shutEntity(entity)
}

// OnEntityMgrEntityFirstAccessComponent 事件处理器：实体管理器中的实体首次访问组件
func (_runtime *RuntimeBehavior) OnEntityMgrEntityFirstAccessComponent(entityMgr runtime.IEntityMgr, entity ec.Entity, component ec.Component) {
	_comp := ec.UnsafeComponent(component)

	if _comp.GetState() != ec.ComponentState_Attach {
		return
	}

	_comp.SetState(ec.ComponentState_Awake)

	if compAwake, ok := component.(LifecycleComponentAwake); ok {
		internal.CallOuterVoid(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
			compAwake.Awake()
		})
	}

	_comp.SetState(ec.ComponentState_Start)
}

// OnEntityMgrEntityAddComponents 事件处理器：实体管理器中的实体添加组件
func (_runtime *RuntimeBehavior) OnEntityMgrEntityAddComponents(entityMgr runtime.IEntityMgr, entity ec.Entity, components []ec.Component) {
	_runtime.addComponents(entity, components)
}

// OnEntityMgrEntityRemoveComponent 事件处理器：实体管理器中的实体删除组件
func (_runtime *RuntimeBehavior) OnEntityMgrEntityRemoveComponent(entityMgr runtime.IEntityMgr, entity ec.Entity, component ec.Component) {
	_runtime.removeComponent(component)
}

// OnEntityDestroySelf 事件处理器：实体销毁自身
func (_runtime *RuntimeBehavior) OnEntityDestroySelf(entity ec.Entity) {
	_runtime.ctx.GetEntityMgr().RemoveEntity(entity.GetId())
}

// OnComponentDestroySelf 事件处理器：组件销毁自身
func (_runtime *RuntimeBehavior) OnComponentDestroySelf(comp ec.Component) {
	comp.GetEntity().RemoveComponentById(comp.GetId())
}

func (_runtime *RuntimeBehavior) addComponents(entity ec.Entity, components []ec.Component) {
	switch entity.GetState() {
	case ec.EntityState_Init, ec.EntityState_Inited, ec.EntityState_Living:
	default:
		return
	}

	for i := range components {
		_runtime.connectComponent(components[i])
	}

	for i := range components {
		_comp := ec.UnsafeComponent(components[i])

		if _comp.GetState() != ec.ComponentState_Awake {
			continue
		}

		if compAwake, ok := components[i].(LifecycleComponentAwake); ok {
			internal.CallOuterVoid(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
				compAwake.Awake()
			})
		}

		_comp.SetState(ec.ComponentState_Start)
	}

	switch entity.GetState() {
	case ec.EntityState_Init, ec.EntityState_Inited, ec.EntityState_Living:
	default:
		return
	}

	for i := range components {
		_comp := ec.UnsafeComponent(components[i])

		if _comp.GetState() != ec.ComponentState_Start {
			continue
		}

		if compStart, ok := components[i].(LifecycleComponentStart); ok {
			internal.CallOuterVoid(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
				compStart.Start()
			})
		}

		_comp.SetState(ec.ComponentState_Living)
	}
}

func (_runtime *RuntimeBehavior) removeComponent(component ec.Component) {
	_runtime.disconnectComponent(component)

	if component.GetState() != ec.ComponentState_Shut {
		return
	}

	if compShut, ok := component.(LifecycleComponentShut); ok {
		internal.CallOuterVoid(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
			compShut.Shut()
		})
	}

	ec.UnsafeComponent(component).SetState(ec.ComponentState_Death)
}

func (_runtime *RuntimeBehavior) connectEntity(entity ec.Entity) {
	if entity.GetState() != ec.EntityState_Entry {
		return
	}

	var hooks [3]event.Hook

	if entityUpdate, ok := entity.(LifecycleEntityUpdate); ok {
		hooks[0] = event.BindEvent[LifecycleEntityUpdate](&_runtime.eventUpdate, entityUpdate)
	}
	if entityLateUpdate, ok := entity.(LifecycleEntityLateUpdate); ok {
		hooks[1] = event.BindEvent[LifecycleEntityLateUpdate](&_runtime.eventLateUpdate, entityLateUpdate)
	}
	hooks[2] = event.BindEvent[ec.EventEntityDestroySelf](ec.UnsafeEntity(entity).EventEntityDestroySelf(), _runtime)

	_runtime.hooksMap[entity.GetId()] = hooks

	entity.RangeComponents(func(comp ec.Component) bool {
		_runtime.connectComponent(comp)
		return true
	})

	ec.UnsafeEntity(entity).SetState(ec.EntityState_Init)
}

func (_runtime *RuntimeBehavior) disconnectEntity(entity ec.Entity) {
	entity.RangeComponents(func(comp ec.Component) bool {
		_runtime.disconnectComponent(comp)
		return true
	})

	entityId := entity.GetId()

	hooks, ok := _runtime.hooksMap[entityId]
	if ok {
		delete(_runtime.hooksMap, entityId)

		for i := range hooks {
			hooks[i].Unbind()
		}
	}

	ec.UnsafeEntity(entity).SetState(ec.EntityState_Shut)
}

func (_runtime *RuntimeBehavior) connectComponent(comp ec.Component) {
	if comp.GetState() != ec.ComponentState_Attach {
		return
	}

	var hooks [3]event.Hook

	if compUpdate, ok := comp.(LifecycleComponentUpdate); ok {
		hooks[0] = event.BindEvent[LifecycleComponentUpdate](&_runtime.eventUpdate, compUpdate)
	}
	if compLateUpdate, ok := comp.(LifecycleComponentLateUpdate); ok {
		hooks[1] = event.BindEvent[LifecycleComponentLateUpdate](&_runtime.eventLateUpdate, compLateUpdate)
	}
	hooks[2] = event.BindEvent[ec.EventComponentDestroySelf](ec.UnsafeComponent(comp).EventComponentDestroySelf(), _runtime)

	_runtime.hooksMap[comp.GetId()] = hooks

	ec.UnsafeComponent(comp).SetState(ec.ComponentState_Awake)
}

func (_runtime *RuntimeBehavior) disconnectComponent(comp ec.Component) {
	compId := comp.GetId()

	hooks, ok := _runtime.hooksMap[compId]
	if ok {
		delete(_runtime.hooksMap, compId)

		for i := range hooks {
			hooks[i].Unbind()
		}
	}

	ec.UnsafeComponent(comp).SetState(ec.ComponentState_Shut)
}

func (_runtime *RuntimeBehavior) initEntity(entity ec.Entity) {
	if entity.GetState() != ec.EntityState_Init {
		return
	}

	if entityInit, ok := entity.(LifecycleEntityInit); ok {
		internal.CallOuterVoid(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
			entityInit.Init()
		})
	}

	if entity.GetState() != ec.EntityState_Init {
		return
	}

	entity.RangeComponents(func(comp ec.Component) bool {
		_comp := ec.UnsafeComponent(comp)

		if _comp.GetState() != ec.ComponentState_Awake {
			return true
		}

		if compAwake, ok := comp.(LifecycleComponentAwake); ok {
			internal.CallOuterVoid(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
				compAwake.Awake()
			})
		}

		_comp.SetState(ec.ComponentState_Start)

		return entity.GetState() == ec.EntityState_Init
	})

	if entity.GetState() != ec.EntityState_Init {
		return
	}

	entity.RangeComponents(func(comp ec.Component) bool {
		_comp := ec.UnsafeComponent(comp)

		if _comp.GetState() != ec.ComponentState_Start {
			return true
		}

		if compStart, ok := comp.(LifecycleComponentStart); ok {
			internal.CallOuterVoid(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
				compStart.Start()
			})
		}

		_comp.SetState(ec.ComponentState_Living)

		return entity.GetState() == ec.EntityState_Init
	})

	if entity.GetState() != ec.EntityState_Init {
		return
	}

	ec.UnsafeEntity(entity).SetState(ec.EntityState_Inited)

	if entityInited, ok := entity.(LifecycleEntityInited); ok {
		internal.CallOuterVoid(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
			entityInited.Inited()
		})
	}

	if entity.GetState() != ec.EntityState_Inited {
		return
	}

	ec.UnsafeEntity(entity).SetState(ec.EntityState_Living)
}

func (_runtime *RuntimeBehavior) shutEntity(entity ec.Entity) {
	if entity.GetState() != ec.EntityState_Shut {
		return
	}

	if entityShut, ok := entity.(LifecycleEntityShut); ok {
		internal.CallOuterVoid(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
			entityShut.Shut()
		})
	}

	entity.RangeComponents(func(comp ec.Component) bool {
		_comp := ec.UnsafeComponent(comp)

		if _comp.GetState() != ec.ComponentState_Shut {
			return true
		}

		if compShut, ok := comp.(LifecycleComponentShut); ok {
			internal.CallOuterVoid(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
				compShut.Shut()
			})
		}

		_comp.SetState(ec.ComponentState_Death)

		return true
	})

	ec.UnsafeEntity(entity).SetState(ec.EntityState_Death)

	if entityDestroy, ok := entity.(LifecycleEntityDestroy); ok {
		internal.CallOuterVoid(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
			entityDestroy.Destroy()
		})
	}
}
