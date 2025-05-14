<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()
const isAdmin = ref(localStorage.getItem('isAdmin') === 'true')

const allLinks = ref([])
const currentPage = ref(1)
const pageSize = ref(8)
const total = ref(0)

const newLink = ref({ image: '', description: '', url: '' })
const dialogVisible = ref(false)

const fetchLinks = async () => {
  const useMock = import.meta.env.VITE_USE_MOCK === 'true'

  if (useMock) {
    allLinks.value = [
      {
        id: 1,
        image: 'https://via.placeholder.com/100',
        description: '示例网站 A',
        url: 'https://example.com/a',
      },
      {
        id: 2,
        image: 'https://via.placeholder.com/100',
        description: '示例网站 B',
        url: 'https://example.com/b',
      },
      {
        id: 3,
        image: 'https://via.placeholder.com/100',
        description: '示例网站 C',
        url: 'https://example.com/c',
      },
      {
        id: 4,
        image: 'https://via.placeholder.com/100',
        description: '示例网站 D',
        url: 'https://example.com/d',
      },
      {
        id: 4,
        image: 'https://via.placeholder.com/100',
        description: '示例网站 D',
        url: 'https://example.com/d',
      },
      {
        id: 5,
        image: 'https://via.placeholder.com/100',
        description: '示例网站 D',
        url: 'https://example.com/d',
      },
      {
        id: 6,
        image: 'https://via.placeholder.com/100',
        description: '示例网站 D',
        url: 'https://example.com/d',
      },
    ]
    total.value = allLinks.value.length
    return
  }

  try {
    const res = await request.get(`${import.meta.env.VITE_API_BASE_URL}/links`)
    console.log(`${import.meta.env.VITE_API_BASE_URL}/links`)
    
    if (Array.isArray(res.data)) {
      allLinks.value = res.data.reverse()
      total.value = allLinks.value.length
    } else {
      console.error('返回的数据不是数组', res.data)
      ElMessage.error('返回的数据格式不正确')
    }
  } catch (e) {
    console.log(`get ${import.meta.env.VITE_API_BASE_URL}/links`)
    console.error('请求出错', e)
    ElMessage.error('获取链接失败')
  }
}


const currentLinks = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return allLinks.value.slice(start, start + pageSize.value)
})

const addLink = async () => {
  const link = newLink.value
  if (!link.image || !link.description || !link.url) {
    ElMessage.error('请填写完整信息！')
    return
  }
  try {
    await request.post('/links', link)
    ElMessage.success('添加成功！')
    dialogVisible.value = false
    newLink.value = { image: '', description: '', url: '' }
    fetchLinks()
  } catch {
    ElMessage.error('添加失败')
  }
}

const deleteLink = async (link) => {
  try {
    await ElMessageBox.confirm('确定要删除此链接吗？', '提示', { type: 'warning' })
    await request.delete(`/links/${link.id}`)
    ElMessage.success('删除成功')
    fetchLinks()
  } catch {
    ElMessage.error('删除失败')
  }
}

const goToLink = (url) => {
  window.open(url, '_blank')
}

const exitLogin = () => {
  localStorage.removeItem('isAdmin')
}

onMounted(fetchLinks)
</script>

<template>
  <el-container>
    <el-main>
      <h1 class="title">💻 Growdu 的计算机世界 🌍</h1>

      <div class="action-buttons" v-if="isAdmin">
        <el-button type="primary" @click="dialogVisible = true">添加链接</el-button>
        <el-button type="danger" @click="exitLogin">退出登录</el-button>
      </div>

      <el-dialog v-model="dialogVisible" title="添加新链接" width="400px">
        <el-form>
          <el-form-item label="图片">
            <el-input v-model="newLink.image" placeholder="图片URL" />
            <el-image v-if="newLink.image" :src="newLink.image" class="preview-image" />
          </el-form-item>
          <el-form-item label="描述">
            <el-input v-model="newLink.description" placeholder="描述" />
          </el-form-item>
          <el-form-item label="链接">
            <el-input v-model="newLink.url" placeholder="URL" />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="addLink">添加</el-button>
        </template>
      </el-dialog>

      <div class="card-container">
        <el-card
          v-for="link in currentLinks"
          :key="link.id"
          class="card"
          @click="goToLink(link.url)"
        >
          <div class="card-content">
            <el-image v-if="link.image" :src="link.image" class="card-image" fit="cover" />
            <p class="card-description">{{ link.description }}</p>
            <el-button
              v-if="isAdmin"
              size="small"
              type="danger"
              @click.stop="deleteLink(link)"
            >删除</el-button>
          </div>
        </el-card>
      </div>

      <div style="text-align:center;margin-top:20px">
        <el-pagination
          background
          layout="prev, pager, next"
          :total="total"
          :page-size="pageSize"
          v-model:current-page="currentPage"
        />
      </div>
    </el-main>
  </el-container>
</template>

<style scoped>
.title {
  text-align: center;
  font-size: 24px;
  margin-bottom: 20px;
}
.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-bottom: 20px;
}
.card-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 20px;
}
.card-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}
.card-image {
  width: 100px;
  height: 100px;
  border-radius: 10px;
}
.card-description {
  font-size: 16px;
  font-weight: bold;
}
</style>
