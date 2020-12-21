import axios from 'axios'
import store from '@/store'
import storage from 'store'
import { redirect } from '@/router/helpers.js'
import notification from 'ant-design-vue/es/notification'
import { VueAxios } from './axios'
import { ACCESS_TOKEN } from '@/store/mutation-types'

// 创建 axios 实例
const request = axios.create({
  // API 请求的默认前缀
  baseURL: process.env.VUE_APP_API_BASE_URL,
  timeout: 650000 // 请求超时时间
})

// 异常拦截处理器
const errorHandler = error => {
  if (error.response) {
    const data = error.response.data
    // 从 localstorage 获取 token
    const token = storage.get(ACCESS_TOKEN)
    if (error.response.status === 403) {
      notification.error({
        message: '403',
        description: error.response.data.detail
      })
      // redirect({ name: '403' })
    }
    if (error.response.status === 404) {
      redirect({ name: '404' })
    }
    // if (error.response.status === 500) {
    //   redirect({ name: '500' })
    // }
    if (error.response.status === 401) {
      notification.error({
        message: 'Unauthorized',
        description: '认证失败，请重新登录'
      })
      if (token) {
        store.dispatch('Logout').then(() => {
          setTimeout(() => {
            window.location.reload()
          }, 1500)
        })
      }
    }
  }
  return Promise.reject(error)
}

// 请求拦截
request.interceptors.request.use(
  config => {
    const token = storage.get(ACCESS_TOKEN)
    // 如果 token 存在
    // 让每个请求携带自定义 token 请根据实际情况自行修改
    if (token) {
      config.headers['Authorization'] = 'JWT ' + token
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截
request.interceptors.response.use(response => {
  return response.data
}, errorHandler)

const installer = {
  vm: {},
  install(Vue) {
    Vue.use(VueAxios, request)
  }
}

export default request

export { installer as VueAxios, request as axios }
