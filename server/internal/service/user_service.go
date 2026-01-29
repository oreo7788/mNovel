package service

import (
	"context"
	"errors"
	"time"

	"yidaiku-server/internal/model"
	"yidaiku-server/internal/repository"
	"yidaiku-server/pkg/auth"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务
type UserService struct {
	userRepo     *repository.UserRepository
	prefRepo     *repository.UserPreferenceRepository
	behaviorRepo *repository.UserBehaviorRepository
	jwtManager   *auth.JWTManager
}

// NewUserService 创建用户服务实例
func NewUserService(jwtManager *auth.JWTManager) *UserService {
	return &UserService{
		userRepo:     repository.NewUserRepository(),
		prefRepo:     repository.NewUserPreferenceRepository(),
		behaviorRepo: repository.NewUserBehaviorRepository(),
		jwtManager:   jwtManager,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	User         *model.User `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int64       `json:"expires_in"` // 秒
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// 检查手机号是否已存在
	existingUser, err := s.userRepo.GetByPhone(ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("该手机号已注册")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &model.User{
		ID:           uuid.New().String(),
		Phone:        req.Phone,
		PasswordHash: string(hashedPassword),
		Nickname:     req.Nickname,
		Status:       1,
		CreatedAt:    time.Now(),
	}

	if user.Nickname == "" {
		user.Nickname = "用户" + req.Phone[len(req.Phone)-4:]
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// 创建默认偏好设置
	pref := &model.UserPreference{
		UserID: user.ID,
	}
	s.prefRepo.CreateOrUpdate(ctx, pref)

	// 生成Token
	return s.generateAuthResponse(user)
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.GetByPhone(ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("密码错误")
	}

	// 检查账户状态
	if user.Status != 1 {
		return nil, errors.New("账户已被禁用")
	}

	return s.generateAuthResponse(user)
}

// RefreshToken 刷新Token
func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return nil, errors.New("用户不存在")
	}

	return s.generateAuthResponse(user)
}

// GetByID 获取用户信息
func (s *UserService) GetByID(ctx context.Context, userID string) (*model.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(ctx context.Context, userID string, updates map[string]interface{}) error {
	// 只允许更新特定字段
	allowedFields := map[string]bool{"nickname": true, "avatar": true}
	filtered := make(map[string]interface{})
	for k, v := range updates {
		if allowedFields[k] {
			filtered[k] = v
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	return s.userRepo.UpdateFields(ctx, userID, filtered)
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("原密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdateFields(ctx, userID, map[string]interface{}{
		"password_hash": string(hashedPassword),
	})
}

// GetPreference 获取用户偏好设置
func (s *UserService) GetPreference(ctx context.Context, userID string) (*model.UserPreference, error) {
	pref, err := s.prefRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if pref == nil {
		// 返回空的偏好设置
		return &model.UserPreference{UserID: userID}, nil
	}
	return pref, nil
}

// UpdatePreference 更新用户偏好设置
func (s *UserService) UpdatePreference(ctx context.Context, pref *model.UserPreference) error {
	return s.prefRepo.CreateOrUpdate(ctx, pref)
}

// RecordBehavior 记录用户行为
func (s *UserService) RecordBehavior(ctx context.Context, userID, outfitID, behaviorType string, metadata model.JSONMap) error {
	behavior := &model.UserBehavior{
		UserID:       userID,
		OutfitID:     outfitID,
		BehaviorType: behaviorType,
		Metadata:     metadata,
	}
	return s.behaviorRepo.Create(ctx, behavior)
}

// GetBehaviorCount 获取用户行为数量
func (s *UserService) GetBehaviorCount(ctx context.Context, userID string) (int64, error) {
	return s.behaviorRepo.CountByType(ctx, userID, "")
}

// generateAuthResponse 生成认证响应
func (s *UserService) generateAuthResponse(user *model.User) (*AuthResponse, error) {
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtManager.AccessTokenTTL().Seconds()),
	}, nil
}
