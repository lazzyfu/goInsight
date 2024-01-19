<template>
  <codemirror ref="myCm" v-model="code" :options="cmOptions" @ready="onCmReady"></codemirror>
</template>

<script>
// mode
import 'codemirror/mode/sql/sql.js'
import 'codemirror/mode/javascript/javascript.js'
// addon
import 'codemirror/addon/selection/active-line'
import 'codemirror/addon/display/autorefresh'
// 提示和自动补全
import 'codemirror/addon/hint/show-hint'
import 'codemirror/addon/hint/show-hint.css'
import 'codemirror/addon/hint/anyword-hint'
import 'codemirror/addon/hint/sql-hint'
import 'codemirror/addon/comment/comment'
import 'codemirror/addon/edit/matchbrackets'
import 'codemirror/addon/edit/closebrackets'
// 编辑器类型
import 'codemirror/keymap/sublime'
// 主题
import 'codemirror/theme/dracula.css'
import 'codemirror/theme/solarized.css'
import 'codemirror/theme/eclipse.css'
import 'codemirror/theme/monokai.css'
import 'codemirror/theme/oceanic-next.css'

export default {
  computed: {
    codemirror() {
      return this.$refs.myCm.codemirror
    },
  },
  data() {
    return {
      code: '',
      cmHeight: 380,
      cmOptions: {
        mode: 'text/x-mysql',
        indentUnit: 2,
        tabSize: 2,
        indentWithTabs: true,
        smartIndent: true,
        autoRefresh: true,
        lineNumbers: true,
        styleActiveLine: true,
        autoCloseBrackets: true,
        matchBrackets: true,
        lineWrapping: true, // 自动换行
        autofocus: true,
        resetSelectionOnContextMenu: false,
        showCursorWhenSelecting: true,
        keyMap: 'sublime', // 编辑器模式
      },
    }
  },
  methods: {
    // 自动补全
    onCmReady(cm) {
      // cmd+enter实现执行SQL
      const _this = this
      // 执行频率限制
      let lastExecTime = null
      cm.setOption('extraKeys', {
        // mac
        'Cmd-Enter': function (cm) {
          // 获取当前时间
          const now = new Date()
          // 比较当前时间与上一次执行时间
          if (lastExecTime && now - lastExecTime < 1000) {
            // 重复点击，提示用户
            _this.$message.warning('请勿重复点击，每秒仅允许执行一次')
            return
          }
          // 记录上一次执行时间
          lastExecTime = now
          // 执行操作
          _this.$nextTick(() => {
            _this.$emit('execEvent')
          })
        },
        // windows
        'Ctrl-Enter': function (cm) {
          // 获取当前时间
          const now = new Date()
          // 比较当前时间与上一次执行时间
          if (lastExecTime && now - lastExecTime < 1000) {
            // 重复点击，提示用户
            _this.$message.warning('请勿重复点击，每秒仅允许执行一次')
            return
          }
          // 记录上一次执行时间
          lastExecTime = now
          // 执行操作
          _this.$nextTick(() => {
            _this.$emit('execEvent')
          })
        },
        // ctrl+k重写，避免和drap快捷键冲突
        'Ctrl-K': function (cm) {
          return false
        },
      })
      // 获取element form表单的高度，将codemirror的高度设置和表单一致
      cm.setSize('height', `${this.cmHeight}px`)
      // 设置自动补全
      cm.on('keypress', () => {
        cm.showHint({ completeSingle: false })
      })
    },
    getValue() {
      return this.codemirror.getValue()
    },
    getSelection() {
      return this.codemirror.getSelection()
    },
    setOption(key, value) {
      this.codemirror.setOption(key, value)
    },
    setValue(value) {
      this.codemirror.setValue(value)
    },
    setTheme(data) {
      this.codemirror.setOption('theme', data)
    },
    saveCache(key) {
      localStorage.setItem(key, this.codemirror.getValue())
    },
    loadCache(key) {
      var cache = localStorage.getItem(key)
      if (cache != null) {
        this.codemirror.setValue(cache)
      }
    },
  },
}
</script>

<style lang="less" scoped>
@font-face {
  font-family: 'JetBrains Mono NL';
  src: url('./font/JetBrainsMonoNL-Light.ttf');
}
::v-deep .CodeMirror {
  font-size: 14px;
  font-family: 'JetBrains Mono NL', monospace;
}
</style>
