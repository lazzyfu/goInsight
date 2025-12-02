<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h1>Ark 运维平台</h1>
      </div>

      <a-form ref="formRef" :model="formState" :rules="rules" layout="vertical" @finish="onSubmit">
        <a-form-item v-if="!uiState.showOTP" label="用户名" name="username">
          <a-input
            v-model:value="formState.username"
            size="large"
            placeholder="请输入用户名"
            :disabled="uiState.loading"
          />
        </a-form-item>

        <a-form-item v-if="!uiState.showOTP" label="密码" name="password">
          <a-input
            v-model:value="formState.password"
            size="large"
            type="password"
            placeholder="请输入密码"
            autocomplete="password"
            :disabled="uiState.loading"
          />
        </a-form-item>

        <a-form-item v-if="uiState.showOTP" label="OTP 验证码" name="otp_code">
          <a-input
            v-model:value="formState.otp_code"
            size="large"
            placeholder="请输入6位OTP验证码"
            :maxlength="6"
            :disabled="uiState.loading"
          />
        </a-form-item>

        <a-form-item>
          <a-button type="primary" html-type="submit" size="large" :loading="uiState.loading" block>
            {{ uiState.loading ? '登录中...' : '登录' }}
          </a-button>
        </a-form-item>
      </a-form>
    </div>

    <!-- 绑定 OTP 子组件 -->
    <BindOTPModal v-model:open="uiState.bindOtpModalOpen" ref="otpModalRef" />
  </div>
</template>

<script setup>
import { Login } from '@/api/login'
import { useUserStore } from '@/store/user'
import { message } from 'ant-design-vue'
import { reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import BindOTPModal from './OTP.vue'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref(null)
const otpModalRef = ref(null)

const uiState = reactive({
  loading: false,
  showOTP: false,
  bindOtpModalOpen: false,
})

const formState = reactive({
  username: '',
  password: '',
  otp_code: '',
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  otp_code: [
    { required: true, message: '请输入OTP验证码', trigger: 'blur' },
    { pattern: /^\d{6}$/, message: '请输入6位数字验证码', trigger: 'blur' },
  ],
}

// 监听 OTP 显示状态，动态管理验证规则
watch(
  () => uiState.showOTP,
  (show) => {
    if (show && formRef.value) {
      formRef.value.clearValidate('otp_code')
    }
  },
)

// 提交登录
const onSubmit = async () => {
  try {
    uiState.loading = true
    const res = await Login(formState).catch((err) => {})

    if (res?.code === '0000') {
      // 登录成功
      localStorage.setItem('onLine', '1')
      userStore.setUserToken(res.data.token)
      message.success('登录成功')
      router.push({ name: 'Root' })
    } else if (res?.code === '4001') {
      // 需要绑定 OTP
      message.warning(res.message || '需要绑定OTP')
      uiState.bindOtpModalOpen = true
      otpModalRef.value?.show(formState)
    } else if (res?.code === '4002') {
      // 需要输入 OTP
      uiState.showOTP = true
      message.info(res.message || '请输入 OTP 验证码')
    }
  } finally {
    uiState.loading = false
  }
}
</script>

<style lang="less" scoped>
.login-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: #f0f2f5 url(@/assets/background.svg);
  background-size: 100%;
  padding: 20px;
}

.login-card {
  width: 100%;
  max-width: 420px;
  padding: 48px 40px;
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
}

.login-header {
  text-align: center;
  margin-bottom: 40px;

  h1 {
    margin: 0 0 8px 0;
    font-size: 28px;
    font-weight: 600;
    color: #1f2937;
  }

  p {
    margin: 0;
    font-size: 14px;
    color: #6b7280;
  }
}

:deep(.ant-form) {
  .ant-form-item {
    margin-bottom: 24px;
  }

  .ant-form-item-label > label {
    font-size: 14px;
    font-weight: 500;
    color: #374151;
  }

  .ant-input {
    border-radius: 8px;
    border: 1px solid #e5e7eb;
    transition: all 0.3s;

    &:hover {
      border-color: #9ca3af;
    }

    &:focus,
    &.ant-input-focused {
      border-color: #667eea;
      box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
    }
  }

  .ant-btn-primary {
    height: 44px;
    border-radius: 8px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border: none;
    font-size: 16px;
    font-weight: 500;
    transition: all 0.3s;

    &:hover:not(:disabled) {
      transform: translateY(-2px);
      box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
    }

    &:active:not(:disabled) {
      transform: translateY(0);
    }
  }
}

// 响应式设计
@media (max-width: 576px) {
  .login-card {
    padding: 32px 24px;
  }

  .login-header h1 {
    font-size: 24px;
  }
}
</style>
