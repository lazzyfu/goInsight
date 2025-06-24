<template>
    <page-header-wrapper v-if="orderDetail" :title="orderDetail.title">
      <template v-slot:content>
        <a-descriptions size="small" :column="3">
          <a-descriptions-item label="申请人">{{ orderDetail.applicant }}</a-descriptions-item>
          <a-descriptions-item label="工单环境">
            <span style="color: red">{{ orderDetail.environment }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="DB类型">{{ orderDetail.db_type }}</a-descriptions-item>
          <a-descriptions-item label="工单类型">{{ orderDetail.sql_type }}</a-descriptions-item>
          <a-descriptions-item label="DB实例">{{ orderDetail.instance }}</a-descriptions-item>
          <a-descriptions-item label="库名">
            <span style="color: blue">{{ orderDetail.schema }}</span>
          </a-descriptions-item>
          <a-descriptions-item v-if="orderDetail.sql_type === 'EXPORT'" label="文件格式">{{
            orderDetail.export_file_format
          }}</a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ orderDetail.created_at }}</a-descriptions-item>
          <a-descriptions-item label="更新时间">{{ orderDetail.updated_at }}</a-descriptions-item>
        </a-descriptions>
      </template>
      <!-- 附加信息 -->
      <a-card :bordered="false" title="附加信息">
        <a-descriptions size="small" :column="1">
          <a-descriptions-item
            v-if="orderDetail.hook_order_id != '00000000-0000-0000-0000-000000000000'"
            label="Hook工单"
          >
            <router-link
              @click.native="refresh"
              target="_blank"
              :to="{ name: 'view.orders.detail', params: { order_id: orderDetail.hook_order_id } }"
              >{{ orderDetail.hook_order_id }}</router-link
            >
          </a-descriptions-item>
          <a-descriptions-item label="工单备注">
            <span v-if="orderDetail.remark">{{ orderDetail.remark }}</span>
            <span v-else>无</span></a-descriptions-item
          >
          <a-descriptions-item label="工单执行人">
            <span v-if="orderDetail.executor">
              <span v-for="k of orderDetail.executor" :key="k">
                <a-tag color="purple">{{ k }}</a-tag>
              </span>
            </span>
            <span v-else> 无 </span>
          </a-descriptions-item>
          <a-descriptions-item label="工单抄送人">
            <span v-if="orderDetail.cc">
              <span v-for="k of orderDetail.cc" :key="k">
                <a-tag>{{ k }}</a-tag>
              </span>
            </span>
            <span v-else> 无 </span>
          </a-descriptions-item>
        </a-descriptions>
      </a-card>
      <!-- 工单状态 -->
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
  
      <!-- 右上角按钮 -->
      <template v-slot:extra>
        <a-button-group style="margin-right: 4px">
          <a-button type="dashed" @click="showActionModal" :disabled="btnStatus.btnDisabled" icon="retweet">
            {{ orderDetail.progress | btnTitle }}
          </a-button>
          <!-- 触发钩子 -->
          <a-button @click="showHookModal" v-if="orderDetail.progress === '已复核'" icon="link">钩子</a-button>
          <!-- 关闭工单 -->
          <a-button type="dashed" @click="showCloseModal" :disabled="btnStatus.closeDisabled" icon="close-circle"
            >关闭工单</a-button
          >
          <!-- 执行工单 -->
          <a-button
            type="dashed"
            @click="generateSqlOrdersTasks"
            v-if="['已批准', '执行中', '已完成', '已复核'].includes(orderDetail.progress)"
            :loading="executeLoading"
            icon="thunderbolt"
            >执行工单</a-button
          >
          <!-- 刷新 -->
          <a-button type="dashed" @click="refresh" :loading="refreshLoading" icon="sync">刷新</a-button>
        </a-button-group>
        <!-- action modal -->
        <a-modal title="请输入附加信息" v-model="actionVisible">
          <a-textarea v-model="confirmMsg" allow-clear rows="3" :autoSize="{ minRows: 3, maxRows: 8 }" />
          <template slot="footer">
            <a-button key="back" @click="handleActionCancel">{{ confirmBtnTips.cancelText }}</a-button>
            <a-button key="submit" type="primary" :loading="loading" @click="handleActionOk">
              {{ confirmBtnTips.okText }}
            </a-button>
          </template>
        </a-modal>
        <!-- close model -->
        <a-modal title="请输入附加信息" v-model="closeVisible">
          <a-textarea v-model="confirmMsg" rows="3" :autoSize="{ minRows: 3, maxRows: 5 }" />
          <template slot="footer">
            <a-button key="back" @click="handleCloseCancel">取消</a-button>
            <a-button key="submit" type="primary" :loading="loading" @click="handleCloseOk">确定</a-button>
          </template>
        </a-modal>
        <!-- hook model -->
        <DetailHookComponent ref="DetailHookComponent" :orderDetail="orderDetail"></DetailHookComponent>
      </template>
      <!-- 任务进度 -->
      <a-card :bordered="false" v-show="statistics.total != 0" title="任务进度" style="margin-top: 12px">
        <a-row :gutter="16">
          <a-col :span="3">
            <a-statistic title="任务数" :value="this.statistics.total" style="margin-right: 50px">
              <template #suffix>
                <a-icon type="flag" />
              </template>
            </a-statistic>
          </a-col>
          <a-col :span="3">
            <a-statistic title="已完成" :value="this.statistics.completed" class="demo-class">
              <template #suffix>
                <a-icon type="flag" />
              </template>
            </a-statistic>
          </a-col>
          <a-col :span="3">
            <a-statistic title="未执行" :value="this.statistics.unexecuted" class="demo-class">
              <template #suffix>
                <a-icon type="flag" />
              </template>
            </a-statistic>
          </a-col>
          <a-col :span="3">
            <a-statistic title="已失败" :value="this.statistics.failed" class="demo-class">
              <template #suffix>
                <a-icon type="flag" />
              </template>
            </a-statistic>
          </a-col>
          <a-col :span="3">
            <a-statistic title="执行中" :value="this.statistics.processing" class="demo-class">
              <template #suffix>
                <a-icon type="flag" />
              </template>
            </a-statistic>
          </a-col>
          <a-col :span="3">
            <a-statistic title="已暂停" :value="this.statistics.paused" class="demo-class">
              <template #suffix>
                <a-icon type="flag" />
              </template>
            </a-statistic>
          </a-col>
        </a-row>
      </a-card>
      <!-- 工单进度 -->
      <DetailProgressComponent ref="DetailProgressComponent" :orderDetail="orderDetail"></DetailProgressComponent>
      <!-- 工单内容 -->
      <div style="margin-top: 12px">
        <DetailCodeMirrorComponent ref="DetailCodeMirrorComponent"></DetailCodeMirrorComponent>
      </div>
    </page-header-wrapper>
  </template>
  
  <script>
  import {
    approveOrdersApi,
    closeOrdersApi,
    feedbackOrdersApi,
    generateTasksApi,
    getOrdersDetailApi,
    previewTasksApi,
    reviewOrdersApi,
  } from '@/api/orders'
  import DetailCodeMirrorComponent from './DetailCodemirror.vue'
  import DetailHookComponent from './DetailHook.vue'
  import DetailProgressComponent from './DetailProgress.vue'
  
  // 模态框按钮
  const ConfirmBtnTips = { okText: '确认', cancelText: '取消', action: '' }
  const CloseBtnTips = { okText: '提交', cancelText: '关闭', action: 'close' }
  const BtnStatus = {
    btnDisabled: false,
    closeDisabled: false,
  }
  
  export default {
    components: {
      DetailCodeMirrorComponent,
      DetailHookComponent,
      DetailProgressComponent,
    },
    data() {
      return {
        orderDetail: null,
        confirmBtnTips: ConfirmBtnTips,
        btnStatus: BtnStatus,
        loading: false,
        executeLoading: false,
        refreshLoading: false,
        actionVisible: false,
        closeVisible: false,
        confirmMsg: '',
        statistics: {},
      }
    },
    methods: {
      getDetail() {
        this.loading = true
        getOrdersDetailApi(this.$route.params.order_id)
          .then((res) => {
            if (res.code === '0000') {
              this.orderDetail = res.data
              this.$nextTick(() => {
                this.$refs.DetailCodeMirrorComponent.setValue(this.orderDetail.content)
              })
            } else {
              this.$message.error(res.message)
            }
          })
          .catch((_error) => {})
          .finally(() => {
            this.loading = false
          })
      },
      // 刷新
      refresh() {
        this.refreshLoading = true
        this.getDetail()
        this.previewTasks()
        this.$refs.DetailProgressComponent.getOpLogs()
        this.refreshLoading = false
      },
      // action modal
      showActionModal() {
        this.actionVisible = true
      },
      hideActionModal() {
        this.actionVisible = false
      },
      handleActionOk() {
        this.hideActionModal()
        this.handleCommit('ok', ConfirmBtnTips.action)
      },
      handleActionCancel() {
        this.hideActionModal()
        this.handleCommit('cancel', ConfirmBtnTips.action)
      },
      handleCommit(btn, action) {
        // console.log(btn, action)
        this.loading = true
        if (this.confirmMsg.length > 256) {
          this.$message.error('提交失败, 消息长度不能超过256个字符')
          this.hideActionModal()
          return false
        }
        // 审批
        if (action === 'approve') {
          // status: pass, reject
          var status = btn === 'ok' ? 'pass' : 'reject'
          var data = {
            status: status,
            msg: this.confirmMsg,
            order_id: this.orderDetail.order_id,
          }
          this.approveOrders(data)
        }
        // 反馈
        if (action === 'feedback') {
          var progress = btn === 'ok' ? '已完成' : '执行中'
          var data = {
            progress: progress,
            msg: this.confirmMsg,
            order_id: this.orderDetail.order_id,
          }
          this.feedbackOrders(data)
        }
        // 复核
        if (action === 'review') {
          var progress = btn === 'ok' ? '已复核' : 'None'
          // 用户点击了关闭模态框
          if (progress === 'None') {
            this.hideActionModal()
            return
          }
          // 用户点击了复核按钮
          var data = {
            msg: this.confirmMsg,
            order_id: this.orderDetail.order_id,
          }
          this.reviewOrders(data)
        }
        // 关闭
        if (action === 'close') {
          // 用户点击了关闭按钮
          var data = {
            msg: this.confirmMsg,
            order_id: this.orderDetail.order_id,
          }
          this.closeOrders(data)
        }
        this.loading = false
        this.confirmMsg = ''
      },
      // 审批
      approveOrders(data) {
        approveOrdersApi(data)
          .then((res) => {
            const messageType = res.code === '0000' ? 'info' : 'error'
            this.$message[messageType](res.message)
          })
          .catch((_error) => {})
          .finally(() => {
            this.hideActionModal()
            this.refresh()
          })
      },
      // 反馈
      feedbackOrders(data) {
        feedbackOrdersApi(data)
          .then((res) => {
            const messageType = res.code === '0000' ? 'info' : 'error'
            this.$message[messageType](res.message)
          })
          .catch((_error) => {})
          .finally(() => {
            this.hideActionModal()
            this.refresh()
          })
      },
      // 复核
      reviewOrders(data) {
        reviewOrdersApi(data)
          .then((res) => {
            const messageType = res.code === '0000' ? 'info' : 'error'
            this.$message[messageType](res.message)
          })
          .catch((_error) => {})
          .finally(() => {
            this.hideActionModal()
            this.refresh()
          })
      },
      // 关闭
      closeOrders(data) {
        closeOrdersApi(data)
          .then((res) => {
            const messageType = res.code === '0000' ? 'info' : 'error'
            this.$message[messageType](res.message)
          })
          .catch((_error) => {})
          .finally(() => {
            this.refresh()
          })
      },
      // hook modal
      showHookModal() {
        this.$refs.DetailHookComponent.showModal()
      },
      // close modal
      showCloseModal() {
        this.closeVisible = true
      },
      hideCloseModal(e) {
        this.closeVisible = false
      },
      handleCloseOk(e) {
        this.handleCommit('ok', CloseBtnTips.action)
        this.hideCloseModal()
      },
      handleCloseCancel() {
        this.hideCloseModal()
      },
      // 生成子任务
      generateSqlOrdersTasks() {
        // console.log('makeOrdersSubTasks ')
        var data = {
          order_id: this.orderDetail.order_id,
        }
        generateTasksApi(data)
          .then((res) => {
            if (res.code === '0000') {
              this.$router.push(`/orders/tasks/${this.orderDetail.order_id}`)
            } else {
              this.$message.error(res.message)
            }
          })
          .catch((_error) => {})
      },
      // 预览任务
      previewTasks() {
        var params = {
          order_id: this.$route.params.order_id,
        }
        previewTasksApi(params)
          .then((res) => {
            this.statistics = res.data
          })
          .catch((_error) => {})
      },
    },
    mounted() {
      this.getDetail()
      this.previewTasks()
    },
    filters: {
      // ('待审核', '已驳回', '已批准', '执行中', '已关闭', '已完成', '已复核')
      btnTitle(progress) {
        if (progress === '待审核') {
          BtnStatus.btnDisabled = false
          BtnStatus.closeDisabled = false
          ConfirmBtnTips.okText = '同意'
          ConfirmBtnTips.cancelText = '驳回'
          ConfirmBtnTips.action = 'approve'
          return '审核'
        } else if (['已批准', '执行中'].includes(progress)) {
          BtnStatus.btnDisabled = false
          BtnStatus.closeDisabled = false
          ConfirmBtnTips.okText = '执行完成'
          ConfirmBtnTips.cancelText = '执行中'
          ConfirmBtnTips.action = 'feedback'
          return '更新状态'
        } else if (progress === '已完成') {
          BtnStatus.btnDisabled = false
          BtnStatus.closeDisabled = true
          ConfirmBtnTips.okText = '确定'
          ConfirmBtnTips.cancelText = '取消'
          ConfirmBtnTips.action = 'review'
          return '复核'
        } else if (['已复核', '已驳回', '已关闭'].includes(progress)) {
          BtnStatus.btnDisabled = true
          BtnStatus.closeDisabled = true
          return '完成'
        }
      },
    },
  }
  </script>
  
  <style lang="less" scoped>
  ::v-deep .ant-pro-page-header-wrap-children-content {
    margin: 12px 0px 0;
  }
  </style>
  