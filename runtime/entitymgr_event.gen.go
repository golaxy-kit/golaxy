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

// Code generated by eventc event; DO NOT EDIT.

package runtime

import (
	event "git.golaxy.org/core/event"
	"git.golaxy.org/core/ec"
)

type iAutoEventEntityMgrAddEntity interface {
	EventEntityMgrAddEntity() event.IEvent
}

func BindEventEntityMgrAddEntity(auto iAutoEventEntityMgrAddEntity, subscriber EventEntityMgrAddEntity, priority ...int32) event.Hook {
	if auto == nil {
		event.Panicf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs)
	}
	return event.Bind[EventEntityMgrAddEntity](auto.EventEntityMgrAddEntity(), subscriber, priority...)
}

func _EmitEventEntityMgrAddEntity(auto iAutoEventEntityMgrAddEntity, entityMgr EntityMgr, entity ec.Entity) {
	if auto == nil {
		event.Panicf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs)
	}
	event.UnsafeEvent(auto.EventEntityMgrAddEntity()).Emit(func(subscriber event.Cache) bool {
		event.Cache2Iface[EventEntityMgrAddEntity](subscriber).OnEntityMgrAddEntity(entityMgr, entity)
		return true
	})
}

func _EmitEventEntityMgrAddEntityWithInterrupt(auto iAutoEventEntityMgrAddEntity, interrupt func(entityMgr EntityMgr, entity ec.Entity) bool, entityMgr EntityMgr, entity ec.Entity) {
	if auto == nil {
		event.Panicf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs)
	}
	event.UnsafeEvent(auto.EventEntityMgrAddEntity()).Emit(func(subscriber event.Cache) bool {
		if interrupt != nil {
			if interrupt(entityMgr, entity) {
				return false
			}
		}
		event.Cache2Iface[EventEntityMgrAddEntity](subscriber).OnEntityMgrAddEntity(entityMgr, entity)
		return true
	})
}

func HandleEventEntityMgrAddEntity(fun func(entityMgr EntityMgr, entity ec.Entity)) EventEntityMgrAddEntityHandler {
	return EventEntityMgrAddEntityHandler(fun)
}

type EventEntityMgrAddEntityHandler func(entityMgr EntityMgr, entity ec.Entity)

func (h EventEntityMgrAddEntityHandler) OnEntityMgrAddEntity(entityMgr EntityMgr, entity ec.Entity) {
	h(entityMgr, entity)
}

type iAutoEventEntityMgrRemoveEntity interface {
	EventEntityMgrRemoveEntity() event.IEvent
}

func BindEventEntityMgrRemoveEntity(auto iAutoEventEntityMgrRemoveEntity, subscriber EventEntityMgrRemoveEntity, priority ...int32) event.Hook {
	if auto == nil {
		event.Panicf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs)
	}
	return event.Bind[EventEntityMgrRemoveEntity](auto.EventEntityMgrRemoveEntity(), subscriber, priority...)
}

func _EmitEventEntityMgrRemoveEntity(auto iAutoEventEntityMgrRemoveEntity, entityMgr EntityMgr, entity ec.Entity) {
	if auto == nil {
		event.Panicf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs)
	}
	event.UnsafeEvent(auto.EventEntityMgrRemoveEntity()).Emit(func(subscriber event.Cache) bool {
		event.Cache2Iface[EventEntityMgrRemoveEntity](subscriber).OnEntityMgrRemoveEntity(entityMgr, entity)
		return true
	})
}

func _EmitEventEntityMgrRemoveEntityWithInterrupt(auto iAutoEventEntityMgrRemoveEntity, interrupt func(entityMgr EntityMgr, entity ec.Entity) bool, entityMgr EntityMgr, entity ec.Entity) {
	if auto == nil {
		event.Panicf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs)
	}
	event.UnsafeEvent(auto.EventEntityMgrRemoveEntity()).Emit(func(subscriber event.Cache) bool {
		if interrupt != nil {
			if interrupt(entityMgr, entity) {
				return false
			}
		}
		event.Cache2Iface[EventEntityMgrRemoveEntity](subscriber).OnEntityMgrRemoveEntity(entityMgr, entity)
		return true
	})
}

func HandleEventEntityMgrRemoveEntity(fun func(entityMgr EntityMgr, entity ec.Entity)) EventEntityMgrRemoveEntityHandler {
	return EventEntityMgrRemoveEntityHandler(fun)
}

type EventEntityMgrRemoveEntityHandler func(entityMgr EntityMgr, entity ec.Entity)

func (h EventEntityMgrRemoveEntityHandler) OnEntityMgrRemoveEntity(entityMgr EntityMgr, entity ec.Entity) {
	h(entityMgr, entity)
}

type iAutoEventEntityMgrEntityAddComponents interface {
	EventEntityMgrEntityAddComponents() event.IEvent
}

func BindEventEntityMgrEntityAddComponents(auto iAutoEventEntityMgrEntityAddComponents, subscriber EventEntityMgrEntityAddComponents, priority ...int32) event.Hook {
	if auto == nil {
		event.Panicf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs)
	}
	return event.Bind[EventEntityMgrEntityAddComponents](auto.EventEntityMgrEntityAddComponents(), subscriber, priority...)
}

func _EmitEventEntityMgrEntityAddComponents(auto iAutoEventEntityMgrEntityAddComponents, entityMgr EntityMgr, entity ec.Entity, components []ec.Component) {
	if auto == nil {
		event.Panicf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs)
	}
	event.UnsafeEvent(auto.EventEntityMgrEntityAddComponents()).Emit(func(subscriber event.Cache) bool {
		event.Cache2Iface[EventEntityMgrEntityAddComponents](subscriber).OnEntityMgrEntityAddComponents(entityMgr, entity, components)
		return true
	})
}

func _EmitEventEntityMgrEntityAddComponentsWithInterrupt(auto iAutoEventEntityMgrEntityAddComponents, interrupt func(entityMgr EntityMgr, entity ec.Entity, components []ec.Component) bool, entityMgr EntityMgr, entity ec.Entity, components []ec.Component) {
	if auto == nil {
		event.Panicf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs)
	}
	event.UnsafeEvent(auto.EventEntityMgrEntityAddComponents()).Emit(func(subscriber event.Cache) bool {
		if interrupt != nil {
			if interrupt(entityMgr, entity, components) {
				return false
			}
		}
		event.Cache2Iface[EventEntityMgrEntityAddComponents](subscriber).OnEntityMgrEntityAddComponents(entityMgr, entity, components)
		return true
	})
}

func HandleEventEntityMgrEntityAddComponents(fun func(entityMgr EntityMgr, entity ec.Entity, components []ec.Component)) EventEntityMgrEntityAddComponentsHandler {
	return EventEntityMgrEntityAddComponentsHandler(fun)
}

type EventEntityMgrEntityAddComponentsHandler func(entityMgr EntityMgr, entity ec.Entity, components []ec.Component)

func (h EventEntityMgrEntityAddComponentsHandler) OnEntityMgrEntityAddComponents(entityMgr EntityMgr, entity ec.Entity, components []ec.Component) {
	h(entityMgr, entity, components)
}

type iAutoEventEntityMgrEntityRemoveComponent interface {
	EventEntityMgrEntityRemoveComponent() event.IEvent
}

func BindEventEntityMgrEntityRemoveComponent(auto iAutoEventEntityMgrEntityRemoveComponent, subscriber EventEntityMgrEntityRemoveComponent, priority ...int32) event.Hook {
	if auto == nil {
		event.Panicf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs)
	}
	return event.Bind[EventEntityMgrEntityRemoveComponent](auto.EventEntityMgrEntityRemoveComponent(), subscriber, priority...)
}

func _EmitEventEntityMgrEntityRemoveComponent(auto iAutoEventEntityMgrEntityRemoveComponent, entityMgr EntityMgr, entity ec.Entity, component ec.Component) {
	if auto == nil {
		event.Panicf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs)
	}
	event.UnsafeEvent(auto.EventEntityMgrEntityRemoveComponent()).Emit(func(subscriber event.Cache) bool {
		event.Cache2Iface[EventEntityMgrEntityRemoveComponent](subscriber).OnEntityMgrEntityRemoveComponent(entityMgr, entity, component)
		return true
	})
}

func _EmitEventEntityMgrEntityRemoveComponentWithInterrupt(auto iAutoEventEntityMgrEntityRemoveComponent, interrupt func(entityMgr EntityMgr, entity ec.Entity, component ec.Component) bool, entityMgr EntityMgr, entity ec.Entity, component ec.Component) {
	if auto == nil {
		event.Panicf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs)
	}
	event.UnsafeEvent(auto.EventEntityMgrEntityRemoveComponent()).Emit(func(subscriber event.Cache) bool {
		if interrupt != nil {
			if interrupt(entityMgr, entity, component) {
				return false
			}
		}
		event.Cache2Iface[EventEntityMgrEntityRemoveComponent](subscriber).OnEntityMgrEntityRemoveComponent(entityMgr, entity, component)
		return true
	})
}

func HandleEventEntityMgrEntityRemoveComponent(fun func(entityMgr EntityMgr, entity ec.Entity, component ec.Component)) EventEntityMgrEntityRemoveComponentHandler {
	return EventEntityMgrEntityRemoveComponentHandler(fun)
}

type EventEntityMgrEntityRemoveComponentHandler func(entityMgr EntityMgr, entity ec.Entity, component ec.Component)

func (h EventEntityMgrEntityRemoveComponentHandler) OnEntityMgrEntityRemoveComponent(entityMgr EntityMgr, entity ec.Entity, component ec.Component) {
	h(entityMgr, entity, component)
}

type iAutoEventEntityMgrEntityFirstAccessComponent interface {
	EventEntityMgrEntityFirstAccessComponent() event.IEvent
}

func BindEventEntityMgrEntityFirstAccessComponent(auto iAutoEventEntityMgrEntityFirstAccessComponent, subscriber EventEntityMgrEntityFirstAccessComponent, priority ...int32) event.Hook {
	if auto == nil {
		event.Panicf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs)
	}
	return event.Bind[EventEntityMgrEntityFirstAccessComponent](auto.EventEntityMgrEntityFirstAccessComponent(), subscriber, priority...)
}

func _EmitEventEntityMgrEntityFirstAccessComponent(auto iAutoEventEntityMgrEntityFirstAccessComponent, entityMgr EntityMgr, entity ec.Entity, component ec.Component) {
	if auto == nil {
		event.Panicf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs)
	}
	event.UnsafeEvent(auto.EventEntityMgrEntityFirstAccessComponent()).Emit(func(subscriber event.Cache) bool {
		event.Cache2Iface[EventEntityMgrEntityFirstAccessComponent](subscriber).OnEntityMgrEntityFirstAccessComponent(entityMgr, entity, component)
		return true
	})
}

func _EmitEventEntityMgrEntityFirstAccessComponentWithInterrupt(auto iAutoEventEntityMgrEntityFirstAccessComponent, interrupt func(entityMgr EntityMgr, entity ec.Entity, component ec.Component) bool, entityMgr EntityMgr, entity ec.Entity, component ec.Component) {
	if auto == nil {
		event.Panicf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs)
	}
	event.UnsafeEvent(auto.EventEntityMgrEntityFirstAccessComponent()).Emit(func(subscriber event.Cache) bool {
		if interrupt != nil {
			if interrupt(entityMgr, entity, component) {
				return false
			}
		}
		event.Cache2Iface[EventEntityMgrEntityFirstAccessComponent](subscriber).OnEntityMgrEntityFirstAccessComponent(entityMgr, entity, component)
		return true
	})
}

func HandleEventEntityMgrEntityFirstAccessComponent(fun func(entityMgr EntityMgr, entity ec.Entity, component ec.Component)) EventEntityMgrEntityFirstAccessComponentHandler {
	return EventEntityMgrEntityFirstAccessComponentHandler(fun)
}

type EventEntityMgrEntityFirstAccessComponentHandler func(entityMgr EntityMgr, entity ec.Entity, component ec.Component)

func (h EventEntityMgrEntityFirstAccessComponentHandler) OnEntityMgrEntityFirstAccessComponent(entityMgr EntityMgr, entity ec.Entity, component ec.Component) {
	h(entityMgr, entity, component)
}
