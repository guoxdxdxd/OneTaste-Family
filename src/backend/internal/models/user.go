package models

import "time"

// User 用户模型
type User struct {
	ID        int64     `json:"user_id" db:"id"`
	Phone     string    `json:"phone" db:"phone"`
	Password  string    `json:"-" db:"password"` // 不序列化到JSON
	Nickname  string    `json:"nickname" db:"nickname"`
	Avatar    string    `json:"avatar" db:"avatar"`
	Status    int       `json:"status" db:"status"` // 1-正常，0-禁用
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// RegisterRequest 注册请求
// @Description 用户注册请求参数
type RegisterRequest struct {
	Phone      string `json:"phone" binding:"required" example:"13800138000"`      // 手机号，11位数字
	Password   string `json:"password" binding:"required,min=6" example:"password123"` // 密码，至少6位
	VerifyCode string `json:"verify_code" binding:"required" example:"123456"`         // 验证码，6位数字
	Nickname   string `json:"nickname" binding:"required" example:"张三"`              // 昵称
}

// LoginRequest 登录请求
// @Description 用户登录请求参数
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required" example:"13800138000"`      // 手机号
	Password string `json:"password" binding:"required" example:"password123"`    // 密码
}

// RegisterResponse 注册响应
// @Description 用户注册成功返回的数据
type RegisterResponse struct {
	UserID int64  `json:"user_id" example:"1"`                                    // 用户ID
	Token  string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT Token
}

// LoginResponse 登录响应
// @Description 用户登录成功返回的数据
type LoginResponse struct {
	UserID    int64  `json:"user_id" example:"1"`                                    // 用户ID
	Token     string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT Token
	ExpiresIn int64  `json:"expires_in" example:"86400"`                              // Token过期时间（秒）
}

// UserInfoResponse 用户信息响应
// @Description 用户信息返回数据
type UserInfoResponse struct {
	UserID    int64          `json:"user_id" example:"1"`              // 用户ID
	Phone     string         `json:"phone" example:"13800138000"`      // 手机号
	Nickname  string         `json:"nickname" example:"张三"`            // 昵称
	Avatar    string         `json:"avatar" example:"https://..."`     // 头像URL
	Membership MembershipInfo `json:"membership"`                      // 会员信息
}

// MembershipInfo 会员信息
// @Description 用户会员信息
type MembershipInfo struct {
	Type      string     `json:"type" example:"free"`       // 会员类型：free-免费版，premium-付费版
	ExpiresAt *time.Time `json:"expires_at" example:"2024-12-31T23:59:59Z"` // 到期时间，免费版为null
}

