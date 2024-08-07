<template>
  <a-card>
    <!-- 搜索框 -->
    <a-form layout="inline" :form="form" @keyup.enter.native="handleSearch">
      <a-form-item>
        <a-input-search allowClear style="width: 250px" placeholder="输入要查询内容" v-decorator="['search']" />
      </a-form-item>
      <a-form-item>
        <span class="table-page-search-submitButtons">
          <a-button type="primary" @click="handleSearch">查询</a-button>
        </span>
      </a-form-item>
    </a-form>
    <!-- 表格 -->
    <a-table
      :columns="tableColumns"
      :rowKey="(record, index) => index"
      :dataSource="tableData"
      :pagination="pagination"
      :loading="loading"
      @change="handleTableChange"
      size="middle"
    >
      <span slot="action" slot-scope="text, record">
        <div class="editable-row-operations">
          <a @click="() => editRow(record)"><span style="color: #409eff">编辑</span></a>
        </div>
      </span>
      <span slot="expandedRowRender" slot-scope="record">
        <pre>{{ JSON.stringify(record.params, null, 2) }}</pre>
      </span>
    </a-table>
    <InspectEditComponent ref="InspectEditComponent" @refreshTable="GetInspectParams"></InspectEditComponent>
  </a-card>
</template>

<script>
import { adminGetInspectParamsApi } from '@/api/inspect';

import InspectEditComponent from './InspectEdit';

const tableColumns = [
  {
    title: '备注',
    dataIndex: 'remark',
    key: 'remark',
    scopedSlots: {
      customRender: 'remark',
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
  components: {
    InspectEditComponent,
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
      form: this.$form.createForm(this, { name: 'inspect' }),
    }
  },
  methods: {
    GetInspectParams() {
      this.loading = true
      const params = {
        page_size: this.pagination.pageSize,
        page: this.pagination.current,
        is_page: true,
        ...this.filters,
      }
      adminGetInspectParamsApi(params)
        .then((res) => {
          if (res.code === '0000') {
            this.pagination.total = res.total
            this.tableData = res.data
          } else {
            this.$message.error(res.message)
          }
        })
        .catch((_error) => {})
        .finally(() => {
          this.loading = false
        })
    },
    editRow(row) {
      this.$refs.InspectEditComponent.showModal(row)
    },
    handleTableChange(pager) {
      this.pagination.current = pager.current
      this.pagination.pageSize = pager.pageSize
      this.GetInspectParams()
    },
    handleSearch(e) {
      e.preventDefault()
      this.form.validateFields((error, values) => {
        if (error) {
          return
        }
        this.filters = {
          environment: values['environment'],
          db_type: values['db_type'],
          search: values['search'],
        }
        this.pagination.current = 1
        this.GetInspectParams()
      })
    },
  },
  mounted() {
    this.GetInspectParams()
  },
}
</script>

<style lang="less" scoped>
::v-deep .ant-card-body {
  padding: 8px;
}

::v-deep .ant-form {
  margin-bottom: 8px;
}
</style>