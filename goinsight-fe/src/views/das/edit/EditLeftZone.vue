<template>
  <div shadow="hover">
    <div class="table-page-search-wrapper">
      <a-form layout="inline" :form="form" @keyup.enter.native="handleSearch">
        <a-form-item>
          <a-select
            placeholder="请选择库名..."
            show-search
            style="width: 90%"
            @change="getTables"
            option-label-prop="label"
            v-model="bindTitle"
          >
            <a-icon slot="suffixIcon" type="database" style="color: green" />
            <a-select-option
              v-for="(s, index) in schemas"
              :key="index"
              :label="`${s.remark}:${s.schema}`"
              :value="`${s.instance_id}#${s.schema}#${s.db_type}`"
            >
              <span v-if="s.db_type.toLowerCase() == 'clickhouse'">
                <svg
                  t="1690439692107"
                  class="icon"
                  viewBox="0 0 1024 1024"
                  version="1.1"
                  xmlns="http://www.w3.org/2000/svg"
                  p-id="19360"
                  width="20"
                  height="20"
                >
                  <path
                    d="M864.213333 591.082667l61.226667 28.053333L517.333333 981.333333 106.666667 619.136l61.226666-28.053333 349.44 311.168 346.88-311.168zM440.832 318.144v385.152h-58.666667V318.144h58.666667z m104.576 0v385.152h-58.666667V318.144h58.666667z m102.037333 0v385.152h-58.666666V318.144h58.666666z m107.136 219.370667v114.773333h-58.666666v-114.773333h58.666666z m-415.786666-153.045334v267.818667h-58.666667V384.469333h58.666667z m415.786666 0v114.773334h-58.666666v-114.773334h58.666666zM517.333333 42.666667l408.106667 359.658666-61.226667 30.592L517.333333 121.749333 167.893333 432.917333 106.666667 402.346667 517.333333 42.666667z"
                    fill="#0e932e"
                    p-id="19361"
                  ></path>
                </svg>
              </span>
              <span v-else-if="s.db_type.toLowerCase() == 'mysql'">
                <svg
                  t="1690439650713"
                  class="icon"
                  viewBox="0 0 1024 1024"
                  version="1.1"
                  xmlns="http://www.w3.org/2000/svg"
                  p-id="14032"
                  width="20"
                  height="20"
                >
                  <path
                    d="M856.368 169.2H577.76a24 24 0 0 0 0 48h278.608a45.968 45.968 0 0 1 45.92 45.92v557.216a45.968 45.968 0 0 1-45.92 45.92h-253.6a1524.576 1524.576 0 0 1-68.304-57.328 24 24 0 1 0-32.096 35.696c8.304 7.456 16.704 14.752 24.784 21.632H299.136a45.968 45.968 0 0 1-45.92-45.92v-41.76a24 24 0 0 0-48 0v41.76a94.032 94.032 0 0 0 93.92 93.92h557.232a94.032 94.032 0 0 0 93.92-93.92V263.088a94.048 94.048 0 0 0-93.92-93.888z"
                    fill="#1296db"
                    p-id="14033"
                  ></path>
                  <path
                    d="M298.496 796.88a24 24 0 0 0 37.312-15.456 895.808 895.808 0 0 1 26.656-102.72 1102.672 1102.672 0 0 0 41.2 57.088 24 24 0 1 0 38-29.328 1078.192 1078.192 0 0 1-65.6-95.36 24 24 0 0 0-43.008 4.352c-1.024 2.816-20.272 56-34.48 112-27.088-39.632-57.504-121.008-6.544-259.984a24 24 0 0 0-6.016-25.696 341.856 341.856 0 0 1-78.8-120.336 233.216 233.216 0 0 0-42.288-76.064c-41.056-49.312-46.944-78.16-40.896-85.504 10.176-12.368 70.528 15.184 119.568 54.544a24 24 0 0 0 26 2.64c1.008-0.496 102.656-50.608 231.184 62.864a664.256 664.256 0 0 1 144.336 194.672 24.112 24.112 0 0 0 25.6 13.76c0.88-0.144 77.744-10.8 151.664 77.584-110.144 19.2-115.2 36.976-118.4 48a24 24 0 0 0 5.744 23.168c32.288 33.792 119.328 136.144 133.984 203.2a24 24 0 0 0 23.424 18.88 24.576 24.576 0 0 0 5.136-0.544 24 24 0 0 0 18.32-28.576c-15.472-70.752-87.888-160.368-124.72-202.416a921.184 921.184 0 0 1 102.72-20.448 24 24 0 0 0 16.608-36.96c-75.616-114.816-165.488-129.6-203.536-130.24a699.504 699.504 0 0 0-149.056-196.048C408.176 134.08 301.008 155.2 262.656 168.288 144.784 77.92 98.624 115.328 87.072 129.344c-26.8 32.48-12.992 81.856 41.04 146.752a185.744 185.744 0 0 1 33.6 60.496 377.6 377.6 0 0 0 80.144 128.768c-79.152 233.728 50.944 327.552 56.64 331.52z"
                    fill="#1296db"
                    p-id="14034"
                  ></path>
                  <path
                    d="M290.048 282.288a24.528 24.528 0 0 0-6.896 16.96 24 24 0 0 0 6.896 16.96 24.128 24.128 0 0 0 34.064 0 24.096 24.096 0 0 0 6.88-16.96 24.512 24.512 0 0 0-6.88-16.96 24.944 24.944 0 0 0-34.064 0z"
                    fill="#1296db"
                    p-id="14035"
                  ></path>
                  <path
                    d="M477.408 757.232a15.904 15.904 0 0 0-2.896-3.52 24.608 24.608 0 0 0-33.92 0 16.528 16.528 0 0 0-3.04 3.52 36.896 36.896 0 0 0-2.24 4.16c-0.48 1.44-0.944 3.04-1.264 4.48a24.464 24.464 0 0 0-0.496 4.8 24.992 24.992 0 0 0 1.76 9.12 26.336 26.336 0 0 0 5.28 7.84 22.4 22.4 0 0 0 7.68 5.12 23.008 23.008 0 0 0 9.28 1.92 24.256 24.256 0 0 0 16.96-7.04 22.704 22.704 0 0 0 5.12-7.84 22.208 22.208 0 0 0 1.92-9.12 21.92 21.92 0 0 0-1.92-9.28 20.8 20.8 0 0 0-2.224-4.16z"
                    fill="#1296db"
                    p-id="14036"
                  ></path>
                </svg>
              </span>
              <span v-else-if="s.db_type.toLowerCase() == 'tidb'">
                <svg
                  t="1690439481617"
                  class="icon"
                  viewBox="0 0 1024 1024"
                  version="1.1"
                  xmlns="http://www.w3.org/2000/svg"
                  p-id="5590"
                  width="20"
                  height="20"
                >
                  <path
                    d="M859.14 292.08L530.38 102.34c-11.27-5.64-24.42-5.64-35.69 0.94L165.92 293.02c-10.33 5.64-17.85 16.91-17.85 30.06v379.48c0 12.21 7.51 23.48 17.85 31l329.7 189.74c10.33 6.58 24.42 6.58 35.69 0l328.76-189.74c10.33-6.58 17.85-17.85 17.85-31V323.08c-0.94-13.15-8.45-25.36-18.78-31z m-52.61 390.76l-294 170.02-294-170.02V343.75l293.06-170.02 294 170.02v339.09h0.94z"
                    fill="#d81e06"
                    p-id="5591"
                  ></path>
                  <path
                    d="M284.28 478.07l123.99-72.33-6.58 309.97 93.93 50.73 1.88-402.97 117.42-63.87-97.69-53.54-232.95 128.68z"
                    fill="#d81e06"
                    p-id="5592"
                  ></path>
                  <path d="M549.16 766.44l96.75-46.03V323.08l-96.75 51.66z" fill="#d81e06" p-id="5593"></path>
                </svg>
              </span>
              {{ s.remark }}:{{ s.schema }}
              <i style="color: #c0c4cc" v-if="s.is_deleted === true">已删除</i>
            </a-select-option>
          </a-select>
          <a-tooltip title="动态刷新库名" placement="top">
            <a-button
              style="width: 10%; font-size: 16px"
              icon="reload"
              type="link"
              :loading="refreshLoading"
              @click="refreshSchemas"
            ></a-button>
          </a-tooltip>
        </a-form-item>
      </a-form>
    </div>
    <a-input-search
      v-if="showSearch"
      style="margin-bottom: 8px"
      placeholder="输入要搜索的表名，然后点击回车"
      @search="onSearch"
    />
    <a-spin :spinning="treeLoading" tip="Loading...">
      <div class="tree-container">
        <div class="block">
          <a-tree :tree-data="treeData" show-line show-icon class="tree filter-tree" :style="{ height: screenHeight }">
            <a-icon slot="child" type="profile" />
            <template #title="{ key: treeKey, title }">
              <span v-if="title.split('#').length === 2">
                <a-tooltip v-if="title.split('#')[1] === 'allow'">
                  <template slot="title">您有该表的访问权限</template>
                  <a-icon type="check-circle" theme="twoTone" two-tone-color="#52c41a" class="icon-align" />
                </a-tooltip>
                <a-tooltip v-else>
                  <template slot="title">您有没有该表的访问权限</template>
                  <a-icon type="close-circle" theme="twoTone" two-tone-color="#eb2f96" class="icon-align" />
                </a-tooltip>
              </span>
              <a-dropdown :trigger="['contextmenu']">
                <span>{{ title.split('#')[0] }}</span>
                <template #overlay>
                  <a-menu @click="({ key: menuKey }) => onContextMenuClick(treeKey, menuKey)">
                    <a-menu-item key="viewTableStruc">查看表结构</a-menu-item>
                    <a-menu-item key="viewTableBase">查看表信息</a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </template>
          </a-tree>
        </div>
      </div>
    </a-spin>

    <template>
      <EditTableInfoComponent
        ref="EditTableInfoComponent"
        :visible="tableInfoVisible"
        @close="tableInfoVisible = false"
      ></EditTableInfoComponent>
    </template>
  </div>
</template>

<script>
import { getSchemasApi, getTablesApi, getUserGrantsApi } from '@/api/das'
import EditTableInfoComponent from './EditTableInfo'

export default {
  props: {
    screenHeight: String,
  },
  components: {
    EditTableInfoComponent,
  },
  data() {
    return {
      form: this.$form.createForm(this),
      service_name: [],
      bindTitle: '',
      schemas: [],
      showSearch: false,
      treeLoading: false,
      treeData: [],
      searchTreeData: [],
      selectedSchema: {},
      refreshLoading: false,
      tableInfoVisible: false,
      selectedKeys: {}, // 设置选中key的节点
    }
  },
  methods: {
    // 刷新schemas列表
    refreshSchemas() {
      this.refreshLoading = true
      this.getSchemas()
      this.refreshLoading = false
      this.$message.info('库列表刷新成功，请展开下拉列表查看')
    },
    // 搜索表
    onSearch(value) {
      if (this.searchTreeData.length === 0) {
        this.searchTreeData = Object.assign([], this.treeData)
      }
      if (!value) {
        this.treeData = this.searchTreeData
      }
      let searchResult = []
      this.treeData
        .map((item) => {
          if (item.title.indexOf(value) > -1) {
            searchResult.push(item)
          }
          return null
        })
        .filter((item, i, self) => item && self.indexOf(item) === i)
      this.treeData = searchResult
    },
    // 获取指定schema的表权限
    async getGrants(params) {
      let data = []
      await getUserGrantsApi(params).then((response) => {
        data = response.data
      })
      return data
    },
    // 获取用户有权限的schemas
    getSchemas() {
      getSchemasApi()
        .then((response) => {
          if (response.code === '0001') {
            this.$notification.error({
              title: '加载失败',
              message: response.message,
            })
            return false
          }
          this.schemas = response.data
        })
        .catch((error) => {
          console.log('error', error)
        })
    },
    // 获取授权的表
    getTables(value) {
      // 将searchTreeData设置为[]
      this.searchTreeData = []
      // display
      this.showSearch = true
      // loading
      this.treeLoading = true
      let vals = value.split('#')
      this.selectedSchema = {
        instance_id: vals[0],
        schema: vals[1],
        db_type: vals[2],
      }
      // 设置选中的schema
      this.$emit('changeSelectedSchema', this.selectedSchema)
      // 请求参数
      let params = {
        instance_id: vals[0],
        schema: vals[1],
      }
      getTablesApi(params)
        .then((response) => {
          if (response.code === '0000') {
            this.getGrants(this.selectedSchema).then((val) => {
              this.renderTree(val, response.data)
            })
          } else {
            this.$notification.error({
              title: '加载失败',
              message: response.message,
            })
            return false
          }
        })
        .finally(() => {
          this.treeLoading = false
        })
    },
    // 渲染tree
    renderTree(grants, data) {
      // 初始化treeData
      let tmpTreeData = []
      // 初始化自动补全
      let tabCompletion = { tables: {} }
      // 检查表是否有权限，并打上标识
      let checkTableRule = function (grants, table) {
        // 正常检查逻辑
        if (grants.tables.length === 1 && grants.tables === '*') {
          return true
        }
        let hasAllow = false
        if (grants.tables[0]['rule'] === 'allow') {
          hasAllow = true
        }
        if (hasAllow === true) {
          for (const [i, v] of grants.tables.entries()) {
            if (v['rule'] === 'allow' && v['table'] === table) {
              return true
            }
          }
          return false
        }
        if (hasAllow === false) {
          for (const [i, v] of grants.tables.entries()) {
            if (v['rule'] === 'deny' && v['table'] === table) {
              return false
            }
          }
          return true
        }
      }
      // 循环数据
      data.forEach(function (row) {
        let tmpColumnsData = []
        let columnsCompletion = []
        // 列节点
        for (let [i, v] of row.columns.split('@@').entries()) {
          let colName = v.split('$$')[0]
          tmpColumnsData.push({
            title: v.replaceAll('$$', ' '),
            key: `${row['table_schema']}#${row['table_name']}#${colName}`,
            scopedSlots: { switcherIcon: 'child' },
            isLeaf: true,
          })
          columnsCompletion.push(colName)
        }
        // 表节点
        var rule = checkTableRule(grants, row.table_name) === true ? 'allow' : 'deny'
        tmpTreeData.push({
          title: `${row.table_name}#${rule}`,
          key: `${row['table_schema']}#${row['table_name']}`,
          scopedSlots: { title: 'custom' },
          children: tmpColumnsData,
        })
        tabCompletion['tables'][row['table_name']] = columnsCompletion
      })
      // 渲染树结构
      this.treeData = tmpTreeData
      // 加载表列自动提示
      this.$emit('setTabCompletion', tabCompletion)
    },
    onContextMenuClick(treeKey, menuKey) {
      if (treeKey.split('#').length >= 2) {
        let vals = treeKey.split('#')
        this.selectedKeys = {
          schema: vals[0],
          table: vals[1],
        }
        if (menuKey === 'viewTableStruc') {
          this.viewTableStruc()
        }
        if (menuKey === 'viewTableBase') {
          this.viewTableBase()
        }
      }
    },
    // 查看表结构
    viewTableStruc() {
      this.tableInfoVisible = true
      const params = {
        ...this.selectedKeys,
        ...this.selectedSchema,
        type: 'structure',
      }
      this.$refs.EditTableInfoComponent.fetchData(params)
    },
    // 查看表元信息
    viewTableBase() {
      this.tableInfoVisible = true
      const params = {
        ...this.selectedKeys,
        ...this.selectedSchema,
        type: 'base',
      }
      this.$refs.EditTableInfoComponent.fetchData(params)
    },
  },
  mounted() {
    this.getSchemas()
  },
}
</script>

<style lang="less" scoped>
.tree {
  overflow: scroll;
  min-height: 100px;
  border-radius: 4px;
  border-left-width: 0px;
  border-right-width: 0px;
}

// 更改树+号的位置，默认没有对齐
/deep/ .ant-tree.ant-tree-show-line li span.ant-tree-switcher {
  line-height: 20px;
}

/deep/.ant-tree li {
  padding: 2px 0;
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

// icon和文字对齐
/deep/.icon-align {
  margin-right: 6px;
  vertical-align: -15%;
}
</style>
