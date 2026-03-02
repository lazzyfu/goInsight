<template>
  <div class="login-page">
    <div class="login-shell">
      <section class="visual-panel" :style="visualStyle">
        <div class="visual-overlay"></div>
        <div class="visual-content">
          <p class="visual-tag">goInsight Database Platform</p>
          <h2>让数据库变更进入可审批、可追踪、可审计的闭环</h2>
          <p class="visual-desc">
            从申请、审批、执行到回溯，统一管理 DDL、DML 与导出任务，提升团队协作效率与变更安全性。
          </p>
          <ul class="visual-points">
            <li>
              <SafetyOutlined />
              工单全流程审批与权限控制
            </li>
            <li>
              <DeploymentUnitOutlined />
              多组织多角色隔离治理
            </li>
            <li>
              <AuditOutlined />
              SQL 语法检查与操作审计
            </li>
          </ul>
        </div>
      </section>
      <section class="form-panel">
        <div class="login-card">
          <div class="login-header">
            <div class="logo-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" />
              </svg>
            </div>
            <h1>{{ appTitle }}</h1>
            <p>请输入账号信息登录控制台</p>
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
            <span>© 2025 goInsight</span>
          </div>
        </div>
      </section>
    </div>
  </div>
  <!-- 绑定 OTP 子组件 -->
  <BindOTPModal v-model:open="uiState.bindOtpModalOpen" ref="otpModalRef" />
</template>

<script setup>
defineOptions({ name: 'UserLogin' })

import { Login } from '@/api/login'
import { useUserStore } from '@/store/user'
import loginVisual from '@/assets/original2.png'
import {
  AuditOutlined,
  DeploymentUnitOutlined,
  LockOutlined,
  SafetyOutlined,
  UserOutlined,
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import BindOTPModal from './OTP.vue'

const appTitle = import.meta.env.VITE_APP_TITLE || '数据库工单平台'
const visualStyle = {
  backgroundImage: `linear-gradient(120deg, rgba(8, 24, 31, 0.72), rgba(9, 30, 42, 0.82)), url(${loginVisual})`,
}

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
.login-page {
  --ink-950: #0f1720;
  --ink-700: #243947;
  --ink-500: #4f6470;
  --line-200: #d8e1e8;
  --surface-100: #f5f8fa;
  --primary-600: #0f766e;
  --primary-700: #0b5f58;
  --primary-050: rgba(15, 118, 110, 0.12);

  min-height: 100vh;
  width: 100vw;
  padding: 0;
  overflow: hidden;
  background:
    radial-gradient(circle at 10% 20%, #dcebf4 0%, transparent 45%),
    radial-gradient(circle at 90% 85%, #d9efe8 0%, transparent 40%),
    linear-gradient(140deg, #e9f0f4 0%, #eef3f4 52%, #f4f7f7 100%);
}

.login-shell {
  width: 100vw;
  max-width: none;
  min-height: 100vh;
  height: 100vh;
  display: grid;
  grid-template-columns: 6fr 4fr;
  border-radius: 0;
  overflow: hidden;
  background: #fff;
  box-shadow: none;
}

.visual-panel {
  position: relative;
  background-size: cover;
  background-position: center;
  padding: 56px 56px 48px;
  display: flex;
  align-items: flex-end;
}

.visual-overlay {
  position: absolute;
  inset: 0;
  background:
    linear-gradient(0deg, rgba(0, 0, 0, 0.25), transparent 32%),
    radial-gradient(circle at 25% 25%, rgba(154, 196, 214, 0.26), transparent 52%);
}

.visual-content {
  position: relative;
  z-index: 1;
  max-width: 500px;
  color: #eff8fb;
}

.visual-tag {
  display: inline-flex;
  padding: 6px 12px;
  margin-bottom: 18px;
  font-size: 12px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  border: 1px solid rgba(239, 248, 251, 0.38);
  border-radius: 999px;
  color: #f6fbfd;
  background: rgba(255, 255, 255, 0.08);
}

.visual-content h2 {
  margin: 0;
  font-size: 34px;
  line-height: 1.22;
  font-weight: 700;
}

.visual-desc {
  margin: 18px 0 24px;
  font-size: 15px;
  line-height: 1.7;
  color: rgba(239, 248, 251, 0.9);
}

.visual-points {
  margin: 0;
  padding: 0;
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.visual-points li {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  color: rgba(239, 248, 251, 0.96);
}

.form-panel {
  background: linear-gradient(180deg, #fbfdfd 0%, #f4f7f8 100%);
  padding: 36px 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow-y: auto;
}

.login-card {
  width: min(440px, 100%);
  padding: 12px 10px;
  background: transparent;
  border: none;
  border-radius: 0;
  backdrop-filter: none;
  box-shadow: none;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 16px;
  padding: 12px;
  background: linear-gradient(135deg, #e4f4f1 0%, #c9ebe4 100%);
  border-radius: 12px;
  color: var(--primary-600);
}

.logo-icon svg {
  width: 100%;
  height: 100%;
}

.login-header h1 {
  margin: 0 0 8px 0;
  font-size: 26px;
  font-weight: 700;
  color: var(--ink-950);
  letter-spacing: -0.02em;
}

.login-header p {
  margin: 0;
  font-size: 13px;
  color: var(--ink-500);
}

.input-icon {
  color: #8ea2b0;
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
  border: 1.5px solid var(--line-200);
  background: #f9fbfc;
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
  border: 1.5px solid var(--line-200);
  background: #f9fbfc;
  transition: all 0.2s ease;
}

:deep(.ant-input-affix-wrapper:hover),
:deep(.ant-input:hover) {
  border-color: #b7c6cf;
  background: #ffffff;
}

:deep(.ant-input-affix-wrapper-focused),
:deep(.ant-input:focus),
:deep(.ant-input-affix-wrapper:focus) {
  border-color: var(--primary-600);
  background: #ffffff;
  box-shadow: 0 0 0 3px var(--primary-050);
}

:deep(.ant-input) {
  font-size: 15px;
  background: transparent;
}

:deep(.ant-input::placeholder) {
  color: #a2b0bb;
}

.submit-item {
  margin-top: 24px;
  margin-bottom: 0;
}

:deep(.ant-btn-primary) {
  height: 48px;
  border-radius: 10px;
  background: linear-gradient(135deg, var(--primary-600) 0%, var(--primary-700) 100%);
  border: none;
  font-size: 15px;
  font-weight: 600;
  box-shadow: 0 8px 18px rgba(15, 118, 110, 0.24);
  transition: all 0.2s ease;
}

:deep(.ant-btn-primary:hover:not(:disabled)) {
  filter: brightness(1.06);
  box-shadow: 0 10px 20px rgba(15, 118, 110, 0.3);
}

:deep(.ant-btn-primary:active:not(:disabled)) {
  filter: brightness(0.94);
}

.login-footer {
  margin-top: 28px;
  text-align: center;
  font-size: 12px;
  color: #9caab5;
}

@media (max-width: 1080px) {
  .login-shell {
    height: auto;
    min-height: 100vh;
    grid-template-columns: 1fr;
  }

  .visual-panel {
    min-height: 300px;
    padding: 36px 30px;
    align-items: flex-end;
  }

  .visual-content h2 {
    font-size: 28px;
  }

  .form-panel {
    padding: 24px;
  }
}

@media (max-width: 640px) {
  .visual-panel {
    min-height: 250px;
    padding: 28px 22px;
  }

  .visual-content h2 {
    font-size: 23px;
  }

  .visual-desc {
    margin-bottom: 16px;
  }

  .visual-points {
    gap: 8px;
  }

  .login-card {
    padding: 6px 2px;
  }

  .login-header h1 {
    font-size: 22px;
  }
}
</style>
