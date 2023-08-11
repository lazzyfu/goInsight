<template>
  <el-card class="box-card">
    <!-- 搜索框 -->
    <a-form layout="inline" :form="form" @keyup.enter.native="handleSearch">
      <a-form-item>
        <a-button type="link" style="width: 100%" @click="showModal">
          新增
        </a-button>
      </a-form-item>
      <a-form-item>
        <a-input
          allowClear
          style="width: 300px"
          placeholder="输入要查询的库名或表名"
          v-decorator="decorator['search']"
        />
      </a-form-item>
      <a-form-item>
        <span class="table-page-search-submitButtons">
          <a-button type="primary" @click="handleSearch">查询</a-button>
        </span>
      </a-form-item>
    </a-form>

    <a-table
      :columns="tableColumns"
      :rowKey="(record, index) => index"
      :dataSource="tableData"
      :pagination="pagination"
      :loading="loading"
      @change="handleTableChange"
      size="middle"
    >
      <span slot="expandedRowRender" slot-scope="record" style="margin: 0">
        <pre class="codeStyle"><code>{{ record.sqltext }}</code></pre>
      </span>
      <span slot="action" slot-scope="text, record">
        <div class="editable-row-operations">
          <a
            type="button"
            v-clipboard:copy="record.sqltext"
            v-clipboard:success="onCopy"
            v-clipboard:error="onError"
          >
            <span style="color: #409eff">拷贝</span></a
          >
          <a-divider type="vertical" />
          <a @click="() => editRow(record)"
            ><span style="color: #409eff">编辑</span></a
          >
          <a-divider type="vertical" />
          <a type="dashed" @click="DeleteConfirm(record.id)"
            ><span style="color: #409eff">删除</span></a
          >
        </div>
      </span>
    </a-table>

    <a-modal v-model="addVisible" title="我的SQL" width="50%" on-ok="handleOk">
      <template slot="footer">
        <a-button
          key="submit"
          type="primary"
          :loading="addLoading"
          @click="handleOk"
        >
          保存
        </a-button>
        <a-button key="back" @click="handleCancel"> 关闭 </a-button>
      </template>
      <!-- 表单 -->
      <el-form
        :model="ruleForm"
        :rules="rules"
        ref="ruleForm"
        size="small"
        label-width="100px"
        class="demo-ruleForm"
      >
        <el-form-item label="标题" prop="title">
          <el-input
            style="width: 90%"
            clearable
            v-model="ruleForm.title"
          ></el-input>
        </el-form-item>
        <el-form-item label="SQL语句" prop="sqltext">
          <el-input
            style="width: 90%"
            :autosize="{ minRows: 15, maxRows: 25 }"
            type="textarea"
            placeholder="请输入SQL内容"
            v-model="ruleForm.sqltext"
          ></el-input>
        </el-form-item>
      </el-form>
    </a-modal>
  </el-card>
</template>

<script>
import {
  getFavorites,
  createFavorites,
  updateFavorites,
  deleteFavorites,
} from '@/api/das';

const ruleForm = {
  id: '',
  title: '',
  sqltext: '',
};

const tableColumns = [
  {
    title: '标题',
    dataIndex: 'title',
    key: 'title',
    scopedSlots: {
      customRender: 'title',
    },
  },
  {
    title: '用户名',
    dataIndex: 'username',
    key: 'username',
    scopedSlots: {
      customRender: 'username',
    },
  },
  {
    title: 'SQL',
    dataIndex: 'sqltext',
    key: 'sqltext',
    scopedSlots: {
      customRender: 'sqltext',
    },
    ellipsis: true,
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    key: 'created_at',
    scopedSlots: {
      customRender: 'created_at',
    },
  },
  {
    title: '更新时间',
    dataIndex: 'updated_at',
    key: 'updated_at',
    scopedSlots: {
      customRender: 'updated_at',
    },
  },
  {
    title: '操作',
    dataIndex: 'action',
    key: 'action',
    scopedSlots: {
      customRender: 'action',
    },
  },
];

export default {
  props: {
    visible: Boolean, // 接收父组件传值
  },
  data() {
    return {
      loading: false,
      addLoading: false,
      addVisible: false,
      tableColumns,
      tableData: [],
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
      ruleForm: Object.assign({}, ruleForm),
      rules: {
        title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
        sqltext: [
          { required: true, message: '请输入SQL内容', trigger: 'blur' },
        ],
      },
      opType: 'create',
    };
  },
  methods: {
    // 加载modal
    showModal() {
      this.addVisible = true;
      (this.ruleForm = Object.assign({}, ruleForm)), (this.opType = 'create');
    },
    // 提交
    handleOk(e) {
      this.addLoading = true;
      this.$refs['ruleForm'].validate((valid) => {
        if (valid) {
          if (this.opType === 'create') {
            createFavorites(this.ruleForm)
              .then((res) => {
                if (res.code === '0000') {
                  this.$message.info(res.message);
                } else {
                  this.$message.error(res.message);
                }
              })
              .finally(() => {
                this.addVisible = false;
                this.addLoading = false;
                this.fetchData();
              });
          }
          if (this.opType === 'update') {
            updateFavorites(this.ruleForm)
              .then((res) => {
                if (res.code === '0000') {
                  this.$message.info(res.message);
                } else {
                  this.$message.error(res.message);
                }
              })
              .finally(() => {
                this.addVisible = false;
                this.addLoading = false;
                this.fetchData();
              });
          }
        }
      });
      this.addLoading = false;
    },
    // 编辑行
    editRow(row) {
      this.showModal();
      this.opType = 'update';
      this.ruleForm.id = row.id;
      this.ruleForm.title = row.title;
      this.ruleForm.sqltext = row.sqltext;
    },
    // 删除确认
    DeleteConfirm(id) {
      const _this = this;
      this.$confirm({
        title: '警告',
        content: '你确定删除？',
        okText: 'Yes',
        okType: 'danger',
        cancelText: 'No',
        onOk() {
          deleteFavorites(id)
            .then((response) => {
              if (response.code === '0001') {
                _this.$message.warning(response.message);
              } else {
                _this.$message.info(response.message);
              }
            })
            .finally(() => {
              _this.fetchData();
            });
        },
        onCancel() {},
      });
    },
    // 关闭
    handleCancel(e) {
      this.addVisible = false;
    },
    // 子组件不能修改父组件传递的visible，因此需要通过this.$emit发射给父组件
    // close自定义即可，父组件需要对应即可
    onClose() {
      this.$emit('close');
    },
    handleTableChange(pager) {
      this.pagination.current = pager.current;
      this.pagination.pageSize = pager.pageSize;
      this.fetchData();
    },
    // 加载收藏的SQL
    fetchData() {
      this.loading = true;
      const params = {
        page_size: this.pagination.pageSize,
        page: this.pagination.current,
        ...this.filters,
      };
      getFavorites(params)
        .then((response) => {
          this.pagination.total = response.total;
          this.tableData = response.data;
        })
        .finally(() => {
          this.loading = false;
        });
    },
    // 搜索
    handleSearch(e) {
      e.preventDefault();
      this.form.validateFields((error, values) => {
        if (error) {
          return;
        }
        this.filters = {
          search: values['search'],
        };
        this.pagination.current = 1;
        this.fetchData();
      });
    },
    // 拷贝
    onCopy: function (e) {
      this.$message.info('拷贝成功');
    },
    onError: function (e) {
      this.$message.error('拷贝失败');
    },
  },
  mounted() {
    this.fetchData();
  },
};
</script>

<style lang='less' scoped>
::v-deep .ant-table {
  font-size: 12px;
}
::v-deep .ant-pagination {
  font-size: 12px;
}
::v-deep .ant-select-sm .ant-select-selection__rendered {
  font-size: 12px;
}
::v-deep .codeStyle {
  tab-size: 4;
  background: #183055;
  color: #e6ecf1;
  padding: 12px 12px 12px 8px;
  direction: ltr;
  text-align: left;
  border: 1px solid #d1d1d1;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  max-width: 100%;
  overflow: auto;
  word-break: normal;
  white-space: pre;
  white-space: pre-wrap;
  word-wrap: break-word;
}

::v-deep .el-card__body {
  padding: 8px;
}
</style>
