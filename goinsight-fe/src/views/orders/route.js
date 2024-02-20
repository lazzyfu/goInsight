const route = [
  {
    name: 'view.orders.list',
    path: '/orders/list',
    component: () => import('@/views/orders/list/index.vue'),
    meta: { title: '工单列表', keepAlive: true, icon: 'unordered-list' },
  },
  {
    name: 'view.orders.detail',
    path: '/orders/detail/:order_id',
    component: () => import('@/views/orders/list/detail.vue'),
    meta: { title: '工单详情', keepAlive: true },
    hidden: true,
  },
  {
    name: 'view.orders.tasks',
    path: '/orders/tasks/:order_id',
    component: () => import('@/views/orders/tasks/index.vue'),
    meta: { title: '工单任务', keepAlive: true },
    hidden: true,
  },
  {
    name: 'view.orders.commit',
    path: '/orders/commit',
    component: () => import('@/views/orders/index.vue'),
    redirect: '/orders/commit/ddl',
    meta: { title: '提交工单', keepAlive: true, icon: 'edit' },
    children: [
      {
        path: 'ddl',
        name: `view.orders.commit.ddl`,
        component: () => import('@/views/orders/commit/index.vue'),
        meta: { title: 'DDL工单' },
      },
      {
        path: 'dml',
        name: `view.orders.commit.dml`,
        component: () => import('@/views/orders/commit/index.vue'),
        meta: { title: 'DML工单' },
      },
      {
        path: 'export',
        name: `view.orders.commit.export`,
        component: () => import('@/views/orders/commit/index.vue'),
        meta: { title: '导出工单' },
      },
    ],
  },
]

export default route
