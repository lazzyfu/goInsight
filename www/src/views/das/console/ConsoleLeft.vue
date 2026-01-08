<template>
  <div class="console-left-root">
    <div class="header">
      <a-select style="width: 90%" @change="loadTablesBySchema" placeholder="选择有权限的库">
        <a-select-option
          v-for="(s, index) in uiData.dbList"
          :key="index"
          :label="`${s.remark}:${s.schema}`"
          :value="`${s.instance_id};${s.schema};${s.db_type}`"
          :disabled="s.is_deleted"
        >
          <ConsoleDbIcon :type="s.db_type.toLowerCase()" /> {{ s.remark }}:{{ s.schema }}
          <i v-if="s.is_deleted" style="color: #c0c4cc">已删除</i>
        </a-select-option>
      </a-select>
    </div>

    <a-input-search
      style="margin: 5px 0px; width: 90%"
      placeholder="输入要搜索的表名"
      @search="handleSearch"
      :disabled="!uiState.isSearchTable"
    />

    <a-spin class="tree-area" :spinning="uiState.isTreeLoading" tip="加载中...">
      <div id="tree-container">
        <div class="block">
          <a-tree
            :tree-data="uiData.searchTreeData"
            show-line
            class="tree filter-tree"
            :defaultExpandAll="true"
          >
            <template #icon="record">
              <span v-if="record.isLeaf">
                <TabletTwoTone />
              </span>
            </template>
            <template #title="{ key: treeKey, title, isLeaf }">
              <PermissionHint
                v-if="title.split('#').length === 2"
                :hasAccess="title.split('#')[1] === 'allow'"
              />
              <a-dropdown :trigger="['contextmenu']">
                <span>{{ title.split('#')[0] }}</span>
                <template #overlay>
                  <a-menu
                    v-if="!isLeaf"
                    @click="({ key: menuKey }) => onContextMenuClick(treeKey, menuKey)"
                  >
                    <a-menu-item key="showTableStructure">查看表结构</a-menu-item>
                    <a-menu-item key="showTableMetadata">查看表信息</a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </template>
          </a-tree>
        </div>
      </div>
    </a-spin>

    <a-drawer v-model:open="uiState.open" width="50%" title="表元信息" placement="right">
      <highlightjs language="sql" :code="uiData.tableInfo" />
    </a-drawer>
  </div>
</template>

<script setup>
import {
    GetPermittedTablesBySchemaApi,
    GetSchemaTablesApi,
    GetSchemasApi,
    GetTableInfoApi,
} from '@/api/das'
import { TabletTwoTone } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { inject, onMounted, reactive } from 'vue'
import PermissionHint from './components/PermissionHint.vue'
import ConsoleDbIcon from './ConsoleDbIcon.vue'

// 共享实例数据
const dasInstanceData = inject('dasInstanceData')

// 状态
const uiState = reactive({
  open: false,
  isTreeLoading: false,
  isSearchTable: false,
})

// 数据
const uiData = reactive({
  tableInfo: '',
  dbList: [],
  selectedTable: undefined,
  searchTreeData: [],
  streeData: [],
  instanceData: { schema: '', instance_id: '', db_type: '', tables: {} },
})

// 获取schema列表
const fetchSchemas = async () => {
  const res = await GetSchemasApi().catch(() => {})
  if (res) {
    uiData.dbList = res.data
  }
}

// 搜索表
const handleSearch = (value) => {
  if (value) {
    let searchResult = []
    uiData.streeData
      .map((item) => {
        if (item.title.indexOf(value) > -1) {
          searchResult.push(item)
        }
        return null
      })
      .filter((item, i, self) => item && self.indexOf(item) === i)
    uiData.searchTreeData = searchResult
  } else {
    uiData.searchTreeData = uiData.streeData
  }
}

// 获取指定实例和schema的表
const loadTablesBySchema = async (value) => {
  const vals = value.split(';')
  console.log('vals: ', vals);
  uiData.instanceData = {
    instance_id: vals[0],
    schema: vals[1],
    db_type: vals[2],
    tables: {},
  }

  uiState.isSearchTable = true
  uiState.isTreeLoading = true

  const res = await GetSchemaTablesApi(uiData.instanceData).catch(() => {})
  if (res) {
    // 获取指定schema的表权限
    const tableRes = await GetPermittedTablesBySchemaApi(uiData.instanceData).catch(() => {})
    if (tableRes) {
      renderTree(tableRes.data, res.data)
      uiState.isTreeLoading = false
    }
  } else {
    renderTree([], [])
    uiState.isTreeLoading = false
    message.error(res?.message)
  }
}

// 生成列节点
const generatorColumnNodes = (columns, tableSchema, tableName) => {
  return columns.split('@@').map((v) => {
    const colName = v.split('$$')[0]
    return {
      title: v.replaceAll('$$', ' '),
      key: `${tableSchema}#${tableName}#${colName}`,
      scopedSlots: { switcherIcon: 'child' },
      isLeaf: true,
    }
  })
}

// 生成表节点
const generatorTableNode = (grants, table, columnNodes) => {
  const rule = checkTableRule(grants, table.table_name) ? 'allow' : 'deny'
  return {
    title: `${table.table_name}#${rule}`,
    key: `${table.table_schema}#${table.table_name}`,
    scopedSlots: { title: 'custom' },
    children: columnNodes,
  }
}

// 渲染树结构
const renderTree = (grants, tableList) => {
  const tmpTreeData = []
  const tables = {}

  tableList.forEach((row) => {
    const columnNodes = generatorColumnNodes(row.columns, row.table_schema, row.table_name)
    const tableNode = generatorTableNode(grants, row, columnNodes)

    tmpTreeData.push(tableNode)
    tables[row.table_name] = columnNodes.map((node) => node.key.split('#')[2])
  })

  uiData.instanceData.tables = { ...tables }
  uiData.streeData = tmpTreeData
  uiData.searchTreeData = tmpTreeData
  dasInstanceData.value = { ...uiData.instanceData }
}

// 检查表是否有权限，并打上标识
const checkTableRule = (grants, table) => {
  // 正常检查逻辑
  if (grants.tables.length === 1 && grants.tables === '*') {
    return true
  }
  var hasAllow = false
  if (grants.tables[0]['rule'] === 'allow') {
    hasAllow = true
  }
  if (hasAllow === true) {
    for (const index in grants.tables) {
      const v = grants.tables[index]
      if (v['rule'] === 'allow' && v['table'] === table) {
        return true
      }
    }
    return false
  }
  if (hasAllow === false) {
    for (const index in grants.tables) {
      const v = grants.tables[index]
      if (v['rule'] === 'deny' && v['table'] === table) {
        return false
      }
    }
    return true
  }
}

// 右键菜单点击事件
const onContextMenuClick = (treeKey, menuKey) => {
  if (treeKey.split('#').length >= 2) {
    let vals = treeKey.split('#')
    uiData.selectedTable = {
      schema: vals[0],
      table: vals[1],
    }
    if (menuKey === 'showTableStructure') {
      showTableStructure()
    } else if (menuKey === 'showTableMetadata') {
      showTableMetadata()
    }
  }
}

// 获取表元信息
const getTableMeta = async (type) => {
  uiState.open = true
  const params = {
    ...uiData.selectedTable,
    type: type,
    instance_id: uiData.instanceData.instance_id,
  }

  const res = await GetTableInfoApi(params).catch(() => {})
  if (res) {
    if (type === 'structure') {
      res.data.forEach((row) => {
        for (const i in row) {
          if (i.toLowerCase() === 'create table' || i.toLowerCase() === 'statement') {
            uiData.tableInfo = row[i]
          }
        }
      })
    }
    if (type === 'base') {
      res.data.forEach(() => {
        var tableBase = []
        res.data.forEach((row) => {
          for (var key in row) {
            tableBase.push(`${key}  ${row[key]}`)
          }
        })
        uiData.tableInfo = tableBase.join('\n')
      })
    }
  }
}

// 查看表结构
const showTableStructure = () => {
  getTableMeta('structure')
}

// 查看表元信息
const showTableMetadata = () => {
  getTableMeta('base')
}

onMounted(() => {
  fetchSchemas()
})
</script>

<style scoped>
:deep(.header) {
  --border: 1;
}

.console-left-root {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.tree-area {
  flex: 1;
  min-height: 0;
}

:deep(#tree-container) {
  overflow: auto;
  overflow-x: hidden;
  height: 100%;
  border-radius: 4px;
  border-left-width: 0px;
  border-right-width: 0px;
}

:deep(.ant-tree.ant-tree-show-line li span.ant-tree-switcher) {
  line-height: 20px;
}

:deep(.ant-tree li) {
  padding: 2px 0;
}
</style>
