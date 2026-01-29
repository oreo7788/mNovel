<template>
  <nav class="top-nav" aria-label="主导航">
    <router-link to="/" :class="['nav-item', { active: route.name === 'home' }]">
      <span class="icon">🏠</span>首页
    </router-link>
    <router-link to="/wardrobe" :class="['nav-item', { active: route.name === 'wardrobe' }]">
      <span class="icon">👔</span>衣橱
    </router-link>
    <router-link to="/recommend" :class="['nav-item', { active: route.name === 'recommend' }]">
      <span class="icon">✨</span>推荐
    </router-link>
    <router-link to="/favorites" :class="['nav-item', { active: route.name === 'favorites' }]">
      <span class="icon">♥</span>收藏
    </router-link>
    <router-link to="/upload" :class="['nav-item', { active: route.name === 'upload' }]">
      <span class="icon">➕</span>上传
    </router-link>
  </nav>
</template>

<script setup>
import { useRoute } from 'vue-router'
const route = useRoute()
</script>
