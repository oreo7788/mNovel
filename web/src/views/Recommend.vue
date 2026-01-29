<template>
  <section class="view-recommend">
    <h1 class="page-title">穿搭推荐</h1>
    <div class="recommend-form">
      <select v-model="occasion" aria-label="场合">
        <option>职场通勤</option>
        <option>约会</option>
        <option>运动</option>
        <option>聚会</option>
        <option>面试</option>
        <option>旅行</option>
      </select>
      <select v-model="weather" aria-label="天气">
        <option>晴天</option>
        <option>多云</option>
        <option>雨天</option>
        <option>高温 &gt;30°C</option>
        <option>低温 &lt;10°C</option>
      </select>
      <button class="btn btn-primary" type="button" @click="generate">生成推荐</button>
    </div>

    <div class="outfit-cards">
      <article v-for="(outfit, idx) in outfits" :key="idx" class="outfit-card card">
        <div class="outfit-header">
          <span class="outfit-title">推荐方案 {{ idx + 1 }}</span>
          <div class="outfit-actions">
            <button type="button" title="收藏" @click="toggleFavorite(idx)">
              {{ favoritedIndices.includes(idx) ? '♥' : '♡' }}
            </button>
            <button type="button" title="分享">↗</button>
          </div>
        </div>
        <div class="outfit-items">
          <div v-for="(piece, i) in outfit.items" :key="i" class="outfit-item">
            <div class="item-thumb">
              <div class="placeholder">{{ piece.role }}</div>
            </div>
            <div class="item-name">{{ piece.name }}</div>
          </div>
        </div>
        <div class="highlights">{{ outfit.highlights }}</div>
      </article>
    </div>
  </section>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const occasion = ref('职场通勤')
const weather = ref('多云')
const favoritedIndices = ref([])

const outfits = ref([
  {
    items: [
      { role: '上衣', name: '深蓝衬衫' },
      { role: '下装', name: '黑色西裤' },
      { role: '鞋', name: '棕色皮鞋' },
    ],
    highlights: '深蓝色衬衫搭配黑色西裤，正式且优雅。点击单品可「换一件」。',
  },
  {
    items: [
      { role: '上衣', name: '白衬衫' },
      { role: '下装', name: '卡其裤' },
      { role: '鞋', name: '小白鞋' },
    ],
    highlights: '白衬衫 + 卡其裤 + 小白鞋，商务休闲，适合周五或轻松会议。',
  },
])

if (route.query.occasion) {
  occasion.value = route.query.occasion
}
if (route.query.daily) {
  // 每日推荐可预填
}

function generate() {
  // 后续对接推荐 API
}

function toggleFavorite(idx) {
  const i = favoritedIndices.value.indexOf(idx)
  if (i >= 0) {
    favoritedIndices.value = favoritedIndices.value.filter((x) => x !== idx)
  } else {
    favoritedIndices.value = [...favoritedIndices.value, idx]
  }
}

onMounted(() => {})
</script>
