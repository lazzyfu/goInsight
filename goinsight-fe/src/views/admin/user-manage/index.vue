<template>
  <a-card>
    <a-tabs type="card" size="small">
      <a-tab-pane key="1" tab="用户">
        <UsersComponent></UsersComponent>
      </a-tab-pane>
      <a-tab-pane key="2" tab="角色">
        <RolesComponent></RolesComponent>
      </a-tab-pane>
      <a-tab-pane key="3" tab="组织">
        <OrganizationsComponent></OrganizationsComponent>
      </a-tab-pane>
    </a-tabs>
  </a-card>
</template>

<script>
// 导入子组件
import UsersComponent from './users'
import RolesComponent from './roles'
import OrganizationsComponent from './organizations'

export default {
  components: {
    UsersComponent,
    RolesComponent,
    OrganizationsComponent
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