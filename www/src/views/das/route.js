const route = {
  name: 'das',
  path: '/das',
  component: () => import('./index.vue'),
  meta: { title: 'SQL查询', icon: 'CodeOutlined', keepAlive: true},
}

export default route
