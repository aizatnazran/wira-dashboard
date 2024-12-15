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
            <h1 class="text-3xl sm:text-4xl font-cinzel text-ac-gold mb-2 break-words">{{ user?.username }}</h1>
            <p class="text-ac-light mb-4 break-words">{{ user?.email }}</p>
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
                    <dd class="mt-1 text-sm text-ac-gold break-words">{{ user?.username }}</dd>
                  </div>
                  <div>
                    <dt class="text-sm font-medium text-ac-light">Email</dt>
                    <dd class="mt-1 text-sm text-ac-gold break-words">{{ user?.email }}</dd>
                  </div>
                  <div>
                    <dt class="text-sm font-medium text-ac-light">Member Since</dt>
                    <dd class="mt-1 text-sm text-ac-gold">{{ user?.created_at ? formatMemberSince(user.created_at) : '' }}</dd>
                  </div>
                </dl>
              </div>
            </div>

            <!-- 2FA Security Section -->
            <div>
              <div class="flex items-center justify-between mb-2">
                <h4 class="text-sm font-medium text-ac-gold">Two-Factor Authentication</h4>
                <span :class="user?.two_factor_enabled ? 'bg-green-600' : 'bg-red-600'" class="px-2 py-1 rounded text-xs text-white">
                  {{ user?.two_factor_enabled ? 'Enabled' : 'Disabled' }}
                </span>
              </div>
              <div class="mt-2 border border-ac-gold/30 rounded-lg p-4 bg-ac-dark">
                <div v-if="!user?.two_factor_enabled" class="space-y-3">
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
  <div class="mt-2 border border-ac-gold/30 rounded-lg p-4 bg-ac-dark flex flex-row gap-6">
    <!-- Left: Character Stats -->
    <div class="flex-1">
      <dl class="space-y-3">
        <div>
          <dt class="text-sm font-medium text-ac-light">Total Characters</dt>
          <dd class="mt-1 text-sm text-ac-gold">{{ characters.length }}</dd>
        </div>
        <div>
          <dt class="text-sm font-medium text-ac-light">Highest Ranked Class</dt>
          <dd class="mt-1 text-sm text-ac-gold">{{ highestRanked.class }} - Rank #{{ highestRanked.rank }}</dd>
        </div>
        <div>
          <dt class="text-sm font-medium text-ac-light">Average Score</dt>
          <dd class="mt-1 text-sm text-ac-gold">{{ averageScore }}</dd>
        </div>
        <div>
          <dt class="text-sm font-medium text-ac-light">Total Victories</dt>
          <dd class="mt-1 text-sm text-ac-gold">{{ playerStats.victories }}</dd>
        </div>
        <div>
          <dt class="text-sm font-medium text-ac-light">Win Rate</dt>
          <dd class="mt-1 text-sm text-ac-gold">{{ playerStats.winRate }}%</dd>
        </div>
        <div>
          <dt class="text-sm font-medium text-ac-light">Battle Rating</dt>
          <dd class="mt-1 text-sm text-ac-gold">{{ playerStats.battleRating }}/100</dd>
        </div>
      </dl>

      <!-- Pentagon Chart -->
      <div class="mt-6">
        <h5 class="text-sm font-medium text-ac-gold mb-4">Combat Attributes</h5>
        <div class="relative w-full h-48">
          <canvas ref="statsChart" class="w-full h-full"></canvas>
        </div>
      </div>
    </div>

    <!-- Right: 3D Model -->
    <div class="flex-1 flex items-center justify-center">
    <div 
    class="w-full h-[500px] sm:h-[400px] md:h-[500px] relative bg-transparent rounded-lg overflow-hidden"
    style="max-height: 80vh;"
  >
    <div ref="modelContainer" class="w-full h-full"></div>
  </div>
</div>
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
import { ref, onMounted, computed, watch, onUnmounted } from 'vue'
import { useStore } from 'vuex'
import { useToast } from 'vue-toastification'
import * as THREE from 'three'
import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls';
import api from '@/api/config'
import QRCodeVue from 'qrcode.vue'
import profileImage from '@/assets/profile.png'
import Chart from 'chart.js/auto'

let camera; // Declare camera globally
let renderer; // Declare renderer globally
let mixer; // Declare mixer globally

const store = useStore()
const toast = useToast()
const loading = ref(false)
const user = ref(null)
const statsChart = ref(null)
let chart = null

const fetchUserProfile = async () => {
  try {
    loading.value = true
    const response = await api.get('/api/profile')
    user.value = response.data
  } catch (error) {
    console.error('Error fetching profile:', error)
    toast.error('Failed to load profile data')
  } finally {
    loading.value = false
  }
}

const highestRanked = computed(() => {
  return {
    class: 'Pahlawan',
    rank: 42
  }
})

const characters = computed(() => {
  return [
    { name: 'Helang_88', class: 'PAHLAWAN', rank: 42, score: 8750 },
    { name: 'UtaraBentara', class: 'RAKSHAK', rank: 156, score: 6420 },
    { name: 'Seri_Bendahara12', class: 'PEMANAH', rank: 89, score: 7340 }
  ]
})

const averageScore = computed(() => {
  const totalScore = characters.value.reduce((acc, char) => acc + char.score, 0)
  return (totalScore / characters.value.length).toFixed(2)
})

const playerStats = computed(() => {
  return {
    victories: 100,
    winRate: 75,
    battleRating: 85,
    combatStats: {
      attack: 85,
      defense: 85,
      speed: 90,
      technique: 85,
      strategy: 80
    }
  }
})

const formatMemberSince = (dateStr) => {
  if (!dateStr) return 'N/A'
  return new Date(dateStr).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const modelContainer = ref(null)
const load3DModel = () => {
  // Scene, Camera, and Renderer
  const scene = new THREE.Scene();
  const camera = new THREE.PerspectiveCamera(
    35,
    modelContainer.value.offsetWidth / modelContainer.value.offsetHeight,
    0.1,
    1000
  );
  const renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true });
  renderer.setSize(modelContainer.value.offsetWidth, modelContainer.value.offsetHeight);
  modelContainer.value.appendChild(renderer.domElement);

  const light = new THREE.AmbientLight(0xffffff, 1.5);
  scene.add(light);

  const controls = new OrbitControls(camera, renderer.domElement); 
  controls.enableDamping = true; 
  controls.dampingFactor = 0.05;
  controls.rotateSpeed = 0.8; 
  controls.enableZoom = false; 
  controls.target.set(0, 0.5, 0); 
  controls.update();

  // Load GLB Model
  const loader = new GLTFLoader();
  loader.load('/models/swordman2.glb', (gltf) => {
    const model = gltf.scene;
    model.position.set(0, -2, 0);
    model.scale.set(2, 2, 2);
    scene.add(model);

    if (gltf.animations && gltf.animations.length > 0) {
      mixer = new THREE.AnimationMixer(model);

      const clip = gltf.animations[0];
      const action = mixer.clipAction(clip);

      action.setLoop(THREE.LoopRepeat, Infinity);
      action.clampWhenFinished = false;
      action.timeScale = 0.5;
      action.play();

      mixer.addEventListener('loop', () => {
        action.time = 0;
      });

      const clock = new THREE.Clock();
      const animate = () => {
        requestAnimationFrame(animate);

        const delta = clock.getDelta();
        mixer.update(delta);

        if (action.time > 3.3) {
          action.time = 0;
        }

        renderer.render(scene, camera);
      };
      animate();
    }
  }, undefined, (error) => {
    console.error('Error loading 3D model:', error);
  });

  camera.position.set(0, 1, -13);
  camera.lookAt(0, 0.5, 0);

  // Animation Loop
  const clock = new THREE.Clock();
  const animate = () => {
    requestAnimationFrame(animate);

    if (mixer) {
      mixer.update(clock.getDelta());
    }

    controls.update(); 
    renderer.render(scene, camera);
  };
  animate();
};

const showEnableDialog = ref(false)
const qrUrl = ref('')
const secret = ref('')
const password = ref('')
const verificationCode = ref('')
const error = ref('')
const step = ref(1)

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
    const response = await api.post('/api/2fa/enable', 
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
    await api.post('/api/2fa/verify', 
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
    await fetchUserProfile()
    closeEnableDialog()
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to verify code'
  }
}

const handleDisable2FA = async () => {
  try {
    await api.post('/api/2fa/disable', {}, {
      headers: {
        Authorization: `Bearer ${store.getters.token}`
      }
    })
    await fetchUserProfile() 
  } catch (err) {
    console.error('Failed to disable 2FA:', err)
  }
}

const initChart = () => {
  if (statsChart.value) {
    const ctx = statsChart.value.getContext('2d')
    const stats = playerStats.value.combatStats
    
    if (chart) {
      chart.destroy()
    }

    chart = new Chart(ctx, {
      type: 'radar',
      data: {
        labels: ['Attack', 'Defense', 'Speed', 'Technique', 'Strategy'],
        datasets: [{
          label: 'Combat Stats',
          data: [stats.attack, stats.defense, stats.speed, stats.technique, stats.strategy],
          backgroundColor: 'rgba(255, 215, 0, 0.2)',
          borderColor: 'rgba(255, 215, 0, 0.8)',
          borderWidth: 2,
          pointBackgroundColor: 'rgba(255, 215, 0, 1)',
        }]
      },
      options: {
        scales: {
          r: {
            angleLines: {
              color: 'rgba(255, 215, 0, 0.1)'
            },
            grid: {
              color: 'rgba(255, 215, 0, 0.1)'
            },
            pointLabels: {
              color: '#D4AF37'
            },
            ticks: {
              color: '#D4AF37',
              backdropColor: 'transparent'
            }
          }
        },
        plugins: {
          legend: {
            display: false
          }
        }
      }
    })
  }
}

const handleResize = () => {
  if (renderer && camera && modelContainer.value) {
    const width = modelContainer.value.offsetWidth;
    const height = modelContainer.value.offsetHeight;

    renderer.setSize(width, height); 
    camera.aspect = width / height; 
    camera.updateProjectionMatrix();
  }
};

window.addEventListener('resize', handleResize);

onMounted(() => {
  load3DModel()
  fetchUserProfile()
  watch(statsChart, () => {
    if (statsChart.value) {
      initChart()
    }
  }, { immediate: true })
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});
</script>
<style scoped>
.bg-transparent {
  background: transparent;
}
</style>