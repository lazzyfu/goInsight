// store/permission.js
import router from '@/router'
import { asyncRoutes } from '@/router/ark'
import { defineStore } from 'pinia'

export const usePermissionStore = defineStore('permission', {
  state: () => ({
    routes: [],
    isRoutesGenerated: false,
  }),

  actions: {
    GenerateRoutes(isSuperuser) {
      if (this.isRoutesGenerated) return

      const accessedRoutes = filterAsyncRoutes(asyncRoutes, isSuperuser)

      accessedRoutes.forEach((route) => {
        router.addRoute(route)
      })

      this.routes = accessedRoutes
      this.isRoutesGenerated = true
    },

    reset() {
      this.routes = []
      this.isRoutesGenerated = false
    },
  },
})

function filterAsyncRoutes(routes, isSuperuser) {
  const res = []

  routes.forEach((route) => {
    const tmp = { ...route }

    if (tmp.meta?.requiresAdmin && !isSuperuser) {
      return
    }

    if (tmp.children) {
      tmp.children = filterAsyncRoutes(tmp.children, isSuperuser)
    }

    res.push(tmp)
  })

  return res
}
