const route = {
  name: 'view.admin',
  path: '/admin',
  icon: 'SettingOutlined',
  component: () => import('./index.vue'),
  redirect: '/admin/users',
  meta: { title: '后台管理', keepAlive: true },
  children: [
    {
      name: `view.admin.perms`,
      path: '/admin/perms',
      icon: 'UserOutlined',
      component: () => import('./perms/index.vue'),
      meta: { title: '权限管理', keepAlive: true },
      children: [
        {
          name: `view.admin.users`,
          path: '/admin/users',
          icon: 'UserOutlined',
          component: () => import('./perms/users/index.vue'),
          meta: { title: '用户管理', keepAlive: true },
        },
        {
          name: `view.admin.roles`,
          path: '/admin/roles',
          icon: 'IdcardOutlined',
          component: () => import('./perms/roles/index.vue'),
          meta: { title: '角色管理', keepAlive: true },
        },
        {
          name: `view.admin.orgs`,
          path: '/admin/orgs',
          icon: 'ApartmentOutlined',
          component: () => import('./perms/orgs/index.vue'),
          meta: { title: '组织管理', keepAlive: true },
        },
      ]
    },
    {
      name: `view.admin.system`,
      path: '/admin/system',
      icon: 'ToolOutlined',
      component: () => import('./system/index.vue'),
      meta: { title: '系统配置', keepAlive: true },
      children: [
        {
          name: `view.admin.environments`,
          path: '/admin/environements',
          icon: 'ClusterOutlined',
          component: () => import('./system/environments/index.vue'),
          meta: { title: '环境管理', keepAlive: true },
        },
        {
          name: `view.admin.dbconfig`,
          path: '/admin/dbconfig',
          icon: 'DatabaseOutlined',
          component: () => import('./system/dbconfig/index.vue'),
          meta: { title: '数据库管理', keepAlive: true },
        },
        {
          name: `view.admin.inspect`,
          path: '/admin/inspect',
          icon: 'CheckCircleOutlined',
          component: () => import('./system/inspect/index.vue'),
          meta: { title: '审核参数', keepAlive: true },
        },
        {
          name: `view.admin.das`,
          path: '/admin/das',
          icon: 'CheckCircleOutlined',
          component: () => import('./system/das/index.vue'),
          meta: { title: '数据访问', keepAlive: true },
        }
      ]
    }
  ]
}
export default route
