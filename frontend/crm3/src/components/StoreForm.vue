<template>
  <form @submit.prevent="handleSubmit" class="store-form">
    <!-- 基本情報 -->
    <div class="form-section">
      <h3 class="section-title">基本情報</h3>

      <div class="form-group">
        <label for="name" class="required">店舗名</label>
        <input
          id="name"
          v-model="formData.name"
          type="text"
          placeholder="店舗名を入力"
          required
          :class="{ 'error': errors.name }"
        />
        <span v-if="errors.name" class="error-text">{{ errors.name }}</span>
      </div>

      <div class="form-group">
        <label for="address">住所</label>
        <textarea
          id="address"
          v-model="formData.address"
          rows="3"
          placeholder="住所を入力"
          :class="{ 'error': errors.address }"
        ></textarea>
        <span v-if="errors.address" class="error-text">{{ errors.address }}</span>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="phone">電話番号</label>
          <input
            id="phone"
            v-model="formData.phone"
            type="tel"
            placeholder="例: 03-1234-5678"
            :class="{ 'error': errors.phone }"
          />
          <span v-if="errors.phone" class="error-text">{{ errors.phone }}</span>
        </div>

        <div class="form-group">
          <label for="email">メールアドレス</label>
          <input
            id="email"
            v-model="formData.email"
            type="email"
            placeholder="例: store@example.com"
            :class="{ 'error': errors.email }"
          />
          <span v-if="errors.email" class="error-text">{{ errors.email }}</span>
        </div>
      </div>
    </div>

    <!-- ステータス設定 -->
    <div class="form-section">
      <h3 class="section-title">ステータス設定</h3>

      <div class="form-group">
        <label for="status" class="required">営業ステータス</label>
        <select
          id="status"
          v-model="formData.status"
          required
          :class="{ 'error': errors.status }"
        >
          <option value="active">営業中</option>
          <option value="inactive">休業中</option>
          <option value="maintenance">メンテナンス中</option>
        </select>
        <span v-if="errors.status" class="error-text">{{ errors.status }}</span>
        <div class="status-description">
          <span v-if="formData.status === 'active'" class="status-info status-info--active">
            通常営業状態です。顧客サービスが提供されます。
          </span>
          <span v-else-if="formData.status === 'inactive'" class="status-info status-info--inactive">
            一時的に休業している状態です。
          </span>
          <span v-else-if="formData.status === 'maintenance'" class="status-info status-info--maintenance">
            メンテナンス中です。営業を一時停止しています。
          </span>
        </div>
      </div>

      <div class="form-group">
        <label for="manager_id">店舗管理者ID</label>
        <input
          id="manager_id"
          v-model.number="formData.manager_id"
          type="number"
          placeholder="管理者IDを入力"
          min="1"
          :class="{ 'error': errors.manager_id }"
        />
        <span v-if="errors.manager_id" class="error-text">{{ errors.manager_id }}</span>
        <div class="help-text">
          店舗を管理する従業員のIDを入力してください（任意）
        </div>
      </div>
    </div>

    <!-- フォームアクション -->
    <div class="form-actions">
      <button type="button" class="btn btn--secondary" @click="$emit('cancel')" :disabled="loading">
        キャンセル
      </button>
      <button type="submit" class="btn btn--primary" :disabled="loading || !isValid">
        <span v-if="loading">保存中...</span>
        <span v-else>{{ isEditing ? '更新' : '作成' }}</span>
      </button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { reactive, computed, onMounted } from 'vue'
import type { Store, StoreCreateRequest, StoreUpdateRequest } from '@/stores/store'

interface Props {
  store?: Store
  loading?: boolean
}

interface Emits {
  submit: [data: StoreCreateRequest | StoreUpdateRequest]
  cancel: []
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<Emits>()

// Form data
const formData = reactive<StoreCreateRequest & { manager_id?: number }>({
  name: '',
  address: '',
  phone: '',
  email: '',
  status: 'active',
  manager_id: undefined
})

// Form validation errors
const errors = reactive<Partial<Record<keyof typeof formData, string>>>({})

// Computed properties
const isEditing = computed(() => !!props.store)

const isValid = computed(() => {
  return formData.name.trim() !== '' &&
         Object.keys(errors).length === 0
})

// Validation rules
const validateName = (value: string): string | null => {
  if (!value.trim()) {
    return '店舗名は必須です'
  }
  if (value.length > 100) {
    return '店舗名は100文字以内で入力してください'
  }
  return null
}

const validateAddress = (value?: string): string | null => {
  if (value && value.length > 500) {
    return '住所は500文字以内で入力してください'
  }
  return null
}

const validatePhone = (value?: string): string | null => {
  if (!value) return null

  const phoneRegex = /^[0-9-+\s()]+$/
  if (!phoneRegex.test(value)) {
    return '有効な電話番号を入力してください'
  }
  if (value.length > 20) {
    return '電話番号は20文字以内で入力してください'
  }
  return null
}

const validateEmail = (value?: string): string | null => {
  if (!value) return null

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(value)) {
    return '有効なメールアドレスを入力してください'
  }
  if (value.length > 320) {
    return 'メールアドレスは320文字以内で入力してください'
  }
  return null
}

const validateManagerId = (value?: number): string | null => {
  if (value !== undefined && (value < 1 || !Number.isInteger(value))) {
    return '有効な管理者IDを入力してください'
  }
  return null
}

// Validation functions
const validateField = (field: keyof typeof formData) => {
  let error: string | null = null

  switch (field) {
    case 'name':
      error = validateName(formData.name)
      break
    case 'address':
      error = validateAddress(formData.address)
      break
    case 'phone':
      error = validatePhone(formData.phone)
      break
    case 'email':
      error = validateEmail(formData.email)
      break
    case 'manager_id':
      error = validateManagerId(formData.manager_id)
      break
  }

  if (error) {
    errors[field] = error
  } else {
    delete errors[field]
  }
}

const validateForm = (): boolean => {
  // Clear previous errors
  Object.keys(errors).forEach(key => delete errors[key as keyof typeof errors])

  // Validate all fields
  validateField('name')
  validateField('address')
  validateField('phone')
  validateField('email')
  validateField('manager_id')

  return Object.keys(errors).length === 0
}

// Form submission
const handleSubmit = () => {
  if (!validateForm()) {
    return
  }

  const submitData: StoreCreateRequest | StoreUpdateRequest = {
    name: formData.name.trim(),
    address: formData.address?.trim() || undefined,
    phone: formData.phone?.trim() || undefined,
    email: formData.email?.trim() || undefined,
    status: formData.status,
    manager_id: formData.manager_id
  }

  emit('submit', submitData)
}

// Watch form changes for real-time validation
const setupValidation = () => {
  // Name validation
  const nameWatcher = () => validateField('name')

  // Address validation
  const addressWatcher = () => validateField('address')

  // Phone validation
  const phoneWatcher = () => validateField('phone')

  // Email validation
  const emailWatcher = () => validateField('email')

  // Manager ID validation
  const managerIdWatcher = () => validateField('manager_id')

  // Set up reactive validation
  setTimeout(() => {
    nameWatcher()
  }, 100)
}

// Initialize form
onMounted(() => {
  if (props.store) {
    // Edit mode: populate with existing data
    formData.name = props.store.name
    formData.address = props.store.address || ''
    formData.phone = props.store.phone || ''
    formData.email = props.store.email || ''
    formData.status = props.store.status
    formData.manager_id = props.store.manager_id
  }

  setupValidation()
})
</script>

<style scoped>
.store-form {
  max-width: 600px;
}

.form-section {
  margin-bottom: 2rem;
}

.section-title {
  margin: 0 0 1rem 0;
  color: #333;
  font-size: 1.125rem;
  font-weight: 600;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid #e5e7eb;
}

.form-group {
  margin-bottom: 1rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

label {
  display: block;
  margin-bottom: 0.25rem;
  font-weight: 500;
  color: #374151;
  font-size: 0.875rem;
}

label.required::after {
  content: ' *';
  color: #ef4444;
}

input,
select,
textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  transition: border-color 0.2s, box-shadow 0.2s;
}

input:focus,
select:focus,
textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

input.error,
select.error,
textarea.error {
  border-color: #ef4444;
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
}

.error-text {
  display: block;
  margin-top: 0.25rem;
  font-size: 0.75rem;
  color: #ef4444;
}

.help-text {
  margin-top: 0.25rem;
  font-size: 0.75rem;
  color: #6b7280;
}

.status-description {
  margin-top: 0.5rem;
}

.status-info {
  display: block;
  padding: 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
}

.status-info--active {
  background: #dcfce7;
  color: #166534;
}

.status-info--inactive {
  background: #fef3c7;
  color: #92400e;
}

.status-info--maintenance {
  background: #fee2e2;
  color: #991b1b;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
  padding-top: 2rem;
  border-top: 1px solid #e5e7eb;
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
  min-width: 100px;
}

.btn--primary {
  background: #3b82f6;
  color: white;
}

.btn--primary:hover:not(:disabled) {
  background: #2563eb;
}

.btn--secondary {
  background: #6b7280;
  color: white;
}

.btn--secondary:hover:not(:disabled) {
  background: #4b5563;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .form-row {
    grid-template-columns: 1fr;
  }

  .form-actions {
    flex-direction: column-reverse;
  }

  .btn {
    width: 100%;
  }
}
</style>