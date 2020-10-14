<template>
  <a-card title="上线版本号">
    <div class="table-page-search-wrapper">
      <a-row :gutter="[8, 16]">
        <a-col :span="10">
          <a-button icon="file-add" @click="addVersion()">新建</a-button>
        </a-col>
        <a-col :span="14">
          <a-row :gutter="[8, 16]">
            <a-form layout="inline" :form="form" @keyup.enter.native="handleSearch">
              <a-col :span="8">
                <a-form-item>
                  <a-input placeholder="输入要检索的内容" v-decorator="decorator['search']" />
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <span class="table-page-search-submitButtons">
                  <a-button type="primary" @click="handleSearch">查询</a-button>
                  <a-button @click="resetForm" style="margin-left: 8px">重置</a-button>
                </span>
              </a-col>
            </a-form>
          </a-row>
        </a-col>
      </a-row>
    </div>
    <a-table
      :columns="table.columns"
      :rowKey="(record) => record.id"
      :dataSource="table.data"
      :pagination="pagination"
      :loading="loading"
      @change="handleTableChange"
      size="middle"
    >
      <template slot="id" slot-scope="text, record">
        <router-link :to="{ name: 'view.sqlorders.version.view', params: { version: record.version } }"
          >详情</router-link
        >
      </template>
      <template v-for="col in ['version', 'expire_time']" :slot="col" slot-scope="text, record">
        <div :key="col">
          <a-input
            v-if="record.editable"
            style="margin: -5px 0"
            :value="text"
            @change="(e) => handleChange(e.target.value, record.key, col)"
          />
          <template v-else>{{ text }}</template>
        </div>
      </template>
      <template slot="key" slot-scope="text, record">
        <div class="editable-row-operations">
          <span v-if="record.editable">
            <a @click="() => save(record.key)">
              <a-tag>保存</a-tag>
            </a>
            <a-popconfirm title="Sure to cancel?" @confirm="() => cancel(record.key)">
              <a>
                <a-tag>关闭</a-tag>
              </a>
            </a-popconfirm>
          </span>
          <span v-else>
            <a :disabled="editingKey !== ''" @click="() => edit(record.key)">编辑</a>
          </span>
          <a type="dashed" style="margin: 10px" @click="DeleteConfirm(text)">删除</a>
        </div>
      </template>
    </a-table>

    <a-modal v-model="visible" title="新建版本">
      <template slot="footer">
        <a-button key="back" @click="handleCancel">关闭</a-button>
        <a-button key="submit" type="primary" @click="handleVersionCommit">提交</a-button>
      </template>
      <el-form :model="ruleForm" :rules="rules" ref="ruleForm" label-width="80px" size="small" class="demo-ruleForm">
        <el-form-item label="版本名称" prop="version">
          <el-input v-model="ruleForm.version"></el-input>
        </el-form-item>

        <el-form-item label="截止日期" prop="expire_time">
          <el-date-picker
            v-model="ruleForm.expire_time"
            type="date"
            value-format="yyyy-MM-dd"
            placeholder="选择日期"
          ></el-date-picker>
        </el-form-item>
      </el-form>
    </a-modal>
  </a-card>
</template>

<script>
import { listReleaseVersions, createReleaseVersions, updateReleaseVersions, deleteReleaseVersions } from '@/api/sql'

export default {
  data() {
    return {
      visible: false,
      versions: [],
      loading: false,
      table: {
        columns: [],
        data: [],
      },
      editingKey: '', // 控制每次只能编辑一行，其他的编辑按钮自动禁用
      pagination: {
        current: 1,
        pageSize: 10,
        total: 0,
        pageSizeOptions: ['10', '20'],
        showSizeChanger: true,
      },
      decorator: {
        search: ['search', { rules: [{ required: false }] }],
      },
      form: this.$form.createForm(this),
      // modal form
      ruleForm: {
        version: '',
        expire_time: '',
      },
      rules: {
        version: [
          { required: true, message: '请输入版本号', trigger: 'blur' },
          {
            min: 3,
            max: 24,
            message: '长度在 3 到 24 个字符',
            trigger: 'blur',
          },
        ],
        expire_time: [
          {
            required: true,
            message: '请选择截止日期',
            trigger: 'change',
          },
        ],
      },
    }
  },
  methods: {
    handleSearch(e) {
      e.preventDefault()
      this.form.validateFields((error, values) => {
        if (error) {
          return
        }
        this.filters = {
          search: values['search'],
        }
        this.pagination.current = 1
        this.fetchData()
      })
    },
    resetForm() {
      this.form.resetFields()
    },
    handleTableChange(pager) {
      this.pagination.current = pager.current
      this.pagination.pageSize = pager.pageSize
      this.fetchData()
    },
    // 获取上线版本号
    fetchData() {
      const params = {
        page_size: this.pagination.pageSize,
        page: this.pagination.current,
        ...this.filters,
      }
      this.loading = true
      listReleaseVersions(params)
        .then((response) => {
          this.pagination.total = response.count
          this.loading = false
          this.table.columns = response.results.columns
          this.table.data = response.results.data
        })
        .finally(() => {
          this.loading = false
        })
    },
    // 删除确认
    DeleteConfirm(id) {
      const _this = this
      this.$confirm({
        title: '警告',
        content: '你确定删除？',
        okText: 'Yes',
        okType: 'danger',
        cancelText: 'No',
        onOk() {
          deleteReleaseVersions(id)
            .then((response) => {
              _this.$message.info(response.message)
              _this.fetchData()
            })
        },
        onCancel() {},
      })
    },
    // 编辑
    handleChange(value, key, column) {
      const newData = [...this.table.data]
      const target = newData.filter((item) => key === item.key)[0]
      if (target) {
        target[column] = value
        this.table.data = newData
      }
    },
    // 编辑状态中的input值
    edit(key) {
      const newData = [...this.table.data]
      const target = newData.filter((item) => key === item.key)[0]
      this.editingKey = key
      if (target) {
        target.editable = true
        this.table.data = newData
      }
    },
    // 编辑确认
    save(key) {
      const newData = [...this.table.data]
      const target = newData.filter((item) => key === item.key)[0]
      if (target) {
        // 删除editable=true这条
        delete target.editable
        // 更新到表格实现实时渲染
        this.table.data = newData
        // 提交数据到后台
        updateReleaseVersions(target)
          .then((response) => {
            if (response.code === '0000') {
              this.$message.info(response.message)
            } else {
              this.$message.error(response.message)
            }
            this.fetchData()
          })
      }
      this.editingKey = ''
    },
    // 确认取消
    cancel(key) {
      const newData = [...this.table.data]
      const target = newData.filter((item) => key === item.key)[0]
      this.editingKey = ''
      if (target) {
        delete target.editable
        this.table.data = newData
      }
    },
    // 新建版本
    addVersion() {
      this.showModal()
    },
    // 模态框
    showModal() {
      this.visible = true
    },
    handleCancel() {
      this.visible = false
      this.$refs['ruleForm'].resetFields()
    },
    handleVersionCommit(e) {
      this.$refs['ruleForm'].validate((valid) => {
        if (valid) {
          const commitData = {
            username: this.$store.getters.userInfo.username,
            ...this.ruleForm,
          }
          createReleaseVersions(commitData)
            .then((response) => {
              if (response.code === '0000') {
                this.$message.info(response.message)
                this.visible = false
                this.fetchData()
              } else {
                this.$message.error(response.message)
              }
            })
            .finally(() => {
              this.visible = false
            })
        } else {
          return false
        }
      })
    },
  },
  mounted() {
    this.fetchData()
  },
}
</script>

<style>
</style>