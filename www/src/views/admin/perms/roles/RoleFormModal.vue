<template>
  <a-modal :open="props.open" :title="props.title" @cancel="handleCancel">
    <template #footer>
      <a-button @click="handleCancel">取消</a-button>
      <a-button type="primary" :loading="uiState.loading" @click="onSubmit">确定</a-button>
    </template>

    <a-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      :label-col="{ span: 4 }"
      :wrapper-col="{ span: 20 }"
      style="margin-top: 24px"
    >
      <a-form-item label="角色名" name="name" has-feedback>
        <a-input v-model:value="formData.name" placeholder="请输入角色名" allow-clear />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { ref, reactive } from 'vue'

// 定义props和emits
const props = defineProps({
  open: Boolean,
  title: String,
})
const emit = defineEmits(['update:open', 'submit'])

// 表单数据
const formData = defineModel('modelValue', {
  type: Object,
  required: true,
})

// 表单引用
const formRef = ref()

// 状态
const uiState = reactive({
  loading: false,
})

// 表单校验规则
const rules = {
  name: [
    {
      required: true,
      message: '角色名不能为空',
      trigger: 'blur',
    },
    {
      min: 2,
      max: 32,
      message: '角色名长度应在2到32个字符之间',
      trigger: 'blur',
    },
    {
      pattern: /^[a-zA-Z0-9\u4e00-\u9fa5_]+$/,
      message: '角色名只能包含字母、数字、中文或下划线',
      trigger: 'blur',
    },
  ],
}

// 取消按钮
const handleCancel = () => {
  emit('update:open', false)
  formRef.value?.resetFields()
}

// 提交表单
const onSubmit = async () => {
  try {
    await formRef.value.validateFields()
    uiState.loading = true
    emit('submit', formData.value)
  } catch {
    // ignore
  } finally {
    uiState.loading = false
  }
}
</script>
