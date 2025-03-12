<template>
  <div class="min-h-screen flex items-center justify-center">
    <div class="card w-96 bg-base-200 shadow-xl">
      <div class="card-body">
        <h2 class="card-title justify-center mb-4">Register</h2>
        <form @submit.prevent="handleRegister">
          <div class="form-control">
            <label class="label">
              <span class="label-text">Username</span>
            </label>
            <input 
              type="text" 
              v-model="username" 
              placeholder="Enter your username" 
              class="input input-bordered" 
              required
            />
          </div>
          <div class="form-control">
            <label class="label">
              <span class="label-text">Email</span>
            </label>
            <input 
              type="email" 
              v-model="email" 
              placeholder="Enter your email" 
              class="input input-bordered" 
              required
            />
          </div>
          <div class="form-control">
            <label class="label">
              <span class="label-text">Password</span>
            </label>
            <input 
              type="password" 
              v-model="password" 
              placeholder="Enter your password" 
              class="input input-bordered" 
              required
            />
          </div>
          <div class="form-control mt-6">
            <button type="submit" class="btn btn-primary" :disabled="loading">
              {{ loading ? 'Loading...' : 'Register' }}
            </button>
          </div>
        </form>
        <div class="text-center mt-4">
          <router-link to="/login" class="link link-primary">
            Already have an account? Login
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const email = ref('')
const password = ref('')
const loading = ref(false)

const handleRegister = async () => {
  try {
    loading.value = true
    await authStore.register(username.value, email.value, password.value)
    router.push('/login')
  } catch (error) {
    console.error('Registration failed:', error)
  } finally {
    loading.value = false
  }
}
</script> 