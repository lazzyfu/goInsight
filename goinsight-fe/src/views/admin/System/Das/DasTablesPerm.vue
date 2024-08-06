<template>
  <div>
    <a-alert message="库表权限说明" type="info" show-icon>
      <p slot="description">
        <ul>
          <li>列表为空，表示用户有当前库所有表的访问权限</li>
          <li>允许用户访问当前库指定的一张或多张表时，请添加allow规则，没有指定allow规则的表不允许访问</li>
          <li>不允许用户访问当前库指定的一张或多张表时，请设置deny规则，没有指定deny规则的表允许访问</li>
          <li>当列表中有allow和deny规则时，仅allow规则生效，deny规则不生效，表示当前用户仅允许访问有allow规则的表</li>
        </ul>
      </p>
    </a-alert>
    <a-card :title="cardTitle" style="margin-top: 8px">
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
        <span slot="rule" slot-scope="text">
          <a-tag v-if="text === 'allow'" color="blue">
            {{ text }}
          </a-tag>
          <a-tag v-else color="red">
            {{ text }}
          </a-tag>
        </span>
        <span slot="action" slot-scope="text, record">
          <div class="editable-row-operations">
            <a type="dashed" @click="DeleteConfirm(record.id)"><span style="color: #409eff">删除</span></a>
          </div>
        </span>
      </a-table>
      <DasTablesPermAddComponent
        :inputData="this.$route.params"
        ref="DasTablesPermAddComponent"
        @refreshTable="getAdminTablesListGrant"
      ></DasTablesPermAddComponent>
    </a-card>
  </div>
</template>

<script>
import { adminGetTablesGrantApi, adminDeleteTablesGrantApi } from '@/api/das'

import DasTablesPermAddComponent from './DasTablesPermAdd'

const tableColumns = [
  {
    title: '规则',
    dataIndex: 'rule',
    key: 'rule',
    scopedSlots: {
      customRender: 'rule',
    },
  },
  {
    title: '用户',
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
    title: '表名',
    dataIndex: 'table',
    key: 'table',
    scopedSlots: {
      customRender: 'table',
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
    DasTablesPermAddComponent,
  },
  data() {
    return {
      loading: false,
      cardTitle: `用户${this.$route.params.username}当前${this.$route.params.schema}库的表权限列表`,
      tableColumns,
      tableData: [],
      pagination: {
        current: 1,
        pageSize: 10,
        total: 0,
        pageSizeOptions: ['10', '20'],
        showSizeChanger: true,
      },
      form: this.$form.createForm(this, { name: 'dasTablesPerm' }),
    }
  },
  methods: {
    getAdminTablesListGrant() {
      this.loading = true
      const params = {
        page_size: this.pagination.pageSize,
        page: this.pagination.current,
        is_page: true,
        instance_id: this.$route.params.instance_id,
        username: this.$route.params.username,
        schema: this.$route.params.schema,
        ...this.filters,
      }
      adminGetTablesGrantApi(params)
        .then((res) => {
          this.pagination.total = res.total
          this.tableData = res.data
        })
        .finally(() => {
          this.loading = false
        })
    },
    onAdd() {
      this.$refs.DasTablesPermAddComponent.showModal(this.$route.params)
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
          adminDeleteTablesGrantApi(id)
            .then((res) => {
              if (res.code === '0001') {
                _this.$message.warning(res.message)
              } else {
                _this.$message.info(res.message)
              }
            })
            .finally(() => {
              _this.getAdminTablesListGrant()
            })
        },
        onCancel() {},
      })
    },
    handleTableChange(pager) {
      this.pagination.current = pager.current
      this.pagination.pageSize = pager.pageSize
      this.getAdminTablesListGrant()
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
        this.getAdminTablesListGrant()
      })
    },
  },
  mounted() {
    this.getAdminTablesListGrant()
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

li {
  list-style-type: circle;
}
li li {
  list-style-type: square;
}
</style>
