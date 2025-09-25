import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import HomeView from '@/views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { requiresAuth: true }
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/auth/LoginView.vue'),
      meta: { requiresGuest: true }
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/auth/RegisterView.vue'),
      meta: { requiresGuest: true }
    },
    {
      path: '/reset-password',
      name: 'password-reset-request',
      component: () => import('@/views/auth/PasswordResetRequest.vue'),
      meta: { requiresGuest: true }
    },
    {
      path: '/reset-password-confirm',
      name: 'password-reset-confirm',
      component: () => import('@/views/auth/PasswordResetConfirm.vue'),
      meta: { requiresGuest: true }
    }
  ]
})

// ナビゲーションガード
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // 初期化がまだの場合は実行
  if (authStore.token && !authStore.user) {
    await authStore.initialize()
  }

  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const requiresGuest = to.matched.some(record => record.meta.requiresGuest)

  if (requiresAuth && !authStore.isAuthenticated) {
    // 認証が必要だがログインしていない場合
    next('/login')
  } else if (requiresGuest && authStore.isAuthenticated) {
    // ゲスト限定ページだが既にログインしている場合
    next('/')
  } else {
    next()
  }
})

export default router
