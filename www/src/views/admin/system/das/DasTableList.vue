<template>
  <a-card title="配置用户表访问权限">
    <!-- 卡片右上角的新增按钮 -->
    <template #extra>
      <a-button type="primary" @click="handleAdd">
        <PlusOutlined />新增表权限
      </a-button>
    </template>

    <!-- 搜索区域 -->
    <div class="search-wrapper">
      <!-- 搜索 -->
      <a-input-search v-model:value="uiData.searchValue" placeholder="搜索环境名..." style="width: 350px"
        @search="handleSearch" />
    </div>

    <!-- 表格 -->
    <div style="margin-top: 12px">
      <a-table size="small" :columns="uiData.tableColumns" :row-key="(record) => record.key"
        :data-source="uiData.tableData" :pagination="pagination" :loading="uiState.loading" @change="handleTableChange"
        :scroll="{ x: 1100 }">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action'">
            <a-popconfirm title="确认删除吗？" ok-text="是" cancel-text="否" @confirm="handleDelete(record)">
              <a>
                <DeleteOutlined /> 删除
              </a>
            </a-popconfirm>
          </template>
        </template>
      </a-table>
    </div>
  </a-card>

  <!-- 新增/编辑弹窗 -->
  <DasTableFormModal :open="uiState.isModalOpen" v-model:modelValue="formState" :tables="uiData.tableList" title="新增表权限"
    @update:open="uiState.isModalOpen = $event" @submit="onSubmit" />
</template>

<script setup>
import {
  createDasTablesGrantApi,
  deleteDasTablesGrantApi,
  getDasTablesGrantApi,
  getDasTablesListApi,
} from '@/api/admin'
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import DasTableFormModal from './DasTableFormModal.vue'

// 路由参数
const route = useRoute()
const routeParams = route.params

// 状态管理
const uiState = reactive({
  loading: false,
  isModalOpen: false,
})

// 数据
const uiData = reactive({
  searchValue: '',
  tableList: [],
  tableData: [],
  tableColumns: [
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
})

// form表单
const defaultForm = {
  tables: [],
  rule: '',
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
  formState.value = { ...defaultForm }
  getTables()
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
    ...routeParams,
  }
  const res = await getDasTablesGrantApi(params).catch(() => { })
  if (res) {
    pagination.total = res.total
    uiData.tableData = res.data
  }
  uiState.loading = false
}

// 获取指定schema的表
const getTables = async () => {
  const res = await getDasTablesListApi(routeParams).catch(() => { })
  uiData.tableList = res.data || []
}

// 提交表单
const onSubmit = useThrottleFn(async (data) => {
  const payload = {
    ...data,
    ...routeParams,
  }
  const res = await createDasTablesGrantApi(payload).catch(() => { })
  if (res) {
    message.success('操作成功')
    uiState.isModalOpen = false
    fetchData()
  }
})

// 删除
const handleDelete = useThrottleFn(async (record) => {
  const res = await deleteDasTablesGrantApi(record.id).catch(() => { })
  if (res) {
    message.info('操作成功')
    fetchData()
  }
})

// 初始化
onMounted(() => {
  fetchData()
})
</script>
