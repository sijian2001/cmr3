import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import StaffForm from '../StaffForm.vue'

describe('StaffForm.vue', () => {
  let wrapper: any

  beforeEach(() => {
    wrapper = mount(StaffForm, {
      props: {
        loading: false
      }
    })
  })

  afterEach(() => {
    wrapper?.unmount()
  })

  describe('Component Rendering', () => {
    it('should render correctly', () => {
      expect(wrapper.find('form').exists()).toBe(true)
      expect(wrapper.find('.staff-form').exists()).toBe(true)
    })

    it('should render all form sections', () => {
      const sections = wrapper.findAll('.form-section')
      expect(sections).toHaveLength(2)

      expect(sections[0].find('.section-title').text()).toBe('基本情報')
      expect(sections[1].find('.section-title').text()).toBe('所属・ステータス設定')
    })

    it('should render all form fields', () => {
      expect(wrapper.find('#name').exists()).toBe(true)
      expect(wrapper.find('#email').exists()).toBe(true)
      expect(wrapper.find('#phone').exists()).toBe(true)
      expect(wrapper.find('#position').exists()).toBe(true)
      expect(wrapper.find('#hire_date').exists()).toBe(true)
      expect(wrapper.find('#store_id').exists()).toBe(true)
      expect(wrapper.find('#status').exists()).toBe(true)
    })

    it('should render required field indicators', () => {
      const requiredLabels = wrapper.findAll('label.required')
      expect(requiredLabels).toHaveLength(5) // name, email, position, hire_date, status
    })

    it('should render position options', () => {
      const positionSelect = wrapper.find('#position')
      const options = positionSelect.findAll('option')

      expect(options).toHaveLength(7) // Including placeholder
      expect(options[0].text()).toBe('役職を選択')
      expect(options[1].text()).toBe('店長')
      expect(options[2].text()).toBe('副店長')
      expect(options[3].text()).toBe('正社員')
      expect(options[4].text()).toBe('アルバイト')
      expect(options[5].text()).toBe('パート')
      expect(options[6].text()).toBe('契約社員')
    })

    it('should render status options', () => {
      const statusSelect = wrapper.find('#status')
      const options = statusSelect.findAll('option')

      expect(options).toHaveLength(3)
      expect(options[0].text()).toBe('在籍')
      expect(options[1].text()).toBe('休職')
      expect(options[2].text()).toBe('休暇中')
    })

    it('should render status description', () => {
      const statusInfo = wrapper.find('.status-info')
      expect(statusInfo.exists()).toBe(true)
      expect(statusInfo.text()).toContain('ステータスについて')
      expect(statusInfo.text()).toContain('在籍:')
      expect(statusInfo.text()).toContain('休職:')
      expect(statusInfo.text()).toContain('休暇中:')
    })

    it('should render form action buttons', () => {
      const buttons = wrapper.findAll('.form-actions .btn')
      expect(buttons).toHaveLength(2)
      expect(buttons[0].text()).toBe('キャンセル')
      expect(buttons[1].text()).toBe('作成')
    })
  })

  describe('Props Handling', () => {
    it('should show "更新" button text in edit mode', async () => {
      const mockStaff = {
        id: 1,
        name: '田中太郎',
        email: 'tanaka@example.com',
        phone: '090-1234-5678',
        position: '店長',
        store_id: 1,
        hire_date: '2023-01-01',
        status: 'active',
        created_at: '2023-01-01T00:00:00Z',
        updated_at: '2023-01-01T00:00:00Z'
      }

      wrapper = mount(StaffForm, {
        props: {
          staff: mockStaff,
          loading: false
        }
      })

      await wrapper.vm.$nextTick()

      const submitButton = wrapper.find('.btn--primary')
      expect(submitButton.text()).toBe('更新')
    })

    it('should show loading state', async () => {
      wrapper = mount(StaffForm, {
        props: {
          loading: true
        }
      })

      const submitButton = wrapper.find('.btn--primary')
      expect(submitButton.text()).toBe('保存中...')
      expect(submitButton.element.disabled).toBe(true)

      const cancelButton = wrapper.find('.btn--secondary')
      expect(cancelButton.element.disabled).toBe(true)
    })

    it('should populate form with staff data in edit mode', async () => {
      const mockStaff = {
        id: 1,
        name: '田中太郎',
        email: 'tanaka@example.com',
        phone: '090-1234-5678',
        position: '店長',
        store_id: 1,
        hire_date: '2023-01-01',
        status: 'active',
        created_at: '2023-01-01T00:00:00Z',
        updated_at: '2023-01-01T00:00:00Z'
      }

      wrapper = mount(StaffForm, {
        props: {
          staff: mockStaff,
          loading: false
        }
      })

      await wrapper.vm.$nextTick()

      expect(wrapper.find('#name').element.value).toBe('田中太郎')
      expect(wrapper.find('#email').element.value).toBe('tanaka@example.com')
      expect(wrapper.find('#phone').element.value).toBe('090-1234-5678')
      expect(wrapper.find('#position').element.value).toBe('店長')
      expect(wrapper.find('#store_id').element.value).toBe('1')
      expect(wrapper.find('#hire_date').element.value).toBe('2023-01-01')
      expect(wrapper.find('#status').element.value).toBe('active')
    })
  })

  describe('Form Validation', () => {
    it('should validate required fields', async () => {
      const submitButton = wrapper.find('.btn--primary')

      // Initially, form should be invalid due to empty required fields
      expect(submitButton.element.disabled).toBe(true)

      // Fill required fields
      await wrapper.find('#name').setValue('田中太郎')
      await wrapper.find('#email').setValue('tanaka@example.com')
      await wrapper.find('#position').setValue('店長')
      await wrapper.find('#hire_date').setValue('2023-01-01')

      await wrapper.vm.$nextTick()

      // Now form should be valid
      expect(submitButton.element.disabled).toBe(false)
    })

    it('should have validation methods', () => {
      // Check if validation methods exist on the component
      expect(typeof wrapper.vm.validateName).toBe('function')
      expect(typeof wrapper.vm.validateEmail).toBe('function')
      expect(typeof wrapper.vm.validatePhone).toBe('function')
      expect(typeof wrapper.vm.validatePosition).toBe('function')
      expect(typeof wrapper.vm.validateHireDate).toBe('function')
    })

    it('should validate form before submission', async () => {
      // Try to submit empty form
      await wrapper.find('form').trigger('submit.prevent')

      // Should not emit submit event for invalid form
      expect(wrapper.emitted('submit')).toBeFalsy()
    })
  })

  describe('Form Submission', () => {
    it('should emit submit event with form data', async () => {
      // Fill form with valid data
      await wrapper.find('#name').setValue('田中太郎')
      await wrapper.find('#email').setValue('tanaka@example.com')
      await wrapper.find('#phone').setValue('090-1234-5678')
      await wrapper.find('#position').setValue('店長')
      await wrapper.find('#store_id').setValue('1')
      await wrapper.find('#hire_date').setValue('2023-01-01')
      await wrapper.find('#status').setValue('active')

      // Submit form
      await wrapper.find('form').trigger('submit.prevent')

      const emitted = wrapper.emitted('submit')
      expect(emitted).toHaveLength(1)

      const submittedData = emitted[0][0]
      expect(submittedData).toEqual({
        name: '田中太郎',
        email: 'tanaka@example.com',
        phone: '090-1234-5678',
        position: '店長',
        store_id: 1,
        hire_date: '2023-01-01',
        status: 'active'
      })
    })

    it('should not submit if form is invalid', async () => {
      // Submit form with empty required fields
      await wrapper.find('form').trigger('submit.prevent')

      // Should not emit submit event
      expect(wrapper.emitted('submit')).toBeFalsy()
    })

    it('should trim text inputs', async () => {
      await wrapper.find('#name').setValue('  田中太郎  ')
      await wrapper.find('#email').setValue('  tanaka@example.com  ')
      await wrapper.find('#phone').setValue('  090-1234-5678  ')
      await wrapper.find('#position').setValue('店長')
      await wrapper.find('#hire_date').setValue('2023-01-01')

      await wrapper.find('form').trigger('submit.prevent')

      const emitted = wrapper.emitted('submit')
      const submittedData = emitted[0][0]

      expect(submittedData.name).toBe('田中太郎')
      expect(submittedData.email).toBe('tanaka@example.com')
      expect(submittedData.phone).toBe('090-1234-5678')
    })

    it('should handle undefined optional fields', async () => {
      await wrapper.find('#name').setValue('田中太郎')
      await wrapper.find('#email').setValue('tanaka@example.com')
      await wrapper.find('#position').setValue('店長')
      await wrapper.find('#hire_date').setValue('2023-01-01')
      // Leave phone and store_id empty

      await wrapper.find('form').trigger('submit.prevent')

      const emitted = wrapper.emitted('submit')
      const submittedData = emitted[0][0]

      expect(submittedData.phone).toBeUndefined()
      expect(submittedData.store_id).toBeUndefined()
    })
  })

  describe('Form Actions', () => {
    it('should emit cancel event', async () => {
      const cancelButton = wrapper.find('.btn--secondary')
      await cancelButton.trigger('click')

      expect(wrapper.emitted('cancel')).toHaveLength(1)
    })
  })

  describe('Initialization', () => {
    it('should set default hire date to today in create mode', () => {
      const today = new Date().toISOString().split('T')[0]
      expect(wrapper.find('#hire_date').element.value).toBe(today)
    })

    it('should not set default hire date in edit mode', async () => {
      const mockStaff = {
        id: 1,
        name: '田中太郎',
        email: 'tanaka@example.com',
        position: '店長',
        hire_date: '2023-01-01',
        status: 'active',
        created_at: '2023-01-01T00:00:00Z',
        updated_at: '2023-01-01T00:00:00Z'
      }

      wrapper = mount(StaffForm, {
        props: {
          staff: mockStaff,
          loading: false
        }
      })

      await wrapper.vm.$nextTick()

      expect(wrapper.find('#hire_date').element.value).toBe('2023-01-01')
    })
  })

  describe('Reactive Form State', () => {
    it('should track form validity state', async () => {
      const submitButton = wrapper.find('.btn--primary')

      // Initially invalid
      expect(submitButton.element.disabled).toBe(true)

      // Fill all required fields
      await wrapper.find('#name').setValue('田中太郎')
      await wrapper.find('#email').setValue('tanaka@example.com')
      await wrapper.find('#position').setValue('店長')
      await wrapper.find('#hire_date').setValue('2023-01-01')

      await wrapper.vm.$nextTick()

      // Should be valid now
      expect(wrapper.vm.isValid).toBe(true)
    })

    it('should have reactive form data', async () => {
      await wrapper.find('#name').setValue('田中太郎')

      expect(wrapper.vm.formData.name).toBe('田中太郎')
    })
  })

  describe('Store Selection', () => {
    it('should render store options', () => {
      const storeSelect = wrapper.find('#store_id')
      const options = storeSelect.findAll('option')

      expect(options).toHaveLength(4) // Including placeholder
      expect(options[0].text()).toBe('店舗を選択（未割り当て）')
      expect(options[1].text()).toBe('新宿店')
      expect(options[2].text()).toBe('渋谷店')
      expect(options[3].text()).toBe('池袋店')
    })

    it('should handle numeric store_id selection', async () => {
      const storeSelect = wrapper.find('#store_id')
      await storeSelect.setValue('2')

      await wrapper.find('#name').setValue('田中太郎')
      await wrapper.find('#email').setValue('tanaka@example.com')
      await wrapper.find('#position').setValue('店長')
      await wrapper.find('#hire_date').setValue('2023-01-01')

      await wrapper.find('form').trigger('submit.prevent')

      const emitted = wrapper.emitted('submit')
      const submittedData = emitted[0][0]

      expect(submittedData.store_id).toBe(2)
    })
  })

  describe('Accessibility', () => {
    it('should have proper labels for form fields', () => {
      expect(wrapper.find('label[for="name"]').text()).toContain('氏名')
      expect(wrapper.find('label[for="email"]').text()).toContain('メールアドレス')
      expect(wrapper.find('label[for="phone"]').text()).toContain('電話番号')
      expect(wrapper.find('label[for="position"]').text()).toContain('役職')
      expect(wrapper.find('label[for="hire_date"]').text()).toContain('入社日')
      expect(wrapper.find('label[for="store_id"]').text()).toContain('所属店舗')
      expect(wrapper.find('label[for="status"]').text()).toContain('ステータス')
    })

    it('should have proper input types', () => {
      expect(wrapper.find('#name').attributes('type')).toBe('text')
      expect(wrapper.find('#email').attributes('type')).toBe('email')
      expect(wrapper.find('#phone').attributes('type')).toBe('tel')
      expect(wrapper.find('#hire_date').attributes('type')).toBe('date')
    })

    it('should have placeholder texts', () => {
      expect(wrapper.find('#name').attributes('placeholder')).toBe('氏名を入力')
      expect(wrapper.find('#email').attributes('placeholder')).toBe('例: staff@example.com')
      expect(wrapper.find('#phone').attributes('placeholder')).toBe('例: 090-1234-5678')
    })
  })

  describe('Error Display', () => {
    it('should have error tracking mechanism', () => {
      // Check if errors reactive object exists
      expect(wrapper.vm.errors).toBeDefined()
      expect(typeof wrapper.vm.errors).toBe('object')
    })

    it('should have error validation functions', () => {
      // Test validation functions exist and work
      const nameError = wrapper.vm.validateName('')
      expect(nameError).toBe('氏名は必須です')

      const emailError = wrapper.vm.validateEmail('')
      expect(emailError).toBe('メールアドレスは必須です')

      const positionError = wrapper.vm.validatePosition('')
      expect(positionError).toBe('役職は必須です')
    })
  })
})