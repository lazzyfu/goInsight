<template>
  <a-card>
    <!-- 搜索框 -->
    <a-form layout="inline" :form="form" @keyup.enter.native="handleSearch">
      <a-form-item>
        <a-switch
          class="cus-switch"
          :checked="onlyMyOrders"
          checked-children="我的工单"
          un-checked-children="我的工单"
          @change="onMyChange"
        />
      </a-form-item>
      <a-form-item>
        <a-select allowClear style="width: 200px" v-decorator="['environment']" placeholder="请选择环境" show-search>
          <a-select-option v-for="(item, index) in environments" :key="index" :label="item.name" :value="item.id">
            {{ item.name }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item>
        <a-select allowClear style="width: 200px" v-decorator="['progress']" placeholder="请选择进度" show-search>
          <a-select-option v-for="s in progs" :key="s" :value="s">{{ s }}</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item>
        <a-input-search allowClear style="width: 300px" placeholder="输入要查询工单标题内容" v-decorator="['search']" />
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
    >
      <span slot="progress" slot-scope="text">
        <!-- ('待审核', '已驳回', '已批准', '执行中', '已关闭', '已完成', '已复核', '已勾住') -->
        <a-tag v-if="text === '待审核'" color="blue">
          {{ text }}
        </a-tag>
        <a-tag v-else-if="['已驳回', '已关闭'].indexOf(text) >= 0" color="red">
          {{ text }}
        </a-tag>
        <a-tag v-else-if="text === '已批准'" color="orange">
          {{ text }}
        </a-tag>
        <a-tag v-else-if="text === '已复核'" color="purple">
          {{ text }}
        </a-tag>
        <a-tag v-else-if="text === '执行中'" color="pink">
          {{ text }}
        </a-tag>
        <a-tag v-else-if="text === '已勾住'">
          {{ text }}
        </a-tag>
        <a-tag v-else color="green">
          {{ text }}
        </a-tag>
      </span>
      <span slot="applicant" slot-scope="text, record">
        <a-tooltip>
          <template slot="title">
            <span v-if="record.is_restrict_access"
              >仅当前工单的提交人/审核人/复核人/抄送人可以查看工单内容（要求工单查看权限）</span
            >
            <span v-else>所有用户均可以查看工单内容（要求工单查看权限）</span>
          </template>
          <a-icon
            type="lock"
            theme="twoTone"
            two-tone-color="#eb2f96"
            v-if="record.is_restrict_access"
          />
          <a-icon type="unlock" v-else />
        </a-tooltip>
        {{ text }}
      </span>
      <span slot="order_title" slot-scope="text, record">
        <router-link :to="{ name: 'view.orders.detail', params: { order_id: record.order_id } }">{{
          text
        }}</router-link>
        {{ text }}
        <br />
        At: {{ record.created_at }}
      </span>
      <span slot="instance" slot-scope="text, record">
        {{ text }}
        <br />
        {{ record.schema }}
      </span>
      <span slot="approver" slot-scope="text">
        <div v-for="tag of text" :key="tag.user + tag.status">
          <span :style="{ color: tag.status === 'pending' ? '#f56c6c' : '#67c23a' }">
            <span>{{ tag.user }}</span>
          </span>
        </div>
      </span>
      <span slot="reviewer" slot-scope="text">
        <div v-for="tag of text" :key="tag.user + tag.status">
          <span :style="{ color: tag.status === 'pending' ? '#f56c6c' : '#67c23a' }">
            <span>{{ tag.user }}</span>
          </span>
        </div>
      </span>
      <span slot="customApplicant">
        <a-tooltip placement="topLeft" title="工单申请人" arrow-point-at-center> 申请人</a-tooltip>
      </span>
    </a-table>
  </a-card>
</template>

<script>
import { getListApi, getEnvironmentsApi } from '@/api/orders'

import { mapActions, mapGetters } from 'vuex'

const tableColumns = [
  {
    title: '进度',
    dataIndex: 'progress',
    key: 'progress',
    width: 100,
    fixed: 'left',
    scopedSlots: {
      customRender: 'progress',
    },
  },
  {
    title: '标题',
    dataIndex: 'order_title',
    key: 'order_title',
    width: 300,
    fixed: 'left',
    ellipsis: true,
    scopedSlots: {
      customRender: 'order_title',
    },
  },
  {
    dataIndex: 'applicant',
    key: 'applicant',
    ellipsis: true,
    slots: { title: 'customApplicant' },
    scopedSlots: {
      customRender: 'applicant',
    },
  },
  {
    title: '环境',
    dataIndex: 'environment',
    key: 'environment',
    ellipsis: true,
    scopedSlots: {
      customRender: 'environment',
    },
  },
  {
    title: '类型',
    dataIndex: 'sql_type',
    key: 'sql_type',
    scopedSlots: {
      customRender: 'sql_type',
    },
  },
  {
    title: '实例/库',
    dataIndex: 'instance',
    key: 'instance',
    scopedSlots: {
      customRender: 'instance',
    },
  },
  {
    title: '审核人',
    dataIndex: 'approver',
    key: 'approver',
    scopedSlots: {
      customRender: 'approver',
    },
  },
  {
    title: '复核人',
    dataIndex: 'reviewer',
    key: 'reviewer',
    scopedSlots: {
      customRender: 'reviewer',
    },
  },
]

export default {
  computed: { ...mapGetters(['sqlState']) },
  data() {
    return {
      loading: false,
      onlyMyOrders: false,
      environments: [],
      progs: ['待审核', '已驳回', '已批准', '执行中', '已关闭', '已完成', '已复核', '已勾住'],
      tableColumns,
      tableData: [],
      pagination: {
        current: 1,
        pageSize: 10,
        total: 0,
        pageSizeOptions: ['10', '20'],
        showSizeChanger: true,
      },
      // name设置表单域内字段id的前缀，解决如果在一个页面中存在多个form.create，但同时他们控件绑定的字段名又有相同的情况
      form: this.$form.createForm(this, { name: 'list' }),
    }
  },
  methods: {
    ...mapActions(['storeMyOrder']),
    // 我的工单，点击显示我的工单
    onMyChange(checked) {
      this.storeMyOrder(checked ? 'ON' : 'OFF')
      this.onlyMyOrders = this.sqlState.orders.my_order == 'ON' ? true : false
      // 重置分页参数
      this.pagination.current = 1
      this.pagination.pageSize = 10
      this.pagination.total = 0
      this.getList()
    },
    // 获取
    getList() {
      this.loading = true
      const params = {
        page_size: this.pagination.pageSize,
        page: this.pagination.current,
        is_page: true,
        only_my_orders: this.onlyMyOrders === true ? 1 : 0,
        ...this.filters,
      }
      getListApi(params)
        .then((response) => {
          if (response != undefined) {
            this.pagination.total = response.total
            this.tableData = response.data
          }
        })
        .catch((_error) => {}) // .catch要加上，否则会抛出Uncaught (in promise) Error
        .finally(() => {
          this.loading = false
        })
    },
    // 获取环境
    getEnvironments() {
      getEnvironmentsApi({ is_page: false })
        .then((res) => {
          this.environments = res.data
        })
        .catch((_error) => {}) //加上_表示可以声明不使用
    },
    // 编辑
    editRow(row) {
      this.$refs.RolesEditComponent.show(row)
    },
    // 新建
    onAdd() {
      this.$refs.RolesAddComponent.show()
    },
    handleTableChange(pager) {
      this.pagination.current = pager.current
      this.pagination.pageSize = pager.pageSize
      this.getList()
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
          ...values,
        }
        this.pagination.current = 1
        this.getList()
      })
    },
  },
  destroyed() {
    // 销毁timer
    clearInterval(this.timer)
  },
  mounted() {
    this.getEnvironments()
    // 加载我的工单状态
    if (this.sqlState.orders.my_order == 'ON') {
      this.onlyMyOrders = true
    }
    this.getList()
    // 每30s刷新一次接口
    if (this.timer) {
      clearInterval(this.timer)
    }
    this.timer = setInterval(() => {
      setTimeout(this.getList(), 0)
    }, 30000)
  },
}
</script>

<style lang="less" scoped>
.ant-card-body {
  padding: 8px;
}

.ant-form {
  margin-bottom: 8px;
}

.ant-tag {
  line-height: 28px;
  font-size: 14px;
}

.cus-switch.ant-switch {
  height: 30px;
  .ant-switch-inner {
    font-size: 14px;
  }
  &.ant-switch::after {
    width: 22px;
    height: 22px;
    top: 3px;
  }
}
</style>
