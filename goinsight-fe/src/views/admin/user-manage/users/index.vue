<template>
  <a-card>
    <!-- 搜索框 -->
    <a-form layout="inline" :form="form" @keyup.enter.native="handleSearch">
      <a-form-item>
        <a-cascader
          v-decorator="['organization_key']"
          style="width: 350px"
          :field-names="{ label: 'title', value: 'key', children: 'children' }"
          :options="organizations"
          change-on-select
          expand-trigger="hover"
          placeholder="请选择组织"
        >
          <a-icon slot="suffixIcon" type="apartment" />
        </a-cascader>
      </a-form-item>
      <a-form-item>
        <a-select allowClear style="width: 200px" v-decorator="['role_id']" placeholder="请选择角色" show-search>
          <a-select-option v-for="(item, index) in roles" :key="index" :label="item.name" :value="item.id">
            {{ item.name }}
          </a-select-option>
          <a-icon slot="suffixIcon" type="user" />
        </a-select>
      </a-form-item>
      <a-form-item>
        <a-input-search allowClear style="width: 300px" placeholder="请输入要查询内容" v-decorator="['search']" />
      </a-form-item>
      <a-form-item>
        <span class="table-page-search-submitButtons">
          <a-button type="primary" @click="handleSearch">查询</a-button>
        </span>
      </a-form-item>
      <a-form-item style="float: right; margin-right: 0px">
        <a-button type="primary" @click="onAdd"> <a-icon type="plus" />新增</a-button>
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
      <span slot="is_active" slot-scope="text">
        <span v-if="text">
          <a-tag color="green">是</a-tag>
        </span>
        <span v-else>
          <a-tag color="volcano">否</a-tag>
        </span>
      </span>
      <span slot="is_superuser" slot-scope="text">
        <span v-if="text">
          <a-tag color="green">是</a-tag>
        </span>
        <span v-else>
          <a-tag color="volcano">否</a-tag>
        </span>
      </span>
      <span slot="is_two_fa" slot-scope="text">
        <span v-if="text">
          <a-tag color="green">是</a-tag>
        </span>
        <span v-else>
          <a-tag color="volcano">否</a-tag>
        </span>
      </span>
      <span slot="action" slot-scope="text, record">
        <div class="editable-row-operations">
          <a @click="() => editRow(record)"><span style="color: #409eff">编辑</span></a>
          <a-divider type="vertical" />
          <a type="dashed" @click="changePassword(record)"><span style="color: #409eff">修改密码</span></a>
          <a-divider type="vertical" />
          <a type="dashed" @click="DeleteConfirm(record.uid)"><span style="color: #409eff">删除</span></a>
        </div>
      </span>
    </a-table>
    <UsersAddComponent ref="UsersAddComponent" @refreshTable="getUsers"></UsersAddComponent>
    <UsersEditComponent ref="UsersEditComponent" @refreshTable="getUsers"></UsersEditComponent>
    <UsersChangePassComponent ref="UsersChangePassComponent" @refreshTable="getUsers"></UsersChangePassComponent>
  </a-card>
</template>

<script>
import { getUsersApi, deleteUsersApi, getOrganizationsApi, getRolesApi } from '@/api/users'

import UsersAddComponent from './UsersAdd'
import UsersEditComponent from './UsersEdit'
import UsersChangePassComponent from './UsersChangePass'

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
]

export default {
  components: {
    UsersAddComponent,
    UsersEditComponent,
    UsersChangePassComponent,
  },
  data() {
    return {
      loading: false,
      tableColumns,
      tableData: [],
      organizations: [],
      roles: [],
      pagination: {
        current: 1,
        pageSize: 10,
        total: 0,
        pageSizeOptions: ['10', '20'],
        showSizeChanger: true,
      },
      form: this.$form.createForm(this, { name: 'users' }),
    }
  },
  methods: {
    getOrganizations() {
      this.loading = true
      const params = {
        is_page: false,
      }
      getOrganizationsApi(params)
        .then((res) => {
          this.organizations = res.data
        })
        .catch((_error) => {})
        .finally(() => {
          this.loading = false
        })
    },
    getRoles() {
      this.loading = true
      const params = {
        is_page: false,
      }
      getRolesApi(params)
        .then((res) => {
          this.roles = res.data
        })
        .catch((_error) => {})
        .finally(() => {
          this.loading = false
        })
    },
    getUsers() {
      this.loading = true
      const params = {
        page_size: this.pagination.pageSize,
        page: this.pagination.current,
        is_page: true,
        ...this.filters,
      }
      getUsersApi(params)
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
      this.$refs.UsersEditComponent.showModal(row)
    },
    onAdd() {
      this.$refs.UsersAddComponent.showModal()
    },
    changePassword(row) {
      this.$refs.UsersChangePassComponent.showModal(row)
    },
    DeleteConfirm(uid) {
      const _this = this
      this.$confirm({
        title: '警告',
        content: '你确定删除？',
        okText: 'Yes',
        okType: 'danger',
        cancelText: 'No',
        onOk() {
          deleteUsersApi(uid)
            .then((res) => {
              if (res.code === '0001') {
                _this.$message.warning(res.message)
              } else {
                _this.$message.info(res.message)
              }
            })
            .catch((_error) => {})
            .finally(() => {
              _this.getUsers()
            })
        },
        onCancel() {},
      })
    },
    handleTableChange(pager) {
      this.pagination.current = pager.current
      this.pagination.pageSize = pager.pageSize
      this.getUsers()
    },
    handleSearch(e) {
      e.preventDefault()
      this.form.validateFields((error, values) => {
        if (error) {
          return
        }
        let organization_key = ''
        if (values['organization_key'] != undefined) {
          organization_key = values['organization_key'][values['organization_key'].length - 1]
        }
        this.filters = {
          organization_key: organization_key,
          role_id: values['role_id'],
          search: values['search'],
        }
        this.pagination.current = 1
        this.getUsers()
      })
    },
  },
  mounted() {
    this.getUsers()
    this.getOrganizations()
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
