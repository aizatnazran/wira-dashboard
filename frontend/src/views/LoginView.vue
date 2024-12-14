<template>
  <div class="min-h-[calc(100vh-4rem)] flex items-center justify-center px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 bg-ac-gray p-8 rounded-lg border border-ac-gold/30">
      <div>
        <h2 class="mt-6 text-center text-3xl font-cinzel text-ac-gold">Sign in to WIRA</h2>
        <p class="mt-2 text-center text-sm text-ac-light">
          Or
          <router-link to="/signup" class="font-medium text-ac-gold hover:text-ac-light">
            create a new account
          </router-link>
        </p>
      </div>
      <form class="mt-8 space-y-6" @submit.prevent="handleSubmit">
        <!-- Username and Password Section -->
        <div class="rounded-md shadow-sm -space-y-px">
          <div>
            <label for="username" class="sr-only">Username</label>
            <input
              id="username"
              v-model="username"
              type="text"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-ac-gold/30 bg-ac-dark text-ac-light rounded-t-md focus:outline-none focus:ring-ac-gold focus:border-ac-gold focus:z-10 sm:text-sm"
              placeholder="Username"
            >
          </div>
          <div>
            <label for="password" class="sr-only">Password</label>
            <input
              id="password"
              v-model="password"
              type="password"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-ac-gold/30 bg-ac-dark text-ac-light rounded-b-md focus:outline-none focus:ring-ac-gold focus:border-ac-gold focus:z-10 sm:text-sm"
              placeholder="Password"
            >
          </div>
        </div>

        <!-- 2FA Section -->
        <div v-if="show2FAInput" class="mt-6 border-t border-ac-gold/30 pt-6">
          <h3 class="text-lg font-semibold text-ac-gold text-center mb-4">Two-Factor Authentication</h3>
          <div>
            <label for="twoFactorCode" class="sr-only">2FA Code</label>
            <input
              id="twoFactorCode"
              v-model="twoFactorCode"
              type="text"
              required
              class="appearance-none rounded-md block w-full px-3 py-2 border border-ac-gold/30 bg-ac-dark text-ac-light focus:outline-none focus:ring-ac-gold focus:border-ac-gold focus:z-10 sm:text-sm"
              placeholder="Enter your 2FA Code"
              maxlength="6"
            >
          </div>
        </div>

        <!-- Submit Button -->
        <div>
          <button
            type="submit"
            :disabled="loading"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-black bg-ac-gold hover:bg-ac-gold/90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-ac-gold"
          >
            <span v-if="!loading">Sign in</span>
            <span v-else class="flex items-center">
              <!-- FontAwesome Spinner Icon -->
              <i class="fa fa-circle-notch fa-spin text-white text-xl mr-2"></i>
              Signing In...
            </span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useStore } from 'vuex'
import Swal from 'sweetalert2'
import { useToast } from 'vue-toastification'
import api from '@/api/config'

export default {
  setup() {
    const store = useStore()
    const toast = useToast()
    const username = ref('')
    const password = ref('')
    const twoFactorCode = ref('')
    const show2FAInput = ref(false)
    const loading = ref(false)

    const validateLoginData = () => {
      if (!username.value || !password.value) {
        Swal.fire({
          title: 'Validation Error',
          text: 'Please fill in all required fields',
          icon: 'error',
          background: '#2A2A2A',
          color: '#F5F5F5',
          confirmButtonColor: '#C6A875'
        })
        return false
      }
      
      if (password.value.length < 6) {
        Swal.fire({
          title: 'Validation Error',
          text: 'Password must be at least 6 characters long',
          icon: 'error',
          background: '#2A2A2A',
          color: '#F5F5F5',
          confirmButtonColor: '#C6A875'
        })
        return false
      }

      return true
    }

    const handleSubmit = async () => {
      try {
        loading.value = true

        if (!validateLoginData()) {
          loading.value = false
          return
        }

        if (!show2FAInput.value) {
          const initialLoginResponse = await api.post('/api/auth/login', {
            username: username.value,
            password: password.value
          })

          if (initialLoginResponse.data.requires_2fa) {
            show2FAInput.value = true
            loading.value = false
            return
          }

          // Use the store's login action with the response data
          await store.dispatch('login', {
            username: username.value,
            password: password.value,
            token: initialLoginResponse.data.token,
            user: initialLoginResponse.data.user,
            sessionID: initialLoginResponse.data.sessionID
          })

          // Show success toast
          toast.success('Welcome back to WIRA!', {
            position: 'top-right',
            timeout: 5000,
            closeOnClick: true,
            pauseOnHover: true,
            draggable: true,
          })

        } else {
          // Second step: Verify 2FA code
          await store.dispatch('verify2FALogin', {
            username: username.value,
            code: twoFactorCode.value
          })

          // Show success toast
          toast.success('Welcome back to WIRA!', {
            position: 'top-right',
            timeout: 5000,
            closeOnClick: true,
            pauseOnHover: true,
            draggable: true,
          })
        }
      } catch (err) {
        // Error handling for invalid credentials or server errors
        const errorMessage = err.response?.data?.error || 'An error occurred during login'
        Swal.fire({
          title: 'Login Failed',
          text: errorMessage,
          icon: 'error',
          confirmButtonText: 'Try Again',
          background: '#2A2A2A',
          color: '#F5F5F5',
          confirmButtonColor: '#C6A875'
        })
      } finally {
        loading.value = false
      }
    }

    return {
      username,
      password,
      twoFactorCode,
      show2FAInput,
      loading,
      handleSubmit
    }
  }
}
</script>
