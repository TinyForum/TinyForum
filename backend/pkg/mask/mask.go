package mask

import (
	"reflect"
)

// Mask 原地修改结构体指针所指对象，将标记了 mask 标签的字段替换为脱敏后的值。
func Mask(v any) error {
	if v == nil {
		return ErrNilPointer
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrNonPointer
	}
	return processStruct(rv.Elem())
}

// MaskCopy 返回脱敏后的深拷贝，原对象不变。
func MaskCopy(v any) (any, error) {
	if v == nil {
		return nil, ErrNilPointer
	}
	rv := reflect.ValueOf(v)
	copied := deepCopy(rv)
	if !copied.CanAddr() {
		// 如果复制结果不可寻址，返回其接口值，但后续 process 可能需要指针
		// 简单处理：若原始是指针，复制后也保持指针
		if rv.Kind() == reflect.Pointer {
			ptr := reflect.New(copied.Type())
			ptr.Elem().Set(copied)
			copied = ptr
		}
	}
	// 确保传入的是可寻址的值或指针
	var target reflect.Value
	if copied.Kind() == reflect.Pointer {
		target = copied.Elem()
	} else {
		target = copied
	}
	if err := processStruct(target); err != nil {
		return nil, err
	}
	return copied.Interface(), nil
}

// deepCopy 递归深拷贝任意值
func deepCopy(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		cp := reflect.New(v.Elem().Type())
		cp.Elem().Set(deepCopy(v.Elem()))
		return cp
	case reflect.Struct:
		cp := reflect.New(v.Type()).Elem()
		for i := 0; i < v.NumField(); i++ {
			cp.Field(i).Set(deepCopy(v.Field(i)))
		}
		return cp
	case reflect.Slice:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		cp := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
		for i := 0; i < v.Len(); i++ {
			cp.Index(i).Set(deepCopy(v.Index(i)))
		}
		return cp
	case reflect.Array:
		cp := reflect.New(v.Type()).Elem()
		for i := 0; i < v.Len(); i++ {
			cp.Index(i).Set(deepCopy(v.Index(i)))
		}
		return cp
	case reflect.Map:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		cp := reflect.MakeMapWithSize(v.Type(), v.Len())
		iter := v.MapRange()
		for iter.Next() {
			cp.SetMapIndex(deepCopy(iter.Key()), deepCopy(iter.Value()))
		}
		return cp
	case reflect.Interface:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		return deepCopy(v.Elem())
	default:
		return v
	}
}
