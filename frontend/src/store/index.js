import { createStore } from 'vuex'
import router from '@/router'
import { useToast } from 'vue-toastification'
import api from '@/api/config'

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
      } else {
        localStorage.removeItem('token')
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
      if (state.sessionCheckInterval) {
        clearInterval(state.sessionCheckInterval)
        state.sessionCheckInterval = null
      }
    }
  },
  actions: {
    async login({ commit, dispatch }, credentials) {
      try {
        const response = await api.post('/api/auth/login', credentials)
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
          try {
            // Set session ID in header
            const headers = { 'X-Session-ID': state.sessionID }
            await api.post('/api/auth/logout', null, { headers })
          } catch (error) {
            // Ignore logout errors since we're clearing the session anyway
            console.debug('Logout request failed:', error)
          }
        }
      } finally {
        commit('clearSession')
        router.push('/login')
      }
    },
    
    async checkSession({ state, dispatch }) {
      try {
        if (!state.sessionID) return
        
        await api.post('/api/auth/validate-session', {
          sessionID: state.sessionID
        })
      } catch (error) {
        // Only log non-401 errors as actual errors
        if (error.response?.status !== 401) {
          console.error('Unexpected session validation error:', error)
        } else {
          console.debug('Session expired naturally')
        }
        dispatch('handleSessionExpired')
      }
    },
    
    startSessionCheck({ dispatch, state }) {
      if (state.sessionCheckInterval) {
        clearInterval(state.sessionCheckInterval)
      }
      
      // Check session every 60 seconds
      state.sessionCheckInterval = setInterval(() => {
        dispatch('checkSession')
      }, 60000)
    },
    
    async handleSessionExpired({ dispatch }) {
      try {
        await dispatch('logout')
        useToast().warning('Your session has expired. Please log in again.', {
          timeout: 5000,
          closeOnClick: true
        })
      } catch (error) {
        console.error('Error handling session expiration:', error)
      }
    },
    
    async verify2FALogin({ commit, dispatch }, { username, code }) {
      try {
        const response = await api.post('/api/auth/2fa/login/verify', { username, code })
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
