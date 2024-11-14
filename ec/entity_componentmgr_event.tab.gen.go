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

// Code generated by eventc eventtab --name=entityComponentMgrEventTab; DO NOT EDIT.

package ec

import (
	event "git.golaxy.org/core/event"
)

type IEntityComponentMgrEventTab interface {
	EventComponentMgrAddComponents() event.IEvent
	EventComponentMgrRemoveComponent() event.IEvent
	EventComponentMgrFirstTouchComponent() event.IEvent
}

var (
	_entityComponentMgrEventTabId = event.DeclareEventTabIdT[entityComponentMgrEventTab]()
	EventComponentMgrAddComponentsId = _entityComponentMgrEventTabId + 0
	EventComponentMgrRemoveComponentId = _entityComponentMgrEventTabId + 1
	EventComponentMgrFirstTouchComponentId = _entityComponentMgrEventTabId + 2
)

type entityComponentMgrEventTab [3]event.Event

func (eventTab *entityComponentMgrEventTab) Init(autoRecover bool, reportError chan error, recursion event.EventRecursion) {
	(*eventTab)[0].Init(autoRecover, reportError, recursion)
	(*eventTab)[1].Init(autoRecover, reportError, recursion)
	(*eventTab)[2].Init(autoRecover, reportError, recursion)
}

func (eventTab *entityComponentMgrEventTab) Open() {
	for i := range *eventTab {
		(*eventTab)[i].Open()
	}
}

func (eventTab *entityComponentMgrEventTab) Close() {
	for i := range *eventTab {
		(*eventTab)[i].Close()
	}
}

func (eventTab *entityComponentMgrEventTab) Clean() {
	for i := range *eventTab {
		(*eventTab)[i].Clean()
	}
}

func (eventTab *entityComponentMgrEventTab) Ctrl() event.IEventCtrl {
	return eventTab
}

func (eventTab *entityComponentMgrEventTab) Event(id uint64) event.IEvent {
	if _entityComponentMgrEventTabId != id & 0xFFFFFFFF00000000 {
		return nil
	}
	pos := id & 0xFFFFFFFF
	if pos >= uint64(len(*eventTab)) {
		return nil
	}
	return &(*eventTab)[pos]
}

func (eventTab *entityComponentMgrEventTab) EventComponentMgrAddComponents() event.IEvent {
	return &(*eventTab)[0]
}

func (eventTab *entityComponentMgrEventTab) EventComponentMgrRemoveComponent() event.IEvent {
	return &(*eventTab)[1]
}

func (eventTab *entityComponentMgrEventTab) EventComponentMgrFirstTouchComponent() event.IEvent {
	return &(*eventTab)[2]
}