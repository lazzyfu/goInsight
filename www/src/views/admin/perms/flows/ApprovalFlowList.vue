<template>
  <div class="flow-page gi-page-shell">
    <div class="page-hero">
      <div class="hero-content">
        <h2>审批流管理中心</h2>
        <p>统一配置审批流阶段、可领取人和用户绑定关系，确保审批路径清晰且可追溯。</p>
        <div class="hero-stats">
          <div class="stat-item">
            <span>流程总数</span>
            <strong>{{ flowStats.total }}</strong>
          </div>
        </div>
      </div>

      <a-space class="hero-actions" wrap>
        <a-button @click="handleBind">
          <DeploymentUnitOutlined />
          绑定流程到用户
        </a-button>
        <a-button type="primary" @click="handleAdd">
          <PlusOutlined />
          新增审批流
        </a-button>
      </a-space>
    </div>

    <PageTableSection class="table-shell">
      <PageToolbar class="toolbar">
        <a-input-search
          v-model:value="uiData.searchValue"
          placeholder="搜索审批流名称、用户名"
          class="gi-toolbar-search"
          @search="handleSearch"
        />
        <div class="toolbar-tags">
          <a-tag color="processing">共 {{ pagination.total }} 条</a-tag>
          <a-tag v-if="uiData.searchValue">检索词：{{ uiData.searchValue }}</a-tag>
        </div>
      </PageToolbar>

      <a-table
        class="flow-table"
        size="middle"
        :columns="uiData.tableColumns"
        :row-key="(record) => record.id"
        :data-source="uiData.tableData"
        :pagination="pagination"
        :loading="uiState.loading"
        :scroll="{ x: 980 }"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'flow'">
            <a-tag color="blue">{{ record.definition.length }} 个阶段</a-tag>
          </template>

          <template v-if="column.key === 'claim_users'">
            <a-tag color="cyan">{{ record.claim_users.length }} 人</a-tag>
          </template>

          <template v-if="column.key === 'action'">
            <a-space :size="14">
              <a @click="handleViewUsers(record)">
                <EyeOutlined />
                查看用户
              </a>
              <a @click="handleEdit(record)">
                <EditOutlined />
                编辑
              </a>
              <a-popconfirm title="确认删除该审批流吗？" ok-text="是" cancel-text="否" @confirm="handleDelete(record)">
                <a class="text-danger">
                  <DeleteOutlined />
                  删除
                </a>
              </a-popconfirm>
            </a-space>
          </template>
        </template>

        <template #expandedRowRender="{ record }">
          <div class="expanded-content">
            <div class="expanded-title">审批流详细阶段</div>
            <ApprovalFlowStages :definition="record.definition" />
          </div>
        </template>
      </a-table>
    </PageTableSection>
  </div>

  <ApprovalFlowFormModal
    :open="uiState.isModalOpen"
    v-model:formData="formState"
    :title="uiState.isEditMode ? '编辑审批流' : '新增审批流'"
    :user-options="uiData.users"
    @update:open="uiState.isModalOpen = $event"
    @submit="onSubmit"
  />

  <BindToUserFormModal
    :open="uiState.isBindModalOpen"
    :flow-options="uiData.flows"
    :user-options="uiData.unBoundUsers"
    @update:open="uiState.isBindModalOpen = $event"
    @submit="onSubmitBind"
  />

  <FlowBoundUsersDetail
    :open="uiState.isViewUsersOpen"
    :flow-id="uiData.viewApprovalFlowID"
    :flow-name="uiData.viewFlowName"
    @update:open="uiState.isViewUsersOpen = $event"
  />
</template>

<script setup>
import {
  bindUsersToApprovalFlowApi,
  createApprovalFlowsApi,
  deleteApprovalFlowsApi,
  getApprovalFlowsApi,
  getApprovalFlowUnboundUsersApi,
  getUsersApi,
  updateApprovalFlowsApi,
} from '@/api/admin'
import {
  DeleteOutlined,
  DeploymentUnitOutlined,
  EditOutlined,
  EyeOutlined,
  PlusOutlined,
} from '@ant-design/icons-vue'
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { computed, onMounted, reactive, ref } from 'vue'
import PageTableSection from '@/components/patterns/PageTableSection.vue'
import PageToolbar from '@/components/patterns/PageToolbar.vue'
import ApprovalFlowFormModal from './ApprovalFlowFormModal.vue'
import ApprovalFlowStages from './ApprovalFlowStages.vue'
import BindToUserFormModal from './BindToUserFormModal.vue'
import FlowBoundUsersDetail from './FlowBoundUsersDetail.vue'

const uiState = reactive({
  loading: false,
  isEditMode: false,
  isModalOpen: false,
  isBindModalOpen: false,
  isViewUsersOpen: false,
})

const uiData = reactive({
  searchValue: '',
  viewApprovalFlowID: null,
  viewFlowName: '',
  users: [],
  unBoundUsers: [],
  flows: [],
  tableData: [],
  tableColumns: [
    { title: '审批流名称', dataIndex: 'name', key: 'name', width: 180, ellipsis: true },
    { title: '流程阶段数', dataIndex: 'definition', key: 'flow', width: 120 },
    { title: '可领取人', dataIndex: 'claim_users', key: 'claim_users', width: 120 },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 176 },
    { title: '更新时间', dataIndex: 'updated_at', key: 'updated_at', width: 176 },
    { title: '操作', dataIndex: 'action', key: 'action', fixed: 'right', width: 230 },
  ],
})

const defaultForm = {
  name: '',
  claim_users: [],
  definition: [
    {
      stage: 1,
      stage_name: '第一阶段审批',
      approvers: [],
      type: 'AND',
    },
  ],
}
const formState = ref({ ...defaultForm })

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  pageSizeOptions: ['10', '20', '50', '100'],
  showSizeChanger: true,
})

const normalizeArrayField = (value) => {
  if (Array.isArray(value)) return value
  if (typeof value === 'string') {
    try {
      const parsed = JSON.parse(value)
      return Array.isArray(parsed) ? parsed : []
    } catch {
      return []
    }
  }
  return []
}

const flowStats = computed(() => ({
  total: pagination.total,
}))

const getUsers = async () => {
  const res = await getUsersApi().catch(() => {})
  if (res && res.data) {
    uiData.users = res.data.map((u) => ({
      label: `${u.nickname || u.username} (${u.username})`,
      value: u.username,
    }))
  }
}

const getUnBoundUsers = async () => {
  const res = await getApprovalFlowUnboundUsersApi().catch(() => {})
  if (res && res.data) {
    uiData.unBoundUsers = res.data.map((u) => ({
      label: `${u.nickname || u.username} (${u.username})`,
      value: u.username,
    }))
  }
}

const fetchData = async () => {
  uiState.loading = true
  const params = {
    page_size: pagination.pageSize,
    page: pagination.current,
    is_page: true,
    search: uiData.searchValue,
  }

  const res = await getApprovalFlowsApi(params).catch(() => {})
  if (res) {
    pagination.total = res.total
    uiData.tableData = (res.data || []).map((item) => ({
      ...item,
      definition: normalizeArrayField(item.definition),
      claim_users: normalizeArrayField(item.claim_users),
    }))

    uiData.flows = (res.data || []).map((flow) => ({
      label: flow.name,
      value: flow.approval_id,
    }))
  }
  uiState.loading = false
}

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

const handleAdd = () => {
  uiState.isEditMode = false
  formState.value = JSON.parse(JSON.stringify(defaultForm))
  uiState.isModalOpen = true
}

const handleEdit = (record) => {
  uiState.isEditMode = true
  formState.value = {
    ...record,
    claim_users: normalizeArrayField(record.claim_users),
    definition: normalizeArrayField(record.definition).length
      ? normalizeArrayField(record.definition)
      : defaultForm.definition,
  }
  uiState.isModalOpen = true
}

const onSubmit = useThrottleFn(async (data) => {
  const payload = { ...data }
  const res = uiState.isEditMode
    ? await updateApprovalFlowsApi(payload).catch(() => {})
    : await createApprovalFlowsApi(payload).catch(() => {})

  if (res) {
    message.success('操作成功')
    uiState.isModalOpen = false
    fetchData()
  }
})

const handleDelete = useThrottleFn(async (record) => {
  const res = await deleteApprovalFlowsApi(record.id).catch(() => {})
  if (res) {
    message.info('操作成功')
    fetchData()
  }
})

const handleBind = () => {
  getUnBoundUsers()
  uiState.isBindModalOpen = true
}

const onSubmitBind = useThrottleFn(async (data) => {
  const res = await bindUsersToApprovalFlowApi(data).catch(() => {})
  if (res) {
    message.success('操作成功')
    uiState.isBindModalOpen = false
  }
})

const handleViewUsers = (record) => {
  uiData.viewApprovalFlowID = record.approval_id
  uiData.viewFlowName = record.name
  uiState.isViewUsersOpen = true
}

onMounted(() => {
  fetchData()
  getUsers()
})
</script>

<style scoped>
.flow-page {
  --flow-bg-soft: #f3f8ff;
  --flow-bg-card: #ffffff;
  --flow-border: #dce6f7;
  --flow-text-main: #16213c;
  --flow-text-sub: #5f6b8a;
  --flow-accent: #1f6feb;
  --flow-shadow-lg: 0 18px 45px -28px rgba(25, 55, 115, 0.45);
  --flow-shadow-sm: 0 10px 24px -22px rgba(17, 35, 78, 0.5);
  font-family: 'Avenir Next', 'PingFang SC', 'Hiragino Sans GB', 'Noto Sans SC', 'Microsoft YaHei', sans-serif;
  background: #ffffff;
  border: 1px solid #eceff5;
  border-radius: 18px;
  padding: 14px;
}

.page-hero {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  border: 1px solid var(--flow-border);
  border-radius: 14px;
  background:
    radial-gradient(circle at 88% 0%, rgba(22, 84, 194, 0.1), rgba(22, 84, 194, 0) 50%),
    var(--flow-bg-card);
  box-shadow: var(--flow-shadow-sm);
  padding: 16px 18px;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 999px;
  color: #1554c2;
  border: 1px solid rgba(31, 111, 235, 0.28);
  background: rgba(31, 111, 235, 0.08);
  font-size: 12px;
  font-weight: 600;
}

.hero-content h2 {
  margin: 8px 0 6px;
  color: var(--flow-text-main);
  font-size: 24px;
  font-weight: 700;
}

.hero-content p {
  margin: 0;
  color: var(--flow-text-sub);
  font-size: 13px;
  line-height: 1.6;
}

.hero-stats {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.stat-item {
  min-width: 110px;
  border: 1px solid rgba(31, 111, 235, 0.14);
  border-radius: 10px;
  padding: 8px 10px;
  background: var(--flow-bg-soft);
}

.stat-item span {
  display: block;
  color: var(--flow-text-sub);
  font-size: 12px;
}

.stat-item strong {
  display: block;
  margin-top: 5px;
  font-size: 20px;
  line-height: 1;
  color: var(--flow-text-main);
}

.table-shell {
  border: 1px solid var(--flow-border);
  border-radius: 14px;
  background: var(--flow-bg-card);
  box-shadow: var(--flow-shadow-lg);
  padding: 12px 12px 10px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 10px;
}

.toolbar-tags {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.flow-table :deep(.ant-table) {
  border: 1px solid #e1eaf9;
  border-radius: 10px;
  overflow: hidden;
}

.flow-table :deep(.ant-table-thead > tr > th) {
  background: #f3f7ff;
  color: #24355f;
  font-weight: 600;
}

.flow-table :deep(.ant-table-tbody > tr > td) {
  border-bottom-color: #edf2fb;
}

.flow-table :deep(.ant-table-expanded-row > td) {
  background: #fbfdff;
}

.expanded-content {
  padding: 4px 2px;
}

.expanded-title {
  margin-bottom: 8px;
  font-size: 13px;
  font-weight: 600;
  color: #23335c;
}

.text-danger {
  color: #e63c4f;
}

@media (max-width: 1200px) {
  .page-hero {
    flex-direction: column;
  }

  .hero-actions {
    width: 100%;
  }
}

@media (max-width: 1080px) {
  .flow-page {
    padding: 12px;
  }

  .hero-content h2 {
    font-size: 21px;
  }

  .table-shell {
    padding: 10px;
  }

}
</style>
