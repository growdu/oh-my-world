<script setup>
import { ref, onMounted } from 'vue'
import { ElMessageBox, ElMessage, ElButton, ElTable, ElTableColumn } from 'element-plus'
import request from '@/utils/request'

const links = ref([])
const loading = ref(false)

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

onMounted(fetchLinks)
</script>

<template>
  <div>
    <h2>卡片管理</h2>
    <el-table :data="links" v-loading="loading" style="width: 100%">
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="category" label="分类" />
      <el-table-column prop="url" label="链接" />
      <el-table-column label="操作">
        <template #default="{ row }">
          <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>