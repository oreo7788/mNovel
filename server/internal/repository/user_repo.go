package repository

import (
	"context"
	"errors"

	"yidaiku-server/internal/model"

	"gorm.io/gorm"
)

// UserRepository 用户数据访问层
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository() *UserRepository {
	return &UserRepository{db: GetDB()}
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// GetByPhone 根据手机号获取用户（空字符串直接返回 nil，不查库）
func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	if phone == "" {
		return nil, nil
	}
	var user model.User
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if email == "" {
		return nil, nil
	}
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// GetByNickname 根据昵称获取用户
func (r *UserRepository) GetByNickname(ctx context.Context, nickname string) (*model.User, error) {
	if nickname == "" {
		return nil, nil
	}
	var user model.User
	err := r.db.WithContext(ctx).Where("nickname = ?", nickname).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// Update 更新用户信息
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// UpdateFields 更新指定字段
func (r *UserRepository) UpdateFields(ctx context.Context, userID string, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("user_id = ?", userID).Updates(fields).Error
}

// Delete 删除用户（软删除）
func (r *UserRepository) Delete(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("user_id = ?", userID).Update("status", 0).Error
}

// UserPreferenceRepository 用户偏好数据访问层
type UserPreferenceRepository struct {
	db *gorm.DB
}

// NewUserPreferenceRepository 创建用户偏好仓库实例
func NewUserPreferenceRepository() *UserPreferenceRepository {
	return &UserPreferenceRepository{db: GetDB()}
}

// GetByUserID 根据用户ID获取偏好设置
func (r *UserPreferenceRepository) GetByUserID(ctx context.Context, userID string) (*model.UserPreference, error) {
	var pref model.UserPreference
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&pref).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &pref, err
}

// CreateOrUpdate 创建或更新用户偏好
func (r *UserPreferenceRepository) CreateOrUpdate(ctx context.Context, pref *model.UserPreference) error {
	// 使用 Upsert
	return r.db.WithContext(ctx).Save(pref).Error
}

// UserBehaviorRepository 用户行为数据访问层
type UserBehaviorRepository struct {
	db *gorm.DB
}

// NewUserBehaviorRepository 创建用户行为仓库实例
func NewUserBehaviorRepository() *UserBehaviorRepository {
	return &UserBehaviorRepository{db: GetDB()}
}

// Create 创建行为记录
func (r *UserBehaviorRepository) Create(ctx context.Context, behavior *model.UserBehavior) error {
	return r.db.WithContext(ctx).Create(behavior).Error
}

// GetByUserID 获取用户行为记录
func (r *UserBehaviorRepository) GetByUserID(ctx context.Context, userID string, limit int) ([]model.UserBehavior, error) {
	var behaviors []model.UserBehavior
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&behaviors).Error
	return behaviors, err
}

// CountByType 统计用户某类型行为数量
func (r *UserBehaviorRepository) CountByType(ctx context.Context, userID, behaviorType string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.UserBehavior{}).
		Where("user_id = ? AND behavior_type = ?", userID, behaviorType).
		Count(&count).Error
	return count, err
}
