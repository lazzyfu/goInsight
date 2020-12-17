<template>
  <a-card
    v-loading="bTableLoading"
    element-loading-text="玩命查询中，请稍后..."
    element-loading-spinner="el-icon-loading"
  >
    <a-row>
      <div class="query-header">
        <span slot="title">
          <span>节点：{{ hostname }}</span>
          <el-divider direction="vertical"></el-divider>
          <span>当前所在库: {{ selectedDB.split("___")[1] }}</span>
          <el-divider direction="vertical"></el-divider>
          <span
            >字符集：
            <el-select
              placeholder="选择字符集"
              style="width: 15%"
              size="mini"
              v-model="selectedCharacter"
            >
              <el-option
                v-for="item in characterSet"
                :key="item.key"
                :label="item.key"
                :value="item.value"
              ></el-option>
            </el-select>
          </span>
          <el-divider direction="vertical"></el-divider>
          <span
            >主题：
            <el-select
              placeholder="选择主题"
              style="width: 15%"
              size="mini"
              v-model="selectedTheme"
              @change="setCodeTheme"
            >
              <el-option
                v-for="item in themes"
                :key="item.key"
                :label="item.key"
                :value="item.value"
              ></el-option>
            </el-select>
          </span>
        </span>
      </div>
    </a-row>
    <div class="box" ref="box">
      <div class="left">
        <a-row>
          <div class="table-page-search-wrapper">
            <a-form
              layout="inline"
              :form="form"
              @keyup.enter.native="handleSearch"
            >
              <a-form-item>
                <a-select
                  placeholder="请选择库名..."
                  size="small"
                  show-search
                  style="width: 100%"
                  @change="getTablesNode"
                >
                  <a-icon
                    slot="suffixIcon"
                    type="database"
                    style="color: green"
                  />
                  <a-select-option
                    v-for="s in databases"
                    :key="s.key"
                    :value="s.key"
                    >{{ s.title }}</a-select-option
                  >
                </a-select>
              </a-form-item>
            </a-form>
          </div>

          <div>
            <a-input-search
              v-if="showSearch"
              size="small"
              style="margin-bottom: 8px"
              placeholder="输入要搜索的表名"
              @change="onSearch"
            />
          </div>

          <div
            v-loading="treeLoading"
            element-loading-text="玩命加载中..."
            element-loading-spinner="el-icon-loading"
          >
            <div class="tree-container">
              <div class="block">
                <el-tree
                  node-key="key"
                  :data="treeData"
                  :props="defaultProps"
                  @node-contextmenu="rightClick"
                  @node-click="closeContextmenu"
                  class="tree filter-tree"
                  :render-content="renderContent"
                  :style="{ height: screenHeight }"
                >
                </el-tree>
              </div>
            </div>
          </div>
          <!-- 右键菜单 -->
          <v-contextmenu ref="contextmenu">
            <v-contextmenu-item @click="viewTableStruc"
              >查看表结构</v-contextmenu-item
            >
            <v-contextmenu-item @click="viewTableBase"
              >查看表信息</v-contextmenu-item
            >
          </v-contextmenu>
        </a-row>
      </div>

      <div class="resize" title="收缩侧边栏">⋮</div>

      <div class="right">
        <a-row>
          <div id="cmHeight">
            <div style="margin-bottom: 6px; margin-top: 4px">
              <el-button-group style="margin-left: 2px; margin-right: 2px">
                <el-tooltip
                  class="item"
                  effect="dark"
                  content="每次仅允许执行一条SQL"
                  placement="top"
                >
                  <el-button
                    size="mini"
                    @click="executeSQL()"
                    icon="el-icon-video-play"
                    >执行SQL</el-button
                  >
                </el-tooltip>

                <el-tooltip
                  class="item"
                  effect="dark"
                  content="格式化SQL"
                  placement="top"
                >
                  <el-button
                    size="mini"
                    @click="formatSQL()"
                    icon="el-icon-edit"
                    >格式化</el-button
                  >
                </el-tooltip>
                <el-tooltip
                  class="item"
                  effect="dark"
                  content="查看执行过的历史SQL"
                  placement="top"
                >
                  <el-button
                    size="mini"
                    @click="getHistorySQL()"
                    icon="el-icon-document"
                    >我的SQL</el-button
                  >
                </el-tooltip>
                <el-button
                  size="mini"
                  @click="loadDBDict()"
                  icon="el-icon-notebook-1"
                  >数据字典</el-button
                >
                <el-button
                  size="mini"
                  @click="dialogVisible = true"
                  icon="el-icon-document"
                  >使用技巧</el-button
                >
              </el-button-group>
            </div>
            <div>
              <codemirror
                ref="myCm1"
                v-model="code1"
                :options="cmOptions1"
                @ready="onCmReady1"
              ></codemirror>
            </div>

            <div style="margin-top: 4px">
              <codemirror
                ref="myCm2"
                v-model="code2"
                :options="cmOptions2"
                @ready="onCmReady2"
              ></codemirror>
            </div>
          </div>
        </a-row>
      </div>
    </div>

    <a-row v-show="visibleResult">
      <a-tabs default-active-key="1">
        <a-tab-pane key="1" tab="数据集">
          <table
            data-click-to-select="true"
            data-show-copy-rows="true"
            data-pagination="true"
            data-page-number="1"
            data-page-size="10"
            data-side-pagination="client"
            data-page-list="[10, 25, 50, 100, 200, All]"
            data-resizable="true"
            id="bTable"
          ></table>
        </a-tab-pane>
      </a-tabs>
    </a-row>

    <!-- 使用小技巧 -->
    <template>
      <el-dialog title="使用小技巧" :visible.sync="dialogVisible">
        <p>
          * 默认返回100条记录，LIMIT
          N最大返回2000条记录。单条SQL查询超时为600秒。已开启SQL审计
        </p>
        <p>* 支持左右拉伸，点击中间蓝色部分即可实现左右拉伸</p>
        <p>* 支持SQL输入框上下拉伸（鼠标置于输入框的右下角进行拖动）</p>
        <p>* 如果SQL没有LIMIT，将会被自动加上LIMIT N</p>
        <p>
          * 支持LIMIT offset, row_count和LIMIT row_count OFFSET
          offset，row_count超过限制，将被重写
        </p>
        <p>* 点击结果集左侧的选择框，即可实现copy</p>
        <p>* 右键点击表名可查看表结构和表基础信息</p>
      </el-dialog>
    </template>

    <!-- 表结构 -->
    <template>
      <drawerTableStructure
        ref="drawerTableStructureChild"
        :visible="tableStructureVisible"
        @close="tableStructureVisible = false"
      ></drawerTableStructure>
    </template>

    <!-- 表基础信息 -->
    <template>
      <drawerTableBaseInfo
        ref="drawerTableBaseInfoChild"
        :visible="tableBaseInfoVisible"
        @close="tableBaseInfoVisible = false"
      ></drawerTableBaseInfo>
    </template>

    <!-- 我的SQL -->
    <template>
      <drawerHistorySQL
        ref="drawerHistorySQLChild"
        :visible="historySqlVisible"
        @close="historySqlVisible = false"
      ></drawerHistorySQL>
    </template>

    <!-- 数据字典 -->
    <template>
      <drawerDBDict
        ref="drawerDBDictChild"
        :visible="DBDictVisible"
        @close="DBDictVisible = false"
      ></drawerDBDict>
    </template>
  </a-card>
</template>

<script>
import { mapActions, mapGetters } from "vuex";

import { getQueryTree, ExecuteQuery, deleteQuerySession } from "@/api/sqlquery";

import sqlFormat from "sql-formatter";

import $ from "jquery";

// 主题
import "codemirror/theme/dracula.css";
import "codemirror/theme/solarized.css";
import "codemirror/theme/eclipse.css";
import "codemirror/theme/monokai.css";
import "codemirror/theme/oceanic-next.css";

// mode
import "codemirror/mode/sql/sql.js";

// addon
import "codemirror/addon/selection/active-line";
import "codemirror/addon/display/autorefresh";

// 提示和自动补全
import "codemirror/addon/hint/show-hint";
import "codemirror/addon/hint/show-hint.css";
import "codemirror/addon/hint/anyword-hint";
import "codemirror/addon/hint/sql-hint";
import "codemirror/addon/comment/comment"
import "codemirror/addon/edit/matchbrackets"
import "codemirror/addon/edit/closebrackets";

// 编辑器类型
import "codemirror/keymap/sublime";

import elementResizeDetectorMaker from "element-resize-detector";

// bootstrap-table
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap-table/dist/bootstrap-table.min.css";
import "bootstrap/dist/js/bootstrap.min.js";
import "bootstrap-table/dist/bootstrap-table.min.js";
import "bootstrap-table/dist/locale/bootstrap-table-zh-CN.min.js";
import "bootstrap-table/dist/extensions/auto-refresh/bootstrap-table-auto-refresh.min.js";
import "bootstrap-table/dist/extensions/copy-rows/bootstrap-table-copy-rows.min.js";
// resize columns width
import "jquery-resizable-columns/dist/jquery.resizableColumns.css";
import "jquery-resizable-columns/dist/jquery.resizableColumns.min.js";
import "bootstrap-table/dist/extensions/resizable/bootstrap-table-resizable.min.js";

// 导入子组件
import drawerHistorySQL from "./myHistorySql";
import drawerDBDict from "./dbDict";
import drawerTableStructure from "./tableStructure";
import drawerTableBaseInfo from "./tableBaseInfo";

// 检查是否https连接
let protocol = "ws://";
if (window.location.protocol === "https:") {
  protocol = "wss://";
}

/**
 * 生成随机字符串
 */
function randomString(e) {
  e = e || 32;
  var t = "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678",
    a = t.length,
    n = "";
  for (let i = 0; i < e; i++) n += t.charAt(Math.floor(Math.random() * a));
  return n;
}

export default {
  components: {
    drawerHistorySQL,
    drawerDBDict,
    drawerTableStructure,
    drawerTableBaseInfo,
  },
  computed: {
    ...mapGetters(["sqlState"]),
    codemirror1() {
      return this.$refs.myCm1.codemirror;
    },
    codemirror2() {
      return this.$refs.myCm2.codemirror;
    },
  },
  data() {
    return {
      defaultProps: {
        children: "children",
        label: "label",
      },
      treeData: [],
      showSearch: false,
      treeLoading: false,
      selectedKeys: [], // 设置选中key的节点
      databases: [],
      hostname: "None",
      selectedDB: "",
      screenHeight: "",
      // codemirror
      code1: "",
      code2: "",
      cm1Height: 450,
      cm2Height: 150,
      cmOptions1: {
        mode: "text/x-mysql",
        indentUnit: 2,
        tabSize: 2,
        indentWithTabs: true,
        smartIndent: true,
        autoRefresh: true,
        lineNumbers: true,
        styleActiveLine: true,
        autoCloseBrackets: true,
        matchBrackets: true,
        lineWrapping: true, // 自动换行
        autofocus: true,
        resetSelectionOnContextMenu: false,
        showCursorWhenSelecting: true,
        keyMap: "sublime", // 编辑器模式
      },
      cmOptions2: {
        mode: "text/x-mysql",
        autoRefresh: true,
        lineWrapping: true, // 自动换行
        readOnly: true,
        placeholder: "SQL执行消息...",
      },
      characterSet: [
        { key: "utf8", value: "utf8" },
        { key: "utf8mb4", value: "utf8mb4" },
        { key: "latin1", value: "latin1" },
        { key: "gbk", value: "gbk" },
      ],
      selectedCharacter: "utf8",
      themes: [
        { key: "default", value: "default" },
        { key: "eclipse", value: "eclipse" },
        { key: "solarized dark", value: "solarized dark" },
        { key: "dracula", value: "dracula" },
        { key: "monokai", value: "monokai" },
        { key: "oceanic-next", value: "oceanic-next" },
      ],
      selectedTheme: "default",
      form: this.$form.createForm(this),
      websocket: {
        path: `${protocol}/${window.location.host}/ws/sql/${this.$store.getters.userInfo.username}/`,
        socket: "",
      },
      visibleResult: false,
      // table
      bTableLoading: false,
      // 表结构
      tableStructureVisible: false,
      // 表基础信息
      tableBaseInfoVisible: false,
      // history sql
      historySqlVisible: false,
      // dbdict
      DBDictVisible: false,
      //dialog
      dialogVisible: false,
    };
  },
  methods: {
    ...mapActions(["storeQueryHash", "storeUserIuput", "storeUserTheme"]),
    // 自动补全
    onCmReady1(cm) {
      // 获取element form表单的高度，将codemirror的高度设置和表单一致
      cm.setSize("height", `${this.cm1Height}px`);
      // 读取保存在本地浏览器的数据
      cm.setValue(this.sqlState.sqlquery.user_input);
      // 设置theme
      cm.setOption("theme", this.sqlState.sqlquery.user_theme);
      this.selectedTheme = this.sqlState.sqlquery.user_theme;
      // 设置自动补全
      cm.on("keypress", () => {
        cm.showHint({ completeSingle: false });
      });
    },
    onCmReady2(cm) {
      // 设置theme
      cm.setOption("theme", this.sqlState.sqlquery.user_theme);
      // 获取element form表单的高度，将codemirror的高度设置和表单一致
      cm.setSize("height", `${this.cm2Height}px`);
    },
    beforeunloadFn() {
      const data = {
        query_hash: this.sqlState.sqlquery.query_hash,
      };
      deleteQuerySession(data).then((response) => {
        console.log("response: ", response);
      });
    },
    // 初始化websocket
    init_websocket() {
      if (typeof WebSocket === "undefined") {
        this.$message.error("您的浏览器不支持websocket");
      }
      // 实例化websocket
      this.websocket.socket = new WebSocket(this.websocket.path);
      // 监控websocket连接
      this.websocket.socket.onopen = this.socketOnOpen;
      // 监听socket错误信息
      this.websocket.socket.onerror = this.socketOnError;
      // 监听socket消息
      this.websocket.socket.onmessage = this.socketOnMessage;
    },
    socketOnOpen() {
      //
    },
    socketOnError() {
      this.init_websocket();
    },
    socketOnMessage(msg) {
      // 接收socket信息
      let result = JSON.parse(msg.data);
      if (result.flag === this.$store.getters.userInfo.username) {
        if (["query"].includes(result.data.type)) {
          this.codemirror2.setValue(result.data.msg);
        }
      }
    },
    socketClose() {
      // 关闭socket
      this.websocket.socket.close();
    },
    renderProcesslist(data) {
      // 渲染执行SQL时，show processlist的输出
      let html = "当前SQL SESSION ID的SHOW PROCESSLIST实时输出：";
      for (let key in data) {
        html += "\n" + key + ": " + data[key];
      }
      return html;
    },
    // 设置主题
    setCodeTheme() {
      this.storeUserTheme(this.selectedTheme);
      this.codemirror1.setOption("theme", this.sqlState.sqlquery.user_theme);
      this.codemirror2.setOption("theme", this.sqlState.sqlquery.user_theme);
    },
    // 格式化SQL
    formatSQL() {
      this.$message.info("正在格式化SQL，请稍等(SQL越大，响应越慢)");
      const sqlContent = this.codemirror1.getValue();
      this.codemirror1.setValue(sqlFormat.format(sqlContent, { indent: "  " }));
    },
    // bootstrap-table
    renderBtable(data) {
      $("#bTable")
        .bootstrapTable("destroy")
        .bootstrapTable({
          columns: data.columns,
          data: data.data,
          locale: "zh-CN",
          search: true,
          showRefresh: true,
          showToggle: true,
          showColumns: true,
          lineWrapping: true,
          matchBrackets: true,
          classes: "table table-hover table-bordered table-striped",
          iconSize: "sm",
          cache: false,
          showFullscreen: true,
          rowStyle: function rowStyle(row, index) {
            return {
              css: { "font-size": "12px" },
            };
          },
        });
      // 更改icon的大小
      $("#bTable").bootstrapTable("refreshOptions", {
        iconSize: "sm",
      });
    },
    // 执行SQL
    executeSQL() {
      // 保存用户的输入到本地浏览器
      this.storeUserIuput(this.codemirror1.getValue());
      // 判断是否选中SQL
      if (!this.codemirror1.somethingSelected()) {
        this.$message.warning("请选中要执行的SQL，单次仅允许执行一条SQL");
        return false;
      }
      // 获取选中的SQL
      const sql = this.codemirror1.getSelection();
      // 保存querHash到本地浏览器
      this.storeQueryHash(randomString(64));
      const data = {
        key: this.selectedDB,
        sql: sql,
        character: this.selectedCharacter,
        query_hash: this.sqlState.sqlquery.query_hash,
      };
      // 设置loading状态
      this.visibleResult = true;
      this.bTableLoading = true;
      // 先销毁表格，避免loading位置异常
      $("#bTable").bootstrapTable("destroy");

      ExecuteQuery(data)
        .then((response) => {
          // 隐藏loading状态
          this.bTableLoading = false;
          // 滚动到页面的底部
          this.$nextTick(() => {
            document.scrollingElement.scrollTop =
              document.scrollingElement.scrollHeight;
          });
          if (response.code == "0001") {
            this.$notify.error({
              title: "错误",
              message: response.message,
            });
            return false;
          }
          // 渲染表格
          this.renderBtable(response.data);
        })
        .catch((error) => {
          this.bTableLoading = false;
          this.$notify.error({
            title: "错误",
            message: error.message,
          });
        });
      // 滚动到页面的底部
      this.$nextTick(() => {
        document.scrollingElement.scrollTop =
          document.scrollingElement.scrollHeight;
      });
    },
    // 获取DB
    getRootNode() {
      getQueryTree({ key: "root" }).then((response) => {
        this.databases = response.data;
      });
    },
    // 渲染tree icon
    renderContent(h, { node, data, store }) {
      return (
        <span>
          <i class={data.icon} style="color: #409EFF"></i>
          <span> {node.label}</span>
        </span>
      );
    },
    // 查找主机名
    findHostName(value) {
      for (let item of this.databases) {
        if (item.key == value) {
          return item.hostname;
        }
      }
    },
    // 切换数据库加载表
    getTablesNode(value) {
      this.showSearch = true;
      this.treeLoading = true;
      this.selectedDB = value;
      this.hostname = this.findHostName(value);
      setTimeout(() => {
        getQueryTree({ key: value })
          .then((response) => {
            if (response.code === "0001") {
              this.$notify.error({
                title: "加载表失败",
                message: response.message,
              });
              return false;
            }
            // 渲染树结构
            this.treeData = response.data.tree_data;
            // 加载表列自动提示
            this.codemirror1.setOption(
              "hintOptions",
              response.data.tab_completion
            );
            // 当库为空
            if (this.treeData.length === 0) {
              this.$notify.warning({
                title: "警告",
                message:
                  "未从当前库" +
                  value.split("___")[1] +
                  "找到表，请确认库不为空",
              });
            }
          })
          .catch((err) => {
            this.$notify.error({
              title: err.response.status,
              message: "请稍后刷新页面重试",
            });
            this.treeData = [];
          })
          .finally(() => {
            this.treeLoading = false;
          });
      }, 500);
    },
    // 右键点击弹出菜单
    rightClick(MouseEvent, object, Node, element) {
      this.closeContextmenu();
      this.selectedKeys = Node.key;
      if (this.selectedKeys.split("___").length === 3) {
        const postition = {
          position: "absolute",
          top: event.pageY,
          left: event.pageX,
        };
        // 显示菜单
        this.$refs.contextmenu.show(postition);
      }
    },
    // 关闭菜单
    closeContextmenu() {
      this.$refs.contextmenu.hide();
    },
    // 搜索表
    onSearch(e) {
      const value = e.target.value;
      if (!value) {
        this.getTablesNode(this.selectedDB);
      }
      const searchResult = [];
      this.treeData
        .map((item) => {
          if (item.label.indexOf(value) > -1) {
            searchResult.push(item);
          }
          return null;
        })
        .filter((item, i, self) => item && self.indexOf(item) === i);
      this.treeData = searchResult;
    },
    // 查看表结构
    viewTableStruc() {
      this.tableStructureVisible = true;
      this.$refs.drawerTableStructureChild.fetchData(this.selectedKeys);
    },
    // 查看表基本信息
    viewTableBase() {
      this.tableBaseInfoVisible = true;
      this.$refs.drawerTableBaseInfoChild.fetchData(this.selectedKeys);
    },
    // 获取历史SQL
    getHistorySQL() {
      this.historySqlVisible = true;
      this.$refs.drawerHistorySQLChild.fetchData();
    },
    // 加载数据字典
    loadDBDict() {
      if (!this.selectedDB) {
        this.$message.error("请先选择左侧的库名", 3);
        return false;
      }
      this.DBDictVisible = true;
      this.$refs.drawerDBDictChild.fetchData(this.selectedDB);
    },
    // 左右拖动
    dragControllerDiv: function () {
      var resize = document.getElementsByClassName("resize");
      var left = document.getElementsByClassName("left");
      var right = document.getElementsByClassName("right");
      var box = document.getElementsByClassName("box");
      for (let i = 0; i < resize.length; i++) {
        // 鼠标按下事件
        resize[i].onmousedown = function (e) {
          //颜色改变提醒
          var startX = e.clientX;
          resize[i].left = resize[i].offsetLeft;
          // 鼠标拖动事件
          document.onmousemove = function (e) {
            var endX = e.clientX;
            var moveLen = resize[i].left + (endX - startX); // （endx-startx）=移动的距离。resize[i].left+移动的距离=左边区域最后的宽度
            var maxT = box[i].clientWidth - resize[i].offsetWidth; // 容器宽度 - 左边区域的宽度 = 右边区域的宽度

            if (moveLen < 80) moveLen = 80; // 左边区域的最小宽度为80px
            if (moveLen > maxT - 300) moveLen = maxT - 300; //右边区域最小宽度为300px
            resize[i].style.left = moveLen; // 设置左侧区域的宽度

            for (let j = 0; j < left.length; j++) {
              left[j].style.width = moveLen + "px";
              right[j].style.width = box[i].clientWidth - moveLen - 10 + "px";
            }
          };
          // 鼠标松开事件
          document.onmouseup = function (evt) {
            document.onmousemove = null;
            document.onmouseup = null;
            resize[i].releaseCapture && resize[i].releaseCapture(); //当你不在需要继续获得鼠标消息就要应该调用ReleaseCapture()释放掉
          };
          resize[i].setCapture && resize[i].setCapture(); //该函数在属于当前线程的指定窗口里设置鼠标捕获
          return false;
        };
      }
    },
  },
  created() {
    // 增加页面刷新和关闭监听
    window.addEventListener("onbeforeunload", this.beforeunloadFn());
  },
  mounted() {
    this.dragControllerDiv();
    this.init_websocket();
    this.getRootNode();
    // yarn add element-resize-detector
    // 实时获取codemirror的高度，并改变左侧表区域的高度
    const _this = this;
    const erd = elementResizeDetectorMaker();
    erd.listenTo(document.getElementById("cmHeight"), (element) => {
      _this.$nextTick(() => {
        this.screenHeight = element.offsetHeight - 74 + "px";
      });
    });
  },
  destroyed() {
    // 关闭websocket
    this.socketClose();
    // 卸载页面刷新和关闭监听
    window.removeEventListener("onbeforeunload", this.beforeunloadFn());
  },
};
</script>

<style lang='less' scoped>
.tree {
  overflow: scroll;
  min-height: 100px;
  font-size: 13px;
  border-radius: 5px;
  border-left-width: 0px;
  border-right-width: 0px;
}

.table-page-search-wrapper {
  .ant-form-inline {
    .ant-form-item {
      display: flex;
      margin-bottom: 8px;
      margin-right: 0;
    }
  }

  .table-page-search-submitButtons {
    display: block;
    margin-bottom: 24px;
    white-space: nowrap;
  }
}

::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-thumb {
  border-radius: 8px;
  background-color: #cbc7cc;
}

::-webkit-scrollbar-track {
  border-radius: 8px;
  background: rgba(0, 0, 0, 0);
}

/* 树形结构节点添加连线 */
.tree /deep/ .el-tree-node {
  position: relative;
  // padding-left: 16px;
}

.tree /deep/ .el-tree-node__children {
  padding-left: 16px;
}

.tree /deep/ .el-tree-node :last-child:before {
  height: 38px;
}

.tree /deep/ .el-tree > .el-tree-node:before {
  border-left: none;
}

.tree-container /deep/ .el-tree > .el-tree-node:after {
  border-top: none;
}

.tree /deep/ .el-tree-node:before {
  content: "";
  left: -4px;
  position: absolute;
  right: auto;
  border-width: 1px;
}

.tree /deep/ .el-tree-node:after {
  content: "";
  left: -4px;
  position: absolute;
  right: auto;
  border-width: 1px;
}
.tree /deep/ .el-tree-node__expand-icon.is-leaf {
  display: none;
}

.tree /deep/ .el-tree-node:before {
  border-left: 1px dashed #b8b9bb;
  bottom: 0px;
  height: 100%;
  top: -26px;
  width: 1px;
}

.tree /deep/ .el-tree-node:after {
  border-top: 1px dashed #b8b9bb;
  height: 20px;
  top: 12px;
  width: 24px;
}

/* 拖拽相关样式 */
/*包围div样式*/
.box {
  width: 100%;
  height: 100%;
  overflow: hidden;
  display: flex;
  display: -webkit-flex; /* Safari */
  justify-content: center;

  .left {
    width: 25%;
    height: 100%;
  }
  /*拖拽区div样式*/
  .resize {
    cursor: col-resize;
    border-radius: 5px;
    margin-bottom: 100px;
    height: 68px;
    width: 6px;
    margin-top: 100px;
    color: white;
    display: flex;
    display: -webkit-flex; /* Safari */
    align-items: center;
    justify-content: center;
    background-color: teal;
  }
  /*右侧div样式*/
  .right {
    width: 75%;
    height: 100%;
  }
}

.query-header {
  width: 100%;
  padding: 4px 6px;
  margin: 0;
  -webkit-box-sizing: border-box;
  box-sizing: border-box;
  border-radius: 4px;
  position: relative;
  // background-color: #f4f4f5;
  background-color: #e4e7ed;
  color: #909399;
  overflow: hidden;
  opacity: 1;
  display: -webkit-box;
  display: -ms-flexbox;
  display: flex;
  -webkit-box-align: center;
  -ms-flex-align: center;
  align-items: center;
  -webkit-transition: opacity 0.2s;
  transition: opacity 0.2s;
}
</style>