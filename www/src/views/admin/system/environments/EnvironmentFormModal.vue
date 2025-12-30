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
      <a-form-item label="环境名" name="name" has-feedback>
        <a-input v-model:value="formData.name" placeholder="请输入环境名" allow-clear />
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
      message: '请输入环境名',
      trigger: ['blur', 'change'],
    },
    {
      validator: (_rule, value) => {
        const v = (value ?? '').trim()
        if (!v) return Promise.reject('环境名不能为空')
        if (v.length < 2 || v.length > 32) return Promise.reject('长度需为 2~32 个字符')
        if (!/^[a-zA-Z0-9\u4e00-\u9fa5_]+$/.test(v)) {
          return Promise.reject('仅可使用字母、数字、中文或下划线')
        }
        if (/^_/.test(v) || /_$/.test(v)) {
          return Promise.reject('不允许以下划线开头或结尾')
        }
        return Promise.resolve()
      },
      trigger: ['blur', 'change'],
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
