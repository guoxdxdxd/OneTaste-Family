package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"onetaste-family/backend/internal/models"
	"onetaste-family/backend/internal/services"
	"onetaste-family/backend/internal/utils"
)

// DishHandler 菜式处理器
type DishHandler struct {
	dishService *services.DishService
}

// NewDishHandler 创建菜式处理器
func NewDishHandler() *DishHandler {
	return &DishHandler{
		dishService: services.NewDishService(),
	}
}

// CreateDish 创建菜式
// @Summary 创建菜式
// @Description 向当前家庭食谱库添加菜式，至少包含一个食材和一个烹饪步骤。需要Bearer Token认证。
// @Tags 菜式
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateDishRequest true "创建菜式请求"
// @Success 200 {object} utils.Response{data=models.DishCreateResponse} "创建成功"
// @Failure 400 {object} utils.Response "参数错误或业务限制"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 404 {object} utils.Response "尚未加入家庭"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /dishes [post]
func (h *DishHandler) CreateDish(c *gin.Context) {
	req, err := utils.BindJSON[models.CreateDishRequest](c)
	if err != nil {
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	resp, err := h.dishService.CreateDish(userID, req)
	if err != nil {
		switch err {
		case services.ErrFamilyNotFound:
			c.JSON(http.StatusNotFound, utils.NotFound("请先创建或加入家庭"))
		case services.ErrInvalidDishName:
			c.JSON(http.StatusBadRequest, utils.BadRequest("菜式名称不合法"))
		case services.ErrDishNameExists:
			c.JSON(http.StatusBadRequest, utils.BadRequest("菜式名称已存在"))
		case services.ErrDishLimitReached:
			c.JSON(http.StatusBadRequest, utils.BadRequest("菜式数量已达上限"))
		case services.ErrInvalidDishIngredients:
			c.JSON(http.StatusBadRequest, utils.BadRequest("请至少填写一个食材"))
		case services.ErrInvalidDishSteps:
			c.JSON(http.StatusBadRequest, utils.BadRequest("请至少填写一个烹饪步骤"))
		default:
			c.JSON(http.StatusInternalServerError, utils.InternalServerError("创建菜式失败"))
		}
		return
	}

	c.JSON(http.StatusOK, utils.SuccessWithMessage("创建成功", resp))
}

// GetDishList 获取菜式列表
// @Summary 获取菜式列表
// @Description 按照家庭返回菜式列表，支持按分类和名称关键字筛选。需要Bearer Token认证。
// @Tags 菜式
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码（默认1）"
// @Param page_size query int false "每页数量（默认20，最大100）"
// @Param category query string false "菜式分类"
// @Param keyword query string false "名称关键字"
// @Success 200 {object} utils.Response{data=models.DishListResponse} "获取成功"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 404 {object} utils.Response "尚未加入家庭"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /dishes [get]
func (h *DishHandler) GetDishList(c *gin.Context) {
	req, err := utils.BindQuery[models.DishListRequest](c)
	if err != nil {
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	resp, err := h.dishService.GetDishList(userID, req)
	if err != nil {
		switch err {
		case services.ErrFamilyNotFound:
			c.JSON(http.StatusNotFound, utils.NotFound("请先创建或加入家庭"))
		default:
			c.JSON(http.StatusInternalServerError, utils.InternalServerError("获取菜式列表失败"))
		}
		return
	}

	c.JSON(http.StatusOK, utils.Success(resp))
}

// GetDishDetail 获取菜式详情
// @Summary 获取菜式详情
// @Description 返回菜式的食材和烹饪步骤信息。需要Bearer Token认证。
// @Tags 菜式
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "菜式ID"
// @Success 200 {object} utils.Response{data=models.DishDetailResponse} "获取成功"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 404 {object} utils.Response "菜式或家庭不存在"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /dishes/{id} [get]
func (h *DishHandler) GetDishDetail(c *gin.Context) {
	uri, err := utils.BindURI[models.DishIDRequest](c)
	if err != nil {
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	resp, err := h.dishService.GetDishDetail(userID, uri.ID)
	if err != nil {
		switch err {
		case services.ErrFamilyNotFound:
			c.JSON(http.StatusNotFound, utils.NotFound("请先创建或加入家庭"))
		case services.ErrDishNotFound:
			c.JSON(http.StatusNotFound, utils.NotFound("菜式不存在或已删除"))
		default:
			c.JSON(http.StatusInternalServerError, utils.InternalServerError("获取菜式详情失败"))
		}
		return
	}

	c.JSON(http.StatusOK, utils.Success(resp))
}

// UpdateDish 更新菜式
// @Summary 更新菜式
// @Description 仅允许菜式创建者或家庭管理员编辑菜式。需要Bearer Token认证。
// @Tags 菜式
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "菜式ID"
// @Param request body models.UpdateDishRequest true "更新菜式请求"
// @Success 200 {object} utils.Response{data=models.DishDetailResponse} "更新成功"
// @Failure 400 {object} utils.Response "参数错误或业务限制"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 403 {object} utils.Response "无权限"
// @Failure 404 {object} utils.Response "菜式或家庭不存在"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /dishes/{id} [put]
func (h *DishHandler) UpdateDish(c *gin.Context) {
	uri, err := utils.BindURI[models.DishIDRequest](c)
	if err != nil {
		return
	}

	req, err := utils.BindJSON[models.UpdateDishRequest](c)
	if err != nil {
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	resp, err := h.dishService.UpdateDish(userID, uri.ID, req)
	if err != nil {
		switch err {
		case services.ErrFamilyNotFound:
			c.JSON(http.StatusNotFound, utils.NotFound("请先创建或加入家庭"))
		case services.ErrDishNotFound:
			c.JSON(http.StatusNotFound, utils.NotFound("菜式不存在或已删除"))
		case services.ErrDishPermissionDenied:
			c.JSON(http.StatusForbidden, utils.Forbidden("仅创建者或家庭管理员可编辑"))
		case services.ErrInvalidDishName:
			c.JSON(http.StatusBadRequest, utils.BadRequest("菜式名称不合法"))
		case services.ErrDishNameExists:
			c.JSON(http.StatusBadRequest, utils.BadRequest("菜式名称已存在"))
		case services.ErrInvalidDishIngredients:
			c.JSON(http.StatusBadRequest, utils.BadRequest("请至少填写一个食材"))
		case services.ErrInvalidDishSteps:
			c.JSON(http.StatusBadRequest, utils.BadRequest("请至少填写一个烹饪步骤"))
		default:
			c.JSON(http.StatusInternalServerError, utils.InternalServerError("更新菜式失败"))
		}
		return
	}

	c.JSON(http.StatusOK, utils.SuccessWithMessage("更新成功", resp))
}

// DeleteDish 删除菜式
// @Summary 删除菜式
// @Description 仅允许菜式创建者或家庭管理员删除菜式。需要Bearer Token认证。
// @Tags 菜式
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "菜式ID"
// @Success 200 {object} utils.Response "删除成功"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 403 {object} utils.Response "无权限"
// @Failure 404 {object} utils.Response "菜式或家庭不存在"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /dishes/{id} [delete]
func (h *DishHandler) DeleteDish(c *gin.Context) {
	uri, err := utils.BindURI[models.DishIDRequest](c)
	if err != nil {
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	if err := h.dishService.DeleteDish(userID, uri.ID); err != nil {
		switch err {
		case services.ErrFamilyNotFound:
			c.JSON(http.StatusNotFound, utils.NotFound("请先创建或加入家庭"))
		case services.ErrDishNotFound:
			c.JSON(http.StatusNotFound, utils.NotFound("菜式不存在或已删除"))
		case services.ErrDishPermissionDenied:
			c.JSON(http.StatusForbidden, utils.Forbidden("仅创建者或家庭管理员可删除"))
		default:
			c.JSON(http.StatusInternalServerError, utils.InternalServerError("删除菜式失败"))
		}
		return
	}

	c.JSON(http.StatusOK, utils.SuccessWithMessage("删除成功", nil))
}
