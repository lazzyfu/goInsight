<template>
  <div v-if="data.showbTable" style="margin-top: 15px; overflow-x: auto">
    <a-table size="small" class="ant-table-striped" bordered :data-source="data.tableData" :columns="tableColumns"
      :rowClassName="setRowClass" :scroll="{ x: 1100 }">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'summary'">
          <ul style="list-style-type: square; padding: 0 0 0 16px; margin: 2px" v-for="(item, index) of record.summary"
            :key="index">
            <li>{{ item }}</li>
          </ul>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup>
import { message } from 'ant-design-vue'
import { reactive } from 'vue'

const tableColumns = [
  {
    title: '级别',
    dataIndex: 'level',
    key: 'level',
    ellipsis: true,
  },
  {
    title: '指纹',
    dataIndex: 'finger_id',
    key: 'finger_id',
    width: '15%',
  },
  {
    title: '语句',
    dataIndex: 'query',
    key: 'query',
    width: '30%',
    ellipsis: true,
  },
  {
    title: '提示',
    dataIndex: 'summary',
    key: 'summary',
    width: '40%',
  },
  {
    title: '类型',
    dataIndex: 'type',
    key: 'type',
  },
]

const data = reactive({
  showbTable: false,
  tableData: [],
})

const setRowClass = (record) => {
  if (record.level === 'INFO') {
    return 'row-level-info'
  }
  if (record.level === 'WARN') {
    return 'row-level-warn'
  }
  if (record.level === 'ERROR') {
    return 'row-level-error'
  }
}

const render = (res) => {
  if (res.status === 0) {
    message.info('语法检查通过，您可以提交工单了，O(∩_∩)O')
  }
  if (res.status === 1) {
    message.error('语法检查未通过，请根据下面输出提示进行更正，(ㄒoㄒ)')
  }
  data.tableData = res.data
  data.showbTable = true
}

defineExpose({
  render,
})
</script>

<style scoped>
:deep(.row-level-info) {
  color: green;
}

:deep(.row-level-warn) {
  color: orange;
}

:deep(.row-level-error) {
  color: red;
}
</style>
