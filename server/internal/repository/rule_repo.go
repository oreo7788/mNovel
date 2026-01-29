package repository

import (
	"context"
	"errors"

	"yidaiku-server/internal/model"

	"gorm.io/gorm"
)

// MatchingRuleRepository 搭配规则数据访问层
type MatchingRuleRepository struct {
	db *gorm.DB
}

// NewMatchingRuleRepository 创建搭配规则仓库实例
func NewMatchingRuleRepository() *MatchingRuleRepository {
	return &MatchingRuleRepository{db: GetDB()}
}

// Create 创建规则
func (r *MatchingRuleRepository) Create(ctx context.Context, rule *model.MatchingRule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

// GetByID 根据ID获取规则
func (r *MatchingRuleRepository) GetByID(ctx context.Context, ruleID string) (*model.MatchingRule, error) {
	var rule model.MatchingRule
	err := r.db.WithContext(ctx).Where("rule_id = ?", ruleID).First(&rule).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &rule, err
}

// Update 更新规则
func (r *MatchingRuleRepository) Update(ctx context.Context, rule *model.MatchingRule) error {
	return r.db.WithContext(ctx).Save(rule).Error
}

// Delete 删除规则
func (r *MatchingRuleRepository) Delete(ctx context.Context, ruleID string) error {
	return r.db.WithContext(ctx).Delete(&model.MatchingRule{}, "rule_id = ?", ruleID).Error
}

// GetAll 获取所有启用的规则
func (r *MatchingRuleRepository) GetAll(ctx context.Context) ([]model.MatchingRule, error) {
	var rules []model.MatchingRule
	err := r.db.WithContext(ctx).
		Where("is_enabled = 1").
		Order("priority DESC, rule_type ASC").
		Find(&rules).Error
	return rules, err
}

// GetByType 根据类型获取规则
func (r *MatchingRuleRepository) GetByType(ctx context.Context, ruleType string) ([]model.MatchingRule, error) {
	var rules []model.MatchingRule
	err := r.db.WithContext(ctx).
		Where("is_enabled = 1 AND rule_type = ?", ruleType).
		Order("priority DESC").
		Find(&rules).Error
	return rules, err
}

// GetColorRules 获取颜色搭配规则
func (r *MatchingRuleRepository) GetColorRules(ctx context.Context) ([]model.MatchingRule, error) {
	return r.GetByType(ctx, model.RuleTypeColor)
}

// GetStyleRules 获取风格搭配规则
func (r *MatchingRuleRepository) GetStyleRules(ctx context.Context) ([]model.MatchingRule, error) {
	return r.GetByType(ctx, model.RuleTypeStyle)
}

// GetOccasionRules 获取场合适配规则
func (r *MatchingRuleRepository) GetOccasionRules(ctx context.Context) ([]model.MatchingRule, error) {
	return r.GetByType(ctx, model.RuleTypeOccasion)
}

// GetSeasonRules 获取季节适配规则
func (r *MatchingRuleRepository) GetSeasonRules(ctx context.Context) ([]model.MatchingRule, error) {
	return r.GetByType(ctx, model.RuleTypeSeason)
}

// SetEnabled 设置规则启用状态
func (r *MatchingRuleRepository) SetEnabled(ctx context.Context, ruleID string, enabled bool) error {
	return r.db.WithContext(ctx).Model(&model.MatchingRule{}).
		Where("rule_id = ?", ruleID).
		Update("is_enabled", enabled).Error
}

// UpdateWeight 更新规则权重
func (r *MatchingRuleRepository) UpdateWeight(ctx context.Context, ruleID string, weight float64) error {
	return r.db.WithContext(ctx).Model(&model.MatchingRule{}).
		Where("rule_id = ?", ruleID).
		Update("weight", weight).Error
}
