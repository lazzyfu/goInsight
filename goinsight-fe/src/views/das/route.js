const route = [
  {
    name: 'view.das',
    path: '/das',
    component: () => import('@/views/das/index.vue'),
    meta: { title: '数据查询', keepAlive: true, icon: 'search' },
  },
]

export default route
