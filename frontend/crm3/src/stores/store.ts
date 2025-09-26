import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface Store {
  id: number
  name: string
  address?: string
  phone?: string
  email?: string
  status: 'active' | 'inactive' | 'maintenance'
  manager_id?: number
  created_at: string
  updated_at: string
}

export interface StoreCreateRequest {
  name: string
  address?: string
  phone?: string
  email?: string
  status: 'active' | 'inactive' | 'maintenance'
  manager_id?: number
}

export interface StoreUpdateRequest {
  name: string
  address?: string
  phone?: string
  email?: string
  status: 'active' | 'inactive' | 'maintenance'
  manager_id?: number
}

export interface StoreSearchParams {
  name?: string
  status?: string
  page?: number
  limit?: number
}

export interface PaginatedStoreResponse {
  stores: Store[]
  total: number
  page: number
  limit: number
  totalPages: number
}

export const useStoreStore = defineStore('store', () => {
  // State
  const stores = ref<Store[]>([])
  const currentStore = ref<Store | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const total = ref(0)
  const currentPage = ref(1)
  const itemsPerPage = ref(10)

  // Computed
  const totalPages = computed(() => Math.ceil(total.value / itemsPerPage.value))

  const activeStores = computed(() =>
    stores.value.filter(store => store.status === 'active')
  )

  const inactiveStores = computed(() =>
    stores.value.filter(store => store.status === 'inactive')
  )

  const maintenanceStores = computed(() =>
    stores.value.filter(store => store.status === 'maintenance')
  )

  const statusCounts = computed(() => ({
    active: activeStores.value.length,
    inactive: inactiveStores.value.length,
    maintenance: maintenanceStores.value.length,
    total: stores.value.length
  }))

  // Actions
  const fetchStores = async (params: StoreSearchParams = {}) => {
    loading.value = true
    error.value = null

    try {
      // デフォルトパラメータ
      const queryParams = {
        page: params.page || currentPage.value,
        limit: params.limit || itemsPerPage.value,
        ...params
      }

      const query = new URLSearchParams()
      Object.entries(queryParams).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          query.append(key, value.toString())
        }
      })

      const response = await fetch(`/api/stores?${query}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to fetch stores')
      }

      const data: PaginatedStoreResponse = await response.json()

      stores.value = data.stores || []
      total.value = data.total
      currentPage.value = data.page

    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch stores'
      stores.value = []
      total.value = 0
    } finally {
      loading.value = false
    }
  }

  const fetchStore = async (id: number) => {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`/api/stores/${id}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to fetch store')
      }

      const data: Store = await response.json()
      currentStore.value = data
      return data

    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch store'
      currentStore.value = null
      throw err
    } finally {
      loading.value = false
    }
  }

  const createStore = async (storeData: StoreCreateRequest) => {
    loading.value = true
    error.value = null

    try {
      const response = await fetch('/api/stores', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(storeData),
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to create store')
      }

      const newStore: Store = await response.json()

      // 現在の一覧に追加
      stores.value.unshift(newStore)
      total.value += 1

      return newStore

    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create store'
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateStore = async (id: number, storeData: StoreUpdateRequest) => {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`/api/stores/${id}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(storeData),
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to update store')
      }

      const updatedStore: Store = await response.json()

      // 一覧内の店舗を更新
      const index = stores.value.findIndex(store => store.id === id)
      if (index !== -1) {
        stores.value[index] = updatedStore
      }

      // 現在の店舗も更新
      if (currentStore.value && currentStore.value.id === id) {
        currentStore.value = updatedStore
      }

      return updatedStore

    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update store'
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteStore = async (id: number) => {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`/api/stores/${id}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to delete store')
      }

      // 一覧から削除
      const index = stores.value.findIndex(store => store.id === id)
      if (index !== -1) {
        stores.value.splice(index, 1)
        total.value -= 1
      }

      // 現在の店舗をクリア
      if (currentStore.value && currentStore.value.id === id) {
        currentStore.value = null
      }

    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to delete store'
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateStoreStatus = async (id: number, status: Store['status']) => {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`/api/stores/${id}/status`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ status }),
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to update store status')
      }

      const updatedStore: Store = await response.json()

      // 一覧内の店舗を更新
      const index = stores.value.findIndex(store => store.id === id)
      if (index !== -1) {
        stores.value[index] = updatedStore
      }

      // 現在の店舗も更新
      if (currentStore.value && currentStore.value.id === id) {
        currentStore.value = updatedStore
      }

      return updatedStore

    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update store status'
      throw err
    } finally {
      loading.value = false
    }
  }

  const searchStores = async (searchParams: StoreSearchParams) => {
    currentPage.value = 1 // 検索時はページをリセット
    await fetchStores(searchParams)
  }

  const setPage = async (page: number) => {
    currentPage.value = page
    await fetchStores()
  }

  const setItemsPerPage = async (limit: number) => {
    itemsPerPage.value = limit
    currentPage.value = 1 // ページサイズ変更時はページをリセット
    await fetchStores()
  }

  const clearError = () => {
    error.value = null
  }

  const clearCurrentStore = () => {
    currentStore.value = null
  }

  const getStoreStatusLabel = (status: Store['status']): string => {
    const statusLabels: Record<Store['status'], string> = {
      active: '営業中',
      inactive: '休業中',
      maintenance: 'メンテナンス中'
    }
    return statusLabels[status] || status
  }

  const getStoreStatusColor = (status: Store['status']): string => {
    const statusColors: Record<Store['status'], string> = {
      active: 'success',
      inactive: 'warning',
      maintenance: 'error'
    }
    return statusColors[status] || 'default'
  }

  return {
    // State
    stores,
    currentStore,
    loading,
    error,
    total,
    currentPage,
    itemsPerPage,

    // Computed
    totalPages,
    activeStores,
    inactiveStores,
    maintenanceStores,
    statusCounts,

    // Actions
    fetchStores,
    fetchStore,
    createStore,
    updateStore,
    deleteStore,
    updateStoreStatus,
    searchStores,
    setPage,
    setItemsPerPage,
    clearError,
    clearCurrentStore,
    getStoreStatusLabel,
    getStoreStatusColor
  }
})