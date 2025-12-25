<template>
  <a-modal :open="props.open" :title="props.title" width="35%" @cancel="handleCancel">
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
      <a-form-item label="描述" name="remark" has-feedback>
        <a-input disabled v-model:value="formData.remark" placeholder="请输入备注" allow-clear />
      </a-form-item>
      <a-form-item
        label="审核参数"
        name="params"
        has-feedback
        help="格式要求为JSON类型，默认为{}，表示继承全局审核参数"
      >
        <a-textarea
          :rows="8"
          v-model:value="formData.params"
          placeholder=" 请输入自定义审核参数，默认为{}"
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
  params: [
    {
      validator: (_rule, value) => {
        const text = (value ?? '').trim()
        if (text === '') return Promise.resolve() // 为空表示继承全局参数
        try {
          const parsed = JSON.parse(text)
          if (parsed === null || typeof parsed !== 'object') {
            return Promise.reject('审核参数必须是有效的 JSON 对象，例如 {}')
          }
          return Promise.resolve()
        } catch {
          return Promise.reject('请输入合法的 JSON，示例：{"enable": true}')
        }
      },
      trigger: ['blur', 'change'],
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
  } catch (err) {
  } finally {
    uiState.loading = false
  }
}
</script>
