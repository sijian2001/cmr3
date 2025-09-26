import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import StoreForm from '../StoreForm.vue'
import type { Store } from '@/stores/store'

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

describe('StoreForm', () => {
  describe('新規作成モード', () => {
    let wrapper: any

    beforeEach(() => {
      wrapper = mount(StoreForm, {
        props: {
          loading: false
        }
      })
    })

    it('正しくレンダリングされる', () => {
      expect(wrapper.find('form').exists()).toBe(true)
      expect(wrapper.find('#name').exists()).toBe(true)
      expect(wrapper.find('#address').exists()).toBe(true)
      expect(wrapper.find('#phone').exists()).toBe(true)
      expect(wrapper.find('#email').exists()).toBe(true)
      expect(wrapper.find('#status').exists()).toBe(true)
      expect(wrapper.find('#manager_id').exists()).toBe(true)
    })

    it('デフォルト値が設定されている', () => {
      const nameInput = wrapper.find('#name')
      const statusSelect = wrapper.find('#status')

      expect(nameInput.element.value).toBe('')
      expect(statusSelect.element.value).toBe('active')
    })

    it('新規作成ボタンが表示される', () => {
      const submitButton = wrapper.find('button[type="submit"]')
      expect(submitButton.text()).toContain('作成')
    })

    it('必須フィールドが空の場合は送信できない', async () => {
      const submitButton = wrapper.find('button[type="submit"]')
      expect(submitButton.attributes('disabled')).toBeDefined()
    })

    it('有効なデータが入力されると送信できる', async () => {
      const nameInput = wrapper.find('#name')
      const statusSelect = wrapper.find('#status')

      await nameInput.setValue('テスト店舗')
      await statusSelect.setValue('active')

      const submitButton = wrapper.find('button[type="submit"]')
      expect(submitButton.attributes('disabled')).toBeUndefined()
    })

    it('正しいデータで送信される', async () => {
      const nameInput = wrapper.find('#name')
      const addressTextarea = wrapper.find('#address')
      const phoneInput = wrapper.find('#phone')
      const emailInput = wrapper.find('#email')
      const statusSelect = wrapper.find('#status')
      const managerIdInput = wrapper.find('#manager_id')

      await nameInput.setValue('新規店舗')
      await addressTextarea.setValue('神奈川県横浜市')
      await phoneInput.setValue('045-1234-5678')
      await emailInput.setValue('new@example.com')
      await statusSelect.setValue('active')
      await managerIdInput.setValue('2')

      const form = wrapper.find('form')
      await form.trigger('submit')

      expect(wrapper.emitted('submit')).toBeTruthy()
      const emittedData = wrapper.emitted('submit')[0][0]

      expect(emittedData).toEqual({
        name: '新規店舗',
        address: '神奈川県横浜市',
        phone: '045-1234-5678',
        email: 'new@example.com',
        status: 'active',
        manager_id: 2
      })
    })

    it('キャンセルボタンがクリックされるとcancelイベントが発生する', async () => {
      const cancelButton = wrapper.find('button[type="button"]')
      await cancelButton.trigger('click')

      expect(wrapper.emitted('cancel')).toBeTruthy()
    })
  })

  describe('編集モード', () => {
    let wrapper: any

    beforeEach(() => {
      wrapper = mount(StoreForm, {
        props: {
          store: mockStore,
          loading: false
        }
      })
    })

    it('編集モードで既存データが設定される', () => {
      const nameInput = wrapper.find('#name')
      const addressTextarea = wrapper.find('#address')
      const phoneInput = wrapper.find('#phone')
      const emailInput = wrapper.find('#email')
      const statusSelect = wrapper.find('#status')
      const managerIdInput = wrapper.find('#manager_id')

      expect(nameInput.element.value).toBe(mockStore.name)
      expect(addressTextarea.element.value).toBe(mockStore.address)
      expect(phoneInput.element.value).toBe(mockStore.phone)
      expect(emailInput.element.value).toBe(mockStore.email)
      expect(statusSelect.element.value).toBe(mockStore.status)
      expect(managerIdInput.element.value).toBe(mockStore.manager_id?.toString())
    })

    it('更新ボタンが表示される', () => {
      const submitButton = wrapper.find('button[type="submit"]')
      expect(submitButton.text()).toContain('更新')
    })

    it('編集されたデータで送信される', async () => {
      const nameInput = wrapper.find('#name')
      const statusSelect = wrapper.find('#status')

      await nameInput.setValue('更新された店舗')
      await statusSelect.setValue('inactive')

      const form = wrapper.find('form')
      await form.trigger('submit')

      expect(wrapper.emitted('submit')).toBeTruthy()
      const emittedData = wrapper.emitted('submit')[0][0]

      expect(emittedData.name).toBe('更新された店舗')
      expect(emittedData.status).toBe('inactive')
    })
  })

  describe('バリデーション', () => {
    let wrapper: any

    beforeEach(() => {
      wrapper = mount(StoreForm, {
        props: {
          loading: false
        }
      })
    })

    it('店舗名が空の場合はエラーが表示される', async () => {
      const nameInput = wrapper.find('#name')
      await nameInput.setValue('')
      await nameInput.trigger('blur')

      // バリデーションロジックは内部で行われるため、実際のエラー表示のテストは困難
      // バリデーション状態のテストは submit 時の disabled 状態で確認
      const submitButton = wrapper.find('button[type="submit"]')
      expect(submitButton.attributes('disabled')).toBeDefined()
    })

    it('無効なメールアドレスの場合はエラーが発生する', async () => {
      const nameInput = wrapper.find('#name')
      const emailInput = wrapper.find('#email')

      await nameInput.setValue('テスト店舗')
      await emailInput.setValue('invalid-email')

      const form = wrapper.find('form')
      await form.trigger('submit')

      // バリデーションエラーの場合、submitイベントは発生しない
      expect(wrapper.emitted('submit')).toBeFalsy()
    })

    it('無効な電話番号の場合はエラーが発生する', async () => {
      const nameInput = wrapper.find('#name')
      const phoneInput = wrapper.find('#phone')

      await nameInput.setValue('テスト店舗')
      await phoneInput.setValue('abc-def-ghij')

      const form = wrapper.find('form')
      await form.trigger('submit')

      // バリデーションエラーの場合、submitイベントは発生しない
      expect(wrapper.emitted('submit')).toBeFalsy()
    })

    it('有効なデータの場合は送信される', async () => {
      const nameInput = wrapper.find('#name')
      const emailInput = wrapper.find('#email')
      const phoneInput = wrapper.find('#phone')

      await nameInput.setValue('テスト店舗')
      await emailInput.setValue('test@example.com')
      await phoneInput.setValue('03-1234-5678')

      const form = wrapper.find('form')
      await form.trigger('submit')

      expect(wrapper.emitted('submit')).toBeTruthy()
    })
  })

  describe('ステータス表示', () => {
    let wrapper: any

    beforeEach(() => {
      wrapper = mount(StoreForm, {
        props: {
          loading: false
        }
      })
    })

    it('activeステータスの説明が正しく表示される', async () => {
      const statusSelect = wrapper.find('#status')
      await statusSelect.setValue('active')

      expect(wrapper.text()).toContain('通常営業状態です。顧客サービスが提供されます。')
    })

    it('inactiveステータスの説明が正しく表示される', async () => {
      const statusSelect = wrapper.find('#status')
      await statusSelect.setValue('inactive')

      expect(wrapper.text()).toContain('一時的に休業している状態です。')
    })

    it('maintenanceステータスの説明が正しく表示される', async () => {
      const statusSelect = wrapper.find('#status')
      await statusSelect.setValue('maintenance')

      expect(wrapper.text()).toContain('メンテナンス中です。営業を一時停止しています。')
    })
  })

  describe('ローディング状態', () => {
    it('ローディング中はボタンが無効になる', () => {
      const wrapper = mount(StoreForm, {
        props: {
          loading: true
        }
      })

      const submitButton = wrapper.find('button[type="submit"]')
      const cancelButton = wrapper.find('button[type="button"]')

      expect(submitButton.attributes('disabled')).toBeDefined()
      expect(cancelButton.attributes('disabled')).toBeDefined()
    })

    it('ローディング中は「保存中...」と表示される', () => {
      const wrapper = mount(StoreForm, {
        props: {
          loading: true
        }
      })

      const submitButton = wrapper.find('button[type="submit"]')
      expect(submitButton.text()).toContain('保存中...')
    })
  })

  describe('データのトリミング', () => {
    it('送信時に文字列フィールドがトリミングされる', async () => {
      const wrapper = mount(StoreForm, {
        props: {
          loading: false
        }
      })

      const nameInput = wrapper.find('#name')
      const addressTextarea = wrapper.find('#address')
      const phoneInput = wrapper.find('#phone')
      const emailInput = wrapper.find('#email')

      await nameInput.setValue('  テスト店舗  ')
      await addressTextarea.setValue('  東京都渋谷区  ')
      await phoneInput.setValue('  03-1234-5678  ')
      await emailInput.setValue('  test@example.com  ')

      const form = wrapper.find('form')
      await form.trigger('submit')

      expect(wrapper.emitted('submit')).toBeTruthy()
      const emittedData = wrapper.emitted('submit')[0][0]

      expect(emittedData.name).toBe('テスト店舗')
      expect(emittedData.address).toBe('東京都渋谷区')
      expect(emittedData.phone).toBe('03-1234-5678')
      expect(emittedData.email).toBe('test@example.com')
    })

    it('空文字列はundefinedとして送信される', async () => {
      const wrapper = mount(StoreForm, {
        props: {
          loading: false
        }
      })

      const nameInput = wrapper.find('#name')
      const addressTextarea = wrapper.find('#address')
      const phoneInput = wrapper.find('#phone')
      const emailInput = wrapper.find('#email')

      await nameInput.setValue('テスト店舗')
      await addressTextarea.setValue('')
      await phoneInput.setValue('')
      await emailInput.setValue('')

      const form = wrapper.find('form')
      await form.trigger('submit')

      expect(wrapper.emitted('submit')).toBeTruthy()
      const emittedData = wrapper.emitted('submit')[0][0]

      expect(emittedData.name).toBe('テスト店舗')
      expect(emittedData.address).toBeUndefined()
      expect(emittedData.phone).toBeUndefined()
      expect(emittedData.email).toBeUndefined()
    })
  })
})