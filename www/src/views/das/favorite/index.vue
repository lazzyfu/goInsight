<template>
  <div>
    <div class="search-wrapper">
      <a-input-search
        v-model:value="searchValue"
        placeholder="输入SQL内容"
        style="width: 350px"
        @search="onSearch"
      />
    </div>
    <div style="margin-top: 14px">
      <a-table
        size="small"
        :columns="columns"
        :row-key="(record) => record.key"
        :data-source="data.dataSource"
        @resizeColumn="handleResizeColumn"
        :pagination="pagination"
        :loading="data.loading"
        @change="handleTableChange"
        :scroll="{ x: 1500 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'username'">
            <div class="user-cell">
              <div class="user-avatar" :style="{ backgroundColor: getUserColor(record.username) }">
                {{ record.username.charAt(0).toUpperCase() }}
              </div>
              <span class="user-name">{{ record.username }}</span>
            </div>
          </template>
          <template v-else-if="column.key === 'sqltext'">
            <div class="sql-cell">
              <code class="sql-preview">{{ record.sqltext }}</code>
              <button @click="showSqlDetail(record)" class="view-btn" title="查看完整SQL">
                <EyeOutlined />
              </button>
            </div>
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
          <template v-else-if="column.key === 'action'">
            <a-space wrap>
              <a-tooltip title="拷贝SQL语句">
                <a-button
                  type="link"
                  block
                  shape="circle"
                  :icon="h(CopyOutlined)"
                  @click="copyRecord(record.sqltext)"
                />
              </a-tooltip>
              <a-tooltip title="编辑">
                <a-button
                  type="link"
                  block
                  shape="circle"
                  :icon="h(EditOutlined)"
                  @click="editRecord(record)"
                />
              </a-tooltip>
              <a-popconfirm
                title="确认删除吗？"
                ok-text="是"
                cancel-text="否"
                @confirm="confirmDelete(record)"
              >
                <a-button block type="link" shape="circle" :icon="h(DeleteOutlined)"> </a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </div>
    <!-- 查看SQL -->
    <a-modal v-model:open="modal.open" title="SQL语句" width="55%" :footer="null" @ok="handleOk">
      <highlightjs language="sql" :code="modal.sqltext" />
    </a-modal>
    <!-- 更新收藏SQL -->
    <FavoritesEdit
      :open="isFavoritesOpen"
      :formState="favoritesFormState"
      :btnType="favoritesBtnType"
      @update:open="isFavoritesOpen = $event"
      @submit="handleFavoritesSubmit"
    />
  </div>
</template>

<script setup>
import { DeleteFavoritesApi, GetFavoritesApi, UpdateFavoritesApi } from '@/api/das'
import FavoritesEdit from '@/views/das/favorite/modal.vue'
import { CopyOutlined, DeleteOutlined, EditOutlined, EyeOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { h, onMounted, reactive, ref } from 'vue'
import useClipboard from 'vue-clipboard3'

const isFavoritesOpen = ref(false)
const favoritesBtnType = ref('更新')
const favoritesFormState = ref({
  id: '',
  title: '',
  sqltext: '',
})

// 复制
const { toClipboard } = useClipboard()

// 搜索
const searchValue = ref('')
const onSearch = (value) => {
  searchValue.value = value
  pagination.current = 1
  fetchData()
}

// table
const columns = ref([
  {
    title: '用户名',
    dataIndex: 'username',
    key: 'username',
    scopedSlots: {
      customRender: 'username',
    },
    fixed: 'left',
  },
  {
    title: '标题',
    dataIndex: 'title',
    key: 'title',
    scopedSlots: {
      customRender: 'title',
    },
    ellipsis: true,
    width: '20%',
  },
  {
    title: 'SQL语句',
    dataIndex: 'sqltext',
    key: 'sqltext',
    scopedSlots: {
      customRender: 'sqltext',
    },
    width: '35%',
  },
  {
    title: '更新时间',
    dataIndex: 'updated_at',
    key: 'updated_at',
    scopedSlots: {
      customRender: 'updated_at',
    },
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    key: 'created_at',
    scopedSlots: {
      customRender: 'created_at',
    },
  },
  {
    title: '操作',
    dataIndex: 'action',
    key: 'action',
    scopedSlots: {
      customRender: 'action',
    },
    fixed: 'right',
  },
])

const modal = reactive({
  open: false,
  sqltext: '',
})

const data = reactive({
  dataSource: [],
  loading: false,
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  pageSizeOptions: ['10', '20', '50', '100'],
  showSizeChanger: true,
})

const fetchData = () => {
  data.loading = true
  const params = {
    page_size: pagination.pageSize,
    page: pagination.current,
    is_page: true,
    search: searchValue.value,
  }
  GetFavoritesApi(params).then((res) => {
    if (res.code === '0000') {
      pagination.total = res.total
      data.dataSource = res.data
    } else {
      message.error(res.message)
    }
  })
  data.loading = false
}

const handleTableChange = (pager) => {
  pagination.current = pager.current
  pagination.pageSize = pager.pageSize
  fetchData()
}

const showSqlDetail = (record) => {
  modal.open = true
  modal.sqltext = record.sqltext
}

const handleOk = (e) => {
  modal.open = false
}

const formatDate = (dateStr) => {
  return dateStr.split(' ')[0]
}

const formatTime = (dateStr) => {
  return dateStr.split(' ')[1]
}

const getUserColor = (username) => {
  const colors = ['#3b82f6', '#8b5cf6', '#f59e0b', '#10b981', '#ef4444', '#6366f1']
  let hash = 0
  for (let i = 0; i < username.length; i++) {
    hash = username.charCodeAt(i) + ((hash << 5) - hash)
  }
  return colors[Math.abs(hash) % colors.length]
}

function handleResizeColumn(w, col) {
  col.width = w
}

const copyRecord = async (value) => {
  try {
    await toClipboard(value)
    message.success('已拷贝到剪贴板')
  } catch (e) {
    message.error(e)
  }
}

const editRecord = (record) => {
  isFavoritesOpen.value = true
  favoritesFormState.value = { ...record }
}

const handleFavoritesSubmit = (data) => {
  UpdateFavoritesApi(data).then((res) => {
    if (res.code == '0000') {
      message.success('更新成功')
    } else {
      message.error(res.message)
    }
    isFavoritesOpen.value = false
    fetchData()
  })
}

const confirmDelete = (val) => {
  DeleteFavoritesApi(val).then((res) => {
    if (res.code === '0000') {
      message.success('删除成功')
      fetchData()
    } else {
      message.error(res.message)
    }
  })
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

.user-name {
  font-weight: 500;
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

.sql-cell {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.sql-preview {
  background: #f8fafc;
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 0.75rem;
  color: #1e293b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.view-btn {
  background: none;
  border: none;
  color: #3b82f6;
  cursor: pointer;
  padding: 0.25rem;
  border-radius: 0.25rem;
  transition: all 0.2s;
}

.view-btn:hover {
  background: #eff6ff;
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
