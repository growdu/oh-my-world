<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'
import { marked } from 'marked'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

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
  let finalUrl = url
  if (!url.startsWith('http://') && !url.startsWith('https://')) {
    finalUrl = 'http://' + url
  }
  //window.location.href = finalUrl
  window.open(finalUrl, '_blank')
}

const exitLogin = () => {
  localStorage.removeItem('isAdmin')
  isAdmin.value = false
  router.push('/')
}

const renderMarkdown = (text) => {
  return marked(text)
}

// 监听路由变化
watch(
  () => router.currentRoute.value,
  () => {
    fetchLinks()
  }
)

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

      <el-dialog v-model="dialogVisible" title="添加新链接" width="800px">
        <el-form>
          <el-form-item label="图片">
            <el-input v-model="newLink.image" placeholder="图片URL" />
            <el-image v-if="newLink.image" :src="newLink.image" class="preview-image" />
          </el-form-item>
          <el-form-item label="描述">
            <div class="markdown-editor-container">
              <MdEditor
                v-model="newLink.description"
                preview
                language="zh-CN"
                :toolbars="[
                  'bold',
                  'underline',
                  'italic',
                  'strikeThrough',
                  'title',
                  'sub',
                  'sup',
                  'quote',
                  'unorderedList',
                  'orderedList',
                  'codeRow',
                  'code',
                  'link',
                  'image',
                  'table',
                  'preview'
                ]"
              />
            </div>
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
            <div class="card-description markdown-body" v-html="renderMarkdown(link.description)"></div>
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
  font-size: 14px;
  line-height: 1.6;
  padding: 10px;
  width: 100%;
}

/* 添加 Markdown 样式 */
:deep(.markdown-body) {
  color: #2c3e50;
}

:deep(.markdown-body h1) {
  font-size: 1.5em;
  margin-bottom: 0.5em;
}

:deep(.markdown-body h2) {
  font-size: 1.3em;
  margin-bottom: 0.4em;
}

:deep(.markdown-body p) {
  margin-bottom: 0.8em;
}

:deep(.markdown-body a) {
  color: #409EFF;
  text-decoration: none;
}

:deep(.markdown-body code) {
  background-color: #f8f8f8;
  padding: 0.2em 0.4em;
  border-radius: 3px;
  font-family: monospace;
}

:deep(.markdown-body pre) {
  background-color: #f8f8f8;
  padding: 1em;
  border-radius: 4px;
  overflow-x: auto;
}

:deep(.markdown-body ul, .markdown-body ol) {
  padding-left: 2em;
  margin-bottom: 1em;
}

:deep(.markdown-body blockquote) {
  border-left: 4px solid #dfe2e5;
  padding-left: 1em;
  margin: 1em 0;
  color: #666;
}

.markdown-editor-container {
  width: 100%;
  min-height: 300px;
}

:deep(.md-editor) {
  height: 300px;
}

:deep(.md-editor-preview) {
  background-color: #fff;
}
</style>
