const route = {
  name: 'orders',
  path: '/orders',
  icon: 'DatabaseOutlined',
  component: () => import('./index.vue'),
  meta: { title: '数据库工单', keepAlive: true },
  children: [
    {
      name: 'orders.commit',
      path: '/orders/commit',
      icon: 'FormOutlined',
      component: () => import('./commit/index.vue'),
      meta: { title: '提交工单', keepAlive: true },
    },
    {
      name: 'orders.export',
      path: '/orders/export',
      icon: 'UnorderedListOutlined',
      component: () => import('./list/index.vue'),
      meta: { title: '工单列表', keepAlive: true },
    },
    {
      name: 'orders.detail',
      path: '/orders/detail/:order_id',
      component: () => import('./detail/index.vue'),
      meta: { title: '工单详情', keepAlive: true, hidden: true },
    }
  ]
}

export default route
