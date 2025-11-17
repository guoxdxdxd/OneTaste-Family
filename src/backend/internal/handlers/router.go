package handlers

import (
	"github.com/gin-gonic/gin"
	"onetaste-family/backend/internal/middleware"
)

// Router 路由注册接口
type Router interface {
	RegisterRoutes(router *gin.RouterGroup)
}

// RegisterAuthRoutes 注册认证相关路由
func RegisterAuthRoutes(api *gin.RouterGroup) {
	authHandler := NewAuthHandler()
	
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
}

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(api *gin.RouterGroup) {
	userHandler := NewUserHandler()
	
	user := api.Group("/user")
	user.Use(middleware.AuthMiddleware()) // 需要认证
	{
		user.GET("/info", userHandler.GetUserInfo)
	}
}

