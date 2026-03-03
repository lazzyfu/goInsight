<template>
  <a-card class="gi-page-card-shell" :title="title" :bordered="bordered">
    <template v-if="$slots.extra" #extra>
      <slot name="extra" />
    </template>
    <slot />
  </a-card>
</template>

<script setup>
defineOptions({ name: 'PageCardShell' })

defineProps({
  title: {
    type: String,
    default: '',
  },
  bordered: {
    type: Boolean,
    default: true,
  },
})
</script>
