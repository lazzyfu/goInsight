<template>
  <a-card title="提交工单">
    <a-row :gutter="{ xs: 8, sm: 16, md: 24, lg: 32 }">
      <a-col :xs="24" :sm="24" :md="24" :lg="10" :xl="8">
        <a-form ref="formRef" :model="formState" :label-col="{ span: 6 }" :wrapper-col="{ span: 18 }" autocomplete="off"
          :rules="rules" @finish="onSubmit">
          <a-form-item label="标题" name="title" has-feedback>
            <a-input v-model:value="formState.title" placeholder="请输入标题" />
          </a-form-item>

          <a-form-item label="备注" name="remark">
            <a-textarea v-model:value="formState.remark" placeholder="请输入工单需求或备注" />
          </a-form-item>

          <a-form-item label="类型" name="sql_type" has-feedback>
            <a-select v-model:value="formState.sql_type" :options="uiData.sqlTypes" allowClear />
          </a-form-item>

          <a-form-item label="环境" name="environment" has-feedback>
            <a-select v-model:value="formState.environment" :options="uiData.environments"
              :field-names="{ label: 'name', value: 'id' }" @change="onEnvChange" allowClear />
          </a-form-item>

          <a-form-item label="数据库" name="db_type" has-feedback>
            <a-select v-model:value="formState.db_type" :options="uiData.dbTypes" @change="onDbTypeChange"
              :disabled="!formState.environment" allowClear />
          </a-form-item>

          <a-form-item label="实例" name="instance_id" has-feedback>
            <a-select v-model:value="formState.instance_id" :options="uiData.instances"
              :field-names="{ label: 'remark', value: 'instance_id' }" @change="onInstanceChange"
              :disabled="!formState.db_type" allowClear />
          </a-form-item>

          <a-form-item label="库名" name="schema" has-feedback>
            <a-select v-model:value="formState.schema" :options="uiData.schemas"
              :field-names="{ label: 'schema', value: 'schema' }" :disabled="!formState.instance_id" allowClear />
          </a-form-item>

          <a-form-item label="抄送" name="cc">
            <a-select v-model:value="formState.cc" mode="multiple" :options="uiData.users"
              :field-names="{ label: 'nick_name', value: 'username' }" allowClear />
          </a-form-item>

          <a-form-item :wrapper-col="{ offset: 4 }">
            <a-button type="primary" html-type="submit" :loading="uiState.loading">提交</a-button>
            <a-button style="margin-left: 10px" @click="resetForm">重置</a-button>
          </a-form-item>
        </a-form>
      </a-col>

      <a-col :xs="24" :sm="24" :md="24" :lg="14" :xl="16">
        <a-space size="small">
          <a-button @click="formatSql">
            <template #icon>
              <CodeOutlined />
            </template>
            格式化
          </a-button>

          <a-button @click="checkSyntax">
            <template #icon>
              <CodeOutlined />
            </template>
            语法检查
          </a-button>
        </a-space>

        <div style="margin-top: 6px">
          <CodeMirror ref="codemirrorRef" />
        </div>
      </a-col>
    </a-row>

    <order-inspect ref="inspectResultTableRef" />
  </a-card>
</template>

<script setup>
import {
  createOrderApi,
  getOrderEnvironmentsApi,
  getOrderInstancesApi,
  getOrderSchemasApi,
  getOrderUsersApi,
  inspectOrderSyntaxApi,
} from '@/api/order'

import CodeMirror from '@/components/edit/Codemirror.vue'
import { CodeOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { debounce } from 'lodash-es'
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import OrderInspect from './OrderInspect.vue'

const router = useRouter()
const formRef = ref()
const codemirrorRef = ref()
const inspectResultTableRef = ref()

// 表单
const formState = reactive({
  title: '',
  remark: '',
  sql_type: '',
  environment: '',
  db_type: '',
  instance_id: '',
  schema: '',
  cc: [],
  content: '',
  export_file_format: 'CSV',
})

// 规则
const rules = {
  title: [
    { required: true, message: '请输入标题', trigger: 'blur' },
    { min: 3, max: 196, message: '标题长度需在 3-196 字符之间', trigger: 'blur' },
  ],
  remark: [
    { required: true, message: '请输入备注', trigger: 'blur' },
    { min: 3, max: 1024, message: '备注长度需在 3-1024 字符之间', trigger: 'blur' },
  ],
  sql_type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  environment: [{ required: true, message: '请选择环境', trigger: 'change' }],
  db_type: [{ required: true, message: '请选择数据库类型', trigger: 'change' }],
  instance_id: [{ required: true, message: '请选择实例', trigger: 'change' }],
  schema: [{ required: true, message: '请选择库名', trigger: 'change' }],
  cc: [],
}

// 状态
const uiState = reactive({
  loading: false,
})

// 数据
const uiData = reactive({
  environments: [],
  instances: [],
  schemas: [],
  users: [],
  dbTypes: [
    { value: 'MySQL', label: 'MySQL' },
    { value: 'TiDB', label: 'TiDB' },
  ],
  sqlTypes: [
    { value: 'DML', label: 'DML' },
    { value: 'DDL', label: 'DDL' },
  ],
})

// 下拉联动
const resetDbLevel = () => {
  formState.db_type = ''
  formState.instance_id = ''
  formState.schema = ''
  uiData.instances = []
  uiData.schemas = []
}

// 环境变化
const onEnvChange = () => {
  resetDbLevel()
}

// 数据库类型变化
const onDbTypeChange = debounce(async () => {
  formState.instance_id = ''
  formState.schema = ''
  uiData.instances = []
  uiData.schemas = []

  if (!formState.db_type) return

  const res = await getOrderInstancesApi({
    id: formState.environment,
    db_type: formState.db_type,
    is_page: false,
  }).catch(() => { })

  uiData.instances = res?.data || []
}, 200)

// 实例变化
const onInstanceChange = debounce(async () => {
  formState.schema = ''
  uiData.schemas = []

  if (!formState.instance_id) return

  const res = await getOrderSchemasApi({
    instance_id: formState.instance_id,
    is_page: false,
  }).catch(() => { })

  uiData.schemas = res?.data || []
}, 200)

// 提交工单
const onSubmit = async () => {
  uiState.loading = true
  formState.content = codemirrorRef.value?.getContent() || ''

  if (!formState.content) {
    message.warning('SQL内容不能为空')
    uiState.loading = false
    return
  }

  const res = await createOrderApi(formState).catch(() => { })
  if (!res) {
    uiState.loading = false
    return message.error('创建失败，请稍后再试')
  }

  message.success('工单创建成功')
  router.push({ name: 'orders.list' })
  uiState.loading = false
}

// 格式化
const formatSql = () => {
  codemirrorRef.value?.formatContent()
}

// 语法检查
const checkSyntax = debounce(async () => {
  const sql = codemirrorRef.value?.getContent() || ''
  if (!sql) return message.warning('SQL内容不能为空')
  if (!formState.environment || !formState.db_type || !formState.instance_id || !formState.schema) {
    return message.warning('请先选择环境、数据库类型、实例、库名')
  }
  const res = await inspectOrderSyntaxApi({
    db_type: formState.db_type,
    sql_type: formState.sql_type,
    instance_id: formState.instance_id,
    schema: formState.schema,
    content: sql,
  }).catch(() => { })
  if (res) inspectResultTableRef.value.render(res.data)
}, 300)

// 重置表单
const resetForm = () => {
  formRef.value?.resetFields()
  codemirrorRef.value?.setContent('')
}

// 初始化
onMounted(async () => {
  const [envRes, userRes] = await Promise.all([
    getOrderEnvironmentsApi().catch(() => { }),
    getOrderUsersApi({ is_page: false }).catch(() => { }),
  ])

  uiData.environments = envRes?.data || []
  uiData.users = userRes?.data || []
})
</script>
