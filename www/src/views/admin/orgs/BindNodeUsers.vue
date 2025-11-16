<template>
  <a-modal :open="open" title="绑定用户" :footer="null" @cancel="handleCancel">
    <a-form ref="formRef" :model="formState" @finish="onSubmit">
      <a-form-item label="用户" name="users" has-feedback>
        <a-select
          ref="select"
          mode="multiple"
          v-model:value="formState.users"
          :options="users"
          :field-names="{ label: 'nick_name', value: 'uid' }"
          allowClear
        >
        </a-select>
      </a-form-item>

      <a-form-item style="text-align: right">
        <a-button @click="handleCancel">取消</a-button>
        <a-button type="primary" html-type="submit" style="margin-left: 10px">确定</a-button>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { getUsersApi } from '@/api/admin'
import { onMounted, reactive, ref } from 'vue'

const emit = defineEmits(['update:open', 'submit'])
const props = defineProps({
  open: Boolean,
  nodeKey: String,
})

const formRef = ref()
const users = ref([])
const formState = reactive({
  users: [],
})

const handleCancel = () => {
  emit('update:open', false)
  formRef.value?.resetFields()
}

const onSubmit = () => {
  const payload = {
    key: props.nodeKey,
    ...formState,
  }
  emit('submit', payload)
}

onMounted(async () => {
  const res = await getUsersApi().catch(() => {})
  users.value = res.data || []
})
</script>
