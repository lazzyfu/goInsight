<template>
  <a-tabs v-model:activeKey="activeKey">
    <a-tab-pane key="1" tab="审批流定义">
      <FlowList />
    </a-tab-pane>
    <a-tab-pane key="2" tab="审批流配置" force-render>审批流配置</a-tab-pane>
  </a-tabs>
</template>
<script setup>
import { ref } from 'vue'
import FlowList from './FlowList.vue'
const activeKey = ref('1')
</script>
