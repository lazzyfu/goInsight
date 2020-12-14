<template>
  <a-drawer
    width="75%"
    :title="dbtitle"
    placement="right"
    :closable="false"
    @close="onClose"
    :visible="visible"
  >
    <div
      v-loading="loading"
      element-loading-text="拼命加载中"
      element-loading-spinner="el-icon-loading"
    >
      <div>
        <h5>一、表名索引</h5>
        <div style="margin-top: 10px" v-html="dict_index"></div>
      </div>
      <div style="margin-top: 20px">
        <h5>二、表结构详情</h5>
        <div v-html="dict_data"></div>
      </div>
    </div>
  </a-drawer>
</template>

<script>
import { getDBDict } from "@/api/sqlquery";

export default {
  props: {
    visible: Boolean, // 接收父组件传值
  },
  data() {
    return {
      loading: false,
      dbtitle: "",
      dict_index: "",
      dict_data: "",
    };
  },
  methods: {
    // 子组件不能修改父组件传递的visible，因此需要通过this.$emit发射给父组件
    // close自定义即可，父组件需要对应即可
    onClose() {
      (this.dbtitle = ""),
        (this.dict_index = ""),
        (this.dict_data = ""),
        this.$emit("close");
    },
    fetchData(selectedDB) {
      this.loading = true;
      const data = {
        key: selectedDB,
      };
      getDBDict(data)
        .then((res) => {
          if (res.code === "0001") {
            this.$notify.error({
              title: "加载数据字典失败",
              message: res.message,
            });
            return false;
          }
          let curDate = new Date();
          let dbtime = curDate.toLocaleDateString().split("/").join("-");
          let db = selectedDB.split("___")[1];
          this.dbtitle = `${db}库数据字典[${dbtime}]`;

          let num = 0;
          let dict_data = "";
          let dict_index = "";
          res.data.forEach(function (row) {
            num += 1;
            let table_name = row[0];
            let table_comment = row[1];
            let create_time = row[2];
            let table_value = row[3].split("<a>");
            let index_value = row[4].split("<a>");

            dict_index += `
              <div style="margin-top: 8px;padding-left: 12px;">
                <a href='#${table_name}'>${num}、${table_name} ............ ${table_comment}</a><br>
              </div>
            `;

            let dict_data_row = `
              <div style="height: auto;overflow: scroll;border:1px solid #ccc;border-radius: 5px;padding:12px;margin:10px">
                <a style='color: black;font-size:14px;font-weight: bold' name='${table_name}'>
                ${num}、表名: ${table_name} 备注: ${table_comment} 创建时间: ${create_time}</a>
              `;
            // 表数据
            let table_tr_html = "";
            let table_row_num = 0;
            table_value.forEach(function (i) {
              table_row_num += 1;
              table_tr_html += `<tr><td>${table_row_num}</td>`;
              i.split("<b>").forEach(function (j) {
                table_tr_html += `<td>${j}</td>`;
              });
              table_tr_html += "</tr>";
            });
            // 表索引
            let index_tr_html = "";
            let index_row_num = 0;
            index_value.forEach(function (i) {
              index_row_num += 1;
              index_tr_html += `<tr><td>${index_row_num}</td>`;
              i.split("<b>").forEach(function (j) {
                index_tr_html += `<td>${j}</td>`;
              });
              index_tr_html += "</tr>";
            });

            dict_data_row += `
              <table class="table table-sm table-hover table-bordered" style='font-size: 8px; margin-top: 12px;padding-top: 5px;margin-bottom: 5px;'>
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
              <table class="table table-sm table-hover table-bordered" style='font-size: 8px; margin-top: 20px;padding-top: 5px'>
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
              `;
            dict_data += dict_data_row;
          });
          this.dict_index = dict_index;
          this.dict_data = dict_data;
        })
        .catch((err) => {
          console.log("err: ", err);
          this.$message.error(err.response.status, 3);
        })
        .finally(() => {
          this.loading = false;
        });
    },
  },
};
</script>

<style>
</style>