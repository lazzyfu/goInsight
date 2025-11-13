<template>
  <a-card size="small">
    <div class="search-wrapper">
      <!-- 搜索 -->
      <a-space>
        <a-input-search
          v-model:value="searchValue"
          placeholder="输入搜索的内容"
          style="width: 350px"
          @search="onSearch"
        />
        <a-button type="primary" @click="onAdd">新增</a-button>
      </a-space>
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
        :loading="loading"
        @change="handleTableChange"
        :scroll="{ x: 1500 }"
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
            <a-space wrap>
              <a-tooltip title="编辑">
                <a-button
                  type="link"
                  block
                  shape="circle"
                  :icon="h(EditOutlined)"
                  @click="onEdit(record)"
                />
              </a-tooltip>
              <a-tooltip title="修改密码">
                <a-button
                  type="link"
                  block
                  shape="circle"
                  :icon="h(CopyOutlined)"
                  @click="changePassword(record)"
                />
              </a-tooltip>
              <a-popconfirm
                title="确认删除吗？"
                ok-text="是"
                cancel-text="否"
                @confirm="onDelete(record)"
              >
                <a-button block type="link" shape="circle" :icon="h(DeleteOutlined)"> </a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </div>
  </a-card>
  <!-- 修改密码 -->
  <ChangePassword
    :open="isChangePasswordModalOpen"
    @update:open="isChangePasswordModalOpen = $event"
    @submit="handleChangePasswordSubmit"
  />
  <!-- 新增/编辑弹窗 -->
  <UserModal
    :open="modalOpen"
    :formState="formState"
    :title="isEditModal ? '编辑用户' : '新增用户'"
    @update:open="modalOpen = $event"
    @submit="handleSubmit"
  />
</template>

<script setup>
import {
  addUsersApi,
  changePasswordApi,
  deleteUsersApi,
  getUsersApi,
  updateUsersApi,
} from '@/api/admin'
import ChangePassword from '@/views/admin/users/components/PasswordModal.vue'
import UserModal from '@/views/admin/users/components/UserModal.vue'
import { CopyOutlined, DeleteOutlined, EditOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { h, onMounted, reactive, ref } from 'vue'

// 状态管理
const modalOpen = ref(false)
const isEditModal = ref(false)
const loading = ref(false)
const searchValue = ref('')
const uid = ref(0)

const defaultForm = {
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

const formState = ref({ ...defaultForm })
const tableData = ref([])
const isChangePasswordModalOpen = ref(false)

// 搜索
const onSearch = (value) => {
  searchValue.value = value
  pagination.current = 1
  fetchData()
}

// 表列
const tableColumns = [
  {
    title: '用户',
    dataIndex: 'username',
    key: 'username',
    scopedSlots: {
      customRender: 'username',
    },
  },
  {
    title: '昵称',
    dataIndex: 'nick_name',
    key: 'nick_name',
    scopedSlots: {
      customRender: 'nick_name',
    },
  },
  {
    title: '角色',
    dataIndex: 'role',
    key: 'role',
    scopedSlots: {
      customRender: 'role',
    },
  },
  {
    title: '激活',
    dataIndex: 'is_active',
    key: 'is_active',
    scopedSlots: {
      customRender: 'is_active',
    },
  },
  {
    title: '2FA认证',
    dataIndex: 'is_two_fa',
    key: 'is_two_fa',
    scopedSlots: {
      customRender: 'is_two_fa',
    },
  },
  {
    title: '管理员',
    dataIndex: 'is_superuser',
    key: 'is_superuser',
    scopedSlots: {
      customRender: 'is_superuser',
    },
  },
  {
    title: '邮箱',
    dataIndex: 'email',
    key: 'email',
    scopedSlots: {
      customRender: 'email',
    },
  },
  {
    title: '手机号',
    dataIndex: 'mobile',
    key: 'mobile',
    scopedSlots: {
      customRender: 'mobile',
    },
  },
  {
    title: '组织',
    dataIndex: 'organization',
    key: 'organization',
    scopedSlots: {
      customRender: 'organization',
    },
  },
  {
    title: '加入时间',
    dataIndex: 'date_joined',
    key: 'date_joined',
    scopedSlots: {
      customRender: 'date_joined',
    },
  },
  {
    title: '操作',
    dataIndex: 'action',
    key: 'action',
    scopedSlots: {
      customRender: 'action',
    },
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
  loading.value = true
  const params = {
    page_size: pagination.pageSize,
    page: pagination.current,
    is_page: true,
    search: searchValue.value,
  }
  const res = await getUsersApi(params).catch(() => {})
  if (res) {
    pagination.total = res.total
    tableData.value = res.data
  }
  loading.value = false
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

// 弹窗逻辑
const onAdd = () => {
  isEditModal.value = false
  formState.value = { ...defaultForm }
  modalOpen.value = true
}

const onEdit = (record) => {
  isEditModal.value = true
  formState.value = { ...record }
  modalOpen.value = true
}

const handleSubmit = async (data) => {
  if (isEditModal.value) {
    const res = await updateUsersApi(data)
    if (res?.code === '0000') {
      message.info('操作成功')
    }
  } else {
    const res = await addUsersApi(data).catch(() => {})
    if (res?.code === '0000') {
      message.info('操作成功')
    }
  }
  modalOpen.value = false
  fetchData()
}

// 修改密码
const changePassword = (record) => {
  uid.value = record.uid
  isChangePasswordModalOpen.value = true
}

const handleChangePasswordSubmit = async (data) => {
  const payload = {
    uid: uid.value,
    ...data,
  }
  const res = await changePasswordApi(payload).catch(() => {})
  if (res?.code === '0000') {
    message.info('操作成功')
  }
  isChangePasswordModalOpen.value = false
}

// 删除用户
const onDelete = async (record) => {
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
