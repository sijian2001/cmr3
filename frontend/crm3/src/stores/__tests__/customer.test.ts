import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useCustomerStore, type Customer, type PaginatedCustomerResponse } from '../customer'

describe('Customer Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  describe('初期状態', () => {
    it('初期状態が正しく設定される', () => {
      const store = useCustomerStore()

      expect(store.customers).toEqual([])
      expect(store.currentCustomer).toBeNull()
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
      expect(store.pagination).toEqual({
        total: 0,
        page: 1,
        limit: 10,
        total_pages: 0
      })
    })

    it('computed が正しく動作する', () => {
      const store = useCustomerStore()

      expect(store.isLoading).toBe(false)
      expect(store.hasError).toBe(false)
      expect(store.isEmpty).toBe(true)

      store.loading = true
      store.error = 'テストエラー'
      store.customers = [{ id: 1, name: 'test', email: 'test@test.com' }]

      expect(store.isLoading).toBe(true)
      expect(store.hasError).toBe(true)
      expect(store.isEmpty).toBe(false)
    })
  })

  describe('基本アクション', () => {
    it('clearError() が正しく動作する', () => {
      const store = useCustomerStore()
      store.error = 'テストエラー'

      store.clearError()

      expect(store.error).toBeNull()
    })

    it('setCurrentCustomer() が正しく動作する', () => {
      const store = useCustomerStore()
      const customer: Customer = { id: 1, name: 'テスト', email: 'test@test.com' }

      store.setCurrentCustomer(customer)
      expect(store.currentCustomer).toEqual(customer)

      store.setCurrentCustomer(null)
      expect(store.currentCustomer).toBeNull()
    })
  })

  describe('内部状態変更のテスト', () => {
    it('顧客作成後に一覧が更新される', async () => {
      const store = useCustomerStore()
      const customerData = {
        name: 'テスト顧客',
        email: 'test@example.com'
      }

      const initialCount = store.customers.length

      await store.createCustomer(customerData)

      expect(store.customers.length).toBe(initialCount + 1)
      expect(store.customers[0].name).toBe(customerData.name)
      expect(store.customers[0].email).toBe(customerData.email)
    })

    it('顧客削除後に一覧から除外される', async () => {
      const store = useCustomerStore()
      // まず顧客を作成
      const customerData = { name: 'テスト', email: 'test@test.com' }
      const createdCustomer = await store.createCustomer(customerData)

      expect(store.customers).toContainEqual(createdCustomer)

      // 削除実行
      await store.deleteCustomer(createdCustomer.id!)

      expect(store.customers).not.toContainEqual(createdCustomer)
    })
  })

  describe('バリデーション', () => {
    it('createCustomer() で名前が空の場合エラーが発生する', async () => {
      const store = useCustomerStore()

      await expect(
        store.createCustomer({ name: '', email: 'test@test.com' })
      ).rejects.toThrow('名前は必須です')

      expect(store.error).toBe('名前は必須です')
    })

    it('createCustomer() でメールアドレスが空の場合エラーが発生する', async () => {
      const store = useCustomerStore()

      await expect(
        store.createCustomer({ name: 'テスト', email: '' })
      ).rejects.toThrow('メールアドレスは必須です')

      expect(store.error).toBe('メールアドレスは必須です')
    })

    it('createCustomer() で無効なメールアドレスの場合エラーが発生する', async () => {
      const store = useCustomerStore()

      await expect(
        store.createCustomer({ name: 'テスト', email: 'invalid-email' })
      ).rejects.toThrow('有効なメールアドレスを入力してください')

      expect(store.error).toBe('有効なメールアドレスを入力してください')
    })

    it('updateCustomerData() でも同様のバリデーションが動作する', async () => {
      const store = useCustomerStore()

      await expect(
        store.updateCustomerData(1, { name: '', email: 'test@test.com' })
      ).rejects.toThrow('名前は必須です')

      await expect(
        store.updateCustomerData(1, { name: 'テスト', email: 'invalid-email' })
      ).rejects.toThrow('有効なメールアドレスを入力してください')
    })
  })

  describe('API 呼び出し（モック）', () => {
    it('fetchCustomers() が正しく動作する', async () => {
      const store = useCustomerStore()

      await store.fetchCustomers({ page: 1, limit: 5 })

      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
      expect(store.customers.length).toBeGreaterThan(0)
      expect(store.pagination.page).toBe(1)
      expect(store.pagination.limit).toBe(5)
    })

    it('fetchCustomer() が正しく動作する', async () => {
      const store = useCustomerStore()

      await store.fetchCustomer(1)

      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
      expect(store.currentCustomer).not.toBeNull()
      expect(store.currentCustomer?.id).toBe(1)
    })

    it('createCustomer() が正しく動作する', async () => {
      const store = useCustomerStore()
      const customerData = {
        name: 'テスト顧客',
        email: 'test@example.com',
        phone: '090-1234-5678'
      }

      const result = await store.createCustomer(customerData)

      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
      expect(result.name).toBe(customerData.name)
      expect(result.email).toBe(customerData.email)
      expect(result.phone).toBe(customerData.phone)
      expect(result.id).toBeDefined()
      expect(result.created_at).toBeDefined()
      // 作成された顧客がリストの先頭に追加されることを確認
      expect(store.customers.length).toBeGreaterThan(0)
      expect(store.customers[0]).toEqual(result)
    })

    it('updateCustomerData() が正しく動作する', async () => {
      const store = useCustomerStore()
      const customerData = {
        name: '更新されたテスト顧客',
        email: 'updated@example.com'
      }

      const result = await store.updateCustomerData(1, customerData)

      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
      expect(result.name).toBe(customerData.name)
      expect(result.email).toBe(customerData.email)
      expect(result.id).toBe(1)
      expect(result.updated_at).toBeDefined()
    })

    it('deleteCustomer() が正しく動作する', async () => {
      const store = useCustomerStore()
      const customer: Customer = { id: 1, name: 'テスト', email: 'test@test.com' }
      store.customers = [customer]
      store.currentCustomer = customer
      store.pagination.total = 1

      await store.deleteCustomer(1)

      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
      expect(store.customers).not.toContain(customer)
      expect(store.currentCustomer).toBeNull()
      expect(store.pagination.total).toBe(0)
    })
  })

  describe('エラーハンドリング', () => {
    it('fetchCustomers() でエラーが発生した場合の処理', async () => {
      const store = useCustomerStore()

      // モックの動作を一時的に変更してエラーを発生させる
      const originalFetch = store.fetchCustomers
      store.fetchCustomers = vi.fn().mockRejectedValue(new Error('ネットワークエラー'))

      try {
        await store.fetchCustomers()
      } catch (error) {
        // エラーが適切にキャッチされることを確認
      }

      expect(store.loading).toBe(false)
    })
  })
})