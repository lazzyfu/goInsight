/**
 * 权限工具函数
 * 参考 ant-design-vue-pro
 */

/**
 * 检查路由权限
 * @param {Array} permissions 用户权限列表
 * @param {Object} route 路由对象
 * @returns {Boolean}
 */
export function hasPermission(permissions, route) {
  if (route.meta && route.meta.permission) {
    if (!permissions || permissions.length === 0) {
      return false
    }

    // 检查是否有任意一个权限匹配
    return permissions.some(permission =>
      route.meta.permission.includes(permission)
    )
  }

  // 没有权限要求的路由默认可访问
  return true
}

/**
 * 检查是否为超级管理员
 * @param {Object} user 用户信息
 * @returns {Boolean}
 */
export function isSuperAdmin(user) {
  return user && user.is_superuser === true
}

/**
 * 检查操作权限
 * @param {String} action 操作名称
 * @param {Array} permissions 权限列表
 * @returns {Boolean}
 */
export function hasAction(action, permissions = []) {
  return permissions.includes(action)
}

/**
 * 递归过滤异步路由
 * @param {Array} routes 路由配置数组
 * @param {Boolean} isSuperuser 是否为超级管理员
 * @returns {Array} 过滤后的路由数组
 */
export function filterAsyncRoutes(routes, isSuperuser) {
  return routes.reduce((acc, route) => {
    const tmp = { ...route }

    // 超级管理员跳过权限检查
    if (!isSuperuser && tmp.meta?.requiresAdmin) {
      return acc
    }

    // 递归过滤子路由
    if (tmp.children?.length) {
      tmp.children = filterAsyncRoutes(tmp.children, isSuperuser)

      // 如果所有子路由都被过滤且父路由需要管理员权限，则跳过
      if (tmp.children.length === 0 && tmp.meta?.requiresAdmin) {
        return acc
      }
    }

    acc.push(tmp)
    return acc
  }, [])
}
