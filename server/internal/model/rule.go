package model

import (
	"time"
)

// MatchingRule 搭配规则
type MatchingRule struct {
	ID        string    `json:"id" gorm:"column:rule_id;primaryKey;type:varchar(64)"`
	RuleType  string    `json:"rule_type" gorm:"column:rule_type;type:varchar(32);index"` // color/style/occasion/season
	Name      string    `json:"name" gorm:"column:name;type:varchar(100)"`
	Condition JSONMap   `json:"condition" gorm:"column:condition;type:json"` // 条件
	Action    JSONMap   `json:"action" gorm:"column:action;type:json"`       // 动作
	Weight    float64   `json:"weight" gorm:"column:weight;type:decimal(3,2);default:1.0"`
	Priority  int       `json:"priority" gorm:"column:priority;type:int;default:0"`
	IsEnabled bool      `json:"is_enabled" gorm:"column:is_enabled;type:tinyint;default:1"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 表名
func (MatchingRule) TableName() string {
	return "matching_rules"
}

// RuleType 规则类型常量
const (
	RuleTypeColor    = "color"
	RuleTypeStyle    = "style"
	RuleTypeOccasion = "occasion"
	RuleTypeSeason   = "season"
)

// ColorMatchRule 颜色搭配规则
type ColorMatchRule struct {
	TopColor    string   `json:"top_color"`
	BottomColors []string `json:"bottom_colors"`
	Score       float64  `json:"score"`
}

// OccasionRule 场合规则
type OccasionRule struct {
	Occasion    string   `json:"occasion"`
	Styles      []string `json:"styles"`        // 适合的风格
	Colors      []string `json:"colors"`        // 适合的颜色
	AvoidItems  []string `json:"avoid_items"`   // 避免的单品
	PreferItems []string `json:"prefer_items"`  // 偏好的单品
}

// 预定义的场合列表
var Occasions = []string{
	"职场通勤",
	"约会",
	"运动",
	"聚会",
	"面试",
	"旅行",
	"周末逛街",
}

// 预定义的风格列表
var Styles = []string{
	"休闲",
	"通勤",
	"运动",
	"正式",
	"甜美",
	"中性",
	"复古",
}

// 预定义的季节列表
var Seasons = []string{
	"春夏",
	"秋冬",
	"四季通用",
}

// 预定义的品类列表
var Categories = []string{
	"上衣",
	"裤子",
	"裙子",
	"鞋子",
	"配饰",
}

// 预定义的颜色列表
var Colors = []string{
	"黑色",
	"白色",
	"灰色",
	"米色",
	"深蓝色",
	"浅蓝色",
	"红色",
	"粉色",
	"绿色",
	"黄色",
	"棕色",
	"卡其色",
	"紫色",
	"橙色",
}
