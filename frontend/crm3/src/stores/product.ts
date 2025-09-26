import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export interface Product {
  id?: number
  name: string
  description?: string
  sku: string
  price: number
  stock_quantity: number
  created_at?: string
  updated_at?: string
}

export interface ProductSearchParams {
  name?: string
  description?: string
  sku?: string
  min_price?: number
  max_price?: number
  page?: number
  limit?: number
}

export interface PaginatedProductResponse {
  products: Product[]
  total: number
  page: number
  limit: number
  total_pages: number
}

export const useProductStore = defineStore('product', () => {
  const products = ref<Product[]>([])
  const currentProduct = ref<Product | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const pagination = ref({
    total: 0,
    page: 1,
    limit: 10,
    total_pages: 0
  })

  // Computed
  const isLoading = computed(() => loading.value)
  const hasError = computed(() => error.value !== null)
  const isEmpty = computed(() => products.value.length === 0)

  // Actions
  const clearError = () => {
    error.value = null
  }

  const setLoading = (value: boolean) => {
    loading.value = value
  }

  const setError = (message: string) => {
    error.value = message
    loading.value = false
  }

  const setProductsFromResponse = (data: PaginatedProductResponse) => {
    products.value = data.products
    pagination.value = {
      total: data.total,
      page: data.page,
      limit: data.limit,
      total_pages: data.total_pages
    }
  }

  const setCurrentProduct = (product: Product | null) => {
    currentProduct.value = product
  }

  const addProduct = (product: Product) => {
    products.value.unshift(product)
    pagination.value.total += 1
  }

  const updateProduct = (updatedProduct: Product) => {
    const index = products.value.findIndex(p => p.id === updatedProduct.id)
    if (index !== -1) {
      products.value[index] = updatedProduct
    }
    if (currentProduct.value?.id === updatedProduct.id) {
      currentProduct.value = updatedProduct
    }
  }

  const removeProduct = (productId: number) => {
    products.value = products.value.filter(p => p.id !== productId)
    pagination.value.total -= 1
    if (currentProduct.value?.id === productId) {
      currentProduct.value = null
    }
  }

  // API Methods (モックとして実装、後でAPI呼び出しに置き換え)
  const fetchProducts = async (params: ProductSearchParams = {}) => {
    try {
      setLoading(true)
      clearError()

      // TODO: 実際のAPI呼び出しに置き換える
      await new Promise(resolve => setTimeout(resolve, 500)) // モック遅延

      const mockData: PaginatedProductResponse = {
        products: [
          {
            id: 1,
            name: 'サンプル製品A',
            description: '高品質な製品A',
            sku: 'SKU-001',
            price: 1500,
            stock_quantity: 100,
            created_at: '2024-01-01'
          },
          {
            id: 2,
            name: 'サンプル製品B',
            description: '人気の製品B',
            sku: 'SKU-002',
            price: 2800,
            stock_quantity: 50,
            created_at: '2024-01-02'
          }
        ],
        total: 2,
        page: params.page || 1,
        limit: params.limit || 10,
        total_pages: 1
      }

      setProductsFromResponse(mockData)
    } catch (err) {
      setError(err instanceof Error ? err.message : '製品一覧の取得に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  const fetchProduct = async (id: number) => {
    try {
      setLoading(true)
      clearError()

      // TODO: 実際のAPI呼び出しに置き換える
      await new Promise(resolve => setTimeout(resolve, 300))

      const mockProduct: Product = {
        id: id,
        name: 'サンプル製品A',
        description: '高品質な製品Aです。長期間の耐久性と優れた性能を提供します。',
        sku: 'SKU-001',
        price: 1500,
        stock_quantity: 100,
        created_at: '2024-01-01'
      }

      setCurrentProduct(mockProduct)
      return mockProduct
    } catch (err) {
      setError(err instanceof Error ? err.message : '製品情報の取得に失敗しました')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const createProduct = async (productData: Omit<Product, 'id' | 'created_at' | 'updated_at'>) => {
    try {
      setLoading(true)
      clearError()

      // バリデーション
      if (!productData.name.trim()) {
        throw new Error('製品名は必須です')
      }
      if (!productData.sku.trim()) {
        throw new Error('SKUは必須です')
      }
      if (productData.price <= 0) {
        throw new Error('価格は正数である必要があります')
      }
      if (productData.stock_quantity < 0) {
        throw new Error('在庫数量は0以上である必要があります')
      }

      // TODO: 実際のAPI呼び出しに置き換える
      await new Promise(resolve => setTimeout(resolve, 500))

      const newProduct: Product = {
        ...productData,
        id: Date.now(), // モックID
        created_at: new Date().toISOString()
      }

      addProduct(newProduct)
      return newProduct
    } catch (err) {
      setError(err instanceof Error ? err.message : '製品の作成に失敗しました')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const updateProductData = async (id: number, productData: Omit<Product, 'id' | 'created_at' | 'updated_at'>) => {
    try {
      setLoading(true)
      clearError()

      // バリデーション
      if (!productData.name.trim()) {
        throw new Error('製品名は必須です')
      }
      if (!productData.sku.trim()) {
        throw new Error('SKUは必須です')
      }
      if (productData.price <= 0) {
        throw new Error('価格は正数である必要があります')
      }
      if (productData.stock_quantity < 0) {
        throw new Error('在庫数量は0以上である必要があります')
      }

      // TODO: 実際のAPI呼び出しに置き換える
      await new Promise(resolve => setTimeout(resolve, 500))

      const updatedProduct: Product = {
        ...productData,
        id: id,
        updated_at: new Date().toISOString()
      }

      updateProduct(updatedProduct)
      return updatedProduct
    } catch (err) {
      setError(err instanceof Error ? err.message : '製品情報の更新に失敗しました')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const deleteProduct = async (id: number) => {
    try {
      setLoading(true)
      clearError()

      // TODO: 実際のAPI呼び出しに置き換える
      await new Promise(resolve => setTimeout(resolve, 300))

      removeProduct(id)
    } catch (err) {
      setError(err instanceof Error ? err.message : '製品の削除に失敗しました')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const adjustStock = async (id: number, quantity: number, operation: 'add' | 'subtract') => {
    try {
      setLoading(true)
      clearError()

      if (quantity <= 0) {
        throw new Error('数量は正数である必要があります')
      }

      // TODO: 実際のAPI呼び出しに置き換える
      await new Promise(resolve => setTimeout(resolve, 300))

      const product = products.value.find(p => p.id === id)
      if (product) {
        const newQuantity = operation === 'add'
          ? product.stock_quantity + quantity
          : Math.max(0, product.stock_quantity - quantity)

        const updatedProduct: Product = {
          ...product,
          stock_quantity: newQuantity,
          updated_at: new Date().toISOString()
        }

        updateProduct(updatedProduct)
      }

      if (currentProduct.value?.id === id) {
        await fetchProduct(id) // 詳細画面のデータも更新
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '在庫調整に失敗しました')
      throw err
    } finally {
      setLoading(false)
    }
  }

  // Helper methods for testing and direct state manipulation
  const setProducts = (productList: Product[], page: number, pageSize: number, total: number) => {
    products.value = productList
    pagination.value = {
      total,
      page,
      limit: pageSize,
      total_pages: Math.ceil(total / pageSize)
    }
  }

  return {
    // State
    products,
    currentProduct,
    loading,
    error,
    pagination,

    // Computed properties (exported as getters for easier access in tests)
    isLoading,
    hasError,
    isEmpty,
    totalCount: computed(() => pagination.value.total),
    currentPage: computed(() => pagination.value.page),
    pageSize: computed(() => pagination.value.limit),
    totalPages: computed(() => pagination.value.total_pages),

    // Actions
    clearError,
    setLoading,
    setError,
    setProducts,
    setCurrentProduct,
    addProduct,
    updateProduct,
    removeProduct,
    fetchProducts,
    fetchProduct,
    createProduct,
    updateProductData,
    deleteProduct,
    adjustStock
  }
})