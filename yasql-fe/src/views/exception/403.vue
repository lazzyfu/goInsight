<template>
  <a-card>
    <a-result status="403" title="403" sub-title="抱歉，您没有权限访问该页面">
      <template #extra>
        <a-button type="primary">
          <router-link :to="{ name: 'menu.home' }">返回个人页面</router-link>
        </a-button>
      </template>
    </a-result>
  </a-card>
</template>

<script>
export default {
  name: 'Exception403',
}
</script>
