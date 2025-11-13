<template>
  <a-modal :open="open" :title="title" :footer="null" @cancel="handleCancel">
    <a-form
      ref="formRef"
      :label-col="{ span: 4 }"
      :wrapper-col="{ span: 18 }"
      :model="localFormState"
      :rules="rules"
      @finish="onSubmit"
    >
      <a-form-item label="用户名" name="username" has-feedback>
        <a-input v-model:value="localFormState.username" placeholder="请输入用户名" allow-clear />
      </a-form-item>

      <a-form-item label="昵称" name="nick_name" has-feedback>
        <a-input v-model:value="localFormState.nick_name" placeholder="请输入昵称" allow-clear />
      </a-form-item>

      <a-form-item v-if="title === '新增用户'" label="密码" name="password" has-feedback>
        <a-input v-model:value="localFormState.password" type="password" placeholder="请输入密码" />
      </a-form-item>

      <a-form-item label="邮箱" name="email" has-feedback>
        <a-input v-model:value="localFormState.email" placeholder="请输入邮箱" allow-clear />
      </a-form-item>

      <a-form-item label="手机号" name="mobile" has-feedback>
        <a-input v-model:value="localFormState.mobile" placeholder="请输入手机号" allow-clear />
      </a-form-item>

      <a-form-item label="角色" name="role" has-feedback>
        <a-select
          v-model:value="localFormState.role"
          :options="roles"
          :field-names="{ label: 'name', value: 'id' }"
          allowClear
        />
      </a-form-item>

      <a-form-item label="激活状态" name="is_active" has-feedback>
        <a-switch v-model:checked="localFormState.is_active" />
      </a-form-item>

      <a-form-item label="开启2FA" name="is_two_fa" has-feedback>
        <a-switch v-model:checked="localFormState.is_two_fa" />
      </a-form-item>

      <a-form-item label="管理员" name="is_superuser" has-feedback>
        <a-switch v-model:checked="localFormState.is_superuser" />
      </a-form-item>

      <a-form-item style="text-align: right">
        <a-button @click="handleCancel">取消</a-button>
        <a-button type="primary" html-type="submit" style="margin-left: 10px">确定</a-button>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { getRolesApi } from '@/api/admin'
import { regEmail, regPassword, regPhone } from '@/utils/validate'
import { onMounted, reactive, ref, watch } from 'vue'

const emit = defineEmits(['update:open', 'submit'])
const props = defineProps({
  open: Boolean,
  title: String,
  formState: Object,
})

// formState父组件传值，子组件修改，需要重新赋值
const localFormState = reactive({ ...props.formState })

watch(
  () => props.formState,
  (newVal) => {
    Object.assign(localFormState, newVal)
  },
  { immediate: true, deep: true },
)

const roles = ref([])
const formRef = ref()

const validatePassword = (_r, v) => {
  if (!v) return Promise.reject('请输入密码')
  if (!regPassword.test(v)) return Promise.reject('密码至少7个字符，包含大小写字母、数字和特殊字符')
  return Promise.resolve()
}

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  nick_name: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
  email: [{ required: true, pattern: regEmail, message: '请输入合法邮箱', trigger: 'blur' }],
  mobile: [{ required: true, pattern: regPhone, message: '请输入合法手机号', trigger: 'blur' }],
  password: [{ required: true, validator: validatePassword, trigger: 'blur' }],
}

const handleCancel = () => {
  emit('update:open', false)
  formRef.value?.resetFields()
}

const onSubmit = () => {
  emit('submit', localFormState)
}

onMounted(async () => {
  const res = await getRolesApi().catch(() => {})
  roles.value = res?.data || []
})
</script>
