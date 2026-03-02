<template>
  <div class="console-left-root" :style="{ height: leftHeight }">
    <div class="fixed-header">
      <a-form layout="vertical" class="left-filter-form">
        <a-form-item label="数据源">
          <a-select class="header-select" @change="loadTablesBySchema" placeholder="选择有权限的数据库">
            <a-select-option
              v-for="(s, index) in uiData.dbList"
              :key="index"
              :label="`${s.remark}:${s.schema}`"
              :value="`${s.instance_id};${s.schema};${s.db_type}`"
              :disabled="s.is_deleted"
            >
              <ConsoleDbIcon :type="s.db_type.toLowerCase()" /> {{ s.remark }}:{{ s.schema }}
              <i v-if="s.is_deleted" class="db-deleted">已删除</i>
            </a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item label="表检索" class="search-item">
          <a-input-search
            class="header-search"
            placeholder="输入要搜索的表名"
            @search="handleSearch"
            :disabled="!uiState.isSearchTable"
          />
        </a-form-item>
      </a-form>
    </div>

    <div class="tree-area">
      <a-empty
        v-if="!uiState.isSearchTable && !uiState.isTreeLoading"
        description="请选择数据库后查看数据表"
      />
      <a-spin v-else :spinning="uiState.isTreeLoading" tip="加载中...">
        <div id="tree-container">
          <a-tree :tree-data="uiData.searchTreeData" show-line class="tree filter-tree" :defaultExpandAll="true">
            <template #icon="record">
              <span v-if="record.isLeaf">
                <TabletTwoTone />
              </span>
            </template>
            <template #title="{ key: treeKey, title, isLeaf }">
              <PermissionHint v-if="title.split('#').length === 2" :hasAccess="title.split('#')[1] === 'allow'" />
              <a-dropdown :trigger="['contextmenu']">
                <span>{{ title.split('#')[0] }}</span>
                <template #overlay>
                  <a-menu v-if="!isLeaf" @click="({ key: menuKey }) => onContextMenuClick(treeKey, menuKey)">
                    <a-menu-item key="showTableStructure">查看表结构</a-menu-item>
                    <a-menu-item key="showTableMetadata">查看表信息</a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </template>
          </a-tree>
        </div>
      </a-spin>
    </div>

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
import { inject, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import PermissionHint from './components/PermissionHint.vue'
import ConsoleDbIcon from './ConsoleDbIcon.vue'

const dasInstanceData = inject('dasInstanceData')

const uiState = reactive({
  open: false,
  isTreeLoading: false,
  isSearchTable: false,
})

const uiData = reactive({
  tableInfo: '',
  dbList: [],
  selectedTable: undefined,
  searchTreeData: [],
  streeData: [],
  instanceData: { schema: '', instance_id: '', db_type: '', tables: {} },
})

const fetchSchemas = async () => {
  const res = await GetSchemasApi().catch(() => {})
  if (res) {
    uiData.dbList = res.data
  }
}

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

const loadTablesBySchema = async (value) => {
  const vals = value.split(';')
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

const generatorTableNode = (grants, table, columnNodes) => {
  const rule = checkTableRule(grants, table.table_name) ? 'allow' : 'deny'
  return {
    title: `${table.table_name}#${rule}`,
    key: `${table.table_schema}#${table.table_name}`,
    scopedSlots: { title: 'custom' },
    children: columnNodes,
  }
}

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

const checkTableRule = (grants, table) => {
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

const showTableStructure = () => {
  getTableMeta('structure')
}

const showTableMetadata = () => {
  getTableMeta('base')
}

const leftHeight = ref('auto')
const resizeObserver = ref(null)

const syncHeightWithRightPanel = () => {
  const rightPanel = document.querySelector('.right-content')
  if (rightPanel) {
    resizeObserver.value = new ResizeObserver((entries) => {
      for (const entry of entries) {
        const height = entry.contentRect.height
        leftHeight.value = `${height}px`
      }
    })
    resizeObserver.value.observe(rightPanel)
  }
}

onMounted(() => {
  fetchSchemas()
  setTimeout(() => {
    syncHeightWithRightPanel()
  }, 100)
})

onBeforeUnmount(() => {
  if (resizeObserver.value) {
    resizeObserver.value.disconnect()
  }
})
</script>

<style scoped>
.console-left-root {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-sizing: border-box;
  background: var(--ant-colorBgContainer, #ffffff);
  border: 1px solid var(--ant-colorSplit, #f0f0f0);
  border-radius: var(--ant-borderRadiusLG, 8px);
}

.fixed-header {
  flex-shrink: 0;
  padding: 10px 10px 8px;
  border-bottom: 1px solid var(--ant-colorSplit, #f0f0f0);
  background: var(--ant-colorFillAlter, #fafafa);
}

.left-filter-form :deep(.ant-form-item) {
  margin-bottom: 8px;
}

.left-filter-form :deep(.search-item) {
  margin-bottom: 0;
}

.header-select,
.header-search {
  width: 100%;
}

.db-deleted {
  color: var(--ant-colorTextQuaternary, rgba(0, 0, 0, 0.25));
  margin-left: 6px;
}

.tree-area {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 8px;
  box-sizing: border-box;
  background: var(--ant-colorBgContainer, #ffffff);
}

:deep(#tree-container) {
  min-height: 100%;
}

:deep(.ant-tree.ant-tree-show-line li span.ant-tree-switcher) {
  line-height: 20px;
}

:deep(.ant-tree li) {
  padding: 2px 0;
}

:deep(.ant-tree-node-content-wrapper:hover) {
  background: var(--ant-colorFillAlter, #fafafa);
}
</style>
