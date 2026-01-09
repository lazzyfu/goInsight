import router from '@/router'
import { useUserStore } from '@/store/user'
import { message } from 'ant-design-vue'
import axios from 'axios'
import Nprogress from 'nprogress'

// 请求头
axios.defaults.headers.post['Content-Type'] = 'application/json;charset=UTF-8'
axios.defaults.headers.put['Content-Type'] = 'application/json;charset=UTF-8'

// 跳转登录页
const toLogin = () => {
  router.push({
    path: '/login',
    query: {
      redirect: router.currentRoute.fullPath,
    },
  })
}

// 统一错误处理
const errorHandle = (response) => {
  const { status, data } = response
  const msg = data?.message || '请求错误'

  switch (status) {
    case 400:
      // Bad Request
      message.error(msg)
      break
    case 401:
      message.error(msg)
      localStorage.clear()
      toLogin()
      break
    case 403:
      message.warning('无权限访问！')
      break
    case 404:
      message.error('请求的资源不存在!')
      break
    case 422:
      message.error('token验证失败!')
      localStorage.clear()
      toLogin()
      break
    case 429:
      message.warning('当前访问太过频繁,请稍等再试!')
      break
    case 500:
      message.error('服务端错误,请稍后再试!')
      break
    case 502:
      message.error('网关错误,请稍后再试!')
      break
    case 504:
      message.error('响应超时,请刷新后再试!')
      break
    default:
      message.error(`服务繁忙(${status})`)
  }
}

// 创建 axios 实例
const service = axios.create({
  timeout: 600000,
  baseURL: '',
})

// 请求拦截器
service.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    const token = userStore?.token
    if (token) config.headers.Authorization = 'JWT ' + token

    Nprogress.start()
    return config
  },
  (error) => Promise.reject(error),
)

//  响应拦截器
service.interceptors.response.use(
  (response) => {
    Nprogress.done()
    if (response.status === 200) {
      if (response?.data?.code) {
        if (response.data.code === '4001' || response.data.code === '4002') {
          return Promise.resolve(response)
        }
        if (response.data.code !== '0000') {
          message.error(response?.data?.message || '请求失败')
          return Promise.reject(response)
        }
      }
      return Promise.resolve(response)
    } else {
      return Promise.reject(response)
    }
  },
  (error) => {
    Nprogress.done()
    const { response } = error
    if (response) {
      errorHandle(response)
      return Promise.reject(response)
    }
    if (!window.navigator.onLine) {
      message.warn('网络已断开,请检查网络后重试！')
      router.push({ name: 'refresh' })
    } else {
      message.error('系统异常')
      return Promise.reject(error)
    }
  },
)

export const get = (url, params) => {
  return new Promise((resolve, reject) => {
    service
      .get(url, {
        params: params,
      })
      .then((res) => {
        resolve(res.data)
      })
      .catch((err) => {
        reject(err.data)
      })
  })
}

export const getBlob = (url, params = {}, config = {}) => {
  return service.get(url, {
    params,
    responseType: 'blob',
    ...config,
  })
}

export function post(url, params, headers = {}) {
  return new Promise((resolve, reject) => {
    service
      .post(url, params, headers)
      .then((res) => {
        resolve(res.data)
      })
      .catch((err) => {
        reject(err.data)
      })
  })
}

export function del(url, params) {
  return new Promise((resolve, reject) => {
    service
      .delete(url, { data: params })
      .then((res) => {
        resolve(res.data)
      })
      .catch((err) => {
        reject(err.data)
      })
  })
}
export function put(url, params) {
  return new Promise((resolve, reject) => {
    service
      .put(url, params)
      .then((res) => {
        resolve(res.data)
      })
      .catch((err) => {
        reject(err.data)
      })
  })
}

export default service
