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
  baseURL: '',
})

function showErrorNotification(error) {
  if (error.response) {
    const message = error.response.status
    const description = error.response.data.msg || 'An unknown error occurred.'
    notification.error({ message, description })
  } else {
    notification.error({ message: 'Network error' })
  }
}

const errorHandler = (error) => {
  if (error.response) {
    console.log('error.response: ', error.response)
    // switch (error.response.status) {
    //   case 403:
    //     showErrorNotification(error)
    //     break
    //   case 404:
    //     redirect({ name: '404' })
    //     break
    //   case 500:
    //     showErrorNotification(error)
    //     break
    //   case 401:
    //     showErrorNotification(error)
    //     break
    //   default:
    //     showErrorNotification(error)
    //     break
    // }
  }
  return Promise.reject(error)
}

// request interceptor
request.interceptors.request.use(
  (config) => {
    const token = storage.get(ACCESS_TOKEN)
    // 如果 token 存在
    // 让每个请求携带自定义 token 请根据实际情况自行修改
    if (token) {
      config.headers['Authorization'] = 'JWT ' + token
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// response interceptor
request.interceptors.response.use((response) => {
  return response.data
}, errorHandler)

const installer = {
  vm: {},
  install(Vue) {
    Vue.use(VueAxios, request)
  },
}

export default request

export { installer as VueAxios, request as axios }
