<template>
  <div>
    <a-card title="SQL执行任务">
      <!-- 搜索框 -->
      <a-form layout="inline" :form="form" @keyup.enter.native="handleSearch">
        <a-form-item>
          <a-button icon="play-circle" @click="executeMTask()">全部执行</a-button>
        </a-form-item>
        <a-form-item>
          <a-select allowClear style="width: 150px" v-decorator="['progress']" placeholder="请选择进度" show-search>
            <a-select-option v-for="s in progs" :key="s" :value="s">{{ s }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-input-search allowClear style="width: 250px" placeholder="输入要查询SQL内容" v-decorator="['search']" />
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
        <!-- 格式化显示SQL -->
        <span slot="sql" slot-scope="text">
          <a-tooltip placement="topLeft">
            <template slot="title">{{ text }}</template>
            <span href="#">{{ text }}</span>
          </a-tooltip>
        </span>
        <span slot="progress" slot-scope="text">
          <span v-if="text === '执行中'" style="color: orange">{{ text }}</span>
          <span v-else-if="text === '已完成'" style="color: green">{{ text }}</span>
          <span v-else-if="text === '已失败'" style="color: red">{{ text }}</span>
          <span v-else-if="text === '已暂停'" style="color: blue">{{ text }}</span>
          <span v-else>{{ text }}</span>
        </span>
        <span slot="action" slot-scope="text, record">
          <div class="editable-row-operations">
            <a @click="() => executeSTask(record)"><span style="color: #409eff">执行</span></a>
            <a-divider type="vertical" />
            <a type="dashed" @click="viewResult(record)"><span style="color: #409eff">结果</span></a>
          </div>
        </span>
      </a-table>
    </a-card>
    <TasksResultComponent ref="TasksResultComponent"></TasksResultComponent>
    <TasksWebsocketComponent ref="TasksWebsocketComponent"></TasksWebsocketComponent>
  </div>
</template>

<script>
import { getTasksApi, executeSingleTaskApi, executeAllTaskApi } from '@/api/orders'

import TasksWebsocketComponent from './components/TasksWebsocket.vue'
import TasksResultComponent from './components/TasksResult.vue'

const tableColumns = [
  {
    title: '进度',
    dataIndex: 'progress',
    key: 'progress',
    width: '10%',
    scopedSlots: {
      customRender: 'progress',
    },
  },
  {
    title: 'TaskID',
    dataIndex: 'task_id',
    key: 'task_id',
    width: '20%',
    scopedSlots: {
      customRender: 'task_id',
    },
  },
  {
    title: 'SQL文本',
    dataIndex: 'sql',
    key: 'sql',
    ellipsis: true,
    width: '35%',
    scopedSlots: {
      customRender: 'sql',
    },
  },
  {
    title: '更新时间',
    dataIndex: 'updated_at',
    key: 'updated_at',
    ellipsis: true,
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
    TasksWebsocketComponent,
    TasksResultComponent,
  },
  data() {
    return {
      loading: false,
      progs: ['未执行', '执行中', '已完成', '已失败', '已暂停'],
      tableColumns,
      tableData: [],
      pagination: {
        current: 1,
        pageSize: 10,
        total: 0,
        pageSizeOptions: ['10', '20', '50', '100'],
        showSizeChanger: true,
      },
      // name设置表单域内字段id的前缀，解决如果在一个页面中存在多个form.create，但同时他们控件绑定的字段名又有相同的情况
      form: this.$form.createForm(this, { name: 'tasks' }),
    }
  },
  methods: {
    handleTableChange(pager) {
      this.pagination.current = pager.current
      this.pagination.pageSize = pager.pageSize
      this.getTasks()
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
        this.getTasks()
      })
    },
    getTasks() {
      this.loading = true
      var params = {
        page_size: this.pagination.pageSize,
        page: this.pagination.current,
        is_page: true,
        order_id: this.$route.params.order_id,
        ...this.filters,
      }
      getTasksApi(params)
        .then((res) => {
          if (res.code === '0000') {
            this.pagination.total = res.total
            this.tableData = res.data
          }
        })
        .catch((_error) => {})
        .finally(() => {
          this.loading = false
        })
    },
    // 执行单个任务
    executeSTask(value) {
      this.$notification.info({
        message: '提示',
        description: '开始执行任务，请查看输出',
      })
      var data = {
        id: value.id,
        order_id: this.$route.params.order_id,
      }
      executeSingleTaskApi(data)
        .then((res) => {
          if (res.code === '0001') {
            this.$notification.error({
              message: '错误',
              description: res.message,
            })
          }
          this.getTasks()
        })
        .catch((_error) => {})
    },
    // 批量执行
    executeMTask() {
      this.$notification.info({
        message: '提示',
        description: '开始执行任务，请查看输出',
      })
      var data = {
        order_id: this.$route.params.order_id,
      }
      executeAllTaskApi(data)
        .then((res) => {
          if (res.code === '0001') {
            this.$notification.error({
              message: '错误',
              description: res.message,
            })
          }
          this.getTasks()
        })
        .catch((_error) => {})
    },
    // view result
    viewResult(data) {
      if (['已完成', '已失败'].indexOf(data.progress) >= 0){
        this.$refs.TasksResultComponent.showModal(data)
      } else {
        this.$notification.warning({
        message: '警告',
        description: '当前任务状态不为【已完成】或【已失败】',
      })
      }
    },
  },
  mounted() {
    this.getTasks()
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
