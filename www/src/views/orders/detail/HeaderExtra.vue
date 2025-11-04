<template>
  <a-space>
    <a-button key="3" v-show="btnOptions.status.btnShow" @click="showBtnModal">
      {{ btnTitle }}</a-button
    >
    <a-button key="2" v-show="props.orderDetail?.progress === 'CLAIMED'" @click="showBtnModal">
      转交</a-button
    >
    <a-button key="1" v-show="!btnOptions.status.closeBtnDisabled" @click="showCloseModal"
      >关闭</a-button
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
    <a-textarea :placeholder="btnOptions.tips.placeholder" v-model:value="confirmMsg" :rows="3" />
  </a-modal>
</template>

<script setup>
import { approvalOrderApi, claimOrderApi, closeOrderApi } from '@/api/order'
import { message } from 'ant-design-vue'
import { computed, reactive, ref, watch } from 'vue'

const props = defineProps({
  orderDetail: Object,
})
const emit = defineEmits(['refresh'])
const confirmMsg = ref('')

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
  status: { btnShow: true, closeBtnDisabled: false }, // 默认显示btn,关闭按钮
})

const getBtnConfig = (progress) => {
  // 统一容错：进度可能为 undefined/null
  const p = String(progress || '').toUpperCase()
  // 默认配置
  const defaultConfig = {
    title: '',
    tips: { okText: '确定', cancelText: '取消', action: 'approval', title: '', placeholder: '' },
    status: { btnShow: true, closeBtnDisabled: false },
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
        status: { btnShow: true, closeBtnDisabled: false },
      }
    case 'APPROVED': // 已审批，待认领
      return {
        title: '认领',
        tips: { okText: '认领', cancelText: '取消', action: 'claim', title: '认领任务' },
        status: { btnShow: true, closeBtnDisabled: false },
      }
    case 'CLAIMED':
    case 'EXECUTING': // 认领或执行中
      return {
        title: '执行',
        tips: { okText: '执行完成', cancelText: '执行中', action: 'execute', title: '执行' },
        status: { btnShow: true, closeBtnDisabled: false },
      }
    case 'COMPLETED': // 执行完成，待复核
      return {
        title: '复核',
        tips: { okText: '确定', cancelText: '取消', action: 'review', title: '复核' },
        status: { btnShow: true, closeBtnDisabled: true }, // 复核时关闭按钮置灰示例
      }
    case 'REJECTED':
    case 'REVIEWED':
    case 'CLOSED':
      return {
        title: '',
        tips: { okText: '确定', cancelText: '取消', action: 'close', title: '' },
        status: { btnShow: false, closeBtnDisabled: true },
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
    btnOptions.status.closeBtnDisabled = cfg.status.closeBtnDisabled
  },
  { immediate: true },
)

const showBtnModal = () => {
  const cfg = getBtnConfig(props.orderDetail?.progress)
  btnOptions.tips = { ...cfg.tips }
  btnOptions.open = true
}

const showCloseModal = () => {
  btnOptions.tips = {
    okText: '确定',
    cancelText: '取消',
    action: 'close',
    title: '确定关闭工单?',
    placeholder: '请输入关闭原因...',
  }
  btnOptions.open = true
}

const RequestApi = async () => {
  let res = null
  switch (btnOptions.tips.action) {
    case 'approval':
      res = await approvalOrderApi({
        order_id: props.orderDetail?.order_id,
        status: btnOptions.tips.currentClick === 'ok' ? 'APPROVED' : 'REJECTED',
        msg: confirmMsg.value,
      }).catch((err) => {})
      if (res?.code === '0000') {
        message.info('操作成功')
        emit('refresh')
      }
      break
    case 'claim':
      res = await claimOrderApi({
        order_id: props.orderDetail?.order_id,
        msg: confirmMsg.value,
      }).catch((err) => {})
      if (res?.code === '0000') {
        message.info('操作成功')
        emit('refresh')
      }
    case 'close':
      if (btnOptions.tips.currentClick === 'cancel') {
        return
      }
      res = await closeOrderApi({
        order_id: props.orderDetail?.order_id,
        msg: confirmMsg.value,
      }).catch((err) => {})
      if (res?.code === '0000') {
        message.info('操作成功')
        emit('refresh')
      }
    default:
      break
  }
}

const handleBtnOk = async () => {
  btnOptions.loading = true
  btnOptions.tips.currentClick = 'ok'
  RequestApi()
  btnOptions.loading = false
  btnOptions.open = false
}
const handleBtnCancel = () => {
  btnOptions.tips.currentClick = 'cancel'
  RequestApi()
  btnOptions.open = false
}
</script>

<!--
| 阶段            | 描述                | 示例触发方     |
| ------------- | ----------------- | --------- |
| **PENDING**   | 待审批（创建后进入此状态）     | 工单提交人     |
| **APPROVED**  | 已批准，待执行           | 审批人同意     |
| **REJECTED**  | 已驳回，流程终止          | 审批人驳回     |
| **CLAIMED**   | 已认领，执行人接单         | 执行人主动认领   |
| **EXECUTING** | 执行中               | 执行人操作     |
| **COMPLETED** | 执行完成，待复核          | 执行人提交结果   |
| **REVIEWED**  | 已复核，流程结束          | 复核人通过     |
| **CLOSED**    | 已关闭，非正常终止（例如人工关闭） | 任意角色（管理方） |



-->
