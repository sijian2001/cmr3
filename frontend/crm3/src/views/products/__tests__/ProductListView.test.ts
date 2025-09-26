import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { setActivePinia, createPinia } from 'pinia'
import ProductListView from '../ProductListView.vue'
import { useProductStore } from '@/stores/product'

// onMountedでfetchProductsが呼ばれないようにする
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
    { path: '/products', component: ProductListView },
    { path: '/products/:id', component: { template: 'Product Detail' } }
  ]
})

// Pinia セットアップ
beforeEach(() => {
  setActivePinia(createPinia())
})

afterEach(() => {
  vi.clearAllMocks()
})

describe('ProductListView', () => {
  describe('コンポーネントの初期表示', () => {
    it('ページタイトルが正しく表示される', async () => {
      const wrapper = mount(ProductListView, {
        global: {
          plugins: [router]
        }
      })

      expect(wrapper.find('h1').text()).toBe('製品管理')
    })

    it('新規登録ボタンが表示される', async () => {
      const wrapper = mount(ProductListView, {
        global: {
          plugins: [router]
        }
      })

      const newButton = wrapper.find('.btn-primary')
      expect(newButton.exists()).toBe(true)
      expect(newButton.text()).toBe('新規製品登録')
    })

    it('検索フォームが表示される', async () => {
      const wrapper = mount(ProductListView, {
        global: {
          plugins: [router]
        }
      })

      // 実際のinputはラベルで確認（ID属性はない）
      expect(wrapper.text()).toContain('製品名')
      expect(wrapper.text()).toContain('最小価格')
      expect(wrapper.text()).toContain('最大価格')
      expect(wrapper.find('.btn-secondary').exists()).toBe(true) // 検索ボタン
      expect(wrapper.find('.btn-outline').exists()).toBe(true) // クリアボタン
    })

    it('製品テーブルが表示される', async () => {
      const wrapper = mount(ProductListView, {
        global: {
          plugins: [router]
        }
      })

      const table = wrapper.find('table')
      expect(table.exists()).toBe(true)

      const headers = wrapper.findAll('th')
      expect(headers).toHaveLength(7) // ID, 製品名, SKU, 価格, 在庫数量, 登録日, 操作
      expect(headers[0].text()).toBe('ID')
      expect(headers[1].text()).toBe('製品名')
      expect(headers[2].text()).toBe('SKU')
      expect(headers[3].text()).toBe('価格')
      expect(headers[4].text()).toBe('在庫数量')
      expect(headers[5].text()).toBe('登録日')
      expect(headers[6].text()).toBe('操作')
    })
  })

  describe('検索機能', () => {
    it('検索ボタンクリック時に検索が実行される', async () => {
      const wrapper = mount(ProductListView, {
        global: {
          plugins: [router]
        }
      })

      const productStore = useProductStore()
      const fetchSpy = vi.spyOn(productStore, 'fetchProducts')

      const inputs = wrapper.findAll('input')
      const searchInput = inputs[0] // 製品名のinput（最初のinput）
      const searchButton = wrapper.find('.btn-secondary') // 検索ボタン

      await searchInput.setValue('テスト')
      await searchButton.trigger('click')

      expect(fetchSpy).toHaveBeenCalledWith({
        name: 'テスト',
        description: '',
        sku: '',
        min_price: undefined,
        max_price: undefined,
        page: 1,
        limit: 10
      })
    })

    it('価格範囲指定で検索が実行される', async () => {
      const wrapper = mount(ProductListView, {
        global: {
          plugins: [router]
        }
      })

      const productStore = useProductStore()
      const fetchSpy = vi.spyOn(productStore, 'fetchProducts')

      const inputs = wrapper.findAll('input')
      const minPriceInput = inputs[3] // 最小価格のinput（4番目のinput）
      const maxPriceInput = inputs[4] // 最大価格のinput（5番目のinput）
      const searchButton = wrapper.find('.btn-secondary') // 検索ボタン

      await minPriceInput.setValue('100')
      await maxPriceInput.setValue('1000')
      await searchButton.trigger('click')

      expect(fetchSpy).toHaveBeenCalledWith({
        name: '',
        description: '',
        sku: '',
        min_price: 100,
        max_price: 1000,
        page: 1,
        limit: 10
      })
    })

    it('クリアボタンで検索条件がリセットされる', async () => {
      const wrapper = mount(ProductListView, {
        global: {
          plugins: [router]
        }
      })

      const inputs = wrapper.findAll('input')
      const searchInput = inputs[0] // 製品名のinput
      const minPriceInput = inputs[3] // 最小価格のinput
      const maxPriceInput = inputs[4] // 最大価格のinput
      const clearButton = wrapper.find('.btn-outline') // クリアボタン

      await searchInput.setValue('テスト')
      await minPriceInput.setValue('100')
      await maxPriceInput.setValue('1000')

      await clearButton.trigger('click')

      expect((searchInput.element as HTMLInputElement).value).toBe('')
      expect((minPriceInput.element as HTMLInputElement).value).toBe('')
      expect((maxPriceInput.element as HTMLInputElement).value).toBe('')
    })
  })

  describe('製品一覧表示', () => {
    it('製品データが正しく表示される', async () => {
      const productStore = useProductStore()
      productStore.setProducts([
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
      ], 1, 10, 1)
      productStore.setLoading(false)
      productStore.clearError() // ローディング状態を解除

      const wrapper = mount(ProductListView, {
        global: {
          plugins: [router]
        }
      })

      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('テスト製品')
      expect(wrapper.text()).toContain('TEST-001')
      expect(wrapper.text()).toContain('¥1,000')
      expect(wrapper.text()).toContain('50')
    })

    it('在庫数量が正しく表示される', async () => {
      const productStore = useProductStore()
      productStore.setProducts([
        {
          id: 1,
          name: '在庫切れ製品',
          description: '',
          sku: 'EMPTY-001',
          price: 1000,
          stock_quantity: 0,
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z'
        },
        {
          id: 2,
          name: '在庫少製品',
          description: '',
          sku: 'LOW-001',
          price: 1000,
          stock_quantity: 5,
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z'
        },
        {
          id: 3,
          name: '在庫十分製品',
          description: '',
          sku: 'NORMAL-001',
          price: 1000,
          stock_quantity: 100,
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z'
        }
      ], 1, 10, 3)
      productStore.setLoading(false)
      productStore.clearError()

      const wrapper = mount(ProductListView, {
        global: {
          plugins: [router]
        }
      })

      await wrapper.vm.$nextTick()

      // 在庫数量が表示されることを確認
      expect(wrapper.text()).toContain('0') // 在庫切れ製品の在庫
      expect(wrapper.text()).toContain('5') // 在庫少製品の在庫
      expect(wrapper.text()).toContain('100') // 在庫十分製品の在庫
    })

    it('製品データがない場合にメッセージが表示される', async () => {
      const productStore = useProductStore()
      productStore.setProducts([], 1, 10, 0)
      productStore.setLoading(false)
      productStore.clearError()

      const wrapper = mount(ProductListView, {
        global: {
          plugins: [router]
        }
      })

      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('製品データがありません')
    })
  })

  describe('ローディング状態', () => {
    it('ローディング中にメッセージが表示される', async () => {
      const productStore = useProductStore()
      productStore.setLoading(true)

      const wrapper = mount(ProductListView, {
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

      const wrapper = mount(ProductListView, {
        global: {
          plugins: [router]
        }
      })

      expect(wrapper.text()).toContain('テストエラー')
    })
  })

  describe('ページネーション', () => {
    it('複数ページの場合にページネーションが表示される', async () => {
      const productStore = useProductStore()
      // 複数ページを作るために最低1つのアイテムが必要
      productStore.setProducts([{
        id: 1,
        name: 'ダミー製品',
        description: '',
        sku: 'DUMMY-001',
        price: 1000,
        stock_quantity: 1,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z'
      }], 1, 10, 100) // 10ページ
      productStore.setLoading(false)
      productStore.clearError()

      const wrapper = mount(ProductListView, {
        global: {
          plugins: [router]
        }
      })

      await wrapper.vm.$nextTick()

      expect(wrapper.find('.pagination').exists()).toBe(true)
      expect(wrapper.text()).toContain('1 / 10')
    })

    it('単一ページの場合はページネーションが非表示', async () => {
      const productStore = useProductStore()
      productStore.setProducts([], 1, 10, 5) // 1ページのみ
      productStore.setLoading(false)
      productStore.clearError()

      const wrapper = mount(ProductListView, {
        global: {
          plugins: [router]
        }
      })

      await wrapper.vm.$nextTick()

      expect(wrapper.find('.pagination').exists()).toBe(false)
    })
  })

  describe('モーダル操作', () => {
    it('新規登録ボタンクリック時にモーダルが開く', async () => {
      const wrapper = mount(ProductListView, {
        global: {
          plugins: [router]
        }
      })

      const newButton = wrapper.find('.btn-primary')
      await newButton.trigger('click')

      const componentInstance = wrapper.vm as any
      expect(componentInstance.showCreateModal).toBe(true)
    })
  })

  describe('ユーティリティ関数', () => {
    it('価格のフォーマットが正しい', async () => {
      const wrapper = mount(ProductListView, {
        global: {
          plugins: [router]
        }
      })

      const componentInstance = wrapper.vm as any
      expect(componentInstance.formatPrice(1000)).toBe('1,000')
      expect(componentInstance.formatPrice(1234567)).toBe('1,234,567')
    })

    it('日付のフォーマットが正しい', async () => {
      const wrapper = mount(ProductListView, {
        global: {
          plugins: [router]
        }
      })

      const dateString = '2024-01-01T00:00:00Z'
      const componentInstance = wrapper.vm as any
      const formatted = componentInstance.formatDate(dateString)
      expect(formatted).toContain('2024')
    })

    it('在庫クラスの判定が正しい（getStockClass）', async () => {
      const wrapper = mount(ProductListView, {
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
      const wrapper = mount(ProductListView, {
        global: {
          plugins: [router]
        }
      })

      expect(wrapper.find('.table-container').exists()).toBe(true)
      expect(wrapper.find('.search-form').exists()).toBe(true)
    })
  })
})