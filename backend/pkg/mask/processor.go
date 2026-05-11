package mask

import (
	"reflect"
)

// processStruct 处理结构体值的所有脱敏字段
func processStruct(v reflect.Value) error {
	if v.Kind() != reflect.Struct {
		return nil
	}
	t := v.Type()
	info := getTypeInfo(t)

	for _, fi := range info.fields {
		// 根据索引路径获取字段值
		field := v.FieldByIndex(fi.index)
		if !field.CanSet() {
			continue
		}
		// 支持 string 和 *string
		switch field.Kind() {
		case reflect.String:
			orig := field.String()
			masked := applyStrategy(fi.name, orig, fi.params)
			field.SetString(masked)
		case reflect.Ptr:
			if field.IsNil() {
				continue
			}
			elem := field.Elem()
			if elem.Kind() == reflect.String {
				orig := elem.String()
				masked := applyStrategy(fi.name, orig, fi.params)
				elem.SetString(masked)
			}
		}
	}
	return nil
}
