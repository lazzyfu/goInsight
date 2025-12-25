<template>
  <a-modal :open="props.open" :title="props.title" @cancel="handleCancel">
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
          :options="uiData.dbTypes"
          allowClear
        >
        </a-select>
      </a-form-item>

      <a-form-item label="用途" name="use_type" has-feedback>
        <a-select
          ref="select"
          v-model:value="formData.use_type"
          :options="uiData.useTypes"
          allowClear
        >
        </a-select>
      </a-form-item>

      <a-form-item label="主机名" name="hostname" has-feedback>
        <a-input v-model:value="formData.hostname" placeholder="请输入主机名" />
      </a-form-item>

      <a-form-item label="端口" name="port" has-feedback>
        <a-input-number v-model:value="formData.port" placeholder="请输入端口" />
      </a-form-item>

      <a-form-item label="用户" name="user" has-feedback>
        <a-input v-model:value="formData.user" placeholder="请输入用户" />
      </a-form-item>

      <a-form-item label="密码" name="password" has-feedback>
        <a-input-password v-model:value="formData.password" placeholder="请输入密码" />
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
        <a-input v-model:value="formData.remark" placeholder="请输入备注" />
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

const uiData = reactive({
  dbTypes: [
    { value: 'MySQL', label: 'MySQL' },
    { value: 'TiDB', label: 'TiDB' },
    { value: 'ClickHouse', label: 'ClickHouse' },
  ],
  useTypes: [
    { value: '查询', label: '查询' },
    { value: '工单', label: '工单' },
  ],
})

// 表单校验规则
const rules = {
  environment: [{ required: true, message: '请选择所属环境', trigger: 'change' }],
  organization_key: [{ required: true, message: '请选择组织/部门', trigger: 'change' }],
  db_type: [{ required: true, message: '请选择数据库类型', trigger: 'change' }],
  use_type: [{ required: true, message: '请选择实例用途', trigger: 'change' }],
  hostname: [
    { required: true, message: '请输入主机名', trigger: 'blur' },
    {
      validator: (_, value) => {
        const v = (value || '').trim()
        if (v.length < 3 || v.length > 255) {
          return Promise.reject(new Error('主机名长度需在 3 到 255 个字符之间'))
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
  user: [
    { required: true, message: '请输入访问数据库的用户名', trigger: 'blur' },
    {
      validator: (_, value) => {
        const v = (value || '').trim()
        if (v.length < 3 || v.length > 32) {
          return Promise.reject(new Error('用户名长度需在 3 到 32 个字符之间'))
        }
        // MySQL user只允许字母数字、下划线、短横线
        const userRe = /^[A-Za-z0-9_-]+$/
        if (!userRe.test(v)) {
          return Promise.reject(new Error('用户名仅支持字母、数字、下划线(_)、短横线(-)'))
        }
        return Promise.resolve()
      },
      trigger: 'blur',
    },
  ],
  password: [
    { required: true, message: '请输入访问数据库的密码', trigger: 'blur' },
    {
      validator: (_, value) => {
        const v = (value || '').trim()
        if (v.length < 8 || v.length > 64) {
          return Promise.reject(new Error('密码长度需在 8 到 64 个字符之间'))
        }
        // MySQL密码支持大多数字符，建议包含大小写字母、数字和特殊字符
        const passwordRe = /^[A-Za-z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]+$/
        if (!passwordRe.test(v)) {
          return Promise.reject(new Error('密码包含不支持的字符'))
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
