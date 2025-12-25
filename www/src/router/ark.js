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
  {
    // 动态路由尚未生成时（首屏刷新/直达深链接）兜底，避免 Vue Router 打印 "No match found" 警告
    // 动态路由生成后会 replace 回真实目标路由
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/403.vue'),
    meta: { title: 'Not Found', hidden: true },
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
