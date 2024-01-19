<template>
  <a-modal title="HOOK工单" v-model="visible" width="40%" on-ok="onSubmit" @cancel="handleCancel">
    <template slot="footer">
      <a-button key="back" @click="handleCancel">取消</a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="onSubmit">确定</a-button>
    </template>
    <a-form :form="form" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
      <a-form-item label="工单ID" v-show="false">
        <a-input v-decorator="['order_id', { rules: [{ required: true }] }]" disabled />
      </a-form-item>
      <a-form-item label="当前工单">
        <a-input v-decorator="['title', { rules: [{ required: true }] }]" disabled />
      </a-form-item>
      <a-form-item label="DB类型">
        <a-input v-decorator="['db_type', { rules: [{ required: true }] }]" disabled />
      </a-form-item>
      <a-form-item label="当前库">
        <a-input v-decorator="['schema', { rules: [{ required: true }] }]" disabled />
      </a-form-item>
      <a-form-item label="目标工单环境" has-feedback>
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
      <a-form-item label="库名" has-feedback>
        <a-select
          v-decorator="['schema', { rules: [{ required: true, message: '请选择数据库' }] }]"
          placeholder="请选择数据库"
          allowClear
          show-search
        >
          <a-select-option v-for="(item, index) in schemas" :key="index" :label="item.schema" :value="item.schema">
            {{ item.schema }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="审核状态">
        <a-switch
          checked-children="重置审核状态为：待审批"
          un-checked-children="继承审核状态为：已批准"
          default-checked
          @change="onRestProgress"
        />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script>
import { getEnvironmentsApi, getInstancesApi, getSchemasApi, hookOrdersApi } from '@/api/orders'

export default {
  props: {
    orderDetail: Object,
  },
  data() {
    return {
      visible: false,
      loading: false,
      environments: [],
      instances: [],
      schemas: [],
      progress: '待审核',
      form: this.$form.createForm(this, { name: 'hook' }),
    }
  },
  methods: {
    showModal() {
      this.visible = true
      this.$nextTick(() => {
        setTimeout(() => {
          this.form.setFieldsValue({
            title: this.orderDetail.title,
            schema: this.orderDetail.schema,
            db_type: this.orderDetail.db_type,
            order_id: this.orderDetail.order_id,
          })
        })
      })
    },
    hideModal() {
      this.visible = false
    },
    handleCancel() {
      this.hideModal()
    },
    // 获取环境
    getEnvironments() {
      getEnvironmentsApi({is_page: false}).then((res) => {
        this.environments = res.data
      })
    },
    // Change环境
    changeEnvs(value) {
      // Change环境时清空指定的字段
      this.form.resetFields(['instance_id', 'schema'])
      // 获取指定环境的实例
      var params = {
        id: value,
        db_type: this.form.getFieldValue('db_type'),
        is_page: false,
      }
      getInstancesApi(params)
        .then((res) => {
          this.instances = res.data
        })
        .catch((_error) => {})
    },
    // Change实例
    changeIns(value) {
      // Change实例时清空指定的字段
      this.form.resetFields(['schema'])
      // 获取指定实例的Schemas
      var params = {
        instance_id: value,
        is_page: false,
      }
      getSchemasApi(params)
        .then((res) => {
          this.schemas = res.data
        })
        .catch((_error) => {})
    },
    // 是否重置审核状态
    onRestProgress(checked) {
      this.progress = checked ? '待审核' : '已批准'
    },
    // onSubmit
    onSubmit(e) {
      this.loading = true
      e.preventDefault()
      this.form.validateFields((err, values) => {
        values['progress'] = this.progress
        if (!err) {
          hookOrdersApi(values)
            .then((res) => {
              if (res.code === '0000') {
                this.$router.push('/orders/list')
              } else {
                this.$message.error(res.message)
              }
            })
            .catch((_error) => {})
        }
      })
      this.loading = false
    },
  },
  mounted() {
    this.getEnvironments()
  },
}
</script>
