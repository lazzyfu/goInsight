<template>
  <a-modal :open="props.open" :title="props.title" @cancel="handleCancel">
    <template #footer>
      <a-button @click="handleCancel">取消</a-button>
      <a-button type="primary" :loading="uiState.loading" @click="onSubmit">确定</a-button>
    </template>

    <a-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      :label-col="{ span: 4 }"
      :wrapper-col="{ span: 20 }"
      style="margin-top: 24px"
    >
      <a-form-item label="表名" name="tables" has-feedback>
        <a-select
          ref="select"
          v-model:value="formData.tables"
          :options="props.tables"
          mode="multiple"
          :field-names="{ label: 'table_name', value: 'table_name' }"
          allowClear
        >
        </a-select>
      </a-form-item>

      <a-form-item label="规则" name="rule" has-feedback>
        <a-select
          ref="select"
          v-model:value="formData.rule"
          :options="[
            { value: 'allow', label: '允许' },
            { value: 'deny', label: '拒绝' },
          ]"
          allowClear
        >
        </a-select>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { ref, reactive } from 'vue'

// 定义props和emits
const props = defineProps({
  open: Boolean,
  title: String,
  tables: Array,
})
const emit = defineEmits(['update:open', 'submit'])

// 表单数据
const formData = defineModel('modelValue', {
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
  tables: [
    {
      required: true,
      message: '请选择至少一个表名',
      trigger: 'change',
    },
  ],
  rule: [
    {
      required: true,
      message: '请选择规则',
      trigger: 'change',
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
    emit('submit', formData.value)
  } catch (err) {
  } finally {
    uiState.loading = false
  }
}
</script>
