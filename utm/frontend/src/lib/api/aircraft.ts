import { apiClient } from './client'
import type { Aircraft } from '@/types'

export const aircraftApi = {
  list: async (): Promise<Aircraft[]> => {
    const { data } = await apiClient.get<Aircraft[]>('/aircraft')
    return data
  },

  get: async (id: string): Promise<Aircraft> => {
    const { data } = await apiClient.get<Aircraft>(`/aircraft/${id}`)
    return data
  },

  create: async (payload: Omit<Aircraft, 'id' | 'owner_id' | 'created_at' | 'owner'>): Promise<Aircraft> => {
    const { data } = await apiClient.post<Aircraft>('/aircraft', payload)
    return data
  },

  update: async (id: string, payload: Partial<Aircraft>): Promise<Aircraft> => {
    const { data } = await apiClient.put<Aircraft>(`/aircraft/${id}`, payload)
    return data
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/aircraft/${id}`)
  },
}
