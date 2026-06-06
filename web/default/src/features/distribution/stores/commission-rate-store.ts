import { create } from 'zustand'
import {
  getAllCommissionRates,
  createCommissionRate,
  updateCommissionRate,
  deleteCommissionRate,
  type CommissionRate,
} from '../api'

interface CommissionRateState {
  rates: CommissionRate[]
  loading: boolean
  error: string | null
  fetchRates: () => Promise<void>
  create: (data: {
    name: string
    type: string
    type_id?: number
    level1_rate: number
    level2_rate: number
    level3_rate: number
    enabled: boolean
  }) => Promise<void>
  update: (id: number, data: {
    name?: string
    level1_rate: number
    level2_rate: number
    level3_rate: number
    enabled: boolean
  }) => Promise<void>
  remove: (id: number) => Promise<void>
}

export const useCommissionRateStore = create<CommissionRateState>((set, get) => ({
  rates: [],
  loading: false,
  error: null,

  fetchRates: async () => {
    set({ loading: true, error: null })
    try {
      const rates = await getAllCommissionRates()
      set({ rates, loading: false })
    } catch (error: any) {
      set({ error: error.message || 'Failed to fetch commission rates', loading: false })
    }
  },

  create: async (data) => {
    set({ loading: true, error: null })
    try {
      await createCommissionRate(data)
      await get().fetchRates()
    } catch (error: any) {
      set({ error: error.message || 'Failed to create commission rate', loading: false })
      throw error
    }
  },

  update: async (id, data) => {
    set({ loading: true, error: null })
    try {
      await updateCommissionRate(id, data)
      await get().fetchRates()
    } catch (error: any) {
      set({ error: error.message || 'Failed to update commission rate', loading: false })
      throw error
    }
  },

  remove: async (id) => {
    set({ loading: true, error: null })
    try {
      await deleteCommissionRate(id)
      await get().fetchRates()
    } catch (error: any) {
      set({ error: error.message || 'Failed to delete commission rate', loading: false })
      throw error
    }
  },
}))
