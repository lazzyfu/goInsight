// eslint-disable-next-line
import { UserLayout, BasicLayout } from '@/layouts'
import Account from '@/views/account/route'
import Admin from '@/views/admin/route'
import Das from '@/views/das/route'
import Orders from '@/views/orders/route'

export const asyncRouterMap = [
  {
    path: '/',
    name: 'menu.home',
    meta: { title: '首页' },
    component: BasicLayout,
    redirect: { name: 'view.account' },
    children: [Account, ...Admin, ...Das, ...Orders],
  },
  {
    path: '*',
    redirect: { name: 'menu.home' },
  },
]

/**
 * 基础路由
 * @type { *[] }
 */
export const constantRouterMap = [
  {
    path: '/user',
    component: UserLayout,
    redirect: '/user/login',
    hidden: true,
    children: [
      {
        path: 'login',
        name: 'login',
        component: () => import(/* webpackChunkName: "user" */ '@/views/account/Login.vue'),
      },
    ],
  },
  {
    name: '404',
    path: '/404',
    component: () => import('@/views/exception/404.vue'),
  },
  {
    name: '403',
    path: '/403',
    component: () => import('@/views/exception/403.vue'),
  },
  {
    name: '500',
    path: '/500',
    component: () => import('@/views/exception/500.vue'),
  },
]
