<template>
  <a-modal :open="props.open" :title="props.title" width="35%" @cancel="handleCancel">
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
      <a-form-item label="描述" name="title" has-feedback>
        <a-input disabled v-model:value="formData.title" placeholder="" />
      </a-form-item>

      <a-form-item label="键" name="key" has-feedback>
        <a-input disabled v-model:value="formData.key" placeholder="" />
      </a-form-item>

      <a-form-item label="类型" name="type" has-feedback>
        <a-input disabled v-model:value="formData.type" placeholder="" />
      </a-form-item>

  <a-form-item label="值" name="value" has-feedback>
        <template v-if="formData.type === 'boolean'">
          <!-- ant-design-vue 的 SelectOption 在部分版本里 value 不支持 boolean，这里统一用字符串承载 -->
          <a-select v-model:value="formData._editValue" placeholder="请选择" style="width: 100%">
            <a-select-option value="true">true</a-select-option>
            <a-select-option value="false">false</a-select-option>
          </a-select>
        </template>

        <template v-else-if="formData.type === 'number'">
          <a-input-number
            v-model:value="formData._editValue"
            style="width: 100%"
            :precision="0"
            placeholder="请输入数字"
          />
        </template>

        <template v-else>
          <a-input v-model:value="formData._editValue" placeholder="请输入" allow-clear />
        </template>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { computed, ref, reactive, watch } from 'vue'

// 定义props和emits
const props = defineProps({
  open: Boolean,
  title: String,
})
const emit = defineEmits(['update:open', 'submit'])

// 表单数据
const formData = defineModel('modelValue', {
  type: Object,
  required: true,
})

// 编辑态的值：根据 type 用 boolean/number/string 承载，保存时再转回字符串提交给后端
const normalizeEditValue = (record) => {
  const type = record?.type
  const raw = record?.value

  if (type === 'boolean') {
    // Select 用字符串承载
    if (raw === true) return 'true'
    if (raw === false) return 'false'
    const s = String(raw ?? '').trim().toLowerCase()
    return s === 'true' ? 'true' : 'false'
  }
  if (type === 'number') {
    if (typeof raw === 'number') return raw
    const n = Number(String(raw ?? '').trim())
    return Number.isFinite(n) ? n : undefined
  }
  return raw ?? ''
}

watch(
  () => props.open,
  (open) => {
    if (!open) return
    // 打开弹窗时回填当前行数据
    // 约定：InspectList 传入的 formData 里包含 { id, title, key, type, value, ... }
    formData.value = {
      ...formData.value,
      _editValue: normalizeEditValue(formData.value),
    }
  },
)

// 表单引用
const formRef = ref()

// 状态
const uiState = reactive({
  loading: false,
})

// 表单校验规则
const rules = computed(() => {
  const type = formData.value?.type
  if (type === 'number') {
    return {
      value: [
        {
          validator: () => {
            const v = formData.value?._editValue
            if (v === undefined || v === null || v === '') return Promise.reject('请输入数字')
            if (typeof v === 'number' && Number.isFinite(v)) return Promise.resolve()
            return Promise.reject('请输入合法数字')
          },
          trigger: ['blur', 'change'],
        },
      ],
    }
  }
  if (type === 'boolean') {
    return {
      value: [
        {
          validator: () => {
            const v = formData.value?._editValue
            if (v === 'true' || v === 'false') return Promise.resolve()
            return Promise.reject('请选择 true/false')
          },
          trigger: ['blur', 'change'],
        },
      ],
    }
  }
  return {
    value: [
      {
        validator: () => {
          // string 允许空字符串
          return Promise.resolve()
        },
        trigger: ['blur', 'change'],
      },
    ],
  }
})

// 取消按钮处理函数
const handleCancel = () => {
  emit('update:open', false)
  formRef.value?.resetFields()
  if (formData.value && '_editValue' in formData.value) {
    // 避免下次打开弹窗残留上一次的临时值
    delete formData.value._editValue
  }
}

// 提交表单
const onSubmit = async () => {
  try {
    await formRef.value.validateFields()
    uiState.loading = true

    const type = formData.value?.type
    const v = formData.value?._editValue
    let valueAsString = ''
  if (type === 'boolean') valueAsString = v === 'true' ? 'true' : 'false'
    else if (type === 'number') valueAsString = String(v)
    else valueAsString = String(v ?? '')

    emit('submit', {
      ...formData.value,
      value: valueAsString,
    })
  } catch {
    // 校验失败时不处理
  } finally {
    uiState.loading = false
  }
}
</script>
