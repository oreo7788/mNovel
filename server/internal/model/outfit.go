package model

import (
	"time"
)

// OutfitRecommendation 穿搭推荐
type OutfitRecommendation struct {
	ID         string    `json:"id" gorm:"column:outfit_id;primaryKey;type:varchar(64)"`
	UserID     string    `json:"user_id" gorm:"column:user_id;type:varchar(20);index"`
	ItemIDs    JSONArray `json:"item_ids" gorm:"column:item_ids;type:json"` // 单品ID数组
	Occasion   string    `json:"occasion" gorm:"column:occasion;type:varchar(32)"`
	Weather    JSONMap   `json:"weather" gorm:"column:weather;type:json"` // 天气信息
	Score      float64   `json:"score" gorm:"column:score;type:decimal(3,2)"`
	Highlights string    `json:"highlights" gorm:"column:highlights;type:text"` // 搭配亮点
	Suitable   JSONArray `json:"suitable" gorm:"column:suitable;type:json"`     // 适用场景
	Source     string    `json:"source" gorm:"column:source;type:varchar(32)"`  // 推荐来源: ai/rule/template
	IsFavorite bool      `json:"is_favorite" gorm:"column:is_favorite;type:tinyint;default:0"`
	IsApplied  bool      `json:"is_applied" gorm:"column:is_applied;type:tinyint;default:0"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime;index"`
}

// TableName 表名
func (OutfitRecommendation) TableName() string {
	return "outfit_recommendations"
}

// OutfitTemplate 穿搭模板（用于冷启动）
type OutfitTemplate struct {
	ID          string    `json:"id" gorm:"column:template_id;primaryKey;type:varchar(64)"`
	Name        string    `json:"name" gorm:"column:name;type:varchar(100)"`
	Description string    `json:"description" gorm:"column:description;type:text"`
	Items       JSONArray `json:"items" gorm:"column:items;type:json"` // 模板单品描述
	Occasion    string    `json:"occasion" gorm:"column:occasion;type:varchar(32)"`
	Style       string    `json:"style" gorm:"column:style;type:varchar(32)"`
	Season      string    `json:"season" gorm:"column:season;type:varchar(16)"`
	ImageURL    string    `json:"image_url" gorm:"column:image_url;type:varchar(512)"`
	SortOrder   int       `json:"sort_order" gorm:"column:sort_order;type:int;default:0"`
	IsEnabled   bool      `json:"is_enabled" gorm:"column:is_enabled;type:tinyint;default:1"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 表名
func (OutfitTemplate) TableName() string {
	return "outfit_templates"
}

// FavoriteOutfit 收藏的穿搭
type FavoriteOutfit struct {
	ID        uint      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID    string    `json:"user_id" gorm:"column:user_id;type:varchar(20);index"`
	OutfitID  string    `json:"outfit_id" gorm:"column:outfit_id;type:varchar(64);index"`
	Notes     string    `json:"notes" gorm:"column:notes;type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

// TableName 表名
func (FavoriteOutfit) TableName() string {
	return "favorite_outfits"
}

// OutfitRecord 穿搭记录/日记
type OutfitRecord struct {
	ID        string    `json:"id" gorm:"column:record_id;primaryKey;type:varchar(64)"`
	UserID    string    `json:"user_id" gorm:"column:user_id;type:varchar(20);index"`
	Date      time.Time `json:"date" gorm:"column:date;type:date;index"`
	ImageURLs JSONArray `json:"image_urls" gorm:"column:image_urls;type:json"`      // 实际穿搭照片
	ItemIDs   JSONArray `json:"item_ids" gorm:"column:item_ids;type:json"`          // 搭配的单品
	OutfitID  string    `json:"outfit_id" gorm:"column:outfit_id;type:varchar(64)"` // 关联的推荐（如有）
	Notes     string    `json:"notes" gorm:"column:notes;type:text"`
	Tags      JSONArray `json:"tags" gorm:"column:tags;type:json"`
	Mood      string    `json:"mood" gorm:"column:mood;type:varchar(32)"` // 心情
	Weather   JSONMap   `json:"weather" gorm:"column:weather;type:json"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 表名
func (OutfitRecord) TableName() string {
	return "outfit_records"
}

// RecommendRequest 推荐请求参数
type RecommendRequest struct {
	Occasion    string          `json:"occasion" binding:"required"`
	Weather     *WeatherInfo    `json:"weather"`
	Preferences *UserPreference `json:"preferences,omitempty"`
}

// WeatherInfo 天气信息
type WeatherInfo struct {
	Temperature int    `json:"temperature"` // 温度(°C)
	Condition   string `json:"condition"`   // 天气状况: 晴天/多云/雨天/雪天
	Humidity    int    `json:"humidity"`    // 湿度(%)
	WindLevel   int    `json:"wind_level"`  // 风力等级
}

// ReplaceItemRequest 换一件请求
type ReplaceItemRequest struct {
	OutfitID  string `json:"outfit_id" binding:"required"`
	OldItemID string `json:"old_item_id" binding:"required"`
	NewItemID string `json:"new_item_id" binding:"required"`
}
