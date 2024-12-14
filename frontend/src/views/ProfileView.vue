<template>
  <div class="min-h-screen p-4">
    <div class="max-w-4xl mx-auto">
      <!-- Profile Header -->
      <div class="bg-ac-dark rounded-lg p-4 sm:p-8">
        <div class="flex flex-col sm:flex-row sm:items-start sm:space-x-8">
          <!-- Profile Image -->
          <div class="flex-shrink-0 mb-4 sm:mb-0">
            <img :src="profileImage" alt="Profile" class="w-24 h-24 sm:w-32 sm:h-32 rounded-full border-4 border-ac-gold mx-auto sm:mx-0" />
          </div>
          
          <!-- Profile Info -->
          <div class="flex-grow text-center sm:text-left">
            <h1 class="text-3xl sm:text-4xl font-cinzel text-ac-gold mb-2 break-words">{{ profile.username }}</h1>
            <p class="text-ac-light mb-4 break-words">{{ profile.email }}</p>
          </div>
        </div>
      </div>

      <!-- Character Stats -->
      <div class="bg-ac-dark rounded-lg p-4 sm:p-6 mt-4">
        <div class="grid grid-cols-1 gap-4 sm:gap-6 sm:grid-cols-2">
          <!-- User Information -->
          <div class="space-y-4 sm:space-y-6">
            <div>
              <h4 class="text-sm font-medium text-ac-gold">Account Information</h4>
              <div class="mt-2 border border-ac-gold/30 rounded-lg p-4 bg-ac-dark">
                <dl class="space-y-3">
                  <div>
                    <dt class="text-sm font-medium text-ac-light">Username</dt>
                    <dd class="mt-1 text-sm text-ac-gold break-words">{{ profile.username }}</dd>
                  </div>
                  <div>
                    <dt class="text-sm font-medium text-ac-light">Email</dt>
                    <dd class="mt-1 text-sm text-ac-gold break-words">{{ profile.email }}</dd>
                  </div>
                  <div>
                    <dt class="text-sm font-medium text-ac-light">Member Since</dt>
                    <dd class="mt-1 text-sm text-ac-gold">{{ formatMemberSince(profile.created_at) }}</dd>
                  </div>
                </dl>
              </div>
            </div>

            <!-- 2FA Security Section -->
            <div>
              <div class="flex items-center justify-between mb-2">
                <h4 class="text-sm font-medium text-ac-gold">Two-Factor Authentication</h4>
                <span :class="profile.two_factor_enabled ? 'bg-green-600' : 'bg-red-600'" class="px-2 py-1 rounded text-xs text-white">
                  {{ profile.two_factor_enabled ? 'Enabled' : 'Disabled' }}
                </span>
              </div>
              <div class="mt-2 border border-ac-gold/30 rounded-lg p-4 bg-ac-dark">
                <div v-if="!profile.two_factor_enabled" class="space-y-3">
                  <p class="text-sm text-ac-light">Enhance your account security by enabling two-factor authentication.</p>
                  <button 
                    @click="showEnableDialog = true"
                    class="w-full px-4 py-2 bg-ac-gold text-ac-dark rounded-md hover:bg-ac-gold/90 transition-colors"
                  >
                    Enable 2FA
                  </button>
                </div>
                <div v-else class="space-y-3">
                  <p class="text-sm text-ac-light">Your account is protected with two-factor authentication.</p>
                  <p class="text-xs text-ac-light/70">For additional security, you'll need to enter a code from your authenticator app when signing in.</p>
                  <button 
                    @click="handleDisable2FA"
                    class="w-full px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition-colors"
                  >
                    Disable 2FA
                  </button>
                </div>
              </div>
            </div>

            <!-- Character Stats -->
            <div>
              <h4 class="text-sm font-medium text-ac-gold">Character Statistics</h4>
              <div class="mt-2 border border-ac-gold/30 rounded-lg p-4 bg-ac-dark">
                <dl class="space-y-3">
                  <div>
                    <dt class="text-sm font-medium text-ac-light">Total Characters</dt>
                    <dd class="mt-1 text-sm text-ac-gold">3</dd>
                  </div>
                  <div>
                    <dt class="text-sm font-medium text-ac-light">Highest Ranked Class</dt>
                    <dd class="mt-1 text-sm text-ac-gold">{{ highestRanked.class }} - Rank #{{ highestRanked.rank }}</dd>
                  </div>
                </dl>
              </div>
            </div>
          </div>

          <!-- Characters List -->
          <div>
            <h4 class="text-sm font-medium text-ac-gold">Your Characters</h4>
            <div class="mt-2 space-y-4">
              <div v-for="(char, index) in characters" :key="index" 
                   class="bg-ac-gray rounded-lg p-4">
                <div class="flex justify-between items-start">
                  <div>
                    <h5 class="text-ac-gold font-cinzel break-words">{{ char.name }}</h5>
                    <p class="text-sm text-ac-light">{{ char.class }}</p>
                  </div>
                  <div class="text-right">
                    <p class="text-ac-gold">Rank #{{ char.rank }}</p>
                    <p class="text-sm text-ac-light">Score: {{ char.score }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Enable 2FA Modal -->
    <div v-if="showEnableDialog" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
      <div class="bg-ac-dark rounded-lg max-w-md w-full p-6">
        <h3 class="text-xl font-cinzel text-ac-gold mb-4">Enable Two-Factor Authentication</h3>
        
        <!-- Step 1: Password Verification -->
        <div v-if="step === 1" class="space-y-4">
          <p class="text-ac-light text-sm">Please enter your password to continue</p>
          <input 
            type="password" 
            v-model="password"
            placeholder="Enter your password"
            class="w-full px-4 py-2 bg-ac-gray border border-ac-gold/30 rounded-md text-ac-light focus:outline-none focus:border-ac-gold"
          />
          <div class="flex justify-end space-x-3">
            <button 
              @click="closeEnableDialog"
              class="px-4 py-2 text-ac-light hover:text-ac-gold transition-colors"
            >
              Cancel
            </button>
            <button 
              @click="verifyPassword"
              :disabled="!password"
              class="px-4 py-2 bg-ac-gold text-ac-dark rounded-md hover:bg-ac-gold/90 transition-colors disabled:opacity-50"
            >
              Continue
            </button>
          </div>
        </div>

        <!-- Step 2: QR Code -->
        <div v-if="step === 2" class="space-y-4">
          <p class="text-ac-light text-sm">Scan this QR code with your authenticator app</p>
          <div class="flex justify-center">
            <QRCodeVue
              :value="qrUrl"
              :size="200"
              level="M"
              render-as="svg"
              class="bg-white p-2 rounded"
            />
          </div>
          <p class="text-ac-light text-sm break-all bg-ac-gray p-2 rounded">
            Secret key: {{ secret }}
          </p>
          <div class="flex justify-end space-x-3">
            <button 
              @click="closeEnableDialog"
              class="px-4 py-2 text-ac-light hover:text-ac-gold transition-colors"
            >
              Cancel
            </button>
            <button 
              @click="step = 3"
              class="px-4 py-2 bg-ac-gold text-ac-dark rounded-md hover:bg-ac-gold/90 transition-colors"
            >
              Next
            </button>
          </div>
        </div>

        <!-- Step 3: Verification -->
        <div v-if="step === 3" class="space-y-4">
          <p class="text-ac-light text-sm">Enter the 6-digit code from your authenticator app</p>
          <input 
            type="text" 
            v-model="verificationCode"
            placeholder="Enter verification code"
            maxlength="6"
            class="w-full px-4 py-2 bg-ac-gray border border-ac-gold/30 rounded-md text-ac-light focus:outline-none focus:border-ac-gold"
          />
          <div class="flex justify-end space-x-3">
            <button 
              @click="closeEnableDialog"
              class="px-4 py-2 text-ac-light hover:text-ac-gold transition-colors"
            >
              Cancel
            </button>
            <button 
              @click="verify2FA"
              :disabled="!verificationCode || verificationCode.length !== 6"
              class="px-4 py-2 bg-ac-gold text-ac-dark rounded-md hover:bg-ac-gold/90 transition-colors disabled:opacity-50"
            >
              Verify
            </button>
          </div>
        </div>

        <div v-if="error" class="mt-4 text-red-500 text-sm text-center">
          {{ error }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useStore } from 'vuex'
import axios from 'axios'
import QRCodeVue from 'qrcode.vue'
import profileImage from '@/assets/profile.png'

const store = useStore()
const profile = ref({
  username: '',
  email: '',
  created_at: '',
  two_factor_enabled: false
})

const showEnableDialog = ref(false)
const step = ref(1)
const password = ref('')
const verificationCode = ref('')
const error = ref('')
const secret = ref('')
const qrUrl = ref('')

const characters = ref([
  {
    name: 'Helang_88',
    class: 'PAHLAWAN',
    rank: 42,
    score: 8750
  },
  {
    name: 'UtaraBentara',
    class: 'RAKSHAK',
    rank: 156,
    score: 6420
  },
  {
    name: 'Seri_Bendahara12',
    class: 'PEMANAH',
    rank: 89,
    score: 7340
  }
])

const highestRanked = computed(() => {
  return {
    class: 'Pahlawan',
    rank: 42
  }
})

const formatMemberSince = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString()
}

const fetchProfile = async () => {
  try {
    const response = await axios.get('http://localhost:8080/api/profile', {
      headers: {
        Authorization: `Bearer ${store.getters.token}`
      }
    })
    profile.value = response.data
  } catch (err) {
    console.error('Failed to fetch profile:', err)
  }
}

const closeEnableDialog = () => {
  showEnableDialog.value = false
  step.value = 1
  password.value = ''
  secret.value = ''
  qrUrl.value = ''
  verificationCode.value = ''
  error.value = ''
}

const verifyPassword = async () => {
  try {
    error.value = ''
    const response = await axios.post('http://localhost:8080/api/2fa/enable', 
      { password: password.value },
      {
        headers: {
          Authorization: `Bearer ${store.getters.token}`
        }
      }
    )
    secret.value = response.data.secret
    qrUrl.value = response.data.qr_url
    step.value = 2
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to verify password'
  }
}

const verify2FA = async () => {
  try {
    error.value = ''
    await axios.post('http://localhost:8080/api/2fa/verify', 
      { 
        code: verificationCode.value,
        secret: secret.value
      },
      {
        headers: {
          Authorization: `Bearer ${store.getters.token}`
        }
      }
    )
    await fetchProfile()
    closeEnableDialog()
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to verify code'
  }
}

const handleDisable2FA = async () => {
  try {
    await axios.post('http://localhost:8080/api/2fa/disable', {}, {
      headers: {
        Authorization: `Bearer ${store.getters.token}`
      }
    })
    await fetchProfile() // Fetch updated profile after disabling 2FA
  } catch (err) {
    console.error('Failed to disable 2FA:', err)
  }
}

onMounted(() => {
  fetchProfile()
})
</script>
