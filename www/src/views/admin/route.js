const route = {
  name: 'view.admin',
  path: '/admin',
  icon: 'SettingOutlined',
  component: () => import('./index.vue'),
  redirect: '/admin/users',
  meta: { title: '后台管理', keepAlive: true },
  children: [
    {
      name: `view.admin.users`,
      path: '/admin/users',
      icon: 'UserOutlined',
      component: () => import('./users/index.vue'),
      meta: { title: '用户管理', keepAlive: true },
    },
    {
      name: `view.admin.roles`,
      path: '/admin/roles',
      icon: 'UserOutlined',
      component: () => import('./roles/index.vue'),
      meta: { title: '角色管理', keepAlive: true },
    },
    {
      name: `view.admin.orgs`,
      path: '/admin/orgs',
      icon: 'UserOutlined',
      component: () => import('./orgs/index.vue'),
      meta: { title: '组织管理', keepAlive: true },
    },
    {
      name: `view.admin.environments`,
      path: '/admin/environements',
      icon: 'UserOutlined',
      component: () => import('./environments/index.vue'),
      meta: { title: '环境管理', keepAlive: true },
    },
    {
      name: `view.admin.dbconfig`,
      path: '/admin/dbconfig',
      icon: 'UserOutlined',
      component: () => import('./dbconfig/index.vue'),
      meta: { title: '数据库管理', keepAlive: true },
    },
    {
      name: `view.admin.inspect`,
      path: '/admin/inspect',
      icon: 'UserOutlined',
      component: () => import('./inspect/index.vue'),
      meta: { title: '审核参数', keepAlive: true },
    },
    // {
    //   path: 'system',
    //   name: `view.admin.systemManage`,
    //   component: () => import('@/views/admin/system/index'),
    //   meta: { title: '系统配置' },
    // },
    // {
    //   path: 'das/detail/:username/:schema/:instance_id',
    //   name: 'view.admin.systemManage.das.detail',
    //   component: () => import('@/views/admin/system/das/DasTablesPerm'),
    //   meta: { title: '权限详情', keepAlive: true },
    //   hidden: true,
    // },
  ]
}
export default route
