import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { mount, VueWrapper } from '@vue/test-utils'
import CustomerFormModal from '../CustomerFormModal.vue'
import type { Customer } from '@/stores/customer'

describe('CustomerFormModal', () => {
  let wrapper: VueWrapper<any>

  const mockCustomer: Customer = {
    id: 1,
    name: '山田太郎',
    email: 'yamada@example.com',
    phone: '090-1234-5678',
    address: '東京都渋谷区1-1-1'
  }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    if (wrapper) {
      wrapper.unmount()
    }
  })

  describe('新規作成モード', () => {
    beforeEach(() => {
      wrapper = mount(CustomerFormModal, {
        props: {
          customer: null
        }
      })
    })

    it('新規作成のタイトルが表示される', () => {
      expect(wrapper.find('.modal-header h2').text()).toBe('新規顧客登録')
    })

    it('フォームが空の状態で表示される', () => {
      const nameInput = wrapper.find('#name')
      const emailInput = wrapper.find('#email')
      const phoneInput = wrapper.find('#phone')
      const addressInput = wrapper.find('#address')

      expect((nameInput.element as HTMLInputElement).value).toBe('')
      expect((emailInput.element as HTMLInputElement).value).toBe('')
      expect((phoneInput.element as HTMLInputElement).value).toBe('')
      expect((addressInput.element as HTMLTextAreaElement).value).toBe('')
    })

    it('保存ボタンのテキストが正しい', () => {
      const saveButton = wrapper.find('.btn-save')
      expect(saveButton.text()).toBe('保存')
    })
  })

  describe('編集モード', () => {
    beforeEach(() => {
      wrapper = mount(CustomerFormModal, {
        props: {
          customer: mockCustomer
        }
      })
    })

    it('編集のタイトルが表示される', () => {
      expect(wrapper.find('.modal-header h2').text()).toBe('顧客情報編集')
    })

    it('顧客データでフォームが初期化される', async () => {
      await wrapper.vm.$nextTick()

      const nameInput = wrapper.find('#name')
      const emailInput = wrapper.find('#email')
      const phoneInput = wrapper.find('#phone')
      const addressInput = wrapper.find('#address')

      expect((nameInput.element as HTMLInputElement).value).toBe(mockCustomer.name)
      expect((emailInput.element as HTMLInputElement).value).toBe(mockCustomer.email)
      expect((phoneInput.element as HTMLInputElement).value).toBe(mockCustomer.phone)
      expect((addressInput.element as HTMLTextAreaElement).value).toBe(mockCustomer.address)
    })
  })

  describe('バリデーション', () => {
    beforeEach(() => {
      wrapper = mount(CustomerFormModal, {
        props: {
          customer: null
        }
      })
    })

    it('必須項目が空の場合エラーが表示される', async () => {
      const form = wrapper.find('form')

      // 名前を空にして送信
      await wrapper.find('#name').setValue('')
      await wrapper.find('#email').setValue('valid@email.com')
      await form.trigger('submit')

      await wrapper.vm.$nextTick()

      const nameError = wrapper.find('.form-group:first-child .error-text')
      expect(nameError.exists()).toBe(true)
      expect(nameError.text()).toBe('名前は必須です')
    })

    it('メールアドレスが空の場合エラーが表示される', async () => {
      const form = wrapper.find('form')

      await wrapper.find('#name').setValue('テスト')
      await wrapper.find('#email').setValue('')
      await form.trigger('submit')

      await wrapper.vm.$nextTick()

      const emailErrors = wrapper.findAll('.error-text')
      const emailError = emailErrors.find(error => error.text() === 'メールアドレスは必須です')
      expect(emailError).toBeTruthy()
    })

    it('無効なメールアドレス形式の場合エラーが表示される', async () => {
      const form = wrapper.find('form')

      await wrapper.find('#name').setValue('テスト')
      await wrapper.find('#email').setValue('invalid-email')
      await form.trigger('submit')

      await wrapper.vm.$nextTick()

      const emailErrors = wrapper.findAll('.error-text')
      const emailError = emailErrors.find(error => error.text() === '有効なメールアドレスを入力してください')
      expect(emailError).toBeTruthy()
    })

    it('無効な電話番号形式の場合エラーが表示される', async () => {
      const form = wrapper.find('form')

      await wrapper.find('#name').setValue('テスト')
      await wrapper.find('#email').setValue('test@example.com')
      await wrapper.find('#phone').setValue('invalid-phone')
      await form.trigger('submit')

      await wrapper.vm.$nextTick()

      const phoneErrors = wrapper.findAll('.error-text')
      const phoneError = phoneErrors.find(error => error.text() === '有効な電話番号を入力してください')
      expect(phoneError).toBeTruthy()
    })

    it('有効な電話番号形式の場合エラーが表示されない', async () => {
      const form = wrapper.find('form')

      await wrapper.find('#name').setValue('テスト')
      await wrapper.find('#email').setValue('test@example.com')
      await wrapper.find('#phone').setValue('090-1234-5678')
      await form.trigger('submit')

      await wrapper.vm.$nextTick()

      const phoneErrors = wrapper.findAll('.error-text')
      const phoneError = phoneErrors.find(error => error.text().includes('電話番号'))
      expect(phoneError).toBeFalsy()
    })
  })

  describe('イベント発火', () => {
    beforeEach(() => {
      wrapper = mount(CustomerFormModal, {
        props: {
          customer: null
        }
      })
    })

    it('キャンセルボタンクリックでcancelイベントが発火される', async () => {
      const cancelButton = wrapper.find('.btn-cancel')
      await cancelButton.trigger('click')

      expect(wrapper.emitted('cancel')).toBeTruthy()
      expect(wrapper.emitted('cancel')).toHaveLength(1)
    })

    it('×ボタンクリックでcancelイベントが発火される', async () => {
      const closeButton = wrapper.find('.btn-close')
      await closeButton.trigger('click')

      expect(wrapper.emitted('cancel')).toBeTruthy()
      expect(wrapper.emitted('cancel')).toHaveLength(1)
    })

    it('モーダル外クリックでcancelイベントが発火される', async () => {
      const overlay = wrapper.find('.modal-overlay')
      await overlay.trigger('click')

      expect(wrapper.emitted('cancel')).toBeTruthy()
      expect(wrapper.emitted('cancel')).toHaveLength(1)
    })

    it('モーダル内容クリックではcancelイベントが発火されない', async () => {
      const modalContent = wrapper.find('.modal-content')
      await modalContent.trigger('click')

      expect(wrapper.emitted('cancel')).toBeFalsy()
    })

    it('有効なフォーム送信でsaveイベントが発火される', async () => {
      const form = wrapper.find('form')

      await wrapper.find('#name').setValue('テスト顧客')
      await wrapper.find('#email').setValue('test@example.com')
      await wrapper.find('#phone').setValue('090-1234-5678')
      await wrapper.find('#address').setValue('テスト住所')

      await form.trigger('submit')
      await wrapper.vm.$nextTick()

      expect(wrapper.emitted('save')).toBeTruthy()
      expect(wrapper.emitted('save')).toHaveLength(1)

      const saveEventData = wrapper.emitted('save')![0][0]
      expect(saveEventData).toEqual({
        name: 'テスト顧客',
        email: 'test@example.com',
        phone: '090-1234-5678',
        address: 'テスト住所'
      })
    })

    it('空の任意項目は undefined として送信される', async () => {
      const form = wrapper.find('form')

      await wrapper.find('#name').setValue('テスト顧客')
      await wrapper.find('#email').setValue('test@example.com')
      // phone と address は空のまま

      await form.trigger('submit')
      await wrapper.vm.$nextTick()

      const saveEventData = wrapper.emitted('save')![0][0]
      expect(saveEventData).toEqual({
        name: 'テスト顧客',
        email: 'test@example.com',
        phone: undefined,
        address: undefined
      })
    })
  })

  describe('UI状態', () => {
    beforeEach(() => {
      wrapper = mount(CustomerFormModal, {
        props: {
          customer: null
        }
      })
    })

    it('送信中はボタンが無効になる', async () => {
      const form = wrapper.find('form')
      const saveButton = wrapper.find('.btn-save')

      await wrapper.find('#name').setValue('テスト')
      await wrapper.find('#email').setValue('test@example.com')

      // isSubmitting を true に設定（実際の送信はモック）
      wrapper.vm.isSubmitting = true
      await wrapper.vm.$nextTick()

      expect(saveButton.attributes('disabled')).toBeDefined()
      expect(saveButton.text()).toBe('保存中...')
    })

    it('エラーがある入力フィールドにerrorクラスが適用される', async () => {
      const form = wrapper.find('form')

      // バリデーションエラーを発生させる
      await wrapper.find('#name').setValue('')
      await form.trigger('submit')
      await wrapper.vm.$nextTick()

      const nameInput = wrapper.find('#name')
      expect(nameInput.classes()).toContain('error')
    })
  })

  describe('アクセシビリティ', () => {
    beforeEach(() => {
      wrapper = mount(CustomerFormModal, {
        props: {
          customer: null
        }
      })
    })

    it('必須フィールドにrequired属性が設定される', () => {
      const nameInput = wrapper.find('#name')
      const emailInput = wrapper.find('#email')

      expect(nameInput.attributes('required')).toBeDefined()
      expect(emailInput.attributes('required')).toBeDefined()
    })

    it('ラベルとinputが正しく関連付けられている', () => {
      const nameLabel = wrapper.find('label[for="name"]')
      const nameInput = wrapper.find('#name')

      expect(nameLabel.exists()).toBe(true)
      expect(nameInput.attributes('id')).toBe('name')
    })

    it('エラーメッセージが適切な場所に表示される', async () => {
      const form = wrapper.find('form')

      await wrapper.find('#name').setValue('')
      await form.trigger('submit')
      await wrapper.vm.$nextTick()

      const nameFormGroup = wrapper.find('.form-group:first-child')
      const errorText = nameFormGroup.find('.error-text')

      expect(errorText.exists()).toBe(true)
      expect(errorText.text()).toBe('名前は必須です')
    })
  })
})