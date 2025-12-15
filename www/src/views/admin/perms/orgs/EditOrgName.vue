<template>
  <a-modal :open="open" title="编辑组织名称" :width="480" centered @cancel="handleCancel">
    <template #footer>
      <a-button @click="handleCancel">取消</a-button>
      <a-button type="primary" :loading="uiState.loading" @click="onSubmit">确定</a-button>
    </template>
    <div class="modal-content">
      <div class="modal-icon edit">
        <EditOutlined />
      </div>
      <p class="modal-desc">修改当前组织的名称</p>
      <a-form ref="formRef" :model="formState" :rules="rules" layout="vertical" class="org-form">
        <a-form-item label="新的组织名称" name="name">
          <a-input
            v-model:value="formState.name"
            placeholder="请输入新的组织名称"
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
import { updateOrganizationsApi } from '@/api/admin'
import { EditOutlined, FolderOutlined } from '@ant-design/icons-vue'
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { reactive, ref } from 'vue'

// 定义props和emits
const props = defineProps({
  open: Boolean,
  nodeKey: String,
})
const emit = defineEmits(['update:open', 'submit', 'refresh'])

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
    { required: true, message: '请输入新的组织名', trigger: 'blur' },
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
    }
  } catch (err) {
  } finally {
    uiState.loading = false
  }
})
</script>

<style scoped>
.modal-content {
  text-align: center;
  padding: 16px 0;
}

.modal-icon {
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #fff7e6 0%, #ffe7ba 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
  font-size: 28px;
  color: #fa8c16;
}

.modal-desc {
  color: #8c8c8c;
  font-size: 14px;
  margin-bottom: 24px;
}

.org-form {
  text-align: left;
}

:deep(.ant-form-item-label > label) {
  font-weight: 500;
}
</style>
