package adapter

// DoubaoAIAdapter 豆包AI适配器
// 注意：AI相关功能暂时不实现，此文件为预留接口
type DoubaoAIAdapter struct {
	apiKey    string
	apiSecret string
	baseURL   string
}

// NewDoubaoAIAdapter 创建豆包AI适配器（暂不实现）
func NewDoubaoAIAdapter(apiKey, apiSecret, baseURL string) *DoubaoAIAdapter {
	return &DoubaoAIAdapter{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   baseURL,
	}
}

// TODO: AI识别和推荐功能将在后续版本实现
// 当前MVP版本使用纯规则引擎进行穿搭推荐
