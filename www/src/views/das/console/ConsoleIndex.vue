<template>
  <a-card>
    <SplitPanel leftWidth="380px">
      <template #left-content>
        <ConsoleLeft />
      </template>
      <template #right-content>
        <ConsoleRight @renderResultTable="renderResultTable" />
      </template>
    </SplitPanel>
  </a-card>

  <div v-if="uiState.showbTable" style="margin-top: 15px; overflow-x: auto">
    <a-tabs default-active-key="1">
      <a-tab-pane key="1" tab="结果集">
        <a-table
          size="small"
          class="ant-table-striped"
          bordered
          :data-source="uiData.tableData"
          :scroll="{ x: 1100 }"
          style="min-width: 100%"
        >
          <a-table-column
            v-for="item in uiData.tableColumns"
            :key="item"
            :title="item"
            :data-index="item"
          />
        </a-table>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script setup>
import SplitPanel from '@/components/panel/index.vue'
import { provide, reactive } from 'vue'
import ConsoleLeft from './ConsoleLeft.vue'
import ConsoleRight from './ConsoleRight.vue'

// 父子组件或深层嵌套组件间的数据共享
const dasInstanceData = reactive({})
provide('dasInstanceData', dasInstanceData)

// 状态
const uiState = reactive({
  showbTable: false,
})

// 数据
const uiData = reactive({
  tableColumns: [],
  tableData: [],
})

// 渲染结果表格
const renderResultTable = (value) => {
  if (value) {
    uiState.showbTable = true
    uiData.tableColumns = value.columns
    uiData.tableData = value.data
  } else {
    uiState.showbTable = false
  }
}
</script>

<style scoped>
:deep(.ant-table-tbody tr:nth-child(2n)) {
  background-color: #fafafa;
}

:deep(.ant-card-body) {
  padding: 10px;
}

:deep(.ant-tabs-nav) {
  margin: 0 0 10px 0;
}
</style>
