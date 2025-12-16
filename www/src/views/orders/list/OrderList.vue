<template>
  <a-card title="历史工单">
    <a-space size="middle" wrap>
      <div class="checkbox-container">
        <a-checkbox v-model:checked="uiState.checked" @change="fetchData">我的工单</a-checkbox>
      </div>
      <!-- 进度筛选 -->
      <a-select v-model:value="uiData.progress" :options="progressOptions" allowClear style="width: 140px"
        placeholder="工单进度" @change="fetchData" />
      <a-input-search v-model:value="uiData.searchValue" placeholder="请输入工单标题、实例、库名等关键词搜索" allowClear
        style="width: 350px" @search="handleSearch" />
    </a-space>

    <div style="margin-top: 12px">
      <a-table size="small" :columns="uiData.tableColumns" :row-key="(record) => record.key"
        :data-source="uiData.tableData" :pagination="pagination" :loading="uiState.loading" @change="handleTableChange"
        :scroll="{ x: 1100 }">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'progress'">
            <template v-if="progressInfo = getProgressAlias(record.progress)">
              <a-tag :color="progressInfo.color">
                {{ progressInfo.text }}
              </a-tag>
            </template>
          </template>
          <template v-if="column.key === 'title'">
            <router-link :to="{ name: 'orders.detail', params: { order_id: record.order_id } }">
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
import { DatabaseOutlined } from '@ant-design/icons-vue'
import { useIntervalFn } from '@vueuse/core'
import { onMounted, reactive } from 'vue'

// 定时器, 10秒刷新一次
const { pause, resume, isActive } = useIntervalFn(
  () => {
    fetchData()
  },
  10000,
  { immediate: true }, // 组件挂载时立刻执行一次
)

const progressOptions = [
  { label: '待审批', value: 'PENDING' },
  { label: '已批准', value: 'APPROVED' },
  { label: '已驳回', value: 'REJECTED' },
  { label: '已认领', value: 'CLAIMED' },
  { label: '执行中', value: 'EXECUTING' },
  { label: '已失败', value: 'FAILED' },
  { label: '已完成', value: 'COMPLETED' },
  { label: '已复核', value: 'REVIEWED' },
  { label: '已撤销', value: 'REVOKED' },
]

// 状态
const uiState = reactive({
  checked: false,
  loading: false,
})

// 数据
const uiData = reactive({
  searchValue: '',
  progress: null,
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
  const params = {
    page_size: pagination.pageSize,
    page: pagination.current,
    is_page: true,
    only_my_orders: uiState.checked,
    search: uiData.searchValue,
    progress: uiData.progress,
  }
  const res = await getOrderListApi(params).catch(() => { })
  if (res) {
    pagination.total = res.total
    uiData.tableData = res.data
  }
  uiState.loading = false
}

// 获取进度别名
const getProgressAlias = (progress) => {
  const statusMap = {
    PENDING: { text: '待审批', color: 'default' },
    APPROVED: { text: '已批准', color: 'blue' },
    REJECTED: { text: '已驳回', color: 'red' },
    CLAIMED: { text: '已认领', color: 'cyan' },
    EXECUTING: { text: '执行中', color: 'orange' },
    FAILED: { text: '已失败', color: 'red' },
    COMPLETED: { text: '已完成', color: 'green' },
    REVIEWED: { text: '已复核', color: 'green' },
    REVOKED: { text: '已撤销', color: 'gray' },
  }
  return statusMap[progress] || { text: progress, color: 'default' }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.checkbox-container {
  height: 32px;
  padding: 0 14px;
  display: inline-flex;
  align-items: center;
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  transition: all 0.2s ease;
  cursor: pointer;
}

.checkbox-container:hover {
  border-color: #bae0ff;
  background: #f0f7ff;
}

/* 选中时的样式 */
.checkbox-container:has(.ant-checkbox-checked) {
  border-color: #91caff;
  background: #e6f4ff;
}
</style>
