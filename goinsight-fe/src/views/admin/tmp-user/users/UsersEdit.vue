<template>
  <a-modal v-model="visible" title="编辑用户" width="55%" on-ok="onSubmit" @cancel="handleCancel">
    <template slot="footer">
      <a-button key="back" @click="handleCancel">取消</a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="onSubmit">确定</a-button>
    </template>
    <!-- 表单 -->
    <a-form :form="form" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
      <a-form-item v-show="false" label="用户ID">
        <a-input v-decorator="['uid']"> </a-input>
      </a-form-item>
      <a-form-item label="用户名">
        <a-input
          v-decorator="[
            'username',
            {
              rules: [{ required: true, min: 2, max: 32, validator: validatorUsername }],
              validateTrigger: 'blur',
            },
          ]"
          autocomplete="current-password"
        >
          <a-icon slot="prefix" type="user" style="color: rgba(0, 0, 0, 0.25)" />
        </a-input>
      </a-form-item>
      <a-form-item label="密码">
        <a-input
          v-decorator="[
            'password',
            {
              rules: [{ required: true, min: 7, max: 128 }],
              validateTrigger: 'blur',
            },
          ]"
          disabled
          autocomplete="current-password"
          type="password"
        >
          <a-icon slot="prefix" type="lock" style="color: rgba(0, 0, 0, 0.25)" />
        </a-input>
      </a-form-item>
      <a-form-item label="昵称">
        <a-input
          v-decorator="[
            'nick_name',
            {
              rules: [{ required: true, min: 2, max: 32 }],
              validateTrigger: 'blur',
            },
          ]"
        >
          <a-icon slot="prefix" type="user" style="color: rgba(0, 0, 0, 0.25)" />
        </a-input>
      </a-form-item>
      <a-form-item label="邮箱">
        <a-input
          v-decorator="[
            'email',
            {
              rules: [
                {
                  type: 'email',
                  message: 'The input is not valid E-mail!',
                },
                {
                  required: true,
                  message: 'Please input your E-mail!',
                },
              ],
            },
          ]"
        >
          <a-icon slot="prefix" type="mail" style="color: rgba(0, 0, 0, 0.25)" />
        </a-input>
      </a-form-item>
      <a-form-item label="手机号">
        <a-input
          v-decorator="[
            'mobile',
            {
              rules: [{ min: 3, max: 32 }],
              validateTrigger: 'blur',
            },
          ]"
        >
          <a-icon slot="prefix" type="phone" style="color: rgba(0, 0, 0, 0.25)" />
        </a-input>
      </a-form-item>
      <a-form-item label="角色">
        <a-select v-decorator="['role_id']" placeholder="请选择角色" show-search>
          <a-select-option v-for="(item, index) in roles" :key="index" :label="item.name" :value="item.id">
            {{ item.name }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="状态" help="该用户是否可以登录站点">
        <a-select
          v-decorator="['is_active', { initialValue: 'YES', rules: [{ required: true, message: '请选择是否激活' }] }]"
          placeholder="请选择是否激活"
        >
          <a-select-option value="YES"> 激活 </a-select-option>
          <a-select-option value="NO"> 禁用 </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="2FA" help="是否开启2次登录验证">
        <a-select
          v-decorator="[
            'is_two_fa',
            { initialValue: 'YES', rules: [{ required: true, message: '是否开启2次登录验证' }] },
          ]"
          placeholder="是否开启2次登录验证"
        >
          <a-select-option value="YES"> 是 </a-select-option>
          <a-select-option value="NO"> 否 </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="管理员" help="该用户是否拥有管理员权限">
        <a-select
          v-decorator="[
            'is_superuser',
            { initialValue: 'NO', rules: [{ required: true, message: '请选择是否为管理员' }] },
          ]"
          placeholder="请选择是否为管理员"
        >
          <a-select-option value="YES"> 是 </a-select-option>
          <a-select-option value="NO"> 否 </a-select-option>
        </a-select>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script>
import { updateUsersApi, getRolesApi } from '@/api/users'

export default {
  data() {
    return {
      visible: false,
      loading: false,
      form: this.$form.createForm(this, { name: 'usersEdit' }),
      roles: [],
      // 自定义验证规则
      validatorUsername: (rule, value, callback) => {
        var reg = /^([a-zA-Z]|[_])([\w|.]){2,15}$/
        if (!value.match(reg)) return callback('用户名必须为英文字符或数字或.的组合')
        callback()
      },
    }
  },
  methods: {
    getRoles() {
      this.loading = true
      const params = {
        is_page: false,
      }
      getRolesApi(params)
        .then((res) => {
          this.roles = res.data
        })
        .catch((_error) => {})
        .finally(() => {
          this.loading = false
        })
    },
    showModal(row) {
      this.getRoles()
      this.form.resetFields()
      this.$nextTick(() => {
        setTimeout(() => {
          const fieldValues = {
            uid: row.uid,
            username: row.username,
            password: row.password,
            email: row.email,
            nick_name: row.nick_name,
            mobile: row.mobile,
            role_id: row.role_id === 0 ? null : row.role_id,
            is_two_fa: row.is_two_fa === true ? 'YES' : 'NO',
            is_superuser: row.is_superuser === true ? 'YES' : 'NO',
            is_active: row.is_active === true ? 'YES' : 'NO',
          }
          this.form.setFieldsValue(fieldValues)
        })
      })
      this.visible = true
    },
    handleCancel(e) {
      this.visible = false
    },
    updateUsers(data) {
      updateUsersApi(data)
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
          // 字符串转换为bool类型
          values['is_two_fa'] = values['is_two_fa'] === 'YES'
          values['is_active'] = values['is_active'] === 'YES'
          values['is_superuser'] = values['is_superuser'] === 'YES'
          this.updateUsers(values)
        }
      })
      this.loading = false
    },
  },
}
</script>
