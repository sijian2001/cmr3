<template>
  <div class="store-list">
    <!-- Header -->
    <div class="store-list__header">
      <h1>店舗管理</h1>
      <button
        class="btn btn--primary"
        @click="showCreateModal = true"
      >
        新規店舗登録
      </button>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stats-card stats-card--primary">
        <h3>営業中</h3>
        <p class="stats-number">{{ statusCounts.active }}</p>
      </div>
      <div class="stats-card stats-card--warning">
        <h3>休業中</h3>
        <p class="stats-number">{{ statusCounts.inactive }}</p>
      </div>
      <div class="stats-card stats-card--error">
        <h3>メンテナンス中</h3>
        <p class="stats-number">{{ statusCounts.maintenance }}</p>
      </div>
      <div class="stats-card">
        <h3>総店舗数</h3>
        <p class="stats-number">{{ statusCounts.total }}</p>
      </div>
    </div>

    <!-- Search Form -->
    <div class="search-form">
      <div class="search-form__group">
        <label for="search-name">店舗名</label>
        <input
          id="search-name"
          v-model="searchForm.name"
          type="text"
          placeholder="店舗名で検索"
          @input="debouncedSearch"
        />
      </div>
      <div class="search-form__group">
        <label for="search-status">ステータス</label>
        <select id="search-status" v-model="searchForm.status" @change="handleSearch">
          <option value="">すべて</option>
          <option value="active">営業中</option>
          <option value="inactive">休業中</option>
          <option value="maintenance">メンテナンス中</option>
        </select>
      </div>
      <div class="search-form__actions">
        <button class="btn btn--secondary" @click="clearSearch">
          クリア
        </button>
      </div>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="error-message" role="alert">
      {{ error }}
      <button @click="clearError" class="error-message__close">×</button>
    </div>

    <!-- Loading Spinner -->
    <div v-if="loading" class="loading-spinner" role="status" aria-label="読み込み中">
      読み込み中...
    </div>

    <!-- Store List -->
    <div v-else class="store-grid">
      <div
        v-for="store in stores"
        :key="store.id"
        class="store-card"
        @click="viewStore(store.id)"
      >
        <div class="store-card__header">
          <h3>{{ store.name }}</h3>
          <span
            :class="['status-badge', `status-badge--${getStoreStatusColor(store.status)}`]"
          >
            {{ getStoreStatusLabel(store.status) }}
          </span>
        </div>

        <div class="store-card__content">
          <div v-if="store.address" class="store-info">
            <span class="store-info__label">住所:</span>
            <span class="store-info__value">{{ store.address }}</span>
          </div>
          <div v-if="store.phone" class="store-info">
            <span class="store-info__label">電話:</span>
            <span class="store-info__value">{{ store.phone }}</span>
          </div>
          <div v-if="store.email" class="store-info">
            <span class="store-info__label">メール:</span>
            <span class="store-info__value">{{ store.email }}</span>
          </div>
          <div class="store-info">
            <span class="store-info__label">作成日:</span>
            <span class="store-info__value">{{ formatDate(store.created_at) }}</span>
          </div>
        </div>

        <div class="store-card__actions">
          <button
            class="btn btn--small btn--secondary"
            @click.stop="editStore(store)"
          >
            編集
          </button>
          <button
            class="btn btn--small btn--primary"
            @click.stop="toggleStoreStatus(store)"
          >
            {{ store.status === 'active' ? '休業' : '営業開始' }}
          </button>
          <button
            class="btn btn--small btn--danger"
            @click.stop="confirmDeleteStore(store)"
          >
            削除
          </button>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="!loading && stores.length === 0" class="empty-state">
      <p>店舗が見つかりません。</p>
      <button class="btn btn--primary" @click="showCreateModal = true">
        最初の店舗を登録
      </button>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="pagination">
      <button
        class="pagination__btn"
        :disabled="currentPage === 1"
        @click="setPage(currentPage - 1)"
      >
        前へ
      </button>

      <span class="pagination__info">
        {{ currentPage }} / {{ totalPages }} ページ
      </span>

      <button
        class="pagination__btn"
        :disabled="currentPage === totalPages"
        @click="setPage(currentPage + 1)"
      >
        次へ
      </button>

      <select
        v-model="itemsPerPage"
        class="pagination__size"
        @change="setItemsPerPage(itemsPerPage)"
      >
        <option :value="10">10件</option>
        <option :value="25">25件</option>
        <option :value="50">50件</option>
        <option :value="100">100件</option>
      </select>
    </div>

    <!-- Create Store Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click="closeCreateModal">
      <div class="modal" @click.stop>
        <div class="modal__header">
          <h2>新規店舗登録</h2>
          <button class="modal__close" @click="closeCreateModal">×</button>
        </div>
        <div class="modal__body">
          <StoreForm
            :loading="loading"
            @submit="handleCreateStore"
            @cancel="closeCreateModal"
          />
        </div>
      </div>
    </div>

    <!-- Edit Store Modal -->
    <div v-if="showEditModal && editingStore" class="modal-overlay" @click="closeEditModal">
      <div class="modal" @click.stop>
        <div class="modal__header">
          <h2>店舗編集</h2>
          <button class="modal__close" @click="closeEditModal">×</button>
        </div>
        <div class="modal__body">
          <StoreForm
            :store="editingStore"
            :loading="loading"
            @submit="handleUpdateStore"
            @cancel="closeEditModal"
          />
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal && deletingStore" class="modal-overlay" @click="closeDeleteModal">
      <div class="modal modal--small" @click.stop>
        <div class="modal__header">
          <h2>店舗削除の確認</h2>
        </div>
        <div class="modal__body">
          <p>「{{ deletingStore.name }}」を削除しますか？</p>
          <p class="warning-text">この操作は取り消せません。</p>
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
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useStoreStore, type Store, type StoreCreateRequest, type StoreUpdateRequest } from '@/stores/store'
import StoreForm from '@/components/StoreForm.vue'
import { debounce } from 'lodash-es'

const router = useRouter()
const storeStore = useStoreStore()

const {
  stores,
  loading,
  error,
  currentPage,
  totalPages,
  statusCounts
} = storeToRefs(storeStore)

const {
  fetchStores,
  createStore,
  updateStore,
  deleteStore,
  updateStoreStatus,
  searchStores,
  setPage,
  setItemsPerPage,
  clearError,
  getStoreStatusLabel,
  getStoreStatusColor
} = storeStore

// Local reactive state
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const editingStore = ref<Store | null>(null)
const deletingStore = ref<Store | null>(null)

const searchForm = reactive({
  name: '',
  status: ''
})

const itemsPerPage = ref(10)

// Computed
const debouncedSearch = debounce(() => {
  handleSearch()
}, 500)

// Methods
const handleSearch = async () => {
  const params = {
    ...searchForm,
    page: 1
  }
  await searchStores(params)
}

const clearSearch = () => {
  searchForm.name = ''
  searchForm.status = ''
  fetchStores()
}

const viewStore = (id: number) => {
  router.push(`/stores/${id}`)
}

const editStore = (store: Store) => {
  editingStore.value = store
  showEditModal.value = true
}

const confirmDeleteStore = (store: Store) => {
  deletingStore.value = store
  showDeleteModal.value = true
}

const toggleStoreStatus = async (store: Store) => {
  try {
    const newStatus = store.status === 'active' ? 'inactive' : 'active'
    await updateStoreStatus(store.id, newStatus)
  } catch (err) {
    // エラーはstoreで処理される
  }
}

const closeCreateModal = () => {
  showCreateModal.value = false
}

const closeEditModal = () => {
  showEditModal.value = false
  editingStore.value = null
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  deletingStore.value = null
}

const handleCreateStore = async (storeData: StoreCreateRequest) => {
  try {
    await createStore(storeData)
    closeCreateModal()
  } catch (err) {
    // エラーはstoreで処理される
  }
}

const handleUpdateStore = async (storeData: StoreUpdateRequest) => {
  if (!editingStore.value) return

  try {
    await updateStore(editingStore.value.id, storeData)
    closeEditModal()
  } catch (err) {
    // エラーはstoreで処理される
  }
}

const handleDeleteStore = async () => {
  if (!deletingStore.value) return

  try {
    await deleteStore(deletingStore.value.id)
    closeDeleteModal()
  } catch (err) {
    // エラーはstoreで処理される
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('ja-JP')
}

// Initialize
onMounted(() => {
  fetchStores()
})
</script>

<style scoped>
.store-list {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.store-list__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.store-list__header h1 {
  margin: 0;
  color: #333;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.stats-card {
  background: white;
  padding: 1.5rem;
  border-radius: 0.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.stats-card h3 {
  margin: 0 0 0.5rem 0;
  font-size: 0.875rem;
  color: #666;
  font-weight: 500;
}

.stats-number {
  margin: 0;
  font-size: 2rem;
  font-weight: bold;
  color: #333;
}

.stats-card--primary .stats-number {
  color: #10b981;
}

.stats-card--warning .stats-number {
  color: #f59e0b;
}

.stats-card--error .stats-number {
  color: #ef4444;
}

.search-form {
  background: white;
  padding: 1.5rem;
  border-radius: 0.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
  display: grid;
  grid-template-columns: 1fr 1fr auto;
  gap: 1rem;
  align-items: end;
}

.search-form__group {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.search-form__group label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #374151;
}

.search-form__group input,
.search-form__group select {
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
}

.search-form__group input:focus,
.search-form__group select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.store-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.store-card {
  background: white;
  border-radius: 0.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.store-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.store-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.store-card__header h3 {
  margin: 0;
  color: #333;
  font-size: 1.125rem;
}

.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 500;
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

.store-card__content {
  margin-bottom: 1rem;
}

.store-info {
  display: flex;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
}

.store-info__label {
  font-weight: 500;
  color: #666;
  min-width: 60px;
}

.store-info__value {
  color: #333;
}

.store-card__actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn--small {
  padding: 0.25rem 0.75rem;
  font-size: 0.75rem;
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

.error-message__close {
  background: none;
  border: none;
  color: #991b1b;
  font-size: 1.25rem;
  cursor: pointer;
  padding: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-spinner {
  text-align: center;
  padding: 2rem;
  color: #666;
}

.empty-state {
  text-align: center;
  padding: 3rem 2rem;
  background: white;
  border-radius: 0.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.empty-state p {
  margin-bottom: 1rem;
  color: #666;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 2rem;
}

.pagination__btn {
  padding: 0.5rem 1rem;
  background: white;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  cursor: pointer;
}

.pagination__btn:hover:not(:disabled) {
  background: #f9fafb;
}

.pagination__btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination__info {
  color: #666;
}

.pagination__size {
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
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
  max-width: 500px;
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
}

@media (max-width: 768px) {
  .search-form {
    grid-template-columns: 1fr;
  }

  .store-grid {
    grid-template-columns: 1fr;
  }

  .store-list__header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>