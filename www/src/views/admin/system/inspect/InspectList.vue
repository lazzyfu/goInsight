<template>
  <a-card title="审核参数配置">
    <!-- 搜索区域 -->
    <div class="search-wrapper">
      <!-- 搜索 -->
      <a-input-search
        v-model:value="uiData.searchValue"
        placeholder="搜索参数名..."
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
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action'">
            <a @click="handleEdit(record)"> <EditOutlined /> 编辑 </a>
          </template>
        </template>
        <template #expandedRowRender="{ record }">
          <p style="margin: 0">
            <highlightjs language="json" :code="JSON.stringify(record.params, null, 2)" />
          </p>
        </template>
      </a-table>
    </div>
  </a-card>

  <!-- 新增/编辑弹窗 -->
  <InspectFormModal
    :open="uiState.isModalOpen"
    title="配置审核参数"
    v-model:modelValue="formState"
    @update:open="uiState.isModalOpen = $event"
    @submit="onSubmit"
  />
</template>

<script setup>
import { getInspectParamsApi, updateInspectParamsApi } from '@/api/admin'
import { EditOutlined } from '@ant-design/icons-vue'
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref } from 'vue'
import InspectFormModal from './InspectFormModal.vue'

// 状态
const uiState = reactive({
  loading: false,
  isModalOpen: false,
})

// 数据
const uiData = reactive({
  searchValue: '',
  tableData: [],
  tableColumns: [
  {
    title: '描述',
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
  },
]
})

// form表单
const defaultForm = {
  params: '{}',
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
  const res = await getInspectParamsApi(params).catch(() => {})
  if (res) {
    pagination.total = res.total
    uiData.tableData = res.data
  }
  uiState.loading = false
}

// 编辑
const handleEdit = (record) => {
  formState.value = {
    ...record,
    params: JSON.stringify(record.params || {}, null, 2), // 仅修改 formState，原 record 保持 Object
  }
  uiState.isModalOpen = true
}

// 提交
const onSubmit = useThrottleFn(async (data) => {
  // 将 inspect_params 转换为 JSON 对象
  const payload = {
    ...data,
    params: JSON.parse(data.params),
  }
  const res = await updateInspectParamsApi(payload).catch(() => {})
  if (res) {
    message.success('操作成功')
    uiState.isModalOpen = false
    fetchData()
  }
})

// 初始化
onMounted(() => {
  fetchData()
})
</script>
