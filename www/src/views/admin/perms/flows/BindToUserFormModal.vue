<template>
  <a-modal
    :open="props.open"
    :footer="null"
    :width="600"
    centered
    class="flow-bind-modal"
    @cancel="handleCancel"
  >
    <div class="modal-shell">
      <div class="modal-head">
        <span class="head-badge">用户绑定</span>
        <h3>绑定审批流到用户</h3>
        <p>选择一个审批流，并把该流程批量分配给尚未绑定的用户。</p>
      </div>

      <a-form ref="formRef" :model="formState" :rules="rules" layout="vertical" class="bind-form">
        <a-form-item label="审批流" name="approval_id" required>
          <a-select
            v-model:value="formState.approval_id"
            placeholder="请选择要分配的审批流"
            :options="props.flowOptions"
            show-search
            :filter-option="filterOption"
            allow-clear
          />
          <template #extra>
            <p>请选择一个已经定义好的审批流程。</p>
          </template>
        </a-form-item>

        <a-form-item label="选择未绑定的用户" name="users" required>
          <a-select
            v-model:value="formState.users"
            mode="multiple"
            placeholder="请选择要绑定的用户"
            :options="props.userOptions"
            show-search
            :filter-option="filterOption"
            allow-clear
          />
          <template #extra>
            <p>这些用户发起审批时，将默认使用以上流程。</p>
          </template>
        </a-form-item>
      </a-form>

      <div class="modal-footer">
        <a-button @click="handleCancel">取消</a-button>
        <a-button type="primary" :loading="uiState.loading" @click="onSubmit">确认绑定</a-button>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { reactive, ref } from 'vue'

const props = defineProps({
  open: Boolean,
  flowOptions: { type: Array, default: () => [] },
  userOptions: { type: Array, default: () => [] },
})
const emit = defineEmits(['update:open', 'submit'])

const formRef = ref()

const uiState = reactive({
  loading: false,
})

const formState = reactive({
  approval_id: undefined,
  users: [],
})

const rules = {
  approval_id: [{ required: true, message: '请选择审批流', trigger: 'change' }],
  users: [
    {
      required: true,
      type: 'array',
      validator: async (rule, value) => {
        if (!value || value.length === 0) {
          throw new Error('请选择至少一个用户进行绑定')
        }
      },
      trigger: 'change',
    },
  ],
}

const filterOption = (input, option) => {
  const label = (option?.label || '').toString().toLowerCase()
  const value = (option?.value || '').toString().toLowerCase()
  return label.includes(input.toLowerCase()) || value.includes(input.toLowerCase())
}

const handleCancel = () => {
  emit('update:open', false)
  formRef.value?.resetFields()
}

const onSubmit = async () => {
  try {
    await formRef.value.validate()
    uiState.loading = true

    emit('submit', {
      users: formState.users,
      approval_id: formState.approval_id,
    })
  } catch {
    // ignore
  } finally {
    uiState.loading = false
  }
}
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
  color: #0f8a54;
  border-radius: 999px;
  border: 1px solid rgba(14, 159, 110, 0.35);
  background: rgba(14, 159, 110, 0.12);
  padding: 3px 10px;
}

.bind-form {
  margin-top: 16px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 6px;
}

:deep(.flow-bind-modal .ant-modal-content) {
  border-radius: 18px;
  padding: 20px 22px;
  background:
    radial-gradient(circle at right top, rgba(14, 159, 110, 0.08), rgba(14, 159, 110, 0) 55%),
    #ffffff;
}

:deep(.flow-bind-modal .ant-form-item-label > label) {
  color: #23335c;
  font-weight: 600;
}

:deep(.flow-bind-modal .ant-form-item-extra) {
  color: #5f6b8a;
}

:deep(.flow-bind-modal .ant-select-selector) {
  border-radius: 10px !important;
  border-color: #d7e2f3 !important;
}

:deep(.flow-bind-modal .ant-select-focused .ant-select-selector) {
  border-color: #1f6feb !important;
  box-shadow: 0 0 0 3px rgba(31, 111, 235, 0.14);
}
</style>
