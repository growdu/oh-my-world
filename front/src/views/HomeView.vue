<template>
  <el-container>
    <el-main>
      <h1 class="title">💻 Growdu 的计算机世界 🌍</h1>

      <!-- 管理按钮 -->
      <!-- <div class="action-buttons">
        <el-button type="primary" class="add-link-btn" @click="dialogVisible = true">
          添加链接
        </el-button>
      </div> -->

      <!-- 添加链接对话框 -->
      <el-dialog v-model="dialogVisible" title="添加新链接" width="400px">
        <el-form>
          <el-form-item label="图片">
            <el-input v-model="newLink.image" placeholder="输入图片 URL"></el-input>
            <el-image v-if="newLink.image" :src="newLink.image" class="preview-image"></el-image>
          </el-form-item>
          <el-form-item label="描述">
            <el-input v-model="newLink.description" placeholder="输入描述"></el-input>
          </el-form-item>
          <el-form-item label="链接">
            <el-input v-model="newLink.url" placeholder="输入网站URL"></el-input>
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
            <el-image v-if="link.image" :src="link.image" class="card-image" fit="cover"></el-image>
            <p class="card-description">{{ link.description }}</p>
          </div>
        </el-card>
      </div>
    </el-main>
  </el-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

const allLinks = ref([])
const newLink = ref({ image: '', description: '', url: '' })
const dialogVisible = ref(false)

// 获取所有链接
const fetchLinks = async () => {
  try {
    const response = await axios.get('http://117.50.177.215:8092/api/v1/links')
    allLinks.value = response.data
  } catch (error) {
    ElMessage.error('无法加载链接')
  }
}

// 添加新链接
const addLink = async () => {
  if (!newLink.value.image || !newLink.value.description || !newLink.value.url) {
    ElMessage.error('请填写完整信息！')
    return
  }
  try {
    await axios.post('http://117.50.177.215:8092/api/v1/links', newLink.value)
    ElMessage.success('链接添加成功！')
    fetchLinks()
    dialogVisible.value = false
    newLink.value = { image: '', description: '', url: '' }
  } catch (error) {
    ElMessage.error('添加链接失败')
  }
}

// 点击卡片跳转
const goToLink = (url) => {
  window.open(url, '_blank')
}

// 初始化时加载链接
onMounted(fetchLinks)
</script>

<style scoped>
.title {
  text-align: center;
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 20px;
}

.action-buttons {
  text-align: right;
  margin-bottom: 20px;
}

.add-link-btn {
  transition:
    background-color 0.3s ease,
    transform 0.3s ease;
}

.add-link-btn:hover {
  background-color: #409eff;
  transform: scale(1.1);
}

.card-container {
  display: grid;
  grid-template-columns: repeat(4, minmax(220px, 1fr));
  gap: 20px;
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

.preview-image {
  width: 100px;
  height: 100px;
  margin-top: 10px;
  border-radius: 10px;
}
</style>
