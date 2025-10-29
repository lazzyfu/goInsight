const route = {
  name: 'orders',
  path: '/orders',
  icon: 'DatabaseOutlined',
  component: () => import('./index.vue'),
  meta: { title: '数据库工单', keepAlive: true },
  children: [
    {
      name: 'orders.create',
      path: '/orders/create',
      icon: 'FormOutlined',
      component: () => import('./create/index.vue'),
      meta: { title: '新建工单', keepAlive: true },
    },
    {
      name: 'orders.list',
      path: '/orders',
      icon: 'UnorderedListOutlined',
      component: () => import('./list/index.vue'),
      meta: { title: '工单列表', keepAlive: true },
    },
    {
      name: 'orders.detail',
      path: '/orders/:order_id',
      component: () => import('./detail/index.vue'),
      meta: { title: '工单详情', keepAlive: true, hidden: true },
    }
  ]
}

export default route
