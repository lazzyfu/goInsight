<template>
  <a-card>
    <!-- 搜索框 -->
    <a-form layout="inline" :form="form" @keyup.enter.native="handleSearch">
      <a-form-item>
        <a-input-search allowClear style="width: 300px" placeholder="输入要查询内容" v-decorator="['search']" />
      </a-form-item>
      <a-form-item>
        <span class="table-page-search-submitButtons">
          <a-button type="primary" @click="handleSearch">查询</a-button>
        </span>
      </a-form-item>
      <a-form-item style="float: right; margin-right: 0px">
        <a-button type="primary" @click="onAdd"> <a-icon type="plus" /> 新增 </a-button>
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
          <a-divider type="vertical" />
          <a type="dashed" @click="DeleteConfirm(record.id)"><span style="color: #409eff">删除</span></a>
        </div>
      </span>
    </a-table>
    <RolesAddComponent ref="RolesAddComponent" @refreshTable="getRoles"></RolesAddComponent>
    <RolesEditComponent ref="RolesEditComponent" @refreshTable="getRoles"></RolesEditComponent>
  </a-card>
</template>

<script>
import { getRolesApi, deleteRolesApi } from '@/api/users'

import RolesAddComponent from './RolesAdd'
import RolesEditComponent from './RolesEdit'

const tableColumns = [
  {
    title: '角色',
    dataIndex: 'name',
    key: 'name',
    scopedSlots: {
      customRender: 'name',
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
    RolesAddComponent,
    RolesEditComponent,
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
      form: this.$form.createForm(this, { name: 'roles' }),
    }
  },
  methods: {
    // 获取
    getRoles() {
      this.loading = true
      const params = {
        page_size: this.pagination.pageSize,
        page: this.pagination.current,
        is_page: true,
        ...this.filters,
      }
      getRolesApi(params)
        .then((res) => {
          this.pagination.total = res.total
          this.tableData = res.data
        })
        .catch((_error) => {})
        .finally(() => {
          this.loading = false
        })
    },
    editRow(row) {
      this.$refs.RolesEditComponent.showModal(row)
    },
    onAdd() {
      this.$refs.RolesAddComponent.showModal()
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
          deleteRolesApi(id)
            .then((res) => {
              if (res.code === '0001') {
                _this.$message.warning(res.message)
              } else {
                _this.$message.info(res.message)
              }
            })
            .catch((_error) => {})
            .finally(() => {
              _this.getRoles()
            })
        },
        onCancel() {},
      })
    },
    handleTableChange(pager) {
      this.pagination.current = pager.current
      this.pagination.pageSize = pager.pageSize
      this.getRoles()
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
        this.getRoles()
      })
    },
  },
  mounted() {
    this.getRoles()
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