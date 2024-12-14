import { createStore } from 'vuex'
import Swal from 'sweetalert2'

export default createStore({
  state: {
    user: null,
    isAuthenticated: false,
  },
  mutations: {
    SET_USER(state, user) {
      state.user = user
      state.isAuthenticated = !!user
    },
    LOGOUT(state) {
      state.user = null
      state.isAuthenticated = false
    }
  },
  actions: {
    async login({ commit }, credentials) {
      try {
        const response = await fetch('http://localhost:8080/api/auth/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(credentials)
        });

        if (!response.ok) {
          throw new Error('Invalid credentials');
        }

        const data = await response.json();
        const user = {
          id: data.user.id,
          username: data.user.username,
          email: data.user.email,
          token: data.token
        };

        commit('SET_USER', user);
        localStorage.setItem('user', JSON.stringify(user));
        localStorage.setItem('token', data.token);
        return true;
      } catch (error) {
        console.error('Login error:', error);
        commit('LOGOUT');
        localStorage.removeItem('user');
        localStorage.removeItem('token');
        return false;
      }
    },
    async logout({ commit }) {
      const result = await Swal.fire({
        title: 'Logout Confirmation',
        text: 'Are you sure you want to logout?',
        icon: 'question',
        showCancelButton: true,
        confirmButtonText: 'Yes',
        cancelButtonText: 'No',
        background: '#2A2A2A',
        color: '#F5F5F5',
        confirmButtonColor: '#C6A875',
        cancelButtonColor: '#6B7280'
      });

      if (result.isConfirmed) {
        commit('LOGOUT')
        localStorage.removeItem('user')
        localStorage.removeItem('token')
      }
    },
    checkAuth({ commit }) {
      const user = localStorage.getItem('user')
      if (user) {
        commit('SET_USER', JSON.parse(user))
      }
    }
  },
  getters: {
    isAuthenticated: state => state.isAuthenticated,
    currentUser: state => state.user,
    token: state => state.user ? state.user.token : null
  }
})
