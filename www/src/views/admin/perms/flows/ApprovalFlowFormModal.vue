<template>
  <a-modal
    :open="props.open"
    :footer="null"
    :width="920"
    centered
    class="flow-form-modal"
    @cancel="handleCancel"
  >
    <div class="modal-shell">
      <div class="modal-head">
        <span class="head-badge">流程配置</span>
        <h3>{{ props.title }}</h3>
        <p>设置审批流名称、可领取人和每个阶段的审批规则，支持会签与或签组合。</p>
      </div>

      <a-form ref="formRef" layout="vertical" :model="formData" class="flow-form">
        <a-row :gutter="12">
          <a-col :span="24">
            <a-form-item
              label="审批流名称"
              name="name"
              required
              :rules="[
                {
                  required: true,
                  message: '请输入审批流名称',
                },
              ]"
            >
              <a-input v-model:value="formData.name" placeholder="请输入审批流名称" allow-clear />
            </a-form-item>
          </a-col>

          <a-col :span="24">
            <a-form-item
              label="可领取人"
              name="claim_users"
              required
              :rules="[
                {
                  required: true,
                  message: '请选择可领取人',
                },
              ]"
            >
              <a-select
                v-model:value="formData.claim_users"
                mode="multiple"
                show-search
                :filter-option="filterUserOption"
                placeholder="请选择可领取人（谁领取谁执行）"
                :options="props.userOptions"
                option-label-prop="label"
                :max-tag-count="4"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <div class="stage-header-bar">
          <div class="stage-title">审批阶段定义</div>
          <a-tag color="processing">共 {{ stageCount }} 个阶段</a-tag>
        </div>

        <div class="stage-list">
          <div v-for="(stage, index) in formData.definition" :key="index" class="stage-item">
            <div class="stage-item-head">
              <div class="stage-index">{{ index + 1 }}</div>
              <a-input
                v-model:value="stage.stage_name"
                placeholder="阶段名称（如：部门经理审批）"
                class="stage-name-input"
              />
              <a-button
                v-if="formData.definition.length > 1"
                type="text"
                danger
                class="delete-stage-btn"
                @click="removeStage(index)"
              >
                <DeleteOutlined />
              </a-button>
            </div>

            <div class="stage-item-body">
              <a-row :gutter="12">
                <a-col :xs="24" :md="8">
                  <a-form-item label="审批类型" required class="form-field">
                    <a-select v-model:value="stage.type">
                      <a-select-option value="AND">会签 (AND)</a-select-option>
                      <a-select-option value="OR">或签 (OR)</a-select-option>
                    </a-select>
                  </a-form-item>
                </a-col>
                <a-col :xs="24" :md="16">
                  <a-form-item label="审批人" required class="form-field">
                    <a-select
                      v-model:value="stage.approvers"
                      mode="multiple"
                      show-search
                      :filter-option="filterUserOption"
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
          <PlusOutlined />
          增加审批阶段
        </a-button>
      </a-form>

      <div class="modal-footer">
        <a-button @click="handleCancel">取消</a-button>
        <a-button type="primary" :loading="uiState.loading" @click="onSubmit">保存审批流</a-button>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { computed, reactive, ref } from 'vue'

const props = defineProps({
  open: Boolean,
  title: String,
  userOptions: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:open', 'submit'])

const formData = defineModel('formData', {
  type: Object,
  required: true,
})

const formRef = ref()

const uiState = reactive({
  loading: false,
})

const stageCount = computed(() => formData.value.definition?.length || 0)

const filterUserOption = (input, option) => {
  if (!option) return false
  const label = (option.label || '').toString().toLowerCase()
  const value = (option.value || '').toString().toLowerCase()
  return label.includes(input.toLowerCase()) || value.includes(input.toLowerCase())
}

const addStage = () => {
  const definition = formData.value.definition || (formData.value.definition = [])
  definition.push({
    stage: definition.length + 1,
    stage_name: `新阶段 ${definition.length + 1}`,
    approvers: [],
    type: 'AND',
  })
}

const removeStage = (index) => {
  const definition = formData.value.definition || []
  definition.splice(index, 1)
  definition.forEach((stage, i) => {
    stage.stage = i + 1
  })
}

const handleCancel = () => {
  emit('update:open', false)
  formRef.value?.resetFields?.()
}

const validateDefinition = async () => {
  const claimUsers = formData.value.claim_users || []
  if (claimUsers.length === 0) {
    throw new Error('请至少选择一个可领取人。')
  }

  const definition = formData.value.definition
  if (!definition || definition.length === 0) {
    throw new Error('请至少配置一个审批阶段。')
  }

  for (const [index, stage] of definition.entries()) {
    if (!stage.stage_name || stage.stage_name.trim() === '') {
      throw new Error(`阶段 ${index + 1}：阶段名称不能为空。`)
    }
    if (!stage.approvers || stage.approvers.length === 0) {
      throw new Error(`阶段 ${index + 1}：至少需要指定一个审批人。`)
    }
  }
}

const onSubmit = async () => {
  try {
    await validateDefinition()
    uiState.loading = true
    emit('submit', formData.value)
  } catch (err) {
    const errorText = err instanceof Error ? err.message : String(err)
    message.error(errorText)
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
  color: #1554c2;
  border-radius: 999px;
  border: 1px solid rgba(31, 111, 235, 0.28);
  background: rgba(31, 111, 235, 0.1);
  padding: 3px 10px;
}

.flow-form {
  margin-top: 16px;
}

.stage-header-bar {
  margin: 2px 0 10px;
  padding: 10px 12px;
  border-radius: 10px;
  border: 1px solid #dbe6f8;
  background: linear-gradient(90deg, #f5f8ff 0%, #ffffff 100%);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.stage-title {
  color: #22335d;
  font-size: 14px;
  font-weight: 700;
}

.stage-list {
  max-height: 420px;
  overflow-y: auto;
  padding-right: 2px;
}

.stage-item {
  border: 1px solid #dfe7f8;
  border-radius: 12px;
  background: #ffffff;
  margin-bottom: 10px;
  overflow: hidden;
}

.stage-item-head {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 12px;
  border-bottom: 1px solid #ecf1fb;
  background: linear-gradient(90deg, #f8fbff 0%, #ffffff 100%);
}

.stage-index {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: #1f6feb;
  color: #fff;
  font-size: 14px;
  font-weight: 700;
}

.stage-name-input {
  flex: 1;
}

.delete-stage-btn {
  width: 30px;
  height: 30px;
  border-radius: 8px;
}

.stage-item-body {
  padding: 12px;
}

.form-field {
  margin-bottom: 0;
}

.add-stage-btn {
  margin-top: 6px;
  height: 40px;
  border: 1px dashed #bfd5fb;
  color: #1f6feb;
  border-radius: 10px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 14px;
}

:deep(.flow-form-modal .ant-modal-content) {
  border-radius: 18px;
  padding: 20px 22px;
  background:
    radial-gradient(circle at right top, rgba(31, 111, 235, 0.08), rgba(31, 111, 235, 0) 55%),
    #ffffff;
}

:deep(.flow-form-modal .ant-form-item-label > label) {
  color: #23335c;
  font-weight: 600;
}

:deep(.flow-form-modal .ant-input-affix-wrapper),
:deep(.flow-form-modal .ant-input),
:deep(.flow-form-modal .ant-select-selector) {
  border-radius: 10px !important;
  border-color: #d7e2f3 !important;
}

:deep(.flow-form-modal .ant-select-focused .ant-select-selector),
:deep(.flow-form-modal .ant-input-affix-wrapper-focused),
:deep(.flow-form-modal .ant-input:focus) {
  border-color: #1f6feb !important;
  box-shadow: 0 0 0 3px rgba(31, 111, 235, 0.14);
}

@media (max-width: 900px) {
  .stage-list {
    max-height: 360px;
  }

  .modal-head h3 {
    font-size: 20px;
  }
}
</style>
