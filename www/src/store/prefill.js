import { defineStore } from 'pinia'

// 仅用于“新建工单页面”的预填数据传递（一次性消费）
export const useOrderCreatePrefillStore = defineStore('orderCreatePrefill', {
  state: () => ({
    createPrefill: null,
  }),

  actions: {
    setCreatePrefill(prefill) {
      this.createPrefill = prefill || null
    },
    consumeCreatePrefill() {
      const prefill = this.createPrefill
      this.createPrefill = null
      return prefill
    },
  },
})

