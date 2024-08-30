<template>
  <a-modal v-model="visible" title="编辑节点" width="55%" on-ok="onSubmit" @cancel="handleCancel">
    <template slot="footer">
      <a-button key="back" @click="handleCancel">取消</a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="onSubmit">确定</a-button>
    </template>
    <!-- 表单 -->
    <a-form :form="form" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
      <a-form-item v-show="false" label="key" has-feedback>
        <a-input v-decorator="['key']"> </a-input>
      </a-form-item>
      <a-form-item label="组织名" has-feedback>
        <a-input
          v-decorator="[
            'name',
            {
              rules: [{ required: true, min: 2, max: 32 }],
              validateTrigger: 'blur',
            },
          ]"
        >
          <a-icon slot="prefix" type="apartment" style="color: rgba(0, 0, 0, 0.25)" />
        </a-input>
      </a-form-item>
    </a-form>
  </a-modal>
</template>
  
  <script>
import { updateOrganizationsApi } from '@/api/users'

export default {
  data() {
    return {
      visible: false,
      loading: false,
      form: this.$form.createForm(this, { name: 'nodeEdit' }),
    }
  },
  methods: {
    showModal(row) {
      this.form.resetFields()
      this.$nextTick(() => {
        setTimeout(() => {
          const fieldValues = {
            key: row.key,
            name: row.title,
          }
          this.form.setFieldsValue(fieldValues)
        })
      })
      this.visible = true
    },
    handleCancel(e) {
      this.visible = false
    },
    updateOrganizations(data) {
      updateOrganizationsApi(data)
        .then((res) => {
          const messageType = res.code === '0000' ? 'info' : 'error'
          this.$message[messageType](res.message)
        })
        .catch((_error) => {})
        .finally(() => {
          this.visible = false
          this.loading = false
          this.$emit('refresh')
        })
    },
    onSubmit(e) {
      this.loading = true
      e.preventDefault()
      this.form.validateFields((err, values) => {
        if (!err) {
          this.updateOrganizations(values)
        }
      })
      this.loading = false
    },
  },
}
</script>