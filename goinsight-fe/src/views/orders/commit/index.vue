<template>
  <div>
    <a-card :title="cardTitle">
      <a-row :gutter="{ xs: 8, sm: 16, md: 24, lg: 32 }">
        <!-- 左侧表单区域 -->
        <a-col class="gutter-row" :span="8">
          <a-card>
            <a-form :form="form" ref="formHeight" :label-col="{ span: 6 }" :wrapper-col="{ span: 18 }">
              <a-form-item label="标题" has-feedback>
                <a-input
                  placeholder="请输入工单标题"
                  v-decorator="[
                    'title',
                    {
                      rules: [{ required: true, min: 5, max: 96 }],
                      validateTrigger: 'blur',
                    },
                  ]"
                >
                </a-input>
              </a-form-item>
              <a-form-item label="备注" has-feedback>
                <a-textarea
                  v-decorator="['remark']"
                  :auto-size="{ minRows: 2, maxRows: 6 }"
                  placeholder="请输入工单需求或备注"
                >
                </a-textarea>
              </a-form-item>
              <a-form-item label="限制访问" help="开启后仅工单的提交人/审核人/复核人/抄送人可以查看工单内容">
                <a-select
                  v-decorator="[
                    'is_restrict_access',
                    { initialValue: 'YES', rules: [{ required: true, message: '是否限制访问' }] },
                  ]"
                  placeholder="是否限制访问"
                >
                  <a-select-option value="YES"> 开启 </a-select-option>
                  <a-select-option value="NO"> 关闭 </a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item label="DB类型" has-feedback>
                <a-select
                  @change="changeDBType"
                  v-decorator="[
                    'db_type',
                    { initialValue: 'MySQL', rules: [{ required: true, message: '请选择DB类型' }] },
                  ]"
                  placeholder="请选择DB类型"
                  allowClear
                  show-search
                >
                  <a-select-option v-for="(item, index) in dbTypes" :key="index" :label="item" :value="item">
                    {{ item }}
                  </a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item label="环境" has-feedback>
                <a-select
                  @change="changeEnvs"
                  v-decorator="['environment', { rules: [{ required: true, message: '请选择工单环境' }] }]"
                  placeholder="请选择工单环境"
                  allowClear
                  show-search
                >
                  <a-select-option
                    v-for="(item, index) in environments"
                    :key="index"
                    :label="item.name"
                    :value="item.id"
                  >
                    {{ item.name }}
                  </a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item label="实例" has-feedback>
                <a-select
                  @change="changeIns"
                  v-decorator="['instance_id', { rules: [{ required: true, message: '请选择数据库实例' }] }]"
                  placeholder="请选择数据库实例"
                  allowClear
                  show-search
                >
                  <a-select-option
                    v-for="(item, index) in instances"
                    :key="index"
                    :label="item.remark"
                    :value="item.instance_id"
                  >
                    {{ item.remark }}
                  </a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item label="库名" has-feedback>
                <a-select
                  v-decorator="['schema', { rules: [{ required: true, message: '请选择数据库' }] }]"
                  placeholder="请选择数据库"
                  allowClear
                  show-search
                >
                  <a-select-option
                    v-for="(item, index) in schemas"
                    :key="index"
                    :label="item.schema"
                    :value="item.schema"
                  >
                    {{ item.schema }}
                  </a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item label="审核人" help="工单审核人，一般为Leader" has-feedback>
                <a-select
                  v-decorator="['approver', { rules: [{ required: true, message: '请选择工单审核人' }] }]"
                  placeholder="请选择工单审核人"
                  mode="multiple"
                  allowClear
                  show-search
                >
                  <a-select-option
                    v-for="(item, index) in users"
                    :key="index"
                    :label="item.username"
                    :value="item.username"
                  >
                    {{ item.username }} <strong>{{ item.nick_name }}</strong>
                  </a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item label="执行人" help="指定工单执行人，一般为数据库管理员" has-feedback>
                <a-select
                  v-decorator="['executor', { rules: [{ required: true, message: '请选择工单执行人' }] }]"
                  placeholder="请选择工单执行人"
                  mode="multiple"
                  allowClear
                  show-search
                >
                  <a-select-option
                    v-for="(item, index) in users"
                    :key="index"
                    :label="item.username"
                    :value="item.username"
                  >
                    {{ item.username }} <strong>{{ item.nick_name }}</strong>
                  </a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item label="复核人" help="工单执行完成后，对结果进行复核的人员" has-feedback>
                <a-select
                  v-decorator="['reviewer', { rules: [{ required: true, message: '请选择工单复核人' }] }]"
                  placeholder="请选择工单复核人"
                  mode="multiple"
                  allowClear
                  show-search
                >
                  <a-select-option
                    v-for="(item, index) in users"
                    :key="index"
                    :label="item.username"
                    :value="item.username"
                  >
                    {{ item.username }} <strong>{{ item.nick_name }}</strong>
                  </a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item label="抄送人" has-feedback>
                <a-select
                  v-decorator="['cc', { rules: [{ required: false, message: '请选择工单抄送人' }] }]"
                  placeholder="请选择工单抄送人"
                  mode="multiple"
                  allowClear
                  show-search
                >
                  <a-select-option
                    v-for="(item, index) in users"
                    :key="index"
                    :label="item.username"
                    :value="item.username"
                  >
                    {{ item.username }} <strong>{{ item.nick_name }}</strong>
                  </a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item :wrapper-col="{ span: 12, offset: 6 }">
                <a-button key="submit" type="primary" :loading="loading" @click="onSubmit"> 提交 </a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>
        <!-- 右侧编辑区域 -->
        <a-col class="gutter-row" :span="16">
          <a-card>
            <a-alert message="支持多条SQL语句，每条SQL语句须以 ; 结尾" banner closable />
            <div style="margin-bottom: 5px; margin-top: 4px">
              <a-button type="dashed" icon="thunderbolt" @click="formatSQL()">格式化</a-button>
              <a-button
                type="dashed"
                icon="safety"
                style="margin-right: 6px"
                :disabled="checkBtnStatus"
                @click="syntaxCheck()"
                >语法检查</a-button
              >
            </div>
            <codemirror ref="myCm" v-model="code" :options="cmOptions" @ready="onCmReady"></codemirror>
          </a-card>
        </a-col>
      </a-row>
    </a-card>
    <a-card v-if="visibleAuditResult">
      <AuditResultComponent ref="AuditResultComponent"></AuditResultComponent>
    </a-card>
  </div>
</template>

<script>
// mode
import 'codemirror/mode/sql/sql.js'
// addon
import 'codemirror/addon/selection/active-line'
import 'codemirror/addon/display/autorefresh'
// 提示和自动补全
import 'codemirror/addon/hint/show-hint'
import 'codemirror/addon/hint/show-hint.css'
import 'codemirror/addon/hint/anyword-hint'
import 'codemirror/addon/hint/sql-hint'
import 'codemirror/addon/comment/comment'
import 'codemirror/addon/edit/matchbrackets'
import 'codemirror/addon/edit/closebrackets'
// 编辑器类型
import 'codemirror/keymap/sublime'

import notification from 'ant-design-vue/es/notification'

// format
import { format } from 'sql-formatter'
import {
  getEnvironmentsApi,
  getInstancesApi,
  getSchemasApi,
  getUsersApi,
  createOrdersApi,
  syntaxCheckApi,
} from '@/api/orders'

import AuditResultComponent from './AuditResult.vue'

const dbTypes = ['MySQL', 'TiDB']
export default {
  components: {
    AuditResultComponent,
  },
  computed: {
    codemirror() {
      return this.$refs.myCm.codemirror
    },
  },
  data() {
    return {
      loading: false,
      checkBtnStatus: false,
      visibleAuditResult: false,
      sqlType: '',
      cardTitle: '',
      dbTypes,
      environments: [],
      instances: [],
      schemas: [],
      users: [],
      form: this.$form.createForm(this, { name: 'commit' }),
      code: '',
      cmOptions: {
        mode: 'text/x-mysql',
        indentUnit: 2,
        tabSize: 2,
        indentWithTabs: true,
        smartIndent: true,
        autoRefresh: true,
        lineNumbers: true,
        styleActiveLine: true,
        autoCloseBrackets: true,
        matchBrackets: true,
        lineWrapping: true, // 自动换行
        autofocus: true,
        resetSelectionOnContextMenu: false,
        showCursorWhenSelecting: true,
        keyMap: 'sublime', // 编辑器模式
      },
    }
  },
  methods: {
    // 获取SQL类型
    getSqlType() {
      var urlSuffix = this.$route.path.split('/').at([-1]).toUpperCase()
      this.sqlType = urlSuffix
      this.cardTitle = `提交${this.sqlType}工单`
    },
    // 自动补全
    onCmReady(cm) {
      // 获取form表单的高度，将codemirror的高度设置和表单一致
      var formHeight = this.$refs.formHeight.$el.offsetHeight - 55
      cm.setSize('height', `${formHeight}px`)
      // 设置自动补全
      cm.on('keypress', () => {
        cm.showHint({ completeSingle: false })
      })
    },
    // 格式化SQL
    formatSQL() {
      var sql = this.codemirror.getValue()
      this.codemirror.setValue(format(sql, { language: 'mysql' }))
    },
    // 获取环境
    getEnvironments() {
      getEnvironmentsApi({ is_page: false })
        .then((res) => {
          this.environments = res.data
        })
        .catch((_error) => {})
    },
    changeDBType() {
      this.form.resetFields(['environment', 'instance_id', 'schema'])
    },
    // Change环境
    changeEnvs(value) {
      // Change环境时清空指定的字段
      this.form.resetFields(['instance_id', 'schema'])
      // 获取指定环境的实例
      var params = {
        id: value,
        db_type: this.form.getFieldValue('db_type'),
        is_page: false,
      }
      getInstancesApi(params)
        .then((res) => {
          this.instances = res.data
        })
        .catch((_error) => {})
    },
    // Change实例
    changeIns(value) {
      // Change实例时清空指定的字段
      this.form.resetFields(['schema'])
      // 获取指定实例的Schemas
      var params = {
        instance_id: value,
        is_page: false,
      }
      getSchemasApi(params)
        .then((res) => {
          this.schemas = res.data
        })
        .catch((_error) => {})
    },
    // 获取审核/复核/抄送人
    getUsers() {
      var params = {
        is_page: false,
      }
      getUsersApi(params)
        .then((res) => {
          this.users = res.data
        })
        .catch((_error) => {})
    },
    // 语法检查
    syntaxCheck() {
      var content = this.codemirror.getValue()
      // 检查提交的内容是否全
      if (content == '') {
        this.$notification.warning({
          message: '警告',
          description: '输入内容不能为空',
        })
        return
      }
      if (this.form.getFieldValue('environment') === undefined) {
        this.$notification.warning({
          message: '警告',
          description: '请选择环境',
        })
        return
      }
      if (this.form.getFieldValue('instance_id') === undefined) {
        this.$notification.warning({
          message: '警告',
          description: '请选择实例',
        })
        return
      }
      if (this.form.getFieldValue('schema') === undefined) {
        this.$notification.warning({
          message: '警告',
          description: '请选择库名',
        })
        return
      }
      this.$notification.info({
        message: '提示',
        description: '开始执行语法检查',
      })
      const data = {
        db_type: this.form.getFieldValue('db_type'),
        sql_type: this.sqlType,
        instance_id: this.form.getFieldValue('instance_id'),
        schema: this.form.getFieldValue('schema'),
        content: this.codemirror.getValue(),
      }
      // 滚动到页面的底部
      this.$nextTick(() => {
        document.scrollingElement.scrollTop = document.scrollingElement.scrollHeight
      })
      this.checkBtnStatus = true
      syntaxCheckApi(data)
        .then((res) => {
          if (res.code === '0000') {
            if (data['sql_type'].toLowerCase() == 'export') {
              this.$notification.success({
                message: '成功',
                description: '语法检查通过，您可以提交工单了，(⊙o⊙)',
              })
            } else {
              this.visibleAuditResult = true
              this.$nextTick(() => {
                this.$refs.AuditResultComponent.renderData(res.data)
              })
            }
          } else {
            this.$notification.error({
              message: '错误',
              description: res.message,
            })
          }
        })
        .catch((_error) => {})
        .finally(() => {
          this.checkBtnStatus = false
        })
    },
    // 提交工单
    onSubmit(e) {
      this.loading = true
      this.disableCommit = false
      e.preventDefault()
      this.form.validateFields((err, values) => {
        if (!err) {
          var content = this.codemirror.getValue()
          if (content != '') {
            values['is_restrict_access'] = values['is_restrict_access'] === 'YES'
            values['content'] = content
            values['sql_type'] = this.sqlType
            createOrdersApi(values)
              .then((res) => {
                if (res.code === '0000') {
                  this.$router.push('/orders/list')
                } else {
                  this.$notification.error({
                    message: '错误',
                    description: '提交的SQL内容不能为空',
                  })
                }
              })
              .catch((_error) => {})
              .finally(() => {
                this.disableCommit = false
              })
          } else {
            this.$notification.error({
              message: '错误',
              description: '提交的SQL内容不能为空',
            })
          }
        }
      })
      this.loading = false
      this.disableCommit = false
    },
  },
  mounted() {
    this.getSqlType()
    this.getEnvironments()
    this.getUsers()
  },
  watch: {
    $route() {
      this.getSqlType()
    },
  },
}
</script>

<style lang="less" scoped>
::v-deep .ant-card-body {
  padding: 12px;
}

::v-deep .ant-form {
  margin-bottom: 6px;
}
</style>
