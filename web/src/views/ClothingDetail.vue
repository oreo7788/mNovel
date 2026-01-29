<template>
  <section class="view-clothing-detail">
    <div class="detail-header">
      <router-link to="/wardrobe" class="back-link">← 返回衣橱</router-link>
    </div>

    <template v-if="loading">
      <div class="detail-loading">
        <p>加载中…</p>
      </div>
    </template>

    <template v-else-if="error">
      <div class="detail-error card">
        <p class="error-text">{{ error }}</p>
        <router-link to="/wardrobe" class="btn btn-primary">返回衣橱</router-link>
      </div>
    </template>

    <template v-else-if="item">
      <div class="detail-layout">
        <div class="detail-image card">
          <div class="image-wrap">
            <img
              v-if="item.image_url"
              :src="item.image_url"
              :alt="item.category || '衣物'"
              class="main-image"
            />
            <div v-else class="image-placeholder">{{ item.category || '图' }}</div>
          </div>
        </div>

        <div class="detail-info">
          <div class="detail-meta card">
            <h1 class="detail-title">{{ categoryLabel }} {{ item.subcategory ? `· ${item.subcategory}` : '' }}</h1>
            <p v-if="item.is_retired" class="retired-badge">已淘汰</p>

            <dl class="detail-fields" v-if="hasAnyMeta">
              <template v-if="item.colors && item.colors.length">
                <dt>颜色</dt>
                <dd class="chips">
                  <span v-for="(c, i) in item.colors" :key="i" class="chip">{{ c }}</span>
                </dd>
              </template>
              <template v-if="item.styles && item.styles.length">
                <dt>风格</dt>
                <dd class="chips">
                  <span v-for="(s, i) in item.styles" :key="i" class="chip">{{ s }}</span>
                </dd>
              </template>
              <template v-if="item.season">
                <dt>季节</dt>
                <dd>{{ item.season }}</dd>
              </template>
              <template v-if="item.tags && item.tags.length">
                <dt>标签</dt>
                <dd class="chips">
                  <span v-for="(t, i) in item.tags" :key="i" class="chip chip-tag">{{ t }}</span>
                </dd>
              </template>
              <template v-if="item.notes">
                <dt>备注</dt>
                <dd class="notes">{{ item.notes }}</dd>
              </template>
              <dt>添加时间</dt>
              <dd class="date">{{ formatDate(item.created_at) }}</dd>
            </dl>
          </div>

          <div class="detail-actions">
            <button type="button" class="btn btn-primary" disabled title="编辑功能开发中">编辑</button>
            <button type="button" class="btn btn-danger" :disabled="deleting" @click="confirmDelete">
              {{ deleting ? '删除中…' : '删除' }}
            </button>
          </div>
        </div>
      </div>
    </template>
  </section>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUser } from '../composables/useUser'
import { getClothingItemApi, deleteClothingItemApi } from '../api'
import { showToast } from '../utils/toast'

const route = useRoute()
const router = useRouter()
const { user, isLoggedIn } = useUser()

const item = ref(null)
const loading = ref(true)
const error = ref('')
const deleting = ref(false)

const categoryMap = {
  top: '上衣',
  pants: '裤子',
  skirt: '裙子',
  shoes: '鞋子',
  acc: '配饰',
}

const categoryLabel = computed(() => {
  if (!item.value) return ''
  const c = (item.value.category || '').toLowerCase()
  return categoryMap[c] || item.value.category || '衣物'
})

const hasAnyMeta = computed(() => {
  const i = item.value
  if (!i) return false
  return (
    (i.colors && i.colors.length) ||
    (i.styles && i.styles.length) ||
    i.season ||
    (i.tags && i.tags.length) ||
    i.notes
  )
})

function formatDate(val) {
  if (!val) return '—'
  const d = new Date(val)
  if (isNaN(d.getTime())) return val
  return d.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

function fetchItem() {
  const id = route.params.id
  if (!id) {
    error.value = '缺少衣物 ID'
    loading.value = false
    return
  }
  if (!isLoggedIn.value || !user.value?.access_token) {
    error.value = '请先登录后查看'
    loading.value = false
    return
  }
  loading.value = true
  error.value = ''
  getClothingItemApi(id, user.value.access_token).then((res) => {
    loading.value = false
    if (res.ok) {
      item.value = res.data
    } else {
      error.value = res.message || '加载失败'
      if (res.message && res.message.includes('不存在')) {
        // 可保持 404 语义
      }
    }
  }).catch((e) => {
    loading.value = false
    error.value = e.message || '网络错误'
  })
}

function confirmDelete() {
  if (!item.value || deleting.value) return
  if (!window.confirm('确定要删除这件衣物吗？删除后可到回收站恢复。')) return
  deleting.value = true
  deleteClothingItemApi(item.value.id, user.value?.access_token, false).then((res) => {
    deleting.value = false
    if (res.ok) {
      showToast('已删除')
      router.replace('/wardrobe')
    } else {
      showToast(res.message || '删除失败')
    }
  }).catch(() => {
    deleting.value = false
    showToast('删除失败')
  })
}

onMounted(fetchItem)
watch(() => route.params.id, fetchItem)
</script>

<style scoped>
.view-clothing-detail {
  max-width: var(--max-width);
  margin: 0 auto;
  padding: var(--page-padding-m);
}

@media (min-width: 768px) {
  .view-clothing-detail {
    padding: var(--page-padding-d);
  }
}

.detail-header {
  margin-bottom: 20px;
}

.back-link {
  color: var(--color-text-secondary);
  text-decoration: none;
  font-size: 14px;
}

.back-link:hover {
  color: var(--color-accent);
}

.detail-loading,
.detail-error {
  text-align: center;
  padding: 48px 24px;
}

.detail-error .error-text {
  margin: 0 0 16px;
  color: var(--color-text-secondary);
}

.detail-layout {
  display: grid;
  gap: 24px;
}

@media (min-width: 768px) {
  .detail-layout {
    grid-template-columns: 1fr 1fr;
    align-items: start;
  }
}

.detail-image.card {
  padding: 0;
  overflow: hidden;
}

.image-wrap {
  aspect-ratio: 1;
  background: var(--color-bg-secondary);
  position: relative;
}

.main-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.image-placeholder {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  color: var(--color-border);
  font-family: var(--font-serif);
}

.detail-meta.card {
  padding: 20px;
}

.detail-title {
  font-family: var(--font-serif);
  font-size: 22px;
  font-weight: 600;
  margin: 0 0 8px;
  color: var(--color-text);
}

.retired-badge {
  display: inline-block;
  font-size: 12px;
  color: var(--color-warning);
  margin: 0 0 16px;
  padding: 2px 8px;
  border: 1px solid var(--color-warning);
  border-radius: var(--radius-chip);
}

.detail-fields {
  margin: 0;
  display: grid;
  gap: 8px 16px;
  grid-template-columns: auto 1fr;
}

.detail-fields dt {
  margin: 0;
  font-size: 13px;
  color: var(--color-text-secondary);
  font-weight: 500;
}

.detail-fields dd {
  margin: 0;
  font-size: 15px;
  color: var(--color-text);
}

.detail-fields .chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.detail-fields .chip {
  font-size: 13px;
  padding: 4px 10px;
  border-radius: var(--radius-chip);
  background: var(--color-bg-secondary);
  color: var(--color-text);
}

.detail-fields .chip-tag {
  background: rgba(196, 165, 116, 0.15);
  color: var(--color-accent);
}

.detail-fields .notes {
  white-space: pre-wrap;
  line-height: 1.5;
}

.detail-fields .date {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.detail-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.detail-actions .btn {
  padding: 10px 20px;
  border-radius: var(--radius-btn);
  font-size: 14px;
  border: none;
  cursor: pointer;
  font-family: inherit;
  text-decoration: none;
  display: inline-block;
  text-align: center;
}

.detail-actions .btn-primary {
  background: var(--color-accent);
  color: #fff;
}

.detail-actions .btn-primary:hover:not(:disabled) {
  background: var(--color-accent-hover);
}

.detail-actions .btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.detail-actions .btn-danger {
  background: transparent;
  color: var(--color-error);
  border: 1px solid var(--color-error);
}

.detail-actions .btn-danger:hover:not(:disabled) {
  background: rgba(184, 84, 80, 0.1);
}

.detail-actions .btn-danger:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
