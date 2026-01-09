const route = {
  name: 'orders',
  path: '/orders',
  redirect: '/orders/list',
  component: () => import('./index.vue'),
  meta: { title: '数据库工单', icon: 'DatabaseOutlined', keepAlive: true },
  children: [
    {
      name: 'orders.create',
      path: '/orders/create',
      component: () => import('./create/OrderCreate.vue'),
      meta: { title: '新建工单', icon: 'FormOutlined', keepAlive: true },
    },
    {
      name: 'orders.list',
      path: '/orders',
      component: () => import('./list/OrderList.vue'),
      meta: { title: '工单列表', icon: 'UnorderedListOutlined', keepAlive: true },
    },
    {
      name: 'orders.detail',
      path: '/orders/:order_id',
      component: () => import('./detail/OrderDetail.vue'),
      meta: { title: '工单详情', keepAlive: true, hidden: true },
    },
    {
      name: 'orders.tasks',
      path: '/orders/tasks/:order_id',
      component: () => import('./tasks/TaskList.vue'),
      meta: { title: '工单任务', keepAlive: true, hidden: true },
    },
    {
      name: 'orders.exportfile.download',
      path: '/orders/tasks/exports/:filename',
      component: () => import('./tasks/ExportfileDownload.vue'),
      meta: { title: '下载导出文件', keepAlive: true, hidden: true },
    },
  ],
}

export default route
