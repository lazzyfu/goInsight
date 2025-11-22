<template>
  <a-modal :open="props.open" :title="props.title" :footer="null" @cancel="handleCancel">
    <a-form
      ref="formRef"
      :label-col="{ span: 4 }"
      :wrapper-col="{ span: 18 }"
      :model="formData"
      :rules="rules"
      @finish="onSubmit"
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

      <a-form-item :wrapper-col="{ offset: 4, span: 18 }" style="text-align: right">
        <a-space>
          <a-button @click="handleCancel">取消</a-button>
          <a-button type="primary" html-type="submit">确定</a-button>
        </a-space>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { ref } from 'vue'

// 定义props和emits
const props = defineProps({
  open: Boolean,
  title: String,
  tables: Array,
})
const emit = defineEmits(['update:open', 'submit'])

// 使用defineModel接收 v-model:modelValue
// 它自动创建了一个名为modelValue的prop，并提供了一个value来读取，以及一个自动触发update:modelValue的setter
const formData = defineModel('modelValue', {
  type: Object,
  required: true,
})

// 表单引用
const formRef = ref()

// 表单校验规则
const rules = {
  tables: [
    {
      required: true,
      message: '不能为空，请输入环境名',
      trigger: 'change',
    },
  ],
  rule: [
    {
      required: true,
      message: '不能为空，请选择规则',
      trigger: 'change',
    },
  ],
}

// 取消按钮处理函数
const handleCancel = () => {
  emit('update:open', false)
  formRef.value?.resetFields()
}

// 提交表单处理函数
const onSubmit = async () => {
  emit('submit', formData.value)
}
</script>
