<template>
  <div class="min-h-screen p-4">
    <div class="container mx-auto">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold">Anime Musim Ini</h1>
        <p class="text-gray-400 mt-2">Temukan anime terbaru berdasarkan musim</p>
      </div>

      <!-- Season Selection -->
      <div class="flex flex-wrap gap-4 mb-8">
        <button
          v-for="season in seasons"
          :key="season.value"
          @click="selectSeason(season)"
          class="btn"
          :class="{ 'btn-primary': selectedSeason === season.value }"
        >
          {{ season.label }}
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center items-center min-h-[400px]">
        <span class="loading loading-spinner loading-lg text-primary"></span>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="alert alert-error">
        {{ error }}
      </div>

      <!-- Content Grid -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div v-for="content in contents" :key="content.id" 
             class="card bg-base-200 shadow-xl overflow-hidden hover:shadow-2xl transition-shadow">
          <figure class="relative aspect-[2/3]">
            <img 
              :src="getImageUrl(content.cover_image)" 
              :alt="content.title"
              class="w-full h-full object-cover"
              @error="handleImageError"
            />
            <div class="absolute top-2 right-2">
              <span class="badge badge-primary">{{ content.type }}</span>
            </div>
          </figure>
          <div class="card-body">
            <h2 class="card-title text-lg">{{ content.title }}</h2>
            <p class="text-sm text-gray-400 line-clamp-2">{{ content.description }}</p>
            <div class="flex flex-wrap gap-1 my-2">
              <span v-for="genre in content.genres" :key="genre.id" class="badge badge-outline badge-sm">
                {{ genre.name }}
              </span>
            </div>
            <div class="card-actions justify-between items-center mt-2">
              <div class="flex items-center">
                <span class="text-yellow-400">★</span>
                <span class="ml-1 text-sm">{{ content.rating || 0 }}/10</span>
              </div>
              <router-link :to="{ name: 'watch', params: { id: content.id }}" class="btn btn-primary btn-sm">
                Watch Now
              </router-link>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="contents.length > 0" class="mt-8 text-center">
        <p class="text-sm text-gray-400 mb-4">
          Showing {{ contents.length }} of {{ totalItems }} items
        </p>
        
        <!-- Load More -->
        <button 
          v-if="hasNextPage" 
          @click="loadMore" 
          class="btn btn-primary"
          :disabled="loading"
        >
          {{ loading ? 'Loading...' : 'Load More' }}
        </button>
      </div>

      <!-- No Results -->
      <div v-else-if="!loading" class="text-center py-12">
        <p class="text-gray-400">No content found for this season</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'

// State
const contents = ref([])
const loading = ref(false)
const error = ref(null)
const currentPage = ref(1)
const totalItems = ref(0)

// const seasons = [
//   { label: 'Winter 2024', value: 'winter_2024' },
//   { label: 'Spring 2024', value: 'spring_2024' },
//   { label: 'Summer 2024', value: 'summer_2024' },
//   { label: 'Fall 2024', value: 'fall_2024' }
// ]

const seasons = ref([])
const selectedSeason = ref('')

// const selectedSeason = ref('winter_2024')

// Methods
const getImageUrl = (cover_image) => {
  if (!cover_image) return
  const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  const cleanPath = cover_image.replace(/\\/g, '/').replace(/^media\//, '')
  return `${baseURL}/media/${cleanPath}`
}

const handleImageError = (e) => {
  return
}

const selectSeason = async (season) => {
  selectedSeason.value = season.value
  currentPage.value = 1
  await fetchContents(true)
}

const fetchContents = async (reset = false) => {
  try {
    loading.value = true
    error.value = null

    const params = {
      page: currentPage.value,
      pageSize: 12,
      season: selectedSeason.value
    }

    const response = await axios.get('/api/contents', { params })
    
    if (reset) {
      contents.value = response.data.contents
    } else {
      contents.value = [...contents.value, ...response.data.contents]
    }
    
    totalItems.value = response.data.total
  } catch (err) {
    console.error('Error fetching contents:', err)
    error.value = 'Failed to load contents'
  } finally {
    loading.value = false
  }
}

const fetchSeasons = async () => {
  try {
    const response = await axios.get('/api/seasons')
    console.log('check season get', response)
    seasons.value = response.data.map(season => ({
      label: `${season.name} ${season.year}`,
      value: season.id
    }))
    selectedSeason.value = response.data[0].id
  } catch (error) {
    console.error('Failed to fetch seasons:', error)
    showErrorMessage('Failed to fetch seasons')
  }
}

const hasNextPage = computed(() => {
  return currentPage.value * 12 < totalItems.value
})

const loadMore = () => {
  if (!loading.value && hasNextPage.value) {
    currentPage.value++
    fetchContents()
  }
}

// Lifecycle
onMounted(() => {
  fetchContents(true)
  fetchSeasons()
})
</script> 