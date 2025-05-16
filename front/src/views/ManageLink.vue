<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import request from '@/utils/request'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

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
const categories = ref([])
const dialogWidth = computed(() => useMarkdown.value ? '80%' : '500px')
const editorHeight = computed(() => useMarkdown.value ? '400px' : '120px')
const editDialogVisible = ref(false)
const editLink = ref({})
const isEditMode = ref(false)
const handleImageError = (e) => {
  e.target.src = 'https://via.placeholder.com/120x80?text=No+Image'
}
const fetchCategories = async () => {
  loading.value = true
  try {
    const res = await request.get('/categories')
    categories.value = res.data || []
    const linkRes = await request.get('/links')
    links.value = linkRes.data
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
    fetchCategories()
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
    fetchCategories()
  } catch (error) {
    ElMessage.error('添加失败：' + error.message)
  }
}
const openEditDialog = (link) => {
  isEditMode.value = true
  editLink.value = { ...link }
  useMarkdown.value = false
  fetchCategories()
  editDialogVisible.value = true
}
const openAddDialog = () => {
  isEditMode.value = false
  newLink.value = { name: '', image: '', description: '', url: '', category: '' }
  useMarkdown.value = false
  fetchCategories()
  dialogVisible.value = true
}
const handleEditLink = async () => {
  if (!editLink.value.name || !editLink.value.url || !editLink.value.category) {
    ElMessage.error('请填写必要的信息！')
    return
  }
  try {
    await request.put(`/links/${editLink.value.id}`, editLink.value)
    ElMessage.success('修改成功！')
    editDialogVisible.value = false
    fetchCategories()
  } catch (error) {
    ElMessage.error('修改失败：' + error.message)
  }
}
onMounted(() => {
  if (localStorage.getItem('isAdmin') !== 'true') {
    router.push('/')
  }
  fetchCategories()
})
</script>

<template>
  <div class="manage-container">
    <div class="manage-header">
      <el-button type="primary" @click="router.push('/')">返回首页</el-button>
      <el-button type="success" @click="openAddDialog">添加新链接</el-button>
      <el-button @click="router.push('/manage-category')">分类管理</el-button>
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
      <el-table-column label="操作" width="160">
        <template #default="{ row }">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加链接对话框 -->
    <el-dialog v-model="dialogVisible" :width="dialogWidth" title="添加新链接">
      <el-form :model="newLink" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="newLink.name" />
        </el-form-item>
        <el-form-item label="图片链接">
          <el-input v-model="newLink.image" placeholder="请输入图片URL" />
        </el-form-item>
        <el-form-item label="图片预览">
          <img v-if="newLink.image" :src="newLink.image" alt="预览" style="max-width:120px;max-height:80px;border-radius:6px;border:1px solid #eee;" @error="handleImageError" />
          <span v-else style="color:#aaa;">暂无图片</span>
        </el-form-item>
        <el-form-item label="URL">
          <el-input v-model="newLink.url" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="newLink.category" placeholder="请选择分类" style="width:100%">
            <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-switch
            v-model="useMarkdown"
            active-text="Markdown编辑器"
            inactive-text="普通文本"
            class="mb-2"
          />
          <el-input
            v-if="!useMarkdown"
            v-model="newLink.description"
            type="textarea"
            :rows="4"
            :style="{height: editorHeight}"
          />
          <md-editor
            v-else
            v-model="newLink.description"
            :style="{ height: editorHeight }"
            :toolbars="['bold','italic','title','link','image','code','preview']"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleAddLink">确认</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 编辑链接对话框 -->
    <el-dialog v-model="editDialogVisible" :width="dialogWidth" title="编辑链接">
      <el-form :model="editLink" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="editLink.name" />
        </el-form-item>
        <el-form-item label="图片链接">
          <el-input v-model="editLink.image" placeholder="请输入图片URL" />
        </el-form-item>
        <el-form-item label="图片预览">
          <img v-if="editLink.image" :src="editLink.image" alt="预览" style="max-width:120px;max-height:80px;border-radius:6px;border:1px solid #eee;" @error="handleImageError" />
          <span v-else style="color:#aaa;">暂无图片</span>
        </el-form-item>
        <el-form-item label="URL">
          <el-input v-model="editLink.url" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="editLink.category" placeholder="请选择分类" style="width:100%">
            <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-switch
            v-model="useMarkdown"
            active-text="Markdown编辑器"
            inactive-text="普通文本"
            class="mb-2"
          />
          <el-input
            v-if="!useMarkdown"
            v-model="editLink.description"
            type="textarea"
            :rows="4"
            :style="{height: editorHeight}"
          />
          <md-editor
            v-else
            v-model="editLink.description"
            :style="{ height: editorHeight }"
            :toolbars="['bold','italic','title','link','image','code','preview']"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleEditLink">保存</el-button>
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