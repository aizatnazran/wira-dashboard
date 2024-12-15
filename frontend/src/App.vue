<template>
  <div class="min-h-screen bg-ac-dark text-ac-light">
    <!-- Navigation -->
    <nav class="bg-ac-gray border-b border-ac-gold/30">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center">
            <div class="flex-shrink-0 flex items-center">
              <img src="@/assets/logo.png" alt="WIRA Logo" class="h-8 w-8 mr-2">
              <router-link to="/" class="text-ac-gold font-cinzel text-xl">WIRA</router-link>
            </div>
          </div>
          <div class="flex items-center space-x-1 sm:space-x-4">
            <template v-if="isAuthenticated">
              <router-link 
                to="/" 
                class="text-ac-gold hover:text-ac-light px-2 sm:px-3 py-2 rounded-md text-sm font-medium"
                :class="{ 'bg-ac-dark': $route.path === '/' }"
              >
                Rankings
              </router-link>
              <router-link 
                to="/classes" 
                class="text-ac-gold hover:text-ac-light px-2 sm:px-3 py-2 rounded-md text-sm font-medium"
                :class="{ 'bg-ac-dark': $route.path === '/classes' }"
              >
                Classes
              </router-link>
              <router-link 
                to="/profile" 
                class="text-ac-gold hover:text-ac-light px-2 sm:px-3 py-2 rounded-md text-sm font-medium"
                :class="{ 'bg-ac-dark': $route.path === '/profile' }"
              >
                Profile
              </router-link>
              <button 
                @click="handleLogout" 
                class="text-ac-gold hover:text-ac-light px-2 sm:px-3 py-2 rounded-md text-sm font-medium"
              >
                Logout
              </button>
            </template>
            <template v-else>
              <router-link 
                to="/login" 
                class="text-ac-gold hover:text-ac-light px-2 sm:px-3 py-2 rounded-md text-sm font-medium"
                :class="{ 'bg-ac-dark': $route.path === '/login' }"
              >
                Login
              </router-link>
              <router-link 
                to="/signup" 
                class="text-ac-gold hover:text-ac-light px-2 sm:px-3 py-2 rounded-md text-sm font-medium"
                :class="{ 'bg-ac-dark': $route.path === '/signup' }"
              >
                Sign Up
              </router-link>
            </template>
          </div>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <router-view></router-view>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useToast } from 'vue-toastification'
import Swal from 'sweetalert2'

const store = useStore()
const toast = useToast()

const isAuthenticated = computed(() => store.state.user !== null)

const handleLogout = async () => {
  const result = await Swal.fire({
    title: 'Logout Confirmation',
    text: 'Are you sure you want to log out?',
    icon: 'question',
    showCancelButton: true,
    confirmButtonText: 'Yes, log out',
    cancelButtonText: 'Cancel',
    background: '#2A2A2A',
    color: '#F5F5F5',
    confirmButtonColor: '#C6A875',
    cancelButtonColor: '#4B5563'
  })

  if (result.isConfirmed) {
    await store.dispatch('logout')
    toast.success('Successfully logged out', {
      position: 'top-right',
      timeout: 5000,
      closeOnClick: true,
      pauseOnHover: true,
      draggable: true,
    })
  }
}

// Check stored user session
onMounted(() => {
  const storedUser = localStorage.getItem('user')
  const storedToken = localStorage.getItem('token')
  const storedSessionID = localStorage.getItem('sessionID')
  
  if (storedUser && storedToken && storedSessionID) {
    store.commit('setUser', JSON.parse(storedUser))
    store.commit('setToken', storedToken)
    store.commit('setSessionID', storedSessionID)
    store.dispatch('startSessionCheck')
  }
})
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Cinzel:wght@400;500;600;700&family=Montserrat:wght@300;400;500;600&display=swap');
@import '@/assets/css/toast.css';

@media (max-width: 640px) {
  .Vue-Toastification__toast {
    font-size: 0.875rem !important;
    padding: 8px 12px !important;
    min-height: 48px !important;
    margin-bottom: 0.5rem !important;
  }
  
  .Vue-Toastification__toast-body {
    margin: 4px 8px !important;
  }
}

.font-cinzel {
  font-family: 'Cinzel', serif;
}

.font-montserrat {
  font-family: 'Montserrat', sans-serif;
}

body {
  margin: 0;
  padding: 0;
}

#app {
  font-family: 'Montserrat', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.dark-mode-popup {
  border: 1px solid #333 !important;
}
</style>
