<template>
  <section class="view-reset">
    <div class="reset-card card">
      <h1 class="reset-title">修改密码</h1>
      <p v-if="resetEmail" class="reset-desc">正在为 {{ resetEmail }} 设置新密码</p>
      <p v-else class="reset-desc">请设置新密码</p>

      <form class="reset-form" @submit.prevent="handleSubmit">
        <div class="form-group">
          <label for="password" class="form-label">新密码</label>
          <input
            id="password"
            v-model="password"
            type="password"
            class="form-input"
            :class="{ 'form-input-error': errors.password }"
            placeholder="请输入新密码（8-50 位）"
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
            placeholder="请再次输入新密码"
            maxlength="50"
            autocomplete="new-password"
            required
            @focus="confirmPasswordFocused = true"
            @blur="onConfirmPasswordBlur"
          />
          <p v-if="confirmPasswordFocused && !errors.confirmPassword" class="form-hint">请再次输入密码</p>
          <p v-if="errors.confirmPassword" class="form-error">{{ errors.confirmPassword }}</p>
        </div>

        <button type="submit" class="btn btn-primary submit-btn" :disabled="submitting">
          {{ submitting ? '提交中…' : '完成修改' }}
        </button>

        <router-link to="/login" class="btn btn-ghost back-login-btn">返回登录</router-link>
      </form>

      <router-link to="/" class="reset-back">返回首页</router-link>
    </div>
  </section>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUser } from '../composables/useUser'
import { showToast } from '../utils/toast'

const router = useRouter()
const route = useRoute()
const { checkResetAllowed, resetPassword } = useUser()

const resetEmail = ref('')
const password = ref('')
const confirmPassword = ref('')
const passwordFocused = ref(false)
const confirmPasswordFocused = ref(false)
const submitting = ref(false)
const errors = reactive({ password: '', confirmPassword: '' })

onMounted(() => {
  const allowed = checkResetAllowed()
  if (!allowed) {
    showToast('请先完成邮箱验证')
    router.replace('/forgot-password')
    return
  }
  resetEmail.value = route.query.email || allowed
})

function validatePassword() {
  const pwd = password.value
  if (!pwd) {
    errors.password = '请输入新密码'
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

function handleSubmit() {
  errors.password = ''
  errors.confirmPassword = ''
  if (!validatePassword() || !validateConfirmPassword()) return

  submitting.value = true
  try {
    const result = resetPassword(password.value)
    if (result.ok) {
      showToast('密码已修改，请使用新密码登录')
      router.replace('/login')
    } else {
      showToast(result.message || '修改失败')
    }
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.view-reset {
  min-height: 60vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 0;
}

.reset-card {
  width: 100%;
  max-width: 400px;
  padding: 32px 24px;
}

.reset-title {
  font-family: var(--font-serif);
  font-weight: 600;
  font-size: 24px;
  margin: 0 0 8px;
  color: var(--color-text);
}

.reset-desc {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin: 0 0 24px;
}

.reset-form {
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

.reset-back {
  display: block;
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: var(--color-text-secondary);
  text-decoration: none;
}

.reset-back:hover {
  color: var(--color-accent);
}
</style>
