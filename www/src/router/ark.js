import Layout from '@/components/layout/Layout.vue'
import ACCOUNT from '@/views/account/route.js'
import ADMIN from '@/views/admin/route.js'
import DAS from '@/views/das/route.js'
import ORDERS from '@/views/orders/route.js'

export const arkRouter = [
  {
    path: '/',
    name: 'Root',
    component: Layout,
    redirect: { name: 'account.basic' },
    meta: { title: '首页' },
    children: [
      ADMIN,
      ACCOUNT,
      DAS,
      ORDERS
    ]
  },
  {
    name: 'Login',
    path: '/login',
    component: () => import('@/views/login/index.vue'),
    meta: {
      title: '用户登录',
      module: "用户登录"
    }
  },
  {
    name: '403',
    path: '/403',
    component: () => import('../views/error/403.vue'),
    meta: {
      title: '当前无权限操作, 请联系管理员!'
    }
  },
]
