package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BindJSON 自动绑定JSON请求体到结构体，并处理验证错误
// 使用泛型，自动推断请求类型
// 如果绑定失败，会自动返回400错误响应，并返回nil和error
// 使用示例：
//
//	req, err := utils.BindJSON[models.LoginRequest](c)
//	if err != nil {
//	    return // 错误已经在BindJSON中处理并返回响应
//	}
func BindJSON[T any](c *gin.Context) (*T, error) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := ExtractValidationErrors(err)
		c.JSON(http.StatusBadRequest, ErrorWithFields(400, "请求参数错误", errors))
		return nil, err
	}
	return &req, nil
}

// BindQuery 绑定查询参数到结构体
// 使用示例：
//
//	req, err := utils.BindQuery[models.ListRequest](c)
//	if err != nil {
//	    return
//	}
func BindQuery[T any](c *gin.Context) (*T, error) {
	var req T
	if err := c.ShouldBindQuery(&req); err != nil {
		errors := ExtractValidationErrors(err)
		c.JSON(http.StatusBadRequest, ErrorWithFields(400, "请求参数错误", errors))
		return nil, err
	}
	return &req, nil
}

// BindURI 绑定URI路径参数到结构体
// 使用示例：
//
//	req, err := utils.BindURI[models.IDRequest](c)
//	if err != nil {
//	    return
//	}
func BindURI[T any](c *gin.Context) (*T, error) {
	var req T
	if err := c.ShouldBindUri(&req); err != nil {
		errors := ExtractValidationErrors(err)
		c.JSON(http.StatusBadRequest, ErrorWithFields(400, "请求参数错误", errors))
		return nil, err
	}
	return &req, nil
}

// BindForm 绑定表单数据到结构体
// 使用示例：
//
//	req, err := utils.BindForm[models.FormRequest](c)
//	if err != nil {
//	    return
//	}
func BindForm[T any](c *gin.Context) (*T, error) {
	var req T
	if err := c.ShouldBind(&req); err != nil {
		errors := ExtractValidationErrors(err)
		c.JSON(http.StatusBadRequest, ErrorWithFields(400, "请求参数错误", errors))
		return nil, err
	}
	return &req, nil
}
