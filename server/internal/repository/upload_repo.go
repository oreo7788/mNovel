package repository

import (
	"context"
	"errors"
	"time"

	"yidaiku-server/internal/model"

	"gorm.io/gorm"
)

// UploadTaskRepository 上传任务数据访问层
type UploadTaskRepository struct {
	db *gorm.DB
}

// NewUploadTaskRepository 创建上传任务仓库实例
func NewUploadTaskRepository() *UploadTaskRepository {
	return &UploadTaskRepository{db: GetDB()}
}

// Create 创建上传任务
func (r *UploadTaskRepository) Create(ctx context.Context, task *model.UploadTask) error {
	return r.db.WithContext(ctx).Create(task).Error
}

// GetByID 根据ID获取上传任务
func (r *UploadTaskRepository) GetByID(ctx context.Context, taskID string) (*model.UploadTask, error) {
	var task model.UploadTask
	err := r.db.WithContext(ctx).Where("task_id = ?", taskID).First(&task).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &task, err
}

// Update 更新上传任务
func (r *UploadTaskRepository) Update(ctx context.Context, task *model.UploadTask) error {
	return r.db.WithContext(ctx).Save(task).Error
}

// UpdateFields 更新指定字段
func (r *UploadTaskRepository) UpdateFields(ctx context.Context, taskID string, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.UploadTask{}).
		Where("task_id = ?", taskID).
		Updates(fields).Error
}

// UpdateChunks 更新已上传的块
func (r *UploadTaskRepository) UpdateChunks(ctx context.Context, taskID string, chunks []int) error {
	return r.db.WithContext(ctx).Model(&model.UploadTask{}).
		Where("task_id = ?", taskID).
		Update("uploaded_chunks", model.JSONIntArray(chunks)).Error
}

// GetPendingTasks 获取用户待处理的上传任务
func (r *UploadTaskRepository) GetPendingTasks(ctx context.Context, userID string) ([]model.UploadTask, error) {
	var tasks []model.UploadTask
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status IN ?", userID, []string{
			model.UploadStatusPending,
			model.UploadStatusUploading,
			model.UploadStatusResuming,
		}).
		Order("created_at ASC").
		Find(&tasks).Error
	return tasks, err
}

// GetByUserAndStatus 根据用户和状态获取任务
func (r *UploadTaskRepository) GetByUserAndStatus(ctx context.Context, userID string, statuses []string) ([]model.UploadTask, error) {
	var tasks []model.UploadTask
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status IN ?", userID, statuses).
		Order("created_at ASC").
		Find(&tasks).Error
	return tasks, err
}

// SetStatus 设置任务状态
func (r *UploadTaskRepository) SetStatus(ctx context.Context, taskID, status string) error {
	return r.db.WithContext(ctx).Model(&model.UploadTask{}).
		Where("task_id = ?", taskID).
		Update("status", status).Error
}

// SetCompleted 设置任务完成
func (r *UploadTaskRepository) SetCompleted(ctx context.Context, taskID, fileURL string) error {
	return r.db.WithContext(ctx).Model(&model.UploadTask{}).
		Where("task_id = ?", taskID).
		Updates(map[string]interface{}{
			"status":   model.UploadStatusCompleted,
			"file_url": fileURL,
		}).Error
}

// SetFailed 设置任务失败
func (r *UploadTaskRepository) SetFailed(ctx context.Context, taskID, errorMsg string) error {
	return r.db.WithContext(ctx).Model(&model.UploadTask{}).
		Where("task_id = ?", taskID).
		Updates(map[string]interface{}{
			"status":    model.UploadStatusFailed,
			"error_msg": errorMsg,
		}).Error
}

// Delete 删除上传任务
func (r *UploadTaskRepository) Delete(ctx context.Context, taskID string) error {
	return r.db.WithContext(ctx).Delete(&model.UploadTask{}, "task_id = ?", taskID).Error
}

// CleanExpired 清理过期的上传任务
func (r *UploadTaskRepository) CleanExpired(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("expires_at < ? AND status IN ?", time.Now(), []string{
			model.UploadStatusPending,
			model.UploadStatusUploading,
		}).
		Delete(&model.UploadTask{})
	return result.RowsAffected, result.Error
}

// GetByFileHash 根据文件哈希获取已完成的任务（用于秒传）
func (r *UploadTaskRepository) GetByFileHash(ctx context.Context, userID, fileHash string) (*model.UploadTask, error) {
	var task model.UploadTask
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND file_hash = ? AND status = ?", userID, fileHash, model.UploadStatusCompleted).
		First(&task).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &task, err
}
