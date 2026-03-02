<template>
  <a-card class="sql-shell" :body-style="{ padding: '14px 16px 16px' }">
    <div class="sql-shell-head">
      <a-typography-title :level="4" style="margin: 0">SQL 查询工作台</a-typography-title>
    </div>

    <a-tabs v-model:activeKey="activeKey" class="sql-tabs">
      <a-tab-pane key="1" tab="编辑器">
        <ConsoleIndex />
      </a-tab-pane>
      <a-tab-pane key="2" tab="收藏SQL">
        <DasFavorite />
      </a-tab-pane>
      <a-tab-pane key="3" tab="历史SQL">
        <DasHistory />
      </a-tab-pane>
    </a-tabs>
  </a-card>
</template>

<script setup>
defineOptions({ name: 'DasIndexView' })

import { ref } from 'vue'
import ConsoleIndex from './console/ConsoleIndex.vue'
import DasFavorite from './favorite/DasFavorite.vue'
import DasHistory from './history/DasHistory.vue'

const activeKey = ref('1')
</script>

<style scoped>
.sql-shell-head {
  margin-bottom: 10px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.sql-tabs :deep(.ant-tabs-nav) {
  margin-bottom: 12px;
}
</style>
