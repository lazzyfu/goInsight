import router from '@/router'
import { asyncRoutes } from '@/router/ark'
import { filterAsyncRoutes } from '@/utils/permission'
import { defineStore } from 'pinia'

export const usePermissionStore = defineStore('permission', {
  state: () => ({
    routes: [],
    isRoutesGenerated: false,
  }),

  getters: {
    /**
     * 获取菜单路由（排除隐藏项）
     */
    menuRoutes: (state) => {
      const root = state.routes.find((r) => r.path === '/')
      return root?.children?.filter(route => !route.meta?.hidden) || []
    },

    /**
     * 检查是否有管理员路由
     */
    hasAdminRoute: (state) => {
      return state.routes.some(route =>
        route.children?.some(child => child.meta?.requiresAdmin)
      )
    },
  },

  actions: {
    /**
     * 生成动态路由
     * @param {Boolean} isSuperuser 是否为超级管理员
     */
    async GenerateRoutes(isSuperuser) {
      if (this.isRoutesGenerated) {
        return
      }

      // 使用工具函数过滤路由
      const accessedRoutes = filterAsyncRoutes(asyncRoutes, isSuperuser)

      // 添加路由到 router
      accessedRoutes.forEach((route) => {
        router.addRoute(route)
      })

      this.routes = accessedRoutes
      this.isRoutesGenerated = true
    },

    /**
     * 重置权限状态
     */
    reset() {
      this.routes = []
      this.isRoutesGenerated = false
    },
  },
})
