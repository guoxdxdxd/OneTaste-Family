package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"onetaste-family/backend/internal/models"
	"onetaste-family/backend/internal/services"
	"onetaste-family/backend/internal/utils"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	userService    *services.UserService
	captchaService *services.CaptchaService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		userService:    services.NewUserService(),
		captchaService: services.NewCaptchaService(),
	}
}

// GetCaptcha 获取图形验证码
// @Summary 获取图形验证码
// @Description 获取一个带干扰的图形验证码图片及对应的编码
// @Tags 用户认证
// @Param width query int false "图片宽度（120-360），默认220"
// @Param height query int false "图片高度（40-160），默认70"
// @Produce json
// @Success 200 {object} utils.Response{data=models.CaptchaResponse} "验证码获取成功"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /auth/captcha [get]
func (h *AuthHandler) GetCaptcha(c *gin.Context) {
	opts, err := utils.BindQuery[models.CaptchaRequest](c)
	if err != nil {
		return
	}

	resp, err := h.captchaService.GenerateCaptcha(c.Request.Context(), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerError("获取验证码失败"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(resp))
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口，支持手机号、密码、验证码、昵称注册。注册成功后返回用户ID和JWT Token。
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "注册请求参数"
// @Success 200 {object} utils.Response{data=models.RegisterResponse} "注册成功"
// @Failure 400 {object} utils.Response "请求参数错误、验证码错误或手机号已注册"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	req, err := utils.BindJSON[models.RegisterRequest](c)
	if err != nil {
		return // 错误已经在BindJSON中处理并返回响应
	}

	resp, err := h.userService.Register(c.Request.Context(), req)
	if err != nil {
		// 根据错误类型返回不同的状态码
		switch err {
		case services.ErrCaptchaExpired:
			c.JSON(http.StatusBadRequest, utils.BadRequest("验证码已过期"))
			return
		case services.ErrInvalidVerifyCode:
			c.JSON(http.StatusBadRequest, utils.BadRequest("验证码不正确"))
			return
		case services.ErrPhoneExists:
			c.JSON(http.StatusBadRequest, utils.BadRequest("手机号已注册"))
			return
		default:
			c.JSON(http.StatusInternalServerError, utils.InternalServerError("注册失败"))
			return
		}
	}

	c.JSON(http.StatusOK, utils.SuccessWithMessage("注册成功", resp))
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口，支持手机号和密码登录。登录成功后返回用户ID、JWT Token和过期时间。
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "登录请求参数"
// @Success 200 {object} utils.Response{data=models.LoginResponse} "登录成功"
// @Failure 400 {object} utils.Response "手机号或密码错误"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	req, err := utils.BindJSON[models.LoginRequest](c)
	if err != nil {
		return // 错误已经在BindJSON中处理并返回响应
	}

	resp, err := h.userService.Login(req)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			c.JSON(http.StatusBadRequest, utils.BadRequest("手机号或密码错误"))
			return
		default:
			c.JSON(http.StatusInternalServerError, utils.InternalServerError("登录失败"))
			return
		}
	}

	c.JSON(http.StatusOK, utils.SuccessWithMessage("登录成功", resp))
}
