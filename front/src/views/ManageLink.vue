<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()
const links = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const useMarkdown = ref(false)
const newLink = ref({
  name: '',
  image: '',
  description: '',
  url: '',
  category: ''
})

const fetchLinks = async () => {
  loading.value = true
  try {
    const res = await request.get('/links')
    links.value = res.data
  } finally {
    loading.value = false
  }
}

const handleDelete = async (link) => {
  try {
    await ElMessageBox.confirm('确定要删除这个链接吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.delete(`/links/${link.id}`)
    ElMessage.success('删除成功')
    fetchLinks()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败：' + error.message)
    }
  }
}

const handleAddLink = async () => {
  if (!newLink.value.name || !newLink.value.url || !newLink.value.category) {
    ElMessage.error('请填写必要的信息！')
    return
  }
  try {
    await request.post('/links', newLink.value)
    ElMessage.success('添加成功！')
    dialogVisible.value = false
    newLink.value = { name: '', image: '', description: '', url: '', category: '' }
    fetchLinks()
  } catch (error) {
    ElMessage.error('添加失败：' + error.message)
  }
}

onMounted(() => {
  if (localStorage.getItem('isAdmin') !== 'true') {
    router.push('/')
  }
  fetchLinks()
})
</script>

<template>
  <div class="manage-container">
    <div class="manage-header">
      <el-button type="primary" @click="router.push('/')">返回首页</el-button>
      <el-button type="success" @click="dialogVisible = true">添加新链接</el-button>
    </div>
    <el-table :data="links" v-loading="loading" class="manage-table" border>
      <el-table-column prop="image" label="图片" width="100">
        <template #default="{ row }">
          <img :src="row.image" alt="图片" style="width:60px;height:40px;object-fit:cover;border-radius:6px;" />
        </template>
      </el-table-column>
      <el-table-column prop="name" label="名称" width="120" />
      <el-table-column prop="category" label="分类" width="100" />
      <el-table-column prop="url" label="链接" width="200">
        <template #default="{ row }">
          <a :href="row.url" target="_blank">{{ row.url }}</a>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" min-width="200" />
      <el-table-column label="操作" width="100">
        <template #default="{ row }">
          <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加链接对话框 -->
    <el-dialog v-model="dialogVisible" title="添加新链接" width="500px">
      <el-form :model="newLink" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="newLink.name" />
        </el-form-item>
        <el-form-item label="图片链接">
          <el-input v-model="newLink.image" />
        </el-form-item>
        <el-form-item label="URL">
          <el-input v-model="newLink.url" />
        </el-form-item>
        <el-form-item label="分类">
          <el-input v-model="newLink.category" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="newLink.description" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleAddLink">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.manage-container {
  max-width: 1100px;
  margin: 40px auto;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.08);
  padding: 32px 24px;
}
.manage-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}
.manage-table {
  margin-top: 12px;
  border-radius: 8px;
  overflow: hidden;
}
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>