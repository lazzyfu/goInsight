const route = [
  {
    name: 'view.admin',
    path: '/admin',
    component: () => import('@/views/admin/index.vue'),
    redirect: '/admin/users',
    meta: { title: 'Admin', keepAlive: true, icon: 'setting' },
    children: [
      {
        path: 'user-manage',
        name: `view.admin.userManage`,
        component: () => import('@/views/admin/user-manage'),
        meta: { title: '用户管理' },
      },
      {
        path: 'system-manage',
        name: `view.admin.systemManage`,
        component: () => import('@/views/admin/system-manage/index.vue'),
        meta: { title: '系统配置' },
      },
      {
        path: 'das/detail/:username/:schema/:instance_id',
        name: 'view.admin.systemManage.das.detail',
        component: () => import('@/views/admin/system-manage/das/DasTablesPerm.vue'),
        meta: { title: '权限详情', keepAlive: true },
        hidden: true,
      },
    ],
  },
]

export default route
