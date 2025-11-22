<template>
  <a-card title="表权限管理">
    <!-- 卡片右上角的新增按钮 -->
    <template #extra>
      <a-button type="primary" @click="handleAdd"><PlusOutlined />新增表权限</a-button>
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
            <a-popconfirm
              title="确认删除吗？"
              ok-text="是"
              cancel-text="否"
              @confirm="handleDelete(record)"
            >
              <a><DeleteOutlined /> 删除</a>
            </a-popconfirm>
          </template>
        </template>
      </a-table>
    </div>
  </a-card>
  <!-- 新增/编辑弹窗 -->
  <DasTableFormModal
    :open="state.isModalOpen"
    v-model:modelValue="formState"
    :tables="state.tableList"
    title="新增表权限"
    @update:open="state.isModalOpen = $event"
    @submit="onSubmit"
  />
</template>

<script setup>
import {
  createTablesGrantApi,
  deleteTablesGrantApi,
  getTablesGrantApi,
  getTablesListApi,
} from '@/api/admin'
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import DasTableFormModal from './DasTableFormModal.vue'

const route = useRoute()
const routeParams = route.params

// 状态管理
const state = reactive({
  loading: false,
  isModalOpen: false,
  tableList: [],
})
const searchValue = ref('')
const defaultForm = {
  tables: [],
  rule: '',
}
const formState = ref({ ...defaultForm })

// 表
const tableData = ref([])
const tableColumns = [
  {
    title: '规则',
    dataIndex: 'rule',
    key: 'rule',
  },
  {
    title: '用户',
    dataIndex: 'username',
    key: 'username',
  },
  {
    title: '库名',
    dataIndex: 'schema',
    key: 'schema',
  },
  {
    title: '表名',
    dataIndex: 'table',
    key: 'table',
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
]

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

// 新增记录
const handleAdd = () => {
  formState.value = { ...defaultForm }
  getTables()
  state.isModalOpen = true
}

// 获取表数据
const fetchData = async () => {
  state.loading = true
  const params = {
    page_size: pagination.pageSize,
    page: pagination.current,
    is_page: true,
    search: searchValue.value,
    ...routeParams,
  }
  const res = await getTablesGrantApi(params)
  if (res) {
    pagination.total = res.total
    tableData.value = res.data
  }
  state.loading = false
}

const getTables = async () => {
  const res = await getTablesListApi(routeParams)
  state.tableList = res.data || []
}

// 提交表单
const onSubmit = async (data) => {
  const payload = {
    ...data,
    ...routeParams,
  }
  const res = await createTablesGrantApi(payload)
  if (res?.code === '0000') {
    message.success('操作成功')
    state.isModalOpen = false
    fetchData()
  }
}

// 删除记录
const handleDelete = async (record) => {
  const res = await deleteTablesGrantApi(record.id).catch(() => {})
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
