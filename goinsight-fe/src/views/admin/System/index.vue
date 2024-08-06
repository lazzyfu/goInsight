<template>
  <a-card>
    <a-tabs type="card" size="small">
      <a-tab-pane key="1" tab="系统环境配置">
        <SystemEnvironmentsComponent></SystemEnvironmentsComponent>
      </a-tab-pane>
      <a-tab-pane key="2" tab="审核参数配置">
        <InspectComponent></InspectComponent>
      </a-tab-pane>
      <a-tab-pane key="3" tab="数据库实例配置">
        <DBConfigsComponent></DBConfigsComponent>
      </a-tab-pane>
      <a-tab-pane key="4" tab="数据查询权限配置">
        <DasComponent></DasComponent>
      </a-tab-pane>
    </a-tabs>
  </a-card>
</template>

<script>
// 导入子组件
import DasComponent from './Das';
import DBConfigsComponent from './DbConfig';
import SystemEnvironmentsComponent from './Environment';
import InspectComponent from './Inspect';

export default {
  components: {
    DBConfigsComponent,
    InspectComponent,
    SystemEnvironmentsComponent,
    DasComponent,
  },
}
</script>

<style lang="less" scoped>
::v-deep .ant-card-body {
  padding: 8px;
}

::v-deep .ant-form {
  margin-bottom: 8px;
}
</style>