package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache Redis缓存客户端
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建Redis缓存客户端
func NewRedisCache(addr, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

// Get 获取缓存
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set 设置缓存
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	var data string
	switch v := value.(type) {
	case string:
		data = v
	default:
		bytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		data = string(bytes)
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

// Delete 删除缓存
func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists 检查缓存是否存在
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	return n > 0, err
}

// GetJSON 获取JSON缓存并解析
func (c *RedisCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

// SetJSON 设置JSON缓存
func (c *RedisCache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

// Keys 获取匹配的键
func (c *RedisCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	return c.client.Keys(ctx, pattern).Result()
}

// DeletePattern 删除匹配的键
func (c *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := c.Keys(ctx, pattern)
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return c.Delete(ctx, keys...)
	}
	return nil
}

// Incr 递增
func (c *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

// Expire 设置过期时间
func (c *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.client.Expire(ctx, key, ttl).Err()
}

// Client 获取原始客户端
func (c *RedisCache) Client() *redis.Client {
	return c.client
}

// Close 关闭连接
func (c *RedisCache) Close() error {
	return c.client.Close()
}
