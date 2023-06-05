package util

import (
	"reflect"
	"strings"
)

// Zero 创建零值
func Zero[T any]() T {
	var zero T
	return zero
}

// New new并初始化
func New[T any](v T) *T {
	var rv T
	rv = v
	return &rv
}

// Bool2Int bool转int
func Bool2Int(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Int2Bool int转bool
func Int2Bool(v int) bool {
	if v != 0 {
		return true
	}
	return false
}

// TypeName 类型名
func TypeName[T any]() string {
	return reflect.TypeOf((*T)(nil)).Elem().Name()
}

// TypeFullName 类型全名
func TypeFullName[T any]() string {
	t := reflect.TypeOf((*T)(nil)).Elem()
	return t.PkgPath() + "/" + t.Name()
}

// TypeOfAnyFullName 类型全名
func TypeOfAnyFullName(i any) string {
	t, ok := i.(reflect.Type)
	if !ok {
		t = reflect.TypeOf(i)
	}
	return t.PkgPath() + "/" + t.Name()
}

// TypeOfFullName 类型全名
func TypeOfFullName(t reflect.Type) string {
	return t.PkgPath() + "/" + t.Name()
}

// WriteFullName 写入类型全名
func WriteFullName(sb strings.Builder, t reflect.Type) {
	sb.WriteString(t.PkgPath())
	sb.WriteString("/")
	sb.WriteString(t.Name())
}
