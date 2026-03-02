<template>
  <div class="org-users" :class="{ 'compact-mode': props.compactMode }">
    <div class="users-header">
      <div>
        <span class="users-caption">成员管理</span>
        <h3>{{ props.nodeName || '当前组织' }}</h3>
        <p>管理组织内用户绑定关系，并分配角色权限。</p>
      </div>
      <a-button type="primary" @click="handleAdd">
        <PlusOutlined />
        绑定用户
      </a-button>
    </div>

    <div class="toolbar">
      <a-input-search
        v-model:value="uiData.searchValue"
        placeholder="搜索用户名、昵称、手机号、邮箱..."
        style="width: 350px; max-width: 100%"
        @search="handleSearch"
      />
      <div class="toolbar-tags">
        <a-tag color="processing">共 {{ pagination.total }} 人</a-tag>
        <a-tag v-if="uiData.searchValue" color="default">检索词：{{ uiData.searchValue }}</a-tag>
      </div>
    </div>

    <a-table
      class="users-table"
      size="middle"
      :columns="uiData.tableColumns"
      :row-key="(record) => `${record.uid}-${record.organization_key}`"
      :data-source="uiData.tableData"
      :pagination="pagination"
      :loading="uiState.loading"
      :scroll="{ x: 760 }"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <a-tooltip title="移除当前用户">
            <a-popconfirm title="确认移除吗？" ok-text="是" cancel-text="否" @confirm="handleDelete(record)">
              <a>
                <DeleteOutlined />
                移除
              </a>
            </a-popconfirm>
          </a-tooltip>
        </template>
      </template>
    </a-table>
  </div>

  <BindOrgUsers
    :open="uiState.isModalOpen"
    :nodeKey="props.nodeKey"
    :nodeName="props.nodeName"
    :users="uiData.users"
    :roles="uiData.roles"
    v-model:modelValue="formState"
    @update:open="uiState.isModalOpen = $event"
    @submit="onSubmit"
  />
</template>

<script setup>
import {
  bindOrganizationsUsersApi,
  deleteOrganizationsUsersApi,
  getOrganizationsUsersApi,
  getRolesApi,
  getUsersApi,
} from '@/api/admin'
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { reactive, ref, watch } from 'vue'
import BindOrgUsers from './BindOrgUsers.vue'

// props
const props = defineProps({
  nodeKey: {
    type: String,
    required: true,
  },
  nodeName: {
    type: String,
    default: '',
  },
  compactMode: {
    type: Boolean,
    default: false,
  },
})

// form表单
const defaultForm = {
  users: [],
  roles: null,
}
const formState = ref({ ...defaultForm })

// 状态
const uiState = reactive({
  loading: false,
  isModalOpen: false,
})

// 数据
const uiData = reactive({
  searchValue: '',
  tableData: [],
  users: [],
  roles: [],
  tableColumns: [
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
    },
    {
      title: '昵称',
      dataIndex: 'nick_name',
      key: 'nick_name',
    },
    {
      title: '角色',
      dataIndex: 'role_name',
      key: 'role_name',
    },
    {
      title: '组织',
      dataIndex: 'organization_name',
      key: 'organization_name',
    },
    {
      title: '操作',
      dataIndex: 'action',
      key: 'action',
      width: 110,
      fixed: 'right',
    },
  ],
})

// 分页
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  pageSizeOptions: ['10', '20', '50', '100'],
  showSizeChanger: true,
})

// 获取列表数据
const fetchData = async () => {
  if (!props.nodeKey) return
  uiState.loading = true
  const params = {
    page_size: pagination.pageSize,
    page: pagination.current,
    is_page: true,
    key: props.nodeKey,
    search: uiData.searchValue,
  }

  const res = await getOrganizationsUsersApi(params).catch(() => {})
  if (res) {
    pagination.total = res.total
    uiData.tableData = res.data
  }
  uiState.loading = false
}

// watch
watch(
  () => props.nodeKey,
  (val) => {
    if (val) {
      pagination.current = 1
      uiData.searchValue = ''
      fetchData()
    } else {
      uiData.tableData = []
      pagination.total = 0
    }
  },
  { immediate: true },
)

// 搜索
const handleSearch = (value) => {
  uiData.searchValue = value
  pagination.current = 1
  fetchData()
}

// 翻页
const handleTableChange = (pager) => {
  pagination.current = pager.current
  pagination.pageSize = pager.pageSize
  fetchData()
}

// 新增
const handleAdd = () => {
  getUsers()
  formState.value = { ...defaultForm }
  uiState.isModalOpen = true
}

// 提交
const onSubmit = useThrottleFn(async (data) => {
  const res = await bindOrganizationsUsersApi(data).catch(() => {})
  if (res) {
    message.success('用户绑定成功')
    uiState.isModalOpen = false
    fetchData()
  }
})

// 删除
const handleDelete = useThrottleFn(async (record) => {
  const payload = {
    key: record.organization_key,
    uid: record.uid,
  }
  const res = await deleteOrganizationsUsersApi(payload).catch(() => {})
  if (res) {
    message.info('用户移除成功')
    fetchData()
  }
})

// 获取所有用户列表
const getUsers = async () => {
  const res = await getUsersApi().catch(() => {})
  uiData.users = res.data || []

  const rolesRes = await getRolesApi().catch(() => {})
  uiData.roles = rolesRes.data || []
}
</script>

<style scoped>
.org-users {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 18px 18px 14px;
  background:
    radial-gradient(circle at right top, rgba(22, 84, 194, 0.08), rgba(22, 84, 194, 0) 40%),
    linear-gradient(180deg, rgba(249, 251, 255, 1) 0%, rgba(255, 255, 255, 1) 35%);
}

.users-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 2px 2px 14px;
}

.users-caption {
  display: inline-block;
  color: #1f6feb;
  font-size: 12px;
  font-weight: 700;
  background: rgba(31, 111, 235, 0.12);
  padding: 3px 8px;
  border-radius: 999px;
}

.users-header h3 {
  margin: 8px 0 6px;
  font-size: 22px;
  color: #16213c;
}

.users-header p {
  margin: 0;
  color: #5f6b8a;
  font-size: 13px;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  flex-wrap: wrap;
  margin: 0 2px 12px;
}

.toolbar-tags {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.users-table {
  flex: 1;
  min-height: 0;
}

.compact-mode {
  padding: 14px 14px 10px;
}

.compact-mode .users-header {
  padding-bottom: 10px;
}

.compact-mode .users-header h3 {
  font-size: 19px;
  margin: 6px 0 4px;
}

.compact-mode .users-header p {
  font-size: 12px;
}

.compact-mode .users-caption {
  font-size: 11px;
  padding: 2px 8px;
}

.compact-mode .toolbar {
  margin-bottom: 8px;
}

:deep(.users-table .ant-table) {
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #e3eaf8;
}

:deep(.users-table .ant-table-thead > tr > th) {
  background: #f3f7ff;
  color: #23335c;
  font-weight: 600;
}

:deep(.users-table .ant-table-tbody > tr > td) {
  border-bottom-color: #edf2fb;
}

:deep(.users-table .ant-table-pagination.ant-pagination) {
  margin-bottom: 0;
}

@media (max-width: 1080px) {
  .org-users {
    padding: 14px;
  }

  .users-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .users-header h3 {
    font-size: 18px;
  }
}
</style>
