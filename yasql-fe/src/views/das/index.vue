<template>
  <el-card class="box-card">
    <a-tabs type="card" size="small">
      <a-tab-pane key="1" tab="编辑器">
        <editComponent ref="editComponentChild"></editComponent
      ></a-tab-pane>
      <a-tab-pane key="3" tab="收藏SQL">
        <favoriteComponent ref="favoriteComponentChild"></favoriteComponent>
      </a-tab-pane>
      <a-tab-pane key="2" tab="历史查询">
        <historyComponent ref="editComponentChild"></historyComponent>
      </a-tab-pane>
    </a-tabs>
  </el-card>
</template>

<script>
// 导入子组件
import editComponent from './edit.vue';
import historyComponent from './history.vue';
import favoriteComponent from './favorite.vue';

export default {
  components: {
    editComponent,
    historyComponent,
    favoriteComponent,
  },
};
</script>

<style  lang='less' scoped>
// 改变tabs底部高度
::v-deep .ant-tabs-bar {
  margin: 0 0 8px 0 !important;
}

/deep/ .el-card__body {
  padding: 8px;
}
</style>