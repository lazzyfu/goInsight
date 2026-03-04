<template>
  <a-config-provider :theme="antdTheme" :locale="zhCN">
    <div id="app" class="gi-app">
      <router-view />
    </div>
  </a-config-provider>
</template>

<script setup>
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import { antdTheme } from '@/design/antdTheme'
</script>

<style scoped>
.gi-app {
  min-height: 100vh;
  background: var(--gi-color-page-bg);
}
</style>
