<template>
  <div class="app">
    <div class="top-bar">
      <AppHeader />
      <BottomNav />
    </div>
    <main class="main-content">
      <RouterView v-slot="{ Component }">
        <component :is="Component" />
      </RouterView>
    </main>
  </div>
</template>

<script setup>
import AppHeader from './components/AppHeader.vue'
import BottomNav from './components/BottomNav.vue'
</script>
