<template>
  <div>
    <a-card :title="cardTitle">
      <a-row>
        <a-col :span="16" :push="8">
          <div style="margin-bottom: 5px">
            <a-button type="dashed" icon="thunderbolt" style="margin-right: 6px" @click="formatSQL()">格式化</a-button>
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
        </a-col>
        <a-col :span="8" :pull="16">
          <el-form :model="ruleForm" :rules="rules" ref="ruleForm" label-width="100px" size="small">
            <el-form-item label="标题" prop="title" key="bindTitle">
              <el-input v-model="ruleForm.title" placeholder="请输入标题" style="width: 95%" />
            </el-form-item>

            <el-form-item label="需求" prop="demand">
              <el-input
                v-model="ruleForm.demand"
                type="textarea"
                :rows="2"
                placeholder="请输入需求描述"
                style="width: 95%"
              />
            </el-form-item>

            <el-form-item key="bindIsHide">
              <template slot="label">
                <span stryle="position: relative">
                  <span>隐藏数据</span>
                  <el-tooltip placement="right-end" effect="light">
                    <div slot="content">
                      <span>
                        开启后
                        <br />仅工单的提交人、审核人、复核人和DBA <br />可以查看工单内容
                      </span>
                    </div>
                    <i class="el-icon-question table-msg" />
                  </el-tooltip>
                </span>
              </template>
              <a-switch style="margin-bottom: 1px" defaultChecked @change="onMyChange" />
            </el-form-item>

            <el-form-item label="版本" v-if="isShow" key="bindVersion">
              <el-select v-model="ruleForm.version" style="width: 95%" placeholder="请选择上线版本" value>
                <el-option v-for="item in versions" :key="item.id" :label="item.version" :value="item.id"></el-option>
              </el-select>
            </el-form-item>

            <el-form-item label="备注" prop="remark" key="bindRemark">
              <el-select v-model="ruleForm.remark" style="width: 95%" placeholder="请选择合适的备注" value>
                <el-option v-for="item in remarks" :key="item" :label="item" :value="item"></el-option>
              </el-select>
            </el-form-item>

            <el-form-item label="DB类别" prop="rds_category" key="bindRdsCategory">
              <el-select
                v-model="ruleForm.rds_category"
                style="width: 95%"
                placeholder="请选择数据库类型"
                @change="changeRdsCategory"
                value
              >
                <el-option
                  v-for="item in rds_category"
                  :key="item.key"
                  :label="item.value"
                  :value="item.key"
                  :disabled="item.disabled"
                ></el-option>
              </el-select>
            </el-form-item>

            <el-form-item label="环境" prop="env_id" key="bindEnvId">
              <el-select
                v-model="ruleForm.env_id"
                style="width: 95%"
                placeholder="请选择工单环境"
                @change="changeEnvs"
                value
              >
                <el-option
                  v-for="item in envs"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                  :disabled="item.disabled"
                ></el-option>
              </el-select>
            </el-form-item>

            <el-form-item label="库名" prop="database" key="bindDatabase">
              <el-select
                v-model="ruleForm.database"
                style="width: 95%"
                clearable
                filterable
                placeholder="请选择数据库"
                value
              >
                <el-option
                  v-for="item in schemas"
                  :key="item.id"
                  :label="`${item.comment}-${item.schema}`"
                  :value="`${item.cid}__${item.schema}`"
                ></el-option>
              </el-select>
            </el-form-item>

            <el-form-item label="文件格式" prop="file_format" v-if="!isShow" key="bindFileFormart">
              <el-select
                v-model="ruleForm.file_format"
                style="width: 95%"
                clearable
                filterable
                placeholder="请选择文件格式"
                value
              >
                <el-option v-for="item in format" :key="item.key" :label="item.value" :value="item.key"></el-option>
              </el-select>
            </el-form-item>

            <el-form-item label="审核人" prop="auditor" key="bindAuditor">
              <el-select
                multiple
                :multiple-limit="3"
                clearable
                filterable
                v-model="ruleForm.auditor"
                style="width: 95%"
                placeholder="请选择工单审核人"
                value
              >
                <el-option
                  v-for="item in users"
                  :key="item.uid"
                  :label="`${item.displayname}[${item.username}]`"
                  :value="item.username"
                ></el-option>
              </el-select>
            </el-form-item>

            <el-form-item label="复核人" prop="reviewer" key="bindReviewer">
              <el-select
                multiple
                :multiple-limit="3"
                clearable
                filterable
                v-model="ruleForm.reviewer"
                style="width: 95%"
                placeholder="请选择工单复核人"
                value
              >
                <el-option
                  v-for="item in users"
                  :key="item.uid"
                  :label="`${item.displayname}[${item.username}]`"
                  :value="item.username"
                ></el-option>
              </el-select>
            </el-form-item>

            <el-form-item label="抄送人" key="bindCc">
              <el-select
                multiple
                :multiple-limit="5"
                clearable
                filterable
                v-model="ruleForm.email_cc"
                style="width: 95%"
                placeholder="请选择工单抄送人"
                value
              >
                <el-option
                  v-for="item in users"
                  :key="item.uid"
                  :label="`${item.displayname}[${item.username}]`"
                  :value="item.username"
                ></el-option>
              </el-select>
            </el-form-item>

            <el-form-item style="text-align: left">
              <el-button type="primary" :loading="isDisabledCommit" @click="submitForm('ruleForm')">提交</el-button>
              <el-button @click="resetForm('ruleForm')">重置</el-button>
            </el-form-item>
          </el-form>
        </a-col>
      </a-row>
    </a-card>
    <a-card v-if="visibleAuditResult">
      <a-row>
        <a-table
          :columns="table.columns"
          :dataSource="table.data"
          :pagination="pagination"
          :loading="tableLoading"
          :rowClassName="setRowClass"
          :rowKey="(record) => record.order_id"
          @change="handleTableChange"
          size="middle"
        >
          <span slot="error_level" slot-scope="text">
            <span v-if="text === 0">成功</span>
            <span v-else-if="text === 1">警告</span>
            <span v-else>错误</span>
          </span>
        </a-table>
      </a-row>
    </a-card>
    <div>
      <a-modal v-model="tidbVisible" title="TiDB注意事项" width="55%" @ok="handleTiDBOk">
        <div style="font-size: 12px">
          <el-divider content-position="left">DML事务</el-divider>
          <p>
            TiDB单条DML语句最大支持的事务为3W,
            若是更新（DELETE/UPDATE）超过了3W条记录，需要拆分为多条SQL语句，每条SQL后面加上LIMIT 20000。
          </p>
          <h3>例子：</h3>
          <h4>原始SQL：</h4>
          <p>UPDATE TEST1 SET NAME='XXX' WHERE I_STATUS = 2;</p>

          <h4>改写为：</h4>
          <p>UPDATE TEST1 SET NAME='XXX' WHERE I_STATUS = 2 LIMIT 20000;</p>
          <p>UPDATE TEST1 SET NAME='XXX' WHERE I_STATUS = 2 LIMIT 20000;</p>

          <el-divider content-position="left">DML备份</el-divider>
          <p>TiDB不支持生成回滚语句、TiDB不支持生成回滚语句、TiDB不支持生成回滚语句</p>

          <el-divider content-position="left">DDL语句</el-divider>
          <p>TiDB的ALTER语句不支持聚合写法，MODIFY/CHANGE/ADD等需要拆分</p>
          <h3>例子：</h3>
          <h4>原始SQL：</h4>
          <p>
            ALTER TABLE TEST1 ADD COL1 CHAR(10) NOT NULL DEFAULT '' COMMENT 'XX',ADD COL2 CHAR(10) NOT NULL DEFAULT ''
            COMMENT 'XX';
          </p>
          <h4>改写为：</h4>
          <p>ALTER TABLE TEST1 ADD COL1 CHAR(10) NOT NULL DEFAULT '' COMMENT 'XX';</p>
          <p>ALTER TABLE TEST1 ADD COL2 CHAR(10) NOT NULL DEFAULT '' COMMENT 'XX';</p>
        </div>
      </a-modal>
    </div>
  </div>
</template>

<script>
import { SqlRemark, rdsCategory, fileFormat } from '@/utils/sql'
import {
  getUsers,
  getDbSchemas,
  getDbEnvironment,
  getReleaseVersions,
  commitSqlOrders,
  incepSyntaxCheck,
} from '@/api/sql'
import sqlFormat from 'sql-formatter'

import 'codemirror/theme/ambiance.css'
import 'codemirror/mode/javascript/javascript'
import 'codemirror/mode/sql/sql.js'
import 'codemirror/addon/selection/active-line'
import 'codemirror/addon/edit/matchbrackets'
import 'codemirror/addon/edit/closebrackets'
import 'codemirror/addon/display/autorefresh'

// 提示和自动补全
import 'codemirror/addon/hint/show-hint'
import 'codemirror/addon/hint/show-hint.css'
import 'codemirror/addon/hint/anyword-hint'
import 'codemirror/addon/hint/sql-hint'

// 编辑器类型
import 'codemirror/keymap/sublime'

export default {
  data() {
    const selectAuditorChecker = (rule, value, callback) => {
      if (value.length < 1) {
        callback(new Error('请至少选择1个审核人'))
      } else {
        callback()
      }
    }
    const selectReviewerChecker = (rule, value, callback) => {
      if (value.length < 1) {
        callback(new Error('请至少选择1个复核人'))
      } else {
        callback()
      }
    }
    const check_rds_category = (rule, value, callback) => {
      if (!value) {
        callback(new Error('请选择工单环境'))
      }
      if (!this.ruleForm.rds_category) {
        callback(new Error('请先选择DB类别'))
      } else {
        callback()
      }
    }
    return {
      isDisabledCommit: false,
      checkBtnStatus: false,
      visibleAuditResult: false,
      tidbVisible: false,
      cardTitle: '',
      sqltype: '',
      code: '',
      schemas: [],
      users: [],
      remarks: SqlRemark,
      format: fileFormat,
      versions: [],
      rds_category: rdsCategory,
      envs: [],
      isShow: true,
      ruleForm: {
        title: '', // 标题
        demand: '', // 需求
        version: '', // 上线版本号
        is_hide: 'OFF', //是否隐藏数据
        remark: '', // 备注
        file_format: 'csv',
        rds_category: 1, // 数据库类别
        env_id: '', // 环境
        database: '', // 库名
        auditor: [], // 审核人
        reviewer: [this.$store.getters.userInfo.username], // 复核人
        email_cc: [], // 抄送人
      },
      rules: {
        title: [
          { required: true, message: '请输入标题', trigger: 'blur' },
          {
            min: 3,
            max: 64,
            message: '长度在 3 到 64 个字符',
            trigger: 'blur',
          },
        ],
        demand: [
          { required: true, message: '请输入需求描述', trigger: 'blur' },
          {
            min: 3,
            max: 256,
            message: '长度在 3 到 256 个字符',
            trigger: 'blur',
          },
        ],
        file_format: [{ required: true, message: '请选择文件格式', trigger: 'change' }],
        remark: [{ required: true, message: '请选择备注', trigger: 'change' }],
        rds_category: [{ required: true, message: '请选择数据库类型', trigger: 'change' }],
        env_id: [{ required: true, validator: check_rds_category, trigger: 'change' }],
        database: [{ required: true, message: '请选择数据库', trigger: 'change' }],
        auditor: [
          {
            required: true,
            validator: selectAuditorChecker,
            trigger: 'change',
          },
        ],
        reviewer: [
          {
            required: true,
            validator: selectReviewerChecker,
            trigger: 'change',
          },
        ],
      },
      // table
      tableLoading: false,
      table: {
        columns: null,
        data: null,
      },
      // pagination
      pagination: {
        current: 1,
        pageSize: 10,
        total: 0,
        pageSizeOptions: ['10', '20', '50'],
        showSizeChanger: true,
      },
      // codemirror
      cmOptions: {
        mode: 'text/x-mysql',
        indentUnit: 2,
        tabSize: 2,
        indentWithTabs: true,
        smartIndent: true,
        autoRefresh: true,
        lineNumbers: true,
        matchBrackets: true, // 括号匹配
        styleActiveLine: true, //背景高亮
        autofocus: true,
        keyMap: 'sublime', // 编辑器模式
        autoCloseBrackets: true,
        autorefresh: true,
        viewportMargin: Infinity, // 解决切换空白的问题
      },
      // 审核规则
      auditRules: '',
    }
  },
  methods: {
    onMyChange(checked) {
      this.ruleForm.is_hide = checked ? 'ON' : 'OFF'
    },
    handleTiDBOk(e) {
      this.tidbVisible = false
    },
    getsqltype() {
      getDbEnvironment.then((response) => {
        const envs = response.data

        // 根据url变动，切换工单类型
        const urlList = this.$route.path.split('/')
        this.sqltype = urlList[urlList.length - 1].toUpperCase() // 此处需要转换为大写
        this.cardTitle = `提交${this.sqltype}工单(每条SQL必须以 ; 结尾)`

        // 当为导出工单时，隐藏部分
        this.isShow = true
        if (this.sqltype === 'EXPORT') {
          this.isShow = false
        }

        // // DDL工单进行过滤处理，不允许直接提交到生产环境
        // //
        // this.ruleForm.env_id = "";
        // this.ruleForm.database = "";
        // envs.map((item) => {
        //   item.disabled = false;
        //   if (this.sqltype === "DDL") {
        //     if (item.id === 2) {   // 数据库中生产环境对应的值为2
        //       item.disabled = true;
        //     }
        //   }
        //   return item;
        // });

        this.envs = envs
      })
    },
    // 操作表格
    handleTableChange(pager) {
      this.pagination.current = pager.current
      this.pagination.pageSize = pager.pageSize
    },
    // 自动补全
    onCmReady(cm) {
      // 获取element form表单的高度，将codemirror的高度设置和表单一致
      const formHeight = this.$refs.ruleForm.$el.offsetHeight - 30
      cm.setSize('height', `${formHeight}px`)
      cm.on('keypress', () => {
        cm.showHint()
      })
    },
    // 获取上线版本号
    getReleaseVersionsList() {
      getReleaseVersions.then((response) => {
        this.versions = response.data
      })
    },
    // 格式化SQL
    formatSQL() {
      const sqlContent = this.codemirror.getValue()
      this.codemirror.setValue(sqlFormat.format(sqlContent, { indent: '  ' }))
    },
    // 语法检查
    syntaxCheck() {
      this.$message.info('正在执行语法检测，请稍等')
      const params = {
        rds_category: this.ruleForm.rds_category,
        database: this.ruleForm.database,
        sqls: this.codemirror.getValue(),
        sql_type: this.sqltype,
      }

      // 滚动到页面的底部
      this.$nextTick(() => {
        document.scrollingElement.scrollTop = document.scrollingElement.scrollHeight
      })

      this.tableLoading = true
      this.checkBtnStatus = true
      this.visibleAuditResult = true

      incepSyntaxCheck(params)
        .then((response) => {
          if (response.code === '0000') {
            if (response.data.data.status === 0) {
              this.$message.success(response.message)
            } else {
              this.$message.error(response.message)
            }
            this.table.data = response.data.data.data
            this.table.columns = response.data.columns
          } else {
            this.$message.error(response.message)
          }
          this.tableLoading = false
          this.checkBtnStatus = false
        })
        .catch((error) => {
          this.$message.error('语法检查失败，错误码: ' + error.response.status)
          this.tableLoading = false
          this.checkBtnStatus = false
        })
    },
    // 设置行的颜色
    setRowClass(record) {
      if (record.error_level === 0) {
        return 'row-info'
      }
      if (record.error_level === 1) {
        return 'row-warn'
      }
      if (record.error_level === 2) {
        return 'row-error'
      }
    },
    // 变更数据库类型，清空已选择的环境和库名
    changeRdsCategory(value) {
      if (value === 2) {
        this.tidbVisible = true
      }
      this.ruleForm.env_id = ''
      this.ruleForm.database = ''
    },
    // 变更环境，获取schemas
    changeEnvs(value) {
      this.ruleForm.database = '' //切换环境时，置空已选择的库名
      const params = {
        env_id: value,
        use_type: 0,
        rds_category: this.ruleForm.rds_category,
      }
      if (this.sqltype === 'EXPORT') {
        params.use_type = 1
      }
      getDbSchemas(params).then((response) => {
        this.schemas = response.data
      })
    },
    // 获取用户
    getUsersList() {
      getUsers.then((response) => {
        this.users = response.data
      })
    },
    // 提交工单
    submitForm(formName) {
      this.isDisabledCommit = true
      setTimeout(() => {
        this.$refs[formName].validate((valid) => {
          if (valid) {
            // 简单判断输入是否为空
            const sqlContent = this.codemirror.getValue()
            if (!sqlContent) {
              this.$message.error('请输入要审核的SQL内容')
              this.isDisabledCommit = false
              return false
            }
            const commitData = {
              sql_type: this.sqltype,
              contents: sqlContent,
              ...this.ruleForm,
            }
            commitSqlOrders(commitData)
              .then((response) => {
                if (response.code === '0000') {
                  this.$router.push('/sqlorders/list')
                } else {
                  this.isDisabledCommit = false
                  this.$message.error(response.message)
                }
              })
              .catch((error) => {
                this.$message.error('提交失败，错误码: ' + error.response.status)
                this.isDisabledCommit = false
              })
          } else {
            this.isDisabledCommit = false
            return false
          }
        })
      }, 500)
    },
    resetForm(formName) {
      this.$refs[formName].resetFields()
      this.codemirror.setValue('')
    },
  },
  mounted() {
    this.getsqltype()
    this.getUsersList()
    this.getReleaseVersionsList()
  },
  computed: {
    codemirror() {
      return this.$refs.myCm.codemirror
    },
  },
  watch: {
    $route() {
      this.getsqltype()
    },
  },
}
</script>

<style>
.row-info {
  color: green;
}
.row-warn {
  color: orange;
}
.row-error {
  color: red;
}
</style>
