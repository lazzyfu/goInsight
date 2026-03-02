<template>
  <a-modal
    :open="props.open"
    title="扫描二维码进行绑定"
    :width="560"
    centered
    :mask-closable="false"
    @cancel="handleCancel"
  >
    <template #footer>
      <a-button @click="handleCancel">关闭</a-button>
      <a-button type="primary" :loading="uiState.loading" @click="onSubmit">
        我已完成扫码绑定
      </a-button>
    </template>

    <div class="qr-code-container">
      <a-spin :spinning="!formState.qrDataUrl" tip="二维码生成中，请稍候...">
        <div class="qr-code-wrapper">
          <img
            v-if="formState.qrDataUrl"
            :src="formState.qrDataUrl"
            alt="二维码"
            class="qr-code-image"
          />
        </div>
      </a-spin>
      <div class="qr-code-tips">
        <a-typography-text type="secondary">
          请使用身份验证器 App 扫描上方二维码，然后输入当前 6 位动态验证码完成绑定。
        </a-typography-text>
      </div>
      <a-form layout="vertical" class="verify-form">
        <a-form-item label="动态验证码" required>
          <a-input
            v-model:value="formState.otpCode"
            placeholder="请输入6位验证码"
            :maxlength="6"
            allow-clear
          />
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { GetOTPAuthURLApi, GetOTPCallbackApi } from '@/api/profile'
import { message, notification } from 'ant-design-vue'
import qrcode from 'qrcode'
import { reactive } from 'vue'

const props = defineProps({
  open: Boolean,
})

const emit = defineEmits(['update:open', 'bound'])

const uiState = reactive({
  loading: false,
})

const formState = reactive({
  params: {},
  callback: '',
  otpAuthUrl: '',
  qrDataUrl: '',
  otpCode: '',
})

const handleCancel = () => {
  formState.otpCode = ''
  emit('update:open', false)
}

const onSubmit = async () => {
  const code = formState.otpCode.trim()
  if (!/^\d{6}$/.test(code)) {
    message.warning('请输入6位数字动态验证码')
    return
  }

  uiState.loading = true
  try {
    const res = await GetOTPCallbackApi({
      ...formState.params,
      callback: formState.callback,
      otp_code: code,
    })

    if (res.code === '0000') {
      notification.success({
        message: '成功',
        description: res.message,
      })
      emit('bound')
    } else {
      notification.error({
        message: '错误',
        description: res.message,
      })
    }
  } catch {
    // error handled by global interceptor
  } finally {
    uiState.loading = false
    formState.otpCode = ''
    emit('update:open', false)
  }
}

const generateQRCode = async (url) => {
  try {
    formState.qrDataUrl = await qrcode.toDataURL(url, {
      width: 200,
      margin: 1,
    })
  } catch (err) {
    console.error('生成二维码失败', err)
    message.error('生成二维码失败')
  }
}

const show = async (params) => {
  formState.params = params || {}
  formState.qrDataUrl = '' // 重置二维码
  formState.otpCode = ''

  try {
    const res = await GetOTPAuthURLApi(params)
    formState.otpAuthUrl = res.data.otpAuthUrl
    formState.callback = res.data.callback

    emit('update:open', true)
    await generateQRCode(formState.otpAuthUrl)
  } catch {
    // error handled by global interceptor
  }
}

defineExpose({
  show,
})
</script>

<style scoped>
.qr-code-container {
  padding: 24px 0;
  text-align: center;
}

.qr-code-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
  margin-bottom: 16px;
}

.qr-code-image {
  width: 200px;
  height: 200px;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  padding: 8px;
  background: #fff;
}

.qr-code-tips {
  margin-top: 8px;
}

.verify-form {
  max-width: 260px;
  margin: 14px auto 0;
  text-align: left;
}
</style>
