<template>
  <div class="product-list">
    <div class="header">
      <h1>製品管理</h1>
      <button @click="showCreateModal = true" class="btn-primary">
        新規製品登録
      </button>
    </div>

    <!-- 検索フォーム -->
    <div class="search-form">
      <div class="form-row">
        <div class="form-group">
          <label>製品名</label>
          <input
            v-model="searchParams.name"
            type="text"
            placeholder="製品名で検索"
            @keyup.enter="handleSearch"
          >
        </div>
        <div class="form-group">
          <label>説明</label>
          <input
            v-model="searchParams.description"
            type="text"
            placeholder="説明で検索"
            @keyup.enter="handleSearch"
          >
        </div>
        <div class="form-group">
          <label>SKU</label>
          <input
            v-model="searchParams.sku"
            type="text"
            placeholder="SKUで検索"
            @keyup.enter="handleSearch"
          >
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label>最小価格</label>
          <input
            v-model.number="searchParams.min_price"
            type="number"
            min="0"
            placeholder="最小価格"
            @keyup.enter="handleSearch"
          >
        </div>
        <div class="form-group">
          <label>最大価格</label>
          <input
            v-model.number="searchParams.max_price"
            type="number"
            min="0"
            placeholder="最大価格"
            @keyup.enter="handleSearch"
          >
        </div>
        <div class="form-group">
          <!-- 空のスペース -->
        </div>
      </div>
      <div class="search-actions">
        <button @click="handleSearch" class="btn-secondary">検索</button>
        <button @click="clearSearch" class="btn-outline">クリア</button>
      </div>
    </div>

    <!-- エラー表示 -->
    <div v-if="productStore.hasError" class="error-message">
      {{ productStore.error }}
      <button @click="productStore.clearError" class="btn-close">×</button>
    </div>

    <!-- ローディング -->
    <div v-if="productStore.isLoading" class="loading">
      読み込み中...
    </div>

    <!-- 製品一覧テーブル -->
    <div v-else class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>製品名</th>
            <th>SKU</th>
            <th>価格</th>
            <th>在庫数量</th>
            <th>登録日</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="productStore.isEmpty">
            <td colspan="7" class="no-data">製品データがありません</td>
          </tr>
          <tr v-else v-for="product in productStore.products" :key="product.id">
            <td>{{ product.id }}</td>
            <td>{{ product.name }}</td>
            <td>{{ product.sku }}</td>
            <td class="price">¥{{ formatPrice(product.price) }}</td>
            <td :class="getStockClass(product.stock_quantity)">
              {{ product.stock_quantity }}
            </td>
            <td>{{ formatDate(product.created_at) }}</td>
            <td class="actions">
              <button @click="viewProduct(product.id!)" class="btn-sm btn-info">詳細</button>
              <button @click="editProduct(product)" class="btn-sm btn-warning">編集</button>
              <button @click="showStockModal(product)" class="btn-sm btn-success">在庫</button>
              <button @click="confirmDelete(product)" class="btn-sm btn-danger">削除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- ページネーション -->
    <div v-if="!productStore.isEmpty" class="pagination">
      <div class="pagination-info">
        全 {{ productStore.pagination.total }} 件中
        {{ (productStore.pagination.page - 1) * productStore.pagination.limit + 1 }} -
        {{ Math.min(productStore.pagination.page * productStore.pagination.limit, productStore.pagination.total) }} 件を表示
      </div>
      <div class="pagination-controls">
        <button
          @click="changePage(productStore.pagination.page - 1)"
          :disabled="productStore.pagination.page <= 1"
          class="btn-outline"
        >
          前へ
        </button>
        <span class="page-info">
          {{ productStore.pagination.page }} / {{ productStore.pagination.total_pages }}
        </span>
        <button
          @click="changePage(productStore.pagination.page + 1)"
          :disabled="productStore.pagination.page >= productStore.pagination.total_pages"
          class="btn-outline"
        >
          次へ
        </button>
      </div>
    </div>

    <!-- 新規作成モーダル -->
    <ProductFormModal
      v-if="showCreateModal"
      :product="null"
      @save="handleSave"
      @cancel="showCreateModal = false"
    />

    <!-- 編集モーダル -->
    <ProductFormModal
      v-if="showEditModal && editingProduct"
      :product="editingProduct"
      @save="handleSave"
      @cancel="showEditModal = false"
    />

    <!-- 在庫調整モーダル -->
    <StockAdjustmentModal
      v-if="showStockAdjustmentModal && stockProduct"
      :product="stockProduct"
      @adjust="handleStockAdjustment"
      @cancel="showStockAdjustmentModal = false"
    />

    <!-- 削除確認モーダル -->
    <ConfirmModal
      v-if="showDeleteModal && deletingProduct"
      title="製品削除の確認"
      :message="`「${deletingProduct.name}」を削除してもよろしいですか？この操作は取り消せません。`"
      @confirm="handleDelete"
      @cancel="showDeleteModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useProductStore, type Product, type ProductSearchParams } from '@/stores/product'
import ProductFormModal from '@/components/products/ProductFormModal.vue'
import StockAdjustmentModal from '@/components/products/StockAdjustmentModal.vue'
import ConfirmModal from '@/components/common/ConfirmModal.vue'

const router = useRouter()
const productStore = useProductStore()

// リアクティブデータ
const searchParams = reactive<ProductSearchParams>({
  name: '',
  description: '',
  sku: '',
  min_price: undefined,
  max_price: undefined,
  page: 1,
  limit: 10
})

const showCreateModal = ref(false)
const showEditModal = ref(false)
const showStockAdjustmentModal = ref(false)
const showDeleteModal = ref(false)
const editingProduct = ref<Product | null>(null)
const stockProduct = ref<Product | null>(null)
const deletingProduct = ref<Product | null>(null)

// メソッド
const formatDate = (dateString?: string) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleDateString('ja-JP')
}

const formatPrice = (price: number) => {
  return price.toLocaleString('ja-JP')
}

const getStockClass = (quantity: number) => {
  if (quantity === 0) return 'stock-empty'
  if (quantity <= 10) return 'stock-low'
  return 'stock-normal'
}

const handleSearch = async () => {
  searchParams.page = 1
  await productStore.fetchProducts({ ...searchParams })
}

const clearSearch = async () => {
  searchParams.name = ''
  searchParams.description = ''
  searchParams.sku = ''
  searchParams.min_price = undefined
  searchParams.max_price = undefined
  searchParams.page = 1
  await productStore.fetchProducts({ ...searchParams })
}

const changePage = async (page: number) => {
  if (page >= 1 && page <= productStore.pagination.total_pages) {
    searchParams.page = page
    await productStore.fetchProducts({ ...searchParams })
  }
}

const viewProduct = (id: number) => {
  router.push(`/products/${id}`)
}

const editProduct = (product: Product) => {
  editingProduct.value = { ...product }
  showEditModal.value = true
}

const showStockModal = (product: Product) => {
  stockProduct.value = product
  showStockAdjustmentModal.value = true
}

const confirmDelete = (product: Product) => {
  deletingProduct.value = product
  showDeleteModal.value = true
}

const handleSave = async (productData: Omit<Product, 'id' | 'created_at' | 'updated_at'>) => {
  try {
    if (showCreateModal.value) {
      await productStore.createProduct(productData)
      showCreateModal.value = false
    } else if (showEditModal.value && editingProduct.value?.id) {
      await productStore.updateProductData(editingProduct.value.id, productData)
      showEditModal.value = false
      editingProduct.value = null
    }
    // 一覧を更新
    await productStore.fetchProducts({ ...searchParams })
  } catch (error) {
    console.error('保存エラー:', error)
  }
}

const handleStockAdjustment = async (quantity: number, operation: 'add' | 'subtract') => {
  if (stockProduct.value?.id) {
    try {
      await productStore.adjustStock(stockProduct.value.id, quantity, operation)
      showStockAdjustmentModal.value = false
      stockProduct.value = null
      // 一覧を更新
      await productStore.fetchProducts({ ...searchParams })
    } catch (error) {
      console.error('在庫調整エラー:', error)
    }
  }
}

const handleDelete = async () => {
  if (deletingProduct.value?.id) {
    try {
      await productStore.deleteProduct(deletingProduct.value.id)
      showDeleteModal.value = false
      deletingProduct.value = null
      // 一覧を更新
      await productStore.fetchProducts({ ...searchParams })
    } catch (error) {
      console.error('削除エラー:', error)
    }
  }
}

// ライフサイクル
onMounted(async () => {
  await productStore.fetchProducts(searchParams)
})
</script>

<style scoped>
.product-list {
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

.price {
  text-align: right;
  font-weight: bold;
}

.stock-normal {
  color: #28a745;
  font-weight: bold;
}

.stock-low {
  color: #ffc107;
  font-weight: bold;
}

.stock-empty {
  color: #dc3545;
  font-weight: bold;
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

.btn-success {
  background: #28a745;
  color: white;
}

.btn-danger {
  background: #dc3545;
  color: white;
}
</style>