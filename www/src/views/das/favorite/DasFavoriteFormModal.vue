<template>
  <a-modal :open="props.open" title="收藏SQL语句" width="45%" @cancel="handleCancel">
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
      <a-form-item label="标题" name="title">
        <a-input v-model:value="formData.title" placeholder="为这条SQL起个名字" allow-clear />
      </a-form-item>
      <a-form-item label="SQL内容" name="sqltext">
        <a-textarea
          v-model:value="formData.sqltext"
          :rows="10"
          placeholder="请输入SQL"
          allow-clear
        />
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
  title: [
    {
      required: true,
      message: '请输入SQL标题',
      trigger: 'blur',
    },
    {
      min: 3,
      max: 100,
      message: '标题长度应为3-100个字符',
      trigger: 'blur',
    },
  ],
  sqltext: [
    {
      required: true,
      message: '请输入SQL语句',
      trigger: 'blur',
    },
    {
      min: 10,
      message: 'SQL语句至少需要10个字符',
      trigger: 'blur',
    },
  ],
}

// 取消按钮处理函数
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
