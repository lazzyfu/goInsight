const route = {
  name: 'das',
  path: '/das',
  icon: 'CodeOutlined',
  component: () => import('./index.vue'),
  meta: { title: 'SQL查询', keepAlive: true},
}

export default route
