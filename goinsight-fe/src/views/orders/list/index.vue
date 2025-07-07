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
        <a-select
          allowClear
          style="width: 250px"
          v-decorator="['environment', { initialValue: '' }]"
          placeholder="请选择环境"
          show-search
        >
          <a-select-option value="">所有环境</a-select-option>
          <a-select-option v-for="(item, index) in environments" :key="index" :label="item.name" :value="item.id">
            {{ item.name }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item>
        <a-select
          allowClear
          style="width: 250px"
          v-decorator="['progress', { initialValue: '' }]"
          placeholder="请选择进度"
          show-search
        >
          <a-select-option value="">所有进度</a-select-option>
          <a-select-option v-for="(item, index) in progs" :key="index" :label="item" :value="item">{{
            item
          }}</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item>
        <a-input-search allowClear style="width: 350px" placeholder="输入要查询工单标题" v-decorator="['search']" />
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
      <!-- 进度插槽 -->
      <span slot="progress" slot-scope="text">
        <a-tag :color="getProgressTagColor(text)" class="progress-tag" circle>
          <a-icon :type="getProgressIcon(text)" />
          {{ text }}
        </a-tag>
      </span>

      <!-- 工单标题插槽 -->
      <span slot="order_title" slot-scope="text, record">
        <a-tooltip>
          <template slot="title">
            <span v-if="record.is_restrict_access">
              仅当前工单的提交人/审核人/复核人/抄送人可以查看工单内容（要求工单查看权限）
            </span>
            <span v-else>所有用户均可以查看工单内容（要求工单查看权限）</span>
          </template>
          <a-icon
            :type="record.is_restrict_access ? 'lock' : 'unlock'"
            :theme="record.is_restrict_access ? 'twoTone' : 'outlined'"
            :two-tone-color="record.is_restrict_access ? '#eb2f96' : '#52c41a'"
            class="access-icon"
          />
        </a-tooltip>
        <router-link :to="{ name: 'view.orders.detail', params: { order_id: record.order_id } }" class="title-link">
          <a-tooltip>
            <template slot="title">{{ text }}</template>
            <span class="title-text">{{ text }}</span>
          </a-tooltip>
        </router-link>
        <div class="title-meta">
          <a-icon type="clock-circle" />
          <span>{{ record.created_at }}</span>
        </div>
      </span>
      <!-- 环境插槽 -->
      <span slot="environment" slot-scope="text">
        <a-tag>
          <a-icon type="cloud" />
          {{ text }}
        </a-tag>
      </span>
      <!-- SQL类型插槽 -->
      <span slot="sql_type" slot-scope="text">
        <span v-if="text === 'EXPORT'">数据导出</span>
        <span v-else>{{ text }}</span>
      </span>
      <!-- 实例/库插槽 -->
      <span slot="instance" slot-scope="text, record" class="instance-cell">
        <div class="instance-info">
          <a-icon type="link" />
          <span class="instance-name">{{ text }}</span>
        </div>
        <div class="schema-info">
          <a-icon type="database" />
          <span class="schema-name">{{ record.schema }}</span>
        </div>
      </span>
      <span slot="approver" slot-scope="text">
        <div v-for="tag of text" :key="tag.user + tag.status">
          <a-icon
            :type="tag.status === 'pending' ? 'clock-circle' : 'check-circle'"
            :style="{ color: tag.status === 'pending' ? '#fa1a16' : '#52c41a' }"
          />
          <span class="user-text" :style="{ color: tag.status === 'pending' ? '#fa1a16' : '#52c41a' }">
            {{ tag.user }}
          </span>
        </div>
      </span>

      <span slot="reviewer" slot-scope="text">
        <div v-for="tag of text" :key="tag.user + tag.status">
          <a-icon
            :type="tag.status === 'pending' ? 'clock-circle' : 'check-circle'"
            :style="{ color: tag.status === 'pending' ? '#fa1a16' : '#52c41a' }"
          />
          <span class="user-text" :style="{ color: tag.status === 'pending' ? '#fa1a16' : '#52c41a' }">
            {{ tag.user }}
          </span>
        </div>
      </span>
      <!-- 自定义申请人标题 -->
      <span slot="customApplicant">
        <a-tooltip placement="topLeft" title="工单申请人" arrow-point-at-center> 申请人</a-tooltip>
      </span>
    </a-table>
  </a-card>
</template>

<script>
import { getEnvironmentsApi, getListApi } from '@/api/orders'

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
    title: '组织',
    dataIndex: 'organization',
    key: 'organization',
    ellipsis: true,
    scopedSlots: {
      customRender: 'organization',
    },
  },
  {
    title: '工单环境',
    dataIndex: 'environment',
    key: 'environment',
    ellipsis: true,
    scopedSlots: {
      customRender: 'environment',
    },
  },
  {
    title: '工单类型',
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
      progs: ['待审核', '已驳回', '已批准', '执行中', '已关闭', '已完成', '已复核'],
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

    // 获取进度标签颜色
    getProgressTagColor(progress) {
      const colorMap = {
        待审核: 'blue',
        已驳回: 'red',
        已关闭: 'red',
        已批准: 'orange',
        已复核: 'purple',
        执行中: 'pink',
        已完成: 'green',
      }
      return colorMap[progress] || 'default'
    },

    // 获取进度图标
    getProgressIcon(progress) {
      const iconMap = {
        待审核: 'clock-circle',
        已驳回: 'close-circle',
        已关闭: 'stop',
        已批准: 'check-circle',
        已复核: 'audit',
        执行中: 'loading',
        已完成: 'check-circle',
      }
      return iconMap[progress] || 'question-circle'
    },

    // 获取数据
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
