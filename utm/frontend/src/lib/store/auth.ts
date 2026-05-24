import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '@/types'

interface AuthState {
  user: User | null
  token: string | null
  setAuth: (user: User, token: string) => void
  logout: () => void
  isAuthenticated: () => boolean
  hasRole: (...roles: User['role'][]) => boolean
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,

      setAuth: (user, token) => {
        localStorage.setItem('utm_token', token)
        set({ user, token })
      },

      logout: () => {
        localStorage.removeItem('utm_token')
        set({ user: null, token: null })
      },

      isAuthenticated: () => !!get().token,

      hasRole: (...roles) => {
        const user = get().user
        return user ? roles.includes(user.role) : false
      },
    }),
    {
      name: 'utm-auth',
      partialize: (state) => ({ user: state.user, token: state.token }),
    }
  )
)
