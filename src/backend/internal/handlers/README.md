# Handlers 路由组织说明

## 路由组织方式

采用**模块化路由注册**的方式，每个业务模块负责注册自己的路由。

### 目录结构

```
internal/handlers/
├── routes.go          # 统一路由注册入口
├── router.go          # 各模块路由注册函数
├── auth_handler.go     # 认证处理器
├── user_handler.go     # 用户处理器
├── family_handler.go   # 家庭处理器（待实现）
├── dish_handler.go     # 菜式处理器（待实现）
└── ...
```

## 如何添加新模块的路由

### 步骤1：在 router.go 中添加路由注册函数

```go
// RegisterFamilyRoutes 注册家庭相关路由
func RegisterFamilyRoutes(api *gin.RouterGroup) {
	familyHandler := NewFamilyHandler()
	
	family := api.Group("/family")
	family.Use(middleware.AuthMiddleware()) // 需要认证
	{
		family.POST("/create", familyHandler.Create)
		family.GET("/info", familyHandler.GetInfo)
		family.POST("/member/invite", familyHandler.InviteMember)
		family.GET("/members", familyHandler.GetMembers)
	}
}
```

### 步骤2：在 routes.go 中注册新模块

```go
func RegisterAllRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		RegisterAuthRoutes(api)
		RegisterUserRoutes(api)
		RegisterFamilyRoutes(api)  // 添加这一行
	}
}
```

## 优势

1. **模块化**：每个模块的路由逻辑集中管理
2. **可维护性**：main.go 保持简洁，只负责应用初始化
3. **可扩展性**：添加新模块只需两步
4. **清晰性**：路由结构一目了然

## 对比

### 之前的方式（所有路由在 main.go）
```go
// main.go 会变得很长
api := r.Group("/api/v1")
{
    auth := api.Group("/auth")
    { /* ... */ }
    user := api.Group("/user")
    { /* ... */ }
    family := api.Group("/family")
    { /* ... */ }
    // ... 更多路由
}
```

### 现在的方式（模块化）
```go
// main.go 保持简洁
handlers.RegisterAllRoutes(r)

// routes.go 统一管理
func RegisterAllRoutes(r *gin.Engine) {
    api := r.Group("/api/v1")
    {
        RegisterAuthRoutes(api)
        RegisterUserRoutes(api)
        RegisterFamilyRoutes(api)
    }
}
```

