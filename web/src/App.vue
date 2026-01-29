<template>
  <div class="app">
    <div class="top-bar">
      <AppHeader />
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
</script>
