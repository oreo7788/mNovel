/**
 * 上传功能 - 先缓存图片到本地（IndexedDB），用于测试上传流程
 */
(function () {
  const DB_NAME = 'ydk_upload';
  const STORE_CACHE = 'image_cache';
  const MAX_FILES = 9;

  let db = null;

  function openDB() {
    return new Promise(function (resolve, reject) {
      if (db) {
        resolve(db);
        return;
      }
      const req = indexedDB.open(DB_NAME, 1);
      req.onerror = function () { reject(req.error); };
      req.onsuccess = function () {
        db = req.result;
        resolve(db);
      };
      req.onupgradeneeded = function (e) {
        const database = e.target.result;
        if (!database.objectStoreNames.contains(STORE_CACHE)) {
          database.createObjectStore(STORE_CACHE, { keyPath: 'id' });
        }
      };
    });
  }

  function generateId() {
    return 'img_' + Date.now() + '_' + Math.random().toString(36).slice(2, 9);
  }

  /**
   * 将选中的图片写入 IndexedDB 本地缓存
   * @param {FileList} files
   * @returns {Promise<{ count: number, ids: string[] }>}
   */
  function cacheImagesToLocal(files) {
    return openDB().then(function (database) {
      const store = database.transaction(STORE_CACHE, 'readwrite').objectStore(STORE_CACHE);
      const list = Array.from(files).slice(0, MAX_FILES);
      const ids = [];

      return new Promise(function (resolve, reject) {
        let done = 0;
        const total = list.length;
        if (total === 0) {
          resolve({ count: 0, ids: [] });
          return;
        }

        list.forEach(function (file) {
          if (!file.type || !file.type.startsWith('image/')) return;
          const id = generateId();
          ids.push(id);
          const record = {
            id: id,
            name: file.name,
            size: file.size,
            type: file.type,
            blob: file,
            createdAt: Date.now()
          };
          const req = store.put(record);
          req.onsuccess = function () {
            done++;
            if (done === total) resolve({ count: total, ids: ids });
          };
          req.onerror = function () { reject(req.error); };
        });
      });
    });
  }

  /**
   * 获取本地已缓存的图片列表
   */
  function getCachedImages() {
    return openDB().then(function (database) {
      return new Promise(function (resolve, reject) {
        const tx = database.transaction(STORE_CACHE, 'readonly');
        const store = tx.objectStore(STORE_CACHE);
        const req = store.getAll();
        req.onsuccess = function () { resolve(req.result || []); };
        req.onerror = function () { reject(req.error); };
      });
    });
  }

  /**
   * 删除单条缓存
   */
  function removeCachedImage(id) {
    return openDB().then(function (database) {
      return new Promise(function (resolve, reject) {
        const tx = database.transaction(STORE_CACHE, 'readwrite');
        const store = tx.objectStore(STORE_CACHE);
        const req = store.delete(id);
        req.onsuccess = function () { resolve(); };
        req.onerror = function () { reject(req.error); };
      });
    });
  }

  /**
   * 清空本地图片缓存
   */
  function clearImageCache() {
    return openDB().then(function (database) {
      return new Promise(function (resolve, reject) {
        const tx = database.transaction(STORE_CACHE, 'readwrite');
        const store = tx.objectStore(STORE_CACHE);
        const req = store.clear();
        req.onsuccess = function () { resolve(); };
        req.onerror = function () { reject(req.error); };
      });
    });
  }

  function showToast(message) {
    var el = document.createElement('div');
    el.className = 'upload-toast';
    el.textContent = message;
    el.setAttribute('role', 'status');
    document.body.appendChild(el);
    setTimeout(function () {
      if (el.parentNode) el.parentNode.removeChild(el);
    }, 2500);
  }

  function renderCachedList(items) {
    var container = document.getElementById('upload-cached-list');
    if (!container) return;
    if (!items || items.length === 0) {
      container.innerHTML = '';
      container.hidden = true;
      return;
    }
    container.hidden = false;
    var title = '<div class="cached-list-title">已缓存到本地（测试用）</div>';
    container.innerHTML = title + items.map(function (item) {
      var url = item.blob ? URL.createObjectURL(item.blob) : '';
      var size = item.size ? (item.size / 1024).toFixed(1) + ' KB' : '';
      return (
        '<div class="cached-item" data-id="' + item.id + '">' +
        '<img src="' + url + '" alt="" class="cached-thumb"/>' +
        '<div class="cached-info">' +
        '<span class="cached-name">' + (item.name || '图片') + '</span>' +
        '<span class="cached-meta">' + size + '</span>' +
        '</div>' +
        '<button type="button" class="btn btn-ghost cached-remove" title="从缓存移除">移除</button>' +
        '</div>'
      );
    }).join('');

    container.querySelectorAll('.cached-remove').forEach(function (btn) {
      btn.addEventListener('click', function () {
        var id = btn.closest('.cached-item').getAttribute('data-id');
        removeCachedImage(id).then(function () {
          refreshCachedList();
        });
      });
    });
  }

  function refreshCachedList() {
    getCachedImages().then(renderCachedList);
  }

  function initUploadZone() {
    var zone = document.querySelector('.upload-zone');
    var input = document.getElementById('upload-file-input');
    var progressWrap = document.querySelector('.upload-progress');
    var progressFill = progressWrap ? progressWrap.querySelector('.progress-fill') : null;
    var progressText = progressWrap ? progressWrap.querySelector('.progress-text') : null;

    if (!zone || !input) return;

    zone.addEventListener('click', function () {
      input.click();
    });

    zone.addEventListener('dragover', function (e) {
      e.preventDefault();
      e.stopPropagation();
      zone.classList.add('upload-zone-dragover');
    });
    zone.addEventListener('dragleave', function (e) {
      e.preventDefault();
      e.stopPropagation();
      zone.classList.remove('upload-zone-dragover');
    });
    zone.addEventListener('drop', function (e) {
      e.preventDefault();
      e.stopPropagation();
      zone.classList.remove('upload-zone-dragover');
      var files = e.dataTransfer && e.dataTransfer.files;
      if (files && files.length) handleFiles(files, progressWrap, progressFill, progressText);
    });

    input.addEventListener('change', function () {
      var files = input.files;
      if (files && files.length) handleFiles(files, progressWrap, progressFill, progressText);
      input.value = '';
    });
  }

  function handleFiles(files, progressWrap, progressFill, progressText) {
    var list = Array.from(files).filter(function (f) {
      return f.type && f.type.startsWith('image/');
    }).slice(0, MAX_FILES);

    if (list.length === 0) {
      showToast('请选择图片文件');
      return;
    }
    if (list.length < files.length) {
      showToast('已忽略非图片或超出 ' + MAX_FILES + ' 张的部分');
    }

    if (progressWrap) progressWrap.hidden = false;
    if (progressText) progressText.textContent = '正在缓存到本地… 0%';
    if (progressFill) progressFill.style.width = '0%';

    cacheImagesToLocal(list).then(function (res) {
      if (progressFill) progressFill.style.width = '100%';
      if (progressText) progressText.textContent = '已缓存 ' + res.count + ' 张图片到本地';
      showToast('已缓存 ' + res.count + ' 张图片到本地');
      refreshCachedList();
      refreshHomeRecentUploads();
      refreshWardrobeGrid();
    }).catch(function (err) {
      if (progressText) progressText.textContent = '缓存失败：' + (err && err.message ? err.message : '未知错误');
      showToast('缓存失败，请重试');
    });
  }

  /**
   * 首页「最近上传」区块：展示已缓存照片，便于在首页看到上传结果
   */
  function refreshHomeRecentUploads() {
    var block = document.getElementById('home-recent-uploads');
    var listEl = document.getElementById('home-recent-uploads-list');
    if (!block || !listEl) return;
    getCachedImages().then(function (items) {
      if (!items || items.length === 0) {
        block.hidden = true;
        return;
      }
      block.hidden = false;
      var recent = items.slice(-3);
      listEl.innerHTML = recent.map(function (item) {
        var url = item.blob ? URL.createObjectURL(item.blob) : '';
        return '<div class="recent-uploads-item"><img src="' + url + '" alt="" class="recent-uploads-thumb"/></div>';
      }).join('');
    });
  }

  /**
   * 衣橱页网格：用已缓存照片动态渲染，这样在衣橱也能看到上传的图
   */
  function refreshWardrobeGrid() {
    var grid = document.getElementById('wardrobe-grid');
    if (!grid) return;
    getCachedImages().then(function (items) {
      var list = items || [];
      var addCardHtml = '<article class="clothing-card card wardrobe-add-card" data-action="upload">' +
        '<div class="thumb"><div class="placeholder">+</div></div>' +
        '<div class="info"><div class="category">添加衣物</div></div></article>';
      var cardsHtml = list.map(function (item) {
        var url = item.blob ? URL.createObjectURL(item.blob) : '';
        var name = (item.name || '图片').replace(/\.[^.]+$/, '');
        return '<article class="clothing-card card" data-id="' + item.id + '">' +
          '<div class="thumb"><img src="' + url + '" alt="" class="wardrobe-thumb-img"/></div>' +
          '<div class="info">' +
          '<div class="category">' + escapeHtml(name) + '</div>' +
          '</div></article>';
      }).join('');
      grid.innerHTML = cardsHtml + addCardHtml;
      var addCard = grid.querySelector('.wardrobe-add-card');
      if (addCard) {
        addCard.addEventListener('click', function () {
          var navUpload = document.querySelector('[data-nav="upload"]');
          if (navUpload) navUpload.click();
        });
      }
    });
  }

  function escapeHtml(s) {
    var div = document.createElement('div');
    div.textContent = s;
    return div.innerHTML;
  }

  function init() {
    initUploadZone();
    refreshCachedList();
    refreshHomeRecentUploads();
    refreshWardrobeGrid();
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }

  window.uploadCache = {
    getCachedImages: getCachedImages,
    clearImageCache: clearImageCache,
    removeCachedImage: removeCachedImage,
    refreshCachedList: refreshCachedList,
    refreshHomeRecentUploads: refreshHomeRecentUploads,
    refreshWardrobeGrid: refreshWardrobeGrid
  };
})();
