<template>
  <div class="components-page-header-responsive" style="border: 1px solid rgb(235, 237, 240)">
    <a-page-header
      :title="orderDetail.title"
      class="site-page-header"
      :avatar="{ src: userStore.avatar }"
      @back="() => $router.go(-1)"
    >
      <template #tags>
        <a-tag color="blue">{{ orderDetail.progress }}</a-tag>
      </template>
      <template #extra>
        <a-button key="3">Operation</a-button>
        <a-button key="2">Operation</a-button>
        <a-button key="1" type="primary">Primary</a-button>
      </template>
      <div class="content">
        <div class="main"></div>
        <a-descriptions size="small" :column="2">
          <a-descriptions-item label="申请人">{{ orderDetail.applicant }}</a-descriptions-item>
          <a-descriptions-item label="工单类型">{{ orderDetail.sql_type }} </a-descriptions-item>
          <a-descriptions-item label="环境">{{ orderDetail.environment }}</a-descriptions-item>
          <a-descriptions-item label="数据库">{{ orderDetail.db_type }}</a-descriptions-item>
          <a-descriptions-item label="实例">{{ orderDetail.instance }}</a-descriptions-item>
          <a-descriptions-item label="库名">{{ orderDetail.schema }}</a-descriptions-item>
          <a-descriptions-item label="更新时间">{{ orderDetail.updated_at }}</a-descriptions-item>
          <a-descriptions-item label="提交时间">{{ orderDetail.created_at }}</a-descriptions-item>
          <a-descriptions-item label="备注">{{ orderDetail.remark }}</a-descriptions-item>
          <a-descriptions-item label="抄送">{{
            orderDetail.cc && orderDetail.cc.length ? orderDetail.cc.join(',') : '无'
          }}</a-descriptions-item>
        </a-descriptions>
      </div>
    </a-page-header>
  </div>
  <div style="margin-top: 12px">
    <a-card title="审批流">
      <a-steps direction="vertical" :current="currentApproverProgress">
        <a-step
          v-for="(item, index) in approvalList"
          :key="index"
          :title="item.title"
          :status="getStepStatus(index, item)"
          :sub-title="item.subTitle"
        >
          <template #description>
            <div v-html="item.description"></div>
          </template>
        </a-step>
      </a-steps>
    </a-card>
  </div>
  <div style="margin-top: 12px">
    <a-card title="工单内容"> <CodeMirror ref="codemirrorRef" /> </a-card>
  </div>
</template>

<script setup>
import { getOrderApprovalApi, getOrderHistoryDetailApi } from '@/api/order'
import CodeMirror from '@/components/edit/Codemirror.vue'
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import { useUserStore } from '@/store/user'
const userStore = useUserStore()

const codemirrorRef = ref(null)
const route = useRoute()
const orderId = route.params.order_id

const orderDetail = ref({})
const approvalList = ref([])
const currentApproverProgress = ref(0)
const hasRejected = ref(false)

const getOrderHistoryDetail = async () => {
  const res = await getOrderHistoryDetailApi({
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
    console.log('res: ', res)
    approvalList.value = formatApprovalData(res.data)
    console.log('approvalList.value: ', approvalList.value)
  }
}

// 格式化审批数据，计算 current 进度
const formatApprovalData = (data) => {
  const grouped = {}
  data.forEach((item) => {
    if (!grouped[item.stage]) {
      grouped[item.stage] = {
        title: item.stage_name,
        subTitle: item.approver_type === 'AND' ? '与签' : '或签',
        items: [],
      }
    }
    grouped[item.stage].items.push(item)
  })

  const result = Object.values(grouped)
  let current = 0

  result.forEach((step) => {
    const items = step.items
    const approvedCount = items.filter((i) => i.status === 'APPROVED').length

    if (step.subTitle === '与签') {
      if (approvedCount === items.length) current++
    } else {
      if (approvedCount > 0) current++
    }

    step.description = items
      .map(
        (i) => `
        <div style="
          display: flex;
          align-items: center;
          gap: 24px;
          border-bottom: 1px solid #f0f0f0;
          padding: 4px 0;
        ">
          <span>审批人：<b>${i.approver}</b></span>
          <span>状态：
            <span style="color: ${
              i.status === 'PENDING' ? '#1890ff' : i.status === 'APPROVED' ? '#52c41a' : '#f5222d'
            };">${
              i.status === 'PENDING' ? '待审批' : i.status === 'APPROVED' ? '已通过' : '已驳回'
            }</span>
          </span>
        </div>
      `,
      )
      .join('')
  })

  currentApproverProgress.value = current
  return result
}

// 根据步骤计算状态
const getStepStatus = (index, step) => {
  if (step.items.some((i) => i.status === 'REJECTED')) return 'error'
  if (index < currentApproverProgress.value) return 'finish'
  if (index === currentApproverProgress.value) return 'process'
  return 'wait'
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
  getOrderHistoryDetail()
  getOrderApproval()
})
</script>

<style scoped>
.components-page-header-responsive {
  padding-bottom: 4px;
}
.components-page-header-responsive tr:last-child td {
  padding-bottom: 0;
}
#components-page-header-responsive .content {
  display: flex;
}
#components-page-header-responsive .ant-statistic-content {
  font-size: 20px;
  line-height: 28px;
}
@media (max-width: 576px) {
  #components-page-header-responsive .content {
    display: block;
  }

  #components-page-header-responsive .main {
    width: 100%;
    margin-bottom: 12px;
  }

  #components-page-header-responsive .extra {
    width: 100%;
    margin-left: 0;
    text-align: left;
  }
}
</style>
