const route = {
  name: 'view.admin',
  path: '/admin',
  component: () => import('./index.vue'),
  redirect: '/admin/users',
  meta: { title: '后台管理', icon: 'SettingOutlined', keepAlive: true },
  children: [
    {
      name: `view.admin.perms`,
      path: '/admin/perms',
      component: () => import('./perms/index.vue'),
      meta: { title: '权限管理', icon: 'UserOutlined', keepAlive: true },
      children: [
        {
          name: `view.admin.users`,
          path: '/admin/users',
          component: () => import('./perms/users/UserList.vue'),
          meta: { title: '用户管理', icon: 'UserOutlined', keepAlive: true },
        },
        {
          name: `view.admin.roles`,
          path: '/admin/roles',
          component: () => import('./perms/roles/RoleList.vue'),
          meta: { title: '角色管理', icon: 'IdcardOutlined', keepAlive: true },
        },
        {
          name: `view.admin.orgs`,
          path: '/admin/orgs',
          component: () => import('./perms/orgs/OrgList.vue'),
          meta: { title: '组织管理', icon: 'ApartmentOutlined', keepAlive: true },
        },
        {
          name: `view.admin.flows`,
          path: '/admin/flows',
          component: () => import('./perms/flows/ApprovalFlowList.vue'),
          meta: { title: '审批流', icon: 'AuditOutlined', keepAlive: true },
        },
      ],
    },
    {
      name: `view.admin.system`,
      path: '/admin/system',
      component: () => import('./system/index.vue'),
      meta: { title: '系统配置', icon: 'ToolOutlined', keepAlive: true },
      children: [
        {
          name: `view.admin.environment`,
          path: '/admin/environement',
          icon: 'ClusterOutlined',
          component: () => import('./system/environments/EnvironmentList.vue'),
          meta: { title: '环境管理', icon: 'ClusterOutlined', keepAlive: true },
        },
        {
          name: `view.admin.instance`,
          path: '/admin/instance',
          component: () => import('./system/instances/InstanceList.vue'),
          meta: { title: '实例配置', icon: 'DatabaseOutlined', keepAlive: true },
        },
        {
          name: `view.admin.inspect`,
          path: '/admin/inspect',
          component: () => import('./system/inspect/InspectList.vue'),
          meta: { title: '审核参数', icon: 'CheckCircleOutlined', keepAlive: true },
        },
        {
          name: `view.admin.das`,
          path: '/admin/das',
          component: () => import('./system/das/DasSchemaList.vue'),
          meta: { title: '数据访问', icon: 'SafetyOutlined', keepAlive: true },
        },
        {
          name: `view.admin.das.tables`,
          path: '/admin/das/:username/:schema/:instance_id',
          component: () => import('./system/das/DasTableList.vue'),
          meta: { title: '数据访问', icon: 'CheckCircleOutlined', keepAlive: true, hidden: true },
        },
      ],
    },
  ],
}
export default route
