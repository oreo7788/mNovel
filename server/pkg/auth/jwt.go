package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTManager JWT管理器
type JWTManager struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// Claims JWT声明
type Claims struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"` // access or refresh
	jwt.RegisteredClaims
}

// NewJWTManager 创建JWT管理器
func NewJWTManager(secret string, accessTTL, refreshTTL time.Duration) *JWTManager {
	return &JWTManager{
		secret:          []byte(secret),
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

// GenerateAccessToken 生成访问令牌
func (m *JWTManager) GenerateAccessToken(userID string) (string, error) {
	return m.generateToken(userID, "access", m.accessTokenTTL)
}

// GenerateRefreshToken 生成刷新令牌
func (m *JWTManager) GenerateRefreshToken(userID string) (string, error) {
	return m.generateToken(userID, "refresh", m.refreshTokenTTL)
}

// generateToken 生成令牌
func (m *JWTManager) generateToken(userID, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "yidaiku",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ValidateToken 验证令牌
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// ValidateAccessToken 验证访问令牌
func (m *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.Type != "access" {
		return nil, errors.New("不是访问令牌")
	}
	return claims, nil
}

// ValidateRefreshToken 验证刷新令牌
func (m *JWTManager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.Type != "refresh" {
		return nil, errors.New("不是刷新令牌")
	}
	return claims, nil
}

// AccessTokenTTL 获取访问令牌有效期
func (m *JWTManager) AccessTokenTTL() time.Duration {
	return m.accessTokenTTL
}

// RefreshTokenTTL 获取刷新令牌有效期
func (m *JWTManager) RefreshTokenTTL() time.Duration {
	return m.refreshTokenTTL
}
