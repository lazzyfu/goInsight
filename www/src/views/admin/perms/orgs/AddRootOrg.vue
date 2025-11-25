<template>
  <a-modal :open="open" title="新增根组织" :width="480" centered @cancel="handleCancel">
    <template #footer>
      <a-button @click="handleCancel">取消</a-button>
      <a-button type="primary" :loading="loading" @click="onSubmit">确定</a-button>
    </template>
    <div class="modal-content">
      <div class="modal-icon">
        <ApartmentOutlined />
      </div>
      <p class="modal-desc">创建一个新的根级组织，作为组织架构的顶层节点</p>
      <a-form ref="formRef" :model="formState" :rules="rules" layout="vertical" class="org-form">
        <a-form-item label="组织名称" name="name">
          <a-input
            v-model:value="formState.name"
            placeholder="请输入组织名称"
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
import { createRootOrganizationsApi } from '@/api/admin'
import { ApartmentOutlined, FolderOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { reactive, ref } from 'vue'

const emit = defineEmits(['update:open', 'submit', 'refresh'])
const props = defineProps({
  open: Boolean,
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
    const res = await createRootOrganizationsApi(formState).catch(() => {})
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
  text-align: center;
  padding: 16px 0;
}

.modal-icon {
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #e6f7ff 0%, #bae7ff 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
  font-size: 28px;
  color: #1890ff;
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
