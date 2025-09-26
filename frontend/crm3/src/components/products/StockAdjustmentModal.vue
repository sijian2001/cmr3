<template>
  <div class="modal-overlay" @click="handleOverlayClick">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h2>在庫調整</h2>
        <button @click="$emit('cancel')" class="btn-close">×</button>
      </div>

      <div class="modal-body">
        <div class="product-info">
          <h3>{{ product.name }}</h3>
          <p class="sku">SKU: {{ product.sku }}</p>
          <p class="current-stock">
            現在の在庫数量: <span :class="getStockClass(product.stock_quantity)">{{ product.stock_quantity }}</span>
          </p>
        </div>

        <form @submit.prevent="handleSubmit" class="adjustment-form">
          <div class="form-group">
            <label for="operation">操作種別</label>
            <div class="radio-group">
              <label class="radio-label">
                <input
                  type="radio"
                  name="operation"
                  value="add"
                  v-model="form.operation"
                >
                <span>入庫 (在庫を増やす)</span>
              </label>
              <label class="radio-label">
                <input
                  type="radio"
                  name="operation"
                  value="subtract"
                  v-model="form.operation"
                >
                <span>出庫 (在庫を減らす)</span>
              </label>
            </div>
          </div>

          <div class="form-group">
            <label for="quantity">数量 <span class="required">*</span></label>
            <input
              id="quantity"
              v-model.number="form.quantity"
              type="number"
              min="1"
              step="1"
              required
              :class="{ 'error': errors.quantity }"
              placeholder="調整する数量を入力"
            >
            <div v-if="errors.quantity" class="error-text">{{ errors.quantity }}</div>
          </div>

          <div class="form-group">
            <label for="reason">理由</label>
            <textarea
              id="reason"
              v-model="form.reason"
              rows="3"
              placeholder="在庫調整の理由を記載（任意）"
            ></textarea>
          </div>

          <div class="preview-section" v-if="form.quantity > 0">
            <h4>調整後の予想在庫数量</h4>
            <div class="stock-preview">
              <span class="current">{{ product.stock_quantity }}</span>
              <span class="operation">{{ form.operation === 'add' ? '+' : '-' }}</span>
              <span class="quantity">{{ form.quantity }}</span>
              <span class="equals">=</span>
              <span :class="getStockClass(getNewStockQuantity())">{{ getNewStockQuantity() }}</span>
            </div>
            <div v-if="form.operation === 'subtract' && getNewStockQuantity() < product.stock_quantity * 0.2" class="warning">
              ⚠️ 在庫が少なくなります
            </div>
          </div>

          <div class="modal-footer">
            <button type="button" @click="$emit('cancel')" class="btn-cancel">
              キャンセル
            </button>
            <button type="submit" :disabled="isSubmitting || !isFormValid" class="btn-save">
              {{ isSubmitting ? '調整中...' : '在庫調整を実行' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import type { Product } from '@/stores/product'

interface Props {
  product: Product
}

interface Emits {
  (event: 'adjust', quantity: number, operation: 'add' | 'subtract'): void
  (event: 'cancel'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// フォームデータ
const form = reactive({
  operation: 'add' as 'add' | 'subtract',
  quantity: 1,
  reason: ''
})

// エラー状態
const errors = reactive({
  quantity: ''
})

const isSubmitting = ref(false)

// Computed
const isFormValid = computed(() => {
  return form.quantity > 0 && !errors.quantity
})

const getNewStockQuantity = () => {
  if (form.operation === 'add') {
    return props.product.stock_quantity + form.quantity
  } else {
    return Math.max(0, props.product.stock_quantity - form.quantity)
  }
}

// メソッド
const getStockClass = (quantity: number) => {
  if (quantity === 0) return 'stock-empty'
  if (quantity <= 10) return 'stock-low'
  return 'stock-normal'
}

const validateForm = (): boolean => {
  errors.quantity = ''

  if (!form.quantity || form.quantity <= 0) {
    errors.quantity = '数量は1以上の整数を入力してください'
    return false
  }

  if (form.operation === 'subtract' && form.quantity > props.product.stock_quantity) {
    errors.quantity = '出庫数量は現在の在庫数量を超えることはできません'
    return false
  }

  return true
}

const handleSubmit = async () => {
  if (!validateForm()) return

  try {
    isSubmitting.value = true
    emit('adjust', form.quantity, form.operation)
  } catch (error) {
    console.error('在庫調整エラー:', error)
  } finally {
    isSubmitting.value = false
  }
}

const handleOverlayClick = (event: MouseEvent) => {
  if (event.target === event.currentTarget) {
    emit('cancel')
  }
}
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

.product-info {
  background: #f8f9fa;
  padding: 15px;
  border-radius: 6px;
  margin-bottom: 20px;
}

.product-info h3 {
  margin: 0 0 5px 0;
  color: #333;
}

.sku {
  margin: 0 0 10px 0;
  color: #666;
  font-size: 14px;
}

.current-stock {
  margin: 0;
  font-weight: bold;
}

.stock-normal {
  color: #28a745;
}

.stock-low {
  color: #ffc107;
}

.stock-empty {
  color: #dc3545;
}

.adjustment-form {
  margin-top: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-weight: bold;
  margin-bottom: 8px;
  color: #555;
}

.required {
  color: #dc3545;
}

.radio-group {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.radio-label {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.radio-label:hover {
  background-color: #f8f9fa;
}

.radio-label input[type="radio"] {
  margin-right: 8px;
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

.form-group input.error {
  border-color: #dc3545;
  box-shadow: 0 0 0 2px rgba(220, 53, 69, 0.25);
}

.error-text {
  color: #dc3545;
  font-size: 12px;
  margin-top: 5px;
}

.preview-section {
  background: #e9f7fe;
  padding: 15px;
  border-radius: 6px;
  margin-bottom: 20px;
}

.preview-section h4 {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 14px;
}

.stock-preview {
  font-size: 18px;
  font-weight: bold;
  display: flex;
  align-items: center;
  gap: 10px;
}

.current {
  color: #666;
}

.operation {
  color: #007bff;
  font-size: 20px;
}

.quantity {
  color: #007bff;
}

.equals {
  color: #666;
}

.warning {
  margin-top: 10px;
  color: #856404;
  background-color: #fff3cd;
  border: 1px solid #ffeaa7;
  padding: 8px 12px;
  border-radius: 4px;
  font-size: 14px;
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
  background: #28a745;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-save:hover {
  background: #218838;
}

.btn-save:disabled {
  background: #6c757d;
  cursor: not-allowed;
}
</style>