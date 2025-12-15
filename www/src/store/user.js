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
  persist: true,
  actions: {
    setUid(uid) {
      this.uid = uid
    },
    setUserName(username) {
      this.username = username
    },
    setNickName(nickname) {
      this.nickname = nickname
    },
    setUserToken(token) {
      this.token = token
    },
    setUserAvatar(avatar) {
      this.avatar = avatar
    },
    setUserEmail(email) {
      this.email = email
    },
    setUserMobile(mobile) {
      this.mobile = mobile
    },
    setUserOrganization(organization) {
      this.organization = organization
    },
    setUserRole(role) {
      this.role = role
    },
    setUserDateJoined(date_joined) {
      this.date_joined = date_joined
    },
    SetIsSuperuser(is_superuser) {
      this.is_superuser = is_superuser
    },

    async getInfo() {
      const res = await GetUserProfileApi()

      if (res && res.code === '0000') {
        this.setUid(res.data.uid)
        this.setUserName(res.data.username)
        this.setNickName(res.data.nick_name)
        this.setUserAvatar(res.data.avatar_file)
        this.setUserEmail(res.data.email)
        this.setUserMobile(res.data.mobile)
        this.setUserOrganization(res.data.organization)
        this.setUserRole(res.data.role)
        this.setUserDateJoined(res.data.date_joined)
        this.SetIsSuperuser(res.data.is_superuser)

        return res.data
      } else {
        // 强制抛出异常，让导航守卫进入 catch 块
        throw new Error(res?.message || '获取用户信息失败/Token失效')
      }
    },

    clear() {
      this.uid = ''
      this.username = ''
      this.nickname = ''
      this.token = ''
      this.avatar = ''
      this.email = ''
      this.mobile = ''
      this.organization = ''
      this.role = ''
      this.date_joined = ''
      this.is_superuser = false
      localStorage.removeItem('onLine')
    },
  },
})
