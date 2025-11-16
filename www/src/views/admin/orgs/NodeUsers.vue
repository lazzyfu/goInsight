<template>
  <!-- 角色管理卡片 -->
  <a-card title="用户管理">
    <!-- 卡片右上角的新增按钮 -->
    <template #extra>
      <a-button type="primary" @click="handleAddRecord"><PlusOutlined />拉选用户</a-button>
    </template>
    <!-- 搜索区域 -->
    <div class="search-wrapper">
      <!-- 搜索 -->
      <a-input-search
        v-model:value="searchValue"
        placeholder="搜索用户..."
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
        :pagination="pagination"
        :loading="state.loading"
        @change="handleTableChange"
        :scroll="{ x: 1300 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action'">
            <a-tooltip title="移除当前用户">
              <a-popconfirm
                title="确认移除吗？"
                ok-text="是"
                cancel-text="否"
                @confirm="handleDeleteRecord(dataRef)"
              >
                <a-button type="text" size="small" danger @click.stop>
                  <template #icon><DeleteOutlined /></template>
                </a-button>
              </a-popconfirm>
            </a-tooltip>
          </template>
        </template>
      </a-table>
    </div>
  </a-card>
  <AddNodeUsersModal
    :open="state.isModalOpen"
    :nodeKey="props.nodeKey"
    @update:open="state.isModalOpen = $event"
    @submit="handleSubmit"
  />
</template>

<script setup>
import {
  bindOrganizationsUsersApi,
  deleteOrganizationsUsersApi,
  getOrganizationsUsersApi,
} from '@/api/admin'
import AddNodeUsersModal from '@/views/admin/orgs/AddNodeUsersModal.vue'
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { reactive, ref, watch } from 'vue'

// 状态管理
const state = reactive({
  loading: false,
  isModalOpen: false,
})

const props = defineProps({
  open: Boolean,
  nodeKey: String,
})
const searchValue = ref('')

watch(
  () => props.nodeKey,
  (val) => {
    if (val) {
      fetchData()
    }
  },
)
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

// 翻页
const handleTableChange = (pager) => {
  pagination.current = pager.current
  pagination.pageSize = pager.pageSize
  fetchData()
}

// 新增记录
const handleAddRecord = () => {
  state.isModalOpen = true
}

const handleSubmit = async (data) => {
  console.log('data: ', data)
  const res = await bindOrganizationsUsersApi(data)
  if (res?.code === '0000') {
    message.success('操作成功')
    state.isModalOpen = false
    fetchData()
  }
}

// 删除记录
const handleDeleteRecord = async (record) => {
  const res = await deleteOrganizationsUsersApi(record.id).catch(() => {})
  if (res?.code === '0000') {
    message.info('操作成功')
    fetchData()
  }
}
</script>
