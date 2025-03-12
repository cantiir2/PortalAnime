<!-- VideoPlayer.vue -->
<template>
  <div class="relative w-full aspect-video bg-black">
    <!-- Video Player -->
    <video
      v-if="type === 'self-hosted'"
      ref="videoRef"
      class="w-full h-full"
      :poster="poster"
      @timeupdate="handleTimeUpdate"
      @loadedmetadata="handleLoadedMetadata"
      @ended="handleEnded"
      @error="handleError"
      @waiting="loading = true"
      @playing="loading = false"
      @pause="isPlaying = false"
      @play="isPlaying = true"
    >
      <source :src="videoUrl" :type="videoType">
      Your browser does not support the video tag.
    </video>

    <!-- Embed Player -->
    <iframe
      v-else-if="type === 'embed'"
      class="w-full h-full"
      :src="embedSource"
      frameborder="0"
      allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
      allowfullscreen
    ></iframe>

    <!-- Server Selection -->
    <div class="absolute top-4 right-4 z-10">
      <div class="dropdown dropdown-end">
        <label tabindex="0" class="btn btn-sm m-1 bg-opacity-50 backdrop-blur-sm">
          {{ currentServer }}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg>
        </label>
        <ul tabindex="0" class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-52">
          <li v-for="server in availableServers" :key="server">
            <a @click="changeServer(server)" :class="{ 'active': currentServer === server }">
              {{ server }}
            </a>
          </li>
        </ul>
      </div>
    </div>

    <!-- Loading Overlay -->
    <div v-if="type === 'self-hosted' && loading" class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50">
      <div class="loading loading-spinner loading-lg text-primary"></div>
    </div>

    <!-- Error Overlay -->
    <div v-if="type === 'self-hosted' && error" class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50">
      <div class="text-error text-center p-4">
        <div class="text-2xl mb-2">Error</div>
        <div>{{ error }}</div>
        <button @click="retryLoad" class="btn btn-primary mt-4">Retry</button>
      </div>
    </div>

    <!-- Video Controls -->
    <div v-if="type === 'self-hosted' && !error" 
         class="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 to-transparent opacity-0 hover:opacity-100 transition-opacity duration-300">
      <!-- Progress Bar -->
      <div class="px-4 py-2">
        <div class="w-full h-1 bg-gray-600 rounded cursor-pointer" 
             @click="handleProgressClick" 
             ref="progressBar">
          <div class="h-full bg-primary rounded" 
               :style="{ width: `${(currentTime / duration) * 100}%` }"></div>
        </div>
      </div>

      <!-- Control Buttons -->
      <div class="flex items-center justify-between p-4">
        <div class="flex items-center space-x-4">
          <!-- Play/Pause -->
          <button @click="togglePlay" class="btn btn-circle btn-sm">
            <svg v-if="isPlaying" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6"/>
            </svg>
            <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/>
            </svg>
          </button>

          <!-- Volume -->
          <button @click="toggleMute" class="btn btn-circle btn-sm">
            <svg v-if="isMuted" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z"/>
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2"/>
            </svg>
            <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.536 8.464a5 5 0 010 7.072m2.828-9.9a9 9 0 010 12.728M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z"/>
            </svg>
          </button>

          <!-- Time Display -->
          <div class="text-sm text-white">
            {{ formatTime(currentTime) }} / {{ formatTime(duration) }}
          </div>
        </div>

        <div class="flex items-center space-x-4">
          <!-- Quality Selector -->
          <div class="dropdown dropdown-top dropdown-end" v-if="qualities.length > 0">
            <button class="btn btn-circle btn-sm">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
              </svg>
            </button>
            <ul class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-52">
              <li v-for="quality in qualities" :key="quality.value">
                <a @click="setQuality(quality.value)" 
                   :class="{ 'active': currentQuality === quality.value }">
                  {{ quality.label }}
                </a>
              </li>
            </ul>
          </div>

          <!-- Playback Speed -->
          <div class="dropdown dropdown-top dropdown-end">
            <button class="btn btn-circle btn-sm">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </button>
            <ul class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-52">
              <li v-for="speed in playbackSpeeds" :key="speed">
                <a @click="setPlaybackSpeed(speed)" 
                   :class="{ 'active': currentPlaybackSpeed === speed }">
                  {{ speed }}x
                </a>
              </li>
            </ul>
          </div>

          <!-- Fullscreen -->
          <button @click="toggleFullscreen" class="btn btn-circle btn-sm">
            <svg v-if="isFullscreen" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 9V4.5M9 9H4.5M9 9L3.75 3.75M9 15v4.5M9 15H4.5M9 15l-5.25 5.25M15 9h4.5M15 9V4.5M15 9l5.25-5.25M15 15h4.5M15 15v4.5M15 15l5.25 5.25"/>
            </svg>
            <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5v-4m0 4h-4m4 0l-5-5"/>
            </svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useFullscreen } from '@vueuse/core'

const props = defineProps({
  source: {
    type: String,
    required: true
  },
  type: {
    type: String,
    default: 'self-hosted',
    validator: (value) => ['self-hosted', 'embed'].includes(value)
  },
  poster: {
    type: String,
    default: ''
  },
  qualities: {
    type: Array,
    default: () => []
  },
  streamLinks: {
    type: Array,
    default: () => []
  },
  initialServer: {
    type: String,
    default: 'Self Hosted'
  }
})

const emit = defineEmits(['progress', 'ended', 'error', 'server-change'])

const videoRef = ref(null)
const progressBar = ref(null)
const loading = ref(true)
const error = ref(null)
const isPlaying = ref(false)
const isMuted = ref(false)
const currentTime = ref(0)
const duration = ref(0)
const volume = ref(1)
const currentQuality = ref(props.qualities[0]?.value || '')
const playbackSpeed = ref(1)
const playbackSpeeds = [0.25, 0.5, 0.75, 1, 1.25, 1.5, 2]
const isLoading = ref(true)

// Fullscreen
const { isFullscreen, toggle: toggleFullscreen } = useFullscreen(videoRef)

const availableServers = computed(() => {
  const servers = props.streamLinks.map(link => {
    if (link.server === 'local') return 'Self Hosted'
    return link.name.charAt(0).toUpperCase() + link.name.slice(1)
  })
  return [...new Set(servers)] // Remove duplicates
})

const currentServer = ref(props.initialServer)

const changeServer = (server) => {
  currentServer.value = server
  emit('server-change', server.toLowerCase() === 'self hosted' ? 'local' : server.toLowerCase())
}

// Computed property for video source URL
const videoUrl = computed(() => {
  if (!props.source) return ''
  
  try {
    // Check if it's already a valid URL
    new URL(props.source)
    return props.source
  } catch {
    // If not a valid URL, construct it with base URL
    const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    return `${baseURL}${props.source.startsWith('/') ? '' : '/'}${props.source}`
  }
})

// Add new computed property for embed source
const embedSource = computed(() => {
  if (!props.source) return ''
  
  // If source contains IFRAME tag, extract src attribute
  if (props.source.toUpperCase().includes('IFRAME')) {
    const srcMatch = props.source.match(/src=["'](.*?)["']/i)
    return srcMatch ? srcMatch[1] : ''
  }
  
  // Otherwise use source directly
  return props.source
})

const videoType = computed(() => {
  if (props.source.includes('api/media/stream')) {
    return 'video/mp4' // Streaming endpoint always returns MP4
  }
  
  const ext = props.source.split('.').pop()?.toLowerCase()
  switch (ext) {
    case 'mp4':
      return 'video/mp4'
    case 'webm':
      return 'video/webm'
    case 'mkv':
      return 'video/x-matroska'
    default:
      return 'video/mp4'
  }
})

// Methods
const handleLoadStart = () => {
  console.log('Video loading started:', props.source)
  isLoading.value = true
  error.value = null
}

const handleLoadedMetadata = () => {
  console.log('Video metadata loaded:', {
    duration: videoRef.value?.duration,
    videoWidth: videoRef.value?.videoWidth,
    videoHeight: videoRef.value?.videoHeight
  })
  isLoading.value = false
  duration.value = videoRef.value?.duration || 0
}

const handleError = (e) => {
  const videoElement = e.target
  console.error('Video error:', {
    error: e.target?.error,
    networkState: videoElement?.networkState,
    readyState: videoElement?.readyState,
    currentSrc: videoElement?.currentSrc,
    source: props.source
  })
  isLoading.value = false
  error.value = e.target?.error?.message || 'Failed to load video'
  emit('error', error.value)
}

const retryLoad = () => {
  if (!videoRef.value) return
  console.log('Retrying video load:', props.source)
  error.value = null
  isLoading.value = true
  videoRef.value.load()
}

const handleTimeUpdate = () => {
  if (!videoRef.value) return
  currentTime.value = videoRef.value.currentTime
  emit('progress', {
    currentTime: currentTime.value,
    duration: duration.value
  })
}

const handleEnded = () => {
  isPlaying.value = false
  emit('ended')
}

const togglePlay = () => {
  if (!videoRef.value) return
  
  if (videoRef.value.paused) {
    videoRef.value.play()
    isPlaying.value = true
  } else {
    videoRef.value.pause()
    isPlaying.value = false
  }
}

const toggleMute = () => {
  if (!videoRef.value) return
  isMuted.value = !isMuted.value
  videoRef.value.muted = isMuted.value
}

const handleProgressClick = (event) => {
  if (!videoRef.value) return
  
  const progressBar = event.currentTarget
  const rect = progressBar.getBoundingClientRect()
  const ratio = (event.clientX - rect.left) / rect.width
  const newTime = ratio * duration.value
  
  videoRef.value.currentTime = newTime
}

const formatTime = (time) => {
  if (!time) return '0:00'
  const minutes = Math.floor(time / 60)
  const seconds = Math.floor(time % 60)
  return `${minutes}:${seconds.toString().padStart(2, '0')}`
}

const setQuality = (quality) => {
  const newQualitySource = props.qualities.find(q => q.value === quality)
  if (newQualitySource) {
    const currentTime = videoRef.value?.currentTime || 0
    currentQuality.value = quality
    videoRef.value.src = newQualitySource.url
    videoRef.value.currentTime = currentTime
    if (isPlaying.value) {
      videoRef.value.play()
    }
  }
}

const setPlaybackSpeed = (speed) => {
  if (!videoRef.value) return
  playbackSpeed.value = speed
  videoRef.value.playbackRate = speed
}

// Keyboard controls
const handleKeydown = (event) => {
  if (!videoRef.value || props.type !== 'self-hosted') return
  
  switch (event.key.toLowerCase()) {
    case ' ':
    case 'k':
      event.preventDefault()
      togglePlay()
      break
    case 'm':
      toggleMute()
      break
    case 'f':
      toggleFullscreen()
      break
    case 'arrowleft':
      videoRef.value.currentTime -= 5
      break
    case 'arrowright':
      videoRef.value.currentTime += 5
      break
  }
}

// Lifecycle
onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
  console.log('Video player mounted:', {
    source: props.source,
    type: props.type,
    fullUrl: videoUrl.value
  })
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})

// Watch for source changes
watch(() => props.source, () => {
  if (videoRef.value) {
    error.value = null
    isLoading.value = true
    videoRef.value.load()
  }
})
</script>

<style scoped>
.dropdown-content {
  visibility: hidden;
  opacity: 0;
  transition: visibility 0s, opacity 0.3s;
}

.dropdown:hover .dropdown-content,
.dropdown:focus-within .dropdown-content {
  visibility: visible;
  opacity: 1;
}

.active {
  @apply bg-primary text-white;
}
</style> 