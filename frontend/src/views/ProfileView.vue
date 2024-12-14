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
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { fetchWithAuth } from '../api'
import profileImage from '@/assets/profile.png'

const profile = ref({
  username: '',
  email: '',
  created_at: ''
})

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
  return characters.value.reduce((prev, current) => {
    return (prev.rank < current.rank) ? prev : current
  })
})

const formatMemberSince = (dateStr) => {
  if (!dateStr) return 'N/A'
  const date = new Date(dateStr)
  return `${date.toLocaleString('default', { month: 'long' })} ${date.getFullYear()}`
}

const fetchProfile = async () => {
  try {
    const data = await fetchWithAuth('/api/profile')
    profile.value = data
  } catch (error) {
    console.error('Error fetching profile:', error)
  }
}

onMounted(() => {
  fetchProfile()
})
</script>