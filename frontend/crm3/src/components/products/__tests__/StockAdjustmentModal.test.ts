import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import StockAdjustmentModal from '../StockAdjustmentModal.vue'
import type { Product } from '@/stores/product'

// Pinia セットアップ
beforeEach(() => {
  setActivePinia(createPinia())
})

afterEach(() => {
  vi.clearAllMocks()
})

describe('StockAdjustmentModal', () => {
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

  describe('コンポーネントの表示', () => {
    it('製品情報が正しく表示される', () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      expect(wrapper.text()).toContain('テスト製品')
      expect(wrapper.text()).toContain('SKU: TEST-001')
      expect(wrapper.text()).toContain('現在の在庫数量: 50')
    })

    it('操作種別のラジオボタンが表示される', () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      const addRadio = wrapper.find('input[value="add"]')
      const subtractRadio = wrapper.find('input[value="subtract"]')

      expect(addRadio.exists()).toBe(true)
      expect(subtractRadio.exists()).toBe(true)
      expect(wrapper.text()).toContain('入庫 (在庫を増やす)')
      expect(wrapper.text()).toContain('出庫 (在庫を減らす)')
    })

    it('数量入力フィールドが表示される', () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      const quantityInput = wrapper.find('#quantity')
      expect(quantityInput.exists()).toBe(true)
      expect(quantityInput.attributes('type')).toBe('number')
      expect(quantityInput.attributes('min')).toBe('1')
      expect(quantityInput.attributes('step')).toBe('1')
    })

    it('理由入力欄が表示される', () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      const reasonTextarea = wrapper.find('#reason')
      expect(reasonTextarea.exists()).toBe(true)
      expect(reasonTextarea.attributes('rows')).toBe('3')
    })
  })

  describe('在庫予想表示', () => {
    it('入庫操作の予想在庫が正しく表示される', async () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      const addRadio = wrapper.find('input[value="add"]')
      const quantityInput = wrapper.find('#quantity')

      await addRadio.setValue(true)
      await quantityInput.setValue('10')

      expect(wrapper.text()).toContain('50')
      expect(wrapper.text()).toContain('+')
      expect(wrapper.text()).toContain('10')
      expect(wrapper.text()).toContain('=60') // スペースなしで検索
    })

    it('出庫操作の予想在庫が正しく表示される', async () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      const subtractRadio = wrapper.find('input[value="subtract"]')
      const quantityInput = wrapper.find('#quantity')

      await subtractRadio.setValue(true)
      await quantityInput.setValue('20')

      expect(wrapper.text()).toContain('50')
      expect(wrapper.text()).toContain('-')
      expect(wrapper.text()).toContain('20')
      expect(wrapper.text()).toContain('=30') // スペースなしで検索
    })

    it('在庫が少なくなる場合に警告が表示される', async () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      const subtractRadio = wrapper.find('input[value="subtract"]')
      const quantityInput = wrapper.find('#quantity')

      await subtractRadio.setValue(true)
      await quantityInput.setValue('45') // 50 - 45 = 5 (50 * 0.2 = 10 未満)

      expect(wrapper.text()).toContain('⚠️ 在庫が少なくなります')
    })
  })

  describe('フォームバリデーション', () => {
    it('数量が0以下の場合にエラーが表示される', async () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      const quantityInput = wrapper.find('#quantity')
      await quantityInput.setValue('0')

      await wrapper.find('form').trigger('submit')

      expect(wrapper.text()).toContain('数量は1以上の整数を入力してください')
    })

    it('出庫数量が現在の在庫を超える場合にエラーが表示される', async () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      const subtractRadio = wrapper.find('input[value="subtract"]')
      const quantityInput = wrapper.find('#quantity')

      await subtractRadio.setValue(true)
      await quantityInput.setValue('100') // 在庫50より多い

      await wrapper.find('form').trigger('submit')

      expect(wrapper.text()).toContain('出庫数量は現在の在庫数量を超えることはできません')
    })

    it('有効な入力の場合はエラーが表示されない', async () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      const addRadio = wrapper.find('input[value="add"]')
      const quantityInput = wrapper.find('#quantity')

      await addRadio.setValue(true)
      await quantityInput.setValue('10')

      await wrapper.find('form').trigger('submit')

      expect(wrapper.text()).not.toContain('数量は1以上の整数を入力してください')
      expect(wrapper.text()).not.toContain('出庫数量は現在の在庫数量を超えることはできません')
    })
  })

  describe('在庫状況のクラス判定', () => {
    it('在庫が0の場合にstock-emptyクラスが適用される', () => {
      const zeroStockProduct: Product = { ...mockProduct, stock_quantity: 0 }
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: zeroStockProduct
        }
      })

      expect(wrapper.find('.stock-empty').exists()).toBe(true)
    })

    it('在庫が10以下の場合にstock-lowクラスが適用される', () => {
      const lowStockProduct: Product = { ...mockProduct, stock_quantity: 5 }
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: lowStockProduct
        }
      })

      expect(wrapper.find('.stock-low').exists()).toBe(true)
    })

    it('在庫が11以上の場合にstock-normalクラスが適用される', () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      expect(wrapper.find('.stock-normal').exists()).toBe(true)
    })
  })

  describe('イベント処理', () => {
    it('キャンセルボタンクリック時にcancelイベントが発火される', async () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      await wrapper.find('.btn-cancel').trigger('click')

      expect(wrapper.emitted('cancel')).toBeTruthy()
    })

    it('×ボタンクリック時にcancelイベントが発火される', async () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      await wrapper.find('.btn-close').trigger('click')

      expect(wrapper.emitted('cancel')).toBeTruthy()
    })

    it('正しいデータでフォーム送信時にadjustイベントが発火される', async () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      const addRadio = wrapper.find('input[value="add"]')
      const quantityInput = wrapper.find('#quantity')

      await addRadio.setValue(true)
      await quantityInput.setValue('10')

      await wrapper.find('form').trigger('submit')

      expect(wrapper.emitted('adjust')).toBeTruthy()
      const adjustEvent = wrapper.emitted('adjust')
      expect(adjustEvent?.[0]).toEqual([10, 'add'])
    })

    it('モーダル外クリック時にcancelイベントが発火される', async () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      const overlay = wrapper.find('.modal-overlay')
      await overlay.trigger('click')

      expect(wrapper.emitted('cancel')).toBeTruthy()
    })
  })

  describe('フォームの状態管理', () => {
    it('送信中は実行ボタンが無効になる', async () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      await wrapper.find('#quantity').setValue('10')

      // 送信中状態をシミュレート
      const vm = wrapper.vm as any
      vm.isSubmitting = true
      await wrapper.vm.$nextTick()

      const submitButton = wrapper.find('.btn-save')
      expect(submitButton.attributes('disabled')).toBeDefined()
      expect(submitButton.text()).toBe('調整中...')
    })

    it('無効なフォームの場合は実行ボタンが無効になる', async () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      await wrapper.find('#quantity').setValue('0') // 無効な値

      const submitButton = wrapper.find('.btn-save')
      expect(submitButton.attributes('disabled')).toBeDefined()
    })

    it('有効なフォームの場合は実行ボタンが有効になる', async () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      await wrapper.find('#quantity').setValue('10') // 有効な値

      const submitButton = wrapper.find('.btn-save')
      expect(submitButton.attributes('disabled')).toBeUndefined()
    })
  })

  describe('初期状態', () => {
    it('デフォルトで入庫操作が選択されている', () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      const addRadio = wrapper.find('input[value="add"]')
      expect((addRadio.element as HTMLInputElement).checked).toBe(true)
    })

    it('デフォルトで数量1が設定されている', () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      const quantityInput = wrapper.find('#quantity')
      expect((quantityInput.element as HTMLInputElement).value).toBe('1')
    })

    it('理由欄は空で開始される', () => {
      const wrapper = mount(StockAdjustmentModal, {
        props: {
          product: mockProduct
        }
      })

      const reasonTextarea = wrapper.find('#reason')
      expect((reasonTextarea.element as HTMLTextAreaElement).value).toBe('')
    })
  })
})