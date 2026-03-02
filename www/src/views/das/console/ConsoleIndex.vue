<template>
  <div class="console-page">
    <div class="console-page-head">
      <div class="head-title">SQL 编辑与执行</div>
      <div class="head-subtitle">左侧选择数据库和数据表，右侧编写 SQL 并查看执行结果。</div>
    </div>

    <div class="console-main">
      <SplitPanel leftWidth="360px">
        <template #left-content>
          <ConsoleLeft />
        </template>

        <template #right-content>
          <div class="console-workspace">
            <ConsoleRight
              @renderResultTable="renderResultTable"
              @renderExecutionMessage="renderExecutionMessage"
            />

            <a-card class="result-card" size="small">
              <template #title>查询结果</template>
              <template #extra>
                <a-space :size="16" class="result-meta">
                  <span>列 {{ uiData.tableColumns.length }}</span>
                  <span>行 {{ uiData.tableData.length }}</span>
                </a-space>
              </template>

              <a-tabs v-model:activeKey="uiData.resultTab" size="small" class="result-tabs">
                <a-tab-pane key="result" tab="结果集">
                  <div class="result-pane">
                    <div ref="resultTableRegionRef" class="result-table-region">
                      <a-empty
                        v-if="uiData.tableColumns.length === 0"
                        description="执行 SQL 后可在这里查看结果集"
                      />
                      <a-table
                        v-else
                        size="small"
                        bordered
                        :row-key="(_, index) => index"
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
                    <a-empty v-if="!uiData.executionMessage" description="暂无执行消息" />
                    <pre v-else class="exec-message">{{ uiData.executionMessage }}</pre>
                  </div>
                </a-tab-pane>
              </a-tabs>
            </a-card>
          </div>
        </template>
      </SplitPanel>
    </div>
  </div>
</template>

<script setup>
import SplitPanel from '@/components/panel/index.vue'
import { computed, nextTick, onBeforeUnmount, onMounted, provide, reactive, ref, watch } from 'vue'
import ConsoleLeft from './ConsoleLeft.vue'
import ConsoleRight from './ConsoleRight.vue'

const dasInstanceData = reactive({})
provide('dasInstanceData', dasInstanceData)

const uiData = reactive({
  tableColumns: [],
  tableData: [],
  executionMessage: '',
  resultTab: 'result',
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
  const bodyHeight = Math.max(220, Math.floor(regionHeight - theadHeight - paginationHeight - reserved))
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

const renderResultTable = (value) => {
  if (value) {
    uiData.tableColumns = value.columns || []
    uiData.tableData = value.data || []
    uiData.resultTab = 'result'
    nextTick(recomputeTableBodyHeight)
  } else {
    uiData.tableColumns = []
    uiData.tableData = []
    nextTick(recomputeTableBodyHeight)
  }
}

const renderExecutionMessage = (value) => {
  uiData.executionMessage = value || ''
  if (uiData.executionMessage) {
    uiData.resultTab = 'message'
  }
}
</script>

<style scoped>
.console-page {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.console-page-head {
  padding: 0 2px;
}

.head-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--ant-colorText, #1f1f1f);
}

.head-subtitle {
  margin-top: 2px;
  color: var(--ant-colorTextSecondary, rgba(0, 0, 0, 0.45));
  font-size: 12px;
}

.console-main :deep(.split-wrapper) {
  height: 78vh;
  max-height: 78vh;
  background: var(--ant-colorBgLayout, #f5f5f5);
}

.console-main :deep(.split-wrapper .left-content),
.console-main :deep(.split-wrapper .right-content) {
  background: var(--ant-colorBgLayout, #f5f5f5);
}

.console-main :deep(.split-wrapper .right-content) {
  padding-top: 5px;
  overflow: hidden;
}

.console-main :deep(.split-wrapper .separator),
.console-main :deep(.split-wrapper .collapsed-handle) {
  background-color: var(--ant-colorFillAlter, #fafafa);
  box-shadow: none;
}

.console-workspace {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  gap: 8px;
}

.result-card {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}

.result-card :deep(.ant-card-head) {
  min-height: 40px;
}

.result-card :deep(.ant-card-body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 8px;
  min-height: 0;
}

.result-meta {
  font-size: 12px;
  color: var(--ant-colorTextSecondary, rgba(0, 0, 0, 0.45));
}

.result-tabs {
  height: 100%;
}

.result-tabs :deep(.ant-tabs-nav) {
  margin-bottom: 8px;
}

.result-tabs :deep(.ant-tabs-content-holder),
.result-tabs :deep(.ant-tabs-content),
.result-tabs :deep(.ant-tabs-tabpane) {
  height: 100%;
}

.result-pane,
.message-pane {
  height: 100%;
  min-height: 0;
}

.result-table-region {
  height: 100%;
  min-height: 0;
  overflow: hidden;
}

.message-pane {
  overflow: auto;
  border: 1px solid var(--ant-colorSplit, #f0f0f0);
  border-radius: 6px;
  padding: 8px;
  background: #fff;
}

.exec-message {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 12px;
  line-height: 1.6;
}

.result-card :deep(.ant-table-wrapper .ant-table-pagination) {
  margin: 0;
  padding: 8px 0 0;
}

.result-card :deep(.ant-table) {
  font-size: 12px;
}

.result-card :deep(.ant-table-thead > tr > th) {
  font-weight: 600;
  background: var(--ant-colorFillAlter, #fafafa);
}
</style>
