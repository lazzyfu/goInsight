<template>
  <div>
    <div class="search-wrapper">
      <a-input-search
        v-model:value="uiData.searchValue"
        placeholder="输入SQL内容"
        style="width: 350px"
        @search="handleSearch"
      />
    </div>
    <div style="margin-top: 14px">
      <a-table
        size="small"
        :columns="uiData.tableColumns"
        :row-key="(record) => record.key"
        :data-source="uiData.tableData"
        :pagination="pagination"
        :loading="uiState.loading"
        @change="handleTableChange"
        :scroll="{ x: 1100 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'sqltext'">
            <a @click="showSqlDetail(record)" title="查看完整SQL">
              <EyeOutlined />
            </a>
            {{ record.sqltext }}
          </template>
          <template v-else-if="column.key === 'updated_at'">
            <div class="time-cell">
              <div class="time-main">{{ formatDate(record.updated_at) }}</div>
              <div class="time-sub">{{ formatTime(record.updated_at) }}</div>
            </div>
          </template>
          <template v-else-if="column.key === 'created_at'">
            <div class="time-cell">
              <div class="time-main">{{ formatDate(record.created_at) }}</div>
              <div class="time-sub">{{ formatTime(record.created_at) }}</div>
            </div>
          </template>
          <template v-if="column.key === 'action'">
            <a-dropdown>
              <EllipsisOutlined />
              <template #overlay>
                <a-menu>
                  <a-menu-item key="1">
                    <a @click="copyRecord(record)"> <CopyOutlined /> 拷贝 </a>
                  </a-menu-item>
                  <a-menu-item key="2">
                    <a @click="handleEdit(record)"> <EditOutlined /> 编辑 </a>
                  </a-menu-item>
                  <a-menu-item key="3">
                    <a-popconfirm
                      title="确认删除吗？"
                      ok-text="是"
                      cancel-text="否"
                      @confirm="handleDelete(record)"
                    >
                      <a><DeleteOutlined /> 删除</a>
                    </a-popconfirm>
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </template>
        </template>
      </a-table>
    </div>
    <!-- 查看SQL -->
    <a-modal
      v-model:open="uiState.open"
      title="SQL语句"
      width="55%"
      :footer="null"
      @cancel="handleCancel"
    >
      <highlightjs language="sql" :code="uiData.sqltext" />
    </a-modal>
    <!-- 更新收藏SQL -->
    <DasFavoriteFormModal
      :open="uiState.isFavoritesOpen"
      v-model:modelValue="formState"
      @update:open="uiState.isFavoritesOpen = $event"
      @submit="onSubmit"
    />
  </div>
</template>

<script setup>
import { DeleteFavoritesApi, GetFavoritesApi, UpdateFavoritesApi } from '@/api/das'
import DasFavoriteFormModal from '@/views/das/favorite/DasFavoriteFormModal.vue'
import {
  CopyOutlined,
  DeleteOutlined,
  EditOutlined,
  EllipsisOutlined,
  EyeOutlined,
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref } from 'vue'
import useClipboard from 'vue-clipboard3'

const { toClipboard } = useClipboard()

// 状态
const uiState = reactive({
  loading: false,
  open: false,
  isFavoritesOpen: false,
})

// 数据
const uiData = reactive({
  searchValue: '',
  sqltext: '',
  tableData: [],
  tableColumns: [
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      fixed: 'left',
      ellipsis: true,
    },
    {
      title: 'SQL语句',
      dataIndex: 'sqltext',
      key: 'sqltext',
      ellipsis: true,
      width: '35%',
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      key: 'updated_at',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
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
  title: '',
  sqltext: '',
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
  const res = await GetFavoritesApi(params).catch(() => {})
  if (res) {
    pagination.total = res.total
    uiData.tableData = res.data
  }
  uiState.loading = false
}

// 查看SQL
const showSqlDetail = (record) => {
  uiState.open = true
  uiData.sqltext = record.sqltext
}

// 关闭查看SQL弹窗
const handleCancel = (e) => {
  uiState.open = false
}

// 格式日期
const formatDate = (dateStr) => {
  return dateStr.split(' ')[0]
}

// 格式时间
const formatTime = (dateStr) => {
  return dateStr.split(' ')[1]
}

// 复制SQL
const copyRecord = async (value) => {
  try {
    await toClipboard(value.sqltext)
    message.success('已拷贝到剪贴板')
  } catch (e) {
    message.error(e)
  }
}

// 编辑
const handleEdit = (record) => {
  formState.value = { ...record }
  uiState.isFavoritesOpen = true
}

// 提交
const onSubmit = async (data) => {
  const res = await UpdateFavoritesApi(data).catch(() => {})
  if (res) {
    message.success('更新成功')
    uiState.isFavoritesOpen = false
    fetchData()
  }
}

// 删除
const handleDelete = async (record) => {
  const res = await DeleteFavoritesApi(record).catch(() => {})
  if (res) {
    message.info('操作成功')
    fetchData()
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.user-cell {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.user-avatar {
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 0.875rem;
}

.schema-tag {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  background: #f0fdf4;
  color: #166534;
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  font-weight: 500;
}

.duration-tag {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  font-weight: 500;
}

.duration-tag.fast {
  background: #f0fdf4;
  color: #166534;
}

.duration-tag.medium {
  background: #fffbeb;
  color: #92400e;
}

.duration-tag.slow {
  background: #fef2f2;
  color: #991b1b;
}

.time-cell {
  font-size: 0.875rem;
}

.time-main {
  font-weight: 500;
  color: #1e293b;
}

.time-sub {
  color: #6b7280;
  font-size: 0.75rem;
}

.status-tag {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  font-weight: 500;
}

.status-tag.success {
  background: #f0fdf4;
  color: #166534;
}

.status-tag.error {
  background: #fef2f2;
  color: #991b1b;
}

.error-message {
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 0.375rem;
  padding: 0.75rem;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 0.875rem;
  color: #991b1b;
}
</style>
