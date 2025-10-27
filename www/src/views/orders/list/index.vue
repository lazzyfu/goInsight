<template>
  <div>
    <div class="search-wrapper">
      <a-input-search
        v-model:value="searchValue"
        placeholder="输入SQL内容"
        style="width: 350px"
        @search="onSearch"
      />
    </div>
    <div style="margin-top: 14px">
      <a-table
        size="small"
        :columns="tableColumns"
        :row-key="(record) => record.key"
        :data-source="data.tableData"
        @resizeColumn="handleResizeColumn"
        :pagination="pagination"
        :loading="data.loading"
        @change="handleTableChange"
        :scroll="{ x: 1500 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'applicant'">
            <div class="user-cell">
              <div class="user-avatar" :style="{ backgroundColor: getUserColor(record.applicant) }">
                {{ record.applicant.charAt(0).toUpperCase() }}
              </div>
              <span class="user-name">{{ record.applicant }}</span>
            </div>
          </template>
          <template v-if="column.key === 'title'">
            <router-link :to="{ name: 'orders.detail', params: { order_id: record.order_id } }">
              {{ record.title }}
            </router-link>
          </template>
          <template v-if="column.key === 'schema'">
            <span class="schema-tag">
              <DatabaseOutlined />
              {{ record.schema }}
            </span>
          </template>
          <template v-else-if="column.key === 'created_at'">
            <div class="time-cell">
              <div class="time-main">{{ formatDate(record.created_at) }}</div>
              <div class="time-sub">{{ formatTime(record.created_at) }}</div>
            </div>
          </template>
        </template>
      </a-table>
    </div>
  </div>
</template>

<script setup>
import { getOrdersHistoryApi } from '@/api/order'
import { onMounted, reactive, ref } from 'vue'

import { DatabaseOutlined } from '@ant-design/icons-vue'

// 搜索
const searchValue = ref('')

const onSearch = (value) => {
  searchValue.value = value
  pagination.current = 1
  fetchData()
}

const tableColumns = ref([
  {
    title: '申请人',
    dataIndex: 'applicant',
    key: 'applicant',
    ellipsis: true,
    fixed: 'left',
    scopedSlots: {
      customRender: 'applicant',
    },
  },
  {
    title: '进度',
    dataIndex: 'progress',
    key: 'progress',
    width: 100,
    fixed: 'left',
    scopedSlots: {
      customRender: 'progress',
    },
  },
  {
    title: '标题',
    dataIndex: 'title',
    key: 'title',
    width: 300,
    fixed: 'left',
    ellipsis: true,
    scopedSlots: {
      customRender: 'title',
    },
  },
  {
    title: '组织',
    dataIndex: 'organization',
    key: 'organization',
    ellipsis: true,
    scopedSlots: {
      customRender: 'organization',
    },
  },
  {
    title: '工单环境',
    dataIndex: 'environment',
    key: 'environment',
    ellipsis: true,
    scopedSlots: {
      customRender: 'environment',
    },
  },
  {
    title: '工单类型',
    dataIndex: 'sql_type',
    key: 'sql_type',
    scopedSlots: {
      customRender: 'sql_type',
    },
  },
  {
    title: '实例/库',
    dataIndex: 'instance',
    key: 'instance',
    scopedSlots: {
      customRender: 'instance',
    },
  },
  {
    title: '提交时间',
    dataIndex: 'created_at',
    key: 'created_at',
    scopedSlots: {
      customRender: 'created_at',
    },
  },
])

const data = reactive({
  tableData: [],
  loading: false,
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  pageSizeOptions: ['10', '20', '50', '100'],
  showSizeChanger: true,
})

const fetchData = async () => {
  data.loading = true
  const params = {
    page_size: pagination.pageSize,
    page: pagination.current,
    is_page: true,
    search: searchValue.value,
  }

  const res = await getOrdersHistoryApi(params)
  console.log('res: ', res)

  if (res) {
    pagination.total = res.total
    data.tableData = res.data
  }
  data.loading = false
}

const handleTableChange = (pager) => {
  pagination.current = pager.current
  pagination.pageSize = pager.pageSize
  fetchData()
}

function handleResizeColumn(w, col) {
  col.width = w
}

const formatDate = (dateStr) => {
  return dateStr.split(' ')[0]
}

const formatTime = (dateStr) => {
  return dateStr.split(' ')[1]
}

const getUserColor = (username) => {
  const colors = ['#3b82f6', '#8b5cf6', '#f59e0b', '#10b981', '#ef4444', '#6366f1']
  let hash = 0
  for (let i = 0; i < username.length; i++) {
    hash = username.charCodeAt(i) + ((hash << 5) - hash)
  }
  return colors[Math.abs(hash) % colors.length]
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

.user-name {
  font-weight: 500;
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
</style>
