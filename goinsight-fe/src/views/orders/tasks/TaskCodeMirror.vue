<template>
  <codemirror
    ref="myCm"
    v-model="code"
    :options="cmOptions"
    @ready="onCmReady"
  ></codemirror>
</template>

<script>
import 'codemirror/addon/comment/comment'
import 'codemirror/addon/display/autorefresh'
import 'codemirror/addon/edit/closebrackets'
import 'codemirror/addon/edit/matchbrackets'
import 'codemirror/addon/mode/overlay'
import 'codemirror/addon/selection/active-line'
import 'codemirror/mode/sql/sql.js'

export default {
  data() {
    return {
      code: '',
      cmOptions: {
        mode: 'text/x-mysql',
        indentUnit: 2,
        tabSize: 2,
        indentWithTabs: true,
        smartIndent: true,
        autoRefresh: true,
        lineNumbers: true,
        lineWrapping: true,
        readOnly: true,
        focuse: false,
      },
    };
  },
  methods: {
    onCmReady(cm) {
      cm.setSize('height', `260px`);
    },
    setValue(value) {
      this.codemirror.setValue(value);
    },
    setSize(width, height) {
      this.codemirror.setSize(width, height);
    }
  },
  computed: {
    codemirror() {
      return this.$refs.myCm.codemirror;
    },
  },
};
</script>