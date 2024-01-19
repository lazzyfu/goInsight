<template>
  <a-modal v-model="visible" title="编辑" width="55%" on-ok="onSubmit" @cancel="handleCancel">
    <template slot="footer">
      <a-button key="back" @click="handleCancel">取消</a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="onSubmit">确定</a-button>
    </template>
    <!-- 表单 -->
    <a-form :form="form" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
      <a-form-item v-show="false" label="ID">
        <a-input v-decorator="['id']"> </a-input>
      </a-form-item>
      <a-form-item label="标题">
        <a-input style="width: 100%" clearable v-decorator="['title']"></a-input>
      </a-form-item>
      <a-form-item label="SQL语句">
        <a-textarea
          style="width: 100%"
          clearable
          :auto-size="{ minRows: 5, maxRows: 25 }"
          v-decorator="['sqltext']"
        ></a-textarea>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script>
import { updateFavoritesApi } from '@/api/das'

export default {
  data() {
    return {
      visible: false,
      loading: false,
      form: this.$form.createForm(this, { name: 'favoriteEdit' }),
    }
  },
  methods: {
    showModal(row) {
      this.form.resetFields()
      this.$nextTick(() => {
        setTimeout(() => {
          const fieldValues = {
            id: row.id,
            title: row.title,
            sqltext: row.sqltext,
          }
          this.form.setFieldsValue(fieldValues)
        })
      })
      this.visible = true
    },
    handleCancel(e) {
      this.visible = false
    },
    updateFavorites(data) {
      updateFavoritesApi(data)
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
          this.updateFavorites(values)
        }
      })
      this.loading = false
    },
  },
}
</script>
