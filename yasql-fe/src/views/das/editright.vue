<template>
  <div>
    <div class="code-header">
      <el-input
        size="mini"
        clearable
        type="text"
        maxlength="128"
        show-word-limit
        prefix-icon="el-icon-attract"
        placeholder="group_concat_max_len=4194304;sql_mode=''"
        v-model="sessionVarsString"
      >
        <template slot="append">
          <el-tooltip
            class="item"
            effect="light"
            content="仅支持MySQL/TiDB/ClickHouse设置session级别的变量，多个变量请使用分号隔开，建议该参数仅在需要的时候指定"
            placement="top"
          >
            <el-button>会话参数</el-button>
          </el-tooltip>
        </template>
      </el-input>
      <el-divider direction="vertical"></el-divider>
      <span>
        字符集
        <el-select
          placeholder="选择字符集,默认utf8"
          size="mini"
          style="width: 15%"
          v-model="selectedCharacterSet"
          @change="saveCharacterToCache"
        >
          <el-option
            v-for="item in characterSets"
            :key="item.key"
            :label="item.key"
            :value="item.value"
          ></el-option>
        </el-select>
      </span>
      <el-divider direction="vertical"></el-divider>
      <span
        >主题
        <el-select
          placeholder="选择主题"
          style="width: 15%"
          size="mini"
          v-model="selectedCodeTheme"
          @change="setCodeTheme"
        >
          <el-option
            v-for="item in CodeThemes"
            :key="item.key"
            :label="item.key"
            :value="item.value"
          ></el-option>
        </el-select>
      </span>
    </div>
    <div id="cmHeight" v-loading="bTableLoading">
      <codeMirrorComponent
        ref="codeMirrorComponentChild"
        @execEvent="ExecEvent"
      ></codeMirrorComponent>
      <div class="el-btn">
        <el-button-group>
          <el-tooltip
            class="item"
            content="每次仅允许执行一条SQL，按command+enter组合键可以快速执行选中的SQL"
            placement="top"
          >
            <el-button
              size="small"
              plain
              @click="executeSQL()"
              icon="el-icon-video-play"
              >执行SQL(command+enter)</el-button
            >
          </el-tooltip>
        </el-button-group>
        <el-button-group class="btn-right-group">
          <el-tooltip
            class="item"
            content="格式化SQL语句（建议每条语句后跟上分号，避免格式化后错乱）"
            placement="top"
          >
            <el-button
              size="small"
              plain
              @click="formatSQL()"
              icon="el-icon-edit"
              >格式化</el-button
            >
          </el-tooltip>
          <el-tooltip
            class="item"
            content="生成当前选中数据库的数据字典（新窗口打开）"
            placement="top"
          >
            <el-button
              size="small"
              plain
              @click="loadDBDict()"
              icon="el-icon-document-copy"
              >数据字典</el-button
            >
          </el-tooltip>
          <el-button
            size="small"
            plain
            @click="tipsVisible = true"
            icon="el-icon-thumb"
            >使用提示</el-button
          >
        </el-button-group>
      </div>
      <el-card class="box-card" shadow="never">
        <span v-html="responseMsg"></span>
      </el-card>
    </div>
    <div v-show="showbTable">
      <btableComponent
        :tabIndex="tabIndex"
        ref="bTableComponentChild"
      ></btableComponent>
    </div>
    <template>
      <userTipsComponent
        :visible="tipsVisible"
        @close="tipsVisible = false"
        ref="userTipsComponentChild"
      ></userTipsComponent>
    </template>
  </div>
</template>

<script>
import sqlFormat from 'sql-formatter';

import { executeMySQLQueryApi, executeClickHouseQueryApi } from '@/api/das';

// 导入子组件
import codeMirrorComponent from './codemirror.vue';
import btableComponent from './btable.vue';
import userTipsComponent from './usetips.vue';

export default {
  props: {
    selectedSchema: Object,
    tabIndex: Number,
  },
  components: {
    codeMirrorComponent,
    btableComponent,
    userTipsComponent,
  },
  computed: {
    characterCacheName: function () {
      return 'userCharacterSet-' + this.tabIndex;
    },
    sessionVarsName: function () {
      return 'userSessionVars-' + this.tabIndex;
    },
    characterSet: function () {
      return {
        character_set_client: this.selectedCharacterSet,
        character_set_connection: this.selectedCharacterSet,
        character_set_results: this.selectedCharacterSet,
      };
    },
    codeThemeName: function () {
      return 'codeTheme-' + this.tabIndex;
    },
  },
  data() {
    return {
      responseMsg: '',
      bTableLoading: false,
      showbTable: false,
      selectedCharacterSet: 'utf8',
      characterSets: [
        { key: 'utf8', value: 'utf8' },
        { key: 'utf8mb4', value: 'utf8mb4' },
        { key: 'latin1', value: 'latin1' },
      ],
      selectedCodeTheme: 'default',
      CodeThemes: [
        { key: 'default', value: 'default' },
        { key: 'eclipse', value: 'eclipse' },
        { key: 'solarized dark', value: 'solarized dark' },
        { key: 'dracula', value: 'dracula' },
        { key: 'monokai', value: 'monokai' },
        { key: 'oceanic-next', value: 'oceanic-next' },
      ],
      sessionVarsString: '',
      tipsVisible: false,
    };
  },
  methods: {
    // 设置tab补全
    setTabCompletion(data) {
      this.$refs.codeMirrorComponentChild.setOption('hintOptions', data);
    },
    // 解析用户输入的session vars
    parserSessionVars() {
      this.savesessionVarsToCache();
      var sessionVars = {};
      if (this.sessionVarsString.length > 0) {
        for (const [i, v] of this.sessionVarsString.split(';').entries()) {
          let sessionVar = v.split('=');
          if (sessionVar.length == 2) {
            sessionVars[sessionVar[0]] = sessionVar[1];
          }
        }
      }
      return sessionVars;
    },
    // 执行mysql/tidb/aurora查询
    executeMySQLQuery(data) {
      // params append字符集 && session variables
      data['params'] = { ...this.characterSet, ...this.parserSessionVars() };
      // 响应消息
      var resMsgs = [];
      this.bTableLoading = true;
      executeMySQLQueryApi(data)
        .then((response) => {
          if (response.code === '0000') {
            resMsgs.push('结果: 执行成功');
            resMsgs.push(`耗时: ${response.data.duration}`);
            resMsgs.push(`SQL: ${response.data.sqltext}`);
            this.$message.success('执行成功');
            // 渲染数据
            this.showbTable = true;
            this.$nextTick(() => {
              this.$refs.bTableComponentChild.renderbTable(response.data);
            });
          } else {
            resMsgs.push('结果: 执行失败');
            resMsgs.push(`错误: ${response.message}`);
            this.showbTable = false;
            if (response.message.includes('sessionid')) {
              this.$message.error(
                '执行失败，认证过期，3s后系统自动刷新页面，请于刷新后重新执行'
              );
              setTimeout(() => {
                // 刷新当前页面
                this.$router.go(0);
              }, 3000);
            } else {
              this.$message.error('执行失败');
            }
          }
          resMsgs.push(`请求ID: ${response.request_id}`);
          this.responseMsg = resMsgs.join('<br>');
        })
        .finally(() => {
          this.bTableLoading = false;
        });
    },
    // 执行clickhouse查询
    executeClickHouseQuery(data) {
      // ClickHouse没有编码的概念
      // params append session variables
      data['params'] = { ...this.parserSessionVars() };
      // 响应消息
      var resMsgs = [];
      this.bTableLoading = true;
      executeClickHouseQueryApi(data)
        .then((response) => {
          if (response.code === '0000') {
            resMsgs.push('结果: 执行成功');
            resMsgs.push(`耗时: ${response.data.duration}`);
            resMsgs.push(`SQL: ${response.data.sqltext}`);
            this.$message.success('执行成功');
            // 渲染数据
            this.showbTable = true;
            this.$nextTick(() => {
              this.$refs.bTableComponentChild.renderbTable(response.data);
            });
          } else {
            resMsgs.push('结果: 执行失败');
            resMsgs.push(`错误: ${response.message}`);
            this.showbTable = false;
            if (response.message.includes('sessionid')) {
              this.$message.error(
                '执行失败，认证过期，3s后系统自动刷新页面，请于刷新后重新执行'
              );
              setTimeout(() => {
                // 刷新当前页面
                this.$router.go(0);
              }, 3000);
            } else {
              this.$message.error('执行失败');
            }
          }
          resMsgs.push(`请求ID: ${response.request_id}`);
          this.responseMsg = resMsgs.join('<br>');
        })
        .finally(() => {
          this.bTableLoading = false;
        });
    },
    // 执行SQL
    executeSQL() {
      /*
      this.selectedSchema = {
         "instance_id": "13e7305e-d8fe-11ed-8e15-126be0261c1a",
          "schema": "bizlic",
          "db_type": "TiDB"
      }
      */
      this.saveCodeToCache();
      // 判断是否选择了DB
      if (Object.keys(this.selectedSchema).length === 0) {
        this.$message.warning('请先选择左侧的库');
        return;
      }
      // 获取选择的SQL
      var sqltext = this.$refs.codeMirrorComponentChild.getSelection();
      // 判断是否输入
      if (sqltext.length == 0) {
        this.$message.warning('请鼠标选中要执行的SQL');
        return;
      }

      /*
      post data:
        {
          "instance_id": "291a892e-d8fb-11ed-8e15-126be0261c1a",
          "schema": "das",
          "db_type": "MySQL",
          "params": {
              "character_set_client": "utf8",
              "character_set_connection": "utf8",
              "character_set_results": "utf8"
          },
          "sqltext": "select * from das_records"
        }
       */

      // post data
      const data = {
        ...this.selectedSchema,
        sqltext: sqltext,
      };

      // 根据db type调用不同的执行接口
      var dbType = this.selectedSchema['db_type'].toLowerCase();
      if (dbType === 'tidb' || dbType === 'mysql') {
        this.executeMySQLQuery(data);
      }
      if (dbType === 'clickhouse') {
        this.executeClickHouseQuery(data);
      }
    },
    // 格式化SQL
    formatSQL() {
      this.saveCodeToCache();
      var sql = this.$refs.codeMirrorComponentChild.getValue();
      this.$refs.codeMirrorComponentChild.setValue(
        sqlFormat.format(sql, {
          indent: '  ',
          tabWidth: 2,
        })
      );
    },
    // 保存到cache
    saveCodeToCache() {
      // 保存用户的输入到本地浏览器
      this.$refs.codeMirrorComponentChild.saveCache(
        `dms-codemirror-${this.tabIndex}`
      );
    },
    // 加载缓存
    loadCodeFromcCache() {
      this.$refs.codeMirrorComponentChild.loadCache(
        `dms-codemirror-${this.tabIndex}`
      );
    },
    // 设置字符集
    saveCharacterToCache() {
      localStorage.setItem(this.characterCacheName, this.selectedCharacterSet);
    },
    // 加载字符集
    loadCharacterFromCache() {
      var cache = localStorage.getItem(this.characterCacheName);
      if (cache != null) {
        this.selectedCharacterSet = cache;
      }
    },
    // 设置sessionVarsString
    savesessionVarsToCache() {
      localStorage.setItem(this.sessionVarsName, this.sessionVarsString);
    },
    // 加载字符集
    loadsessionVarsFromCache() {
      var cache = localStorage.getItem(this.sessionVarsName);
      if (cache != null) {
        this.sessionVarsString = cache;
      }
    },
    // 设置codemirror主题
    setCodeTheme() {
      localStorage.setItem(this.codeThemeName, this.selectedCodeTheme);
      this.$nextTick(() => {
        this.$refs.codeMirrorComponentChild.setTheme(this.selectedCodeTheme);
      });
    },
    // 加载codemirror主题
    loadCodeThemeFromCache() {
      var cache = localStorage.getItem(this.codeThemeName);
      if (cache != null) {
        this.selectedCodeTheme = cache;
        this.$nextTick(() => {
          this.$refs.codeMirrorComponentChild.setTheme(cache);
        });
      }
    },
    // 加载数据字典
    loadDBDict() {
      // 判断是否选择了DB
      if (Object.keys(this.selectedSchema).length === 0) {
        this.$message.warning('请先选择左侧的库');
        return;
      } else {
        const routeData = this.$router.resolve({
          path: '/sqlquery/das/dbdict',
          query: this.selectedSchema,
        });
        window.open(routeData.href, '_blank');
      }
    },
    ExecEvent() {
      this.executeSQL();
    },
  },
  mounted() {
    this.loadCodeFromcCache();
    this.loadCharacterFromCache();
    this.loadCodeThemeFromCache();
    this.loadsessionVarsFromCache();
  },
};
</script>

<style lang='less' scoped>
.box-card {
  font-size: 12px;
  margin-top: 4px;
  height: 100px;
  width: 100%;
  overflow: auto;
  zoom: 1;
  white-space: normal;
  word-break: break-all;
}
.btn-right-group {
  float: right;
}
.code-header {
  display: block;
  margin-bottom: 4px;
}
.el-input {
  width: 55%;
}
.el-btn {
  margin-top: 4px;
}
</style>