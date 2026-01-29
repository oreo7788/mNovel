package repository

import (
	"context"
	"errors"
	"time"

	"yidaiku-server/internal/model"

	"gorm.io/gorm"
)

// OutfitRepository 穿搭推荐数据访问层
type OutfitRepository struct {
	db *gorm.DB
}

// NewOutfitRepository 创建穿搭仓库实例
func NewOutfitRepository() *OutfitRepository {
	return &OutfitRepository{db: GetDB()}
}

// Create 创建推荐
func (r *OutfitRepository) Create(ctx context.Context, outfit *model.OutfitRecommendation) error {
	return r.db.WithContext(ctx).Create(outfit).Error
}

// BatchCreate 批量创建推荐
func (r *OutfitRepository) BatchCreate(ctx context.Context, outfits []model.OutfitRecommendation) error {
	return r.db.WithContext(ctx).Create(&outfits).Error
}

// GetByID 根据ID获取推荐
func (r *OutfitRepository) GetByID(ctx context.Context, outfitID string) (*model.OutfitRecommendation, error) {
	var outfit model.OutfitRecommendation
	err := r.db.WithContext(ctx).Where("outfit_id = ?", outfitID).First(&outfit).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &outfit, err
}

// Update 更新推荐
func (r *OutfitRepository) Update(ctx context.Context, outfit *model.OutfitRecommendation) error {
	return r.db.WithContext(ctx).Save(outfit).Error
}

// UpdateFields 更新指定字段
func (r *OutfitRepository) UpdateFields(ctx context.Context, outfitID string, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.OutfitRecommendation{}).
		Where("outfit_id = ?", outfitID).
		Updates(fields).Error
}

// GetUserRecommendations 获取用户的推荐列表
func (r *OutfitRepository) GetUserRecommendations(ctx context.Context, userID string, page, pageSize int) ([]model.OutfitRecommendation, int64, error) {
	var outfits []model.OutfitRecommendation
	var total int64

	query := r.db.WithContext(ctx).Model(&model.OutfitRecommendation{}).
		Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&outfits).Error

	return outfits, total, err
}

// GetRecentRecommendations 获取最近N天的推荐（用于去重）
func (r *OutfitRepository) GetRecentRecommendations(ctx context.Context, userID string, days int) ([]model.OutfitRecommendation, error) {
	var outfits []model.OutfitRecommendation
	startDate := time.Now().AddDate(0, 0, -days)
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND created_at >= ?", userID, startDate).
		Find(&outfits).Error
	return outfits, err
}

// SetFavorite 设置收藏状态
func (r *OutfitRepository) SetFavorite(ctx context.Context, outfitID string, favorite bool) error {
	return r.db.WithContext(ctx).Model(&model.OutfitRecommendation{}).
		Where("outfit_id = ?", outfitID).
		Update("is_favorite", favorite).Error
}

// SetApplied 设置应用状态
func (r *OutfitRepository) SetApplied(ctx context.Context, outfitID string, applied bool) error {
	return r.db.WithContext(ctx).Model(&model.OutfitRecommendation{}).
		Where("outfit_id = ?", outfitID).
		Update("is_applied", applied).Error
}

// OutfitTemplateRepository 穿搭模板数据访问层
type OutfitTemplateRepository struct {
	db *gorm.DB
}

// NewOutfitTemplateRepository 创建穿搭模板仓库实例
func NewOutfitTemplateRepository() *OutfitTemplateRepository {
	return &OutfitTemplateRepository{db: GetDB()}
}

// GetAll 获取所有启用的模板
func (r *OutfitTemplateRepository) GetAll(ctx context.Context) ([]model.OutfitTemplate, error) {
	var templates []model.OutfitTemplate
	err := r.db.WithContext(ctx).
		Where("is_enabled = 1").
		Order("sort_order ASC").
		Find(&templates).Error
	return templates, err
}

// GetByOccasion 根据场合获取模板
func (r *OutfitTemplateRepository) GetByOccasion(ctx context.Context, occasion string) ([]model.OutfitTemplate, error) {
	var templates []model.OutfitTemplate
	err := r.db.WithContext(ctx).
		Where("is_enabled = 1 AND occasion = ?", occasion).
		Order("sort_order ASC").
		Find(&templates).Error
	return templates, err
}

// GetByStyle 根据风格获取模板
func (r *OutfitTemplateRepository) GetByStyle(ctx context.Context, style string) ([]model.OutfitTemplate, error) {
	var templates []model.OutfitTemplate
	err := r.db.WithContext(ctx).
		Where("is_enabled = 1 AND style = ?", style).
		Order("sort_order ASC").
		Find(&templates).Error
	return templates, err
}

// FavoriteOutfitRepository 收藏穿搭数据访问层
type FavoriteOutfitRepository struct {
	db *gorm.DB
}

// NewFavoriteOutfitRepository 创建收藏仓库实例
func NewFavoriteOutfitRepository() *FavoriteOutfitRepository {
	return &FavoriteOutfitRepository{db: GetDB()}
}

// Create 创建收藏
func (r *FavoriteOutfitRepository) Create(ctx context.Context, favorite *model.FavoriteOutfit) error {
	return r.db.WithContext(ctx).Create(favorite).Error
}

// Delete 删除收藏
func (r *FavoriteOutfitRepository) Delete(ctx context.Context, userID, outfitID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND outfit_id = ?", userID, outfitID).
		Delete(&model.FavoriteOutfit{}).Error
}

// GetByUser 获取用户的收藏列表
func (r *FavoriteOutfitRepository) GetByUser(ctx context.Context, userID string, page, pageSize int) ([]model.FavoriteOutfit, int64, error) {
	var favorites []model.FavoriteOutfit
	var total int64

	query := r.db.WithContext(ctx).Model(&model.FavoriteOutfit{}).
		Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&favorites).Error

	return favorites, total, err
}

// Exists 检查是否已收藏
func (r *FavoriteOutfitRepository) Exists(ctx context.Context, userID, outfitID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.FavoriteOutfit{}).
		Where("user_id = ? AND outfit_id = ?", userID, outfitID).
		Count(&count).Error
	return count > 0, err
}

// OutfitRecordRepository 穿搭记录数据访问层
type OutfitRecordRepository struct {
	db *gorm.DB
}

// NewOutfitRecordRepository 创建穿搭记录仓库实例
func NewOutfitRecordRepository() *OutfitRecordRepository {
	return &OutfitRecordRepository{db: GetDB()}
}

// Create 创建记录
func (r *OutfitRecordRepository) Create(ctx context.Context, record *model.OutfitRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

// GetByID 根据ID获取记录
func (r *OutfitRecordRepository) GetByID(ctx context.Context, recordID string) (*model.OutfitRecord, error) {
	var record model.OutfitRecord
	err := r.db.WithContext(ctx).Where("record_id = ?", recordID).First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &record, err
}

// Update 更新记录
func (r *OutfitRecordRepository) Update(ctx context.Context, record *model.OutfitRecord) error {
	return r.db.WithContext(ctx).Save(record).Error
}

// Delete 删除记录
func (r *OutfitRecordRepository) Delete(ctx context.Context, recordID string) error {
	return r.db.WithContext(ctx).Delete(&model.OutfitRecord{}, "record_id = ?", recordID).Error
}

// GetByUserAndDate 获取用户某日期的记录
func (r *OutfitRecordRepository) GetByUserAndDate(ctx context.Context, userID string, date time.Time) (*model.OutfitRecord, error) {
	var record model.OutfitRecord
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND date = ?", userID, date.Format("2006-01-02")).
		First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &record, err
}

// GetByUser 获取用户的穿搭记录
func (r *OutfitRecordRepository) GetByUser(ctx context.Context, userID string, page, pageSize int) ([]model.OutfitRecord, int64, error) {
	var records []model.OutfitRecord
	var total int64

	query := r.db.WithContext(ctx).Model(&model.OutfitRecord{}).
		Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("date DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error

	return records, total, err
}

// GetByUserAndDateRange 获取用户某时间段的记录
func (r *OutfitRecordRepository) GetByUserAndDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]model.OutfitRecord, error) {
	var records []model.OutfitRecord
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("date DESC").
		Find(&records).Error
	return records, err
}
