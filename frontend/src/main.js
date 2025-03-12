import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import axios from 'axios'

// Configure axios
axios.defaults.baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
axios.defaults.headers.common['Content-Type'] = 'application/json'
axios.defaults.withCredentials = true
axios.defaults.timeout = 30000 // 30 detik timeout

// Add axios interceptors for better error handling and token management
axios.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      console.log('Token being sent:', token)
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

axios.interceptors.response.use(
  response => response,
  error => {
    console.error('Axios Error:', {
      message: error.message,
      response: error.response?.data,
      status: error.response?.status,
      config: error.config
    })
    
    if (error.response?.status === 403) {
      // Cek token dan role
      const token = localStorage.getItem('token')
      const user = JSON.parse(localStorage.getItem('user') || '{}')
      console.log('Current token:', token)
      console.log('Current user:', user)
      
      if (!token || !user.role || user.role !== 'admin') {
        console.log('Invalid admin access, redirecting to login...')
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        router.push('/login')
      }
    }
    
    if (error.response?.status === 401) {
      console.log('Token expired, redirecting to login...')
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      router.push('/login')
    }
    
    return Promise.reject(error)
  }
)

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
