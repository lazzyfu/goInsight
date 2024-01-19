<template>
  <a-modal v-model="visible" title="新增数据库实例" width="55%" on-ok="onSubmit" @cancel="handleCancel">
    <template slot="footer">
      <a-button key="back" @click="handleCancel">取消</a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="onSubmit">确定</a-button>
    </template>
    <!-- 表单 -->
    <a-form :form="form" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
      <a-form-item label="环境" has-feedback>
        <a-select
          v-decorator="['environment', { rules: [{ required: true, message: '请选择环境' }] }]"
          placeholder="请选择环境"
          allowClear
          show-search
        >
          <a-select-option v-for="(item, index) in environments" :key="index" :label="item.name" :value="item.id">
            {{ item.name }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="组织" has-feedback>
        <a-cascader
          v-decorator="['organization_key', { rules: [{ required: true, message: '请选择组织' }] }]"
          :field-names="{ label: 'title', value: 'key', children: 'children' }"
          :options="organizations"
          change-on-select
          expand-trigger="hover"
          placeholder="请选择组织"
        >
        </a-cascader>
      </a-form-item>
      <a-form-item label="类型" has-feedback>
        <a-select
          v-decorator="['db_type', { initialValue: 'MySQL', rules: [{ required: true, message: '请选择数据库类型' }] }]"
          placeholder="请选择数据库类型"
        >
          <a-select-option value="MySQL"> MySQL </a-select-option>
          <a-select-option value="TiDB"> TiDB </a-select-option>
          <a-select-option value="ClickHouse"> ClickHouse </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="用途" has-feedback>
        <a-select
          v-decorator="['use_type', { initialValue: '工单', rules: [{ required: true, message: '请选择数据库用途' }] }]"
          placeholder="请选择数据库用途"
        >
          <a-select-option value="查询"> 查询 </a-select-option>
          <a-select-option value="工单"> 工单 </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="主机名" has-feedback>
        <a-input
          v-decorator="[
            'hostname',
            {
              rules: [{ required: true, min: 3, max: 128, message: '请输入主机名' }],
              validateTrigger: 'blur',
            },
          ]"
        >
        </a-input>
      </a-form-item>
      <a-form-item label="端口" has-feedback>
        <a-input-number
          v-decorator="[
            'port',
            {
              type: 'number',
              rules: [{ required: true, message: '请输入端口' }],
              validateTrigger: 'blur',
            },
          ]"
        >
        </a-input-number>
      </a-form-item>
      <a-form-item label="备注" has-feedback>
        <a-input
          v-decorator="[
            'remark',
            {
              rules: [{ required: true, min: 3, max: 256, message: '请输入备注' }],
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
import { getOrganizationsApi } from '@/api/users'
import { adminCreateDBConfigApi, adminGetEnvironmentsApi } from '@/api/common'

export default {
  data() {
    return {
      visible: false,
      loading: false,
      environments: [],
      organizations: [],
      form: this.$form.createForm(this, { name: 'dbconfigAdd' }),
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
    async getEnvironments() {
      try {
        const res = await adminGetEnvironmentsApi({ is_page: false })
        this.environments = res.data
      } catch (error) {
        this.$message.error('Failed to fetch environments:', error)
      }
    },
    async getOrganizations(){
      try {
        const res = await getOrganizationsApi({ is_page: false })
        this.organizations = res.data
      } catch (error){
        this.$message.error('Failed to fetch organizations:', error)
      }
    },
    createDBConfig(data) {
      adminCreateDBConfigApi(data)
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
          // 将 values['port'] 转换为数字类型
          values['port'] = parseInt(values['port'], 10)
          // 确保转换成功后再进行后续操作
          if (!isNaN(values['port'])) {
            this.createDBConfig(values)
          }
        }
      })
      this.loading = false
    },
  },
  mounted() {
    this.getEnvironments()
    this.getOrganizations()
  },
}
</script>
