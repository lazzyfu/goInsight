<template>
  <a-card class="order-list-card" title="历史工单">
    <div class="toolbar">
      <div class="toolbar-left">
        <div class="checkbox-container">
          <a-checkbox v-model:checked="uiState.checked" @change="handleMyOrdersChange">我的工单</a-checkbox>
        </div>
        <a-select
          v-model:value="uiData.progress"
          :options="progressOptions"
          allowClear
          class="toolbar-progress-select"
          placeholder="工单进度"
          @change="handleProgressChange"
        />
      </div>
      <a-input-search
        v-model:value="uiData.searchValue"
        placeholder="请输入工单标题、实例、库名等关键词搜索"
        allowClear
        class="toolbar-search"
        @search="handleSearch"
      />
    </div>

    <div v-if="uiState.checked" class="overview-row">
      <div class="overview-item">
        <span class="overview-label">我的工单</span>
        <strong>{{ uiData.myOrderStats.total }}</strong>
      </div>
      <div class="overview-item pending">
        <span class="overview-label">待审批</span>
        <strong>{{ uiData.myOrderStats.pending }}</strong>
      </div>
      <div class="overview-item running">
        <span class="overview-label">执行中</span>
        <strong>{{ uiData.myOrderStats.executing }}</strong>
      </div>
      <div class="overview-item failed">
        <span class="overview-label">失败</span>
        <strong>{{ uiData.myOrderStats.failed }}</strong>
      </div>
    </div>

    <div class="table-wrapper">
      <a-table
        size="middle"
        :columns="uiData.tableColumns"
        :row-key="(record) => record.order_id"
        :data-source="uiData.tableData"
        :pagination="pagination"
        :loading="uiState.loading"
        :row-class-name="() => 'order-row'"
        @change="handleTableChange"
        :scroll="{ x: 1100 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'progress'">
            <a-tag :color="getOrderStatusMeta(record.progress).color" class="progress-tag">
              {{ getOrderStatusMeta(record.progress).text }}
            </a-tag>
          </template>
          <template v-if="column.key === 'title'">
            <router-link class="title-link" :to="{ name: 'orders.detail', params: { order_id: record.order_id } }">
              {{ record.title }}
            </router-link>
          </template>
          <template v-if="column.key === 'instance'">
            {{ record.instance }}
          </template>
          <template v-if="column.key === 'schema'">
            <DatabaseOutlined />
            {{ record.schema }}
          </template>
        </template>
      </a-table>
    </div>
  </a-card>
</template>

<script setup>
import { getOrderListApi } from '@/api/order'
import { ORDER_PROGRESS_OPTIONS, getOrderStatusMeta } from '@/views/orders/shared/orderStatusMeta'
import { DatabaseOutlined } from '@ant-design/icons-vue'
import { useIntervalFn } from '@vueuse/core'
import { onMounted, reactive } from 'vue'

// 定时器, 10秒刷新一次
useIntervalFn(
  () => {
    refreshPageData()
  },
  10000,
  { immediate: false },
)

const progressOptions = ORDER_PROGRESS_OPTIONS

// 状态
const uiState = reactive({
  checked: false,
  loading: false,
})

// 数据
const uiData = reactive({
  searchValue: '',
  progress: null,
  myOrderStats: {
    total: 0,
    pending: 0,
    executing: 0,
    failed: 0,
  },
  tableData: [],
  tableColumns: [
    {
      title: '进度',
      dataIndex: 'progress',
      key: 'progress',
      width: 100,
      fixed: 'left',
    },
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      width: 300,
      fixed: 'left',
      ellipsis: true,
    },
    {
      title: '申请人',
      dataIndex: 'applicant',
      key: 'applicant',
      ellipsis: true,
    },
    {
      title: '组织',
      dataIndex: 'organization',
      key: 'organization',
      ellipsis: true,
    },
    {
      title: '工单环境',
      dataIndex: 'environment',
      key: 'environment',
      ellipsis: true,
    },
    {
      title: '工单类型',
      dataIndex: 'sql_type',
      key: 'sql_type',
    },
    {
      title: '实例',
      dataIndex: 'instance',
      key: 'instance',
    },
    {
      title: '库名',
      dataIndex: 'schema',
      key: 'schema',
    },
    {
      title: '提交时间',
      dataIndex: 'created_at',
      key: 'created_at',
    },
  ],
})

// 搜索
const handleSearch = (value) => {
  uiData.searchValue = value
  pagination.current = 1
  fetchData()
}

const handleMyOrdersChange = async () => {
  pagination.current = 1
  await refreshPageData()
}

const handleProgressChange = () => {
  pagination.current = 1
  fetchData()
}

// 分页
const handleTableChange = (pager) => {
  pagination.current = pager.current
  pagination.pageSize = pager.pageSize
  fetchData()
}

// 翻页
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  pageSizeOptions: ['10', '20', '50', '100'],
  showSizeChanger: true,
})

// 获取列表数据
const fetchData = async () => {
  uiState.loading = true
  try {
    const params = {
      page_size: pagination.pageSize,
      page: pagination.current,
      is_page: true,
      only_my_orders: uiState.checked,
      search: uiData.searchValue,
      progress: uiData.progress,
    }
    const res = await getOrderListApi(params).catch(() => {})
    if (res?.data) {
      pagination.total = res.total
      uiData.tableData = res.data
    } else {
      pagination.total = 0
      uiData.tableData = []
    }
  } finally {
    uiState.loading = false
  }
}

const resetMyOrderStats = () => {
  uiData.myOrderStats.total = 0
  uiData.myOrderStats.pending = 0
  uiData.myOrderStats.executing = 0
  uiData.myOrderStats.failed = 0
}

const fetchMyOrderStats = async () => {
  if (!uiState.checked) {
    resetMyOrderStats()
    return
  }

  const res = await getOrderListApi({
    is_page: false,
    only_my_orders: true,
  }).catch(() => {})

  const rows = Array.isArray(res?.data) ? res.data : []
  uiData.myOrderStats.total = rows.length
  uiData.myOrderStats.pending = rows.filter((item) => item.progress === 'PENDING').length
  uiData.myOrderStats.executing = rows.filter((item) => item.progress === 'EXECUTING').length
  uiData.myOrderStats.failed = rows.filter((item) => item.progress === 'FAILED').length
}

const refreshPageData = async () => {
  await fetchData()
  await fetchMyOrderStats()
}

onMounted(async () => {
  await refreshPageData()
})
</script>

<style scoped>
.order-list-card :deep(.ant-card-body) {
  padding-top: 14px;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.toolbar-search {
  width: 360px;
  max-width: 100%;
}

.toolbar-progress-select {
  width: 140px;
}

.checkbox-container {
  height: 32px;
  padding: 0 14px;
  display: inline-flex;
  align-items: center;
  background: #fafbfd;
  border: 1px solid #dbe5ec;
  border-radius: 8px;
  transition: all 0.2s ease;
  cursor: pointer;
}

.checkbox-container:hover {
  border-color: #b5cad7;
  background: #f3f8fb;
}

.checkbox-container:has(.ant-checkbox-checked) {
  border-color: #8bb0c7;
  background: #eaf3f8;
}

.overview-row {
  margin-top: 12px;
  display: grid;
  grid-template-columns: repeat(4, minmax(120px, 1fr));
  gap: 10px;
}

.overview-item {
  padding: 10px 12px;
  border-radius: 10px;
  border: 1px solid #e5edf3;
  background: #fafcfd;
  display: flex;
  align-items: baseline;
  justify-content: space-between;
}

.overview-item strong {
  font-size: 18px;
  color: #1f2d38;
}

.overview-label {
  font-size: 12px;
  color: #6d7f8b;
}

.overview-item.pending {
  border-color: #f5deb1;
  background: #fff9ed;
}

.overview-item.running {
  border-color: #ffd7b5;
  background: #fff4eb;
}

.overview-item.failed {
  border-color: #ffc7c7;
  background: #fff1f1;
}

.table-wrapper {
  margin-top: 12px;
}

.title-link {
  color: #0d5f8c;
  font-weight: 500;
}

.progress-tag {
  border-radius: 999px;
  padding-inline: 10px;
}

:deep(.order-row:hover > td) {
  background: #f7fafc !important;
}

@media (max-width: 860px) {
  .overview-row {
    grid-template-columns: repeat(2, minmax(120px, 1fr));
  }
}

@media (max-width: 600px) {
  .overview-row {
    grid-template-columns: 1fr;
  }
}
</style>
