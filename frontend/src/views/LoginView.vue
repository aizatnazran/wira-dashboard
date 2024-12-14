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
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import Swal from 'sweetalert2'
import axios from 'axios'
import { useToast } from 'vue-toastification'

export default {
  setup() {
    const router = useRouter()
    const store = useStore()
    const toast = useToast()
    const username = ref('')
    const password = ref('')
    const twoFactorCode = ref('')
    const show2FAInput = ref(false)
    const loading = ref(false)

    const handleSubmit = async () => {
      try {
        loading.value = true

        if (!show2FAInput.value) {
          // First step: Login with username and password
          const response = await axios.post('http://localhost:8080/api/auth/login', {
            username: username.value,
            password: password.value
          })

          if (response.data.requires_2fa) {
            show2FAInput.value = true
            return
          }

          const user = {
            id: response.data.user.id,
            username: response.data.user.username,
            email: response.data.user.email,
            token: response.data.token
          }

          store.commit('SET_USER', user)
          localStorage.setItem('user', JSON.stringify(user))
          localStorage.setItem('token', response.data.token)

          // Show success toast
          toast.success('Welcome back to WIRA!', {
            position: 'top-right',
            timeout: 5000,
            closeOnClick: true,
            pauseOnHover: true,
            draggable: true,
          })

          router.push('/')
        } else {
          // Second step: Verify 2FA code
          const response = await axios.post('http://localhost:8080/api/auth/2fa/login/verify', {
            username: username.value,
            code: twoFactorCode.value
          })

          const user = {
            id: response.data.user.id,
            username: response.data.user.username,
            email: response.data.user.email,
            token: response.data.token
          }

          store.commit('SET_USER', user)
          localStorage.setItem('user', JSON.stringify(user))
          localStorage.setItem('token', response.data.token)

          // Show success toast
          toast.success('Welcome back to WIRA!', {
            position: 'top-right',
            timeout: 5000,
            closeOnClick: true,
            pauseOnHover: true,
            draggable: true,
          })

          router.push('/')
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
        store.commit('LOGOUT')
        localStorage.removeItem('user')
        localStorage.removeItem('token')
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
