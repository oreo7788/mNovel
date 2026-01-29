/**
 * 上传缓存 - IndexedDB 本地图片缓存，用于测试上传流程
 */
const DB_NAME = 'ydk_upload'
const STORE_CACHE = 'image_cache'
const MAX_FILES = 9

let db = null

function openDB() {
  return new Promise((resolve, reject) => {
    if (db) {
      resolve(db)
      return
    }
    const req = indexedDB.open(DB_NAME, 1)
    req.onerror = () => reject(req.error)
    req.onsuccess = () => {
      db = req.result
      resolve(db)
    }
    req.onupgradeneeded = (e) => {
      const database = e.target.result
      if (!database.objectStoreNames.contains(STORE_CACHE)) {
        database.createObjectStore(STORE_CACHE, { keyPath: 'id' })
      }
    }
  })
}

function generateId() {
  return 'img_' + Date.now() + '_' + Math.random().toString(36).slice(2, 9)
}

/**
 * 将选中的图片写入 IndexedDB 本地缓存
 */
export function cacheImagesToLocal(files) {
  return openDB().then((database) => {
    const store = database.transaction(STORE_CACHE, 'readwrite').objectStore(STORE_CACHE)
    const list = Array.from(files).slice(0, MAX_FILES)
    const ids = []

    return new Promise((resolve, reject) => {
      let done = 0
      const total = list.length
      if (total === 0) {
        resolve({ count: 0, ids: [] })
        return
      }

      list.forEach((file) => {
        if (!file.type || !file.type.startsWith('image/')) return
        const id = generateId()
        ids.push(id)
        const record = {
          id,
          name: file.name,
          size: file.size,
          type: file.type,
          blob: file,
          createdAt: Date.now(),
        }
        const req = store.put(record)
        req.onsuccess = () => {
          done++
          if (done === total) resolve({ count: total, ids })
        }
        req.onerror = () => reject(req.error)
      })
    })
  })
}

/**
 * 获取本地已缓存的图片列表
 */
export function getCachedImages() {
  return openDB().then((database) => {
    return new Promise((resolve, reject) => {
      const tx = database.transaction(STORE_CACHE, 'readonly')
      const store = tx.objectStore(STORE_CACHE)
      const req = store.getAll()
      req.onsuccess = () => resolve(req.result || [])
      req.onerror = () => reject(req.error)
    })
  })
}

/**
 * 删除单条缓存
 */
export function removeCachedImage(id) {
  return openDB().then((database) => {
    return new Promise((resolve, reject) => {
      const tx = database.transaction(STORE_CACHE, 'readwrite')
      const store = tx.objectStore(STORE_CACHE)
      const req = store.delete(id)
      req.onsuccess = () => resolve()
      req.onerror = () => reject(req.error)
    })
  })
}

/**
 * 清空本地图片缓存
 */
export function clearImageCache() {
  return openDB().then((database) => {
    return new Promise((resolve, reject) => {
      const tx = database.transaction(STORE_CACHE, 'readwrite')
      const store = tx.objectStore(STORE_CACHE)
      const req = store.clear()
      req.onsuccess = () => resolve()
      req.onerror = () => reject(req.error)
    })
  })
}

export const MAX_UPLOAD_FILES = MAX_FILES
