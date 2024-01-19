<template>
  <a-card>
    <!-- 搜索框 -->
    <a-form layout="inline" :form="form" @keyup.enter.native="handleSearch">
      <a-form-item>
        <a-input-search allowClear style="width: 300px" placeholder="输入要查询的库名或表名" v-decorator="['search']" />
      </a-form-item>
      <a-form-item style="float: right; margin-right: 0px">
        <a-button type="primary" @click="onAdd"> <a-icon type="plus" />新增</a-button>
      </a-form-item>
      <a-form-item>
        <span class="table-page-search-submitButtons">
          <a-button type="primary" @click="handleSearch">查询</a-button>
        </span>
      </a-form-item>
    </a-form>

    <a-table
      :columns="tableColumns"
      :rowKey="(record, index) => index"
      :dataSource="tableData"
      :pagination="pagination"
      :loading="loading"
      @change="handleTableChange"
      size="middle"
    >
      <span slot="expandedRowRender" slot-scope="record" style="margin: 0">
        <pre class="codeStyle"><code>{{ record.sqltext }}</code></pre>
      </span>
      <span slot="action" slot-scope="text, record">
        <div class="editable-row-operations">
          <a type="button" v-clipboard:copy="record.sqltext" v-clipboard:success="onCopy" v-clipboard:error="onError">
            <span style="color: #409eff">拷贝</span></a
          >
          <a-divider type="vertical" />
          <a @click="() => editRow(record)"><span style="color: #409eff">编辑</span></a>
          <a-divider type="vertical" />
          <a type="dashed" @click="DeleteConfirm(record.id)"><span style="color: #409eff">删除</span></a>
        </div>
      </span>
    </a-table>
    <FavoriteAddComponent ref="FavoriteAddComponent" @refreshTable="fetchData"></FavoriteAddComponent>
    <FavoriteEditComponent ref="FavoriteEditComponent" @refreshTable="fetchData"></FavoriteEditComponent>
  </a-card>
</template>

<script>
import { getFavoritesApi, deleteFavoritesApi } from '@/api/das'

import FavoriteAddComponent from './FavoriteAdd'
import FavoriteEditComponent from './FavoriteEdit'

const tableColumns = [
  {
    title: '标题',
    dataIndex: 'title',
    key: 'title',
    scopedSlots: {
      customRender: 'title',
    },
  },
  {
    title: '用户名',
    dataIndex: 'username',
    key: 'username',
    scopedSlots: {
      customRender: 'username',
    },
  },
  {
    title: 'SQL',
    dataIndex: 'sqltext',
    key: 'sqltext',
    scopedSlots: {
      customRender: 'sqltext',
    },
    ellipsis: true,
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
    title: '更新时间',
    dataIndex: 'updated_at',
    key: 'updated_at',
    scopedSlots: {
      customRender: 'updated_at',
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
]

export default {
  props: {
    visible: Boolean, // 接收父组件传值
  },
  components: { FavoriteAddComponent, FavoriteEditComponent },
  data() {
    return {
      loading: false,
      tableColumns,
      tableData: [],
      pagination: {
        current: 1,
        pageSize: 10,
        total: 0,
        pageSizeOptions: ['10', '20'],
        showSizeChanger: true,
      },
      form: this.$form.createForm(this),
    }
  },
  methods: {
    onAdd() {
      this.$refs.FavoriteAddComponent.showModal()
    },
    editRow(row) {
      this.$refs.FavoriteEditComponent.showModal(row)
    },
    DeleteConfirm(id) {
      const _this = this
      this.$confirm({
        title: '警告',
        content: '你确定删除？',
        okText: 'Yes',
        okType: 'danger',
        cancelText: 'No',
        onOk() {
          deleteFavoritesApi(id)
            .then((res) => {
              const messageType = res.code === '0000' ? 'info' : 'error'
              _this.$message[messageType](res.message)
            })
            .finally(() => {
              _this.fetchData()
            })
        },
        onCancel() {},
      })
    },
    handleTableChange(pager) {
      this.pagination.current = pager.current
      this.pagination.pageSize = pager.pageSize
      this.fetchData()
    },
    // 加载收藏的SQL
    fetchData() {
      this.loading = true
      const params = {
        page_size: this.pagination.pageSize,
        page: this.pagination.current,
        is_page: true,
        ...this.filters,
      }
      getFavoritesApi(params)
        .then((res) => {
          this.pagination.total = res.total
          this.tableData = res.data
        })
        .finally(() => {
          this.loading = false
        })
    },
    // 搜索
    handleSearch(e) {
      e.preventDefault()
      this.form.validateFields((error, values) => {
        if (error) {
          return
        }
        this.filters = {
          search: values['search'],
        }
        this.pagination.current = 1
        this.fetchData()
      })
    },
    // 拷贝
    onCopy: function (e) {
      this.$message.info('拷贝成功')
    },
    onError: function (e) {
      this.$message.error('拷贝失败')
    },
  },
  mounted() {
    this.fetchData()
  },
}
</script>

<style lang='less' scoped>
::v-deep .ant-table {
  font-size: 12px;
}
::v-deep .ant-pagination {
  font-size: 12px;
}
::v-deep .ant-select-sm .ant-select-selection__rendered {
  font-size: 12px;
}
::v-deep .codeStyle {
  tab-size: 4;
  background: #183055;
  color: #e6ecf1;
  padding: 12px 12px 12px 8px;
  direction: ltr;
  text-align: left;
  border: 1px solid #d1d1d1;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  max-width: 100%;
  overflow: auto;
  word-break: normal;
  white-space: pre;
  white-space: pre-wrap;
  word-wrap: break-word;
}

::v-deep .ant-card-body {
  padding: 8px;
}

::v-deep .ant-form {
  margin-bottom: 8px;
}
</style>
