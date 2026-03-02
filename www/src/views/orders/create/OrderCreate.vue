<template>
  <a-card class="order-create-card">
    <template #title>
      <div class="title-wrap">
        <span class="title-main">提交工单</span>
        <span class="title-sub">完善工单信息，编写 SQL 并通过语法检查后再提交审批</span>
      </div>
    </template>

    <a-row class="create-row" :gutter="{ xs: 12, sm: 16, md: 20, lg: 24 }">
        <a-col class="pane-col" :xs="24" :sm="24" :md="24" :lg="10" :xl="8">
          <div class="panel form-panel">
            <a-alert
              message="请完整填写环境、实例与库名信息，避免提交到错误数据源"
              type="info"
              show-icon
              banner
              class="form-tip"
            />

            <a-form ref="formRef" :model="formState" layout="vertical" autocomplete="off" :rules="rules" @finish="onSubmit">
              <a-form-item label="标题" name="title" has-feedback>
                <a-input v-model:value="formState.title" placeholder="请输入标题" />
              </a-form-item>

              <a-form-item label="备注" name="remark">
                <a-textarea
                  v-model:value="formState.remark"
                  placeholder="请输入工单需求或备注"
                  :auto-size="{ minRows: 3, maxRows: 6 }"
                />
              </a-form-item>

              <a-form-item label="类型" name="sql_type" has-feedback>
                <a-select v-model:value="formState.sql_type" :options="uiData.sqlTypes" placeholder="请选择工单类型" allowClear />
              </a-form-item>

              <a-form-item label="环境" name="environment" has-feedback>
                <a-select
                  v-model:value="formState.environment"
                  show-search
                  :filter-option="filterEnvOption"
                  :options="uiData.environments"
                  :field-names="{ label: 'name', value: 'id' }"
                  @change="onEnvChange"
                  placeholder="请选择环境"
                  allowClear
                />
              </a-form-item>

              <a-form-item label="数据库" name="db_type" has-feedback>
                <a-select
                  v-model:value="formState.db_type"
                  :options="uiData.dbTypes"
                  @change="onDbTypeChange"
                  :disabled="!formState.environment"
                  placeholder="请选择数据库类型"
                  allowClear
                />
              </a-form-item>

              <a-form-item label="实例" name="instance_id" has-feedback>
                <a-select
                  v-model:value="formState.instance_id"
                  show-search
                  :filter-option="filterInstanceOption"
                  :options="uiData.instances"
                  :field-names="{ label: 'remark', value: 'instance_id' }"
                  @change="onInstanceChange"
                  :disabled="!formState.db_type"
                  placeholder="请选择实例"
                  allowClear
                />
              </a-form-item>

              <a-form-item label="库名" name="schema" has-feedback>
                <a-select
                  v-model:value="formState.schema"
                  show-search
                  :filter-option="filterSchemaOption"
                  :options="uiData.schemas"
                  :field-names="{ label: 'schema', value: 'schema' }"
                  :disabled="!formState.instance_id"
                  placeholder="请选择库名"
                  allowClear
                />
              </a-form-item>

              <a-form-item label="抄送" name="cc">
                <a-select
                  v-model:value="formState.cc"
                  mode="multiple"
                  show-search
                  :filter-option="filterCcOption"
                  :options="uiData.users"
                  :field-names="{ label: 'nick_name', value: 'username' }"
                  placeholder="可选，通知相关同学"
                  allowClear
                />
              </a-form-item>

              <a-form-item class="form-actions">
                <a-space>
                  <a-button type="primary" html-type="submit" :loading="uiState.loading">提交</a-button>
                  <a-button @click="resetForm">重置</a-button>
                </a-space>
              </a-form-item>
            </a-form>
          </div>
        </a-col>

        <a-col class="pane-col" :xs="24" :sm="24" :md="24" :lg="14" :xl="16">
          <div class="panel editor-panel">
            <div class="editor-header">
              <div class="editor-title">
                <span class="editor-title-main">SQL 编辑区</span>
                <span class="editor-title-sub">建议先格式化，再做语法检查</span>
              </div>
              <a-space size="small" wrap>
                <a-button @click="formatSql">
                  <template #icon>
                    <CodeOutlined />
                  </template>
                  格式化
                </a-button>

                <a-button type="primary" ghost @click="checkSyntax">
                  <template #icon>
                    <CodeOutlined />
                  </template>
                  语法检查
                </a-button>
              </a-space>
            </div>

            <div class="editor-shell">
              <CodeMirror ref="codemirrorRef" :height="'100%'" />
            </div>
          </div>
        </a-col>
    </a-row>

    <div class="inspect-wrapper">
      <order-inspect ref="inspectResultTableRef" v-model:modelValue="formState" />
    </div>
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
import { useOrderCreatePrefillStore } from '@/store/prefill'
import { CodeOutlined } from '@ant-design/icons-vue'
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
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
    { value: 'DML', label: 'DML工单' },
    { value: 'DDL', label: 'DDL工单' },
    { value: 'EXPORT', label: '数据导出工单' },
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

const filterEnvOption = (input, option) => {
  return option?.name?.toLowerCase().includes(input.toLowerCase())
}

const filterInstanceOption = (input, option) => {
  return option?.remark?.toLowerCase().includes(input.toLowerCase())
}

const filterSchemaOption = (input, option) => {
  return option?.schema?.toLowerCase().includes(input.toLowerCase())
}

const filterCcOption = (input, option) => {
  return option?.nick_name?.toLowerCase().includes(input.toLowerCase())
}

// 环境变化
const onEnvChange = () => {
  resetDbLevel()
}

// 数据库类型变化
const onDbTypeChange = useThrottleFn(async () => {
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
})

// 实例变化
const onInstanceChange = useThrottleFn(async () => {
  formState.schema = ''
  uiData.schemas = []

  if (!formState.instance_id) return

  const res = await getOrderSchemasApi({
    instance_id: formState.instance_id,
    is_page: false,
  }).catch(() => { })

  uiData.schemas = res?.data || []
})

// 提交工单
const onSubmit = useThrottleFn(async () => {
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
})

// 格式化
const formatSql = () => {
  codemirrorRef.value?.formatContent()
}

// 语法检查
const checkSyntax = useThrottleFn(async () => {
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
})

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

  // 如果存在来自 Pinia store 的预填数据，则读取并填充表单（一次性消费）
  const orderCreatePrefillStore = useOrderCreatePrefillStore()
  const prefill = orderCreatePrefillStore.consumeCreatePrefill()
  if (prefill) {
    formState.title = prefill.title || ''
    formState.remark = prefill.remark || ''
    formState.sql_type = prefill.sql_type || ''
    if (prefill.cc) formState.cc = prefill.cc
    if (prefill.content) codemirrorRef.value?.setContent(prefill.content)
  }
})
</script>

<style scoped>
.order-create-card :deep(.ant-card-body) {
  padding-top: 14px;
}

.title-wrap {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.title-main {
  font-size: 17px;
  font-weight: 600;
  color: #1f2d38;
}

.title-sub {
  font-size: 12px;
  color: #71818c;
}

.create-row {
  align-items: stretch;
}

.pane-col {
  display: flex;
}

.pane-col > .panel {
  width: 100%;
}

.panel {
  height: 100%;
  padding: 14px;
  border-radius: 12px;
  border: 1px solid #e5edf4;
  background: #ffffff;
  box-shadow: 0 6px 18px rgba(16, 24, 40, 0.03);
}

.form-tip {
  margin-bottom: 14px;
}

.form-panel :deep(.ant-form-item) {
  margin-bottom: 14px;
}

.form-panel :deep(.ant-form-item-label > label) {
  font-size: 13px;
  font-weight: 500;
}

.form-panel :deep(.ant-input),
.form-panel :deep(.ant-input-affix-wrapper),
.form-panel :deep(.ant-select-selector) {
  border-radius: 8px;
}

.form-actions {
  margin-top: 4px;
  margin-bottom: 0;
}

.editor-panel {
  display: flex;
  flex-direction: column;
}

.editor-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 12px;
}

.editor-title {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.editor-title-main {
  font-size: 15px;
  font-weight: 600;
  color: #1f2d38;
}

.editor-title-sub {
  font-size: 12px;
  color: #71818c;
}

.editor-shell {
  flex: 1 1 auto;
  min-height: 0;
  border-radius: 10px;
  border: 1px solid #e5edf4;
  overflow: hidden;
}

.inspect-wrapper {
  margin-top: 14px;
}

@media (max-width: 768px) {
  .panel {
    padding: 12px;
  }
}
</style>
