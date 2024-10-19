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

package ec

import (
	event "git.golaxy.org/core/event"
)

type iAutoEventEntityDestroySelf interface {
	EventEntityDestroySelf() event.IEvent
}

func BindEventEntityDestroySelf(auto iAutoEventEntityDestroySelf, subscriber EventEntityDestroySelf, priority ...int32) event.Hook {
	if auto == nil {
		event.Panicf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs)
	}
	return event.Bind[EventEntityDestroySelf](auto.EventEntityDestroySelf(), subscriber, priority...)
}

func _EmitEventEntityDestroySelf(auto iAutoEventEntityDestroySelf, entity Entity) {
	if auto == nil {
		event.Panicf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs)
	}
	event.UnsafeEvent(auto.EventEntityDestroySelf()).Emit(func(subscriber event.Cache) bool {
		event.Cache2Iface[EventEntityDestroySelf](subscriber).OnEntityDestroySelf(entity)
		return true
	})
}

func _EmitEventEntityDestroySelfWithInterrupt(auto iAutoEventEntityDestroySelf, interrupt func(entity Entity) bool, entity Entity) {
	if auto == nil {
		event.Panicf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs)
	}
	event.UnsafeEvent(auto.EventEntityDestroySelf()).Emit(func(subscriber event.Cache) bool {
		if interrupt != nil {
			if interrupt(entity) {
				return false
			}
		}
		event.Cache2Iface[EventEntityDestroySelf](subscriber).OnEntityDestroySelf(entity)
		return true
	})
}

func HandleEventEntityDestroySelf(fun func(entity Entity)) EventEntityDestroySelfHandler {
	return EventEntityDestroySelfHandler(fun)
}

type EventEntityDestroySelfHandler func(entity Entity)

func (h EventEntityDestroySelfHandler) OnEntityDestroySelf(entity Entity) {
	h(entity)
}
