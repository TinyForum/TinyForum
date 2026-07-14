package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// FlattenConfig 将嵌套的配置结构展平为键值对
func FlattenConfig(prefix string, data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	flattenConfig(prefix, data, result)
	return result
}

func flattenConfig(prefix string, value interface{}, result map[string]interface{}) {
	if value == nil {
		return
	}

	v := reflect.ValueOf(value)

	// 处理指针
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		for _, key := range v.MapKeys() {
			newKey := key.String()
			if prefix != "" {
				newKey = prefix + "." + newKey
			}
			flattenConfig(newKey, v.MapIndex(key).Interface(), result)
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			newKey := fmt.Sprintf("%s.%d", prefix, i)
			flattenConfig(newKey, v.Index(i).Interface(), result)
		}

	default:
		// 基本类型直接存储
		result[prefix] = value
	}
}

// SetNestedValue 设置嵌套值
func SetNestedValue(data map[string]interface{}, key string, value string) error {
	parts := strings.Split(key, ".")

	// 如果只有一层，直接设置
	if len(parts) == 1 {
		data[parts[0]] = convertValue(value)
		return nil
	}

	// 递归创建嵌套结构
	current := data
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		if _, exists := current[part]; !exists {
			current[part] = make(map[string]interface{})
		}

		next, ok := current[part].(map[string]interface{})
		if !ok {
			// 如果存在但不是 map，覆盖
			current[part] = make(map[string]interface{})
			next = current[part].(map[string]interface{})
		}
		current = next
	}

	// 设置最终值
	lastPart := parts[len(parts)-1]
	current[lastPart] = convertValue(value)

	return nil
}

// convertValue 转换字符串值为对应类型
func convertValue(value string) interface{} {
	// 尝试转换为布尔值
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}

	// 尝试转换为数字
	if val, err := strconv.Atoi(value); err == nil {
		return val
	}
	if val, err := strconv.ParseFloat(value, 64); err == nil {
		return val
	}

	// 默认为字符串
	return value
}

// FlattenConfigToString 将嵌套配置展平为字符串键值对
func FlattenConfigToString(data interface{}) map[string]string {
	flat := FlattenConfig("", data)
	result := make(map[string]string)

	for key, val := range flat {
		result[key] = formatValue(val)
	}

	return result
}

func formatValue(val interface{}) string {
	if val == nil {
		return ""
	}

	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", val)
	}
}
