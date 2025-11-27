<template>
  <a-card title="审批流管理" class="flow-manager-card">
    <template #extra>
      <a-space>
        <a-button @click="handleBind"> <DeploymentUnitOutlined /> 绑定流程到用户 </a-button>
        <a-button type="primary" @click="handleAdd"> <PlusOutlined /> 新增审批流 </a-button>
      </a-space>
    </template>

    <div class="search-wrapper">
      <a-input-search
        v-model:value="uiData.searchValue"
        placeholder="搜索审批流名称..."
        style="width: 350px"
        @search="handleSearch"
      />
    </div>

    <div style="margin-top: 16px">
      <a-table
        size="middle"
        :columns="uiData.tableColumns"
        :row-key="(record) => record.id"
        :data-source="uiData.tableData"
        :pagination="pagination"
        :loading="uiState.loading"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'flow'">
            <a-tag color="blue">{{ record.definition.length }} 个阶段</a-tag>
          </template>

          <template v-if="column.key === 'action'">
            <a-space>
              <a @click="handleViewUsers(record)"> <EyeOutlined /> 查看用户 </a>

              <a @click="handleEdit(record)"> <EditOutlined /> 编辑 </a>
              <a-popconfirm
                title="确认删除该审批流吗？"
                ok-text="是"
                cancel-text="否"
                @confirm="handleDelete(record)"
              >
                <a class="text-danger"> <DeleteOutlined /> 删除 </a>
              </a-popconfirm>
            </a-space>
          </template>
        </template>

        <template #expandedRowRender="{ record }">
          <div style="padding: 10px 0 10px 50px">
            <h4>审批流详细阶段</h4>
            <ApprovalFlowStages :definition="record.definition" />
          </div>
        </template>
      </a-table>
    </div>
  </a-card>

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
    :user-options="uiData.users"
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
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref } from 'vue'
import ApprovalFlowFormModal from './ApprovalFlowFormModal.vue'
import ApprovalFlowStages from './ApprovalFlowStages.vue'
import BindToUserFormModal from './BindToUserFormModal.vue'
import FlowBoundUsersDetail from './FlowBoundUsersDetail.vue' // 导入新组件

const uiState = reactive({
  loading: false,
  isEditMode: false,
  isModalOpen: false,
  isBindModalOpen: false,
  isViewUsersOpen: false,
})

// 数据
const uiData = reactive({
  searchValue: '',
  viewApprovalFlowID: null,
  viewFlowName: '',
  users: [],
  flows: [],
  tableData: [],
  tableColumns: [
    { title: '审批流名称', dataIndex: 'name', key: 'name', width: 150 },
    { title: '流程阶段数', dataIndex: 'definition', key: 'flow', width: 120 },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at' },
    { title: '更新时间', dataIndex: 'updated_at', key: 'updated_at' },
    { title: '操作', dataIndex: 'action', key: 'action', fixed: 'right', width: 250 },
  ],
})

// form表单
const defaultForm = {
  name: '',
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

// 获取用户
const getUsers = async () => {
  const res = await getUsersApi().catch(() => {})
  if (res && res.data) {
    uiData.users = res.data.map((u) => ({
      label: `${u.nickname || u.username} (${u.username})`,
      value: u.username,
    }))
  }
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
  const res = await getApprovalFlowsApi(params)
  if (res) {
    pagination.total = res.total
    uiData.tableData = res.data.map((item) => ({
      ...item,
      definition: Array.isArray(item.definition)
        ? item.definition
        : item.definition
          ? JSON.parse(item.definition)
          : [],
    }))

    uiData.flows = res.data.map((flow) => ({
      label: flow.name,
      value: flow.approval_id,
    }))
  }
  uiState.loading = false
}

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

// 新增
const handleAdd = () => {
  uiState.isEditMode = false
  formState.value = JSON.parse(JSON.stringify(defaultForm))
  uiState.isModalOpen = true
}

// 编辑
const handleEdit = (record) => {
  uiState.isEditMode = true
  const definition = Array.isArray(record.definition)
    ? record.definition
    : typeof record.definition === 'string'
      ? JSON.parse(record.definition)
      : []

  formState.value = {
    ...record,
    definition: definition.length > 0 ? definition : defaultForm.definition,
  }
  uiState.isModalOpen = true
}

// 提交（
const onSubmit = async (data) => {
  const payload = { ...data }
  const res = uiState.isEditMode
    ? await updateApprovalFlowsApi(payload).catch(() => {})
    : await createApprovalFlowsApi(payload).catch(() => {})

  if (res?.code === '0000') {
    message.success('操作成功')
    uiState.isModalOpen = false
    fetchData()
  }
}

const handleDelete = async (record) => {
  const res = await deleteApprovalFlowsApi(record.id).catch(() => {})
  if (res?.code === '0000') {
    message.info('操作成功')
    fetchData()
  }
}

// 打开绑定审批流模态框
const handleBind = () => {
  if (uiData.flows.length === 0) {
    message.warning('当前没有可用的审批流，请先创建。')
    return
  }
  if (uiData.users.length === 0) {
    message.warning('用户列表未加载，请稍候重试。')
    return
  }
  uiState.isBindModalOpen = true
}

// 提交绑定审批流到用户的请求
const onSubmitBind = async (data) => {
  const res = await bindUsersToApprovalFlowApi(data).catch(() => {})
  if (res?.code === '0000') {
    message.success('操作成功')
    uiState.isBindModalOpen = false
  }
}

// 查看绑定了该审批流的用户列表
const handleViewUsers = (record) => {
  uiData.viewApprovalFlowID = record.approval_id
  uiData.viewFlowName = record.name
  uiState.isViewUsersOpen = true
}

onMounted(() => {
  fetchData()
  getUsers() // 获取用户列表
})
</script>

<style scoped>
.flow-manager-card {
  min-height: 80vh;
}

.text-danger {
  color: #ff4d4f;
}
</style>
