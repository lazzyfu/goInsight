<template>
  <a-drawer
    width="65%"
    placement="right"
    :closable="false"
    @close="onClose"
    :visible="visible"
  >
    <codemirror
      ref="myCm3"
      v-model="code3"
      :options="cmOptions3"
      @ready="onCmReady3"
    ></codemirror>
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
      code3: "",
      cmOptions3: {
        mode: "text/x-mysql",
        autoRefresh: true,
        lineWrapping: true, // 自动换行
        readOnly: true,
      },
    };
  },
  methods: {
    // 子组件不能修改父组件传递的visible，因此需要通过this.$emit发射给父组件
    // close自定义即可，父组件需要对应即可
    onClose() {
      this.$emit("close");
    },
    onCmReady3(cm) {
      cm.setSize("height", "auto");
    },
    fetchData(selectedKeys) {
      const data = {
        key: selectedKeys,
        type: "table_structure",
      };
      getTableInfo(data)
        .then((res) => {
          if (res.code === "0001") {
            this.$notify.error({
              title: "加载表结构失败",
              message: res.message,
            });
            return false;
          }

          this.codemirror3.setValue(res.data);
        })
        .catch((err) => {
          this.$message.error(err.response.status, 3);
        });
    },
  },
  computed: {
    codemirror3() {
      return this.$refs.myCm3.codemirror;
    },
  },
};
</script>

<style>
</style>