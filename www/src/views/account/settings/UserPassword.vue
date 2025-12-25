<template>
  <a-modal title="修改密码" v-model:open="modalOpen" width="50%" :footer="null" @cancel="handleCancel">
    <a-form ref="formRef" :model="formState" :rules="rules" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }"
      autocomplete="off" @finish="onSubmit">
      <a-form-item label="当前密码" has-feedback name="old_password">
        <a-input-password v-model:value="formState.old_password" type="password" autocomplete="off">
          <LockOutlined />
        </a-input-password>
      </a-form-item>
      <a-form-item label="新密码" has-feedback name="new_password">
        <a-input-password v-model:value="formState.new_password" autocomplete="off" type="password">
        </a-input-password>
      </a-form-item>
      <a-form-item label="确认密码" has-feedback name="confirm_password">
        <a-input-password v-model:value="formState.confirm_password" autocomplete="off" type="password">
        </a-input-password>
      </a-form-item>
      <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
        <a-button type="primary" html-type="submit">提交</a-button>
        <a-button style="margin-left: 10px" @click="resetForm">重置</a-button>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { ChangePasswordApi } from '@/api/profile'
import { regPassword } from '@/utils/validate'
import { LockOutlined } from '@ant-design/icons-vue'
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { computed, reactive, ref } from 'vue'

const props = defineProps({
  open: Boolean,
  title: String,
})
const emit = defineEmits(['update:open'])

// 用 computed 包一层，避免直接修改 props
const modalOpen = computed({
  get: () => props.open,
  set: (val) => emit('update:open', val),
})

const formRef = ref()
const formState = reactive({
  old_password: '',
  new_password: '',
  confirm_password: '',
})

const validateOldPass = async (_rule, value) => {
  if (!value) return Promise.reject('请输入密码')
  if (value.length < 1) return Promise.reject('至少1个字符')
}

const validateNewPass = async (_rule, value) => {
  if (!value) return Promise.reject('请输入密码')
  if (!regPassword.test(value))
    return Promise.reject('密码至少7个字符,必须包含大写字母、小写字母、数字和特殊字符')
  if (formState.confirm_password) formRef.value.validateFields('confirm_password')
  return Promise.resolve()
}

const validateVerifyPass = async (_rule, value) => {
  if (!value) return Promise.reject('请输入密码')
  if (!regPassword.test(value))
    return Promise.reject('密码至少7个字符,必须包含大写字母、小写字母、数字和特殊字符')
  if (value !== formState.new_password) return Promise.reject('两次输入的密码不一致')
  return Promise.resolve()
}

const rules = {
  old_password: [
    {
      required: true,
      validator: validateOldPass,
      trigger: 'change',
    },
  ],
  new_password: [
    {
      required: true,
      validator: validateNewPass,
      trigger: 'change',
    },
  ],
  confirm_password: [
    {
      required: true,
      validator: validateVerifyPass,
      trigger: 'change',
    },
  ],
}

// 取消按钮处理函数
const handleCancel = () => {
  modalOpen.value = false
  formRef.value?.resetFields()
}

// 提交表单
const onSubmit = useThrottleFn(async (values) => {
  const res = await ChangePasswordApi(values).catch(() => { })
  if (res) {
    message.info("密码修改成功")
    handleCancel()
    location.reload() // 刷新页面，此时token已经过期，需要重新登录
  }
})

const resetForm = () => {
  formRef.value.resetFields()
}
</script>
