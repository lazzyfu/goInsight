<template>
  <a-modal :open="open" title="编辑审核参数" width="50%" :footer="null" @cancel="handleCancel">
    <a-form
      ref="formRef"
      :label-col="{ span: 4 }"
      :wrapper-col="{ span: 18 }"
      :model="localFormState"
      :rules="rules"
      @finish="onSubmit"
    >
      <a-form-item label="描述" name="remark" has-feedback>
        <a-input
          disabled
          v-model:value="localFormState.remark"
          placeholder="请输入备注"
          allow-clear
        />
      </a-form-item>
      <a-form-item
        label="审核参数"
        name="params"
        has-feedback
        help="格式要求为JSON类型，默认为{}，表示继承全局审核参数"
      >
        <a-textarea
          :rows="5"
          v-model:value="localFormState.params"
          placeholder=" 请输入自定义审核参数，默认为{}"
        />
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
import { reactive, ref, watch } from 'vue'

const emit = defineEmits(['update:open', 'submit'])
const props = defineProps({
  open: Boolean,
  title: String,
  formState: Object,
})

const formRef = ref()

// formState父组件传值，子组件修改，需要重新赋值
const localFormState = reactive({ ...props.formState })

watch(
  () => props.formState,
  (newVal) => {
    Object.assign(localFormState, newVal)
  },
  { immediate: true, deep: true },
)

const rules = {
  name: [
    { required: true, message: '请输入角色名', trigger: 'blur' },
    { min: 2, max: 32, message: '角色名长度应在2~32个字符', trigger: 'blur' },
    {
      pattern: /^[a-zA-Z0-9\u4e00-\u9fa5_]+$/,
      message: '角色名只能包含字母、数字、中文或下划线',
      trigger: 'blur',
    },
  ],
}

const handleCancel = () => {
  emit('update:open', false)
  formRef.value?.resetFields()
}

const onSubmit = () => {
  emit('submit', localFormState)
}
</script>
