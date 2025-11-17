package utils

import (
	"regexp"
)

// ValidatePhone 验证手机号格式（中国大陆手机号）
func ValidatePhone(phone string) bool {
	// 匹配11位数字，以1开头
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	return matched
}

// ValidatePassword 验证密码强度
// 至少6位字符
func ValidatePassword(password string) bool {
	return len(password) >= 6
}

