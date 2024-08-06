<template>
    <a-modal v-model="visible" title="编辑审核参数" width="55%" on-ok="onSubmit" @cancel="handleCancel">
      <template slot="footer">
        <a-button key="back" @click="handleCancel">取消</a-button>
        <a-button key="submit" type="primary" :loading="loading" @click="onSubmit">确定</a-button>
      </template>
      <!-- 表单 -->
      <a-form :form="form" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
        <a-form-item v-show="false" label="ID">
          <a-input v-decorator="['id']"> </a-input>
        </a-form-item>
        <a-form-item label="审核参数" help="格式要求为JSON类型，禁止修改KEY，可修改Value！！！" has-feedback>
        <a-textarea
          :auto-size="{ minRows: 3, maxRows: 5 }"
          placeholder="请输入自定义审核参数，默认为{}"
          v-decorator="[
            'params',
            {
              initialValue: '{}',
              rules: [
                {
                  required: true,
                  message: '请输入自定义审核参数，默认为{}',
                  validator: validatorInspectParams,
                },
              ],
            },
          ]"
        >
        </a-textarea>
      </a-form-item>
        <a-form-item label="备注" help="非必要不要修改" has-feedback>
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
  import { adminUpdateInspectParamsApi } from '@/api/inspect'
  
  export default {
    data() {
      return {
        visible: false,
        loading: false,
        form: this.$form.createForm(this, { name: 'inspectEdit' }),
        validatorInspectParams: (rule, value, callback) => {
          try {
            JSON.parse(value)
          } catch (error) {
            return callback('请输入正确的JSON格式')
          }
          callback()
        },
      }
    },
    methods: {
      showModal(row) {
        this.form.resetFields()
        this.$nextTick(() => {
          setTimeout(() => {
            const fieldValues = {
              id: row.id,
              params: JSON.stringify(row.params),
              remark: row.remark,
            }
            this.form.setFieldsValue(fieldValues)
          })
        })
        this.visible = true
      },
      handleCancel(e) {
        this.visible = false
      },
      UpdateInspectParams(data) {
        adminUpdateInspectParamsApi(data)
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
            // 将 values['params'] 转换为 JSON 对象
            values['params'] = JSON.parse(values['params'])
            this.UpdateInspectParams(values)
          }
        })
        this.loading = false
      },
    },
  }
  </script>
  