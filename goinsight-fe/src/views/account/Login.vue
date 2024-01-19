<template>
  <div class="main">
    <a-form id="formLogin" class="user-layout-login" ref="formLogin" :form="form" @submit="handleSubmit">
      <a-alert v-if="isLoginError" type="error" showIcon style="margin-bottom: 24px" message="" />
      <a-form-item v-show="!showOTP">
        <a-input
          size="large"
          type="text"
          placeholder="用户名"
          v-decorator="[
            'username',
            {
              rules: [{ required: true, message: '请输入帐户名或邮箱地址' }, { validator: handleUsernameOrEmail }],
              validateTrigger: 'change',
            },
          ]"
          autocomplete="new-password"
        >
          <a-icon slot="prefix" type="user" :style="{ color: 'rgba(0,0,0,.25)' }" />
        </a-input>
      </a-form-item>
      <a-form-item v-show="!showOTP">
        <a-input-password
          size="large"
          placeholder="密码"
          autocomplete="new-password"
          v-decorator="['password', { rules: [{ required: true, message: '请输入密码' }], validateTrigger: 'blur' }]"
        >
          <a-icon slot="prefix" type="lock" :style="{ color: 'rgba(0,0,0,.25)' }" />
        </a-input-password>
      </a-form-item>
      <a-form-item v-if="showOTP">
        <a-input
          size="large"
          v-decorator="['otp_code', { rules: [{ required: true, message: '请输入6位OTP码' }] }]"
          placeholder="请输入6位OTP码"
        >
          <a-icon slot="prefix" type="qrcode" :style="{ color: 'rgba(0,0,0,.25)' }" />
        </a-input>
      </a-form-item>
      <a-form-item style="margin-top: 24px">
        <a-button
          size="large"
          type="primary"
          htmlType="submit"
          class="login-button"
          :loading="state.loginBtn"
          :disabled="state.loginBtn"
        >
          <span v-if="!showOTP">{{ $t('user.login.login') }}</span>
          <span v-else>{{ $t('user.login.verify') }}</span>
        </a-button>
      </a-form-item>
    </a-form>
    <OTPComponent ref="OTPComponent"></OTPComponent>
  </div>
</template>

<script>
import { mapActions } from 'vuex'
import { timeFix } from '@/utils/util'

import OTPComponent from './OTP'

export default {
  components: { OTPComponent },
  data() {
    return {
      loginBtn: false,
      isLoginError: false,
      showOTP: false,
      form: this.$form.createForm(this),
      state: {
        time: 60,
        loginBtn: false,
      },
    }
  },
  methods: {
    ...mapActions(['Login', 'Logout']),
    // handler
    handleUsernameOrEmail(rule, value, callback) {
      const { state } = this
      const regex = /^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+((\.[a-zA-Z0-9_-]{2,3}){1,2})$/
      if (regex.test(value)) {
        state.loginType = 0
      } else {
        state.loginType = 1
      }
      callback()
    },
    handleSubmit(e) {
      e.preventDefault()
      const {
        form: { validateFields },
        state,
        Login,
      } = this

      state.loginBtn = true

      const validateFieldsKey = ['username', 'password', 'otp_code']

      validateFields(
        validateFieldsKey,
        {
          force: true,
        },
        // 等待token設置完成
        async (err, values) => {
          if (!err) {
            const loginParams = {
              ...values,
            }
            delete loginParams.username
            loginParams.username = values.username
            loginParams.password = values.password
            loginParams.otp_code = values.otp_code

            try {
              const res = await Login(loginParams)
              this.loginSuccess(res)
            } catch (err) {
              this.requestFailed(err)
            } finally {
              setTimeout(() => {
                state.loginBtn = false
              }, 600)
            }
          } else {
            setTimeout(() => {
              state.loginBtn = false
            }, 600)
          }
        }
      )
    },
    loginSuccess(res) {
      if (res.code === 200) {
        this.$router.push({
          path: '/',
        })
        // 延迟 1 秒显示欢迎信息
        setTimeout(() => {
          this.$notification.success({
            message: '欢迎',
            description: `${timeFix()}，欢迎回来`,
          })
        }, 1000)
      } else {
        this.$notification['error']({
          message: '错误提示',
          description: res.message,
        })
      }
    },
    requestFailed(err, values) {
      if (err.response.config.url === '/api/v1/user/login') {
        if (err.response.data.status === 'otp_required') {
          this.showOTP = true
        } else if (err.response.data.status === 'otp_rebind') {
          this.$notification['warning']({
            message: '警告',
            description: '您需要绑定OTP验证码，请扫码绑定',
            duration: 4,
          })
          var parmas = {
            username: this.form.getFieldValue('username'),
            password: this.form.getFieldValue('password'),
          }
          this.$refs.OTPComponent.showModal(parmas)
        } else {
          this.$notification['error']({
            message: '错误',
            description: ((err.response || {}).data || {}).message || '请求出现错误，请稍后再试',
            duration: 4,
          })
        }
      }
    },
  },
}
</script>

<style lang="less" scoped>
.user-layout-login {
  label {
    font-size: 14px;
  }

  .getCaptcha {
    display: block;
    width: 100%;
    height: 40px;
  }

  .forge-password {
    font-size: 14px;
  }

  button.login-button {
    padding: 0 15px;
    font-size: 16px;
    height: 40px;
    width: 100%;
  }

  .user-login-other {
    text-align: left;
    margin-top: 24px;
    line-height: 22px;

    .item-icon {
      font-size: 24px;
      color: rgba(0, 0, 0, 0.2);
      margin-left: 16px;
      vertical-align: middle;
      cursor: pointer;
      transition: color 0.3s;

      &:hover {
        color: #1890ff;
      }
    }

    .register {
      float: right;
    }
  }
}
</style>
