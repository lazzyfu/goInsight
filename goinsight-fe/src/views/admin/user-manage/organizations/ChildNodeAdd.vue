<template>
  <a-modal v-model="visible" title="新增子节点" width="55%" on-ok="onSubmit" @cancel="handleCancel">
    <template slot="footer">
      <a-button key="back" @click="handleCancel">取消</a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="onSubmit">确定</a-button>
    </template>
    <!-- 表单 -->
    <a-form :form="form" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
      <a-form-item v-show="false" label="父key">
        <a-input v-decorator="['parent_node_key']"> </a-input>
      </a-form-item>
      <a-form-item label="父节点">
        <a-input v-decorator="['parent_node_name']" disabled>
          <a-icon slot="prefix" type="apartment" style="color: rgba(0, 0, 0, 0.25)" />
        </a-input>
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
import { createChildOrganizationsApi } from '@/api/users'

export default {
  data() {
    return {
      visible: false,
      loading: false,
      form: this.$form.createForm(this, { name: 'childNodesAdd' }),
    }
  },
  methods: {
    showModal(row) {
      this.form.resetFields()
      this.$nextTick(() => {
        setTimeout(() => {
          const fieldValues = {
            parent_node_key: row.key,
            parent_node_name: row.title,
          }
          this.form.setFieldsValue(fieldValues)
        })
      })
      this.visible = true
    },
    handleCancel(e) {
      this.visible = false
    },
    createChildOrganizations(data) {
      createChildOrganizationsApi(data)
        .then((res) => {
          const messageType = res.code === '0000' ? 'info' : 'error'
          this.$message[messageType](res.message)
        })
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
          this.createChildOrganizations(values)
        }
      })
      this.loading = false
    },
  },
}
</script>
