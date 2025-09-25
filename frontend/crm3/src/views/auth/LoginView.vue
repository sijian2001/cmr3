<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          ログイン
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          CRM3システムにサインインしてください
        </p>
      </div>

      <form class="mt-8 space-y-6" @submit.prevent="handleSubmit">
        <div class="rounded-md shadow-sm space-y-4">
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-1">
              メールアドレス
            </label>
            <input
              id="email"
              v-model="email"
              type="email"
              required
              class="block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              :class="{
                'border-red-300 focus:border-red-500 focus:ring-red-500': !!authStore.error,
              }"
              placeholder="あなたのメールアドレス"
              :disabled="authStore.loading"
            />
          </div>

          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-1">
              パスワード
            </label>
            <input
              id="password"
              v-model="password"
              type="password"
              required
              class="block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              :class="{
                'border-red-300 focus:border-red-500 focus:ring-red-500': !!authStore.error,
              }"
              placeholder="パスワード"
              :disabled="authStore.loading"
            />
          </div>
        </div>

        <div v-if="authStore.error" class="rounded-md bg-red-50 p-4">
          <div class="flex">
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">
                ログインに失敗しました
              </h3>
              <div class="mt-2 text-sm text-red-700">
                {{ authStore.error }}
              </div>
            </div>
          </div>
        </div>

        <div>
          <button
            type="submit"
            :disabled="authStore.loading || !email || !password"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="authStore.loading" class="flex items-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              ログイン中...
            </span>
            <span v-else>ログイン</span>
          </button>
        </div>

        <div class="flex items-center justify-between">
          <div class="text-sm">
            <router-link
              to="/reset-password"
              class="text-blue-600 hover:text-blue-500 font-medium"
            >
              パスワードを忘れた方はこちら
            </router-link>
          </div>
          <div class="text-sm">
            <router-link
              to="/register"
              class="text-blue-600 hover:text-blue-500 font-medium"
            >
              新規登録
            </router-link>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')

const handleSubmit = async () => {
  authStore.clearError()

  try {
    await authStore.login({
      email: email.value,
      password: password.value
    })

    // ログイン成功時はダッシュボードにリダイレクト
    router.push('/')
  } catch (error) {
    // エラーはストアで処理される
    console.error('Login failed:', error)
  }
}
</script>

<style scoped>
/* カスタムスタイルが必要な場合はここに追加 */
</style>