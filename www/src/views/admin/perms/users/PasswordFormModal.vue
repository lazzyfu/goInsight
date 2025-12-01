<template>
  <a-modal :open="props.open" :title="props.title" @cancel="handleCancel">
    <template #footer>
      <a-button @click="handleCancel">取消</a-button>
      <a-button type="primary" :loading="uiState.loading" @click="onSubmit">确定</a-button>
    </template>

    <a-form
      ref="formRef"
      :model="formState"
      :rules="rules"
      :label-col="{ span: 4 }"
      :wrapper-col="{ span: 20 }"
    >
      <a-form-item label="新密码" has-feedback name="password">
        <a-input v-model:value="formState.password" autocomplete="off" type="password"> </a-input>
      </a-form-item>
      <a-form-item label="确认密码" has-feedback name="verify_password">
        <a-input v-model:value="formState.verify_password" autocomplete="off" type="password">
        </a-input>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { regPassword } from '@/utils/validate'
import { reactive, ref } from 'vue'

// 定义props和emits
const props = defineProps({
  open: Boolean,
  title: String,
})
const emit = defineEmits(['update:open', 'submit'])

// 表单数据
const formState = reactive({
  password: '',
  verify_password: '',
})

// 表单引用
const formRef = ref()

// 状态
const uiState = reactive({
  loading: false,
})

// 有效性验证
const validateNewPass = async (_rule, value) => {
  if (!value) return Promise.reject('请输入密码')
  if (!regPassword.test(value))
    return Promise.reject('密码至少7个字符,必须包含大写字母、小写字母、数字和特殊字符')
  if (formState.verify_password) formRef.value.validateFields('verify_password')
  return Promise.resolve()
}

const validateVerifyPass = async (_rule, value) => {
  if (!value) return Promise.reject('请输入密码')
  if (!regPassword.test(value))
    return Promise.reject('密码至少7个字符,必须包含大写字母、小写字母、数字和特殊字符')
  if (value !== formState.password) return Promise.reject('两次输入的密码不一致')
  return Promise.resolve()
}

// 表单校验规则
const rules = {
  password: [
    {
      required: true,
      validator: validateNewPass,
      trigger: 'change',
    },
  ],
  verify_password: [
    {
      required: true,
      validator: validateVerifyPass,
      trigger: 'change',
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
    emit('submit', formState)
  } catch (err) {
  } finally {
    uiState.loading = false
  }
}
</script>
