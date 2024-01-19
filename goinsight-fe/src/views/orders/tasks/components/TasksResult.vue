<template>
  <a-modal v-model="visible" title="执行详情" width="50%">
    <template slot="footer">
      <a-button key="back" @click="handleCancel"> 关闭 </a-button>
    </template>
    <div v-if="executeResult != null">
      <a-card size="small">
        <a-row :gutter="16">
          <a-col :span="8">
            <a-statistic title="执行耗时" :value="executeResult.execute_cost_time" class="demo-class"> </a-statistic>
          </a-col>
          <a-col :span="8">
            <a-statistic title="备份耗时" :value="executeResult.backup_cost_time" class="demo-class"> </a-statistic>
          </a-col>
          <a-col :span="8">
            <a-statistic title="影响行数" :value="executeResult.affected_rows" style="margin-right: 50px">
            </a-statistic>
          </a-col>
        </a-row>
      </a-card>
      <a-card v-if="this.executeResult.error != null" size="small" title="错误信息" style="margin-top: 6px">
        <codemirror ref="myCmErr" v-model="codeErr" :options="cmOptions" @ready="onCmReadyErr"></codemirror>
      </a-card>
      <a-card size="small" title="执行日志" style="margin-top: 6px">
        <codemirror ref="myCmLog" v-model="codeLog" :options="cmOptions" @ready="onCmReadyLog"></codemirror>
      </a-card>
      <a-card size="small" title="回滚SQL" style="margin-top: 6px">
        <codemirror ref="myCmRbsql" v-model="codeRbsql" :options="cmOptions" @ready="onCmReadyRbsql"></codemirror>
      </a-card>
    </div>
  </a-modal>
</template>

<script>
import 'codemirror/mode/sql/sql.js'
import 'codemirror/addon/selection/active-line'
import 'codemirror/addon/display/autorefresh'
import 'codemirror/addon/comment/comment'
import 'codemirror/addon/edit/matchbrackets'
import 'codemirror/addon/edit/closebrackets'
import 'codemirror/addon/mode/overlay'

export default {
  data() {
    return {
      visible: false,
      executeResult: {},
      codeLog: '',
      codeErr: '',
      codeRbsql: '',
      cmOptions: {
        mode: 'text/x-mysql',
        indentUnit: 2,
        tabSize: 2,
        indentWithTabs: true,
        smartIndent: true,
        autoRefresh: true,
        lineNumbers: true,
        lineWrapping: true,
        readOnly: true,
        focuse: false,
      },
    }
  },
  methods: {
    onCmReadyRbsql(cm) {
      cm.setSize('height', `200px`)
    },
    onCmReadyLog(cm) {
      cm.setSize('height', `200px`)
    },
    onCmReadyErr(cm) {
      cm.setSize('height', `100px`)
    },
    // show modal
    showModal(data) {
      this.visible = true
      this.executeResult = data
      console.log(data)
      if (data != null) {
        this.$nextTick(() => {
          this.codemirrorRbsql.setValue(data.rollback_sql)
          this.codemirrorLog.setValue(data.execute_log)
          if (data.error != null) {
            this.codemirrorErr.setValue(data.error)
          }
        })
      }
    },
    // close modal
    handleCancel(e) {
      this.visible = false
    },
  },
  computed: {
    codemirrorRbsql() {
      return this.$refs.myCmRbsql.codemirror
    },
    codemirrorLog() {
      return this.$refs.myCmLog.codemirror
    },
    codemirrorErr() {
      return this.$refs.myCmErr.codemirror
    },
  },
}
</script>