/**
 * Vue Router 配置
 * 路由守卫逻辑已移至 @/permission.js
 */
import { createRouter, createWebHistory } from 'vue-router'
import { staticRoutes } from './ark'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: staticRoutes,
})

export default router
