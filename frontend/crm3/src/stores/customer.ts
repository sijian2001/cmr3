import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export interface Customer {
  id?: number
  name: string
  email: string
  phone?: string
  address?: string
  created_at?: string
  updated_at?: string
}

export interface CustomerSearchParams {
  name?: string
  email?: string
  phone?: string
  page?: number
  limit?: number
}

export interface PaginatedCustomerResponse {
  customers: Customer[]
  total: number
  page: number
  limit: number
  total_pages: number
}

export const useCustomerStore = defineStore('customer', () => {
  const customers = ref<Customer[]>([])
  const currentCustomer = ref<Customer | null>(null)
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
  const isEmpty = computed(() => customers.value.length === 0)

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

  const setCustomers = (data: PaginatedCustomerResponse) => {
    customers.value = data.customers
    pagination.value = {
      total: data.total,
      page: data.page,
      limit: data.limit,
      total_pages: data.total_pages
    }
  }

  const setCurrentCustomer = (customer: Customer | null) => {
    currentCustomer.value = customer
  }

  const addCustomer = (customer: Customer) => {
    customers.value.unshift(customer)
    pagination.value.total += 1
  }

  const updateCustomer = (updatedCustomer: Customer) => {
    const index = customers.value.findIndex(c => c.id === updatedCustomer.id)
    if (index !== -1) {
      customers.value[index] = updatedCustomer
    }
    if (currentCustomer.value?.id === updatedCustomer.id) {
      currentCustomer.value = updatedCustomer
    }
  }

  const removeCustomer = (customerId: number) => {
    customers.value = customers.value.filter(c => c.id !== customerId)
    pagination.value.total -= 1
    if (currentCustomer.value?.id === customerId) {
      currentCustomer.value = null
    }
  }

  // API Methods (モックとして実装、後でAPI呼び出しに置き換え)
  const fetchCustomers = async (params: CustomerSearchParams = {}) => {
    try {
      setLoading(true)
      clearError()

      // TODO: 実際のAPI呼び出しに置き換える
      await new Promise(resolve => setTimeout(resolve, 500)) // モック遅延

      const mockData: PaginatedCustomerResponse = {
        customers: [
          { id: 1, name: '山田太郎', email: 'yamada@example.com', phone: '090-1234-5678', created_at: '2024-01-01' },
          { id: 2, name: '佐藤花子', email: 'sato@example.com', phone: '090-8765-4321', created_at: '2024-01-02' }
        ],
        total: 2,
        page: params.page || 1,
        limit: params.limit || 10,
        total_pages: 1
      }

      setCustomers(mockData)
    } catch (err) {
      setError(err instanceof Error ? err.message : '顧客一覧の取得に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  const fetchCustomer = async (id: number) => {
    try {
      setLoading(true)
      clearError()

      // TODO: 実際のAPI呼び出しに置き換える
      await new Promise(resolve => setTimeout(resolve, 300))

      const mockCustomer: Customer = {
        id: id,
        name: '山田太郎',
        email: 'yamada@example.com',
        phone: '090-1234-5678',
        address: '東京都渋谷区1-1-1',
        created_at: '2024-01-01'
      }

      setCurrentCustomer(mockCustomer)
      return mockCustomer
    } catch (err) {
      setError(err instanceof Error ? err.message : '顧客情報の取得に失敗しました')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const createCustomer = async (customerData: Omit<Customer, 'id' | 'created_at' | 'updated_at'>) => {
    try {
      setLoading(true)
      clearError()

      // バリデーション
      if (!customerData.name.trim()) {
        throw new Error('名前は必須です')
      }
      if (!customerData.email.trim()) {
        throw new Error('メールアドレスは必須です')
      }
      if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(customerData.email)) {
        throw new Error('有効なメールアドレスを入力してください')
      }

      // TODO: 実際のAPI呼び出しに置き換える
      await new Promise(resolve => setTimeout(resolve, 500))

      const newCustomer: Customer = {
        ...customerData,
        id: Date.now(), // モックID
        created_at: new Date().toISOString()
      }

      addCustomer(newCustomer)
      return newCustomer
    } catch (err) {
      setError(err instanceof Error ? err.message : '顧客の作成に失敗しました')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const updateCustomerData = async (id: number, customerData: Omit<Customer, 'id' | 'created_at' | 'updated_at'>) => {
    try {
      setLoading(true)
      clearError()

      // バリデーション
      if (!customerData.name.trim()) {
        throw new Error('名前は必須です')
      }
      if (!customerData.email.trim()) {
        throw new Error('メールアドレスは必須です')
      }
      if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(customerData.email)) {
        throw new Error('有効なメールアドレスを入力してください')
      }

      // TODO: 実際のAPI呼び出しに置き換える
      await new Promise(resolve => setTimeout(resolve, 500))

      const updatedCustomer: Customer = {
        ...customerData,
        id: id,
        updated_at: new Date().toISOString()
      }

      updateCustomer(updatedCustomer)
      return updatedCustomer
    } catch (err) {
      setError(err instanceof Error ? err.message : '顧客情報の更新に失敗しました')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const deleteCustomer = async (id: number) => {
    try {
      setLoading(true)
      clearError()

      // TODO: 実際のAPI呼び出しに置き換える
      await new Promise(resolve => setTimeout(resolve, 300))

      removeCustomer(id)
    } catch (err) {
      setError(err instanceof Error ? err.message : '顧客の削除に失敗しました')
      throw err
    } finally {
      setLoading(false)
    }
  }

  return {
    // State
    customers,
    currentCustomer,
    loading,
    error,
    pagination,

    // Getters
    isLoading,
    hasError,
    isEmpty,

    // Actions
    clearError,
    setCurrentCustomer,
    fetchCustomers,
    fetchCustomer,
    createCustomer,
    updateCustomerData,
    deleteCustomer
  }
})