# 参数绑定工具使用指南

## 概述

`binding.go` 提供了基于 Go 泛型的参数绑定工具，可以自动绑定请求参数并处理验证错误，类似于 C# 的 `[FromBody]` 特性。

## 优势

1. **代码简洁**：一行代码完成参数绑定和错误处理
2. **类型安全**：利用泛型在编译期确保类型正确
3. **统一错误处理**：自动处理验证错误并返回统一格式的响应
4. **支持多种绑定方式**：JSON、Query、URI、Form

## 使用方法

### 1. 绑定 JSON 请求体（最常用）

```go
func (h *AuthHandler) Login(c *gin.Context) {
    // 一行代码完成参数绑定和错误处理
    req, err := utils.BindJSON[models.LoginRequest](c)
    if err != nil {
        return // 错误已经在BindJSON中处理并返回响应
    }

    // 使用 req 进行业务处理
    resp, err := h.userService.Login(req)
    // ...
}
```

**对比传统方式：**
```go
// 传统方式（需要5-6行代码）
var req models.LoginRequest
if err := c.ShouldBindJSON(&req); err != nil {
    errors := utils.ExtractValidationErrors(err)
    c.JSON(http.StatusBadRequest, utils.ErrorWithFields(400, "请求参数错误", errors))
    return
}

// 新方式（只需2-3行代码）
req, err := utils.BindJSON[models.LoginRequest](c)
if err != nil {
    return
}
```

### 2. 绑定查询参数

```go
// 请求模型
type ListRequest struct {
    Page     int    `form:"page" binding:"required,min=1"`
    PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
    Category string `form:"category"`
}

func (h *DishHandler) GetDishList(c *gin.Context) {
    req, err := utils.BindQuery[ListRequest](c)
    if err != nil {
        return
    }

    // 使用 req.Page, req.PageSize, req.Category
    // ...
}
```

### 3. 绑定 URI 路径参数

```go
// 请求模型
type IDRequest struct {
    ID int64 `uri:"id" binding:"required,min=1"`
}

func (h *DishHandler) GetDish(c *gin.Context) {
    req, err := utils.BindURI[IDRequest](c)
    if err != nil {
        return
    }

    // 使用 req.ID
    // ...
}
```

### 4. 绑定表单数据

```go
// 请求模型
type FormRequest struct {
    Name  string `form:"name" binding:"required"`
    Email string `form:"email" binding:"required,email"`
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
    req, err := utils.BindForm[FormRequest](c)
    if err != nil {
        return
    }

    // 使用 req.Name, req.Email
    // ...
}
```

## 完整示例

### 注册接口示例

```go
// auth_handler.go
func (h *AuthHandler) Register(c *gin.Context) {
    // 1. 绑定JSON参数
    req, err := utils.BindJSON[models.RegisterRequest](c)
    if err != nil {
        return // 参数验证失败，已自动返回400错误
    }

    // 2. 调用业务逻辑
    resp, err := h.userService.Register(req)
    if err != nil {
        // 3. 处理业务错误
        switch err {
        case services.ErrInvalidVerifyCode:
            c.JSON(http.StatusBadRequest, utils.BadRequest("验证码错误"))
            return
        case services.ErrPhoneExists:
            c.JSON(http.StatusBadRequest, utils.BadRequest("手机号已注册"))
            return
        default:
            c.JSON(http.StatusInternalServerError, utils.InternalServerError("注册失败"))
            return
        }
    }

    // 4. 返回成功响应
    c.JSON(http.StatusOK, utils.SuccessWithMessage("注册成功", resp))
}
```

### 带查询参数的列表接口示例

```go
// dish_handler.go
type DishListRequest struct {
    Page     int    `form:"page" binding:"required,min=1"`
    PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
    Category string `form:"category"`
}

func (h *DishHandler) GetDishList(c *gin.Context) {
    req, err := utils.BindQuery[DishListRequest](c)
    if err != nil {
        return
    }

    dishes, total, err := h.dishService.GetList(req.Page, req.PageSize, req.Category)
    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.InternalServerError("获取列表失败"))
        return
    }

    c.JSON(http.StatusOK, utils.Success(gin.H{
        "dishes": dishes,
        "total":  total,
        "page":   req.Page,
        "page_size": req.PageSize,
    }))
}
```

### 带URI参数的详情接口示例

```go
// dish_handler.go
type DishIDRequest struct {
    ID int64 `uri:"id" binding:"required,min=1"`
}

func (h *DishHandler) GetDish(c *gin.Context) {
    req, err := utils.BindURI[DishIDRequest](c)
    if err != nil {
        return
    }

    dish, err := h.dishService.GetByID(req.ID)
    if err != nil {
        if err == services.ErrDishNotFound {
            c.JSON(http.StatusNotFound, utils.NotFound("菜式不存在"))
            return
        }
        c.JSON(http.StatusInternalServerError, utils.InternalServerError("获取菜式失败"))
        return
    }

    c.JSON(http.StatusOK, utils.Success(dish))
}
```

## 注意事项

1. **错误处理**：`BindJSON`、`BindQuery` 等函数在绑定失败时会自动返回400错误响应，你只需要检查 `err != nil` 并 `return` 即可。

2. **类型推断**：使用泛型时，Go 会自动推断类型，无需手动指定指针类型。

3. **验证规则**：参数验证通过 Gin 的 `binding` 标签实现，确保在模型定义中添加正确的验证规则。

4. **性能**：泛型在编译期展开，运行时性能与手动绑定相同，无额外开销。

## 迁移指南

### 从传统方式迁移

**迁移前：**
```go
func (h *Handler) SomeMethod(c *gin.Context) {
    var req SomeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        errors := utils.ExtractValidationErrors(err)
        c.JSON(http.StatusBadRequest, utils.ErrorWithFields(400, "请求参数错误", errors))
        return
    }
    // ...
}
```

**迁移后：**
```go
func (h *Handler) SomeMethod(c *gin.Context) {
    req, err := utils.BindJSON[SomeRequest](c)
    if err != nil {
        return
    }
    // ...
}
```

## 支持的绑定方式

| 函数 | 绑定方式 | 适用场景 |
|------|---------|---------|
| `BindJSON` | JSON 请求体 | POST/PUT 请求，Content-Type: application/json |
| `BindQuery` | URL 查询参数 | GET 请求，如 `?page=1&size=10` |
| `BindURI` | URI 路径参数 | 路径中的参数，如 `/api/v1/dishes/:id` |
| `BindForm` | 表单数据 | POST 请求，Content-Type: application/x-www-form-urlencoded |

## 最佳实践

1. **为每个接口创建专门的请求模型**，不要复用业务模型
2. **在模型字段上添加合适的验证标签**，如 `binding:"required,min=1"`
3. **统一使用泛型绑定函数**，保持代码风格一致
4. **错误处理时直接 return**，不要重复处理已处理的错误

