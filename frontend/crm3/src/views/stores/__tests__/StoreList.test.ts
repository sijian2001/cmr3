import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import StoreList from '../StoreList.vue'
import { useStoreStore, type Store } from '@/stores/store'
import StoreForm from '@/components/StoreForm.vue'

// Mock debounce function
vi.mock('lodash-es', () => ({
  debounce: (fn: Function) => fn
}))

const mockStores: Store[] = [
  {
    id: 1,
    name: 'テスト店舗1',
    address: '東京都渋谷区',
    phone: '03-1234-5678',
    email: 'test1@example.com',
    status: 'active',
    manager_id: 1,
    created_at: '2023-01-01T00:00:00Z',
    updated_at: '2023-01-01T00:00:00Z'
  },
  {
    id: 2,
    name: 'テスト店舗2',
    address: '大阪府大阪市',
    phone: '06-1234-5678',
    email: 'test2@example.com',
    status: 'inactive',
    manager_id: 2,
    created_at: '2023-01-02T00:00:00Z',
    updated_at: '2023-01-02T00:00:00Z'
  },
  {
    id: 3,
    name: 'テスト店舗3',
    address: '愛知県名古屋市',
    phone: '052-1234-5678',
    email: 'test3@example.com',
    status: 'maintenance',
    manager_id: 3,
    created_at: '2023-01-03T00:00:00Z',
    updated_at: '2023-01-03T00:00:00Z'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/stores', name: 'stores', component: { template: '<div>Stores</div>' } },
    { path: '/stores/:id', name: 'store-detail', component: { template: '<div>Store Detail</div>' } }
  ]
})

describe('StoreList', () => {
  let pinia: any

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)
    vi.clearAllMocks()
  })

  const createWrapper = (options = {}) => {
    return mount(StoreList, {
      global: {
        plugins: [pinia, router],
        stubs: {
          StoreForm: true
        }
      },
      ...options
    })
  }

  describe('初期表示', () => {
    it('正しくレンダリングされる', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('h1').text()).toBe('店舗管理')
      expect(wrapper.find('.btn--primary').text()).toBe('新規店舗登録')
    })

    it('マウント時にfetchStoresが呼ばれる', () => {
      const storeStore = useStoreStore()
      const fetchStoresSpy = vi.spyOn(storeStore, 'fetchStores').mockImplementation(() => Promise.resolve())

      createWrapper()

      expect(fetchStoresSpy).toHaveBeenCalled()
    })
  })

  describe('統計カード', () => {
    it('ステータス別の統計が正しく表示される', async () => {
      const storeStore = useStoreStore()
      storeStore.stores = mockStores

      const wrapper = createWrapper()

      // 統計の確認
      const statsCards = wrapper.findAll('.stats-card')
      expect(statsCards).toHaveLength(4)

      // 営業中の店舗数
      expect(statsCards[0].find('.stats-number').text()).toBe('1')
      // 休業中の店舗数
      expect(statsCards[1].find('.stats-number').text()).toBe('1')
      // メンテナンス中の店舗数
      expect(statsCards[2].find('.stats-number').text()).toBe('1')
      // 総店舗数
      expect(statsCards[3].find('.stats-number').text()).toBe('3')
    })
  })

  describe('検索機能', () => {
    it('検索フォームが表示される', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('#search-name').exists()).toBe(true)
      expect(wrapper.find('#search-status').exists()).toBe(true)
    })

    it('店舗名で検索できる', async () => {
      const storeStore = useStoreStore()
      const searchStoresSpy = vi.spyOn(storeStore, 'searchStores').mockImplementation(() => Promise.resolve())

      const wrapper = createWrapper()
      const searchInput = wrapper.find('#search-name')

      await searchInput.setValue('テスト')
      await searchInput.trigger('input')

      expect(searchStoresSpy).toHaveBeenCalledWith({
        name: 'テスト',
        status: '',
        page: 1
      })
    })

    it('ステータスで検索できる', async () => {
      const storeStore = useStoreStore()
      const searchStoresSpy = vi.spyOn(storeStore, 'searchStores').mockImplementation(() => Promise.resolve())

      const wrapper = createWrapper()
      const statusSelect = wrapper.find('#search-status')

      await statusSelect.setValue('active')

      expect(searchStoresSpy).toHaveBeenCalledWith({
        name: '',
        status: 'active',
        page: 1
      })
    })

    it('検索クリアボタンが動作する', async () => {
      const storeStore = useStoreStore()
      const fetchStoresSpy = vi.spyOn(storeStore, 'fetchStores').mockImplementation(() => Promise.resolve())

      const wrapper = createWrapper()
      const searchInput = wrapper.find('#search-name')
      const statusSelect = wrapper.find('#search-status')
      const clearButton = wrapper.find('.btn--secondary')

      // 検索条件を設定
      await searchInput.setValue('テスト')
      await statusSelect.setValue('active')

      // クリアボタンをクリック
      await clearButton.trigger('click')

      expect(searchInput.element.value).toBe('')
      expect(statusSelect.element.value).toBe('')
      expect(fetchStoresSpy).toHaveBeenCalled()
    })
  })

  describe('店舗一覧表示', () => {
    it('店舗が正しく表示される', () => {
      const storeStore = useStoreStore()
      storeStore.stores = mockStores

      const wrapper = createWrapper()

      const storeCards = wrapper.findAll('.store-card')
      expect(storeCards).toHaveLength(3)

      // 最初の店舗の確認
      const firstCard = storeCards[0]
      expect(firstCard.find('h3').text()).toBe('テスト店舗1')
      expect(firstCard.text()).toContain('東京都渋谷区')
      expect(firstCard.text()).toContain('03-1234-5678')
      expect(firstCard.text()).toContain('test1@example.com')
      expect(firstCard.find('.status-badge--success').exists()).toBe(true)
    })

    it('店舗カードクリックで詳細画面に遷移する', async () => {
      const storeStore = useStoreStore()
      storeStore.stores = mockStores

      const pushSpy = vi.spyOn(router, 'push')
      const wrapper = createWrapper()

      const firstCard = wrapper.find('.store-card')
      await firstCard.trigger('click')

      expect(pushSpy).toHaveBeenCalledWith('/stores/1')
    })
  })

  describe('店舗操作', () => {
    it('編集ボタンクリックで編集モーダルが開く', async () => {
      const storeStore = useStoreStore()
      storeStore.stores = mockStores

      const wrapper = createWrapper()

      const editButton = wrapper.findAll('.btn--small.btn--secondary')[0]
      await editButton.trigger('click')

      expect(wrapper.find('.modal').exists()).toBe(true)
      expect(wrapper.find('.modal__header h2').text()).toBe('店舗編集')
    })

    it('ステータス切り替えボタンが動作する', async () => {
      const storeStore = useStoreStore()
      storeStore.stores = mockStores
      const updateStoreStatusSpy = vi.spyOn(storeStore, 'updateStoreStatus').mockImplementation(() => Promise.resolve(mockStores[0]))

      const wrapper = createWrapper()

      // 最初の店舗（active）のステータス切り替えボタンをクリック
      const statusToggleButton = wrapper.findAll('.btn--small.btn--primary')[0]
      await statusToggleButton.trigger('click')

      expect(updateStoreStatusSpy).toHaveBeenCalledWith(1, 'inactive')
    })

    it('削除ボタンクリックで削除確認モーダルが開く', async () => {
      const storeStore = useStoreStore()
      storeStore.stores = mockStores

      const wrapper = createWrapper()

      const deleteButton = wrapper.findAll('.btn--small.btn--danger')[0]
      await deleteButton.trigger('click')

      expect(wrapper.find('.modal--small').exists()).toBe(true)
      expect(wrapper.find('.modal__header h2').text()).toBe('店舗削除の確認')
      expect(wrapper.text()).toContain('「テスト店舗1」を削除しますか？')
    })
  })

  describe('新規作成', () => {
    it('新規店舗登録ボタンクリックで作成モーダルが開く', async () => {
      const wrapper = createWrapper()

      const createButton = wrapper.find('.btn--primary')
      await createButton.trigger('click')

      expect(wrapper.find('.modal').exists()).toBe(true)
      expect(wrapper.find('.modal__header h2').text()).toBe('新規店舗登録')
    })

    it('新規店舗が正常に作成される', async () => {
      const storeStore = useStoreStore()
      const createStoreSpy = vi.spyOn(storeStore, 'createStore').mockImplementation(() => Promise.resolve(mockStores[0]))

      const wrapper = createWrapper()

      // 作成モーダルを開く
      const createButton = wrapper.find('.btn--primary')
      await createButton.trigger('click')

      // StoreFormコンポーネントからsubmitイベントを発火
      const storeForm = wrapper.findComponent({ name: 'StoreForm' })
      await storeForm.vm.$emit('submit', {
        name: '新規店舗',
        status: 'active'
      })

      expect(createStoreSpy).toHaveBeenCalledWith({
        name: '新規店舗',
        status: 'active'
      })
    })

    it('作成モーダルがキャンセルできる', async () => {
      const wrapper = createWrapper()

      // 作成モーダルを開く
      const createButton = wrapper.find('.btn--primary')
      await createButton.trigger('click')

      expect(wrapper.find('.modal').exists()).toBe(true)

      // キャンセルボタンをクリック
      const cancelButton = wrapper.find('.modal__close')
      await cancelButton.trigger('click')

      expect(wrapper.find('.modal').exists()).toBe(false)
    })
  })

  describe('編集', () => {
    it('店舗が正常に更新される', async () => {
      const storeStore = useStoreStore()
      storeStore.stores = mockStores
      const updateStoreSpy = vi.spyOn(storeStore, 'updateStore').mockImplementation(() => Promise.resolve(mockStores[0]))

      const wrapper = createWrapper()

      // 編集モーダルを開く
      const editButton = wrapper.findAll('.btn--small.btn--secondary')[0]
      await editButton.trigger('click')

      // StoreFormコンポーネントからsubmitイベントを発火
      const storeForm = wrapper.findComponent({ name: 'StoreForm' })
      await storeForm.vm.$emit('submit', {
        name: '更新された店舗',
        status: 'inactive'
      })

      expect(updateStoreSpy).toHaveBeenCalledWith(1, {
        name: '更新された店舗',
        status: 'inactive'
      })
    })

    it('編集モーダルがキャンセルできる', async () => {
      const storeStore = useStoreStore()
      storeStore.stores = mockStores

      const wrapper = createWrapper()

      // 編集モーダルを開く
      const editButton = wrapper.findAll('.btn--small.btn--secondary')[0]
      await editButton.trigger('click')

      expect(wrapper.find('.modal').exists()).toBe(true)

      // StoreFormコンポーネントからcancelイベントを発火
      const storeForm = wrapper.findComponent({ name: 'StoreForm' })
      await storeForm.vm.$emit('cancel')

      expect(wrapper.find('.modal').exists()).toBe(false)
    })
  })

  describe('削除', () => {
    it('店舗が正常に削除される', async () => {
      const storeStore = useStoreStore()
      storeStore.stores = mockStores
      const deleteStoreSpy = vi.spyOn(storeStore, 'deleteStore').mockImplementation(() => Promise.resolve())

      const wrapper = createWrapper()

      // 削除確認モーダルを開く
      const deleteButton = wrapper.findAll('.btn--small.btn--danger')[0]
      await deleteButton.trigger('click')

      // 削除実行ボタンをクリック
      const confirmDeleteButton = wrapper.findAll('.btn--danger')[1] // モーダル内の削除ボタン
      await confirmDeleteButton.trigger('click')

      expect(deleteStoreSpy).toHaveBeenCalledWith(1)
    })

    it('削除がキャンセルできる', async () => {
      const storeStore = useStoreStore()
      storeStore.stores = mockStores

      const wrapper = createWrapper()

      // 削除確認モーダルを開く
      const deleteButton = wrapper.findAll('.btn--small.btn--danger')[0]
      await deleteButton.trigger('click')

      expect(wrapper.find('.modal--small').exists()).toBe(true)

      // キャンセルボタンをクリック
      const cancelButton = wrapper.find('.modal__footer .btn--secondary')
      await cancelButton.trigger('click')

      expect(wrapper.find('.modal--small').exists()).toBe(false)
    })
  })

  describe('ページネーション', () => {
    it('ページネーションが表示される', () => {
      const storeStore = useStoreStore()
      storeStore.totalPages = 3
      storeStore.currentPage = 2

      const wrapper = createWrapper()

      const pagination = wrapper.find('.pagination')
      expect(pagination.exists()).toBe(true)
      expect(pagination.text()).toContain('2 / 3 ページ')
    })

    it('前へボタンが動作する', async () => {
      const storeStore = useStoreStore()
      storeStore.totalPages = 3
      storeStore.currentPage = 2
      const setPageSpy = vi.spyOn(storeStore, 'setPage').mockImplementation(() => Promise.resolve())

      const wrapper = createWrapper()

      const prevButton = wrapper.find('.pagination__btn')
      await prevButton.trigger('click')

      expect(setPageSpy).toHaveBeenCalledWith(1)
    })

    it('次へボタンが動作する', async () => {
      const storeStore = useStoreStore()
      storeStore.totalPages = 3
      storeStore.currentPage = 2
      const setPageSpy = vi.spyOn(storeStore, 'setPage').mockImplementation(() => Promise.resolve())

      const wrapper = createWrapper()

      const nextButton = wrapper.findAll('.pagination__btn')[1]
      await nextButton.trigger('click')

      expect(setPageSpy).toHaveBeenCalledWith(3)
    })

    it('ページサイズ変更が動作する', async () => {
      const storeStore = useStoreStore()
      storeStore.totalPages = 3
      const setItemsPerPageSpy = vi.spyOn(storeStore, 'setItemsPerPage').mockImplementation(() => Promise.resolve())

      const wrapper = createWrapper()

      const pageSizeSelect = wrapper.find('.pagination__size')
      await pageSizeSelect.setValue('25')

      expect(setItemsPerPageSpy).toHaveBeenCalledWith(25)
    })
  })

  describe('ローディング状態', () => {
    it('ローディング中はスピナーが表示される', () => {
      const storeStore = useStoreStore()
      storeStore.loading = true

      const wrapper = createWrapper()

      expect(wrapper.find('.loading-spinner').exists()).toBe(true)
      expect(wrapper.find('.loading-spinner').text()).toBe('読み込み中...')
    })
  })

  describe('エラー状態', () => {
    it('エラー時はエラーメッセージが表示される', () => {
      const storeStore = useStoreStore()
      storeStore.error = 'サーバーエラーが発生しました'

      const wrapper = createWrapper()

      expect(wrapper.find('.error-message').exists()).toBe(true)
      expect(wrapper.find('.error-message').text()).toContain('サーバーエラーが発生しました')
    })

    it('エラーメッセージをクリアできる', async () => {
      const storeStore = useStoreStore()
      storeStore.error = 'サーバーエラーが発生しました'
      const clearErrorSpy = vi.spyOn(storeStore, 'clearError')

      const wrapper = createWrapper()

      const closeButton = wrapper.find('.error-message__close')
      await closeButton.trigger('click')

      expect(clearErrorSpy).toHaveBeenCalled()
    })
  })

  describe('空の状態', () => {
    it('店舗がない場合は空の状態が表示される', () => {
      const storeStore = useStoreStore()
      storeStore.stores = []
      storeStore.loading = false

      const wrapper = createWrapper()

      expect(wrapper.find('.empty-state').exists()).toBe(true)
      expect(wrapper.find('.empty-state').text()).toContain('店舗が見つかりません。')
    })

    it('空の状態から新規作成ボタンが動作する', async () => {
      const storeStore = useStoreStore()
      storeStore.stores = []
      storeStore.loading = false

      const wrapper = createWrapper()

      const createButton = wrapper.find('.empty-state .btn--primary')
      await createButton.trigger('click')

      expect(wrapper.find('.modal').exists()).toBe(true)
    })
  })
})