<template>
  <div class="flow-stages-detail">
    <a-list :data-source="definition" size="small" :split="true" :header="false">
      <template #renderItem="{ item, index }">
        <a-list-item class="stage-list-item">
          <div class="stage-info">
            <a-tag color="processing">阶段 {{ index + 1 }}</a-tag>
            <span class="stage-name">{{ item.stage_name }}</span>
          </div>

          <div class="stage-type">
            <a-tag :color="item.type === 'AND' ? 'blue' : 'green'">
              {{ item.type === 'AND' ? '会签 (AND)' : '或签 (OR)' }}
            </a-tag>
          </div>

          <div class="stage-approvers">
            <span class="approvers-count">
              <UserOutlined /> 共 {{ item.approvers.length }} 人
            </span>
            <a-tooltip placement="topLeft" :title="item.approvers.join('; ')">
              <span class="approvers-list">
                {{ displayApprovers(item.approvers) }}
              </span>
            </a-tooltip>
          </div>
        </a-list-item>
      </template>
    </a-list>
  </div>
</template>

<script setup>
import { UserOutlined } from '@ant-design/icons-vue'

// props
defineProps({
  definition: {
    type: Array,
    default: () => [],
    required: true,
  },
})

// 格式化审批人列表，只显示前三个，其余用...表示
const displayApprovers = (approvers) => {
  if (!approvers || approvers.length === 0) {
    return '暂无审批人'
  }
  const limit = 3
  if (approvers.length <= limit) {
    return approvers.join(', ')
  }
  return approvers.slice(0, limit).join(', ') + ` 等 ${approvers.length} 人...`
}
</script>

<style scoped>
.flow-stages-detail {
  background: #fcfcfc;
  border: 1px solid #f0f0f0;
  border-radius: 4px;
}
.stage-list-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px;
}
.stage-list-item:not(:last-child) {
  border-bottom: 1px solid #f0f0f0;
}
.stage-info {
  flex: 3;
  min-width: 200px;
  display: flex;
  align-items: center;
}
.stage-name {
  font-weight: 500;
  margin-left: 12px;
  color: #333;
}
.stage-type {
  flex: 2;
  min-width: 150px;
}
.stage-approvers {
  flex: 5;
  min-width: 300px;
  display: flex;
  align-items: center;
  color: rgba(0, 0, 0, 0.65);
}
.approvers-count {
  margin-right: 15px;
  color: #1890ff;
  font-weight: 500;
}
.approvers-list {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 70%;
}
</style>
