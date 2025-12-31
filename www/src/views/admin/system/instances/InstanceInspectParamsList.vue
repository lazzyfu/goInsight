<template>
  <a-card title="实例自定义审核参数（优先级大于全局审核参数）">
    <!-- 卡片右上角的新增按钮 -->
    <template #extra>
      <a-button type="primary" @click="handleAdd">
        <PlusOutlined />新增
      </a-button>
    </template>

    <!-- 搜索区域 -->
    <div class="search-wrapper">
      <!-- 搜索 -->
      <a-input-search v-model:value="uiData.searchValue" placeholder="搜索审核参数名..." style="width: 350px"
        @search="handleSearch" />
    </div>

    <!-- 表格 -->
    <div style="margin-top: 12px">
      <a-table size="small" :columns="uiData.tableColumns" :row-key="(record) => record.key"
        :data-source="uiData.tableData" :pagination="pagination" :loading="uiState.loading" @change="handleTableChange"
        :scroll="{ x: 1100 }">
        <template #bodyCell="{ column, record }">
          <!-- 类型转换为中文显示 -->
          <template v-if="column.key === 'type'">
            <template v-if="record.type === 'string'">字符串</template>
            <template v-else-if="record.type === 'number'">数字</template>
            <template v-else-if="record.type === 'boolean'">布尔值</template>
            <template v-else>{{ record.type }}</template>
          </template>
          <template v-if="column.key === 'action'">
            <a @click="handleEdit(record)">
              <EditOutlined /> 编辑
            </a>
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

  <!-- 编辑弹窗 -->
  <InstanceInspectParamsFormModal
    :open="uiState.isModalOpen"
    v-model:modelValue="formState"
    :inspect-params="uiData.inspectParamsList"
    :title="uiState.isEditMode ? '编辑审核参数' : '新增审核参数'"
    @update:open="uiState.isModalOpen = $event"
    @submit="onSubmit"
  />
</template>

<script setup>
import {
  createInstanceInspectParamsApi,
  deleteInstanceInspectParamsApi,
  getInspectParamsApi,
  getInstanceInspectParamsApi,
  updateInstanceInspectParamsApi,
} from '@/api/admin'
import { DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import InstanceInspectParamsFormModal from './InstanceInspectParamsFormModal.vue'

// 路由参数
const route = useRoute()
const routeParams = route.params

// 状态管理
const uiState = reactive({
  loading: false,
  isModalOpen: false,
  isEditMode: false,
})

// 数据
const uiData = reactive({
  searchValue: '',
  inspectParamsList: [],
  tableData: [],
  tableColumns: [
    {
      title: '描述',
      dataIndex: 'title',
      key: 'title',
      fixed: 'left',
      width: '40%',
      ellipsis: true,
    },
    {
      title: '值',
      dataIndex: 'value',
      width: '20%',
      key: 'value',
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      width: '20%',
      key: 'updated_at',
    },
    {
      title: '操作',
      dataIndex: 'action',
      key: 'action',
    },
  ]
})

// form表单
const defaultForm = {
  id: undefined,
  title: '',
  key: '',
  type: '',
  value: '',
  _editValue: '',
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
  getInspectParams()
  uiState.isModalOpen = true
}

// 编辑
const handleEdit = (record) => {
  uiState.isEditMode = true
  formState.value = {
    ...defaultForm,
    ...record,
  }
  uiState.isModalOpen = true
}

// 获取实例的审核参数
const fetchData = async () => {
  uiState.loading = true
  const params = {
    page_size: pagination.pageSize,
    page: pagination.current,
    is_page: true,
    search: uiData.searchValue,
    ...routeParams,
  }
  const res = await getInstanceInspectParamsApi(params).catch(() => { })
  if (res) {
    pagination.total = res.total
    uiData.tableData = res.data
  }
  uiState.loading = false
}

// 获取所有的审核参数
const getInspectParams = async () => {
  const res = await getInspectParamsApi(routeParams).catch(() => { })
  uiData.inspectParamsList = res.data || []
}

// 提交表单
const onSubmit = useThrottleFn(async (data) => {
  const payload = {
    ...data,
    ...routeParams,
  }
  const res = uiState.isEditMode
    ? await updateInstanceInspectParamsApi(payload).catch(() => { })
    : await createInstanceInspectParamsApi(payload).catch(() => { })
  if (res) {
    message.success('操作成功')
    uiState.isModalOpen = false
    fetchData()
  }
})

// 删除
const handleDelete = useThrottleFn(async (record) => {
  const res = await deleteInstanceInspectParamsApi(record.id).catch(() => { })
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
