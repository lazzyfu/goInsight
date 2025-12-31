<template>
  <a-card title="全局审核参数配置">
    <!-- 搜索区域 -->
    <div class="search-wrapper">
      <!-- 搜索 -->
      <a-input-search v-model:value="uiData.searchValue" placeholder="搜索参数名..." style="width: 350px"
        @search="handleSearch" />
    </div>

    <!-- 表格 -->
    <div style="margin-top: 12px">
      <a-table size="small" :columns="uiData.tableColumns" :row-key="(record) => record.id" :scroll="{ x: 1100 }"
        :data-source="uiData.tableData" :pagination="pagination" :loading="uiState.loading" @change="handleTableChange">
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
          </template>
        </template>
      </a-table>
    </div>
  </a-card>

  <!-- 新增/编辑弹窗 -->
  <InspectParamsFormModal :open="uiState.isModalOpen" title="配置审核参数" v-model:modelValue="formState"
    @update:open="uiState.isModalOpen = $event" @submit="onSubmit" />
</template>

<script setup>
import { getInspectParamsApi, updateInspectParamsApi } from '@/api/admin'
import { EditOutlined } from '@ant-design/icons-vue'
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref } from 'vue'
import InspectParamsFormModal from './InspectParamsFormModal.vue'

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
  title: '',
  key: '',
  type: 'string',
  value: '',
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
  const res = await getInspectParamsApi(params).catch(() => { })
  if (res) {
    pagination.total = res.total
    uiData.tableData = res.data
  }
  uiState.loading = false
}

// 编辑
const handleEdit = (record) => {
  formState.value = {
    ...defaultForm,
    ...record,
  }
  uiState.isModalOpen = true
}

// 提交
const onSubmit = useThrottleFn(async (data) => {
  const payload = { ...data }
  // 移除弹窗内部使用的临时字段
  delete payload._editValue

  const res = await updateInspectParamsApi(payload).catch(() => { })
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
