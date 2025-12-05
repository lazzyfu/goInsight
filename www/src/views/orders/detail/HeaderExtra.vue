<template>
  <!-- 工单操作按钮 -->
  <a-space wrap>
    <!-- 主操作按钮 -->
    <a-button v-if="uiState.showMainButton" @click="openMainModal">
      {{ uiData.btnTitle }}
    </a-button>

    <!-- 执行任务按钮 -->
    <a-button v-if="uiState.showExecuteButton" @click="generateExecuteTask"> 执行 </a-button>

    <!-- 已完成按钮 -->
    <a-button v-if="uiState.showCompleteFailButton" @click="openCompleteModal"> 已完成 </a-button>

    <!-- 已失败按钮 -->
    <a-button v-if="uiState.showCompleteFailButton" @click="openFailModal"> 已失败 </a-button>

    <!-- 转交按钮 -->
    <a-button v-if="uiState.showTransferButton" @click="openTransferModal"> 转交 </a-button>

    <!-- 撤销按钮 -->
    <a-button v-if="!uiState.revokeDisabled" @click="openRevokeModal"> 撤销 </a-button>
  </a-space>

  <!-- 操作弹窗 -->
  <a-modal :open="uiData.modalOpen" :title="uiData.modalTitle" width="45%" @cancel="handleCancel">
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
      <a-form-item v-if="uiData.modalAction === 'transfer'" label="新执行人" name="newExecutor"
        :rules="[{ required: true, message: '请选择新执行人' }]">
        <a-select v-model:value="uiData.formData.newExecutor" :options="uiData.userList"
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
  getOrderUsersApi,
  reviewOrderApi,
  revokeOrderApi,
  transferOrderApi,
} from '@/api/order'
import { message } from 'ant-design-vue'
import { reactive, ref, watch } from 'vue'

const props = defineProps({ orderDetail: Object })
const emit = defineEmits(['refresh'])
const formRef = ref()

// UI 状态：只存 bool
const uiState = reactive({
  showMainButton: true,
  showExecuteButton: false,
  showCompleteFailButton: false,
  showTransferButton: false,
  revokeDisabled: false,
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
  formData: { newExecutor: '', confirmMsg: '' },
  userList: [],
})

// 工单状态配置
const getStatusConfig = (progress) => {
  const p = String(progress || '').toUpperCase()
  const map = {
    PENDING: { title: '审批', showMain: true, revokeDisabled: false, action: 'approval' },
    APPROVED: { title: '认领', showMain: true, revokeDisabled: false, action: 'claim' },
    CLAIMED: { showMain: false, revokeDisabled: false, action: 'claim' },
    EXECUTING: { showMain: false, revokeDisabled: false, action: 'execute' },
    COMPLETED: { title: '复核', showMain: true, revokeDisabled: true, action: 'review' },
    REJECTED: { showMain: false, revokeDisabled: true },
    REVIEWED: { showMain: false, revokeDisabled: true },
    FAILED: { showMain: false, revokeDisabled: true },
    REVOKED: { showMain: false, revokeDisabled: true },
  }
  return map[p] || { showMain: false, revokeDisabled: true }
}

// 初始化按钮显示状态
const updateUIState = () => {
  const cfg = getStatusConfig(props.orderDetail?.progress)
  uiData.btnTitle = cfg.title || ''
  uiState.showMainButton = cfg.showMain ?? false
  uiState.revokeDisabled = cfg.revokeDisabled ?? true
  uiState.showExecuteButton = ['CLAIMED', 'EXECUTING', 'COMPLETED', 'FAILED', 'REVIEWED'].includes(
    props.orderDetail?.progress,
  )
  uiState.showCompleteFailButton = ['CLAIMED', 'EXECUTING'].includes(props.orderDetail?.progress)
  uiState.showTransferButton = ['CLAIMED', 'EXECUTING'].includes(props.orderDetail?.progress)
}

// 监听工单进度变化
watch(() => props.orderDetail?.progress, updateUIState, { immediate: true })

// 打开 modal 通用方法
const openModal = (action, title, okText = '确定', placeholder = '') => {
  uiData.modalOpen = true
  uiData.modalAction = action
  uiData.modalTitle = title
  uiData.modalOkText = okText
  uiData.modalCancelText = '取消'
  uiData.modalPlaceholder = placeholder
}

// 打开各操作 modal
const openMainModal = () =>
  openModal(getStatusConfig(props.orderDetail?.progress).action, uiData.btnTitle)
const openCompleteModal = () => openModal('complete', '确认执行完成？', '提交')
const openFailModal = () => openModal('fail', '确认执行失败？', '提交', '请输入失败原因…')
const openTransferModal = async () => {
  await fetchUsers()
  openModal('transfer', '转交工单给其他执行人？', '提交', '请输入转交原因…')
}
const openRevokeModal = () => openModal('revoke', '确定撤销工单？', '确定', '请输入撤销原因…')

// API 映射
const apiMap = {
  approval: approvalOrderApi,
  claim: claimOrderApi,
  transfer: transferOrderApi,
  revoke: revokeOrderApi,
  complete: completeOrderApi,
  review: reviewOrderApi,
  fail: failOrderApi,
}

// 提交 modal 操作
const handleSubmit = async () => {
  uiData.modalLoading = true
  const payload = { order_id: props.orderDetail?.order_id, msg: uiData.formData.confirmMsg }
  try {
    let res
    if (uiData.modalAction === 'approval')
      res = await apiMap.approval({ ...payload, status: 'APPROVED' })
    else if (uiData.modalAction === 'transfer')
      res = await apiMap.transfer({ ...payload, new_executor: uiData.formData.newExecutor })
    else res = await apiMap[uiData.modalAction](payload)

    if (res?.code === '0000') {
      message.success('操作成功')
      emit('refresh')
    } else message.error(res?.message || '操作失败')
  } finally {
    uiData.modalLoading = false
    uiData.modalOpen = false
    formRef.value?.resetFields()
  }
}

// 取消 modal
const handleCancel = () => {
  uiData.modalOpen = false
  formRef.value?.resetFields()
}

// 获取用户列表
const fetchUsers = async () => {
  const res = await getOrderUsersApi().catch(() => null)
  if (res) uiData.userList = res.data
}

// 占位：生成执行任务
const generateExecuteTask = () => message.info('生成执行任务功能尚未实现')
</script>
