/**
 * 后端 API 请求封装
 * 环境变量 VITE_API_BASE_URL 指定后端地址，默认 http://localhost:8080
 */

function getBaseUrl() {
  return import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
}

/**
 * 通用请求：返回 { ok, data, message }，与后端 response.Response 对应
 * @param {string} path - 接口路径，如 /api/v1/auth/register
 * @param {RequestInit} options - fetch 选项
 */
export async function request(path, options = {}) {
  const url = getBaseUrl() + path
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  }
  try {
    const res = await fetch(url, { ...options, headers })
    const body = await res.json().catch(() => ({}))
    const code = body.code
    const message = body.message ?? '请求失败'
    const data = body.data

    if (res.ok && code === 0) {
      return { ok: true, data, message }
    }
    return { ok: false, message, data }
  } catch (err) {
    return { ok: false, message: err.message || '网络错误' }
  }
}

/**
 * 用户注册
 * @param {string} nickname - 昵称（前端传用户名）
 * @param {string} password - 密码
 * @param {string} email - 邮箱
 * @returns {Promise<{ok: boolean, data?: object, message?: string}>}
 */
export function registerApi(nickname, password, email) {
  return request('/api/v1/auth/register', {
    method: 'POST',
    body: JSON.stringify({
      nickname: nickname || '',
      password: password || '',
      email: (email || '').trim(),
    }),
  })
}

/**
 * 用户登录（支持用户名/昵称、邮箱或手机号）
 * @param {string} account - 用户名(昵称)、邮箱或手机号
 * @param {string} password - 密码
 * @returns {Promise<{ok: boolean, data?: object, message?: string}>}
 */
export function loginApi(account, password) {
  const value = (account || '').trim()
  const body = { password: password || '' }
  if (/^1[3-9]\d{9}$/.test(value)) {
    body.phone = value
  } else if (/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
    body.email = value
  } else {
    body.nickname = value
  }
  return request('/api/v1/auth/login', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

/**
 * 带 Token 的请求（用于需登录的接口）
 * @param {string} path - 接口路径
 * @param {RequestInit} options - fetch 选项（可含 body、method 等）
 * @param {string} accessToken - JWT access_token
 */
export async function requestWithToken(path, options = {}, accessToken) {
  const url = (import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080') + path
  const headers = {
    ...options.headers,
  }
  if (typeof accessToken === 'string' && accessToken.length > 0) {
    headers['Authorization'] = 'Bearer ' + accessToken
  }
  try {
    const res = await fetch(url, { ...options, headers })
    const body = await res.json().catch(() => ({}))
    const code = body.code
    const message = body.message ?? '请求失败'
    const data = body.data
    if (res.ok && code === 0) {
      return { ok: true, data, message }
    }
    return { ok: false, message, data }
  } catch (err) {
    return { ok: false, message: err.message || '网络错误' }
  }
}

/**
 * 简单上传图片到七牛（后端转发）
 * @param {File} file - 图片文件
 * @param {string} accessToken - JWT access_token
 * @returns {Promise<{ok: boolean, data?: { url: string, filename: string }, message?: string}>}
 */
export function uploadSimpleApi(file, accessToken) {
  const form = new FormData()
  form.append('file', file)
  return requestWithToken('/api/v1/upload/simple', {
    method: 'POST',
    body: form,
    headers: {}, // 不设 Content-Type，让浏览器自动带 multipart boundary
  }, accessToken)
}
