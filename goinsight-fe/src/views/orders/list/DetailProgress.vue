<template>
  <a-card :bordered="false" style="margin-top: 12px" title="工单进度">
    <a-steps :current="currentStatus" size="small">
      <a-step title="创建工单"></a-step>
      <a-step title="待审核"></a-step>
      <a-step title="已批准/已驳回"></a-step>
      <a-step title="执行中"></a-step>
      <a-step title="已完成/已关闭"></a-step>
      <a-step title="已复核"></a-step>
      <a-step title="已勾住"></a-step>
    </a-steps>

    <a-row :gutter="8" style="margin-top: 24px">
      <a-col class="gutter-row" :span="opLogs.length != 0 ? 10 : 24">
        <a-card type="inner" ref="leftHeight">
          <a-descriptions size="small" :col="4" v-for="k of orderDetail.approver" :key="`auditor` + k.user + k.status">
            <a-descriptions-item label="审核人">
              {{ k.user }}
              <span v-if="k.is_superuser === 1">(超级审核人)</span>
            </a-descriptions-item>
            <a-descriptions-item label="状态">
              <a-tag v-if="k.status === 'pending'" color="blue">未审核</a-tag>
              <a-tag v-else-if="k.status === 'pass'" color="green">已审核</a-tag>
              <a-tag v-else color="red">已驳回</a-tag>
            </a-descriptions-item>
          </a-descriptions>
          <a-divider style="margin: 5px" />
          <a-descriptions size="small" :col="4" v-for="k of orderDetail.reviewer" :key="`reviewer` + k.user + k.status">
            <a-descriptions-item label="复核人">{{ k.user }}</a-descriptions-item>
            <a-descriptions-item label="状态">
              <a-tag v-if="k.status === 'pending'" color="blue">未复核</a-tag>
              <a-tag v-else color="green">已复核</a-tag>
            </a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>
      <a-col :span="opLogs.length != 0 ? 14 : 0">
        <a-card class="box-card" type="inner" :style="{ height: leftHeight }">
          <span v-html="opLogs"></span>
        </a-card>
      </a-col>
    </a-row>
  </a-card>
</template>


<script>
import { getOpLogsApi } from '@/api/orders'

import elementResizeDetectorMaker from 'element-resize-detector'
const erd = elementResizeDetectorMaker()

export default {
  props: {
    orderDetail: Object,
  },
  data() {
    return {
      opLogs: '',
      leftHeight: '',
    }
  },
  methods: {
    getOpLogs() {
      var params = { order_id: this.orderDetail.order_id }
      getOpLogsApi(params)
        .then((res) => {
          var msgs = []
          if (res.code === '0000') {
            res.data.forEach(function (val) {
              msgs.push(`[${val.updated_at}] ${val.msg}`)
            })
            this.opLogs = msgs.join('<br>')
          }
        })
        .catch((_error) => {})
    },
    // 实时获取高度，并改变右侧区域的高度
    getLeftHeight() {
      erd.listenTo(this.$refs.leftHeight.$el, (element) => {
        this.leftHeight = element.offsetHeight + 'px'
      })
    },
  },
  computed: {
    currentStatus() {
      if (this.orderDetail.progress === '待审核') {
        return 1
      } else if (this.orderDetail.progress === '已驳回') {
        return 2
      } else if (this.orderDetail.progress === '已批准') {
        return 2
      } else if (this.orderDetail.progress === '执行中') {
        return 3
      } else if (this.orderDetail.progress === '已关闭') {
        return 4
      } else if (this.orderDetail.progress === '已完成') {
        return 4
      } else if (this.orderDetail.progress === '已复核') {
        return 5
      } else {
        return 6
      }
    },
  },
  mounted() {
    this.getOpLogs()
    this.getLeftHeight()
  },
}
</script>

<style lang='less' scoped>
.box-card {
  height: 100%;
  overflow: auto;
  zoom: 1;
  white-space: normal;
  word-break: break-all;
}
</style>