<template>
  <a-modal :open="open" :footer="null" :width="620" centered class="org-modal" @cancel="handleCancel">
    <div class="modal-shell">
      <div class="modal-head">
        <span class="head-badge">成员绑定</span>
        <h3>绑定组织成员</h3>
        <p>将用户加入当前组织并授予角色，后续可在成员列表中随时移除。</p>
      </div>

      <div class="node-info" v-if="props.nodeName">
        <span>目标组织</span>
        <strong>{{ props.nodeName }}</strong>
      </div>

      <a-form ref="formRef" :model="formData" :rules="rules" layout="vertical" class="bind-form">
        <a-form-item label="选择用户" name="users" required>
          <a-select
            v-model:value="formData.users"
            mode="multiple"
            placeholder="请选择用户"
            :options="props.users"
            :field-names="{ label: 'username', value: 'uid', children: 'children' }"
            show-search
            allow-clear
            class="modal-field-full-width"
          />
        </a-form-item>

        <a-form-item label="选择角色" name="roles" required>
          <a-select
            v-model:value="formData.roles"
            placeholder="请选择角色"
            :options="props.roles"
            :field-names="{ label: 'name', value: 'id', children: 'children' }"
            show-search
            allow-clear
            class="modal-field-full-width"
          />
        </a-form-item>
      </a-form>

      <div class="selected-info" v-if="formData.users.length > 0">
        <InfoCircleOutlined />
        已选择 {{ formData.users.length }} 个用户
      </div>

      <div class="modal-footer">
        <a-button @click="handleCancel">取消</a-button>
        <a-button type="primary" :loading="uiState.loading" @click="onSubmit">确认绑定</a-button>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { InfoCircleOutlined } from '@ant-design/icons-vue'
import { reactive, ref } from 'vue'

const emit = defineEmits(['update:open', 'submit'])
const props = defineProps({
  open: Boolean,
  nodeKey: String,
  nodeName: {
    type: String,
    default: '',
  },
  users: {
    type: Array,
    default: () => [],
  },
  roles: {
    type: Array,
    default: () => [],
  },
})

const formRef = ref()

const uiState = reactive({
  loading: false,
})

const rules = {
  users: [{ required: true, type: 'array', min: 1, message: '请至少选择一个用户', trigger: 'change' }],
  roles: [{ required: true, message: '请选择角色', trigger: 'change' }],
}

const formData = defineModel('modelValue', {
  type: Object,
  required: true,
})

const handleCancel = () => {
  emit('update:open', false)
  formRef.value?.resetFields()
}

const onSubmit = async () => {
  try {
    await formRef.value.validateFields()
    uiState.loading = true
    const payload = {
      key: props.nodeKey,
      ...formData.value,
    }
    emit('submit', payload)
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
  background: rgba(14, 159, 110, 0.11);
  padding: 3px 10px;
}

.node-info {
  margin-top: 14px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  border-left: 3px solid rgba(14, 159, 110, 0.6);
  background: #f2fbf7;
  border-radius: 8px;
  padding: 10px 12px;
}

.node-info span {
  color: #517868;
  font-size: 12px;
}

.node-info strong {
  color: #1e4235;
  font-size: 14px;
}

.bind-form {
  margin-top: 16px;
}

.selected-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 11px 12px;
  background: #f2fbf7;
  border-radius: 10px;
  color: #128555;
  border: 1px solid #bfead7;
  font-size: 13px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 14px;
}

:deep(.org-modal .ant-modal-content) {
  border-radius: 18px;
  padding: 20px 22px;
  background:
    radial-gradient(circle at right top, rgba(14, 159, 110, 0.08), rgba(14, 159, 110, 0) 54%),
    #ffffff;
}

:deep(.org-modal .ant-modal-close-x) {
  color: #7182aa;
}

:deep(.org-modal .ant-form-item-label > label) {
  color: #23335c;
  font-weight: 600;
}

:deep(.org-modal .ant-select-selector) {
  border-radius: 10px !important;
  border-color: #d7e2f3 !important;
}

:deep(.org-modal .ant-select-focused .ant-select-selector) {
  border-color: #1f6feb !important;
  box-shadow: 0 0 0 3px rgba(31, 111, 235, 0.14);
}

:deep(.org-modal .ant-select-selection-item) {
  background: #eef5ff;
  border-color: #c3d8ff;
}
</style>
