<template>
  <a-modal v-model="visible" title="请使用google或第三方身份验证器扫码绑定" @cancel="handleCancel">
    <template slot="footer">
      <a-button key="back" @click="handleCancel"> 关闭 </a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="handleOk"> 我已扫码绑定 </a-button>
    </template>

    <div class="qr-code-container" flex="dir:center cross:center">
      <img
        v-if="otpAuthUrl"
        :src="generateQRCode(otpAuthUrl)"
        alt="QR Code"
        width="100%"
        height="100%"
        margin="0 auto"
      />
    </div>
  </a-modal>
</template>

<script>
import { GetOTPAuthURLApi, GetOTPCallbackApi } from '@/api/profile'
import qrcode from 'qrcode'

export default {
  data() {
    return {
      visible: false,
      loading: false,
      parmas: {},
      callback: '',
      otpAuthUrl: '', // URL for Google Authenticator setup
    }
  },
  methods: {
    showModal(parmas) {
      this.parmas = parmas
      this.fetchOtpAuthUrl(parmas)
    },
    handleCancel(e) {
      this.visible = false
    },
    handleOk(e) {
      this.loading = true
      this.parmas['callback'] = this.callback
      GetOTPCallbackApi(this.parmas)
        .then((res) => {
          if (res.code === '0000') {
            this.$notification.success({
              message: '成功',
              description: res.message,
            })
          } else {
            this.$notification.error({
              message: '错误',
              description: res.message,
            })
          }
        })
        .catch((_error) => {})
        .finally(() => {
          this.visible = false
          this.loading = false
        })
    },
    generateQRCode(url) {
      // 使用 qrcode.js 生成二维码的 base64 数据
      let qrDataUrl = ''
      qrcode.toDataURL(url, (err, dataUrl) => {
        if (err) {
          console.error('Error generating QR code:', err)
        } else {
          qrDataUrl = dataUrl
        }
      })
      return qrDataUrl
    },
    fetchOtpAuthUrl(parmas) {
      GetOTPAuthURLApi(parmas)
        .then((res) => {
          this.visible = true
          this.otpAuthUrl = res.otpAuthUrl
          this.callback = res.callback
        })
        .catch((error) => {
          this.$message.error(error.response.data.message)
        })
    },
  },
}
</script>
