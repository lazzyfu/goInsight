<template>
  <div class="order-detail-page gi-page-shell">
    <div class="components-page-header-responsive">
      <a-page-header
        :title="orderDetail.title"
        class="site-page-header"
        :avatar="{ src: userStore.avatar || '/avatar.png' }"
        @back="() => $router.go(-1)"
      >
        <template #tags>
          <template v-if="orderDetail.progress">
            <a-tag :color="getOrderStatusMeta(orderDetail.progress).color">
              {{ getOrderStatusMeta(orderDetail.progress).text }}
            </a-tag>
          </template>
        </template>
        <template #extra>
          <header-extra :order-detail="orderDetail" @refresh="refresh" />
        </template>
        <header-content :order-detail="orderDetail" />
      </a-page-header>
    </div>

    <a-card size="small" title="审批流" class="detail-section-card">
      <approval-steps :approval-status="approvalStatus" />
    </a-card>

    <a-card size="small" title="操作日志" class="detail-section-card">
      <div class="order-logs-container">
        <a-timeline>
          <a-timeline-item v-for="(item, index) in orderLogs" :key="index">
            {{ item.created_at }} {{ item.msg }}
          </a-timeline-item>
        </a-timeline>
      </div>
    </a-card>

    <a-card size="small" title="工单内容" class="detail-section-card">
      <CodeMirror ref="codemirrorRef" :height="'500px'" />
    </a-card>
  </div>
</template>

<script setup>
import { getOrderApprovalStatusApi, getOrderDetailApi, getOrderLogsApi } from '@/api/order'
import CodeMirror from '@/components/edit/Codemirror.vue'
import { useUserStore } from '@/store/user'
import { getOrderStatusMeta } from '@/views/orders/shared/orderStatusMeta'
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
  }).catch(() => {})
  if (res) {
    orderDetail.value = res.data
  }
}

const getOrderApprovalStatus = async () => {
  const res = await getOrderApprovalStatusApi({
    order_id: orderId,
  }).catch(() => {})
  if (res) {
    approvalStatus.value = res.data
  }
}

const getOrderLogs = async () => {
  const res = await getOrderLogsApi({
    order_id: orderId,
  }).catch(() => {})
  if (res) {
    orderLogs.value = res.data
  }
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
.order-detail-page {
  gap: 0;
}

.components-page-header-responsive {
  border: 1px solid rgb(235, 237, 240);
  border-radius: 8px 8px 0 0;
  background: var(--gi-color-container-bg);
  box-shadow: var(--gi-shadow-sm);
}

.detail-section-card {
  margin-top: var(--gi-spacing-ssm);
  border-radius: var(--gi-radius-card);
}

.order-logs-container {
  max-height: 260px;
  padding: var(--gi-spacing-ssm);
  overflow-y: auto;
  overflow-x: hidden;
  word-wrap: break-word;
}

.order-logs-container :deep(.ant-timeline-item-content) {
  color: var(--gi-color-text-secondary);
}
</style>
