<template>
  <a-modal :open="open" title="绑定用户" :footer="null" @cancel="handleCancel">
    <a-form
      ref="formRef"
      :label-col="{ span: 4 }"
      :wrapper-col="{ span: 18 }"
      :model="formState"
      @finish="onSubmit"
    >
      <a-form-item label="用户" name="users" has-feedback>
        <a-select
          ref="select"
          mode="multiple"
          v-model:value="formState.users"
          :options="props.users"
          :field-names="{ label: 'nick_name', value: 'uid' }"
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
import { reactive, ref } from 'vue'

const emit = defineEmits(['update:open', 'submit'])
const props = defineProps({
  open: Boolean,
  nodeKey: String,
  users: Array,
})

const formRef = ref()
const formState = reactive({
  users: [],
})

// 关闭弹窗
const handleCancel = () => {
  emit('update:open', false)
  formRef.value?.resetFields()
}

// 提交表单
const onSubmit = async () => {
  const payload = {
    key: props.nodeKey,
    ...formState,
  }
  emit('submit', payload)
}
</script>
