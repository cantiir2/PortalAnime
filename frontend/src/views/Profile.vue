<template>
  <div class="min-h-screen p-4">
    <div class="container mx-auto max-w-4xl">
      <div class="card bg-base-200 shadow-xl">
        <div class="card-body">
          <h2 class="card-title mb-4">Profile Settings</h2>
          <form @submit.prevent="handleUpdateProfile">
            <div class="form-control">
              <label class="label">
                <span class="label-text">Username</span>
              </label>
              <input 
                type="text" 
                v-model="profile.username" 
                class="input input-bordered" 
              />
            </div>
            <div class="form-control">
              <label class="label">
                <span class="label-text">Email</span>
              </label>
              <input 
                type="email" 
                v-model="profile.email" 
                class="input input-bordered" 
              />
            </div>
            <div class="form-control mt-6">
              <button type="submit" class="btn btn-primary" :disabled="loading">
                {{ loading ? 'Saving...' : 'Save Changes' }}
              </button>
            </div>
          </form>

          <div class="divider"></div>

          <h3 class="text-lg font-semibold mb-4">Change Password</h3>
          <form @submit.prevent="handleChangePassword">
            <div class="form-control">
              <label class="label">
                <span class="label-text">Current Password</span>
              </label>
              <input 
                type="password" 
                v-model="passwordForm.currentPassword" 
                class="input input-bordered" 
              />
            </div>
            <div class="form-control">
              <label class="label">
                <span class="label-text">New Password</span>
              </label>
              <input 
                type="password" 
                v-model="passwordForm.newPassword" 
                class="input input-bordered" 
              />
            </div>
            <div class="form-control mt-6">
              <button type="submit" class="btn btn-secondary" :disabled="passwordLoading">
                {{ passwordLoading ? 'Changing...' : 'Change Password' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()

const profile = ref({
  username: authStore.user?.username || '',
  email: authStore.user?.email || ''
})

const passwordForm = ref({
  currentPassword: '',
  newPassword: ''
})

const loading = ref(false)
const passwordLoading = ref(false)

const handleUpdateProfile = async () => {
  try {
    loading.value = true
    await authStore.updateProfile(profile.value)
  } catch (error) {
    console.error('Failed to update profile:', error)
  } finally {
    loading.value = false
  }
}

const handleChangePassword = async () => {
  try {
    passwordLoading.value = true
    await authStore.changePassword(
      passwordForm.value.currentPassword,
      passwordForm.value.newPassword
    )
    passwordForm.value = {
      currentPassword: '',
      newPassword: ''
    }
  } catch (error) {
    console.error('Failed to change password:', error)
  } finally {
    passwordLoading.value = false
  }
}
</script> 