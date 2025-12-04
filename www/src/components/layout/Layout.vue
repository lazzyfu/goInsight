<template>
  <div>
    <a-layout class="layout" style="min-height: 100vh">
      <a-layout-header class="header">
        <a-row type="flex">
          <a-col :flex="12">
            <img src="@/assets/logo.svg" width="100" height="60" />
          </a-col>
          <a-col :flex="4">
            <menu-unfold-outlined v-if="data.collapsed" @click="clickScope" />
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
                  <a-avatar v-else src="/assets/account.png" icon="user" />
                </a-badge>
                <span style="padding-left: 8px"> {{ userStore.nickname }}</span>
              </span>
              <template v-slot:overlay>
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
            :items="data.menus"
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
import { GetUserProfileApi } from '@/api/login'
import router from '@/router'
import { useAsyncRouterStore } from '@/store/static-router'
import { useUserStore } from '@/store/user'
import {
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  SafetyOutlined,
  UserOutlined,
} from '@ant-design/icons-vue'
import { defineAsyncComponent, h, onMounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'

import UserPassword from '@/views/account/settings/UserPassword.vue'

const route = useRoute()

// 动态加载图标组件
const renderIcon = (iconName) => {
  if (!iconName) return undefined
  const iconComponent = defineAsyncComponent(() =>
    import('@ant-design/icons-vue').then((module) => module[iconName]),
  )
  return h(iconComponent)
}

const userStore = useUserStore()
const asyncRouterStore = useAsyncRouterStore()

const state = reactive({
  openPasswordModal: false,
})

const data = reactive({
  menus: [],
  items: [],
  openKeys: [],
  collapsed: false,
})

// 转换路由数据为菜单项
const transformRoutesToMenuItems = (routes) => {
  return (
    routes
      .filter((route) => !route.meta?.hidden)
      .map((route) => ({
        key: route.path,
        label: route.meta?.title,
        title: route.meta?.title,
        icon: route.icon ? renderIcon(route.icon) : undefined,
        children: route.children?.length ? transformRoutesToMenuItems(route.children) : undefined,
      }))
  )
}

const initializeLayoutData = async () => {
  await GetUserProfileApi().then((res) => {
    if (res.code === '0000') {
      userStore.setUid(res.data.uid)
      userStore.setUserName(res.data.username)
      userStore.setNickName(res.data.nick_name)
      userStore.setUserAvatar(res.data.avatar_file)
      userStore.setUserEmail(res.data.email)
      userStore.setUserMobile(res.data.mobile)
      userStore.setUserOrganization(res.data.organization)
      userStore.setUserRole(res.data.role)
      userStore.setUserDateJoined(res.data.date_joined)
    }
  })

  const routes = asyncRouterStore.addRouters.find((item) => item.path === '/')
  const rootRoutes = (routes && routes.children) || []

  data.menus = transformRoutesToMenuItems(rootRoutes)
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

const collapsed = ref(false)
const clickScope = () => {
  data.collapsed = !data.collapsed
  collapsed.value = data.collapsed
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
  userStore.clear()
  router.push({ name: 'Login' })
}

onMounted(async () => {
  await initializeLayoutData()
})
</script>

<style scoped>
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
