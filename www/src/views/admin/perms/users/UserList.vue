<template>
  <a-card title="用户管理">
    <template #extra>
      <a-button type="primary" @click="handleAdd">
        <PlusOutlined />新增用户
      </a-button>
    </template>

    <div class="search-wrapper">
      <a-space>
        <a-cascader v-model:value="uiData.searchOrganizationKey"
          :field-names="{ label: 'title', value: 'key', children: 'children' }" :options="uiData.organizations"
          change-on-select expand-trigger="hover" placeholder="请选择组织">
        </a-cascader>

        <!-- 搜索 -->
        <a-input-search v-model:value="uiData.searchValue" placeholder="搜索用户名、昵称、手机号、邮箱..." style="width: 350px"
          @search="handleSearch" />
      </a-space>
    </div>

    <!-- 表格 -->
    <div style="margin-top: 12px">
      <a-table size="small" :columns="uiData.tableColumns" :row-key="(record) => record.key"
        :data-source="uiData.tableData" :pagination="pagination" :loading="uiState.loading" @change="handleTableChange"
        :scroll="{ x: 1100 }">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'is_active'">
            <span v-if="record.is_active">
              <a-tag color="green">是</a-tag>
            </span>
            <span v-else>
              <a-tag color="volcano">否</a-tag>
            </span>
          </template>
          <template v-if="column.key === 'is_superuser'">
            <span v-if="record.is_superuser">
              <a-tag color="green">是</a-tag>
            </span>
            <span v-else>
              <a-tag color="volcano">否</a-tag>
            </span>
          </template>
          <template v-if="column.key === 'is_two_fa'">
            <span v-if="record.is_two_fa">
              <a-tag color="green">是</a-tag>
            </span>
            <span v-else>
              <a-tag color="volcano">否</a-tag>
            </span>
          </template>
          <template v-if="column.key === 'action'">
            <a-dropdown>
              <EllipsisOutlined />
              <template #overlay>
                <a-menu>
                  <a-menu-item key="1">
                    <a @click="handleEdit(record)">
                      <EditOutlined /> 编辑
                    </a>
                  </a-menu-item>
                  <a-menu-item key="2">
                    <a-popconfirm title="确认删除吗？" ok-text="是" cancel-text="否" @confirm="handleDelete(record)">
                      <a>
                        <DeleteOutlined /> 删除
                      </a>
                    </a-popconfirm>
                  </a-menu-item>
                  <a-menu-item key="3">
                    <a @click="handleResetPassword(record)">
                      <KeyOutlined /> 重置密码
                    </a>
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </template>
        </template>
      </a-table>
    </div>
  </a-card>

  <!-- 重置密码 -->
  <PasswordFormModal :open="uiState.passwordModalOpen" :title="uiData.passwordFormTitle"
    @update:open="uiState.passwordModalOpen = $event" @submit="handleResetPasswordSubmit" />

  <!-- 新增/编辑弹窗 -->
  <UserFormModal :open="uiState.userModalOpen" v-model:modelValue="formState" :roles="uiData.roles"
    :title="uiState.isEditMode ? '编辑用户' : '新增用户'" @update:open="uiState.userModalOpen = $event" @submit="onSubmit" />
</template>

<script setup>
import {
  addUsersApi,
  deleteUsersApi,
  getOrganizationsApi,
  getRolesApi,
  getUsersApi,
  ResetPasswordApi,
  updateUsersApi,
} from '@/api/admin'
import {
  DeleteOutlined,
  EditOutlined,
  EllipsisOutlined,
  KeyOutlined,
  PlusOutlined,
} from '@ant-design/icons-vue'
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref } from 'vue'
import PasswordFormModal from './PasswordFormModal.vue'
import UserFormModal from './UserFormModal.vue'

// 状态管理
const uiState = reactive({
  loading: false,
  isEditMode: false,
  userModalOpen: false,
  passwordModalOpen: false,
})

// 数据
const uiData = reactive({
  searchValue: '',
  searchOrganizationKey: '',
  roles: [],
  organizations: [],
  passwordFormTitle: '',
  tableData: [],
  tableColumns: [
    {
      title: '用户',
      dataIndex: 'username',
      key: 'username',
      fixed: 'left',
    },
    {
      title: '昵称',
      dataIndex: 'nick_name',
      key: 'nick_name',
    },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
    },
    {
      title: '激活',
      dataIndex: 'is_active',
      key: 'is_active',
    },
    {
      title: '2FA认证',
      dataIndex: 'is_two_fa',
      key: 'is_two_fa',
    },
    {
      title: '管理员',
      dataIndex: 'is_superuser',
      key: 'is_superuser',
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: '手机号',
      dataIndex: 'mobile',
      key: 'mobile',
    },
    {
      title: '组织',
      dataIndex: 'organization',
      key: 'organization',
    },
    {
      title: '加入时间',
      dataIndex: 'date_joined',
      key: 'date_joined',
    },
    {
      title: '操作',
      key: 'action',
      fixed: 'right',
      width: 120,
    },
  ],
})

// uid
const uid = ref(0)

// form表单
const defaultUserForm = {
  username: '',
  password: '',
  nick_name: '',
  email: '',
  mobile: '',
  role_id: null,
  is_active: true,
  is_two_fa: true,
  is_superuser: false,
}
const formState = ref({ ...defaultUserForm })

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

// 获取角色
const getRoles = async () => {
  const res = await getRolesApi().catch(() => { })
  uiData.roles = res?.data || []
}

// 获取组织
const getOrganizations = async () => {
  const res = await getOrganizationsApi().catch(() => { })
  uiData.organizations = res?.data || []
}

// 获取列表数据
const fetchData = async () => {
  uiState.loading = true
  let organization_key = ''
  if (uiData.searchOrganizationKey != undefined) {
    organization_key = uiData.searchOrganizationKey[uiData.searchOrganizationKey.length - 1]
  }
  const params = {
    page_size: pagination.pageSize,
    page: pagination.current,
    is_page: true,
    search: uiData.searchValue,
    organization_key: organization_key,
  }
  const res = await getUsersApi(params).catch(() => { })
  if (res) {
    pagination.total = res.total
    uiData.tableData = res.data
  }
  uiState.loading = false
}

// 新增
const handleAdd = () => {
  uiState.isEditMode = false
  formState.value = { ...defaultUserForm }
  getRoles()
  uiState.userModalOpen = true
}

// 编辑
const handleEdit = (record) => {
  uiState.isEditMode = true
  formState.value = { ...record }
  getRoles()
  uiState.userModalOpen = true
}

// 提交
const onSubmit = useThrottleFn(async (data) => {
  const res = uiState.isEditMode
    ? await updateUsersApi(data).catch(() => { })
    : await addUsersApi(data).catch(() => { })
  if (res) {
    message.success('操作成功')
    uiState.userModalOpen = false
    fetchData()
  }
})

// 重置密码
const handleResetPassword = (record) => {
  uid.value = record.uid
  uiData.passwordFormTitle = `重置用户${record.username}密码`
  uiState.passwordModalOpen = true
}

// 提交重置密码
const handleResetPasswordSubmit = useThrottleFn(async (data) => {
  const payload = {
    uid: uid.value,
    ...data,
  }
  const res = await ResetPasswordApi(payload).catch(() => { })
  if (res) {
    message.info('操作成功')
  }
  uiState.passwordModalOpen = false
})

// 删除
const handleDelete = useThrottleFn(async (record) => {
  const res = await deleteUsersApi(record.uid).catch(() => { })
  if (res) {
    message.info('操作成功')
    fetchData()
  }
})

// 初始化
onMounted(() => {
  fetchData()
  getOrganizations()
})
</script>
