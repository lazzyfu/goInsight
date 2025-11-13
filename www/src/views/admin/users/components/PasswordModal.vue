<template>
  <a-modal :open="open" title="修改密码" width="45%" :footer="null" @cancel="handleCancel">
    <a-form ref="formRef" layout="vertical" :model="formState" :rules="rules" @finish="onSubmit">
      <a-form-item label="新密码" has-feedback name="password">
        <a-input v-model:value="formState.password" autocomplete="off" type="password"> </a-input>
      </a-form-item>
      <a-form-item label="确认密码" has-feedback name="verify_password">
        <a-input v-model:value="formState.verify_password" autocomplete="off" type="password">
        </a-input>
      </a-form-item>
      <a-form-item style="text-align: right">
        <a-button @click="handleCancel">取消</a-button>
        <a-button type="primary" html-type="submit" style="margin-left: 10px"> 确认 </a-button>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { regPassword } from '@/utils/validate'
import { reactive, ref } from 'vue'

const emit = defineEmits(['update:open', 'submit'])

const props = defineProps({
  open: Boolean,
})

// 不需要父组件传值
const formState = reactive({
  password: '',
  verify_password: '',
})

// 表单
const formRef = ref()

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

// 关闭窗口
const handleCancel = () => {
  formRef.value.resetFields()
  emit('update:open', false)
}

// 提交
const onSubmit = async () => {
  emit('submit', formState)
}
</script>
