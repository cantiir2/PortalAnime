<template>
  <div class="min-h-screen p-4">
    <div class="container mx-auto">
      <!-- Search and Filters -->
      <div class="mb-8 space-y-4">
        <!-- Search Bar -->
        <div class="form-control">
          <div class="input-group">
            <input
              type="text"
              v-model="searchQuery"
              placeholder="Search anime..."
              class="input input-bordered w-full"
              @input="handleSearch"
            />
            <button class="btn btn-square">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Filters -->
        <div class="flex flex-wrap gap-4">
          <!-- Genre Filter -->
          <div class="dropdown">
            <label tabindex="0" class="btn m-1">
              {{ selectedGenre?.name || 'Genre' }}
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </label>
            <ul tabindex="0" class="dropdown-content menu p-2 shadow bg-base-200 rounded-box w-52 max-h-60 overflow-y-auto" style="z-index: 999;">
              <li><button @click="handleGenreChange(null)">All Genres</button></li>
              <li v-for="genre in genres" :key="genre.id">
                <button @click="handleGenreChange(genre)">{{ genre.name }}</button>
              </li>
            </ul>
          </div>

          <!-- Category Filter -->
          <div class="dropdown">
            <label tabindex="0" class="btn m-1">
              {{ selectedCategory?.name || 'Category' }}
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </label>
            <ul tabindex="0" class="dropdown-content menu p-2 shadow bg-base-200 rounded-box w-52 max-h-60 overflow-y-auto" style="z-index: 999;">
              <li><button @click="handleCategoryChange(null)">All Categories</button></li>
              <li v-for="category in categories" :key="category.id">
                <button @click="handleCategoryChange(category)">{{ category.name }}</button>
              </li>
            </ul>
          </div>
        </div>
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
        <div v-for="content in contents" :key="content.id" class="card bg-base-200 shadow-xl overflow-hidden hover:shadow-2xl transition-shadow">
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

      <!-- Pagination Info -->
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
        <p class="text-gray-400">No content found</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import debounce from 'lodash/debounce'

// State
const contents = ref([])
const genres = ref([])
const categories = ref([])
const selectedGenre = ref(null)
const selectedCategory = ref(null)
const searchQuery = ref('')
const currentPage = ref(1)
const totalItems = ref(0)
const loading = ref(false)
const error = ref(null)

// Methods
const getImageUrl = (cover_image) => {
  if (!cover_image) return '/placeholder.png'
  const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  const cleanPath = cover_image.replace(/\\/g, '/').replace(/^media\//, '')
  return `${baseURL}/media/${cleanPath}`
}

// const handleImageError = (e) => {
//   console.error('Error loading image:', e.target.src)
//   e.target.src = '/placeholder.png'
// }

const fetchGenres = async () => {
  try {
    const response = await axios.get('/api/genres')
    genres.value = response.data
  } catch (err) {
    console.error('Error fetching genres:', err)
  }
}

const fetchCategories = async () => {
  try {
    const response = await axios.get('/api/categories')
    categories.value = response.data
  } catch (err) {
    console.error('Error fetching categories:', err)
  }
}

const fetchContents = async (page = 1, reset = false) => {
  try {
    loading.value = true
    error.value = null

    const params = {
      page,
      pageSize: 12
    }

    if (searchQuery.value) {
      params.q = searchQuery.value
    }

    if (selectedGenre.value) {
      const response = await axios.get(`/api/contents/genre/${selectedGenre.value.id}`, { params })
      if (reset) {
        contents.value = response.data.contents
      } else {
        contents.value = [...contents.value, ...response.data.contents]
      }
      totalItems.value = response.data.total
      currentPage.value = page
      return
    }

    if (selectedCategory.value) {
      const response = await axios.get(`/api/contents/category/${selectedCategory.value.id}`, { params })
      if (reset) {
        contents.value = response.data.contents
      } else {
        contents.value = [...contents.value, ...response.data.contents]
      }
      totalItems.value = response.data.total
      currentPage.value = page
      return
    }

    const response = await axios.get('/api/contents', { params })
    if (reset) {
      contents.value = response.data.contents
    } else {
      contents.value = [...contents.value, ...response.data.contents]
    }
    totalItems.value = response.data.total
    currentPage.value = page
  } catch (err) {
    console.error('Error fetching contents:', err)
    error.value = 'Failed to load contents'
  } finally {
    loading.value = false
  }
}

const handleSearch = debounce(() => {
  currentPage.value = 1
  fetchContents(1, true)
}, 300)

const handleGenreChange = (genre) => {
  selectedGenre.value = genre
  selectedCategory.value = null // Reset category when genre changes
  currentPage.value = 1
  fetchContents(1, true)
}

const handleCategoryChange = (category) => {
  selectedCategory.value = category
  selectedGenre.value = null // Reset genre when category changes
  currentPage.value = 1
  fetchContents(1, true)
}

const hasNextPage = computed(() => {
  const pageSize = 12
  return currentPage.value * pageSize < totalItems.value
})

const loadMore = () => {
  if (!loading.value && hasNextPage.value) {
    fetchContents(currentPage.value + 1)
  }
}

// Lifecycle
onMounted(async () => {
  await Promise.all([
    fetchGenres(),
    fetchCategories(),
    fetchContents(1, true)
  ])
})
</script> 