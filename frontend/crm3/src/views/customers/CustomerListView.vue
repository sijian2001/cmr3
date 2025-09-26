<template>
  <div class="customer-list">
    <div class="header">
      <h1>顧客管理</h1>
      <button @click="showCreateModal = true" class="btn-primary">
        新規顧客登録
      </button>
    </div>

    <!-- 検索フォーム -->
    <div class="search-form">
      <div class="form-row">
        <div class="form-group">
          <label>名前</label>
          <input
            v-model="searchParams.name"
            type="text"
            placeholder="名前で検索"
            @keyup.enter="handleSearch"
          >
        </div>
        <div class="form-group">
          <label>メールアドレス</label>
          <input
            v-model="searchParams.email"
            type="email"
            placeholder="メールアドレスで検索"
            @keyup.enter="handleSearch"
          >
        </div>
        <div class="form-group">
          <label>電話番号</label>
          <input
            v-model="searchParams.phone"
            type="tel"
            placeholder="電話番号で検索"
            @keyup.enter="handleSearch"
          >
        </div>
      </div>
      <div class="search-actions">
        <button @click="handleSearch" class="btn-secondary">検索</button>
        <button @click="clearSearch" class="btn-outline">クリア</button>
      </div>
    </div>

    <!-- エラー表示 -->
    <div v-if="customerStore.hasError" class="error-message">
      {{ customerStore.error }}
      <button @click="customerStore.clearError" class="btn-close">×</button>
    </div>

    <!-- ローディング -->
    <div v-if="customerStore.isLoading" class="loading">
      読み込み中...
    </div>

    <!-- 顧客一覧テーブル -->
    <div v-else class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>名前</th>
            <th>メールアドレス</th>
            <th>電話番号</th>
            <th>登録日</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="customerStore.isEmpty">
            <td colspan="6" class="no-data">顧客データがありません</td>
          </tr>
          <tr v-else v-for="customer in customerStore.customers" :key="customer.id">
            <td>{{ customer.id }}</td>
            <td>{{ customer.name }}</td>
            <td>{{ customer.email }}</td>
            <td>{{ customer.phone || '-' }}</td>
            <td>{{ formatDate(customer.created_at) }}</td>
            <td class="actions">
              <button @click="viewCustomer(customer.id!)" class="btn-sm btn-info">詳細</button>
              <button @click="editCustomer(customer)" class="btn-sm btn-warning">編集</button>
              <button @click="confirmDelete(customer)" class="btn-sm btn-danger">削除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- ページネーション -->
    <div v-if="!customerStore.isEmpty" class="pagination">
      <div class="pagination-info">
        全 {{ customerStore.pagination.total }} 件中
        {{ (customerStore.pagination.page - 1) * customerStore.pagination.limit + 1 }} -
        {{ Math.min(customerStore.pagination.page * customerStore.pagination.limit, customerStore.pagination.total) }} 件を表示
      </div>
      <div class="pagination-controls">
        <button
          @click="changePage(customerStore.pagination.page - 1)"
          :disabled="customerStore.pagination.page <= 1"
          class="btn-outline"
        >
          前へ
        </button>
        <span class="page-info">
          {{ customerStore.pagination.page }} / {{ customerStore.pagination.total_pages }}
        </span>
        <button
          @click="changePage(customerStore.pagination.page + 1)"
          :disabled="customerStore.pagination.page >= customerStore.pagination.total_pages"
          class="btn-outline"
        >
          次へ
        </button>
      </div>
    </div>

    <!-- 新規作成モーダル -->
    <CustomerFormModal
      v-if="showCreateModal"
      :customer="null"
      @save="handleSave"
      @cancel="showCreateModal = false"
    />

    <!-- 編集モーダル -->
    <CustomerFormModal
      v-if="showEditModal && editingCustomer"
      :customer="editingCustomer"
      @save="handleSave"
      @cancel="showEditModal = false"
    />

    <!-- 削除確認モーダル -->
    <ConfirmModal
      v-if="showDeleteModal && deletingCustomer"
      title="顧客削除の確認"
      :message="`「${deletingCustomer.name}」を削除してもよろしいですか？この操作は取り消せません。`"
      @confirm="handleDelete"
      @cancel="showDeleteModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCustomerStore, type Customer, type CustomerSearchParams } from '@/stores/customer'
import CustomerFormModal from '@/components/customers/CustomerFormModal.vue'
import ConfirmModal from '@/components/common/ConfirmModal.vue'

const router = useRouter()
const customerStore = useCustomerStore()

// リアクティブデータ
const searchParams = reactive<CustomerSearchParams>({
  name: '',
  email: '',
  phone: '',
  page: 1,
  limit: 10
})

const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const editingCustomer = ref<Customer | null>(null)
const deletingCustomer = ref<Customer | null>(null)

// メソッド
const formatDate = (dateString?: string) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleDateString('ja-JP')
}

const handleSearch = async () => {
  searchParams.page = 1
  await customerStore.fetchCustomers({ ...searchParams })
}

const clearSearch = async () => {
  searchParams.name = ''
  searchParams.email = ''
  searchParams.phone = ''
  searchParams.page = 1
  await customerStore.fetchCustomers({ ...searchParams })
}

const changePage = async (page: number) => {
  if (page >= 1 && page <= customerStore.pagination.total_pages) {
    searchParams.page = page
    await customerStore.fetchCustomers({ ...searchParams })
  }
}

const viewCustomer = (id: number) => {
  router.push(`/customers/${id}`)
}

const editCustomer = (customer: Customer) => {
  editingCustomer.value = { ...customer }
  showEditModal.value = true
}

const confirmDelete = (customer: Customer) => {
  deletingCustomer.value = customer
  showDeleteModal.value = true
}

const handleSave = async (customerData: Omit<Customer, 'id' | 'created_at' | 'updated_at'>) => {
  try {
    if (showCreateModal.value) {
      await customerStore.createCustomer(customerData)
      showCreateModal.value = false
    } else if (showEditModal.value && editingCustomer.value?.id) {
      await customerStore.updateCustomerData(editingCustomer.value.id, customerData)
      showEditModal.value = false
      editingCustomer.value = null
    }
    // 一覧を更新
    await customerStore.fetchCustomers({ ...searchParams })
  } catch (error) {
    console.error('保存エラー:', error)
  }
}

const handleDelete = async () => {
  if (deletingCustomer.value?.id) {
    try {
      await customerStore.deleteCustomer(deletingCustomer.value.id)
      showDeleteModal.value = false
      deletingCustomer.value = null
      // 一覧を更新
      await customerStore.fetchCustomers({ ...searchParams })
    } catch (error) {
      console.error('削除エラー:', error)
    }
  }
}

// ライフサイクル
onMounted(async () => {
  await customerStore.fetchCustomers(searchParams)
})
</script>

<style scoped>
.customer-list {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h1 {
  margin: 0;
  color: #333;
}

.search-form {
  background: #f9f9f9;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 15px;
  margin-bottom: 15px;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-group label {
  font-weight: bold;
  margin-bottom: 5px;
  color: #555;
}

.form-group input {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.search-actions {
  display: flex;
  gap: 10px;
}

.error-message {
  background: #fee;
  border: 1px solid #fcc;
  color: #c00;
  padding: 10px;
  border-radius: 4px;
  margin-bottom: 15px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.btn-close {
  background: none;
  border: none;
  color: #c00;
  font-size: 18px;
  cursor: pointer;
  padding: 0;
  width: 20px;
  height: 20px;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #666;
}

.table-container {
  overflow-x: auto;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.data-table th {
  background: #f8f9fa;
  font-weight: bold;
  color: #555;
}

.data-table tr:hover {
  background: #f5f5f5;
}

.no-data {
  text-align: center;
  color: #666;
  font-style: italic;
}

.actions {
  white-space: nowrap;
}

.actions button {
  margin-right: 5px;
}

.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 20px;
  padding: 15px 0;
}

.pagination-info {
  color: #666;
  font-size: 14px;
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 15px;
}

.page-info {
  font-weight: bold;
  color: #333;
}

/* ボタンスタイル */
.btn-primary {
  background: #007bff;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-primary:hover {
  background: #0056b3;
}

.btn-secondary {
  background: #6c757d;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-outline {
  background: transparent;
  color: #6c757d;
  border: 1px solid #6c757d;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-outline:hover {
  background: #6c757d;
  color: white;
}

.btn-outline:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-sm {
  padding: 4px 8px;
  font-size: 12px;
  border: none;
  border-radius: 3px;
  cursor: pointer;
}

.btn-info {
  background: #17a2b8;
  color: white;
}

.btn-warning {
  background: #ffc107;
  color: #212529;
}

.btn-danger {
  background: #dc3545;
  color: white;
}
</style>