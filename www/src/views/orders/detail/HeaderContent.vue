<template>
  <a-descriptions size="small" :column="2">
    <a-descriptions-item label="申请人">{{ orderDetail.applicant }}</a-descriptions-item>
    <a-descriptions-item label="工单类型">{{ orderDetail.sql_type }} </a-descriptions-item>
    <a-descriptions-item label="环境">{{ orderDetail.environment }}</a-descriptions-item>
    <a-descriptions-item label="数据库">{{ orderDetail.db_type }}</a-descriptions-item>
    <a-descriptions-item label="实例">{{ orderDetail.instance }}</a-descriptions-item>
    <a-descriptions-item label="库名">{{ orderDetail.schema }}</a-descriptions-item>
    <a-descriptions-item label="认领人">{{ orderDetail.claimer }}</a-descriptions-item>
    <a-descriptions-item label="可领取人">{{ claimUsersText }}</a-descriptions-item>
    <a-descriptions-item label="执行人">{{ orderDetail.executor }}</a-descriptions-item>
    <a-descriptions-item label="更新时间">{{ orderDetail.updated_at }}</a-descriptions-item>
    <a-descriptions-item label="提交时间">{{ orderDetail.created_at }}</a-descriptions-item>
    <a-descriptions-item label="备注">{{ orderDetail.remark }}</a-descriptions-item>
    <a-descriptions-item label="抄送">{{
      orderDetail.cc && orderDetail.cc.length ? orderDetail.cc.join(',') : '无'
    }}</a-descriptions-item>
  </a-descriptions>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  orderDetail: Object,
})

const claimUsersText = computed(() => {
  const raw = props.orderDetail?.claim_users
  if (!raw) return '无'
  if (Array.isArray(raw)) return raw.join(', ')
  if (typeof raw === 'string') {
    try {
      const users = JSON.parse(raw)
      return Array.isArray(users) && users.length > 0 ? users.join(', ') : '无'
    } catch {
      return raw
    }
  }
  return '无'
})
</script>
