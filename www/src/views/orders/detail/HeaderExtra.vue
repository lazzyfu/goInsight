<template>
  <!-- 工单操作按钮 -->
  <a-space wrap>
    <!-- 主操作按钮 -->
    <a-button v-if="showMainButton" @click="openMainModal">
      {{ uiData.btnTitle }}
    </a-button>

    <!-- 执行任务按钮 -->
    <a-tooltip title="点击加载执行任务" placement="top">
      <a-button v-if="showExecuteButton" @click="genOrderTasks" :loading="uiState.genOrderTasksLoading">
        执行
      </a-button>
    </a-tooltip>

    <!-- 已完成按钮 -->
    <a-tooltip title="标记当前工单为完成状态" placement="top">
      <a-button v-if="showCompleteFailButton" @click="openCompleteModal"> 已完成 </a-button>
    </a-tooltip>

    <!-- 已失败按钮 -->
    <a-tooltip title="标记当前工单为失败状态" placement="top">
      <a-button v-if="showCompleteFailButton" @click="openFailModal"> 已失败 </a-button>
    </a-tooltip>

    <!-- 转交按钮 -->
    <a-tooltip title="转交当前工单给其他人执行" placement="top">
      <a-button v-if="showTransferButton" @click="openTransferModal"> 转交 </a-button>
    </a-tooltip>

    <!-- 撤销按钮 -->
    <a-tooltip title="撤销当前工单" placement="top">
      <a-button v-if="!revokeDisabled" @click="openRevokeModal"> 撤销 </a-button>
    </a-tooltip>

    <!-- 复制工单 -->
    <a-tooltip title="快速复制当前工单为一个新的工单" placement="top">
      <a-button @click="handleCopyAsNewOrder"> 复制工单 </a-button>
    </a-tooltip>
  </a-space>

  <!-- 操作弹窗 -->
  <a-modal :open="uiData.modalOpen" :title="uiData.modalTitle" @cancel="handleCancel">
    <template #footer>
      <!-- 取消按钮 -->
      <a-button @click="handleCancel">{{ uiData.modalCancelText }}</a-button>
      <!-- 确定按钮 -->
      <a-button type="primary" :loading="uiData.modalLoading" @click="handleSubmit">
        {{ uiData.modalOkText }}
      </a-button>
    </template>

    <!-- 表单 -->
    <a-form ref="formRef" :model="uiData.formData" layout="vertical">
      <!-- 转交操作新执行人 -->
      <a-form-item v-if="uiData.modalAction === 'transfer'" label="新执行人" name="newClaimer"
        :rules="[{ required: true, message: '请选择新执行人' }]">
        <a-select v-model:value="uiData.formData.newClaimer" :options="uiData.userList"
          :field-names="{ label: 'nick_name', value: 'username' }" allowClear style="width: 100%" />
      </a-form-item>

      <!-- 附加信息 -->
      <a-form-item label="附加信息" name="confirmMsg">
        <a-textarea v-model:value="uiData.formData.confirmMsg" :placeholder="uiData.modalPlaceholder" :rows="3"
          allow-clear style="width: 100%" />
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
  formData: { newClaimer: '', confirmMsg: '' },
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
const openMainModal = () => openModal(currentStatusConfig.value.action, uiData.btnTitle)
const openCompleteModal = () => openModal('complete', '确认执行完成？', '提交')
const openFailModal = () => openModal('fail', '确认执行失败？', '提交', '请输入失败原因…')
const openRevokeModal = () => openModal('revoke', '确定撤销工单？', '确定', '请输入撤销原因…')
const openTransferModal = async () => {
  await fetchUsers()
  openModal('transfer', '转交工单给其他执行人？', '提交', '请输入转交原因…')
}

// API 映射表
const API_ACTION_MAP = {
  approval: (payload) => approvalOrderApi({ ...payload, status: 'APPROVED' }),
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
}

// 提交操作
const handleSubmit = useThrottleFn(async () => {
  uiData.modalLoading = true

  const payload = {
    order_id: props.orderDetail?.order_id,
    msg: uiData.formData.confirmMsg,
  }

  try {
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
