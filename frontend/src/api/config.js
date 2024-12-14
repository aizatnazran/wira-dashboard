// API configuration
const API_CONFIG = {
  development: {
    baseURL: 'http://localhost:8080',
  },
  production: {
    baseURL: process.env.VUE_APP_API_URL || 'https://wira.aizat.dev/api',
  }
}

// Get current environment
const env = process.env.NODE_ENV || 'development'

// Export the configuration for the current environment
export const config = API_CONFIG[env]

// Create and configure axios instance
import axios from 'axios'

const api = axios.create({
  baseURL: config.baseURL,
  withCredentials: true,
})

// Add request interceptor to add auth token
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export default api
