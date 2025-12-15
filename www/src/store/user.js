// store/user.js
import { GetUserProfileApi } from '@/api/login'
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  state: () => ({
    uid: '',
    username: '',
    nickname: '',
    token: localStorage.getItem('onLine') || '',
    avatar: '',
    email: '',
    mobile: '',
    organization: '',
    role: '',
    date_joined: '',
    is_superuser: undefined, // 初始为 undefined，用于判断是否已加载
  }),

  persist: {
    key: 'user-store',
    storage: localStorage,
    paths: ['uid', 'username', 'nickname', 'token', 'avatar', 'email', 'mobile', 'organization', 'role', 'date_joined', 'is_superuser'],
  },

  getters: {
    /**
     * 用户全名（用于显示）
     */
    displayName: (state) => state.nickname || state.username || '未知用户',

    /**
     * 是否已登录
     */
    isLoggedIn: (state) => !!state.token,

    /**
     * 用户信息是否已完整加载
     */
    isInfoLoaded: (state) => state.is_superuser !== undefined && state.is_superuser !== null,
  },

  actions: {
    /**
     * 批量设置用户信息
     * @param {Object} userInfo 用户信息对象
     */
    setUserInfo(userInfo) {
      this.uid = userInfo.uid || ''
      this.username = userInfo.username || ''
      this.nickname = userInfo.nick_name || ''
      this.avatar = userInfo.avatar_file || ''
      this.email = userInfo.email || ''
      this.mobile = userInfo.mobile || ''
      this.organization = userInfo.organization || ''
      this.role = userInfo.role || ''
      this.date_joined = userInfo.date_joined || ''
      this.is_superuser = userInfo.is_superuser
    },

    /**
     * 设置 Token
     * @param {String} token 用户 token
     */
    setToken(token) {
      this.token = token
      localStorage.setItem('onLine', token)
    },

    /**
     * 获取用户信息
     * @returns {Promise<Object>} 用户信息
     */
    async getInfo() {
      try {
        const res = await GetUserProfileApi()

        if (res?.code === '0000') {
          this.setUserInfo(res.data)
          return res.data
        }

        throw new Error(res?.message || '获取用户信息失败/Token失效')
      } catch (error) {
        console.error('[User] 获取用户信息失败:', error)
        throw error
      }
    },

    /**
     * 清除用户信息
     */
    clear() {
      this.$reset() // 使用 Pinia 的 $reset 方法重置状态
      this.token = '' // $reset 后需要手动清除 token
      localStorage.removeItem('onLine')
      localStorage.removeItem('user-store')
    },
  },
})

