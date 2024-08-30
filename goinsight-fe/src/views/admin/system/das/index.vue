<template>
  <a-card>
    <!-- 搜索框 -->
    <a-form layout="inline" :form="form" @keyup.enter.native="handleSearch">
      <a-form-item>
        <a-select allowClear style="width: 200px" v-decorator="['environment']" placeholder="请选择环境" show-search>
          <a-select-option v-for="(item, index) in environments" :key="index" :label="item.name" :value="item.name">
            {{ item.name }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item>
        <a-input-search allowClear style="width: 300px" placeholder="输入要查询内容" v-decorator="['search']" />
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
      <span slot="schema" slot-scope="text, record">
        <router-link
          :to="{
            name: 'view.admin.systemManage.das.detail',
            params: { username: record.username, schema: record.schema, instance_id: record.instance_id },
          }"
          >{{ text }}</router-link
        >
      </span>
      <span slot="hostname" slot-scope="text, record"> {{ text }}:{{ record.port }} </span>
      <span slot="action" slot-scope="text, record">
        <div class="editable-row-operations">
          <a type="dashed" @click="DeleteConfirm(record.id)"><span style="color: #409eff">删除</span></a>
        </div>
      </span>
    </a-table>
    <DasAddComponent ref="DasAddComponent" @refreshTable="getAdminGrants"></DasAddComponent>
  </a-card>
</template>

<script>
import { adminGetEnvironmentsApi } from '@/api/common'
import { adminGetSchemasListGrantApi, adminDeleteSchemasGrantApi } from '@/api/das'

import DasAddComponent from './DasAdd'

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
    title: '库名',
    dataIndex: 'schema',
    key: 'schema',
    scopedSlots: {
      customRender: 'schema',
    },
  },
  {
    title: '环境',
    dataIndex: 'environment',
    key: 'environment',
    scopedSlots: {
      customRender: 'environment',
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
    title: '实例ID',
    dataIndex: 'instance_id',
    key: 'instance_id',
    scopedSlots: {
      customRender: 'instance_id',
    },
  },
  {
    title: '主机',
    dataIndex: 'hostname',
    key: 'hostname',
    scopedSlots: {
      customRender: 'hostname',
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
  components: { DasAddComponent },
  data() {
    return {
      loading: false,
      environments: [],
      tableColumns,
      tableData: [],
      pagination: {
        current: 1,
        pageSize: 10,
        total: 0,
        pageSizeOptions: ['10', '20'],
        showSizeChanger: true,
      },
      form: this.$form.createForm(this, { name: 'das' }),
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
    getAdminGrants() {
      this.loading = true
      const params = {
        page_size: this.pagination.pageSize,
        page: this.pagination.current,
        is_page: true,
        ...this.filters,
      }
      adminGetSchemasListGrantApi(params)
        .then((res) => {
          if (res.code === "0000") {
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
    // 新建
    onAdd() {
      this.$refs.DasAddComponent.showModal()
    },
    // 删除确认
    DeleteConfirm(id) {
      const _this = this
      this.$confirm({
        title: '警告',
        content: '你确定删除？这将会删除当前的库表权限！！！',
        okText: 'Yes',
        okType: 'danger',
        cancelText: 'No',
        onOk() {
          adminDeleteSchemasGrantApi(id)
            .then((res) => {
              if (res.code === '0001') {
                _this.$message.warning(res.message)
              } else {
                _this.$message.info(res.message)
              }
            })
            .catch((_error) => {})
            .finally(() => {
              _this.getAdminGrants()
            })
        },
        onCancel() {},
      })
    },
    handleTableChange(pager) {
      this.pagination.current = pager.current
      this.pagination.pageSize = pager.pageSize
      this.getAdminGrants()
    },
    // 搜索
    handleSearch(e) {
      e.preventDefault()
      this.form.validateFields((error, values) => {
        if (error) {
          return
        }
        this.filters = {
          environment: values['environment'],
          search: values['search'],
        }
        this.pagination.current = 1
        this.getAdminGrants()
      })
    },
  },
  mounted() {
    this.getAdminGrants()
    this.getEnvironments()
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
