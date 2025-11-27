<template>
  <a-card title="角色管理">
    <!-- 卡片右上角的新增按钮 -->
    <template #extra>
      <a-button type="primary" @click="handleAdd"> <PlusOutlined />新增角色 </a-button>
    </template>

    <!-- 搜索区域 -->
    <div class="search-wrapper">
      <!-- 搜索 -->
      <a-input-search
        v-model:value="uiData.searchValue"
        placeholder="搜索角色名..."
        style="width: 350px"
        @search="handleSearch"
      />
    </div>

    <!-- 表格 -->
    <div style="margin-top: 12px">
      <a-table
        size="small"
        :columns="uiData.tableColumns"
        :row-key="(record) => record.key"
        :data-source="uiData.tableData"
        :pagination="pagination"
        :loading="uiState.loading"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action'">
            <a-space>
              <a @click="handleEdit(record)"> <EditOutlined /> 编辑 </a>
              <a-popconfirm
                title="确认删除吗？"
                ok-text="是"
                cancel-text="否"
                @confirm="handleDelete(record)"
              >
                <a> <DeleteOutlined /> 删除 </a>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </div>
  </a-card>

  <!-- 新增/编辑弹窗 -->
  <RoleFormModal
    :open="uiState.isModalOpen"
    v-model:modelValue="formState"
    :title="uiState.isEditMode ? '编辑角色' : '新增角色'"
    @update:open="uiState.isModalOpen = $event"
    @submit="onSubmit"
  />
</template>

<script setup>
import { createRolesApi, deleteRolesApi, getRolesApi, updateRolesApi } from '@/api/admin'
import { DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref } from 'vue'
import RoleFormModal from './RoleFormModal.vue'

// 状态
const uiState = reactive({
  loading: false,
  isEditMode: false,
  isModalOpen: false,
})

// 数据
const uiData = reactive({
  searchValue: '',
  tableData: [],
  tableColumns: [
    {
      title: '角色',
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
      fixed: 'right',
    },
  ],
})

// form表单
const defaultForm = {
  name: '',
}
const formState = ref({ ...defaultForm })

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
  formState.value = { ...defaultForm }
  uiState.isModalOpen = true
}

// 编辑
const handleEdit = (record) => {
  uiState.isEditMode = true
  formState.value = { ...record }
  uiState.isModalOpen = true
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
  const res = await getRolesApi(params).catch(() => {})
  if (res) {
    pagination.total = res.total
    uiData.tableData = res.data
  }
  uiState.loading = false
}

// 提交表单
const onSubmit = async (data) => {
  const res = uiState.isEditMode
    ? await updateRolesApi(data).catch(() => {})
    : await createRolesApi(data).catch(() => {})
  if (res?.code === '0000') {
    message.success('操作成功')
    uiState.isModalOpen = false
    fetchData()
  }
}

// 删除记录
const handleDelete = async (record) => {
  const res = await deleteRolesApi(record.id).catch(() => {})
  if (res?.code === '0000') {
    message.info('操作成功')
    fetchData()
  }
}

// 初始化
onMounted(() => {
  fetchData()
})
</script>
