import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueJsx(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  // server: {
  //   host: '0.0.0.0', // 允许外部访问
  //   port: 5173,       // 或你的端口
  //   allowedHosts: ['growdu.cn'], // ✅ 添加你要允许的域名
  //   // 如果你希望完全开放，也可以写成：
  //   // allowedHosts: 'all'
  // }
})
