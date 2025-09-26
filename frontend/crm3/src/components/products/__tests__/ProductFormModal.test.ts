import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import ProductFormModal from '../ProductFormModal.vue'
import type { Product } from '@/stores/product'

// Pinia セットアップ
beforeEach(() => {
  setActivePinia(createPinia())
})

afterEach(() => {
  vi.clearAllMocks()
})

describe('ProductFormModal', () => {
  describe('コンポーネントの表示', () => {
    it('新規作成モードで正しく表示される', () => {
      const wrapper = mount(ProductFormModal, {
        props: {
          product: null
        }
      })

      expect(wrapper.find('h2').text()).toBe('新規製品登録')
      expect(wrapper.find('#name').exists()).toBe(true)
      expect(wrapper.find('#sku').exists()).toBe(true)
      expect(wrapper.find('#price').exists()).toBe(true)
      expect(wrapper.find('#stock_quantity').exists()).toBe(true)
      expect(wrapper.find('#description').exists()).toBe(true)
    })

    it('編集モードで正しく表示される', async () => {
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

      const wrapper = mount(ProductFormModal, {
        props: {
          product: mockProduct
        }
      })

      // onMountedの処理が完了するまで待つ
      await wrapper.vm.$nextTick()

      expect(wrapper.find('h2').text()).toBe('製品情報編集')
      expect((wrapper.find('#name').element as HTMLInputElement).value).toBe('テスト製品')
      expect((wrapper.find('#sku').element as HTMLInputElement).value).toBe('TEST-001')
      expect((wrapper.find('#price').element as HTMLInputElement).value).toBe('1000')
      expect((wrapper.find('#stock_quantity').element as HTMLInputElement).value).toBe('50')
      expect((wrapper.find('#description').element as HTMLTextAreaElement).value).toBe('テスト用の製品です')
    })

    it('必須フィールドが正しくマークされている', () => {
      const wrapper = mount(ProductFormModal, {
        props: {
          product: null
        }
      })

      const requiredFields = wrapper.findAll('.required')
      expect(requiredFields).toHaveLength(4) // name, sku, price, stock_quantity

      expect(wrapper.find('label[for="name"]').text()).toContain('*')
      expect(wrapper.find('label[for="sku"]').text()).toContain('*')
      expect(wrapper.find('label[for="price"]').text()).toContain('*')
      expect(wrapper.find('label[for="stock_quantity"]').text()).toContain('*')
    })
  })

  describe('フォームバリデーション', () => {
    it('製品名が空の場合にエラーが表示される', async () => {
      const wrapper = mount(ProductFormModal, {
        props: {
          product: null
        }
      })

      await wrapper.find('form').trigger('submit')

      expect(wrapper.find('.error-text').text()).toContain('製品名は必須です')
    })

    it('SKUが空の場合にエラーが表示される', async () => {
      const wrapper = mount(ProductFormModal, {
        props: {
          product: null
        }
      })

      const nameInput = wrapper.find('#name')
      await nameInput.setValue('テスト製品')

      await wrapper.find('form').trigger('submit')

      expect(wrapper.text()).toContain('SKUは必須です')
    })

    it('価格が0以下の場合にエラーが表示される', async () => {
      const wrapper = mount(ProductFormModal, {
        props: {
          product: null
        }
      })

      const nameInput = wrapper.find('#name')
      const skuInput = wrapper.find('#sku')
      const priceInput = wrapper.find('#price')
      const stockInput = wrapper.find('#stock_quantity')

      await nameInput.setValue('テスト製品')
      await skuInput.setValue('TEST-001')
      await priceInput.setValue('0')
      await stockInput.setValue('10')

      await wrapper.find('form').trigger('submit')

      expect(wrapper.text()).toContain('価格は正数である必要があります')
    })

    it('在庫数量が負数の場合にエラーが表示される', async () => {
      const wrapper = mount(ProductFormModal, {
        props: {
          product: null
        }
      })

      const nameInput = wrapper.find('#name')
      const skuInput = wrapper.find('#sku')
      const priceInput = wrapper.find('#price')
      const stockInput = wrapper.find('#stock_quantity')

      await nameInput.setValue('テスト製品')
      await skuInput.setValue('TEST-001')
      await priceInput.setValue('1000')
      await stockInput.setValue('-1')

      await wrapper.find('form').trigger('submit')

      expect(wrapper.text()).toContain('在庫数量は0以上である必要があります')
    })

    it('文字数制限が正しく機能する', async () => {
      const wrapper = mount(ProductFormModal, {
        props: {
          product: null
        }
      })

      const nameInput = wrapper.find('#name')
      const longName = 'あ'.repeat(101) // 101文字

      await nameInput.setValue(longName)
      await wrapper.find('form').trigger('submit')

      expect(wrapper.text()).toContain('製品名は100文字以内で入力してください')
    })
  })

  describe('イベント処理', () => {
    it('キャンセルボタンクリック時にcancelイベントが発火される', async () => {
      const wrapper = mount(ProductFormModal, {
        props: {
          product: null
        }
      })

      await wrapper.find('.btn-cancel').trigger('click')

      expect(wrapper.emitted('cancel')).toBeTruthy()
    })

    it('×ボタンクリック時にcancelイベントが発火される', async () => {
      const wrapper = mount(ProductFormModal, {
        props: {
          product: null
        }
      })

      await wrapper.find('.btn-close').trigger('click')

      expect(wrapper.emitted('cancel')).toBeTruthy()
    })

    it('正しいデータでフォーム送信時にsaveイベントが発火される', async () => {
      const wrapper = mount(ProductFormModal, {
        props: {
          product: null
        }
      })

      const nameInput = wrapper.find('#name')
      const skuInput = wrapper.find('#sku')
      const priceInput = wrapper.find('#price')
      const stockInput = wrapper.find('#stock_quantity')
      const descriptionInput = wrapper.find('#description')

      await nameInput.setValue('テスト製品')
      await skuInput.setValue('TEST-001')
      await priceInput.setValue('1000')
      await stockInput.setValue('50')
      await descriptionInput.setValue('テスト用の製品です')

      await wrapper.find('form').trigger('submit')

      expect(wrapper.emitted('save')).toBeTruthy()
      const saveEvent = wrapper.emitted('save')
      expect(saveEvent?.[0][0]).toEqual({
        name: 'テスト製品',
        sku: 'TEST-001',
        price: 1000,
        stock_quantity: 50,
        description: 'テスト用の製品です'
      })
    })

    it('モーダル外クリック時にcancelイベントが発火される', async () => {
      const wrapper = mount(ProductFormModal, {
        props: {
          product: null
        }
      })

      const overlay = wrapper.find('.modal-overlay')
      // オーバーレイをクリック（バブリングを防ぐため直接イベントを発火）
      await overlay.trigger('click')

      expect(wrapper.emitted('cancel')).toBeTruthy()
    })
  })

  describe('フォームの状態管理', () => {
    it('送信中は保存ボタンが無効になる', async () => {
      const wrapper = mount(ProductFormModal, {
        props: {
          product: null
        }
      })

      // フォームに有効なデータを入力
      await wrapper.find('#name').setValue('テスト製品')
      await wrapper.find('#sku').setValue('TEST-001')
      await wrapper.find('#price').setValue('1000')
      await wrapper.find('#stock_quantity').setValue('50')

      // データを変更して送信中状態をシミュレート（Vue 3 Composition APIでは直接アクセス）
      const vm = wrapper.vm as any
      vm.isSubmitting = true
      await wrapper.vm.$nextTick()

      const saveButton = wrapper.find('.btn-save')
      expect(saveButton.attributes('disabled')).toBeDefined()
      expect(saveButton.text()).toBe('保存中...')
    })

    it('データが空の場合は説明フィールドがundefinedになる', async () => {
      const wrapper = mount(ProductFormModal, {
        props: {
          product: null
        }
      })

      await wrapper.find('#name').setValue('テスト製品')
      await wrapper.find('#sku').setValue('TEST-001')
      await wrapper.find('#price').setValue('1000')
      await wrapper.find('#stock_quantity').setValue('50')
      // descriptionは空のまま

      await wrapper.find('form').trigger('submit')

      expect(wrapper.emitted('save')).toBeTruthy()
      const saveEvent = wrapper.emitted('save')
      expect(saveEvent?.[0][0].description).toBeUndefined()
    })
  })

  describe('レスポンシブ対応', () => {
    it('モバイル表示クラスが存在する', () => {
      const wrapper = mount(ProductFormModal, {
        props: {
          product: null
        }
      })

      expect(wrapper.html()).toContain('form-row')
    })
  })
})