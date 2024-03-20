<template>
  <a-modal v-model="visible" title="执行详情" width="50%">
    <template slot="footer">
      <a-button key="back" @click="handleCancel"> 关闭 </a-button>
    </template>

    <a-card size="small" title="执行信息">
      <a-row :gutter="16">
        <a-col :span="8">
          <a-statistic title="执行耗时" :value="executeResult.execute_cost_time" class="demo-class"> </a-statistic>
        </a-col>
        <a-col :span="8">
          <a-statistic title="备份耗时" :value="executeResult.backup_cost_time" class="demo-class"> </a-statistic>
        </a-col>
        <a-col :span="8">
          <a-statistic title="影响行数" :value="executeResult.affected_rows" style="margin-right: 50px"> </a-statistic>
        </a-col>
      </a-row>
    </a-card>

    <a-card v-show="sql_type === 'EXPORT'" size="small" title="导出文件信息" style="margin-top: 8px">
      <a-row :gutter="16">
        <a-col :span="24">
          <a-statistic title="文件名" :value="executeResult.file_name"> </a-statistic>
        </a-col>
        <a-col :span="24">
          <a-statistic title="文件大小（字节）" :value="executeResult.file_size"> </a-statistic>
        </a-col>
        <a-col :span="24">
          <a-statistic title="导出行数" :value="executeResult.export_rows"> </a-statistic>
        </a-col>
        <a-col :span="24">
          <a-statistic title="文件加密秘钥" :value="executeResult.encryption_key"> </a-statistic>
        </a-col>
        <a-col :span="24">
          <a-statistic title="文件下载路径" :value="executeResult.download_url"> </a-statistic>
        </a-col>
      </a-row>
    </a-card>

    <a-card v-show="executeResult.error != ''" size="small" title="错误信息" style="margin-top: 6px">
      <codemirror ref="myCmErr" v-model="codeErr" :options="cmOptions" @ready="onCmReadyErr"></codemirror>
    </a-card>
    <a-card v-show="executeResult.execute_log != ''" size="small" title="执行日志" style="margin-top: 6px">
      <codemirror ref="myCmLog" v-model="codeLog" :options="cmOptions" @ready="onCmReadyLog"></codemirror>
    </a-card>
    <a-card v-show="executeResult.rollback_sql != ''" size="small" title="回滚SQL" style="margin-top: 6px">
      <codemirror ref="myCmRbsql" v-model="codeRbsql" :options="cmOptions" @ready="onCmReadyRbsql"></codemirror>
    </a-card>
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
      sql_type: '',
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
      cm.setSize('height', `250px`)
    },
    onCmReadyLog(cm) {
      cm.setSize('height', `300px`)
    },
    onCmReadyErr(cm) {
      cm.setSize('height', `100px`)
    },
    // show modal
    showModal(data) {
      this.visible = true
      this.executeResult = data.result
      this.sql_type = data.sql_type

      if (this.executeResult != null) {
        this.$nextTick(() => {
          this.codemirrorRbsql.setValue(this.executeResult.rollback_sql)
          this.codemirrorLog.setValue(this.executeResult.execute_log)
          if (this.executeResult.error != null) {
            this.codemirrorErr.setValue(this.executeResult.error)
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


<style lang="less" scoped>
/deep/.ant-statistic-content {
  font-size: 14px;
}
</style>