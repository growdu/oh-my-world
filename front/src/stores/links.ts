import { defineStore } from 'pinia'
import { ref } from 'vue'
import { API_BASE_URL } from '../config/env'

interface Link {
  id: number
  name: string
  url: string
  description: string
  categoryId: number
  visitCount: number
}

export const useLinksStore = defineStore('links', () => {
  const links = ref<Link[]>([])

  const fetchLinks = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/links`)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      const data = await response.json()
      links.value = data
    } catch (error) {
      console.error('Error fetching links:', error)
      throw error // 重新抛出错误以便上层处理
    }
  }

  const incrementVisitCount = async (linkId: number) => {
    try {
      const response = await fetch(`${API_BASE_URL}/links/${linkId}/visit`, { 
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        }
      })
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      const link = links.value.find(l => l.id === linkId)
      if (link) {
        link.visitCount++
      }
    } catch (error) {
      console.error('Error incrementing visit count:', error)
      throw error
    }
  }

  return {
    links,
    fetchLinks,
    incrementVisitCount
  }
}) 