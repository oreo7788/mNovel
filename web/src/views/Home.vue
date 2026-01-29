<template>
  <section class="view-home">
    <div class="weather-card card">
      <div class="location">北京 · 当前定位</div>
      <div class="temp">18°C</div>
      <div class="condition">多云</div>
      <div class="hint">适合薄外套，早晚略凉</div>
    </div>

    <div v-show="recentList.length > 0" class="recent-uploads card">
      <div class="recent-uploads-title">最近上传</div>
      <div class="recent-uploads-list">
        <div v-for="item in recentList" :key="item.id" class="recent-uploads-item">
          <img :src="item.url" alt="" class="recent-uploads-thumb" />
        </div>
      </div>
      <router-link to="/wardrobe" class="recent-uploads-link">去衣橱查看</router-link>
    </div>

    <div class="daily-entry-card card" role="button" tabindex="0" @click="goDailyRecommend">
      <div class="label">每日推荐</div>
      <div class="title">今日推荐</div>
      <div class="desc">点击查看今日穿搭推荐（已根据当前天气生成）</div>
    </div>

    <div class="occasion-section">
      <div class="section-label">选择场合，一键生成推荐</div>
      <div class="occasion-chips">
        <span
          v-for="occ in occasions"
          :key="occ"
          class="chip"
          :class="{ active: selectedOccasion === occ }"
          @click="toggleOccasion(occ)"
        >
          {{ occ }}
        </span>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCachedImages } from '../composables/useUploadCache'

const router = useRouter()
const cachedItems = ref([])
const selectedOccasion = ref(null)

const occasions = ['通勤', '约会', '运动', '聚会', '面试', '旅行', '周末逛街']

const recentList = computed(() => {
  const list = [...cachedItems.value].slice(-3)
  return list.map((item) => ({
    ...item,
    url: item.blob ? URL.createObjectURL(item.blob) : '',
  }))
})

function loadCached() {
  getCachedImages().then((items) => {
    cachedItems.value = items || []
  })
}

function toggleOccasion(occ) {
  selectedOccasion.value = selectedOccasion.value === occ ? null : occ
  if (selectedOccasion.value) {
    router.push({ path: '/recommend', query: { occasion: occ } })
  }
}

function goDailyRecommend() {
  router.push({ path: '/recommend', query: { daily: '1' } })
}

onMounted(loadCached)
</script>
