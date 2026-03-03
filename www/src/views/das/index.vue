<template>
  <a-card class="sql-shell" :body-style="{ padding: '16px 16px 14px' }">
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
  border-radius: 12px;
  border: 1px solid var(--ant-colorBorderSecondary, #f0f0f0);
  box-shadow: 0 2px 8px rgb(0 0 0 / 5%);
  background: var(--ant-colorBgContainer, #ffffff);
}

.sql-tabs :deep(.ant-tabs-nav) {
  margin-bottom: 10px;
}

.sql-tabs :deep(.ant-tabs-tab) {
  border-radius: 8px 8px 0 0;
  padding-inline: 12px;
}

.tab-pane-content {
  min-height: 660px;
}

@media (max-width: 900px) {
  .tab-pane-content {
    min-height: 520px;
  }
}
</style>
