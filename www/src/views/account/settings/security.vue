<template>
  <a-list item-layout="horizontal" :data-source="items">
    <template #renderItem="{ item }">
      <a-list-item>
        <a-list-item-meta>
          <template #title>
            <a>{{ item.title }}</a>
          </template>
        </a-list-item-meta>
        <template v-if="item.actions">
          <a slot="actions" @click="item.actions.callback">{{ item.actions.title }}</a>
        </template>
      </a-list-item>
    </template>

    <a-modal title="修改密码" v-model:open="open" width="50%" :footer="null">
      <a-form ref="formRef" :model="formState" :rules="rules" @finish="handleFinish" :label-col="{ span: 4 }"
        :wrapper-col="{ span: 18 }" autocomplete="off">
        <a-form-item label="当前密码" has-feedback name="old_password">
          <a-input v-model:value="formState.old_password" type="password" autocomplete="off">
            <LockOutlined />
          </a-input>
        </a-form-item>
        <a-form-item label="新密码" has-feedback name="new_password">
          <a-input v-model:value="formState.new_password" autocomplete="off" type="password">
          </a-input>
        </a-form-item>
        <a-form-item label="确认密码" has-feedback name="confirm_password">
          <a-input v-model:value="formState.confirm_password" autocomplete="off" type="password">
          </a-input>
        </a-form-item>
        <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
          <a-button type="primary" html-type="submit">提交</a-button>
          <a-button style="margin-left: 10px" @click="resetForm">重置</a-button>
        </a-form-item>
      </a-form>
    </a-modal>
  </a-list>
</template>

<script setup>
import { ChangePasswordApi } from '@/api/profile';
import { regPassword } from '@/utils/validate';
import { LockOutlined } from '@ant-design/icons-vue';
import { message } from 'ant-design-vue';
import { reactive, ref } from 'vue';

const open = ref(false)
const loading = ref(false)
const formRef = ref();
const formState = reactive({
  old_password: '',
  new_password: '',
  confirm_password: '',
})

const items = [
  {
    title: '账户密码',
    actions: {
      title: '修改',
      callback: () => {
        showModal()
      },
    },
  },
]

const validateOldPass = async (_rule, value) => {
  if (!value) return Promise.reject('请输入密码');
  if (value.length < 1) return Promise.reject('至少1个字符');
}

const validateNewPass = async (_rule, value) => {
  if (!value) return Promise.reject('请输入密码');
  if (!regPassword.test(value)) return Promise.reject('密码至少7个字符,必须包含大写字母、小写字母、数字和特殊字符');
  if (formState.confirm_password) formRef.value.validateFields('confirm_password');
  return Promise.resolve();
};

const validateVerifyPass = async (_rule, value) => {
  if (!value) return Promise.reject('请输入密码');
  if (!regPassword.test(value)) return Promise.reject('密码至少7个字符,必须包含大写字母、小写字母、数字和特殊字符');
  if (value !== formState.new_password) return Promise.reject('两次输入的密码不一致');
  return Promise.resolve();
};

const rules = {
  old_password: [
    {
      required: true,
      validator: validateOldPass,
      trigger: 'change',
    },
  ],
  new_password: [
    {
      required: true,
      validator: validateNewPass,
      trigger: 'change',
    },
  ],
  confirm_password: [
    {
      required: true,
      validator: validateVerifyPass,
      trigger: 'change',
    },
  ],
};

const showModal = () => {
  open.value = true
}

const handleCancel = (e) => {
  open.value = false
}

const handleFinish = values => {
  loading.value = true
  ChangePasswordApi(values).then((res) => {
    if (res.code === '0000') {
      message.info(res.message)
      handleCancel()
      location.reload() // 刷新页面，此时token已经过期，需要重新登录
      loading.value = false
    } else {
      message.error(res.message)
      loading.value = false
    }
  })
  loading.value = false
};
const resetForm = () => {
  formRef.value.resetFields();
};
</script>
