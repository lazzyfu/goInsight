<template>
  <a-modal :open="open" title="绑定用户" :width="520" centered @cancel="handleCancel">
    <template #footer>
      <a-button @click="handleCancel">取消</a-button>
      <a-button type="primary" :loading="loading" @click="onSubmit">确定</a-button>
    </template>
    <div class="modal-content">
      <div class="modal-icon">
        <UserAddOutlined />
      </div>
      <p class="modal-desc">选择要添加到当前组织的用户</p>
      <a-form ref="formRef" :model="formState" layout="vertical" class="bind-form">
        <a-form-item label="选择用户" name="users">
          <a-select
            v-model:value="formState.users"
            mode="multiple"
            placeholder="请选择用户"
            :options="userOptions"
            show-search
            :filter-option="filterOption"
            allow-clear
            style="width: 100%"
          >
            <template #option="{ value, label }">
              <div class="user-option">
                <a-avatar :size="24" class="option-avatar">
                  {{ label?.charAt(0)?.toUpperCase() }}
                </a-avatar>
                <span>{{ label }}</span>
              </div>
            </template>
          </a-select>
        </a-form-item>
      </a-form>
      <div class="selected-info" v-if="formState.users.length > 0">
        <InfoCircleOutlined />
        已选择 {{ formState.users.length }} 个用户
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { InfoCircleOutlined, UserAddOutlined } from '@ant-design/icons-vue'
import { computed, reactive, ref } from 'vue'

const emit = defineEmits(['update:open', 'submit'])
const props = defineProps({
  open: Boolean,
  nodeKey: String,
  users: Array,
})

const formRef = ref()
const loading = ref(false)

const formState = reactive({
  users: [],
})

const userOptions = computed(() => {
  return (props.users || []).map((user) => ({
    value: user.uid || user.id,
    label: user.username || user.name,
  }))
})

const filterOption = (input, option) => {
  return option.label.toLowerCase().includes(input.toLowerCase())
}

const handleCancel = () => {
  emit('update:open', false)
  formRef.value?.resetFields()
  formState.users = []
}

const onSubmit = async () => {
  if (formState.users.length === 0) {
    return
  }
  loading.value = true
  const payload = {
    key: props.nodeKey,
    ...formState,
  }
  emit('submit', payload)
  loading.value = false
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
  background: linear-gradient(135deg, #f6ffed 0%, #b7eb8f 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
  font-size: 28px;
  color: #52c41a;
}

.modal-desc {
  color: #8c8c8c;
  font-size: 14px;
  margin-bottom: 24px;
}

.bind-form {
  text-align: left;
}

.user-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.option-avatar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  font-size: 12px;
}

.selected-info {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px;
  background: #f6ffed;
  border-radius: 8px;
  color: #52c41a;
  font-size: 14px;
  margin-top: 16px;
}

:deep(.ant-form-item-label > label) {
  font-weight: 500;
}

:deep(.ant-select-selection-item) {
  background: #f0f5ff;
  border-color: #adc6ff;
}
</style>
