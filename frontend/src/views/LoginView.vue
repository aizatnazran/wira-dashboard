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
      <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
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
            <label for="password" class="sr-only">Password</label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-ac-gold/30 bg-ac-dark text-ac-light rounded-b-md focus:outline-none focus:ring-ac-gold focus:border-ac-gold focus:z-10 sm:text-sm"
              placeholder="Password"
            >
          </div>
        </div>

        <div>
          <button
            type="submit"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-black bg-ac-gold hover:bg-ac-gold/90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-ac-gold"
          >
            Sign in
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import { useToast } from 'vue-toastification'

export default {
  name: 'LoginView',
  setup() {
    const store = useStore()
    const router = useRouter()
    const toast = useToast()
    
    const form = ref({
      username: '',
      password: ''
    })

    const handleLogin = async () => {
      try {
        const success = await store.dispatch('login', form.value)
        if (success) {
          toast.success('Welcome back to WIRA!', {
            position: 'top-right',
            timeout: 3000,
            closeOnClick: true,
            pauseOnHover: true,
            draggable: true,
          })
          router.push('/')
        } else {
          Swal.fire({
            title: 'Login Failed',
            text: 'Invalid username or password',
            icon: 'error',
            confirmButtonText: 'Try Again',
            background: '#2A2A2A',
            color: '#F5F5F5',
            confirmButtonColor: '#C6A875'
          })
        }
      } catch (error) {
        console.error('Login error:', error)
      }
    }

    return {
      form,
      handleLogin
    }
  }
}
</script>
