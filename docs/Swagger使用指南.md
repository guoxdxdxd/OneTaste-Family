# Swagger API 文档使用指南

## 概述

项目已集成 Swagger UI，可以自动生成和展示 API 接口文档。

## 快速开始

### 1. 访问文档

启动服务器后，在浏览器中访问：

```
http://localhost:8080/swagger/index.html
```

### 2. 生成文档

当添加新接口或修改现有接口后，需要重新生成文档：

```bash
cd src/backend
./scripts/swagger.sh
```

或者手动生成：

```bash
swag init -g cmd/main.go -o docs/swagger
```

## 功能特性

✅ **自动生成**：根据代码注释自动生成文档  
✅ **在线测试**：可以直接在 Swagger UI 中测试接口  
✅ **类型安全**：自动识别请求和响应类型  
✅ **认证支持**：支持 Bearer Token 认证测试  
✅ **示例数据**：自动显示请求和响应示例  

## 使用示例

### 在 Swagger UI 中测试接口

1. **打开 Swagger UI**：访问 `http://localhost:8080/swagger/index.html`

2. **测试注册接口**：
   - 找到 "用户认证" 分组
   - 点击 "POST /auth/register"
   - 点击 "Try it out"
   - 填写请求参数：
     ```json
     {
       "phone": "13800138000",
       "password": "password123",
       "verify_code": "123456",
       "nickname": "测试用户"
     }
     ```
   - 点击 "Execute"
   - 查看响应结果

3. **测试需要认证的接口**：
   - 先调用登录接口获取 Token
   - 点击页面右上角的 "Authorize" 按钮
   - 输入 Token（格式：`Bearer {token}` 或直接输入 token）
   - 点击 "Authorize"
   - 现在可以测试需要认证的接口了

## 代码注释规范

### 接口注释

```go
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
    // 实现代码
}
```

### 模型注释

```go
// RegisterRequest 注册请求
// @Description 用户注册请求参数
type RegisterRequest struct {
    Phone      string `json:"phone" example:"13800138000"`      // 手机号，11位数字
    Password   string `json:"password" example:"password123"`   // 密码，至少6位
    VerifyCode string `json:"verify_code" example:"123456"`      // 验证码，6位数字
    Nickname   string `json:"nickname" example:"张三"`          // 昵称
}
```

## 常用注释标签说明

| 标签 | 说明 | 示例 |
|------|------|------|
| `@Summary` | 接口简要说明 | `@Summary 用户注册` |
| `@Description` | 接口详细描述 | `@Description 用户注册接口的详细说明` |
| `@Tags` | 接口分组标签 | `@Tags 用户认证` |
| `@Accept` | 接受的请求类型 | `@Accept json` |
| `@Produce` | 返回的数据类型 | `@Produce json` |
| `@Param` | 请求参数说明 | `@Param request body models.RegisterRequest true "注册请求"` |
| `@Success` | 成功响应说明 | `@Success 200 {object} utils.Response{data=models.RegisterResponse}` |
| `@Failure` | 失败响应说明 | `@Failure 400 {object} utils.Response` |
| `@Router` | 路由路径和方法 | `@Router /auth/register [post]` |
| `@Security` | 安全认证 | `@Security BearerAuth` |

## 注意事项

### 1. 路由路径

`@Router` 中的路径**不需要**包含 `/api/v1` 前缀，因为已在 `main.go` 中通过 `@BasePath` 定义：

```go
// @BasePath  /api/v1
```

### 2. 响应格式

使用统一的响应格式时，使用以下语法：

```go
// @Success 200 {object} utils.Response{data=models.RegisterResponse}
```

### 3. 认证接口

需要认证的接口必须添加：

```go
// @Security BearerAuth
```

### 4. 示例值

在结构体字段中使用 `example` 标签：

```go
Phone string `json:"phone" example:"13800138000"`
```

### 5. 更新文档

每次修改接口后，记得重新生成文档：

```bash
./scripts/swagger.sh
```

## 文件结构

```
src/backend/
├── cmd/
│   └── main.go              # 包含 Swagger 通用配置
├── internal/
│   ├── handlers/            # 包含接口注释
│   ├── models/             # 包含模型注释
│   └── utils/              # 包含工具类注释
├── docs/
│   └── swagger/            # 生成的 Swagger 文档
│       ├── docs.go         # 生成的文档代码
│       ├── swagger.json    # JSON 格式文档
│       └── swagger.yaml    # YAML 格式文档
└── scripts/
    └── swagger.sh          # 文档生成脚本
```

## 故障排查

### 问题1：文档不显示

**原因**：文档未生成或生成失败

**解决**：
1. 检查 `docs/swagger/` 目录是否存在
2. 运行 `./scripts/swagger.sh` 重新生成
3. 检查代码中的 Swagger 注释格式是否正确

### 问题2：接口不显示

**原因**：注释格式错误或缺少必要标签

**解决**：
1. 检查是否包含 `@Router` 标签
2. 检查路由路径是否正确（不包含 `/api/v1` 前缀）
3. 检查方法是否正确（`[get]`, `[post]` 等）

### 问题3：认证测试失败

**原因**：Token 格式错误或已过期

**解决**：
1. 确保 Token 格式正确：`Bearer {token}` 或直接输入 token
2. 检查 Token 是否过期
3. 重新登录获取新 Token

## 更多资源

- [Swag 官方文档](https://github.com/swaggo/swag)
- [Swagger 规范](https://swagger.io/specification/)
- [Gin Swagger 示例](https://github.com/swaggo/gin-swagger)

