<template>
  <a-card>
    <a-form layout="inline" :form="form" @keyup.enter.native="handleSearch">
      <a-form-item>
        <a-input-search
          allowClear
          style="width: 300px; margin-bottom: 10px"
          placeholder="输入要查询的库名或表名"
          v-decorator="['search']"
        />
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
      <span slot="expandedRowRender" slot-scope="record">
        <ul>
          <li>
            原始SQL
            <pre class="codeStyle"><code>{{ record.sqltext }}</code></pre>
          </li>
          <li>
            重写SQL
            <pre class="codeStyle"><code>{{ record.rewrite_sqltext }}</code></pre>
          </li>
          <li>
            错误
            <pre class="codeStyle"><code>{{ record.error_msg }}</code></pre>
          </li>
        </ul>
      </span>
      <span slot="action" slot-scope="text, record">
        <div class="editable-row-operations">
          <a type="button" v-clipboard:copy="record.sqltext" v-clipboard:success="onCopy" v-clipboard:error="onError">
            <span style="color: #409eff">原始SQL</span>
          </a>
        </div>
      </span>
    </a-table>
  </a-card>
</template>

<script>
import { getHistory } from '@/api/das'

const tableColumns = [
  {
    title: '用户名',
    dataIndex: 'username',
    key: 'username',
    scopedSlots: {
      customRender: 'username',
    },
  },
  {
    title: '实例ID',
    dataIndex: 'instance_id',
    key: 'instance_id',
    scopedSlots: {
      customRender: 'instance_id',
    },
    ellipsis: true,
  },
  {
    title: '库',
    dataIndex: 'schema',
    key: 'schema',
    scopedSlots: {
      customRender: 'schema',
    },
  },
  {
    title: '表',
    dataIndex: 'tables',
    key: 'tables',
    scopedSlots: {
      customRender: 'tables',
    },
  },
  {
    title: '原始SQL',
    dataIndex: 'sqltext',
    key: 'sqltext',
    scopedSlots: {
      customRender: 'sqltext',
    },
    ellipsis: true,
  },
  {
    title: '重写SQL',
    dataIndex: 'rewrite_sqltext',
    key: 'rewrite_sqltext',
    scopedSlots: {
      customRender: 'rewrite_sqltext',
    },
    ellipsis: true,
  },
  {
    title: '参数',
    dataIndex: 'params',
    key: 'params',
    scopedSlots: {
      customRender: 'params',
    },
  },
  {
    title: '返回行数',
    dataIndex: 'return_rows',
    key: 'return_rows',
    scopedSlots: {
      customRender: 'return_rows',
    },
  },
  {
    title: '耗时(ms)',
    dataIndex: 'duration',
    key: 'duration',
    scopedSlots: {
      customRender: 'duration',
    },
  },
  {
    title: '错误',
    dataIndex: 'error_msg',
    key: 'error_msg',
    scopedSlots: {
      customRender: 'error_msg',
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
    title: '拷贝',
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
    // 子组件不能修改父组件传递的visible，因此需要通过this.$emit发射给父组件
    // close自定义即可，父组件需要对应即可
    onClose() {
      this.$emit('close')
    },
    handleTableChange(pager) {
      this.pagination.current = pager.current
      this.pagination.pageSize = pager.pageSize
      this.fetchData()
    },
    fetchData() {
      this.loading = true
      const params = {
        page_size: this.pagination.pageSize,
        page: this.pagination.current,
        is_page: true,
        ...this.filters,
      }
      getHistory(params)
        .then((response) => {
          this.pagination.total = response.total
          this.tableData = response.data
        })
        .finally(() => {
          this.loading = false
        })
    },
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
    resetForm() {
      this.form.resetFields()
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
  padding: 4px;
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

li {
  list-style-type: circle;
}

::v-deep .ant-card-body {
  padding: 8px;
}

::v-deep .ant-form {
  margin-bottom: 8px;
}
</style>