import apiClient from '@/utils/api'

export const authService = {
  async register(username, password, email) {
    return await apiClient.post('/auth/register', {
      username,
      password,
      email
    })
  },

  async login(username, password) {
    return await apiClient.post('/auth/login', {
      username,
      password
    })
  },

  async getCurrentUser() {
    return await apiClient.get('/auth/me')
  }
}

export const repositoryService = {
  async getRepositories() {
    return await apiClient.get('/repos')
  },

  async getRepository(id) {
    return await apiClient.get(`/repos/${id}`)
  },

  async createRepository(data) {
    return await apiClient.post('/repos', data)
  },

  async updateRepository(id, data) {
    return await apiClient.put(`/repos/${id}`, data)
  },

  async deleteRepository(id) {
    return await apiClient.delete(`/repos/${id}`)
  },

  async getCollaborators(id) {
    return await apiClient.get(`/repos/${id}/collaborators`)
  },

  async addCollaborator(id, data) {
    return await apiClient.post(`/repos/${id}/collaborators`, data)
  },

  async removeCollaborator(id, userId) {
    return await apiClient.delete(`/repos/${id}/collaborators/${userId}`)
  }
}

export const activityService = {
  async getActivities() {
    return await apiClient.get('/activities')
  },

  async getRepositoryActivities(repoId) {
    return await apiClient.get(`/activities/repo/${repoId}`)
  }
}
