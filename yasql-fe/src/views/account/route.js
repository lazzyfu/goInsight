const route = {
  path: '/account',
  name: 'account',
  component: () => import('@/views/account/settings/Index'),
  meta: { title: '个人设置', hideHeader: true, icon: 'user' },
  redirect: '/account/settings/base',
  hidden: true,
  children: [
    {
      path: '/account/settings/base',
      name: 'BaseSettings',
      component: () => import('@/views/account/settings/BaseSetting'),
      meta: { title: '基本设置', hidden: true }
    },
    {
      path: '/account/settings/security',
      name: 'SecuritySettings',
      component: () => import('@/views/account/settings/Security'),
      meta: { title: '安全设置', hidden: true, keepAlive: true }
    },
    // {
    //   path: '/account/settings/custom',
    //   name: 'CustomSettings',
    //   component: () => import('@/views/account/settings/Custom'),
    //   meta: { title: '个性化设置', hidden: true, keepAlive: true }
    // },
    {
      path: '/account/settings/binding',
      name: 'BindingSettings',
      component: () => import('@/views/account/settings/Binding'),
      meta: { title: '账户绑定', hidden: true, keepAlive: true }
    }
    // {
    //   path: '/account/settings/notification',
    //   name: 'NotificationSettings',
    //   component: () => import('@/views/account/settings/Notification'),
    //   meta: { title: '新消息通知', hidden: true, keepAlive: true }
    // }
  ]
}

export default route
