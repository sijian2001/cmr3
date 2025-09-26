<template>
  <form @submit.prevent="handleSubmit" class="staff-form">
    <!-- 基本情報 -->
    <div class="form-section">
      <h3 class="section-title">基本情報</h3>

      <div class="form-group">
        <label for="name" class="required">氏名</label>
        <input
          id="name"
          v-model="formData.name"
          type="text"
          placeholder="氏名を入力"
          required
          :class="{ 'error': errors.name }"
        />
        <span v-if="errors.name" class="error-text">{{ errors.name }}</span>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="email" class="required">メールアドレス</label>
          <input
            id="email"
            v-model="formData.email"
            type="email"
            placeholder="例: staff@example.com"
            required
            :class="{ 'error': errors.email }"
          />
          <span v-if="errors.email" class="error-text">{{ errors.email }}</span>
        </div>

        <div class="form-group">
          <label for="phone">電話番号</label>
          <input
            id="phone"
            v-model="formData.phone"
            type="tel"
            placeholder="例: 090-1234-5678"
            :class="{ 'error': errors.phone }"
          />
          <span v-if="errors.phone" class="error-text">{{ errors.phone }}</span>
        </div>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="position" class="required">役職</label>
          <select
            id="position"
            v-model="formData.position"
            required
            :class="{ 'error': errors.position }"
          >
            <option value="">役職を選択</option>
            <option value="店長">店長</option>
            <option value="副店長">副店長</option>
            <option value="正社員">正社員</option>
            <option value="アルバイト">アルバイト</option>
            <option value="パート">パート</option>
            <option value="契約社員">契約社員</option>
          </select>
          <span v-if="errors.position" class="error-text">{{ errors.position }}</span>
        </div>

        <div class="form-group">
          <label for="hire_date" class="required">入社日</label>
          <input
            id="hire_date"
            v-model="formData.hire_date"
            type="date"
            required
            :class="{ 'error': errors.hire_date }"
          />
          <span v-if="errors.hire_date" class="error-text">{{ errors.hire_date }}</span>
        </div>
      </div>
    </div>

    <!-- 所属・ステータス設定 -->
    <div class="form-section">
      <h3 class="section-title">所属・ステータス設定</h3>

      <div class="form-row">
        <div class="form-group">
          <label for="store_id">所属店舗</label>
          <select
            id="store_id"
            v-model.number="formData.store_id"
            :class="{ 'error': errors.store_id }"
          >
            <option value="">店舗を選択（未割り当て）</option>
            <option
              v-for="store in availableStores"
              :key="store.id"
              :value="store.id"
            >
              {{ store.name }}
            </option>
          </select>
          <span v-if="errors.store_id" class="error-text">{{ errors.store_id }}</span>
        </div>

        <div class="form-group">
          <label for="status" class="required">ステータス</label>
          <select
            id="status"
            v-model="formData.status"
            required
            :class="{ 'error': errors.status }"
          >
            <option value="active">在籍</option>
            <option value="inactive">休職</option>
            <option value="on_leave">休暇中</option>
          </select>
          <span v-if="errors.status" class="error-text">{{ errors.status }}</span>
        </div>
      </div>

      <div class="status-description">
        <div class="status-info">
          <h4>ステータスについて</h4>
          <ul>
            <li><strong>在籍:</strong> 通常勤務状態です。</li>
            <li><strong>休職:</strong> 一時的に業務から離れている状態です。</li>
            <li><strong>休暇中:</strong> 有給休暇や特別休暇を取得している状態です。</li>
          </ul>
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
import { reactive, computed, onMounted, ref } from 'vue'
import type { Staff, StaffCreateRequest, StaffUpdateRequest } from '@/stores/staff'

interface Props {
  staff?: Staff
  loading?: boolean
}

interface Emits {
  submit: [data: StaffCreateRequest | StaffUpdateRequest]
  cancel: []
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<Emits>()

// Available stores (mock data for now, should be fetched from stores API)
const availableStores = ref([
  { id: 1, name: '新宿店' },
  { id: 2, name: '渋谷店' },
  { id: 3, name: '池袋店' }
])

// Form data
const formData = reactive<StaffCreateRequest & { store_id?: number }>({
  name: '',
  email: '',
  phone: '',
  position: '',
  store_id: undefined,
  hire_date: '',
  status: 'active'
})

// Form validation errors
const errors = reactive<Partial<Record<keyof typeof formData, string>>>({})

// Computed properties
const isEditing = computed(() => !!props.staff)

const isValid = computed(() => {
  return formData.name.trim() !== '' &&
         formData.email.trim() !== '' &&
         formData.position !== '' &&
         formData.hire_date !== '' &&
         Object.keys(errors).length === 0
})

// Validation rules
const validateName = (value: string): string | null => {
  if (!value.trim()) {
    return '氏名は必須です'
  }
  if (value.length > 100) {
    return '氏名は100文字以内で入力してください'
  }
  return null
}

const validateEmail = (value: string): string | null => {
  if (!value.trim()) {
    return 'メールアドレスは必須です'
  }

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(value)) {
    return '有効なメールアドレスを入力してください'
  }
  if (value.length > 320) {
    return 'メールアドレスは320文字以内で入力してください'
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

const validatePosition = (value: string): string | null => {
  if (!value) {
    return '役職は必須です'
  }
  return null
}

const validateHireDate = (value: string): string | null => {
  if (!value) {
    return '入社日は必須です'
  }

  const hireDate = new Date(value)
  const today = new Date()

  // 未来の日付はチェック（入社日は過去または今日まで）
  if (hireDate > today) {
    return '入社日は今日以前の日付を入力してください'
  }

  // 1900年以前はエラー
  if (hireDate.getFullYear() < 1900) {
    return '有効な入社日を入力してください'
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
    case 'email':
      error = validateEmail(formData.email)
      break
    case 'phone':
      error = validatePhone(formData.phone)
      break
    case 'position':
      error = validatePosition(formData.position)
      break
    case 'hire_date':
      error = validateHireDate(formData.hire_date)
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
  validateField('email')
  validateField('phone')
  validateField('position')
  validateField('hire_date')

  return Object.keys(errors).length === 0
}

// Form submission
const handleSubmit = () => {
  if (!validateForm()) {
    return
  }

  const submitData: StaffCreateRequest | StaffUpdateRequest = {
    name: formData.name.trim(),
    email: formData.email.trim(),
    phone: formData.phone?.trim() || undefined,
    position: formData.position,
    store_id: formData.store_id,
    hire_date: formData.hire_date,
    status: formData.status
  }

  emit('submit', submitData)
}

// Initialize form
onMounted(() => {
  if (props.staff) {
    // Edit mode: populate with existing data
    formData.name = props.staff.name
    formData.email = props.staff.email
    formData.phone = props.staff.phone || ''
    formData.position = props.staff.position
    formData.store_id = props.staff.store_id
    formData.hire_date = props.staff.hire_date
    formData.status = props.staff.status
  } else {
    // Create mode: set default hire date to today
    const today = new Date().toISOString().split('T')[0]
    formData.hire_date = today
  }
})
</script>

<style scoped>
.staff-form {
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
select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  transition: border-color 0.2s, box-shadow 0.2s;
}

input:focus,
select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

input.error,
select.error {
  border-color: #ef4444;
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
}

.error-text {
  display: block;
  margin-top: 0.25rem;
  font-size: 0.75rem;
  color: #ef4444;
}

.status-description {
  margin-top: 1rem;
}

.status-info {
  background: #f9fafb;
  padding: 1rem;
  border-radius: 0.375rem;
  border-left: 4px solid #3b82f6;
}

.status-info h4 {
  margin: 0 0 0.5rem 0;
  font-size: 0.875rem;
  color: #374151;
}

.status-info ul {
  margin: 0;
  padding-left: 1.5rem;
  font-size: 0.8125rem;
  color: #6b7280;
}

.status-info li {
  margin-bottom: 0.25rem;
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