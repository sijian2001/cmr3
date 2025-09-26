<template>
  <div class="customer-detail">
    <div class="header">
      <button @click="goBack" class="btn-back">← 戻る</button>
      <h1>顧客詳細</h1>
      <div class="actions">
        <button @click="editCustomer" class="btn-primary">編集</button>
        <button @click="confirmDelete" class="btn-danger">削除</button>
      </div>
    </div>

    <!-- エラー表示 -->
    <div v-if="customerStore.hasError" class="error-message">
      {{ customerStore.error }}
    </div>

    <!-- ローディング -->
    <div v-if="customerStore.isLoading" class="loading">
      読み込み中...
    </div>

    <!-- 顧客情報表示 -->
    <div v-else-if="customerStore.currentCustomer" class="customer-info">
      <div class="info-section">
        <h2>基本情報</h2>
        <div class="info-grid">
          <div class="info-item">
            <label>ID</label>
            <div class="value">{{ customerStore.currentCustomer.id }}</div>
          </div>
          <div class="info-item">
            <label>名前</label>
            <div class="value">{{ customerStore.currentCustomer.name }}</div>
          </div>
          <div class="info-item">
            <label>メールアドレス</label>
            <div class="value">{{ customerStore.currentCustomer.email }}</div>
          </div>
          <div class="info-item">
            <label>電話番号</label>
            <div class="value">{{ customerStore.currentCustomer.phone || '-' }}</div>
          </div>
          <div class="info-item full-width">
            <label>住所</label>
            <div class="value">{{ customerStore.currentCustomer.address || '-' }}</div>
          </div>
          <div class="info-item">
            <label>登録日</label>
            <div class="value">{{ formatDate(customerStore.currentCustomer.created_at) }}</div>
          </div>
          <div class="info-item">
            <label>更新日</label>
            <div class="value">{{ formatDate(customerStore.currentCustomer.updated_at) }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 顧客が見つからない場合 -->
    <div v-else class="not-found">
      <h2>顧客が見つかりません</h2>
      <p>指定された顧客が存在しないか、削除された可能性があります。</p>
    </div>

    <!-- 編集モーダル -->
    <CustomerFormModal
      v-if="showEditModal && customerStore.currentCustomer"
      :customer="customerStore.currentCustomer"
      @save="handleSave"
      @cancel="showEditModal = false"
    />

    <!-- 削除確認モーダル -->
    <ConfirmModal
      v-if="showDeleteModal && customerStore.currentCustomer"
      title="顧客削除の確認"
      :message="`「${customerStore.currentCustomer.name}」を削除してもよろしいですか？この操作は取り消せません。`"
      @confirm="handleDelete"
      @cancel="showDeleteModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCustomerStore, type Customer } from '@/stores/customer'
import CustomerFormModal from '@/components/customers/CustomerFormModal.vue'
import ConfirmModal from '@/components/common/ConfirmModal.vue'

const route = useRoute()
const router = useRouter()
const customerStore = useCustomerStore()

const showEditModal = ref(false)
const showDeleteModal = ref(false)

// メソッド
const formatDate = (dateString?: string) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('ja-JP')
}

const goBack = () => {
  router.push('/customers')
}

const editCustomer = () => {
  showEditModal.value = true
}

const confirmDelete = () => {
  showDeleteModal.value = true
}

const handleSave = async (customerData: Omit<Customer, 'id' | 'created_at' | 'updated_at'>) => {
  try {
    if (customerStore.currentCustomer?.id) {
      await customerStore.updateCustomerData(customerStore.currentCustomer.id, customerData)
      showEditModal.value = false
      // 詳細情報を再取得
      await customerStore.fetchCustomer(customerStore.currentCustomer.id)
    }
  } catch (error) {
    console.error('更新エラー:', error)
  }
}

const handleDelete = async () => {
  if (customerStore.currentCustomer?.id) {
    try {
      await customerStore.deleteCustomer(customerStore.currentCustomer.id)
      showDeleteModal.value = false
      // 一覧画面に戻る
      router.push('/customers')
    } catch (error) {
      console.error('削除エラー:', error)
    }
  }
}

// ライフサイクル
onMounted(async () => {
  const customerId = Number(route.params.id)
  if (customerId) {
    try {
      await customerStore.fetchCustomer(customerId)
    } catch (error) {
      console.error('顧客取得エラー:', error)
    }
  }
})
</script>

<style scoped>
.customer-detail {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 15px;
  border-bottom: 1px solid #eee;
}

.header h1 {
  margin: 0;
  color: #333;
  flex: 1;
  text-align: center;
}

.btn-back {
  background: transparent;
  color: #007bff;
  border: 1px solid #007bff;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  text-decoration: none;
}

.btn-back:hover {
  background: #007bff;
  color: white;
}

.actions {
  display: flex;
  gap: 10px;
}

.error-message {
  background: #fee;
  border: 1px solid #fcc;
  color: #c00;
  padding: 15px;
  border-radius: 4px;
  margin-bottom: 20px;
  text-align: center;
}

.loading {
  text-align: center;
  padding: 60px;
  color: #666;
  font-size: 16px;
}

.customer-info {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  overflow: hidden;
}

.info-section {
  padding: 30px;
}

.info-section h2 {
  margin: 0 0 20px 0;
  color: #333;
  font-size: 20px;
  padding-bottom: 10px;
  border-bottom: 2px solid #007bff;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.info-item {
  display: flex;
  flex-direction: column;
}

.info-item.full-width {
  grid-column: 1 / -1;
}

.info-item label {
  font-weight: bold;
  color: #555;
  margin-bottom: 5px;
  font-size: 14px;
}

.info-item .value {
  color: #333;
  font-size: 16px;
  padding: 8px 0;
  border-bottom: 1px solid #eee;
  min-height: 24px;
}

.not-found {
  text-align: center;
  padding: 60px 20px;
  color: #666;
}

.not-found h2 {
  color: #333;
  margin-bottom: 15px;
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

.btn-danger {
  background: #dc3545;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-danger:hover {
  background: #c82333;
}

/* レスポンシブ対応 */
@media (max-width: 768px) {
  .customer-detail {
    padding: 10px;
  }

  .header {
    flex-direction: column;
    gap: 15px;
  }

  .header h1 {
    text-align: center;
  }

  .actions {
    justify-content: center;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .info-section {
    padding: 20px;
  }
}
</style>