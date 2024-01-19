<template>
  <a-modal v-model="visible" title="编辑角色" width="55%" on-ok="onSubmit" @cancel="handleCancel">
    <template slot="footer">
      <a-button key="back" @click="handleCancel">取消</a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="onSubmit">确定</a-button>
    </template>
    <!-- 表单 -->
    <a-form :form="form" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
      <a-form-item v-show="false" label="ID">
        <a-input v-decorator="['id']"> </a-input>
      </a-form-item>
      <a-form-item label="角色名" has-feedback>
        <a-input
          v-decorator="[
            'name',
            {
              rules: [{ required: true, min: 2, max: 32 }],
              validateTrigger: 'blur',
            },
          ]"
        >
          <a-icon slot="prefix" type="heart" style="color: rgba(0, 0, 0, 0.25)" />
        </a-input>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script>
import { updateRolesApi } from '@/api/users'

export default {
  data() {
    return {
      visible: false,
      loading: false,
      form: this.$form.createForm(this, { name: 'rolesEdit' }),
    }
  },
  methods: {
    showModal(row) {
      this.form.resetFields()
      this.$nextTick(() => {
        setTimeout(() => {
          const fieldValues = {
            id: row.id,
            name: row.name,
          }
          this.form.setFieldsValue(fieldValues)
        })
      })
      this.visible = true
    },
    handleCancel(e) {
      this.visible = false
    },
    updateRoles(data) {
      updateRolesApi(data)
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
          this.updateRoles(values)
        }
      })
      this.loading = false
    },
  },
}
</script>