<template>
  <div class="components-page-header-responsive">
    <a-page-header
      :title="orderDetail.title"
      class="site-page-header"
      :avatar="{ src: userStore.avatar }"
      @back="() => $router.go(-1)"
    >
      <template #tags>
        <a-tag :color="progressInfo.color">{{ progressInfo.text }}</a-tag>
      </template>
      <template #extra>
        <header-extra :order-detail="orderDetail" @refresh="refresh" />
      </template>
      <header-content :order-detail="orderDetail" />
    </a-page-header>
  </div>
  <a-card title="审批流" style="margin-top: 12px">
    <approval-steps :approval-data="approvalData" />
  </a-card>
  <a-card title="工单内容" style="margin-top: 12px">
    <CodeMirror ref="codemirrorRef" />
  </a-card>
</template>

<script setup>
import { getOrderApprovalApi, getOrderDetailApi } from '@/api/order'
import CodeMirror from '@/components/edit/Codemirror.vue'
import { useUserStore } from '@/store/user'
import { onMounted, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import ApprovalSteps from './ApprovalSteps.vue'
import HeaderContent from './HeaderContent.vue'
import HeaderExtra from './HeaderExtra.vue'

const userStore = useUserStore()
const codemirrorRef = ref(null)
const route = useRoute()
const orderId = route.params.order_id
const orderDetail = ref({})
const approvalData = ref([])
const progressInfo = reactive({ text: '', color: '' })

const getOrderDetail = async () => {
  const res = await getOrderDetailApi({
    order_id: orderId,
  }).catch((err) => {})
  if (res) {
    console.log('res: ', res)
    orderDetail.value = res.data
  }
}

const getOrderApproval = async () => {
  const res = await getOrderApprovalApi({
    order_id: orderId,
  }).catch((err) => {})
  if (res) {
    console.log('res: ', res);
    approvalData.value = res.data
  }
}

const getProgressAlias = (progress) => {
  const statusMap = {
    PENDING: { text: '待审批', color: 'default' },
    APPROVED: { text: '已批准', color: 'blue' },
    REJECTED: { text: '已驳回', color: 'red' },
    CLAIMED: { text: '已认领', color: 'cyan' },
    EXECUTING: { text: '执行中', color: 'orange' },
    COMPLETED: { text: '已完成', color: 'green' },
    REVIEWED: { text: '已复核', color: 'green' },
    CLOSED: { text: '已关闭', color: 'gray' },
  }
  return statusMap[progress] || { text: progress, color: 'default' }
}

const refresh = () => {
  getOrderDetail()
  getOrderApproval()
}

watch(
  () => orderDetail.value.progress,
  (newValue) => {
    const p = getProgressAlias(newValue)
    progressInfo.text = p.text
    progressInfo.color = p.color
  },
)
watch(
  () => orderDetail.value.content,
  (newValue) => {
    if (codemirrorRef.value) {
      codemirrorRef.value.setHeight(650)
      codemirrorRef.value.setContent(newValue)
      codemirrorRef.value.setReadonly(true)
    }
  },
  { immediate: true },
)

onMounted(async () => {
  getOrderDetail()
  getOrderApproval()
})
</script>

<style scoped>
.components-page-header-responsive {
  border: 1px solid rgb(235, 237, 240);
  border-radius: 8px 8px 0 0;
}
</style>
