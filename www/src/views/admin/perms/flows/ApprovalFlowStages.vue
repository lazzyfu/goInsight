<template>
  <div class="flow-stages-detail">
    <a-list :data-source="definition" size="small" :split="false" :header="false">
      <template #renderItem="{ item, index }">
        <a-list-item class="stage-list-item">
          <div class="stage-main">
            <span class="stage-order">阶段 {{ index + 1 }}</span>
            <span class="stage-name">{{ item.stage_name }}</span>
            <a-tag :color="item.type === 'AND' ? 'blue' : 'green'" class="type-tag">
              {{ item.type === 'AND' ? '会签 (AND)' : '或签 (OR)' }}
            </a-tag>
          </div>

          <div class="stage-approvers">
            <span class="approvers-count">
              <UserOutlined />
              共 {{ item.approvers.length }} 人
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

defineProps({
  definition: {
    type: Array,
    default: () => [],
    required: true,
  },
})

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
  background: #ffffff;
  border: 1px solid #dfe8f9;
  border-radius: 10px;
}

.stage-list-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
}

.stage-list-item + .stage-list-item {
  border-top: 1px dashed #e6ecf8;
}

.stage-main {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  min-width: 220px;
}

.stage-order {
  display: inline-flex;
  align-items: center;
  height: 24px;
  padding: 0 8px;
  border-radius: 999px;
  background: #eaf3ff;
  color: #1f6feb;
  font-size: 12px;
  font-weight: 600;
}

.stage-name {
  color: #22304f;
  font-weight: 600;
}

.type-tag {
  margin-inline-end: 0;
}

.stage-approvers {
  min-width: 220px;
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
  color: #5f6b8a;
}

.approvers-count {
  color: #1f6feb;
  font-weight: 500;
  white-space: nowrap;
}

.approvers-list {
  max-width: 360px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

@media (max-width: 900px) {
  .stage-list-item {
    flex-direction: column;
    align-items: flex-start;
  }

  .stage-approvers {
    width: 100%;
    justify-content: flex-start;
  }

  .approvers-list {
    max-width: 100%;
  }
}
</style>
