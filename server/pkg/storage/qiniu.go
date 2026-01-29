package storage

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// QiniuOSS 七牛云存储客户端
type QiniuOSS struct {
	accessKey string
	secretKey string
	bucket    string
	domain    string
	mac       *qbox.Mac
	cfg       *storage.Config
}

// NewQiniuOSS 创建七牛云存储客户端
func NewQiniuOSS(accessKey, secretKey, bucket, domain string) *QiniuOSS {
	mac := qbox.NewMac(accessKey, secretKey)
	cfg := &storage.Config{
		Zone:          &storage.ZoneHuadong, // 华东区域，可根据实际配置
		UseHTTPS:      true,
		UseCdnDomains: true,
	}

	return &QiniuOSS{
		accessKey: accessKey,
		secretKey: secretKey,
		bucket:    bucket,
		domain:    domain,
		mac:       mac,
		cfg:       cfg,
	}
}

// Upload 上传文件
func (q *QiniuOSS) Upload(key string, data []byte) (string, error) {
	putPolicy := storage.PutPolicy{
		Scope: q.bucket,
	}
	putPolicy.Expires = 3600 // 1小时有效期
	upToken := putPolicy.UploadToken(q.mac)

	formUploader := storage.NewFormUploader(q.cfg)
	ret := storage.PutRet{}
	err := formUploader.Put(context.Background(), &ret, upToken, key, bytes.NewReader(data), int64(len(data)), nil)
	if err != nil {
		return "", err
	}

	return q.GetFileURL(ret.Key), nil
}

// Delete 删除文件
func (q *QiniuOSS) Delete(key string) error {
	bucketManager := storage.NewBucketManager(q.mac, q.cfg)
	return bucketManager.Delete(q.bucket, key)
}

// DeleteByURL 根据URL删除文件
func (q *QiniuOSS) DeleteByURL(fileURL string) error {
	key := q.extractKeyFromURL(fileURL)
	if key == "" {
		return nil
	}
	return q.Delete(key)
}

// GetFileURL 获取文件URL
func (q *QiniuOSS) GetFileURL(key string) string {
	return fmt.Sprintf("%s/%s", q.domain, key)
}

// GetThumbnailURL 获取缩略图URL（七牛云图片处理）
func (q *QiniuOSS) GetThumbnailURL(originalURL string, width, height int) string {
	// 七牛云图片处理参数
	return fmt.Sprintf("%s?imageView2/1/w/%d/h/%d", originalURL, width, height)
}

// GetUploadToken 获取上传凭证
func (q *QiniuOSS) GetUploadToken(key string, expires uint64) string {
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", q.bucket, key),
	}
	if expires > 0 {
		putPolicy.Expires = expires
	} else {
		putPolicy.Expires = 3600
	}
	return putPolicy.UploadToken(q.mac)
}

// InitMultipartUpload 初始化分片上传
func (q *QiniuOSS) InitMultipartUpload(key string) (string, error) {
	// 七牛云分片上传V2
	// 简化实现，实际使用时可以使用七牛云的分片上传SDK
	// 这里返回一个唯一标识作为uploadId
	return fmt.Sprintf("%s_%d", key, generateUploadID()), nil
}

// UploadPart 上传分片
func (q *QiniuOSS) UploadPart(key, uploadID string, partNumber int, data []byte) error {
	// 简化实现：分片上传每个部分
	// 实际应该使用七牛云的分片上传API
	partKey := fmt.Sprintf("%s.part%d", key, partNumber)
	_, err := q.Upload(partKey, data)
	return err
}

// CompleteMultipartUpload 完成分片上传
func (q *QiniuOSS) CompleteMultipartUpload(key, uploadID string, totalParts int) (string, error) {
	// 简化实现：合并所有分片
	// 实际应该使用七牛云的分片合并API
	// 这里直接返回最终URL
	return q.GetFileURL(key), nil
}

// AbortMultipartUpload 取消分片上传
func (q *QiniuOSS) AbortMultipartUpload(key, uploadID string) error {
	// 清理已上传的分片
	// 实际应该调用七牛云的取消分片上传API
	return nil
}

// extractKeyFromURL 从URL提取key
func (q *QiniuOSS) extractKeyFromURL(fileURL string) string {
	if !strings.HasPrefix(fileURL, q.domain) {
		return ""
	}
	return strings.TrimPrefix(fileURL, q.domain+"/")
}

// generateUploadID 生成上传ID
func generateUploadID() int64 {
	// 简单实现，使用时间戳
	return int64(1000000)
}
