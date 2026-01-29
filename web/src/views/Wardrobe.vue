<template>
  <section class="view-wardrobe">
    <h1 class="page-title">我的衣橱</h1>
    <div class="wardrobe-toolbar">
      <div class="left">
        <button class="btn btn-primary" type="button" @click="router.push('/upload')">上传衣物</button>
        <router-link to="/upload" class="btn btn-secondary">上传队列</router-link>
      </div>
      <div class="view-toggle">
        <button
          type="button"
          :aria-pressed="layout === 'grid'"
          :class="{ active: layout === 'grid' }"
          @click="layout = 'grid'"
        >
          网格
        </button>
        <button
          type="button"
          :aria-pressed="layout === 'list'"
          :class="{ active: layout === 'list' }"
          @click="layout = 'list'"
        >
          列表
        </button>
      </div>
    </div>

    <div class="filter-bar">
      <span
        v-for="f in filters"
        :key="f.id"
        class="chip"
        :class="{ active: activeFilter === f.id }"
        @click="activeFilter = f.id"
      >
        {{ f.label }}
      </span>
    </div>

    <div
      class="wardrobe-grid"
      :class="{ 'wardrobe-list': layout === 'list' }"
      :data-layout="layout"
    >
      <article
        v-for="item in cachedList"
        :key="item.id"
        class="clothing-card card"
        @click="openItem(item)"
      >
        <div class="thumb">
          <img v-if="item.url" :src="item.url" alt="" />
          <div v-else class="placeholder">{{ item.category || '图' }}</div>
        </div>
        <div class="info">
          <div class="category">{{ item.displayName }}</div>
        </div>
      </article>
      <article class="clothing-card card wardrobe-add-card" @click="router.push('/upload')">
        <div class="thumb"><div class="placeholder">+</div></div>
        <div class="info"><div class="category">添加衣物</div></div>
      </article>
    </div>
  </section>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCachedImages } from '../composables/useUploadCache'

const router = useRouter()
const rawList = ref([])
const layout = ref('grid')
const activeFilter = ref('all')

const filters = [
  { id: 'all', label: '全部' },
  { id: 'retired', label: '已淘汰(0)' },
  { id: 'top', label: '上衣' },
  { id: 'pants', label: '裤子' },
  { id: 'skirt', label: '裙子' },
  { id: 'shoes', label: '鞋子' },
  { id: 'acc', label: '配饰' },
  { id: 'casual', label: '休闲' },
  { id: 'commute', label: '通勤' },
  { id: 'ss', label: '春夏' },
  { id: 'aw', label: '秋冬' },
]

const cachedList = computed(() => {
  return (rawList.value || []).map((item) => {
    const name = (item.name || '图片').replace(/\.[^.]+$/, '')
    return {
      ...item,
      url: item.blob ? URL.createObjectURL(item.blob) : '',
      displayName: name,
    }
  })
})

function loadCached() {
  getCachedImages().then((items) => {
    rawList.value = items || []
  })
}

function openItem(item) {
  // 后续可做详情/编辑
}

onMounted(loadCached)
</script>
