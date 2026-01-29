package service

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"sort"
	"time"

	"yidaiku-server/internal/adapter"
	"yidaiku-server/internal/model"
	"yidaiku-server/internal/repository"

	"github.com/google/uuid"
)

// OutfitService 穿搭推荐服务
type OutfitService struct {
	outfitRepo   *repository.OutfitRepository
	templateRepo *repository.OutfitTemplateRepository
	favoriteRepo *repository.FavoriteOutfitRepository
	recordRepo   *repository.OutfitRecordRepository
	ruleRepo     *repository.MatchingRuleRepository
	clothingRepo *repository.ClothingRepository
	userRepo     *repository.UserPreferenceRepository
	doubaoAI     *adapter.DoubaoAIAdapter
}

// NewOutfitService 创建穿搭服务实例
func NewOutfitService(doubaoAI *adapter.DoubaoAIAdapter) *OutfitService {
	return &OutfitService{
		outfitRepo:   repository.NewOutfitRepository(),
		templateRepo: repository.NewOutfitTemplateRepository(),
		favoriteRepo: repository.NewFavoriteOutfitRepository(),
		recordRepo:   repository.NewOutfitRecordRepository(),
		ruleRepo:     repository.NewMatchingRuleRepository(),
		clothingRepo: repository.NewClothingRepository(),
		userRepo:     repository.NewUserPreferenceRepository(),
		doubaoAI:     doubaoAI,
	}
}

// RecommendResult 推荐结果
type RecommendResult struct {
	Outfits      []OutfitDetail `json:"outfits"`
	IsColdStart  bool           `json:"is_cold_start"`
	MissingItems []string       `json:"missing_items,omitempty"` // 缺失的品类
	Templates    []model.OutfitTemplate `json:"templates,omitempty"` // 冷启动时的模板
}

// OutfitDetail 穿搭详情
type OutfitDetail struct {
	ID         string               `json:"id"`
	Items      []model.ClothingItem `json:"items"`
	Score      float64              `json:"score"`
	Highlights string               `json:"highlights"`
	Occasion   string               `json:"occasion"`
	Weather    *model.WeatherInfo   `json:"weather,omitempty"`
	Source     string               `json:"source"` // ai/rule/template
}

// Recommend 生成穿搭推荐
func (s *OutfitService) Recommend(ctx context.Context, userID string, req *model.RecommendRequest) (*RecommendResult, error) {
	// 获取用户衣橱
	wardrobe, err := s.clothingRepo.GetUserWardrobe(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 检查冷启动
	if len(wardrobe) < 10 {
		return s.coldStartRecommend(ctx, userID, req, wardrobe)
	}

	// 检查是否缺少关键品类
	missingCategories := s.checkMissingCategories(wardrobe)
	if len(missingCategories) > 0 {
		return s.coldStartRecommend(ctx, userID, req, wardrobe)
	}

	// 获取用户偏好
	preferences, _ := s.userRepo.GetByUserID(ctx, userID)

	// 生成推荐
	return s.generateRecommendations(ctx, userID, wardrobe, req, preferences)
}

// coldStartRecommend 冷启动推荐
func (s *OutfitService) coldStartRecommend(ctx context.Context, userID string, req *model.RecommendRequest, wardrobe []model.ClothingItem) (*RecommendResult, error) {
	result := &RecommendResult{
		IsColdStart:  true,
		MissingItems: s.checkMissingCategories(wardrobe),
	}

	// 获取匹配的模板
	templates, err := s.templateRepo.GetByOccasion(ctx, req.Occasion)
	if err != nil {
		return nil, err
	}
	if len(templates) == 0 {
		templates, _ = s.templateRepo.GetAll(ctx)
	}
	result.Templates = templates

	// 尝试基于已有衣物生成部分匹配的推荐
	if len(wardrobe) > 0 {
		outfits := s.generatePartialMatchOutfits(ctx, wardrobe, req)
		result.Outfits = outfits
	}

	return result, nil
}

// generateRecommendations 生成推荐
func (s *OutfitService) generateRecommendations(ctx context.Context, userID string, wardrobe []model.ClothingItem, req *model.RecommendRequest, preferences *model.UserPreference) (*RecommendResult, error) {
	result := &RecommendResult{
		IsColdStart: false,
	}

	// 获取最近7天的推荐（用于去重）
	recentOutfits, _ := s.outfitRepo.GetRecentRecommendations(ctx, userID, 7)
	recentItemSets := s.extractItemSets(recentOutfits)

	// 获取搭配规则
	rules, _ := s.ruleRepo.GetAll(ctx)

	// 筛选候选衣物
	candidates := s.filterCandidates(wardrobe, req)

	// 生成搭配组合
	combinations := s.generateCombinations(candidates)

	// 评分和排序
	scoredOutfits := s.scoreOutfits(combinations, req, preferences, rules)

	// 去重
	scoredOutfits = s.deduplicateOutfits(scoredOutfits, recentItemSets)

	// 保证多样性
	scoredOutfits = s.ensureDiversity(scoredOutfits)

	// 取Top 5
	if len(scoredOutfits) > 5 {
		scoredOutfits = scoredOutfits[:5]
	}

	// 尝试调用豆包AI增强推荐（如果配置了）
	if s.doubaoAI != nil && len(scoredOutfits) > 0 {
		s.enhanceWithAI(ctx, scoredOutfits, req, preferences)
	}

	// 保存推荐结果
	for _, outfit := range scoredOutfits {
		s.saveRecommendation(ctx, userID, outfit, req)
	}

	result.Outfits = scoredOutfits
	return result, nil
}

// filterCandidates 筛选候选衣物
func (s *OutfitService) filterCandidates(wardrobe []model.ClothingItem, req *model.RecommendRequest) map[string][]model.ClothingItem {
	candidates := make(map[string][]model.ClothingItem)

	for _, item := range wardrobe {
		// 根据天气筛选季节
		if req.Weather != nil {
			if !s.isSeasonMatch(item.Season, req.Weather.Temperature) {
				continue
			}
		}

		// 按品类分组
		candidates[item.Category] = append(candidates[item.Category], item)
	}

	return candidates
}

// isSeasonMatch 检查季节是否匹配
func (s *OutfitService) isSeasonMatch(season string, temperature int) bool {
	if season == "四季通用" {
		return true
	}
	if temperature > 20 && season == "春夏" {
		return true
	}
	if temperature < 20 && season == "秋冬" {
		return true
	}
	// 温度在15-25之间，两种季节都可以
	if temperature >= 15 && temperature <= 25 {
		return true
	}
	return false
}

// generateCombinations 生成搭配组合
func (s *OutfitService) generateCombinations(candidates map[string][]model.ClothingItem) [][]model.ClothingItem {
	var combinations [][]model.ClothingItem

	tops := candidates["上衣"]
	bottoms := append(candidates["裤子"], candidates["裙子"]...)
	shoes := candidates["鞋子"]
	accessories := candidates["配饰"]

	// 基础组合：上衣 + 下装 + 鞋子
	for _, top := range tops {
		for _, bottom := range bottoms {
			for _, shoe := range shoes {
				combo := []model.ClothingItem{top, bottom, shoe}
				combinations = append(combinations, combo)

				// 可选添加配饰
				if len(accessories) > 0 {
					accessory := accessories[rand.Intn(len(accessories))]
					comboWithAccessory := append([]model.ClothingItem{}, combo...)
					comboWithAccessory = append(comboWithAccessory, accessory)
					combinations = append(combinations, comboWithAccessory)
				}
			}
		}
	}

	// 限制组合数量
	if len(combinations) > 100 {
		rand.Shuffle(len(combinations), func(i, j int) {
			combinations[i], combinations[j] = combinations[j], combinations[i]
		})
		combinations = combinations[:100]
	}

	return combinations
}

// scoreOutfits 对搭配进行评分
func (s *OutfitService) scoreOutfits(combinations [][]model.ClothingItem, req *model.RecommendRequest, preferences *model.UserPreference, rules []model.MatchingRule) []OutfitDetail {
	var outfits []OutfitDetail

	for _, combo := range combinations {
		score := 0.0
		highlights := ""

		// 1. 场合适配分数 (30%)
		occasionScore := s.scoreOccasion(combo, req.Occasion, rules)
		score += occasionScore * 0.3

		// 2. 天气适配分数 (20%)
		if req.Weather != nil {
			weatherScore := s.scoreWeather(combo, req.Weather)
			score += weatherScore * 0.2
		} else {
			score += 0.8 * 0.2 // 无天气信息时给默认分
		}

		// 3. 用户偏好分数 (30%)
		if preferences != nil {
			prefScore := s.scorePreference(combo, preferences)
			score += prefScore * 0.3
		} else {
			score += 0.7 * 0.3
		}

		// 4. 基础搭配规则分数 (20%)
		ruleScore, ruleHighlights := s.scoreRules(combo, rules)
		score += ruleScore * 0.2
		if ruleHighlights != "" {
			highlights = ruleHighlights
		}

		// 生成默认亮点
		if highlights == "" {
			highlights = s.generateHighlights(combo, req.Occasion)
		}

		outfits = append(outfits, OutfitDetail{
			ID:         uuid.New().String(),
			Items:      combo,
			Score:      score,
			Highlights: highlights,
			Occasion:   req.Occasion,
			Weather:    req.Weather,
			Source:     "rule",
		})
	}

	// 按分数排序
	sort.Slice(outfits, func(i, j int) bool {
		return outfits[i].Score > outfits[j].Score
	})

	return outfits
}

// scoreOccasion 场合适配评分
func (s *OutfitService) scoreOccasion(items []model.ClothingItem, occasion string, rules []model.MatchingRule) float64 {
	score := 0.5 // 基础分

	// 找到对应场合的规则
	for _, rule := range rules {
		if rule.RuleType == model.RuleTypeOccasion {
			var cond map[string]interface{}
			condBytes, _ := json.Marshal(rule.Condition)
			json.Unmarshal(condBytes, &cond)

			if cond["occasion"] == occasion {
				// 检查风格匹配
				if styles, ok := rule.Action["styles"].([]interface{}); ok {
					for _, item := range items {
						for _, style := range []string(item.Styles) {
							for _, s := range styles {
								if s == style {
									score += 0.1
								}
							}
						}
					}
				}
			}
		}
	}

	if score > 1.0 {
		score = 1.0
	}
	return score
}

// scoreWeather 天气适配评分
func (s *OutfitService) scoreWeather(items []model.ClothingItem, weather *model.WeatherInfo) float64 {
	score := 0.5

	for _, item := range items {
		if s.isSeasonMatch(item.Season, weather.Temperature) {
			score += 0.1
		}
	}

	if score > 1.0 {
		score = 1.0
	}
	return score
}

// scorePreference 用户偏好评分
func (s *OutfitService) scorePreference(items []model.ClothingItem, preferences *model.UserPreference) float64 {
	score := 0.5

	prefStyles := map[string]bool{}
	for _, s := range []string(preferences.Styles) {
		prefStyles[s] = true
	}

	prefColors := map[string]bool{}
	for _, c := range []string(preferences.Colors) {
		prefColors[c] = true
	}

	for _, item := range items {
		// 检查风格偏好
		for _, style := range []string(item.Styles) {
			if prefStyles[style] {
				score += 0.1
			}
		}
		// 检查颜色偏好
		for _, color := range []string(item.Colors) {
			if prefColors[color] {
				score += 0.05
			}
		}
	}

	if score > 1.0 {
		score = 1.0
	}
	return score
}

// scoreRules 规则评分
func (s *OutfitService) scoreRules(items []model.ClothingItem, rules []model.MatchingRule) (float64, string) {
	score := 0.5
	highlights := ""

	// 检查颜色搭配
	colors := make(map[string]bool)
	for _, item := range items {
		for _, c := range []string(item.Colors) {
			colors[c] = true
		}
	}

	// 中性色搭配加分
	neutralColors := []string{"黑色", "白色", "灰色", "米色"}
	neutralCount := 0
	for _, nc := range neutralColors {
		if colors[nc] {
			neutralCount++
		}
	}
	if neutralCount >= 2 {
		score += 0.15
		highlights = "中性色搭配，简约大方"
	}

	// 颜色数量不超过3种
	if len(colors) <= 3 {
		score += 0.1
	}

	// 检查风格一致性
	styleCount := make(map[string]int)
	for _, item := range items {
		for _, s := range []string(item.Styles) {
			styleCount[s]++
		}
	}
	maxStyleCount := 0
	dominantStyle := ""
	for style, count := range styleCount {
		if count > maxStyleCount {
			maxStyleCount = count
			dominantStyle = style
		}
	}
	if maxStyleCount >= len(items)/2 {
		score += 0.1
		if highlights == "" {
			highlights = dominantStyle + "风格统一，整体协调"
		}
	}

	if score > 1.0 {
		score = 1.0
	}
	return score, highlights
}

// generateHighlights 生成搭配亮点
func (s *OutfitService) generateHighlights(items []model.ClothingItem, occasion string) string {
	// 简单的亮点生成逻辑
	if len(items) >= 2 {
		return items[0].Category + "搭配" + items[1].Category + "，适合" + occasion
	}
	return "经典搭配，适合" + occasion
}

// deduplicateOutfits 去重
func (s *OutfitService) deduplicateOutfits(outfits []OutfitDetail, recentItemSets []map[string]bool) []OutfitDetail {
	var result []OutfitDetail

	for _, outfit := range outfits {
		itemSet := make(map[string]bool)
		for _, item := range outfit.Items {
			itemSet[item.ID] = true
		}

		// 检查是否与最近推荐重复
		isDuplicate := false
		for _, recent := range recentItemSets {
			matchCount := 0
			for id := range itemSet {
				if recent[id] {
					matchCount++
				}
			}
			// 如果超过一半的单品重复，认为是重复推荐
			if matchCount > len(itemSet)/2 {
				isDuplicate = true
				break
			}
		}

		if !isDuplicate {
			result = append(result, outfit)
		}
	}

	return result
}

// ensureDiversity 保证多样性
func (s *OutfitService) ensureDiversity(outfits []OutfitDetail) []OutfitDetail {
	if len(outfits) <= 3 {
		return outfits
	}

	var result []OutfitDetail
	usedTopIDs := make(map[string]bool)

	for _, outfit := range outfits {
		// 检查上衣是否已使用
		topUsed := false
		for _, item := range outfit.Items {
			if item.Category == "上衣" && usedTopIDs[item.ID] {
				topUsed = true
				break
			}
		}

		if !topUsed {
			result = append(result, outfit)
			for _, item := range outfit.Items {
				if item.Category == "上衣" {
					usedTopIDs[item.ID] = true
				}
			}
		}

		if len(result) >= 5 {
			break
		}
	}

	return result
}

// extractItemSets 提取单品集合
func (s *OutfitService) extractItemSets(outfits []model.OutfitRecommendation) []map[string]bool {
	var sets []map[string]bool
	for _, outfit := range outfits {
		set := make(map[string]bool)
		for _, id := range []string(outfit.ItemIDs) {
			set[id] = true
		}
		sets = append(sets, set)
	}
	return sets
}

// checkMissingCategories 检查缺失的品类
func (s *OutfitService) checkMissingCategories(wardrobe []model.ClothingItem) []string {
	categories := make(map[string]bool)
	for _, item := range wardrobe {
		categories[item.Category] = true
	}

	required := []string{"上衣", "鞋子"}
	hasBottom := categories["裤子"] || categories["裙子"]

	var missing []string
	for _, req := range required {
		if !categories[req] {
			missing = append(missing, req)
		}
	}
	if !hasBottom {
		missing = append(missing, "下装(裤子/裙子)")
	}

	return missing
}

// generatePartialMatchOutfits 生成部分匹配的搭配
func (s *OutfitService) generatePartialMatchOutfits(ctx context.Context, wardrobe []model.ClothingItem, req *model.RecommendRequest) []OutfitDetail {
	var outfits []OutfitDetail

	// 简单的部分匹配逻辑
	if len(wardrobe) >= 2 {
		outfit := OutfitDetail{
			ID:         uuid.New().String(),
			Items:      wardrobe[:min(3, len(wardrobe))],
			Score:      0.5,
			Highlights: "基于您现有的衣物推荐",
			Occasion:   req.Occasion,
			Source:     "partial",
		}
		outfits = append(outfits, outfit)
	}

	return outfits
}

// enhanceWithAI 使用AI增强推荐
func (s *OutfitService) enhanceWithAI(ctx context.Context, outfits []OutfitDetail, req *model.RecommendRequest, preferences *model.UserPreference) {
	// TODO: 调用豆包AI获取更好的搭配亮点描述
	// 这里保留接口，后续可以接入
}

// saveRecommendation 保存推荐结果
func (s *OutfitService) saveRecommendation(ctx context.Context, userID string, outfit OutfitDetail, req *model.RecommendRequest) {
	itemIDs := make([]string, len(outfit.Items))
	for i, item := range outfit.Items {
		itemIDs[i] = item.ID
	}

	var weatherData model.JSONMap
	if req.Weather != nil {
		weatherData = model.JSONMap{
			"temperature": req.Weather.Temperature,
			"condition":   req.Weather.Condition,
		}
	}

	recommendation := &model.OutfitRecommendation{
		ID:         outfit.ID,
		UserID:     userID,
		ItemIDs:    model.JSONArray(itemIDs),
		Occasion:   req.Occasion,
		Weather:    weatherData,
		Score:      outfit.Score,
		Highlights: outfit.Highlights,
		Source:     outfit.Source,
		CreatedAt:  time.Now(),
	}

	s.outfitRepo.Create(ctx, recommendation)
}

// ReplaceItem 换一件功能
func (s *OutfitService) ReplaceItem(ctx context.Context, userID string, req *model.ReplaceItemRequest) (*OutfitDetail, error) {
	// 获取原推荐
	outfit, err := s.outfitRepo.GetByID(ctx, req.OutfitID)
	if err != nil || outfit == nil {
		return nil, errors.New("推荐不存在")
	}
	if outfit.UserID != userID {
		return nil, errors.New("无权限操作")
	}

	// 获取新单品
	newItem, err := s.clothingRepo.GetByID(ctx, req.NewItemID)
	if err != nil || newItem == nil {
		return nil, errors.New("新单品不存在")
	}
	if newItem.UserID != userID {
		return nil, errors.New("无权限使用此单品")
	}

	// 获取所有单品
	items, err := s.clothingRepo.GetByIDs(ctx, []string(outfit.ItemIDs))
	if err != nil {
		return nil, err
	}

	// 替换单品
	var newItems []model.ClothingItem
	for _, item := range items {
		if item.ID == req.OldItemID {
			newItems = append(newItems, *newItem)
		} else {
			newItems = append(newItems, item)
		}
	}

	// 重新评分
	rules, _ := s.ruleRepo.GetAll(ctx)
	preferences, _ := s.userRepo.GetByUserID(ctx, userID)

	score := 0.6 // 基础分
	_, highlights := s.scoreRules(newItems, rules)
	if preferences != nil {
		score = s.scorePreference(newItems, preferences)
	}

	if highlights == "" {
		highlights = "已重新优化搭配"
	}

	// 更新推荐
	newItemIDs := make([]string, len(newItems))
	for i, item := range newItems {
		newItemIDs[i] = item.ID
	}

	outfit.ItemIDs = model.JSONArray(newItemIDs)
	outfit.Score = score
	outfit.Highlights = highlights
	s.outfitRepo.Update(ctx, outfit)

	return &OutfitDetail{
		ID:         outfit.ID,
		Items:      newItems,
		Score:      score,
		Highlights: highlights,
		Occasion:   outfit.Occasion,
		Source:     "replaced",
	}, nil
}

// SetFavorite 设置收藏
func (s *OutfitService) SetFavorite(ctx context.Context, userID, outfitID string, favorite bool) error {
	if favorite {
		// 检查是否已收藏
		exists, _ := s.favoriteRepo.Exists(ctx, userID, outfitID)
		if exists {
			return nil
		}
		fav := &model.FavoriteOutfit{
			UserID:   userID,
			OutfitID: outfitID,
		}
		if err := s.favoriteRepo.Create(ctx, fav); err != nil {
			return err
		}
	} else {
		if err := s.favoriteRepo.Delete(ctx, userID, outfitID); err != nil {
			return err
		}
	}
	return s.outfitRepo.SetFavorite(ctx, outfitID, favorite)
}

// SetApplied 设置应用状态
func (s *OutfitService) SetApplied(ctx context.Context, userID, outfitID string) error {
	return s.outfitRepo.SetApplied(ctx, outfitID, true)
}

// GetFavorites 获取收藏列表
func (s *OutfitService) GetFavorites(ctx context.Context, userID string, page, pageSize int) ([]OutfitDetail, int64, error) {
	favorites, total, err := s.favoriteRepo.GetByUser(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var details []OutfitDetail
	for _, fav := range favorites {
		outfit, err := s.outfitRepo.GetByID(ctx, fav.OutfitID)
		if err != nil || outfit == nil {
			continue
		}
		items, _ := s.clothingRepo.GetByIDs(ctx, []string(outfit.ItemIDs))

		details = append(details, OutfitDetail{
			ID:         outfit.ID,
			Items:      items,
			Score:      outfit.Score,
			Highlights: outfit.Highlights,
			Occasion:   outfit.Occasion,
			Source:     outfit.Source,
		})
	}

	return details, total, nil
}

// GetTemplates 获取穿搭模板
func (s *OutfitService) GetTemplates(ctx context.Context, occasion string) ([]model.OutfitTemplate, error) {
	if occasion != "" {
		return s.templateRepo.GetByOccasion(ctx, occasion)
	}
	return s.templateRepo.GetAll(ctx)
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
