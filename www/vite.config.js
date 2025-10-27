import { fileURLToPath, URL } from 'node:url';

import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';

// https://vite.dev/config/
export default defineConfig({
  transpileDependencies: true,
  plugins: [
    vue(),
    // vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
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
  chainWebpack: config => {
    config.module
      .rule('svg-icons')
      .test(/\.svg$/)
      .include
      .add(path.resolve(__dirname, 'src/assets/icons'))  // 只处理这个目录下的 SVG 文件
      .end()
      .use('url-loader')
      .loader('url-loader')
      .options({
        limit: 4096,
        name: 'img/[name].[hash:8].[ext]',
        fallback: 'file-loader'
      });

    // 修改已存在的 svg 规则，使其不处理特定目录下的 SVG 文件
    const svgRule = config.module.rule('svg');
    svgRule.exclude.add(path.resolve(__dirname, 'src/assets/icons'));
  },
})
