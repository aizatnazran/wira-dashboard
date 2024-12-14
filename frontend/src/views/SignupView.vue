<template>
  <div class="min-h-[calc(100vh-4rem)] flex items-center justify-center px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 bg-ac-gray p-8 rounded-lg border border-ac-gold/30">
      <div>
        <h2 class="mt-6 text-center text-3xl font-cinzel text-ac-gold">Create your WIRA account</h2>
        <p class="mt-2 text-center text-sm text-ac-light">
          Or
          <router-link to="/login" class="font-medium text-ac-gold hover:text-ac-light">
            sign in to your account
          </router-link>
        </p>
      </div>
      <form class="mt-8 space-y-6" @submit.prevent="handleSubmit">
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
            <label for="email" class="sr-only">Email</label>
            <input
              id="email"
              v-model="email"
              type="email"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-ac-gold/30 bg-ac-dark text-ac-light focus:outline-none focus:ring-ac-gold focus:border-ac-gold focus:z-10 sm:text-sm"
              placeholder="Email"
            >
          </div>
          <div>
            <label for="password" class="sr-only">Password</label>
            <input
              id="password"
              v-model="password"
              type="password"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-ac-gold/30 bg-ac-dark text-ac-light focus:outline-none focus:ring-ac-gold focus:border-ac-gold focus:z-10 sm:text-sm"
              placeholder="Password"
            >
          </div>
          <div>
            <label for="confirmPassword" class="sr-only">Confirm Password</label>
            <input
              id="confirmPassword"
              v-model="confirmPassword"
              type="password"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-ac-gold/30 bg-ac-dark text-ac-light rounded-b-md focus:outline-none focus:ring-ac-gold focus:border-ac-gold focus:z-10 sm:text-sm"
              placeholder="Confirm Password"
            >
          </div>
        </div>

        <div>
          <button
            type="submit"
            :disabled="loading"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-black bg-ac-gold hover:bg-ac-gold/90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-ac-gold"
          >
            <span v-if="!loading">Create Account</span>
            <span v-else class="flex items-center">
              <!-- Stable FontAwesome Spinner -->
              <i class="fa fa-circle-notch fa-spin text-white text-xl mr-2"></i>
              Creating Account...
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
import { useToast } from 'vue-toastification'
import api from '@/api/config'

export default {
  name: 'SignupView',
  setup() {
    const router = useRouter()
    const toast = useToast()

    const username = ref('')
    const email = ref('')
    const password = ref('')
    const confirmPassword = ref('')
    const loading = ref(false)

    const validateForm = () => {
      // Username validation
      if (!username.value || username.value.length < 3) {
        toast.error('Username must be at least 3 characters long')
        return false
      }

      // Email validation
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if (!email.value || !emailRegex.test(email.value)) {
        toast.error('Please enter a valid email address')
        return false
      }

      // Password validation
      if (!password.value || password.value.length < 6) {
        toast.error('Password must be at least 6 characters long')
        return false
      }

      return true
    }

    const handleSubmit = async () => {
      try {
        loading.value = true

        if (!validateForm()) {
          loading.value = false
          return
        }

        if (password.value !== confirmPassword.value) {
          toast.error('The passwords you entered do not match')
          loading.value = false
          return
        }

        const response = await api.post('/api/auth/register', {
          username: username.value,
          password: password.value,
          email: email.value
        })

        if (response.status === 201) {
          toast.success('Registration successful! Please log in.')
          router.push('/login')
        }
      } catch (error) {
        console.error('Registration error:', error)
        if (error.response?.status === 409) {
          const message = error.response.data?.message || 'Username or email already exists'
          toast.error(message)
        } else {
          toast.error('Registration failed. Please try again.')
        }
      } finally {
        loading.value = false
      }
    }

    return {
      username,
      email,
      password,
      confirmPassword,
      loading,
      handleSubmit
    }
  }
}
</script>
