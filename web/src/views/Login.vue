<template>
  <section class="view-login">
    <div class="login-card card">
      <h1 class="login-title">登录</h1>
      <p class="login-desc">登录后同步衣橱与收藏，享受个性化推荐</p>

      <form class="login-form" @submit.prevent="handleSubmit">
        <div class="form-group">
          <label for="username" class="form-label">用户名</label>
          <input
            id="username"
            v-model="username"
            type="text"
            class="form-input"
            placeholder="请输入用户名"
            maxlength="20"
            autocomplete="username"
            required
          />
        </div>

        <div class="form-group">
          <label for="password" class="form-label">密码</label>
          <input
            id="password"
            v-model="password"
            type="password"
            class="form-input"
            placeholder="请输入密码"
            autocomplete="current-password"
            required
          />
        </div>

        <div class="form-options">
          <label class="form-checkbox">
            <input v-model="autoLogin" type="checkbox" class="form-checkbox-input" />
            <span class="form-checkbox-text">自动登录</span>
          </label>
          <router-link to="/forgot-password" class="forgot-password-link">找回密码</router-link>
        </div>

        <button type="submit" class="btn btn-primary login-btn" :disabled="submitting">
          {{ submitting ? '登录中…' : '登录' }}
        </button>

        <router-link to="/register" class="btn btn-ghost register-btn">去注册</router-link>
      </form>

      <router-link to="/" class="login-back">返回首页</router-link>
    </div>
  </section>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUser } from '../composables/useUser'
import { showToast } from '../utils/toast'

const router = useRouter()
const route = useRoute()
const { isLoggedIn, login } = useUser()

const username = ref('')
const password = ref('')
const autoLogin = ref(true)
const submitting = ref(false)

onMounted(() => {
  if (isLoggedIn.value) {
    const redirect = route.query.redirect || '/'
    router.replace(redirect)
  }
})

async function handleSubmit() {
  const name = username.value.trim()
  const pwd = password.value
  if (!name || !pwd) return

  submitting.value = true
  try {
    const ok = await login(name, pwd, autoLogin.value)
    if (ok) {
      showToast('登录成功')
      const redirect = route.query.redirect || '/'
      router.replace(redirect)
    } else {
      showToast('用户名或密码错误')
    }
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.view-login {
  min-height: 60vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 0;
}

.login-card {
  width: 100%;
  max-width: 400px;
  padding: 32px 24px;
}

.login-title {
  font-family: var(--font-serif);
  font-weight: 600;
  font-size: 24px;
  margin: 0 0 8px;
  color: var(--color-text);
}

.login-desc {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin: 0 0 24px;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text);
}

.form-input {
  padding: 12px 14px;
  font-size: 15px;
  font-family: inherit;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-btn);
  background: #fff;
  color: var(--color-text);
  transition: border-color 0.2s;
}

.form-input::placeholder {
  color: var(--color-text-secondary);
}

.form-input:focus {
  outline: none;
  border-color: var(--color-accent);
}

.form-options {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.form-checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
  color: var(--color-text);
}

.forgot-password-link {
  font-size: 14px;
  color: var(--color-accent);
  text-decoration: none;
}

.forgot-password-link:hover {
  text-decoration: underline;
}

.form-checkbox-input {
  width: 18px;
  height: 18px;
  accent-color: var(--color-accent);
  cursor: pointer;
}

.form-checkbox-text {
  user-select: none;
}

.login-btn {
  width: 100%;
  padding: 12px 20px;
  font-size: 16px;
  margin-top: 8px;
}

.login-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.register-btn {
  width: 100%;
  padding: 12px 20px;
  font-size: 15px;
  text-decoration: none;
  text-align: center;
}

.login-back {
  display: block;
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: var(--color-text-secondary);
  text-decoration: none;
}

.login-back:hover {
  color: var(--color-accent);
}
</style>
