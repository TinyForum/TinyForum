// response/validation.go
package response

import (
	"github.com/go-playground/validator/v10"
)

// ParseValidationError 将 validator 的校验错误转换为 ValidationError 数组
func ParseValidationError(err error) []ValidationError {
	if err == nil {
		return nil
	}

	var errors []ValidationError

	// 处理 gin 的 validator 错误
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			errors = append(errors, ValidationError{
				Field:   e.Field(),
				Message: getValidationMessage(e),
			})
		}
		return errors
	}

	// 处理普通错误
	return []ValidationError{
		{Field: "request", Message: err.Error()},
	}
}

// getValidationMessage 根据校验标签返回友好的中文提示
func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "此字段不能为空"
	case "email":
		return "邮箱格式不正确"
	case "min":
		return "长度不能小于 " + e.Param()
	case "max":
		return "长度不能大于 " + e.Param()
	case "gte":
		return "必须大于等于 " + e.Param()
	case "lte":
		return "必须小于等于 " + e.Param()
	default:
		return "校验失败: " + e.Tag()
	}
}

// SimpleValidationError 快速创建单个校验错误
func SimpleValidationError(field, message string) []ValidationError {
	return []ValidationError{
		{Field: field, Message: message},
	}
}
