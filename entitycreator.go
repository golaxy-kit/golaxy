package golaxy

import (
	"errors"
	"fmt"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/pt"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
)

// NewEntityCreator 创建实体构建器
func NewEntityCreator(runtimeCtx runtime.Context, options ...pt.EntityOption) EntityCreator {
	if runtimeCtx == nil {
		panic("nil runtimeCtx")
	}

	opts := pt.EntityOptions{}
	pt.WithEntityOption{}.Default()(&opts)

	for i := range options {
		options[i](&opts)
	}

	if opts.FaceCache == nil {
		opts.FaceCache = runtimeCtx.GetFaceCache()
	}
	if opts.HookCache == nil {
		opts.HookCache = runtimeCtx.GetHookCache()
	}

	return EntityCreator{
		runtimeCtx: runtimeCtx,
		options:    opts,
		inited:     true,
	}
}

type EntityCreator struct {
	runtimeCtx runtime.Context
	options    pt.EntityOptions
	inited     bool
}

// Spawn 创建实体
func (creator EntityCreator) Spawn() (ec.Entity, error) {
	return creator.spawn(nil)
}

// SpawnWithID 使用指定ID创建实体
func (creator EntityCreator) SpawnWithID(id ec.ID) (ec.Entity, error) {
	return creator.spawn(func(options *pt.EntityOptions) {
		options.PersistID = id
	})
}

func (creator EntityCreator) spawn(modifyOptions func(options *pt.EntityOptions)) (ec.Entity, error) {
	if !creator.inited {
		return nil, errors.New("not inited")
	}

	runtimeCtx := creator.runtimeCtx
	serviceCtx := runtimeCtx.GetServiceCtx()

	entityLib := service.UnsafeContext(serviceCtx).GetOptions().EntityLib
	if entityLib == nil {
		return nil, errors.New("nil entityLib")
	}

	entityPt, ok := entityLib.Get(creator.options.Prototype)
	if !ok {
		return nil, fmt.Errorf("entity '%s' not registered", creator.options.Prototype)
	}

	if modifyOptions != nil {
		modifyOptions(&creator.options)
	}

	entity := entityPt.UnsafeConstruct(creator.options)

	if err := runtimeCtx.GetEntityMgr().AddEntity(entity, creator.options.Accessibility); err != nil {
		return nil, fmt.Errorf("add entity to runtime context failed, %v", err)
	}

	return entity, nil
}
