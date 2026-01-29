package handler

import (
	"strconv"

	"yidaiku-server/internal/middleware"
	"yidaiku-server/internal/model"
	"yidaiku-server/internal/service"
	"yidaiku-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// OutfitHandler 穿搭推荐处理器
type OutfitHandler struct {
	outfitService *service.OutfitService
}

// NewOutfitHandler 创建穿搭处理器
func NewOutfitHandler(outfitService *service.OutfitService) *OutfitHandler {
	return &OutfitHandler{outfitService: outfitService}
}

// Recommend 生成穿搭推荐
// @Summary 生成穿搭推荐
// @Tags 推荐
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body model.RecommendRequest true "推荐请求"
// @Success 200 {object} response.Response{data=service.RecommendResult}
// @Router /api/v1/outfits/recommend [post]
func (h *OutfitHandler) Recommend(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	var req model.RecommendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := h.outfitService.Recommend(c.Request.Context(), userID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// ReplaceItem 换一件
// @Summary 换一件功能
// @Tags 推荐
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body model.ReplaceItemRequest true "替换请求"
// @Success 200 {object} response.Response{data=service.OutfitDetail}
// @Router /api/v1/outfits/replace [post]
func (h *OutfitHandler) ReplaceItem(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	var req model.ReplaceItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result, err := h.outfitService.ReplaceItem(c.Request.Context(), userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// SetFavorite 收藏/取消收藏推荐
// @Summary 收藏/取消收藏推荐
// @Tags 推荐
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "推荐ID"
// @Param request body map[string]bool true "收藏状态"
// @Success 200 {object} response.Response
// @Router /api/v1/outfits/{id}/favorite [put]
func (h *OutfitHandler) SetFavorite(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	outfitID := c.Param("id")
	if outfitID == "" {
		response.BadRequest(c, "缺少推荐ID")
		return
	}

	var req struct {
		Favorite bool `json:"favorite"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.outfitService.SetFavorite(c.Request.Context(), userID, outfitID, req.Favorite); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	msg := "已收藏"
	if !req.Favorite {
		msg = "已取消收藏"
	}
	response.SuccessWithMessage(c, msg, nil)
}

// SetApplied 标记为已应用
// @Summary 标记推荐为已应用
// @Tags 推荐
// @Security Bearer
// @Produce json
// @Param id path string true "推荐ID"
// @Success 200 {object} response.Response
// @Router /api/v1/outfits/{id}/apply [post]
func (h *OutfitHandler) SetApplied(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	outfitID := c.Param("id")
	if outfitID == "" {
		response.BadRequest(c, "缺少推荐ID")
		return
	}

	if err := h.outfitService.SetApplied(c.Request.Context(), userID, outfitID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "已标记为今日穿搭", nil)
}

// GetFavorites 获取收藏列表
// @Summary 获取收藏的穿搭列表
// @Tags 推荐
// @Security Bearer
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/v1/outfits/favorites [get]
func (h *OutfitHandler) GetFavorites(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	outfits, total, err := h.outfitService.GetFavorites(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, outfits, total, page, pageSize)
}

// GetTemplates 获取穿搭模板
// @Summary 获取穿搭模板（用于冷启动）
// @Tags 推荐
// @Produce json
// @Param occasion query string false "场合筛选"
// @Success 200 {object} response.Response{data=[]model.OutfitTemplate}
// @Router /api/v1/outfits/templates [get]
func (h *OutfitHandler) GetTemplates(c *gin.Context) {
	occasion := c.Query("occasion")

	templates, err := h.outfitService.GetTemplates(c.Request.Context(), occasion)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, templates)
}

// GetOccasions 获取场合列表
// @Summary 获取场合列表
// @Tags 推荐
// @Produce json
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/v1/outfits/occasions [get]
func (h *OutfitHandler) GetOccasions(c *gin.Context) {
	response.Success(c, model.Occasions)
}
