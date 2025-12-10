<template>
  <a-modal
    :open="open"
    title="任务详情"
    width="50%"
    @cancel="handleCancel"
  >
    <template #footer>
      <a-button @click="handleCancel">取消</a-button>
    </template>

    <!-- 执行信息 -->
    <a-card size="small" title="执行信息" class="mb-4">
      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :sm="12" :md="8">
          <a-statistic title="执行耗时" :value="formData.execute_cost_time" />
        </a-col>
        <a-col :xs="24" :sm="12" :md="8">
          <a-statistic title="备份耗时" :value="formData.backup_cost_time" />
        </a-col>
        <a-col :xs="24" :sm="12" :md="8">
          <a-statistic title="影响行数" :value="formData.affected_rows" />
        </a-col>
      </a-row>
    </a-card>

    <!-- 导出文件信息 -->
    <a-card
      v-if="formData.sqlType === 'EXPORT'"
      size="small"
      title="导出文件信息"
      class="mb-4"
    >
      <a-row :gutter="[16, 16]">
        <a-col :xs="24">
          <a-statistic title="文件名" :value="formData.file_name" />
        </a-col>
        <a-col :xs="24">
          <a-statistic title="文件大小（字节）" :value="formData.file_size" />
        </a-col>
        <a-col :xs="24">
          <a-statistic title="导出行数" :value="formData.export_rows" />
        </a-col>
        <a-col :xs="24">
          <a-statistic title="文件加密秘钥" :value="formData.encryption_key" />
        </a-col>
        <a-col :xs="24">
          <a-statistic title="文件下载路径" :value="formData.download_url" />
        </a-col>
      </a-row>
    </a-card>

    <!-- 错误信息 -->
    <a-card
      v-if="formData.error"
      size="small"
      title="错误信息"
      class="mb-4"
    >
      <CodeMirror :initVal="formData.error" />
    </a-card>

    <!-- 执行日志 -->
    <a-card
      v-if="formData.execute_log"
      size="small"
      title="执行日志"
      class="mb-4"
    >
      <CodeMirror :initVal="formData.execute_log" />
    </a-card>

    <!-- 回滚 SQL -->
    <a-card
      v-if="formData.rollback_sql"
      size="small"
      title="回滚 SQL"
      class="mb-4"
    >
      <CodeMirror :initVal="formData.rollback_sql" />
    </a-card>
  </a-modal>
</template>

<script setup>
import CodeMirror from '@/components/edit/Codemirror.vue'

const props = defineProps({
  open: Boolean,
})

const emit = defineEmits(['update:open'])

// 任务详情数据
const formData = defineModel('modelValue', {
  type: Object,
  required: true
})

// 关闭
const handleCancel = () => {
  emit('update:open', false)
}
</script>

<style scoped>
/* 统一卡片间距，避免散乱 class */
.mb-4 {
  margin-bottom: 16px;
}
</style>
