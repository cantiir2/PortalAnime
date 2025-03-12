<template>
  <div v-if="isAdmin" class="mt-6 space-y-4">
    <div class="bg-white rounded-lg shadow p-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4">Upload Video</h3>
      <VideoUploader 
        :content-id="contentId" 
        @upload-success="handleVideoUploadSuccess"
        @upload-error="handleVideoUploadError"
      />
    </div>
  </div>
</template>

<script>
import VideoUploader from '@/components/VideoUploader.vue'

export default {
  components: {
    VideoUploader
  },
  data() {
    return {
      isAdmin: true, // Assuming isAdmin is true for the template
      contentId: this.$route.params.id // Assuming contentId is passed via route
    }
  },
  methods: {
    handleVideoUploadSuccess(result) {
      // Refresh content data to show new video
      this.fetchContent()
      this.toast.success('Video berhasil diupload')
    },
    handleVideoUploadError(error) {
      console.error('Video upload failed:', error)
      this.toast.error('Gagal mengupload video: ' + error.message)
    },
    fetchContent() {
      // Implement the logic to fetch content data
    }
  }
}
</script>

<style>
  /* Add your styles here */
</style> 