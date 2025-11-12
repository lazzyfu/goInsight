<template>
  <a-card size="small">
    <div class="search-wrapper">
      <a-space>
        <a-input-search
          v-model:value="searchValue"
          placeholder="输入搜索的内容"
          style="width: 350px"
          @search="onSearch"
        />
      </a-space>
    </div>
    <div style="margin-top: 12px">
      <a-table
        size="small"
        :columns="tableColumns"
        :row-key="(record) => record.key"
        :data-source="data.tableData"
        @resizeColumn="handleResizeColumn"
        :pagination="pagination"
        :loading="data.loading"
        @change="handleTableChange"
        :scroll="{ x: 1500 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'is_active'">
            <span v-if="record.is_active">
              <a-tag color="green">是</a-tag>
            </span>
            <span v-else>
              <a-tag color="volcano">否</a-tag>
            </span>
          </template>
          <template v-if="column.key === 'is_superuser'">
            <span v-if="record.is_superuser">
              <a-tag color="green">是</a-tag>
            </span>
            <span v-else>
              <a-tag color="volcano">否</a-tag>
            </span>
          </template>
          <template v-if="column.key === 'is_two_fa'">
            <span v-if="record.is_two_fa">
              <a-tag color="green">是</a-tag>
            </span>
            <span v-else>
              <a-tag color="volcano">否</a-tag>
            </span>
          </template>
          <template v-if="column.key === 'action'">
            <a-space wrap>
              <a-tooltip title="编辑">
                <a-button
                  type="link"
                  block
                  shape="circle"
                  :icon="h(EditOutlined)"
                  @click="editRecord(record)"
                />
              </a-tooltip>
              <a-tooltip title="修改密码">
                <a-button
                  type="link"
                  block
                  shape="circle"
                  :icon="h(CopyOutlined)"
                  @click="changePassword(record)"
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
  </a-card>
  <ChangePassword :open="isChangePasswordOpen" :uid="currentUid" @update:open="isChangePasswordOpen = $event" />
</template>

<script setup>
import { getUsersApi } from '@/api/admin'
import {
  CopyOutlined,
  DeleteOutlined,
  EditOutlined
} from '@ant-design/icons-vue'
import { h, onMounted, reactive, ref } from 'vue'
import ChangePassword from './ChangePassword.vue'

// 当前用户UID
const currentUid = ref()
const isChangePasswordOpen = ref(false)

// 搜索
const searchValue = ref('')
const onSearch = (value) => {
  searchValue.value = value
  pagination.current = 1
  fetchData()
}

// 表列
const tableColumns = ref([
  {
    title: '用户',
    dataIndex: 'username',
    key: 'username',
    scopedSlots: {
      customRender: 'username',
    },
  },
  {
    title: '昵称',
    dataIndex: 'nick_name',
    key: 'nick_name',
    scopedSlots: {
      customRender: 'nick_name',
    },
  },
  {
    title: '角色',
    dataIndex: 'role',
    key: 'role',
    scopedSlots: {
      customRender: 'role',
    },
  },
  {
    title: '激活',
    dataIndex: 'is_active',
    key: 'is_active',
    scopedSlots: {
      customRender: 'is_active',
    },
  },
  {
    title: '2FA认证',
    dataIndex: 'is_two_fa',
    key: 'is_two_fa',
    scopedSlots: {
      customRender: 'is_two_fa',
    },
  },
  {
    title: '管理员',
    dataIndex: 'is_superuser',
    key: 'is_superuser',
    scopedSlots: {
      customRender: 'is_superuser',
    },
  },
  {
    title: '邮箱',
    dataIndex: 'email',
    key: 'email',
    scopedSlots: {
      customRender: 'email',
    },
  },
  {
    title: '手机号',
    dataIndex: 'mobile',
    key: 'mobile',
    scopedSlots: {
      customRender: 'mobile',
    },
  },
  {
    title: '组织',
    dataIndex: 'organization',
    key: 'organization',
    scopedSlots: {
      customRender: 'organization',
    },
  },
  {
    title: '加入时间',
    dataIndex: 'date_joined',
    key: 'date_joined',
    scopedSlots: {
      customRender: 'date_joined',
    },
  },
  {
    title: '操作',
    dataIndex: 'action',
    key: 'action',
    scopedSlots: {
      customRender: 'action',
    },
  },
])

// 表数据
const data = reactive({
  tableData: [],
  loading: false,
})

// 分页
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  pageSizeOptions: ['10', '20', '50', '100'],
  showSizeChanger: true,
})

// 获取表数据
const fetchData = async () => {
  data.loading = true
  const params = {
    page_size: pagination.pageSize,
    page: pagination.current,
    is_page: true,
    search: searchValue.value,
  }
  const res = await getUsersApi(params).catch(()=> {})
  if (res) {
    pagination.total = res.total
    data.tableData = res.data
  }
  data.loading = false
}

// 翻页
const handleTableChange = (pager) => {
  pagination.current = pager.current
  pagination.pageSize = pager.pageSize
  fetchData()
}

function handleResizeColumn(w, col) {
  col.width = w
}

// 编辑用户信息
const editRecord = (record) => {
  console.log('record: ', record)
}

// 修改密码
const changePassword = (record) => {
  console.log('record: ', record)
  isChangePasswordOpen.value = true
  currentUid.value = record.uid
}

// 删除用户
const confirmDelete = (record) => {
  console.log('record: ', record)
}

onMounted(() => {
  fetchData()
})
</script>
