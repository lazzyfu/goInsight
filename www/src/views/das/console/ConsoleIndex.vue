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

          <div class="console-result">
            <a-tabs default-active-key="result">
              <a-tab-pane key="result" tab="结果集">
                <div class="result-pane">
                  <div ref="resultTableRegionRef" class="result-table-region">
                    <a-table
                      size="small"
                      bordered
                      :data-source="uiData.tableData"
                      :scroll="tableScroll"
                      style="min-width: 100%"
                    >
                      <a-table-column
                        v-for="item in uiData.tableColumns"
                        :key="item"
                        :title="item"
                        :data-index="item"
                      />
                    </a-table>
                  </div>
                </div>
              </a-tab-pane>
              <a-tab-pane key="message" tab="执行消息">
                <div class="message-pane">
                  <div class="exec-message" v-html="uiData.executionMessage" />
                </div>
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
import { computed, nextTick, onBeforeUnmount, onMounted, provide, reactive, ref, watch } from 'vue'
import ConsoleLeft from './ConsoleLeft.vue'
import ConsoleRight from './ConsoleRight.vue'

// 父子组件或深层嵌套组件间的数据共享
const dasInstanceData = reactive({})
provide('dasInstanceData', dasInstanceData)

// 数据
const uiData = reactive({
  tableColumns: [],
  tableData: [],
  executionMessage: '',
})

const resultTableRegionRef = ref(null)
const tableBodyScrollY = ref(null)

const tableScroll = computed(() => ({
  x: 1100,
  y: tableBodyScrollY.value ?? 360,
}))

const recomputeTableBodyHeight = () => {
  const regionEl = resultTableRegionRef.value
  if (!regionEl) return

  const regionHeight = regionEl.clientHeight
  if (!regionHeight) return

  const theadEl = regionEl.querySelector('.ant-table-thead')
  const paginationEl = regionEl.querySelector('.ant-table-pagination')

  const theadHeight = theadEl ? theadEl.getBoundingClientRect().height : 0
  const paginationHeight = paginationEl ? paginationEl.getBoundingClientRect().height : 0

  const reserved = 8
  const bodyHeight = Math.max(240, Math.floor(regionHeight - theadHeight - paginationHeight - reserved))
  tableBodyScrollY.value = bodyHeight
}

let tableRegionResizeObserver
const handleWindowResize = () => {
  nextTick(recomputeTableBodyHeight)
}

onMounted(() => {
  if (typeof window !== 'undefined') {
    window.addEventListener('resize', handleWindowResize)
  }

  if (typeof ResizeObserver !== 'undefined') {
    tableRegionResizeObserver = new ResizeObserver(() => {
      nextTick(recomputeTableBodyHeight)
    })
  }
})

onBeforeUnmount(() => {
  tableRegionResizeObserver?.disconnect?.()
  if (typeof window !== 'undefined') {
    window.removeEventListener('resize', handleWindowResize)
  }
})

watch(
  () => resultTableRegionRef.value,
  async (el) => {
    tableRegionResizeObserver?.disconnect?.()
    if (el) {
      tableRegionResizeObserver?.observe?.(el)
      await nextTick()
      recomputeTableBodyHeight()
    }
  },
  { immediate: true },
)

watch(
  () => [uiData.tableColumns.length, uiData.tableData.length],
  async () => {
    await nextTick()
    recomputeTableBodyHeight()
  },
  { immediate: true },
)

// 渲染结果表格
const renderResultTable = (value) => {
  if (value) {
    uiData.tableColumns = value.columns
    uiData.tableData = value.data
    nextTick(recomputeTableBodyHeight)
  } else {
    uiData.tableColumns = []
    uiData.tableData = []
    nextTick(recomputeTableBodyHeight)
  }
}

const renderExecutionMessage = (value) => {
  uiData.executionMessage = value || ''
}
</script>

<style scoped>
:deep(.ant-card-body) {
  padding: 10px;
}

:deep(.ant-tabs-nav) {
  margin: 0 0 10px 0;
}

/* 不分页：限制左右面板最大高度，超出滚动 */
:deep(.split-wrapper) {
  height: 82vh;
  max-height: 82vh;
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
  padding: 8px;
  box-sizing: border-box;
  border: 1px solid var(--ant-colorSplit, #f0f0f0);
  border-radius: var(--ant-borderRadiusLG, 8px);
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
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.console-result :deep(.ant-tabs-content) {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.console-result :deep(.ant-tabs-tabpane) {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.result-pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.result-table-region {
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

.message-pane {
  flex: 1;
  min-height: 0;
  overflow: auto;
}

/* 预留底部空间：避免分页默认 margin 导致被裁切 */
.console-result :deep(.ant-table-wrapper .ant-table-pagination) {
  margin: 0;
  padding: 8px 0 0;
}

/* 结果集表格字体更小（仅影响 Console 结果区） */
.console-result :deep(.ant-table) {
  font-size: 12px;
}

.console-result :deep(.ant-table-thead > tr > th),
.console-result :deep(.ant-table-tbody > tr > td) {
  font-size: 12px;
}

/* 表头：加粗 + 浅色背景 */
.console-result :deep(.ant-table-thead > tr > th) {
  font-weight: 600;
  background: var(--ant-colorFillAlter, #fafafa);
}

.exec-message {
  white-space: pre-wrap;
  font-size: 12px;
}
</style>
