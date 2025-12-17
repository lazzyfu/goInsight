<template>
  <a-card title="任务列表">
    <a-space size="middle" wrap>
      <a-button type="primary" @click="executeBatchTasks">
        <PlayCircleOutlined />
        全部执行
      </a-button>

      <!-- 进度筛选 -->
      <a-select v-model:value="uiData.progress" :options="progressOptions" allowClear style="width: 140px"
        placeholder="任务进度" @change="fetchData" />

      <!-- 搜索框 -->
      <a-input-search v-model:value="uiData.searchValue" placeholder="请输入 SQL 文本" style="width: 260px" allowClear
        @search="handleSearch" />
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
          <template v-else-if="column.key === 'sql'">
            <a @click="showSqlDetail(record)" title="查看完整SQL">
              <EyeOutlined />
            </a>
            {{ record.sql }}
          </template>
          <template v-if="column.key === 'action'">
            <a-space>
              <a-tooltip title="执行当前任务">
                <a @click="executeTask(record)">
                  <PlayCircleOutlined />
                  执行
                </a>
              </a-tooltip>

              <a @click="showTaskResult(record)">
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
  <a-modal v-model:open="uiState.open" title="SQL语句" width="55%" :footer="null" @cancel="handleCancel">
    <highlightjs language="sql" :code="uiData.sql" />
  </a-modal>

  <TaskResult :open="uiState.taskResultOpen" v-model:modelValue="taskResultData"
    @update:open="uiState.taskResultOpen = $event">
  </TaskResult>
  <TaskStream />
</template>

<script setup>
import { executeTaskApi, executebatchTasksApi, getOrderTasksApi } from '@/api/order'
import {
  EyeOutlined,
  FileSearchOutlined,
  PlayCircleOutlined
} from '@ant-design/icons-vue'
import { useIntervalFn, useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import TaskResult from './TaskResult.vue'
import TaskStream from './TaskStream.vue'

const route = useRoute()
const orderID = route.params.order_id
const taskResultData = ref({})
const progressOptions = [
  { value: 'PENDING', label: '待执行' },
  { value: 'EXECUTING', label: '执行中' },
  { value: 'COMPLETED', label: '已完成' },
  { value: 'FAILED', label: '已失败' },
]

// 定时器, 10秒刷新一次
const { pause, resume, isActive } = useIntervalFn(
  () => {
    fetchData()
  },
  10000,
  { immediate: true }  // 组件挂载时立刻执行一次
)

// 状态
const uiState = reactive({
  loading: false,
  open: false,
  taskResultOpen: false,
})

// 数据
const uiData = reactive({
  searchValue: '',
  progress: null,
  sql: '',
  tableData: [],
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
    progress: uiData.progress,
    order_id: orderID,
  }
  const res = await getOrderTasksApi(params).catch(() => { })
  if (res) {
    pagination.total = res.total
    uiData.tableData = res.data
  }
  uiState.loading = false
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
  uiState.open = true
  uiData.sql = record.sql
}

// 关闭查看SQL弹窗
const handleCancel = (e) => {
  uiState.open = false
}

// 执行单个任务
const executeTask = useThrottleFn(async (record) => {
  message.success(`开始执行任务: ${record.task_id}`)
  const res = await executeTaskApi({ order_id: orderID, task_id: record.task_id }).catch(() => { })
  if (res) {
    fetchData()
  }
})

// 批量执行
const executeBatchTasks = useThrottleFn(async () => {
  message.success('开始批量执行任务')
  const res = await executebatchTasksApi({ order_id: orderID }).catch(() => { })
  if (res) {
    fetchData()
  }
})

// 查看任务结果
const showTaskResult = (record) => {
  taskResultData.value = record || {}
  uiState.taskResultOpen = true
}

onMounted(() => {
  fetchData()
})

</script>
