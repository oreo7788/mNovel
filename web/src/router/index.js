import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  { path: '/', name: 'home', component: () => import('../views/Home.vue'), meta: { title: '首页' } },
  { path: '/login', name: 'login', component: () => import('../views/Login.vue'), meta: { title: '登录' } },
  { path: '/register', name: 'register', component: () => import('../views/Register.vue'), meta: { title: '注册' } },
  { path: '/wardrobe', name: 'wardrobe', component: () => import('../views/Wardrobe.vue'), meta: { title: '衣橱' } },
  { path: '/recommend', name: 'recommend', component: () => import('../views/Recommend.vue'), meta: { title: '推荐' } },
  { path: '/favorites', name: 'favorites', component: () => import('../views/Favorites.vue'), meta: { title: '收藏' } },
  { path: '/upload', name: 'upload', component: () => import('../views/Upload.vue'), meta: { title: '上传衣物' } },
  { path: '/settings', name: 'settings', component: () => import('../views/Settings.vue'), meta: { title: '设置' } },
  { path: '/profile', name: 'profile', component: () => import('../views/Profile.vue'), meta: { title: '个人中心' } },
  { path: '/terms', name: 'terms', component: () => import('../views/Terms.vue'), meta: { title: '服务条款' } },
  { path: '/forgot-password', name: 'forgot-password', component: () => import('../views/ForgotPassword.vue'), meta: { title: '找回密码' } },
  { path: '/reset-password', name: 'reset-password', component: () => import('../views/ResetPassword.vue'), meta: { title: '修改密码' } },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

router.afterEach((to) => {
  document.title = to.meta.title ? `${to.meta.title} - 衣搭库` : '衣搭库'
})

export default router
