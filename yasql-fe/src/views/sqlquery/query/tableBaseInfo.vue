<template>
  <a-drawer
    width="65%"
    placement="right"
    :closable="false"
    @close="onClose"
    :visible="visible"
  >
    <div v-html="tableInfoTable"></div>
  </a-drawer>
</template>

<script>
import { getTableInfo } from "@/api/sqlquery";

export default {
  props: {
    visible: Boolean, // 接收父组件传值
  },
  data() {
    return {
      tableInfoTable: "",
    };
  },
  methods: {
    // 子组件不能修改父组件传递的visible，因此需要通过this.$emit发射给父组件
    // close自定义即可，父组件需要对应即可
    onClose() {
      this.$emit("close");
    },
    fetchData(selectedKeys) {
      const data = {
        key: selectedKeys,
        type: "table_base",
      };
      getTableInfo(data)
        .then((res) => {
          if (res.code === "0001") {
            this.$notify.error({
              title: "加载表信息失败",
              message: res.message,
            });
            return false;
          }
          this.tableInfoTable = `<table class="table table-hover">
        ${res.data}
        </table>`;
        })
        .catch((err) => {
          this.$message.error(err.response.status, 3);
        });
    },
  },
};
</script>

<style>
</style>