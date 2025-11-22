<template>
  <a-card title="角色管理">
    <!-- 搜索区域 -->
    <div class="search-wrapper">
      <!-- 搜索 -->
      <a-input-search
        v-model:value="searchValue"
        placeholder="搜索参数名..."
        style="width: 350px"
        @search="handleSearch"
      />
    </div>
    <!-- 表格 -->
    <div style="margin-top: 12px">
      <a-table
        size="small"
        :columns="tableColumns"
        :row-key="(record) => record.id"
        :data-source="tableData"
        :pagination="pagination"
        :loading="state.loading"
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
    :open="state.isModalOpen"
    title="配置审核参数"
    v-model:modelValue="formState"
    @update:open="state.isModalOpen = $event"
    @submit="onSubmit"
  />
</template>

<script setup>
import { getInspectParamsApi, updateInspectParamsApi } from '@/api/admin'
import { EditOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref } from 'vue'
import InspectFormModal from './InspectFormModal.vue'

// TODO 修复JSON重复转义BUG

// 状态管理
const state = reactive({
  loading: false,
  isModalOpen: false,
})
const searchValue = ref('')
const defaultForm = {
  params: '{}',
  remark: '',
}
const formState = ref({ ...defaultForm })

// 表
const tableData = ref([])
const tableColumns = [
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

// 获取表数据
const fetchData = async () => {
  state.loading = true
  const params = {
    page_size: pagination.pageSize,
    page: pagination.current,
    is_page: true,
    search: searchValue.value,
  }
  const res = await getInspectParamsApi(params)
  if (res) {
    pagination.total = res.total
    tableData.value = res.data
  }
  state.loading = false
}

// 编辑记录
const handleEdit = (record) => {
  formState.value = {
    ...record,
    params: JSON.stringify(record.params || {}, null, 2), // 仅修改 formState，原 record 保持 Object
  }
  state.isModalOpen = true
}

// 提交表单
const onSubmit = async (data) => {
  // 将 inspect_params 转换为 JSON 对象
  const payload = {
    ...data,
    params: JSON.parse(data.params),
  }
  const res = await updateInspectParamsApi(payload)
  if (res?.code === '0000') {
    message.success('操作成功')
    state.isModalOpen = false
    fetchData()
  }
}

// 生命周期
onMounted(() => {
  fetchData()
})
</script>
