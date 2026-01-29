package handler

import (
	"strconv"

	"yidaiku-server/internal/middleware"
	"yidaiku-server/internal/model"
	"yidaiku-server/internal/service"
	"yidaiku-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// ClothingHandler 衣物处理器
type ClothingHandler struct {
	clothingService *service.ClothingService
}

// NewClothingHandler 创建衣物处理器
func NewClothingHandler(clothingService *service.ClothingService) *ClothingHandler {
	return &ClothingHandler{clothingService: clothingService}
}

// Create 创建衣物
// @Summary 创建衣物（手动分类）
// @Tags 衣橱
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body service.CreateClothingRequest true "衣物信息"
// @Success 200 {object} response.Response{data=model.ClothingItem}
// @Router /api/v1/wardrobe/items [post]
func (h *ClothingHandler) Create(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	var req service.CreateClothingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	item, err := h.clothingService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, item)
}

// GetByID 获取衣物详情
// @Summary 获取衣物详情
// @Tags 衣橱
// @Security Bearer
// @Produce json
// @Param id path string true "衣物ID"
// @Success 200 {object} response.Response{data=model.ClothingItem}
// @Router /api/v1/wardrobe/items/{id} [get]
func (h *ClothingHandler) GetByID(c *gin.Context) {
	itemID := c.Param("id")
	if itemID == "" {
		response.BadRequest(c, "缺少衣物ID")
		return
	}

	item, err := h.clothingService.GetByID(c.Request.Context(), itemID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if item == nil {
		response.NotFound(c, "衣物不存在")
		return
	}

	response.Success(c, item)
}

// Update 更新衣物
// @Summary 更新衣物信息
// @Tags 衣橱
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "衣物ID"
// @Param request body service.UpdateClothingRequest true "更新信息"
// @Success 200 {object} response.Response{data=model.ClothingItem}
// @Router /api/v1/wardrobe/items/{id} [put]
func (h *ClothingHandler) Update(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	itemID := c.Param("id")
	if itemID == "" {
		response.BadRequest(c, "缺少衣物ID")
		return
	}

	var req service.UpdateClothingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	item, err := h.clothingService.Update(c.Request.Context(), userID, itemID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, item)
}

// Delete 删除衣物
// @Summary 删除衣物（软删除）
// @Tags 衣橱
// @Security Bearer
// @Produce json
// @Param id path string true "衣物ID"
// @Success 200 {object} response.Response
// @Router /api/v1/wardrobe/items/{id} [delete]
func (h *ClothingHandler) Delete(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	itemID := c.Param("id")
	if itemID == "" {
		response.BadRequest(c, "缺少衣物ID")
		return
	}

	// 检查是否永久删除
	permanent := c.Query("permanent") == "true"

	var err error
	if permanent {
		err = h.clothingService.PermanentDelete(c.Request.Context(), userID, itemID)
	} else {
		err = h.clothingService.Delete(c.Request.Context(), userID, itemID)
	}

	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// List 获取衣物列表
// @Summary 获取衣物列表
// @Tags 衣橱
// @Security Bearer
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param category query string false "品类筛选"
// @Param style query string false "风格筛选"
// @Param season query string false "季节筛选"
// @Param color query string false "颜色筛选"
// @Param tag query string false "标签筛选"
// @Param is_retired query bool false "是否已淘汰"
// @Param keyword query string false "关键词搜索"
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/v1/wardrobe/items [get]
func (h *ClothingHandler) List(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	// 解析筛选参数
	filter := &model.ClothingFilter{
		UserID: userID,
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	filter.Page = page
	filter.PageSize = pageSize

	// 品类
	if category := c.Query("category"); category != "" {
		filter.Categories = []string{category}
	}

	// 风格
	if style := c.Query("style"); style != "" {
		filter.Styles = []string{style}
	}

	// 季节
	if season := c.Query("season"); season != "" {
		filter.Seasons = []string{season}
	}

	// 颜色
	if color := c.Query("color"); color != "" {
		filter.Colors = []string{color}
	}

	// 标签
	if tag := c.Query("tag"); tag != "" {
		filter.Tags = []string{tag}
	}

	// 已淘汰状态
	if isRetiredStr := c.Query("is_retired"); isRetiredStr != "" {
		isRetired := isRetiredStr == "true"
		filter.IsRetired = &isRetired
	}

	// 关键词
	filter.Keyword = c.Query("keyword")

	items, total, err := h.clothingService.List(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, items, total, page, pageSize)
}

// SetRetired 设置淘汰状态
// @Summary 设置衣物淘汰状态
// @Tags 衣橱
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "衣物ID"
// @Param request body map[string]bool true "淘汰状态"
// @Success 200 {object} response.Response
// @Router /api/v1/wardrobe/items/{id}/retired [put]
func (h *ClothingHandler) SetRetired(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	itemID := c.Param("id")
	if itemID == "" {
		response.BadRequest(c, "缺少衣物ID")
		return
	}

	var req struct {
		Retired bool `json:"retired"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.clothingService.SetRetired(c.Request.Context(), userID, itemID, req.Retired); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	msg := "已标记为淘汰"
	if !req.Retired {
		msg = "已恢复为可用"
	}
	response.SuccessWithMessage(c, msg, nil)
}

// BatchSetRetired 批量设置淘汰状态
// @Summary 批量设置淘汰状态
// @Tags 衣橱
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "批量操作"
// @Success 200 {object} response.Response
// @Router /api/v1/wardrobe/items/batch/retired [put]
func (h *ClothingHandler) BatchSetRetired(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	var req struct {
		ItemIDs []string `json:"item_ids" binding:"required"`
		Retired bool     `json:"retired"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.clothingService.BatchSetRetired(c.Request.Context(), userID, req.ItemIDs, req.Retired); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "批量操作成功", nil)
}

// BatchDelete 批量删除
// @Summary 批量删除衣物
// @Tags 衣橱
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body map[string][]string true "批量删除"
// @Success 200 {object} response.Response
// @Router /api/v1/wardrobe/items/batch [delete]
func (h *ClothingHandler) BatchDelete(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	var req struct {
		ItemIDs []string `json:"item_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.clothingService.BatchDelete(c.Request.Context(), userID, req.ItemIDs); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "批量删除成功", nil)
}

// GetStats 获取衣橱统计
// @Summary 获取衣橱统计
// @Tags 衣橱
// @Security Bearer
// @Produce json
// @Success 200 {object} response.Response{data=model.ClothingStats}
// @Router /api/v1/wardrobe/stats [get]
func (h *ClothingHandler) GetStats(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	stats, err := h.clothingService.GetStats(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// Restore 恢复已删除的衣物
// @Summary 恢复已删除的衣物
// @Tags 衣橱
// @Security Bearer
// @Produce json
// @Param id path string true "衣物ID"
// @Success 200 {object} response.Response
// @Router /api/v1/wardrobe/items/{id}/restore [post]
func (h *ClothingHandler) Restore(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	itemID := c.Param("id")
	if itemID == "" {
		response.BadRequest(c, "缺少衣物ID")
		return
	}

	if err := h.clothingService.Restore(c.Request.Context(), userID, itemID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "恢复成功", nil)
}

// GetCategories 获取品类列表
// @Summary 获取品类列表
// @Tags 衣橱
// @Produce json
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/v1/wardrobe/categories [get]
func (h *ClothingHandler) GetCategories(c *gin.Context) {
	response.Success(c, model.Categories)
}

// GetStyles 获取风格列表
// @Summary 获取风格列表
// @Tags 衣橱
// @Produce json
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/v1/wardrobe/styles [get]
func (h *ClothingHandler) GetStyles(c *gin.Context) {
	response.Success(c, model.Styles)
}

// GetSeasons 获取季节列表
// @Summary 获取季节列表
// @Tags 衣橱
// @Produce json
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/v1/wardrobe/seasons [get]
func (h *ClothingHandler) GetSeasons(c *gin.Context) {
	response.Success(c, model.Seasons)
}

// GetColors 获取颜色列表
// @Summary 获取颜色列表
// @Tags 衣橱
// @Produce json
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/v1/wardrobe/colors [get]
func (h *ClothingHandler) GetColors(c *gin.Context) {
	response.Success(c, model.Colors)
}
