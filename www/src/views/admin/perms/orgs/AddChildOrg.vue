<template>
  <a-modal :open="open" title="新增子组织" :width="480" centered @cancel="handleCancel">
    <template #footer>
      <a-button @click="handleCancel">取消</a-button>
      <a-button type="primary" :loading="uiState.loading" @click="onSubmit">确定</a-button>
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

// 定义props和emits
const emit = defineEmits(['update:open', 'submit', 'refresh'])
const props = defineProps({
  open: Boolean,
  parent_node_key: String,
  parent_node_name: String,
})

// 表单引用
const formRef = ref()

// 状态
const uiState = reactive({
  loading: false,
})

// 表单数据
const formState = reactive({
  name: '',
})

// 表单校验规则
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
    }
  } catch (err) {
  } finally {
    uiState.loading = false
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
