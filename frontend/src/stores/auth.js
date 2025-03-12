import { defineStore } from 'pinia'
import axios from 'axios'

export const useAuthStore = defineStore('auth', {
  state: () => {
    // Safely parse user data
    let user = null
    try {
      const userStr = localStorage.getItem('user')
      if (userStr) {
        user = JSON.parse(userStr)
      }
    } catch (e) {
      console.error('Failed to parse user data:', e)
      localStorage.removeItem('user')
    }

    return {
      token: localStorage.getItem('token') || null,
      user,
      loading: false,
      error: null
    }
  },

  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) => state.user?.role === 'admin'
  },

  actions: {
    async login(email, password) {
      try {
        this.loading = true
        this.error = null
        
        console.log('Login attempt:', { email, baseURL: axios.defaults.baseURL })
        
        const response = await axios.post('/api/auth/login', { 
          email, 
          password 
        }, {
          headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json'
          }
        })

        if (response.data?.token) {
          this.token = response.data.token
          this.user = response.data.user
          localStorage.setItem('token', this.token)
          localStorage.setItem('user', JSON.stringify(this.user))
          axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
          
          console.log('Login successful:', { user: this.user })
          return response.data
        } else {
          throw new Error('Invalid response format')
        }
      } catch (error) {
        console.error('Login error:', {
          message: error.message,
          response: error.response?.data,
          status: error.response?.status
        })
        this.error = error.response?.data?.error || error.message || 'Login failed'
        throw error
      } finally {
        this.loading = false
      }
    },

    async register(username, email, password) {
      try {
        console.log('Attempting registration...', { username, email })
        this.loading = true
        this.error = null
        
        const response = await axios.post('/api/auth/register', {
          username,
          email,
          password
        })
        
        console.log('Registration response:', response.data)
        return response.data
      } catch (error) {
        console.error('Registration error:', {
          message: error.message,
          response: error.response?.data,
          status: error.response?.status
        })
        this.error = error.response?.data?.error || error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    logout() {
      this.token = null
      this.user = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      delete axios.defaults.headers.common['Authorization']
    }
  }
})