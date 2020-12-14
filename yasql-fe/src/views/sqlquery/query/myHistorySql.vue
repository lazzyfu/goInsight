<template>
  <a-drawer
    width="65%"
    placement="right"
    :closable="false"
    @close="onClose"
    :visible="visible"
  >
    <div class="table-page-search-wrapper">
      <a-row :gutter="[8, 16]">
        <a-form laout="inline" :form="form" @keyup.enter.native="handleSearch">
          <a-col :span="6">
            <a-form-item>
              <a-input
                placeholder="输入要查询的库名或表名"
                v-decorator="decorator['search']"
              />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item>
              <span class="table-page-search-submitButtons">
                <a-button type="primary" @click="handleSearch">查询</a-button>
                <a-button @click="resetForm" style="margin-left: 8px"
                  >重置</a-button
                >
              </span>
            </a-form-item>
          </a-col>
        </a-form>
      </a-row>
    </div>

    <a-table
      :columns="table.columns"
      :rowKey="(record, index) => index"
      :dataSource="table.data"
      :pagination="pagination"
      :loading="loading"
      @change="handleTableChange"
      size="small"
    >
      <p slot="expandedRowRender" slot-scope="record" style="margin: 0">
        {{ record.query_sql }}
      </p>
    </a-table>
  </a-drawer>
</template>

<script>
import { getHistorySql } from "@/api/sqlquery";

export default {
  props: {
    visible: Boolean, // 接收父组件传值
  },
  data() {
    return {
      loading: false,
      table: {
        columns: null,
        data: null,
      },
      pagination: {
        current: 1,
        pageSize: 10,
        total: 0,
        pageSizeOptions: ["10", "20"],
        showSizeChanger: true,
      },
      decorator: {
        search: ["search", { rules: [{ required: false }] }],
      },
      form: this.$form.createForm(this),
    };
  },
  methods: {
    // 子组件不能修改父组件传递的visible，因此需要通过this.$emit发射给父组件
    // close自定义即可，父组件需要对应即可
    onClose() {
      this.$emit("close");
    },
    handleTableChange(pager) {
      this.pagination.current = pager.current;
      this.pagination.pageSize = pager.pageSize;
      this.fetchData();
    },
    fetchData() {
      const params = {
        page_size: this.pagination.pageSize,
        page: this.pagination.current,
        ...this.filters,
      };
      this.loading = true;

      getHistorySql(params)
        .then((response) => {
          this.pagination.total = response.count;
          this.loading = false;
          this.table.columns = response.results.columns;
          this.table.data = response.results.data;
        })
        .catch((error) => {
          this.$message.error(error.response.data.detail, 5);
        })
        .finally(() => {
          this.loading = false;
        });
    },
    handleSearch(e) {
      e.preventDefault();
      this.form.validateFields((error, values) => {
        if (error) {
          return;
        }
        this.filters = {
          search: values["search"],
        };
        this.pagination.current = 1;
        this.fetchData();
      });
    },
    resetForm() {
      this.form.resetFields();
    },
  },
};
</script>

<style>
</style>