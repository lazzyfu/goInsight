<template>
  <a-tabs v-model:activeKey="activeKey" >
    <a-tab-pane key="1" tab="编辑器">
      <ConsoleEdit />
    </a-tab-pane>
    <a-tab-pane key="2" tab="收藏SQL">
      <FavoriteSQL />
    </a-tab-pane>
    <a-tab-pane key="3" tab="历史SQL">
      <HistorySQL />
    </a-tab-pane>
  </a-tabs>
</template>

<script setup>
import { ref } from 'vue';
import ConsoleEdit from './console/index.vue';
import FavoriteSQL from './favorite/index.vue';
import HistorySQL from './history/index.vue';

const activeKey = ref('1');


</script>

<style scoped>
body .anticon {
  transform: translate(-0.5px, -3px);
}

</style>
