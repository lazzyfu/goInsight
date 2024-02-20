<template>
  <a-card>
    <a-tabs size="small">
      <a-tab-pane key="1">
        <span slot="tab" >
          <a-icon type="form" />
          编辑器
        </span>
        <EditComponent></EditComponent
      ></a-tab-pane>
      <a-tab-pane key="2">
        <span slot="tab" >
          <a-icon type="heart" />
          收藏SQL
        </span>
        <FavoriteComponent></FavoriteComponent>
      </a-tab-pane>
      <a-tab-pane key="3">
        <span slot="tab" >
          <a-icon type="history" />
          历史查询
        </span>
        <HistoryComponent></HistoryComponent>
      </a-tab-pane>
    </a-tabs>
  </a-card>
</template>

<script>
// 导入子组件
import EditComponent from './edit/Index'
import HistoryComponent from './history/Index'
import FavoriteComponent from './favorite/Index'

export default {
  components: {
    EditComponent,
    HistoryComponent,
    FavoriteComponent,
  },
}
</script>

<style scoped>
body .anticon {
  transform: translate(-0.5px, -3px);
}
</style>
