<template>
  <div>
    <div class="search-wrapper">
      <a-input-search v-model:value="uiData.searchValue" placeholder="输入要查询的SQL内容" style="width: 350px"
        @search="handleSearch" />
    </div>
    <div style="margin-top: 14px">
      <a-table size="small" :columns="uiData.tableColumns" :row-key="(record) => record.key"
        :data-source="uiData.tableData" :pagination="pagination" :loading="uiState.loading" @change="handleTableChange"
        :scroll="{ x: 1100 }">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'username'">
            {{ record.username }}
          </template>
          <template v-if="column.key === 'schema'">
            <span class="schema-tag">
              <DatabaseOutlined />
              {{ record.schema }}
            </span>
          </template>
          <template v-else-if="column.key === 'sqltext'">
            <a @click="showSqlDetail(record)" title="查看完整SQL">
              <EyeOutlined />
            </a>
            {{ record.sqltext }}
          </template>
          <template v-else-if="column.key === 'duration'">
            <span class="duration-tag" :class="getDurationClass(record.duration)">
              <ClockCircleOutlined />
              {{ record.duration }} ms
            </span>
          </template>
          <template v-else-if="column.key === 'error_msg'">
            <span class="status-tag" :class="record.error_msg ? 'error' : 'success'">
              <CheckSquareOutlined v-if="!record.error_msg" />
              <CloseCircleOutlined v-else />
              <a-tooltip>
                <template v-if="record.error_msg" #title>{{ record.error_msg }}</template>
                {{ record.error_msg ? '执行失败' : '执行成功' }}
              </a-tooltip>
            </span>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space wrap>
              <a-tooltip title="拷贝SQL语句">
                <a-button type="link" block shape="circle" :icon="h(CopyOutlined)"
                  @click="copyRecord(record.sqltext)" />
              </a-tooltip>
            </a-space>
          </template>
        </template>
      </a-table>
    </div>

    <a-modal v-model:open="uiState.open" title="SQL详情" width="55%" :footer="null" @ok="handleOk">
      <a-card title="原始SQL" size="small">
        <highlightjs language="sql" :code="uiData.sqltext" />
      </a-card>
      <a-card title="重写SQL" style="margin-top: 10px" size="small">
        <highlightjs language="sql" :code="uiData.rewrite_sqltext" />
      </a-card>
      <a-card title="参数" style="margin-top: 10px" size="small">
        <highlightjs language="sql" :code="uiData.params" />
      </a-card>
    </a-modal>
  </div>
</template>
<script setup>
import { GetHistoryApi } from '@/api/das'
import {
  CheckSquareOutlined,
  ClockCircleOutlined,
  CloseCircleOutlined,
  CopyOutlined,
  DatabaseOutlined,
  EyeOutlined,
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { h, onMounted, reactive } from 'vue'
import useClipboard from 'vue-clipboard3'

const { toClipboard } = useClipboard()

// 状态
const uiState = reactive({
  open: false,
  loading: false,
})

// 数据
const uiData = reactive({
  searchValue: '',
  sqltext: '',
  rewrite_sqltext: '',
  params: '',
  tableData: [],
  tableColumns: [
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
      fixed: 'left',
    },
    {
      title: '实例ID',
      dataIndex: 'instance_id',
      key: 'instance_id',
      ellipsis: true,
    },
    {
      title: '库名',
      dataIndex: 'schema',
      key: 'schema',
    },
    {
      title: '表名',
      dataIndex: 'tables',
      key: 'tables',
      ellipsis: true,
    },
    {
      title: 'SQL语句',
      dataIndex: 'sqltext',
      key: 'sqltext',
      ellipsis: true,
      width: 350,
    },
    {
      title: '返回行数',
      dataIndex: 'return_rows',
      key: 'return_rows',
    },
    {
      title: '耗时',
      dataIndex: 'duration',
      key: 'duration',
    },
    {
      title: '执行状态',
      dataIndex: 'error_msg',
      key: 'error_msg',
    },
    {
      title: '执行时间',
      dataIndex: 'created_at',
      key: 'created_at',
    },
    {
      title: '操作',
      dataIndex: 'action',
      key: 'action',
      fixed: 'right',
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
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  pageSizeOptions: ['10', '20', '50', '100'],
  showSizeChanger: true,
})

// 翻页
const handleTableChange = (pager) => {
  pagination.current = pager.current
  pagination.pageSize = pager.pageSize
  fetchData()
}

// 获取列表数据
const fetchData = async () => {
  uiState.loading = true
  const params = {
    page_size: pagination.pageSize,
    page: pagination.current,
    is_page: true,
    search: uiData.searchValue,
  }
  const res = await GetHistoryApi(params).catch(() => { })
  if (res) {
    pagination.total = res.total
    uiData.tableData = res.data
  }
  uiState.loading = false
}

// 显示SQL详情
const showSqlDetail = (record) => {
  uiState.open = true
  uiData.sqltext = record.sqltext
  uiData.rewrite_sqltext = record.rewrite_sqltext
  uiData.params = JSON.stringify(record.params)
}

// 确认按钮处理函数
const handleOk = (e) => {
  uiState.open = false
}

// 获取耗时标签样式
const getDurationClass = (duration) => {
  if (duration > 1000) return 'slow'
  if (duration > 100) return 'medium'
  return 'fast'
}

// 复制到剪贴板
const copyRecord = async (value) => {
  try {
    await toClipboard(value)
    message.success('已拷贝到剪贴板')
  } catch (e) {
    message.error(e)
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.user-cell {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.user-avatar {
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 0.875rem;
}

.schema-tag {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  background: #f0fdf4;
  color: #166534;
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  font-weight: 500;
}

.duration-tag {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  font-weight: 500;
}

.duration-tag.fast {
  background: #f0fdf4;
  color: #166534;
}

.duration-tag.medium {
  background: #fffbeb;
  color: #92400e;
}

.duration-tag.slow {
  background: #fef2f2;
  color: #991b1b;
}

.time-cell {
  font-size: 0.875rem;
}

.time-main {
  font-weight: 500;
  color: #1e293b;
}

.time-sub {
  color: #6b7280;
  font-size: 0.75rem;
}

.status-tag {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  font-weight: 500;
}

.status-tag.success {
  background: #f0fdf4;
  color: #166534;
}

.status-tag.error {
  background: #fef2f2;
  color: #991b1b;
}

.error-message {
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 0.375rem;
  padding: 0.75rem;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 0.875rem;
  color: #991b1b;
}
</style>
