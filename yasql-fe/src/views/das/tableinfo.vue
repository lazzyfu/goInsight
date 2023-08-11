<template>
  <a-drawer
    width="65%"
    placement="right"
    :closable="false"
    @close="onClose"
    :visible="visible"
  >
    <codemirror
      ref="myCm"
      v-model="code"
      :options="cmOptions"
      @ready="onCmReady"
    ></codemirror>
  </a-drawer>
</template>

<script>
import { getTableInfo } from '@/api/das';

export default {
  props: {
    visible: Boolean, // 接收父组件传值
  },
  data() {
    return {
      code: '',
      cmOptions: {
        mode: 'text/x-mysql',
        autoRefresh: true,
        lineWrapping: true,
        readOnly: true,
        lineNumbers: true,
      },
    };
  },
  methods: {
    onClose() {
      this.$emit('close');
    },
    onCmReady(cm) {
      cm.setSize('height', 'auto');
    },
    fetchData(params) {
      getTableInfo(params)
        .then((response) => {
          if (response.code === '0001') {
            this.$message.error(response.message);
            this.$emit('close');
            return false;
          } else {
            if (params['type'] === 'structure') {
              response.data.forEach((element) => {
                for (const i in element) {
                  if (
                    i.toLowerCase() === 'create table' ||
                    i.toLowerCase() === 'statement'
                  ) {
                    this.codemirror.setValue(element[i]);
                  }
                }
              });
            }
            if (params['type'] === 'base') {
              var tableBase = [];
              response.data.forEach((element) => {
                for (var key in element) {
                  tableBase.push(`${key}  ${element[key]}`);
                }
              });
              this.codemirror.setValue(tableBase.join('\n'));
            }
          }
        })
        .catch((e) => {});
    },
  },
  computed: {
    codemirror() {
      return this.$refs.myCm.codemirror;
    },
  },
};
</script>

<style>
</style>