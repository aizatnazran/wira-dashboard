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
      <div class="mb-6 flex flex-col md:flex-row gap-4">
        <!-- Search Input -->
        <div class="flex-1">
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
      </div>

      <!-- Rankings Table -->
      <div class="bg-ac-gray rounded-lg overflow-hidden">
        <div class="overflow-x-auto">
          <table class="min-w-full table-auto">
            <thead class="bg-ac-dark">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-ac-gold uppercase tracking-wider">Rank</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-ac-gold uppercase tracking-wider">Player</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-ac-gold uppercase tracking-wider">Class</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-ac-gold uppercase tracking-wider">Score</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-ac-dark">
              <template v-if="loading">
                <tr v-for="i in itemsPerPage" :key="i">
                  <td colspan="4" class="px-6 py-4">
                    <div class="animate-pulse flex space-x-4">
                      <div class="h-4 bg-ac-dark rounded w-3/4"></div>
                    </div>
                  </td>
                </tr>
              </template>
              <template v-else>
                <tr v-for="rank in rankings" :key="rank.char_id" class="hover:bg-ac-dark/50">
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-ac-light">{{ rank.rank }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-ac-light">{{ rank.username }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-ac-light">{{ rank.class_name }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-ac-light">{{ rank.reward_score }}</td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Pagination -->
      <div class="mt-6 flex justify-between items-center">
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
import { fetchWithAuth } from '../api'

// Pagination state
const currentPage = ref(1)
const itemsPerPage = ref(10)
const totalItems = ref(0)
const searchQuery = ref('')
const selectedClass = ref('all')

// Rankings data
const rankings = ref([])
const loading = ref(false)

// Classes grouped by race
const classGroups = {
  Human: ['PAHLAWAN', 'PENDEKAR', 'PEMANAH', 'PENGAMAL'],
  Numah: ['KSHATRIYA', 'VYAPARI', 'RAKSHAK', 'VAIDYA']
}

// Fetch rankings with pagination and filtering
const fetchRankings = async () => {
  loading.value = true
  try {
    const classParam = selectedClass.value === 'all' ? '' : `&class=${selectedClass.value}`
    const data = await fetchWithAuth(`/api/rankings?page=${currentPage.value}&limit=${itemsPerPage.value}&search=${searchQuery.value}${classParam}`)
    rankings.value = data.rankings
    totalItems.value = data.total
  } catch (error) {
    console.error('Error fetching rankings:', error)
  } finally {
    loading.value = false
  }
}

const debouncedSearch = debounce(() => {
  fetchRankings()
}, 300)

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