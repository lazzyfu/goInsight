const route = {
  path: '/account',
  name: 'view.account',
  component: () => import('@/views/account/settings/Index.vue'),
  meta: { title: '个人中心',  icon: 'user' },
  redirect: '/account/settings/basic',
  hideChildrenInMenu: true,
  // hidden: true,
  children: [
    {
      path: '/account/settings/basic',
      name: 'BasicSettings',
      component: () => import('@/views/account/settings/BasicSetting.vue'),
      meta: { title: '基本设置', keepAlive: true  },
      hidden: true,
    },
    {
      path: '/account/settings/security',
      name: 'SecuritySettings',
      component: () => import('@/views/account/settings/Security.vue'),
      meta: { title: '安全设置',  keepAlive: true },
      hidden: true,
    },
  ],
}
export default route
