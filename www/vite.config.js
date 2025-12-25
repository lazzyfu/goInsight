import { fileURLToPath, URL } from 'node:url'

import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

// 保留 fileURLToPath 用于 alias '@' 指向 src

// https://vite.dev/config/
export default defineConfig({
  transpileDependencies: true,
  plugins: [
    vue(),
    // vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  // 本地开发代理
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8083',
        changeOrigin: true,
      },
      '/media': {
        target: 'http://localhost:8083',
        changeOrigin: true,
      },
      '/ws': {
        target: 'ws://127.0.0.1:8083',
        ws: true,
        secure: false,
      },
    }
  },
  // NOTE:
  // Vite 的配置不支持 chainWebpack（这是 vue-cli/webpack-chain 的 API），放在这里不会生效。
  // 之前的 svg-icons 规则也依赖 webpack loader（url-loader/file-loader），在 Vite 下同样不会起作用。
  // 如需 svg sprite/自动注册图标，请引入相应的 Vite 插件（例如 vite-plugin-svg-icons 等）。
})
