import { ref, computed, readonly } from 'vue'
import { defineStore } from 'pinia'
import {
  authService,
  type LoginRequest,
  type RegisterRequest,
  type PasswordResetRequest,
  type PasswordResetConfirmRequest,
  type UserResponse
} from '@/services/auth'

export const useAuthStore = defineStore('auth', () => {
  // ステート
  const user = ref<UserResponse | null>(null)
  const token = ref<string | null>(authService.getToken())
  const loading = ref(false)
  const error = ref<string | null>(null)

  // コンピューテッド
  const isAuthenticated = computed(() => !!token.value && !!user.value)

  // アクション
  async function login(credentials: LoginRequest) {
    loading.value = true
    error.value = null

    try {
      const response = await authService.login(credentials)

      // トークンとユーザー情報を保存
      token.value = response.token
      user.value = response.user
      authService.saveToken(response.token)

      return response
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'ログインに失敗しました'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function register(userData: RegisterRequest) {
    loading.value = true
    error.value = null

    try {
      const response = await authService.register(userData)

      // 登録後自動ログイン
      token.value = response.token
      user.value = response.user
      authService.saveToken(response.token)

      return response
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'ユーザー登録に失敗しました'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function requestPasswordReset(data: PasswordResetRequest) {
    loading.value = true
    error.value = null

    try {
      const response = await authService.requestPasswordReset(data)
      return response
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'パスワードリセット要求に失敗しました'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function confirmPasswordReset(data: PasswordResetConfirmRequest) {
    loading.value = true
    error.value = null

    try {
      const response = await authService.confirmPasswordReset(data)
      return response
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'パスワードリセットに失敗しました'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchProfile() {
    if (!token.value) return

    loading.value = true
    error.value = null

    try {
      const profile = await authService.getProfile()
      user.value = profile
      return profile
    } catch (err) {
      // トークンが無効な場合はログアウト
      logout()
      error.value = err instanceof Error ? err.message : 'プロフィールの取得に失敗しました'
      throw err
    } finally {
      loading.value = false
    }
  }

  function logout() {
    user.value = null
    token.value = null
    authService.removeToken()
    error.value = null
  }

  // 初期化（アプリ起動時にトークンがあればプロフィールを取得）
  async function initialize() {
    if (token.value) {
      try {
        await fetchProfile()
      } catch {
        // トークンが無効な場合は静かにログアウト
        logout()
      }
    }
  }

  // エラークリア
  function clearError() {
    error.value = null
  }

  return {
    // ステート
    user: readonly(user),
    token: readonly(token),
    loading: readonly(loading),
    error: readonly(error),

    // コンピューテッド
    isAuthenticated,

    // アクション
    login,
    register,
    requestPasswordReset,
    confirmPasswordReset,
    fetchProfile,
    logout,
    initialize,
    clearError
  }
})