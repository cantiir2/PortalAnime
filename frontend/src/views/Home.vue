<template>
  <div class="min-h-screen">
    <!-- Debug Info (development only) -->
    <div v-if="error" class="bg-red-500 text-white p-4">
      Error: {{ error }}
    </div>

    <!-- Hero Section -->
    <section class="hero min-h-[60vh] bg-base-200">
      <div class="hero-content text-center">
        <div class="max-w-md">
          <h1 class="text-5xl font-bold">{{ appName }}</h1>
          <p class="py-6">{{ appDescription }}</p>
          <router-link to="/browse" class="btn btn-primary">Start Watching</router-link>
        </div>
      </div>
    </section>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <span class="loading loading-spinner loading-lg text-primary"></span>
    </div>

    <!-- Content Sections -->
    <template v-else>
      <!-- Featured Content -->
      <section class="py-12">
        <div class="container mx-auto px-4">
          <h2 class="text-2xl font-bold mb-6">Featured Anime</h2>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <div v-for="content in featuredContent" :key="content.id" 
                 class="card bg-base-200 shadow-xl">
              <figure>
                <img :src="getImageUrl(content.cover_image)" :alt="content.title" 
                     class="w-full h-48 object-cover" />
              </figure>
              <div class="card-body">
                <h3 class="card-title">{{ content.title }}</h3>
                <p class="text-sm">{{ content.description }}</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Latest Updates -->
      <section class="py-12 bg-base-200">
        <div class="container mx-auto px-4">
          <h2 class="text-2xl font-bold mb-6">Latest Updates</h2>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <div v-for="content in latestUpdates" :key="content.id" 
                 class="card bg-base-100 shadow-xl">
              <figure>
                <img :src="getImageUrl(content.cover_image)" :alt="content.title" 
                     class="w-full h-48 object-cover" />
              </figure>
              <div class="card-body">
                <h3 class="card-title">{{ content.title }}</h3>
                <p class="text-sm">{{ content.description }}</p>
              </div>
            </div>
          </div>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import axios from 'axios'

const appName = import.meta.env.VITE_APP_NAME || 'AnimePortal'
const appDescription = import.meta.env.VITE_APP_DESCRIPTION || 'Your Ultimate Anime Streaming Platform'

const featuredContent = ref([])
const latestUpdates = ref([])
const loading = ref(true)
const error = ref(null)

const getImageUrl = (cover_image) => {
  if (!cover_image) return '/placeholder.png'
  const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  const cleanPath = cover_image.replace(/\\/g, '/').replace(/^media\//, '')
  return `${baseURL}/media/${cleanPath}`
}

onMounted(async () => {
  console.log('Home view mounted')
  try {
    // Fetch featured content
    const featuredResponse = await axios.get('/api/contents', {
      params: { page: 1, pageSize: 4 }
    })
    console.log('Featured content response:', featuredResponse.data)
    featuredContent.value = featuredResponse.data.contents

    // Fetch latest updates
    const latestResponse = await axios.get('/api/contents', {
      params: { page: 1, pageSize: 4, sort: 'created_at' }
    })
    console.log('Latest updates response:', latestResponse.data)
    latestUpdates.value = latestResponse.data.contents
  } catch (err) {
    console.error('Home view error:', {
      message: err.message,
      response: err.response?.data,
      status: err.response?.status
    })
    error.value = err.message
  } finally {
    loading.value = false
  }
})
</script> 