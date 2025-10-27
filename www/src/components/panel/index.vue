<template>
  <div class="split-wrapper">
    <div ref="scalable" class="scalable" :style="{ width: leftWidth }">
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
import { throttle } from 'lodash-es'
import { computed, ref } from 'vue'

export default {
  props: ['leftWidth'],
  setup(props) {
    const leftWidth = computed(() => {
      if (props.leftWidth) {
        return props.leftWidth
      } else {
        return "250px"
      }
    })
    const scalable = ref()
    // 拖拽中
    let startX
    let startWidth
    const onDrag = throttle(function (e) {
      scalable.value && (scalable.value.style.width = `${startWidth + e.clientX - startX}px`)
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
      scalable.value && (startWidth = parseInt(window.getComputedStyle(scalable.value).width, 10))
      document.documentElement.style.userSelect = 'none'
      document.documentElement.addEventListener('mousemove', onDrag)
      document.documentElement.addEventListener('mouseup', dragEnd)
    }

    return {
      leftWidth,
      scalable,
      startDrag
    }
  }
}
</script>

<style lang="less">
@import '@/styles/theme.less';

@classNames: split-wrapper, separator;
.themeBgColor(@classNames);

.split-wrapper {
  display: flex;
  width: 100%;
  height: 100%;

  .scalable {
    position: relative;
    min-width: 100px;
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
      box-shadow:
        -4px -2px 4px -5px rgb(0 0 0 / 35%),
        4px 3px 4px -5px rgb(0 0 0 / 35%);
      cursor: col-resize;

      i {
        width: 2px;
        height: 24px;
        margin: 0 1px;
        background-color: #e9e9e9;
      }
    }
  }

  .right-content {
    background-color: #ffffff;
    flex: 1;
  }

  .left-content,
  .right-content {
    overflow: auto;
  }
}
</style>
