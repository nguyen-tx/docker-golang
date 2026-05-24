import axios from 'axios'

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080/api/v1'

export const apiClient = axios.create({
  baseURL: API_URL,
  headers: { 'Content-Type': 'application/json' },
})

// Attach JWT token từ localStorage
apiClient.interceptors.request.use((config) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('utm_token')
    if (token) config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Xử lý 401 - redirect về login
apiClient.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401 && typeof window !== 'undefined') {
      localStorage.removeItem('utm_token')
      window.location.href = '/auth/login'
    }
    return Promise.reject(err)
  }
)
