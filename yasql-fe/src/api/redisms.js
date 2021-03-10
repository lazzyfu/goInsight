import axios from '@/utils/request'

export default {
  getRedisIns() {
    return axios.request({
      method: 'get',
      url: '/v1/redis/list',
    })
  },
  getRedisCmdList() {
    return axios.request({
      method: 'get',
      url: '/v1/redis/cmds',
    })
  },
  execRedisCmd(data) {
    return axios.request({
      method: 'post',
      url: '/v1/redis/exec_cmd',
      data: data,
    })
  },
  getRedisCmd(pk, db) {
    return axios.request({
      method: 'get',
      url: `/v1/redis/${pk}/metrics`,
      params: {"db": db},
    })
  },
  getRedisHealth(pk, option) {
    return axios.request({
      method: 'get',
      url: `/v1/redis/${pk}/check`,
      params: {"option": option},
    })
  }
}