<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'
import axios from 'axios'
import {
  UserIcon,
  Cog6ToothIcon,
  ArrowRightOnRectangleIcon
} from '@heroicons/vue/24/outline'

const router = useRouter()
const authStore = useAuthStore()

// Add debug logging
onMounted(async () => {
  console.log('App mounted')
  console.log('Base URL:', axios.defaults.baseURL)
  console.log('Auth status:', authStore.isAuthenticated)
  
  // Test API connection
  try {
    const response = await axios.get('/api/contents')
    console.log('API Test Response:', response.data)
  } catch (error) {
    console.error('API Test Error:', {
      message: error.message,
      response: error.response?.data,
      status: error.response?.status,
      headers: error.response?.headers
    })
  }
})

// Add logout handler
const handleLogout = async () => {
  console.log('Logging out...')
  await authStore.logout()
  router.push('/login')
}
</script>

<template>
  <div class="min-h-screen bg-base-100 text-white">
    <!-- Navigation -->
    <nav class="bg-base-200 shadow-lg">
      <div class="container mx-auto px-4">
        <div class="flex items-center justify-between h-16">
          <!-- Logo -->
          <router-link to="/" class="flex items-center space-x-2">
            <span class="text-xl font-bold text-primary">AnimePortal</span>
          </router-link>

          <!-- Navigation Links -->
          <div class="hidden md:flex items-center space-x-4">
            <router-link
              to="/browse"
              class="px-3 py-2 rounded-md text-sm font-medium hover:bg-base-300"
              :class="{ 'text-primary': $route.path === '/browse' }"
            >
              Browse
            </router-link>
            <router-link
              to="/season"
              class="px-3 py-2 rounded-md text-sm font-medium hover:bg-base-300"
              :class="{ 'text-primary': $route.path === '/Season' }"
            >
              Musim
            </router-link>
          </div>

          <!-- User Menu -->
          <div class="flex items-center space-x-4">
            <template v-if="authStore.isAuthenticated">
              <div class="dropdown dropdown-end">
                <label tabindex="0" class="btn btn-ghost btn-circle avatar">
                  <div class="w-10 rounded-full">
                    <img
                      :src="authStore.user?.avatar || 'https://ui-avatars.com/api/?name=' + authStore.user?.username"
                      :alt="authStore.user?.username"
                    >
                  </div>
                </label>
                <ul tabindex="0" class="mt-3 p-2 shadow menu menu-sm dropdown-content bg-base-200 rounded-box w-52" style="z-index: 999;">
                 <li>
                    <router-link to="/profile">
                      <UserIcon class="w-4 h-4" />
                      Profile
                    </router-link>
                  </li>
                  <li v-if="authStore.isAdmin">
                    <router-link to="/admin">
                      <Cog6ToothIcon class="w-4 h-4" />
                      Admin
                    </router-link>
                  </li>
                  <li>
                    <button @click="handleLogout">
                      <ArrowRightOnRectangleIcon class="w-4 h-4" />
                      Logout
                    </button>
                  </li>
                </ul>
              </div>
            </template>
            <template v-else>
              <router-link
                to="/login"
                class="btn btn-primary btn-sm"
              >
                Login
              </router-link>
              <router-link
                to="/register"
                class="btn btn-ghost btn-sm"
              >
                Register
              </router-link>
            </template>
          </div>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <router-view v-slot="{ Component }">
      <transition name="fade" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </div>
</template>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
