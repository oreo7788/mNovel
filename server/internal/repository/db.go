package repository

import (
	"fmt"
	"time"

	"yidaiku-server/internal/config"
	"yidaiku-server/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// 北京时区（东八区），全项目 created_at/updated_at 统一使用
var beijingLoc = mustLoadLocation("Asia/Shanghai")

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic("加载时区失败: " + err.Error())
	}
	return loc
}

// BeijingNow 返回当前北京时间，用于业务层需要显式设置时间时与 GORM NowFunc 一致
func BeijingNow() time.Time {
	return time.Now().In(beijingLoc)
}

// InitDB 初始化数据库连接
func InitDB(cfg *config.DatabaseConfig) error {
	var err error

	// 配置GORM：created_at/updated_at 一律使用北京时间
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().In(beijingLoc)
		},
	}

	// 连接数据库
	db, err = gorm.Open(mysql.Open(cfg.DSN()), gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层*sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return db
}

// autoMigrate 自动迁移数据库表结构
func autoMigrate() error {
	return db.AutoMigrate(
		&model.User{},
		&model.UserPreference{},
		&model.UserBehavior{},
		&model.ClothingItem{},
		&model.OutfitRecommendation{},
		&model.OutfitTemplate{},
		&model.FavoriteOutfit{},
		&model.OutfitRecord{},
		&model.MatchingRule{},
		&model.UploadTask{},
	)
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
