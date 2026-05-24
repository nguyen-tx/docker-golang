import { apiClient } from './client'
import type { LoginResponse, User } from '@/types'

export const authApi = {
  login: async (email: string, password: string): Promise<LoginResponse> => {
    const { data } = await apiClient.post<LoginResponse>('/auth/login', { email, password })
    return data
  },

  register: async (payload: {
    email: string
    password: string
    full_name: string
    phone: string
    organization?: string
    license_no?: string
  }): Promise<User> => {
    const { data } = await apiClient.post<User>('/auth/register', payload)
    return data
  },

  me: async (): Promise<User> => {
    const { data } = await apiClient.get<User>('/auth/me')
    return data
  },
}
