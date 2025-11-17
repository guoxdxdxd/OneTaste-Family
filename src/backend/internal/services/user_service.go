package services

import (
	"errors"
	"fmt"

	"onetaste-family/backend/internal/models"
	"onetaste-family/backend/internal/repositories"
	"onetaste-family/backend/internal/utils"
)

var (
	// ErrInvalidVerifyCode 验证码错误
	ErrInvalidVerifyCode = errors.New("invalid verify code")
	// ErrPhoneExists 手机号已存在
	ErrPhoneExists = errors.New("phone already exists")
	// ErrInvalidCredentials 用户名或密码错误
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// UserService 用户业务逻辑层
type UserService struct {
	userRepo *repositories.UserRepository
}

// NewUserService 创建用户Service
func NewUserService() *UserService {
	return &UserService{
		userRepo: repositories.NewUserRepository(),
	}
}

// Register 用户注册
func (s *UserService) Register(req *models.RegisterRequest) (*models.RegisterResponse, error) {
	// 验证手机号格式
	if !utils.ValidatePhone(req.Phone) {
		return nil, fmt.Errorf("invalid phone format")
	}

	// 验证密码强度
	if !utils.ValidatePassword(req.Password) {
		return nil, fmt.Errorf("password must be at least 6 characters")
	}

	// 验证验证码（暂时使用简单验证，后续可接入短信服务）
	if !s.validateVerifyCode(req.Phone, req.VerifyCode) {
		return nil, ErrInvalidVerifyCode
	}

	// 检查手机号是否已存在
	exists, err := s.userRepo.ExistsByPhone(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("failed to check phone exists: %w", err)
	}
	if exists {
		return nil, ErrPhoneExists
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 创建用户
	user := &models.User{
		Phone:    req.Phone,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Status:   1, // 正常状态
	}

	if err := s.userRepo.Create(user); err != nil {
		if err == repositories.ErrUserExists {
			return nil, ErrPhoneExists
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 生成JWT Token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.RegisterResponse{
		UserID: user.ID,
		Token:  token,
	}, nil
}

// Login 用户登录
func (s *UserService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByPhone(req.Phone)
	if err != nil {
		if err == repositories.ErrUserNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	// 生成JWT Token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		UserID:    user.ID,
		Token:     token,
		ExpiresIn: utils.GetTokenExpiration(),
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(userID int64) (*models.UserInfoResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == repositories.ErrUserNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// TODO: 获取会员信息（后续实现）
	membership := models.MembershipInfo{
		Type:      "free",
		ExpiresAt: nil,
	}

	return &models.UserInfoResponse{
		UserID:    user.ID,
		Phone:     user.Phone,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Membership: membership,
	}, nil
}

// validateVerifyCode 验证验证码
// TODO: 后续接入短信服务，从Redis获取验证码
// 目前为了开发方便，使用固定验证码 "123456" 或从环境变量读取
func (s *UserService) validateVerifyCode(phone, code string) bool {
	// 开发环境：允许使用固定验证码
	// 生产环境：从Redis获取验证码并验证
	// 这里先实现一个简单版本，后续可以接入短信服务
	if code == "123456" {
		return true
	}
	// TODO: 从Redis获取验证码并验证
	return false
}

