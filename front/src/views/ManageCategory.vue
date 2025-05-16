<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const categories = ref([])
const loading = ref(false)
const newCategory = ref('')
const editCategoryDialog = ref(false)
const editCategoryData = ref({ id: null, name: '' })

const fetchCategories = async () => {
  loading.value = true
  try {
    const res = await request.get('/categories')
    categories.value = res.data
  } finally {
    loading.value = false
  }
}
const addCategory = async () => {
  if (!newCategory.value) return
  await request.post('/categories', { name: newCategory.value })
  newCategory.value = ''
  fetchCategories()
}
const deleteCategory = async (cat) => {
  await ElMessageBox.confirm('确定删除该分类？', '提示', { type: 'warning' })
  await request.delete(`/categories/${cat.id}`)
  fetchCategories()
}
const openEditCategory = (cat) => {
  editCategoryData.value = { ...cat }
  editCategoryDialog.value = true
}
const handleEditCategory = async () => {
  await request.put(`/categories/${editCategoryData.value.id}`, { name: editCategoryData.value.name })
  editCategoryDialog.value = false
  fetchCategories()
}
onMounted(fetchCategories)
</script>

<template>
  <div style="max-width:600px;margin:40px auto;">
    <h2>分类管理</h2>
    <el-table :data="categories" v-loading="loading">
      <el-table-column prop="name" label="分类名" />
      <el-table-column label="操作">
        <template #default="{ row }">
          <el-button size="small" @click="openEditCategory(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="deleteCategory(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <div style="margin:16px 0;">
      <el-input v-model="newCategory" placeholder="新分类名" style="width:200px;margin-right:8px;" />
      <el-button @click="addCategory" type="primary">添加</el-button>
    </div>
    <el-dialog v-model="editCategoryDialog" title="编辑分类">
      <el-input v-model="editCategoryData.name" />
      <template #footer>
        <el-button @click="editCategoryDialog = false">取消</el-button>
        <el-button type="primary" @click="handleEditCategory">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.el-table {
  margin-bottom: 16px;
}
</style> 