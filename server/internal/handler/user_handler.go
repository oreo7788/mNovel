package handler

import (
	"errors"

	"yidaiku-server/internal/middleware"
	"yidaiku-server/internal/model"
	"yidaiku-server/internal/service"
	"yidaiku-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register 用户注册
// @Summary 用户注册
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response{data=service.AuthResponse}
// @Router /api/v1/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrPhoneAlreadyRegistered) {
			response.Conflict(c, err.Error())
			return
		}
		if errors.Is(err, service.ErrInvalidPhoneFormat) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// Login 用户登录
// @Summary 用户登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=service.AuthResponse}
// @Router /api/v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// RefreshToken 刷新Token
// @Summary 刷新Token
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body map[string]string true "刷新令牌"
// @Success 200 {object} response.Response{data=service.AuthResponse}
// @Router /api/v1/auth/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result, err := h.userService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetProfile 获取用户信息
// @Summary 获取当前用户信息
// @Tags 用户
// @Security Bearer
// @Produce json
// @Success 200 {object} response.Response{data=model.User}
// @Router /api/v1/user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil || user == nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}

// UpdateProfile 更新用户信息
// @Summary 更新用户信息
// @Tags 用户
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "更新信息"
// @Success 200 {object} response.Response
// @Router /api/v1/user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.userService.UpdateProfile(c.Request.Context(), userID, updates); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Tags 用户
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body map[string]string true "密码信息"
// @Success 200 {object} response.Response
// @Router /api/v1/user/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.userService.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "密码修改成功", nil)
}

// GetPreference 获取用户偏好设置
// @Summary 获取用户偏好设置
// @Tags 用户
// @Security Bearer
// @Produce json
// @Success 200 {object} response.Response{data=model.UserPreference}
// @Router /api/v1/user/preference [get]
func (h *UserHandler) GetPreference(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	pref, err := h.userService.GetPreference(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, pref)
}

// UpdatePreference 更新用户偏好设置
// @Summary 更新用户偏好设置
// @Tags 用户
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body model.UserPreference true "偏好设置"
// @Success 200 {object} response.Response
// @Router /api/v1/user/preference [put]
func (h *UserHandler) UpdatePreference(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	var pref model.UserPreference
	if err := c.ShouldBindJSON(&pref); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	pref.UserID = userID
	if err := h.userService.UpdatePreference(c.Request.Context(), &pref); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "偏好设置已保存", nil)
}
