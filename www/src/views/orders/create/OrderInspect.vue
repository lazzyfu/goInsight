<template>
  <div v-if="data.showbTable" style="margin-top: 15px; overflow-x: auto">
    <!-- Added statistics summary bar -->
    <div class="stats-bar">
      <span class="stats-item">
        共 <strong>{{ data.tableData.length }}</strong> 条记录
      </span>
      <span v-if="data.summaryStats.error > 0" class="stats-item error-stat">
        <span class="stat-icon">✕</span> {{ data.summaryStats.error }} 个错误
      </span>
      <span v-if="data.summaryStats.warn > 0" class="stats-item warn-stat">
        <span class="stat-icon">⚠</span> {{ data.summaryStats.warn }} 个警告
      </span>
    </div>

    <a-table size="small" class="ant-table-striped" bordered :data-source="data.tableData" :columns="tableColumns"
      :rowClassName="setRowClass" :scroll="{ x: 1100 }"
      :pagination="{ pageSize: 10, showSizeChanger: true, showTotal: (total) => `共 ${total} 条` }">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'summary'">
          <div v-if="record.summary && record.summary.length > 0" class="summary-container">
            <!-- Added expand/collapse for many items -->
            <div v-for="(item, index) in getVisibleSummary(record)" :key="index" class="summary-item">
              <span :class="`summary-level-${item.level?.toLowerCase()}`">
                {{ item.level }}：{{ item.message }}
              </span>
            </div>
            <a v-if="record.summary.length > 3 && !record.expanded" class="expand-link" @click="toggleExpand(record)">
              + 展开全部 ({{ record.summary.length - 3 }} 条)
            </a>
            <a v-if="record.expanded" class="expand-link" @click="toggleExpand(record)">
              - 收起
            </a>
          </div>
          <div v-else class="empty-summary">暂无检查提示</div>
        </template>

        <template v-if="column.key === 'level'">
          <a-tag :color="getLevelColor(record.level)">
            {{ record.level }}
          </a-tag>
        </template>

        <!-- Added copy button for SQL -->
        <template v-if="column.key === 'query'">
          <div class="query-wrapper">
            <a-tooltip :title="record.query" placement="topLeft">
              <div class="query-cell">{{ record.query }}</div>
            </a-tooltip>
          </div>
        </template>

        <!-- Highlight zero affected rows -->
        <template v-if="column.key === 'affected_rows'">
          <span :class="{ 'zero-rows': record.affected_rows === 0 }">
            {{ record.affected_rows }}
          </span>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup>
import { message } from 'ant-design-vue'
import { reactive, computed } from 'vue'

const formData = defineModel('modelValue', {
  type: Object,
  required: true,
})

const tableColumns = [
  {
    title: '指纹',
    dataIndex: 'finger_id',
    key: 'finger_id',
    width: '15%',
    ellipsis: true,
  },
  {
    title: 'SQL语句',
    dataIndex: 'query',
    key: 'query',
    width: '30%',
  },
  {
    title: '操作类型',
    dataIndex: 'type',
    key: 'type',
    width: '120px',
  },
  {
    title: '影响行数',
    dataIndex: 'affected_rows',
    key: 'affected_rows',
    width: '100px',
    align: 'center',
  },
  {
    title: '检查提示',
    dataIndex: 'summary',
    key: 'summary',
    width: '35%',
  },
]

const data = reactive({
  showbTable: false,
  tableData: [],
  summaryStats: {
    error: 0,
    warn: 0,
  },
})

const getLevelColor = (level) => {
  const colorMap = {
    INFO: 'green',
    WARN: 'orange',
    ERROR: 'red',
  }
  return colorMap[level] || 'default'
}

const setRowClass = (record) => {
  return ''
}

const getVisibleSummary = (record) => {
  if (!record.summary) return []
  if (record.expanded || record.summary.length <= 3) {
    return record.summary
  }
  return record.summary.slice(0, 3)
}

const toggleExpand = (record) => {
  record.expanded = !record.expanded
}

const calculateStats = (tableData) => {
  const stats = { error: 0, warn: 0 }
  tableData.forEach((row) => {
    if (row.summary && Array.isArray(row.summary)) {
      row.summary.forEach((item) => {
        const level = item.level?.toLowerCase()
        if (level === 'error') stats.error++
        else if (level === 'warn') stats.warn++
      })
    }
  })
  return stats
}

const render = (res) => {
  if (formData.value.sql_type === 'EXPORT') {
    message.info('数据导出工单无需语法检查，可以直接提交工单')
    data.tableData = []
    data.showbTable = false
    return
  }
  if (res?.status === 0) {
    message.success('✓ 语法检查通过，您可以提交工单了 O(∩_∩)O')
  }
  if (res?.status === 1) {
    message.error('✗ 语法检查未通过，请根据下面提示进行更正 (ㄒoㄒ)')
  }

  const tableData = (res?.data || []).map(item => ({
    ...item,
    expanded: false,
  }))

  data.tableData = tableData
  data.summaryStats = calculateStats(tableData)
  data.showbTable = true
}

defineExpose({
  render,
})
</script>

<style scoped>
/* Added statistics bar styling */
.stats-bar {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 12px 16px;
  background: #fafafa;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  margin-bottom: 12px;
  font-size: 13px;
}

.stats-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.stat-icon {
  font-size: 14px;
}

.error-stat {
  color: #ff4d4f;
  font-weight: 500;
}

.warn-stat {
  color: #faad14;
  font-weight: 500;
}

.info-stat {
  color: #52c41a;
  font-weight: 500;
}

:deep(.row-level-info) {
  color: #52c41a;
}

:deep(.row-level-warn) {
  color: #faad14;
}

:deep(.row-level-error) {
  color: #ff4d4f;
}

.summary-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 4px 0;
}

.summary-item {
  line-height: 1.8;
  padding-left: 12px;
  position: relative;
}

.summary-item::before {
  content: "•";
  position: absolute;
  left: 0;
  font-weight: bold;
}

.summary-level-info {
  color: #52c41a;
  font-weight: 500;
}

.summary-level-info::before {
  color: #52c41a;
}

.summary-level-warn {
  color: #faad14;
  font-weight: 500;
}

.summary-level-warn::before {
  color: #faad14;
}

.summary-level-error {
  color: #ff4d4f;
  font-weight: 500;
}

.summary-level-error::before {
  color: #ff4d4f;
}

/* Added expand link styling */
.expand-link {
  color: #1890ff;
  font-size: 12px;
  cursor: pointer;
  padding-left: 12px;
  user-select: none;
}

.expand-link:hover {
  text-decoration: underline;
}

.empty-summary {
  color: #bfbfbf;
  font-style: italic;
}

/* Enhanced SQL cell with copy button */
.query-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.query-cell {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 12px;
  padding: 4px 0;
  color: rgba(0, 0, 0, 0.85);
}

.copy-btn {
  flex-shrink: 0;
  padding: 0 4px;
  height: 24px;
  opacity: 0.6;
  transition: opacity 0.2s;
}

.copy-btn:hover {
  opacity: 1;
}

/* Highlight zero affected rows */
.zero-rows {
  color: #bfbfbf;
  font-style: italic;
}

:deep(.ant-table-striped .ant-table-tbody > tr:nth-child(2n)) {
  background-color: #fafafa;
}
</style>
