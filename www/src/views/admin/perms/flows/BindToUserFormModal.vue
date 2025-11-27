<template>
  <a-modal :open="props.open" title="绑定审批流到用户" :width="560" centered destroyOnClose @cancel="handleCancel">
    <template #footer>
      <a-button @click="handleCancel">取消</a-button>
      <a-button type="primary" :loading="loading" @click="onSubmit">确定绑定</a-button>
    </template>

    <a-form ref="formRef" :model="formState" :rules="rules" layout="vertical" class="bind-form">
      <a-form-item label="审批流" name="approval_id">
        <a-select v-model:value="formState.approval_id" placeholder="请选择要分配的审批流" :options="props.flowOptions"
          show-search :filter-option="filterOption" allow-clear style="width: 100%" />
        <template #extra>
          <p>请选择一个已经定义好的审批流程。</p>
        </template>
      </a-form-item>

      <a-form-item label="选择用户" name="users">
        <a-select v-model:value="formState.users" mode="multiple" placeholder="请选择要绑定的用户" :options="props.userOptions"
          show-search :filter-option="filterOption" allow-clear style="width: 100%" />
        <template #extra>
          <p>这些用户发起审批时，将默认使用以上流程。</p>
        </template>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { reactive, ref } from 'vue'

// 定义props和emits
const props = defineProps({
  open: Boolean,
  flowOptions: { type: Array, default: () => [] },
  userOptions: { type: Array, default: () => [] },
})
const emit = defineEmits(['update:open', 'submit'])

// 表单引用
const formRef = ref()

// 状态
const uiState = reactive({
  loading: false
})

// form表单
const formState = reactive({
  approval_id: undefined,
  users: [],
})

// 表单校验规则
const rules = {
  approval_id: [{ required: true, message: '请选择审批流', trigger: 'change' }],
  users: [
    {
      required: true,
      type: 'array',
      validator: async (rule, value) => {
        if (!value || value.length === 0) {
          throw new Error('请选择至少一个用户进行绑定')
        }
      },
      trigger: 'change',
    },
  ],
}

// 搜索逻辑
const filterOption = (input, option) => {
  return (
    option.label.toLowerCase().includes(input.toLowerCase()) ||
    option.value.toLowerCase().includes(input.toLowerCase())
  )
}

// 取消按钮
const handleCancel = () => {
  emit('update:open', false)
}

// 提交表单
const onSubmit = async () => {
  try {
    await formRef.value.validate()
    uiState.loading = true

    const payload = {
      users: formState.users,
      approval_id: formState.approval_id,
    }

    emit('submit', payload)
  } catch (err) {
  } finally {
    uiState.loading = false
  }
}
</script>

<style scoped>
.bind-form {
  padding: 16px 0;
}
</style>
