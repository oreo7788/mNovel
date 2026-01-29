package model

import (
	"time"
)

// UploadTask 上传任务
type UploadTask struct {
	ID             string    `json:"id" gorm:"column:task_id;primaryKey;type:varchar(64)"`
	UserID         string    `json:"user_id" gorm:"column:user_id;type:varchar(20);index"`
	FileName       string    `json:"file_name" gorm:"column:file_name;type:varchar(255)"`
	FileSize       int64     `json:"file_size" gorm:"column:file_size;type:bigint"`
	FileHash       string    `json:"file_hash" gorm:"column:file_hash;type:varchar(64)"`
	TotalChunks    int       `json:"total_chunks" gorm:"column:total_chunks;type:int"`
	UploadedChunks JSONArray `json:"uploaded_chunks" gorm:"column:uploaded_chunks;type:json"`              // 已上传块索引数组
	UploadID       string    `json:"upload_id" gorm:"column:upload_id;type:varchar(64)"`                   // 七牛云返回的upload_id
	Status         string    `json:"status" gorm:"column:status;type:varchar(16);default:'pending';index"` // pending/uploading/completed/failed
	FileURL        string    `json:"file_url" gorm:"column:file_url;type:varchar(512)"`                    // 上传完成后的文件URL
	ErrorMsg       string    `json:"error_msg" gorm:"column:error_msg;type:varchar(255)"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	ExpiresAt      time.Time `json:"expires_at" gorm:"column:expires_at"`
}

// TableName 表名
func (UploadTask) TableName() string {
	return "upload_tasks"
}

// UploadTaskStatus 上传任务状态常量
const (
	UploadStatusPending   = "pending"
	UploadStatusUploading = "uploading"
	UploadStatusResuming  = "resuming"
	UploadStatusCompleted = "completed"
	UploadStatusFailed    = "failed"
)

// InitUploadRequest 初始化上传请求
type InitUploadRequest struct {
	FileName    string `json:"file_name" binding:"required"`
	FileSize    int64  `json:"file_size" binding:"required,gt=0"`
	FileHash    string `json:"file_hash" binding:"required"`
	TotalChunks int    `json:"total_chunks" binding:"required,gt=0"`
}

// InitUploadResponse 初始化上传响应
type InitUploadResponse struct {
	TaskID    string    `json:"task_id"`
	UploadID  string    `json:"upload_id"`
	ChunkSize int64     `json:"chunk_size"`
	ExpiresAt time.Time `json:"expires_at"`
}

// UploadChunkRequest 上传分块请求
type UploadChunkRequest struct {
	TaskID     string `form:"task_id" binding:"required"`
	ChunkIndex int    `form:"chunk_index" binding:"required,gte=0"`
	ChunkHash  string `form:"chunk_hash" binding:"required"`
}

// UploadChunkResponse 上传分块响应
type UploadChunkResponse struct {
	Status         string  `json:"status"`
	UploadedChunks []int   `json:"uploaded_chunks"`
	Progress       float64 `json:"progress"`
}

// CompleteUploadRequest 完成上传请求
type CompleteUploadRequest struct {
	TaskID   string `json:"task_id" binding:"required"`
	FileHash string `json:"file_hash" binding:"required"`
}

// CompleteUploadResponse 完成上传响应
type CompleteUploadResponse struct {
	FileURL string `json:"file_url"`
	FileID  string `json:"file_id"`
}

// UploadProgress 上传进度
type UploadProgress struct {
	TaskID         string  `json:"task_id"`
	FileName       string  `json:"file_name"`
	Status         string  `json:"status"`
	TotalChunks    int     `json:"total_chunks"`
	UploadedChunks int     `json:"uploaded_chunks"`
	Progress       float64 `json:"progress"` // 0-100
	ErrorMsg       string  `json:"error_msg,omitempty"`
}
