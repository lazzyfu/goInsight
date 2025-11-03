<template>
  <div class="approval-steps-container">
    <a-steps direction="vertical" :current="currentProgress">
      <a-step
        v-for="(item, index) in approvalStages"
        :key="index"
        :title="item.title"
        :status="getStepStatus(index, item)"
        :sub-title="item.subTitle"
      >
        <template #description>
          <div class="approvals-list">
            <div v-for="(i, index) in item.items" :key="index" class="approval-item">
              <div v-if="index > 0" class="approval-divider"></div>
              <div class="approval-info-grid">
                <div class="approval-info-item">
                  <span class="info-label">审批人：</span>
                  <span class="info-value">{{ i.approver }}</span>
                </div>
                <div class="approval-info-item">
                  <span class="info-label">状态：</span>
                  <a-tag :color="getStatusColor(i.approval_status)" size="small">
                    {{ getStatusText(i.approval_status) }}
                  </a-tag>
                </div>
                <div class="approval-info-item">
                  <span class="info-label">审批时间：</span>
                  <span class="info-value">{{ i.approval_at }}</span>
                </div>
              </div>
              <div class="approval-msg">
                <span class="info-label">备注：</span>
                <span class="info-value">{{ i.msg }}</span>
              </div>
            </div>
          </div>
        </template>
      </a-step>
    </a-steps>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  approvalStatus: Array,
})

const currentProgress = ref(0)
const approvalStages = ref([])

// 格式化审批数据，计算 current 进度
const formatApprovalStatus = (data) => {
  const grouped = {}
  data.forEach((item) => {
    if (!grouped[item.stage]) {
      grouped[item.stage] = {
        title: item.stage_name,
        subTitle: item.approval_type === 'AND' ? '会签' : '或签（任一审批人通过即可）',
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

    if (step.subTitle === '会签') {
      if (approvedCount === items.length) current++
    } else {
      if (approvedCount > 0) current++
    }
  })

  currentProgress.value = current
  return result
}

const getStatusText = (status) => {
  const textMap = {
    APPROVED: '已通过',
    REJECTED: '已驳回',
    PENDING: '待审批',
  }
  return textMap[status] || status
}

const getStatusColor = (status) => {
  const colorMap = {
    APPROVED: '#52c41a',
    REJECTED: '#f5222d',
    PENDING: '',
  }
  return colorMap[status] || ''
}

watch(
  () => props.approvalStatus,
  (newVal) => {
    approvalStages.value = formatApprovalStatus(newVal)
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

<style scoped>
.approval-steps-container {
  padding: 12px;
  background: #fff;
}

.approvals-list {
  padding-top: 8px;
}

.approval-item {
  margin-bottom: 12px;
}

.approval-item:last-child {
  margin-bottom: 0;
}

.approval-divider {
  height: 1px;
  background-color: #f0f0f0;
  margin: 12px 0;
}

.approval-info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 8px 24px;
  align-items: center;
}

.approval-info-item {
  display: flex;
  align-items: center;
  min-width: 0;
}

.info-label {
  font-size: 14px;
  white-space: nowrap;
  flex-shrink: 0;
}

.info-value {
  font-size: 14px;
  word-break: break-word;
}

.approval-msg {
  display: flex;
  gap: 4px;
  margin-top: 8px;
  padding-top: 8px;
}

.approval-msg .info-label {
  flex-shrink: 0;
}

.approval-msg .info-value {
  flex: 1;
  word-break: break-word;
  line-height: 1.6;
}

.approval-type-hint {
  margin-top: 8px;
  color: rgba(0, 0, 0, 0.45);
  font-size: 12px;
  font-style: italic;
}

@media (max-width: 768px) {
  .approval-info-grid {
    grid-template-columns: 1fr;
    gap: 8px;
  }
}
</style>
