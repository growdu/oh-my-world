<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'
import { marked } from 'marked'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { Swiper, SwiperSlide } from 'swiper/vue'
import { EffectCoverflow, Mousewheel, Autoplay } from 'swiper/modules'
import 'swiper/css'
import 'swiper/css/effect-coverflow'

const router = useRouter()
const isAdmin = ref(localStorage.getItem('isAdmin') === 'true')

const allLinks = ref([])
const currentPage = ref(1)
const pageSize = ref(8)
const total = ref(0)

const newLink = ref({ image: '', description: '', url: '', category: '' })
const dialogVisible = ref(false)

const swiperModules = [EffectCoverflow, Mousewheel, Autoplay]

const categories = ref([
  { name: '编程语言', count: 0 },
  { name: '开发工具', count: 0 },
  { name: '框架', count: 0 },
  { name: '数据库', count: 0 },
  { name: '运维', count: 0 },
  { name: '前端', count: 0 },
  { name: '后端', count: 0 },
  { name: '人工智能', count: 0 },
])

const popularLinks = computed(() => {
  return [...allLinks.value]
    .sort((a, b) => (b.visits || 0) - (a.visits || 0))
    .slice(0, 5)
})

const selectedCategory = ref('全部')

const filteredLinks = computed(() => {
  if (selectedCategory.value === '全部') return allLinks.value
  return allLinks.value.filter(link => link.category === selectedCategory.value)
})

const updateCategoryCounts = () => {
  categories.value.forEach(category => {
    category.count = allLinks.value.filter(link => link.category === category.name).length
  })
}

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
    updateCategoryCounts()
    return
  }

  try {
    const res = await request.get(`${import.meta.env.VITE_API_BASE_URL}/links`)
    console.log(`${import.meta.env.VITE_API_BASE_URL}/links`)
    
    if (Array.isArray(res.data)) {
      allLinks.value = res.data.reverse()
      total.value = allLinks.value.length
      updateCategoryCounts()
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
    newLink.value = { image: '', description: '', url: '', category: '' }
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
    <!-- 左侧分类栏 -->
    <el-aside width="200px" class="category-aside">
      <h3 class="aside-title">分类导航</h3>
      <el-menu
        :default-active="selectedCategory"
        class="category-menu"
        @select="(index) => selectedCategory = index"
      >
        <el-menu-item index="全部">
          全部
          <span class="category-count">{{ allLinks.length }}</span>
        </el-menu-item>
        <el-menu-item 
          v-for="category in categories" 
          :key="category.name"
          :index="category.name"
        >
          {{ category.name }}
          <span class="category-count">{{ category.count }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-main>
      <h1 class="title">💻 Growdu 的计算机世界 🌍</h1>

      <div class="action-buttons" v-if="isAdmin">
        <el-button type="primary" @click="dialogVisible = true">添加链接</el-button>
        <el-button type="danger" @click="exitLogin">退出登录</el-button>
      </div>

      <div class="swiper-container">
        <swiper
          :modules="swiperModules"
          :effect="'coverflow'"
          :grab-cursor="true"
          :centeredSlides="true"
          :slidesPerView="'auto'"
          :mousewheel="true"
          :autoplay="{
            delay: 3000,
            disableOnInteraction: false,
          }"
          :coverflowEffect="{
            rotate: 50,
            stretch: 0,
            depth: 100,
            modifier: 1,
            slideShadows: true
          }"
          class="mySwiper"
        >
          <swiper-slide v-for="link in filteredLinks" :key="link.id" class="card-slide">
            <div class="card-wrapper" @click="goToLink(link.url)">
              <div class="card-content">
                <div class="image-container">
                  <el-image :src="link.image" fit="cover" />
                </div>
                <div class="content-container">
                  <div class="card-description markdown-body" v-html="renderMarkdown(link.description)"></div>
                  <el-button
                    v-if="isAdmin"
                    size="small"
                    type="danger"
                    class="delete-button"
                    @click.stop="deleteLink(link)"
                  >删除</el-button>
                </div>
              </div>
            </div>
          </swiper-slide>
        </swiper>
      </div>

      <el-dialog v-model="dialogVisible" title="添加新链接" width="800px">
        <el-form>
          <el-form-item label="图片">
            <el-input v-model="newLink.image" placeholder="图片URL" />
            <el-image v-if="newLink.image" :src="newLink.image" class="preview-image" />
          </el-form-item>
          <el-form-item label="分类">
            <el-select v-model="newLink.category" placeholder="请选择分类">
              <el-option
                v-for="category in categories"
                :key="category.name"
                :label="category.name"
                :value="category.name"
              />
            </el-select>
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
    </el-main>

    <!-- 右侧流行度栏 -->
    <el-aside width="200px" class="popular-aside">
      <h3 class="aside-title">热门推荐</h3>
      <div class="popular-list">
        <div 
          v-for="(link, index) in popularLinks" 
          :key="link.id"
          class="popular-item"
          @click="goToLink(link.url)"
        >
          <span class="popular-rank" :class="{'top-three': index < 3}">{{ index + 1 }}</span>
          <el-image 
            :src="link.image" 
            class="popular-image"
            fit="cover"
          />
          <div class="popular-info">
            <div class="popular-title">{{ link.description.split('\n')[0] }}</div>
            <div class="popular-visits">访问次数: {{ link.visits || 0 }}</div>
          </div>
        </div>
      </div>
    </el-aside>
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

.swiper-container {
  width: 60%;
  margin: 0 auto;
  padding: 50px 0;
}

.mySwiper {
  width: 100%;
  padding-top: 20px;
  padding-bottom: 50px;
}

.card-slide {
  background-position: center;
  background-size: cover;
  width: 500px;
  height: 600px;
}

.card-wrapper {
  width: 100%;
  height: 100%;
  background: #fff;
  border-radius: 15px;
  box-shadow: 0 8px 20px rgba(0,0,0,0.1);
  overflow: hidden;
  transition: transform 0.3s ease;
  cursor: pointer;
}

.card-wrapper:hover {
  transform: translateY(-5px);
}

.card-content {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.image-container {
  height: 60%;
  overflow: hidden;
}

.image-container :deep(.el-image) {
  width: 100%;
  height: 100%;
  transition: transform 0.3s ease;
}

.card-wrapper:hover .image-container :deep(.el-image) {
  transform: scale(1.05);
}

.content-container {
  height: 40%;
  padding: 20px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  background: #fff;
}

.card-description {
  flex: 1;
  overflow-y: auto;
  padding-right: 10px;
}

.delete-button {
  align-self: flex-end;
  margin-top: 10px;
}

/* 自定义滚动条样式 */
.card-description::-webkit-scrollbar {
  width: 4px;
}

.card-description::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 2px;
}

.card-description::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 2px;
}

.card-description::-webkit-scrollbar-thumb:hover {
  background: #555;
}

/* Markdown 样式保持不变 */
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

.category-aside,
.popular-aside {
  background: #f5f7fa;
  padding: 20px 0;
  height: 100vh;
  position: fixed;
  top: 0;
  overflow-y: auto;
}

.category-aside {
  left: 0;
  border-right: 1px solid #e6e6e6;
}

.popular-aside {
  right: 0;
  border-left: 1px solid #e6e6e6;
}

.aside-title {
  text-align: center;
  margin-bottom: 20px;
  color: #409EFF;
  font-size: 18px;
}

.category-menu {
  border-right: none;
}

.category-count {
  float: right;
  color: #909399;
  font-size: 12px;
}

.popular-list {
  padding: 0 10px;
}

.popular-item {
  display: flex;
  align-items: center;
  padding: 10px;
  margin-bottom: 10px;
  background: white;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.popular-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 12px 0 rgba(0,0,0,.1);
}

.popular-rank {
  width: 24px;
  height: 24px;
  line-height: 24px;
  text-align: center;
  background: #909399;
  color: white;
  border-radius: 50%;
  margin-right: 10px;
  font-size: 12px;
}

.popular-rank.top-three {
  background: #409EFF;
}

.popular-image {
  width: 40px;
  height: 40px;
  border-radius: 4px;
  margin-right: 10px;
}

.popular-info {
  flex: 1;
  overflow: hidden;
}

.popular-title {
  font-size: 14px;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.popular-visits {
  font-size: 12px;
  color: #909399;
}

/* 调整主内容区域的边距以适应两侧栏 */
.el-main {
  margin: 0 200px;
  padding: 20px 40px;
}

/* 调整 swiper 容器的宽度 */
.swiper-container {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
}
</style>
