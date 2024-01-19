<template>
  <a-modal v-model="visible" title="新增表权限" width="55%" on-ok="onSubmit" @cancel="handleCancel">
    <template slot="footer">
      <a-button key="back" @click="handleCancel">取消</a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="onSubmit">确定</a-button>
    </template>
    <!-- 表单 -->
    <a-form :form="form" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
      <a-form-item label="表名" has-feedback>
        <a-select
          v-decorator="['tables', { rules: [{ required: true, message: '请选择表名' }] }]"
          placeholder="请选择表名"
          allowClear
          mode="multiple"
          show-search
        >
          <a-select-option v-for="(item, index) in tables" :key="index" :label="item" :value="item">
            {{ item }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="规则" has-feedback>
        <a-select
          v-decorator="['rule', { rules: [{ required: true, message: '请选择规则' }] }]"
          placeholder="请选择规则"
          allowClear
        >
          <a-select-option v-for="(item, index) in Rules" :key="index" :label="item" :value="item">
            {{ item }}
          </a-select-option>
        </a-select>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script>
import { adminGetTablesListApi, adminCreateTablesGrantApi } from '@/api/das'

const Rules = ['allow', 'deny']

export default {
  props: {
    inputData: Object,
  },
  data() {
    return {
      visible: false,
      loading: false,
      tables: [],
      Rules,
      form: this.$form.createForm(this, { name: 'dasTablesPermAdd' }),
    }
  },
  methods: {
    showModal(data) {
      this.form.resetFields()
      this.getTables(data)
      this.visible = true
    },
    handleCancel(e) {
      this.visible = false
    },
    getTables(value) {
      this.tables = []
      const params = {
        instance_id: value.instance_id,
        schema: value.schema,
      }
      adminGetTablesListApi(params)
        .then((res) => {
          if (res.code === '0001') {
            this.$notify.error({
              title: '加载表失败，请稍后重试',
              message: res.message,
            })
            return false
          }
          res.data.forEach((item) => {
            this.tables.push(item.table_name)
          })
        })
        .catch((_error) => {})
    },
    createAdminTablesGrants(data) {
      var newData = {
        ...data,
        ...this.inputData,
      }
      adminCreateTablesGrantApi(newData)
        .then((res) => {
          const messageType = res.code === '0000' ? 'info' : 'error'
          this.$message[messageType](res.message)
        })
        .catch((_error) => {})
        .finally(() => {
          this.handleCancel()
          this.$emit('refreshTable')
        })
    },
    onSubmit(e) {
      this.loading = true
      e.preventDefault()
      this.form.validateFields((err, values) => {
        if (!err) {
          this.createAdminTablesGrants(values)
        }
      })
      this.loading = false
    },
  },
}
</script>