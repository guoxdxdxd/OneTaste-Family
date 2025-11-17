package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"onetaste-family/backend/internal/utils"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, utils.Unauthorized("未提供认证token"))
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, utils.Unauthorized("认证格式错误"))
			c.Abort()
			return
		}

		// 解析token
		token := parts[1]
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.Unauthorized("无效的token"))
			c.Abort()
			return
		}

		// 将用户ID存储到上下文
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

