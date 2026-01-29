package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID           string     `json:"id" gorm:"column:user_id;primaryKey;type:varchar(64)"`
	Phone        string     `json:"phone" gorm:"column:phone;type:varchar(20);uniqueIndex"`
	Email        string     `json:"email" gorm:"column:email;type:varchar(100);index"`
	PasswordHash string     `json:"-" gorm:"column:password_hash;type:varchar(255)"`
	Nickname     string     `json:"nickname" gorm:"column:nickname;type:varchar(50)"`
	Avatar       string     `json:"avatar" gorm:"column:avatar;type:varchar(512)"`
	Status       int        `json:"status" gorm:"column:status;type:tinyint;default:1"` // 1:正常 0:禁用
	IsMember     bool       `json:"is_member" gorm:"column:is_member;type:tinyint;default:0"`
	MemberExpire *time.Time `json:"member_expire" gorm:"column:member_expire"`
	CreatedAt    time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}

// UserPreference 用户偏好设置
type UserPreference struct {
	UserID    string    `json:"user_id" gorm:"column:user_id;primaryKey;type:varchar(64)"`
	Styles    JSONArray `json:"styles" gorm:"column:styles;type:json"`       // 风格偏好数组
	Occasions JSONArray `json:"occasions" gorm:"column:occasions;type:json"` // 场合偏好数组
	Colors    JSONArray `json:"colors" gorm:"column:colors;type:json"`       // 颜色偏好数组
	Height    int       `json:"height" gorm:"column:height;type:int"`        // 身高(cm)
	Weight    int       `json:"weight" gorm:"column:weight;type:int"`        // 体重(kg)
	Principles JSONArray `json:"principles" gorm:"column:principles;type:json"` // 穿搭原则
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 表名
func (UserPreference) TableName() string {
	return "user_preferences"
}

// UserBehavior 用户行为记录
type UserBehavior struct {
	ID           uint      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID       string    `json:"user_id" gorm:"column:user_id;type:varchar(64);index"`
	OutfitID     string    `json:"outfit_id" gorm:"column:outfit_id;type:varchar(64);index"`
	BehaviorType string    `json:"behavior_type" gorm:"column:behavior_type;type:varchar(32)"` // favorite/apply/dislike/replace
	Metadata     JSONMap   `json:"metadata" gorm:"column:metadata;type:json"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

// TableName 表名
func (UserBehavior) TableName() string {
	return "user_behaviors"
}
