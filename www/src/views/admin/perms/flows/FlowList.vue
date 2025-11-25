<template>
  <a-card title="审批流管理" class="flow-manager-card">
    <template #extra>
      <a-button type="primary" @click="handleAdd"> <PlusOutlined /> 新增审批流 </a-button>
    </template>

    <div class="search-wrapper">
      <a-input-search
        v-model:value="searchValue"
        placeholder="搜索审批流名称..."
        style="width: 350px"
        @search="handleSearch"
      />
    </div>

    <div style="margin-top: 16px">
      <a-table
        size="middle"
        :columns="tableColumns"
        :row-key="(record) => record.id"
        :data-source="tableData"
        :pagination="pagination"
        :loading="state.loading"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'flow'">
            <a-tag color="blue">{{ record.definition.length }} 个阶段</a-tag>
          </template>

          <template v-if="column.key === 'action'">
            <a-space>
              <a @click="handleEdit(record)"> <EditOutlined /> 编辑 </a>
              <a-popconfirm
                title="确认删除该审批流吗？"
                ok-text="是"
                cancel-text="否"
                @confirm="handleDelete(record)"
              >
                <a class="text-danger"><DeleteOutlined /> 删除</a>
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
    :open="state.isModalOpen"
    v-model:formData="formState"
    :title="state.isEditMode ? '编辑审批流' : '新增审批流'"
    :user-options="state.users"
    @update:open="state.isModalOpen = $event"
    @submit="onSubmit"
  />
</template>

<script setup>
import {
  createApprovalFlowsApi,
  deleteApprovalFlowsApi,
  getApprovalFlowsApi,
  getUsersApi,
  updateApprovalFlowsApi,
} from '@/api/admin'
import { DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref } from 'vue'
import ApprovalFlowFormModal from './ApprovalFlowFormModal.vue'
import ApprovalFlowStages from './ApprovalFlowStages.vue'

const state = reactive({
  loading: false,
  isEditMode: false,
  isModalOpen: false,
  users: [],
})
const searchValue = ref('')

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

// 表格列定义 (保持不变)
const tableColumns = [
  { title: '审批流名称', dataIndex: 'name', key: 'name', width: 150 },
  { title: '流程阶段数', dataIndex: 'definition', key: 'flow', width: 120 },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at' },
  { title: '更新时间', dataIndex: 'updated_at', key: 'updated_at' },
  { title: '操作', dataIndex: 'action', key: 'action', fixed: 'right', width: 120 },
]
const tableData = ref([])

// 搜索
const handleSearch = (value) => {
  searchValue.value = value
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

// 获取用户列表的方法
const getUsers = async () => {
  const res = await getUsersApi().catch(() => {})
  if (res && res.data) {
    // 假设后端返回的用户数据是 [{ username: 'zhangsan', nickname: '张三' }, ...]
    // 转换为 Ant Design Select Options 格式: [{ label: '张三 (zhangsan)', value: 'zhangsan' }]
    state.users = res.data.map((u) => ({
      label: `${u.nickname || u.username} (${u.username})`,
      value: u.username,
    }))
  }
}

// 获取审批流列表的方法
const fetchData = async () => {
  state.loading = true
  const params = {
    page_size: pagination.pageSize,
    page: pagination.current,
    is_page: true,
    search: searchValue.value,
  }
  const res = await getApprovalFlowsApi(params)
  if (res) {
    pagination.total = res.total
    tableData.value = res.data.map((item) => ({
      ...item,
      definition: Array.isArray(item.definition)
        ? item.definition
        : item.definition
          ? JSON.parse(item.definition)
          : [],
    }))
  }
  state.loading = false
}

const handleAdd = () => {
  state.isEditMode = false
  formState.value = JSON.parse(JSON.stringify(defaultForm))
  state.isModalOpen = true
}

const handleEdit = (record) => {
  state.isEditMode = true
  const definition = Array.isArray(record.definition)
    ? record.definition
    : typeof record.definition === 'string'
      ? JSON.parse(record.definition)
      : []

  formState.value = {
    ...record,
    definition: definition.length > 0 ? definition : defaultForm.definition,
  }
  state.isModalOpen = true
}

const onSubmit = async (data) => {
  const payload = { ...data }
  const res = state.isEditMode
    ? await updateApprovalFlowsApi(payload)
    : await createApprovalFlowsApi(payload)

  if (res?.code === '0000') {
    message.success('操作成功')
    state.isModalOpen = false
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

// 生命周期
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
