<template>
  <a-tabs v-model:activeKey="activeKey">
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
</template>

<script setup>
import { ref } from 'vue'
import ConsoleIndex from './console/ConsoleIndex.vue'
import DasFavorite from './favorite/DasFavorite.vue'
import DasHistory from './history/DasHistory.vue'

const activeKey = ref('1')
</script>

<style scoped>
body .anticon {
  transform: translate(-0.5px, -3px);
}
</style>
