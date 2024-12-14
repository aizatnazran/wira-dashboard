// API configuration
const API_CONFIG = {
  development: {
    baseURL: 'http://localhost:8080',
  },
  docker: {
    baseURL: '/api', 
  },
  production: {
    baseURL: '', 
  }
}

// Get current environment
const env = process.env.VUE_APP_ENV || process.env.NODE_ENV || 'development'

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
