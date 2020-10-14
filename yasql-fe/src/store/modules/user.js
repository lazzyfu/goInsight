import storage from 'store'
import { login, getInfo, logout } from '@/api/login'
import { ACCESS_TOKEN } from '@/store/mutation-types'
import { welcome } from '@/utils/util'

const user = {
  state: {
    token: '',
    name: '',
    welcome: '',
    avatar: '',
    roles: [],
    info: {}
  },

  mutations: {
    SET_TOKEN: (state, token) => {
      state.token = token
    },
    SET_NAME: (state, { name, welcome }) => {
      state.name = name
      state.welcome = welcome
    },
    SET_AVATAR: (state, avatar) => {
      state.avatar = avatar
    },
    SET_ROLES: (state, roles) => {
      state.roles = roles
    },
    SET_INFO: (state, info) => {
      state.info = info
    }
  },
  

  actions: {
    // 登录
    Login({ commit }, userInfo) {
      return new Promise((resolve, reject) => {
        login(userInfo)
          .then(response => {
            if (response.code === '0000') {
              const result = response.data
              storage.set(ACCESS_TOKEN, result.token, 7 * 24 * 60 * 60 * 1000)
              commit('SET_TOKEN', result.token)
            }
            resolve(response)
          })
          .catch(error => {
            reject(error)
          })
      })
    },

    // 获取用户信息
    GetInfo({ commit }) {
      return new Promise((resolve, reject) => {
        try {
          getInfo()
            .then(response => {
              const result = response.data
              commit('SET_INFO', result)
              commit('SET_NAME', { name: result.username, welcome: welcome() })
              commit('SET_AVATAR', result.avatar)

              resolve(response)
            })
            .catch(error => {
              const errors = [401, 403]
              if (!error || !error.response || errors.includes(error.response.status)) {
                commit('SET_NAME', { name: 'None', welcome: welcome() })
                commit('SET_AVATAR', '')
                reject('unauth')
              }
              reject('unauth')
            })
        } catch (e) {
          // console.log('e: ', e)
        }
      })
    },

    

    // 登出
    Logout({ commit, state }) {
      return new Promise(resolve => {
        logout(state.token)
          .then(() => {
            resolve()
          })
          .catch(() => {
            resolve()
          })
          .finally(() => {
            commit('SET_TOKEN', '')
            commit('SET_ROLES', [])
            storage.remove(ACCESS_TOKEN)
          })
      })
    }
  }
}

export default user
