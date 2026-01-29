import { ref, computed } from 'vue'

const STORAGE_KEY = 'ydk_user'

function loadUser() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    return raw ? JSON.parse(raw) : null
  } catch {
    return null
  }
}

const user = ref(loadUser())

export function useUser() {
  const isLoggedIn = computed(() => !!user.value)
  const avatar = computed(() => user.value?.avatar ?? '')
  const nickname = computed(() => user.value?.nickname ?? '')

  function login(payload = {}) {
    const { nickname: n = '用户', avatar: a = '' } = payload
    const letter = (n || '用').slice(0, 1)
    user.value = { nickname: n, avatar: a || defaultAvatar(letter) }
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(user.value))
    } catch (_) {}
  }

  function logout() {
    user.value = null
    try {
      localStorage.removeItem(STORAGE_KEY)
    } catch (_) {}
  }

  return { user, isLoggedIn, avatar, nickname, login, logout }
}

function defaultAvatar(letter) {
  const colors = ['C4A574', '5B7C9D', '4A7C59', 'B85450', '8E8C89']
  const c = colors[Math.floor(Math.random() * colors.length)]
  return `data:image/svg+xml,${encodeURIComponent(
    `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64"><rect width="64" height="64" fill="#${c}"/><text x="32" y="40" font-size="28" fill="white" text-anchor="middle" font-family="sans-serif">${letter}</text></svg>`
  )}`
}
