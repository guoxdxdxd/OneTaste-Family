package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"onetaste-family/backend/internal/models"
	"onetaste-family/backend/internal/services"
	"onetaste-family/backend/internal/utils"
)

// IngredientHandler 基础食材用户侧处理器
type IngredientHandler struct {
	ingredientService *services.IngredientService
}

// NewIngredientHandler 创建处理器
func NewIngredientHandler() *IngredientHandler {
	return &IngredientHandler{
		ingredientService: services.NewIngredientService(),
	}
}

// SearchIngredients 食材模糊搜索
// @Summary 食材模糊搜索
// @Description 根据关键字返回最多10条可用基础食材
// @Tags 食材
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param keyword query string true "关键字"
// @Param limit query int false "返回条数 (1-20, 默认10)"
// @Success 200 {object} utils.Response{data=map[string]interface{}} "搜索成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /ingredients/search [get]
func (h *IngredientHandler) SearchIngredients(c *gin.Context) {
	req, err := utils.BindQuery[models.IngredientSearchQuery](c)
	if err != nil {
		return
	}

	items, err := h.ingredientService.SearchIngredients(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerError("搜索食材失败"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(gin.H{"items": items}))
}

// GetIngredientsByCategory 分类分页查询
// @Summary 按分类分页查询食材
// @Description 用户根据分类浏览基础食材
// @Tags 食材
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category query string true "食材分类"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param keyword query string false "分类内模糊搜索"
// @Success 200 {object} utils.Response{data=models.IngredientCategoryListResponse} "查询成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /ingredients/by-category [get]
func (h *IngredientHandler) GetIngredientsByCategory(c *gin.Context) {
	req, err := utils.BindQuery[models.IngredientCategoryQuery](c)
	if err != nil {
		return
	}

	resp, err := h.ingredientService.GetIngredientsByCategory(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerError("获取食材列表失败"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(resp))
}
