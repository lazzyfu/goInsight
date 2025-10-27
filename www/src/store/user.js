import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  state: () => ({
    uid: '',
    username: '',
    nickname: '',
    token: '',
    avatar: '',
    email: '',
    mobile: '',
    organization: '',
    role: '',
    date_joined: ''
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
    }
  }
})
