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
        <extra-btn :order-detail="orderDetail" />
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
      <approver-steps :approval-list="approvalList" />
    </a-card>
  </div>
  <div style="margin-top: 12px">
    <a-card title="工单内容">
      <CodeMirror ref="codemirrorRef" />
    </a-card>
  </div>
</template>

<script setup>
import { getOrderApprovalApi, getOrderDetailApi } from '@/api/order'
import CodeMirror from '@/components/edit/Codemirror.vue'
import { useUserStore } from '@/store/user'
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import ExtraBtn from './extra.vue'
import ApproverSteps from './steps.vue'

const userStore = useUserStore()
const codemirrorRef = ref(null)
const route = useRoute()
const orderId = route.params.order_id
const orderDetail = ref({})
const approvalList = ref([])

const getOrderHistoryDetail = async () => {
  const res = await getOrderDetailApi({
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
    approvalList.value = res.data
  }
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
