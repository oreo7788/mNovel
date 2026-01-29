import { ref, computed } from 'vue'
import { registerApi, loginApi } from '../api'

const STORAGE_KEY = 'ydk_user'
const USERS_KEY = 'ydk_users'
const RESET_CODE_KEY = 'ydk_reset_code'
const RESET_ALLOWED_KEY = 'ydk_reset_allowed'
const CODE_EXPIRY_MS = 5 * 60 * 1000
const RESET_ALLOWED_MS = 15 * 60 * 1000

function loadUser() {
  try {
    const raw = sessionStorage.getItem(STORAGE_KEY) ?? localStorage.getItem(STORAGE_KEY)
    return raw ? JSON.parse(raw) : null
  } catch {
    return null
  }
}

function loadUsers() {
  try {
    const raw = localStorage.getItem(USERS_KEY)
    return raw ? JSON.parse(raw) : {}
  } catch {
    return {}
  }
}

const user = ref(loadUser())

export function useUser() {
  const isLoggedIn = computed(() => !!user.value)
  const avatar = computed(() => user.value?.avatar ?? '')
  const nickname = computed(() => user.value?.nickname ?? '')

  /** 用户名/邮箱/手机号+密码登录，调用后端接口；remember 为 true 时持久化到 localStorage，否则仅 sessionStorage */
  async function login(account, password, remember = true) {
    const result = await loginApi((account || '').trim(), password || '')
    if (!result.ok) return false
    const data = result.data || {}
    const u = data.user || {}
    const letter = (u.nickname || u.email || '用').toString().slice(0, 1)
    user.value = {
      id: u.id,
      nickname: u.nickname || '',
      email: u.email || '',
      avatar: u.avatar || defaultAvatar(letter),
      access_token: data.access_token,
    }
    try {
      const json = JSON.stringify(user.value)
      if (remember) {
        localStorage.setItem(STORAGE_KEY, json)
        sessionStorage.removeItem(STORAGE_KEY)
      } else {
        sessionStorage.setItem(STORAGE_KEY, json)
        localStorage.removeItem(STORAGE_KEY)
      }
    } catch (_) {}
    return true
  }

  /** 注册：调用后端接口，传昵称(用户名)、密码、邮箱 */
  async function register(username, password, email) {
    const result = await registerApi(
      (username || '').trim(),
      password || '',
      (email || '').trim()
    )
    return result
  }

  function logout() {
    user.value = null
    try {
      localStorage.removeItem(STORAGE_KEY)
      sessionStorage.removeItem(STORAGE_KEY)
    } catch (_) {}
  }

  /** 根据邮箱查找对应用户名（本地用户表） */
  function getUsernameByEmail(email) {
    const users = loadUsers()
    const em = (email || '').trim().toLowerCase()
    for (const [name, data] of Object.entries(users)) {
      if ((data.email || '').trim().toLowerCase() === em) return name
    }
    return null
  }

  /** 发送找回密码验证码（演示：生成 6 位码存 sessionStorage） */
  function requestResetCode(email) {
    const em = (email || '').trim()
    if (!em) return { ok: false, message: '请输入邮箱' }
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(em)) return { ok: false, message: '请输入有效的邮箱地址' }
    const username = getUsernameByEmail(em)
    if (!username) return { ok: false, message: '该邮箱未注册' }
    const code = String(Math.floor(100000 + Math.random() * 900000))
    const payload = { email: em, code, expiry: Date.now() + CODE_EXPIRY_MS }
    try {
      sessionStorage.setItem(RESET_CODE_KEY, JSON.stringify(payload))
    } catch (_) {
      return { ok: false, message: '发送失败' }
    }
    return { ok: true, code }
  }

  /** 校验验证码，通过则允许进入重置密码页 */
  function verifyResetCode(email, code) {
    const em = (email || '').trim()
    const c = (code || '').trim()
    if (!em || !c) return { ok: false, message: '请填写邮箱和验证码' }
    try {
      const raw = sessionStorage.getItem(RESET_CODE_KEY)
      if (!raw) return { ok: false, message: '验证码已过期，请重新获取' }
      const { email: storedEmail, code: storedCode, expiry } = JSON.parse(raw)
      if (storedEmail !== em) return { ok: false, message: '邮箱与获取验证码时不一致' }
      if (Date.now() > expiry) return { ok: false, message: '验证码已过期，请重新获取' }
      if (storedCode !== c) return { ok: false, message: '验证码错误' }
      sessionStorage.removeItem(RESET_CODE_KEY)
      sessionStorage.setItem(RESET_ALLOWED_KEY, JSON.stringify({ email: em, until: Date.now() + RESET_ALLOWED_MS }))
    } catch (_) {
      return { ok: false, message: '验证失败' }
    }
    return { ok: true }
  }

  /** 检查当前是否允许重置密码（用于重置页进入校验） */
  function checkResetAllowed() {
    try {
      const raw = sessionStorage.getItem(RESET_ALLOWED_KEY)
      if (!raw) return null
      const { email, until } = JSON.parse(raw)
      if (Date.now() > until) {
        sessionStorage.removeItem(RESET_ALLOWED_KEY)
        return null
      }
      return email
    } catch (_) {
      return null
    }
  }

  /** 重置密码（需先通过验证码校验） */
  function resetPassword(newPassword) {
    const allowed = checkResetAllowed()
    if (!allowed) return { ok: false, message: '请先完成邮箱验证' }
    const pwd = (newPassword || '').trim()
    if (pwd.length < 8 || pwd.length > 50) return { ok: false, message: '密码长度在8到50个字符' }
    const users = loadUsers()
    const username = getUsernameByEmail(allowed)
    if (!username) return { ok: false, message: '用户不存在' }
    users[username].password = pwd
    try {
      localStorage.setItem(USERS_KEY, JSON.stringify(users))
      sessionStorage.removeItem(RESET_ALLOWED_KEY)
    } catch (_) {
      return { ok: false, message: '保存失败' }
    }
    return { ok: true }
  }

  return {
    user,
    isLoggedIn,
    avatar,
    nickname,
    login,
    register,
    logout,
    getUsernameByEmail,
    requestResetCode,
    verifyResetCode,
    checkResetAllowed,
    resetPassword,
  }
}

function defaultAvatar(letter) {
  const colors = ['C4A574', '5B7C9D', '4A7C59', 'B85450', '8E8C89']
  const c = colors[Math.floor(Math.random() * colors.length)]
  return `data:image/svg+xml,${encodeURIComponent(
    `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64"><rect width="64" height="64" fill="#${c}"/><text x="32" y="40" font-size="28" fill="white" text-anchor="middle" font-family="sans-serif">${letter}</text></svg>`
  )}`
}
