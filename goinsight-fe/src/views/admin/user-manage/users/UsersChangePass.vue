<template>
  <a-modal v-model="visible" :title="title" width="55%" on-ok="onSubmit" @cancel="handleCancel">
    <template slot="footer">
      <a-button key="back" @click="handleCancel">取消</a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="onSubmit">确定</a-button>
    </template>
    <!-- 表单 -->
    <a-form :form="form" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
      <a-form-item v-show="false" label="用户ID">
        <a-input v-decorator="['uid']"> </a-input>
      </a-form-item>
      <a-form-item label="新密码" has-feedback>
        <a-input
          v-decorator="[
            'password',
            {
              rules: [{ required: true, min: 2, max: 32, validator: validatorPass }],
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
              rules: [{ required: true, min: 2, max: 32, validator: validatorVerifyPass }],
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
</template>

<script>
import { changeUsersPassApi } from '@/api/users'

export default {
  data() {
    return {
      visible: false,
      loading: false,
      title: '',
      form: this.$form.createForm(this, { name: 'usersChangePassword' }),
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
        if (value && value !== form.getFieldValue('password')) {
          callback('您输入的两个密码不一致')
        }
        callback()
      },
    }
  },
  methods: {
    showModal(row) {
      this.form.resetFields()
      this.title = '修改用户' + row.username + '的密码'
      this.visible = true
      this.$nextTick(() => {
        setTimeout(() => {
          this.form.setFieldsValue({
            uid: row.uid,
          })
        })
      })
    },
    handleCancel(e) {
      this.visible = false
    },
    changeUsersPass(data) {
      changeUsersPassApi(data)
        .then((res) => {
          const messageType = res.code === '0000' ? 'info' : 'error'
          this.$message[messageType](res.message)
        })
        .catch((_error) => {})
        .finally(() => {
          this.visible = false
          this.loading = false
          this.$emit('refreshTable')
        })
    },
    onSubmit(e) {
      this.loading = true
      e.preventDefault()
      this.form.validateFields((err, values) => {
        if (!err) {
          this.changeUsersPass(values)
        }
      })
      this.loading = false
    },
  },
}
</script>