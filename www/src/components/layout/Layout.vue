<template>
  <a-layout class="layout" :style="layoutVars">
    <a-layout-sider
      class="layout-sider"
      :width="layoutConfig.sidebarExpandedWidth"
      :collapsed-width="uiState.isMobile ? 0 : layoutConfig.sidebarCollapsedWidth"
      breakpoint="lg"
      v-model:collapsed="uiState.collapsed"
      :trigger="null"
      collapsible
      @breakpoint="handleBreakpoint"
    >
      <div class="logo-wrap">
        <img class="logo" src="@/assets/logo.svg" />
      </div>
      <a-menu
        class="layout-menu"
        theme="dark"
        mode="inline"
        :openKeys="uiData.openKeys"
        :selectedKeys="[route.path]"
        :items="menuItems"
        @select="select"
        @openChange="openChange"
      />
    </a-layout-sider>

    <a-layout class="layout-main">
      <a-layout-header class="layout-header">
        <div class="header-left">
          <MenuUnfoldOutlined v-if="uiState.collapsed" class="trigger" @click="toggleCollapse" />
          <MenuFoldOutlined v-else class="trigger" @click="toggleCollapse" />
          <a-breadcrumb class="layout-breadcrumb">
            <a-breadcrumb-item v-for="(item, index) in route.matched" :key="item.name">
              <router-link
                v-if="item.meta.title && index !== route.matched.length - 1"
                :to="{ path: item.path === '' ? '/' : item.path }"
              >
                {{ item.meta.title }}
              </router-link>
              <span v-else>{{ item.meta.title }}</span>
            </a-breadcrumb-item>
          </a-breadcrumb>
        </div>

        <div class="header-right">
          <a-dropdown placement="bottomRight">
            <span class="user-info">
              <a-avatar v-if="userStore.avatar" :src="userStore.avatar" />
              <a-avatar v-else src="/avatar.png" />
              <span v-if="!uiState.isMobile" class="nickname">{{ userStore.nickname }}</span>
            </span>
            <template #overlay>
              <a-menu>
                <a-menu-item @click="userCenter">
                  <UserOutlined />
                  个人中心
                </a-menu-item>
                <a-menu-item @click="uiState.openPasswordModal = true">
                  <SafetyOutlined />
                  修改密码
                </a-menu-item>
                <a-menu-item @click="Logout">
                  <LogoutOutlined />
                  注销登录
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>

      <a-layout-content class="layout-content">
        <div class="layout-content-inner">
          <router-view />
        </div>
      </a-layout-content>
    </a-layout>
  </a-layout>

  <UserPassword :open="uiState.openPasswordModal" @update:open="uiState.openPasswordModal = $event" />
</template>

<script setup>
defineOptions({ name: 'GoLayout' })

import { layoutConfig } from '@/components/layout/layoutConfig'
import router from '@/router'
import { usePermissionStore } from '@/store/permission'
import { useUserStore } from '@/store/user'
import UserPassword from '@/views/account/settings/UserPassword.vue'
import {
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  SafetyOutlined,
  UserOutlined,
} from '@ant-design/icons-vue'
import { computed, defineAsyncComponent, h, onBeforeUnmount, onMounted, reactive } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const userStore = useUserStore()
const permissionStore = usePermissionStore()

const renderIcon = (iconName) => {
  if (!iconName) return undefined
  const iconComponent = defineAsyncComponent(() =>
    import('@ant-design/icons-vue').then((module) => module[iconName]),
  )
  return h(iconComponent)
}

const uiState = reactive({
  collapsed: false,
  isMobile: false,
  openPasswordModal: false,
})

const uiData = reactive({
  openKeys: [],
})

const layoutVars = computed(() => ({
  '--gi-layout-header-height': `${layoutConfig.headerHeight}px`,
  '--gi-layout-sider-width': `${layoutConfig.sidebarExpandedWidth}px`,
  '--gi-layout-sider-width-collapsed': `${layoutConfig.sidebarCollapsedWidth}px`,
}))

const transformRoutesToMenuItems = (routes) =>
  routes
    .filter((route) => !route.meta?.hidden)
    .map((route) => ({
      key: route.path,
      label: route.meta?.title,
      title: route.meta?.title,
      icon: route.meta?.icon ? renderIcon(route.meta.icon) : undefined,
      children: route.children?.length ? transformRoutesToMenuItems(route.children) : undefined,
    }))

const menuItems = computed(() => transformRoutesToMenuItems(permissionStore.menuRoutes))

const initializeLayoutData = async () => {
  const storedOpenKeys = sessionStorage.getItem('openKeys')
  if (storedOpenKeys) {
    uiData.openKeys = JSON.parse(storedOpenKeys)
  }
}

const syncViewport = () => {
  const isMobile = window.innerWidth < layoutConfig.mobileBreakpoint
  uiState.isMobile = isMobile
  if (isMobile) {
    uiState.collapsed = true
  }
}

const handleBreakpoint = (broken) => {
  uiState.isMobile = broken
  if (broken) {
    uiState.collapsed = true
  }
}

const select = (value) => {
  router.push({ path: value.key })
  if (uiState.isMobile) {
    uiState.collapsed = true
  }
}

const openChange = (openKeys) => {
  uiData.openKeys = openKeys || []
  sessionStorage.setItem('openKeys', JSON.stringify(uiData.openKeys))
}

const toggleCollapse = () => {
  uiState.collapsed = !uiState.collapsed
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
  userStore.clear()
  permissionStore.reset()
  router.push({ name: 'Login' }).then(() => {
    window.location.reload()
  })
}

onMounted(async () => {
  await initializeLayoutData()
  syncViewport()
  window.addEventListener('resize', syncViewport)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', syncViewport)
})
</script>

<style scoped>
.layout {
  min-height: 100vh;
  background: var(--gi-color-page-bg);
}

.layout-sider {
  height: 100vh;
  overflow-y: auto;
  box-shadow: 2px 0 10px rgba(5, 17, 25, 0.18);
}

.layout-sider :deep(.ant-layout-sider-children) {
  display: flex;
  flex-direction: column;
}

.logo-wrap {
  height: 72px;
  padding: var(--gi-spacing-md) var(--gi-spacing-sm);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo {
  width: auto;
  height: 38px;
  max-width: 170px;
  object-fit: contain;
}

.layout-menu {
  flex: 1;
  padding-top: var(--gi-spacing-sm);
}

.layout-main {
  min-width: 0;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.layout-header {
  height: var(--gi-layout-header-height);
  line-height: var(--gi-layout-header-height);
  padding-inline: var(--gi-spacing-md);
  border-bottom: 1px solid color-mix(in srgb, var(--gi-color-border), #ffffff 24%);
  background: var(--gi-color-container-bg);
  box-shadow: var(--gi-shadow-sm);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--gi-spacing-sm);
  position: sticky;
  top: 0;
  z-index: 20;
}

.header-left,
.header-right {
  display: flex;
  align-items: center;
}

.header-left {
  min-width: 0;
  gap: var(--gi-spacing-sm);
}

.layout-breadcrumb {
  min-width: 0;
}

.layout-breadcrumb :deep(.ant-breadcrumb) {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.trigger {
  font-size: 18px;
  cursor: pointer;
  color: var(--gi-color-text-secondary);
  transition: color var(--gi-duration-fast) ease;
}

.trigger:hover {
  color: var(--gi-color-primary);
}

.user-info {
  display: inline-flex;
  align-items: center;
  gap: var(--gi-spacing-sm);
  cursor: pointer;
}

.nickname {
  max-width: 180px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  color: var(--gi-color-text-primary);
}

.layout-content {
  flex: 1;
  padding: var(--gi-spacing-md);
  background: var(--gi-color-page-bg);
  overflow: auto;
}

.layout-content-inner {
  min-height: calc(
    100vh - var(--gi-layout-header-height) - var(--gi-spacing-md) - var(--gi-spacing-md)
  );
}

@media (min-width: 1440px) {
  .layout-content {
    padding: var(--gi-spacing-xl);
  }

  .layout-content-inner {
    min-height: calc(
      100vh - var(--gi-layout-header-height) - var(--gi-spacing-xl) - var(--gi-spacing-xl)
    );
  }
}

@media (max-width: 1023px) {
  .layout-header {
    padding-inline: var(--gi-spacing-sm);
  }

  .layout-content {
    padding: var(--gi-spacing-sm);
  }

  .layout-content-inner {
    min-height: calc(
      100vh - var(--gi-layout-header-height) - var(--gi-spacing-sm) - var(--gi-spacing-sm)
    );
  }
}

@media (max-width: 767px) {
  .layout-content {
    padding: var(--gi-spacing-ssm);
  }

  .layout-content-inner {
    min-height: calc(
      100vh - var(--gi-layout-header-height) - var(--gi-spacing-ssm) - var(--gi-spacing-ssm)
    );
  }
}
</style>
