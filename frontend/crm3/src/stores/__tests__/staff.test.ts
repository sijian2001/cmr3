import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useStaffStore } from '../staff'

// Fetch のモック
const mockFetch = vi.fn()
global.fetch = mockFetch

// LocalStorage のモック
const localStorageMock = (() => {
  let store: { [key: string]: string } = {}
  return {
    getItem: (key: string) => store[key] || null,
    setItem: (key: string, value: string) => { store[key] = value.toString() },
    clear: () => { store = {} }
  }
})()
Object.defineProperty(window, 'localStorage', { value: localStorageMock })

describe('Staff Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    localStorageMock.setItem('token', 'test-token')
  })

  describe('State Management', () => {
    it('should initialize with default state', () => {
      const store = useStaffStore()

      expect(store.staff).toEqual([])
      expect(store.currentStaff).toBeNull()
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
      expect(store.total).toBe(0)
      expect(store.currentPage).toBe(1)
      expect(store.itemsPerPage).toBe(10)
    })
  })

  describe('Computed Properties', () => {
    it('should calculate totalPages correctly', () => {
      const store = useStaffStore()

      store.total = 25
      store.itemsPerPage = 10
      expect(store.totalPages).toBe(3)

      store.total = 30
      store.itemsPerPage = 10
      expect(store.totalPages).toBe(3)
    })

    it('should filter staff by status correctly', () => {
      const store = useStaffStore()

      store.staff = [
        { id: 1, name: '田中太郎', email: 'tanaka@example.com', position: '店長', status: 'active', hire_date: '2023-01-01', created_at: '2023-01-01T00:00:00Z', updated_at: '2023-01-01T00:00:00Z' },
        { id: 2, name: '佐藤花子', email: 'sato@example.com', position: '正社員', status: 'inactive', hire_date: '2023-02-01', created_at: '2023-02-01T00:00:00Z', updated_at: '2023-02-01T00:00:00Z' },
        { id: 3, name: '鈴木次郎', email: 'suzuki@example.com', position: 'アルバイト', status: 'on_leave', hire_date: '2023-03-01', created_at: '2023-03-01T00:00:00Z', updated_at: '2023-03-01T00:00:00Z' }
      ]

      expect(store.activeStaff).toHaveLength(1)
      expect(store.activeStaff[0].name).toBe('田中太郎')

      expect(store.inactiveStaff).toHaveLength(1)
      expect(store.inactiveStaff[0].name).toBe('佐藤花子')

      expect(store.onLeaveStaff).toHaveLength(1)
      expect(store.onLeaveStaff[0].name).toBe('鈴木次郎')
    })

    it('should calculate status counts correctly', () => {
      const store = useStaffStore()

      store.staff = [
        { id: 1, name: '田中太郎', email: 'tanaka@example.com', position: '店長', status: 'active', hire_date: '2023-01-01', created_at: '2023-01-01T00:00:00Z', updated_at: '2023-01-01T00:00:00Z' },
        { id: 2, name: '佐藤花子', email: 'sato@example.com', position: '正社員', status: 'active', hire_date: '2023-02-01', created_at: '2023-02-01T00:00:00Z', updated_at: '2023-02-01T00:00:00Z' },
        { id: 3, name: '鈴木次郎', email: 'suzuki@example.com', position: 'アルバイト', status: 'inactive', hire_date: '2023-03-01', created_at: '2023-03-01T00:00:00Z', updated_at: '2023-03-01T00:00:00Z' },
        { id: 4, name: '山田三郎', email: 'yamada@example.com', position: 'パート', status: 'on_leave', hire_date: '2023-04-01', created_at: '2023-04-01T00:00:00Z', updated_at: '2023-04-01T00:00:00Z' }
      ]

      expect(store.statusCounts.active).toBe(2)
      expect(store.statusCounts.inactive).toBe(1)
      expect(store.statusCounts.on_leave).toBe(1)
      expect(store.statusCounts.total).toBe(4)
    })

    it('should group staff by store correctly', () => {
      const store = useStaffStore()

      store.staff = [
        { id: 1, name: '田中太郎', email: 'tanaka@example.com', position: '店長', status: 'active', hire_date: '2023-01-01', created_at: '2023-01-01T00:00:00Z', updated_at: '2023-01-01T00:00:00Z', store: { id: 1, name: '新宿店' } },
        { id: 2, name: '佐藤花子', email: 'sato@example.com', position: '正社員', status: 'active', hire_date: '2023-02-01', created_at: '2023-02-01T00:00:00Z', updated_at: '2023-02-01T00:00:00Z', store: { id: 1, name: '新宿店' } },
        { id: 3, name: '鈴木次郎', email: 'suzuki@example.com', position: 'アルバイト', status: 'active', hire_date: '2023-03-01', created_at: '2023-03-01T00:00:00Z', updated_at: '2023-03-01T00:00:00Z' }
      ]

      expect(store.staffByStore['新宿店']).toHaveLength(2)
      expect(store.staffByStore['未割り当て']).toHaveLength(1)
    })
  })

  describe('fetchStaff', () => {
    it('should fetch staff successfully', async () => {
      const mockResponse = {
        staff: [
          { id: 1, name: '田中太郎', email: 'tanaka@example.com', position: '店長', status: 'active', hire_date: '2023-01-01', created_at: '2023-01-01T00:00:00Z', updated_at: '2023-01-01T00:00:00Z' }
        ],
        total: 1,
        page: 1,
        limit: 10,
        totalPages: 1
      }

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse
      })

      const store = useStaffStore()
      await store.fetchStaff()

      expect(store.staff).toEqual(mockResponse.staff)
      expect(store.total).toBe(1)
      expect(store.currentPage).toBe(1)
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
    })

    it('should handle search parameters', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ staff: [], total: 0, page: 1, limit: 10, totalPages: 0 })
      })

      const store = useStaffStore()
      await store.fetchStaff({ name: '田中', position: '店長', status: 'active' })

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringMatching(/name=%E7%94%B0%E4%B8%AD.*position=%E5%BA%97%E9%95%B7.*status=active/),
        expect.any(Object)
      )
    })

    it('should handle fetch error', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Failed to fetch staff' })
      })

      const store = useStaffStore()
      await store.fetchStaff()

      expect(store.error).toBe('Failed to fetch staff')
      expect(store.staff).toEqual([])
      expect(store.total).toBe(0)
    })
  })

  describe('fetchStaffMember', () => {
    it('should fetch single staff member successfully', async () => {
      const mockStaff = { id: 1, name: '田中太郎', email: 'tanaka@example.com', position: '店長', status: 'active', hire_date: '2023-01-01', created_at: '2023-01-01T00:00:00Z', updated_at: '2023-01-01T00:00:00Z' }

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockStaff
      })

      const store = useStaffStore()
      const result = await store.fetchStaffMember(1)

      expect(result).toEqual(mockStaff)
      expect(store.currentStaff).toEqual(mockStaff)
      expect(store.error).toBeNull()
    })

    it('should handle fetch single staff error', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Staff not found' })
      })

      const store = useStaffStore()

      await expect(store.fetchStaffMember(999)).rejects.toThrow('Staff not found')
      expect(store.currentStaff).toBeNull()
      expect(store.error).toBe('Staff not found')
    })
  })

  describe('createStaffMember', () => {
    it('should create staff member successfully', async () => {
      const newStaffData = {
        name: '新人太郎',
        email: 'newbie@example.com',
        position: 'アルバイト',
        status: 'active' as const,
        hire_date: '2023-06-01'
      }

      const createdStaff = {
        id: 4,
        ...newStaffData,
        created_at: '2023-06-01T00:00:00Z',
        updated_at: '2023-06-01T00:00:00Z'
      }

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => createdStaff
      })

      const store = useStaffStore()
      const result = await store.createStaffMember(newStaffData)

      expect(result).toEqual(createdStaff)
      expect(store.staff[0]).toEqual(createdStaff)
      expect(store.total).toBe(1)
    })

    it('should handle create staff error', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Email already exists' })
      })

      const store = useStaffStore()

      await expect(store.createStaffMember({
        name: '重複太郎',
        email: 'duplicate@example.com',
        position: 'アルバイト',
        status: 'active',
        hire_date: '2023-06-01'
      })).rejects.toThrow('Email already exists')
    })
  })

  describe('updateStaffMember', () => {
    it('should update staff member successfully', async () => {
      const store = useStaffStore()
      store.staff = [
        { id: 1, name: '田中太郎', email: 'tanaka@example.com', position: '店長', status: 'active', hire_date: '2023-01-01', created_at: '2023-01-01T00:00:00Z', updated_at: '2023-01-01T00:00:00Z' }
      ]

      const updateData = {
        name: '田中次郎',
        email: 'tanaka@example.com',
        position: '副店長',
        status: 'active' as const,
        hire_date: '2023-01-01'
      }

      const updatedStaff = {
        id: 1,
        ...updateData,
        created_at: '2023-01-01T00:00:00Z',
        updated_at: '2023-06-01T00:00:00Z'
      }

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => updatedStaff
      })

      const result = await store.updateStaffMember(1, updateData)

      expect(result).toEqual(updatedStaff)
      expect(store.staff[0]).toEqual(updatedStaff)
    })

    it('should update current staff if matches', async () => {
      const store = useStaffStore()
      store.currentStaff = { id: 1, name: '田中太郎', email: 'tanaka@example.com', position: '店長', status: 'active', hire_date: '2023-01-01', created_at: '2023-01-01T00:00:00Z', updated_at: '2023-01-01T00:00:00Z' }

      const updatedStaff = {
        id: 1,
        name: '田中次郎',
        email: 'tanaka@example.com',
        position: '副店長',
        status: 'active' as const,
        hire_date: '2023-01-01',
        created_at: '2023-01-01T00:00:00Z',
        updated_at: '2023-06-01T00:00:00Z'
      }

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => updatedStaff
      })

      await store.updateStaffMember(1, {
        name: '田中次郎',
        email: 'tanaka@example.com',
        position: '副店長',
        status: 'active',
        hire_date: '2023-01-01'
      })

      expect(store.currentStaff).toEqual(updatedStaff)
    })
  })

  describe('deleteStaffMember', () => {
    it('should delete staff member successfully', async () => {
      const store = useStaffStore()
      store.staff = [
        { id: 1, name: '田中太郎', email: 'tanaka@example.com', position: '店長', status: 'active', hire_date: '2023-01-01', created_at: '2023-01-01T00:00:00Z', updated_at: '2023-01-01T00:00:00Z' },
        { id: 2, name: '佐藤花子', email: 'sato@example.com', position: '正社員', status: 'active', hire_date: '2023-02-01', created_at: '2023-02-01T00:00:00Z', updated_at: '2023-02-01T00:00:00Z' }
      ]
      store.total = 2

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({})
      })

      await store.deleteStaffMember(1)

      expect(store.staff).toHaveLength(1)
      expect(store.staff[0].id).toBe(2)
      expect(store.total).toBe(1)
    })

    it('should clear current staff if deleted', async () => {
      const store = useStaffStore()
      store.currentStaff = { id: 1, name: '田中太郎', email: 'tanaka@example.com', position: '店長', status: 'active', hire_date: '2023-01-01', created_at: '2023-01-01T00:00:00Z', updated_at: '2023-01-01T00:00:00Z' }

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({})
      })

      await store.deleteStaffMember(1)

      expect(store.currentStaff).toBeNull()
    })
  })

  describe('updateStaffStatus', () => {
    it('should update staff status successfully', async () => {
      const store = useStaffStore()
      store.staff = [
        { id: 1, name: '田中太郎', email: 'tanaka@example.com', position: '店長', status: 'active', hire_date: '2023-01-01', created_at: '2023-01-01T00:00:00Z', updated_at: '2023-01-01T00:00:00Z' }
      ]

      const updatedStaff = {
        ...store.staff[0],
        status: 'inactive' as const,
        updated_at: '2023-06-01T00:00:00Z'
      }

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => updatedStaff
      })

      const result = await store.updateStaffStatus(1, 'inactive')

      expect(result.status).toBe('inactive')
      expect(store.staff[0].status).toBe('inactive')
    })
  })

  describe('assignToStore', () => {
    it('should assign staff to store successfully', async () => {
      const store = useStaffStore()
      store.staff = [
        { id: 1, name: '田中太郎', email: 'tanaka@example.com', position: '店長', status: 'active', hire_date: '2023-01-01', created_at: '2023-01-01T00:00:00Z', updated_at: '2023-01-01T00:00:00Z' }
      ]

      const updatedStaff = {
        ...store.staff[0],
        store_id: 1,
        store: { id: 1, name: '新宿店' },
        updated_at: '2023-06-01T00:00:00Z'
      }

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => updatedStaff
      })

      const result = await store.assignToStore(1, 1)

      expect(result.store_id).toBe(1)
      expect(result.store?.name).toBe('新宿店')
      expect(store.staff[0].store_id).toBe(1)
    })
  })

  describe('searchStaff', () => {
    it('should reset page and search staff', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ staff: [], total: 0, page: 1, limit: 10, totalPages: 0 })
      })

      const store = useStaffStore()
      store.currentPage = 3

      await store.searchStaff({ name: '田中' })

      expect(store.currentPage).toBe(1)
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringMatching(/name=%E7%94%B0%E4%B8%AD/),
        expect.any(Object)
      )
    })
  })

  describe('pagination methods', () => {
    it('should set page and fetch staff', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ staff: [], total: 0, page: 2, limit: 10, totalPages: 0 })
      })

      const store = useStaffStore()
      await store.setPage(2)

      expect(store.currentPage).toBe(2)
    })

    it('should set items per page and reset page', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ staff: [], total: 0, page: 1, limit: 20, totalPages: 0 })
      })

      const store = useStaffStore()
      store.currentPage = 3
      await store.setItemsPerPage(20)

      expect(store.itemsPerPage).toBe(20)
      expect(store.currentPage).toBe(1)
    })
  })

  describe('utility methods', () => {
    it('should clear error', () => {
      const store = useStaffStore()
      store.error = 'Some error'

      store.clearError()

      expect(store.error).toBeNull()
    })

    it('should clear current staff', () => {
      const store = useStaffStore()
      store.currentStaff = { id: 1, name: '田中太郎', email: 'tanaka@example.com', position: '店長', status: 'active', hire_date: '2023-01-01', created_at: '2023-01-01T00:00:00Z', updated_at: '2023-01-01T00:00:00Z' }

      store.clearCurrentStaff()

      expect(store.currentStaff).toBeNull()
    })

    it('should return correct status labels', () => {
      const store = useStaffStore()

      expect(store.getStaffStatusLabel('active')).toBe('在籍')
      expect(store.getStaffStatusLabel('inactive')).toBe('休職')
      expect(store.getStaffStatusLabel('on_leave')).toBe('休暇中')
    })

    it('should return correct status colors', () => {
      const store = useStaffStore()

      expect(store.getStaffStatusColor('active')).toBe('success')
      expect(store.getStaffStatusColor('inactive')).toBe('warning')
      expect(store.getStaffStatusColor('on_leave')).toBe('info')
    })
  })

  describe('error handling', () => {
    it('should handle network errors', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network error'))

      const store = useStaffStore()
      await store.fetchStaff()

      expect(store.error).toBe('Network error')
      expect(store.loading).toBe(false)
    })

    it('should handle JSON parsing errors', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        json: async () => { throw new Error('Invalid JSON') }
      })

      const store = useStaffStore()
      await store.fetchStaff()

      expect(store.error).toBe('Invalid JSON')
      expect(store.loading).toBe(false)
    })
  })

  describe('authentication', () => {
    it('should include authorization header in requests', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ staff: [], total: 0, page: 1, limit: 10, totalPages: 0 })
      })

      const store = useStaffStore()
      await store.fetchStaff()

      expect(mockFetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            'Authorization': 'Bearer test-token'
          })
        })
      )
    })
  })
})