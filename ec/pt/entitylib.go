/*
 * This file is part of Golaxy Distributed Service Development Framework.
 *
 * Golaxy Distributed Service Development Framework is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * Golaxy Distributed Service Development Framework is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with Golaxy Distributed Service Development Framework. If not, see <http://www.gnu.org/licenses/>.
 *
 * Copyright (c) 2024 pangdogs.
 */

package pt

import (
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/utils/exception"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/core/utils/types"
	"reflect"
	"slices"
	"sync"
)

// EntityLib 实体原型库
type EntityLib interface {
	iEntityLib
	EntityPTProvider

	// GetComponentLib 获取组件原型库
	GetComponentLib() ComponentLib
	// Declare 声明实体原型
	Declare(prototype any, comps ...any) ec.EntityPT
	// Redeclare 重声明实体原型
	Redeclare(prototype any, comps ...any) ec.EntityPT
	// Undeclare 取消声明实体原型
	Undeclare(prototype string)
	// Get 获取实体原型
	Get(prototype string) (ec.EntityPT, bool)
	// Range 遍历所有已注册的实体原型
	Range(fun generic.Func1[ec.EntityPT, bool])
	// ReversedRange 反向遍历所有已注册的实体原型
	ReversedRange(fun generic.Func1[ec.EntityPT, bool])
}

type iEntityLib interface {
	setCallback(declareCB, redeclareCB, undeclareCB generic.Action1[ec.EntityPT])
}

// NewEntityLib 创建实体原型库
func NewEntityLib(compLib ComponentLib) EntityLib {
	if compLib == nil {
		exception.Panicf("%w: %w: compLib is nil", ErrPt, exception.ErrArgs)
	}

	return &_EntityLib{
		compLib:     compLib,
		entityIndex: map[string]*_Entity{},
	}
}

type _EntityLib struct {
	sync.RWMutex
	compLib                             ComponentLib
	entityIndex                         map[string]*_Entity
	entityList                          []*_Entity
	declareCB, redeclareCB, undeclareCB generic.Action1[ec.EntityPT]
}

// GetEntityLib 获取实体原型库
func (lib *_EntityLib) GetEntityLib() EntityLib {
	return lib
}

// GetComponentLib 获取组件原型库
func (lib *_EntityLib) GetComponentLib() ComponentLib {
	return lib.compLib
}

// Declare 声明实体原型
func (lib *_EntityLib) Declare(prototype any, comps ...any) ec.EntityPT {
	entityPT := lib.declare(false, prototype, comps...)
	lib.declareCB.UnsafeCall(entityPT)
	return entityPT
}

// Redeclare 重声明实体原型
func (lib *_EntityLib) Redeclare(prototype any, comps ...any) ec.EntityPT {
	entityPT := lib.declare(true, prototype, comps...)
	lib.redeclareCB.UnsafeCall(entityPT)
	return entityPT
}

// Undeclare 取消声明实体原型
func (lib *_EntityLib) Undeclare(prototype string) {
	entityPT, ok := lib.undeclare(prototype)
	if !ok {
		return
	}
	lib.undeclareCB.UnsafeCall(entityPT)
}

// Get 获取实体原型
func (lib *_EntityLib) Get(prototype string) (ec.EntityPT, bool) {
	lib.RLock()
	defer lib.RUnlock()

	entityPT, ok := lib.entityIndex[prototype]
	if !ok {
		return nil, false
	}

	return entityPT, ok
}

// Range 遍历所有已注册的实体原型
func (lib *_EntityLib) Range(fun generic.Func1[ec.EntityPT, bool]) {
	lib.RLock()
	copied := slices.Clone(lib.entityList)
	lib.RUnlock()

	for i := range copied {
		if !fun.UnsafeCall(copied[i]) {
			return
		}
	}
}

// ReversedRange 反向遍历所有已注册的实体原型
func (lib *_EntityLib) ReversedRange(fun generic.Func1[ec.EntityPT, bool]) {
	lib.RLock()
	copied := slices.Clone(lib.entityList)
	lib.RUnlock()

	for i := len(copied) - 1; i >= 0; i-- {
		if !fun.UnsafeCall(copied[i]) {
			return
		}
	}
}

func (lib *_EntityLib) setCallback(declareCB, redeclareCB, undeclareCB generic.Action1[ec.EntityPT]) {
	lib.Lock()
	defer lib.Unlock()

	lib.declareCB = declareCB
	lib.redeclareCB = redeclareCB
	lib.undeclareCB = undeclareCB
}

func (lib *_EntityLib) declare(re bool, prototype any, comps ...any) ec.EntityPT {
	if prototype == nil {
		exception.Panicf("%w: %w: prototype is nil", ErrPt, exception.ErrArgs)
	}

	if slices.Contains(comps, nil) {
		exception.Panicf("%w: %w: comps contains nil", ErrPt, exception.ErrArgs)
	}

	lib.Lock()
	defer lib.Unlock()

	var entityAtti EntityAttribute

	switch v := prototype.(type) {
	case EntityAttribute:
		entityAtti = v
	case *EntityAttribute:
		entityAtti = *v
	case string:
		entityAtti = EntityAttribute{Prototype: v}
	default:
		exception.Panicf("%w: invalid prototype type: %T", ErrPt, prototype)
	}

	if entityAtti.Prototype == "" {
		exception.Panicf("%w: prototype can't empty", ErrPt)
	}

	entityPT := &_Entity{
		prototype:                  entityAtti.Prototype,
		scope:                      entityAtti.Scope,
		componentNameIndexing:      entityAtti.ComponentNameIndexing,
		componentAwakeOnFirstTouch: entityAtti.ComponentAwakeOnFirstTouch,
		componentUniqueID:          entityAtti.ComponentUniqueID,
		extra:                      entityAtti.Extra,
	}

	if entityAtti.Instance != nil {
		instanceRT, ok := entityAtti.Instance.(reflect.Type)
		if !ok {
			instanceRT = reflect.TypeOf(entityAtti.Instance)
		}

		for instanceRT.Kind() == reflect.Pointer {
			instanceRT = instanceRT.Elem()
		}

		if instanceRT.Name() == "" {
			exception.Panicf("%w: anonymous entity instance not allowed", ErrPt)
		}

		if !reflect.PointerTo(instanceRT).Implements(reflect.TypeFor[ec.Entity]()) {
			exception.Panicf("%w: entity instance %q not implement ec.Entity", ErrPt, types.FullNameRT(instanceRT))
		}

		entityPT.instanceRT = instanceRT
	}

	for i, comp := range comps {
		builtin := ec.BuiltinComponent{
			Offset: i,
		}

	retry:
		switch v := comp.(type) {
		case ComponentAttribute:
			builtin.Name = v.Name
			builtin.Removable = v.Removable
			builtin.Extra = v.Extra
			comp = v.Instance
			goto retry
		case *ComponentAttribute:
			comp = *v
			goto retry
		case string:
			compPT, ok := lib.compLib.Get(v)
			if !ok {
				exception.Panicf("%w: entity %q builtin component %q was not declared", ErrPt, prototype, v)
			}
			builtin.PT = compPT
		default:
			if v == nil {
				exception.Panicf("%w: entity %q builtin component is nil", ErrPt, prototype)
			}
			builtin.PT = lib.compLib.Declare(v)
		}

		if builtin.Name == "" {
			builtin.Name = types.NameRT(builtin.PT.InstanceRT().Elem())
		}

		entityPT.components = append(entityPT.components, builtin)
	}

	if _, ok := lib.entityIndex[entityAtti.Prototype]; ok {
		if re {
			lib.entityList = slices.DeleteFunc(lib.entityList, func(pt *_Entity) bool {
				return pt.prototype == prototype
			})
		} else {
			exception.Panicf("%w: entity %q is already declared", ErrPt, prototype)
		}
	}

	lib.entityIndex[entityAtti.Prototype] = entityPT
	lib.entityList = append(lib.entityList, entityPT)

	return entityPT
}

func (lib *_EntityLib) undeclare(prototype string) (ec.EntityPT, bool) {
	lib.Lock()
	defer lib.Unlock()

	entityPT, ok := lib.entityIndex[prototype]
	if !ok {
		return nil, false
	}

	delete(lib.entityIndex, prototype)

	lib.entityList = slices.DeleteFunc(lib.entityList, func(pt *_Entity) bool {
		return pt.prototype == prototype
	})

	return entityPT, true
}
