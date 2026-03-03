<template>
  <div class="console-page">
    <div class="console-main">
      <SplitPanel leftWidth="420px">
        <template #left-content>
          <ConsoleLeft />
        </template>

        <template #right-content>
          <div ref="consoleWorkspaceRef" class="console-workspace">
            <ConsoleRight
              :editor-height="uiData.editorHeight"
              @renderResultTable="renderResultTable"
              @renderExecutionMessage="renderExecutionMessage"
            />

            <div
              class="workspace-resizer"
              role="separator"
              aria-label="调整编辑区与结果区高度"
              tabindex="0"
              @mousedown="startWorkspaceResize"
              @keydown.up.prevent="adjustEditorHeight(-20)"
              @keydown.down.prevent="adjustEditorHeight(20)"
            ></div>

            <a-card ref="resultCardRef" class="result-card" size="small">
              <template #title>查询结果</template>
              <template #extra>
                <a-space :size="14" class="result-meta" wrap>
                  <span>列 {{ uiData.tableColumns.length }}</span>
                  <span>行 {{ uiData.tableData.length }}</span>
                  <a-button type="link" size="small" class="fullscreen-trigger" @click="toggleResultFullscreen">
                    {{ uiData.resultFullscreen ? '退出全屏' : '全屏显示' }}
                  </a-button>
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
                        row-key="__goinsightRowKey"
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
import { message } from 'ant-design-vue'
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
  resultFullscreen: false,
  editorHeight: 280,
})

const MIN_EDITOR_HEIGHT = 200
const MIN_RESULT_HEIGHT = 280

const consoleWorkspaceRef = ref(null)
const resultCardRef = ref(null)
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
  nextTick(() => {
    uiData.editorHeight = clampEditorHeight(uiData.editorHeight)
    recomputeTableBodyHeight()
  })
}

const getMaxEditorHeight = () => {
  const workspaceEl = consoleWorkspaceRef.value
  if (!workspaceEl) return 460
  const workspaceHeight = workspaceEl.getBoundingClientRect().height
  const reservedHeight = 80
  const maxHeight = Math.floor(workspaceHeight - MIN_RESULT_HEIGHT - reservedHeight)
  return Math.max(MIN_EDITOR_HEIGHT, maxHeight)
}

const clampEditorHeight = (value) => {
  const normalized = Number.isFinite(value) ? value : MIN_EDITOR_HEIGHT
  return Math.min(getMaxEditorHeight(), Math.max(MIN_EDITOR_HEIGHT, Math.floor(normalized)))
}

const adjustEditorHeight = (delta) => {
  uiData.editorHeight = clampEditorHeight(uiData.editorHeight + delta)
  nextTick(recomputeTableBodyHeight)
}

let resizeStartY = 0
let resizeStartHeight = 0
let isResizingWorkspace = false

const handleWorkspaceResizeMove = (event) => {
  if (!isResizingWorkspace) return
  const deltaY = event.clientY - resizeStartY
  uiData.editorHeight = clampEditorHeight(resizeStartHeight + deltaY)
  nextTick(recomputeTableBodyHeight)
}

const stopWorkspaceResize = () => {
  if (!isResizingWorkspace) return
  isResizingWorkspace = false
  document.documentElement.style.userSelect = ''
  document.removeEventListener('mousemove', handleWorkspaceResizeMove)
  document.removeEventListener('mouseup', stopWorkspaceResize)
}

const startWorkspaceResize = (event) => {
  if (typeof document === 'undefined') return
  isResizingWorkspace = true
  resizeStartY = event.clientY
  resizeStartHeight = uiData.editorHeight
  document.documentElement.style.userSelect = 'none'
  document.addEventListener('mousemove', handleWorkspaceResizeMove)
  document.addEventListener('mouseup', stopWorkspaceResize)
}

const getResultCardElement = () => {
  if (!resultCardRef.value) return null
  return resultCardRef.value.$el || resultCardRef.value
}

const getFullscreenElement = () => {
  if (typeof document === 'undefined') return null
  return document.fullscreenElement || document.webkitFullscreenElement || null
}

const exitBrowserFullscreen = async () => {
  if (typeof document === 'undefined') return
  if (typeof document.exitFullscreen === 'function') {
    await document.exitFullscreen()
    return
  }
  if (typeof document.webkitExitFullscreen === 'function') {
    document.webkitExitFullscreen()
  }
}

const requestElementFullscreen = async (el) => {
  if (typeof el.requestFullscreen === 'function') {
    await el.requestFullscreen()
    return
  }
  if (typeof el.webkitRequestFullscreen === 'function') {
    el.webkitRequestFullscreen()
    return
  }
  throw new Error('fullscreen-not-supported')
}

const syncFullscreenState = () => {
  const resultCardEl = getResultCardElement()
  uiData.resultFullscreen = !!(resultCardEl && getFullscreenElement() === resultCardEl)
  nextTick(recomputeTableBodyHeight)
}

const toggleResultFullscreen = async () => {
  const resultCardEl = getResultCardElement()
  if (!resultCardEl) return

  try {
    const fullscreenEl = getFullscreenElement()
    if (fullscreenEl === resultCardEl) {
      await exitBrowserFullscreen()
      return
    }
    if (fullscreenEl) {
      await exitBrowserFullscreen()
    }
    await requestElementFullscreen(resultCardEl)
  } catch (error) {
    message.warning('当前浏览器不支持结果集全屏显示')
  }
}

onMounted(() => {
  if (typeof window !== 'undefined') {
    window.addEventListener('resize', handleWindowResize)
  }
  if (typeof document !== 'undefined') {
    document.addEventListener('fullscreenchange', syncFullscreenState)
    document.addEventListener('webkitfullscreenchange', syncFullscreenState)
  }

  if (typeof ResizeObserver !== 'undefined') {
    tableRegionResizeObserver = new ResizeObserver(() => {
      nextTick(recomputeTableBodyHeight)
    })
  }

  nextTick(() => {
    uiData.editorHeight = clampEditorHeight(uiData.editorHeight)
    recomputeTableBodyHeight()
  })
})

onBeforeUnmount(() => {
  tableRegionResizeObserver?.disconnect?.()
  if (typeof window !== 'undefined') {
    window.removeEventListener('resize', handleWindowResize)
  }
  if (typeof document !== 'undefined') {
    stopWorkspaceResize()
  }
  if (typeof document !== 'undefined') {
    document.removeEventListener('fullscreenchange', syncFullscreenState)
    document.removeEventListener('webkitfullscreenchange', syncFullscreenState)
    if (getFullscreenElement() === getResultCardElement()) {
      exitBrowserFullscreen().catch(() => {})
    }
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

watch(
  () => consoleWorkspaceRef.value,
  async (el) => {
    if (!el) return
    await nextTick()
    uiData.editorHeight = clampEditorHeight(uiData.editorHeight)
    recomputeTableBodyHeight()
  },
  { immediate: true },
)

const renderResultTable = (value) => {
  if (value) {
    uiData.tableColumns = value.columns || []
    uiData.tableData = (value.data || []).map((row, index) => ({
      ...(row || {}),
      __goinsightRowKey: index + 1,
    }))
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
  --console-radius: 10px;
  --console-border: 1px solid var(--ant-colorBorderSecondary, #f0f0f0);
  --console-muted: var(--ant-colorTextSecondary, rgba(0, 0, 0, 0.45));
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.console-main :deep(.split-wrapper) {
  height: clamp(680px, 82vh, 900px);
  max-height: clamp(680px, 82vh, 900px);
  border-radius: 12px;
  border: var(--console-border);
  background: var(--ant-colorBgLayout, #f5f5f5);
  padding: 8px;
  gap: 8px;
}

.console-main :deep(.split-wrapper .left-content),
.console-main :deep(.split-wrapper .right-content) {
  background: transparent;
}

.console-main :deep(.split-wrapper .right-content) {
  padding: 0;
  overflow: hidden;
}

.console-main :deep(.split-wrapper .separator),
.console-main :deep(.split-wrapper .collapsed-handle) {
  background-color: var(--ant-colorBgContainer, #ffffff);
  border-left: 1px solid var(--ant-colorBorderSecondary, #f0f0f0);
  box-shadow: none;
}

.console-workspace {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  gap: 8px;
}

.workspace-resizer {
  height: 10px;
  border-radius: 999px;
  border: 1px solid var(--ant-colorBorderSecondary, #f0f0f0);
  background: var(--ant-colorBgContainer, #ffffff);
  cursor: row-resize;
  flex-shrink: 0;
  position: relative;
}

.workspace-resizer::before {
  content: '';
  position: absolute;
  left: 50%;
  top: 50%;
  width: 48px;
  height: 2px;
  transform: translate(-50%, -50%);
  background: var(--ant-colorTextQuaternary, rgba(0, 0, 0, 0.25));
  box-shadow: 0 -4px 0 var(--ant-colorTextQuaternary, rgba(0, 0, 0, 0.25)),
    0 4px 0 var(--ant-colorTextQuaternary, rgba(0, 0, 0, 0.25));
}

.workspace-resizer:hover,
.workspace-resizer:focus-visible {
  background: var(--ant-colorFillAlter, #fafafa);
  outline: none;
}

.result-card {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 280px;
  border-radius: var(--console-radius);
  border: var(--console-border);
  box-shadow: 0 6px 16px rgb(15 23 42 / 5%);
}

.result-card :deep(.ant-card-head) {
  min-height: 42px;
  border-bottom-color: var(--ant-colorBorderSecondary, #f0f0f0);
  background: var(--ant-colorFillAlter, #fafafa);
}

.result-card :deep(.ant-card-body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 10px;
  min-height: 0;
}

.result-meta {
  font-size: 12px;
  color: var(--console-muted);
}

.fullscreen-trigger {
  padding-inline: 0;
  height: auto;
}

.result-tabs {
  height: 100%;
}

.result-tabs :deep(.ant-tabs-nav) {
  margin-bottom: 10px;
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
  border: var(--console-border);
  border-radius: 8px;
  padding: 10px;
  background: var(--ant-colorBgContainer, #ffffff);
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

.result-card:fullscreen,
.result-card:-webkit-full-screen {
  width: 100%;
  height: 100%;
  border: 0;
  border-radius: 0;
  box-shadow: none;
  margin: 0;
  background: var(--ant-colorBgContainer, #ffffff);
}

.result-card:fullscreen :deep(.ant-card-head),
.result-card:-webkit-full-screen :deep(.ant-card-head) {
  border-bottom-color: var(--ant-colorBorderSecondary, #f0f0f0);
  padding-inline: 14px;
}

.result-card:fullscreen :deep(.ant-card-body),
.result-card:-webkit-full-screen :deep(.ant-card-body) {
  padding: 12px;
}

@media (max-width: 1200px) {
  .workspace-resizer {
    display: none;
  }

  .console-main :deep(.split-wrapper) {
    height: auto;
    max-height: none;
    min-height: 0;
    padding: 8px;
    display: block;
  }

  .console-main :deep(.split-wrapper .scalable) {
    width: 100% !important;
    max-width: 100%;
  }

  .console-main :deep(.split-wrapper .separator),
  .console-main :deep(.split-wrapper .collapsed-handle) {
    display: none !important;
  }

  .console-main :deep(.split-wrapper .right-content) {
    padding-top: 10px;
    min-height: 500px;
  }
}

@media (max-width: 768px) {
  .console-main :deep(.split-wrapper) {
    padding: 6px;
    gap: 6px;
  }

  .result-card :deep(.ant-card-head) {
    padding-inline: 10px;
  }

  .result-card :deep(.ant-card-body) {
    padding: 8px;
  }
}
</style>
