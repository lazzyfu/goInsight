<template>
  <a-card title="用户管理">
    <template #extra>
      <a-button type="primary" @click="handleAdd"><PlusOutlined />绑定用户</a-button>
    </template>
    <div class="search-wrapper">
      <a-input-search
        v-model:value="searchValue"
        placeholder="搜索用户..."
        style="width: 350px"
        size="small"
        @search="handleSearch"
      />
    </div>
    <div style="margin-top: 12px">
      <a-table
        size="small"
        :columns="tableColumns"
        :row-key="(record) => record.uid"
        :data-source="tableData"
        :pagination="pagination"
        :loading="state.loading"
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
                <a><DeleteOutlined /> 移除</a>
              </a-popconfirm>
            </a-tooltip>
          </template>
        </template>
      </a-table>
    </div>
  </a-card>
  <BindOrgUsers
    :open="state.isModalOpen"
    :nodeKey="props.nodeKey"
    :users="state.allUsers"
    @update:open="state.isModalOpen = $event"
    @submit="onSubmit"
  />
</template>

<script setup>
import {
  bindOrganizationsUsersApi,
  deleteOrganizationsUsersApi,
  getOrganizationsUsersApi,
  getUsersApi,
} from '@/api/admin'
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref, watch } from 'vue'
import BindOrgUsers from './BindOrgUsers.vue'

const props = defineProps({
  nodeKey: {
    type: String,
    required: true,
  },
})

// --- 状态和分页定义 ---
const state = reactive({
  loading: false,
  isModalOpen: false,
  allUsers: [],
})

const searchValue = ref('')
const tableData = ref([])

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  pageSizeOptions: ['10', '20', '50', '100'],
  showSizeChanger: true,
})

const tableColumns = [
  {
    title: '用户名',
    dataIndex: 'username',
    key: 'username',
    width: 150,
  },
  {
    title: '昵称',
    dataIndex: 'nick_name',
    key: 'nick_name',
    width: 150,
  },
  {
    title: '操作',
    dataIndex: 'action',
    key: 'action',
    width: 80,
  },
]

// --- 核心函数定义 (放在 watch 之前，修复 ReferenceError) ---

// 获取当前组织的用户列表
const fetchData = async () => {
  if (!props.nodeKey) return
  state.loading = true
  const params = {
    page_size: pagination.pageSize,
    page: pagination.current,
    is_page: true,
    key: props.nodeKey,
    search: searchValue.value,
  }

  const res = await getOrganizationsUsersApi(params)
  if (res) {
    pagination.total = res.total
    tableData.value = res.data
  }
  state.loading = false
}

// 获取所有用户列表
const getAllUsers = async () => {
  const res = await getUsersApi().catch(() => {})
  state.allUsers = res.data || []
}

// --- watch 监听 ---

watch(
  () => props.nodeKey,
  (val) => {
    if (val) {
      // 切换组织时，重置分页和搜索条件
      pagination.current = 1
      searchValue.value = ''
      fetchData() // ✅ 此时 fetchData 已定义
    } else {
      tableData.value = []
      pagination.total = 0
    }
  },
  { immediate: true },
)

// --- 操作逻辑 ---

const handleSearch = (value) => {
  searchValue.value = value
  pagination.current = 1
  fetchData()
}

const handleTableChange = (pager) => {
  pagination.current = pager.current
  pagination.pageSize = pager.pageSize
  fetchData()
}

const handleAdd = () => {
  state.isModalOpen = true
}

const onSubmit = async (data) => {
  const res = await bindOrganizationsUsersApi(data).catch(() => {})
  if (res?.code === '0000') {
    message.success('用户绑定成功')
    state.isModalOpen = false
    fetchData()
  }
}

const handleDelete = async (record) => {
  const payload = {
    key: props.nodeKey,
    uid: record.uid,
  }
  const res = await deleteOrganizationsUsersApi(payload).catch(() => {})
  if (res?.code === '0000') {
    message.info('用户移除成功')
    fetchData()
  }
}

// --- 生命周期 ---

onMounted(() => {
  getAllUsers()
})
</script>
