<template>
  <a-modal :open="props.open" :title="props.title" width="50%" @cancel="handleCancel">
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
  environment: [{ required: true, message: '请选择环境', trigger: 'blur' }],
  organization_key: [{ required: true, message: '请选择组织', trigger: 'blur' }],
  db_type: [{ required: true, message: '请选择数据库类型', trigger: 'blur' }],
  use_type: [{ required: true, message: '请选择用途', trigger: 'blur' }],
  hostname: [
    { required: true, message: '请输入主机名', trigger: 'blur' },
    { min: 3, max: 256, message: '长度应在3~256个字符', trigger: 'blur' },
  ],
  port: [{ required: true, message: '请输入端口', trigger: 'blur' }],
  inspect_params: [
    {
      validator: (_, value) => {
        try {
          JSON.parse(value)
        } catch (error) {
          return Promise.reject(new Error('请输入合法的JSON格式'))
        }
        return Promise.resolve()
      },
      trigger: 'blur',
    },
  ],
  remark: [
    { required: true, message: '请输入备注', trigger: 'blur' },
    { min: 3, max: 256, message: '长度应在3~256个字符', trigger: 'blur' },
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
