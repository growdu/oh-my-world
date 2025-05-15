import { defineStore } from 'pinia'
import { ref } from 'vue'
import { API_BASE_URL } from '../config/env'

interface Category {
  id: number
  name: string
}

export const useCategoriesStore = defineStore('categories', () => {
  const categories = ref<Category[]>([])

  const fetchCategories = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/categories`)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      const data = await response.json()
      categories.value = data
    } catch (error) {
      console.error('Error fetching categories:', error)
      throw error
    }
  }

  return {
    categories,
    fetchCategories
  }
}) 