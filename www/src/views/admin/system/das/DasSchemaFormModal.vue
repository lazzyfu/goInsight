<template>
  <a-modal :open="open" title="数据访问权限" :footer="null" @cancel="handleCancel">
    <a-form
      ref="formRef"
      :label-col="{ span: 4 }"
      :wrapper-col="{ span: 18 }"
      :model="formState"
      :rules="rules"
      @finish="onSubmit"
    >
      <a-form-item label="用户" name="username" has-feedback>
        <a-select
          ref="select"
          v-model:value="formState.username"
          :options="state.users"
          :field-names="{ label: 'nick_name', value: 'username' }"
          allowClear
        >
        </a-select>
      </a-form-item>

      <a-form-item label="类型" name="db_type" has-feedback>
        <a-select
          ref="select"
          v-model:value="formState.db_type"
          :options="[
            { value: 'MySQL', label: 'MySQL' },
            { value: 'TiDB', label: 'TiDB' },
            { value: 'ClickHouse', label: 'ClickHouse' },
          ]"
          allowClear
        >
        </a-select>
      </a-form-item>

      <a-form-item label="环境" name="environment" has-feedback>
        <a-select
          ref="select"
          v-model:value="formState.environment"
          :options="state.environments"
          :field-names="{ label: 'name', value: 'id' }"
          @change="changeEnv"
          allowClear
        ></a-select>
      </a-form-item>

      <a-form-item label="实例" name="instance_id" has-feedback>
        <a-select
          ref="select"
          v-model:value="formState.instance_id"
          :options="state.instances"
          :field-names="{ label: 'remark', value: 'instance_id' }"
          @change="changeInstance"
          allowClear
        ></a-select>
      </a-form-item>

      <a-form-item label="库名" name="schema" has-feedback>
        <a-select
          ref="select"
          v-model:value="formState.schema"
          :options="state.schemas"
          :field-names="{ label: 'schema', value: 'schema' }"
          allowClear
        ></a-select>
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
import {
  getEnvironmentsApi,
  getInstancesListApi,
  getSchemasListApi,
  getUsersApi,
} from '@/api/admin'
import { onMounted, reactive, ref } from 'vue'

const emit = defineEmits(['update:open', 'submit'])
const props = defineProps({
  open: Boolean,
})

const formRef = ref()

const state = reactive({
  users: [],
  instances: [],
  schemas: [],
})

// formState父组件传值，子组件修改，需要重新赋值
const formState = reactive({
  username: '',
  db_type: '',
  environment: '',
  instance_id: '',
  schema: '',
})

const rules = {
  username: [{ required: true, message: '请选择用户', trigger: 'blur' }],
  db_type: [{ required: true, message: '请选择数据库类型', trigger: 'blur' }],
  environment: [{ required: true, message: '请选择环境', trigger: 'blur' }],
  instance: [{ required: true, message: '请选择实例', trigger: 'blur' }],
  schema: [{ required: true, message: '请选择库名', trigger: 'blur' }],
}

const getUsers = async () => {
  const res = await getUsersApi().catch(() => {})
  state.users = res.data || []
}

const getEnvironments = async () => {
  const res = await getEnvironmentsApi().catch(() => {})
  state.environments = res.data || []
}

const changeEnv = async (value) => {
  const payload = {
    id: value,
    db_type: formState.db_type,
    is_page: false,
  }
  const res = await getInstancesListApi(payload).catch(() => {})
  state.instances = res.data || []
}

const changeInstance = async (value) => {
  const payload = {
    instance_id: value,
    is_page: false,
  }
  const res = await getSchemasListApi(payload).catch(() => {})
  state.schemas = res.data || []
}

const handleCancel = () => {
  emit('update:open', false)
  formRef.value?.resetFields()
}

const onSubmit = () => {
  emit('submit', formState)
}

onMounted(async () => {
  await getUsers()
  await getEnvironments()
})
</script>
