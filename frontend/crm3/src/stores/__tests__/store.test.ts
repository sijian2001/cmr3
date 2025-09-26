import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useStoreStore, type Store, type StoreCreateRequest, type StoreUpdateRequest } from '../store'

// Mock fetch
global.fetch = vi.fn()

const mockStore: Store = {
  id: 1,
  name: 'テスト店舗',
  address: '東京都渋谷区',
  phone: '03-1234-5678',
  email: 'test@example.com',
  status: 'active',
  manager_id: 1,
  created_at: '2023-01-01T00:00:00Z',
  updated_at: '2023-01-01T00:00:00Z'
}

const mockStoreList = [
  mockStore,
  {
    id: 2,
    name: 'テスト店舗2',
    address: '大阪府大阪市',
    phone: '06-1234-5678',
    email: 'test2@example.com',
    status: 'inactive' as const,
    manager_id: 2,
    created_at: '2023-01-02T00:00:00Z',
    updated_at: '2023-01-02T00:00:00Z'
  }
]

describe('Store Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    localStorage.setItem('token', 'test-token')
  })

  describe('初期状態', () => {
    it('初期値が正しく設定されている', () => {
      const store = useStoreStore()

      expect(store.stores).toEqual([])
      expect(store.currentStore).toBeNull()
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
      expect(store.total).toBe(0)
      expect(store.currentPage).toBe(1)
      expect(store.itemsPerPage).toBe(10)
    })

    it('computed values が正しく動作する', () => {
      const store = useStoreStore()
      store.stores = mockStoreList

      expect(store.activeStores).toHaveLength(1)
      expect(store.inactiveStores).toHaveLength(1)
      expect(store.maintenanceStores).toHaveLength(0)
      expect(store.statusCounts).toEqual({
        active: 1,
        inactive: 1,
        maintenance: 0,
        total: 2
      })
    })
  })

  describe('fetchStores', () => {
    it('成功時に店舗一覧を取得する', async () => {
      const mockResponse = {
        stores: mockStoreList,
        total: 2,
        page: 1,
        limit: 10,
        totalPages: 1
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse
      })

      const store = useStoreStore()
      await store.fetchStores()

      expect(store.stores).toEqual(mockStoreList)
      expect(store.total).toBe(2)
      expect(store.currentPage).toBe(1)
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
    })

    it('検索パラメータありで店舗一覧を取得する', async () => {
      const mockResponse = {
        stores: [mockStore],
        total: 1,
        page: 1,
        limit: 10,
        totalPages: 1
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse
      })

      const store = useStoreStore()
      await store.fetchStores({
        name: 'テスト',
        status: 'active'
      })

      expect(global.fetch).toHaveBeenCalledWith(
        '/api/stores?page=1&limit=10&name=%E3%83%86%E3%82%B9%E3%83%88&status=active',
        {
          headers: {
            'Authorization': 'Bearer test-token',
            'Content-Type': 'application/json'
          }
        }
      )
    })

    it('エラー時にエラーメッセージを設定する', async () => {
      ;(global.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Server error' })
      })

      const store = useStoreStore()
      await store.fetchStores()

      expect(store.stores).toEqual([])
      expect(store.total).toBe(0)
      expect(store.loading).toBe(false)
      expect(store.error).toBe('Server error')
    })

    it('ネットワークエラー時にエラーメッセージを設定する', async () => {
      ;(global.fetch as any).mockRejectedValueOnce(new Error('Network error'))

      const store = useStoreStore()
      await store.fetchStores()

      expect(store.error).toBe('Network error')
      expect(store.loading).toBe(false)
    })
  })

  describe('fetchStore', () => {
    it('成功時に店舗詳細を取得する', async () => {
      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockStore
      })

      const store = useStoreStore()
      const result = await store.fetchStore(1)

      expect(result).toEqual(mockStore)
      expect(store.currentStore).toEqual(mockStore)
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()

      expect(global.fetch).toHaveBeenCalledWith(
        '/api/stores/1',
        {
          headers: {
            'Authorization': 'Bearer test-token',
            'Content-Type': 'application/json'
          }
        }
      )
    })

    it('エラー時にエラーを投げる', async () => {
      ;(global.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Store not found' })
      })

      const store = useStoreStore()

      await expect(store.fetchStore(999)).rejects.toThrow('Store not found')
      expect(store.currentStore).toBeNull()
      expect(store.error).toBe('Store not found')
    })
  })

  describe('createStore', () => {
    it('成功時に新しい店舗を作成する', async () => {
      const createData: StoreCreateRequest = {
        name: '新規店舗',
        address: '神奈川県横浜市',
        phone: '045-1234-5678',
        email: 'new@example.com',
        status: 'active',
        manager_id: 3
      }

      const createdStore: Store = {
        id: 3,
        ...createData,
        created_at: '2023-01-03T00:00:00Z',
        updated_at: '2023-01-03T00:00:00Z'
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => createdStore
      })

      const store = useStoreStore()
      const result = await store.createStore(createData)

      expect(result).toEqual(createdStore)
      expect(store.stores[0]).toEqual(createdStore)
      expect(store.total).toBe(1)
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()

      expect(global.fetch).toHaveBeenCalledWith(
        '/api/stores',
        {
          method: 'POST',
          headers: {
            'Authorization': 'Bearer test-token',
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(createData)
        }
      )
    })

    it('エラー時にエラーを投げる', async () => {
      const createData: StoreCreateRequest = {
        name: '新規店舗',
        status: 'active'
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Validation error' })
      })

      const store = useStoreStore()

      await expect(store.createStore(createData)).rejects.toThrow('Validation error')
      expect(store.error).toBe('Validation error')
    })
  })

  describe('updateStore', () => {
    it('成功時に店舗を更新する', async () => {
      const updateData: StoreUpdateRequest = {
        name: '更新された店舗',
        address: '埼玉県さいたま市',
        phone: '048-1234-5678',
        email: 'updated@example.com',
        status: 'inactive',
        manager_id: 2
      }

      const updatedStore: Store = {
        ...mockStore,
        ...updateData,
        updated_at: '2023-01-04T00:00:00Z'
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => updatedStore
      })

      const store = useStoreStore()
      store.stores = [mockStore]
      store.currentStore = mockStore

      const result = await store.updateStore(1, updateData)

      expect(result).toEqual(updatedStore)
      expect(store.stores[0]).toEqual(updatedStore)
      expect(store.currentStore).toEqual(updatedStore)
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()

      expect(global.fetch).toHaveBeenCalledWith(
        '/api/stores/1',
        {
          method: 'PUT',
          headers: {
            'Authorization': 'Bearer test-token',
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(updateData)
        }
      )
    })

    it('エラー時にエラーを投げる', async () => {
      const updateData: StoreUpdateRequest = {
        name: '更新された店舗',
        status: 'active'
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Store not found' })
      })

      const store = useStoreStore()

      await expect(store.updateStore(999, updateData)).rejects.toThrow('Store not found')
      expect(store.error).toBe('Store not found')
    })
  })

  describe('deleteStore', () => {
    it('成功時に店舗を削除する', async () => {
      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true
      })

      const store = useStoreStore()
      store.stores = [mockStore]
      store.currentStore = mockStore
      store.total = 1

      await store.deleteStore(1)

      expect(store.stores).toEqual([])
      expect(store.total).toBe(0)
      expect(store.currentStore).toBeNull()
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()

      expect(global.fetch).toHaveBeenCalledWith(
        '/api/stores/1',
        {
          method: 'DELETE',
          headers: {
            'Authorization': 'Bearer test-token',
            'Content-Type': 'application/json'
          }
        }
      )
    })

    it('エラー時にエラーを投げる', async () => {
      ;(global.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Store not found' })
      })

      const store = useStoreStore()

      await expect(store.deleteStore(999)).rejects.toThrow('Store not found')
      expect(store.error).toBe('Store not found')
    })
  })

  describe('updateStoreStatus', () => {
    it('成功時に店舗ステータスを更新する', async () => {
      const updatedStore: Store = {
        ...mockStore,
        status: 'maintenance',
        updated_at: '2023-01-05T00:00:00Z'
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => updatedStore
      })

      const store = useStoreStore()
      store.stores = [mockStore]
      store.currentStore = mockStore

      const result = await store.updateStoreStatus(1, 'maintenance')

      expect(result).toEqual(updatedStore)
      expect(store.stores[0]).toEqual(updatedStore)
      expect(store.currentStore).toEqual(updatedStore)

      expect(global.fetch).toHaveBeenCalledWith(
        '/api/stores/1/status',
        {
          method: 'POST',
          headers: {
            'Authorization': 'Bearer test-token',
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ status: 'maintenance' })
        }
      )
    })
  })

  describe('ユーティリティメソッド', () => {
    it('searchStores がページをリセットして検索する', async () => {
      const mockResponse = {
        stores: mockStoreList,
        total: 2,
        page: 1,
        limit: 10,
        totalPages: 1
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse
      })

      const store = useStoreStore()
      store.currentPage = 3 // 現在のページを3に設定

      await store.searchStores({ name: 'テスト' })

      expect(store.currentPage).toBe(1) // ページがリセットされている
    })

    it('setPage がページを設定してfetchStoresを呼ぶ', async () => {
      const mockResponse = {
        stores: mockStoreList,
        total: 2,
        page: 2,
        limit: 10,
        totalPages: 1
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse
      })

      const store = useStoreStore()
      await store.setPage(2)

      expect(store.currentPage).toBe(2)
    })

    it('setItemsPerPage がページサイズを設定してページをリセットする', async () => {
      const mockResponse = {
        stores: mockStoreList,
        total: 2,
        page: 1,
        limit: 25,
        totalPages: 1
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse
      })

      const store = useStoreStore()
      store.currentPage = 3

      await store.setItemsPerPage(25)

      expect(store.itemsPerPage).toBe(25)
      expect(store.currentPage).toBe(1) // ページがリセットされている
    })

    it('clearError がエラーをクリアする', () => {
      const store = useStoreStore()
      store.error = 'Test error'

      store.clearError()

      expect(store.error).toBeNull()
    })

    it('clearCurrentStore が現在の店舗をクリアする', () => {
      const store = useStoreStore()
      store.currentStore = mockStore

      store.clearCurrentStore()

      expect(store.currentStore).toBeNull()
    })

    it('getStoreStatusLabel が正しいラベルを返す', () => {
      const store = useStoreStore()

      expect(store.getStoreStatusLabel('active')).toBe('営業中')
      expect(store.getStoreStatusLabel('inactive')).toBe('休業中')
      expect(store.getStoreStatusLabel('maintenance')).toBe('メンテナンス中')
    })

    it('getStoreStatusColor が正しい色を返す', () => {
      const store = useStoreStore()

      expect(store.getStoreStatusColor('active')).toBe('success')
      expect(store.getStoreStatusColor('inactive')).toBe('warning')
      expect(store.getStoreStatusColor('maintenance')).toBe('error')
    })
  })
})