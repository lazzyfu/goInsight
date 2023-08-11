
<template>
  <el-dialog
    title="使用提示"
    :visible.sync="visible"
    :before-close="close"
    width="50%"
    center
  >
    <ul>
      <li>
        使用限制：
        <ul>
          <li>DMS已开启审计功能</li>
          <li>单条SQL查询超时为600秒</li>
          <li>默认返回100条记录，LIMIT N最大返回100条数据</li>
          <li>当前平台仅支持查询，数据变更请使用数据库工单系统</li>
        </ul>
      </li>
      <li>
        快捷键：
        <ul>
          <li>command+L组合键可以快速选中当前光标所在的行</li>
          <li>command+enter组合键可以快速执行选中的SQL</li>
          <li>esc键可对结果表格快速切入切出全屏模式</li>
          <li>ctrl+k组合键可快速折叠左侧库表区域</li>
        </ul>
      </li>
      <li>
        使用技巧：
        <ul>
          <li>左侧表显示区域支持左右拉伸</li>
          <li>支持SQL输入框上下拉伸（鼠标置于输入框的右下角进行拖动）</li>
        </ul>
      </li>
    </ul>
  </el-dialog>
</template>

<script>
export default {
  props: {
    visible: Boolean,
  },
  data() {
    return {};
  },
  methods: {
    close() {
      this.$emit('close');
    },
  },
};
</script>


<style scoped>
li {
  list-style-type: circle;
}
li li {
  list-style-type: square;
}
</style>