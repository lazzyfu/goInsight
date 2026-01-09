<template>
  <a-card title="配置用户库访问权限">
    <!-- 卡片右上角的新增按钮 -->
    <template #extra>
      <a-button type="primary" @click="handleAdd">
        <PlusOutlined />新增库访问权限
      </a-button>
    </template>

    <!-- 搜索区域 -->
    <div class="search-wrapper">
      <!-- 搜索 -->
      <a-input-search v-model:value="uiData.searchValue" placeholder="搜索..." style="width: 350px"
        @search="handleSearch" />
    </div>

    <!-- 表格 -->
    <div style="margin-top: 12px">
      <a-table size="middle" :columns="uiData.tableColumns" :row-key="(record) => record.key"
        :data-source="uiData.tableData" :pagination="pagination" :loading="uiState.loading" @change="handleTableChange"
        :scroll="{ x: 1100 }">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'schema'">
            <router-link :to="{
              name: 'view.admin.das.tables',
              params: {
                username: record.username,
                schema: record.schema,
                instance_id: record.instance_id,
              },
            }">
              {{ record.schema }}
            </router-link>
          </template>
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

  <!-- 新增弹窗 -->
  <DasSchemaFormModal :users="uiData.users" :environments="uiData.environments" :open="uiState.isModalOpen"
    title="新增库权限" @update:open="uiState.isModalOpen = $event" @submit="onSubmit" />
</template>

<script setup>
import {
    createDasSchemasGrantApi,
    deleteDasSchemasGrantApi,
    getDasSchemasListGrantApi,
    getEnvironmentsApi,
    getUsersApi,
} from '@/api/admin'
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { onMounted, reactive } from 'vue'
import DasSchemaFormModal from './DasSchemaFormModal.vue'

// 状态管理
const uiState = reactive({
  loading: false,
  isModalOpen: false,
})

// 数据
const uiData = reactive({
  searchValue: '',
  users: [],
  environments: [],
  tableData: [],
  tableColumns: [
    {
      title: '用户',
      dataIndex: 'username',
      key: 'username',
      fixed: 'left',
    },
    {
      title: '库名',
      dataIndex: 'schema',
      key: 'schema',
      fixed: 'left',
    },
    {
      title: '环境',
      dataIndex: 'environment',
      key: 'environment',
    },
    {
      title: '类型',
      dataIndex: 'db_type',
      key: 'db_type',
    },
    {
      title: '实例ID',
      dataIndex: 'instance_id',
      key: 'instance_id',
    },
    {
      title: '主机',
      dataIndex: 'hostname',
      key: 'hostname',
    },
    {
      title: '备注',
      dataIndex: 'remark',
      key: 'remark',
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

// 获取列表数据
const fetchData = async () => {
  uiState.loading = true
  const params = {
    page_size: pagination.pageSize,
    page: pagination.current,
    is_page: true,
    search: uiData.searchValue,
  }
  const res = await getDasSchemasListGrantApi(params).catch(() => { })
  if (res) {
    pagination.total = res.total
    uiData.tableData = res.data
  }
  uiState.loading = false
}

// 新增
const handleAdd = () => {
  getUsers()
  getEnvironments()
  uiState.isModalOpen = true
}

// 提交
const onSubmit = useThrottleFn(async (data) => {
  const res = await createDasSchemasGrantApi(data).catch(() => { })
  if (res) {
    message.success('操作成功')
    uiState.isModalOpen = false
    fetchData()
  }
})

// 删除
const handleDelete = useThrottleFn(async (record) => {
  const res = await deleteDasSchemasGrantApi(record.id).catch(() => { })
  if (res) {
    message.info('操作成功')
    fetchData()
  }
})

// 获取用户
const getUsers = async () => {
  const res = await getUsersApi().catch(() => { })
  uiData.users = res.data || []
}

// 获取环境
const getEnvironments = async () => {
  const res = await getEnvironmentsApi().catch(() => { })
  uiData.environments = res.data || []
}

// 初始化
onMounted(() => {
  fetchData()
})
</script>
