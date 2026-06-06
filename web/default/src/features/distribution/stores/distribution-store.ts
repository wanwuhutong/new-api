import { create } from 'zustand'
import {
  getAllDistributors,
  getDistributorStatistics,
  createDistributor,
  updateDistributor,
  deleteDistributor,
  settleDistributorCommission,
  type Distributor,
  type DistributorStatistics,
} from '../api'

interface DistributionState {
  distributors: Distributor[]
  total: number
  page: number
  pageSize: number
  statistics: DistributorStatistics | null
  loading: boolean
  error: string | null
  fetchDistributors: (page?: number, pageSize?: number) => Promise<void>
  fetchStatistics: () => Promise<void>
  create: (data: { user_id: number; level: number; parent_id?: number }) => Promise<void>
  update: (id: number, data: { level?: number; status?: number }) => Promise<void>
  remove: (id: number) => Promise<void>
  settle: (id: number) => Promise<void>
  setPage: (page: number) => void
  setPageSize: (pageSize: number) => void
}

export const useDistributionStore = create<DistributionState>((set, get) => ({
  distributors: [],
  total: 0,
  page: 1,
  pageSize: 20,
  statistics: null,
  loading: false,
  error: null,

  fetchDistributors: async (page?: number, pageSize?: number) => {
    set({ loading: true, error: null })
    try {
      const p = page ?? get().page
      const ps = pageSize ?? get().pageSize
      const response = await getAllDistributors(p, ps)
      set({
        distributors: response.distributors,
        total: response.total,
        page: p,
        pageSize: ps,
        loading: false,
      })
    } catch (error: any) {
      set({ error: error.message || 'Failed to fetch distributors', loading: false })
    }
  },

  fetchStatistics: async () => {
    try {
      const statistics = await getDistributorStatistics()
      set({ statistics })
    } catch (error: any) {
      set({ error: error.message || 'Failed to fetch statistics' })
    }
  },

  create: async (data) => {
    set({ loading: true, error: null })
    try {
      await createDistributor(data)
      await get().fetchDistributors()
      await get().fetchStatistics()
    } catch (error: any) {
      set({ error: error.message || 'Failed to create distributor', loading: false })
      throw error
    }
  },

  update: async (id, data) => {
    set({ loading: true, error: null })
    try {
      await updateDistributor(id, data)
      await get().fetchDistributors()
    } catch (error: any) {
      set({ error: error.message || 'Failed to update distributor', loading: false })
      throw error
    }
  },

  remove: async (id) => {
    set({ loading: true, error: null })
    try {
      await deleteDistributor(id)
      await get().fetchDistributors()
      await get().fetchStatistics()
    } catch (error: any) {
      set({ error: error.message || 'Failed to delete distributor', loading: false })
      throw error
    }
  },

  settle: async (id) => {
    set({ loading: true, error: null })
    try {
      await settleDistributorCommission(id)
      await get().fetchDistributors()
      await get().fetchStatistics()
    } catch (error: any) {
      set({ error: error.message || 'Failed to settle commission', loading: false })
      throw error
    }
  },

  setPage: (page) => {
    set({ page })
    get().fetchDistributors(page)
  },

  setPageSize: (pageSize) => {
    set({ pageSize, page: 1 })
    get().fetchDistributors(1, pageSize)
  },
}))
