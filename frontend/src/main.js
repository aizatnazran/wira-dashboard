import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import axios from 'axios'
import Toast from 'vue-toastification'
import 'vue-toastification/dist/index.css'
import './assets/tailwind.css'
import '@fortawesome/fontawesome-free/css/all.css';

// Configure axios defaults
axios.defaults.baseURL = 'http://localhost:8080'
axios.defaults.withCredentials = true

const app = createApp(App)

const toastOptions = {
  position: "top-right",
  timeout: 3000,
  closeOnClick: true,
  pauseOnFocusLoss: true,
  pauseOnHover: true,
  draggable: true,
  draggablePercent: 0.6,
  showCloseButtonOnHover: false,
  hideProgressBar: true,
  closeButton: "button",
  icon: true,
  rtl: false,
  toastClassName: "bg-ac-gray border border-ac-gold/30 text-ac-light",
}

app.use(router)
app.use(store)
app.use(Toast, toastOptions)
app.mount('#app')
