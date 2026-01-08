<template>
  <div class="split-wrapper" :class="{ collapsed: isCollapsed }">
    <!-- 左侧收起后显示拖拽把手，可拖回展开 -->
    <div v-show="isCollapsed" class="collapsed-handle" @mousedown="startDrag"><i></i><i></i></div>

    <div v-show="!isCollapsed" ref="scalable" class="scalable" :style="{ width: leftWidthPx }">
      <div class="left-content">
        <slot name="left-content"> 左边内容区 </slot>
      </div>
      <div ref="separator" class="separator" @mousedown="startDrag"><i></i><i></i></div>
    </div>
    <div class="right-content">
      <slot name="right-content"> 右边内容区 </slot>
    </div>
  </div>
</template>

<script>
export default {
  name: 'SplitPanel',
}
</script>

<script setup>
import { throttle } from 'lodash-es';
import { computed, ref, watch } from 'vue';

const props = defineProps({
  leftWidth: {
    type: String,
    default: '',
  },
  // 可拖到最左侧收起
  collapsible: {
    type: Boolean,
    default: true,
  },
  // 左侧最小宽度（未收起状态）
  minLeftWidth: {
    type: Number,
    default: 100,
  },
  // 小于该阈值则自动收起
  collapseThreshold: {
    type: Number,
    default: 60,
  },
})

const computedLeftWidth = computed(() => props.leftWidth || '250px')

// 由 Vue 控制宽度，避免 slot 更新时覆盖拖拽结果
const leftWidthPx = ref(computedLeftWidth.value)
const isCollapsed = computed(() => leftWidthPx.value === '0px')

watch(
  () => props.leftWidth,
  (v) => {
    // 外部更新 leftWidth 时同步，但不主动“展开”已收起状态
    if (!isCollapsed.value) {
      leftWidthPx.value = v || '250px'
    }
  }
)

const scalable = ref()

// 拖拽中
let startX
let startWidth

const setLeftWidth = (nextWidth) => {
  if (props.collapsible && nextWidth <= props.collapseThreshold) {
    leftWidthPx.value = '0px'
    return
  }
  const minW = props.minLeftWidth
  const clamped = Math.max(minW, nextWidth)
  leftWidthPx.value = `${clamped}px`
}

const onDrag = throttle(function (e) {
  const nextWidth = startWidth + e.clientX - startX
  setLeftWidth(nextWidth)
}, 20)

// 拖拽结束
const dragEnd = () => {
  document.documentElement.style.userSelect = 'unset'
  document.documentElement.removeEventListener('mousemove', onDrag)
  document.documentElement.removeEventListener('mouseup', dragEnd)
}

// 鼠标按下
const startDrag = (e) => {
  startX = e.clientX
  // 使用受控宽度作为起始值（收起时为 0）
  startWidth = parseInt(leftWidthPx.value, 10) || 0
  document.documentElement.style.userSelect = 'none'
  document.documentElement.addEventListener('mousemove', onDrag)
  document.documentElement.addEventListener('mouseup', dragEnd)
}
</script>

<style lang="less">
@import '@/styles/theme.less';

@classNames: split-wrapper, separator, collapsed-handle;
.themeBgColor(@classNames);

.split-wrapper {
  position: relative;
  display: flex;
  width: 100%;
  height: 100%;

  .scalable {
    position: relative;
    min-width: 0;
    max-width: 50vw;
    overflow: auto;

    .left-content {
      height: 100%;
      padding: 5px;
    }

    .separator {
      display: flex;
      position: absolute;
      top: 0;
      right: 0;
      align-items: center;
      justify-content: center;
      width: 14px;
      height: 100%;
      background-color: #f7f7f7;
      box-shadow:
        -4px -2px 4px -5px rgb(0 0 0 / 35%),
        4px 3px 4px -5px rgb(0 0 0 / 35%);
      cursor: col-resize;

      i {
        width: 2px;
        height: 24px;
        margin: 0 1px;
        background-color: #c0c4cc;
      }
    }
  }

  .right-content {
    background-color: #ffffff;
    flex: 1;
    padding-left: 8px;
    box-sizing: border-box;
  }

  .collapsed-handle {
    display: flex;
    position: absolute;
    top: 0;
    left: 0;
    align-items: center;
    justify-content: center;
    width: 16px;
    height: 100%;
    box-sizing: border-box;
    padding: 0 2px;
    border-right: 1px solid #e9e9e9;
    background-color: #f7f7f7;
    cursor: col-resize;
    z-index: 3;

    i {
      width: 2px;
      height: 24px;
      margin: 0 1px;
      background-color: #c0c4cc;
      opacity: 0.7;
    }

    &:hover {
      filter: brightness(0.96);

      i {
        opacity: 1;
      }
    }
  }

  &.collapsed {
    .right-content {
      /* 预留把手宽度，避免遮挡右侧内容（按钮/表格等） */
      padding-left: 22px;
      box-sizing: border-box;
    }
  }

  .left-content,
  .right-content {
    overflow: auto;
  }
}
</style>
