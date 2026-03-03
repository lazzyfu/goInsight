<template>
  <a-card class="sql-shell" :body-style="{ padding: '18px 18px 16px' }">
    <a-tabs v-model:activeKey="activeKey" class="sql-tabs">
      <a-tab-pane key="1" tab="编辑器">
        <div class="tab-pane-content">
          <ConsoleIndex />
        </div>
      </a-tab-pane>
      <a-tab-pane key="2" tab="收藏SQL">
        <div class="tab-pane-content">
          <DasFavorite />
        </div>
      </a-tab-pane>
      <a-tab-pane key="3" tab="历史SQL">
        <div class="tab-pane-content">
          <DasHistory />
        </div>
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
.sql-shell {
  border-radius: 14px;
  border: 1px solid var(--ant-colorBorderSecondary, #f0f0f0);
  box-shadow: 0 8px 24px rgb(15 23 42 / 6%);
  background: var(--ant-colorBgContainer, #ffffff);
}

.sql-tabs :deep(.ant-tabs-nav) {
  margin-bottom: 12px;
}

.sql-tabs :deep(.ant-tabs-tab) {
  border-radius: 999px;
  padding-inline: 14px;
}

.tab-pane-content {
  min-height: 640px;
}

@media (max-width: 900px) {
  .tab-pane-content {
    min-height: 520px;
  }
}
</style>
