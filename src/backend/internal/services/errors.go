package services

import "errors"

var (
	// ErrInvalidVerifyCode 验证码错误
	ErrInvalidVerifyCode = errors.New("invalid verify code")
	// ErrCaptchaExpired 验证码过期
	ErrCaptchaExpired = errors.New("captcha expired")
	// ErrPhoneExists 手机号已存在
	ErrPhoneExists = errors.New("phone already exists")
	// ErrInvalidCredentials 用户名或密码错误
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user not found")
)
