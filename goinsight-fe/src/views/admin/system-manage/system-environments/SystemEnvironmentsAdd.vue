<template>
  <a-modal v-model="visible" title="新增环境" width="40%" on-ok="onSubmit" @cancel="handleCancel">
    <template slot="footer">
      <a-button key="back" @click="handleCancel">取消</a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="onSubmit">确定</a-button>
    </template>
    <!-- 表单 -->
    <a-form :form="form" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
      <a-form-item label="环境" has-feedback>
        <a-input
          v-decorator="[
            'name',
            {
              rules: [{ required: true, min: 3, max: 128, message: '请输入环境名' }],
              validateTrigger: 'blur',
            },
          ]"
        >
        </a-input>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script>
import { adminCreateEnvironmentsApi } from '@/api/common'

export default {
  data() {
    return {
      visible: false,
      loading: false,
      form: this.$form.createForm(this, { name: 'systemEnvironmentAdd' }),
    }
  },
  methods: {
    showModal() {
      this.form.resetFields()
      this.visible = true
    },
    handleCancel(e) {
      this.visible = false
    },
    createEnvironments(data) {
      adminCreateEnvironmentsApi(data)
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
          this.createEnvironments(values)
        }
      })
      this.loading = false
    },
  },
}
</script>