<template>
  <a-card :title="cardTitle">
    <a-table
      :columns="table.columns"
      :rowKey="(record) => record.id"
      :dataSource="table.data"
      :pagination="pagination"
      :loading="loading"
      @change="handleTableChange"
      size="middle"
    >
      <span slot="escape_title" slot-scope="text, record">
        <router-link :to="{ name: 'view.sqlorders.detail', params: { order_id: record.order_id } }">{{
          text
        }}</router-link>
      </span>
    </a-table>
  </a-card>
</template>

<script>
import { viewReleaseVersions } from '@/api/sql'
import { orderProgress } from '@/utils/sql'

// 转换
const progressList = {}
orderProgress.forEach((item) => {
  progressList[item.key] = item.value
})
const filterCol = ['id', 'escape_title', 'order_id', 'applicant']

export default {
  data() {
    return {
      cardTitle: this.$route.params.version,
      loading: false,
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
    }
  },
  methods: {
    handleTableChange(pager) {
      this.pagination.current = pager.current
      this.pagination.pageSize = pager.pageSize
    },
    progressFormat(data) {
      // 格式化进度
      Object.keys(data).forEach((key) => {
        if (!filterCol.includes(key)) {
          if (data[key] === -1) {
            data[key] = '-'
          } else {
            data[key] = progressList[data[key].toString()]
          }
        }
      })
      return data
    },
    runViewReleaseVersions() {
      this.loading = true
      viewReleaseVersions(this.$route.params.version).then((response) => {
        this.table.columns = response.data.columns
        this.table.data = response.data.data

        this.table.data.map((item) => {
          item = this.progressFormat(item)
        })
      })
      this.loading = false
    },
  },
  created() {
    this.runViewReleaseVersions()
  },
}
</script>

<style>
</style>