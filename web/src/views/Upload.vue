<template>
  <section class="view-upload">
    <h1 class="page-title">上传衣物</h1>
    <input
      ref="fileInput"
      type="file"
      accept="image/*"
      multiple
      hidden
      @change="onFileChange"
    />
    <div
      class="upload-zone"
      :class="{ 'upload-zone-dragover': isDragging }"
      role="button"
      tabindex="0"
      @click="fileInput?.click()"
      @dragover.prevent="isDragging = true"
      @dragleave.prevent="isDragging = false"
      @drop.prevent="onDrop"
    >
      <div>点击或拖拽照片到此处</div>
      <div class="hint">最多 {{ MAX_FILES }} 张/次，支持从相册选择或拍照（先缓存到本地）</div>
    </div>

    <div v-show="uploading" class="upload-progress">
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: progress + '%' }"></div>
      </div>
      <div class="progress-text">{{ progressText }}</div>
    </div>

    <div v-show="cachedList.length > 0" class="cached-list" aria-label="已缓存到本地的图片">
      <div class="cached-list-title">已缓存到本地（测试用）</div>
      <div
        v-for="item in cachedList"
        :key="item.id"
        class="cached-item"
      >
        <img :src="item.url" alt="" class="cached-thumb" />
        <div class="cached-info">
          <span class="cached-name">{{ item.name || '图片' }}</span>
          <span class="cached-meta">{{ item.sizeText }}</span>
        </div>
        <button type="button" class="btn btn-ghost cached-remove" title="从缓存移除" @click="removeItem(item.id)">
          移除
        </button>
      </div>
    </div>

    <div class="result-cards">
      <div class="result-card card">
        <div class="result-thumb">
          <div class="placeholder" style="height:100%;display:flex;align-items:center;justify-content:center;font-size:12px">图</div>
        </div>
        <div class="result-fields">
          <div class="field"><span class="label">品类</span>上衣</div>
          <div class="field"><span class="label">颜色</span>深蓝、白</div>
          <div class="field"><span class="label">风格</span>通勤、正式</div>
          <div class="field"><span class="label">季节</span>四季通用</div>
          <button class="btn btn-ghost" type="button" style="margin-top:8px">编辑</button>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import {
  getCachedImages,
  removeCachedImage,
  cacheImagesToLocal,
  MAX_UPLOAD_FILES,
} from '../composables/useUploadCache'
import { showToast } from '../utils/toast'

const MAX_FILES = MAX_UPLOAD_FILES
const fileInput = ref(null)
const isDragging = ref(false)
const uploading = ref(false)
const progress = ref(0)
const progressText = ref('')
const rawCached = ref([])

const cachedList = computed(() => {
  return (rawCached.value || []).map((item) => ({
    ...item,
    url: item.blob ? URL.createObjectURL(item.blob) : '',
    sizeText: item.size ? (item.size / 1024).toFixed(1) + ' KB' : '',
  }))
})

function loadCached() {
  getCachedImages().then((items) => {
    rawCached.value = items || []
  })
}

function onFileChange(e) {
  const files = e.target.files
  if (files?.length) handleFiles(files)
  e.target.value = ''
}

function onDrop(e) {
  isDragging.value = false
  const files = e.dataTransfer?.files
  if (files?.length) handleFiles(files)
}

function handleFiles(files) {
  const list = Array.from(files)
    .filter((f) => f.type?.startsWith('image/'))
    .slice(0, MAX_FILES)

  if (list.length === 0) {
    showToast('请选择图片文件')
    return
  }
  if (list.length < files.length) {
    showToast('已忽略非图片或超出 ' + MAX_FILES + ' 张的部分')
  }

  uploading.value = true
  progress.value = 0
  progressText.value = '正在缓存到本地… 0%'

  cacheImagesToLocal(list)
    .then((res) => {
      progress.value = 100
      progressText.value = '已缓存 ' + res.count + ' 张图片到本地'
      showToast('已缓存 ' + res.count + ' 张图片到本地')
      loadCached()
    })
    .catch((err) => {
      progressText.value = '缓存失败：' + (err?.message || '未知错误')
      showToast('缓存失败，请重试')
    })
    .finally(() => {
      uploading.value = false
    })
}

function removeItem(id) {
  removeCachedImage(id).then(loadCached)
}

onMounted(loadCached)
</script>
