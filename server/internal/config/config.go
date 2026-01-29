package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Qiniu    QiniuConfig    `mapstructure:"qiniu"`
	DoubaoAI DoubaoAIConfig `mapstructure:"doubao_ai"`
	Weather  WeatherConfig  `mapstructure:"weather"`
	Upload   UploadConfig   `mapstructure:"upload"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"` // debug, release, test
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

// DSN 获取数据库连接字符串
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.Database)
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// Addr 获取Redis地址
func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret           string        `mapstructure:"secret"`
	AccessTokenTTL   time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL  time.Duration `mapstructure:"refresh_token_ttl"`
}

// QiniuConfig 七牛云配置
type QiniuConfig struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Domain    string `mapstructure:"domain"`
}

// DoubaoAIConfig 豆包AI配置
type DoubaoAIConfig struct {
	APIKey    string        `mapstructure:"api_key"`
	APISecret string        `mapstructure:"api_secret"`
	BaseURL   string        `mapstructure:"base_url"`
	Timeout   time.Duration `mapstructure:"timeout"`
}

// WeatherConfig 天气API配置
type WeatherConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	MaxFileSize    int64 `mapstructure:"max_file_size"`     // 最大文件大小（字节）
	ChunkSize      int64 `mapstructure:"chunk_size"`        // 分块大小（字节）
	MaxConcurrent  int   `mapstructure:"max_concurrent"`    // 最大并发数
	AllowedFormats []string `mapstructure:"allowed_formats"` // 允许的文件格式
}

var globalConfig *Config

// Load 加载配置
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 设置默认值
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	globalConfig = config
	return config, nil
}

// Get 获取全局配置
func Get() *Config {
	return globalConfig
}

// setDefaults 设置默认配置值
func setDefaults() {
	// Server
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")

	// Database
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)

	// Redis
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)

	// JWT
	viper.SetDefault("jwt.access_token_ttl", "24h")
	viper.SetDefault("jwt.refresh_token_ttl", "168h") // 7天

	// DoubaoAI
	viper.SetDefault("doubao_ai.base_url", "https://api.doubao.com/v1")
	viper.SetDefault("doubao_ai.timeout", "30s")

	// Upload
	viper.SetDefault("upload.max_file_size", 10*1024*1024) // 10MB
	viper.SetDefault("upload.chunk_size", 1*1024*1024)     // 1MB
	viper.SetDefault("upload.max_concurrent", 3)
	viper.SetDefault("upload.allowed_formats", []string{"jpg", "jpeg", "png", "webp"})
}
