/**
 * 路由权限控制
 * 参考 ant-design-vue-pro 的设计模式，适配 Vue3 + Pinia
 */
import { usePermissionStore } from '@/store/permission'
import { useUserStore } from '@/store/user'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import router from './router'

NProgress.configure({ showSpinner: false })

// 免登录白名单
const WHITE_LIST = ['/login', '/403']
// 默认路由
const DEFAULT_ROUTE = '/das'

/**
 * 路由前置守卫
 * 处理登录验证、动态路由生成和权限控制
 */
router.beforeEach(async (to, from, next) => {
  NProgress.start()

  const userStore = useUserStore()
  const permissionStore = usePermissionStore()

  // 白名单路由直接放行（登录页、403页面等）
  if (WHITE_LIST.includes(to.path)) {
    next()
    NProgress.done()
    return
  }

  // 未登录用户重定向到登录页
  if (!userStore.isLoggedIn) {
    next({ path: '/login', query: { redirect: to.fullPath } })
    NProgress.done()
    return
  }

  // 已登录用户访问登录页，重定向到默认页面
  if (to.path === '/login') {
    next({ path: DEFAULT_ROUTE })
    NProgress.done()
    return
  }

  // 首次访问时生成动态路由
  if (!permissionStore.isRoutesGenerated) {
    try {
      // 确保用户信息已加载
      if (!userStore.isInfoLoaded) {
        await userStore.getInfo()
      }

      // 根据用户权限生成可访问的路由
      await permissionStore.GenerateRoutes(userStore.is_superuser)

      // 如果是刷新/深链接访问，在动态路由加入之前可能先命中了兜底路由（NotFound）。
      // 此时需要 replace 回原始地址，才能匹配到新加入的真实路由（例如 /orders/:order_id）。
      if (to.name === 'NotFound') {
        next({ path: to.fullPath, replace: true })
        return
      }

      // 重新导航到目标路由，确保动态添加的路由生效
      // 使用 replace: true 避免在浏览器历史记录中留下重复条目
      next({ ...to, replace: true })
      return
    } catch (error) {
      console.error('[Permission] 路由生成失败:', error)
      // 清除用户状态并重定向到登录页
      userStore.clear()
      permissionStore.reset()
      next({ path: '/login', query: { redirect: to.fullPath } })
      NProgress.done()
      return
    }
  }

  // 验证管理员权限（如果路由需要管理员权限）
  if (to.meta?.requiresAdmin && !userStore.is_superuser) {
    next({ path: '/403' })
    NProgress.done()
    return
  }

  // 通过所有验证，允许访问
  next()
})

/**
 * 路由后置守卫
 */
router.afterEach(() => {
  NProgress.done()
})
