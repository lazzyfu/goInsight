<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <div class="logo-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" />
          </svg>
        </div>
        <h1>Ark 运维平台</h1>
        <p>安全登录以访问控制台</p>
      </div>

      <a-form
        ref="formRef"
        :model="formState"
        autocomplete="off"
        :rules="rules"
        layout="vertical"
        @finish="onSubmit"
      >
        <a-form-item v-if="!uiState.showOTP" name="username">
          <a-input v-model:value="formState.username" size="large" placeholder="请输入用户名">
            <template #prefix>
              <UserOutlined class="input-icon" />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item v-if="!uiState.showOTP" name="password">
          <a-input-password
            v-model:value="formState.password"
            size="large"
            placeholder="请输入密码"
          >
            <template #prefix>
              <LockOutlined class="input-icon" />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item v-if="uiState.showOTP" name="otp_code">
          <a-input
            v-model:value="formState.otp_code"
            size="large"
            placeholder="请输入6位OTP验证码"
            :maxlength="6"
          >
            <template #prefix>
              <SafetyOutlined class="input-icon" />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item class="submit-item">
          <a-button type="primary" html-type="submit" size="large" :loading="uiState.loading" block>
            {{ uiState.loading ? '登录中...' : '登录' }}
          </a-button>
        </a-form-item>
      </a-form>

      <div class="login-footer">
        <span>© 2025 Ark Platform</span>
      </div>
    </div>
  </div>
  <!-- 绑定 OTP 子组件 -->
  <BindOTPModal v-model:open="uiState.bindOtpModalOpen" ref="otpModalRef" />
</template>

<script setup>
defineOptions({ name: 'UserLogin' })

import { Login } from '@/api/login'
import { useUserStore } from '@/store/user'
import { LockOutlined, SafetyOutlined, UserOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import BindOTPModal from './OTP.vue'

const formRef = ref(null)
const router = useRouter()
const userStore = useUserStore()
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

watch(
  () => uiState.showOTP,
  (show) => {
    if (show && formRef.value) {
      formRef.value.clearValidate('otp_code')
    }
  },
)

const onSubmit = async () => {
  try {
    uiState.loading = true
    const res = await Login(formState).catch(() => {})

    if (res?.code === '0000') {
      // 登录成功
      localStorage.setItem('onLine', '1')
      userStore.setToken(res.data.token)
      message.success('登录成功')
      router.push('/')
    } else if (res?.code === '4001') {
      // 需要绑定 OTP
      message.warning(res.message || '需要绑定OTP')
      uiState.bindOtpModalOpen = true
      otpModalRef.value?.show(formState)
    } else if (res?.code === '4002') {
      // 需要输入 OTP
      uiState.showOTP = true
    }
  } finally {
    uiState.loading = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(145deg, #f5f7fa 0%, #e4e8ec 100%);
  padding: 20px;
}

.login-card {
  width: 100%;
  max-width: 400px;
  padding: 40px 36px;
  background: #ffffff;
  border-radius: 16px;
  box-shadow:
    0 4px 6px -1px rgba(0, 0, 0, 0.05),
    0 10px 15px -3px rgba(0, 0, 0, 0.08);
}

.login-header {
  text-align: center;
  margin-bottom: 36px;
}

.logo-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 16px;
  padding: 12px;
  background: #f0f5ff;
  border-radius: 12px;
  color: #4f6ef7;
}

.logo-icon svg {
  width: 100%;
  height: 100%;
}

.login-header h1 {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
  color: #1a1a2e;
  letter-spacing: -0.5px;
}

.login-header p {
  margin: 0;
  font-size: 14px;
  color: #8b8fa3;
}

.input-icon {
  color: #b0b5c3;
  font-size: 16px;
}

:deep(.ant-form-item) {
  margin-bottom: 20px;
}

:deep(.ant-input-affix-wrapper) {
  display: flex;
  align-items: center;
  height: 48px;
  padding: 0 14px;
  border-radius: 10px;
  border: 1.5px solid #e8eaef;
  background: #fafbfc;
  transition: all 0.2s ease;
}

:deep(.ant-input-affix-wrapper .ant-input-prefix) {
  display: flex;
  align-items: center;
  margin-inline-end: 10px;
}

:deep(.ant-input-affix-wrapper .ant-input) {
  height: auto;
  padding: 0;
  border: none;
  background: transparent;
}

:deep(.ant-input:not(.ant-input-affix-wrapper .ant-input)) {
  height: 48px;
  border-radius: 10px;
  border: 1.5px solid #e8eaef;
  background: #fafbfc;
  transition: all 0.2s ease;
}

:deep(.ant-input-affix-wrapper:hover),
:deep(.ant-input:hover) {
  border-color: #c5cad6;
  background: #ffffff;
}

:deep(.ant-input-affix-wrapper-focused),
:deep(.ant-input:focus),
:deep(.ant-input-affix-wrapper:focus) {
  border-color: #4f6ef7;
  background: #ffffff;
  box-shadow: 0 0 0 3px rgba(79, 110, 247, 0.08);
}

:deep(.ant-input) {
  font-size: 15px;
  background: transparent;
}

:deep(.ant-input::placeholder) {
  color: #b0b5c3;
}

.submit-item {
  margin-top: 28px;
  margin-bottom: 0;
}

:deep(.ant-btn-primary) {
  height: 48px;
  border-radius: 10px;
  background: #4f6ef7;
  border: none;
  font-size: 15px;
  font-weight: 500;
  box-shadow: 0 2px 8px rgba(79, 110, 247, 0.25);
  transition: all 0.2s ease;
}

:deep(.ant-btn-primary:hover:not(:disabled)) {
  background: #3d5ce8;
  box-shadow: 0 4px 12px rgba(79, 110, 247, 0.35);
}

:deep(.ant-btn-primary:active:not(:disabled)) {
  background: #3451d1;
}

.login-footer {
  margin-top: 32px;
  text-align: center;
  font-size: 12px;
  color: #b0b5c3;
}

@media (max-width: 480px) {
  .login-card {
    padding: 32px 24px;
  }

  .login-header h1 {
    font-size: 22px;
  }
}
</style>
