<template>
  <div ref="editor" :initVal="initVal" style="height: 100%"></div>
</template>

<script setup>
import { autocompletion, completionKeymap } from '@codemirror/autocomplete'
import { historyKeymap, indentWithTab } from '@codemirror/commands'
import { MySQL, sql, StandardSQL } from '@codemirror/lang-sql'
import { Compartment, EditorState } from '@codemirror/state'
import { EditorView, keymap } from '@codemirror/view'
import { basicSetup } from 'codemirror'
import { format } from 'sql-formatter'
import { onBeforeUnmount, onMounted, ref } from 'vue'

const editor = ref(null) // template ref
const editorView = ref(null) // 编辑器实例

// 定义props
const props = defineProps({
  initVal: {
    type: String,
    default: 'SELECT * FROM ',
  },
})

const languageCompartment = new Compartment()
const editableCompartment = new Compartment()
const readonlyCompartment = new Compartment()

// 初始化扩展
const fixedExtensions = [
  basicSetup, // basicSetup 包含了 history, defaultKeymap 等基础配置
  keymap.of([
    // ...completionKeymap 包含了处理自动补全的快捷键 (如回车、上下箭头)
    ...completionKeymap,
    // ...historyKeymap 提供了撤销/重做功能
    ...historyKeymap,
    indentWithTab,
  ]),
  // 启用自动补全
  autocompletion(),
]

// 初始化编辑器
const initEditor = () => {
  if (editorView.value) {
    editorView.value.destroy()
  }
  const startState = EditorState.create({
    doc: props.initVal,
    extensions: [
      fixedExtensions,
      languageCompartment.of(sql({ dialect: StandardSQL })),
      readonlyCompartment.of(EditorState.readOnly.of(false)),
      editableCompartment.of(EditorView.editable.of(true)),
    ],
  })
  editorView.value = new EditorView({
    state: startState,
    parent: editor.value,
  })
  editorView.value.dom.style.border = '1px solid #f0f0f0'
  editorView.value.dom.style.height = '470px'
}

// 设置为只读，但可选中、可复制
const setReadonly = (readonly) => {
  if (!editorView.value) return
  editorView.value.dispatch({
    effects: readonlyCompartment.reconfigure(EditorState.readOnly.of(readonly)),
  })
}

// 设置编辑器高度
const setHeight = (height) => {
  if (!editorView.value) return
  // 支持数值或字符串两种传入形式
  const value = typeof height === 'number' ? `${height}px` : height
  editorView.value.dom.style.height = value
}

// 设置自动补全
const setCompletion = (completionData) => {
  if (!editorView.value) return

  editorView.value.dispatch({
    effects: languageCompartment.reconfigure(sql({ dialect: MySQL, schema: completionData })),
  })
}

// 格式化内容
const formatContent = () => {
  if (!editorView.value) return
  const formatted = format(editorView.value.state.doc.toString())
  editorView.value.dispatch({
    changes: { from: 0, to: editorView.value.state.doc.length, insert: formatted },
  })
}

// 获取全部内容
const getContent = () => {
  return editorView.value.state.doc.toString()
}

// 获取选中内容
const getSelectedText = () => {
  let text = ''
  for (let range of editorView.value.state.selection.ranges) {
    if (!range.empty) {
      text += editorView.value.state.doc.sliceString(range.from, range.to)
    }
  }
  return text
}

// 设置内容
const setContent = (content) => {
  if (!editorView.value) return
  editorView.value.dispatch({
    changes: { from: 0, to: editorView.value.state.doc.length, insert: content },
  })
}

onMounted(() => {
  initEditor()
})

onBeforeUnmount(() => {
  if (editorView.value) {
    editorView.value.destroy()
  }
})

defineExpose({
  formatContent,
  getContent,
  getSelectedText,
  setCompletion,
  setContent,
  editorView,
  setReadonly,
  setHeight,
})
</script>

<style scoped>
:deep(.cm-gutters) {
  color: #9f9d9d;
  background-color: #f7f7f7;
}

:deep(.cm-editor.cm-focused) {
  outline: none !important; /* 去掉浏览器默认的焦点高亮，CodeMirror 6 的容器 <div class="cm-editor"> 在获得焦点时，默认会被加上 */
}
</style>
