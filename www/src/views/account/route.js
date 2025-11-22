const route = {
  name: 'account',
  path: '/account',
  icon: 'UserOutlined',
  component: () => import('./index.vue'),
  meta: { title: '个人中心', keepAlive: true },
  children: [
    {
      name: 'account.basic',
      path: '/account/basic',
      icon: 'SettingOutlined',
      component: () => import('./settings/UserBasic.vue'),
      meta: { title: '基本设置', keepAlive: true },
    },
    {
      name: 'account.security',
      path: '/account/security',
      icon: 'SafetyOutlined',
      component: () => import('./settings/UserSecurity.vue'),
      meta: { title: '安全设置', keepAlive: true },
    },
  ],
}

export default route
