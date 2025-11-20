<template>
  <a-modal :open="props.open" :title="props.title" :footer="null" @cancel="handleCancel">
    <a-form
      ref="formRef"
      :label-col="{ span: 4 }"
      :wrapper-col="{ span: 18 }"
      :model="formData"
      :rules="rules"
      @finish="onSubmit"
    >
      <a-form-item label="角色名" name="name" has-feedback>
        <a-input v-model:value="formData.name" placeholder="请输入角色名" allow-clear />
      </a-form-item>

      <a-form-item :wrapper-col="{ offset: 4, span: 18 }" style="text-align: right">
        <a-space>
          <a-button @click="handleCancel">取消</a-button>
          <a-button type="primary" html-type="submit">确定</a-button>
        </a-space>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { ref } from 'vue'

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

// 表单校验规则
const rules = {
  name: [
    {
      required: true,
      message: '不能为空，请输入角色名',
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

const handleCancel = () => {
  emit('update:open', false)
  formRef.value?.resetFields()
}

const onSubmit = async () => {
  emit('submit', formData.value)
}
</script>
