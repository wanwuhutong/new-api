import { api, handleResponse } from '@/lib/api'

export interface Distributor {
  id: number
  user_id: number
  username: string
  display_name: string
  level: number
  parent_id: number
  parent_username?: string
  status: number
  total_commission: number
  pending_commission: number
  available_commission: number
  direct_users: number
  child_distributors: number
  created_at: number
}

export interface DistributorListResponse {
  distributors: Distributor[]
  total: number
  page: number
  page_size: number
}

export interface CommissionRate {
  id: number
  name: string
  type: string
  type_id: number
  level1_rate: number
  level2_rate: number
  level3_rate: number
  enabled: boolean
  created_at: number
  updated_at: number
}

export interface CommissionLog {
  id: number
  distributor_id: number
  username: string
  user_display_name: string
  order_id: string
  order_type: string
  amount: number
  commission: number
  level: number
  status: number
  remark: string
  created_at: number
  settled_at?: number
}

export interface CommissionLogListResponse {
  logs: CommissionLog[]
  total: number
  page: number
  page_size: number
}

export interface DistributorDashboard {
  total_commission: number
  pending_commission: number
  available_commission: number
  withdrawn_commission: number
  total_users: number
  level1_users: number
  level2_users: number
  level3_users: number
  total_distributors: number
  level1_distributors: number
  level2_distributors: number
  level3_distributors: number
}

export interface DistributorStatistics {
  total_distributors: number
  active_distributors: number
  total_commission_paid: number
  total_commission_pending: number
  total_users_invited: number
}

export interface CreateDistributorRequest {
  user_id: number
  level: number
  parent_id?: number
  status?: number
}

export interface UpdateDistributorRequest {
  level?: number
  status?: number
}

export interface CreateCommissionRateRequest {
  name: string
  type: string
  type_id?: number
  level1_rate: number
  level2_rate: number
  level3_rate: number
  enabled: boolean
}

export interface UpdateCommissionRateRequest {
  name?: string
  level1_rate: number
  level2_rate: number
  level3_rate: number
  enabled: boolean
}

// Admin APIs
export const getAllDistributors = async (
  page = 1,
  pageSize = 20
): Promise<DistributorListResponse> => {
  const response = await api.get('/distribution/admin/', {
    params: { page, page_size: pageSize },
  })
  return handleResponse(response)
}

export const getDistributor = async (id: number): Promise<Distributor> => {
  const response = await api.get(`/distribution/admin/${id}`)
  return handleResponse(response)
}

export const createDistributor = async (
  data: CreateDistributorRequest
): Promise<Distributor> => {
  const response = await api.post('/distribution/admin/', data)
  return handleResponse(response)
}

export const updateDistributor = async (
  id: number,
  data: UpdateDistributorRequest
): Promise<void> => {
  const response = await api.put(`/distribution/admin/${id}`, data)
  return handleResponse(response)
}

export const deleteDistributor = async (id: number): Promise<void> => {
  const response = await api.delete(`/distribution/admin/${id}`)
  return handleResponse(response)
}

export const settleDistributorCommission = async (
  id: number
): Promise<void> => {
  const response = await api.post(`/distribution/admin/${id}/settle`)
  return handleResponse(response)
}

export const getDistributorStatistics = async (): Promise<DistributorStatistics> => {
  const response = await api.get('/distribution/admin/statistics')
  return handleResponse(response)
}

export const getAllCommissionLogs = async (
  page = 1,
  pageSize = 20
): Promise<CommissionLogListResponse> => {
  const response = await api.get('/distribution/admin/logs', {
    params: { page, page_size: pageSize },
  })
  return handleResponse(response)
}

// Commission Rate APIs
export const getAllCommissionRates = async (): Promise<CommissionRate[]> => {
  const response = await api.get('/commission-rate/')
  return handleResponse(response)
}

export const createCommissionRate = async (
  data: CreateCommissionRateRequest
): Promise<CommissionRate> => {
  const response = await api.post('/commission-rate/', data)
  return handleResponse(response)
}

export const updateCommissionRate = async (
  id: number,
  data: UpdateCommissionRateRequest
): Promise<void> => {
  const response = await api.put(`/commission-rate/${id}`, data)
  return handleResponse(response)
}

export const deleteCommissionRate = async (id: number): Promise<void> => {
  const response = await api.delete(`/commission-rate/${id}`)
  return handleResponse(response)
}

// Distributor Self-service APIs
export const getMyDistributorInfo = async (): Promise<Distributor> => {
  const response = await api.get('/distribution/self')
  return handleResponse(response)
}

export const getDistributorDashboard = async (): Promise<DistributorDashboard> => {
  const response = await api.get('/distribution/dashboard')
  return handleResponse(response)
}

export const getMyCommissionLogs = async (
  page = 1,
  pageSize = 20
): Promise<CommissionLogListResponse> => {
  const response = await api.get('/distribution/logs', {
    params: { page, page_size: pageSize },
  })
  return handleResponse(response)
}

export const getChildDistributors = async (
  page = 1,
  pageSize = 20
): Promise<DistributorListResponse> => {
  const response = await api.get('/distribution/children', {
    params: { page, page_size: pageSize },
  })
  return handleResponse(response)
}

export const getDistributorDirectUsers = async (
  page = 1,
  pageSize = 20
): Promise<{ users: any[]; total: number; page: number; page_size: number }> => {
  const response = await api.get('/distribution/users', {
    params: { page, page_size: pageSize },
  })
  return handleResponse(response)
}
