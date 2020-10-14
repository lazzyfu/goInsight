<template>
  <a-card>
    <a-result status="404" title="404" sub-title="抱歉，您访问的页面不存在.">
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
  name: 'Exception404',
}
</script>
