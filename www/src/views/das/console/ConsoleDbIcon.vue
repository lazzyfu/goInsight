<template>
  <component :is="getDbIcon" />
</template>

<script setup>
import { computed } from 'vue';
import IconClickHouse from './components/IconClickHouse.vue';
import IconMySQL from './components/IconMySQL.vue';
import IconTiDB from './components/IconTiDB.vue';

const props = defineProps({
  type: {
    type: String,
    required: true,
    validator: (value) => ['clickhouse', 'mysql', 'tidb'].includes(value),
  },
});

const getDbIcon = computed(() => {
  switch (props.type) {
    case 'clickhouse':
      return IconClickHouse;
    case 'mysql':
      return IconMySQL;
    case 'tidb':
      return IconTiDB;
    default:
      return null;
  }
});
</script>
