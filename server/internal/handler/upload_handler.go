package handler

import (
	"io"

	"yidaiku-server/internal/middleware"
	"yidaiku-server/internal/model"
	"yidaiku-server/internal/service"
	"yidaiku-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// UploadHandler 上传处理器
type UploadHandler struct {
	uploadService *service.UploadService
}

// NewUploadHandler 创建上传处理器
func NewUploadHandler(uploadService *service.UploadService) *UploadHandler {
	return &UploadHandler{uploadService: uploadService}
}

// InitUpload 初始化上传任务
// @Summary 初始化分块上传任务
// @Tags 上传
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body model.InitUploadRequest true "初始化请求"
// @Success 200 {object} response.Response{data=model.InitUploadResponse}
// @Router /api/v1/upload/init [post]
func (h *UploadHandler) InitUpload(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	var req model.InitUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := h.uploadService.InitUpload(c.Request.Context(), userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// UploadChunk 上传分块
// @Summary 上传文件分块
// @Tags 上传
// @Security Bearer
// @Accept multipart/form-data
// @Produce json
// @Param task_id formData string true "任务ID"
// @Param chunk_index formData int true "分块索引"
// @Param chunk_hash formData string true "分块哈希"
// @Param chunk formData file true "分块数据"
// @Success 200 {object} response.Response{data=model.UploadChunkResponse}
// @Router /api/v1/upload/chunk [post]
func (h *UploadHandler) UploadChunk(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	var req model.UploadChunkRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取上传的文件块
	file, err := c.FormFile("chunk")
	if err != nil {
		response.BadRequest(c, "缺少文件数据")
		return
	}

	f, err := file.Open()
	if err != nil {
		response.InternalError(c, "读取文件失败")
		return
	}
	defer f.Close()

	chunkData, err := io.ReadAll(f)
	if err != nil {
		response.InternalError(c, "读取文件数据失败")
		return
	}

	result, err := h.uploadService.UploadChunk(c.Request.Context(), userID, &req, chunkData)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// CompleteUpload 完成上传
// @Summary 完成分块上传
// @Tags 上传
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body model.CompleteUploadRequest true "完成请求"
// @Success 200 {object} response.Response{data=model.CompleteUploadResponse}
// @Router /api/v1/upload/complete [post]
func (h *UploadHandler) CompleteUpload(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	var req model.CompleteUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result, err := h.uploadService.CompleteUpload(c.Request.Context(), userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetProgress 获取上传进度
// @Summary 获取上传任务进度
// @Tags 上传
// @Security Bearer
// @Produce json
// @Param task_id path string true "任务ID"
// @Success 200 {object} response.Response{data=model.UploadProgress}
// @Router /api/v1/upload/progress/{task_id} [get]
func (h *UploadHandler) GetProgress(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	taskID := c.Param("task_id")
	if taskID == "" {
		response.BadRequest(c, "缺少任务ID")
		return
	}

	progress, err := h.uploadService.GetTaskProgress(c.Request.Context(), userID, taskID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, progress)
}

// GetPendingTasks 获取待处理的上传任务
// @Summary 获取待处理的上传任务列表
// @Tags 上传
// @Security Bearer
// @Produce json
// @Success 200 {object} response.Response{data=[]model.UploadProgress}
// @Router /api/v1/upload/pending [get]
func (h *UploadHandler) GetPendingTasks(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	tasks, err := h.uploadService.GetPendingTasks(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, tasks)
}

// CancelUpload 取消上传
// @Summary 取消上传任务
// @Tags 上传
// @Security Bearer
// @Produce json
// @Param task_id path string true "任务ID"
// @Success 200 {object} response.Response
// @Router /api/v1/upload/{task_id} [delete]
func (h *UploadHandler) CancelUpload(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	taskID := c.Param("task_id")
	if taskID == "" {
		response.BadRequest(c, "缺少任务ID")
		return
	}

	if err := h.uploadService.CancelUpload(c.Request.Context(), userID, taskID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "已取消上传", nil)
}

// RetryUpload 重试上传
// @Summary 重试失败的上传任务
// @Tags 上传
// @Security Bearer
// @Produce json
// @Param task_id path string true "任务ID"
// @Success 200 {object} response.Response
// @Router /api/v1/upload/{task_id}/retry [post]
func (h *UploadHandler) RetryUpload(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	taskID := c.Param("task_id")
	if taskID == "" {
		response.BadRequest(c, "缺少任务ID")
		return
	}

	if err := h.uploadService.RetryUpload(c.Request.Context(), userID, taskID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "已开始重试", nil)
}

// SimpleUpload 简单上传（小文件）
// @Summary 简单上传（小文件，不分块）
// @Tags 上传
// @Security Bearer
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/upload/simple [post]
func (h *UploadHandler) SimpleUpload(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "缺少文件")
		return
	}

	f, err := file.Open()
	if err != nil {
		response.InternalError(c, "打开文件失败")
		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		response.InternalError(c, "读取文件失败")
		return
	}

	url, err := h.uploadService.SimpleUpload(c.Request.Context(), userID, file.Filename, data)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"url":      url,
		"filename": file.Filename,
	})
}
