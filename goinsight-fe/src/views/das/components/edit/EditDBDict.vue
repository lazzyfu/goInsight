<template>
  <a-drawer :title="dbtitle" width="75%" :closable="false" :visible="visible" @close="onClose">
    <a-spin tip="Loading..." :spinning="loading">
      <div class="spin-content">
        <div>
          <h3>一、表名索引</h3>
          <div style="margin-top: 10px" v-html="dictIndex"></div>
        </div>
        <div style="margin-top: 20px">
          <h3>二、表结构详情</h3>
          <div v-html="dictData"></div>
        </div>
      </div>
    </a-spin>
  </a-drawer>
</template>

<script>
import 'bootstrap/dist/css/bootstrap.min.css'
import { getDBDictApi } from '@/api/das'

export default {
  data() {
    return {
      loading: false,
      visible: false,
      dbtitle: '',
      dictIndex: '',
      dictData: '',
    }
  },
  methods: {
    show(data) {
      this.fetchData(data)
      this.loading = true
      this.visible = true
    },
    onClose() {
      this.visible = false
      ;(this.dbtitle = ''), (this.dictIndex = ''), (this.dictData = ''), this.$emit('close')
    },

    fetchData(data) {
      getDBDictApi(data).then((response) => {
        if (response.code === '0000') {
          let curDate = new Date()
          let dbtime = curDate.toLocaleDateString().split('/').join('-')
          this.dbtitle = `${data.schema}库数据字典[${dbtime}]`
          this.format(response.data)
          this.loading = false
        } else {
          this.onClose()
          this.$message.error(response.message)
          this.loading = false
        }
      })
    },
    format(data) {
      let num = 0
      let dictData = ''
      let dictIndex = ''
      data.forEach(function (row) {
        num += 1
        let tableName = row.TABLE_NAME
        let tableComment = row.TABLE_COMMENT === '' ? 'None' : row.TABLE_COMMENT
        let createTime = row.CREATE_TIME
        let columns = row.COLUMNS_INFO.split('<a>')
        let indexes = row.INDEXES_INFO.split('<a>')

        dictIndex += `
              <div style="margin-top: 8px;padding-left: 12px;">
                <a href='#${tableName}'>${num}、${tableName} ............ ${tableComment}</a><br>
              </div>
            `

        let dictData_row = `
              <div style="height: auto;overflow: scroll;border:1px solid #ccc;border-radius: 5px;padding:12px;margin:10px">
                <a style='color: black;font-size:14px;font-weight: bold' name='${tableName}'>
                ${num}、表名: ${tableName} 备注: ${tableComment} 创建时间: ${createTime}</a>
              `
        // 表数据
        let table_tr_html = ''
        let table_row_num = 0
        columns.forEach(function (i) {
          table_row_num += 1
          table_tr_html += `<tr><td>${table_row_num}</td>`
          i.split('<b>').forEach(function (j) {
            table_tr_html += `<td>${j}</td>`
          })
          table_tr_html += '</tr>'
        })
        // 表索引
        let index_tr_html = ''
        let index_row_num = 0
        indexes.forEach(function (i) {
          index_row_num += 1
          index_tr_html += `<tr><td>${index_row_num}</td>`
          i.split('<b>').forEach(function (j) {
            index_tr_html += `<td>${j}</td>`
          })
          index_tr_html += '</tr>'
        })

        dictData_row += `
              <table class="table table-sm table-hover table-bordered" style='font-size: 12px; margin-top: 12px;padding-top: 5px;margin-bottom: 5px;'>
                <thead>
                  <tr style='background: #e4dede;'>
                    <th>序列</th>
                    <th>列名</th>
                    <th>数据类型</th>
                    <th>可空</th>
                    <th>默认值</th>
                    <th>字符集</th>
                    <th>校对规则</th>
                    <th>备注</th>
                  </tr>
                </thead>
                <tbody>
                  ${table_tr_html}
                </tbody>
              </table>
              <table class="table table-sm table-hover table-bordered" style='font-size: 12px; margin-top: 20px;padding-top: 5px'>
                <thead>
                  <tr style='background: #e4dede;'>
                    <th>序列</th>
                    <th>索引名</th>
                    <th>唯一</th>
                    <th>基数</th>
                    <th>类型</th>
                    <th>包含字段</th>
                  </tr>
                </thead>
                <tbody>
                  ${index_tr_html}
                </tbody>
              </table>
              </div>
              `
        dictData += dictData_row
      })
      this.dictIndex = dictIndex
      this.dictData = dictData
    },
  },
}
</script>