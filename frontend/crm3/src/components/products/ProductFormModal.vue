<template>
  <div class="modal-overlay" @click="handleOverlayClick">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h2>{{ product ? '製品情報編集' : '新規製品登録' }}</h2>
        <button @click="$emit('cancel')" class="btn-close">×</button>
      </div>

      <form @submit.prevent="handleSubmit" class="modal-body">
        <div class="form-group">
          <label for="name">製品名 <span class="required">*</span></label>
          <input
            id="name"
            v-model="form.name"
            type="text"
            required
            :class="{ 'error': errors.name }"
            placeholder="製品名を入力"
          >
          <div v-if="errors.name" class="error-text">{{ errors.name }}</div>
        </div>

        <div class="form-group">
          <label for="sku">SKU <span class="required">*</span></label>
          <input
            id="sku"
            v-model="form.sku"
            type="text"
            required
            :class="{ 'error': errors.sku }"
            placeholder="SKU-001"
          >
          <div v-if="errors.sku" class="error-text">{{ errors.sku }}</div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="price">価格 <span class="required">*</span></label>
            <input
              id="price"
              v-model.number="form.price"
              type="number"
              min="0"
              step="1"
              required
              :class="{ 'error': errors.price }"
              placeholder="1000"
            >
            <div v-if="errors.price" class="error-text">{{ errors.price }}</div>
          </div>

          <div class="form-group">
            <label for="stock_quantity">在庫数量 <span class="required">*</span></label>
            <input
              id="stock_quantity"
              v-model.number="form.stock_quantity"
              type="number"
              min="0"
              step="1"
              required
              :class="{ 'error': errors.stock_quantity }"
              placeholder="100"
            >
            <div v-if="errors.stock_quantity" class="error-text">{{ errors.stock_quantity }}</div>
          </div>
        </div>

        <div class="form-group">
          <label for="description">説明</label>
          <textarea
            id="description"
            v-model="form.description"
            rows="4"
            :class="{ 'error': errors.description }"
            placeholder="製品の詳細説明を入力（任意）"
          ></textarea>
          <div v-if="errors.description" class="error-text">{{ errors.description }}</div>
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
import type { Product } from '@/stores/product'

interface Props {
  product: Product | null
}

interface Emits {
  (event: 'save', data: Omit<Product, 'id' | 'created_at' | 'updated_at'>): void
  (event: 'cancel'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// フォームデータ
const form = reactive({
  name: '',
  description: '',
  sku: '',
  price: 0,
  stock_quantity: 0
})

// エラー状態
const errors = reactive({
  name: '',
  description: '',
  sku: '',
  price: '',
  stock_quantity: ''
})

const isSubmitting = ref(false)

// バリデーション
const validateForm = (): boolean => {
  // エラーをクリア
  Object.keys(errors).forEach(key => {
    errors[key as keyof typeof errors] = ''
  })

  let isValid = true

  // 製品名のバリデーション
  if (!form.name.trim()) {
    errors.name = '製品名は必須です'
    isValid = false
  } else if (form.name.trim().length > 100) {
    errors.name = '製品名は100文字以内で入力してください'
    isValid = false
  }

  // SKUのバリデーション
  if (!form.sku.trim()) {
    errors.sku = 'SKUは必須です'
    isValid = false
  } else if (form.sku.trim().length > 50) {
    errors.sku = 'SKUは50文字以内で入力してください'
    isValid = false
  }

  // 価格のバリデーション
  if (form.price === null || form.price === undefined || isNaN(form.price)) {
    errors.price = '価格は必須です'
    isValid = false
  } else if (form.price <= 0) {
    errors.price = '価格は正数である必要があります'
    isValid = false
  } else if (form.price > 99999999) {
    errors.price = '価格が大きすぎます'
    isValid = false
  }

  // 在庫数量のバリデーション
  if (form.stock_quantity === null || form.stock_quantity === undefined || isNaN(form.stock_quantity)) {
    errors.stock_quantity = '在庫数量は必須です'
    isValid = false
  } else if (form.stock_quantity < 0) {
    errors.stock_quantity = '在庫数量は0以上である必要があります'
    isValid = false
  } else if (form.stock_quantity > 999999) {
    errors.stock_quantity = '在庫数量が大きすぎます'
    isValid = false
  }

  // 説明のバリデーション（任意だが、入力された場合はチェック）
  if (form.description && form.description.length > 1000) {
    errors.description = '説明は1000文字以内で入力してください'
    isValid = false
  }

  return isValid
}

// フォーム送信
const handleSubmit = async () => {
  if (!validateForm()) return

  try {
    isSubmitting.value = true

    const productData = {
      name: form.name.trim(),
      description: form.description.trim() || undefined,
      sku: form.sku.trim(),
      price: form.price,
      stock_quantity: form.stock_quantity
    }

    emit('save', productData)
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
  if (props.product) {
    form.name = props.product.name || ''
    form.description = props.product.description || ''
    form.sku = props.product.sku || ''
    form.price = props.product.price || 0
    form.stock_quantity = props.product.stock_quantity || 0
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
  max-width: 600px;
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

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
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

/* レスポンシブ対応 */
@media (max-width: 768px) {
  .modal-content {
    width: 95%;
    margin: 10px;
  }

  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>