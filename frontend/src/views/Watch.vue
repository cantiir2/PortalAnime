<!-- Watch.vue -->
<template>
  <div class="min-h-screen p-4">
    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <span class="loading loading-spinner loading-lg text-primary"></span>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="flex flex-col items-center justify-center py-12">
      <div class="alert alert-error mb-4">
        {{ error }}
      </div>
      <div class="flex gap-4">
        <button @click="loadContent" class="btn btn-primary">
          Try Again
        </button>
        <router-link to="/" class="btn btn-ghost">
          Go Home
        </router-link>
      </div>
    </div>

    <!-- Content -->
    <div v-else-if="content" class="container mx-auto">
      <!-- Header Info -->
      <div class="container mx-auto px-4 py-4">
        <h1 class="text-2xl font-bold">{{ content?.title }}</h1>
        <div class="text-sm text-gray-400 mb-2">
          {{ currentEpisode ? `Season ${currentEpisode.season_number} Episode ${currentEpisode.episode_number}` : '' }}
          {{ currentEpisode?.title ? `- ${currentEpisode.title}` : '' }}
        </div>
      </div>

      <!-- Video Player Section -->
      <div class="w-full bg-black">
        <div class="container mx-auto px-4">
          <!-- Video Player -->
          <div class="aspect-video relative">
            <VideoPlayer
              v-if="currentVideo"
              :source="currentVideo.url"
              :type="currentVideo.type"
              :thumbnail-url="currentVideo.thumbnailUrl"
              :qualities="currentVideo.qualities"
              :initial-server="currentServer"
              :stream-links="currentEpisode?.stream_links || content?.stream_links || []"
              @progress="handleProgress"
              @ended="handleEnded"
              @server-change="handleServerChange"
            />
            
            <!-- Loading State -->
            <div v-else class="w-full h-full flex items-center justify-center bg-base-300">
              <span class="loading loading-spinner loading-lg text-primary"></span>
            </div>
          </div>
          
          <!-- Stream Options -->
          <div class="bg-base-200 p-4 rounded-b-lg">
            <!-- Server Selection -->
            <!-- <div class="mb-4">
              <h3 class="text-lg font-semibold mb-2">Server</h3>
              <div class="flex flex-wrap gap-2">
                <button
                  v-for="server in streamServers"
                  :key="server.value"
                  @click="selectServer(server)"
                  class="btn"
                  :class="{ 'btn-primary': currentServer?.value === server.value }"
                >
                  {{ server.label }}
                </button>
              </div>
            </div> -->
            
            <!-- Quality Selection -->
            <div v-if="currentServer">
              <h3 class="text-lg font-semibold mb-2">Quality</h3>
              <div class="flex flex-wrap gap-2">
                <button
                  v-for="quality in currentServer.qualities"
                  :key="quality.label"
                  @click="selectQuality(quality)"
                  class="btn"
                  :class="{ 'btn-secondary': currentQuality?.label === quality.label, 'btn-outline': currentQuality?.label !== quality.label }"
                >
                  {{ quality.label }}
                </button>
              </div>
            </div>
            
            <!-- Download Links -->
            <div v-if="downloadLinks.length > 0" class="mt-4">
              <h3 class="text-lg font-semibold mb-2">Download</h3>
              <!-- Group by Quality -->
              <div v-for="quality in downloadQualities" :key="quality" class="mb-4">
                <h4 class="text-base font-medium text-gray-400 mb-2">Quality: {{ quality }}</h4>
                <div class="flex flex-wrap gap-2">
                  <a
                    v-for="link in getDownloadLinksByQuality(quality)"
                    :key="link.id"
                    :href="isValidUrl(link.url) ? link.url : `https://${link.url}`"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="btn btn-outline hover:btn-primary transition-colors duration-200"
                  >
                    <span class="mr-2">{{ link.name }}</span>
                    <span class="text-xs opacity-75">{{ formatFileSize(link.size) }}</span>
                  </a>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Content Information -->
      <div class="container mx-auto px-4 py-6">
        <div class="flex flex-col md:flex-row gap-6">
          <!-- Main Content -->
          <div class="flex-1">
            <div class="mb-6">
              <h2 class="text-xl font-semibold mb-2">Synopsis</h2>
              <p class="text-gray-300">{{ content?.description }}</p>
            </div>

            <!-- Episodes List -->
            <div v-if="episodes.length > 0" class="mt-6">
              <h2 class="text-xl font-semibold mb-4">Episodes</h2>
              
              <!-- Episode Navigation -->
              <div class="flex justify-between mb-4">
                <button 
                  v-if="previousEpisode"
                  @click="playEpisode(previousEpisode)"
                  class="btn btn-primary"
                >
                  Previous Episode
                </button>
                <div></div> <!-- Spacer -->
                <button 
                  v-if="nextEpisode"
                  @click="playEpisode(nextEpisode)"
                  class="btn btn-primary"
                >
                  Next Episode
                </button>
              </div>
              
              <!-- Episode Grid -->
              <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-2">
                <button
                  v-for="episode in episodes"
                  :key="episode.id"
                  @click="playEpisode(episode)"
                  class="btn"
                  :class="{ 'btn-primary': currentEpisode?.id === episode.id, 'btn-outline': currentEpisode?.id !== episode.id }"
                >
                  {{ episode.episode_number }}
                </button>
              </div>
            </div>
          </div>

          <!-- Sidebar -->
          <div class="w-full md:w-80">
            <!-- Content Info -->
            <div class="card bg-base-200 shadow-xl overflow-hidden">
              <figure>
                <img
                  :src="getImageUrl(content?.cover_image)"
                  :alt="content?.title"
                  class="w-full h-auto object-cover"
                  @error="handleImageError"
                  @load="() => console.log('Image loaded successfully')"
                >
              </figure>
              <div class="card-body">
                <h3 class="card-title">{{ content?.title }}</h3>
                
                <div class="flex items-center gap-2 mb-2">
                  <span class="badge badge-primary">{{ content?.type }}</span>
                  <div class="flex items-center">
                    <span class="text-yellow-400">★</span>
                    <span class="ml-1">{{ content?.rating || 0 }}/10</span>
                  </div>
                </div>
                
                <!-- Genres -->
                <div class="mb-2">
                  <h4 class="text-sm font-semibold mb-1">Genres</h4>
                  <div class="flex flex-wrap gap-1">
                    <span
                      v-for="genre in content?.genres"
                      :key="genre.id"
                      class="badge badge-outline"
                    >
                      {{ genre.name }}
                    </span>
                  </div>
                </div>
                
                <!-- Release Date -->
                <div v-if="content?.release_date" class="mb-2">
                  <h4 class="text-sm font-semibold mb-1">Release Date</h4>
                  <p>{{ formatDate(content.release_date) }}</p>
                </div>
              </div>
            </div>

            <!-- Related Content -->
            <div v-if="relatedContent.length > 0" class="mt-6">
              <h3 class="text-lg font-semibold mb-2">You might also like</h3>
              <div class="space-y-4">
                <router-link
                  v-for="item in relatedContent"
                  :key="item.id"
                  :to="{ name: 'watch', params: { id: item.id }}"
                  class="card card-side bg-base-200 hover:bg-base-300 transition-colors"
                >
                  <figure class="w-24">
                    <img
                      :src="item.cover_image"
                      :alt="item.title"
                      class="w-full h-full object-cover"
                    >
                  </figure>
                  <div class="card-body p-4">
                    <h4 class="card-title text-sm">{{ item.title }}</h4>
                    <p class="text-xs text-gray-400">{{ item.type }}</p>
                  </div>
                </router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useContentStore } from '../stores/content'
import VideoPlayer from '../components/VideoPlayer.vue'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const contentStore = useContentStore()


// State
const content = ref(null)
const episodes = ref([])
const relatedContent = ref([])
const currentVideo = ref(null)
const currentEpisode = ref(null)
const streamServers = ref([
  { label: 'Self Hosted', value: 'self-hosted' },
  { label: 'MP4Upload', value: 'mp4upload' },
  { label: 'Streamtape', value: 'streamtape' }
])
const downloadLinks = ref([])
const currentServer = ref('Self Hosted')
const currentQuality = ref(null)
const loading = ref(true)
const error = ref(null)

// Computed
const nextEpisode = computed(() => {
  if (!currentEpisode.value || !episodes.value.length) return null
  
  const currentIndex = episodes.value.findIndex(ep => ep.id === currentEpisode.value.id)
  if (currentIndex < episodes.value.length - 1) {
    return episodes.value[currentIndex + 1]
  }
  return null
})

const previousEpisode = computed(() => {
  if (!currentEpisode.value || !episodes.value.length) return null
  
  const currentIndex = episodes.value.findIndex(ep => ep.id === currentEpisode.value.id)
  if (currentIndex > 0) {
    return episodes.value[currentIndex - 1]
  }
  return null
})

const downloadQualities = computed(() => {
  if (!downloadLinks.value) return []
  return [...new Set(downloadLinks.value.map(link => link.quality))].sort()
})

const getDownloadLinksByQuality = (quality) => {
  console.log('Check DownloadLinks', downloadLinks)
  return downloadLinks.value?.filter(link => link.quality === quality) || []
}

const formatFileSize = (bytes) => {
  if (!bytes) return ''
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
  if (bytes === 0) return '0 Byte'
  const i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)))
  return Math.round(bytes / Math.pow(1024, i), 2) + ' ' + sizes[i]
}

// Methods

// Tambahkan fungsi-fungsi ini di dalam script setup di Watch.vue
const loadEpisodes = async () => {
  try {
    console.log('Loading episodes for content:', content.value.id)
    //const response = await axios.get(`/api/contents/${content.value.id}/episodes`)
    
    // Update episodes from response
    episodes.value = content.value.episodes || []
    console.log('Episodes:', episodes.value)

    // Sort episodes by season and episode number
    episodes.value.sort((a, b) => {
      if (a.season_number === b.season_number) {
        return a.episode_number - b.episode_number
      }
      return a.season_number - b.season_number
    })

    // Set first episode as current if no episode is selected
    if (episodes.value.length > 0) {
      const episodeId = route.query.episode ? parseInt(route.query.episode) : null
      console.log('Episode ID:', episodeId)
      if (episodeId) {
        const episode = episodes.value.find(ep => ep.episode_number === episodeId)
        if (episode) {
          playEpisode(episode)
        } else {
          playEpisode(episodes.value[0])
        }
      } else {
        playEpisode(episodes.value[0])
      }
    }

    console.log('Episodes loaded:', episodes.value.length)
  } catch (err) {
    console.error('Error loading episodes:', err)
    error.value = 'Failed to load episodes'
  }
}

const loadRelatedContent = async () => {
  try {
    // Ambil konten dengan genre yang sama
    if (content.value.genres && content.value.genres.length > 0) {
      const genreId = content.value.genres[0].id
      const response = await axios.get(`/api/contents/genre/${genreId}`, {
        params: {
          page: 1,
          pageSize: 4
        }
      })
      
      // Filter out current content
      relatedContent.value = response.data.contents.filter(
        item => item.id !== content.value.id
      ).slice(0, 3)
      console.log('Related content:', relatedContent.value)
    }
  } catch (err) {
    console.error('Error loading related content:', err)
    // Tidak perlu menampilkan error ke user untuk related content
    relatedContent.value = []
  }
}

// Modifikasi fungsi loadContent untuk menggunakan type yang dinamis
const loadContent = async () => {
  try {
    loading.value = true
    error.value = null
    
    const contentId = route.params.id
    console.log('Loading content with ID:', contentId)
    
    if (!contentId) {
      error.value = 'Invalid content ID'
      return
    }
    
    const response = await axios.get(`/api/contents/${contentId}`)
    content.value = response.data
    
    if (!content.value) {
      error.value = 'Content not found'
      return
    }
    
    console.log('Content loaded successfully:', content.value)
    
    // Load episodes jika content type adalah series atau memiliki episodes
    if (content.value.episodes?.length > 0 || ['Series', 'Anime', 'TV'].includes(content.value.type)) {
      await loadEpisodes()
    }
    
    // Load related content
    await loadRelatedContent()
    
  } catch (err) {
    console.error('Error loading content:', err)
    error.value = err.response?.data?.error || 'Failed to load content'
  } finally {
    loading.value = false
  }
}

const playEpisode = (episode) => {
  currentEpisode.value = episode
  
  // Update URL with episode ID
  router.replace({ 
    name: 'watch', 
    params: { id: content.value.id },
    query: { episode: episode.id }
  })
  
  // Find stream links for this episode
  const episodeStreamLinks = content.value.stream_links?.filter(link => 
    link.episode_number === episode.episode_number && 
    link.season_number === episode.season_number
  ) || []
  
  console.log('Stream links for episode:', episodeStreamLinks)
  
  // Set download links for the episode
  downloadLinks.value = content.value.download_links?.filter(link => 
    link.episode_number === episode.episode_number && 
    link.season_number === episode.season_number
  ) || []
  
  // Sort download links by quality (highest to lowest)
  downloadLinks.value.sort((a, b) => {
    const getQualityNumber = (quality) => parseInt(quality.replace(/[^0-9]/g, '')) || 0
    return getQualityNumber(b.quality) - getQualityNumber(a.quality)
  })
  
  // Initialize server selection if we have stream links
  if (episodeStreamLinks.length > 0) {
    // Try to keep the current server if it exists for this episode
    const currentServerExists = episodeStreamLinks.some(link => 
      (link.server === 'local' && currentServer.value.toLowerCase() === 'self hosted') ||
      (link.server === 'external' && link.name.toLowerCase() === currentServer.value.toLowerCase()) ||
      (link.server.toLowerCase() === currentServer.value.toLowerCase())
    )
    
    if (!currentServerExists) {
      // Set default server based on first stream link
      const firstLink = episodeStreamLinks[0]
      if (firstLink.server === 'local') {
        currentServer.value = 'Self Hosted'
      } else if (firstLink.server === 'external') {
        currentServer.value = firstLink.name.charAt(0).toUpperCase() + firstLink.name.slice(1)
      } else {
        currentServer.value = firstLink.server.charAt(0).toUpperCase() + firstLink.server.slice(1)
      }
      console.log('Set initial server to:', currentServer.value)
    }
  }
  
  // Update the current video based on the selected server
  updateCurrentVideo()
}

const setupSelfHostedVideo = (episode) => {
  const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  const streamUrl = `${baseURL}/api/media/stream/${content.value.id}/episodes/${episode.id}`
  console.log('Setting up self-hosted video with URL:', streamUrl)
  
  currentVideo.value = {
    url: streamUrl,
    type: 'self-hosted',
    thumbnailUrl: episode.thumbnail_url || content.value.cover_image,
    qualities: []
  }
  
  // Set video qualities if available
  if (content.value.stream_links) {
    const qualities = content.value.stream_links
      .filter(link => 
        link.episode_number === episode.episode_number && 
        link.season_number === episode.season_number &&
        link.type === 'self-hosted'
      )
      .map(link => ({
        label: link.quality,
        value: link.quality,
        url: `${streamUrl}?quality=${link.quality}`
      }))
    
    if (qualities.length > 0) {
      currentVideo.value.qualities = qualities
    }
  }
}

const selectServer = (server) => {
  currentServer.value = server
  updateCurrentVideo()
}

const selectQuality = (quality) => {
  currentQuality.value = quality
  currentVideo.value = {
    url: quality.url,
    type: quality.type,
    thumbnailUrl: currentEpisode.value?.thumbnail_url || content.value?.cover_image
  }
}

const handleProgress = ({ currentTime, duration }) => {
  if (!currentEpisode.value) return
  
  // Update watch progress
  // contentStore.updateWatchProgress(
  //   content.value.id,
  //   currentEpisode.value.id,
  //   Math.floor(currentTime)
  // )
}

const handleEnded = () => {
  if (nextEpisode.value) {
    playEpisode(nextEpisode.value)
  }
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const getImageUrl = (cover_image) => {
  if (!cover_image) return null
  const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  const cleanPath = cover_image.replace(/\\/g, '/').replace(/^media\//, '')
  console.log('Clean path:', cleanPath)
  return `${baseURL}/media/${cleanPath}`
}

const handleImageError = (e) => {
  console.error('Error loading image:', e.target.src)
  // Set a fallback image
  e.target.src = '/placeholder.png'
}

const isValidUrl = (url) => {
  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}

const formatExternalUrl = (url) => {
  // Jika tidak memiliki protokol, tambahkan https
  if (!url.startsWith('http://') && !url.startsWith('https://')) {
    return `https://${url}`
  }
  return url
}

// Fungsi untuk menangani perubahan server
const handleServerChange = (server) => {
  console.log('Server changed to:', server)
  
  // Konversi nama server jika diperlukan
  if (server === 'local') {
    currentServer.value = 'Self Hosted'
  } else {
    // Capitalize first letter
    currentServer.value = server.charAt(0).toUpperCase() + server.slice(1)
  }
  
  console.log('Current server updated to:', currentServer.value)
  updateCurrentVideo()
}

// Fungsi untuk memperbarui video yang sedang diputar berdasarkan server yang dipilih
const updateCurrentVideo = () => {
  if (!content.value || !currentEpisode.value) return
  
  console.log('Updating video with server:', currentServer.value)
  console.log('Available stream links:', content.value.stream_links)
  console.log('Current episode:', currentEpisode.value)
  
  // Normalisasi nama server untuk perbandingan
  const normalizedServerName = currentServer.value.toLowerCase()
  console.log('Normalized server name:', normalizedServerName)
  
  // Cari stream link yang sesuai dengan episode dan server yang dipilih
  const streamLink = content.value.stream_links.find(link => {
    const matchesEpisode = link.episode_number === currentEpisode.value.episode_number && 
                          link.season_number === currentEpisode.value.season_number
    
    // Normalisasi nama server dari link untuk perbandingan
    const linkServerName = link.server ? link.server.toLowerCase() : ''
    const linkName = link.name ? link.name.toLowerCase() : ''
    
    console.log(`Checking link: ep=${link.episode_number}, season=${link.season_number}, server=${linkServerName}, name=${linkName}`)
    
    const matchesServer = linkServerName === normalizedServerName || 
                         linkName === normalizedServerName ||
                         (normalizedServerName === 'self hosted' && linkServerName === 'local')
    
    return matchesEpisode && matchesServer
  })
  
  console.log('Found stream link:', streamLink)
  
  if (streamLink) {
    // Update video source
    currentVideo.value = {
      url: streamLink.url,
      type: streamLink.type,
      thumbnailUrl: currentEpisode.value.thumbnail_url || content.value.cover_image
    }
    console.log('Video updated to:', currentVideo.value)
  } else {
    console.log('No matching stream link found, trying fallback...')
    
    // Fallback ke stream link pertama jika tidak ada yang cocok
    const fallbackLinks = content.value.stream_links.filter(link => 
      link.episode_number === currentEpisode.value.episode_number && 
      link.season_number === currentEpisode.value.season_number
    )
    
    console.log('Available fallback links:', fallbackLinks)
    
    if (fallbackLinks.length > 0) {
      const fallbackLink = fallbackLinks[0]
      console.log('Using fallback link:', fallbackLink)
      
      currentVideo.value = {
        url: fallbackLink.url,
        type: fallbackLink.type,
        thumbnailUrl: currentEpisode.value.thumbnail_url || content.value.cover_image
      }
      
      // Update current server to match the fallback link
      let serverName = fallbackLink.name
      if (fallbackLink.server === 'local') {
        serverName = 'Self Hosted'
      } else if (fallbackLink.server && fallbackLink.server !== 'external') {
        serverName = fallbackLink.server
      }
      
      // Capitalize first letter
      currentServer.value = serverName.charAt(0).toUpperCase() + serverName.slice(1)
      console.log('Updated current server to:', currentServer.value)
    } else {
      console.error('No stream links found for current episode')
    }
  }
}

// Lifecycle
onMounted(() => {
  console.log('Watch component mounted, initializing with server:', currentServer.value)
  loadContent()
})

// Watch for route changes
watch(
  () => route.params.id,
  () => {
    loadContent()
  }
)

watch(
  () => route.query.episode,
  (newEpisodeId) => {
    console.log('New episode ID:', newEpisodeId)
    if (newEpisodeId && episodes.value.length > 0) {
      const episode = episodes.value.find(ep => ep.episode_number === parseInt(newEpisodeId))
      console.log('Found episode:', episode)
      if (episode && (!currentEpisode.value || currentEpisode.value.id !== episode.id)) {
        playEpisode(episode)
      }
    }
  }
)

watch(currentEpisode, () => {
  updateCurrentVideo()
}, { immediate: true })
</script>

<style scoped>
/* Custom styles for the watch page */
.episode-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(60px, 1fr));
  gap: 0.5rem;
}
</style> 