<template>
  <div class="product-detail">
    <div class="header">
      <button @click="goBack" class="btn-back">← 戻る</button>
      <h1>製品詳細</h1>
      <div class="actions">
        <button @click="editProduct" class="btn-primary">編集</button>
        <button @click="showStockModal" class="btn-success">在庫調整</button>
        <button @click="confirmDelete" class="btn-danger">削除</button>
      </div>
    </div>

    <!-- エラー表示 -->
    <div v-if="productStore.hasError" class="error-message">
      {{ productStore.error }}
    </div>

    <!-- ローディング -->
    <div v-if="productStore.isLoading" class="loading">
      読み込み中...
    </div>

    <!-- 製品情報表示 -->
    <div v-else-if="productStore.currentProduct" class="product-info">
      <div class="info-section">
        <h2>基本情報</h2>
        <div class="info-grid">
          <div class="info-item">
            <label>ID</label>
            <div class="value">{{ productStore.currentProduct.id }}</div>
          </div>
          <div class="info-item">
            <label>製品名</label>
            <div class="value">{{ productStore.currentProduct.name }}</div>
          </div>
          <div class="info-item">
            <label>SKU</label>
            <div class="value">{{ productStore.currentProduct.sku }}</div>
          </div>
          <div class="info-item">
            <label>価格</label>
            <div class="value price">¥{{ formatPrice(productStore.currentProduct.price) }}</div>
          </div>
          <div class="info-item">
            <label>在庫数量</label>
            <div class="value" :class="getStockClass(productStore.currentProduct.stock_quantity)">
              {{ productStore.currentProduct.stock_quantity }}
              <span class="stock-status">({{ getStockStatus(productStore.currentProduct.stock_quantity) }})</span>
            </div>
          </div>
          <div class="info-item">
            <label>在庫評価額</label>
            <div class="value">¥{{ formatPrice(productStore.currentProduct.price * productStore.currentProduct.stock_quantity) }}</div>
          </div>
          <div class="info-item full-width">
            <label>説明</label>
            <div class="value description">{{ productStore.currentProduct.description || '-' }}</div>
          </div>
          <div class="info-item">
            <label>登録日</label>
            <div class="value">{{ formatDate(productStore.currentProduct.created_at) }}</div>
          </div>
          <div class="info-item">
            <label>更新日</label>
            <div class="value">{{ formatDate(productStore.currentProduct.updated_at) }}</div>
          </div>
        </div>
      </div>

      <!-- 在庫状況セクション -->
      <div class="info-section">
        <h2>在庫状況</h2>
        <div class="stock-overview">
          <div class="stock-card" :class="getStockClass(productStore.currentProduct.stock_quantity)">
            <div class="stock-number">{{ productStore.currentProduct.stock_quantity }}</div>
            <div class="stock-label">現在在庫</div>
          </div>
          <div class="stock-actions">
            <button @click="quickStockAdjustment('add', 1)" class="btn-sm btn-success">+1</button>
            <button @click="quickStockAdjustment('add', 10)" class="btn-sm btn-success">+10</button>
            <button @click="quickStockAdjustment('subtract', 1)" :disabled="productStore.currentProduct.stock_quantity < 1" class="btn-sm btn-warning">-1</button>
            <button @click="quickStockAdjustment('subtract', 10)" :disabled="productStore.currentProduct.stock_quantity < 10" class="btn-sm btn-warning">-10</button>
            <button @click="showStockModal" class="btn-sm btn-info">詳細調整</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 製品が見つからない場合 -->
    <div v-else class="not-found">
      <h2>製品が見つかりません</h2>
      <p>指定された製品が存在しないか、削除された可能性があります。</p>
    </div>

    <!-- 編集モーダル -->
    <ProductFormModal
      v-if="showEditModal && productStore.currentProduct"
      :product="productStore.currentProduct"
      @save="handleSave"
      @cancel="showEditModal = false"
    />

    <!-- 在庫調整モーダル -->
    <StockAdjustmentModal
      v-if="showStockAdjustmentModal && productStore.currentProduct"
      :product="productStore.currentProduct"
      @adjust="handleStockAdjustment"
      @cancel="showStockAdjustmentModal = false"
    />

    <!-- 削除確認モーダル -->
    <ConfirmModal
      v-if="showDeleteModal && productStore.currentProduct"
      title="製品削除の確認"
      :message="`「${productStore.currentProduct.name}」を削除してもよろしいですか？この操作は取り消せません。`"
      @confirm="handleDelete"
      @cancel="showDeleteModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProductStore, type Product } from '@/stores/product'
import ProductFormModal from '@/components/products/ProductFormModal.vue'
import StockAdjustmentModal from '@/components/products/StockAdjustmentModal.vue'
import ConfirmModal from '@/components/common/ConfirmModal.vue'

const route = useRoute()
const router = useRouter()
const productStore = useProductStore()

const showEditModal = ref(false)
const showStockAdjustmentModal = ref(false)
const showDeleteModal = ref(false)

// メソッド
const formatDate = (dateString?: string) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('ja-JP')
}

const formatPrice = (price: number) => {
  return price.toLocaleString('ja-JP')
}

const getStockClass = (quantity: number) => {
  if (quantity === 0) return 'stock-empty'
  if (quantity <= 10) return 'stock-low'
  return 'stock-normal'
}

const getStockStatus = (quantity: number) => {
  if (quantity === 0) return '在庫切れ'
  if (quantity <= 10) return '在庫少'
  return '在庫十分'
}

const goBack = () => {
  router.push('/products')
}

const editProduct = () => {
  showEditModal.value = true
}

const showStockModal = () => {
  showStockAdjustmentModal.value = true
}

const confirmDelete = () => {
  showDeleteModal.value = true
}

const quickStockAdjustment = async (operation: 'add' | 'subtract', quantity: number) => {
  if (productStore.currentProduct?.id) {
    try {
      await productStore.adjustStock(productStore.currentProduct.id, quantity, operation)
      // 詳細情報を再取得
      await productStore.fetchProduct(productStore.currentProduct.id)
    } catch (error) {
      console.error('在庫調整エラー:', error)
    }
  }
}

const handleSave = async (productData: Omit<Product, 'id' | 'created_at' | 'updated_at'>) => {
  try {
    if (productStore.currentProduct?.id) {
      await productStore.updateProductData(productStore.currentProduct.id, productData)
      showEditModal.value = false
      // 詳細情報を再取得
      await productStore.fetchProduct(productStore.currentProduct.id)
    }
  } catch (error) {
    console.error('更新エラー:', error)
  }
}

const handleStockAdjustment = async (quantity: number, operation: 'add' | 'subtract') => {
  if (productStore.currentProduct?.id) {
    try {
      await productStore.adjustStock(productStore.currentProduct.id, quantity, operation)
      showStockAdjustmentModal.value = false
      // 詳細情報を再取得
      await productStore.fetchProduct(productStore.currentProduct.id)
    } catch (error) {
      console.error('在庫調整エラー:', error)
    }
  }
}

const handleDelete = async () => {
  if (productStore.currentProduct?.id) {
    try {
      await productStore.deleteProduct(productStore.currentProduct.id)
      showDeleteModal.value = false
      // 一覧画面に戻る
      router.push('/products')
    } catch (error) {
      console.error('削除エラー:', error)
    }
  }
}

// ライフサイクル
onMounted(async () => {
  const productId = Number(route.params.id)
  if (productId) {
    try {
      await productStore.fetchProduct(productId)
    } catch (error) {
      console.error('製品取得エラー:', error)
    }
  }
})
</script>

<style scoped>
.product-detail {
  padding: 20px;
  max-width: 1000px;
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

.product-info {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.info-section {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  overflow: hidden;
}

.info-section h2 {
  margin: 0;
  padding: 20px;
  background: #f8f9fa;
  color: #333;
  font-size: 20px;
  border-bottom: 1px solid #dee2e6;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  padding: 30px;
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

.price {
  font-weight: bold;
  color: #28a745;
}

.description {
  line-height: 1.6;
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

.stock-status {
  font-size: 12px;
  opacity: 0.8;
}

.stock-overview {
  padding: 30px;
  display: flex;
  align-items: center;
  gap: 30px;
}

.stock-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  border-radius: 8px;
  background: #f8f9fa;
  border: 2px solid;
  min-width: 120px;
}

.stock-card.stock-normal {
  border-color: #28a745;
  background: #d4edda;
}

.stock-card.stock-low {
  border-color: #ffc107;
  background: #fff3cd;
}

.stock-card.stock-empty {
  border-color: #dc3545;
  background: #f8d7da;
}

.stock-number {
  font-size: 36px;
  font-weight: bold;
  margin-bottom: 5px;
}

.stock-label {
  font-size: 14px;
  opacity: 0.8;
}

.stock-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
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

.btn-success {
  background: #28a745;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-success:hover {
  background: #218838;
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

.btn-sm {
  padding: 6px 12px;
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

.btn-sm:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* レスポンシブ対応 */
@media (max-width: 768px) {
  .product-detail {
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
    padding: 20px;
  }

  .stock-overview {
    flex-direction: column;
    gap: 20px;
  }

  .stock-actions {
    justify-content: center;
  }
}
</style>