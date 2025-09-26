import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { setActivePinia, createPinia } from 'pinia'
import ProductDetailView from '../ProductDetailView.vue'
import { useProductStore } from '@/stores/product'

// ルートパラメータのモック
vi.mock('vue-router', async () => {
  const actual = await vi.importActual('vue-router')
  return {
    ...actual,
    useRoute: () => ({
      params: { id: '1' }
    }),
    useRouter: () => ({
      push: vi.fn()
    })
  }
})

// onMountedでfetchProductが呼ばれないようにする
const originalOnMounted = vi.hoisted(() => vi.fn())
vi.mock('vue', async () => {
  const actual = await vi.importActual('vue')
  return {
    ...actual,
    onMounted: originalOnMounted
  }
})

// モックルーター
const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: { template: 'Home' } },
    { path: '/products', component: { template: 'Product List' } },
    { path: '/products/:id', component: ProductDetailView }
  ]
})

// Pinia セットアップ
beforeEach(() => {
  setActivePinia(createPinia())
})

afterEach(() => {
  vi.clearAllMocks()
})

describe('ProductDetailView', () => {
  const mockProduct = {
    id: 1,
    name: 'テスト製品',
    description: 'テスト用の製品です',
    sku: 'TEST-001',
    price: 1000,
    stock_quantity: 50,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z'
  }

  describe('コンポーネントの初期表示', () => {
    it('ページタイトルが正しく表示される', async () => {
      const productStore = useProductStore()
      productStore.setCurrentProduct(mockProduct)

      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      expect(wrapper.find('h1').text()).toBe('製品詳細')
    })

    it('戻るボタンが表示される', async () => {
      const productStore = useProductStore()
      productStore.setCurrentProduct(mockProduct)

      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      const backButton = wrapper.find('.btn-back')
      expect(backButton.exists()).toBe(true)
      expect(backButton.text()).toContain('戻る')
    })

    it('操作ボタンが表示される', async () => {
      const productStore = useProductStore()
      productStore.setCurrentProduct(mockProduct)

      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      expect(wrapper.find('.btn-primary').text()).toBe('編集')
      expect(wrapper.find('.btn-success').text()).toBe('在庫調整')
      expect(wrapper.find('.btn-danger').text()).toBe('削除')
    })
  })

  describe('製品情報の表示', () => {
    it('基本情報が正しく表示される', async () => {
      const productStore = useProductStore()
      productStore.setCurrentProduct(mockProduct)
      productStore.setLoading(false)
      productStore.clearError()
      productStore.clearError() // ローディングを解除

      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('1') // ID
      expect(wrapper.text()).toContain('テスト製品') // 製品名
      expect(wrapper.text()).toContain('TEST-001') // SKU
      expect(wrapper.text()).toContain('¥1,000') // 価格
      expect(wrapper.text()).toContain('50') // 在庫数量
      expect(wrapper.text()).toContain('テスト用の製品です') // 説明
    })

    it('在庫評価額が正しく計算される', async () => {
      const productStore = useProductStore()
      productStore.setCurrentProduct(mockProduct)
      productStore.setLoading(false)
      productStore.clearError()
      productStore.clearError()

      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('¥50,000') // 1000 * 50
    })

    it('在庫ステータスが正しく表示される', async () => {
      const productStore = useProductStore()
      productStore.setCurrentProduct(mockProduct)
      productStore.setLoading(false)
      productStore.clearError()
      productStore.clearError()

      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('在庫十分')
    })

    it('在庫が少ない場合の表示', async () => {
      const lowStockProduct = { ...mockProduct, stock_quantity: 5 }
      const productStore = useProductStore()
      productStore.setCurrentProduct(lowStockProduct)
      productStore.setLoading(false)
      productStore.clearError()

      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('在庫少')
    })

    it('在庫切れの場合の表示', async () => {
      const emptyStockProduct = { ...mockProduct, stock_quantity: 0 }
      const productStore = useProductStore()
      productStore.setCurrentProduct(emptyStockProduct)
      productStore.setLoading(false)
      productStore.clearError()

      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('在庫切れ')
    })
  })

  describe('在庫状況セクション', () => {
    it('在庫カードが表示される', async () => {
      const productStore = useProductStore()
      productStore.setCurrentProduct(mockProduct)
      productStore.setLoading(false)
      productStore.clearError()
      productStore.clearError()

      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      await wrapper.vm.$nextTick()

      expect(wrapper.find('.stock-card').exists()).toBe(true)
      expect(wrapper.text()).toContain('現在在庫')
    })

    it('クイック調整ボタンが表示される', async () => {
      const productStore = useProductStore()
      productStore.setCurrentProduct(mockProduct)
      productStore.setLoading(false)
      productStore.clearError()
      productStore.clearError()

      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      await wrapper.vm.$nextTick()

      const actionButtons = wrapper.findAll('.stock-actions .btn-sm')
      expect(actionButtons.length).toBeGreaterThanOrEqual(5)
      expect(wrapper.text()).toContain('+1')
      expect(wrapper.text()).toContain('+10')
      expect(wrapper.text()).toContain('-1')
      expect(wrapper.text()).toContain('-10')
      expect(wrapper.text()).toContain('詳細調整')
    })

    it('在庫が不足している場合は減算ボタンが無効になる', async () => {
      const lowStockProduct = { ...mockProduct, stock_quantity: 5 }
      const productStore = useProductStore()
      productStore.setCurrentProduct(lowStockProduct)
      productStore.setLoading(false)
      productStore.clearError()

      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      await wrapper.vm.$nextTick()

      const minus10Button = wrapper.findAll('.btn-sm').find(btn =>
        btn.text().includes('-10')
      )
      expect(minus10Button?.attributes('disabled')).toBeDefined()
    })
  })

  describe('ローディング状態', () => {
    it('ローディング中にメッセージが表示される', async () => {
      const productStore = useProductStore()
      productStore.setLoading(true)

      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      expect(wrapper.text()).toContain('読み込み中...')
    })
  })

  describe('エラー状態', () => {
    it('エラー時にメッセージが表示される', async () => {
      const productStore = useProductStore()
      productStore.setError('テストエラー')

      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      expect(wrapper.text()).toContain('テストエラー')
    })
  })

  describe('製品が見つからない場合', () => {
    it('製品が存在しない場合にメッセージが表示される', async () => {
      const productStore = useProductStore()
      productStore.setCurrentProduct(null)
      productStore.setLoading(false)
      productStore.clearError()

      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      expect(wrapper.text()).toContain('製品が見つかりません')
      expect(wrapper.text()).toContain('指定された製品が存在しないか、削除された可能性があります')
    })
  })

  describe('モーダル表示制御', () => {
    it('編集ボタンクリック時に編集モーダルが表示される', async () => {
      const productStore = useProductStore()
      productStore.setCurrentProduct(mockProduct)
      productStore.setLoading(false)
      productStore.clearError()
      productStore.clearError()

      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      const editButton = wrapper.find('.btn-primary')
      await editButton.trigger('click')

      expect(wrapper.vm.showEditModal).toBe(true)
    })

    it('在庫調整ボタンクリック時に在庫調整モーダルが表示される', async () => {
      const productStore = useProductStore()
      productStore.setCurrentProduct(mockProduct)
      productStore.setLoading(false)
      productStore.clearError()
      productStore.clearError()

      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      const stockButton = wrapper.find('.btn-success')
      await stockButton.trigger('click')

      expect(wrapper.vm.showStockAdjustmentModal).toBe(true)
    })

    it('削除ボタンクリック時に削除確認モーダルが表示される', async () => {
      const productStore = useProductStore()
      productStore.setCurrentProduct(mockProduct)
      productStore.setLoading(false)
      productStore.clearError()
      productStore.clearError()

      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      const deleteButton = wrapper.find('.btn-danger')
      await deleteButton.trigger('click')

      expect(wrapper.vm.showDeleteModal).toBe(true)
    })
  })

  describe('ユーティリティ関数', () => {
    it('価格のフォーマットが正しい', async () => {
      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      const componentInstance = wrapper.vm as any
      expect(componentInstance.formatPrice(1000)).toBe('1,000')
      expect(componentInstance.formatPrice(1234567)).toBe('1,234,567')
    })

    it('日付のフォーマットが正しい', async () => {
      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      const dateString = '2024-01-01T00:00:00Z'
      const componentInstance = wrapper.vm as any
      const formatted = componentInstance.formatDate(dateString)
      expect(formatted).toContain('2024')
    })

    it('未定義の日付の場合は"-"が返される', async () => {
      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      const componentInstance = wrapper.vm as any
      expect(componentInstance.formatDate(undefined)).toBe('-')
    })

    it('在庫ステータスの判定が正しい', async () => {
      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      const componentInstance = wrapper.vm as any
      expect(componentInstance.getStockStatus(0)).toBe('在庫切れ')
      expect(componentInstance.getStockStatus(5)).toBe('在庫少')
      expect(componentInstance.getStockStatus(50)).toBe('在庫十分')
    })

    it('在庫クラスの判定が正しい', async () => {
      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      const componentInstance = wrapper.vm as any
      expect(componentInstance.getStockClass(0)).toBe('stock-empty')
      expect(componentInstance.getStockClass(5)).toBe('stock-low')
      expect(componentInstance.getStockClass(50)).toBe('stock-normal')
    })
  })

  describe('レスポンシブ対応', () => {
    it('レスポンシブ用のクラスが適用されている', async () => {
      const productStore = useProductStore()
      productStore.setCurrentProduct(mockProduct)

      const wrapper = mount(ProductDetailView, {
        global: {
          plugins: [router]
        }
      })

      await wrapper.vm.$nextTick()

      expect(wrapper.find('.product-detail').exists()).toBe(true)
      expect(wrapper.find('.info-section').exists()).toBe(true)
      expect(wrapper.find('.info-grid').exists()).toBe(true)
    })
  })
})