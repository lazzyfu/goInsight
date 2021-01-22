const route = {
  name: 'redisms',
  path: 'redisms',
  component: () => import('./index.vue'),
  redirect: { name: 'redisms.list' },
  meta: { title: 'Redis管理', icon: 'sync' },
  children: [
    {
      name: 'redisms.list',
      path: 'list',
      component: () => import('./List/index.vue'),
      meta: { title: 'Redis查询', icon: 'search' },
    }
  ],
}

export default route