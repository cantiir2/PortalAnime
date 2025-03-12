<template>
  <div class="w-full">
    <div 
      class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-primary-500 transition-colors cursor-pointer"
      :class="{ 'border-primary-500': isDragging }"
      @dragenter.prevent="isDragging = true"
      @dragleave.prevent="isDragging = false"
      @dragover.prevent
      @drop.prevent="handleDrop"
      @click="$refs.fileInput.click()"
    >
      <div v-if="!selectedFile && !uploading">
        <svg class="mx-auto h-12 w-12 text-gray-400" stroke="currentColor" fill="none" viewBox="0 0 48 48">
          <path d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4-4m4-4h8m-4-4v8m-12 4h.02" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        <div class="mt-4">
          <span class="text-sm font-medium text-primary-600">
            Klik untuk upload
          </span>
          <span class="text-sm text-gray-500">
            atau drag and drop
          </span>
        </div>
        <p class="mt-1 text-xs text-gray-500">
          MP4 up to 500MB
        </p>
      </div>

      <div v-else-if="selectedFile && !uploading" class="text-left">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-2">
            <svg class="h-6 w-6 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span class="text-sm font-medium text-gray-900">{{ selectedFile.name }}</span>
          </div>
          <button 
            @click.stop="clearFile" 
            class="text-sm font-medium text-primary-600 hover:text-primary-500"
          >
            Remove
          </button>
        </div>
        <div class="mt-4 flex justify-end">
          <button 
            @click.stop="uploadFile" 
            class="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            Upload Video
          </button>
        </div>
      </div>

      <div v-else class="text-center">
        <div class="flex justify-center">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
        </div>
        <p class="mt-2 text-sm text-gray-500">Uploading video...</p>
      </div>
    </div>

    <input
      ref="fileInput"
      type="file"
      accept="video/mp4"
      class="hidden"
      @change="handleFileSelect"
    >
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useContentStore } from '@/stores/content'

const props = defineProps({
  contentId: {
    type: Number,
    required: true
  },
  episodeId: {
    type: Number,
    default: null
  }
})

const emit = defineEmits(['upload-success', 'upload-error'])

const contentStore = useContentStore()
const fileInput = ref(null)
const selectedFile = ref(null)
const isDragging = ref(false)
const uploading = ref(false)

const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (file && file.type === 'video/mp4' && file.size <= 500 * 1024 * 1024) {
    selectedFile.value = file
  } else {
    alert('Please select an MP4 file under 500MB')
  }
}

const handleDrop = (event) => {
  isDragging.value = false
  const file = event.dataTransfer.files[0]
  if (file && file.type === 'video/mp4' && file.size <= 500 * 1024 * 1024) {
    selectedFile.value = file
  } else {
    alert('Please drop an MP4 file under 500MB')
  }
}

const clearFile = () => {
  selectedFile.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const uploadFile = async () => {
  if (!selectedFile.value) return

  uploading.value = true
  try {
    const result = await contentStore.uploadVideo(
      props.contentId,
      props.episodeId,
      selectedFile.value
    )
    emit('upload-success', result)
    clearFile()
  } catch (error) {
    emit('upload-error', error)
  } finally {
    uploading.value = false
  }
}
</script> 