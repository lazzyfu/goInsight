<template>
  <a-card size="small">
    <a-row :gutter="{ xs: 8, sm: 16, md: 24, lg: 32 }">
      <a-col class="gutter-row" :span="8">
        <a-form
          ref="formRef"
          :model="formState"
          name="basic"
          :label-col="{ span: 4 }"
          :wrapper-col="{ span: 20 }"
          autocomplete="off"
          @finish="onSubmit"
        >
          <a-form-item
            label="标题"
            name="title"
            has-feedback
            :rules="[{ required: true, min: 3, max: 196 }]"
          >
            <a-input v-model:value="formState.title" />
          </a-form-item>
          <a-form-item
            label="备注"
            name="remark"
            has-feedback
            :rules="[{ required: true, min: 3, max: 1024 }]"
          >
            <a-textarea v-model:value="formState.remark" placeholder=" 请输入工单需求或备注" />
          </a-form-item>
          <a-form-item label="类型" name="sql_type" has-feedback :rules="[{ required: true }]">
            <a-select
              ref="select"
              v-model:value="formState.sql_type"
              :options="[
                { value: 'DML', label: 'DML' },
                { value: 'DDL', label: 'DDL' },
              ]"
              allowClear
            >
            </a-select>
          </a-form-item>
          <a-form-item label="环境" name="environment" has-feedback :rules="[{ required: true }]">
            <a-select
              ref="select"
              v-model:value="formState.environment"
              :options="environments"
              :field-names="{ label: 'name', value: 'id' }"
              @change="handleChangeEnv"
              allowClear
            ></a-select>
          </a-form-item>
          <a-form-item label="数据库" name="db_type" has-feedback :rules="[{ required: true }]">
            <a-select
              ref="select"
              v-model:value="formState.db_type"
              :options="dbTypes"
              @change="handleChangeDbType"
              :disabled="!formState.environment"
              allowClear
            ></a-select>
          </a-form-item>
          <a-form-item label="实例" name="instance_id" has-feedback :rules="[{ required: true }]">
            <a-select
              ref="select"
              v-model:value="formState.instance_id"
              :options="instances"
              :field-names="{ label: 'remark', value: 'instance_id' }"
              @change="handleChangeIns"
              :disabled="!formState.db_type"
              allowClear
            >
            </a-select>
          </a-form-item>
          <a-form-item label="库名" name="schema" has-feedback :rules="[{ required: true }]">
            <a-select
              ref="select"
              v-model:value="formState.schema"
              :options="schemas"
              :field-names="{ label: 'schema', value: 'schema' }"
              :disabled="!formState.instance_id"
              allowClear
            >
            </a-select>
          </a-form-item>

          <a-form-item label="抄送" name="cc" has-feedback>
            <a-select
              ref="select"
              mode="multiple"
              v-model:value="formState.cc"
              :options="users"
              :field-names="{ label: 'nick_name', value: 'username' }"
              allowClear
            >
            </a-select>
          </a-form-item>

          <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
            <a-button type="primary" html-type="submit">提交</a-button>
            <a-button style="margin-left: 10px" @click="resetForm">重置</a-button>
          </a-form-item>
        </a-form>
      </a-col>
      <a-col class="gutter-row" :span="16">
        <a-space size="small">
          <a-button @click="formatSqlContent()">
            <template #icon>
              <CodeOutlined />
            </template>
            格式化
          </a-button>
          <a-button @click="checkSyntax()">
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
const codemirrorRef = ref(null)
const inspectResultTableRef = ref(null)

// 表单数据
const formState = reactive({
  title: '',
  remark: '',
  environment: '',
  db_type: '',
  sql_type: '',
  instance_id: '',
  schema: '',
  cc: [],
  content: '',
  export_file_format: 'CSV',
})
const formRef = ref()

// 选项数据
const dbTypes = ref([
  { value: 'MySQL', label: 'MySQL' },
  { value: 'TiDB', label: 'TiDB' },
])

const environments = ref([])
const instances = ref([])
const schemas = ref([])
const users = ref([])

const handleChangeEnv = (value) => {
  formState.db_type = ''
  formState.instance_id = ''
  formState.schema = ''
  instances.value = []
  schemas.value = []
}

const handleChangeDbType = debounce(async (value) => {
  formState.instance_id = ''
  formState.schema = ''
  instances.value = []
  schemas.value = []
  if (!value) return

  const res = await getOrderInstancesApi({
    id: formState.environment,
    db_type: formState.db_type,
    is_page: false,
  }).catch((err) => {})
  if (res) {
    instances.value = res?.data || []
  }
})

const handleChangeIns = debounce(async (value) => {
  formState.schema = ''
  schemas.value = []
  if (!value) return
  const res = await getOrderSchemasApi({ instance_id: value, is_page: false })
  if (res) {
    schemas.value = res.data || []
  }
})

// 提交工单
const onSubmit = debounce(async () => {
  formState.content = codemirrorRef.value.getContent()
  if (formState.content.length == 0) {
    message.warning('SQL内容不能为空')
    return
  }

  const res = await createOrderApi(formState).catch((err) => {})
  if (res) {
    message.success('工单创建成功')
    router.push({ name: 'orders.list' })
  }
})

const resetForm = () => {
  formRef.value.resetFields()
}

// 格式化
const formatSqlContent = () => {
  codemirrorRef.value.formatContent()
}

// 检查语法
const checkSyntax = debounce(async () => {
  const sqltext = codemirrorRef.value.getContent()
  if (sqltext.length == 0) {
    message.warning('SQL内容不能为空')
    return
  }

  if (!formState.environment || !formState.db_type || !formState.instance_id || !formState.schema) {
    message.warning('请先选择环境、数据库类型、实例、库名')
    return
  }

  const data = {
    db_type: formState.db_type,
    sql_type: formState.sql_type,
    instance_id: formState.instance_id,
    schema: formState.schema,
    content: sqltext,
  }

  const res = await inspectOrderSyntaxApi(data).catch((err) => {})
  if (res) {
    inspectResultTableRef.value.render(res.data)
  }
})

onMounted(async () => {
  const [envRes, userRes] = await Promise.all([
    getOrderEnvironmentsApi(),
    getOrderUsersApi({ is_page: false }),
  ])

  environments.value = envRes.data || []
  users.value = userRes.data || []
})
</script>
