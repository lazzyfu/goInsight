<template>
  <div>
    <a-card title="SQL执行任务">
      <div class="table-page-search-wrapper">
        <a-row :gutter="[8, 8]">
          <a-col :span="8">
            <a-button icon="play-circle" @click="executeMTask($route.params.task_id)">全部执行</a-button>
          </a-col>
          <a-col :span="16">
            <a-row :gutter="[8, 8]">
              <a-form layout="inline" :form="form" @keyup.enter.native="handleSearch">
                <a-col :span="8">
                  <a-form-item>
                    <a-select placeholder="状态" v-decorator="decorator['progress']">
                      <a-select-option v-for="s in progress" :key="s.key" :value="s.key">{{ s.value }}</a-select-option>
                    </a-select>
                  </a-form-item>
                </a-col>
                <a-col :span="8">
                  <a-form-item>
                    <a-input placeholder="输入要查询的工单内容" v-decorator="decorator['search']" />
                  </a-form-item>
                </a-col>
                <a-col :span="8">
                  <a-button type="primary" @click="handleSearch">查询</a-button>
                  <a-button @click="resetForm" style="margin-left: 6px">重置</a-button>
                </a-col>
              </a-form>
            </a-row>
          </a-col>
        </a-row>
      </div>

      <a-table
        :columns="table.columns"
        :rowKey="(record) => record.id"
        :dataSource="table.data"
        :pagination="pagination"
        :loading="loading"
        @change="handleTableChange"
        size="middle"
      >
        <!-- 设置序号 -->
        <span slot="num" slot-scope="text, record, index">
          <span v-text="index + 1"></span>
        </span>
        <!-- 格式化显示SQL -->
        <span slot="sql" slot-scope="text">
          <a-tooltip placement="topLeft">
            <template slot="title">{{ text }}</template>
            <span href="#">{{ text }}</span>
          </a-tooltip>
        </span>
        <span slot="execute" slot-scope="text, record">
          <a-icon type="right-circle" @click="executeSTask(record)" />
        </span>
        <span slot="progress" slot-scope="text">
          <span v-if="text === '处理中'" style="color: orange">{{ text }}</span>
          <span v-else-if="text === '已完成'" style="color: green">{{ text }}</span>
          <span v-else-if="text === '失败'" style="color: red">{{ text }}</span>
          <span v-else-if="text === '暂停'" style="color: blue">{{ text }}</span>
          <span v-else>{{ text }}</span>
        </span>
        <span slot="ghost_pause" slot-scope="text, record">
          <a-icon type="pause-circle" @click="ghostTaskThrottle(record, 'pause')" />
        </span>
        <span slot="ghost_recovery" slot-scope="text, record">
          <a-icon type="play-circle" @click="ghostTaskThrottle(record, 'recovery')" />
        </span>
        <span slot="result" slot-scope="text, record">
          <a-icon type="eye" @click="showModal(record)" />
        </span>
      </a-table>
      <div v-show="executeMsgVisible">
        <el-divider content-position="center">执行输出</el-divider>
        <codemirror ref="myCm3" v-model="code3" :options="executeMsgOptions" @ready="onCmReady3"></codemirror>

        <a-modal v-model="modalVisible" title="任务执行详情" width="70%">
          <template slot="footer">
            <a-button key="back" @click="handleCancel">Close</a-button>
          </template>
          <div v-loading="resultLoading" element-loading-text="玩命加载中..." element-loading-spinner="el-icon-loading">
            <div>
              <h5>-> 执行日志</h5>
              <codemirror ref="myCm2" v-model="code2" :options="cmOptions" @ready="onCmReady2"></codemirror>
            </div>
            <div style="margin-top: 15px">
              <h5>-> 回滚SQL</h5>
              <codemirror ref="myCm1" v-model="code1" :options="cmOptions" @ready="onCmReady1"></codemirror>
            </div>
          </div>
        </a-modal>
      </div>
    </a-card>
  </div>
</template>

<script>
import moment from 'moment'

import 'codemirror/mode/sql/sql.js'
import 'codemirror/addon/display/autorefresh'

import { getSqlOrdersTasksList, executeSingleTask, executeMultiTask, getTasksResult, TaskThrottle } from '@/api/sql'
import { taskProgress } from '@/utils/sql'

let protocol = 'ws://'
if (window.location.protocol === 'https:') {
  protocol = 'wss://'
}

export default {
  name: 'tasks-list',
  data() {
    return {
      modalLoading: false,
      modalVisible: false,
      executeMsgVisible: false,
      resultLoading: false,
      timer: '',
      websocket: {
        path: `${protocol}/${window.location.host}/ws/sql/${this.$route.params.task_id}/`,
        // path: `${protocol}/127.0.0.1:8000/ws/sql/${this.$route.params.task_id}/`,
        socket: '',
      },
      loading: false,
      progress: taskProgress,
      table: {
        columns: null,
        data: null,
      },
      pagination: {
        current: 1,
        pageSize: 10,
        total: 0,
        pageSizeOptions: ['10', '20'],
        showSizeChanger: true,
      },
      decorator: {
        progress: ['progress', { rules: [{ required: false }] }],
        search: ['search', { rules: [{ required: false }] }],
      },
      code1: '',
      code2: '',
      code3: '',
      cmOptions: {
        mode: 'text/x-mysql',
        indentUnit: 2,
        tabSize: 2,
        indentWithTabs: true,
        smartIndent: true,
        autoRefresh: true,
        lineWrapping: true,
        viewportMargin: Infinity,
        readOnly: true,
        autofocus: false,
      },
      executeMsgOptions: {
        smartIndent: true,
        autoRefresh: true,
        lineWrapping: true,
        viewportMargin: Infinity,
        readOnly: true,
        autofocus: false,
      },
      form: this.$form.createForm(this),
    }
  },
  methods: {
    // 初始化websocket
    init_websocket() {
      if (typeof WebSocket === 'undefined') {
        this.$message.error('您的浏览器不支持websocket')
      }
      // 实例化websocket
      this.websocket.socket = new WebSocket(this.websocket.path)
      // 监控websocket连接
      this.websocket.socket.onopen = this.socketOnOpen
      // 监听socket错误信息
      this.websocket.socket.onerror = this.socketOnError
      // 监听socket消息
      this.websocket.socket.onmessage = this.socketOnMessage
    },
    socketOnOpen() {},
    socketOnError() {
      this.init_websocket()
    },
    socketOnMessage(msg) {
      // 接收socket信息
      let result = JSON.parse(msg.data)
      if (result.flag === this.$route.params.task_id) {
        this.executeMsgVisible = true
        if (result.data.type === 'processlist') {
          this.codemirror3.setValue(this.renderProcesslist(result.data.data))
        }
        if (result.data.type === 'execute') {
          this.codemirror3.setValue(result.data.data)
        }
        if (['export', 'ghost'].includes(result.data.type)) {
          // 追加显示
          this.codemirror3.replaceRange(result.data.data, this.codemirror3.getCursor(this.codemirror3.lastLine()))
          // 自动滚动到行的末尾
          this.codemirror3.setCursor(this.codemirror3.lastLine())
        }
      }
    },
    socketClose() {
      // 关闭socket
      this.websocket.socket.close()
    },
    renderProcesslist(data) {
      // 渲染执行SQL时，show processlist的输出
      let html = '当前SQL SESSION ID的SHOW PROCESSLIST实时输出：'
      for (let key in data) {
        html += '\n' + key + ': ' + data[key]
      }
      return html
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
        task_id: this.$route.params.task_id,
        ...this.filters,
      }
      this.loading = true
      getSqlOrdersTasksList(params)
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
    handleSearch(e) {
      e.preventDefault()
      this.form.validateFields((error, values) => {
        if (error) {
          return
        }
        this.filters = {
          progress: values['progress'],
          search: values['search'],
        }
        this.pagination.current = 1
        this.fetchData()
      })
    },
    resetForm() {
      this.form.resetFields()
    },
    executeSTask(value) {
      // 执行单个任务
      const data = { id: value.id }
      executeSingleTask(data).then((response) => {
        if (response.code === '0000') {
          this.$message.info(response.message)
        } else {
          this.$message.error(response.message)
        }
      })
    },
    executeMTask(value) {
      // 执行全部任务
      const data = { task_id: value }
      executeMultiTask(data).then((response) => {
        if (response.code === '0000') {
          this.$message.info(response.message)
        } else {
          this.$message.error(response.message)
        }
      })
    },
    ghostTaskThrottle(value, op) {
      // 节流任务
      const data = { id: value.id, op: op }
      TaskThrottle(data).then((response) => {
        if (response.code === '0000') {
          this.$message.info(response.message)
        } else {
          this.$message.error(response.message)
        }
      })
    },
    getTResult(value) {
      // 获取结果
      getTasksResult(value.id)
        .then((response) => {
          this.codemirror1.setValue(response.data.rollback_sql)
          this.codemirror2.setValue(response.data.execute_log)
          setTimeout(() => {
            this.resultLoading = false
          }, 1000)
        })
        .catch((err) => {
          setTimeout(() => {
            this.resultLoading = false
          }, 1000)
        })
    },
    showModal(value) {
      this.modalVisible = true
      this.resultLoading = true
      this.getTResult(value)
    },
    handleCancel(e) {
      this.modalVisible = false
    },
    onCmReady1(cm) {
      cm.setSize('height', `250px`)
    },
    onCmReady2(cm) {
      cm.setSize('height', `250px`)
    },
    onCmReady3(cm) {
      cm.setSize('height', `450px`)
    },
  },
  destroyed() {
    // 关闭websocket
    this.socketClose()
    // 销毁timer
    clearInterval(this.timer)
  },
  computed: {
    codemirror1() {
      return this.$refs.myCm1.codemirror
    },
    codemirror2() {
      return this.$refs.myCm2.codemirror
    },
    codemirror3() {
      return this.$refs.myCm3.codemirror
    },
  },
  mounted() {
    this.fetchData()
    this.init_websocket()
    // 每10s刷新一次接口
    if (this.timer) {
      clearInterval(this.timer)
    }
    this.timer = setInterval(() => {
      setTimeout(this.fetchData(), 0)
    }, 10000)
  },
}
</script>

<style>
.CodeMirror {
  border: 2px solid #eee;
  font-size: 11px;
  min-height: 100px;
  max-height: 800px;
  /* 支持上下拉伸 */
  resize: vertical;
  overflow: y !important;
}
.CodeMirror pre.CodeMirror-placeholder {
  color: #999;
}
</style>
