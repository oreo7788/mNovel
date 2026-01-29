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

// InitDB 初始化数据库连接
func InitDB(cfg *config.DatabaseConfig) error {
	var err error
	
	// 配置GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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
