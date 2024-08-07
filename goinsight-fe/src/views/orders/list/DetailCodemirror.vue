<template>
  <a-card :bordered="false" title="工单内容">
    <codemirror ref="myCm" v-model="code" :options="cmOptions" @ready="onCmReady"></codemirror>
  </a-card>
</template>

<script>
// CodeMirror
import CodeMirror from 'codemirror'
import 'codemirror/mode/sql/sql.js'
import 'codemirror/addon/selection/active-line'
import 'codemirror/addon/display/autorefresh'
import 'codemirror/addon/comment/comment'
import 'codemirror/addon/edit/matchbrackets'
import 'codemirror/addon/edit/closebrackets'
import 'codemirror/addon/mode/overlay'

// 定义高亮的关键字
let keyword = new RegExp('^create |^drop |^truncate |^update |^delete |^insert into|on update|unique|change|modify|rename', 'i')
CodeMirror.defineMode('highlightText', function (config, parserConfig) {
  let searchOverlay = {
    token: function (stream, state) {
      if (stream.match(keyword)) {
        return 'highlightText'
      }
      while (stream.next() != null && !stream.match(keyword, false)) {
        return null
      }
      return null
    },
  }
  return CodeMirror.overlayMode(CodeMirror.getMode(config, parserConfig.backdrop || 'text/x-mysql'), searchOverlay)
})

export default {
  data() {
    return {
      code: '',
      cmOptions: {
        mode: 'highlightText', // 这里要换成highlightText
        indentUnit: 2,
        tabSize: 2,
        indentWithTabs: true,
        smartIndent: true,
        autoRefresh: true,
        lineNumbers: true,
        styleActiveLine: true,
        lineWrapping: true,
        showCursorWhenSelecting: true,
        readOnly: true,
        focuse: false,
      },
    }
  },
  methods: {
    onCmReady(cm) {
      cm.setSize('height', `550px`)
    },
    setValue(value) {
      this.codemirror.setValue(value)
    },
  },
  computed: {
    codemirror() {
      return this.$refs.myCm.codemirror
    },
  },
}
</script>

