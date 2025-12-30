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
      <a-form-item label="用户名" name="username" has-feedback>
        <a-input v-model:value="formData.username" placeholder="请输入用户名" allow-clear />
      </a-form-item>
      <a-form-item label="昵称" name="nick_name" has-feedback>
        <a-input v-model:value="formData.nick_name" placeholder="请输入昵称" allow-clear />
      </a-form-item>
      <a-form-item v-if="isCreate" label="密码" name="password" has-feedback>
        <a-input-password v-model:value="formData.password" type="password" placeholder="请输入密码" autocomplete="off" />
      </a-form-item>
      <a-form-item label="邮箱" name="email" has-feedback>
        <a-input v-model:value="formData.email" placeholder="请输入邮箱" allow-clear />
      </a-form-item>
      <a-form-item label="手机号" name="mobile" has-feedback>
        <a-input v-model:value="formData.mobile" placeholder="请输入手机号" allow-clear />
      </a-form-item>
      <a-form-item label="角色" name="role_id" has-feedback>
        <a-select
          v-model:value="formData.role_id"
          :options="props.roles"
          :field-names="{ label: 'name', value: 'id' }"
          allowClear
        />
      </a-form-item>
      <a-form-item label="激活状态" name="is_active" has-feedback>
        <a-switch v-model:checked="formData.is_active" />
      </a-form-item>
      <a-form-item label="开启2FA" name="is_two_fa" has-feedback>
        <a-switch v-model:checked="formData.is_two_fa" />
      </a-form-item>
      <a-form-item label="管理员" name="is_superuser" has-feedback>
        <a-switch v-model:checked="formData.is_superuser" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import {} from '@/api/admin'
import { regEmail, regPassword, regPhone } from '@/utils/validate'
import { computed, ref, reactive } from 'vue'

// 定义props和emits
const props = defineProps({
  open: Boolean,
  title: String,
  roles: {
    type: Array,
    default: () => [],
  },
})
const emit = defineEmits(['update:open', 'submit'])

// 表单数据
const formData = defineModel('modelValue', {
  type: Object,
  required: true,
})

// 判断是否是新增用户
const isCreate = computed(() => props.title === '新增用户')

// 表单引用
const formRef = ref()

// 状态
const uiState = reactive({
  loading: false,
})

// 表单校验规则
const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 32, message: '用户名长度为3-32个字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名仅支持字母、数字和下划线', trigger: 'blur' }
  ],
  nick_name: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 32, message: '昵称长度为2-32个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { pattern: regEmail, message: '请输入合法的邮箱地址', trigger: 'blur' }
  ],
  mobile: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: regPhone, message: '请输入合法的手机号码', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    {
      validator: (_, value) => {
        // 编辑用户不必填密码
        if (!isCreate.value && !value) return Promise.resolve()
        if (!value) return Promise.reject(new Error('请输入密码'))
        if (value.length < 8) return Promise.reject(new Error('密码长度至少8个字符'))
        if (!regPassword.test(value))
          return Promise.reject(new Error('密码必须包含大小写字母、数字和特殊字符'))
        return Promise.resolve()
      },
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
