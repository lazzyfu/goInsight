<template>
  <a-card class="box-card">
    <!-- 搜索框 -->
    <a-form layout="inline" :form="form" @keyup.enter.native="handleSearch">
      <a-form-item>
        <a-input-search allowClear style="width: 300px" placeholder="请输入要查询的用户名" v-decorator="['search']" />
      </a-form-item>
      <a-form-item>
        <span class="table-page-search-submitButtons">
          <a-button type="primary" @click="handleSearch">查询</a-button>
        </span>
      </a-form-item>
      <a-form-item style="float: right; margin-right: 0px">
        <a-button type="primary" @click="onAdd"> <a-icon type="plus" />绑定用户</a-button>
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
          <a type="dashed" @click="DeleteConfirm(record)"><span style="color: #409eff">删除</span></a>
        </div>
      </span>
    </a-table>
    <NodeUsersAdd ref="NodeUsersAdd" @refresh="getOrganizationsUsers"></NodeUsersAdd>
  </a-card>
</template>
  
  <script>
import { getOrganizationsUsersApi, deleteOrganizationsUsersApi } from '@/api/users'

import NodeUsersAdd from './NodeUsersAdd'

const tableColumns = [
  {
    title: '用户',
    dataIndex: 'username',
    key: 'username',
    scopedSlots: {
      customRender: 'username',
    },
  },
  {
    title: '组织',
    dataIndex: 'organization_name',
    key: 'organization_name',
    scopedSlots: {
      customRender: 'organization_name',
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
    NodeUsersAdd,
  },
  props: {
    nodeKey: String,
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
      form: this.$form.createForm(this, { name: 'nodeUsers' }),
    }
  },
  methods: {
    getOrganizationsUsers() {
      this.loading = true
      const params = {
        page_size: this.pagination.pageSize,
        page: this.pagination.current,
        is_page: true,
        key: this.nodeKey,
        ...this.filters,
      }
      getOrganizationsUsersApi(params)
        .then((res) => {
          this.pagination.total = res.total
          this.tableData = res.data
        })
        .finally(() => {
          this.loading = false
        })
    },
    onAdd() {
      this.$refs.NodeUsersAdd.showModal(this.nodeKey)
    },
    DeleteConfirm(row) {
      const _this = this
      this.$confirm({
        title: '警告',
        content: '你确定删除？',
        okText: 'Yes',
        okType: 'danger',
        cancelText: 'No',
        onOk() {
          const data = {
            key: row.organization_key,
            uid: row.uid,
          }
          deleteOrganizationsUsersApi(data)
            .then((res) => {
              if (res.code === '0001') {
                _this.$message.warning(res.message)
              } else {
                _this.$message.info(res.message)
              }
            })
            .catch((_error) => {})
            .finally(() => {
              _this.getOrganizationsUsers()
            })
        },
        onCancel() {},
      })
    },
    handleTableChange(pager) {
      this.pagination.current = pager.current
      this.pagination.pageSize = pager.pageSize
      this.getOrganizationsUsers()
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
        this.getOrganizationsUsers()
      })
    },
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