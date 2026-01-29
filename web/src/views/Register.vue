<template>
  <section class="view-register">
    <div class="register-card card">
      <h1 class="register-title">注册</h1>
      <p class="register-desc">填写以下信息完成注册</p>

      <form class="register-form" @submit.prevent="handleSubmit">
        <div class="form-group">
          <label for="username" class="form-label">用户名</label>
          <input
            id="username"
            v-model="username"
            type="text"
            class="form-input"
            :class="{ 'form-input-error': errors.username }"
            placeholder="请输入用户名"
            maxlength="15"
            autocomplete="username"
            required
            @focus="usernameFocused = true"
            @blur="onUsernameBlur"
          />
          <p v-if="usernameFocused && !errors.username" class="form-hint">用户名由3到15个字符组成</p>
          <p v-if="errors.username" class="form-error">{{ errors.username }}</p>
        </div>

        <div class="form-group">
          <label for="password" class="form-label">密码</label>
          <input
            id="password"
            v-model="password"
            type="password"
            class="form-input"
            :class="{ 'form-input-error': errors.password }"
            placeholder="请输入密码"
            maxlength="50"
            autocomplete="new-password"
            required
            @focus="passwordFocused = true"
            @blur="onPasswordBlur"
          />
          <p v-if="passwordFocused && !errors.password" class="form-hint">密码长度在8到50个字符</p>
          <p v-if="errors.password" class="form-error">{{ errors.password }}</p>
        </div>

        <div class="form-group">
          <label for="confirmPassword" class="form-label">确认密码</label>
          <input
            id="confirmPassword"
            v-model="confirmPassword"
            type="password"
            class="form-input"
            :class="{ 'form-input-error': errors.confirmPassword }"
            placeholder="请再次输入密码"
            maxlength="50"
            autocomplete="new-password"
            required
            @focus="confirmPasswordFocused = true"
            @blur="onConfirmPasswordBlur"
          />
          <p v-if="confirmPasswordFocused && !errors.confirmPassword" class="form-hint">请再次输入密码</p>
          <p v-if="errors.confirmPassword" class="form-error">{{ errors.confirmPassword }}</p>
        </div>

        <div class="form-group">
          <label for="email" class="form-label">邮箱</label>
          <input
            id="email"
            v-model="email"
            type="email"
            class="form-input"
            :class="{ 'form-input-error': errors.email }"
            placeholder="请输入邮箱"
            autocomplete="email"
            required
          />
          <p v-if="errors.email" class="form-error">{{ errors.email }}</p>
        </div>

        <label class="form-checkbox">
          <input v-model="agreeTerms" type="checkbox" class="form-checkbox-input" />
          <span class="form-checkbox-text">同意<router-link to="/terms" class="terms-link" target="_blank">网站服务条款</router-link></span>
        </label>
        <p v-if="errors.agreeTerms" class="form-error">{{ errors.agreeTerms }}</p>

        <button type="submit" class="btn btn-primary register-btn" :disabled="submitting">
          {{ submitting ? '注册中…' : '注册' }}
        </button>

        <router-link to="/login" class="btn btn-ghost to-login-btn">去登录</router-link>
      </form>

      <router-link to="/" class="register-back">返回首页</router-link>
    </div>
  </section>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUser } from '../composables/useUser'
import { showToast } from '../utils/toast'

const router = useRouter()
const { register } = useUser()

const username = ref('')
const usernameFocused = ref(false)
const passwordFocused = ref(false)
const confirmPasswordFocused = ref(false)
const password = ref('')
const confirmPassword = ref('')
const email = ref('')
const agreeTerms = ref(false)
const submitting = ref(false)
const errors = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  email: '',
  agreeTerms: '',
})

function validateUsername() {
  const name = username.value.trim()
  if (!name) {
    errors.username = '请输入用户名'
    return false
  }
  if (name.length < 3 || name.length > 15) {
    errors.username = '用户名由3到15个字符组成'
    return false
  }
  errors.username = ''
  return true
}

function onUsernameBlur() {
  usernameFocused.value = false
  if (username.value.trim()) validateUsername()
}

function validatePassword() {
  const pwd = password.value
  if (!pwd) {
    errors.password = '请输入密码'
    return false
  }
  if (pwd.length < 8 || pwd.length > 50) {
    errors.password = '密码长度在8到50个字符'
    return false
  }
  errors.password = ''
  return true
}

function validateConfirmPassword() {
  const pwd = password.value
  const confirm = confirmPassword.value
  if (!confirm) {
    errors.confirmPassword = '请再次输入密码'
    return false
  }
  if (pwd !== confirm) {
    errors.confirmPassword = '两次输入的密码不一致'
    return false
  }
  errors.confirmPassword = ''
  return true
}

function onPasswordBlur() {
  passwordFocused.value = false
  if (password.value) validatePassword()
}

function onConfirmPasswordBlur() {
  confirmPasswordFocused.value = false
  if (confirmPassword.value || password.value) validateConfirmPassword()
}

function validate() {
  let valid = true
  errors.username = ''
  errors.password = ''
  errors.confirmPassword = ''
  errors.email = ''
  errors.agreeTerms = ''

  const name = username.value.trim()
  if (!name) {
    errors.username = '请输入用户名'
    valid = false
  } else if (name.length < 3 || name.length > 15) {
    errors.username = '用户名由3到15个字符组成'
    valid = false
  }

  const pwd = password.value
  if (!pwd) {
    errors.password = '请输入密码'
    valid = false
  } else if (pwd.length < 8 || pwd.length > 50) {
    errors.password = '密码长度在8到50个字符'
    valid = false
  }

  const confirm = confirmPassword.value
  if (!confirm) {
    errors.confirmPassword = '请再次输入密码'
    valid = false
  } else if (pwd !== confirm) {
    errors.confirmPassword = '两次输入的密码不一致'
    valid = false
  }

  const em = email.value.trim()
  if (!em) {
    errors.email = '请输入邮箱'
    valid = false
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(em)) {
    errors.email = '请输入有效的邮箱地址'
    valid = false
  }

  if (!agreeTerms.value) {
    errors.agreeTerms = '请先同意网站服务条款'
    valid = false
  }

  return valid
}

async function handleSubmit() {
  if (!validate()) return

  submitting.value = true
  try {
    const result = await register(username.value.trim(), password.value, email.value.trim())
    if (result.ok) {
      showToast('注册成功，请登录')
      router.replace('/login')
    } else {
      showToast(result.message || '注册失败')
    }
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.view-register {
  min-height: 60vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 0;
}

.register-card {
  width: 100%;
  max-width: 400px;
  padding: 32px 24px;
}

.register-title {
  font-family: var(--font-serif);
  font-weight: 600;
  font-size: 24px;
  margin: 0 0 8px;
  color: var(--color-text);
}

.register-desc {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin: 0 0 24px;
}

.register-form {
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

.form-input-error:focus {
  border-color: var(--color-error);
}

.form-hint {
  font-size: 12px;
  color: var(--color-text-secondary);
  margin: 0;
}

.form-error {
  font-size: 12px;
  color: var(--color-error);
  margin: 0;
}

.form-checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
  color: var(--color-text);
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

.terms-link {
  color: var(--color-accent);
  text-decoration: none;
  margin: 0 2px;
}

.terms-link:hover {
  text-decoration: underline;
}

.register-btn {
  width: 100%;
  padding: 12px 20px;
  font-size: 16px;
  margin-top: 8px;
}

.register-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.to-login-btn {
  width: 100%;
  padding: 12px 20px;
  font-size: 15px;
  text-decoration: none;
  text-align: center;
}

.register-back {
  display: block;
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: var(--color-text-secondary);
  text-decoration: none;
}

.register-back:hover {
  color: var(--color-accent);
}
</style>
