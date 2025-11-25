<template>
  <a-modal :open="props.open" :title="props.title" @cancel="handleCancel">
    <template #footer>
      <a-button @click="handleCancel">取消</a-button>
      <a-button type="primary" :loading="loading" @click="onSubmit">确定</a-button>
    </template>
    <a-form
      ref="formRef"
      :label-col="{ span: 4 }"
      :wrapper-col="{ span: 20 }"
      :model="formData"
      :rules="rules"
    >
      <a-form-item label="环境名" name="name" has-feedback>
        <a-input v-model:value="formData.name" placeholder="请输入环境名" allow-clear />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { ref } from 'vue'

// 定义props和emits
const props = defineProps({
  open: Boolean,
  title: String,
})
const emit = defineEmits(['update:open', 'submit'])

// 使用defineModel接收 v-model:modelValue
// 它自动创建了一个名为modelValue的prop，并提供了一个value来读取，以及一个自动触发update:modelValue的setter
const formData = defineModel('modelValue', {
  type: Object,
  required: true,
})

// 表单引用
const formRef = ref()
const loading = ref(false)

// 表单校验规则
const rules = {
  name: [
    {
      required: true,
      message: '不能为空，请输入环境名',
      trigger: 'blur',
    },
    {
      min: 2,
      max: 32,
      message: `长度应在2~32个字符`,
      trigger: 'blur',
    },
    {
      pattern: /^[a-zA-Z0-9\u4e00-\u9fa5_]+$/,
      message: '只能包含字母、数字、中文或下划线',
      trigger: 'blur',
    },
  ],
}

// 取消按钮处理函数
const handleCancel = () => {
  emit('update:open', false)
  formRef.value?.resetFields()
}

// 提交表单处理函数
const onSubmit = async () => {
  try {
    await formRef.value.validateFields()
    loading.value = true
    emit('submit', formData.value)
  } catch (err) {
  } finally {
    loading.value = false
  }
}
</script>
