import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface Staff {
  id: number
  name: string
  email: string
  phone?: string
  position: string
  store_id?: number
  hire_date: string
  status: 'active' | 'inactive' | 'on_leave'
  created_at: string
  updated_at: string
  // 関連データ
  store?: {
    id: number
    name: string
  }
}

export interface StaffCreateRequest {
  name: string
  email: string
  phone?: string
  position: string
  store_id?: number
  hire_date: string
  status: 'active' | 'inactive' | 'on_leave'
}

export interface StaffUpdateRequest {
  name: string
  email: string
  phone?: string
  position: string
  store_id?: number
  hire_date: string
  status: 'active' | 'inactive' | 'on_leave'
}

export interface StaffSearchParams {
  name?: string
  email?: string
  position?: string
  store_id?: number
  status?: string
  page?: number
  limit?: number
}

export interface PaginatedStaffResponse {
  staff: Staff[]
  total: number
  page: number
  limit: number
  totalPages: number
}

export const useStaffStore = defineStore('staff', () => {
  // State
  const staff = ref<Staff[]>([])
  const currentStaff = ref<Staff | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const total = ref(0)
  const currentPage = ref(1)
  const itemsPerPage = ref(10)

  // Computed
  const totalPages = computed(() => Math.ceil(total.value / itemsPerPage.value))

  const activeStaff = computed(() =>
    staff.value.filter(s => s.status === 'active')
  )

  const inactiveStaff = computed(() =>
    staff.value.filter(s => s.status === 'inactive')
  )

  const onLeaveStaff = computed(() =>
    staff.value.filter(s => s.status === 'on_leave')
  )

  const statusCounts = computed(() => ({
    active: activeStaff.value.length,
    inactive: inactiveStaff.value.length,
    on_leave: onLeaveStaff.value.length,
    total: staff.value.length
  }))

  const staffByStore = computed(() => {
    const grouped: Record<string, Staff[]> = {}
    staff.value.forEach(s => {
      const storeKey = s.store?.name || '未割り当て'
      if (!grouped[storeKey]) {
        grouped[storeKey] = []
      }
      grouped[storeKey].push(s)
    })
    return grouped
  })

  // Actions
  const fetchStaff = async (params: StaffSearchParams = {}) => {
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

      const response = await fetch(`/api/staff?${query}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to fetch staff')
      }

      const data: PaginatedStaffResponse = await response.json()

      staff.value = data.staff || []
      total.value = data.total
      currentPage.value = data.page

    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch staff'
      staff.value = []
      total.value = 0
    } finally {
      loading.value = false
    }
  }

  const fetchStaffMember = async (id: number) => {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`/api/staff/${id}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to fetch staff member')
      }

      const data: Staff = await response.json()
      currentStaff.value = data
      return data

    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch staff member'
      currentStaff.value = null
      throw err
    } finally {
      loading.value = false
    }
  }

  const createStaffMember = async (staffData: StaffCreateRequest) => {
    loading.value = true
    error.value = null

    try {
      const response = await fetch('/api/staff', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(staffData),
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to create staff member')
      }

      const newStaff: Staff = await response.json()

      // 現在の一覧に追加
      staff.value.unshift(newStaff)
      total.value += 1

      return newStaff

    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create staff member'
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateStaffMember = async (id: number, staffData: StaffUpdateRequest) => {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`/api/staff/${id}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(staffData),
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to update staff member')
      }

      const updatedStaff: Staff = await response.json()

      // 一覧内のスタッフを更新
      const index = staff.value.findIndex(s => s.id === id)
      if (index !== -1) {
        staff.value[index] = updatedStaff
      }

      // 現在のスタッフも更新
      if (currentStaff.value && currentStaff.value.id === id) {
        currentStaff.value = updatedStaff
      }

      return updatedStaff

    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update staff member'
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteStaffMember = async (id: number) => {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`/api/staff/${id}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to delete staff member')
      }

      // 一覧から削除
      const index = staff.value.findIndex(s => s.id === id)
      if (index !== -1) {
        staff.value.splice(index, 1)
        total.value -= 1
      }

      // 現在のスタッフをクリア
      if (currentStaff.value && currentStaff.value.id === id) {
        currentStaff.value = null
      }

    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to delete staff member'
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateStaffStatus = async (id: number, status: Staff['status']) => {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`/api/staff/${id}/status`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ status }),
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to update staff status')
      }

      const updatedStaff: Staff = await response.json()

      // 一覧内のスタッフを更新
      const index = staff.value.findIndex(s => s.id === id)
      if (index !== -1) {
        staff.value[index] = updatedStaff
      }

      // 現在のスタッフも更新
      if (currentStaff.value && currentStaff.value.id === id) {
        currentStaff.value = updatedStaff
      }

      return updatedStaff

    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update staff status'
      throw err
    } finally {
      loading.value = false
    }
  }

  const assignToStore = async (id: number, storeId: number) => {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`/api/staff/${id}/store`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ store_id: storeId }),
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to assign staff to store')
      }

      const updatedStaff: Staff = await response.json()

      // 一覧内のスタッフを更新
      const index = staff.value.findIndex(s => s.id === id)
      if (index !== -1) {
        staff.value[index] = updatedStaff
      }

      // 現在のスタッフも更新
      if (currentStaff.value && currentStaff.value.id === id) {
        currentStaff.value = updatedStaff
      }

      return updatedStaff

    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to assign staff to store'
      throw err
    } finally {
      loading.value = false
    }
  }

  const searchStaff = async (searchParams: StaffSearchParams) => {
    currentPage.value = 1 // 検索時はページをリセット
    await fetchStaff(searchParams)
  }

  const setPage = async (page: number) => {
    currentPage.value = page
    await fetchStaff()
  }

  const setItemsPerPage = async (limit: number) => {
    itemsPerPage.value = limit
    currentPage.value = 1 // ページサイズ変更時はページをリセット
    await fetchStaff()
  }

  const clearError = () => {
    error.value = null
  }

  const clearCurrentStaff = () => {
    currentStaff.value = null
  }

  const getStaffStatusLabel = (status: Staff['status']): string => {
    const statusLabels: Record<Staff['status'], string> = {
      active: '在籍',
      inactive: '休職',
      on_leave: '休暇中'
    }
    return statusLabels[status] || status
  }

  const getStaffStatusColor = (status: Staff['status']): string => {
    const statusColors: Record<Staff['status'], string> = {
      active: 'success',
      inactive: 'warning',
      on_leave: 'info'
    }
    return statusColors[status] || 'default'
  }

  return {
    // State
    staff,
    currentStaff,
    loading,
    error,
    total,
    currentPage,
    itemsPerPage,

    // Computed
    totalPages,
    activeStaff,
    inactiveStaff,
    onLeaveStaff,
    statusCounts,
    staffByStore,

    // Actions
    fetchStaff,
    fetchStaffMember,
    createStaffMember,
    updateStaffMember,
    deleteStaffMember,
    updateStaffStatus,
    assignToStore,
    searchStaff,
    setPage,
    setItemsPerPage,
    clearError,
    clearCurrentStaff,
    getStaffStatusLabel,
    getStaffStatusColor
  }
})