package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"onetaste-family/backend/internal/models"
	"onetaste-family/backend/internal/services"
	"onetaste-family/backend/internal/utils"
)

// MenuHandler 菜单处理器
type MenuHandler struct {
	menuService *services.MenuService
}

// NewMenuHandler 创建菜单处理器
func NewMenuHandler() *MenuHandler {
	return &MenuHandler{
		menuService: services.NewMenuService(),
	}
}

// CreateMenu 创建菜单
// @Summary 创建菜单
// @Description 为某一天某一餐创建菜单，支持日期、餐次（早餐/午餐/晚餐）、菜式列表输入。同一日期同一餐次只能有一个菜单，可以覆盖已有菜单。需要Bearer Token认证。
// @Tags 菜单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateMenuRequest true "创建菜单请求"
// @Success 200 {object} utils.Response{data=models.MenuCreateResponse} "创建成功"
// @Failure 400 {object} utils.Response "参数错误或业务限制"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 404 {object} utils.Response "尚未加入家庭或菜式不存在"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /menus [post]
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	req, err := utils.BindJSON[models.CreateMenuRequest](c)
	if err != nil {
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	resp, err := h.menuService.CreateMenu(userID, req)
	if err != nil {
		switch err {
		case services.ErrFamilyNotFound:
			c.JSON(http.StatusNotFound, utils.NotFound("请先创建或加入家庭"))
		case services.ErrInvalidMenuDate:
			c.JSON(http.StatusBadRequest, utils.BadRequest("日期格式错误，请使用YYYY-MM-DD格式"))
		case services.ErrInvalidMealType:
			c.JSON(http.StatusBadRequest, utils.BadRequest("餐次类型错误，必须是breakfast、lunch或dinner"))
		case services.ErrInvalidDishIDs:
			c.JSON(http.StatusBadRequest, utils.BadRequest("请至少选择一个菜式"))
		case services.ErrDishNotFound:
			c.JSON(http.StatusNotFound, utils.NotFound("菜式不存在或已删除"))
		case services.ErrDishNotInFamily:
			c.JSON(http.StatusBadRequest, utils.BadRequest("菜式不属于当前家庭"))
		default:
			c.JSON(http.StatusInternalServerError, utils.InternalServerError("创建菜单失败"))
		}
		return
	}

	c.JSON(http.StatusOK, utils.SuccessWithMessage("创建成功", resp))
}

// GetDailyMenu 获取每日菜单
// @Summary 获取每日菜单
// @Description 获取某一天的三餐菜单（早餐、午餐、晚餐）。需要Bearer Token认证。
// @Tags 菜单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param date query string true "日期，格式：YYYY-MM-DD"
// @Success 200 {object} utils.Response{data=models.DailyMenuResponse} "获取成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 404 {object} utils.Response "尚未加入家庭"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /menus/daily [get]
func (h *MenuHandler) GetDailyMenu(c *gin.Context) {
	req, err := utils.BindQuery[models.DailyMenuRequest](c)
	if err != nil {
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	resp, err := h.menuService.GetDailyMenu(userID, req)
	if err != nil {
		switch err {
		case services.ErrFamilyNotFound:
			c.JSON(http.StatusNotFound, utils.NotFound("请先创建或加入家庭"))
		case services.ErrInvalidMenuDate:
			c.JSON(http.StatusBadRequest, utils.BadRequest("日期格式错误，请使用YYYY-MM-DD格式"))
		default:
			c.JSON(http.StatusInternalServerError, utils.InternalServerError("获取每日菜单失败"))
		}
		return
	}

	c.JSON(http.StatusOK, utils.Success(resp))
}

// GetWeeklyMenu 获取每周菜单
// @Summary 获取每周菜单
// @Description 获取从开始日期起一周的菜单列表。需要Bearer Token认证。
// @Tags 菜单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string true "开始日期，格式：YYYY-MM-DD"
// @Success 200 {object} utils.Response{data=models.WeeklyMenuResponse} "获取成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 404 {object} utils.Response "尚未加入家庭"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /menus/weekly [get]
func (h *MenuHandler) GetWeeklyMenu(c *gin.Context) {
	req, err := utils.BindQuery[models.WeeklyMenuRequest](c)
	if err != nil {
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	resp, err := h.menuService.GetWeeklyMenu(userID, req)
	if err != nil {
		switch err {
		case services.ErrFamilyNotFound:
			c.JSON(http.StatusNotFound, utils.NotFound("请先创建或加入家庭"))
		case services.ErrInvalidMenuDate:
			c.JSON(http.StatusBadRequest, utils.BadRequest("日期格式错误，请使用YYYY-MM-DD格式"))
		default:
			c.JSON(http.StatusInternalServerError, utils.InternalServerError("获取每周菜单失败"))
		}
		return
	}

	c.JSON(http.StatusOK, utils.Success(resp))
}

// UpdateMenu 更新菜单
// @Summary 更新菜单
// @Description 更新已有菜单，支持添加菜式、删除菜式、修改日期或餐次。需要Bearer Token认证。
// @Tags 菜单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "菜单ID"
// @Param request body models.UpdateMenuRequest true "更新菜单请求"
// @Success 200 {object} utils.Response{data=models.MenuUpdateResponse} "更新成功"
// @Failure 400 {object} utils.Response "参数错误或业务限制"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 404 {object} utils.Response "菜单或家庭不存在"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /menus/{id} [put]
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	uri, err := utils.BindURI[models.MenuIDRequest](c)
	if err != nil {
		return
	}

	req, err := utils.BindJSON[models.UpdateMenuRequest](c)
	if err != nil {
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	resp, err := h.menuService.UpdateMenu(userID, uri.ID, req)
	if err != nil {
		switch err {
		case services.ErrFamilyNotFound:
			c.JSON(http.StatusNotFound, utils.NotFound("请先创建或加入家庭"))
		case services.ErrMenuNotFound:
			c.JSON(http.StatusNotFound, utils.NotFound("菜单不存在或已删除"))
		case services.ErrInvalidMenuDate:
			c.JSON(http.StatusBadRequest, utils.BadRequest("日期格式错误，请使用YYYY-MM-DD格式"))
		case services.ErrInvalidMealType:
			c.JSON(http.StatusBadRequest, utils.BadRequest("餐次类型错误，必须是breakfast、lunch或dinner"))
		case services.ErrDishNotFound:
			c.JSON(http.StatusNotFound, utils.NotFound("菜式不存在或已删除"))
		case services.ErrDishNotInFamily:
			c.JSON(http.StatusBadRequest, utils.BadRequest("菜式不属于当前家庭"))
		default:
			c.JSON(http.StatusInternalServerError, utils.InternalServerError("更新菜单失败"))
		}
		return
	}

	c.JSON(http.StatusOK, utils.SuccessWithMessage("更新成功", resp))
}

