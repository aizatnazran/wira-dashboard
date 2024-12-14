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
      <form class="mt-8 space-y-6" @submit.prevent="handleSignup">
        <div class="rounded-md shadow-sm -space-y-px">
          <div>
            <label for="username" class="sr-only">Username</label>
            <input
              id="username"
              v-model="form.username"
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
              v-model="form.email"
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
              v-model="form.password"
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
              v-model="form.confirmPassword"
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
            :disabled="isSubmitting"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-black bg-ac-gold hover:bg-ac-gold/90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-ac-gold"
          >
            <span v-if="!isSubmitting">Create Account</span>
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
import Swal from 'sweetalert2'
import { useToast } from 'vue-toastification'

export default {
  name: 'SignupView',
  setup() {
    const router = useRouter()
    const toast = useToast()

    const form = ref({
      username: '',
      email: '',
      password: '',
      confirmPassword: ''
    })

    const isSubmitting = ref(false)

    const handleSignup = async () => {
      if (form.value.password !== form.value.confirmPassword) {
        Swal.fire({
          title: 'Password Mismatch',
          text: 'The passwords you entered do not match',
          icon: 'error',
          confirmButtonText: 'Try Again',
          background: '#2A2A2A',
          color: '#F5F5F5',
          confirmButtonColor: '#C6A875'
        })
        return
      }

      isSubmitting.value = true

      try {
        const response = await fetch('http://localhost:8080/api/auth/register', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            username: form.value.username,
            password: form.value.password,
            email: form.value.email
          })
        });

        if (!response.ok) {
          const data = await response.json()
          if (data.error === 'Username already exists') {
            throw new Error('The username is already taken. Please choose another one.')
          }
          if (data.error === 'Email already exists') {
            throw new Error('An account with this email already exists. Please use a different email.')
          }
          throw new Error(data.error || 'Failed to create user')
        }

        toast.success('Account created successfully!', {
          position: 'top-right',
          timeout: 3000,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
        })
        router.push('/login')
      } catch (error) {
        console.error('Signup error:', error)
        Swal.fire({
          title: 'Signup Failed',
          text: error.message || 'An error occurred during signup',
          icon: 'error',
          confirmButtonText: 'Try Again',
          background: '#2A2A2A',
          color: '#F5F5F5',
          confirmButtonColor: '#C6A875'
        })
      } finally {
        isSubmitting.value = false
      }
    }

    return {
      form,
      isSubmitting,
      handleSignup
    }
  }
}
</script>
