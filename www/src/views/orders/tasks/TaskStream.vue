<template>
  <a-card title="执行输出" v-show="uiState.open" class="mt-2">
    <CodeMirror ref="cmRef" />
  </a-card>
</template>

<script setup>
import CodeMirror from '@/components/edit/Codemirror.vue'
import { message } from 'ant-design-vue'
import { onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'

// 状态
const uiState = reactive({
  open: false,
})

const cmRef = ref(null)
const route = useRoute()
const orderID = route.params.order_id

// WebSocket URL
let protocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://'
const wsURL = `${protocol}${window.location.host}/ws/${orderID}`

// WebSocket 实例
const ws = ref(null)

// 心跳
let heartbeatTimer = null
let heartbeatInterval = 15000 // 15s 发一次 ping

// 重连控制
let reconnectTimer = null
let reconnecting = false
let reconnectInterval = 3000 // 3 秒重试一次

// 初始化 WS
const initWebsocket = () => {
  if (typeof WebSocket === 'undefined') {
    message.error('您的浏览器不支持 WebSocket')
    return
  }

  clearReconnectState()

  ws.value = new WebSocket(wsURL)
  ws.value.onopen = onOpen
  ws.value.onerror = onError
  ws.value.onclose = onClose
  ws.value.onmessage = onMessage
}

// WebSocket 打开
const onOpen = () => {
  startHeartbeat()
}

// WebSocket 错误
const onError = () => {
  tryReconnect()
}

// WebSocket 关闭
const onClose = () => {
  tryReconnect()
}

// WebSocket 接收消息
const onMessage = (msg) => {
  // 后端可能发心跳 pong
  if (msg.data === 'pong') return

  const result = JSON.parse(msg.data)
  uiState.open = true

  if (result.type === 'processlist') {
    cmRef.value.setContent(renderProcesslist(result.data))
  } else if (result.type === 'ghost') {
    cmRef.value.appendContent(result.data)
  } else {
    cmRef.value.setContent(result.data)
  }
}

// 渲染 processlist
const renderProcesslist = (data) => {
  return Object.keys(data)
    .map((key) => `${key}: ${data[key]}`)
    .join('\n')
}

// 心跳机制
const startHeartbeat = () => {
  stopHeartbeat()
  heartbeatTimer = setInterval(() => {
    if (ws.value && ws.value.readyState === WebSocket.OPEN) {
      ws.value.send('ping') // 发送心跳
    }
  }, heartbeatInterval)
}

const stopHeartbeat = () => {
  if (heartbeatTimer) {
    clearInterval(heartbeatTimer)
    heartbeatTimer = null
  }
}

// 自动重连机制
const tryReconnect = () => {
  if (reconnecting) return

  reconnecting = true
  reconnectTimer = setTimeout(() => {
    initWebsocket()
    reconnecting = false
  }, reconnectInterval)
}

const clearReconnectState = () => {
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
  reconnecting = false
}

// 关闭 WS
const closeWS = () => {
  stopHeartbeat()
  clearReconnectState()
  if (ws.value) {
    ws.value.close()
    ws.value = null
  }
}

// 生命周期
onMounted(() => {
  initWebsocket()
})

onBeforeUnmount(() => {
  closeWS()
})
</script>

<style scoped>
.mt-2 {
  margin-top: 8px;
}
</style>
