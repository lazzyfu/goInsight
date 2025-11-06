<template>
  <a-card size="small">
    <div class="search-wrapper">
      <a-space>
        <a-switch
          v-model:checked="state.checked"
          checked-children="我的工单"
          un-checked-children="所有工单"
          @change="fetchData"
        />
        <a-input-search
          v-model:value="searchValue"
          placeholder="输入SQL内容"
          style="width: 350px"
          @search="onSearch"
        />
      </a-space>
    </div>
    <div style="margin-top: 12px">
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

// 搜索
const searchValue = ref('')
const onSearch = (value) => {
  searchValue.value = value
  pagination.current = 1
  fetchData()
}

// 我的工单开关
const state = reactive({
  checked: false,
})

const tableColumns = ref([
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
    title: '申请人',
    dataIndex: 'applicant',
    key: 'applicant',
    ellipsis: true,
    scopedSlots: {
      customRender: 'applicant',
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
    title: '实例',
    dataIndex: 'instance',
    key: 'instance',
    scopedSlots: {
      customRender: 'instance',
    },
  },
  {
    title: '库名',
    dataIndex: 'schema',
    key: 'schema',
    scopedSlots: {
      customRender: 'schema',
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
    only_my_orders: state.checked,
    search: searchValue.value,
  }

  const res = await getOrderListApi(params)
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

const getProgressAlias = (progress) => {
  const statusMap = {
    PENDING: { text: '待审批', color: 'default' },
    APPROVED: { text: '已批准', color: 'blue' },
    REJECTED: { text: '已驳回', color: 'red' },
    CLAIMED: { text: '已认领', color: 'cyan' },
    EXECUTING: { text: '执行中', color: 'orange' },
    COMPLETED: { text: '已完成', color: 'green' },
    REVIEWED: { text: '已复核', color: 'green' },
    CLOSED: { text: '已关闭', color: 'gray' },
  }
  return statusMap[progress] || { text: progress, color: 'default' }
}
onMounted(() => {
  fetchData()
})
</script>
