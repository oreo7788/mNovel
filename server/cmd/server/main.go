package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"yidaiku-server/internal/config"
	"yidaiku-server/internal/handler"
	"yidaiku-server/internal/repository"
	"yidaiku-server/internal/service"
	"yidaiku-server/pkg/auth"
	"yidaiku-server/pkg/cache"
	"yidaiku-server/pkg/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	if err := repository.InitDB(&cfg.Database); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer repository.CloseDB()

	// 初始化Redis缓存（可选）
	var redisCache *cache.RedisCache
	if cfg.Redis.Host != "" {
		redisCache, err = cache.NewRedisCache(cfg.Redis.Addr(), cfg.Redis.Password, cfg.Redis.DB)
		if err != nil {
			log.Printf("警告: Redis连接失败，将不使用缓存: %v", err)
		} else {
			defer redisCache.Close()
		}
	}

	// 初始化七牛云OSS（可选）
	var ossClient *storage.QiniuOSS
	if cfg.Qiniu.AccessKey != "" && cfg.Qiniu.SecretKey != "" {
		ossClient = storage.NewQiniuOSS(
			cfg.Qiniu.AccessKey,
			cfg.Qiniu.SecretKey,
			cfg.Qiniu.Bucket,
			cfg.Qiniu.Domain,
		)
	}

	// 初始化JWT管理器
	jwtManager := auth.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)

	// 初始化服务
	userService := service.NewUserService(jwtManager)
	clothingService := service.NewClothingService(ossClient)
	outfitService := service.NewOutfitService(nil) // AI适配器暂时不使用
	uploadService := service.NewUploadService(ossClient, &cfg.Upload, repository.NewClothingRepository())

	// 初始化处理器
	userHandler := handler.NewUserHandler(userService)
	clothingHandler := handler.NewClothingHandler(clothingService)
	outfitHandler := handler.NewOutfitHandler(outfitService)
	uploadHandler := handler.NewUploadHandler(uploadService)

	// 创建Gin引擎
	engine := gin.New()

	// 设置路由
	router := handler.NewRouter(
		engine,
		jwtManager,
		userHandler,
		clothingHandler,
		outfitHandler,
		uploadHandler,
	)
	router.Setup()

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 启动服务器（非阻塞）
	go func() {
		log.Printf("服务器启动在端口 %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("服务器关闭异常: %v", err)
	}

	log.Println("服务器已关闭")
}
