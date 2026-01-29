package service

import (
	"context"
	"errors"

	"yidaiku-server/internal/model"
	"yidaiku-server/internal/repository"
	"yidaiku-server/pkg/storage"

	"github.com/google/uuid"
)

// ClothingService 衣物服务
type ClothingService struct {
	clothingRepo *repository.ClothingRepository
	ossClient    *storage.QiniuOSS
}

// NewClothingService 创建衣物服务实例
func NewClothingService(ossClient *storage.QiniuOSS) *ClothingService {
	return &ClothingService{
		clothingRepo: repository.NewClothingRepository(),
		ossClient:    ossClient,
	}
}

// CreateClothingRequest 创建衣物请求（MVP版本：手动分类）
type CreateClothingRequest struct {
	ImageURL    string   `json:"image_url" binding:"required"`
	Category    string   `json:"category" binding:"required"`
	Subcategory string   `json:"subcategory"`
	Colors      []string `json:"colors"`
	Styles      []string `json:"styles"`
	Season      string   `json:"season"`
	Tags        []string `json:"tags"`
	Notes       string   `json:"notes"`
}

// UpdateClothingRequest 更新衣物请求
type UpdateClothingRequest struct {
	Category    *string  `json:"category"`
	Subcategory *string  `json:"subcategory"`
	Colors      []string `json:"colors"`
	Styles      []string `json:"styles"`
	Season      *string  `json:"season"`
	Tags        []string `json:"tags"`
	Notes       *string  `json:"notes"`
}

// Create 创建衣物（MVP版本：用户手动填写属性）
func (s *ClothingService) Create(ctx context.Context, userID string, req *CreateClothingRequest) (*model.ClothingItem, error) {
	// 验证品类
	if !isValidCategory(req.Category) {
		return nil, errors.New("无效的品类")
	}

	// 验证季节
	if req.Season != "" && !isValidSeason(req.Season) {
		return nil, errors.New("无效的季节")
	}

	// 生成缩略图URL（如果OSS支持图片处理）
	thumbnailURL := req.ImageURL
	if s.ossClient != nil {
		thumbnailURL = s.ossClient.GetThumbnailURL(req.ImageURL, 256, 256)
	}

	item := &model.ClothingItem{
		ID:           uuid.New().String(),
		UserID:       userID,
		ImageURL:     req.ImageURL,
		ThumbnailURL: thumbnailURL,
		Category:     req.Category,
		Subcategory:  req.Subcategory,
		Colors:       model.JSONArray(req.Colors),
		Styles:       model.JSONArray(req.Styles),
		Season:       req.Season,
		Tags:         model.JSONArray(req.Tags),
		Notes:        req.Notes,
	}

	if err := s.clothingRepo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

// GetByID 获取衣物详情
func (s *ClothingService) GetByID(ctx context.Context, itemID string) (*model.ClothingItem, error) {
	return s.clothingRepo.GetByID(ctx, itemID)
}

// GetByIDs 批量获取衣物
func (s *ClothingService) GetByIDs(ctx context.Context, itemIDs []string) ([]model.ClothingItem, error) {
	return s.clothingRepo.GetByIDs(ctx, itemIDs)
}

// Update 更新衣物信息
func (s *ClothingService) Update(ctx context.Context, userID, itemID string, req *UpdateClothingRequest) (*model.ClothingItem, error) {
	item, err := s.clothingRepo.GetByID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("衣物不存在")
	}
	if item.UserID != userID {
		return nil, errors.New("无权限修改此衣物")
	}

	// 更新字段
	if req.Category != nil {
		if !isValidCategory(*req.Category) {
			return nil, errors.New("无效的品类")
		}
		item.Category = *req.Category
	}
	if req.Subcategory != nil {
		item.Subcategory = *req.Subcategory
	}
	if req.Colors != nil {
		item.Colors = model.JSONArray(req.Colors)
	}
	if req.Styles != nil {
		item.Styles = model.JSONArray(req.Styles)
	}
	if req.Season != nil {
		if !isValidSeason(*req.Season) {
			return nil, errors.New("无效的季节")
		}
		item.Season = *req.Season
	}
	if req.Tags != nil {
		item.Tags = model.JSONArray(req.Tags)
	}
	if req.Notes != nil {
		item.Notes = *req.Notes
	}

	if err := s.clothingRepo.Update(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

// Delete 删除衣物（软删除）
func (s *ClothingService) Delete(ctx context.Context, userID, itemID string) error {
	item, err := s.clothingRepo.GetByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("衣物不存在")
	}
	if item.UserID != userID {
		return errors.New("无权限删除此衣物")
	}

	return s.clothingRepo.SoftDelete(ctx, itemID)
}

// PermanentDelete 永久删除衣物
func (s *ClothingService) PermanentDelete(ctx context.Context, userID, itemID string) error {
	item, err := s.clothingRepo.GetByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("衣物不存在")
	}
	if item.UserID != userID {
		return errors.New("无权限删除此衣物")
	}

	// 删除OSS文件
	if s.ossClient != nil && item.ImageURL != "" {
		_ = s.ossClient.DeleteByURL(item.ImageURL)
	}

	return s.clothingRepo.HardDelete(ctx, itemID)
}

// Restore 恢复已删除的衣物
func (s *ClothingService) Restore(ctx context.Context, userID, itemID string) error {
	return s.clothingRepo.Restore(ctx, itemID)
}

// SetRetired 设置淘汰状态
func (s *ClothingService) SetRetired(ctx context.Context, userID, itemID string, retired bool) error {
	item, err := s.clothingRepo.GetByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("衣物不存在")
	}
	if item.UserID != userID {
		return errors.New("无权限操作此衣物")
	}

	return s.clothingRepo.SetRetired(ctx, itemID, retired)
}

// BatchSetRetired 批量设置淘汰状态
func (s *ClothingService) BatchSetRetired(ctx context.Context, userID string, itemIDs []string, retired bool) error {
	// 验证所有衣物属于该用户
	items, err := s.clothingRepo.GetByIDs(ctx, itemIDs)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.UserID != userID {
			return errors.New("存在无权限操作的衣物")
		}
	}

	return s.clothingRepo.BatchSetRetired(ctx, itemIDs, retired)
}

// BatchDelete 批量删除（软删除）
func (s *ClothingService) BatchDelete(ctx context.Context, userID string, itemIDs []string) error {
	// 验证所有衣物属于该用户
	items, err := s.clothingRepo.GetByIDs(ctx, itemIDs)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.UserID != userID {
			return errors.New("存在无权限操作的衣物")
		}
	}

	return s.clothingRepo.BatchSoftDelete(ctx, itemIDs)
}

// List 获取衣物列表
func (s *ClothingService) List(ctx context.Context, filter *model.ClothingFilter) ([]model.ClothingItem, int64, error) {
	// 设置默认分页
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	return s.clothingRepo.List(ctx, filter)
}

// GetWardrobe 获取用户可用衣橱（用于推荐）
func (s *ClothingService) GetWardrobe(ctx context.Context, userID string) ([]model.ClothingItem, error) {
	return s.clothingRepo.GetUserWardrobe(ctx, userID)
}

// GetStats 获取衣橱统计
func (s *ClothingService) GetStats(ctx context.Context, userID string) (*model.ClothingStats, error) {
	return s.clothingRepo.GetStats(ctx, userID)
}

// GetCount 获取用户衣物数量
func (s *ClothingService) GetCount(ctx context.Context, userID string) (int64, error) {
	return s.clothingRepo.CountByUser(ctx, userID)
}

// GetByCategory 获取用户某品类的衣物
func (s *ClothingService) GetByCategory(ctx context.Context, userID, category string) ([]model.ClothingItem, error) {
	return s.clothingRepo.GetByUserAndCategory(ctx, userID, category)
}

// 验证品类是否有效
func isValidCategory(category string) bool {
	for _, c := range model.Categories {
		if c == category {
			return true
		}
	}
	return false
}

// 验证季节是否有效
func isValidSeason(season string) bool {
	for _, s := range model.Seasons {
		if s == season {
			return true
		}
	}
	return false
}
