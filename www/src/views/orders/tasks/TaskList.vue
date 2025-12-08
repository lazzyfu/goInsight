<template>
  <a-card title="任务列表">
    <div class="search-wrapper">
      <div class="search-row">
        <a-tooltip title="依次执行所有任务">
          <a-button type="primary" @click="executeBatchTasks">
            <PlayCircleOutlined />
            全部执行
          </a-button>
        </a-tooltip>

        <a-input-search
          v-model:value="uiData.searchValue"
          placeholder="请输入SQL文本进行搜索"
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
        :scroll="{ x: 1100 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'progress'">
            <template v-if="progressInfo = getProgressAlias(record.progress)">
              <a-tag :color="progressInfo.color">
                {{ progressInfo.text }}
              </a-tag>
            </template>
          </template>
          <template v-else-if="column.key === 'sql'">
            <a @click="showSqlDetail(record)" title="查看完整SQL">
              <EyeOutlined />
            </a>
            {{ record.sql }}
          </template>
          <template v-if="column.key === 'action'">
            <a-space>
              <a-tooltip title="执行当前任务">
                <a @click="handleEdit(record)">
                  <PlayCircleOutlined />
                  执行
                </a>
              </a-tooltip>

              <a @click="showResult(record)">
                <FileSearchOutlined />
                结果
              </a>
            </a-space>
          </template>
        </template>
      </a-table>
    </div>
  </a-card>
  <!-- 查看SQL -->
  <a-modal
    v-model:open="uiState.open"
    title="SQL语句"
    width="55%"
    :footer="null"
    @cancel="handleCancel"
  >
    <highlightjs language="sql" :code="uiData.sql" />
  </a-modal>
</template>

<script setup>
import { getOrderTasksApi } from '@/api/order'
import { EyeOutlined, FileSearchOutlined, PlayCircleOutlined } from '@ant-design/icons-vue'
import { onMounted, reactive } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const orderId = route.params.order_id

// 状态
const uiState = reactive({
  loading: false,
  open: false,
})

// 数据
const uiData = reactive({
  searchValue: '',
  tableData: [],
  sql: '',
  tableColumns: [
    {
      title: '进度',
      dataIndex: 'progress',
      key: 'progress',
      fixed: 'left',
      width: 100,
    },
    {
      title: '任务ID',
      dataIndex: 'task_id',
      key: 'task_id',
    },
    {
      title: 'SQL文本',
      dataIndex: 'sql',
      key: 'sql',
      ellipsis: true,
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      key: 'updated_at',
      ellipsis: true,
    },
    {
      title: '操作',
      dataIndex: 'action',
      key: 'action',
      fixed: 'right',
      width: 150,
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
    search: uiData.searchValue,
    order_id: orderId,
  }
  const res = await getOrderTasksApi(params).catch(() => {})
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
  // 'PENDING', 'EXECUTING', 'COMPLETED', 'FAILED', 'PAUSED'
  const statusMap = {
    PENDING: { text: '待执行', color: 'default' },
    EXECUTING: { text: '执行中', color: 'orange' },
    PAUSED: { text: '已暂停', color: 'gray' },
    FAILED: { text: '已失败', color: 'red' },
    COMPLETED: { text: '已完成', color: 'green' },
  }
  return statusMap[progress] || { text: progress, color: 'default' }
}

// 查看SQL
const showSqlDetail = (record) => {
  console.log('record: ', record)
  uiState.open = true
  uiData.sql = record.sql
}

// 关闭查看SQL弹窗
const handleCancel = (e) => {
  uiState.open = false
}

// 批量执行
const executeBatchTasks = () => {
  console.log('执行全部任务')
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
