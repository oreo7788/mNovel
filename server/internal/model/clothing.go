package model

import (
	"time"
)

// ClothingItem 衣物模型
type ClothingItem struct {
	ID           string    `json:"id" gorm:"column:item_id;primaryKey;type:varchar(64)"`
	UserID       string    `json:"user_id" gorm:"column:user_id;type:varchar(64);index"`
	ImageURL     string    `json:"image_url" gorm:"column:image_url;type:varchar(512)"`
	ThumbnailURL string    `json:"thumbnail_url" gorm:"column:thumbnail_url;type:varchar(512)"`
	Category     string    `json:"category" gorm:"column:category;type:varchar(32);index"` // 上衣、裤子、裙子、鞋子、配饰
	Subcategory  string    `json:"subcategory" gorm:"column:subcategory;type:varchar(32)"` // 子品类
	Colors       JSONArray `json:"colors" gorm:"column:colors;type:json"`                  // 颜色数组
	Styles       JSONArray `json:"styles" gorm:"column:styles;type:json"`                  // 风格数组
	Season       string    `json:"season" gorm:"column:season;type:varchar(16)"`           // 春夏、秋冬、四季通用
	Tags         JSONArray `json:"tags" gorm:"column:tags;type:json"`                      // 自定义标签
	Notes        string    `json:"notes" gorm:"column:notes;type:text"`                    // 备注
	IsDeleted    bool      `json:"is_deleted" gorm:"column:is_deleted;type:tinyint;default:0;index"`
	IsRetired    bool      `json:"is_retired" gorm:"column:is_retired;type:tinyint;default:0;index"` // 已淘汰
	DeletedAt    *time.Time `json:"deleted_at" gorm:"column:deleted_at"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 表名
func (ClothingItem) TableName() string {
	return "clothing_items"
}

// ClothingColor 衣物颜色信息
type ClothingColor struct {
	Name       string  `json:"name"`
	RGB        []int   `json:"rgb,omitempty"`
	Proportion float64 `json:"proportion,omitempty"` // 颜色占比
}

// ClothingFilter 衣物筛选条件
type ClothingFilter struct {
	UserID     string   `json:"user_id"`
	Categories []string `json:"categories,omitempty"`
	Styles     []string `json:"styles,omitempty"`
	Seasons    []string `json:"seasons,omitempty"`
	Colors     []string `json:"colors,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	IsRetired  *bool    `json:"is_retired,omitempty"`
	IsDeleted  *bool    `json:"is_deleted,omitempty"`
	Keyword    string   `json:"keyword,omitempty"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
}

// ClothingStats 衣物统计
type ClothingStats struct {
	Total         int            `json:"total"`
	ByCategory    map[string]int `json:"by_category"`
	BySeason      map[string]int `json:"by_season"`
	RetiredCount  int            `json:"retired_count"`
}
