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
		auth.GET("/captcha", authHandler.GetCaptcha)
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

// RegisterFamilyRoutes 注册家庭相关路由
func RegisterFamilyRoutes(api *gin.RouterGroup) {
	familyHandler := NewFamilyHandler()

	family := api.Group("/family")
	family.Use(middleware.AuthMiddleware()) // 需要认证
	{
		family.POST("/create", familyHandler.CreateFamily)
		family.GET("/info", familyHandler.GetFamilyInfo)
		family.POST("/member/invite", familyHandler.JoinFamilyViaInvite)
		family.GET("/members", familyHandler.GetFamilyMembers)
	}
}

// RegisterDishRoutes 注册菜式相关路由
func RegisterDishRoutes(api *gin.RouterGroup) {
	dishHandler := NewDishHandler()

	dishes := api.Group("/dishes")
	dishes.Use(middleware.AuthMiddleware())
	{
		dishes.POST("", dishHandler.CreateDish)
		dishes.GET("", dishHandler.GetDishList)
		dishes.GET("/:id", dishHandler.GetDishDetail)
		dishes.PUT("/:id", dishHandler.UpdateDish)
		dishes.DELETE("/:id", dishHandler.DeleteDish)
	}
}

// RegisterMediaRoutes 注册媒体上传路由
func RegisterMediaRoutes(api *gin.RouterGroup) {
	mediaHandler := NewMediaHandler()

	media := api.Group("/media")
	media.Use(middleware.AuthMiddleware())
	{
		media.POST("/upload", mediaHandler.UploadMedia)
	}
}
