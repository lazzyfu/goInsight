<template>
  <a-card title="环境管理">
    <!-- 卡片右上角的新增按钮 -->
    <template #extra>
      <a-button type="primary" @click="handleAddRecord"><PlusOutlined />新增环境</a-button>
    </template>
    <!-- 搜索区域 -->
    <div class="search-wrapper">
      <!-- 搜索 -->
      <a-input-search
        v-model:value="searchValue"
        placeholder="搜索环境名..."
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
            <a-space>
              <a @click="handleEditRecord(record)"> <EditOutlined /> 编辑 </a>
              <a-popconfirm
                title="确认删除吗？"
                ok-text="是"
                cancel-text="否"
                @confirm="handleDeleteRecord(record)"
              >
                <a><DeleteOutlined /> 删除</a>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </div>
  </a-card>
  <!-- 新增/编辑弹窗 -->
  <Modal
    :open="state.isModalOpen"
    :formState="formState"
    :title="state.isEditModal ? '编辑环境' : '新增环境'"
    @update:open="state.isModalOpen = $event"
    @submit="handleSubmit"
  />
</template>

<script setup>
import {
  createEnvironmentsApi,
  deleteEnvironmentsApi,
  getEnvironmentsApi,
  updateEnvironmentsApi,
} from '@/api/admin'
import { DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref } from 'vue'
import Modal from './components/Modal.vue'

// 状态管理
const state = reactive({
  loading: false,
  isEditModal: false,
  isModalOpen: false,
})
const searchValue = ref('')
const defaultForm = {
  name: '',
}
const formState = ref({ ...defaultForm })

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
    title: '环境',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    key: 'created_at',
  },
  {
    title: '更新时间',
    dataIndex: 'updated_at',
    key: 'updated_at',
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
  state.loading = true
  const params = {
    page_size: pagination.pageSize,
    page: pagination.current,
    is_page: true,
    search: searchValue.value,
  }
  const res = await getEnvironmentsApi(params)
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
  state.isEditModal = false
  formState.value = { ...defaultForm }
  state.isModalOpen = true
}

// 编辑记录
const handleEditRecord = (record) => {
  state.isEditModal = true
  formState.value = { ...record }
  state.isModalOpen = true
}

const handleSubmit = async (data) => {
  const res = state.isEditModal
    ? await updateEnvironmentsApi(data)
    : await createEnvironmentsApi(data)
  if (res?.code === '0000') {
    message.success('操作成功')
    state.isModalOpen = false
    fetchData()
  }
}

// 删除记录
const handleDeleteRecord = async (record) => {
  const res = await deleteEnvironmentsApi(record.id).catch(() => {})
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
