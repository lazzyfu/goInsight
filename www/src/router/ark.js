import Layout from '@/components/layout/Layout.vue'
import { markRaw } from 'vue'

import ACCOUNT from '@/views/account/route'
import ADMIN from '@/views/admin/route'
import DAS from '@/views/das/route'
import ORDERS from '@/views/orders/route'

export const staticRoutes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/Login.vue'),
    meta: { title: '用户登录', hidden: true },
  },
  {
    path: '/403',
    name: '403',
    component: () => import('@/views/error/403.vue'),
    meta: { title: '403', hidden: true },
  },
]

export const asyncRoutes = [
  {
    path: '/',
    name: 'Root',
    component: markRaw(Layout),
    redirect: '/das', // 重定向到 SQL查询页面
    children: [
      // 管理员路由 - 标记需要超级管理员权限
      {
        ...ADMIN,
        meta: { ...ADMIN.meta, requiresAdmin: true },
      },
      // 普通用户路由
      ACCOUNT,
      DAS,
      ORDERS,
    ],
  },
]
