<template>
  <a-row>
    <a-table
      :columns="tableColumns"
      :dataSource="tableData"
      :pagination="pagination"
      :loading="tableLoading"
      :rowClassName="setRowClass"
      :rowKey="(record, index) => index"
      @change="handleTableChange"
      size="middle"
    >
      <template slot="summary" slot-scope="text">
        <ul
          style="list-style-type: square; padding: 0 0 0 16px; margin: 2px"
          v-for="(item, index) of text"
          :key="index"
        >
          <li>{{ item }}</li>
        </ul>
      </template>
    </a-table>
  </a-row>
</template>

<script>
const tableColumns = [
  {
    title: '指纹',
    dataIndex: 'finger_id',
    key: 'finger_id',
    width: '15%',
    scopedSlots: {
      customRender: 'finger_id',
    },
  },
  {
    title: '级别',
    dataIndex: 'level',
    key: 'level',
    ellipsis: true,
    scopedSlots: {
      customRender: 'level',
    },
  },
  {
    title: '查询',
    dataIndex: 'query',
    key: 'query',
    width: '25%',
    ellipsis: true,
    scopedSlots: {
      customRender: 'query',
    },
  },
  {
    title: '提示',
    dataIndex: 'summary',
    key: 'summary',
    width: '50%',
    scopedSlots: {
      customRender: 'summary',
    },
  },
  {
    title: '类型',
    dataIndex: 'type',
    key: 'type',
    scopedSlots: {
      customRender: 'type',
    },
  },
]
export default {
  data() {
    return {
      // table
      tableLoading: false,
      tableColumns,
      tableData: [],
      // pagination
      pagination: {
        current: 1,
        pageSize: 10,
        total: 0,
        pageSizeOptions: ['10', '20', '50'],
        showSizeChanger: true,
      },
    }
  },
  methods: {
    // 设置行的颜色
    setRowClass(record) {
      if (record.level === 'INFO') {
        return 'row-level-info'
      }
      if (record.level === 'WARN') {
        return 'row-level-warn'
      }
      if (record.level === 'ERROR') {
        return 'row-level-error'
      }
    },
    renderData(res) {
      if (res.status === 0) {
        this.$notification.success({
          message: '成功',
          description: '语法检查通过，您可以提交工单了，O(∩_∩)O',
        })
      }
      if (res.status === 1) {
        this.$notification.warning({
          message: '警告',
          description: '语法检查未通过，请根据下面输出提示进行更正，(ㄒoㄒ)',
        })
      }
      this.tableData = res.data
    },
    // 操作表格
    handleTableChange(pager) {
      this.pagination.current = pager.current
      this.pagination.pageSize = pager.pageSize
    },
  },
}
</script>

<style>
.row-level-info {
  color: green;
}
.row-level-warn {
  color: orange;
}
.row-level-error {
  color: red;
}
</style>
