import { createStore } from 'vuex'
import axios from 'axios'
import router from '@/router'
import { useToast } from 'vue-toastification'

export default createStore({
  state: {
    user: null,
    token: null,
    sessionID: null,
    sessionCheckInterval: null
  },
  getters: {
    isAuthenticated: state => !!state.user,
    token: state => state.token,
    user: state => state.user
  },
  mutations: {
    setUser(state, user) {
      state.user = user
      if (user) {
        localStorage.setItem('user', JSON.stringify(user))
      } else {
        localStorage.removeItem('user')
      }
    },
    setToken(state, token) {
      state.token = token
      if (token) {
        localStorage.setItem('token', token)
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
      } else {
        localStorage.removeItem('token')
        delete axios.defaults.headers.common['Authorization']
      }
    },
    setSessionID(state, sessionID) {
      state.sessionID = sessionID
      if (sessionID) {
        localStorage.setItem('sessionID', sessionID)
      } else {
        localStorage.removeItem('sessionID')
      }
    },
    clearSession(state) {
      state.user = null
      state.token = null
      state.sessionID = null
      localStorage.removeItem('user')
      localStorage.removeItem('token')
      localStorage.removeItem('sessionID')
      delete axios.defaults.headers.common['Authorization']
      if (state.sessionCheckInterval) {
        clearInterval(state.sessionCheckInterval)
        state.sessionCheckInterval = null
      }
    }
  },
  actions: {
    async login({ commit, dispatch }, credentials) {
      try {
        const response = await axios.post('/api/auth/login', credentials)
        const { token, user, sessionID } = response.data
        
        commit('setUser', user)
        commit('setToken', token)
        commit('setSessionID', sessionID)
        
        // Start session check
        dispatch('startSessionCheck')
        
        router.push('/')
        return true
      } catch (error) {
        console.error('Login failed:', error)
        throw error
      }
    },
    
    async logout({ commit, state }) {
      try {
        if (state.sessionID) {
          // Set session ID in header
          const headers = { 'X-Session-ID': state.sessionID }
          await axios.post('/api/auth/logout', null, { headers })
        }
      } catch (error) {
        console.error('Logout error:', error)
      } finally {
        commit('clearSession')
        router.push('/login')
      }
    },
    
    async checkSession({ state, dispatch }) {
      try {
        if (!state.sessionID) return
        
        await axios.post('/api/auth/validate-session', {
          sessionID: state.sessionID
        })
      } catch (error) {
        console.error('Session validation failed:', error)
        dispatch('handleSessionExpired')
      }
    },
    
    startSessionCheck({ dispatch, state }) {
      if (state.sessionCheckInterval) {
        clearInterval(state.sessionCheckInterval)
      }
      
      // Check session every minute
      state.sessionCheckInterval = setInterval(() => {
        dispatch('checkSession')
      }, 60000) // 1 minute
    },
    
    async handleSessionExpired({ dispatch }) {
      await dispatch('logout')
      useToast().error('Session expired. Please log in again.', {
        timeout: false, // Toast will not auto-close
        closeOnClick: false // Cannot be closed by clicking
      })
    },
    
    async verify2FALogin({ commit, dispatch }, { username, code }) {
      try {
        const response = await axios.post('/api/auth/2fa/login/verify', { username, code })
        const { token, user, sessionID } = response.data
        
        commit('setUser', user)
        commit('setToken', token)
        commit('setSessionID', sessionID)
        
        // Start session check
        dispatch('startSessionCheck')
        
        router.push('/')
        return true
      } catch (error) {
        console.error('2FA login failed:', error)
        throw error
      }
    }
  }
})
