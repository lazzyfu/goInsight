import NProgress from 'nprogress';
import { createRouter, createWebHistory } from 'vue-router';
import { arkRouter } from './ark.js';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: arkRouter
})

router.beforeEach((to, from, next) => {
  if (!localStorage.getItem("onLine")) {
    if (to.path !== '/login') {
      return next('/login')
    }
  }
  NProgress.start()
  next()
})

router.afterEach(() => {
  NProgress.done()
})

export default router
