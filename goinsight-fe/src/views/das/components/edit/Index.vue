<template>
  <div>
    <a-card>
      <a-row type="flex" class="fold-left-box">
        <a-col class="fold-left-box-left" :style="{ width: asideWidth + 'px' }" v-show="asideWidth > 0">
          <EditLeftZoneComponent
            @changeSelectedSchema="ChangeSelectedSchema"
            @setTabCompletion="SetTabCompletion"
            :screenHeight="screenHeight"
          ></EditLeftZoneComponent>
        </a-col>
        <a-col class="fold-left-box-line" :style="{ cursor: asideWidth === 0 ? '' : 'col-resize' }">
          <div class="fold-left-box-line-button">
            <a-button
              shape="circle"
              icon="right-circle"
              type="link"
              v-if="asideWidth === 0"
              @click="foldLeft"
            ></a-button>
          </div>
        </a-col>
        <a-col class="fold-left-box-main" ref="screenHeight">
          <a-tabs v-model="activeKey" type="editable-card" size="small" @edit="onEdit" @change="changeTab">
            <a-tab-pane v-for="pane in panes" :key="pane.key" :tab="pane.title" :closable="pane.closable">
              <EditRightZoneComponent
                :selectedSchema="selectedSchema"
                :tabIndex="tabIndex"
                @showBootstrapTable="showBootstrapTable"
                @renderBootstrapTable="renderBootstrapTable"
                ref="EditRightZoneComponent"
              ></EditRightZoneComponent>
            </a-tab-pane>
          </a-tabs>
        </a-col>
      </a-row>
    </a-card>
    <div v-show="showbTable">
      <EditBootstrapTableComponent ref="EditBootstrapTableComponent"></EditBootstrapTableComponent>
    </div>
  </div>
</template>

<script>
// 导入子组件
import EditLeftZoneComponent from './EditLeftZone.vue'
import EditRightZoneComponent from './EditRightZone.vue'
import elementResizeDetectorMaker from 'element-resize-detector'
import EditBootstrapTableComponent from './EditBootstrapTable'
const erd = elementResizeDetectorMaker()

export default {
  components: {
    EditLeftZoneComponent,
    EditRightZoneComponent,
    EditBootstrapTableComponent,
  },
  data() {
    const panes = [
      {
        title: 'Tab 1',
        content: 'Content of Tab 1',
        key: '1',
        closable: false,
      },
    ]
    return {
      showbTable: false,
      asideWidth: 450,
      moveLen: 0, // 拖拽距离
      selectedSchema: {},
      screenHeight: '',
      activeKey: panes[0].key,
      panes,
      newTabIndex: 2,
      tabCompletion: {},
    }
  },
  computed: {
    tabIndex() {
      // 返回activeKey在panes中的索引位置
      return this.panes.findIndex((v) => {
        return v.key == this.activeKey
      })
    },
  },
  methods: {
    showBootstrapTable(value) {
      this.showbTable = value
    },
    renderBootstrapTable(data) {
      this.$nextTick(() => {
        this.$refs.EditBootstrapTableComponent.renderbTable(data)
      })
    },
    // 切换Schema
    ChangeSelectedSchema(data) {
      this.selectedSchema = data
    },
    // 设置自动补全
    SetTabCompletion(data) {
      this.tabCompletion = data
      this.$nextTick(() => {
        this.$refs.EditRightZoneComponent[this.tabIndex].setTabCompletion(this.tabCompletion)
      })
    },
    foldLeft() {
      this.asideWidth = this.asideWidth === 0 ? (this.moveLen <= 64 ? 426 : this.moveLen) : 0
    },
    // 左右Drag
    bindDrop() {
      const _this = this
      var drag = document.getElementsByClassName('fold-left-box-line')
      var left = document.getElementsByClassName('fold-left-box-left')
      var right = document.getElementsByClassName('fold-left-box-main')
      var box = document.getElementsByClassName('fold-left-box')
      for (let i = 0; i < drag.length; i++) {
        // 鼠标按下事件
        drag[i].onmousedown = function (e) {
          var startX = e.clientX
          drag[i].left = drag[i].offsetLeft
          // 鼠标拖动事件
          document.onmousemove = function (e) {
            var endX = e.clientX
            // 移动的距离（endx-startx）+ drag[i].left+移动的距离=左边区域最后的宽度
            var moveLen = drag[i].left + (endX - startX)
            // 设置右边区域宽度
            if (moveLen > box[i].clientWidth * 0.55) moveLen = box[i].clientWidth * 0.55
            // 设置左侧区域的宽度
            drag[i].style.left = moveLen
            for (let j = 0; j < left.length; j++) {
              left[j].style.width = moveLen + 'px'
              right[j].style.width = box[i].clientWidth - moveLen + 'px'
            }
            // 记录拖拽距离
            _this.moveLen = moveLen
            // 拖拽时距离最左边很近时自动贴合
            _this.asideWidth = moveLen
            if (moveLen < 64) {
              document.onmouseup()
              _this.asideWidth = 0
            }
          }
          // 鼠标松开事件
          document.onmouseup = function () {
            document.onmousemove = null
            document.onmouseup = null
            drag[i].releaseCapture && drag[i].releaseCapture()
          }
          drag[i].setCapture && drag[i].setCapture()
          return false
        }
      }
    },
    // 编辑tabs
    onEdit(targetKey, action) {
      this[action](targetKey)
    },
    // 增加tab
    add() {
      const panes = this.panes
      const activeKey = this.newTabIndex++
      panes.push({
        title: `New Tab ${activeKey}`,
        content: `Content of new Tab ${activeKey}`,
        key: activeKey,
      })
      this.panes = panes
    },
    // 切换tab
    changeTab() {
      // 自动设置补全
      this.$nextTick(() => {
        this.$refs.EditRightZoneComponent[this.tabIndex].setTabCompletion(this.tabCompletion)
      })
    },
    // 移除tab
    remove(targetKey) {
      let activeKey = this.activeKey
      let lastIndex
      this.panes.forEach((pane, i) => {
        if (pane.key === targetKey) {
          lastIndex = i - 1
        }
      })
      const panes = this.panes.filter((pane) => pane.key !== targetKey)
      if (panes.length && activeKey === targetKey) {
        if (lastIndex >= 0) {
          activeKey = panes[lastIndex].key
        } else {
          activeKey = panes[0].key
        }
      }
      this.panes = panes
      this.activeKey = activeKey
    },
    // 实时获取高度，并改变左侧区域的高度
    getLeftHeight() {
      erd.listenTo(this.$refs.screenHeight.$el, (element) => {
        var screenHeight = element.offsetHeight - 80
        this.screenHeight = screenHeight + 'px'
      })
    },
  },
  mounted() {
    this.bindDrop()
    this.getLeftHeight()
  },
}
</script>

<style  lang='less' scoped>
// 左侧div样式
.fold-left-box {
  height: 100%;
  overflow: hidden;
  display: flex;
  display: -webkit-flex;
  justify-content: center;

  .fold-left-box-left {
    height: 100%;
    padding-right: 6px;
    overflow: hidden;
  }

  .fold-left-box-line {
    width: 4px;
    display: flex;
    position: relative;
    -webkit-box-pack: center;
    background-size: cover;
    background-position: center;
    display: -webkit-flex;
    align-items: center;
    justify-content: center;
    background-color: #f5f6f7;

    .fold-left-box-line-button {
      position: absolute;
      top: 50%;
      right: -10px;
    }
  }
  .fold-left-box-line:active {
    background-color: #005fb8;
    width: 5px;
  }
  .fold-left-box-main {
    height: 100%;
    flex: 1;
    padding-left: 6px;
    overflow: hidden;
  }
}

// 改变tabs底部高度
/deep/ .ant-tabs-bar {
  margin: 0 0 8px 0 !important;
}

/deep/ .ant-card-body {
  padding: 8px;
}
</style>