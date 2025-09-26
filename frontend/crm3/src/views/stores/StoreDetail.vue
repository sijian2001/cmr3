<template>
  <div class="store-detail">
    <!-- Header -->
    <div class="store-detail__header">
      <div class="breadcrumb">
        <router-link to="/stores" class="breadcrumb__link">店舗一覧</router-link>
        <span class="breadcrumb__separator">›</span>
        <span class="breadcrumb__current">
          {{ currentStore?.name || 'Loading...' }}
        </span>
      </div>

      <div class="header-actions">
        <button
          class="btn btn--secondary"
          @click="$router.push('/stores')"
        >
          戻る
        </button>
        <button
          v-if="currentStore"
          class="btn btn--primary"
          @click="editStore"
        >
          編集
        </button>
      </div>
    </div>

    <!-- Loading Spinner -->
    <div v-if="loading" class="loading-spinner" role="status" aria-label="読み込み中">
      読み込み中...
    </div>

    <!-- Error Message -->
    <div v-else-if="error" class="error-message" role="alert">
      {{ error }}
      <button @click="loadStore" class="btn btn--secondary">再読み込み</button>
    </div>

    <!-- Store Details -->
    <div v-else-if="currentStore" class="store-content">
      <!-- Basic Information Card -->
      <div class="info-card">
        <div class="info-card__header">
          <h2>基本情報</h2>
          <span
            :class="['status-badge', 'status-badge--large', `status-badge--${getStoreStatusColor(currentStore.status)}`]"
          >
            {{ getStoreStatusLabel(currentStore.status) }}
          </span>
        </div>

        <div class="info-grid">
          <div class="info-item">
            <span class="info-item__label">店舗名</span>
            <span class="info-item__value">{{ currentStore.name }}</span>
          </div>

          <div class="info-item" v-if="currentStore.address">
            <span class="info-item__label">住所</span>
            <span class="info-item__value">{{ currentStore.address }}</span>
          </div>

          <div class="info-item" v-if="currentStore.phone">
            <span class="info-item__label">電話番号</span>
            <span class="info-item__value">
              <a :href="`tel:${currentStore.phone}`" class="phone-link">
                {{ currentStore.phone }}
              </a>
            </span>
          </div>

          <div class="info-item" v-if="currentStore.email">
            <span class="info-item__label">メールアドレス</span>
            <span class="info-item__value">
              <a :href="`mailto:${currentStore.email}`" class="email-link">
                {{ currentStore.email }}
              </a>
            </span>
          </div>

          <div class="info-item" v-if="currentStore.manager_id">
            <span class="info-item__label">管理者ID</span>
            <span class="info-item__value">{{ currentStore.manager_id }}</span>
          </div>

          <div class="info-item">
            <span class="info-item__label">作成日</span>
            <span class="info-item__value">{{ formatDateTime(currentStore.created_at) }}</span>
          </div>

          <div class="info-item">
            <span class="info-item__label">更新日</span>
            <span class="info-item__value">{{ formatDateTime(currentStore.updated_at) }}</span>
          </div>
        </div>
      </div>

      <!-- Status Management Card -->
      <div class="info-card">
        <div class="info-card__header">
          <h2>ステータス管理</h2>
        </div>

        <div class="status-management">
          <div class="status-current">
            <span class="status-current__label">現在のステータス:</span>
            <span
              :class="['status-badge', `status-badge--${getStoreStatusColor(currentStore.status)}`]"
            >
              {{ getStoreStatusLabel(currentStore.status) }}
            </span>
          </div>

          <div class="status-actions">
            <button
              v-if="currentStore.status !== 'active'"
              class="btn btn--success"
              @click="updateStatus('active')"
              :disabled="loading"
            >
              営業開始
            </button>
            <button
              v-if="currentStore.status !== 'inactive'"
              class="btn btn--warning"
              @click="updateStatus('inactive')"
              :disabled="loading"
            >
              休業
            </button>
            <button
              v-if="currentStore.status !== 'maintenance'"
              class="btn btn--error"
              @click="updateStatus('maintenance')"
              :disabled="loading"
            >
              メンテナンス
            </button>
          </div>

          <div class="status-description">
            <div v-if="currentStore.status === 'active'" class="status-info status-info--active">
              <h4>営業中</h4>
              <p>この店舗は通常営業中です。顧客サービスが提供されています。</p>
            </div>
            <div v-else-if="currentStore.status === 'inactive'" class="status-info status-info--inactive">
              <h4>休業中</h4>
              <p>この店舗は一時的に休業中です。営業再開の予定については店舗にお問い合わせください。</p>
            </div>
            <div v-else-if="currentStore.status === 'maintenance'" class="status-info status-info--maintenance">
              <h4>メンテナンス中</h4>
              <p>この店舗はメンテナンス中です。営業を一時停止しています。</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Actions Card -->
      <div class="info-card">
        <div class="info-card__header">
          <h2>操作</h2>
        </div>

        <div class="actions-grid">
          <button class="action-btn action-btn--primary" @click="editStore">
            <span class="action-btn__icon">✏️</span>
            <span class="action-btn__text">編集</span>
          </button>

          <button class="action-btn action-btn--secondary" @click="viewStaff">
            <span class="action-btn__icon">👥</span>
            <span class="action-btn__text">従業員一覧</span>
          </button>

          <button class="action-btn action-btn--warning" @click="confirmDelete">
            <span class="action-btn__icon">🗑️</span>
            <span class="action-btn__text">削除</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Edit Store Modal -->
    <div v-if="showEditModal" class="modal-overlay" @click="closeEditModal">
      <div class="modal" @click.stop>
        <div class="modal__header">
          <h2>店舗編集</h2>
          <button class="modal__close" @click="closeEditModal">×</button>
        </div>
        <div class="modal__body">
          <StoreForm
            :store="currentStore || undefined"
            :loading="loading"
            @submit="handleUpdateStore"
            @cancel="closeEditModal"
          />
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="modal-overlay" @click="closeDeleteModal">
      <div class="modal modal--small" @click.stop>
        <div class="modal__header">
          <h2>店舗削除の確認</h2>
        </div>
        <div class="modal__body">
          <p>「{{ currentStore?.name }}」を削除しますか？</p>
          <p class="warning-text">
            この操作は取り消せません。店舗に関連するデータもすべて削除されます。
          </p>
        </div>
        <div class="modal__footer">
          <button class="btn btn--secondary" @click="closeDeleteModal">
            キャンセル
          </button>
          <button class="btn btn--danger" @click="handleDeleteStore" :disabled="loading">
            削除
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useStoreStore, type Store, type StoreUpdateRequest } from '@/stores/store'
import StoreForm from '@/components/StoreForm.vue'

const route = useRoute()
const router = useRouter()
const storeStore = useStoreStore()

const {
  currentStore,
  loading,
  error
} = storeToRefs(storeStore)

const {
  fetchStore,
  updateStore,
  updateStoreStatus,
  deleteStore,
  clearCurrentStore,
  clearError,
  getStoreStatusLabel,
  getStoreStatusColor
} = storeStore

// Local state
const showEditModal = ref(false)
const showDeleteModal = ref(false)

// Methods
const loadStore = async () => {
  const storeId = parseInt(route.params.id as string)
  if (isNaN(storeId)) {
    router.push('/stores')
    return
  }

  try {
    await fetchStore(storeId)
  } catch (err) {
    // Error is handled by the store
  }
}

const editStore = () => {
  showEditModal.value = true
}

const updateStatus = async (status: Store['status']) => {
  if (!currentStore.value) return

  try {
    await updateStoreStatus(currentStore.value.id, status)
  } catch (err) {
    // Error is handled by the store
  }
}

const confirmDelete = () => {
  showDeleteModal.value = true
}

const viewStaff = () => {
  router.push(`/stores/${currentStore.value?.id}/staff`)
}

const closeEditModal = () => {
  showEditModal.value = false
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
}

const handleUpdateStore = async (storeData: StoreUpdateRequest) => {
  if (!currentStore.value) return

  try {
    await updateStore(currentStore.value.id, storeData)
    closeEditModal()
  } catch (err) {
    // Error is handled by the store
  }
}

const handleDeleteStore = async () => {
  if (!currentStore.value) return

  try {
    await deleteStore(currentStore.value.id)
    closeDeleteModal()
    router.push('/stores')
  } catch (err) {
    // Error is handled by the store
  }
}

const formatDateTime = (dateString: string) => {
  return new Date(dateString).toLocaleString('ja-JP')
}

// Lifecycle
onMounted(() => {
  loadStore()
})

onUnmounted(() => {
  clearCurrentStore()
  clearError()
})
</script>

<style scoped>
.store-detail {
  max-width: 1000px;
  margin: 0 auto;
  padding: 2rem;
}

.store-detail__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
}

.breadcrumb__link {
  color: #3b82f6;
  text-decoration: none;
}

.breadcrumb__link:hover {
  text-decoration: underline;
}

.breadcrumb__separator {
  color: #6b7280;
}

.breadcrumb__current {
  color: #374151;
  font-weight: 500;
}

.header-actions {
  display: flex;
  gap: 1rem;
}

.store-content {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.info-card {
  background: white;
  border-radius: 0.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 2rem;
}

.info-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e5e7eb;
}

.info-card__header h2 {
  margin: 0;
  color: #333;
  font-size: 1.25rem;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.info-item__label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #6b7280;
}

.info-item__value {
  font-size: 1rem;
  color: #374151;
}

.phone-link,
.email-link {
  color: #3b82f6;
  text-decoration: none;
}

.phone-link:hover,
.email-link:hover {
  text-decoration: underline;
}

.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 500;
}

.status-badge--large {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
}

.status-badge--success {
  background: #dcfce7;
  color: #166534;
}

.status-badge--warning {
  background: #fef3c7;
  color: #92400e;
}

.status-badge--error {
  background: #fee2e2;
  color: #991b1b;
}

.status-management {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.status-current {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.status-current__label {
  font-weight: 500;
  color: #374151;
}

.status-actions {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.status-description {
  margin-top: 1rem;
}

.status-info {
  padding: 1rem;
  border-radius: 0.375rem;
  border-left: 4px solid;
}

.status-info--active {
  background: #f0f9ff;
  border-color: #10b981;
}

.status-info--inactive {
  background: #fffbeb;
  border-color: #f59e0b;
}

.status-info--maintenance {
  background: #fef2f2;
  border-color: #ef4444;
}

.status-info h4 {
  margin: 0 0 0.5rem 0;
  font-size: 1rem;
  font-weight: 600;
}

.status-info p {
  margin: 0;
  font-size: 0.875rem;
  color: #6b7280;
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  padding: 1.5rem;
  border: none;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
  text-decoration: none;
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.action-btn--primary {
  background: #3b82f6;
  color: white;
}

.action-btn--secondary {
  background: #6b7280;
  color: white;
}

.action-btn--warning {
  background: #ef4444;
  color: white;
}

.action-btn__icon {
  font-size: 1.5rem;
}

.action-btn__text {
  font-size: 0.875rem;
  font-weight: 500;
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn--primary {
  background: #3b82f6;
  color: white;
}

.btn--primary:hover {
  background: #2563eb;
}

.btn--secondary {
  background: #6b7280;
  color: white;
}

.btn--secondary:hover {
  background: #4b5563;
}

.btn--success {
  background: #10b981;
  color: white;
}

.btn--success:hover {
  background: #059669;
}

.btn--warning {
  background: #f59e0b;
  color: white;
}

.btn--warning:hover {
  background: #d97706;
}

.btn--error {
  background: #ef4444;
  color: white;
}

.btn--error:hover {
  background: #dc2626;
}

.btn--danger {
  background: #ef4444;
  color: white;
}

.btn--danger:hover {
  background: #dc2626;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.loading-spinner {
  text-align: center;
  padding: 3rem;
  color: #666;
}

.error-message {
  background: #fee2e2;
  color: #991b1b;
  padding: 1rem;
  border-radius: 0.375rem;
  margin-bottom: 1rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: white;
  border-radius: 0.5rem;
  max-width: 600px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
}

.modal--small {
  max-width: 400px;
}

.modal__header {
  padding: 1.5rem;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal__header h2 {
  margin: 0;
  color: #333;
}

.modal__close {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #666;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal__body {
  padding: 1.5rem;
}

.modal__footer {
  padding: 1.5rem;
  border-top: 1px solid #e5e7eb;
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
}

.warning-text {
  color: #ef4444;
  font-size: 0.875rem;
  margin-top: 0.5rem;
}

@media (max-width: 768px) {
  .store-detail__header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }

  .header-actions {
    justify-content: stretch;
  }

  .header-actions .btn {
    flex: 1;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .status-actions {
    flex-direction: column;
  }

  .actions-grid {
    grid-template-columns: 1fr;
  }
}
</style>