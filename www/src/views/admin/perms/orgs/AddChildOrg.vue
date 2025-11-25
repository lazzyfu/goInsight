<template>
  <a-modal :open="open" title="新增子组织" :width="480" centered @cancel="handleCancel">
    <template #footer>
      <a-button @click="handleCancel">取消</a-button>
      <a-button type="primary" :loading="loading" @click="onSubmit">确定</a-button>
    </template>
    <div class="modal-content">
      <div class="parent-info">
        <span class="parent-label">父级组织</span>
        <a-tag color="blue" class="parent-tag">
          <template #icon><FolderOutlined /></template>
          {{ parent_node_name }}
        </a-tag>
      </div>
      <a-form ref="formRef" :model="formState" :rules="rules" layout="vertical" class="org-form">
        <a-form-item label="组织名称" name="name">
          <a-input
            v-model:value="formState.name"
            placeholder="请输入子组织名称"
            :maxlength="32"
            show-count
          >
            <template #prefix>
              <FolderOutlined style="color: #bfbfbf" />
            </template>
          </a-input>
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { createChildOrganizationsApi } from '@/api/admin'
import { FolderOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { reactive, ref } from 'vue'

const emit = defineEmits(['update:open', 'submit', 'refresh'])
const props = defineProps({
  open: Boolean,
  parent_node_key: String,
  parent_node_name: String,
})

const formRef = ref()
const loading = ref(false)

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
  try {
    await formRef.value.validateFields()
    loading.value = true
    const payload = {
      parent_node_key: props.parent_node_key,
      parent_node_name: props.parent_node_name,
      ...formState,
    }
    const res = await createChildOrganizationsApi(payload).catch(() => {})
    if (res?.code === '0000') {
      message.success('创建成功')
      emit('update:open', false)
      emit('refresh')
      formRef.value?.resetFields()
    }
  } catch (err) {
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.modal-content {
  padding: 8px 0;
}

.parent-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
  margin-bottom: 24px;
}

.parent-label {
  font-size: 14px;
  color: #8c8c8c;
}

.parent-tag {
  font-size: 14px;
}

.org-form {
  text-align: left;
}

:deep(.ant-form-item-label > label) {
  font-weight: 500;
}
</style>
