<template>
  <a-card>
    <page-header-wrapper v-if="orderDetail" :title="orderDetail.title">
      <template v-slot:content>
        <a-descriptions size="small" :column="isMobile ? 1 : 2">
          <a-descriptions-item label="申请人">{{ orderDetail.applicant }}</a-descriptions-item>
          <a-descriptions-item label="备注">{{ orderDetail.remark }}</a-descriptions-item>
          <a-descriptions-item label="工单环境">
            <span style="color: red">{{ orderDetail.env_id }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="工单类型">{{ orderDetail.sql_type }}</a-descriptions-item>
          <a-descriptions-item label="DB类型">{{ orderDetail.display_rds_category }}</a-descriptions-item>
          <a-descriptions-item label="DB实例">{{ orderDetail.host + ':' + orderDetail.port }}</a-descriptions-item>
          <a-descriptions-item label="库名">
            <span style="color: blue">{{ orderDetail.database }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ orderDetail.created_at }}</a-descriptions-item>
          <a-descriptions-item label="需求">{{ orderDetail.demand }}</a-descriptions-item>
        </a-descriptions>
      </template>

      <template v-slot:extraContent>
        <a-row class="status-list">
          <a-col :xs="24" :sm="24">
            <div class="text">状态</div>
            <div class="heading">
              <h2>{{ orderDetail.progress }}</h2>
            </div>
          </a-col>
        </a-row>
      </template>

      <!-- actions -->
      <template v-slot:extra>
        <a-button-group style="margin-right: 4px">
          <a-button type="dashed" @click="showModal" :disabled="btnStatus.btnDisabled" icon="retweet">
            {{ orderDetail.progress | btnTitle }}
          </a-button>

          <a-button @click="showHookModal" v-if="orderDetail.progress === '已复核'" icon="link">钩子</a-button>

          <a-button type="dashed" @click="showCloseModal" :disabled="btnStatus.closeDisabled" icon="close-circle"
            >关闭工单</a-button
          >

          <a-button type="dashed" @click="generateSqlOrdersTasks" :loading="executeLoading" icon="thunderbolt"
            >执行工单</a-button
          >
          <a-button
            type="dashed"
            v-if="['处理中', '已完成', '已复核', '已勾住'].includes(orderDetail.progress)"
            @click="showTasksDrawer"
            icon="eye"
            >子任务详情</a-button
          >
          <a-button type="dashed" @click="refresh" :loading="loading" icon="sync">刷新</a-button>
        </a-button-group>
        <!-- 其他操作model -->
        <a-modal title="请输入[可选]" v-model="visible">
          <a-textarea v-model="confirmMsg" rows="3" :autoSize="{ minRows: 3, maxRows: 5 }" />
          <template slot="footer">
            <a-button key="back" @click="handleCancel">{{ confirmBtnTips.cancelText }}</a-button>
            <a-button key="submit" type="primary" :loading="loading" @click="handleOk">
              {{ confirmBtnTips.okText }}
            </a-button>
          </template>
        </a-modal>
        <!-- close model -->
        <a-modal title="请输入[可选]" v-model="closeVisible">
          <a-textarea v-model="confirmMsg" rows="3" :autoSize="{ minRows: 3, maxRows: 5 }" />
          <template slot="footer">
            <a-button key="back" @click="handleCloseCancel">关闭</a-button>
            <a-button key="submit" type="primary" :loading="loading" @click="handleCloseOk">提交</a-button>
          </template>
        </a-modal>
        <!-- hook model -->
        <a-modal title="请选择目标环境[钩子]" v-model="hookVisible" width="35%">
          <el-form :model="ruleForm" ref="ruleForm" label-width="150px" size="small">
            <el-form-item label="当前工单">
              <el-input v-model="ruleForm.title" readonly placeholder="请输入标题" style="width: 95%" />
            </el-form-item>

            <el-form-item label="当前库">
              <el-input v-model="ruleForm.current_database" readonly placeholder="请输入需求描述" style="width: 95%" />
            </el-form-item>

            <el-form-item label="目标环境">
              <el-select
                v-model="ruleForm.env_id"
                style="width: 95%"
                placeholder="请选择工单环境"
                @change="changeEnvs"
                value
              >
                <el-option
                  v-for="item in sql_envs"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                  :disabled="item.disabled"
                ></el-option>
              </el-select>
            </el-form-item>

            <el-form-item label="目标库">
              <el-select
                v-model="ruleForm.database"
                style="width: 95%"
                clearable
                filterable
                placeholder="请选择目标数据库"
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

            <el-form-item label="审核状态">
              <a-switch
                checked-children="重置审核状态为：待审批"
                un-checked-children="继承审核状态为：已批准"
                default-checked
                @change="onRestChange"
              />
            </el-form-item>

            <el-form-item label="备注">
              <el-select v-model="ruleForm.remark" style="width: 95%" placeholder="请选择合适的备注" value>
                <el-option v-for="item in remarks" :key="item" :label="item" :value="item"></el-option>
              </el-select>
            </el-form-item>
          </el-form>
          <template slot="footer">
            <a-button key="back" @click="handleHookCancel">关闭</a-button>
            <a-button key="submit" type="primary" :loading="hookLoading" @click="handleHookOk">提交</a-button>
          </template>
        </a-modal>
      </template>

      <a-card :bordered="false" title="事件进度" style="margin-top: -10px">
        <a-steps :current="currentStatus" size="small">
          <a-step title="创建工单"></a-step>
          <a-step title="审核中"></a-step>
          <a-step title="已审核"></a-step>
          <a-step title="处理中"></a-step>
          <a-step title="已完成"></a-step>
          <a-step title="已复核"></a-step>
          <a-step title="已勾住"></a-step>
        </a-steps>

        <a-card type="inner" title="事件状态" style="margin-top: 18px">
          <a-descriptions
            size="small"
            :col="4"
            v-for="k of JSON.parse(orderDetail.auditor)"
            :key="`auditor` + k.user + k.status"
          >
            <a-descriptions-item label="审核人">
              {{ k.user }}
              <span v-if="k.is_superuser === 1">(超级审核人)</span>
            </a-descriptions-item>
            <a-descriptions-item label="状态">
              <a-tag v-if="k.status === 0" color="red">未审核</a-tag>
              <a-tag v-else-if="k.status === 1" color="green">已审核</a-tag>
              <a-tag v-else color="blue">已驳回</a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="审核时间">{{ k.time }}</a-descriptions-item>
            <a-descriptions-item label="附加信息">{{ k.msg }}</a-descriptions-item>
          </a-descriptions>

          <a-divider style="margin: 5px" />
          <a-descriptions
            size="small"
            :col="4"
            v-for="k of JSON.parse(orderDetail.reviewer)"
            :key="`reviewer` + k.user + k.status"
          >
            <a-descriptions-item label="复核人">{{ k.user }}</a-descriptions-item>
            <a-descriptions-item label="状态">
              <a-tag v-if="k.status === 0" color="red">未复核</a-tag>
              <a-tag v-else color="green">已复核</a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="复核时间">{{ k.time }}</a-descriptions-item>
            <a-descriptions-item label="附加信息">{{ k.msg }}</a-descriptions-item>
          </a-descriptions>

          <a-divider style="margin: 5px" />
          <a-descriptions
            size="small"
            :col="4"
            v-for="k of JSON.parse(orderDetail.executor)"
            :key="`executor` + k.user + k.status"
          >
            <a-descriptions-item label="执行人">{{ k.user }}</a-descriptions-item>
            <a-descriptions-item label="状态">-</a-descriptions-item>
            <a-descriptions-item label="执行时间">{{ k.time }}</a-descriptions-item>
            <a-descriptions-item label="附加信息">{{ k.msg }}</a-descriptions-item>
          </a-descriptions>

          <div v-if="orderDetail.progress === '已关闭'">
            <a-divider style="margin: 5px" />
            <a-descriptions
              size="small"
              :col="4"
              v-for="k of JSON.parse(orderDetail.closer)"
              :key="`closer` + k.user + k.status"
            >
              <a-descriptions-item label="关闭人">{{ k.user }}</a-descriptions-item>
              <a-descriptions-item label="状态">-</a-descriptions-item>
              <a-descriptions-item label="关闭时间">{{ k.time }}</a-descriptions-item>
              <a-descriptions-item label="附加信息">{{ k.msg }}</a-descriptions-item>
            </a-descriptions>
          </div>
        </a-card>
      </a-card>

      <a-card title="工单内容" :bordered="false" style="margin-top: -10px">
        <codemirror ref="myCm" v-model="code" :options="cmOptions" @ready="onCmReady"></codemirror>
      </a-card>

      <a-drawer
        title="子任务执行详情"
        width="60%"
        placement="right"
        :closable="false"
        :visible="visibleDrawer"
        @close="onCloseDrawer"
      >
        <drawerTasksPreview :order_id="this.orderDetail.id" />
      </a-drawer>
    </page-header-wrapper>
  </a-card>
</template>

<script>
import { baseMixin } from '@/store/app-mixin'
import { SqlRemark } from '@/utils/sql'

import {
  getSqlOrdersDetail,
  opSqlOrders,
  getDbSchemas,
  generateSqlOrdersExecuteTasks,
  HookSqlOrders,
  getDbEnvironment,
} from '@/api/sql'

import 'codemirror/mode/sql/sql.js'
import 'codemirror/theme/ambiance.css'
import 'codemirror/addon/display/autorefresh'

import drawerTasksPreview from './preview.vue'

// 模态框按钮
const ConfirmBtnTips = { okText: '确认', cancelText: '取消', action: '' }
const CloseBtnTips = { okText: '提交', cancelText: '关闭', action: 'close' }
const BtnStatus = {
  btnDisabled: false,
  closeDisabled: false,
}

export default {
  name: 'detail',
  mixins: [baseMixin],
  components: {
    drawerTasksPreview,
  },
  data() {
    return {
      visible: false,
      visibleDrawer: false,
      closeVisible: false,
      loading: false,
      executeLoading: false,
      // 钩子
      envs: [],
      hookLoading: false,
      hookVisible: false,
      btnStatus: BtnStatus,
      confirmMsg: '',
      orderDetail: '',
      toUrl: '',
      confirmBtnTips: ConfirmBtnTips,
      code: '',
      cmOptions: {
        mode: 'text/x-mysql',
        indentUnit: 2,
        tabSize: 2,
        indentWithTabs: true,
        smartIndent: true,
        autoRefresh: true,
        lineNumbers: true,
        readOnly: true,
        focuse: false,
      },
      schemas: [],
      sql_envs: [],
      resetAuditStatus: 'ON',
      remarks: SqlRemark,
      ruleForm: {
        title: '',
        current_database: '',
        env_id: '',
        database: '',
        remark: '', // 备注
      },
    }
  },
  methods: {
    getEnvs() {
      getDbEnvironment.then((response) => {
        this.envs = response.data
      })
    },
    showTasksDrawer() {
      this.visibleDrawer = true
    },
    onCloseDrawer() {
      this.visibleDrawer = false
    },
    getOrderDetail() {
      getSqlOrdersDetail(this.$route.params.order_id)
        .then((response) => {
          this.orderDetail = response.data
        })
    },
    // 变更环境，获取schemas
    changeEnvs(value) {
      this.ruleForm.database = '' //切换环境时，置空已选择的库名
      const params = {
        env_id: value,
        use_type: 0,
        rds_category: this.orderDetail.rds_category,
      }
      getDbSchemas(params).then((response) => {
        this.schemas = response.data
      })
    },
    onCmReady(cm) {
      cm.setSize('height', `550px`)
      cm.setValue(this.orderDetail.contents)
    },
    refresh() {
      this.loading = true
      this.getOrderDetail()
      setTimeout(() => {
        this.loading = false
      }, 1000)
    },
    handleCommit(btn, action) {
      if (this.confirmMsg.length > 128) {
        this.$message.error('提交失败，消息长度不能超过128个字符')
        this.hideModal()
        return false
      }
      const commitData = {
        action: action,
        msg: this.confirmMsg,
        btn: btn,
        pk: this.orderDetail.id,
      }
      opSqlOrders(commitData)
        .then((response) => {
          if (response.code === '0000') {
            this.$message.success(response.message)
          } else {
            this.$message.error(JSON.stringify(response.message))
          }
        })
        .finally(() => {
          this.hideModal()
          this.refresh()
        })
    },
    generateSqlOrdersTasks() {
      this.executeLoading = true
      const data = { id: this.orderDetail.id }
      generateSqlOrdersExecuteTasks(data)
        .then((response) => {
          if (response.code === '0000') {
            this.$router.push(`/sqlorders/tasks/${response.data}`)
          } else {
            this.$message.error(response.message)
          }
        })
        .finally(() => {
          this.executeLoading = false
        })
    },
    showModal() {
      this.visible = true
    },
    hideModal() {
      this.visible = false
    },
    handleOk(e) {
      this.handleCommit('ok', ConfirmBtnTips.action)
      this.hideModal()
    },
    handleCancel() {
      this.handleCommit('cancel', ConfirmBtnTips.action)
      this.hideModal()
    },
    // hook
    showHookModal() {
      this.ruleForm.title = this.orderDetail.title
      this.ruleForm.current_database = this.orderDetail.database
      // 删除当前环境
      this.sql_envs = []
      this.envs.map((item) => {
        if (item.name != this.orderDetail.env_id) {
          this.sql_envs.push(item)
        }
      })
      this.hookVisible = true
    },
    hideHookModal() {
      this.hookVisible = false
    },
    handleHookOk() {
      const data = {
        id: this.orderDetail.id,
        reset: this.resetAuditStatus,
        ...this.ruleForm,
      }
      HookSqlOrders(data).then((response) => {
        if (response.code === '0000') {
          this.$router.push(`/sqlorders/list`)
        } else {
          this.$message.error(response.message)
        }
      })
    },
    handleHookCancel() {
      this.hideHookModal()
    },
    // close
    showCloseModal() {
      this.closeVisible = true
    },
    hideCloseModal() {
      this.closeVisible = false
    },
    handleCloseOk(e) {
      this.handleCommit('ok', CloseBtnTips.action)
      this.hideCloseModal()
    },
    handleCloseCancel() {
      this.hideCloseModal()
    },
    onRestChange(checked) {
      this.resetAuditStatus = checked ? 'ON' : 'OFF'
    },
  },
  mounted() {
    this.getEnvs()
    this.getOrderDetail()
  },
  filters: {
    btnTitle(progress) {
      if (progress === '待审核') {
        BtnStatus.btnDisabled = false
        BtnStatus.closeDisabled = false
        ConfirmBtnTips.okText = '同意'
        ConfirmBtnTips.cancelText = '驳回'
        ConfirmBtnTips.action = 'approve' // 审核
        return '审核'
      } else if (['已批准', '处理中'].includes(progress)) {
        BtnStatus.btnDisabled = false
        BtnStatus.closeDisabled = false
        ConfirmBtnTips.okText = '处理中'
        ConfirmBtnTips.cancelText = '执行完成'
        ConfirmBtnTips.action = 'feedback' // 审核
        return '反馈'
      } else if (progress === '已完成') {
        BtnStatus.btnDisabled = false
        BtnStatus.closeDisabled = true
        ConfirmBtnTips.okText = '已复核'
        ConfirmBtnTips.cancelText = '关闭窗口'
        ConfirmBtnTips.action = 'review' // 审核
        return '复核'
      } else if (['已关闭', '已勾住', '已驳回', '已复核'].includes(progress)) {
        BtnStatus.btnDisabled = true
        BtnStatus.closeDisabled = true
        return '完成'
      }
    },
  },
  computed: {
    codemirror() {
      return this.$refs.myCm.codemirror
    },
    currentStatus() {
      if (this.orderDetail.progress === '已驳回') {
        return 2 // currentStatus=2表示已审核
      } else if (this.orderDetail.progress === '已关闭') {
        return 4 // currentStatus=4表示已完成
      } else if (this.orderDetail.progress === '待审核') {
        return 1 // currentStatus=1表示审核中
      } else if (this.orderDetail.progress === '已批准') {
        return 2 // currentStatus=2表示已审核
      } else if (this.orderDetail.progress === '处理中') {
        return 3 // currentStatus=3表示处理中
      } else if (this.orderDetail.progress === '已完成') {
        return 4 // currentStatus=4表示已完成
      } else if (this.orderDetail.progress === '已复核') {
        return 5 // currentStatus=5表示已核对
      } else {
        return 6 // currentStatus=6表示已勾住
      }
    },
  },
}
</script>

<style lang="less" scoped>
.CodeMirror {
  border: 2px solid #eee;
  font-family: 'JetBrains Mono NL', Menlo, Monaco, Consolas, 'Lucida Console', 'Courier New', monospace;
  min-height: 100px;
  /* 支持上下拉伸 */
  resize: vertical;
  overflow: y !important;
}
.CodeMirror pre.CodeMirror-placeholder {
  color: #999;
}
.mobile {
  .detail-layout {
    margin-left: unset;
  }
  .status-list {
    text-align: left;
  }
}
</style>
