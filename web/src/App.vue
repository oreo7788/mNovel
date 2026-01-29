<template>
  <div class="app">
    <AppHeader />
    <main class="main-content">
      <RouterView v-slot="{ Component }">
        <component :is="Component" />
      </RouterView>
    </main>
    <BottomNav />
  </div>
</template>

<script setup>
import AppHeader from './components/AppHeader.vue'
import BottomNav from './components/BottomNav.vue'
</script>
