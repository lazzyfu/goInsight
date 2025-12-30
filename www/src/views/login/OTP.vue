<template>
  <a-modal
    :open="props.open"
    title="扫描二维码进行绑定"
    width="35%"
    centered
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
          请使用身份验证器 App 扫描上方二维码以完成绑定
        </a-typography-text>
      </div>
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

const emit = defineEmits(['update:open'])

const uiState = reactive({
  loading: false,
})

const formState = reactive({
  params: {},
  callback: '',
  otpAuthUrl: '',
  qrDataUrl: '',
})

const handleCancel = () => {
  emit('update:open', false)
}

const onSubmit = async () => {
  uiState.loading = true
  try {
    const res = await GetOTPCallbackApi({
      ...formState.params,
      callback: formState.callback,
    })

    if (res.code === '0000') {
      notification.success({
        message: '成功',
        description: res.message,
      })
    } else {
      notification.error({
        message: '错误',
        description: res.message,
      })
    }
  } catch {
    message.error('请求失败')
  } finally {
    uiState.loading = false
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

  try {
    const res = await GetOTPAuthURLApi(params)
    formState.otpAuthUrl = res.otpAuthUrl
    formState.callback = res.callback

    emit('update:open', true)
    await generateQRCode(res.otpAuthUrl)
  } catch (e) {
    message.error(e?.response?.data?.message || '获取失败')
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
</style>
