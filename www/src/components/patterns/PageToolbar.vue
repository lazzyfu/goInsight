<template>
  <div class="gi-page-toolbar gi-pattern-toolbar">
    <div class="gi-page-toolbar-left">
      <slot name="left">
        <slot />
      </slot>
    </div>
    <div v-if="$slots.right" class="gi-page-toolbar-right">
      <slot name="right" />
    </div>
  </div>
</template>

<script setup>
defineOptions({ name: 'PageToolbar' })
</script>
