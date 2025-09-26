import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useProductStore, type Product } from '../product'

// Pinia セットアップ
beforeEach(() => {
  setActivePinia(createPinia())
})

afterEach(() => {
  vi.clearAllMocks()
})

describe('Product Store', () => {
  describe('初期状態', () => {
    it('初期状態が正しく設定されている', () => {
      const productStore = useProductStore()

      expect(productStore.products).toEqual([])
      expect(productStore.currentProduct).toBeNull()
      expect(productStore.isLoading).toBe(false)
      expect(productStore.error).toBeNull()
      expect(productStore.hasError).toBe(false)
      expect(productStore.totalCount).toBe(0)
      expect(productStore.currentPage).toBe(1)
      expect(productStore.pageSize).toBe(10)
      expect(productStore.totalPages).toBe(0)
    })
  })

  describe('データ操作', () => {
    it('製品データを正しく設定できる', () => {
      const productStore = useProductStore()
      const mockProducts: Product[] = [
        {
          id: 1,
          name: 'テスト製品',
          description: 'テスト用の製品です',
          sku: 'TEST-001',
          price: 1000,
          stock_quantity: 50,
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z'
        }
      ]

      productStore.setProducts(mockProducts, 1, 10, 100)

      expect(productStore.products).toEqual(mockProducts)
      expect(productStore.totalCount).toBe(100)
      expect(productStore.currentPage).toBe(1)
      expect(productStore.pageSize).toBe(10)
      expect(productStore.totalPages).toBe(10)
    })

    it('現在の製品を正しく設定できる', () => {
      const productStore = useProductStore()
      const mockProduct: Product = {
        id: 1,
        name: 'テスト製品',
        description: 'テスト用の製品です',
        sku: 'TEST-001',
        price: 1000,
        stock_quantity: 50,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z'
      }

      productStore.setCurrentProduct(mockProduct)
      expect(productStore.currentProduct).toEqual(mockProduct)
    })

    it('製品を追加できる', () => {
      const productStore = useProductStore()
      const existingProduct: Product = {
        id: 1,
        name: '既存製品',
        description: '既存の製品です',
        sku: 'EXIST-001',
        price: 500,
        stock_quantity: 30,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z'
      }

      const newProduct: Product = {
        id: 2,
        name: '新規製品',
        description: '新規の製品です',
        sku: 'NEW-001',
        price: 1000,
        stock_quantity: 50,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z'
      }

      productStore.setProducts([existingProduct], 1, 10, 1)
      productStore.addProduct(newProduct)

      expect(productStore.products).toHaveLength(2)
      expect(productStore.products[0]).toEqual(newProduct) // unshiftされるので最初に来る
      expect(productStore.totalCount).toBe(2)
    })

    it('製品を更新できる', () => {
      const productStore = useProductStore()
      const originalProduct: Product = {
        id: 1,
        name: 'オリジナル製品',
        description: 'オリジナルの製品です',
        sku: 'ORIG-001',
        price: 500,
        stock_quantity: 30,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z'
      }

      const updatedProduct: Product = {
        ...originalProduct,
        name: '更新された製品',
        price: 800,
        updated_at: '2024-01-02T00:00:00Z'
      }

      productStore.setProducts([originalProduct], 1, 10, 1)
      productStore.updateProduct(updatedProduct)

      expect(productStore.products[0]).toEqual(updatedProduct)
      expect(productStore.products[0].name).toBe('更新された製品')
      expect(productStore.products[0].price).toBe(800)
    })

    it('製品を削除できる', () => {
      const productStore = useProductStore()
      const product1: Product = {
        id: 1,
        name: '製品1',
        description: '製品1です',
        sku: 'PROD-001',
        price: 500,
        stock_quantity: 30,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z'
      }

      const product2: Product = {
        id: 2,
        name: '製品2',
        description: '製品2です',
        sku: 'PROD-002',
        price: 1000,
        stock_quantity: 50,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z'
      }

      productStore.setProducts([product1, product2], 1, 10, 2)
      productStore.removeProduct(1)

      expect(productStore.products).toHaveLength(1)
      expect(productStore.products[0]).toEqual(product2)
      expect(productStore.totalCount).toBe(1)
    })
  })

  describe('エラーハンドリング', () => {
    it('エラーを正しく設定できる', () => {
      const productStore = useProductStore()
      const errorMessage = 'テストエラー'

      productStore.setError(errorMessage)

      expect(productStore.error).toBe(errorMessage)
      expect(productStore.hasError).toBe(true)
    })

    it('エラーをクリアできる', () => {
      const productStore = useProductStore()

      productStore.setError('テストエラー')
      expect(productStore.hasError).toBe(true)

      productStore.clearError()
      expect(productStore.error).toBeNull()
      expect(productStore.hasError).toBe(false)
    })
  })

  describe('ローディング状態', () => {
    it('ローディング状態を正しく管理できる', () => {
      const productStore = useProductStore()

      expect(productStore.isLoading).toBe(false)

      productStore.setLoading(true)
      expect(productStore.isLoading).toBe(true)

      productStore.setLoading(false)
      expect(productStore.isLoading).toBe(false)
    })
  })

  describe('在庫調整', () => {
    it('在庫を正しく調整できる（加算）', async () => {
      const productStore = useProductStore()
      const mockProduct: Product = {
        id: 1,
        name: 'テスト製品',
        description: 'テスト用の製品です',
        sku: 'TEST-001',
        price: 1000,
        stock_quantity: 50,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z'
      }

      productStore.setProducts([mockProduct], 1, 10, 1)

      await productStore.adjustStock(1, 10, 'add')

      expect(productStore.products[0].stock_quantity).toBe(60)
    })

    it('在庫を正しく調整できる（減算）', async () => {
      const productStore = useProductStore()
      const mockProduct: Product = {
        id: 1,
        name: 'テスト製品',
        description: 'テスト用の製品です',
        sku: 'TEST-001',
        price: 1000,
        stock_quantity: 50,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z'
      }

      productStore.setProducts([mockProduct], 1, 10, 1)

      await productStore.adjustStock(1, 20, 'subtract')

      expect(productStore.products[0].stock_quantity).toBe(30)
    })

    it('在庫が0未満にならないように制限される', async () => {
      const productStore = useProductStore()
      const mockProduct: Product = {
        id: 1,
        name: 'テスト製品',
        description: 'テスト用の製品です',
        sku: 'TEST-001',
        price: 1000,
        stock_quantity: 10,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z'
      }

      productStore.setProducts([mockProduct], 1, 10, 1)

      await productStore.adjustStock(1, 20, 'subtract')

      expect(productStore.products[0].stock_quantity).toBe(0)
    })

    it('無効な数量でエラーが発生する', async () => {
      const productStore = useProductStore()
      const mockProduct: Product = {
        id: 1,
        name: 'テスト製品',
        description: 'テスト用の製品です',
        sku: 'TEST-001',
        price: 1000,
        stock_quantity: 50,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z'
      }

      productStore.setProducts([mockProduct], 1, 10, 1)

      await expect(productStore.adjustStock(1, 0, 'add')).rejects.toThrow('数量は正数である必要があります')
      await expect(productStore.adjustStock(1, -5, 'add')).rejects.toThrow('数量は正数である必要があります')
    })
  })

  describe('製品検索とフィルタリング', () => {
    it('fetchProducts が正しく呼び出される', async () => {
      const productStore = useProductStore()

      await productStore.fetchProducts(1, 10, 'テスト', 100, 2000)

      expect(productStore.isLoading).toBe(false)
    })

    it('fetchProduct が正しく呼び出される', async () => {
      const productStore = useProductStore()

      await productStore.fetchProduct(1)

      expect(productStore.isLoading).toBe(false)
    })

    it('createProduct が正しく呼び出される', async () => {
      const productStore = useProductStore()
      const productData = {
        name: '新規製品',
        description: '新規の製品です',
        sku: 'NEW-001',
        price: 1000,
        stock_quantity: 50
      }

      await productStore.createProduct(productData)

      expect(productStore.isLoading).toBe(false)
    })

    it('updateProductData が正しく呼び出される', async () => {
      const productStore = useProductStore()
      const productData = {
        name: '更新製品',
        description: '更新された製品です',
        sku: 'UPD-001',
        price: 1500,
        stock_quantity: 75
      }

      await productStore.updateProductData(1, productData)

      expect(productStore.isLoading).toBe(false)
    })

    it('deleteProduct が正しく呼び出される', async () => {
      const productStore = useProductStore()

      await productStore.deleteProduct(1)

      expect(productStore.isLoading).toBe(false)
    })
  })
})