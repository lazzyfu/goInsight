<template>
  <a-list itemLayout="horizontal" :dataSource="data">
    <a-list-item slot="renderItem" slot-scope="item, index" :key="index">
      <a-list-item-meta>
        <a slot="title">{{ item.title }}</a>
        <span slot="description">
          <span class="security-list-description">{{ item.description }}</span>
          <span v-if="item.value"> : </span>
          <span class="security-list-value">{{ item.value }}</span>
        </span>
      </a-list-item-meta>
      <template v-if="item.actions">
        <a slot="actions" @click="item.actions.callback">{{ item.actions.title }}</a>
      </template>
    </a-list-item>

    <a-modal title="修改密码" v-model="visible" width="40%" on-ok="onSubmit" @cancel="handleCancel">
      <template slot="footer">
        <a-button key="submit" type="primary" :loading="loading" @click="onSubmit"> 保存 </a-button>
        <a-button key="back" @click="handleCancel"> 关闭 </a-button>
      </template>
      <a-form :form="form" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
        <a-form-item label="当前密码" has-feedback>
          <a-input v-decorator="['current_password']" autocomplete="new-password" type="password">
            <a-icon slot="prefix" type="lock" style="color: rgba(0, 0, 0, 0.25)" />
          </a-input>
        </a-form-item>
        <a-form-item label="新密码" has-feedback>
          <a-input
            v-decorator="[
              'new_password',
              {
                rules: [{ required: true, min: 7, max: 32, validator: validatorPass }],
                validateTrigger: 'blur',
              },
            ]"
            autocomplete="new-password"
            type="password"
          >
            <a-icon slot="prefix" type="lock" style="color: rgba(0, 0, 0, 0.25)" />
          </a-input>
        </a-form-item>
        <a-form-item label="确认密码" has-feedback>
          <a-input
            v-decorator="[
              'verify_password',
              {
                rules: [{ required: true, min: 7, max: 32, validator: validatorVerifyPass }],
                validateTrigger: 'blur',
              },
            ]"
            autocomplete="new-password"
            type="password"
          >
            <a-icon slot="prefix" type="lock" style="color: rgba(0, 0, 0, 0.25)" />
          </a-input>
        </a-form-item>
      </a-form>
    </a-modal>
  </a-list>
</template>

<script>
import { ChangePassword } from '@/api/profile'

export default {
  data() {
    return {
      form: this.$form.createForm(this),
      data: [
        {
          title: '账户密码',
          description: '当前密码强度',
          value: '强',
          actions: {
            title: '修改',
            callback: () => {
              this.showModal()
            },
          },
        },
      ],
      loading: false,
      visible: false,
      // 自定义验证规则
      validatorPass: (rule, value, callback) => {
        var reg = /^.*(?=.{7,})(?=.*\d)(?=.*[A-Z])(?=.*[a-z])(?=.*[!@#$%^&*? ]).*$/
        if (!value.match(reg))
          return callback('密码强度不满足要求，最少7位，包括至少1个大写字母，1个小写字母，1个数字，1个特殊字符')
        callback()
      },
      validatorVerifyPass: (rule, value, callback) => {
        var reg = /^.*(?=.{7,})(?=.*\d)(?=.*[A-Z])(?=.*[a-z])(?=.*[!@#$%^&*? ]).*$/
        if (!value.match(reg))
          return callback('密码强度不满足要求，最少7位，包括至少1个大写字母，1个小写字母，1个数字，1个特殊字符')
        const form = this.form
        if (value && value !== form.getFieldValue('new_password')) {
          callback('您输入的两个密码不一致')
        }
        callback()
      },
    }
  },
  methods: {
    showModal() {
      this.visible = true
    },
    // close modal
    handleCancel(e) {
      this.visible = false
    },
    onSubmit(e) {
      e.preventDefault()
      this.form.validateFields((err, values) => {
        if (!err) {
          ChangePassword(values).then((res) => {
            if (res.code === '0000') {
              this.$message.info(res.message)
              this.handleCancel()
              location.reload() // 刷新页面，此时token已经过期，需要重新登录
            } else {
              this.$message.error(res.message)
            }
          })
        }
      })
    },
  },
}
</script>

<style scoped></style>