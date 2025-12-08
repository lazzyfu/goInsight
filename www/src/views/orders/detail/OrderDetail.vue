<template>
  <div class="components-page-header-responsive">
    <a-page-header
      :title="orderDetail.title"
      class="site-page-header"
      :avatar="{ src: userStore.avatar }"
      @back="() => $router.go(-1)"
    >
      <template #tags>
        <template v-if="progressInfo = getProgressAlias(orderDetail.progress)">
          <a-tag :color="progressInfo.color">
            {{ progressInfo.text }}
          </a-tag>
        </template>
      </template>
      <template #extra>
        <header-extra :order-detail="orderDetail" @refresh="refresh" />
      </template>
      <header-content :order-detail="orderDetail" />
    </a-page-header>
  </div>

  <a-card size="small" title="审批流" style="margin-top: 12px">
    <approval-steps :approval-status="approvalStatus" />
  </a-card>

  <a-card size="small" title="操作日志" style="margin-top: 12px">
    <div
      style="
        max-height: 260px;
        padding: 12px;
        overflow-y: auto;
        overflow-x: hidden;
        word-wrap: break-word;
      "
    >
      <a-timeline>
        <a-timeline-item v-for="(item, index) in orderLogs" :key="index">
          {{ item.created_at }} {{ item.msg }}
        </a-timeline-item>
      </a-timeline>
    </div>
  </a-card>

  <a-card size="small" title="工单内容" style="margin-top: 12px">
    <CodeMirror ref="codemirrorRef" />
  </a-card>
</template>

<script setup>
import { getOrderApprovalStatusApi, getOrderDetailApi, getOrderLogsApi } from '@/api/order'
import CodeMirror from '@/components/edit/Codemirror.vue'
import { useUserStore } from '@/store/user'
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import ApprovalSteps from './ApprovalSteps.vue'
import HeaderContent from './HeaderContent.vue'
import HeaderExtra from './HeaderExtra.vue'

const userStore = useUserStore()
const codemirrorRef = ref(null)
const route = useRoute()
const orderId = route.params.order_id
const orderDetail = ref({})
const approvalStatus = ref([])
const orderLogs = ref([])

const getOrderDetail = async () => {
  const res = await getOrderDetailApi({
    order_id: orderId,
  }).catch((err) => {})
  if (res) {
    orderDetail.value = res.data
  }
}

const getOrderApprovalStatus = async () => {
  const res = await getOrderApprovalStatusApi({
    order_id: orderId,
  }).catch((err) => {})
  if (res) {
    approvalStatus.value = res.data
  }
}

const getOrderLogs = async () => {
  const res = await getOrderLogsApi({
    order_id: orderId,
  }).catch((err) => {})
  if (res) {
    orderLogs.value = res.data
    console.log('orderLogs.value: ', orderLogs.value)
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
    FAILED: { text: '已失败', color: 'red' },
    REVIEWED: { text: '已复核', color: 'green' },
    REVOKED: { text: '已撤销', color: 'gray' },
  }
  return statusMap[progress] || { text: progress, color: 'default' }
}

const refresh = () => {
  getOrderDetail()
  getOrderApprovalStatus()
  getOrderLogs()
}

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
  getOrderApprovalStatus()
  getOrderLogs()
})
</script>

<style scoped>
.components-page-header-responsive {
  border: 1px solid rgb(235, 237, 240);
  border-radius: 8px 8px 0 0;
}
</style>
