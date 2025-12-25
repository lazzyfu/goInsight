const route = {
  name: 'account',
  path: '/account',
  component: () => import('./index.vue'),
  meta: { title: '个人中心', icon: 'UserOutlined', keepAlive: true, hidden: true },
  children: [
    {
      name: 'account.basic',
      path: '/account/basic',
      component: () => import('./settings/UserBasic.vue'),
      meta: { title: '基本设置', icon: 'SettingOutlined', keepAlive: true, hidden: true },
    },
  ],
}

export default route
