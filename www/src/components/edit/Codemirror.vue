<template>
  <!-- resize 的唯一尺寸源 -->
  <div ref="wrapper" class="cm-wrapper" :initVal="props.initVal" :height="props.height">
    <!-- CodeMirror 挂载点 -->
    <div ref="editor" style="height: 100%"></div>
  </div>
</template>

<script setup>
defineOptions({ name: 'GoCodemirror' })

import { autocompletion, completionKeymap } from '@codemirror/autocomplete'
import { historyKeymap, indentWithTab } from '@codemirror/commands'
import { MySQL, sql, StandardSQL } from '@codemirror/lang-sql'
import { Compartment, EditorState } from '@codemirror/state'
import { EditorView, keymap } from '@codemirror/view'
import { basicSetup } from 'codemirror'
import { format } from 'sql-formatter'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'

/* refs */
const wrapper = ref(null)
const editor = ref(null)
const editorView = ref(null)
let resizeObserver = null

/* props（原样保留） */
const props = defineProps({
  initVal: {
    type: String,
    default: '',
  },
  height: {
    type: String,
    default: '200px',
  },
})

/* compartments */
const languageCompartment = new Compartment()
const editableCompartment = new Compartment()
const readonlyCompartment = new Compartment()

/* 固定扩展 */
const fixedExtensions = [
  basicSetup,
  keymap.of([
    ...completionKeymap,
    ...historyKeymap,
    indentWithTab,
  ]),
  autocompletion(),
  EditorView.lineWrapping,
]

/* 初始化编辑器（initVal 只在初始化时生效，行为不变） */
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
}

/* 设置只读 */
const setReadonly = (readonly) => {
  if (!editorView.value) return
  editorView.value.dispatch({
    effects: readonlyCompartment.reconfigure(EditorState.readOnly.of(readonly)),
  })
}

/* 设置高度（外部强制覆盖） */
const setHeight = (height) => {
  if (!wrapper.value) return
  wrapper.value.style.height = height
  editorView.value?.requestMeasure()
}

/* 设置自动补全 */
const setCompletion = (completionData) => {
  if (!editorView.value) return
  editorView.value.dispatch({
    effects: languageCompartment.reconfigure(
      sql({ dialect: MySQL, schema: completionData }),
    ),
  })
}

/* 格式化 */
const formatContent = () => {
  if (!editorView.value) return
  const formatted = format(editorView.value.state.doc.toString())
  editorView.value.dispatch({
    changes: {
      from: 0,
      to: editorView.value.state.doc.length,
      insert: formatted,
    },
  })
}

/* 获取内容 */
const getContent = () => {
  return editorView.value?.state.doc.toString() || ''
}

/* 获取选中内容 */
const getSelectedText = () => {
  if (!editorView.value) return ''
  let text = ''
  for (const range of editorView.value.state.selection.ranges) {
    if (!range.empty) {
      text += editorView.value.state.doc.sliceString(range.from, range.to)
    }
  }
  return text
}

/* 设置内容 */
const setContent = (content) => {
  if (!editorView.value) return
  editorView.value.dispatch({
    changes: {
      from: 0,
      to: editorView.value.state.doc.length,
      insert: content,
    },
  })
}


/* lifecycle */
onMounted(() => {
  // height 初始生效
  if (wrapper.value) {
    wrapper.value.style.height = props.height
  }

  initEditor()

  // resize 时通知 CodeMirror 重算
  resizeObserver = new ResizeObserver(() => {
    editorView.value?.requestMeasure()
  })

  resizeObserver.observe(wrapper.value)
})

onBeforeUnmount(() => {
  resizeObserver?.disconnect()
  editorView.value?.destroy()
})

/* 外部动态修改 height 时仍然生效 */
watch(
  () => props.height,
  (val) => {
    if (!wrapper.value) return
    wrapper.value.style.height = val
    editorView.value?.requestMeasure()
  },
)

/* expose（原样保留） */
defineExpose({
  formatContent,
  getContent,
  getSelectedText,
  setCompletion,
  setContent,
  setReadonly,
  setHeight,
  editorView,
})
</script>

<style scoped>
.cm-wrapper {
  resize: vertical;
  overflow: auto;
  border: 1px solid #f0f0f0;
}

/* CodeMirror 填满容器 */
.cm-wrapper :deep(.cm-editor) {
  height: 100%;
}

/* gutter 样式 */
:deep(.cm-gutters) {
  color: #9f9d9d;
  background-color: #f7f7f7;
}

/* 去焦点 outline */
:deep(.cm-editor.cm-focused) {
  outline: none !important;
}
</style>
