<template>
  <a-modal :open="open" :footer="null" :width="560" centered class="org-modal" @cancel="handleCancel">
    <div class="modal-shell">
      <div class="modal-head">
        <span class="head-badge">根节点</span>
        <h3>新增根组织</h3>
        <p>创建组织树的顶层入口，后续所有部门节点都将从这里扩展。</p>
      </div>

      <div class="modal-body">
        <div class="modal-icon">
          <ApartmentOutlined />
        </div>

        <a-form ref="formRef" :model="formState" :rules="rules" layout="vertical" class="org-form">
          <a-form-item label="组织名称" name="name">
            <a-input v-model:value="formState.name" placeholder="请输入组织名称" :maxlength="32" show-count>
              <template #prefix>
                <FolderOutlined class="input-icon" />
              </template>
            </a-input>
          </a-form-item>
        </a-form>
      </div>

      <div class="modal-footer">
        <a-button @click="handleCancel">取消</a-button>
        <a-button type="primary" :loading="uiState.loading" @click="onSubmit">创建组织</a-button>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { createRootOrganizationsApi } from '@/api/admin'
import { ApartmentOutlined, FolderOutlined } from '@ant-design/icons-vue'
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { reactive, ref } from 'vue'

const emit = defineEmits(['update:open', 'submit', 'refresh'])
defineProps({
  open: Boolean,
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
    const res = await createRootOrganizationsApi(formState).catch(() => {})
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

.modal-body {
  margin-top: 18px;
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
  color: #1f6feb;
  font-size: 26px;
  background: linear-gradient(140deg, rgba(31, 111, 235, 0.24), rgba(31, 111, 235, 0.06));
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
