const route = {
  name: 'view.sqlorders',
  path: 'sqlorders',
  redirect: { name: 'view.sqlorders.list' },
  component: () => import('./index.vue'),
  meta: { title: 'DB工单', keepAlive: true, icon: 'database' },
  children: [
    {
      name: 'view.sqlorders.commit.version',
      path: '/sqlorders/version',
      component: () => import('./version.vue'),
      meta: { title: '上线版本', keepAlive: true, icon: 'fork' }
    },
    {
      name: 'view.sqlorders.list',
      path: '/sqlorders/list',
      component: () => import('./list.vue'),
      meta: { title: '工单列表', keepAlive: true, icon: 'history' }
    },
    {
      name: 'view.sqlorders.detail',
      path: '/sqlorders/detail/:order_id',
      hidden: true,
      component: () => import('./detail.vue'),
      meta: { title: '工单详情', keepAlive: true }
    },
    {
      name: 'view.sqlorders.commit.ddl',
      path: '/sqlorders/commit/ddl',
      component: () => import('./commit.vue'),
      meta: { title: '提交DDL工单', keepAlive: true, icon: 'retweet' }
    },
    {
      name: 'view.sqlorders.commit.dml',
      path: '/sqlorders/commit/dml',
      component: () => import('./commit.vue'),
      meta: { title: '提交DML工单', keepAlive: true, icon: 'swap' }
    },
    {
      name: 'view.sqlorders.commit.export',
      path: '/sqlorders/commit/export',
      component: () => import('./commit.vue'),
      meta: { title: '提交导出工单', keepAlive: true, icon: 'export' }
    },
    {
      name: 'view.sqlorders.tasks',
      path: '/sqlorders/tasks/:task_id',
      hidden: true,
      component: () => import('./task.vue'),
      meta: { title: '工单任务', hidden: true }
    },
    {
      name: 'view.sqlorders.export.download',
      path: '/sqlorders/export/download/:base64_filename',
      hidden: true,
      component: () => import('./download.vue'),
      meta: { title: '下载导出文件', hidden: true }
    },
    {
      name: 'view.sqlorders.version.view',
      path: '/sqlorders/version/view/:version',
      hidden: true,
      component: () => import('./versionView.vue'),
      meta: { title: '版本详情', hidden: true }
    }
  ]
}

export default route
