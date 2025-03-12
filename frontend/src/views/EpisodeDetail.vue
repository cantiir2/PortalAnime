<template>
  <div class="container mx-auto p-4">
    <div class="bg-white rounded-lg shadow p-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4">Upload Video Episode</h3>
      <VideoUploader 
        :content-id="contentId"
        :episode-id="episodeId"
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
      contentId: '',
      episodeId: ''
    }
  },
  methods: {
    handleVideoUploadSuccess(result) {
      // Refresh episode data to show new video
      this.fetchEpisode()
      this.toast.success('Video berhasil diupload')
    },
    handleVideoUploadError(error) {
      console.error('Video upload failed:', error)
      this.toast.error('Gagal mengupload video: ' + error.message)
    }
  }
}
</script>

<style>
  /* Add your styles here */
</style> 