<template>
  <a-modal :open="open" :footer="null" :width="560" centered class="org-modal" @cancel="handleCancel">
    <div class="modal-shell">
      <div class="modal-head">
        <span class="head-badge">子节点</span>
        <h3>新增子组织</h3>
        <p>将在当前组织下创建一个新的层级节点，用于细化部门和成员归属。</p>
      </div>

      <div class="parent-info">
        <span class="parent-label">父级组织</span>
        <a-tag color="processing" class="parent-tag">
          <template #icon>
            <FolderOutlined />
          </template>
          {{ parent_node_name }}
        </a-tag>
      </div>

      <a-form ref="formRef" :model="formState" :rules="rules" layout="vertical" class="org-form">
        <a-form-item label="组织名称" name="name">
          <a-input v-model:value="formState.name" placeholder="请输入子组织名称" :maxlength="32" show-count>
            <template #prefix>
              <FolderOutlined class="input-icon" />
            </template>
          </a-input>
        </a-form-item>
      </a-form>

      <div class="modal-footer">
        <a-button @click="handleCancel">取消</a-button>
        <a-button type="primary" :loading="uiState.loading" @click="onSubmit">创建子组织</a-button>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { createChildOrganizationsApi } from '@/api/admin'
import { FolderOutlined } from '@ant-design/icons-vue'
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { reactive, ref } from 'vue'

const emit = defineEmits(['update:open', 'submit', 'refresh'])
const props = defineProps({
  open: Boolean,
  parent_node_key: String,
  parent_node_name: String,
})

const formRef = ref()

const uiState = reactive({
  loading: false,
})

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

const onSubmit = useThrottleFn(async () => {
  try {
    await formRef.value.validateFields()
    uiState.loading = true
    const payload = {
      parent_node_key: props.parent_node_key,
      parent_node_name: props.parent_node_name,
      ...formState,
    }
    const res = await createChildOrganizationsApi(payload).catch(() => {})
    if (res) {
      message.success('创建成功')
      emit('update:open', false)
      emit('refresh')
      formRef.value?.resetFields()
    }
  } catch {
    // ignore
  } finally {
    uiState.loading = false
  }
})
</script>

<style scoped>
.modal-shell {
  padding: 8px 4px 2px;
}

.modal-head h3 {
  margin: 10px 0 6px;
  font-size: 24px;
  color: #16213c;
}

.modal-head p {
  margin: 0;
  color: #5f6b8a;
  line-height: 1.7;
}

.head-badge {
  display: inline-flex;
  align-items: center;
  font-size: 12px;
  font-weight: 700;
  color: #1f6feb;
  border-radius: 999px;
  border: 1px solid rgba(31, 111, 235, 0.24);
  background: rgba(31, 111, 235, 0.1);
  padding: 3px 10px;
}

.parent-info {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 16px;
  padding: 12px 14px;
  border-radius: 12px;
  border: 1px solid #dbe7f9;
  background: linear-gradient(90deg, #f7faff 0%, #fefeff 100%);
}

.parent-label {
  font-size: 13px;
  color: #5f6b8a;
  font-weight: 600;
}

.parent-tag {
  font-size: 13px;
  padding-inline: 8px;
}

.org-form {
  margin-top: 16px;
}

.input-icon {
  color: #8a96b0;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 8px;
}

:deep(.org-modal .ant-modal-content) {
  border-radius: 18px;
  padding: 20px 22px;
  background:
    radial-gradient(circle at right top, rgba(31, 111, 235, 0.1), rgba(31, 111, 235, 0) 50%),
    #ffffff;
}

:deep(.org-modal .ant-modal-close-x) {
  color: #7182aa;
}

:deep(.org-modal .ant-form-item-label > label) {
  color: #23335c;
  font-weight: 600;
}

:deep(.org-modal .ant-input-affix-wrapper) {
  border-radius: 10px;
  border-color: #d7e2f3;
}

:deep(.org-modal .ant-input-affix-wrapper-focused) {
  border-color: #1f6feb;
  box-shadow: 0 0 0 3px rgba(31, 111, 235, 0.14);
}
</style>
