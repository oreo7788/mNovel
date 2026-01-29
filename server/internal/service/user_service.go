package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"regexp"
	"strings"

	"yidaiku-server/internal/model"
	"yidaiku-server/internal/repository"
	"yidaiku-server/pkg/auth"

	"golang.org/x/crypto/bcrypt"
)

// 注册相关错误，便于 Handler 区分 HTTP 状态码
var (
	ErrPhoneAlreadyRegistered = errors.New("该手机号已注册")
	ErrEmailAlreadyRegistered = errors.New("该邮箱已注册")
	ErrInvalidPhoneFormat     = errors.New("手机号格式不正确，应为11位数字")
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

// RegisterRequest 注册请求（手机号可选；支持昵称+邮箱注册）
type RegisterRequest struct {
	Phone    string `json:"phone"` // 可选
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"` // 可选，前端传用户名
	Email    string `json:"email"`    // 可选
}

// LoginRequest 登录请求（手机号、邮箱、昵称三选一，与密码一起使用）
type LoginRequest struct {
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	User         *model.User `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int64       `json:"expires_in"` // 秒
}

// 中国大陆手机号：1 开头，共 11 位数字
var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

// userID 字符集：大小写字母 + 数字
const userIDCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const userIDLength = 10
const userIDMaxRetries = 10

// generateUniqueUserID 生成由大小写字母和数字组成的唯一 user_id（长度 10）
func (s *UserService) generateUniqueUserID(ctx context.Context) (string, error) {
	for i := 0; i < userIDMaxRetries; i++ {
		id, err := randomAlphanumeric(userIDLength)
		if err != nil {
			return "", err
		}
		existing, err := s.userRepo.GetByID(ctx, id)
		if err != nil {
			return "", err
		}
		if existing == nil {
			return id, nil
		}
	}
	return "", errors.New("生成唯一 user_id 失败，请重试")
}

// randomAlphanumeric 生成长度为 n 的随机字符串（字符集：a-zA-Z0-9）
func randomAlphanumeric(n int) (string, error) {
	b := make([]byte, n)
	charsetLen := big.NewInt(int64(len(userIDCharset)))
	for i := 0; i < n; i++ {
		idx, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		b[i] = userIDCharset[idx.Int64()]
	}
	return string(b), nil
}

// Register 用户注册（不要求绑定手机号；支持昵称+邮箱）
func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// 若填写了手机号，则校验格式并检查是否已注册
	if req.Phone != "" {
		if !phoneRegex.MatchString(req.Phone) {
			return nil, ErrInvalidPhoneFormat
		}
		existingUser, err := s.userRepo.GetByPhone(ctx, req.Phone)
		if err != nil {
			return nil, err
		}
		if existingUser != nil {
			return nil, ErrPhoneAlreadyRegistered
		}
	}
	// 若填写了邮箱，检查是否已注册
	if req.Email != "" {
		existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		if existingUser != nil {
			return nil, ErrEmailAlreadyRegistered
		}
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userID, err := s.generateUniqueUserID(ctx)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		UserID:       userID,
		PasswordHash: string(hashedPassword),
		Nickname:     req.Nickname,
		Email:        req.Email,
		Status:       1,
	}

	// 手机号：有则写入，无则保持 null（不绑定）
	if req.Phone != "" {
		user.Phone = &req.Phone
	}

	// 默认昵称：有手机用后4位，有邮箱用邮箱前缀，否则用 user_id 后4位
	if user.Nickname == "" {
		if req.Phone != "" {
			user.Nickname = "用户" + req.Phone[len(req.Phone)-4:]
		} else if req.Email != "" {
			at := strings.Index(req.Email, "@")
			if at > 0 && at <= 20 {
				user.Nickname = req.Email[:at]
			} else {
				user.Nickname = "用户" + userID[len(userID)-4:]
			}
		} else {
			if len(userID) >= 4 {
				user.Nickname = "用户" + userID[len(userID)-4:]
			} else {
				user.Nickname = "用户"
			}
		}
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// 创建默认偏好设置
	pref := &model.UserPreference{
		UserID: user.UserID,
	}
	s.prefRepo.CreateOrUpdate(ctx, pref)

	// 生成Token
	return s.generateAuthResponse(user)
}

// Login 用户登录（支持手机号、邮箱或昵称）
func (s *UserService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	var user *model.User
	var err error
	switch {
	case strings.TrimSpace(req.Phone) != "":
		user, err = s.userRepo.GetByPhone(ctx, strings.TrimSpace(req.Phone))
	case strings.TrimSpace(req.Email) != "":
		user, err = s.userRepo.GetByEmail(ctx, strings.TrimSpace(req.Email))
	case strings.TrimSpace(req.Nickname) != "":
		user, err = s.userRepo.GetByNickname(ctx, strings.TrimSpace(req.Nickname))
	default:
		return nil, errors.New("请填写手机号、邮箱或用户名")
	}
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("用户不存在或密码错误")
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
	accessToken, err := s.jwtManager.GenerateAccessToken(user.UserID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.UserID)
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
