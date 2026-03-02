<template>
  <a-modal :open="open" :footer="null" :width="560" centered class="org-modal" @cancel="handleCancel">
    <div class="modal-shell">
      <div class="modal-head">
        <span class="head-badge">重命名</span>
        <h3>编辑组织名称</h3>
        <p>更新组织展示名称，不会影响节点层级关系和已绑定成员。</p>
      </div>

      <div class="node-info" v-if="nodeName">
        <span>当前名称</span>
        <strong>{{ nodeName }}</strong>
      </div>

      <div class="modal-body">
        <div class="modal-icon">
          <EditOutlined />
        </div>
        <a-form ref="formRef" :model="formState" :rules="rules" layout="vertical" class="org-form">
          <a-form-item label="新的组织名称" name="name">
            <a-input v-model:value="formState.name" placeholder="请输入新的组织名称" :maxlength="32" show-count>
              <template #prefix>
                <FolderOutlined class="input-icon" />
              </template>
            </a-input>
          </a-form-item>
        </a-form>
      </div>

      <div class="modal-footer">
        <a-button @click="handleCancel">取消</a-button>
        <a-button type="primary" :loading="uiState.loading" @click="onSubmit">保存修改</a-button>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { updateOrganizationsApi } from '@/api/admin'
import { EditOutlined, FolderOutlined } from '@ant-design/icons-vue'
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { reactive, ref } from 'vue'

const props = defineProps({
  open: Boolean,
  nodeKey: String,
  nodeName: {
    type: String,
    default: '',
  },
})
const emit = defineEmits(['update:open', 'submit', 'refresh'])

const formRef = ref()

const uiState = reactive({
  loading: false,
})

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

const onSubmit = useThrottleFn(async () => {
  try {
    await formRef.value.validateFields()
    uiState.loading = true
    const payload = {
      key: props.nodeKey,
      ...formState,
    }
    const res = await updateOrganizationsApi(payload).catch(() => {})
    if (res) {
      message.success('更新成功')
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
  color: #b26800;
  border-radius: 999px;
  border: 1px solid rgba(250, 140, 22, 0.4);
  background: rgba(250, 140, 22, 0.12);
  padding: 3px 10px;
}

.node-info {
  margin-top: 16px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  border-left: 3px solid rgba(250, 140, 22, 0.65);
  background: #fffaf2;
  border-radius: 8px;
  padding: 10px 12px;
}

.node-info span {
  color: #7f6750;
  font-size: 12px;
}

.node-info strong {
  color: #412f1b;
  font-size: 14px;
}

.modal-body {
  margin-top: 16px;
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.modal-icon {
  width: 54px;
  height: 54px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: #c97802;
  font-size: 24px;
  background: linear-gradient(140deg, rgba(250, 140, 22, 0.24), rgba(250, 140, 22, 0.06));
}

.org-form {
  flex: 1;
}

.input-icon {
  color: #8a96b0;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 10px;
}

:deep(.org-modal .ant-modal-content) {
  border-radius: 18px;
  padding: 20px 22px;
  background:
    radial-gradient(circle at right top, rgba(250, 140, 22, 0.08), rgba(250, 140, 22, 0) 54%),
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
