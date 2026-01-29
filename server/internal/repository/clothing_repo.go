package repository

import (
	"context"
	"errors"

	"yidaiku-server/internal/model"

	"gorm.io/gorm"
)

// ClothingRepository 衣物数据访问层
type ClothingRepository struct {
	db *gorm.DB
}

// NewClothingRepository 创建衣物仓库实例
func NewClothingRepository() *ClothingRepository {
	return &ClothingRepository{db: GetDB()}
}

// Create 创建衣物
func (r *ClothingRepository) Create(ctx context.Context, item *model.ClothingItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// GetByID 根据ID获取衣物
func (r *ClothingRepository) GetByID(ctx context.Context, itemID string) (*model.ClothingItem, error) {
	var item model.ClothingItem
	err := r.db.WithContext(ctx).
		Where("item_id = ? AND is_deleted = 0", itemID).
		First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}

// GetByIDs 批量获取衣物
func (r *ClothingRepository) GetByIDs(ctx context.Context, itemIDs []string) ([]model.ClothingItem, error) {
	var items []model.ClothingItem
	err := r.db.WithContext(ctx).
		Where("item_id IN ? AND is_deleted = 0", itemIDs).
		Find(&items).Error
	return items, err
}

// Update 更新衣物
func (r *ClothingRepository) Update(ctx context.Context, item *model.ClothingItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

// UpdateFields 更新指定字段
func (r *ClothingRepository) UpdateFields(ctx context.Context, itemID string, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.ClothingItem{}).
		Where("item_id = ?", itemID).
		Updates(fields).Error
}

// SoftDelete 软删除衣物
func (r *ClothingRepository) SoftDelete(ctx context.Context, itemID string) error {
	return r.db.WithContext(ctx).Model(&model.ClothingItem{}).
		Where("item_id = ?", itemID).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": gorm.Expr("NOW()"),
		}).Error
}

// HardDelete 硬删除衣物
func (r *ClothingRepository) HardDelete(ctx context.Context, itemID string) error {
	return r.db.WithContext(ctx).Delete(&model.ClothingItem{}, "item_id = ?", itemID).Error
}

// Restore 恢复已删除的衣物
func (r *ClothingRepository) Restore(ctx context.Context, itemID string) error {
	return r.db.WithContext(ctx).Model(&model.ClothingItem{}).
		Where("item_id = ?", itemID).
		Updates(map[string]interface{}{
			"is_deleted": false,
			"deleted_at": nil,
		}).Error
}

// SetRetired 设置淘汰状态
func (r *ClothingRepository) SetRetired(ctx context.Context, itemID string, retired bool) error {
	return r.db.WithContext(ctx).Model(&model.ClothingItem{}).
		Where("item_id = ?", itemID).
		Update("is_retired", retired).Error
}

// BatchSetRetired 批量设置淘汰状态
func (r *ClothingRepository) BatchSetRetired(ctx context.Context, itemIDs []string, retired bool) error {
	return r.db.WithContext(ctx).Model(&model.ClothingItem{}).
		Where("item_id IN ?", itemIDs).
		Update("is_retired", retired).Error
}

// BatchSoftDelete 批量软删除
func (r *ClothingRepository) BatchSoftDelete(ctx context.Context, itemIDs []string) error {
	return r.db.WithContext(ctx).Model(&model.ClothingItem{}).
		Where("item_id IN ?", itemIDs).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": gorm.Expr("NOW()"),
		}).Error
}

// List 根据筛选条件获取衣物列表
func (r *ClothingRepository) List(ctx context.Context, filter *model.ClothingFilter) ([]model.ClothingItem, int64, error) {
	var items []model.ClothingItem
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ClothingItem{}).
		Where("user_id = ?", filter.UserID)

	// 默认不显示已删除
	if filter.IsDeleted == nil || !*filter.IsDeleted {
		query = query.Where("is_deleted = 0")
	}

	// 筛选已淘汰状态
	if filter.IsRetired != nil {
		query = query.Where("is_retired = ?", *filter.IsRetired)
	}

	// 品类筛选
	if len(filter.Categories) > 0 {
		query = query.Where("category IN ?", filter.Categories)
	}

	// 季节筛选
	if len(filter.Seasons) > 0 {
		query = query.Where("season IN ?", filter.Seasons)
	}

	// 风格筛选 (JSON数组包含)
	for _, style := range filter.Styles {
		query = query.Where("JSON_CONTAINS(styles, ?)", `"`+style+`"`)
	}

	// 颜色筛选
	for _, color := range filter.Colors {
		query = query.Where("JSON_CONTAINS(colors, ?)", `"`+color+`"`)
	}

	// 标签筛选
	for _, tag := range filter.Tags {
		query = query.Where("JSON_CONTAINS(tags, ?)", `"`+tag+`"`)
	}

	// 关键词搜索
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("notes LIKE ? OR JSON_SEARCH(tags, 'one', ?) IS NOT NULL", keyword, filter.Keyword)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (filter.Page - 1) * filter.PageSize
	if offset < 0 {
		offset = 0
	}

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&items).Error

	return items, total, err
}

// GetUserWardrobe 获取用户可用衣橱（用于推荐）
func (r *ClothingRepository) GetUserWardrobe(ctx context.Context, userID string) ([]model.ClothingItem, error) {
	var items []model.ClothingItem
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_deleted = 0 AND is_retired = 0", userID).
		Find(&items).Error
	return items, err
}

// CountByUser 统计用户衣物数量
func (r *ClothingRepository) CountByUser(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ClothingItem{}).
		Where("user_id = ? AND is_deleted = 0", userID).
		Count(&count).Error
	return count, err
}

// GetStats 获取用户衣橱统计
func (r *ClothingRepository) GetStats(ctx context.Context, userID string) (*model.ClothingStats, error) {
	stats := &model.ClothingStats{
		ByCategory: make(map[string]int),
		BySeason:   make(map[string]int),
	}

	// 总数
	var total int64
	r.db.WithContext(ctx).Model(&model.ClothingItem{}).
		Where("user_id = ? AND is_deleted = 0", userID).
		Count(&total)
	stats.Total = int(total)

	// 已淘汰数量
	var retiredCount int64
	r.db.WithContext(ctx).Model(&model.ClothingItem{}).
		Where("user_id = ? AND is_deleted = 0 AND is_retired = 1", userID).
		Count(&retiredCount)
	stats.RetiredCount = int(retiredCount)

	// 按品类统计
	type CategoryCount struct {
		Category string
		Count    int
	}
	var categoryCounts []CategoryCount
	r.db.WithContext(ctx).Model(&model.ClothingItem{}).
		Select("category, COUNT(*) as count").
		Where("user_id = ? AND is_deleted = 0", userID).
		Group("category").
		Find(&categoryCounts)
	for _, cc := range categoryCounts {
		stats.ByCategory[cc.Category] = cc.Count
	}

	// 按季节统计
	var seasonCounts []CategoryCount
	r.db.WithContext(ctx).Model(&model.ClothingItem{}).
		Select("season as category, COUNT(*) as count").
		Where("user_id = ? AND is_deleted = 0", userID).
		Group("season").
		Find(&seasonCounts)
	for _, sc := range seasonCounts {
		stats.BySeason[sc.Category] = sc.Count
	}

	return stats, nil
}

// GetByUserAndCategory 获取用户某品类的衣物
func (r *ClothingRepository) GetByUserAndCategory(ctx context.Context, userID, category string) ([]model.ClothingItem, error) {
	var items []model.ClothingItem
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND category = ? AND is_deleted = 0 AND is_retired = 0", userID, category).
		Find(&items).Error
	return items, err
}
