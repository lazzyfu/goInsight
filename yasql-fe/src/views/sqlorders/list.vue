<template>
  <a-card title="SQL工单">
      <div class="table-page-search-wrapper">
        <a-form layout="inline" :form="form" @keyup.enter.native="handleSearch">
          <a-row :gutter="[8, 8]">
            <a-col :md="2" :sm="24">
              <a-form-item>
                <a-switch
                  style="margin-bottom: 1px"
                  checked-children="我的工单"
                  un-checked-children="我的工单"
                  @change="onMyChange"
                />
              </a-form-item>
            </a-col>

            <a-col :md="3" :sm="24">
              <a-form-item>
                <a-select placeholder="环境" v-decorator="decorator['env']">
                  <a-select-option v-for="s in envs" :key="s.id" :value="s.id">{{ s.name }}</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :md="3" :sm="24">
              <a-form-item>
                <a-select placeholder="状态" v-decorator="decorator['progress']">
                  <a-select-option v-for="s in progress" :key="s.key" :value="s.key">{{ s.value }}</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :md="4" :sm="24">
              <a-form-item>
                <a-input placeholder="输入要查询的工单内容" v-decorator="decorator['search']" />
              </a-form-item>
            </a-col>
            <a-col :md="5" :sm="24">
              <a-form-item>
                <a-range-picker v-decorator="decorator['created_at']" />
              </a-form-item>
            </a-col>
            <a-col :md="4" :sm="24">
              <span class="table-page-search-submitButtons">
                <a-button type="primary" @click="handleSearch">查询</a-button>
                <a-button @click="resetForm" style="margin-left: 8px">重置</a-button>
              </span>
            </a-col>
          </a-row>
        </a-form>
      </div>
      <a-table
        :columns="table.columns"
        :rowKey="(record) => record.id"
        :dataSource="table.data"
        :pagination="pagination"
        :loading="loading"
        @change="handleTableChange"
        size="middle"
        :scroll="{ x: 1100 }"
      >
        <span slot="progress" slot-scope="text">
          <div v-for="tag of progress" :key="tag.value">
            <el-button size="small" round plain :type="tag.color" v-if="tag.value === text">{{ text }}</el-button>
          </div>
        </span>
        <span slot="applicant" slot-scope="text, record">
          <el-tooltip placement="right-end" effect="light">
            <div slot="content">
              <span v-if="record.is_hide === 'ON'">
                仅有查看权限且仅工单的提交人、审核人、复核人和DBA可以查看工单内容
              </span>
              <span v-else> 有查看权限的用户可以查看当前工单内容 </span>
            </div>
            <i class="el-icon-lock table-msg" v-if="record.is_hide === 'ON'" style="color: #52c41a" />
            <i class="el-icon-lock table-msg" v-else />
          </el-tooltip>
          {{ text }}
        </span>
        <span slot="department" slot-scope="text">
          <div v-for="dept of text.split(',')" :key="dept">
            <span>{{ dept }}</span>
          </div>
        </span>
        <span slot="escape_title" slot-scope="text, record">
          <router-link :to="{ name: 'view.sqlorders.detail', params: { order_id: record.order_id } }">{{
            text
          }}</router-link>
          <br />
          At: {{ record.created_at }}
        </span>
        <span slot="host" slot-scope="text, record">
          {{ record.host }}:{{ record.port }}
          <br />
          {{ record.database }}
        </span>
        <template slot="version" slot-scope="text">
          <span v-if="text">
            <router-link :to="{ name: 'view.sqlorders.version.view', params: { version: text } }">{{
              text
            }}</router-link>
          </span>
          <span v-else>-</span>
        </template>
        <span slot="auditor" slot-scope="text">
          <div v-for="tag of JSON.parse(text)" :key="tag.user + tag.status">
            <span :style="{ color: tag.status === 0 ? '#f56c6c' : '#67c23a' }">
              <span v-if="tag.display_name">{{ tag.display_name }}</span>
              <span v-else>{{ tag.user }}</span>
            </span>
          </div>
        </span>

        <span slot="reviewer" slot-scope="text">
          <div v-for="tag of JSON.parse(text)" :key="`reviewer_` + tag.user + tag.status">
            <span :style="{ color: tag.status === 0 ? '#f56c6c' : '#67c23a' }">
              <span v-if="tag.display_name">{{ tag.display_name }}</span>
              <span v-else>{{ tag.user }}</span>
            </span>
          </div>
        </span>
      </a-table>
  </a-card>
</template>

<script>
import moment from 'moment'
import { getSqlOrdersList, getDbEnvironment } from '@/api/sql'
import { orderProgress } from '@/utils/sql'

export default {
  data() {
    return {
      loading: false,
      timer: '',
      username: '',
      envs: 0,
      progress: orderProgress,
      confirmMsg: '',
      visible: false,
      pagination: {
        current: 1,
        pageSize: 10,
        total: 0,
        pageSizeOptions: ['5', '10', '20'],
        showSizeChanger: true,
      },
      filter: {},
      filterDisabledBtn: ['已勾住', '未通过', '已关闭'],
      table: {
        columns: null,
        data: null,
      },
      decorator: {
        env: ['env', { rules: [{ required: false }] }],
        progress: ['progress', { rules: [{ required: false }] }],
        search: ['search', { rules: [{ required: false }] }],
        created_at: ['created_at', { rules: [{ required: false }] }],
      },
      form: this.$form.createForm(this),
    }
  },
  methods: {
    // 获取工单环境
    getEnvs() {
      getDbEnvironment.then((response) => {
        this.envs = response.data
      })
    },
    // 我的工单，点击显示我的工单
    onMyChange(checked) {
      this.resetAuditStatus = checked ? 'ON' : 'OFF'
      this.username = ''
      if (this.resetAuditStatus == 'ON') {
        this.username = this.$store.getters.userInfo.username
      }
      // 重置分页参数
      this.pagination.current = 1
      this.pagination.pageSize = 10
      this.pagination.total = 0
      this.fetchData()
    },
    handleTableChange(pager) {
      this.pagination.current = pager.current
      this.pagination.pageSize = pager.pageSize
      this.fetchData()
    },
    fetchData() {
      const params = {
        page_size: this.pagination.pageSize,
        page: this.pagination.current,
        username: this.username,
        ...this.filters,
      }
      this.loading = true
      getSqlOrdersList(params)
        .then((response) => {
          this.pagination.total = response.count
          this.loading = false
          this.table.columns = response.results.columns
          this.table.data = response.results.data
        })
        .finally(() => {
          this.loading = false
        })
    },
    // 搜索
    handleSearch(e) {
      e.preventDefault()
      this.form.validateFields((error, values) => {
        if (error) {
          return
        }
        this.filters = {
          progress: values['progress'],
          search: values['search'],
          env: values['env'],
          start_created_at: values['created_at'] ? moment(values['created_at'][0]).format('YYYY-MM-DD') : undefined,
          end_created_at: values['created_at'] ? moment(values['created_at'][1]).format('YYYY-MM-DD') : undefined,
        }
        this.pagination.current = 1
        this.fetchData()
      })
    },
    resetForm() {
      this.form.resetFields()
    },
  },
  destroyed() {
    // 销毁timer
    clearInterval(this.timer)
  },
  mounted() {
    this.getEnvs()
    this.fetchData()
    // 每5s刷新一次接口
    if (this.timer) {
      clearInterval(this.timer)
    }
    this.timer = setInterval(() => {
      setTimeout(this.fetchData(), 0)
    }, 30000)
  },
}
</script>
<style>
.ant-table {
  font-size: 13px;
}
.ant-pagination {
  font-size: 13px;
}
.ant-table-pagination.ant-pagination {
  font-size: 13px;
}
</style>
