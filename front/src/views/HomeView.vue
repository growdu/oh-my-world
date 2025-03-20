<template>
  <el-container>
    <el-main>
      <h1 class="title">💻 Growdu 的计算机世界 🌍</h1>

      <!-- 管理按钮 -->
      <div class="action-buttons">
        <el-button type="primary" class="add-link-btn" @click="dialogVisible = true">
          添加链接
        </el-button>
      </div>

      <!-- 添加链接对话框 -->
      <el-dialog v-model="dialogVisible" title="添加新链接" width="400px">
        <el-form>
          <el-form-item label="图片">
            <el-upload
              class="upload-demo"
              action=""
              :show-file-list="false"
              :before-upload="() => false"
              :on-change="handleUpload"
            >
              <el-button type="primary">选择图片</el-button>
            </el-upload>
            <el-image v-if="newLink.image" :src="newLink.image" class="preview-image"></el-image>
          </el-form-item>
          <el-form-item label="描述">
            <el-input v-model="newLink.description" placeholder="输入描述"></el-input>
          </el-form-item>
          <el-form-item label="链接">
            <el-input v-model="newLink.url" placeholder="输入URL"></el-input>
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="addLink">添加</el-button>
        </template>
      </el-dialog>

      <!-- 卡片列表 -->
      <div class="card-container">
        <el-card v-for="link in allLinks" :key="link.id" class="card" @click="goToLink(link.url)">
          <div class="card-content">
            <el-image :src="link.image" class="card-image"></el-image>
            <p class="card-description">{{ link.description }}</p>
          </div>
        </el-card>
      </div>
    </el-main>
  </el-container>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const links = ref([
  {
    id: 1,
    image: 'https://via.placeholder.com/100',
    description: 'Google',
    url: 'https://www.google.com',
  },
  {
    id: 2,
    image: 'https://via.placeholder.com/100',
    description: 'Baidu',
    url: 'https://www.baidu.com',
  },
  {
    id: 3,
    image: 'https://via.placeholder.com/100',
    description: 'GitHub',
    url: 'https://github.com',
  },
  {
    id: 4,
    image: 'https://via.placeholder.com/100',
    description: 'YouTube',
    url: 'https://www.youtube.com',
  },
  {
    id: 5,
    image: 'https://via.placeholder.com/100',
    description: 'Facebook',
    url: 'https://www.facebook.com',
  },
  {
    id: 6,
    image: 'https://via.placeholder.com/100',
    description: 'Twitter',
    url: 'https://www.twitter.com',
  },
  {
    id: 7,
    image: 'https://via.placeholder.com/100',
    description: 'LinkedIn',
    url: 'https://www.linkedin.com',
  },
  {
    id: 8,
    image: 'https://via.placeholder.com/100',
    description: 'Instagram',
    url: 'https://www.instagram.com',
  },
  {
    id: 9,
    image: 'https://via.placeholder.com/100',
    description: 'Reddit',
    url: 'https://www.reddit.com',
  },
  {
    id: 10,
    image: 'https://via.placeholder.com/100',
    description: 'Netflix',
    url: 'https://www.netflix.com',
  },
])

const allLinks = ref([...links.value])

const newLink = ref({ image: '', description: '', url: '' })
const dialogVisible = ref(false)

// 处理图片上传
const handleUpload = (file) => {
  const reader = new FileReader()
  reader.readAsDataURL(file.raw)
  reader.onload = () => {
    newLink.value.image = reader.result // Base64 存储图片
  }
}

// 添加新链接
const addLink = () => {
  if (!newLink.value.image || !newLink.value.description || !newLink.value.url) {
    ElMessage.error('请填写完整信息！')
    return
  }
  const id = allLinks.value.length + 1
  allLinks.value.unshift({ id, ...newLink.value })
  newLink.value = { image: '', description: '', url: '' }
  dialogVisible.value = false
  ElMessage.success('链接添加成功！')
}

// 点击卡片跳转
const goToLink = (url) => {
  window.open(url, '_blank')
}
</script>

<style scoped>
/******** 标题样式 ********/
.title {
  text-align: center;
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 20px;
}

/******** 操作按钮 ********/
.action-buttons {
  text-align: right;
  margin-bottom: 20px;
}

/******** 管理按钮 ********/
.add-link-btn {
  transition:
    background-color 0.3s ease,
    transform 0.3s ease;
}

.add-link-btn:hover {
  background-color: #409eff;
  transform: scale(1.1);
}

/******** 卡片布局 ********/
.card-container {
  display: grid;
  grid-template-columns: repeat(4, minmax(220px, 1fr)); /* 自动调整列数 */
  gap: 20px;
  width: 100%;
  margin: 0 auto;
}

.card {
  width: 100%;
  padding: 10px;
  text-align: center;
  cursor: pointer;
  transition:
    transform 0.3s ease,
    box-shadow 0.3s ease;
}

.card:hover {
  transform: translateY(-5px);
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
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

/******** 预览图片 ********/
.preview-image {
  width: 100px;
  height: 100px;
  margin-top: 10px;
  border-radius: 10px;
}

/******** 弹窗样式 ********/
.el-dialog {
  transition: transform 0.3s ease-in-out;
}

.el-dialog__wrapper {
  overflow: hidden;
}
</style>
