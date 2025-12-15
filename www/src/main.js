import { setupDirectives } from '@/directives'; // 全局指令
import Antd from 'ant-design-vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import { createApp } from 'vue'
import { VueCookieNext } from 'vue-cookie-next'
import App from './App.vue'
import './permission'; // 路由权限控制
import router from './router'

// 引入antd样式
import "ant-design-vue/dist/reset.css"
import "epic-designer/dist/style.css"
import { setupAntd } from "epic-designer/dist/ui/antd"
setupAntd()

// 引入vue-cropper
import VueCodemirror from 'vue-codemirror'
import VueCropper from 'vue-cropper'
import 'vue-cropper/dist/index.css'

// 引入highlight.js
import hljsVuePlugin from '@highlightjs/vue-plugin'; // 支持vue3的组件
import 'highlight.js/lib/common'; // 依赖包
import 'highlight.js/styles/atom-one-dark.css'; // 样式

const app = createApp(App)

// pinia
const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);
app.use(pinia)

// 注册全局指令
setupDirectives(app)

// 过滤 Vue Router 的 "No match found" 警告（动态路由生成前的正常警告）
const originalWarn = console.warn
console.warn = function (...args) {
  if (typeof args[0] === 'string' && args[0].includes('No match found for location')) {
    return
  }
  originalWarn.apply(console, args)
}

app
  .use(router)
  .use(Antd)
  .use(VueCookieNext)
  .use(VueCropper)
  .use(VueCodemirror)
  .use(hljsVuePlugin)
  .mount('#app')
