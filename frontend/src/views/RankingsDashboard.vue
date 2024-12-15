<template>
  <div class="min-h-screen p-4">
    <div class="max-w-7xl mx-auto">
      <!-- Hero Section -->
      <div class="relative bg-ac-dark py-16">
        <div class="bg-hero-pattern absolute inset-0 opacity-50"></div>
        <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h1 class="text-4xl font-cinzel text-ac-gold text-center mb-4">WIRA Rankings</h1>
          <p class="text-ac-light text-center max-w-2xl mx-auto">Explore the rankings of the most skilled warriors across eight unique classes in the world of WIRA.</p>
        </div>
      </div>

<!-- Search and Filter Section -->
<div class="mt-8 mb-6 flex flex-col md:flex-row gap-4 items-center justify-between">
  <!-- Search Input -->
  <div class="flex-1 relative">
    <input
      v-model="searchQuery"
      type="text"
      placeholder="Search players..."
      class="w-full px-4 py-2 rounded-lg bg-ac-gray text-ac-light border border-ac-gold focus:outline-none focus:ring-2 focus:ring-ac-gold"
    />
  </div>

  <!-- Class Filter -->
  <div class="w-full md:w-64">
    <select
      v-model="selectedClass"
      class="w-full px-4 py-2 rounded-lg bg-ac-gray text-ac-light border border-ac-gold focus:outline-none focus:ring-2 focus:ring-ac-gold"
    >
      <option value="all">All Classes</option>
      <optgroup label="Human">
        <option v-for="className in classGroups.Human" :key="className" :value="className">
          {{ className }}
        </option>
      </optgroup>
      <optgroup label="Numah">
        <option v-for="className in classGroups.Numah" :key="className" :value="className">
          {{ className }}
        </option>
      </optgroup>
    </select>
  </div>

  <!-- Reset Cache Button -->
  <button 
  @click="handleResetCache" 
  class="flex items-center justify-center w-10 h-10 rounded-lg bg-ac-gray border border-ac-gold text-white hover:text-ac-gold transition-colors duration-300"
  title="Reset Cache"
>
  <i 
    class="fas fa-redo transform transition-transform duration-300 hover:rotate-180"
  ></i>
</button>
</div>


      <!-- Rankings Table -->
      <div class="mt-8 flex flex-col">
        <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
          <div class="inline-block min-w-full py-2 align-middle md:px-6 lg:px-8">
            <div class="overflow-hidden shadow ring-1 ring-black ring-opacity-5 md:rounded-lg">
              <table class="min-w-full divide-y divide-ac-gold/30">
                <thead class="bg-ac-gray">
                  <tr>
                    <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-ac-gold sm:pl-6">RANK</th>
                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-ac-gold">PLAYER</th>
                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-ac-gold">CLASS</th>
                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-ac-gold">SCORE</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-ac-gold/30 bg-ac-dark">
                  <template v-if="loading">
                    <tr v-for="i in itemsPerPage" :key="i">
                      <td colspan="4" class="px-6 py-4">
                        <div class="animate-pulse flex space-x-4">
                          <div class="h-4 bg-ac-gold/20 rounded w-3/4"></div>
                        </div>
                      </td>
                    </tr>
                  </template>
                  <template v-else-if="rankings && rankings.length > 0">
                    <tr
                      v-for="(ranking, index) in rankings"
                      :key="ranking.char_id"
                      :class="{'bg-ac-gold/10': (currentPage - 1) * itemsPerPage + index + 1 <= 3}"
                    >
                      <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm text-ac-gold sm:pl-6">
                        <img 
                          v-if="(currentPage - 1) * itemsPerPage + index + 1 === 1" 
                          src="@/assets/medals/gold-medal.png" 
                          alt="Gold Medal" 
                          class="w-8 h-8 inline-block"
                        />
                        <img 
                          v-else-if="(currentPage - 1) * itemsPerPage + index + 1 === 2" 
                          src="@/assets/medals/silver-medal.png" 
                          alt="Silver Medal" 
                          class="w-8 h-8 inline-block"
                        />
                        <img 
                          v-else-if="(currentPage - 1) * itemsPerPage + index + 1 === 3" 
                          src="@/assets/medals/bronze-medal.png" 
                          alt="Bronze Medal" 
                          class="w-8 h-8 inline-block"
                        />
                        <span v-else>
                          {{ (currentPage - 1) * itemsPerPage + index + 1 }}
                        </span>
                      </td>
                      <td class="whitespace-nowrap px-3 py-4 text-sm text-ac-light">
                        {{ ranking.username }}
                      </td>
                      <td class="whitespace-nowrap px-3 py-4 text-sm text-ac-light">
                        {{ ranking.class_name }}
                      </td>
                      <td class="whitespace-nowrap px-3 py-4 text-sm text-ac-light">
                        {{ ranking.reward_score }}
                      </td>
                    </tr>
                  </template>
                  <tr v-else>
                    <td colspan="4" class="px-3 py-8 text-center text-ac-light">
                      <div class="flex flex-col items-center justify-center space-y-2">
                        <i class="fas fa-search text-ac-gold text-2xl mb-2"></i>
                        <p class="text-lg font-medium text-ac-gold">No results found</p>
                        <p class="text-sm text-ac-light" v-if="searchQuery">
                          No matches found for "{{ searchQuery }}". Try adjusting your search.
                        </p>
                        <p class="text-sm text-ac-light" v-else>
                          No rankings available at the moment.
                        </p>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="rankings && rankings.length > 0" class="flex items-center justify-between border-t border-ac-gold/30 bg-ac-dark px-4 py-3 sm:px-6">
        <div class="text-sm text-ac-light">
          Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, totalItems) }} of {{ totalItems }} results
        </div>
        <div class="flex gap-2">
          <button
            @click="currentPage--"
            :disabled="currentPage === 1"
            class="px-4 py-2 rounded-lg bg-ac-dark text-ac-gold disabled:opacity-50 disabled:cursor-not-allowed hover:bg-ac-dark/80"
          >
            Previous
          </button>
          <button
            @click="currentPage++"
            :disabled="currentPage >= totalPages"
            class="px-4 py-2 rounded-lg bg-ac-dark text-ac-gold disabled:opacity-50 disabled:cursor-not-allowed hover:bg-ac-dark/80"
          >
            Next
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import debounce from 'lodash/debounce'
import api from '@/api/config'
import { useToast } from 'vue-toastification'
import Swal from 'sweetalert2'

const toast = useToast()
const searchQuery = ref('')
const selectedClass = ref('all')
const currentPage = ref(1)
const itemsPerPage = ref(10)
const totalItems = ref(0)
const rankings = ref([]) 
const loading = ref(false)

const classGroups = {
  Human: ['PAHLAWAN', 'PENDEKAR', 'PEMANAH', 'PENGAMAL'],
  Numah: ['KSHATRIYA', 'VYAPARI', 'RAKSHAK', 'VAIDYA']
}

const fetchRankings = async () => {
  try {
    loading.value = true
    const params = {
      page: currentPage.value,
      limit: itemsPerPage.value,
      search: searchQuery.value,
      class: selectedClass.value === 'all' ? undefined : selectedClass.value
    }

    const response = await api.get('/api/rankings', { params })
    rankings.value = response.data.rankings || []
    totalItems.value = response.data.total || 0
  } catch (error) {
    console.error('Error fetching rankings:', error)
    rankings.value = []
    totalItems.value = 0
  } finally {
    loading.value = false
  }
}

const debouncedSearch = debounce(() => {
  fetchRankings()
}, 300)

const handleResetCache = async () => {
  const result = await Swal.fire({
    title: 'Clear Cache',
    text: 'Are you sure you want to clear the cache?',
    icon: 'warning',
    showCancelButton: true,
    confirmButtonText: 'Yes, clear it',
    cancelButtonText: 'Cancel',
    background: '#2A2A2A',
    color: '#F5F5F5',
    confirmButtonColor: '#C6A875',
    cancelButtonColor: '#4B5563'
  })

  if (result.isConfirmed) {
    try {
      const response = await api.post('/api/cache/clear')

      if (response.status === 200) {
        toast.success('Cache cleared successfully', {
          position: 'top-right',
          timeout: 5000,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
        })
      } else {
        toast.error('Failed to clear cache')
      }
    } catch (error) {
      console.error('Error clearing cache:', error)
      toast.error('Failed to clear cache')
    }
  }
}

// Trigger rankings refresh
watch([currentPage, itemsPerPage], fetchRankings)
watch(selectedClass, () => {
  currentPage.value = 1 
  fetchRankings()
})
watch(searchQuery, () => {
  currentPage.value = 1 
  debouncedSearch()
})

const totalPages = computed(() => Math.ceil(totalItems.value / itemsPerPage.value))

onMounted(fetchRankings)
</script>
