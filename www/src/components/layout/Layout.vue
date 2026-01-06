<template>
  <a-layout class="layout" style="min-height: 100vh">
    <!-- 左侧siderbar -->
    <a-layout-sider breakpoint="lg" collapsed-width="0" v-model:collapsed="uiState.collapsed" :trigger="null"
      collapsible>
      <!-- 左上角logo -->
      <img class="app-logo" src="@/assets/logo.svg" />
      <!-- 菜单 -->
      <a-menu theme="dark" mode="inline" :openKeys="uiData.openKeys" :selectedKeys="[route.path]" :items="menuItems"
        @select="select" @openChange="openChange">
      </a-menu>
    </a-layout-sider>

    <a-layout>
      <!-- 折叠图标 -->
      <a-layout-header class="layout-header">
        <div class="header-left">
          <MenuUnfoldOutlined v-if="uiState.collapsed" class="trigger" @click="toggleCollapse" />
          <MenuFoldOutlined v-else class="trigger" @click="toggleCollapse" />
          <!-- 面包屑 -->
          <a-breadcrumb>
            <a-breadcrumb-item v-for="(item, index) in route.matched" :key="item.name">
              <router-link v-if="item.meta.title && index !== route.matched.length - 1"
                :to="{ path: item.path === '' ? '/' : item.path }">{{ item.meta.title }}
              </router-link>
              <span v-else>{{ item.meta.title }}</span>
            </a-breadcrumb-item>
          </a-breadcrumb>
        </div>
        <!-- 右上角用户 -->
        <div class="header-right">
          <a-dropdown placement="bottomRight">
            <span class="user-info">
              <a-avatar v-if="userStore.avatar" :src="userStore.avatar" />
              <a-avatar v-else src="/avatar.png" />
              <span class="nickname">{{ userStore.nickname }}</span>
            </span>
            <template #overlay>
              <a-menu>
                <a-menu-item @click="userCenter">
                  <UserOutlined /> 个人中心
                </a-menu-item>
                <a-menu-item @click="uiState.openPasswordModal = true">
                  <SafetyOutlined /> 修改密码
                </a-menu-item>
                <a-menu-item @click="Logout">
                  <LogoutOutlined /> 注销登录
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>

      <!-- 内容区域 -->
      <a-layout-content :style="{ margin: '24px 16px', padding: '24px', background: '#fff', minHeight: '280px' }">
        <router-view></router-view>
      </a-layout-content>
    </a-layout>
  </a-layout>

  <UserPassword :open="uiState.openPasswordModal" @update:open="uiState.openPasswordModal = $event"></UserPassword>
</template>

<script setup>
defineOptions({ name: 'GoLayout' })

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

const uiState = reactive({
  collapsed: false,
  openPasswordModal: false,
})

const uiData = reactive({
  openKeys: [],
})

// 转换路由数据为菜单项
const transformRoutesToMenuItems = (routes) => {
  return routes
    .filter((route) => !route.meta?.hidden) // 过滤掉 hidden 的路由
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
})

const initializeLayoutData = async () => {
  const storedOpenKeys = sessionStorage.getItem('openKeys')
  if (storedOpenKeys && storedOpenKeys !== null) {
    uiData.openKeys = JSON.parse(storedOpenKeys)
  }
}

const select = (value) => {
  router.push({ path: value.key })
}

const openChange = (openKeys) => {
  if (openKeys?.length > 0) {
    uiData.openKeys = openKeys
    sessionStorage.setItem('openKeys', JSON.stringify(uiData.openKeys))
  }
}

// toggle
const toggleCollapse = () => {
  uiState.collapsed = !uiState.collapsed
}

// 跳转到用户中心
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
.layout-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  /* 左右分开 */
  height: 48px;
  padding: 0 16px;
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
}

.header-left,
.header-right {
  display: flex;
  align-items: center;
}

.trigger {
  font-size: 18px;
  margin-right: 16px;
  cursor: pointer;
}

.trigger:hover {
  color: #1890ff;
}

.nickname {
  margin-left: 8px;
}

.app-logo {
  height: 60px;
  width: auto;
  /* 自动缩放宽度 */
  max-width: 182px;
  object-fit: contain;
  display: block;
}


.site-layout .site-layout-background {
  background: #fff;
}
</style>
