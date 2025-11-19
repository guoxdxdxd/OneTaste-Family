package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"onetaste-family/backend/internal/models"
	"onetaste-family/backend/internal/services"
	"onetaste-family/backend/internal/utils"
)

// FamilyHandler 家庭相关处理器
type FamilyHandler struct {
	familyService *services.FamilyService
}

// NewFamilyHandler 创建家庭处理器
func NewFamilyHandler() *FamilyHandler {
	return &FamilyHandler{
		familyService: services.NewFamilyService(),
	}
}

// CreateFamily 创建家庭
// @Summary 创建家庭
// @Description 创建一个新的家庭，创建者自动成为owner并加入家庭。需要Bearer Token认证。
// @Tags 家庭
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateFamilyRequest true "创建家庭请求参数"
// @Success 200 {object} utils.Response{data=models.FamilyResponse} "创建成功"
// @Failure 400 {object} utils.Response "请求参数错误或已加入家庭"
// @Failure 401 {object} utils.Response "未授权或Token无效"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /family/create [post]
func (h *FamilyHandler) CreateFamily(c *gin.Context) {
	req, err := utils.BindJSON[models.CreateFamilyRequest](c)
	if err != nil {
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	resp, err := h.familyService.CreateFamily(req, userID)
	if err != nil {
		switch err {
		case services.ErrInvalidFamilyName, services.ErrInvalidFamilyDescription:
			c.JSON(http.StatusBadRequest, utils.BadRequest("家庭名称或描述不合法"))
			return
		case services.ErrUserAlreadyInFamily:
			c.JSON(http.StatusBadRequest, utils.BadRequest("您已经创建或加入了一个家庭"))
			return
		default:
			c.JSON(http.StatusInternalServerError, utils.InternalServerError("创建家庭失败"))
			return
		}
	}

	c.JSON(http.StatusOK, utils.SuccessWithMessage("创建成功", resp))
}

// GetFamilyInfo 获取家庭信息
// @Summary 获取家庭信息
// @Description 获取当前用户所属家庭信息，包括成员数量与菜式数量。需要Bearer Token认证。
// @Tags 家庭
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=models.FamilyInfoResponse} "获取成功"
// @Failure 401 {object} utils.Response "未授权或Token无效"
// @Failure 404 {object} utils.Response "尚未创建或加入家庭"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /family/info [get]
func (h *FamilyHandler) GetFamilyInfo(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	info, err := h.familyService.GetFamilyInfo(userID)
	if err != nil {
		switch err {
		case services.ErrFamilyNotFound:
			c.JSON(http.StatusNotFound, utils.NotFound("您还没有创建或加入任何家庭"))
			return
		default:
			c.JSON(http.StatusInternalServerError, utils.InternalServerError("获取家庭信息失败"))
			return
		}
	}

	c.JSON(http.StatusOK, utils.Success(info))
}

// JoinFamilyViaInvite 扫码加入家庭
// @Summary 扫码加入家庭
// @Description 用户在前端点击“同意”后调用该接口，加入邀请人所在家庭。需要Bearer Token认证。
// @Tags 家庭
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.FamilyInviteRequest true "扫码加入家庭请求参数"
// @Success 200 {object} utils.Response{data=models.FamilyJoinResponse} "加入成功"
// @Failure 400 {object} utils.Response "参数错误或业务限制"
// @Failure 401 {object} utils.Response "未授权或Token无效"
// @Failure 404 {object} utils.Response "家庭不存在或邀请无效"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /family/member/invite [post]
func (h *FamilyHandler) JoinFamilyViaInvite(c *gin.Context) {
	req, err := utils.BindJSON[models.FamilyInviteRequest](c)
	if err != nil {
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	resp, err := h.familyService.JoinFamilyViaInvite(req, userID)
	if err != nil {
		switch err {
		case services.ErrInvalidInviteAction, services.ErrFamilyNameMismatch:
			c.JSON(http.StatusBadRequest, utils.BadRequest("邀请信息不正确"))
			return
		case services.ErrUserAlreadyInFamily:
			c.JSON(http.StatusBadRequest, utils.BadRequest("您已经创建或加入了一个家庭"))
			return
		case services.ErrInviterNotInFamily:
			c.JSON(http.StatusBadRequest, utils.BadRequest("邀请人信息无效"))
			return
		case services.ErrFamilyMemberLimitReached:
			c.JSON(http.StatusBadRequest, utils.BadRequest("家庭成员数量已达上限"))
			return
		case services.ErrFamilyNotFound:
			c.JSON(http.StatusNotFound, utils.NotFound("家庭不存在或已解散"))
			return
		default:
			c.JSON(http.StatusInternalServerError, utils.InternalServerError("加入家庭失败"))
			return
		}
	}

	c.JSON(http.StatusOK, utils.SuccessWithMessage("加入家庭成功", resp))
}

// GetFamilyMembers 获取家庭成员列表
// @Summary 获取家庭成员列表
// @Description 返回当前用户所属家庭的成员列表。需要Bearer Token认证。
// @Tags 家庭
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]models.FamilyMemberInfo} "获取成功"
// @Failure 401 {object} utils.Response "未授权或Token无效"
// @Failure 404 {object} utils.Response "尚未创建或加入家庭"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /family/members [get]
func (h *FamilyHandler) GetFamilyMembers(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	members, err := h.familyService.GetFamilyMembers(userID)
	if err != nil {
		switch err {
		case services.ErrFamilyNotFound:
			c.JSON(http.StatusNotFound, utils.NotFound("您还没有创建或加入任何家庭"))
			return
		default:
			c.JSON(http.StatusInternalServerError, utils.InternalServerError("获取家庭成员失败"))
			return
		}
	}

	c.JSON(http.StatusOK, utils.Success(members))
}

func getUserIDFromContext(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.Unauthorized("未授权"))
		return 0, false
	}

	id, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.Unauthorized("无效的用户ID"))
		return 0, false
	}

	return id, true
}
