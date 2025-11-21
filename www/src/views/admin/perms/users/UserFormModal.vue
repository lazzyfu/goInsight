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
      <a-form-item label="用户名" name="username" has-feedback>
        <a-input v-model:value="formData.username" placeholder="请输入用户名" allow-clear />
      </a-form-item>
      <a-form-item label="昵称" name="nick_name" has-feedback>
        <a-input v-model:value="formData.nick_name" placeholder="请输入昵称" allow-clear />
      </a-form-item>
      <a-form-item v-if="isCreate" label="密码" name="password" has-feedback>
        <a-input v-model:value="formData.password" type="password" placeholder="请输入密码" />
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
import {} from '@/api/admin'
import { regEmail, regPassword, regPhone } from '@/utils/validate'
import { computed, ref } from 'vue'

const props = defineProps({
  open: Boolean,
  title: String,
  roles: {
    type: Array,
    default: () => [],
  },
})
const emit = defineEmits(['update:open', 'submit'])

// 使用defineModel接收 v-model:modelValue
// 它自动创建了一个名为modelValue的prop，并提供了一个value来读取，以及一个自动触发update:modelValue的setter
const formData = defineModel('modelValue', {
  type: Object,
  required: true,
})

// 判断是否是新增用户
const isCreate = computed(() => props.title === '新增用户')

// 表单引用
const formRef = ref()

const rules = {
  username: [{ required: true, min: 2, max: 32, message: '请输入用户名', trigger: 'blur' }],
  nick_name: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
  email: [{ required: true, pattern: regEmail, message: '请输入合法邮箱', trigger: 'blur' }],
  mobile: [{ required: true, pattern: regPhone, message: '请输入合法手机号', trigger: 'blur' }],
  password: [
    { required: true },
    {
      validator: (_, value) => {
        // 编辑用户不必填密码
        if (!isCreate.value && !value) return Promise.resolve()
        if (!value) return Promise.reject(new Error('请输入密码'))
        if (!regPassword.test(value))
          return Promise.reject(new Error('密码至少7个字符，包含大小写字母、数字和特殊字符'))
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
  emit('submit', formData.value)
}
</script>
