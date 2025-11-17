package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"onetaste-family/backend/internal/services"
	"onetaste-family/backend/internal/utils"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(),
	}
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的基本信息和会员状态。需要Bearer Token认证。
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=models.UserInfoResponse} "获取成功"
// @Failure 401 {object} utils.Response "未授权或Token无效"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /user/info [get]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.Unauthorized("未授权"))
		return
	}

	// 类型断言
	id, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.Unauthorized("无效的用户ID"))
		return
	}

	info, err := h.userService.GetUserInfo(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerError("获取用户信息失败"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(info))
}

