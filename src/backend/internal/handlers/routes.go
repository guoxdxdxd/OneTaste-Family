package handlers

import (
	"github.com/gin-gonic/gin"
)

// RegisterAllRoutes 注册所有路由
// 这是统一的路由注册入口，所有模块的路由都在这里注册
func RegisterAllRoutes(r *gin.Engine) {
	// API v1 路由组
	api := r.Group("/api/v1")
	{
		// 注册各个模块的路由
		RegisterAuthRoutes(api)   // 认证路由（不需要认证）
		RegisterUserRoutes(api)   // 用户路由（需要认证）
		RegisterFamilyRoutes(api) // 家庭路由
		RegisterDishRoutes(api)   // 菜式路由
		RegisterMediaRoutes(api)  // 文件上传路由
		// 后续添加新模块时，只需要在这里添加一行即可
		// RegisterDishRoutes(api)     // 菜式路由
		// RegisterMenuRoutes(api)     // 菜单路由
		// RegisterShoppingRoutes(api) // 购物清单路由
	}
}
