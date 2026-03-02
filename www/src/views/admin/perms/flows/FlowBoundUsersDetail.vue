<template>
  <a-modal :open="props.open" :footer="null" :width="760" centered class="flow-users-modal" @cancel="handleCancel">
    <div class="modal-shell">
      <div class="modal-head">
        <span class="head-badge">用户明细</span>
        <h3>{{ props.flowName ? `审批流：${props.flowName}` : '审批流用户列表' }}</h3>
        <p>查看当前审批流已绑定的用户，并可直接解除绑定关系。</p>
      </div>

      <div class="toolbar">
        <a-input-search
          v-model:value="uiData.searchValue"
          placeholder="搜索用户名、昵称、手机号、邮箱..."
          style="width: 350px"
          @search="handleSearch"
        />
        <a-tag color="processing">共 {{ pagination.total }} 人</a-tag>
      </div>

      <a-table
        class="users-table"
        size="middle"
        :columns="uiData.tableColumns"
        :row-key="(record) => record.key"
        :data-source="uiData.tableData"
        :pagination="pagination"
        :loading="uiState.loading"
        :scroll="{ x: 640 }"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action'">
            <a-popconfirm title="确认删除吗？" ok-text="是" cancel-text="否" @confirm="handleDelete(record)">
              <a class="text-danger">
                <DeleteOutlined />
                删除
              </a>
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
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { reactive, watch } from 'vue'

const props = defineProps({
  open: Boolean,
  flowId: { type: [Number, String], default: null },
  flowName: { type: String, default: '' },
})
const emit = defineEmits(['update:open'])

const uiState = reactive({
  loading: false,
})

const uiData = reactive({
  searchValue: '',
  tableData: [],
  tableColumns: [
    {
      title: '用户',
      dataIndex: 'username',
      key: 'username',
      width: 180,
    },
    {
      title: '绑定时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 220,
    },
    {
      title: '操作',
      dataIndex: 'action',
      key: 'action',
      width: 90,
      fixed: 'right',
    },
  ],
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  pageSizeOptions: ['10', '20', '50', '100'],
  showSizeChanger: true,
})

const handleSearch = (value) => {
  uiData.searchValue = value
  pagination.current = 1
  fetchData()
}

const handleTableChange = (pager) => {
  pagination.current = pager.current
  pagination.pageSize = pager.pageSize
  fetchData()
}

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
      pagination.current = 1
      uiData.searchValue = ''
      fetchData()
    }
  },
)

const handleCancel = () => {
  emit('update:open', false)
}

const handleDelete = useThrottleFn(async (record) => {
  const res = await deleteUsersFromApprovalFlowApi(record).catch(() => {})
  if (res) {
    message.info('操作成功')
    fetchData()
  }
})
</script>

<style scoped>
.modal-shell {
  padding: 8px 2px 4px;
}

.modal-head h3 {
  margin: 10px 0 6px;
  font-size: 22px;
  color: #16213c;
}

.modal-head p {
  margin: 0;
  color: #5f6b8a;
  line-height: 1.7;
}

.head-badge {
  display: inline-flex;
  align-items: center;
  font-size: 12px;
  font-weight: 700;
  color: #1f6feb;
  border-radius: 999px;
  border: 1px solid rgba(31, 111, 235, 0.28);
  background: rgba(31, 111, 235, 0.08);
  padding: 3px 10px;
}

.toolbar {
  margin: 14px 0 10px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.users-table :deep(.ant-table) {
  border: 1px solid #e1eaf9;
  border-radius: 10px;
  overflow: hidden;
}

.users-table :deep(.ant-table-thead > tr > th) {
  background: #f3f7ff;
  color: #24355f;
  font-weight: 600;
}

.users-table :deep(.ant-table-tbody > tr > td) {
  border-bottom-color: #edf2fb;
}

.text-danger {
  color: #e63c4f;
}

:deep(.flow-users-modal .ant-modal-content) {
  border-radius: 18px;
  padding: 20px 22px;
  background:
    radial-gradient(circle at right top, rgba(31, 111, 235, 0.08), rgba(31, 111, 235, 0) 52%),
    #ffffff;
}

</style>
