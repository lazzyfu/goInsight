<template>
  <a-modal :open="props.open" :title="props.title" width="45%" @cancel="handleCancel">
    <template #footer>
      <a-button @click="handleCancel">取消</a-button>
      <a-button type="primary" :loading="uiState.loading" @click="onSubmit"> 确定 </a-button>
    </template>
    <a-form ref="formRef" layout="vertical" :model="formData" style="margin-top: 24px">
      <a-form-item
        label="审批流名称"
        name="name"
        :rules="[
          {
            required: true,
            message: '请输入审批流名称',
          }
        ]"
        has-feedback
      >
        <a-input v-model:value="formData.name" placeholder="请输入审批流名称" allow-clear />
      </a-form-item>

      <div class="divider-section">
        <div class="divider-title">
          <span class="divider-icon">📋</span>
          <span>审批阶段定义</span>
          <a-badge
            :count="formData.definition ? formData.definition.length : 0"
            :number-style="{ backgroundColor: '#1890ff' }"
            style="margin-left: 12px"
          />
        </div>
      </div>

      <a-form-item
        label=""
        :label-col="{ span: 0 }"
        :wrapper-col="{ span: 24 }"
        class="dynamic-definition-item"
      >
        <div class="stage-list">
          <div v-for="(stage, index) in formData.definition" :key="index" class="stage-item">
            <div class="stage-header">
              <div class="stage-number">{{ index + 1 }}</div>
              <a-input
                v-model:value="stage.stage_name"
                placeholder="阶段名称 (如：部门经理审批)"
                class="stage-name-input"
              />
              <a-button
                v-if="formData.definition.length > 1"
                type="text"
                danger
                @click="removeStage(index)"
                class="delete-btn"
              >
                <DeleteOutlined />
              </a-button>
            </div>

            <div class="stage-content">
              <a-row :gutter="16">
                <a-col :span="8">
                  <a-form-item
                    label="审批类型"
                    :rules="[
                      {
                        required: true,
                        message: '请选择审批类型',
                      }
                    ]"
                    class="form-field"
                  >
                    <a-select v-model:value="stage.type" style="width: 100%">
                      <a-select-option value="AND">
                        <span class="option-text">🤝 AND (会签)</span>
                      </a-select-option>
                      <a-select-option value="OR">
                        <span class="option-text">✅ OR (或签)</span>
                      </a-select-option>
                    </a-select>
                  </a-form-item>
                </a-col>
                <a-col :span="16">
                  <a-form-item
                    label="审批人"
                    :rules="[
                      {
                        required: true,
                        message: '请选择审批人',
                      }
                    ]"
                    class="form-field"
                  >
                    <a-select
                      v-model:value="stage.approvers"
                      mode="multiple"
                      show-search
                      :filter-option="filterUserOption"
                      style="width: 100%"
                      placeholder="请选择审批人"
                      :options="props.userOptions"
                      option-label-prop="label"
                      :max-tag-count="3"
                    />
                  </a-form-item>
                </a-col>
              </a-row>
            </div>
          </div>
        </div>

        <a-button type="dashed" block class="add-stage-btn" @click="addStage">
          <PlusOutlined /> 增加审批阶段
        </a-button>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { ref, reactive } from 'vue'
import { message } from 'ant-design-vue'

// props
const props = defineProps({
  open: Boolean,
  title: String,
  userOptions: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:open', 'submit'])

// 表单数据 (保持你现有的 defineModel 用法)
const formData = defineModel('formData', {
  type: Object,
  required: true,
})

// 表单引用
const formRef = ref()

// 状态
const uiState = reactive({
  loading: false,
})

// 过滤器：增加健壮性，防止 option 或 label 未定义时报错
const filterUserOption = (input, option) => {
  if (!option) return false
  const label = (option.label || '').toString()
  const value = (option.value || '').toString()
  return label.toLowerCase().includes(input.toLowerCase()) || value.toLowerCase().includes(input.toLowerCase())
}

const addStage = () => {
  const definition = formData.value.definition || []
  const newStage = {
    stage: definition.length + 1,
    stage_name: `新阶段 ${definition.length + 1}`,
    approvers: [],
    type: 'AND',
  }
  definition.push(newStage)
}

const removeStage = (index) => {
  const definition = formData.value.definition
  definition.splice(index, 1)
  definition.forEach((stage, i) => {
    stage.stage = i + 1
  })
}

// 取消按钮
const handleCancel = () => {
  emit('update:open', false)
  // 如果表单实例存在则重置（保持兼容）
  formRef.value?.resetFields?.()
}

// 自定义校验（对动态数组字段使用自定义校验更可控）
const validateDefinition = async () => {
  const definition = formData.value.definition
  if (!definition || definition.length === 0) {
    return Promise.reject('请至少配置一个审批阶段。')
  }
  for (const [index, stage] of definition.entries()) {
    if (!stage.stage_name || stage.stage_name.trim() === '') {
      return Promise.reject(`阶段 ${index + 1}：阶段名称不能为空。`)
    }
    if (!stage.approvers || stage.approvers.length === 0) {
      return Promise.reject(`阶段 ${index + 1}：至少需要指定一个审批人。`)
    }
  }
  return true
}

// 提交表单
const onSubmit = async () => {
  try {
    await validateDefinition()
    uiState.loading = true
    emit('submit', formData.value)
  } catch (err) {
    message.error(err)
  } finally {
    uiState.loading = false
  }
}
</script>

<style scoped>
.dynamic-definition-item :deep(.ant-form-item-control-input-content) {
  padding: 0;
}

.divider-section {
  margin: 32px 0 24px;
  padding: 0 0 16px;
  border-bottom: 2px solid #f0f0f0;
}

.divider-title {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
  color: #262626;
}

.divider-icon {
  font-size: 20px;
  margin-right: 8px;
}

.stage-list {
  max-height: 500px;
  overflow-y: auto;
  padding: 4px;
}

.stage-list::-webkit-scrollbar {
  width: 6px;
}

.stage-list::-webkit-scrollbar-track {
  background: #f5f5f5;
  border-radius: 3px;
}

.stage-list::-webkit-scrollbar-thumb {
  background: #d9d9d9;
  border-radius: 3px;
}

.stage-list::-webkit-scrollbar-thumb:hover {
  background: #bfbfbf;
}

.stage-item {
  border: 2px solid #f0f0f0;
  border-radius: 12px;
  padding: 0;
  margin-bottom: 16px;
  background: linear-gradient(135deg, #ffffff 0%, #fafafa 100%);
  transition: all 0.3s ease;
  overflow: hidden;
}

.stage-item:hover {
  border-color: #1890ff;
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.15);
  transform: translateY(-2px);
}

.stage-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  background: linear-gradient(90deg, #f5f7fa 0%, #ffffff 100%);
  border-bottom: 2px solid #f0f0f0;
}

.stage-number {
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%);
  color: white;
  font-weight: 700;
  font-size: 16px;
  border-radius: 50%;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.3);
}

.stage-name-input {
  flex: 1;
  border-radius: 8px;
}

.delete-btn {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.delete-btn:hover {
  background: #fff1f0;
  transform: scale(1.1);
}

.stage-content {
  padding: 20px;
}

.form-field {
  margin-bottom: 0;
}

.field-label {
  display: block;
  margin-bottom: 8px;
  color: #595959;
  font-size: 14px;
  font-weight: 500;
}

.option-text {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.add-stage-btn {
  margin-top: 16px;
  height: 48px;
  border: 2px dashed #d9d9d9;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 500;
  color: #595959;
  transition: all 0.3s ease;
}

.add-stage-btn:hover {
  border-color: #1890ff;
  color: #1890ff;
  background: #f0f7ff;
  transform: translateY(-1px);
}

/* Ant Design 组件样式覆盖 */
:deep(.ant-input),
:deep(.ant-select-selector) {
  border-radius: 8px;
}

:deep(.ant-badge-count) {
  box-shadow: none;
  font-weight: 600;
}

:deep(.ant-modal-header) {
  border-bottom: 2px solid #f0f0f0;
  padding: 20px 24px;
}

:deep(.ant-modal-title) {
  font-size: 18px;
  font-weight: 600;
  color: #262626;
}

:deep(.ant-modal-footer) {
  border-top: 2px solid #f0f0f0;
  padding: 16px 24px;
}
</style>
