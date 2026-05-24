import { apiClient } from './client'
import type { AuthorizationRequest, AuthorizationStatus } from '@/types'

export const authorizationApi = {
  list: async (status?: AuthorizationStatus): Promise<AuthorizationRequest[]> => {
    const params = status ? { status } : {}
    const { data } = await apiClient.get<AuthorizationRequest[]>('/authorizations', { params })
    return data
  },

  get: async (id: string): Promise<AuthorizationRequest> => {
    const { data } = await apiClient.get<AuthorizationRequest>(`/authorizations/${id}`)
    return data
  },

  create: async (payload: {
    aircraft_id: string
    zone_id?: string
    purpose: string
    flight_area_json: string
    planned_alt_m: number
    planned_start_at: string
    planned_end_at: string
    pilot_notes?: string
  }): Promise<AuthorizationRequest> => {
    const { data } = await apiClient.post<AuthorizationRequest>('/authorizations', payload)
    return data
  },

  review: async (id: string, action: 'approve' | 'reject', review_notes?: string): Promise<AuthorizationRequest> => {
    const { data } = await apiClient.post<AuthorizationRequest>(`/authorizations/${id}/review`, {
      action,
      review_notes,
    })
    return data
  },

  cancel: async (id: string): Promise<void> => {
    await apiClient.post(`/authorizations/${id}/cancel`)
  },
}
