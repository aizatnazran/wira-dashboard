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
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { computed } from 'vue'
import { useToast } from 'vue-toastification'

const router = useRouter()
const store = useStore()
const toast = useToast()
const isAuthenticated = computed(() => store.getters.isAuthenticated)

const handleLogout = async () => {
  try {
    await store.dispatch('logout')
    if (!store.getters.isAuthenticated) {
      router.push('/login')
      toast.success('Successfully logged out')
    }
  } catch (error) {
    console.error('Logout failed:', error)
    toast.error('Logout failed. Please try again.')
  }
}

store.dispatch('checkAuth')
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Cinzel:wght@400;500;600;700&family=Montserrat:wght@300;400;500;600&display=swap');

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
</style>
