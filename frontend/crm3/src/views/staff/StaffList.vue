<template>
  <div class="staff-list">
    <!-- Header -->
    <div class="staff-list__header">
      <h1>店員管理</h1>
      <button
        class="btn btn--primary"
        @click="showCreateModal = true"
      >
        新規店員登録
      </button>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stats-card stats-card--primary">
        <h3>在籍</h3>
        <p class="stats-number">{{ statusCounts.active }}</p>
      </div>
      <div class="stats-card stats-card--warning">
        <h3>休職</h3>
        <p class="stats-number">{{ statusCounts.inactive }}</p>
      </div>
      <div class="stats-card stats-card--info">
        <h3>休暇中</h3>
        <p class="stats-number">{{ statusCounts.on_leave }}</p>
      </div>
      <div class="stats-card">
        <h3>総店員数</h3>
        <p class="stats-number">{{ statusCounts.total }}</p>
      </div>
    </div>

    <!-- Search Form -->
    <div class="search-form">
      <div class="search-form__group">
        <label for="search-name">氏名</label>
        <input
          id="search-name"
          v-model="searchForm.name"
          type="text"
          placeholder="氏名で検索"
          @input="debouncedSearch"
        />
      </div>
      <div class="search-form__group">
        <label for="search-email">メール</label>
        <input
          id="search-email"
          v-model="searchForm.email"
          type="email"
          placeholder="メールアドレスで検索"
          @input="debouncedSearch"
        />
      </div>
      <div class="search-form__group">
        <label for="search-position">役職</label>
        <input
          id="search-position"
          v-model="searchForm.position"
          type="text"
          placeholder="役職で検索"
          @input="debouncedSearch"
        />
      </div>
      <div class="search-form__group">
        <label for="search-status">ステータス</label>
        <select id="search-status" v-model="searchForm.status" @change="handleSearch">
          <option value="">すべて</option>
          <option value="active">在籍</option>
          <option value="inactive">休職</option>
          <option value="on_leave">休暇中</option>
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

    <!-- Staff List -->
    <div v-else class="staff-grid">
      <div
        v-for="member in staff"
        :key="member.id"
        class="staff-card"
        @click="viewStaff(member.id)"
      >
        <div class="staff-card__header">
          <div class="staff-card__name">
            <h3>{{ member.name }}</h3>
            <span class="staff-card__position">{{ member.position }}</span>
          </div>
          <span
            :class="['status-badge', `status-badge--${getStaffStatusColor(member.status)}`]"
          >
            {{ getStaffStatusLabel(member.status) }}
          </span>
        </div>

        <div class="staff-card__content">
          <div class="staff-info">
            <span class="staff-info__label">メール:</span>
            <span class="staff-info__value">{{ member.email }}</span>
          </div>
          <div v-if="member.phone" class="staff-info">
            <span class="staff-info__label">電話:</span>
            <span class="staff-info__value">{{ member.phone }}</span>
          </div>
          <div class="staff-info">
            <span class="staff-info__label">入社日:</span>
            <span class="staff-info__value">{{ formatDate(member.hire_date) }}</span>
          </div>
          <div v-if="member.store" class="staff-info">
            <span class="staff-info__label">所属店舗:</span>
            <span class="staff-info__value">{{ member.store.name }}</span>
          </div>
          <div v-else class="staff-info">
            <span class="staff-info__label">所属店舗:</span>
            <span class="staff-info__value staff-info__value--unassigned">未割り当て</span>
          </div>
        </div>

        <div class="staff-card__actions">
          <button
            class="btn btn--small btn--secondary"
            @click.stop="editStaff(member)"
          >
            編集
          </button>
          <button
            v-if="member.status === 'active'"
            class="btn btn--small btn--warning"
            @click.stop="updateStatus(member, 'inactive')"
          >
            休職
          </button>
          <button
            v-else-if="member.status === 'inactive'"
            class="btn btn--small btn--primary"
            @click.stop="updateStatus(member, 'active')"
          >
            復職
          </button>
          <button
            v-if="member.status !== 'on_leave'"
            class="btn btn--small btn--info"
            @click.stop="updateStatus(member, 'on_leave')"
          >
            休暇
          </button>
          <button
            class="btn btn--small btn--danger"
            @click.stop="confirmDeleteStaff(member)"
          >
            削除
          </button>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="!loading && staff.length === 0" class="empty-state">
      <p>店員が見つかりません。</p>
      <button class="btn btn--primary" @click="showCreateModal = true">
        最初の店員を登録
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

    <!-- Create Staff Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click="closeCreateModal">
      <div class="modal" @click.stop>
        <div class="modal__header">
          <h2>新規店員登録</h2>
          <button class="modal__close" @click="closeCreateModal">×</button>
        </div>
        <div class="modal__body">
          <StaffForm
            :loading="loading"
            @submit="handleCreateStaff"
            @cancel="closeCreateModal"
          />
        </div>
      </div>
    </div>

    <!-- Edit Staff Modal -->
    <div v-if="showEditModal && editingStaff" class="modal-overlay" @click="closeEditModal">
      <div class="modal" @click.stop>
        <div class="modal__header">
          <h2>店員編集</h2>
          <button class="modal__close" @click="closeEditModal">×</button>
        </div>
        <div class="modal__body">
          <StaffForm
            :staff="editingStaff"
            :loading="loading"
            @submit="handleUpdateStaff"
            @cancel="closeEditModal"
          />
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal && deletingStaff" class="modal-overlay" @click="closeDeleteModal">
      <div class="modal modal--small" @click.stop>
        <div class="modal__header">
          <h2>店員削除の確認</h2>
        </div>
        <div class="modal__body">
          <p>「{{ deletingStaff.name }}」を削除しますか？</p>
          <p class="warning-text">この操作は取り消せません。</p>
        </div>
        <div class="modal__footer">
          <button class="btn btn--secondary" @click="closeDeleteModal">
            キャンセル
          </button>
          <button class="btn btn--danger" @click="handleDeleteStaff" :disabled="loading">
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
import { useStaffStore, type Staff, type StaffCreateRequest, type StaffUpdateRequest } from '@/stores/staff'
import StaffForm from '@/components/StaffForm.vue'
import { debounce } from 'lodash-es'

const router = useRouter()
const staffStore = useStaffStore()

const {
  staff,
  loading,
  error,
  currentPage,
  totalPages,
  statusCounts
} = storeToRefs(staffStore)

const {
  fetchStaff,
  createStaffMember,
  updateStaffMember,
  deleteStaffMember,
  updateStaffStatus,
  searchStaff,
  setPage,
  setItemsPerPage,
  clearError,
  getStaffStatusLabel,
  getStaffStatusColor
} = staffStore

// Local reactive state
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const editingStaff = ref<Staff | null>(null)
const deletingStaff = ref<Staff | null>(null)

const searchForm = reactive({
  name: '',
  email: '',
  position: '',
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
  await searchStaff(params)
}

const clearSearch = () => {
  searchForm.name = ''
  searchForm.email = ''
  searchForm.position = ''
  searchForm.status = ''
  fetchStaff()
}

const viewStaff = (id: number) => {
  router.push(`/staff/${id}`)
}

const editStaff = (member: Staff) => {
  editingStaff.value = member
  showEditModal.value = true
}

const confirmDeleteStaff = (member: Staff) => {
  deletingStaff.value = member
  showDeleteModal.value = true
}

const updateStatus = async (member: Staff, status: Staff['status']) => {
  try {
    await updateStaffStatus(member.id, status)
  } catch (err) {
    // エラーはstoreで処理される
  }
}

const closeCreateModal = () => {
  showCreateModal.value = false
}

const closeEditModal = () => {
  showEditModal.value = false
  editingStaff.value = null
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  deletingStaff.value = null
}

const handleCreateStaff = async (staffData: StaffCreateRequest) => {
  try {
    await createStaffMember(staffData)
    closeCreateModal()
  } catch (err) {
    // エラーはstoreで処理される
  }
}

const handleUpdateStaff = async (staffData: StaffUpdateRequest) => {
  if (!editingStaff.value) return

  try {
    await updateStaffMember(editingStaff.value.id, staffData)
    closeEditModal()
  } catch (err) {
    // エラーはstoreで処理される
  }
}

const handleDeleteStaff = async () => {
  if (!deletingStaff.value) return

  try {
    await deleteStaffMember(deletingStaff.value.id)
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
  fetchStaff()
})
</script>

<style scoped>
.staff-list {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.staff-list__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.staff-list__header h1 {
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

.stats-card--info .stats-number {
  color: #3b82f6;
}

.search-form {
  background: white;
  padding: 1.5rem;
  border-radius: 0.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
  display: grid;
  grid-template-columns: 1fr 1fr 1fr 1fr auto;
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

.staff-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.staff-card {
  background: white;
  border-radius: 0.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.staff-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.staff-card__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
}

.staff-card__name h3 {
  margin: 0 0 0.25rem 0;
  color: #333;
  font-size: 1.125rem;
}

.staff-card__position {
  font-size: 0.875rem;
  color: #666;
  font-weight: 500;
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

.status-badge--info {
  background: #dbeafe;
  color: #1e40af;
}

.staff-card__content {
  margin-bottom: 1rem;
}

.staff-info {
  display: flex;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
}

.staff-info__label {
  font-weight: 500;
  color: #666;
  min-width: 80px;
}

.staff-info__value {
  color: #333;
}

.staff-info__value--unassigned {
  color: #ef4444;
  font-style: italic;
}

.staff-card__actions {
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

.btn--warning {
  background: #f59e0b;
  color: white;
}

.btn--warning:hover {
  background: #d97706;
}

.btn--info {
  background: #3b82f6;
  color: white;
}

.btn--info:hover {
  background: #2563eb;
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

  .staff-grid {
    grid-template-columns: 1fr;
  }

  .staff-list__header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>