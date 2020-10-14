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

    <a-modal title="修改密码" v-model="visible">
      <el-form :model="ruleForm" :layout="formLayout" :rules="rules" ref="ruleForm" size="small">
        <el-form-item label="当前密码" prop="current_password">
          <el-input v-model="ruleForm.current_password" type="password" placeholder="当前密码" />
        </el-form-item>

        <el-form-item label="新密码" prop="new_password">
          <el-input v-model="ruleForm.new_password" type="password" placeholder="新密码" />
        </el-form-item>

        <el-form-item label="确认密码" prop="verify_password">
          <el-input v-model="ruleForm.verify_password" type="password" placeholder="确认密码" />
        </el-form-item>
      </el-form>
      <template slot="footer">
        <a-button key="back" @click="handleCancel">关闭</a-button>
        <a-button key="submit" type="primary" :loading="loading" @click="handleOk">提交</a-button>
      </template>
    </a-modal>
  </a-list>
</template>

<script>
import { ChangePassword } from '@/api/user'

export default {
  data() {
    const handlePasswordLevel = (rule, value, callback) => {
      if (value.trim().length < 7 || value.trim().length > 30) {
        callback(new Error('密码长度在 7 到 30 个字符'))
      }

      let level = 0

      // 判断这个字符串中有没有数字
      if (/[0-9]/.test(value)) {
        level++
      }
      // 判断字符串中有没有字母
      if (/[a-zA-Z]/.test(value)) {
        level++
      }
      // 判断字符串中有没有特殊符号
      if (/[^0-9a-zA-Z_]/.test(value)) {
        level++
      }

      if (level === 3) {
        callback()
      } else {
        callback(new Error('密码强度不够，密码必须包含[数字、字符、特殊字符]'))
      }
    }
    const handlePasswordCheck = (rule, value, callback) => {
      if (!value) {
        callback(new Error('请输入密码'))
      }
      if (value && this.ruleForm.new_password && value.trim() !== this.ruleForm.new_password.trim()) {
        callback(new Error('两次密码不一致'))
      }
      callback()
    }
    return {
      data: [
        {
          title: '账户密码',
          description: '当前密码强度',
          value: '强',
          actions: {
            title: '修改',
            callback: () => {
              this.showModal()
            }
          }
        }
      ],
      loading: false,
      visible: false,
      formLayout: 'vertical',
      ruleForm: {
        current_password: '',
        new_password: '',
        verify_password: ''
      },
      rules: {
        current_password: [{ required: true, trigger: 'blur' }],
        new_password: [{ required: true, validator: handlePasswordLevel, trigger: 'blur' }],
        verify_password: [{ required: true, validator: handlePasswordCheck, trigger: 'blur' }]
      }
    }
  },
  methods: {
    showModal() {
      this.visible = true
    },
    hideModal() {
      this.visible = false
    },
    handleCancel() {
      this.hideModal()
    },
    handleOk() {
      this.showModal()
      this.$refs['ruleForm'].validate(valid => {
        if (valid) {
          const commitData = { ...this.ruleForm }
          ChangePassword(commitData).then(response => {
            console.log('response: ', response)
            if (response.code === '0000') {
              this.$message.info(response.message)
              this.handleCancel()
              location.reload()  // 刷新页面，此时token已经过期，需要重新登录
            } else {
              this.$message.error(response.message)
            }
          })
        } else {
          return false
        }
      })
    }
  }
}
</script>

<style scoped></style>
