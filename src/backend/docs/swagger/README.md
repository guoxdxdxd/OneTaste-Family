# Swagger API 文档

## 访问文档

启动服务器后，访问以下地址查看 Swagger 文档：

```
http://localhost:8080/swagger/index.html
```

## 生成文档

### 方式1：使用脚本（推荐）

```bash
./scripts/swagger.sh
```

### 方式2：手动生成

```bash
swag init -g cmd/main.go -o docs/swagger
```

## 更新文档

当添加新的接口或修改现有接口时，需要重新生成文档：

1. 确保代码中的 Swagger 注释已更新
2. 运行生成脚本：`./scripts/swagger.sh`
3. 重启服务器（如果需要）

## Swagger 注释说明

### 接口注释示例

```go
// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口的详细描述
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "注册请求参数"
// @Success 200 {object} utils.Response{data=models.RegisterResponse} "注册成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
    // ...
}
```

### 模型注释示例

```go
// RegisterRequest 注册请求
// @Description 用户注册请求参数
type RegisterRequest struct {
    Phone      string `json:"phone" example:"13800138000"`      // 手机号
    Password   string `json:"password" example:"password123"`    // 密码
    VerifyCode string `json:"verify_code" example:"123456"`     // 验证码
    Nickname   string `json:"nickname" example:"张三"`           // 昵称
}
```

## 常用注释标签

- `@Summary` - 接口简要说明
- `@Description` - 接口详细描述
- `@Tags` - 接口分组标签
- `@Accept` - 接受的请求类型（json, form-data等）
- `@Produce` - 返回的数据类型（json等）
- `@Param` - 请求参数说明
- `@Success` - 成功响应说明
- `@Failure` - 失败响应说明
- `@Router` - 路由路径和方法
- `@Security` - 安全认证（如BearerAuth）

## 注意事项

1. **路由路径**：`@Router` 中的路径不需要包含 `/api/v1` 前缀，因为已在 `@BasePath` 中定义
2. **认证接口**：需要认证的接口使用 `@Security BearerAuth`
3. **响应格式**：使用 `utils.Response{data=具体类型}` 格式定义响应
4. **示例值**：在结构体字段中使用 `example` 标签提供示例值

## 更多信息

- [Swag 官方文档](https://github.com/swaggo/swag)
- [Swagger 规范](https://swagger.io/specification/)

