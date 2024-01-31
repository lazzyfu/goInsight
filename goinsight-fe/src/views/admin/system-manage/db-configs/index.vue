<template>
  <a-card>
    <!-- 搜索框 -->
    <a-form layout="inline" :form="form" @keyup.enter.native="handleSearch">
      <a-form-item>
        <a-cascader
          v-decorator="['organization_key']"
          style="width: 200px"
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
        <a-select allowClear style="width: 200px" v-decorator="['environment']" placeholder="请选择环境" show-search>
          <a-select-option v-for="(item, index) in environments" :key="index" :label="item.name" :value="item.name">
            {{ item.name }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item>
        <a-select allowClear style="width: 200px" v-decorator="['db_type']" placeholder="请选择DB类型" show-search>
          <a-select-option v-for="(item, index) in dbTypes" :key="index" :label="item" :value="item">
            {{ item }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item>
        <a-input-search allowClear style="width: 250px" placeholder="输入要查询内容" v-decorator="['search']" />
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
      <span slot="use_type" slot-scope="text">
        <a-tag v-if="text === '工单'" color="blue">
          {{ text }}
        </a-tag>
        <a-tag v-else color="purple">{{ text }}</a-tag>
      </span>
      <span slot="action" slot-scope="text, record">
        <div class="editable-row-operations">
          <a @click="() => editRow(record)"><span style="color: #409eff">编辑</span></a>
          <a-divider type="vertical" />
          <a type="dashed" @click="DeleteConfirm(record.id)"><span style="color: #409eff">删除</span></a>
        </div>
      </span>
      <span slot="expandedRowRender" slot-scope="record">
        <pre>{{ JSON.stringify(record.inspect_params, null, 2) }}</pre>
      </span>
    </a-table>
    <DBConfigsAddComponent ref="DBConfigsAddComponent" @refreshTable="getDBConfig"></DBConfigsAddComponent>
    <DBConfigsEditComponent ref="DBConfigsEditComponent" @refreshTable="getDBConfig"></DBConfigsEditComponent>
  </a-card>
</template>

<script>
import { getOrganizationsApi } from '@/api/users'
import { adminGetDBConfigApi, adminDeleteDBConfigApi, adminGetEnvironmentsApi } from '@/api/common'

import DBConfigsAddComponent from './DBConfigsAdd'
import DBConfigsEditComponent from './DBConfigsEdit'

const tableColumns = [
  {
    title: '用途',
    dataIndex: 'use_type',
    key: 'use_type',
    scopedSlots: {
      customRender: 'use_type',
    },
  },
  {
    title: '环境',
    dataIndex: 'environment_name',
    key: 'environment_name',
    scopedSlots: {
      customRender: 'environment_name',
    },
  },
  {
    title: '类型',
    dataIndex: 'db_type',
    key: 'db_type',
    scopedSlots: {
      customRender: 'db_type',
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
    title: '主机名',
    dataIndex: 'hostname',
    key: 'hostname',
    scopedSlots: {
      customRender: 'hostname',
    },
  },
  {
    title: '端口',
    dataIndex: 'port',
    key: 'port',
    scopedSlots: {
      customRender: 'port',
    },
  },
  {
    title: '备注',
    dataIndex: 'remark',
    key: 'remark',
    scopedSlots: {
      customRender: 'remark',
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

const dbTypes = ['MySQL', 'TiDB', 'ClickHouse']

export default {
  components: {
    DBConfigsAddComponent,
    DBConfigsEditComponent,
  },
  data() {
    return {
      loading: false,
      tableColumns,
      tableData: [],
      environments: [],
      organizations: [],
      dbTypes,
      pagination: {
        current: 1,
        pageSize: 10,
        total: 0,
        pageSizeOptions: ['10', '20'],
        showSizeChanger: true,
      },
      form: this.$form.createForm(this, { name: 'dbconfig' }),
    }
  },
  methods: {
    async getEnvironments() {
      try {
        const res = await adminGetEnvironmentsApi({ is_page: false })
        this.environments = res.data
      } catch (error) {
        this.$message.error('Failed to fetch environments:', error)
      }
    },
    async getOrganizations() {
      try {
        const res = await getOrganizationsApi({ is_page: false })
        this.organizations = res.data
      } catch (error) {
        this.$message.error('Failed to fetch organizations:', error)
      }
    },
    getDBConfig() {
      this.loading = true
      const params = {
        page_size: this.pagination.pageSize,
        page: this.pagination.current,
        is_page: true,
        ...this.filters,
      }
      adminGetDBConfigApi(params)
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
      this.$refs.DBConfigsEditComponent.showModal(row)
    },
    onAdd() {
      this.$refs.DBConfigsAddComponent.showModal()
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
          adminDeleteDBConfigApi(id)
            .then((res) => {
              if (res.code === '0001') {
                _this.$message.warning(res.message)
              } else {
                _this.$message.info(res.message)
              }
            })
            .catch((_error) => {})
            .finally(() => {
              _this.getDBConfig()
            })
        },
        onCancel() {},
      })
    },
    handleTableChange(pager) {
      this.pagination.current = pager.current
      this.pagination.pageSize = pager.pageSize
      this.getDBConfig()
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
        this.getDBConfig()
      })
    },
  },
  mounted() {
    this.getDBConfig()
    this.getEnvironments()
    this.getOrganizations()
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