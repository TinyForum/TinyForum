package common

import (
	"reflect"
	"strconv"
)

// ApplyDefaults 通过反射将结构体中为零值的字段设置为 default tag 指定的默认值。
// 支持基础类型（bool, int, uint, string, float），不支持复杂类型（切片、map 的默认值通常无意义）。
// 参数必须为指针类型，否则不会修改原值。
func ApplyDefaults[T any](req *T) {
	if req == nil {
		return
	}
	val := reflect.ValueOf(req).Elem()
	if val.Kind() != reflect.Struct {
		return
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		if !field.CanSet() {
			continue
		}
		if !isZeroValue(field) {
			continue
		}
		defaultVal := fieldType.Tag.Get("default")
		if defaultVal == "" {
			continue
		}
		setFieldFromString(field, defaultVal)
	}
}

// isZeroValue 判断字段是否是 Go 零值
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		// 复杂类型（如结构体）不自动填充默认值
		return false
	}
}

// setFieldFromString 将字符串 default 值设置到字段中
func setFieldFromString(field reflect.Value, defaultVal string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(defaultVal)
	case reflect.Bool:
		if b, err := strconv.ParseBool(defaultVal); err == nil {
			field.SetBool(b)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if i, err := strconv.ParseInt(defaultVal, 10, 64); err == nil {
			field.SetInt(i)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if u, err := strconv.ParseUint(defaultVal, 10, 64); err == nil {
			field.SetUint(u)
		}
	case reflect.Float32, reflect.Float64:
		if f, err := strconv.ParseFloat(defaultVal, 64); err == nil {
			field.SetFloat(f)
		}
	}
}
