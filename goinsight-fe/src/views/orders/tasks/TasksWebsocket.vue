<template>
  <a-card :title="bindTitle" v-show="executeMsgVisible" style="margin-top: 8px">
    <codemirror ref="myCm" v-model="code" :options="cmOptions" @ready="onCmReady"></codemirror>
  </a-card>
</template>


<script>
import 'codemirror/mode/sql/sql.js'
import 'codemirror/addon/display/autorefresh'

// websocket
let protocol = 'ws://'
if (window.location.protocol === 'https:') {
  protocol = 'wss://'
}

export default {
  data() {
    return {
      bindTitle: '',
      websocket: {
        path: `${protocol}/${window.location.host}/ws/${this.$route.params.order_id}`,
        socket: '',
      },
      executeMsgVisible: false,
      code: '',
      cmOptions: {
        indentUnit: 2,
        tabSize: 2,
        indentWithTabs: true,
        smartIndent: true,
        autoRefresh: true,
        lineWrapping: true,
        readOnly: true,
        focuse: false,
      },
    }
  },
  methods: {
    onCmReady(cm) {
      cm.setSize('height', `450px`)
    },
    // 初始化websocket
    initWebsocket() {
      if (typeof WebSocket === 'undefined') {
        this.$message.error('您的浏览器不支持websocket')
      }
      // 实例化websocket
      this.websocket.socket = new WebSocket(this.websocket.path)
      // 监控websocket连接
      this.websocket.socket.onopen = this.socketOnOpen
      // 监听socket错误信息
      this.websocket.socket.onerror = this.socketOnError
      // 监听socket消息
      this.websocket.socket.onmessage = this.socketOnMessage
    },
    socketOnOpen() {
      //
    },
    socketOnError() {
      setTimeout(() => {
        this.initWebsocket()
      }, 3000)
    },
    socketOnMessage(msg) {
      // 接收socket信息
      var result = JSON.parse(msg.data)
      this.executeMsgVisible = true
      if (result.type === 'processlist') {
        this.codemirror.setValue(this.renderProcesslist(result.data))
      } else if (result.type === 'ghost') {
        // 追加显示
        let pos = this.codemirror.getCursor(this.codemirror.lastLine())
        this.codemirror.replaceRange(result.data, pos)
        // 自动滚动到行的末尾
        this.$nextTick(() => {
          this.codemirror.setCursor(this.codemirror.lastLine())
        })
      } else  {
        this.codemirror.setValue(result.data)
      }
    },
    socketClose() {
      // 关闭socket
      this.websocket.socket.close()
    },
    renderProcesslist(data) {
      this.bindTitle = '当前SQL SESSION ID的SHOW PROCESSLIST输出'
      let html = ''
      for (let key in data) {
        html += key + ': ' + data[key] + '\n'
      }
      return html
    },
  },
  destroyed() {
    // 关闭websocket
    this.socketClose()
  },
  computed: {
    codemirror() {
      return this.$refs.myCm.codemirror
    },
  },
  mounted() {
    this.initWebsocket()
  },
}
</script>