<template>
  <div class="min-h-screen p-4">
    <!-- Error Alert -->
    <div v-if="showError" class="alert alert-error shadow-lg mb-4">
      <div>
        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current flex-shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span>{{ error }}</span>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <h1 class="text-2xl font-bold mb-6">Admin Dashboard</h1>

      <!-- Tabs -->
      <div class="tabs tabs-boxed mb-6">
        <button 
          class="tab" 
          :class="{ 'tab-active': activeTab === 'contents' }"
          @click="activeTab = 'contents'"
        >
          Contents
        </button>
        <button 
          class="tab" 
          :class="{ 'tab-active': activeTab === 'genres' }"
          @click="activeTab = 'genres'"
        >
          Genres
        </button>
        <button 
          class="tab" 
          :class="{ 'tab-active': activeTab === 'categories' }"
          @click="activeTab = 'categories'"
        >
          Categories
        </button>
        <button 
          class="tab" 
          :class="{ 'tab-active': activeTab === 'seasons' }"
          @click="activeTab = 'seasons'"
        >
          Seasons
        </button>
      </div>

      <!-- Content Management -->
      <div v-if="activeTab === 'contents'" class="space-y-6">
        <div class="flex justify-between items-center">
          <h2 class="text-xl font-semibold">Content List</h2>
          <button class="btn btn-primary" @click="openContentModal">
            Add New Content
          </button>
        </div>

        <!-- Content Table -->
        <div class="overflow-x-auto bg-base-200 rounded-lg">
          <table class="table w-full">
            <thead>
              <tr>
                <th>No</th>
                <th>Title</th>
                <th>Type</th>
                <th>Release Date</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(content, index) in contents" :key="content.id">
                <td>{{ index + 1 }}</td>
                <td>{{ content.title }}</td>
                <td>
                  <span class="badge" :class="{
                    'badge-primary': content.type === 'Movie' || content.type === 'Series',
                    'badge-secondary': content.type === 'Anime'
                  }">
                    {{ content.type }}
                  </span>
                </td>
                <td>{{ formatDate(content.release_date) }}</td>
                <td>
                  <div class="flex gap-2">
                    <button 
                      class="btn btn-sm btn-info"
                      @click="editContent(content)"
                    >
                      Edit
                    </button>
                    <button 
                      class="btn btn-sm btn-error"
                      @click="deleteContent(content.id)"
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Seasons Management -->
      <div v-if="activeTab === 'seasons'" class="space-y-6">
        <div class="flex justify-between items-center">
          <h2 class="text-xl font-semibold">Seasons</h2>
          <button class="btn btn-primary" @click="showAddSeasonModal = true">
            Add New Season
          </button>
        </div>

        <!-- Seasons Table -->
        <div class="overflow-x-auto">
          <table class="table w-full">
            <thead>
              <tr>
                <th>ID</th>
                <th>Name</th>
                <th>Year</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="season in seasons" :key="season.id">
                <td>{{ season.id }}</td>
                <td>{{ season.name }}</td>
                <td>{{ season.year }}</td>
                <td>
                  <span class="badge" :class="{
                    'badge-success': season.status === 'Active',
                    'badge-warning': season.status === 'Coming Soon',
                    'badge-error': season.status === 'Ended'
                  }">
                    {{ season.status }}
                  </span>
                </td>
                <td>
                  <div class="flex gap-2">
                    <button 
                      class="btn btn-sm btn-info"
                      @click="editSeason(season)"
                    >
                      Edit
                    </button>
                    <button 
                      class="btn btn-sm btn-error"
                      @click="deleteSeason(season.id)"
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Add/Edit Season Modal -->
      <dialog id="seasonModal" class="modal" :class="{ 'modal-open': showAddSeasonModal }">
        <div class="modal-box">
          <h3 class="font-bold text-lg mb-4">
            {{ editingSeasonId ? 'Edit Season' : 'Add New Season' }}
          </h3>
          <form @submit.prevent="saveSeason">
            <div class="form-control">
              <label class="label">
                <span class="label-text">Name</span>
              </label>
              <select v-model="seasonForm.name" class="select select-bordered w-full">
                <option value="Winter">Winter</option>
                <option value="Spring">Spring</option>
                <option value="Summer">Summer</option>
                <option value="Fall">Fall</option>
              </select>
            </div>

            <div class="form-control mt-4">
              <label class="label">
                <span class="label-text">Year</span>
              </label>
              <input 
                type="number" 
                v-model="seasonForm.year" 
                class="input input-bordered"
                min="2000"
                :max="currentYear + 1"
              />
            </div>

            <div class="form-control mt-4">
              <label class="label">
                <span class="label-text">Status</span>
              </label>
              <select v-model="seasonForm.status" class="select select-bordered w-full">
                <option value="Coming Soon">Coming Soon</option>
                <option value="Active">Active</option>
                <option value="Ended">Ended</option>
              </select>
            </div>

            <div class="modal-action">
              <button type="button" class="btn" @click="closeSeasonModal">Cancel</button>
              <button type="submit" class="btn btn-primary">Save</button>
            </div>
          </form>
        </div>
      </dialog>

      <!-- Genres Management -->
      <div v-if="activeTab === 'genres'" class="space-y-6">
        <div class="flex justify-between items-center">
          <h2 class="text-xl font-semibold">Genres</h2>
          <button class="btn btn-primary" @click="showAddGenreModal = true">
            Add New Genre
          </button>
        </div>

        <!-- Genres Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          <div v-for="genre in genres" :key="genre.id" class="card bg-base-200">
            <div class="card-body">
              <h3 class="card-title">{{ genre.name }}</h3>
              <p>{{ genre.description }}</p>
              <div class="card-actions justify-end">
                <button 
                  class="btn btn-sm btn-info"
                  @click="editGenre(genre)"
                >
                  Edit
                </button>
                <button 
                  class="btn btn-sm btn-error"
                  @click="deleteGenre(genre.id)"
                >
                  Delete
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Categories Management -->
      <div v-if="activeTab === 'categories'" class="space-y-6">
        <div class="flex justify-between items-center">
          <h2 class="text-xl font-semibold">Categories</h2>
          <button class="btn btn-primary" @click="showAddCategoryModal = true">
            Add New Category
          </button>
        </div>

        <!-- Categories Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          <div v-for="category in categories" :key="category.id" class="card bg-base-200">
            <div class="card-body">
              <h3 class="card-title">{{ category.name }}</h3>
              <p>{{ category.description }}</p>
              <div class="card-actions justify-end">
                <button 
                  class="btn btn-sm btn-info"
                  @click="editCategory(category)"
                >
                  Edit
                </button>
                <button 
                  class="btn btn-sm btn-error"
                  @click="deleteCategory(category.id)"
                >
                  Delete
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Content Modal -->
    <dialog :class="{ 'modal modal-open': showAddContentModal }" @click="showAddContentModal = false">
      <div class="modal-box max-w-3xl" @click.stop>
        <h3 class="font-bold text-lg">
          {{ editingContent ? 'Edit Content' : 'Add New Content' }}
        </h3>
        
        <div class="space-y-4 mt-4">
          <!-- Cover Image -->
          <div class="form-control">
            <label class="label">Cover Image</label>
            <div class="flex items-center space-x-4">
              <div class="w-32 h-32 relative">
                <img 
                  v-if="coverImagePreview" 
                  :src="coverImagePreview" 
                  class="w-full h-full object-cover rounded-lg"
                  alt="Cover preview"
                />
                <div 
                  v-else 
                  class="w-full h-full border-2 border-dashed border-gray-400 rounded-lg flex items-center justify-center"
                >
                  <span class="text-gray-400">No image</span>
                </div>
              </div>
              <input 
                type="file" 
                ref="coverImageInput"
                @change="handleCoverImageChange" 
                accept="image/*"
                class="file-input file-input-bordered w-full max-w-xs"
              />
            </div>
          </div>

          <!-- Basic Info -->
          <div class="form-control">
            <label class="label">Title</label>
            <input 
              v-model="contentForm.title" 
              type="text" 
              class="input input-bordered" 
              required
            />
          </div>

          <div class="form-control">
            <label class="label">Description</label>
            <textarea 
              v-model="contentForm.description" 
              class="textarea textarea-bordered h-24"
              placeholder="Enter content description"
            ></textarea>
          </div>

          <div class="form-control">
            <label class="label">Release Date</label>
            <input 
              v-model="contentForm.releaseDate" 
              type="date" 
              class="input input-bordered" 
              required
            />
          </div>

          <div class="form-control">
            <label class="label">Rating</label>
            <input 
              v-model="contentForm.rating" 
              type="number" 
              min="0" 
              max="10" 
              step="0.1"
              class="input input-bordered" 
              placeholder="0-10"
            />
          </div>

          <!-- Type Selection from Categories -->
          <div class="form-control">
            <label class="label">Type</label>
            <select 
              v-model="contentForm.type" 
              class="select select-bordered w-full"
              required
            >
              <option value="" disabled>Select type</option>
              <option v-for="category in categories" :key="category.id" :value="category.name">
                {{ category.name }}
              </option>
            </select>
          </div>

          <!-- Season -->
          <div class="form-control">
            <label class="label">
              <span class="label-text">Season</span>
            </label>
            <select v-model="contentForm.seasonId" class="select select-bordered w-full">
              <option value="">Select Season</option>
              <option v-for="season in seasons" :key="season.id" :value="season.id">
                {{ season.name }} {{ season.year }}
              </option>
            </select>
          </div>

          <!-- Multiple Genre Selection -->
          <div class="form-control">
            <label class="label">Genres</label>
            <div class="dropdown">
              <label tabindex="0" class="btn btn-bordered w-full justify-between">
                {{ selectedGenres.length ? `${selectedGenres.length} selected` : 'Select genres' }}
                <svg class="w-4 h-4 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </label>
              <div tabindex="0" class="dropdown-content bg-base-200 w-full p-2 rounded-box shadow-lg max-h-60 overflow-y-auto" style="z-index: 999;">
                <div v-for="genre in genres" :key="genre.id" class="form-control">
                  <label class="label cursor-pointer">
                    <span class="label-text">{{ genre.name }}</span>
                    <input 
                      type="checkbox" 
                      :value="genre.id"
                      v-model="contentForm.genreIds" 
                      class="checkbox" 
                    />
                  </label>
                </div>
              </div>
            </div>
          </div>

          <!-- Episode List Section -->
          <div class="space-y-4">
            <h3 class="text-lg font-semibold">Episodes</h3>
            
            <!-- Episode List -->
            <div v-if="contentForm.episodes.length > 0" class="space-y-4">
              <div v-for="(episode, index) in contentForm.episodes" :key="index" 
                   class="card bg-base-200 p-4">
                <div class="flex justify-between items-center mb-4">
                  <h4 class="text-lg font-medium">
                    Season {{ episode.seasonNumber }} Episode {{ episode.episodeNumber }}: {{ episode.title }}
                  </h4>
                  <div class="flex gap-2">
                    <button type="button" 
                            @click="episodeManagement.openEditEpisodeModal(episode, index)" 
                            class="btn btn-sm btn-primary">
                      Edit
                    </button>
                    <button type="button" 
                            @click="episodeManagement.removeEpisode(index)" 
                            class="btn btn-sm btn-error">
                      Remove
                    </button>
                  </div>
                </div>
                
                <!-- Episode Info Preview -->
                <div class="grid gap-4">
                  <p class="text-sm">{{ episode.description }}</p>
                  
                  <!-- Stream Links Preview -->
                  <div v-if="episode.streamLinks.length > 0">
                    <h5 class="font-medium text-sm mb-2">Stream Links: {{ episode.streamLinks.length }}</h5>
                  </div>
                  
                  <!-- Download Links Preview -->
                  <div v-if="episode.downloadLinks.length > 0">
                    <h5 class="font-medium text-sm mb-2">Download Links: {{ episode.downloadLinks.length }}</h5>
                  </div>
                </div>
              </div>
            </div>

            <!-- Add Episode Button -->
            <button type="button" 
                    @click="episodeManagement.openAddEpisodeModal" 
                    class="btn btn-primary">
              Add Episode
            </button>
          </div>
        </div>

        <div class="modal-action">
          <button type="button" class="btn" @click="showAddContentModal = false">Cancel</button>
          <button type="button" class="btn btn-primary" @click="handleContentSubmit">Save</button>
        </div>
      </div>
    </dialog>

    <!-- Modal untuk Genre -->
    <dialog :class="{ 'modal modal-open': showAddGenreModal }" @click="showAddGenreModal = false">
      <div class="modal-box" @click.stop>
        <h3 class="font-bold text-lg mb-4">
          {{ editingGenre ? 'Edit Genre' : 'Add New Genre' }}
        </h3>
        <form @submit.prevent="handleGenreSubmit" class="space-y-4">
          <div class="form-control">
            <label class="label">Name</label>
            <input 
              v-model="genreForm.name" 
              type="text" 
              class="input input-bordered" 
              required
            />
          </div>
          <div class="form-control">
            <label class="label">Description</label>
            <textarea 
              v-model="genreForm.description" 
              class="textarea textarea-bordered"
              rows="3"
            ></textarea>
          </div>
          <div class="modal-action">
            <button type="button" class="btn" @click="showAddGenreModal = false">Cancel</button>
            <button type="submit" class="btn btn-primary">Save</button>
          </div>
        </form>
      </div>
    </dialog>

    <!-- Modal untuk Category -->
    <dialog :class="{ 'modal modal-open': showAddCategoryModal }" @click="showAddCategoryModal = false">
      <div class="modal-box" @click.stop>
        <h3 class="font-bold text-lg mb-4">
          {{ editingCategory ? 'Edit Category' : 'Add New Category' }}
        </h3>
        <form @submit.prevent="handleCategorySubmit" class="space-y-4">
          <div class="form-control">
            <label class="label">Name</label>
            <input 
              v-model="categoryForm.name" 
              type="text" 
              class="input input-bordered" 
              required
            />
          </div>
          <div class="form-control">
            <label class="label">Description</label>
            <textarea 
              v-model="categoryForm.description" 
              class="textarea textarea-bordered"
              rows="3"
            ></textarea>
          </div>
          <div class="modal-action">
            <button type="button" class="btn" @click="showAddCategoryModal = false">Cancel</button>
            <button type="submit" class="btn btn-primary">Save</button>
          </div>
        </form>
      </div>
    </dialog>

    <!-- Episode Modal -->
    <dialog :class="{ 'modal modal-open': showEpisodeModal }" @click="episodeManagement.closeEpisodeModal">
      <div class="modal-box" @click.stop>
        <h3 class="font-bold text-lg mb-4">
          {{ editingEpisodeIndex !== null ? 'Edit Episode' : 'Add New Episode' }}
        </h3>
        <form @submit.prevent="episodeManagement.saveEpisode">
          <div class="space-y-4">
            <!-- Episode Info -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="form-control">
                <label class="label">Season Number</label>
                <input 
                  v-model="episodeForm.season_number" 
                  type="number" 
                  min="1" 
                  class="input input-bordered" 
                  required
                />
              </div>
              
              <div class="form-control">
                <label class="label">Episode Number</label>
                <input 
                  v-model="episodeForm.episode_number" 
                  type="number" 
                  min="1" 
                  class="input input-bordered" 
                  required
                />
              </div>
            </div>
            
            <div class="form-control">
              <label class="label">Title</label>
              <input 
                v-model="episodeForm.title" 
                type="text" 
                class="input input-bordered" 
                required
              />
            </div>
            
            <div class="form-control">
              <label class="label">Description</label>
              <textarea 
                v-model="episodeForm.description" 
                class="textarea textarea-bordered h-24"
              ></textarea>
            </div>
            
            <!-- Stream Links -->
            <div class="space-y-2">
              <div class="flex justify-between items-center">
                <h4 class="font-medium">Stream Links</h4>
                <button 
                  type="button"
                  @click="episodeManagement.addStreamLink"
                  class="btn btn-sm btn-primary"
                >
                  Add Stream Link
                </button>
              </div>
              
              <div v-for="(link, index) in episodeForm.streamLinks" :key="index" class="card bg-base-200 p-4">
                <div class="flex justify-between items-start mb-2">
                  <h5 class="font-medium">Stream Link #{{ index + 1 }}</h5>
                  <button 
                    type="button"
                    @click="episodeManagement.removeStreamLink(index)"
                    class="btn btn-sm btn-error"
                  >
                    Remove
                  </button>
                </div>
                
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-2">
                  <div class="form-control">
                    <label class="label">Name</label>
                    <input 
                      v-model="link.name" 
                      type="text" 
                      class="input input-bordered" 
                      placeholder="Server name"
                      required
                    />
                  </div>
                  
                  <div class="form-control">
                    <label class="label">Quality</label>
                    <input 
                      v-model="link.quality" 
                      type="text" 
                      class="input input-bordered" 
                      placeholder="e.g. 720p"
                      required
                    />
                  </div>
                </div>
                <div class="form-control">
                  <label class="label">Type</label>
                  <select v-model="link.type" class="select select-bordered" required>
                    <option value="embed">Embedded Link</option>
                    <option value="self-hosted">Self-hosted Video</option>
                  </select>
                </div>
                
                <!-- URL Input - Hanya tampilkan untuk embed -->
                <div v-if="link.type === 'embed'" class="form-control">
                  <label class="label">Embed URL</label>
                  <input 
                    v-model="link.url" 
                    type="text" 
                    class="input input-bordered" 
                    placeholder="Paste embed URL here"
                    required
                  />
                </div>
                
                <!-- File Upload - Hanya tampilkan untuk self-hosted -->
                <div v-if="link.type === 'self-hosted'" class="form-control">
                  <label class="label">Video File</label>
                  <input 
                    type="file" 
                    @change="(e) => episodeManagement.handleVideoUpload(e, index)"
                    accept="video/*"
                    class="file-input file-input-bordered w-full"
                    required
                  />
                  <!-- Preview URL jika sudah diupload -->
                  <div v-if="link.url" class="mt-2 text-sm text-success">
                    Video uploaded: {{ link.url }}
                  </div>
                </div>
              </div>
            </div>
            
            <!-- Download Links -->
            <div class="space-y-2">
              <div class="flex justify-between items-center">
                <h4 class="font-medium">Download Links</h4>
                <button 
                  type="button"
                  @click="episodeManagement.addDownloadLink"
                  class="btn btn-sm btn-primary"
                >
                  Add Download Link
                </button>
              </div>
              
              <div v-for="(link, index) in episodeForm.downloadLinks" :key="index" class="card bg-base-200 p-4">
                <div class="flex justify-between items-start mb-2">
                  <h5 class="font-medium">Download Link #{{ index + 1 }}</h5>
                  <button 
                    type="button"
                    @click="episodeManagement.removeDownloadLink(index)"
                    class="btn btn-sm btn-error"
                  >
                    Remove
                  </button>
                </div>
                
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-2">
                  <div class="form-control">
                    <label class="label">Name</label>
                    <input 
                      v-model="link.name" 
                      type="text" 
                      class="input input-bordered" 
                      placeholder="Server name"
                      required
                    />
                  </div>
                  
                  <div class="form-control">
                    <label class="label">Quality</label>
                    <input 
                      v-model="link.quality" 
                      type="text" 
                      class="input input-bordered" 
                      placeholder="e.g. 720p"
                      required
                    />
                  </div>
                </div>
                
                <div class="form-control">
                  <label class="label">URL</label>
                  <input 
                    v-model="link.url" 
                    type="text" 
                    class="input input-bordered" 
                    placeholder="Download URL"
                    required
                  />
                </div>
              </div>
            </div>
          </div>
          
          <div class="modal-action">
            <button type="button" @click="episodeManagement.closeEpisodeModal" class="btn">
              Cancel
            </button>
            <button type="submit" class="btn btn-primary">
              Save Episode
            </button>
          </div>
        </form>
      </div>
    </dialog>
  </div>
</template>

<style scoped>
/* Tambahkan CSS untuk memastikan layout responsif */
@media (min-width: 1024px) {
  .container {
    max-width: 1280px;
  }
}

.modal-box {
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
}

/* Perbaikan untuk table di mobile */
.table {
  @apply min-w-full;
}

/* Perbaikan untuk form fields */
.form-control {
  @apply w-full;
}

/* Perbaikan untuk input groups */
.input-group {
  @apply flex flex-col sm:flex-row gap-2;
}
</style>

<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import router from '@/router'

// State
const activeTab = ref('contents')
const contents = ref([])
const genres = ref([])
const categories = ref([])

const showAddContentModal = ref(false)
const showAddGenreModal = ref(false)
const showAddCategoryModal = ref(false)
const showEpisodeModal = ref(false)
const showAddSeasonModal = ref(false)

const contentForm = ref({
  title: '',
  description: '',
  type: '',
  releaseDate: '',
  rating: 0,
  genreIds: [],
  episodes: [],
  coverImage: null,
  seasonId: ''
})

const editingContent = ref(null)
const editingGenre = ref(null)
const editingCategory = ref(null)
const editingEpisodeIndex = ref(null)
const editingSeasonId = ref(null)

const coverImageInput = ref(null)
const coverImagePreview = ref(null)
const selectedGenres = computed(() => {
  return contentForm.value.genreIds.map(id => 
    genres.value.find(g => g.id === id)
  ).filter(Boolean)
})

const episodeForm = ref({
  title: '',
  description: '',
  episodeNumber: 1,
  seasonNumber: 1,
  streamLinks: [],
  downloadLinks: []
})

const seasonForm = ref({
  name: '',
  year: '',
  status: ''
})

const tempVideoFiles = ref(new Map()) // Menyimpan file video sementara dengan format: "seasonNumber-episodeNumber-index" => File

const editContent = (content) => {
  console.log('Editing content:', content)
  editingContent.value = content
  
  // Proses episodes, stream_links, dan download_links
  let processedEpisodes = []
  
  if (content.episodes && content.episodes.length > 0) {
    processedEpisodes = content.episodes.map(episode => {
      // Cari stream links dan download links untuk episode ini
      const episodeStreamLinks = content.stream_links?.filter(link => 
        link.episode_number === episode.episode_number && 
        link.season_number === episode.season_number
      ) || []
      
      const episodeDownloadLinks = content.download_links?.filter(link => 
        link.episode_number === episode.episode_number && 
        link.season_number === episode.season_number
      ) || []
      
      return {
        ...episode,
        streamLinks: episodeStreamLinks.map(link => ({
          name: link.name,
          quality: link.quality,
          url: link.url,
          type: link.type
        })),
        downloadLinks: episodeDownloadLinks.map(link => ({
          name: link.name,
          quality: link.quality,
          url: link.url
        }))
      }
    })
  }
  
  contentForm.value = {
    title: content.title || '',
    description: content.description || '',
    type: content.type || '',
    releaseDate: content.release_date ? new Date(content.release_date).toISOString().split('T')[0] : '',
    rating: content.rating || 0,
    genreIds: content.genres?.map(g => g.id) || [],
    episodes: processedEpisodes,
    coverImage: null,
    seasonId: content.season_id || ''
  }
  console.log('check content', content)
  coverImagePreview.value = getImageUrl(content.cover_image)
  showAddContentModal.value = true
}

const getImageUrl = (cover_image) => {
  if (!cover_image) return
  const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  const cleanPath = cover_image.replace(/\\/g, '/').replace(/^media\//, '')
  return `${baseURL}/media/${cleanPath}`
}

const deleteContent = async (id) => {
  try {
    // Pastikan id valid
    if (!id || isNaN(id)) {
      showErrorMessage('Invalid content ID')
      return
    }

    if (!confirm('Are you sure you want to delete this content?')) {
      return
    }

    console.log('Attempting to delete content:', id)
    
    // Gunakan endpoint yang sesuai dengan backend
    const response = await axios.delete(`/api/contents/${id}`, getAxiosConfig())
    
    if (response.status === 200) {
      await fetchContents()
      showErrorMessage('Content deleted successfully')
    }
  } catch (error) {
    console.error('Failed to delete content:', error)
    showErrorMessage(error.response?.data?.error || 'Failed to delete content')
  }
}

const handleContentSubmit = async () => {
  try {
    console.log('Submitting content form:', contentForm.value)
    
    const formData = new FormData()
    
    // Append basic content info
    formData.append('title', contentForm.value.title)
    formData.append('description', contentForm.value.description)
    formData.append('type', contentForm.value.type)
    
    if (contentForm.value.rating) {
      formData.append('rating', contentForm.value.rating)
    }
    
    if (contentForm.value.releaseDate) {
      formData.append('releaseDate', new Date(contentForm.value.releaseDate).toISOString())
    }

    if (contentForm.value.seasonId) {
      formData.append('season_id', contentForm.value.seasonId)
    }
    
    // Append genre IDs
    if (contentForm.value.genreIds && contentForm.value.genreIds.length > 0) {
      contentForm.value.genreIds.forEach((id, index) => {
        formData.append(`genreIds[]`, id.toString())
      })
      console.log('Genre IDs being sent:', contentForm.value.genreIds)
    }
    
    if (contentForm.value.coverImage) {
      formData.append('coverImage', contentForm.value.coverImage)
    }

    // Proses episodes dan video files
    const processedEpisodes = contentForm.value.episodes.map(episode => {
      const processedEpisode = { ...episode }
      
      // Proses stream links
      processedEpisode.streamLinks = episode.streamLinks.map((link, index) => {
        const processedLink = { ...link }
        
        // Jika ini adalah self-hosted video, tambahkan file ke FormData
        if (link.type === 'self-hosted' && link.tempKey) {
          const videoFile = tempVideoFiles.value.get(link.tempKey)
          if (videoFile) {
            const fieldName = `video_${episode.seasonNumber}_${episode.episodeNumber}_${index}`
            formData.append(fieldName, videoFile)
            processedLink.videoField = fieldName // Referensi ke field name di FormData
          }
          delete processedLink.tempKey // Hapus temporary key
        }
        
        return processedLink
      })
      
      return processedEpisode
    })
    
    formData.append('episodes', JSON.stringify(processedEpisodes))
    
    // Log FormData untuk debugging
    for (let [key, value] of formData.entries()) {
      console.log(`${key}: ${value instanceof File ? value.name : value}`)
    }
    
    const response = await axios.post('/api/contents/create', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    
    console.log('Content created successfully:', response.data)
    
    // Reset form dan temporary files
    resetContentForm()
    tempVideoFiles.value.clear()
    showAddContentModal.value = false
    await fetchContents()
    
    showErrorMessage('Content created successfully!')
  } catch (error) {
    console.error('Failed to create content:', error)
    showErrorMessage(error.response?.data?.error || 'Failed to create content')
  }
}

// Fungsi untuk reset form
const resetContentForm = () => {
  contentForm.value = {
    title: '',
    description: '',
    type: '',
    releaseDate: '',
    rating: 0,
    genreIds: [],
    episodes: [],
    coverImage: null,
    seasonId: ''
  }
  coverImagePreview.value = null
}

const openContentModal = () => {
  editingContent.value = null
  resetContentForm()
  showAddContentModal.value = true
}

const genreForm = ref({
  name: '',
  description: ''
})

const categoryForm = ref({
  name: '',
  description: ''
})

// Methods
const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const fetchContents = async () => {
  try {
    const response = await axios.get('/api/contents', getAxiosConfig())
    contents.value = response.data.contents
  } catch (error) {
    console.error('Failed to fetch contents:', error)
    showErrorMessage('Failed to fetch contents')
  }
}

const fetchGenres = async () => {
  try {
    const response = await axios.get('/api/genres', getAxiosConfig())
    genres.value = response.data
  } catch (error) {
    console.error('Failed to fetch genres:', error)
    showErrorMessage('Failed to fetch genres')
  }
}

const fetchCategories = async () => {
  try {
    const response = await axios.get('/api/categories', getAxiosConfig())
    categories.value = response.data
  } catch (error) {
    console.error('Failed to fetch categories:', error)
    showErrorMessage('Failed to fetch categories')
  }
}

const handleGenreSubmit = async () => {
  try {
    const user = JSON.parse(localStorage.getItem('user'))
    const token = localStorage.getItem('token')
    
    console.log('Current user:', user)
    console.log('Current token:', token)
    console.log('Submitting genre:', genreForm.value)
    
    if (!token || !user || user.role !== 'admin') {
      showErrorMessage('You do not have admin privileges')
      router.push('/login')
      return
    }
    
    const config = getAxiosConfig()
    console.log('Request config:', config)
    
    if (editingGenre.value) {
      const response = await axios.put(
        `/api/genres/${editingGenre.value.id}`, 
        genreForm.value,
        config
      )
      console.log('Genre updated:', response)
    } else {
      const response = await axios.post('/api/genres', genreForm.value, config)
      console.log('Genre created:', response)
    }
    
    showAddGenreModal.value = false
    await fetchGenres()
    genreForm.value = { name: '', description: '' }
    showErrorMessage('Genre saved successfully!')
  } catch (error) {
    console.error('Failed to save genre:', {
      error,
      response: error.response?.data,
      status: error.response?.status,
      headers: error.config?.headers
    })
    
    if (error.response?.status === 403) {
      showErrorMessage('Session expired. Please login again.')
      // localStorage.removeItem('token')
      // localStorage.removeItem('user')
      // forceLogout()
    } else {
      showErrorMessage(error.response?.data?.error || 'Failed to save genre')
    }
  }
}
const forceLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  router.push('/login')
}
const handleCategorySubmit = async () => {
  try {
    const user = JSON.parse(localStorage.getItem('user'))
    const token = localStorage.getItem('token')
    
    console.log('Current user:', user)
    console.log('Current token:', token)
    console.log('Submitting category:', categoryForm.value)
    
    if (!token || !user || user.role !== 'admin') {
      showErrorMessage('You do not have admin privileges')
      router.push('/login')
      return
    }
    
    const config = getAxiosConfig()
    console.log('Request config:', config)
    
    if (editingCategory.value) {
      const response = await axios.put(
        `/api/categories/${editingCategory.value.id}`, 
        categoryForm.value,
        config
      )
      console.log('Category updated:', response.data)
    } else {
      const response = await axios.post('/api/categories', categoryForm.value, config)
      console.log('Category created:', response.data)
    }
    
    showAddCategoryModal.value = false
    await fetchCategories()
    categoryForm.value = { name: '', description: '' }
    showErrorMessage('Category saved successfully!')
  } catch (error) {
    console.error('Failed to save category:', {
      error,
      response: error.response?.data,
      status: error.response?.status,
      headers: error.config?.headers
    })
    
    if (error.response?.status === 403) {
      showErrorMessage('Session expired. Please login again.')
      // localStorage.removeItem('token')
      // localStorage.removeItem('user')
      // router.push('/login')
    } else {
      showErrorMessage(error.response?.data?.error || 'Failed to save category')
    }
  }
}

// Methods untuk genre
const editGenre = (genre) => {
  editingGenre.value = genre
  genreForm.value = { ...genre }
  showAddGenreModal.value = true
}

const deleteGenre = async (id) => {
  if (confirm('Are you sure you want to delete this genre?')) {
    try {
      await axios.delete(`/api/genres/${id}`, getAxiosConfig())
      await fetchGenres()
    } catch (error) {
      console.error('Failed to delete genre:', error)
    }
  }
}

// Methods untuk category
const editCategory = (category) => {
  editingCategory.value = category
  categoryForm.value = { ...category }
  showAddCategoryModal.value = true
}

const deleteCategory = async (id) => {
  if (confirm('Are you sure you want to delete this category?')) {
    try {
      await axios.delete(`/api/categories/${id}`, getAxiosConfig())
      await fetchCategories()
    } catch (error) {
      console.error('Failed to delete category:', error)
    }
  }
}

// Lifecycle
const initializeContentForm = () => {
  contentForm.value = {
    title: '',
    description: '',
    type: '',
    releaseDate: '',
    rating: 0,
    genreIds: [],
    episodes: [],
    coverImage: null,
    seasonId: ''
  }
}

onMounted(async () => {
  initializeContentForm()
  try {
    const user = JSON.parse(localStorage.getItem('user'))
    const token = localStorage.getItem('token')
    
    console.log('Mounting Admin component:', { user, token })
    
    if (!token || !user) {
      showErrorMessage('No authentication data found')
      router.push('/login')
      return
    }
    
    if (user.role !== 'admin') {
      showErrorMessage('You do not have admin privileges')
      router.push('/')
      return
    }

    // Verifikasi akses admin
    try {
      await axios.get('/api/admin/verify', getAxiosConfig())
      console.log('Admin access verified')
    } catch (error) {
      console.error('Admin verification failed:', error)
      if (error.response?.status === 403) {
        showErrorMessage('Admin access required. Please login again.')
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        router.push('/login')
        return
      }
    }
    
    // Fetch data secara parallel
    await Promise.all([
      fetchContents(),
      fetchGenres(),
      fetchCategories(),
      fetchSeasons()
    ])
  } catch (error) {
    console.error('Failed to initialize admin page:', error)
    showErrorMessage(error.message)
  }
})

// Tambahkan fungsi helper untuk mendapatkan config axios
const getAxiosConfig = () => {
  const token = localStorage.getItem('token')
  const user = JSON.parse(localStorage.getItem('user'))
  
  if (!token || !user) {
    showErrorMessage('No authentication data found')
    router.push('/login')
    throw new Error('No authentication data found')
  }
  
  if (user.role !== 'admin') {
    showErrorMessage('You do not have admin privileges')
    router.push('/')
    throw new Error('Admin privileges required')
  }
  
  return {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
      'Accept': 'application/json'
    },
    withCredentials: true
  }
}

// State untuk error handling
const error = ref(null)
const showError = ref(false)

// Fungsi untuk menampilkan error
const showErrorMessage = (message) => {
  error.value = message
  showError.value = true
  setTimeout(() => {
    showError.value = false
  }, 5000)
}

const handleCoverImageChange = (event) => {
  const file = event.target.files[0]
  if (file) {
    contentForm.value.coverImage = file
    coverImagePreview.value = URL.createObjectURL(file)
  }
}

// Fungsi untuk Episode Management
const episodeManagement = {
  // 1. Membuka modal untuk menambah episode baru
  openAddEpisodeModal() {
    editingEpisodeIndex.value = null
    episodeForm.value = {
      title: '',
      description: '',
      episodeNumber: contentForm.value.episodes.length + 1,
      seasonNumber: 1,
      streamLinks: [],
      downloadLinks: []
    }
    showEpisodeModal.value = true
  },

  // 2. Membuka modal untuk mengedit episode
  openEditEpisodeModal(episode, index) {
    editingEpisodeIndex.value = index
    episodeForm.value = JSON.parse(JSON.stringify(episode)) // Deep copy
    showEpisodeModal.value = true
  },

  // 3. Menutup modal episode
  closeEpisodeModal() {
    showEpisodeModal.value = false
    editingEpisodeIndex.value = null
    episodeForm.value = {
      title: '',
      description: '',
      episodeNumber: 1,
      seasonNumber: 1,
      streamLinks: [],
      downloadLinks: []
    }
  },

  // 4. Menyimpan episode (baik tambah baru atau edit)
  async saveEpisode() {
    try {
      const episodeData = {
        title: episodeForm.value.title,
        description: episodeForm.value.description,
        episodeNumber: episodeForm.value.episodeNumber,
        seasonNumber: episodeForm.value.seasonNumber,
        streamLinks: episodeForm.value.streamLinks || [],
        downloadLinks: episodeForm.value.downloadLinks || []
      }
      
      if (editingEpisodeIndex.value !== null) {
        // Update existing episode
        contentForm.value.episodes[editingEpisodeIndex.value] = episodeData
      } else {
        // Add new episode
        if (!Array.isArray(contentForm.value.episodes)) {
          contentForm.value.episodes = []
        }
        contentForm.value.episodes.push(episodeData)
      }
      
      this.closeEpisodeModal()
      this.reorderEpisodes()
    } catch (error) {
      console.error('Failed to save episode:', error)
      showErrorMessage('Failed to save episode')
    }
  },

  // 5. Menghapus episode
  async removeEpisode(index) {
    try {
      if (confirm('Are you sure you want to remove this episode?')) {
        contentForm.value.episodes.splice(index, 1)
        await this.reorderEpisodes()
      }
    } catch (error) {
      console.error('Failed to remove episode:', error)
      showErrorMessage('Failed to remove episode')
    }
  },

  // 6. Mengurutkan ulang nomor episode
  reorderEpisodes() {
    try {
      if (!Array.isArray(contentForm.value.episodes)) {
        contentForm.value.episodes = []
        return
      }
      
      contentForm.value.episodes.sort((a, b) => {
        if (a.seasonNumber === b.seasonNumber) {
          return a.episodeNumber - b.episodeNumber
        }
        return a.seasonNumber - b.seasonNumber
      })
    } catch (error) {
      console.error('Failed to reorder episodes:', error)
      showErrorMessage('Failed to reorder episodes')
    }
  },

  // 7. Manajemen Stream Links
  addStreamLink() {
    episodeForm.value.streamLinks.push({
      name: '',
      quality: '',
      url: '',
      type: 'embed' // Default ke embed
    })
  },

  removeStreamLink(index) {
    episodeForm.value.streamLinks.splice(index, 1)
  },

  // 8. Manajemen Download Links
  addDownloadLink() {
    episodeForm.value.downloadLinks.push({
      name: '',
      quality: '',
      url: ''
    })
  },

  removeDownloadLink(index) {
    episodeForm.value.downloadLinks.splice(index, 1)
  },

  // 9. Menangani upload video
  async handleVideoUpload(event, index) {
    const file = event.target.files[0]
    if (!file) return
    
    try {
      // Buat preview URL untuk file
      const previewUrl = URL.createObjectURL(file)
      
      // Simpan file sementara dengan key unik
      const key = `${episodeForm.value.seasonNumber}-${episodeForm.value.episodeNumber}-${index}`
      tempVideoFiles.value.set(key, file)
      
      // Update URL dengan preview sementara
      episodeForm.value.streamLinks[index].url = previewUrl
      episodeForm.value.streamLinks[index].tempKey = key // Tambahkan key untuk referensi nanti
      
      showErrorMessage('Video file selected and ready to upload')
    } catch (error) {
      console.error('Failed to process video:', error)
      showErrorMessage('Failed to process video file')
      episodeForm.value.streamLinks[index].url = ''
    }
  }
}

const seasons = ref([])

const editSeason = (season) => {
  editingSeasonId.value = season.id
  seasonForm.value = {
    name: season.name,
    year: season.year,
    status: season.status
  }
  showAddSeasonModal.value = true
}

const deleteSeason = async (id) => {
  if (confirm('Are you sure you want to delete this season?')) {
    try {
      await axios.delete(`/api/seasons/${id}`, getAxiosConfig())
      await fetchSeasons()
    } catch (error) {
      console.error('Failed to delete season:', error)
    }
  }
}

const saveSeason = async () => {
  try {
    const user = JSON.parse(localStorage.getItem('user'))
    const token = localStorage.getItem('token')
    
    console.log('Current user:', user)
    console.log('Current token:', token)
    console.log('Submitting season:', seasonForm.value)
    
    if (!token || !user || user.role !== 'admin') {
      showErrorMessage('You do not have admin privileges')
      router.push('/login')
      return
    }
    
    const config = getAxiosConfig()
    console.log('Request config:', config)
    
    if (editingSeasonId.value) {
      const response = await axios.put(
        `/api/seasons/${editingSeasonId.value}`, 
        seasonForm.value,
        config
      )
      console.log('Season updated:', response)
    } else {
      const response = await axios.post('/api/seasons', seasonForm.value, config)
      console.log('Season created:', response)
    }
    
    showAddSeasonModal.value = false
    await fetchSeasons()
    seasonForm.value = { name: '', year: '', status: '' }
    showErrorMessage('Season saved successfully!')
  } catch (error) {
    console.error('Failed to save season:', {
      error,
      response: error.response?.data,
      status: error.response?.status,
      headers: error.config?.headers
    })
    
    if (error.response?.status === 403) {
      showErrorMessage('Session expired. Please login again.')
      // localStorage.removeItem('token')
      // localStorage.removeItem('user')
      // router.push('/login')
    } else {
      showErrorMessage(error.response?.data?.error || 'Failed to save season')
    }
  }
}

const closeSeasonModal = () => {
  showAddSeasonModal.value = false
  editingSeasonId.value = null
  seasonForm.value = { name: '', year: '', status: '' }
}

const fetchSeasons = async () => {
  try {
    const response = await axios.get('/api/seasons', getAxiosConfig())
    seasons.value = response.data
  } catch (error) {
    console.error('Failed to fetch seasons:', error)
    showErrorMessage('Failed to fetch seasons')
  }
}
</script>