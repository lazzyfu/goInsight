/**
 * Action 权限指令 (Vue3 版本)
 * 参考 ant-design-vue-pro，适配 Vue3 + Pinia
 *
 * 使用方法:
 * <a-button v-action:add>添加</a-button>
 * <a-button v-action:edit>编辑</a-button>
 * <a-button v-action:delete>删除</a-button>
 *
 * 当用户没有对应权限时，按钮会被隐藏
 */

const action = {
  mounted(el, binding, vnode) {
    const actionName = binding.arg

    // 获取 store（Vue3 方式）
    const { $pinia } = vnode.appContext.app.config.globalProperties
    if (!$pinia) return

    // 动态导入 store
    import('@/store/user').then(({ useUserStore }) => {
      const userStore = useUserStore($pinia)

      // 超级管理员拥有所有权限
      if (userStore.is_superuser) {
        return
      }

      // 从路由元信息中获取权限配置
      const permissions = binding.instance.$route.meta?.permissions || []

      // 检查是否有对应的操作权限
      if (actionName && !permissions.includes(actionName)) {
        // 没有权限则移除元素
        el.parentNode && el.parentNode.removeChild(el)
      }
    })
  }
}

export default action
