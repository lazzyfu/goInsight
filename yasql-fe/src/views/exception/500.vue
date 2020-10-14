<template>
  <a-card>
    <a-result status="500" title="500" sub-title="Sorry, the server is reporting an error.">
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
  name: 'Exception500',
}
</script>
