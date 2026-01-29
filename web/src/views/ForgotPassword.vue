<template>
  <section class="view-forgot">
    <div class="forgot-card card">
      <h1 class="forgot-title">找回密码</h1>
      <p class="forgot-desc">输入注册邮箱，获取验证码后验证并重置密码</p>

      <form class="forgot-form" @submit.prevent="handleSubmit">
        <div class="form-group">
          <label for="email" class="form-label">邮箱</label>
          <input
            id="email"
            v-model="email"
            type="email"
            class="form-input"
            :class="{ 'form-input-error': errors.email }"
            placeholder="请输入注册邮箱"
            autocomplete="email"
            required
          />
          <p v-if="errors.email" class="form-error">{{ errors.email }}</p>
        </div>

        <div class="form-group">
          <label for="code" class="form-label">验证码</label>
          <div class="form-row">
            <input
              id="code"
              v-model="code"
              type="text"
              class="form-input form-input-code"
              :class="{ 'form-input-error': errors.code }"
              placeholder="请输入验证码"
              maxlength="6"
              autocomplete="one-time-code"
              required
            />
            <button
              type="button"
              class="btn btn-secondary send-code-btn"
              :disabled="countdown > 0 || sendingCode"
              @click="sendCode"
            >
              {{ countdown > 0 ? `${countdown}s 后重发` : sendingCode ? '发送中…' : '获取验证码' }}
            </button>
          </div>
          <p v-if="errors.code" class="form-error">{{ errors.code }}</p>
        </div>

        <button type="submit" class="btn btn-primary submit-btn" :disabled="submitting">
          {{ submitting ? '验证中…' : '下一步，修改密码' }}
        </button>

        <router-link to="/login" class="btn btn-ghost back-login-btn">返回登录</router-link>
      </form>

      <router-link to="/" class="forgot-back">返回首页</router-link>
    </div>
  </section>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUser } from '../composables/useUser'
import { showToast } from '../utils/toast'

const router = useRouter()
const {
  requestResetCode,
  verifyResetCode,
} = useUser()

const email = ref('')
const code = ref('')
const sendingCode = ref(false)
const countdown = ref(0)
const submitting = ref(false)
const errors = reactive({ email: '', code: '' })

let countdownTimer = null

function clearCountdown() {
  if (countdownTimer) {
    clearInterval(countdownTimer)
    countdownTimer = null
  }
}

function sendCode() {
  errors.email = ''
  errors.code = ''
  const em = email.value.trim()
  if (!em) {
    errors.email = '请输入邮箱'
    return
  }
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(em)) {
    errors.email = '请输入有效的邮箱地址'
    return
  }
  sendingCode.value = true
  const result = requestResetCode(em)
  sendingCode.value = false
  if (result.ok) {
    showToast('验证码已发送，请查收')
    if (import.meta.env?.DEV && result.code) {
      console.log('[找回密码 演示] 验证码:', result.code)
    }
    countdown.value = 60
    countdownTimer = setInterval(() => {
      countdown.value -= 1
      if (countdown.value <= 0) clearCountdown()
    }, 1000)
  } else {
    showToast(result.message || '发送失败')
  }
}

function handleSubmit() {
  errors.email = ''
  errors.code = ''
  const em = email.value.trim()
  const c = code.value.trim()
  if (!em) {
    errors.email = '请输入邮箱'
    return
  }
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(em)) {
    errors.email = '请输入有效的邮箱地址'
    return
  }
  if (!c) {
    errors.code = '请输入验证码'
    return
  }
  submitting.value = true
  try {
    const result = verifyResetCode(em, c)
    if (result.ok) {
      showToast('验证成功')
      router.replace({ path: '/reset-password', query: { email: em } })
    } else {
      errors.code = result.message || '验证失败'
      showToast(result.message || '验证失败')
    }
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.view-forgot {
  min-height: 60vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 0;
}

.forgot-card {
  width: 100%;
  max-width: 400px;
  padding: 32px 24px;
}

.forgot-title {
  font-family: var(--font-serif);
  font-weight: 600;
  font-size: 24px;
  margin: 0 0 8px;
  color: var(--color-text);
}

.forgot-desc {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin: 0 0 24px;
}

.forgot-form {
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

.form-input-error {
  border-color: var(--color-error);
}

.form-row {
  display: flex;
  gap: 12px;
  align-items: stretch;
}

.form-input-code {
  flex: 1;
  min-width: 0;
}

.send-code-btn {
  flex-shrink: 0;
  white-space: nowrap;
}

.send-code-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.form-error {
  font-size: 12px;
  color: var(--color-error);
  margin: 0;
}

.submit-btn {
  width: 100%;
  padding: 12px 20px;
  font-size: 16px;
  margin-top: 8px;
}

.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.back-login-btn {
  width: 100%;
  padding: 12px 20px;
  font-size: 15px;
  text-decoration: none;
  text-align: center;
}

.forgot-back {
  display: block;
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: var(--color-text-secondary);
  text-decoration: none;
}

.forgot-back:hover {
  color: var(--color-accent);
}
</style>
