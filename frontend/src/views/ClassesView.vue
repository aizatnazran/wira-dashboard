<template>
  <div class="space-y-6 min-h-screen p-1">
    <!-- Hero Section -->
    <div class="relative bg-ac-dark py-16">
      <div class="bg-hero-pattern absolute inset-0 opacity-50"></div>
      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <h1 class="text-4xl font-cinzel text-ac-gold text-center mb-4">Character Classes</h1>
        <p class="text-ac-light text-center max-w-2xl mx-auto">Discover the unique abilities and playstyles of each character class in WIRA.</p>
      </div>
    </div>

    <!-- Classes Grid -->
    <div class="grid grid-cols-1 gap-6">
      <!-- Human Classes -->
      <div v-for="race in ['HUMAN', 'NUMAH']" :key="race" class="mb-12">
        <h2 class="text-3xl font-cinzel text-ac-gold mb-6">{{ race }}</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div v-for="classInfo in getClassesByRace(race)" :key="classInfo.id" 
               class="bg-ac-gray rounded-lg overflow-hidden border border-ac-gold/30 hover:border-ac-gold transition-colors duration-300 relative">
            <div class="p-6">
              <CombatTypeIcon :combat-type="classInfo.combat_type.toUpperCase()" />
              <h3 class="text-xl font-cinzel text-ac-gold mb-1">{{ classInfo.name }}</h3>
              <p class="text-sm text-ac-gold/70 mb-3">{{ classInfo.title }}</p>
              <p class="text-ac-light text-sm mb-4">{{ classInfo.description }}</p>
              
              <!-- Combat Type -->
              <div class="mb-4">
                <span class="text-sm font-semibold text-ac-gold">Type: </span>
                <span class="text-ac-light">{{ classInfo.combat_type }}</span>
              </div>
              
              <!-- Stats -->
              <div class="space-y-2">
                <div v-for="(stat, statName) in {
                  'Damage': classInfo.damage,
                  'Defense': classInfo.defense,
                  'Difficulty': classInfo.difficulty,
                  'Speed': classInfo.speed
                }" :key="statName" class="relative pt-1">
                  <div class="flex items-center justify-between">
                    <div>
                      <span class="text-xs font-semibold inline-block text-ac-gold">
                        {{ statName }}
                      </span>
                    </div>
                    <div>
                      <span class="text-xs font-semibold inline-block text-ac-light">
                        {{ stat }}/100
                      </span>
                    </div>
                  </div>
                  <div class="overflow-hidden h-2 text-xs flex rounded bg-ac-dark">
                    <div :style="{ width: stat + '%' }" 
                         class="shadow-none flex flex-col text-center whitespace-nowrap text-white justify-center bg-ac-gold">
                    </div>
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
import { ref, onMounted } from 'vue'
import { useToast } from 'vue-toastification'
import api from '@/api/config'
import CombatTypeIcon from '@/components/CombatTypeIcon.vue'

const toast = useToast()
const classes = ref([])
const loading = ref(true)
const error = ref(null)

const raceMap = {
  'HUMAN': 1,
  'NUMAH': 2
}

const fetchClasses = async () => {
  try {
    loading.value = true
    const response = await api.get('/api/classes')
    classes.value = response.data
  } catch (err) {
    console.error('Error fetching classes:', err)
    error.value = 'Failed to load classes'
    toast.error('Failed to load classes')
  } finally {
    loading.value = false
  }
}

const getClassesByRace = (raceName) => {
  const raceId = raceMap[raceName]
  return classes.value.filter(c => c.race_id === raceId)
}

onMounted(() => {
  fetchClasses()
})
</script>
