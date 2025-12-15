// router/index.js
import { usePermissionStore } from '@/store/permission'
import { useUserStore } from '@/store/user'
import NProgress from 'nprogress'
import { createRouter, createWebHistory } from 'vue-router'
import { staticRoutes } from './ark'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: staticRoutes,
})

const whiteList = ['/login', '/403']

router.beforeEach(async (to, from, next) => {
  NProgress.start()

  const userStore = useUserStore()
  const permissionStore = usePermissionStore()

  // 未登录
  if (!userStore.token) {
    if (whiteList.includes(to.path)) {
      next()
    } else {
      next('/login')
    }
    NProgress.done()
    return
  }

  // 已登录禁止访问 login
  if (to.path === '/login') {
    next('/')
    NProgress.done()
    return
  }

  // 动态路由只生成一次
  if (!permissionStore.isRoutesGenerated) {
    try {
      if (userStore.is_superuser === undefined) {
        await userStore.getInfo()
      }

      permissionStore.GenerateRoutes(userStore.is_superuser)

      next({ ...to, replace: true })
      return
    } catch (e) {
      userStore.clear()
      permissionStore.reset()
      next('/login')
      return
    }
  }

  // admin 兜底
  if (to.meta.requiresAdmin && !userStore.is_superuser) {
    next('/403')
    return
  }

  next()
})

router.afterEach(() => {
  NProgress.done()
})

export default router
