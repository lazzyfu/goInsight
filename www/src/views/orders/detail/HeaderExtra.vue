<template>
  <a-space>
    <a-button key="3" v-show="btnOptions.status.btnShow" @click="showBtnModal">
      {{ btnTitle }}</a-button
    >
    <a-button
      key="3"
      v-show="
        props.orderDetail?.progress === 'CLAIMED' || props.orderDetail?.progress === 'EXECUTING'
      "
      @click="showCompleteModal"
    >
      已完成</a-button
    >
    <a-button
      key="3"
      v-show="
        props.orderDetail?.progress === 'CLAIMED' || props.orderDetail?.progress === 'EXECUTING'
      "
      @click="showFailModal"
    >
      已失败</a-button
    >
    <a-button key="2" v-show="props.orderDetail?.progress === 'CLAIMED'" @click="showTransferModal">
      转交</a-button
    >
    <a-button key="1" v-show="!btnOptions.status.revokeBtnDisabled" @click="showRevokeModal"
      >撤销</a-button
    >
  </a-space>

  <!-- 审批等操作模态框 -->
  <a-modal v-model:open="btnOptions.open" :title="btnOptions.tips.title" @ok="handleBtnOk">
    <template #footer>
      <a-button key="back" @click="handleBtnCancel">{{ btnOptions.tips.cancelText }}</a-button>
      <a-button key="submit" type="primary" :loading="btnOptions.loading" @click="handleBtnOk">{{
        btnOptions.tips.okText
      }}</a-button>
    </template>

    <a-form
      ref="formRef"
      layout="vertical"
      :model="formState"
      style="margin-top: 24px"
      autocomplete="off"
    >
      <a-form-item
        v-show="btnOptions.tips.action === 'transfer'"
        label="新执行人"
        :rules="
          btnOptions.tips.action === 'transfer'
            ? [{ required: true, message: '请选择新执行人' }]
            : []
        "
        has-feedback
        name="newExecutor"
      >
        <a-select
          v-model:value="formState.newExecutor"
          :options="users"
          :field-names="{ label: 'nick_name', value: 'username' }"
          allowClear
        />
      </a-form-item>
      <a-form-item label="附加信息" has-feedback name="confirmMsg">
        <a-textarea
          :placeholder="btnOptions.tips.placeholder"
          v-model:value="formState.confirmMsg"
          :rows="3"
          allow-clear
        />
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
import { computed, reactive, ref, watch } from 'vue'

const props = defineProps({
  orderDetail: Object,
})
const emit = defineEmits(['refresh'])

const users = ref([])
const formRef = ref()
const formState = reactive({
  newExecutor: '',
  confirmMsg: '',
})

const btnOptions = reactive({
  loading: false,
  open: false,
  tips: {
    okText: '确定',
    cancelText: '取消',
    action: 'approval',
    title: '',
    placeholder: '',
    currentClick: '', // 当前点击的按钮
  },
  status: { btnShow: true, revokeBtnDisabled: false }, // 默认显示btn,撤销按钮
})

const getBtnConfig = (progress) => {
  // 统一容错：进度可能为 undefined/null
  const p = String(progress || '').toUpperCase()
  // 默认配置
  const defaultConfig = {
    title: '',
    tips: { okText: '确定', cancelText: '取消', action: 'approval', title: '', placeholder: '' },
    status: { btnShow: true, revokeBtnDisabled: false },
  }

  switch (p) {
    case 'PENDING': // 待审批
      return {
        title: '审批',
        tips: {
          okText: '同意',
          cancelText: '驳回',
          action: 'approval',
          title: '请审批',
          placeholder: '请输入审批意见...',
        },
        status: { btnShow: true, revokeBtnDisabled: false },
      }
    case 'APPROVED': // 已审批，待认领
      return {
        title: '认领',
        tips: { okText: '认领', cancelText: '取消', action: 'claim', title: '认领任务' },
        status: { btnShow: true, revokeBtnDisabled: false },
      }
    case 'CLAIMED':
    case 'EXECUTING': // 认领或执行中
      return {
        title: '执行',
        tips: { okText: '执行完成', cancelText: '执行中', action: 'execute', title: '执行工单' },
        status: { btnShow: true, revokeBtnDisabled: false },
      }
    case 'COMPLETED': // 执行完成，待复核
      return {
        title: '复核',
        tips: { okText: '确定', cancelText: '取消', action: 'review', title: '复核' },
        status: { btnShow: true, revokeBtnDisabled: true }, // 复核时关闭按钮置灰
      }
    case 'REJECTED':
    case 'REVIEWED':
    case 'REVOKED':
    case 'FAILED':
      return {
        title: '',
        tips: { okText: '确定', cancelText: '取消', action: 'close', title: '' },
        status: { btnShow: false, revokeBtnDisabled: true },
      }
    default:
      return defaultConfig
  }
}

const btnTitle = computed(() => {
  const p = props.orderDetail?.progress
  const cfg = getBtnConfig(p)
  return cfg.title
})

watch(
  () => props.orderDetail?.progress,
  (newProgress) => {
    const cfg = getBtnConfig(newProgress)
    btnOptions.status.btnShow = cfg.status.btnShow
    btnOptions.status.revokeBtnDisabled = cfg.status.revokeBtnDisabled
  },
  { immediate: true },
)

const showBtnModal = async () => {
  const cfg = getBtnConfig(props.orderDetail?.progress)
  btnOptions.tips = { ...cfg.tips }
  btnOptions.open = true
}

// 手动更新为完成
const showCompleteModal = async () => {
  btnOptions.tips = {
    okText: '提交',
    cancelText: '取消',
    action: 'complete',
    title: '确认执行完成？',
    placeholder: '',
  }
  btnOptions.open = true
}

// 手动更新为失败
const showFailModal = async () => {
  btnOptions.tips = {
    okText: '提交',
    cancelText: '取消',
    action: 'fail',
    title: '确认执行失败？',
    placeholder: '请输入失败原因...',
  }
  btnOptions.open = true
}

const showTransferModal = async () => {
  btnOptions.tips = {
    okText: '提交',
    cancelText: '取消',
    action: 'transfer',
    title: '转交工单给其他执行人？',
    placeholder: '请输入转交工单原因...',
  }
  await getUsers()
  btnOptions.open = true
}

const showRevokeModal = async () => {
  btnOptions.tips = {
    okText: '确定',
    cancelText: '取消',
    action: 'revoke',
    title: '确定撤销工单？',
    placeholder: '请输入撤销原因...',
  }
  btnOptions.open = true
}

const getUsers = async () => {
  const res = await getOrderUsersApi().catch((err) => {})
  if (res) {
    users.value = res.data
  }
}

const RequestApi = async () => {
  const { action, currentClick } = btnOptions.tips
  const order_id = props.orderDetail?.order_id
  const payload = { order_id, msg: formState.confirmMsg }

  // 对于认领/撤销，点击取消意味着不发请求，直接关闭弹窗
  if (
    (action === 'claim' ||
      action === 'revoke' ||
      action === 'transfer' ||
      action === 'complete' ||
      action === 'fail') &&
    currentClick === 'cancel'
  ) {
    btnOptions.open = false
    formState.confirmMsg = ''
    return
  }

  btnOptions.loading = true
  try {
    let res = null
    switch (action) {
      case 'approval':
        res = await approvalOrderApi({
          ...payload,
          status: currentClick === 'ok' ? 'APPROVED' : 'REJECTED',
        })
        break
      case 'claim':
        res = await claimOrderApi(payload)
        break
      case 'transfer':
        res = await transferOrderApi({
          ...payload,
          new_executor: formState.newExecutor,
        })
        break
      case 'revoke':
        res = await revokeOrderApi(payload)
        break
      case 'complete':
        res = await completeOrderApi(payload)
        break
      case 'review':
        res = await reviewOrderApi(payload)
        break
      case 'fail':
        res = await failOrderApi(payload)
        break
      default:
        return
    }

    if (res?.code === '0000') {
      message.info('操作成功')
      formState.confirmMsg = ''
      emit('refresh')
    } else {
      // 提示服务端返回的错误信息
      message.error(res?.message || '操作失败')
    }
  } catch (err) {
  } finally {
    btnOptions.open = false
    btnOptions.loading = false
  }
}

const handleBtnOk = async () => {
  btnOptions.tips.currentClick = 'ok'
  await RequestApi()
}
const handleBtnCancel = async () => {
  btnOptions.tips.currentClick = 'cancel'
  await RequestApi()
}
</script>
