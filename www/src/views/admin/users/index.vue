<template>
  <a-card title="用户管理">
    <template #extra>
      <a-button type="primary" @click="handleAddUser"><PlusOutlined />新增用户</a-button>
    </template>
    <div class="search-wrapper">
      <!-- 搜索 -->
      <a-input-search
        v-model:value="searchValue"
        placeholder="搜索用户名、手机号、邮箱..."
        style="width: 350px"
        @search="handleSearch"
      />
    </div>
    <!-- 表格 -->
    <div style="margin-top: 12px">
      <a-table
        size="small"
        :columns="tableColumns"
        :row-key="(record) => record.key"
        :data-source="tableData"
        @resizeColumn="handleResizeColumn"
        :pagination="pagination"
        :loading="state.loading"
        @change="handleTableChange"
        :scroll="{ x: 1300 }"
      >
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
                    <a @click="handleEditUser(record)"> <EditOutlined /> 编辑 </a>
                  </a-menu-item>
                  <a-menu-item key="2">
                    <a-popconfirm
                      title="确认删除吗？"
                      ok-text="是"
                      cancel-text="否"
                      @confirm="handleDeleteUser(record)"
                    >
                      <a><DeleteOutlined /> 删除</a>
                    </a-popconfirm>
                  </a-menu-item>
                  <a-menu-item key="3">
                    <a @click="handleResetPassword(record)"><KeyOutlined /> 重置密码</a>
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
  <ResetPasswordModal
    :open="isResetPasswordModalOpen"
    @update:open="isResetPasswordModalOpen = $event"
    @submit="handleResetPasswordSubmit"
  />
  <!-- 新增/编辑弹窗 -->
  <UserModal
    :open="state.userModalOpen"
    :formState="userFormState"
    :title="state.isEditUserModal ? '编辑用户' : '新增用户'"
    @update:open="state.userModalOpen = $event"
    @submit="handleUserSubmit"
  />
</template>

<script setup>
import {
  addUsersApi,
  deleteUsersApi,
  getUsersApi,
  ResetPasswordApi,
  updateUsersApi,
} from '@/api/admin'
import ResetPasswordModal from '@/views/admin/users/components/PasswordModal.vue'
import UserModal from '@/views/admin/users/components/UserModal.vue'
import {
  DeleteOutlined,
  EditOutlined,
  EllipsisOutlined,
  KeyOutlined,
  PlusOutlined,
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref } from 'vue'

// 状态管理
const state = reactive({
  loading: false,
  isEditUserModal: false,
  userModalOpen: false,
})

const searchValue = ref('')

const defaultUserForm = {
  username: '',
  password: '',
  nick_name: '',
  email: '',
  mobile: '',
  role: '',
  is_active: true,
  is_two_fa: true,
  is_superuser: false,
}

const uid = ref(0)
const userFormState = ref({ ...defaultUserForm })
const isResetPasswordModalOpen = ref(false)

// 搜索
const handleSearch = (value) => {
  searchValue.value = value
  pagination.current = 1
  fetchData()
}

// 表
const tableData = ref([])
const tableColumns = [
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
]

// 分页
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  pageSizeOptions: ['10', '20', '50', '100'],
  showSizeChanger: true,
})

// 获取表数据
const fetchData = async () => {
  state.loading = true
  const params = {
    page_size: pagination.pageSize,
    page: pagination.current,
    is_page: true,
    search: searchValue.value,
  }
  const res = await getUsersApi(params)
  if (res) {
    pagination.total = res.total
    tableData.value = res.data
  }
  state.loading = false
}

// 翻页
const handleTableChange = (pager) => {
  pagination.current = pager.current
  pagination.pageSize = pager.pageSize
  fetchData()
}

function handleResizeColumn(w, col) {
  col.width = w
}

// 新增用户
const handleAddUser = () => {
  state.isEditUserModal = false
  userFormState.value = { ...defaultUserForm }
  state.userModalOpen = true
}

// 编辑用户
const handleEditUser = (record) => {
  state.isEditUserModal = true
  userFormState.value = { ...record }
  state.userModalOpen = true
}

const handleUserSubmit = async (data) => {
  const res = state.isEditUserModal ? await updateUsersApi(data) : await addUsersApi(data)
  if (res?.code === '0000') {
    message.success('操作成功')
    state.userModalOpen = false
    fetchData()
  }
}

// 重置密码
const handleResetPassword = (record) => {
  uid.value = record.uid
  isResetPasswordModalOpen.value = true
}

const handleResetPasswordSubmit = async (data) => {
  const payload = {
    uid: uid.value,
    ...data,
  }
  const res = await ResetPasswordApi(payload).catch(() => {})
  if (res?.code === '0000') {
    message.info('操作成功')
  }
  isResetPasswordModalOpen.value = false
}

// 删除用户
const handleDeleteUser = async (record) => {
  const res = await deleteUsersApi(record.uid).catch(() => {})
  if (res?.code === '0000') {
    message.info('操作成功')
    fetchData()
  }
}

// 生命周期
onMounted(() => {
  fetchData()
})
</script>
