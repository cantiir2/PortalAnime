import { defineStore } from 'pinia'
import axios from 'axios'

export const useContentStore = defineStore('content', {
  state: () => ({
    contents: [],
    currentContent: null,
    currentEpisode: null,
    loading: false,
    error: null,
    filters: {
      type: null,
      genre: null,
      category: null,
      search: ''
    },
    pagination: {
      page: 1,
      pageSize: 12,
      total: 0
    }
  }),

  getters: {
    filteredContents: (state) => state.contents,
    hasNextPage: (state) => state.pagination.page * state.pagination.pageSize < state.pagination.total
  },

  actions: {
    async fetchContents() {
      try {
        this.loading = true
        this.error = null
        const params = {
          page: this.pagination.page,
          pageSize: this.pagination.pageSize,
          ...this.filters
        }
        const response = await axios.get('/api/contents', { params })
        this.contents = response.data.contents
        this.pagination.total = response.data.total
      } catch (error) {
        this.error = error.response?.data?.message || 'An error occurred'
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchContentById(id) {
      try {
        this.loading = true
        this.error = null
        
        console.log('Fetching content with ID:', id)
        
        if (!id) {
          this.error = 'Invalid content ID'
          console.error('Invalid content ID:', id)
          return
        }
        
        const response = await axios.get(`/api/contents/${id}`)
        this.currentContent = response.data
        
        console.log('Content fetched successfully:', this.currentContent)
      } catch (error) {
        console.error('Error fetching content:', {
          message: error.message,
          response: error.response?.data,
          status: error.response?.status
        })
        
        this.error = error.response?.data?.error || error.message || 'Failed to fetch content'
        this.currentContent = null
      } finally {
        this.loading = false
      }
    },

    async fetchEpisodeById(contentId, episodeId) {
      try {
        this.loading = true
        this.error = null
        const response = await axios.get(`/api/contents/${contentId}/episodes/${episodeId}`)
        this.currentEpisode = response.data
      } catch (error) {
        this.error = error.response?.data?.message || 'An error occurred'
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateWatchProgress(contentId, episodeId, progress) {
      try {
        await axios.post('/api/watch-history', {
          contentId,
          episodeId,
          progress
        })
      } catch (error) {
        console.error('Failed to update watch progress:', error)
      }
    },

    async uploadVideo(contentId, episodeId, file) {
      try {
        const formData = new FormData();
        formData.append('video', file);

        const url = episodeId 
          ? `/api/media/content/${contentId}/episodes/${episodeId}/video`
          : `/api/media/content/${contentId}/video`;

        const response = await axios.post(url, formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        });

        return response.data;
      } catch (error) {
        console.error('Failed to upload video:', error);
        throw error;
      }
    },

    setFilters(newFilters) {
      this.filters = { ...this.filters, ...newFilters }
      this.pagination.page = 1
      this.fetchContents()
    },

    nextPage() {
      if (this.hasNextPage) {
        this.pagination.page++
        this.fetchContents()
      }
    },

    resetFilters() {
      this.filters = {
        type: null,
        genre: null,
        category: null,
        search: ''
      }
      this.pagination.page = 1
      this.fetchContents()
    }
  }
}) 