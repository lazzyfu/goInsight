<template>
  <a-modal v-model="visible" title="新增数据库访问权限" width="55%" on-ok="onSubmit" @cancel="handleCancel">
    <template slot="footer">
      <a-button key="back" @click="handleCancel">取消</a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="onSubmit">确定</a-button>
    </template>
    <!-- 表单 -->
    <a-form :form="form" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
      <a-form-item label="用户" help="需要开通访问权限的用户" has-feedback>
        <a-select
          v-decorator="['username', { rules: [{ required: true, message: '请选择用户' }] }]"
          placeholder="请选择用户"
          allowClear
          show-search
        >
          <a-select-option v-for="(item, index) in users" :key="index" :label="item.username" :value="item.username">
            {{ item.username }} <strong>{{ item.nick_name }}</strong>
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="DB类型" has-feedback>
        <a-select
          @change="changeDBType"
          v-decorator="['db_type', { initialValue: 'MySQL', rules: [{ required: true, message: '请选择DB类型' }] }]"
          placeholder="请选择DB类型"
          allowClear
          show-search
        >
          <a-select-option v-for="(item, index) in dbTypes" :key="index" :label="item" :value="item">
            {{ item }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="环境" has-feedback>
        <a-select
          @change="changeEnvs"
          v-decorator="['environment', { rules: [{ required: true, message: '请选择工单环境' }] }]"
          placeholder="请选择工单环境"
          allowClear
          show-search
        >
          <a-select-option v-for="(item, index) in environments" :key="index" :label="item.name" :value="item.id">
            {{ item.name }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="实例" has-feedback>
        <a-select
          @change="changeIns"
          v-decorator="['instance_id', { rules: [{ required: true, message: '请选择数据库实例' }] }]"
          placeholder="请选择数据库实例"
          allowClear
          show-search
        >
          <a-select-option
            v-for="(item, index) in instances"
            :key="index"
            :label="item.remark"
            :value="item.instance_id"
          >
            {{ item.remark }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="库名" help="需要开通访问权限的库名" has-feedback>
        <a-select
          v-decorator="['schema', { rules: [{ required: true, message: '请选择数据库' }] }]"
          placeholder="请选择数据库"
          allowClear
          show-search
          @change="getTables"
        >
          <a-select-option v-for="(item, index) in schemas" :key="index" :label="item.schema" :value="item.schema">
            {{ item.schema }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item
        label="表名"
        help="需要开通访问权限的表，可不选。不选，默认开通当前库的访问权限；选择指定的表，仅开通当前库指定表的访问权限"
        has-feedback
      >
        <a-select
          v-decorator="['tables', { rules: [{ required: false, message: '请选择表名' }] }]"
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
    </a-form>
  </a-modal>
</template>

<script>
import { getUsersApi } from '@/api/users'
import { adminGetEnvironmentsApi } from '@/api/common'
import {
  adminGetInstancesListApi,
  adminGetSchemasListApi,
  adminGetTablesListApi,
  adminCreateSchemasGrantApi,
} from '@/api/das'

const dbTypes = ['MySQL', 'TiDB', 'ClickHouse']

export default {
  data() {
    return {
      visible: false,
      loading: false,
      users: [],
      environments: [],
      instances: [],
      schemas: [],
      tables: [],
      dbTypes,
      form: this.$form.createForm(this, { name: 'dasAdd' }),
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
    getUsers() {
      getUsersApi({ is_page: false })
        .then((res) => {
          this.users = res.data
        })
        .catch((_error) => {})
    },
    async getEnvironments() {
      try {
        const res = await adminGetEnvironmentsApi({ is_page: false })
        this.environments = res.data
      } catch (error) {
        this.$message.error('Failed to fetch environments:', error)
      }
    },
    changeDBType() {
      this.form.resetFields(['environment', 'instance_id', 'schema'])
    },
    changeEnvs(value) {
      // Change环境时清空指定的字段
      this.form.resetFields(['instance_id', 'schema'])
      // 获取指定环境的实例
      var params = {
        id: value,
        db_type: this.form.getFieldValue('db_type'),
        is_page: false,
      }
      adminGetInstancesListApi(params)
        .then((res) => {
          this.instances = res.data
        })
        .catch((_error) => {})
    },
    changeIns(value) {
      // Change实例时清空指定的字段
      this.form.resetFields(['schema'])
      // 获取指定实例的Schemas
      var params = {
        instance_id: value,
        is_page: false,
      }
      adminGetSchemasListApi(params)
        .then((res) => {
          this.schemas = res.data
        })
        .catch((_error) => {})
    },
    getTables(value) {
      // 切换库名，将选择的表清空
      this.form.resetFields(['tables'])
      const params = {
        instance_id: this.form.getFieldValue('instance_id'),
        schema: value,
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
    
    createAdminGrants(data) {
      adminCreateSchemasGrantApi(data)
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
          if (typeof values['port'] === 'string') {
            values['port'] = parseInt(values['port'])
          }
          this.createAdminGrants(values)
        }
      })
      this.loading = false
    },
  },
  mounted() {
    this.getUsers()
    this.getEnvironments()
  },
}
</script>
