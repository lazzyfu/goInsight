<template>
  <a-modal :open="open" title="新增组织" :footer="null" @cancel="handleCancel">
    <a-form
      ref="formRef"
      :model="formState"
      :rules="rules"
      @finish="onSubmit"
    >
      <a-form-item label="组织名" name="name" has-feedback>
        <a-input v-model:value="formState.name" placeholder="请输入组织名" allow-clear />
      </a-form-item>

      <a-form-item style="text-align: right">
        <a-button @click="handleCancel">取消</a-button>
        <a-button type="primary" html-type="submit" style="margin-left: 10px">确定</a-button>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { createRootOrganizationsApi } from '@/api/admin'
import { message } from 'ant-design-vue'
import { reactive, ref } from 'vue'

const emit = defineEmits(['update:open', 'submit'])
const props = defineProps({
  open: Boolean,
})

const formRef = ref()
const formState = reactive({
  name: '',
})

const rules = {
  name: [
    { required: true, message: '请输入组织名', trigger: 'blur' },
    { min: 2, max: 32, message: '长度应在2~32个字符', trigger: 'blur' },
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
  const res = await createRootOrganizationsApi(formState).catch((err) => {})
  if (res?.code === '0000') {
    message.success('操作成功')
    emit('update:open', false)
  }
}
</script>
