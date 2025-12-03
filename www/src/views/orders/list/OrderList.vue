<template>
  <a-card title="历史工单">
    <div class="search-wrapper">
      <div class="search-row">
        <a-switch
          v-model:checked="uiState.checked"
          checked-children="我的工单"
          un-checked-children="所有工单"
          @change="fetchData"
        />
        <a-input-search
          v-model:value="uiData.searchValue"
          placeholder="请输入工单标题、实例、库名等关键词搜索"
          style="width: 350px"
          @search="handleSearch"
        />
      </div>
    </div>

    <div style="margin-top: 12px">
      <a-table
        size="small"
        :columns="uiData.tableColumns"
        :row-key="(record) => record.key"
        :data-source="uiData.tableData"
        :pagination="pagination"
        :loading="uiState.loading"
        @change="handleTableChange"
        :scroll="{ x: 1500 }"
      >
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
          <template v-else-if="column.key === 'created_at'">
            <div class="time-cell">
              <div class="time-main"><FieldTimeOutlined /> {{ formatDate(record.created_at) }}</div>
              <div class="time-sub">{{ formatTime(record.created_at) }}</div>
            </div>
          </template>
        </template>
      </a-table>
    </div>
  </a-card>
</template>

<script setup>
import { getOrderListApi } from '@/api/order'
import { DatabaseOutlined, FieldTimeOutlined } from '@ant-design/icons-vue'
import { onMounted, reactive, ref } from 'vue'

// 状态
const uiState = reactive({
  checked: false,
  loading: false,
})

// 数据
const uiData = reactive({
  searchValue: '',
  tableData: [],
  tableColumns: ref([
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
  ]),
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
  }
  const res = await getOrderListApi(params).catch(() => {})
  if (res) {
    pagination.total = res.total
    uiData.tableData = res.data
  }
  uiState.loading = false
}

// 格式化日期时间
const formatDate = (dateStr) => {
  return dateStr.split(' ')[0]
}

// 格式化时间
const formatTime = (dateStr) => {
  return dateStr.split(' ')[1]
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
.search-row {
  display: flex;
  align-items: center;
  gap: 12px;
}
</style>
