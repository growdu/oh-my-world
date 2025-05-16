<script setup>
import { ref, computed, onMounted, watch, onUnmounted, nextTick } from 'vue'
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
import { ElDialog, ElForm, ElFormItem, ElInput, ElSelect, ElOption } from 'element-plus'

const router = useRouter()
const isAdmin = ref(localStorage.getItem('isAdmin') === 'true')
const loading = ref(true)
const currentIndex = ref(0)

const allLinks = ref([])
const total = ref(0)

const visibleLinks = computed(() => {
  return filteredLinks.value
})

const getCardStyle = (index) => {
  const diff = index - currentIndex.value
  
  if (Math.abs(diff) >= 3) {
    return {
      display: 'none'
    }
  }

  let transform = ''
  let opacity = 1
  let filter = 'blur(0px)'
  let zIndex = 0

  switch (diff) {
    case -2:
      transform = 'translateX(-120%) translateZ(-200px) rotateY(25deg)'
      opacity = 0.6
      filter = 'blur(2px)'
      break
    case -1:
      transform = 'translateX(-60%) translateZ(-100px) rotateY(15deg)'
      opacity = 0.8
      filter = 'blur(1px)'
      break
    case 0:
      transform = 'translateZ(0) rotateY(0deg)'
      opacity = 1
      filter = 'blur(0)'
      zIndex = 1
      break
    case 1:
      transform = 'translateX(60%) translateZ(-100px) rotateY(-15deg)'
      opacity = 0.8
      filter = 'blur(1px)'
      break
    case 2:
      transform = 'translateX(120%) translateZ(-200px) rotateY(-25deg)'
      opacity = 0.6
      filter = 'blur(2px)'
      break
  }

  return {
    transform,
    opacity,
    filter,
    zIndex
  }
}

const handleCardClick = async (link, index) => {
  currentIndex.value = index
  await new Promise(resolve => setTimeout(resolve, 600))
  goToLink(link)
}

const prevCard = () => {
  if (currentIndex.value > 0) {
    currentIndex.value--
  }
}

const nextCard = () => {
  if (currentIndex.value < allLinks.value.length - 1) {
    currentIndex.value++
  }
}

const handleKeydown = (e) => {
  if (e.key === 'ArrowLeft') {
    e.preventDefault()
    prevCard()
  } else if (e.key === 'ArrowRight') {
    e.preventDefault()
    nextCard()
  }
}

const swiperModules = [EffectCoverflow, Mousewheel, Autoplay]

const categories = ref([{ name: '全部', count: 0 }])
const fetchCategories = async () => {
  try {
    const res = await request.get('/categories')
    if (Array.isArray(res.data)) {
      categories.value = [{ name: '全部', count: allLinks.value.length }, ...res.data.map(cat => ({ name: cat.name, count: 0 }))]
      updateCategoryCounts()
    }
  } catch (e) {
    ElMessage.error('获取分类失败')
  }
}

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
  const all = allLinks.value
  const categoryMap = new Map()

  all.forEach(link => {
    const cat = link.category || '未分类'
    categoryMap.set(cat, (categoryMap.get(cat) || 0) + 1)
  })

  categories.value = [{ name: '全部', count: all.length }]
  for (const [name, count] of categoryMap.entries()) {
    categories.value.push({ name, count })
  }
}

const fetchLinks = async () => {
  loading.value = true
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
    loading.value = false
    return
  }

  try {
    const res = await request.get(`${import.meta.env.VITE_API_BASE_URL}/links`)
    console.log(`${import.meta.env.VITE_API_BASE_URL}/links`)
    
    if (Array.isArray(res.data)) {
      allLinks.value = res.data.reverse()
      total.value = allLinks.value.length
      await fetchCategories()
      updateCategoryCounts()
    } else {
      console.error('返回的数据不是数组', res.data)
      ElMessage.error('返回的数据格式不正确')
    }
  } catch (e) {
    console.log(`get ${import.meta.env.VITE_API_BASE_URL}/links`)
    console.error('请求出错', e)
    ElMessage.error('获取链接失败')
  } finally {
    loading.value = false
  }
}

const goToLink = async (link) => {
  try {
    await request.post(`/links/${link.id}/visit`)
  } catch (error) {
    console.error('Failed to update visit count:', error)
  }
  let finalUrl = link.url
  if (!finalUrl.startsWith('http://') && !finalUrl.startsWith('https://')) {
    finalUrl = 'http://' + finalUrl
  }
  window.open(finalUrl, '_blank')
}

const handleLogout = () => {
  ElMessageBox.confirm('确认退出登录？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    localStorage.removeItem('isAdmin')
    isAdmin.value = false
    ElMessage.success('已退出登录')
  }).catch(() => {})
}

const renderMarkdown = (text) => {
  return marked(text)
}

watch(
  () => router.currentRoute.value,
  () => {
    fetchLinks()
  }
)

// 节流函数
const throttle = (fn, delay) => {
  let last = 0
  return (...args) => {
    const now = Date.now()
    if (now - last > delay) {
      last = now
      fn.apply(null, args)
    }
  }
}

// 处理滚轮事件
const handleWheel = throttle((e) => {
  e.preventDefault()
  if (e.deltaY > 0) {
    nextCard()
  } else {
    prevCard()
  }
}, 200) // 200ms 节流

// 处理图片加载错误
const handleImageError = (e) => {
  // 设置默认图片
  e.target.src = 'https://via.placeholder.com/300x200?text=Image+Not+Found'
}

onMounted(() => {
  fetchLinks()
  window.addEventListener('keydown', handleKeydown)
  // 添加滚轮事件监听
  const content = document.querySelector('.content')
  if (content) {
    content.addEventListener('wheel', handleWheel, { passive: false })
  }

  if (localStorage.getItem('isAdmin') !== 'true') {
    router.push('/')
  }
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
  // 移除滚轮事件监听
  const content = document.querySelector('.content')
  if (content) {
    content.removeEventListener('wheel', handleWheel)
  }
})
</script>

<template>
  <div class="home-container" tabindex="0" @keydown="handleKeydown">
    <header class="main-header">
      <div class="header-left">
        <div class="logo-title">
          <img class="growdu-logo" src="/logo.svg" alt="growdu logo" />
          <div class="title-group">
            <h1 class="site-title">growdu</h1>
            <p class="site-explore">持续成长进步，探索计算机世界,一切从好玩开始。</p>
          </div>
        </div>
      </div>

      <div class="header-actions" v-if="isAdmin">
        <button class="action-btn logout-btn" @click="handleLogout">
          <i class="el-icon-switch-button">⇲</i>
          退出登录
        </button>
        <button class="action-btn manage-btn" @click="router.push('/manage')">
          <i class="el-icon-setting">⚙</i>
          卡片管理
        </button>
      </div>
      <!-- 右上角GitHub链接 -->
      <a class="github-link" href="https://github.com/growdu" target="_blank" title="GitHub">
        <svg height="28" viewBox="0 0 16 16" width="28" style="vertical-align:middle;"><path fill="currentColor" d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.01.08-2.11 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.91.08 2.11.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.19 0 .21.15.46.55.38A8.013 8.013 0 0 0 16 8c0-4.42-3.58-8-8-8z"></path></svg>
      </a>
    </header>
    
    <div class="main-content">
      <aside class="sidebar-left">
        <nav class="categories">
          <h2>分类</h2>
          <ul>
            <li v-for="category in categories" :key="category.name">
              <a :class="{ active: selectedCategory === category.name }" @click="selectedCategory = category.name">
                {{ category.name }}
                <span class="category-count">({{ category.count }})</span>
              </a>
            </li>
          </ul>
        </nav>
      </aside>

      <main class="content">
        <div class="cards-container">
          <template v-if="loading">
            <div v-for="n in 5" :key="n" class="link-card skeleton">
              <div class="card-content">
                <div class="skeleton-title"></div>
                <div class="skeleton-description"></div>
                <div class="card-footer">
                  <div class="skeleton-visits"></div>
                  <div class="skeleton-category"></div>
                </div>
              </div>
            </div>
          </template>
          <template v-else>
            <button class="nav-button prev" 
              @click.stop="prevCard" 
              :disabled="currentIndex === 0">
              <i class="el-icon-arrow-left">◀</i>
            </button>
            
            <div v-for="(link, index) in visibleLinks" 
              :key="link.id" 
              class="link-card"
              :class="{ active: index === currentIndex }"
              :style="getCardStyle(index)"
              @click.stop="handleCardClick(link, index)">
              <div class="card-content">
                <div class="card-image">
                  <img 
                    :src="link.image" 
                    :alt="link.name"
                    @error="handleImageError"
                    loading="lazy"
                  />
                </div>
                <div class="card-info">
                  <h3>{{ link.name }}</h3>
                  <p class="link-description-full" v-html="renderMarkdown(link.description)"></p>
                  <div class="card-footer">
                    <span class="visit-count">访问: {{ link.visits || 0 }}</span>
                    <span class="category-tag">{{ link.category }}</span>
                  </div>
                </div>
              </div>
            </div>

            <button class="nav-button next" 
              @click.stop="nextCard" 
              :disabled="currentIndex >= allLinks.length - 1">
              <i class="el-icon-arrow-right">▶</i>
            </button>
          </template>
        </div>
        <!-- 评论功能挂载点 -->
        <!-- <div id="comments" class="comments-section"> -->
          <!-- 这里可以挂载评论组件，如valine、giscus等 -->
        <!-- </div> -->
        <div class="scroll-hint" :class="{ visible: !loading && allLinks.length > 0 }">
          <i>↕</i>
          <span>滚动鼠标滚轮切换卡片</span>
        </div>
      </main>

      <aside class="sidebar-right">
        <section class="hot-links">
          <h2>热门推荐</h2>
          <ul>
            <li v-for="link in popularLinks" :key="link.id">
              <a @click="goToLink(link)">{{ link.name }}</a>
              <span class="visit-count">({{ link.visits || 0 }})</span>
            </li>
          </ul>
        </section>
      </aside>
    </div>

    <!-- 页脚备案和版权信息 -->
    <footer class="main-footer">
      <div>
        <a href="https://beian.miit.gov.cn/" target="_blank">蜀ICP备2025139986号</a>
        &nbsp;|&nbsp;
        Copyright © {{ new Date().getFullYear() }} growdu
      </div>
    </footer>
  </div>
</template>

<style>
.home-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  outline: none;
}

.main-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 1.2rem 2rem 1.2rem 2rem;
  background: rgba(255, 255, 255, 0.1);
  margin-bottom: 1.5rem;
  position: relative;
}

.header-left {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.2rem;
}

.logo-title {
  display: flex;
  align-items: center;
  gap: 0.7rem;
}

.growdu-logo {
  width: 38px;
  height: 38px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  background: #fff;
}

.title-group {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.1rem;
}

/* .site-title {
  font-size: 1.7rem;
  font-style: italic;
  color: #2c3e50;
  margin: 0;
  font-weight: 800;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
  line-height: 1.1;
} */

.site-title {
  font-size: 2.5rem;
  font-style: italic;
  color: #2c3e50;
  margin: 0;
  font-weight: 800;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
  line-height: 1.2;
}

.site-subtitle {
  font-size: 1.1rem;
  color: #34495e;
  margin: 0.5rem 0 0;
  font-style: italic;
  font-weight: 300;
  line-height: 1.2;
}

.site-explore {
  font-size: 1.02rem;
  color: #1a73e8;
  font-weight: 500;
  margin: 0 0 0 2px;
  letter-spacing: 1px;
}

.site-slogan {
  font-size: 1.08rem;
  color: #34495e;
  margin: 0.2rem 0 0 2px;
  font-style: italic;
  font-weight: 400;
  line-height: 1.2;
}

.main-content {
  display: grid;
  grid-template-columns: 180px 1fr 180px;
  gap: 2rem;
  max-width: 1600px;
  margin: 0 auto;
  padding: 0 2rem 2rem;
  min-height: calc(100vh - 150px);
}

.sidebar-left, .sidebar-right {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  padding: 1rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  align-self: start;
  position: sticky;
  top: 2rem;
  max-height: calc(100vh - 4rem);
  overflow-y: auto;
}

.content {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 0 1rem;
  perspective: 2000px;
  min-height: 600px;
  /* 添加平滑滚动效果 */
  scroll-behavior: smooth;
  /* 禁用默认滚动 */
  overflow: hidden;
}

.cards-container {
  position: relative;
  width: 100%;
  max-width: 1200px;
  height: 500px;
  margin: 0 auto;
  display: flex;
  justify-content: center;
  align-items: center;
  transform-style: preserve-3d;
  perspective: 2000px;
}

.link-card {
  position: absolute;
  width: 400px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  padding: 0; /* 移除内边距以适应图片 */
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
  transition: all 0.3s cubic-bezier(0.23, 1, 0.32, 1);
  cursor: pointer;
  backdrop-filter: blur(10px);
  transform-style: preserve-3d;
  backface-visibility: hidden;
  will-change: transform, opacity;
  left: 50%;
  margin-left: -200px;
  overflow: hidden; /* 防止图片溢出圆角 */
  z-index: 0;
}

.link-card:nth-child(1),
.link-card:nth-child(2),
.link-card:nth-child(3),
.link-card:nth-child(4),
.link-card:nth-child(5) {
  transform: none;
  opacity: 1;
  filter: none;
}

.link-card::before {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: 12px;
  padding: 2px;
  background: linear-gradient(
    45deg,
    rgba(0, 0, 0, 0.1),
    rgba(255, 255, 255, 0.5),
    rgba(0, 0, 0, 0.1)
  );
  -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  -webkit-mask-composite: xor;
  mask-composite: exclude;
  pointer-events: none;
}

.link-card.slide-enter-active,
.link-card.slide-leave-active {
  transition: all 0.6s cubic-bezier(0.23, 1, 0.32, 1);
}

.link-card.slide-enter-from {
  transform: translateX(100%) translateZ(-200px) rotateY(-25deg);
  opacity: 0;
}

.link-card.slide-leave-to {
  transform: translateX(-100%) translateZ(-200px) rotateY(25deg);
  opacity: 0;
}

@keyframes cardRotate {
  0% {
    transform: rotateY(0deg);
  }
  100% {
    transform: rotateY(360deg);
  }
}

.link-card:hover {
  box-shadow: 0 8px 24px rgba(0,0,0,0.08);
}

.card-content {
  display: flex;
  flex-direction: column;
  height: 100%;
  transform-style: preserve-3d;
}

.card-image {
  width: 100%;
  height: 200px;
  overflow: hidden;
  position: relative;
  transform-style: preserve-3d;
  transform: translateZ(20px);
}

.card-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.card-info {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  flex: 1;
  gap: 1rem;
  background: linear-gradient(
    to bottom,
    rgba(255, 255, 255, 0.95),
    rgba(255, 255, 255, 0.98)
  );
  transform: translateZ(30px);
}

.link-card:hover .card-image img {
  transform: scale(1.1);
}

.link-card h3 {
  margin: 0;
  color: #2c3e50;
  font-size: 1.5rem;
  font-weight: 600;
  line-height: 1.3;
  transform: translateZ(40px);
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
}

.link-description-full {
  color: #666;
  line-height: 1.6;
  margin: 0;
  font-size: 1rem;
  white-space: pre-line;
  word-break: break-all;
  overflow: visible;
  display: block;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
  transform: translateZ(45px);
}

.categories h2, .hot-links h2 {
  color: #2c3e50;
  margin: 0 0 1rem;
  font-size: 1.1rem;
  font-weight: 600;
}

.categories ul, .hot-links ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.categories li a {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0.75rem;
  margin: 0.25rem 0;
  border-radius: 6px;
  color: #2c3e50;
  text-decoration: none;
  transition: all 0.2s;
  cursor: pointer;
  font-size: 0.9rem;
}

.categories li a:hover, .categories li a.active {
  background-color: #e9ecef;
  color: #1a73e8;
}

.category-count {
  font-size: 0.85rem;
  color: #666;
}

.hot-links li {
  padding: 0.5rem 0;
  border-bottom: 1px solid #eee;
  font-size: 0.9rem;
}

.hot-links li:last-child {
  border-bottom: none;
}

.hot-links a {
  color: #2c3e50;
  text-decoration: none;
  cursor: pointer;
  transition: color 0.2s;
}

.hot-links a:hover {
  color: #1a73e8;
}

.visit-count {
  color: #666;
  font-size: 0.85rem;
}

.category-tag {
  background: rgba(233, 236, 239, 0.8);
  padding: 0.3rem 0.8rem;
  border-radius: 6px;
  font-size: 0.85rem;
  color: #2c3e50;
  backdrop-filter: blur(4px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  transform: translateZ(40px);
}

/* 添加鼠标移动时的 3D 效果 */
@media (hover: hover) {
  .link-card {
    transition: transform 0.5s cubic-bezier(0.23, 1, 0.32, 1);
  }
  
  .link-card:hover {
    transition: transform 0.5s cubic-bezier(0.23, 1, 0.32, 1);
  }
  
  .card-content, .link-card h3, .link-description-full, .card-footer, .category-tag {
    transition: transform 0.5s cubic-bezier(0.23, 1, 0.32, 1);
  }
}

/* 添加卡片的光泽效果 */
.link-card::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: 12px;
  background: linear-gradient(
    125deg,
    transparent 0%,
    transparent 40%,
    rgba(255, 255, 255, 0.1) 45%,
    rgba(255, 255, 255, 0.2) 50%,
    rgba(255, 255, 255, 0.1) 55%,
    transparent 60%,
    transparent 100%
  );
  transition: transform 0.5s cubic-bezier(0.23, 1, 0.32, 1);
  pointer-events: none;
}

.link-card:hover::after {
  transform: translateX(100%) rotate(45deg);
}

@media (max-width: 1200px) {
  .main-content {
    grid-template-columns: 180px 1fr;
  }
  .sidebar-right {
    display: none;
  }
}

@media (max-width: 768px) {
  .main-content {
    grid-template-columns: 1fr;
    padding: 0 1rem;
  }
  .sidebar-left {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    top: auto;
    z-index: 100;
    border-radius: 12px 12px 0 0;
    max-height: 60vh;
    transform: translateY(calc(100% - 50px));
    transition: transform 0.3s ease;
  }
  .sidebar-left:hover,
  .sidebar-left:focus-within {
    transform: translateY(0);
  }
  .sidebar-left::before {
    content: '分类';
    display: block;
    text-align: center;
    padding: 0.5rem;
    font-weight: 600;
    color: #2c3e50;
    border-bottom: 1px solid #eee;
  }
  .cards-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    height: auto;
    padding: 1rem 0;
  }
  .link-card {
    position: static;
    width: 100%;
    max-width: 400px;
    margin: 1rem 0;
    transform: none !important;
    opacity: 1 !important;
    filter: none !important;
    left: auto;
    margin-left: 0;
    box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  }
  .nav-button {
    display: none;
  }
  .card-image {
    height: 180px;
  }
  .card-info {
    padding: 1rem;
  }
}

@media (max-width: 480px) {
  .card-image {
    height: 140px;
  }
  .card-info {
    padding: 0.7rem;
  }
  .link-card h3 {
    font-size: 1.1rem;
  }
  .link-description-full {
    font-size: 0.82rem;
    -webkit-line-clamp: 2;
  }
  .category-tag {
    font-size: 0.7rem;
    padding: 0.15rem 0.5rem;
  }
}

.nav-button {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  background: rgba(255, 255, 255, 0.9);
  border: none;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 10;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.nav-button:hover {
  background: rgba(255, 255, 255, 1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-50%) scale(1.1);
}

.nav-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.nav-button.prev {
  left: 20px;
}

.nav-button.next {
  right: 20px;
}

@keyframes shimmer {
  0% {
    background-position: -200% 0;
  }
  100% {
    background-position: 200% 0;
  }
}

.skeleton {
  pointer-events: none;
}

.skeleton-title,
.skeleton-description,
.skeleton-visits,
.skeleton-category {
  background: linear-gradient(
    90deg,
    rgba(255, 255, 255, 0.1) 25%,
    rgba(255, 255, 255, 0.3) 50%,
    rgba(255, 255, 255, 0.1) 75%
  );
  background-size: 200% 100%;
  animation: shimmer 2s infinite;
  border-radius: 4px;
}

.skeleton-title {
  height: 28px;
  width: 80%;
  margin-bottom: 1rem;
}

.skeleton-description {
  height: 80px;
  width: 100%;
  margin-bottom: 1rem;
}

.skeleton-visits {
  height: 20px;
  width: 60px;
}

.skeleton-category {
  height: 20px;
  width: 80px;
}

.empty-card {
  pointer-events: none;
  opacity: 0 !important;
}

/* 添加激活状态样式 */
.link-card.active {
  z-index: 2 !important;
  animation: cardRotate 1s cubic-bezier(0.23, 1, 0.32, 1);
  transform: scale(1.05) translateY(-10px) !important;
  box-shadow: 
    0 15px 35px rgba(0, 0, 0, 0.2),
    0 5px 15px rgba(0, 0, 0, 0.1),
    0 0 0 1px rgba(255, 255, 255, 0.1);
}

/* 添加滚动提示 */
.scroll-hint {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  color: rgba(44, 62, 80, 0.6);
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.scroll-hint.visible {
  opacity: 1;
}

.scroll-hint i {
  font-size: 1.2rem;
  animation: scrollHint 2s infinite;
}

@keyframes scrollHint {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-5px);
  }
}

/* 添加图片加载动画 */
@keyframes imageLoading {
  0% {
    background-position: -200% 0;
  }
  100% {
    background-position: 200% 0;
  }
}

.card-image::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    rgba(255, 255, 255, 0.1) 25%,
    rgba(255, 255, 255, 0.3) 50%,
    rgba(255, 255, 255, 0.1) 75%
  );
  background-size: 200% 100%;
  animation: imageLoading 1.5s infinite;
  z-index: -1;
}

/* 优化卡片悬停效果 */
.link-card:hover {
  transform: scale(1.05) translateY(-10px);
  box-shadow: 
    0 15px 35px rgba(0, 0, 0, 0.2),
    0 5px 15px rgba(0, 0, 0, 0.1);
  z-index: 2;
}

.link-card:hover::after {
  opacity: 1;
}

.link-card::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(
    to bottom,
    transparent,
    rgba(255, 255, 255, 0.1)
  );
  opacity: 0;
  transition: opacity 0.3s ease;
}

.header-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 1rem;
}

.action-btn {
  padding: 8px 16px;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: all 0.3s ease;
}

.add-btn {
  background: #1a73e8;
  color: white;
}

.add-btn:hover {
  background: #1557b0;
}

.logout-btn {
  background: #dc3545;
  color: white;
}

.logout-btn:hover {
  background: #bb2d3b;
}

/* 修改卡片操作按钮样式 */
.card-actions {
  position: absolute;
  top: 8px;
  right: 8px;
  display: flex;
  gap: 8px;
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: 10;
}

.link-card:hover .card-actions {
  opacity: 1;
}

.action-button {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  background: rgba(255, 255, 255, 0.9);
  color: #666;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  backdrop-filter: blur(4px);
}

.action-button:hover {
  transform: scale(1.1);
}

.action-button.delete {
  color: #dc3545;
}

.action-button.delete:hover {
  background: #dc3545;
  color: white;
}

/* 调整对话框样式 */
:deep(.el-dialog) {
  border-radius: 12px;
  overflow: hidden;
}

:deep(.el-dialog__header) {
  margin: 0;
  padding: 1rem 1.5rem;
  background: #f8f9fa;
  border-bottom: 1px solid #eee;
}

:deep(.el-dialog__body) {
  padding: 1.5rem;
}

:deep(.el-dialog__footer) {
  padding: 1rem 1.5rem;
  border-top: 1px solid #eee;
}

:deep(.md-editor) {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  margin-top: 8px;
}

:deep(.md-editor-toolbar) {
  border-bottom: 1px solid #dcdfe6;
}

:deep(.md-editor-preview) {
  padding: 16px;
  background-color: #fff;
}

:deep(.el-switch) {
  margin-bottom: 12px;
}

:deep(.el-dialog) {
  display: flex;
  flex-direction: column;
  max-height: 90vh;
  margin-top: 5vh !important;
}

:deep(.el-dialog__body) {
  flex: 1;
  overflow: auto;
  padding: 20px 30px;
}

:deep(.el-form) {
  height: 100%;
}

:deep(.el-form-item:last-child) {
  margin-bottom: 0;
  height: calc(100% - 240px);
  display: flex;
  flex-direction: column;
}

:deep(.el-form-item:last-child .el-form-item__content) {
  flex: 1;
  margin-left: 0 !important;
}

@media (max-width: 768px) {
  .header-actions {
    margin-top: 0.5rem;
  }

  :deep(.el-dialog) {
    width: 90% !important;
    margin: 0 auto;
  }
}

.github-link {
  position: absolute;
  top: 24px;
  right: 32px;
  color: #333;
  opacity: 0.8;
  transition: opacity 0.2s;
  z-index: 100;
}

.github-link:hover {
  opacity: 1;
  color: #1a73e8;
}

.comments-section {
  margin: 48px auto 0 auto;
  max-width: 800px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
  padding: 32px 24px 24px 24px;
}

.main-footer {
  width: 100%;
  text-align: center;
  padding: 24px 0 12px 0;
  color: #888;
  font-size: 0.95rem;
  background: transparent;
  letter-spacing: 1px;
  margin-top: 48px;
}
</style>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.mb-2 {
  margin-bottom: 8px;
}

:deep(.v-md-editor) {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}

:deep(.el-dialog__body) {
  padding: 20px 30px;
}

:deep(.el-form-item__content) {
  display: flex;
  flex-direction: column;
}
</style>
