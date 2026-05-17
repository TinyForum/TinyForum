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
		field := v.FieldByIndex(fi.index)
		if !field.CanSet() {
			continue
		}
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

		case reflect.Slice, reflect.Array:
			for i := 0; i < field.Len(); i++ {
				elem := field.Index(i)
				if elem.Kind() == reflect.String {
					if elem.CanSet() {
						orig := elem.String()
						masked := applyStrategy(fi.name, orig, fi.params)
						elem.SetString(masked)
					}
				} else if elem.Kind() == reflect.Ptr && elem.Type().Elem().Kind() == reflect.String && !elem.IsNil() {
					orig := elem.Elem().String()
					masked := applyStrategy(fi.name, orig, fi.params)
					elem.Elem().SetString(masked)
				}
			}
		}
	}
	return nil
}
