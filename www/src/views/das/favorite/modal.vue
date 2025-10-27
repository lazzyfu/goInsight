<template>
  <a-modal :open="open" title="收藏SQL语句" width="45%" :footer="null" @cancel="handleCancel">
    <a-form layout="vertical" :model="formState" :rules="rules" @finish="onSubmit">
      <a-form-item label="标题" name="title">
        <a-input v-model:value="formState.title" placeholder="为这条SQL起个名字" allow-clear />
      </a-form-item>
      <a-form-item label="SQL内容" name="sqltext">
        <a-textarea
          v-model:value="formState.sqltext"
          :rows="10"
          placeholder="请输入SQL"
          allow-clear
        />
      </a-form-item>
      <a-form-item style="text-align: right">
        <a-button @click="handleCancel">取消</a-button>
        <a-button type="primary" html-type="submit" style="margin-left: 10px">
          <template #icon>
            <StarOutlined />
          </template>
          {{ btnType }}
        </a-button>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { StarOutlined } from '@ant-design/icons-vue'

const emit = defineEmits(['update:open', 'submit'])

const rules = {
  title: [
    {
      required: true,
      message: '为这条SQL起个名字',
      trigger: 'change',
    },
    {
      min: 3,
      max: 256,
      message: '最少3个字符，最多256个字符',
      trigger: 'blur',
    },
  ],
  sqltext: [
    {
      required: true,
      message: '请输入SQL',
      trigger: 'change',
    },
    {
      min: 5,
      message: '最少5个字符',
      trigger: 'blur',
    },
  ],
}

const props = defineProps({
  open: Boolean,
  formState: {
    type: Object,
    default: () => ({
      title: '',
      sqltext: '',
    }),
  },
  btnType: '',
})

const handleCancel = () => {
  emit('update:open', false)
}

const onSubmit = () => {
  emit('submit', props.formState)
  emit('update:open', false)
}
</script>
