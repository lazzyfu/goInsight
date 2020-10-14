.<template>
  <div>
    <a-row :gutter="16">
      <a-col :span="6">
        <a-statistic title="子任务总数量" :value="this.statistics.total" style="margin-right: 50px">
          <template #suffix>
            <a-icon type="flag" />
          </template>
        </a-statistic>
      </a-col>
      <a-col :span="6">
        <a-statistic title="已完成" :value="this.statistics.progress_1" class="demo-class">
          <template #suffix>
            <a-icon type="flag" />
          </template>
        </a-statistic>
      </a-col>
      <a-col :span="6">
        <a-statistic title="未执行" :value="this.statistics.progress_0" class="demo-class">
          <template #suffix>
            <a-icon type="flag" />
          </template>
        </a-statistic>
      </a-col>
      <a-col :span="6">
        <a-statistic title="失败" :value="this.statistics.progress_3" class="demo-class">
          <template #suffix>
            <a-icon type="flag" />
          </template>
        </a-statistic>
      </a-col>
    </a-row>
    <el-divider></el-divider>
    <div class="table-page-search-wrapper">
      <a-row :gutter="[8, 16]">
        <a-form laout="inline" :form="form" @keyup.enter.native="handleSearch">
          <a-col :span="4">
            <a-form-item>
              <a-select placeholder="状态" v-decorator="decorator['progress']">
                <a-select-option v-for="s in progress" :key="s.key" :value="s.key">{{ s.value }}</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="6">
            <a-form-item>
              <a-input placeholder="输入要查询的工单内容" v-decorator="decorator['search']" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item>
              <span class="table-page-search-submitButtons">
                <a-button type="primary" @click="handleSearch">查询</a-button>
                <a-button @click="resetForm" style="margin-left: 8px">重置</a-button>
              </span>
            </a-form-item>
          </a-col>
        </a-form>
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
      <span slot="progress" slot-scope="text">
        <span v-if="text === '处理中'" style="color: orange">{{ text }}</span>
        <span v-else-if="text === '已完成'" style="color: green">{{ text }}</span>
        <span v-else-if="text === '失败'" style="color: red">{{ text }}</span>
        <span v-else-if="text === '暂停'" style="color: blue">{{ text }}</span>
        <span v-else>{{ text }}</span>
      </span>
      <span slot="result" slot-scope="text, record">
        <a-icon type="eye" @click="showModal(record)" />
      </span>
    </a-table>

    <div>
      <a-modal v-model="modalVisible" title="任务执行详情" width="55%">
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
  </div>
</template>

<script>
import 'codemirror/mode/sql/sql.js'
import 'codemirror/addon/display/autorefresh'

import { getSqlOrdersTasksPreviewList, getTasksResult, getSqlOrdersTaskId } from '@/api/sql'
import { taskProgress } from '@/utils/sql'

export default {
  name: 'tasks-list',
  props: {
    order_id: Number,
  },
  data() {
    return {
      modalLoading: false,
      modalVisible: false,
      resultLoading: false,
      timer: '',
      loading: false,
      progress: taskProgress,
      statistics: {
        total: 0,
        progress_0: 0,
        progress_1: 0,
        progress_3: 0,
      },
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
      form: this.$form.createForm(this),
    }
  },
  methods: {
    handleTableChange(pager) {
      this.pagination.current = pager.current
      this.pagination.pageSize = pager.pageSize
      this.fetchData()
    },
    fetchData() {
      getSqlOrdersTaskId({ order_id: this.order_id }).then((response) => {
        this.task_id = response.data
        const params = {
          page_size: this.pagination.pageSize,
          page: this.pagination.current,
          task_id: this.task_id,
          ...this.filters,
        }
        this.loading = true
        getSqlOrdersTasksPreviewList(params)
          .then((response) => {
            this.pagination.total = response.count
            this.loading = false
            this.table.columns = response.results.columns
            this.table.data = response.results.data.data

            this.statistics.total = response.results.data.total
            this.statistics.progress_0 = response.results.data.progress_0
            this.statistics.progress_1 = response.results.data.progress_1
            this.statistics.progress_3 = response.results.data.progress_3
          })
          .finally(() => {
            this.loading = false
          })
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
    getTResult(value) {
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
  },

  computed: {
    codemirror1() {
      return this.$refs.myCm1.codemirror
    },
    codemirror2() {
      return this.$refs.myCm2.codemirror
    },
  },
  destroyed() {
    // 销毁timer
    clearInterval(this.timer)
  },
  mounted() {
    this.fetchData()
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
