<template>
  <div class="min-h-screen flex items-center justify-center">
    <div class="card w-96 bg-base-200 shadow-xl">
      <div class="card-body">
        <h2 class="card-title justify-center mb-4">Login</h2>
        
        <!-- Alert Error -->
        <div v-if="authStore.error" class="alert alert-error mb-4">
          {{ authStore.error }}
        </div>

        <form @submit.prevent="handleLogin">
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
            <label class="label">
              <a href="#" class="label-text-alt link link-hover">Forgot password?</a>
            </label>
          </div>
          
          <div class="form-control mt-6">
            <button 
              type="submit" 
              class="btn btn-primary"
              :disabled="authStore.loading"
            >
              {{ authStore.loading ? 'Logging in...' : 'Login' }}
            </button>
          </div>
        </form>

        <div class="divider">OR</div>

        <div class="text-center">
          <router-link to="/register" class="link link-primary">
            Create new account
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

const email = ref('')
const password = ref('')

const handleLogin = async () => {
  try {
    await authStore.login(email.value, password.value)
    router.push('/')
  } catch (error) {
    console.error('Login failed:', error)
  }
}
</script> 