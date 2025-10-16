import { defineStore } from 'pinia'
import { ref } from 'vue'
import { repositoryService } from '@/services'

export const useRepositoryStore = defineStore('repository', () => {
  const repositories = ref([])
  const currentRepository = ref(null)
  const loading = ref(false)
  const error = ref(null)

  async function fetchRepositories() {
    loading.value = true
    error.value = null
    try {
      const data = await repositoryService.getRepositories()
      repositories.value = data
      return data
    } catch (err) {
      error.value = err
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchRepository(id) {
    loading.value = true
    error.value = null
    try {
      const data = await repositoryService.getRepository(id)
      currentRepository.value = data
      return data
    } catch (err) {
      error.value = err
      throw err
    } finally {
      loading.value = false
    }
  }

  async function createRepository(repoData) {
    loading.value = true
    error.value = null
    try {
      const data = await repositoryService.createRepository(repoData)
      repositories.value.push(data)
      return data
    } catch (err) {
      error.value = err
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateRepository(id, repoData) {
    loading.value = true
    error.value = null
    try {
      const data = await repositoryService.updateRepository(id, repoData)
      const index = repositories.value.findIndex(r => r.id === id)
      if (index !== -1) {
        repositories.value[index] = data
      }
      if (currentRepository.value?.id === id) {
        currentRepository.value = data
      }
      return data
    } catch (err) {
      error.value = err
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteRepository(id) {
    loading.value = true
    error.value = null
    try {
      await repositoryService.deleteRepository(id)
      repositories.value = repositories.value.filter(r => r.id !== id)
      if (currentRepository.value?.id === id) {
        currentRepository.value = null
      }
    } catch (err) {
      error.value = err
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    repositories,
    currentRepository,
    loading,
    error,
    fetchRepositories,
    fetchRepository,
    createRepository,
    updateRepository,
    deleteRepository
  }
})
