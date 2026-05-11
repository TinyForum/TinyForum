package mask

import (
	"reflect"
	"sync"
)

// fieldInfo 存储字段的脱敏元信息
type fieldInfo struct {
	index  []int             // 嵌套字段索引路径
	name   string            // 策略名
	params map[string]string // 策略参数
}

// typeInfo 存储一个类型的脱敏字段列表
type typeInfo struct {
	fields []fieldInfo
}

var (
	typeCache sync.Map // map[reflect.Type]*typeInfo
)

// getTypeInfo 获取或构造类型的脱敏信息
func getTypeInfo(t reflect.Type) *typeInfo {
	if cached, ok := typeCache.Load(t); ok {
		return cached.(*typeInfo)
	}
	info := &typeInfo{}
	collectFields(t, nil, info)
	typeCache.Store(t, info)
	return info
}

// collectFields 递归收集带有 mask 标签的字段
func collectFields(t reflect.Type, parentIndex []int, info *typeInfo) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("mask")
		if tag == "-" {
			continue
		}
		index := make([]int, len(parentIndex)+1)
		copy(index, parentIndex)
		index[len(parentIndex)] = i

		if tag != "" && tag != "-" {
			// 当前字段有 mask 标签
			name, params := parseTag(tag)
			if name != "" {
				info.fields = append(info.fields, fieldInfo{
					index:  index,
					name:   name,
					params: params,
				})
			}
		}

		// 如果是结构体或结构体指针，继续递归
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() == reflect.Struct {
			collectFields(fieldType, index, info)
		}
	}
}
