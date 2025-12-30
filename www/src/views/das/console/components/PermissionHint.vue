<template>
  <span v-if="hasAccess" class="permission-hint">
    <a-tooltip>
      <template #title>您有该表的数据访问权限</template>
      <TableOutlined :style="{ fontSize: '14px', color: '#52c41a' }" />
    </a-tooltip>
  </span>
  <span v-else class="permission-hint">
    <a-tooltip>
      <template #title>您没有该表的数据访问权限</template>
      <CloseCircleTwoTone two-tone-color="#eb2f96" />
    </a-tooltip>
  </span>
</template>

<script setup>
import { CloseCircleTwoTone, TableOutlined } from '@ant-design/icons-vue';

defineProps({
  hasAccess: {
    type: Boolean,
    required: true,
  },
})
</script>

<style scoped>
.permission-hint {
  margin-right: 4px;
}
</style>
