<template>
  <section class="view-settings">
    <h1 class="page-title">设置</h1>
    <div class="card settings-card">
      <p class="settings-placeholder">偏好、推送、账号等设置页（占位）</p>
      <router-link to="/terms" class="settings-link">服务条款</router-link>
    </div>
  </section>
</template>

<script setup>
</script>

<style scoped>
.settings-card {
  padding: 20px;
}

.settings-placeholder {
  color: var(--color-text-secondary);
  margin: 0 0 16px;
}

.settings-link {
  font-size: 14px;
  color: var(--color-accent);
  text-decoration: none;
}

.settings-link:hover {
  text-decoration: underline;
}
</style>
