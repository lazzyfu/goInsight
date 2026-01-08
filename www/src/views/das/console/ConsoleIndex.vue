<template>
  <a-card>
    <SplitPanel leftWidth="380px">
      <template #left-content>
        <ConsoleLeft />
      </template>
      <template #right-content>
        <div class="console-right">
          <ConsoleRight
            @renderResultTable="renderResultTable"
            @renderExecutionMessage="renderExecutionMessage"
          />

          <div v-if="uiState.showbTable || uiData.executionMessage" class="console-result">
            <a-tabs default-active-key="result">
              <a-tab-pane key="result" tab="结果集">
                <a-table
                  v-if="uiState.showbTable"
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
              <a-tab-pane key="message" tab="执行消息">
                <div class="exec-message" v-html="uiData.executionMessage" />
              </a-tab-pane>
            </a-tabs>
          </div>
        </div>
      </template>
    </SplitPanel>
  </a-card>
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
  executionMessage: '',
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

const renderExecutionMessage = (value) => {
  uiData.executionMessage = value || ''
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

/* 不分页：限制左右面板最大高度，超出滚动 */
:deep(.split-wrapper) {
  height: 80vh;
  max-height: 80vh;
}

/* 右侧不整体滚动：让结果集区域单独滚动 */
:deep(.split-wrapper .right-content) {
  overflow: hidden;
}

.console-right {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.console-result {
  margin-top: 15px;
  flex: 1;
  min-height: 0;
  overflow: hidden;
  min-height: 240px;
}

/* 只滚动内容区：Tab 标题栏固定 */
.console-result :deep(.ant-tabs) {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.console-result :deep(.ant-tabs-content-holder) {
  flex: 1;
  min-height: 0;
  overflow: auto;
}

/* 结果集表格字体更小（仅影响 Console 结果区） */
.console-result :deep(.ant-table) {
  font-size: 12px;
}

.console-result :deep(.ant-table-thead > tr > th),
.console-result :deep(.ant-table-tbody > tr > td) {
  font-size: 12px;
}

.exec-message {
  white-space: pre-wrap;
  font-size: 12px;
}
</style>
