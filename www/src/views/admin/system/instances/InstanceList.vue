<template>
  <a-card title="数据库实例管理">
    <!-- 卡片右上角的新增按钮 -->
    <template #extra>
      <a-button type="primary" @click="handleAdd"> <PlusOutlined />新增数据库实例 </a-button>
    </template>

    <!-- 搜索区域 -->
    <div class="search-wrapper">
      <!-- 搜索 -->
      <a-input-search
        v-model:value="uiData.searchValue"
        placeholder="搜索实例..."
        style="width: 350px"
        @search="handleSearch"
      />
    </div>

    <!-- 表格 -->
    <div style="margin-top: 12px">
      <a-table
        size="small"
        :columns="uiData.tableColumns"
        :row-key="(record) => record.id"
        :data-source="uiData.tableData"
        :pagination="pagination"
        :loading="uiState.loading"
        @change="handleTableChange"
        :scroll="{ x: 1300 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'inspect_params'">
            <pre>{{ JSON.stringify(record.inspect_params, null, 2) }}</pre>
          </template>
          <template v-if="column.key === 'action'">
            <a-dropdown>
              <EllipsisOutlined />
              <template #overlay>
                <a-menu>
                  <a-menu-item key="1">
                    <a @click="handleEdit(record)"> <EditOutlined /> 编辑 </a>
                  </a-menu-item>
                  <a-menu-item key="2">
                    <a-popconfirm
                      title="确认删除吗？"
                      ok-text="是"
                      cancel-text="否"
                      @confirm="handleDelete(record)"
                    >
                      <a> <DeleteOutlined /> 删除 </a>
                    </a-popconfirm>
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </template>
        </template>
        <template #expandedRowRender="{ record }">
          <p style="margin: 0">
            <highlightjs language="json" :code="JSON.stringify(record.inspect_params, null, 2)" />
          </p>
        </template>
      </a-table>
    </div>
  </a-card>

  <!-- 新增/编辑弹窗 -->
  <InstanceFormModal
    :open="uiState.isModalOpen"
    v-model:modelValue="formState"
    :environments="uiData.environments"
    :organizations="uiData.environments"
    :title="uiState.isEditMode ? '编辑数据库实例' : '新增数据库实例'"
    @update:open="uiState.isModalOpen = $event"
    @submit="onSubmit"
  />
</template>

<script setup>
import {
  createDBConfigApi,
  deleteDBConfigApi,
  getDBConfigApi,
  getEnvironmentsApi,
  getOrganizationsApi,
  updateDBConfigApi,
} from '@/api/admin'
import { DeleteOutlined, EditOutlined, EllipsisOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref } from 'vue'
import InstanceFormModal from './InstanceFormModal.vue'

// 状态
const uiState = reactive({
  loading: false,
  isEditMode: false,
  isModalOpen: false,
})

const uiData = reactive({
  searchValue: '',
  environments: [],
  organizations: [],
  tableData: [],
  tableColumns: [
    {
      title: '描述',
      dataIndex: 'remark',
      key: 'remark',
      fixed: 'left',
    },

    {
      title: '主机名',
      dataIndex: 'hostname',
      key: 'hostname',
    },
    {
      title: '端口',
      dataIndex: 'port',
      key: 'port',
    },
    {
      title: '用途',
      dataIndex: 'use_type',
      key: 'use_type',
    },
    {
      title: '环境',
      dataIndex: 'environment_name',
      key: 'environment_name',
    },
    {
      title: '类型',
      dataIndex: 'db_type',
      key: 'db_type',
    },
    {
      title: '组织',
      dataIndex: 'organization_name',
      key: 'organization_name',
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
  environment: '',
  organization_key: '',
  db_type: '',
  use_type: '',
  hostname: '',
  port: 3306,
  inspect_params: '{}',
  remark: '',
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

// 获取列表数据
const fetchData = async () => {
  uiState.loading = true
  const params = {
    page_size: pagination.pageSize,
    page: pagination.current,
    is_page: true,
    search: uiData.searchValue,
  }
  const res = await getDBConfigApi(params).catch(() => {})
  if (res) {
    pagination.total = res.total
    uiData.tableData = res.data
  }
  uiState.loading = false
}

// 新增
const handleAdd = () => {
  uiState.isEditMode = false
  formState.value = { ...defaultForm }
  // 在打开 Modal 之前加载数据
  getEnvironments()
  getOrganizations()
  uiState.isModalOpen = true
}

// 编辑
const handleEdit = (record) => {
  uiState.isEditMode = true
  record.organization_key = record.organization_path

  formState.value = {
    ...record,
    inspect_params: JSON.stringify(record.inspect_params || {}, null, 2), // 仅修改 formState，原 record 保持 Object
  }

  // 在打开 Modal 之前加载数据
  getEnvironments()
  getOrganizations()
  uiState.isModalOpen = true
}

// 提交（新增或编辑）
const onSubmit = async (data) => {
  // 将 inspect_params 转换为 JSON 对象
  const payload = {
    ...data,
    inspect_params: JSON.parse(data.inspect_params),
  }
  const res = uiState.isEditMode
    ? await updateDBConfigApi(payload).catch(() => {})
    : await createDBConfigApi(payload).catch(() => {})
  if (res?.code === '0000') {
    message.success('操作成功')
    uiState.isModalOpen = false
    fetchData()
  }
}

// 删除
const handleDelete = async (record) => {
  const res = await deleteDBConfigApi(record.id).catch(() => {})
  if (res?.code === '0000') {
    message.info('操作成功')
    fetchData()
  }
}

// 获取环境
const getEnvironments = async () => {
  const res = await getEnvironmentsApi().catch(() => {})
  uiData.environments = res?.data || []
}

// 获取组织
const getOrganizations = async () => {
  const res = await getOrganizationsApi().catch(() => {})
  uiData.environments = res?.data || []
}

// 初始化
onMounted(() => {
  fetchData()
})
</script>
