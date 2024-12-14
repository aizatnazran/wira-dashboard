import { createRouter, createWebHistory } from 'vue-router'
import RankingsDashboard from '../views/RankingsDashboard.vue'
import store from '../store'

const routes = [
  {
    path: '/',
    name: 'Rankings',
    component: RankingsDashboard,
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/LoginView.vue'),
    meta: { requiresGuest: true }
  },
  {
    path: '/signup',
    name: 'Signup',
    component: () => import('../views/SignupView.vue'),
    meta: { requiresGuest: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('../views/ProfileView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/classes',
    name: 'Classes',
    component: () => import('../views/ClassesView.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

router.beforeEach((to, from, next) => {
  const isAuthenticated = store.getters.isAuthenticated

  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login')
  } else if (to.meta.requiresGuest && isAuthenticated) {
    next('/')
  } else {
    next()
  }
})

export default router
