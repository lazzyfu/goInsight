<template>
  <a-modal
    :open="props.open"
    :title="`审批流：${props.flowName}`"
    :width="650"
    :footer="null"
    @cancel="handleCancel"
  >
    <div class="search-wrapper">
      <!-- 搜索 -->
      <a-input-search
        v-model:value="uiData.searchValue"
        placeholder="搜索用户名..."
        style="width: 350px"
        @search="handleSearch"
      />
    </div>

    <!-- 表格 -->
    <div style="margin-top: 12px">
      <a-table
        size="small"
        :columns="uiData.tableColumns"
        :row-key="(record) => record.key"
        :data-source="uiData.tableData"
        :pagination="pagination"
        :loading="uiState.loading"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action'">
            <a-popconfirm
              title="确认删除吗？"
              ok-text="是"
              cancel-text="否"
              @confirm="handleDelete(record)"
            >
              <a> <DeleteOutlined /> 删除 </a>
            </a-popconfirm>
          </template>
        </template>
      </a-table>
    </div>
  </a-modal>
</template>

<script setup>
import { deleteUsersFromApprovalFlowApi, getApprovalFlowUsersApi } from '@/api/admin'
import { DeleteOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { reactive, watch } from 'vue'

// 定义props和emits
const props = defineProps({
  open: Boolean,
  flowId: { type: [Number, String], default: null }, // 审批流ID
  flowName: { type: String, default: '' }, // 审批流名称
})
const emit = defineEmits(['update:open'])

// 状态
const uiState = reactive({
  loading: false,
})

// 数据
const uiData = reactive({
  searchValue: '',
  tableData: [],
  tableColumns: [
    {
      title: '用户',
      dataIndex: 'username',
      key: 'username',
    },
    {
      title: '绑定时间',
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
    approval_id: props.flowId,
  }
  const res = await getApprovalFlowUsersApi(params).catch(() => {})
  if (res) {
    pagination.total = res.total
    uiData.tableData = res.data
  }
  uiState.loading = false
}

watch(
  () => props.open,
  (newVal) => {
    if (newVal && props.flowId) {
      fetchData()
    }
  },
)

// 取消按钮
const handleCancel = () => {
  emit('update:open', false)
}

// 删除
const handleDelete = async (record) => {
  const res = await deleteUsersFromApprovalFlowApi(record).catch(() => {})
  if (res?.code === '0000') {
    message.info('操作成功')
    fetchData()
  }
}
</script>

<style scoped>
.user-list-container {
  max-height: 500px;
  overflow-y: auto;
  padding: 10px;
}

.list-header {
  padding: 8px 16px;
  background: #f0f5ff;
  border: 1px solid #bae7ff;
  border-radius: 4px;
  font-weight: 500;
  color: #1890ff;
  font-size: 14px;
}

.list-header strong {
  font-size: 16px;
  margin: 0 4px;
}

.user-name {
  font-weight: 600;
  color: #333;
}
</style>
