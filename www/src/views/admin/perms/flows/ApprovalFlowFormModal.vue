<template>
  <a-modal
    :open="props.open"
    :title="props.title"
    width="800px"
    destroyOnClose
    @cancel="handleCancel"
  >
    <template #footer>
      <a-button @click="handleCancel">取消</a-button>
      <a-button type="primary" :loading="uiState.loading" @click="onSubmit">确定</a-button>
    </template>
    <a-form
      ref="formRef"
      :label-col="{ span: 4 }"
      :wrapper-col="{ span: 20 }"
      :model="formData"
      :rules="rules"
    >
      <a-form-item label="审批流名称" name="name" has-feedback>
        <a-input v-model:value="formData.name" placeholder="请输入审批流名称" allow-clear />
      </a-form-item>

      <a-divider orientation="left">
        审批阶段定义 (共 {{ formData.definition ? formData.definition.length : 0 }} 个阶段)
      </a-divider>

      <a-form-item
        label=""
        :label-col="{ span: 0 }"
        :wrapper-col="{ span: 24 }"
        class="dynamic-definition-item"
      >
        <div class="stage-list">
          <div v-for="(stage, index) in formData.definition" :key="index" class="stage-item">
            <div class="stage-header">
              <span class="stage-title">阶段 {{ index + 1 }}</span>
              <a-input
                v-model:value="stage.stage_name"
                placeholder="阶段名称 (如：部门经理审批)"
                style="width: 250px; margin-right: 16px"
              />
              <a-button
                v-if="formData.definition.length > 1"
                type="danger"
                size="small"
                @click="removeStage(index)"
                danger
              >
                <DeleteOutlined /> 删除
              </a-button>
            </div>
            <a-row :gutter="16">
              <a-col :span="6">
                <a-form-item label="审批类型" class="inner-form-item">
                  <a-select v-model:value="stage.type">
                    <a-select-option value="AND">AND (会签)</a-select-option>
                    <a-select-option value="OR">OR (或签)</a-select-option>
                  </a-select>
                </a-form-item>
              </a-col>
              <a-col :span="18">
                <a-form-item label="审批人 (用户名列表)" class="inner-form-item">
                  <a-select
                    v-model:value="stage.approvers"
                    mode="multiple"
                    show-search
                    :filter-option="filterUserOption"
                    style="width: 100%"
                    placeholder="请选择审批人"
                    :options="props.userOptions"
                    option-label-prop="label"
                  />
                </a-form-item>
              </a-col>
            </a-row>
          </div>
        </div>
        <a-button type="dashed" style="width: 100%; margin-top: 16px" @click="addStage">
          <PlusOutlined /> 增加审批阶段
        </a-button>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { ref, reactive } from 'vue'

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

// 表单数据
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

// 表单校验规则
const rules = {
  name: [
    {
      required: true,
      message: '请输入审批流名称',
      trigger: 'blur',
    },
  ],
}

const filterUserOption = (input, option) => {
  return (
    option.label.toLowerCase().includes(input.toLowerCase()) ||
    option.value.toLowerCase().includes(input.toLowerCase())
  )
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
}

const validateDefinition = async () => {
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
  return true
}

// 提交表单
const onSubmit = async () => {
  try {
    await validateDefinition()
    uiState.loading = true
    emit('submit', formData.value)
  } catch (err) {
  } finally {
    uiState.loading = false
  }
}
</script>

<style scoped>
.dynamic-definition-item :deep(.ant-form-item-control-input-content) {
  padding: 0;
}

.stage-list {
  max-height: 400px;
  overflow-y: auto;
  padding-right: 8px;
}
.stage-item {
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  padding: 16px;
  margin-bottom: 16px;
  background-color: #fafafa;
}
.stage-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px dashed #e0e0e0;
}
.stage-title {
  font-weight: 600;
  color: #333;
  margin-right: 16px;
  white-space: nowrap;
}

/* 调整动态阶段内嵌的 form-item 样式，去除上下边距 */
.inner-form-item :deep(.ant-form-item) {
  margin-bottom: 0px;
}
</style>
