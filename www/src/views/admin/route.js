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
