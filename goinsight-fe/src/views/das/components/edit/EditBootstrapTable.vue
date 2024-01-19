<template>
  <a-tabs default-active-key="1">
    <a-tab-pane key="1" tab="结果集(按esc键可快速切入切出全屏模式)">
      <table
        data-click-to-select="true"
        data-show-copy-rows="true"
        data-pagination="true"
        data-page-number="1"
        data-page-size="10"
        data-side-pagination="client"
        data-page-list="[10, 20, 50, 100]"
        data-resizable="true"
        id="btTable"
      ></table>
    </a-tab-pane>
  </a-tabs>
</template>

<script>
import $ from 'jquery'

import 'bootstrap/dist/css/bootstrap.min.css'
// bootstrap-table
import 'bootstrap-table/dist/bootstrap-table.min.css'
import 'bootstrap/dist/js/bootstrap.min.js'
// 必须导入icons，否则toolbar icon不显示
import 'bootstrap-icons/font/bootstrap-icons.css'
import 'bootstrap-table/dist/bootstrap-table.min.js'
import 'bootstrap-table/dist/locale/bootstrap-table-zh-CN.min.js'
import 'bootstrap-table/dist/extensions/auto-refresh/bootstrap-table-auto-refresh.min.js'
import 'bootstrap-table/dist/extensions/copy-rows/bootstrap-table-copy-rows.min.js'
// resize columns width
import 'jquery-resizable-columns/dist/jquery.resizableColumns.css'
import 'jquery-resizable-columns/dist/jquery.resizableColumns.min.js'
import 'bootstrap-table/dist/extensions/resizable/bootstrap-table-resizable.min.js'
// jump to
import 'bootstrap-table/dist/extensions/page-jump-to/bootstrap-table-page-jump-to.min.css'
import 'bootstrap-table/dist/extensions/page-jump-to/bootstrap-table-page-jump-to.min.js'
// sticky-header
import 'bootstrap-table/dist/extensions/sticky-header/bootstrap-table-sticky-header.min.css'
import 'bootstrap-table/dist/extensions/sticky-header/bootstrap-table-sticky-header.min.js'

export default {
  data() {
    return {}
  },
  methods: {
    // 渲染列
    renderColumns(columns) {
      var cols = [
        {
          field: 'state',
          checkbox: true,
        },
      ]
      for (const [i, v] of columns.entries()) {
        cols.push({
          field: v,
          title: v,
          escape: true,
        })
      }
      return cols
    },
    // 渲染表格
    renderbTable(data) {
      $('#btTable')
        .bootstrapTable('destroy')
        .bootstrapTable({
          columns: this.renderColumns(data.columns),
          data: data.data,
          locale: 'zh-CN',
          search: true,
          showRefresh: true,
          showToggle: true,
          showColumns: true,
          lineWrapping: true,
          matchBrackets: true,
          showCopyRows: true,
          showJumpTo: true,
          stickyHeader: true,
          classes: 'table table-hover table-bordered table-striped',
          iconSize: 'sm',
          cache: false,
          showFullscreen: true,
          rowStyle: function rowStyle(row, index) {
            return {
              css: { 'font-size': '12px' },
            }
          },
        })

      // 更改icon的大小
      $('#btTable').bootstrapTable('refreshOptions', {
        iconSize: 'sm',
      })
    },
    escEvent(e) {
      // 按esc键可以在全屏和非全屏模式下切换
      if (e.key === 'Escape') {
        $('#btTable').bootstrapTable('toggleFullscreen')
      }
    },
  },
  mounted() {
    window.addEventListener('keyup', this.escEvent, { passive: true })
  },
  destroyed() {
    window.removeEventListener('keyup', this.escEvent)
  },
}
</script>


