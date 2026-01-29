package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"yidaiku-server/internal/config"
	"yidaiku-server/internal/model"
	"yidaiku-server/internal/repository"
	"yidaiku-server/pkg/storage"

	"github.com/google/uuid"
)

// UploadService 上传服务
type UploadService struct {
	uploadRepo *repository.UploadTaskRepository
	ossClient  *storage.QiniuOSS
	config     *config.UploadConfig
}

// NewUploadService 创建上传服务实例
func NewUploadService(ossClient *storage.QiniuOSS, cfg *config.UploadConfig) *UploadService {
	return &UploadService{
		uploadRepo: repository.NewUploadTaskRepository(),
		ossClient:  ossClient,
		config:     cfg,
	}
}

// InitUpload 初始化上传任务
func (s *UploadService) InitUpload(ctx context.Context, userID string, req *model.InitUploadRequest) (*model.InitUploadResponse, error) {
	// 验证文件大小
	if req.FileSize > s.config.MaxFileSize {
		return nil, errors.New("文件大小超过限制")
	}

	// 检查是否已有相同文件（秒传）
	existingTask, err := s.uploadRepo.GetByFileHash(ctx, userID, req.FileHash)
	if err == nil && existingTask != nil {
		return &model.InitUploadResponse{
			TaskID:    existingTask.ID,
			UploadID:  existingTask.UploadID,
			ChunkSize: s.config.ChunkSize,
			ExpiresAt: existingTask.ExpiresAt,
		}, nil
	}

	// 创建上传任务
	taskID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	// 初始化七牛云分片上传
	uploadID := ""
	if s.ossClient != nil {
		key := s.generateFileKey(userID, taskID, req.FileName)
		uploadID, err = s.ossClient.InitMultipartUpload(key)
		if err != nil {
			return nil, errors.New("初始化上传失败: " + err.Error())
		}
	}

	task := &model.UploadTask{
		ID:             taskID,
		UserID:         userID,
		FileName:       req.FileName,
		FileSize:       req.FileSize,
		FileHash:       req.FileHash,
		TotalChunks:    req.TotalChunks,
		UploadedChunks: model.JSONArray{},
		UploadID:       uploadID,
		Status:         model.UploadStatusPending,
		ExpiresAt:      expiresAt,
	}

	if err := s.uploadRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	return &model.InitUploadResponse{
		TaskID:    taskID,
		UploadID:  uploadID,
		ChunkSize: s.config.ChunkSize,
		ExpiresAt: expiresAt,
	}, nil
}

// UploadChunk 上传分块
func (s *UploadService) UploadChunk(ctx context.Context, userID string, req *model.UploadChunkRequest, chunkData []byte) (*model.UploadChunkResponse, error) {
	// 获取任务
	task, err := s.uploadRepo.GetByID(ctx, req.TaskID)
	if err != nil || task == nil {
		return nil, errors.New("上传任务不存在")
	}
	if task.UserID != userID {
		return nil, errors.New("无权限操作此任务")
	}
	if task.Status == model.UploadStatusCompleted {
		return nil, errors.New("任务已完成")
	}

	// 验证分块索引
	if req.ChunkIndex < 0 || req.ChunkIndex >= task.TotalChunks {
		return nil, errors.New("无效的分块索引")
	}

	// 验证分块哈希
	hash := md5.Sum(chunkData)
	chunkHash := hex.EncodeToString(hash[:])
	if chunkHash != req.ChunkHash {
		return nil, errors.New("分块数据校验失败")
	}

	// 检查是否已上传
	chunkIdxStr := strconv.Itoa(req.ChunkIndex)
	for _, idx := range task.UploadedChunks {
		if idx == chunkIdxStr {
			return s.buildChunkResponse(task)
		}
	}
	uploadedChunks := task.UploadedChunks

	// 上传到七牛云
	if s.ossClient != nil {
		key := s.generateFileKey(userID, task.ID, task.FileName)
		err = s.ossClient.UploadPart(key, task.UploadID, req.ChunkIndex+1, chunkData)
		if err != nil {
			return nil, errors.New("上传分块失败: " + err.Error())
		}
	}

	// 更新任务状态
	task.Status = model.UploadStatusUploading
	uploadedChunks = append(uploadedChunks, chunkIdxStr)
	task.UploadedChunks = uploadedChunks

	if err := s.uploadRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return s.buildChunkResponse(task)
}

// CompleteUpload 完成上传
func (s *UploadService) CompleteUpload(ctx context.Context, userID string, req *model.CompleteUploadRequest) (*model.CompleteUploadResponse, error) {
	// 获取任务
	task, err := s.uploadRepo.GetByID(ctx, req.TaskID)
	if err != nil || task == nil {
		return nil, errors.New("上传任务不存在")
	}
	if task.UserID != userID {
		return nil, errors.New("无权限操作此任务")
	}

	// 验证文件哈希
	if task.FileHash != req.FileHash {
		return nil, errors.New("文件校验失败")
	}

	// 检查所有分块是否已上传
	if len(task.UploadedChunks) < task.TotalChunks {
		return nil, errors.New("还有分块未上传完成")
	}

	// 完成七牛云分片上传
	fileURL := ""
	if s.ossClient != nil {
		key := s.generateFileKey(userID, task.ID, task.FileName)
		url, err := s.ossClient.CompleteMultipartUpload(key, task.UploadID, task.TotalChunks)
		if err != nil {
			return nil, errors.New("完成上传失败: " + err.Error())
		}
		fileURL = url
	}

	// 更新任务状态
	if err := s.uploadRepo.SetCompleted(ctx, task.ID, fileURL); err != nil {
		return nil, err
	}

	return &model.CompleteUploadResponse{
		FileURL: fileURL,
		FileID:  task.ID,
	}, nil
}

// GetTaskProgress 获取任务进度
func (s *UploadService) GetTaskProgress(ctx context.Context, userID, taskID string) (*model.UploadProgress, error) {
	task, err := s.uploadRepo.GetByID(ctx, taskID)
	if err != nil || task == nil {
		return nil, errors.New("上传任务不存在")
	}
	if task.UserID != userID {
		return nil, errors.New("无权限查看此任务")
	}

	uploadedCount := len(task.UploadedChunks)
	progress := float64(uploadedCount) / float64(task.TotalChunks) * 100

	return &model.UploadProgress{
		TaskID:         task.ID,
		FileName:       task.FileName,
		Status:         task.Status,
		TotalChunks:    task.TotalChunks,
		UploadedChunks: uploadedCount,
		Progress:       progress,
		ErrorMsg:       task.ErrorMsg,
	}, nil
}

// GetPendingTasks 获取待处理的上传任务
func (s *UploadService) GetPendingTasks(ctx context.Context, userID string) ([]model.UploadProgress, error) {
	tasks, err := s.uploadRepo.GetPendingTasks(ctx, userID)
	if err != nil {
		return nil, err
	}

	var progresses []model.UploadProgress
	for _, task := range tasks {
		uploadedCount := len(task.UploadedChunks)
		progress := float64(uploadedCount) / float64(task.TotalChunks) * 100

		progresses = append(progresses, model.UploadProgress{
			TaskID:         task.ID,
			FileName:       task.FileName,
			Status:         task.Status,
			TotalChunks:    task.TotalChunks,
			UploadedChunks: uploadedCount,
			Progress:       progress,
			ErrorMsg:       task.ErrorMsg,
		})
	}

	return progresses, nil
}

// CancelUpload 取消上传
func (s *UploadService) CancelUpload(ctx context.Context, userID, taskID string) error {
	task, err := s.uploadRepo.GetByID(ctx, taskID)
	if err != nil || task == nil {
		return errors.New("上传任务不存在")
	}
	if task.UserID != userID {
		return errors.New("无权限操作此任务")
	}

	// 取消七牛云分片上传
	if s.ossClient != nil && task.UploadID != "" {
		key := s.generateFileKey(userID, task.ID, task.FileName)
		s.ossClient.AbortMultipartUpload(key, task.UploadID)
	}

	return s.uploadRepo.Delete(ctx, taskID)
}

// RetryUpload 重试上传
func (s *UploadService) RetryUpload(ctx context.Context, userID, taskID string) error {
	task, err := s.uploadRepo.GetByID(ctx, taskID)
	if err != nil || task == nil {
		return errors.New("上传任务不存在")
	}
	if task.UserID != userID {
		return errors.New("无权限操作此任务")
	}

	if task.Status != model.UploadStatusFailed {
		return errors.New("任务状态不允许重试")
	}

	return s.uploadRepo.SetStatus(ctx, taskID, model.UploadStatusResuming)
}

// CleanExpiredTasks 清理过期任务
func (s *UploadService) CleanExpiredTasks(ctx context.Context) (int64, error) {
	return s.uploadRepo.CleanExpired(ctx)
}

// SimpleUpload 简单上传（小文件，不分块）
func (s *UploadService) SimpleUpload(ctx context.Context, userID, fileName string, data []byte) (string, error) {
	if int64(len(data)) > s.config.MaxFileSize {
		return "", errors.New("文件大小超过限制")
	}

	// 直接上传到七牛云
	if s.ossClient != nil {
		key := s.generateSimpleFileKey(userID, fileName)
		url, err := s.ossClient.Upload(key, data)
		if err != nil {
			return "", errors.New("上传失败: " + err.Error())
		}
		return url, nil
	}

	return "", errors.New("存储服务未配置")
}

// buildChunkResponse 构建分块响应
func (s *UploadService) buildChunkResponse(task *model.UploadTask) (*model.UploadChunkResponse, error) {
	uploadedCount := len(task.UploadedChunks)
	progress := float64(uploadedCount) / float64(task.TotalChunks)

	// 转换已上传的块索引（JSON 存的是字符串数组）
	var uploadedIndices []int
	for _, chunk := range task.UploadedChunks {
		idx, err := strconv.Atoi(chunk)
		if err == nil {
			uploadedIndices = append(uploadedIndices, idx)
		}
	}

	return &model.UploadChunkResponse{
		Status:         task.Status,
		UploadedChunks: uploadedIndices,
		Progress:       progress,
	}, nil
}

// generateFileKey 生成文件存储路径
func (s *UploadService) generateFileKey(userID, taskID, fileName string) string {
	return "users/" + userID + "/items/" + taskID + "/" + fileName
}

// generateSimpleFileKey 生成简单文件存储路径
func (s *UploadService) generateSimpleFileKey(userID, fileName string) string {
	timestamp := time.Now().UnixNano()
	return "users/" + userID + "/uploads/" + strconv.FormatInt(timestamp, 10) + "_" + fileName
}
