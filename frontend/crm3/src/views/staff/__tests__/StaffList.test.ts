import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import StaffList from '../StaffList.vue'
import { useStaffStore } from '../../../stores/staff'

// Modal components mock
vi.mock('../../../components/Modal.vue', () => ({
  default: {
    name: 'Modal',
    template: '<div class="modal"><slot /></div>',
    props: ['modelValue'],
    emits: ['update:modelValue']
  }
}))

// StaffForm component mock
vi.mock('../../../components/StaffForm.vue', () => ({
  default: {
    name: 'StaffForm',
    template: '<div class="staff-form">Staff Form</div>',
    props: ['staff', 'loading'],
    emits: ['submit', 'cancel']
  }
}))

// Router mock
const mockRouter = {
  push: vi.fn()
}

vi.mock('vue-router', () => ({
  useRouter: () => mockRouter
}))

describe('StaffList.vue', () => {
  let wrapper: any
  let store: any

  const mockStaff = [
    {
      id: 1,
      name: '田中太郎',
      email: 'tanaka@example.com',
      phone: '090-1234-5678',
      position: '店長',
      status: 'active',
      hire_date: '2023-01-01',
      created_at: '2023-01-01T00:00:00Z',
      updated_at: '2023-01-01T00:00:00Z',
      store: { id: 1, name: '新宿店' }
    },
    {
      id: 2,
      name: '佐藤花子',
      email: 'sato@example.com',
      position: '正社員',
      status: 'inactive',
      hire_date: '2023-02-01',
      created_at: '2023-02-01T00:00:00Z',
      updated_at: '2023-02-01T00:00:00Z'
    },
    {
      id: 3,
      name: '鈴木次郎',
      email: 'suzuki@example.com',
      position: 'アルバイト',
      status: 'on_leave',
      hire_date: '2023-03-01',
      created_at: '2023-03-01T00:00:00Z',
      updated_at: '2023-03-01T00:00:00Z',
      store: { id: 2, name: '渋谷店' }
    }
  ]

  beforeEach(() => {
    const pinia = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        staff: {
          staff: mockStaff,
          total: mockStaff.length,
          currentPage: 1,
          itemsPerPage: 10,
          loading: false,
          error: null
        }
      }
    })

    wrapper = mount(StaffList, {
      global: {
        plugins: [pinia]
      }
    })

    store = useStaffStore()
    store.fetchStaff = vi.fn()
    store.deleteStaffMember = vi.fn()
    store.updateStaffStatus = vi.fn()
    store.createStaffMember = vi.fn()
    store.updateStaffMember = vi.fn()
    store.setPage = vi.fn()
    store.setItemsPerPage = vi.fn()
    store.searchStaff = vi.fn()
  })

  describe('Component Rendering', () => {
    it('should render correctly', () => {
      expect(wrapper.find('h1').text()).toBe('スタッフ管理')
      expect(wrapper.find('.staff-list').exists()).toBe(true)
    })

    it('should display staff statistics', () => {
      const statsCards = wrapper.findAll('.stats-card')
      expect(statsCards).toHaveLength(4)

      expect(statsCards[0].find('.stat-value').text()).toBe('3')
      expect(statsCards[0].find('.stat-label').text()).toBe('総スタッフ数')

      expect(statsCards[1].find('.stat-value').text()).toBe('1')
      expect(statsCards[1].find('.stat-label').text()).toBe('在籍')

      expect(statsCards[2].find('.stat-value').text()).toBe('1')
      expect(statsCards[2].find('.stat-label').text()).toBe('休職')

      expect(statsCards[3].find('.stat-value').text()).toBe('1')
      expect(statsCards[3].find('.stat-label').text()).toBe('休暇中')
    })

    it('should display staff table', () => {
      const table = wrapper.find('.staff-table')
      expect(table.exists()).toBe(true)

      const rows = wrapper.findAll('tbody tr')
      expect(rows).toHaveLength(3)

      // Check first staff row
      expect(rows[0].find('td:nth-child(2)').text()).toBe('田中太郎')
      expect(rows[0].find('td:nth-child(3)').text()).toBe('tanaka@example.com')
      expect(rows[0].find('td:nth-child(5)').text()).toBe('店長')
      expect(rows[0].find('td:nth-child(6)').text()).toBe('新宿店')
      expect(rows[0].find('.status-badge.success').text()).toBe('在籍')
    })

    it('should show unassigned store label', () => {
      const rows = wrapper.findAll('tbody tr')
      expect(rows[1].find('td:nth-child(6)').text()).toBe('未割り当て')
    })
  })

  describe('Search Functionality', () => {
    it('should handle search input', async () => {
      const searchInput = wrapper.find('input[placeholder="スタッフ名で検索..."]')

      await searchInput.setValue('田中')

      // Wait for debounce
      await new Promise(resolve => setTimeout(resolve, 350))

      expect(store.searchStaff).toHaveBeenCalledWith({ name: '田中' })
    })

    it('should handle advanced search', async () => {
      // Open advanced search
      const advancedToggle = wrapper.find('.advanced-search-toggle')
      await advancedToggle.trigger('click')

      const emailInput = wrapper.find('input[placeholder="メールアドレスで検索..."]')
      const positionSelect = wrapper.find('select:nth-child(1)')
      const statusSelect = wrapper.find('select:nth-child(2)')

      await emailInput.setValue('tanaka@example.com')
      await positionSelect.setValue('店長')
      await statusSelect.setValue('active')

      const searchButton = wrapper.find('.btn--primary')
      await searchButton.trigger('click')

      expect(store.searchStaff).toHaveBeenCalledWith({
        name: '',
        email: 'tanaka@example.com',
        position: '店長',
        status: 'active'
      })
    })

    it('should clear search filters', async () => {
      // Open advanced search
      const advancedToggle = wrapper.find('.advanced-search-toggle')
      await advancedToggle.trigger('click')

      const clearButton = wrapper.find('.btn--secondary')
      await clearButton.trigger('click')

      expect(store.fetchStaff).toHaveBeenCalled()
    })
  })

  describe('CRUD Operations', () => {
    it('should open create modal', async () => {
      const createButton = wrapper.find('.btn--primary')
      await createButton.trigger('click')

      expect(wrapper.vm.showCreateModal).toBe(true)
    })

    it('should handle staff creation', async () => {
      const staffData = {
        name: '新人太郎',
        email: 'newbie@example.com',
        position: 'アルバイト',
        status: 'active',
        hire_date: '2023-06-01'
      }

      store.createStaffMember.mockResolvedValueOnce({ id: 4, ...staffData })

      await wrapper.vm.handleCreateStaff(staffData)

      expect(store.createStaffMember).toHaveBeenCalledWith(staffData)
      expect(wrapper.vm.showCreateModal).toBe(false)
    })

    it('should handle staff creation error', async () => {
      store.createStaffMember.mockRejectedValueOnce(new Error('Email already exists'))
      console.error = vi.fn() // Mock console.error

      await wrapper.vm.handleCreateStaff({
        name: '重複太郎',
        email: 'duplicate@example.com',
        position: 'アルバイト',
        status: 'active',
        hire_date: '2023-06-01'
      })

      expect(console.error).toHaveBeenCalled()
    })

    it('should open edit modal', async () => {
      const editButton = wrapper.find('.btn--edit')
      await editButton.trigger('click')

      expect(wrapper.vm.showEditModal).toBe(true)
      expect(wrapper.vm.editingStaff).toEqual(mockStaff[0])
    })

    it('should handle staff update', async () => {
      const updatedData = {
        name: '田中次郎',
        email: 'tanaka@example.com',
        position: '副店長',
        status: 'active',
        hire_date: '2023-01-01'
      }

      wrapper.vm.editingStaff = mockStaff[0]
      store.updateStaffMember.mockResolvedValueOnce({ id: 1, ...updatedData })

      await wrapper.vm.handleUpdateStaff(updatedData)

      expect(store.updateStaffMember).toHaveBeenCalledWith(1, updatedData)
      expect(wrapper.vm.showEditModal).toBe(false)
      expect(wrapper.vm.editingStaff).toBeNull()
    })

    it('should handle staff deletion with confirmation', async () => {
      window.confirm = vi.fn().mockReturnValue(true)
      store.deleteStaffMember.mockResolvedValueOnce(undefined)

      const deleteButton = wrapper.find('.btn--danger')
      await deleteButton.trigger('click')

      expect(window.confirm).toHaveBeenCalledWith('本当にこのスタッフを削除しますか？')
      expect(store.deleteStaffMember).toHaveBeenCalledWith(1)
    })

    it('should cancel deletion when not confirmed', async () => {
      window.confirm = vi.fn().mockReturnValue(false)

      const deleteButton = wrapper.find('.btn--danger')
      await deleteButton.trigger('click')

      expect(store.deleteStaffMember).not.toHaveBeenCalled()
    })
  })

  describe('Status Management', () => {
    it('should handle status change', async () => {
      store.updateStaffStatus.mockResolvedValueOnce({ ...mockStaff[0], status: 'inactive' })

      const statusButton = wrapper.find('.status-action-btn')
      await statusButton.trigger('click')

      expect(store.updateStaffStatus).toHaveBeenCalledWith(1, 'inactive')
    })

    it('should handle status change error', async () => {
      store.updateStaffStatus.mockRejectedValueOnce(new Error('Status change failed'))
      console.error = vi.fn()

      await wrapper.vm.handleStatusChange(mockStaff[0], 'inactive')

      expect(console.error).toHaveBeenCalled()
    })

    it('should display correct status action buttons', () => {
      const rows = wrapper.findAll('tbody tr')

      // Active staff should show "休職にする" and "休暇にする" buttons
      const activeStaffButtons = rows[0].findAll('.status-action-btn')
      expect(activeStaffButtons[0].text()).toBe('休職にする')
      expect(activeStaffButtons[1].text()).toBe('休暇にする')

      // Inactive staff should show "復職させる" button
      const inactiveStaffButtons = rows[1].findAll('.status-action-btn')
      expect(inactiveStaffButtons[0].text()).toBe('復職させる')

      // On leave staff should show "復職させる" button
      const onLeaveStaffButtons = rows[2].findAll('.status-action-btn')
      expect(onLeaveStaffButtons[0].text()).toBe('復職させる')
    })
  })

  describe('Pagination', () => {
    it('should handle page change', async () => {
      const pageButton = wrapper.find('[data-testid="page-2"]')
      await pageButton.trigger('click')

      expect(store.setPage).toHaveBeenCalledWith(2)
    })

    it('should handle items per page change', async () => {
      const itemsSelect = wrapper.find('.items-per-page select')
      await itemsSelect.setValue('20')

      expect(store.setItemsPerPage).toHaveBeenCalledWith(20)
    })

    it('should display pagination info', () => {
      const paginationInfo = wrapper.find('.pagination-info')
      expect(paginationInfo.text()).toContain('1-3 / 3件')
    })
  })

  describe('Loading and Error States', () => {
    it('should show loading state', async () => {
      store.loading = true
      await wrapper.vm.$nextTick()

      expect(wrapper.find('.loading-spinner').exists()).toBe(true)
    })

    it('should show error state', async () => {
      store.error = 'Failed to load staff'
      await wrapper.vm.$nextTick()

      expect(wrapper.find('.error-message').exists()).toBe(true)
      expect(wrapper.find('.error-message').text()).toContain('Failed to load staff')
    })

    it('should show empty state', async () => {
      store.staff = []
      store.total = 0
      await wrapper.vm.$nextTick()

      expect(wrapper.find('.empty-state').exists()).toBe(true)
      expect(wrapper.find('.empty-state').text()).toContain('スタッフが見つかりませんでした')
    })
  })

  describe('Lifecycle', () => {
    it('should fetch staff on mount', () => {
      expect(store.fetchStaff).toHaveBeenCalled()
    })
  })

  describe('Reactive Data', () => {
    it('should update when store data changes', async () => {
      store.staff = [
        {
          id: 4,
          name: '新人太郎',
          email: 'newbie@example.com',
          position: 'アルバイト',
          status: 'active',
          hire_date: '2023-06-01',
          created_at: '2023-06-01T00:00:00Z',
          updated_at: '2023-06-01T00:00:00Z'
        }
      ]
      store.total = 1

      await wrapper.vm.$nextTick()

      const rows = wrapper.findAll('tbody tr')
      expect(rows).toHaveLength(1)
      expect(rows[0].find('td:nth-child(2)').text()).toBe('新人太郎')
    })
  })

  describe('Modal Management', () => {
    it('should close modals on cancel', async () => {
      wrapper.vm.showCreateModal = true
      await wrapper.vm.handleCreateCancel()

      expect(wrapper.vm.showCreateModal).toBe(false)
    })

    it('should close edit modal and clear editing staff', async () => {
      wrapper.vm.showEditModal = true
      wrapper.vm.editingStaff = mockStaff[0]

      await wrapper.vm.handleEditCancel()

      expect(wrapper.vm.showEditModal).toBe(false)
      expect(wrapper.vm.editingStaff).toBeNull()
    })
  })

  describe('Computed Properties', () => {
    it('should compute filtered items count correctly', () => {
      expect(wrapper.vm.filteredItemsCount).toBe(3)
    })

    it('should compute start and end item numbers correctly', () => {
      expect(wrapper.vm.startItem).toBe(1)
      expect(wrapper.vm.endItem).toBe(3)
    })
  })

  describe('Table Sorting', () => {
    it('should display sortable table headers', () => {
      const headers = wrapper.findAll('th')
      expect(headers[1].text()).toBe('氏名')
      expect(headers[2].text()).toBe('メールアドレス')
      expect(headers[3].text()).toBe('電話番号')
      expect(headers[4].text()).toBe('役職')
      expect(headers[5].text()).toBe('所属店舗')
      expect(headers[6].text()).toBe('ステータス')
      expect(headers[7].text()).toBe('入社日')
    })
  })

  describe('Accessibility', () => {
    it('should have proper ARIA labels', () => {
      const createButton = wrapper.find('.btn--primary')
      expect(createButton.attributes('aria-label')).toBe('新しいスタッフを追加')

      const searchInput = wrapper.find('input[placeholder="スタッフ名で検索..."]')
      expect(searchInput.attributes('aria-label')).toBe('スタッフ名で検索')
    })

    it('should have proper table structure', () => {
      const table = wrapper.find('table')
      expect(table.exists()).toBe(true)

      const thead = wrapper.find('thead')
      expect(thead.exists()).toBe(true)

      const tbody = wrapper.find('tbody')
      expect(tbody.exists()).toBe(true)
    })
  })
})