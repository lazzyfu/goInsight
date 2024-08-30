<template>
  <a-modal v-model="visible" title="新增用户" width="55%" on-ok="onSubmit" @cancel="handleCancel">
    <template slot="footer">
      <a-button key="back" @click="handleCancel">取消</a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="onSubmit">确定</a-button>
    </template>
    <!-- 表单 -->
    <a-form :form="form" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
      <a-form-item v-show="false" label="节点Key">
        <a-input v-decorator="['key']"> </a-input>
      </a-form-item>
      <a-form-item label="用户">
        <a-select mode="multiple" v-decorator="['users']" placeholder="请选择用户" show-search>
          <a-select-option v-for="(item, index) in users" :key="index" :label="item.username" :value="item.uid">
            {{ item.username }}
          </a-select-option>
        </a-select>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script>
import { bindOrganizationsUsersApi, getUsersApi } from '@/api/users'

export default {
  data() {
    return {
      visible: false,
      loading: false,
      users: [],
      form: this.$form.createForm(this, { name: 'nodeUsersAdd' }),
    }
  },
  methods: {
    showModal(nodeKey) {
      this.form.resetFields()
      this.$nextTick(() => {
        setTimeout(() => {
          const fieldValues = {
            key: nodeKey,
          }
          this.form.setFieldsValue(fieldValues)
        })
      })
      this.getUsers()
      this.visible = true
    },
    handleCancel(e) {
      this.visible = false
    },
    getUsers() {
      this.loading = true
      const params = {
        is_page: false,
      }
      getUsersApi(params)
        .then((res) => {
          this.users = res.data
        })
        .catch((_error) => {})
        .finally(() => {
          this.loading = false
        })
    },
    bindOrganizationsUsers(data) {
      bindOrganizationsUsersApi(data)
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
          this.bindOrganizationsUsers(values)
        }
      })
      this.loading = false
    },
  },
}
</script>
