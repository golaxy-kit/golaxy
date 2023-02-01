package golaxy

import (
	"github.com/golaxy-kit/golaxy/ec"
	"github.com/golaxy-kit/golaxy/internal"
	"github.com/golaxy-kit/golaxy/localevent"
	"github.com/golaxy-kit/golaxy/runtime"
	"github.com/golaxy-kit/golaxy/util"
)

// NewRuntime 创建运行时
func NewRuntime(runtimeCtx runtime.Context, options ...RuntimeOption) Runtime {
	opts := RuntimeOptions{}
	WithRuntimeOption{}.Default()(&opts)

	for i := range options {
		options[i](&opts)
	}

	return UnsafeNewRuntime(runtimeCtx, opts)
}

func UnsafeNewRuntime(runtimeCtx runtime.Context, options RuntimeOptions) Runtime {
	if !options.Inheritor.IsNil() {
		options.Inheritor.Iface.init(runtimeCtx, &options)
		return options.Inheritor.Iface
	}

	runtime := &RuntimeBehavior{}
	runtime.init(runtimeCtx, &options)

	return runtime.opts.Inheritor.Iface
}

// Runtime 运行时接口
type Runtime interface {
	_Runtime
	internal.Running

	// GetRuntimeCtx 获取运行时上下文
	GetRuntimeCtx() runtime.Context
}

type _Runtime interface {
	init(runtimeCtx runtime.Context, opts *RuntimeOptions)
	getOptions() *RuntimeOptions
}

type _HookKey struct {
	ID ec.ID
	SN int64
}

// RuntimeBehavior 运行时行为，在需要扩展运行时能力时，匿名嵌入至运行时结构体中
type RuntimeBehavior struct {
	opts            RuntimeOptions
	ctx             runtime.Context
	hooksMap        map[_HookKey][3]localevent.Hook
	processQueue    chan func()
	eventUpdate     localevent.Event
	eventLateUpdate localevent.Event
}

// GetRuntimeCtx 获取运行时上下文
func (_runtime *RuntimeBehavior) GetRuntimeCtx() runtime.Context {
	return _runtime.ctx
}

func (_runtime *RuntimeBehavior) init(runtimeCtx runtime.Context, opts *RuntimeOptions) {
	if runtimeCtx == nil {
		panic("nil runtimeCtx")
	}

	if opts == nil {
		panic("nil opts")
	}

	if !internal.UnsafeContext(runtimeCtx).Paired() {
		panic("runtime context already paired")
	}

	_runtime.opts = *opts

	if _runtime.opts.Inheritor.IsNil() {
		_runtime.opts.Inheritor = util.NewFace[Runtime](_runtime)
	}

	_runtime.ctx = runtimeCtx
	_runtime.hooksMap = make(map[_HookKey][3]localevent.Hook)

	_runtime.eventUpdate.Init(runtimeCtx.GetAutoRecover(), runtimeCtx.GetReportError(), localevent.EventRecursion_Disallow, runtimeCtx.GetHookCache(), _runtime)
	_runtime.eventLateUpdate.Init(runtimeCtx.GetAutoRecover(), runtimeCtx.GetReportError(), localevent.EventRecursion_Disallow, runtimeCtx.GetHookCache(), _runtime)

	if opts.EnableAutoRun {
		_runtime.opts.Inheritor.Iface.Run()
	}
}

func (_runtime *RuntimeBehavior) getOptions() *RuntimeOptions {
	return &_runtime.opts
}

// OnEntityMgrAddEntity 事件回调：实体管理器添加实体
func (_runtime *RuntimeBehavior) OnEntityMgrAddEntity(entityMgr runtime.IEntityMgr, entity ec.Entity) {
	_runtime.connectEntity(entity)
	_runtime.initEntity(entity)
}

// OnEntityMgrRemoveEntity 事件回调：实体管理器删除实体
func (_runtime *RuntimeBehavior) OnEntityMgrRemoveEntity(entityMgr runtime.IEntityMgr, entity ec.Entity) {
	_runtime.disconnectEntity(entity)
	_runtime.shutEntity(entity)
}

// OnEntityMgrEntityFirstAccessComponent 事件回调：实体管理器中的实体首次访问组件
func (_runtime *RuntimeBehavior) OnEntityMgrEntityFirstAccessComponent(entityMgr runtime.IEntityMgr, entity ec.Entity, component ec.Component) {
	_comp := ec.UnsafeComponent(component)

	if _comp.GetState() != ec.ComponentState_Attach {
		return
	}

	_comp.SetState(ec.ComponentState_Awake)

	if compAwake, ok := component.(_ComponentAwake); ok {
		internal.CallOuterNoRet(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
			compAwake.Awake()
		})
	}

	_comp.SetState(ec.ComponentState_Start)
}

// OnEntityMgrEntityAddComponents 事件回调：实体管理器中的实体添加组件
func (_runtime *RuntimeBehavior) OnEntityMgrEntityAddComponents(entityMgr runtime.IEntityMgr, entity ec.Entity, components []ec.Component) {
	_runtime.addComponents(entity, components)
}

// OnEntityMgrEntityRemoveComponent 事件回调：实体管理器中的实体删除组件
func (_runtime *RuntimeBehavior) OnEntityMgrEntityRemoveComponent(entityMgr runtime.IEntityMgr, entity ec.Entity, component ec.Component) {
	_runtime.removeComponent(component)
}

// OnEntityDestroySelf 事件回调：实体销毁自身
func (_runtime *RuntimeBehavior) OnEntityDestroySelf(entity ec.Entity) {
	_runtime.ctx.GetEntityMgr().RemoveEntity(entity.GetID())
}

// OnComponentDestroySelf 事件回调：组件销毁自身
func (_runtime *RuntimeBehavior) OnComponentDestroySelf(comp ec.Component) {
	comp.GetEntity().RemoveComponentByID(comp.GetID())
}

func (_runtime *RuntimeBehavior) addComponents(entity ec.Entity, components []ec.Component) {
	if entity.GetState() != ec.EntityState_Init && entity.GetState() != ec.EntityState_Living {
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

		if compAwake, ok := components[i].(_ComponentAwake); ok {
			internal.CallOuterNoRet(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
				compAwake.Awake()
			})
		}

		_comp.SetState(ec.ComponentState_Start)
	}

	if entity.GetState() != ec.EntityState_Init && entity.GetState() != ec.EntityState_Living {
		return
	}

	for i := range components {
		_comp := ec.UnsafeComponent(components[i])

		if _comp.GetState() != ec.ComponentState_Start {
			continue
		}

		if compStart, ok := components[i].(_ComponentStart); ok {
			internal.CallOuterNoRet(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
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

	if compShut, ok := component.(_ComponentShut); ok {
		internal.CallOuterNoRet(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
			compShut.Shut()
		})
	}

	ec.UnsafeComponent(component).SetState(ec.ComponentState_Death)
}

func (_runtime *RuntimeBehavior) connectEntity(entity ec.Entity) {
	if entity.GetState() != ec.EntityState_Entry {
		return
	}

	var hooks [3]localevent.Hook

	if entityUpdate, ok := entity.(_EntityUpdate); ok {
		hooks[0] = localevent.BindEvent[_EntityUpdate](&_runtime.eventUpdate, entityUpdate)
	}
	if entityLateUpdate, ok := entity.(_EntityLateUpdate); ok {
		hooks[1] = localevent.BindEvent[_EntityLateUpdate](&_runtime.eventLateUpdate, entityLateUpdate)
	}
	hooks[2] = localevent.BindEvent[ec.EventEntityDestroySelf](ec.UnsafeEntity(entity).EventEntityDestroySelf(), _runtime)

	_runtime.hooksMap[_HookKey{
		ID: entity.GetID(),
		SN: entity.GetSerialNo(),
	}] = hooks

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

	hookKey := _HookKey{
		ID: entity.GetID(),
		SN: entity.GetSerialNo(),
	}

	hooks, ok := _runtime.hooksMap[hookKey]
	if ok {
		delete(_runtime.hooksMap, hookKey)

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

	var hooks [3]localevent.Hook

	if compUpdate, ok := comp.(_ComponentUpdate); ok {
		hooks[0] = localevent.BindEvent[_ComponentUpdate](&_runtime.eventUpdate, compUpdate)
	}
	if compLateUpdate, ok := comp.(_ComponentLateUpdate); ok {
		hooks[1] = localevent.BindEvent[_ComponentLateUpdate](&_runtime.eventLateUpdate, compLateUpdate)
	}
	hooks[2] = localevent.BindEvent[ec.EventComponentDestroySelf](ec.UnsafeComponent(comp).EventComponentDestroySelf(), _runtime)

	_runtime.hooksMap[_HookKey{
		ID: comp.GetID(),
		SN: comp.GetSerialNo(),
	}] = hooks

	ec.UnsafeComponent(comp).SetState(ec.ComponentState_Awake)
}

func (_runtime *RuntimeBehavior) disconnectComponent(comp ec.Component) {
	hookKey := _HookKey{
		ID: comp.GetID(),
		SN: comp.GetSerialNo(),
	}

	hooks, ok := _runtime.hooksMap[hookKey]
	if ok {
		delete(_runtime.hooksMap, hookKey)

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

	if entityInit, ok := entity.(_EntityInit); ok {
		internal.CallOuterNoRet(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
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

		if compAwake, ok := comp.(_ComponentAwake); ok {
			internal.CallOuterNoRet(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
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

		if compStart, ok := comp.(_ComponentStart); ok {
			internal.CallOuterNoRet(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
				compStart.Start()
			})
		}

		_comp.SetState(ec.ComponentState_Living)

		return entity.GetState() == ec.EntityState_Init
	})

	if entity.GetState() != ec.EntityState_Init {
		return
	}

	if entityInitFin, ok := entity.(_EntityInitFin); ok {
		internal.CallOuterNoRet(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
			entityInitFin.InitFin()
		})
	}

	if entity.GetState() != ec.EntityState_Init {
		return
	}

	ec.UnsafeEntity(entity).SetState(ec.EntityState_Living)
}

func (_runtime *RuntimeBehavior) shutEntity(entity ec.Entity) {
	if entity.GetState() != ec.EntityState_Shut {
		return
	}

	if entityShut, ok := entity.(_EntityShut); ok {
		internal.CallOuterNoRet(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
			entityShut.Shut()
		})
	}

	entity.RangeComponents(func(comp ec.Component) bool {
		_comp := ec.UnsafeComponent(comp)

		if _comp.GetState() != ec.ComponentState_Shut {
			return true
		}

		if compShut, ok := comp.(_ComponentShut); ok {
			internal.CallOuterNoRet(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
				compShut.Shut()
			})
		}

		_comp.SetState(ec.ComponentState_Death)

		return true
	})

	if entityShutFin, ok := entity.(_EntityShutFin); ok {
		internal.CallOuterNoRet(_runtime.ctx.GetAutoRecover(), _runtime.ctx.GetReportError(), func() {
			entityShutFin.ShutFin()
		})
	}

	ec.UnsafeEntity(entity).SetState(ec.EntityState_Death)
}
