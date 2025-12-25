<template>
  <div>
    <a-layout class="layout" style="min-height: 100vh">
      <a-layout-header class="header">
        <a-row type="flex">
          <a-col :flex="12">
            <img src="@/assets/logo.svg" width="100" height="60" />
          </a-col>
          <a-col :flex="4">
            <menu-unfold-outlined v-if="data.collapsed" @click="toggleCollapse" />
            <menu-fold-outlined v-else @click="toggleCollapse" />
          </a-col>
          <a-col :flex="200">
            <a-breadcrumb style="margin-top: 22px">
              <a-breadcrumb-item v-for="(item, index) in route.matched" :key="item.name">
                <router-link
                  v-if="item.meta.title && index !== route.matched.length - 1"
                  :to="{ path: item.path === '' ? '/' : item.path }"
                  >{{ item.meta.title }}
                </router-link>
                <span v-else>{{ item.meta.title }}</span>
              </a-breadcrumb-item>
            </a-breadcrumb>
          </a-col>
          <a-col>
            <a-dropdown>
              <span class="ant-dropdown-link">
                <a-badge>
                  <a-avatar v-if="userStore.avatar" :src="userStore.avatar" icon="user" />
                  <a-avatar v-else src="/avatar.png" icon="user" />
                </a-badge>
                <span style="padding-left: 8px"> {{ userStore.nickname }}</span>
              </span>
              <template #overlay>
                <a-menu>
                  <a-menu-item @click="userCenter"> <UserOutlined /> 个人中心 </a-menu-item>
                  <a-menu-item @click="state.openPasswordModal = true">
                    <SafetyOutlined /> 修改密码
                  </a-menu-item>
                  <a-menu-item @click="Logout()"> <LogoutOutlined /> 注销登录 </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </a-col>
        </a-row>
      </a-layout-header>
      <a-layout
        :style="{
          background: '#fff',
          padding: '14px',
          marginTop: '60px',
          marginLeft: data.collapsed ? '80px' : '280px',
          minHeight: '360px',
          transition: 'margin-left 0.2s'
        }"
      >
        <a-layout-sider
          class="layout-sider"
          :collapsed="data.collapsed"
          :collapsed-width="80"
          :width="280"
        >
          <div style="padding-top: 20px"></div>
          <a-menu
            theme="light"
            mode="inline"
            :openKeys="data.openKeys"
            :selectedKeys="[route.path]"
            :items="menuItems"
            @select="select"
            @openChange="openChange"
          >
          </a-menu>
        </a-layout-sider>
        <a-layout-content>
          <router-view></router-view>
        </a-layout-content>
      </a-layout>
    </a-layout>

    <UserPassword
      :open="state.openPasswordModal"
      @update:open="state.openPasswordModal = $event"
    ></UserPassword>
  </div>
</template>

<script setup>
import router from '@/router'
import { usePermissionStore } from '@/store/permission'
import { useUserStore } from '@/store/user'
import {
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  SafetyOutlined,
  UserOutlined,
} from '@ant-design/icons-vue'
import { computed, defineAsyncComponent, h, onMounted, reactive } from 'vue'
import { useRoute } from 'vue-router'

import UserPassword from '@/views/account/settings/UserPassword.vue'

const route = useRoute()
const userStore = useUserStore()
const permissionStore = usePermissionStore()

// 动态加载图标组件
const renderIcon = (iconName) => {
  if (!iconName) return undefined
  const iconComponent = defineAsyncComponent(() =>
    import('@ant-design/icons-vue').then((module) => module[iconName]),
  )
  return h(iconComponent)
}

const state = reactive({
  openPasswordModal: false,
})

const data = reactive({
  openKeys: [],
  collapsed: false,
})

// 转换路由数据为菜单项
const transformRoutesToMenuItems = (routes) => {
  return routes
    .filter(route => !route.meta?.hidden) // 过滤掉 hidden 的路由
    .map((route) => ({
      key: route.path,
      label: route.meta?.title,
      title: route.meta?.title,
      icon: route.meta?.icon ? renderIcon(route.meta.icon) : undefined,
      children: route.children?.length ? transformRoutesToMenuItems(route.children) : undefined,
    }))
}

// 使用 permission store 的 getter
const menuItems = computed(() => {
  return transformRoutesToMenuItems(permissionStore.menuRoutes)
});


const initializeLayoutData = async () => {
    const storedOpenKeys = sessionStorage.getItem('openKeys')
    if (storedOpenKeys && storedOpenKeys !== null) {
      data.openKeys = JSON.parse(storedOpenKeys)
    }
}

const select = (value) => {
  router.push({ path: value.key })
}

const openChange = (openKeys) => {
  if (openKeys?.length > 0) {
    data.openKeys = openKeys
    sessionStorage.setItem('openKeys', JSON.stringify(data.openKeys))
  }
}

const toggleCollapse = () => {
  data.collapsed = !data.collapsed
}

const userCenter = () => {
  router.push({
    name: 'account.basic',
    query: {
      type: 'info',
    },
  })
}

const Logout = () => {
  const permissionStore = usePermissionStore()

  // 清除用户信息
  userStore.clear()

  // 重置权限状态
  permissionStore.reset()

  // 跳转到登录页
  router.push({ name: 'Login' }).then(() => {
    // 刷新页面以清除动态路由
    window.location.reload()
  })
}

// onMounted 仅用于初始化非响应式数据（如 openKeys）
onMounted(async () => {
  await initializeLayoutData()
})
</script>

<style scoped>
/* 样式保持不变 */
.layout-sider {
  overflow: auto;
  height: calc(100vh - 60px);
  position: fixed;
  top: 60px;
  left: 0;
  background: #fff;
  z-index: 100;
  border-right: 1px solid rgb(235, 237, 240);
  transition: all 0.2s;
}

.header {
  position: fixed;
  right: 0;
  top: 0;
  left: 0;
  background: #fff;
  z-index: 999;
  box-shadow: 0 2px 4px 0 var(--cb-color-shadow, rgba(0, 0, 0, 0.16));
  padding: 0 30px;
}

.ant-row {
  display: flex;
  justify-content: flex-start;
}
</style>
