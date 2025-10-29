<template>
  <a-steps direction="vertical" :current="currentProgress">
    <a-step
      v-for="(item, index) in approvals"
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
</template>

<script setup>
import { ref, watch } from 'vue'

const currentProgress = ref(0)
const approvals = ref([])

const props = defineProps({
  approvalList: Array,
})

// 格式化审批数据，计算 current 进度
const formatApprovalData = (data) => {
  const grouped = {}
  data.forEach((item) => {
    if (!grouped[item.stage]) {
      grouped[item.stage] = {
        title: item.stage_name,
        subTitle: item.approval_type === 'AND' ? '与签' : '或签',
        items: [],
      }
    }
    grouped[item.stage].items.push(item)
  })

  const result = Object.values(grouped)
  let current = 0

  result.forEach((step) => {
    const items = step.items
    const approvedCount = items.filter((i) => i.approval_status === 'APPROVED').length

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
              i.approval_status === 'PENDING' ? '' : i.approval_status === 'APPROVED' ? '#52c41a' : '#f5222d'
            };">${
              i.approval_status === 'PENDING' ? '待审批' : i.approval_status === 'APPROVED' ? '已通过' : '已驳回'
            }</span>
          </span>
          <span>备注：<b>${i.approval_remark}</b></span>
          <span>审批时间：<b>${i.approval_time}</b></span>
        </div>
      `,
      )
      .join('')
  })

  currentProgress.value = current
  return result
}

watch(
  () => props.approvalList,
  (newVal) => {
    approvals.value = formatApprovalData(newVal)
  },
  { immediate: true },
)

// 根据步骤计算状态
const getStepStatus = (index, step) => {
  if (step.items.some((i) => i.approval_status === 'REJECTED')) return 'error'
  if (index < currentProgress.value) return 'finish'
  if (index === currentProgress.value) return 'process'
  return 'wait'
}
</script>
