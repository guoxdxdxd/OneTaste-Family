package utils

// Response 统一响应格式
// @Description 统一的API响应格式
type Response struct {
	Code    int         `json:"code" example:"200"`                    // 响应状态码
	Message string      `json:"message" example:"success"`              // 响应消息
	Data    interface{} `json:"data,omitempty"`                        // 响应数据
	Errors  []FieldError `json:"errors,omitempty"`                     // 字段错误列表（仅在错误时返回）
}

// FieldError 字段错误
// @Description 字段验证错误信息
type FieldError struct {
	Field   string `json:"field" example:"phone"`                      // 错误字段名
	Message string `json:"message" example:"手机号格式不正确"`            // 错误消息
}

// Success 成功响应
func Success(data interface{}) *Response {
	return &Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(message string, data interface{}) *Response {
	return &Response{
		Code:    200,
		Message: message,
		Data:    data,
	}
}

// Error 错误响应
func Error(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
	}
}

// ErrorWithFields 带字段错误的响应
func ErrorWithFields(code int, message string, errors []FieldError) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Errors:  errors,
	}
}

// BadRequest 400错误
func BadRequest(message string) *Response {
	return Error(400, message)
}

// Unauthorized 401错误
func Unauthorized(message string) *Response {
	return Error(401, message)
}

// Forbidden 403错误
func Forbidden(message string) *Response {
	return Error(403, message)
}

// NotFound 404错误
func NotFound(message string) *Response {
	return Error(404, message)
}

// InternalServerError 500错误
func InternalServerError(message string) *Response {
	return Error(500, message)
}

// ExtractValidationErrors 提取验证错误
func ExtractValidationErrors(err error) []FieldError {
	var errors []FieldError
	
	if err == nil {
		return errors
	}

	// 尝试转换为gin的验证错误
	if validationErrors, ok := err.(interface {
		Error() string
	}); ok {
		// 如果是gin的binding错误，提取字段信息
		// 这里简化处理，返回通用错误信息
		// 如果需要更详细的字段错误，可以使用validator库
		errors = append(errors, FieldError{
			Field:   "request",
			Message: validationErrors.Error(),
		})
	} else {
		errors = append(errors, FieldError{
			Field:   "general",
			Message: err.Error(),
		})
	}
	
	return errors
}

