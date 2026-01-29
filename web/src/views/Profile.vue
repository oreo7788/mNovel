<template>
  <section class="view-profile">
    <h1 class="page-title">个人中心</h1>
    <div class="card profile-card">
      <template v-if="isLoggedIn">
        <div class="profile-header">
          <img :src="avatar" alt="头像" class="profile-avatar" />
          <p class="profile-nickname">{{ nickname }}</p>
        </div>
        <button type="button" class="btn btn-secondary" @click="logout">退出登录</button>
      </template>
      <template v-else>
        <p class="profile-tip">登录后同步衣橱与收藏</p>
        <button type="button" class="btn btn-primary" @click="mockLogin">登录 / 注册</button>
      </template>
    </div>
  </section>
</template>

<script setup>
import { useUser } from '../composables/useUser'

const { isLoggedIn, avatar, nickname, login, logout } = useUser()

function mockLogin() {
  login({ nickname: '用户' + Math.floor(Math.random() * 1000) })
}
</script>

<style scoped>
.profile-card {
  padding: 24px;
  text-align: center;
}
.profile-header {
  margin-bottom: 20px;
}
.profile-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
  display: block;
  margin: 0 auto 12px;
}
.profile-nickname {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  color: var(--color-text);
}
.profile-tip {
  color: var(--color-text-secondary);
  margin: 0 0 16px;
}
.btn {
  padding: 10px 20px;
  border-radius: var(--radius-btn);
  font-size: 14px;
  border: none;
  cursor: pointer;
  font-family: inherit;
}
.btn-primary {
  background: var(--color-accent);
  color: #fff;
}
.btn-primary:hover {
  background: var(--color-accent-hover);
}
.btn-secondary {
  background: var(--color-bg-secondary);
  color: var(--color-text);
}
.btn-secondary:hover {
  background: var(--color-border);
}
</style>
