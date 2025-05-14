<template>
  <el-card class="login-card">
    <el-form :model="loginForm">
      <el-form-item label="用户">
        <el-input v-model="loginForm.username"></el-input>
      </el-form-item>
      <el-form-item label="密码">
        <el-input v-model="loginForm.password" type="password"></el-input>
      </el-form-item>
      <el-button type="primary" @click="login">登录</el-button>
    </el-form>
  </el-card>
</template>

<script setup>
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()
const loginForm = reactive({ username: '', password: '' })

const login = () => {
  if (loginForm.username === 'admin' && loginForm.password === '123456') {
    localStorage.setItem('isAdmin', 'true')
    ElMessage.success('登录成功')
    router.push('/')
  } else {
    ElMessage.error('用户名或密码错误')
  }
}
</script>

<style scoped>
.login-card {
  width: 300px;
  margin: 100px auto;
  padding: 20px;
  text-align: center;
}
</style>
