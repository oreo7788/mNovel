package handler

import (
	"yidaiku-server/internal/middleware"
	"yidaiku-server/pkg/auth"
	"yidaiku-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// Router 路由管理器
type Router struct {
	engine          *gin.Engine
	jwtManager      *auth.JWTManager
	userHandler     *UserHandler
	clothingHandler *ClothingHandler
	outfitHandler   *OutfitHandler
	uploadHandler   *UploadHandler
}

// NewRouter 创建路由管理器
func NewRouter(
	engine *gin.Engine,
	jwtManager *auth.JWTManager,
	userHandler *UserHandler,
	clothingHandler *ClothingHandler,
	outfitHandler *OutfitHandler,
	uploadHandler *UploadHandler,
) *Router {
	return &Router{
		engine:          engine,
		jwtManager:      jwtManager,
		userHandler:     userHandler,
		clothingHandler: clothingHandler,
		outfitHandler:   outfitHandler,
		uploadHandler:   uploadHandler,
	}
}

// Setup 设置路由
func (r *Router) Setup() {
	// 全局中间件
	r.engine.Use(middleware.CORSMiddleware())
	r.engine.Use(middleware.RecoveryMiddleware())
	r.engine.Use(middleware.RequestLogMiddleware())

	// 健康检查
	r.engine.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"})
	})

	// API v1
	v1 := r.engine.Group("/api/v1")
	{
		// 认证相关（无需登录）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.userHandler.Register)
			auth.POST("/login", r.userHandler.Login)
			auth.POST("/refresh", r.userHandler.RefreshToken)
		}

		// 用户相关（需要登录）
		user := v1.Group("/user")
		user.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			user.GET("/profile", r.userHandler.GetProfile)
			user.PUT("/profile", r.userHandler.UpdateProfile)
			user.PUT("/password", r.userHandler.ChangePassword)
			user.GET("/preference", r.userHandler.GetPreference)
			user.PUT("/preference", r.userHandler.UpdatePreference)
		}

		// 衣橱相关（需要登录）
		wardrobe := v1.Group("/wardrobe")
		wardrobe.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			// 衣物CRUD
			wardrobe.POST("/items", r.clothingHandler.Create)
			wardrobe.GET("/items", r.clothingHandler.List)
			wardrobe.GET("/items/:id", r.clothingHandler.GetByID)
			wardrobe.PUT("/items/:id", r.clothingHandler.Update)
			wardrobe.DELETE("/items/:id", r.clothingHandler.Delete)

			// 衣物状态管理
			wardrobe.PUT("/items/:id/retired", r.clothingHandler.SetRetired)
			wardrobe.POST("/items/:id/restore", r.clothingHandler.Restore)

			// 批量操作
			wardrobe.PUT("/items/batch/retired", r.clothingHandler.BatchSetRetired)
			wardrobe.DELETE("/items/batch", r.clothingHandler.BatchDelete)

			// 统计
			wardrobe.GET("/stats", r.clothingHandler.GetStats)

			// 选项列表（无需登录也可访问）
			wardrobe.GET("/categories", r.clothingHandler.GetCategories)
			wardrobe.GET("/styles", r.clothingHandler.GetStyles)
			wardrobe.GET("/seasons", r.clothingHandler.GetSeasons)
			wardrobe.GET("/colors", r.clothingHandler.GetColors)
		}

		// 穿搭推荐（需要登录）
		outfits := v1.Group("/outfits")
		outfits.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			// 推荐
			outfits.POST("/recommend", r.outfitHandler.Recommend)
			outfits.POST("/replace", r.outfitHandler.ReplaceItem)

			// 收藏
			outfits.GET("/favorites", r.outfitHandler.GetFavorites)
			outfits.PUT("/:id/favorite", r.outfitHandler.SetFavorite)
			outfits.POST("/:id/apply", r.outfitHandler.SetApplied)

			// 模板（无需登录也可访问）
			outfits.GET("/templates", r.outfitHandler.GetTemplates)
			outfits.GET("/occasions", r.outfitHandler.GetOccasions)
		}

		// 上传相关（需要登录）
		upload := v1.Group("/upload")
		upload.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			// 分块上传
			upload.POST("/init", r.uploadHandler.InitUpload)
			upload.POST("/chunk", r.uploadHandler.UploadChunk)
			upload.POST("/complete", r.uploadHandler.CompleteUpload)

			// 简单上传
			upload.POST("/simple", r.uploadHandler.SimpleUpload)

			// 任务管理
			upload.GET("/pending", r.uploadHandler.GetPendingTasks)
			upload.GET("/progress/:task_id", r.uploadHandler.GetProgress)
			upload.DELETE("/:task_id", r.uploadHandler.CancelUpload)
			upload.POST("/:task_id/retry", r.uploadHandler.RetryUpload)
		}
	}

	// 静态资源（如需要）
	// r.engine.Static("/static", "./static")
}
