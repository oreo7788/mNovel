package adapter

import (
	"context"
)

// DoubaoAIAdapter 豆包AI适配器（预留接口，暂不实现）
type DoubaoAIAdapter struct {
	// 预留字段
}

// NewDoubaoAIAdapter 创建豆包AI适配器（暂未接入，返回空实例）
func NewDoubaoAIAdapter(apiKey, apiSecret, baseURL string) *DoubaoAIAdapter {
	return &DoubaoAIAdapter{}
}

// RecognizeClothing 衣物识别（后续版本实现）
func (a *DoubaoAIAdapter) RecognizeClothing(ctx context.Context, imageData []byte) (map[string]interface{}, error) {
	// TODO: 后续版本接入豆包AI衣物识别
	return nil, nil
}

// RecommendOutfit 穿搭推荐增强（后续版本实现）
func (a *DoubaoAIAdapter) RecommendOutfit(ctx context.Context, wardrobe []map[string]interface{}, occasion string, weather map[string]interface{}) (map[string]interface{}, error) {
	// TODO: 后续版本接入豆包AI推荐增强
	return nil, nil
}
