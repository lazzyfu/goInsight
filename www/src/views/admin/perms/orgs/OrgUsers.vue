<template>
  <a-card title="用户管理">
    <template #extra>
      <a-button type="primary" @click="handleAdd"> <PlusOutlined />绑定用户 </a-button>
    </template>

    <div class="search-wrapper">
      <a-input-search
        v-model:value="uiData.searchValue"
        placeholder="搜索用户..."
        style="width: 350px"
        @search="handleSearch"
      />
    </div>

    <div style="margin-top: 12px">
      <a-table
        size="small"
        :columns="uiData.tableColumns"
        :row-key="(record) => record.uid"
        :data-source="uiData.tableData"
        :pagination="pagination"
        :loading="uiState.loading"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action'">
            <a-tooltip title="移除当前用户">
              <a-popconfirm
                title="确认移除吗？"
                ok-text="是"
                cancel-text="否"
                @confirm="handleDelete(record)"
              >
                <a> <DeleteOutlined /> 移除 </a>
              </a-popconfirm>
            </a-tooltip>
          </template>
        </template>
      </a-table>
    </div>
  </a-card>
  <BindOrgUsers
    :open="uiState.isModalOpen"
    :nodeKey="props.nodeKey"
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

// 分页
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
    key: props.nodeKey,
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

  const rolesRes = await getRolesApi().catch(() => { })
  uiData.roles = rolesRes.data || []
}
</script>
