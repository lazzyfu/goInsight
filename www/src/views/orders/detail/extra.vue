<template>
  <a-space>
    <a-button key="3" v-show="btnStatus.btnShow"> {{ btnTitle }}</a-button>
    <a-button key="2" :disabled="btnStatus.closeBtnDisabled">关闭</a-button>
    <a-button key="1" type="primary">Primary</a-button>
  </a-space>
</template>

<script setup>
import { computed, reactive } from 'vue'

const props = defineProps({
  orderDetail: Object,
})

// 默认显示btn,关闭按钮
const btnStatus = reactive({ btnShow: true, closeBtnDisabled: false })

const btnTitle = computed(() => {
  const progress = props.orderDetail.progress
  switch (progress) {
    case 'PENDING':
      return '审批'
    case 'APPROVED':
      return '认领'
    case ('REJECTED', 'REVIEWED', 'CLOSED'):
      btnStatus.btnShow = false
      btnStatus.closeBtnDisabled = true
      return '已驳回'
    case ('CLAIMED', 'EXECUTING'):
      return '执行'
    case 'COMPLETED':
      return '复核'
  }
})
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


| 枚举        | 显示文案 | 显示颜色    |
| --------- | ---- | ------- |
| PENDING   | 待审批  | default |
| APPROVED  | 已批准  | blue    |
| REJECTED  | 已驳回  | red     |
| CLAIMED   | 已认领  | cyan    |
| EXECUTING | 执行中  | orange  |
| COMPLETED | 已完成  | green   |
| REVIEWED  | 已复核  | green   |
| CLOSED    | 已关闭  | gray    |
-->
