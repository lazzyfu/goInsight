<template>
  <a-modal :open="open" title="编辑节点名" :footer="null" @cancel="handleCancel">
    <a-form
      ref="formRef"
      :label-col="{ span: 4 }"
      :wrapper-col="{ span: 18 }"
      :model="formState"
      :rules="rules"
      @finish="onSubmit"
    >
      <a-form-item label="组织名" name="name" has-feedback>
        <a-input v-model:value="formState.name" placeholder="请输入新的组织名" allow-clear />
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
import { updateOrganizationsApi } from '@/api/admin'
import { message } from 'ant-design-vue'
import { reactive, ref } from 'vue'

const emit = defineEmits(['update:open', 'submit'])
const props = defineProps({
  open: Boolean,
  nodeKey: String,
})

const formRef = ref()
const formState = reactive({
  name: '',
})

const rules = {
  name: [
    { required: true, message: '请输入新的组织名', trigger: 'blur' },
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
  const payload = {
    key: props.nodeKey,
    ...formState,
  }
  const res = await updateOrganizationsApi(payload).catch(() => {})
  if (res?.code === '0000') {
    message.success('操作成功')
    emit('update:open', false)
  }
}
</script>
