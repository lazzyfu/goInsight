<template>
  <!-- 工单操作按钮 -->
  <a-space class="header-action-group">
    <!-- 主操作按钮 -->
    <a-button v-if="showMainButton" type="primary" size="small" @click="openMainModal">
      {{ uiData.btnTitle }}
    </a-button>

    <!-- 执行任务按钮 -->
    <a-tooltip title="点击加载执行任务" placement="top">
      <a-button v-if="showExecuteButton" size="small" @click="genOrderTasks" :loading="uiState.genOrderTasksLoading">
        执行
      </a-button>
    </a-tooltip>

    <!-- 已完成按钮 -->
    <a-tooltip title="标记当前工单为完成状态" placement="top">
      <a-button v-if="showCompleteFailButton" size="small" @click="openCompleteModal"> 已完成 </a-button>
    </a-tooltip>

    <!-- 已失败按钮 -->
    <a-tooltip title="标记当前工单为失败状态" placement="top">
      <a-button v-if="showCompleteFailButton" size="small" @click="openFailModal"> 已失败 </a-button>
    </a-tooltip>

    <!-- 转交按钮 -->
    <a-tooltip title="转交当前工单给其他人执行" placement="top">
      <a-button v-if="showTransferButton" size="small" @click="openTransferModal"> 转交 </a-button>
    </a-tooltip>

    <!-- 撤销按钮 -->
    <a-tooltip title="撤销当前工单" placement="top">
      <a-button v-if="!revokeDisabled" size="small" @click="openRevokeModal"> 撤销 </a-button>
    </a-tooltip>

    <!-- 复制工单 -->
    <a-tooltip title="快速复制当前工单为一个新的工单" placement="top">
      <a-button size="small" @click="handleCopyAsNewOrder"> 复制工单 </a-button>
    </a-tooltip>
  </a-space>

  <!-- 操作弹窗 -->
  <a-modal :open="uiData.modalOpen" :title="uiData.modalTitle" @cancel="handleCancel">
    <template #footer>
      <template v-if="uiData.modalAction === 'approval'">
        <a-button danger :loading="uiData.modalLoading" @click="handleSubmit('REJECTED')">驳回</a-button>
        <a-button type="primary" :loading="uiData.modalLoading" @click="handleSubmit('APPROVED')">通过</a-button>
      </template>
      <a-button v-else type="primary" :loading="uiData.modalLoading" @click="handleSubmit()">
        {{ uiData.modalOkText }}
      </a-button>
    </template>

    <!-- 表单 -->
    <a-form ref="formRef" :model="uiData.formData" layout="vertical">
      <!-- 转交操作新执行人 -->
      <a-form-item v-if="uiData.modalAction === 'transfer'" label="新执行人" name="newClaimer"
        :rules="[{ required: true, message: '请选择新执行人' }]">
        <a-select v-model:value="uiData.formData.newClaimer" :options="uiData.userList"
          :field-names="{ label: 'nick_name', value: 'username' }" allowClear class="modal-field-full-width" />
      </a-form-item>

      <!-- 附加信息 -->
      <a-form-item
        :label="uiData.modalAction === 'approval' ? '审批意见' : '附加信息'"
        name="confirmMsg"
        :rules="[
          { max: 256, message: '附加信息不能超过256个字符' },
          { validator: validateConfirmMsg, trigger: 'blur' },
        ]"
      >
        <a-textarea
          v-model:value="uiData.formData.confirmMsg"
          :placeholder="confirmMsgPlaceholder"
          :rows="3"
          allow-clear class="modal-field-full-width" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import {
    approvalOrderApi,
    claimOrderApi,
    completeOrderApi,
    failOrderApi,
    generateOrderTasksApi,
    getOrderUsersApi,
    reviewOrderApi,
    revokeOrderApi,
    transferOrderApi
} from '@/api/order'
import { useOrderCreatePrefillStore } from '@/store/prefill'
import { useUserStore } from '@/store/user'
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { computed, reactive, ref, toRaw, watch } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({ orderDetail: Object })
const emit = defineEmits(['refresh'])
const formRef = ref()
const userStore = useUserStore()
const router = useRouter()
const orderCreatePrefillStore = useOrderCreatePrefillStore()

// 工单状态常量
const ORDER_STATUS = {
  PENDING: 'PENDING',
  APPROVED: 'APPROVED',
  CLAIMED: 'CLAIMED',
  EXECUTING: 'EXECUTING',
  COMPLETED: 'COMPLETED',
  REJECTED: 'REJECTED',
  REVIEWED: 'REVIEWED',
  FAILED: 'FAILED',
  REVOKED: 'REVOKED',
}

// 可执行任务的状态列表
const EXECUTABLE_STATUSES = [
  ORDER_STATUS.CLAIMED,
  ORDER_STATUS.EXECUTING,
  ORDER_STATUS.COMPLETED,
  ORDER_STATUS.FAILED,
  ORDER_STATUS.REVIEWED,
]

// 可完成/失败的状态列表
const COMPLETABLE_STATUSES = [ORDER_STATUS.CLAIMED, ORDER_STATUS.EXECUTING]

// UI 状态
const uiState = reactive({
  genOrderTasksLoading: false,
})

// UI 数据
const uiData = reactive({
  btnTitle: '',
  modalOpen: false,
  modalLoading: false,
  modalAction: '',
  modalTitle: '',
  modalOkText: '确定',
  modalCancelText: '取消',
  modalPlaceholder: '',
  formData: { newClaimer: '', confirmMsg: '', approvalStatus: 'APPROVED' },
  userList: [],
})

// 工单状态配置映射
const STATUS_CONFIG_MAP = {
  [ORDER_STATUS.PENDING]: { title: '审批', showMain: true, revokeDisabled: false, action: 'approval' },
  [ORDER_STATUS.APPROVED]: { title: '认领', showMain: true, revokeDisabled: false, action: 'claim' },
  [ORDER_STATUS.CLAIMED]: { showMain: false, revokeDisabled: false, action: 'claim' },
  [ORDER_STATUS.EXECUTING]: { showMain: false, revokeDisabled: false, action: 'execute' },
  [ORDER_STATUS.COMPLETED]: { title: '复核', showMain: true, revokeDisabled: true, action: 'review' },
  [ORDER_STATUS.REJECTED]: { showMain: false, revokeDisabled: true },
  [ORDER_STATUS.REVIEWED]: { showMain: false, revokeDisabled: true },
  [ORDER_STATUS.FAILED]: { showMain: false, revokeDisabled: true },
  [ORDER_STATUS.REVOKED]: { showMain: false, revokeDisabled: true },
}

// 获取当前工单状态配置
const currentStatusConfig = computed(() => {
  const progress = props.orderDetail?.progress?.toUpperCase()
  return STATUS_CONFIG_MAP[progress] || { showMain: false, revokeDisabled: true }
})

// 按钮显示状态计算属性
const showMainButton = computed(() => currentStatusConfig.value.showMain ?? false)
const revokeDisabled = computed(() => currentStatusConfig.value.revokeDisabled ?? true)
const showExecuteButton = computed(() => EXECUTABLE_STATUSES.includes(props.orderDetail?.progress))
const showCompleteFailButton = computed(() => COMPLETABLE_STATUSES.includes(props.orderDetail?.progress))
const showTransferButton = computed(
  () =>
    props.orderDetail?.progress === ORDER_STATUS.CLAIMED &&
    props.orderDetail?.claimer === userStore.username,
)

// 更新按钮标题
watch(currentStatusConfig, (config) => {
  uiData.btnTitle = config.title || ''
}, { immediate: true })

// 打开弹窗通用方法
const openModal = (action, title, okText = '确定', placeholder = '') => {
  Object.assign(uiData, {
    modalOpen: true,
    modalAction: action,
    modalTitle: title,
    modalOkText: okText,
    modalCancelText: '取消',
    modalPlaceholder: placeholder,
  })
}

// 打开各操作弹窗
const openMainModal = () => {
  if (currentStatusConfig.value.action === 'approval') {
    uiData.formData.approvalStatus = 'APPROVED'
  }
  openModal(currentStatusConfig.value.action, uiData.btnTitle)
}
const openCompleteModal = () => openModal('complete', '确认执行完成？', '提交')
const openFailModal = () => openModal('fail', '确认执行失败？', '提交', '请输入失败原因…')
const openRevokeModal = () => openModal('revoke', '确定撤销工单？', '确定', '请输入撤销原因…')
const openTransferModal = async () => {
  await fetchUsers()
  openModal('transfer', '转交工单给其他执行人？', '提交', '请输入转交原因…')
}

// API 映射表
const API_ACTION_MAP = {
  approval: (payload) => approvalOrderApi({ ...payload, status: uiData.formData.approvalStatus }),
  claim: claimOrderApi,
  transfer: (payload) => transferOrderApi({ ...payload, new_claimer: uiData.formData.newClaimer }),
  revoke: revokeOrderApi,
  complete: completeOrderApi,
  review: reviewOrderApi,
  fail: failOrderApi,
}

// 关闭弹窗并重置表单
const closeModal = () => {
  uiData.modalOpen = false
  formRef.value?.resetFields()
  uiData.formData.approvalStatus = 'APPROVED'
  uiData.formData.newClaimer = ''
  uiData.formData.confirmMsg = ''
}

const confirmMsgPlaceholder = computed(() => {
  if (uiData.modalAction === 'approval') {
    return uiData.formData.approvalStatus === 'REJECTED'
      ? '请填写驳回原因（必填）'
      : '可选：填写审批说明'
  }
  return uiData.modalPlaceholder
})

const validateConfirmMsg = async (_rule, value) => {
  if (
    uiData.modalAction === 'approval' &&
    uiData.formData.approvalStatus === 'REJECTED' &&
    !String(value || '').trim()
  ) {
    throw new Error('驳回时必须填写原因')
  }
}

// 提交操作
const handleSubmit = useThrottleFn(async (approvalStatus) => {
  if (uiData.modalAction === 'approval' && approvalStatus) {
    uiData.formData.approvalStatus = approvalStatus
  }

  uiData.modalLoading = true

  const payload = {
    order_id: props.orderDetail?.order_id,
    msg: uiData.formData.confirmMsg,
  }

  try {
    await formRef.value?.validate?.()
    const apiAction = API_ACTION_MAP[uiData.modalAction]
    const res = await apiAction(payload).catch(() => null)

    if (res) {
      message.success('操作成功')
      emit('refresh')
      closeModal()
    }
  } finally {
    uiData.modalLoading = false
  }
})

// 取消操作
const handleCancel = closeModal

// 获取用户列表
const fetchUsers = async () => {
  const claimUsers = props.orderDetail?.claim_users
  if (Array.isArray(claimUsers) && claimUsers.length > 0) {
    uiData.userList = claimUsers.map((username) => ({ username, nick_name: username }))
    return
  }
  if (typeof claimUsers === 'string' && claimUsers.length > 0) {
    try {
      const parsed = JSON.parse(claimUsers)
      if (Array.isArray(parsed) && parsed.length > 0) {
        uiData.userList = parsed.map((username) => ({ username, nick_name: username }))
        return
      }
    } catch {
      // Ignore parse failure and fallback to full user list.
    }
  }
  const res = await getOrderUsersApi().catch(() => null)
  if (res) uiData.userList = res.data
}

// 生成执行任务
const genOrderTasks = useThrottleFn(async () => {
  uiState.genOrderTasksLoading = true
  message.info('正在加载执行任务，请稍候...')

  try {
    const res = await generateOrderTasksApi({
      order_id: props.orderDetail?.order_id,
    }).catch(() => null)

    if (res) {
      router.push({
        name: 'orders.tasks',
        params: { order_id: props.orderDetail?.order_id },
      })
    }
  } finally {
    uiState.genOrderTasksLoading = false
  }
})

const clonePlainObject = (value) => {
  if (!value) return value
  try {
    // 优先使用 structuredClone（更健壮），否则退化为 JSON 克隆
    if (typeof structuredClone === 'function') return structuredClone(value)
    return JSON.parse(JSON.stringify(value))
  } catch {
    return null
  }
}

// 从 orderDetail 生成“新建工单”的预填数据（只保留表单需要的字段）
const toCreateOrderPrefill = (orderDetail) => {
  if (!orderDetail) return null
  const raw = clonePlainObject(toRaw(orderDetail))
  if (!raw) return null
  return {
    title: raw.title || '',
    remark: raw.remark || '',
    sql_type: raw.sql_type || '',
    cc: raw.cc || [],
    content: raw.content || '',
  }
}

// 复制当前工单为新工单：设置预填并跳转到创建页
const handleCopyAsNewOrder = () => {
  orderCreatePrefillStore.setCreatePrefill(toCreateOrderPrefill(props.orderDetail))
  router.push({ name: 'orders.create' })
}
</script>

<style scoped>
.header-action-group {
  align-items: center;
  flex-wrap: nowrap;
  gap: var(--gi-spacing-xs) var(--gi-spacing-sm);
}

.header-action-group :deep(.ant-btn) {
  border-radius: var(--gi-radius-md);
}

.modal-field-full-width {
  width: 100%;
}

@media (max-width: 767px) {
  .header-action-group {
    width: 100%;
  }
}
</style>
