<template>
  <div>
    <div class="code-header">
      <a-input
        style="width: 50%"
        allowClear
        type="text"
        addon-before="会话参数"
        placeholder="group_concat_max_len=4194304;sql_mode=''"
        v-model="sessionVarsString"
      >
      </a-input>
      <a-divider type="vertical"></a-divider>
      <span>
        字符集
        <a-select
          placeholder="选择字符集,默认utf8"
          style="width: 15%"
          v-model="selectedCharacterSet"
          @change="saveCharacterToCache"
        >
          <a-select-option v-for="item in characterSets" :key="item.key" :value="item.value">{{
            item.key
          }}</a-select-option>
        </a-select>
      </span>
      <a-divider type="vertical"></a-divider>
      <span
        >主题
        <a-select placeholder="选择主题" style="width: 15%" v-model="selectedCodeTheme" @change="setCodeTheme">
          <a-select-option v-for="item in CodeThemes" :key="item.key" :value="item.value">{{
            item.key
          }}</a-select-option>
        </a-select>
      </span>
    </div>
    <a-spin :spinning="bTableLoading" tip="Loading..." id="cmHeight">
      <EditCodeMirrorComponent ref="EditCodeMirrorComponent" @execEvent="ExecEvent"></EditCodeMirrorComponent>
      <div style="margin-top: 4px">
        <a-button-group>
          <a-tooltip>
            <template slot="title">每次仅允许执行一条SQL，Mac按command+enter组合键，Windows按ctrl+enter组合键可以快速执行选中的SQL</template>
            <a-button @click="executeSQL()" icon="thunderbolt">执行SQL</a-button>
          </a-tooltip>
        </a-button-group>
        <a-button-group class="btn-right-group">
          <a-tooltip>
            <template slot="title">格式化SQL语句（建议每条语句后跟上分号，避免格式化后错乱）</template>
            <a-button @click="formatSQL()" icon="edit">格式化</a-button>
          </a-tooltip>
          <a-tooltip>
            <template slot="title">生成当前选中数据库的数据字典（新窗口打开）</template>
            <a-button @click="loadDBDict()" icon="copy">数据字典</a-button>
          </a-tooltip>
        </a-button-group>
      </div>
      <a-card class="box-card">
        <span v-html="responseMsg"></span>
      </a-card>
    </a-spin>
    <EditDBDictVue ref="EditDBDictVue"></EditDBDictVue>
  </div>
</template>

<script>
import { format } from 'sql-formatter'
import { executeMySQLQueryApi, executeClickHouseQueryApi } from '@/api/das'

// 导入子组件
import EditCodeMirrorComponent from './EditCodeMirror'
import EditDBDictVue from './EditDBDict'

export default {
  props: {
    selectedSchema: Object,
    tabIndex: Number,
  },
  components: {
    EditCodeMirrorComponent,
    EditDBDictVue,
  },
  computed: {
    characterCacheName: function () {
      return 'userCharacterSet-' + this.tabIndex
    },
    sessionVarsName: function () {
      return 'userSessionVars-' + this.tabIndex
    },
    characterSet: function () {
      return {
        character_set_client: this.selectedCharacterSet,
        character_set_connection: this.selectedCharacterSet,
        character_set_results: this.selectedCharacterSet,
      }
    },
    codeThemeName: function () {
      return 'codeTheme-' + this.tabIndex
    },
  },
  data() {
    return {
      responseMsg: '',
      bTableLoading: false,
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
    }
  },
  methods: {
    // 设置tab补全
    setTabCompletion(data) {
      this.$nextTick(() => {
        this.$refs.EditCodeMirrorComponent.setOption('hintOptions', data)
      })
    },
    // 解析用户输入的session vars
    parserSessionVars() {
      this.savesessionVarsToCache()
      var sessionVars = {}
      if (this.sessionVarsString.length > 0) {
        for (const [i, v] of this.sessionVarsString.split(';').entries()) {
          let sessionVar = v.split('=')
          if (sessionVar.length == 2) {
            sessionVars[sessionVar[0]] = sessionVar[1]
          }
        }
      }
      return sessionVars
    },
    // 执行mysql/tidb/aurora查询
    executeMySQLQuery(data) {
      // params append字符集 && session variables
      data['params'] = { ...this.characterSet, ...this.parserSessionVars() }
      // 响应消息
      var resMsgs = []
      this.bTableLoading = true
      executeMySQLQueryApi(data)
        .then((response) => {
          if (response.code === '0000') {
            resMsgs.push('结果: 执行成功')
            resMsgs.push(`耗时: ${response.data.duration}`)
            resMsgs.push(`SQL: ${response.data.sqltext}`)
            this.$message.success('执行成功')
            // 渲染数据
            this.$emit('showBootstrapTable', true)
            this.$emit('renderBootstrapTable', response.data)
          } else {
            resMsgs.push('结果: 执行失败')
            resMsgs.push(`错误: ${response.message}`)
            this.$emit('showBootstrapTable', false)
            this.$emit('close')
            if (response.message.includes('sessionid')) {
              this.$message.error('执行失败，认证过期，3s后系统自动刷新页面，请于刷新后重新执行')
              setTimeout(() => {
                // 刷新当前页面
                this.$router.go(0)
              }, 3000)
            } else {
              this.$message.error('执行失败')
            }
          }
          resMsgs.push(`请求ID: ${response.request_id}`)
          this.responseMsg = resMsgs.join('<br>')
        })
        .finally(() => {
          this.bTableLoading = false
        })
    },
    // 执行clickhouse查询
    executeClickHouseQuery(data) {
      // ClickHouse没有编码的概念
      // params append session variables
      data['params'] = { ...this.parserSessionVars() }
      // 响应消息
      var resMsgs = []
      this.bTableLoading = true
      executeClickHouseQueryApi(data)
        .then((response) => {
          if (response.code === '0000') {
            resMsgs.push('结果: 执行成功')
            resMsgs.push(`耗时: ${response.data.duration}`)
            resMsgs.push(`SQL: ${response.data.sqltext}`)
            this.$message.success('执行成功')
            // 渲染数据
            this.$emit('showBootstrapTable', true)
            this.$emit('renderBootstrapTable', response.data)
          } else {
            resMsgs.push('结果: 执行失败')
            resMsgs.push(`错误: ${response.message}`)
            this.$emit('showBootstrapTable', false)
            if (response.message.includes('sessionid')) {
              this.$message.error('执行失败，认证过期，3s后系统自动刷新页面，请于刷新后重新执行')
              setTimeout(() => {
                // 刷新当前页面
                this.$router.go(0)
              }, 3000)
            } else {
              this.$message.error('执行失败')
            }
          }
          resMsgs.push(`请求ID: ${response.request_id}`)
          this.responseMsg = resMsgs.join('<br>')
        })
        .finally(() => {
          this.bTableLoading = false
        })
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
      this.saveCodeToCache()
      // 判断是否选择了DB
      if (Object.keys(this.selectedSchema).length === 0) {
        this.$message.warning('请先选择左侧的库')
        return
      }
      // 获取选择的SQL
      var sqltext = this.$refs.EditCodeMirrorComponent.getSelection()
      // 判断是否输入
      if (sqltext.length == 0) {
        this.$message.warning('请鼠标选中要执行的SQL')
        return
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
      }

      // 根据db type调用不同的执行接口
      var dbType = this.selectedSchema['db_type'].toLowerCase()
      if (dbType === 'tidb' || dbType === 'mysql') {
        this.executeMySQLQuery(data)
      }
      if (dbType === 'clickhouse') {
        this.executeClickHouseQuery(data)
      }
    },
    // 格式化SQL
    formatSQL() {
      this.saveCodeToCache()
      var sql = this.$refs.EditCodeMirrorComponent.getValue()
      this.$nextTick(() => {
        this.$refs.EditCodeMirrorComponent.setValue(format(sql, { language: 'mysql' }))
      })
    },
    // 保存到cache
    saveCodeToCache() {
      // 保存用户的输入到本地浏览器
      this.$nextTick(() => {
        this.$refs.EditCodeMirrorComponent.saveCache(`dms-codemirror-${this.tabIndex}`)
      })
    },
    // 加载缓存
    loadCodeFromcCache() {
      this.$nextTick(() => {
        this.$refs.EditCodeMirrorComponent.loadCache(`dms-codemirror-${this.tabIndex}`)
      })
    },
    // 设置字符集
    saveCharacterToCache() {
      localStorage.setItem(this.characterCacheName, this.selectedCharacterSet)
    },
    // 加载字符集
    loadCharacterFromCache() {
      var cache = localStorage.getItem(this.characterCacheName)
      if (cache != null) {
        this.selectedCharacterSet = cache
      }
    },
    // 设置sessionVarsString
    savesessionVarsToCache() {
      localStorage.setItem(this.sessionVarsName, this.sessionVarsString)
    },
    // 加载字符集
    loadSessionVarsFromCache() {
      var cache = localStorage.getItem(this.sessionVarsName)
      if (cache != null) {
        this.sessionVarsString = cache
      }
    },
    // 设置codemirror主题
    setCodeTheme() {
      localStorage.setItem(this.codeThemeName, this.selectedCodeTheme)
      this.$nextTick(() => {
        this.$refs.EditCodeMirrorComponent.setTheme(this.selectedCodeTheme)
      })
    },
    // 加载codemirror主题
    loadCodeThemeFromCache() {
      var cache = localStorage.getItem(this.codeThemeName)
      if (cache != null) {
        this.selectedCodeTheme = cache
        this.$nextTick(() => {
          this.$refs.EditCodeMirrorComponent.setTheme(cache)
        })
      }
    },
    // 加载数据字典
    loadDBDict() {
      // 判断是否选择了DB
      if (Object.keys(this.selectedSchema).length === 0) {
        this.$message.warning('请先选择左侧的库')
        return
      } else {
        this.$refs.EditDBDictVue.show(this.selectedSchema)
      }
    },
    ExecEvent() {
      this.executeSQL()
    },
  },
  mounted() {
    this.loadCodeFromcCache()
    this.loadCharacterFromCache()
    this.loadCodeThemeFromCache()
    this.loadSessionVarsFromCache()
  },
}
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
</style>