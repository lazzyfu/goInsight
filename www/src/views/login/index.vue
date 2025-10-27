<template>
  <div class="login-box">
    <div class="login-logo">
      <h1 class="mb-0 ml-2 text-3xl font-bold">Ark 运维平台</h1>
    </div>
    <a-form layout="horizontal" :model="data.formInline" @submit.prevent="handleSubmit">
      <a-form-item>
        <a-input v-model:value="data.formInline.username" size="large" placeholder="admin">
          <template #prefix><user-outlined type="user" /></template>
        </a-input>
      </a-form-item>
      <a-form-item>
        <a-input v-model:value="data.formInline.password" size="large" type="password" placeholder="a123456"
          autocomplete="password">
          <template #prefix><lock-outlined type="user" /></template>
        </a-input>
      </a-form-item>
      <a-form-item>
        <a-button type="primary" html-type="submit" size="large" :loading="data.loading" block>
          登录
        </a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup>
import { Login } from '@/api/login'
import { useUserStore } from '@/store/user'
import { LockOutlined, UserOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { reactive } from 'vue'
import { useRouter } from 'vue-router'

const userStore = useUserStore()

const data = reactive({
  loading: false,
  formInline: {
    username: 'admin',
    password: 'admin',
  },
})

const router = useRouter()

const handleSubmit = async () => {
  const { username, password } = data.formInline
  if (username.trim() == '' || password.trim() == '') {
    return message.warning('用户名或密码不能为空！')
  }
  message.loading('登录中...', 0)
  data.loading = true

  Login(data.formInline).then(res => {
    if (res.code === '0000') {
      localStorage.setItem("onLine", 1)
      userStore.setUserToken(res.data.token)
      message.success('登录成功！')
      router.push({ name: 'Root' })
    } else {
      message.warning(res.message || '登录失败，请检查用户名和密码！')
    }
  })

  data.loading = false
  message.destroy()
}
</script>

<style lang="less" scoped>
.login-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100vw;
  height: 100vh;
  background: #f0f2f5 url(@/assets/background.svg);
  background-size: 100%;

  .login-logo {
    display: flex;
    align-items: center;
    margin-top: 200px;
    // margin-bottom: 30px;

    .svg-icon {
      font-size: 48px;
    }
  }

  .desc {
    font-size: 14px;
    color: rgba(0, 0, 0, 0.45);
    margin-bottom: 60px;
  }

  :deep(.ant-form) {
    width: 400px;

    .ant-col {
      width: 100%;
    }

    .ant-form-item-label {
      padding-right: 6px;
    }
  }
}
</style>
