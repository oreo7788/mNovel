<template>
  <header class="app-header">
    <div class="app-header-inner">
      <router-link to="/" class="logo">衣搭库</router-link>
      <div class="header-actions">
        <router-link v-if="!isLoggedIn" to="/profile" class="header-login">登录/注册</router-link>
        <router-link v-else to="/profile" class="header-avatar-link" aria-label="个人中心">
          <img :src="avatar" alt="" class="header-avatar" />
        </router-link>
        <div class="header-dropdown" ref="dropdownRef">
          <button
            type="button"
            class="header-action-trigger"
            :class="{ open: dropdownOpen }"
            aria-haspopup="true"
            :aria-expanded="dropdownOpen"
            @click="dropdownOpen = !dropdownOpen"
          >
            操作
          </button>
          <Transition name="dropdown">
            <div v-show="dropdownOpen" class="header-dropdown-menu" role="menu">
              <router-link to="/wardrobe" class="header-dropdown-item" role="menuitem" @click="closeDropdown">
                👔 衣橱
              </router-link>
              <router-link to="/recommend" class="header-dropdown-item" role="menuitem" @click="closeDropdown">
                ✨ 推荐
              </router-link>
              <router-link to="/favorites" class="header-dropdown-item" role="menuitem" @click="closeDropdown">
                ♥ 收藏
              </router-link>
              <router-link to="/upload" class="header-dropdown-item" role="menuitem" @click="closeDropdown">
                ➕ 上传
              </router-link>
            </div>
          </Transition>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useUser } from '../composables/useUser'

const { isLoggedIn, avatar } = useUser()
const dropdownOpen = ref(false)
const dropdownRef = ref(null)

function closeDropdown() {
  dropdownOpen.value = false
}

function onClickOutside(e) {
  if (dropdownRef.value && !dropdownRef.value.contains(e.target)) {
    dropdownOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', onClickOutside)
})
onUnmounted(() => {
  document.removeEventListener('click', onClickOutside)
})
</script>
