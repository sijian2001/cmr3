<template>
  <div class="modal-overlay" @click="handleOverlayClick">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h2>{{ customer ? '顧客情報編集' : '新規顧客登録' }}</h2>
        <button @click="$emit('cancel')" class="btn-close">×</button>
      </div>

      <form @submit.prevent="handleSubmit" class="modal-body">
        <div class="form-group">
          <label for="name">名前 <span class="required">*</span></label>
          <input
            id="name"
            v-model="form.name"
            type="text"
            required
            :class="{ 'error': errors.name }"
            placeholder="顧客名を入力"
          >
          <div v-if="errors.name" class="error-text">{{ errors.name }}</div>
        </div>

        <div class="form-group">
          <label for="email">メールアドレス <span class="required">*</span></label>
          <input
            id="email"
            v-model="form.email"
            type="email"
            required
            :class="{ 'error': errors.email }"
            placeholder="example@example.com"
          >
          <div v-if="errors.email" class="error-text">{{ errors.email }}</div>
        </div>

        <div class="form-group">
          <label for="phone">電話番号</label>
          <input
            id="phone"
            v-model="form.phone"
            type="tel"
            :class="{ 'error': errors.phone }"
            placeholder="090-1234-5678"
          >
          <div v-if="errors.phone" class="error-text">{{ errors.phone }}</div>
        </div>

        <div class="form-group">
          <label for="address">住所</label>
          <textarea
            id="address"
            v-model="form.address"
            rows="3"
            :class="{ 'error': errors.address }"
            placeholder="住所を入力"
          ></textarea>
          <div v-if="errors.address" class="error-text">{{ errors.address }}</div>
        </div>

        <div class="modal-footer">
          <button type="button" @click="$emit('cancel')" class="btn-cancel">
            キャンセル
          </button>
          <button type="submit" :disabled="isSubmitting" class="btn-save">
            {{ isSubmitting ? '保存中...' : '保存' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import type { Customer } from '@/stores/customer'

interface Props {
  customer: Customer | null
}

interface Emits {
  (event: 'save', data: Omit<Customer, 'id' | 'created_at' | 'updated_at'>): void
  (event: 'cancel'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// フォームデータ
const form = reactive({
  name: '',
  email: '',
  phone: '',
  address: ''
})

// エラー状態
const errors = reactive({
  name: '',
  email: '',
  phone: '',
  address: ''
})

const isSubmitting = ref(false)

// バリデーション
const validateForm = (): boolean => {
  // エラーをクリア
  Object.keys(errors).forEach(key => {
    errors[key as keyof typeof errors] = ''
  })

  let isValid = true

  // 名前のバリデーション
  if (!form.name.trim()) {
    errors.name = '名前は必須です'
    isValid = false
  }

  // メールアドレスのバリデーション
  if (!form.email.trim()) {
    errors.email = 'メールアドレスは必須です'
    isValid = false
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    errors.email = '有効なメールアドレスを入力してください'
    isValid = false
  }

  // 電話番号のバリデーション（任意だが、入力された場合はチェック）
  if (form.phone && !/^[\d\-\(\)\+\s]+$/.test(form.phone)) {
    errors.phone = '有効な電話番号を入力してください'
    isValid = false
  }

  return isValid
}

// フォーム送信
const handleSubmit = async () => {
  if (!validateForm()) return

  try {
    isSubmitting.value = true

    const customerData = {
      name: form.name.trim(),
      email: form.email.trim(),
      phone: form.phone.trim() || undefined,
      address: form.address.trim() || undefined
    }

    emit('save', customerData)
  } catch (error) {
    console.error('フォーム送信エラー:', error)
  } finally {
    isSubmitting.value = false
  }
}

// モーダル外クリック時の処理
const handleOverlayClick = (event: MouseEvent) => {
  if (event.target === event.currentTarget) {
    emit('cancel')
  }
}

// 初期化
onMounted(() => {
  if (props.customer) {
    form.name = props.customer.name || ''
    form.email = props.customer.email || ''
    form.phone = props.customer.phone || ''
    form.address = props.customer.address || ''
  }
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #eee;
}

.modal-header h2 {
  margin: 0;
  color: #333;
  font-size: 18px;
}

.btn-close {
  background: none;
  border: none;
  font-size: 24px;
  color: #666;
  cursor: pointer;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-close:hover {
  color: #333;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-weight: bold;
  margin-bottom: 5px;
  color: #555;
}

.required {
  color: #dc3545;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
}

.form-group input.error,
.form-group textarea.error {
  border-color: #dc3545;
  box-shadow: 0 0 0 2px rgba(220, 53, 69, 0.25);
}

.error-text {
  color: #dc3545;
  font-size: 12px;
  margin-top: 5px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.btn-cancel {
  background: transparent;
  color: #6c757d;
  border: 1px solid #6c757d;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-cancel:hover {
  background: #6c757d;
  color: white;
}

.btn-save {
  background: #007bff;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-save:hover {
  background: #0056b3;
}

.btn-save:disabled {
  background: #6c757d;
  cursor: not-allowed;
}
</style>