<template>
  <a-card title="执行输出" v-show="uiState.open" class="log-card" :bodyStyle="{ padding: 0 }">
    <div class="terminal-wrapper">
      <div ref="termRef" class="xterm-container"></div>
    </div>
  </a-card>
</template>

<script setup>
import { FitAddon } from '@xterm/addon-fit'
import { Terminal } from '@xterm/xterm'
import '@xterm/xterm/css/xterm.css'
import { nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

// 状态与路由
const route = useRoute()
const orderID = route.params.order_id
const termRef = ref(null)

const uiState = reactive({
  open: false,
})

// 核心变量（非响应式，保证性能）
let term = null
let fitAddon = null
let ws = null
let flushTimer = null
let outputBuffer = ''
let isInitializing = false

// 执行流控制
let currentExecutionID = null

// Terminal 初始化
const initTerm = async () => {
  if (term || isInitializing) return
  isInitializing = true

  await nextTick()

  term = new Terminal({
    cursorBlink: true,
    convertEol: true,
    fontSize: 13,
    disableStdin: true,
    scrollback: 10000,
  })

  fitAddon = new FitAddon()
  term.loadAddon(fitAddon)
  term.open(termRef.value)

  setTimeout(() => {
    fitAddon.fit()
  }, 100)

  window.addEventListener('resize', onResize)
  isInitializing = false
}

const onResize = () => fitAddon && fitAddon.fit()

// Flush（缓冲写入，防止卡死）
const startFlush = () => {
  if (flushTimer) return

  flushTimer = setInterval(() => {
    if (!term || !outputBuffer) return

    const buffer = term.buffer.active
    const isAtBottom = (buffer.baseY - buffer.viewportY) <= term.rows

    term.write(outputBuffer)
    outputBuffer = ''

    if (isAtBottom) {
      term.scrollToBottom()
    }
  }, 30)
}

// WebSocket 逻辑
const initWebsocket = () => {
  closeWs()

  const protocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://'
  const wsURL = `${protocol}${window.location.host}/ws/${orderID}`
  ws = new WebSocket(wsURL)

  ws.onopen = () => {
    startFlush()
  }

  ws.onmessage = (msg) => {
    if (msg.data === 'pong') return

    let result
    try {
      result = JSON.parse(msg.data)
    } catch {
      outputBuffer += msg.data
      return
    }

    const { execution_id, type, data } = result

    // 自动打开面板
    if (!uiState.open) {
      uiState.open = true
    }

    // execution_id 变化：新的一次执行
    if (execution_id && execution_id !== currentExecutionID) {
      currentExecutionID = execution_id

      // 清空历史执行
      outputBuffer = ''
      if (term) term.reset()
    }

    // 根据 type 决定渲染策略
    switch (type) {
      case 'log':
      case 'gh-ost':
        outputBuffer += data
        break

      case 'processlist':
        // 每条都是快照：清空 + 重绘
        outputBuffer = ''
        if (term) term.reset()
        outputBuffer = formatProcessList(data)
        break

      default:
        // 兜底
        outputBuffer += String(data ?? '')
    }
  }

  ws.onclose = (e) => {
    if (e.code !== 1000) {
      setTimeout(initWebsocket, 3000)
    }
  }
}

const closeWs = () => {
  if (ws) {
    ws.close(1000)
    ws = null
  }
}

// 工具函数
const formatProcessList = (data) => {
  return Object.entries(data)
    .map(([k, v]) => `\x1b[32m${k}\x1b[0m: ${v}`)
    .join('\r\n') + '\r\n'
}

// 生命周期
watch(() => uiState.open, (val) => {
  if (val) initTerm()
})

onMounted(() => {
  initWebsocket()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', onResize)
  if (flushTimer) clearInterval(flushTimer)
  closeWs()
  if (term) term.dispose()
})
</script>

<style scoped>
.log-card {
  margin-top: 10px;
  border-radius: 6px;
  overflow: hidden;
}

.terminal-wrapper {
  background-color: #181818;
  padding: 12px;
}

.xterm-container {
  height: 400px;
  width: 100%;
}

:deep(.xterm-viewport::-webkit-scrollbar) {
  width: 10px;
}

:deep(.xterm-viewport::-webkit-scrollbar-track) {
  background: #181818;
}

:deep(.xterm-viewport::-webkit-scrollbar-thumb) {
  background: #444;
  border-radius: 5px;
}

:deep(.xterm-viewport::-webkit-scrollbar-thumb:hover) {
  background: #666;
}
</style>
