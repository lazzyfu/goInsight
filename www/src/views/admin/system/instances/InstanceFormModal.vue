<template>
  <a-modal :open="props.open" :title="props.title" width="35%" @cancel="handleCancel">
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
      style="margin-top: 24px"
    >
      <a-form-item label="环境" name="environment" has-feedback>
        <a-select
          ref="select"
          v-model:value="formData.environment"
          :options="props.environments"
          :field-names="{ label: 'name', value: 'id' }"
          allowClear
        ></a-select>
      </a-form-item>
      <a-form-item label="组织" name="organization_key" has-feedback>
        <a-cascader
          v-model:value="formData.organization_key"
          :field-names="{ label: 'title', value: 'key', children: 'children' }"
          :options="props.organizations"
          change-on-select
          expand-trigger="hover"
          placeholder="请选择组织"
        >
        </a-cascader>
      </a-form-item>

      <a-form-item label="类型" name="db_type" has-feedback>
        <a-select
          ref="select"
          v-model:value="formData.db_type"
          :options="[
            { value: 'MySQL', label: 'MySQL' },
            { value: 'TiDB', label: 'TiDB' },
            { value: 'ClickHouse', label: 'ClickHouse' },
          ]"
          allowClear
        >
        </a-select>
      </a-form-item>

      <a-form-item label="用途" name="use_type" has-feedback>
        <a-select
          ref="select"
          v-model:value="formData.use_type"
          :options="[
            { value: '查询', label: '查询' },
            { value: '工单', label: '工单' },
          ]"
          allowClear
        >
        </a-select>
      </a-form-item>

      <a-form-item label="主机名" name="hostname" has-feedback>
        <a-input v-model:value="formData.hostname" placeholder="请输入主机名" allow-clear />
      </a-form-item>

      <a-form-item label="端口" name="port" has-feedback>
        <a-input-number v-model:value="formData.port" placeholder="请输入端口" allow-clear />
      </a-form-item>

      <a-form-item
        label="审核参数"
        name="inspect_params"
        has-feedback
        help="格式要求为JSON类型，默认为{}，表示继承全局审核参数"
      >
        <a-textarea
          :rows="8"
          v-model:value="formData.inspect_params"
          placeholder=" 请输入自定义审核参数，默认为{}"
        />
      </a-form-item>

      <a-form-item label="备注" name="remark" has-feedback>
        <a-input v-model:value="formData.remark" placeholder="请输入备注" allow-clear />
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
  environments: Array,
  organizations: Array,
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
  environment: [{ required: true, message: '请选择所属环境', trigger: 'change' }],
  organization_key: [{ required: true, message: '请选择组织/部门', trigger: 'change' }],
  db_type: [{ required: true, message: '请选择数据库类型', trigger: 'change' }],
  use_type: [{ required: true, message: '请选择实例用途', trigger: 'change' }],
  hostname: [
    { required: true, message: '请输入主机名', trigger: ['blur', 'change'] },
    {
      validator: (_, value) => {
        const v = (value || '').trim()
        if (v.length < 3 || v.length > 256) {
          return Promise.reject(new Error('主机名长度需在 3 到 256 个字符之间'))
        }
        // 允许字母数字、点、短横线，且不能包含空格
        const hostRe = /^[A-Za-z0-9.-]+$/
        if (!hostRe.test(v)) {
          return Promise.reject(new Error('主机名仅支持字母、数字、点(.)、短横线(-)'))
        }
        return Promise.resolve()
      },
      trigger: 'blur',
    },
  ],
  port: [
    { required: true, message: '请输入端口号', trigger: ['blur', 'change'] },
    {
      validator: (_, value) => {
        const num = Number(value)
        if (!Number.isInteger(num)) {
          return Promise.reject(new Error('端口号必须为整数'))
        }
        if (num < 1 || num > 65535) {
          return Promise.reject(new Error('端口号范围为 1-65535'))
        }
        return Promise.resolve()
      },
      trigger: 'blur',
    },
  ],
  inspect_params: [
    {
      validator: (_, value) => {
        const v = (value || '').trim()
        if (!v) return Promise.resolve() // 空表示继承全局参数
        try {
          const parsed = JSON.parse(v)
          if (typeof parsed !== 'object' || Array.isArray(parsed)) {
            return Promise.reject(new Error('审核参数需为 JSON 对象，如 {}'))
          }
        } catch {
          return Promise.reject(new Error('请输入合法的 JSON 格式，如 {"key":"value"}'))
        }
        return Promise.resolve()
      },
      trigger: 'blur',
    },
  ],
  remark: [
    { required: true, message: '请输入备注', trigger: ['blur', 'change'] },
    {
      validator: (_, value) => {
        const v = (value || '').trim()
        if (v.length < 3 || v.length > 256) {
          return Promise.reject(new Error('备注长度需在 3 到 256 个字符之间'))
        }
        return Promise.resolve()
      },
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
    emit('submit', formData.value)
  } catch (err) {
  } finally {
    uiState.loading = false
  }
}
</script>
