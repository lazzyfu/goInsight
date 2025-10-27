<template>
  <div class="header">
    <a-select v-model:value="data.schema" style="width: 90%" @change="loadTablesBySchema">
      <a-select-option v-for="(s, index) in data.dbList" :key="index" :label="`${s.remark}:${s.schema}`"
        :value="`${s.instance_id};${s.schema};${s.db_type}`" :disabled="s.is_deleted">
        <DbIcon :type="s.db_type.toLowerCase()" /> {{ s.remark }}:{{ s.schema }}
        <i v-if="s.is_deleted" style="color: #c0c4cc">已删除</i>
      </a-select-option>
    </a-select>
  </div>
  <a-input-search style="margin: 5px 0px;width: 90%" placeholder="输入要搜索的表名" @search="handleSearch"
    :disabled="!isSearchTableEnabled" />
  <a-spin :spinning="isTreeLoading" tip="加载中...">
    <div id="tree-container">
      <div class="block">
        <a-tree :tree-data="data.searchTreeData" show-line class="tree filter-tree" :defaultExpandAll="true">
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
    </div>
  </a-spin>
  <a-drawer v-model:open="open" width="50%" title="表元信息" placement="right">
    <highlightjs language="sql" :code="data.tableInfo" />
  </a-drawer>
</template>

<script setup>
import { GetPermittedTablesBySchemaApi, GetSchemaTablesApi, GetSchemasApi, GetTableInfoApi } from "@/api/das"
import { TabletTwoTone } from "@ant-design/icons-vue"
import { message } from "ant-design-vue"
import { inject, onMounted, reactive, ref } from 'vue'
import PermissionHint from './components/PermissionHint.vue'
import DbIcon from './icon.vue'

const dasInstanceData = inject('dasInstanceData')

const open = ref(false)
const isTreeLoading = ref(false)
const isSearchTableEnabled = ref(false)

const data = reactive({
  tableInfo: '',
  dbList: [],
  selectedTable: undefined,
  searchTreeData: [],
  streeData: [],
})

const instanceData = reactive({ schema: '', instance_id: '', tables: {} })

const fetchSchemas = async () => {
  const resp = await GetSchemasApi()
  if (resp.code === '0000') {
    data.dbList = resp.data
  } else {
    message.error(resp.message)
  }
}

const handleSearch = (value) => {
  if (value) {
    let searchResult = []
    data.streeData.map((item) => {
      if (item.title.indexOf(value) > -1) {
        searchResult.push(item)
      }
      return null
    }).filter((item, i, self) => item && self.indexOf(item) === i)
    data.searchTreeData = searchResult
  } else {
    data.searchTreeData = data.streeData
  }
}

// 获取指定实例和schema的表
const loadTablesBySchema = async (value) => {
  const vals = value.split(';')
  instanceData.instance_id = vals[0]
  instanceData.schema = vals[1]
  instanceData.db_type = vals[2]
  const payLoad = {
    "instance_id": instanceData.instance_id,
    "schema": instanceData.schema,
  }
  isSearchTableEnabled.value = true
  isTreeLoading.value = true
  const resp = await GetSchemaTablesApi(payLoad)
  if (resp.code === '0000') {
    // 获取指定schema的表权限
    GetPermittedTablesBySchemaApi(payLoad).then((val) => {
      renderTree(val.data, resp.data)
      isTreeLoading.value = false
    })
  } else {
    renderTree([], [])
    isTreeLoading.value = false
    message.error(resp.message)
  }
}

const generatorColumnNodes = (columns, tableSchema, tableName) => {
  return columns.split('@@').map(v => {
    const colName = v.split('$$')[0];
    return {
      title: v.replaceAll('$$', ' '),
      key: `${tableSchema}#${tableName}#${colName}`,
      scopedSlots: { switcherIcon: 'child' },
      isLeaf: true,
    };
  });
};

const generatorTableNode = (grants, table, columnNodes) => {
  const rule = checkTableRule(grants, table.table_name) ? 'allow' : 'deny';
  return {
    title: `${table.table_name}#${rule}`,
    key: `${table.table_schema}#${table.table_name}`,
    scopedSlots: { title: 'custom' },
    children: columnNodes,
  };
};

const renderTree = (grants, tableList) => {
  const tmpTreeData = [];
  const tables = {};

  tableList.forEach(row => {
    const columnNodes = generatorColumnNodes(row.columns, row.table_schema, row.table_name);
    const tableNode = generatorTableNode(grants, row, columnNodes);

    tmpTreeData.push(tableNode);
    tables[row.table_name] = columnNodes.map(node => node.key.split('#')[2]);
  });

  instanceData.tables = { ...tables };
  data.streeData = tmpTreeData;
  data.searchTreeData = tmpTreeData;
  dasInstanceData.value = { ...instanceData };
};


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

const onContextMenuClick = (treeKey, menuKey) => {
  if (treeKey.split('#').length >= 2) {
    let vals = treeKey.split('#')
    data.selectedTable = {
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
const getTableMeta = (type) => {
  open.value = true
  const params = {
    ...data.selectedTable,
    type: type,
    instance_id: instanceData.instance_id,
  }

  GetTableInfoApi(params).then((res) => {
    if (res.code === '0000') {
      if (type === 'structure') {
        res.data.forEach((element) => {
          for (const i in element) {
            if (
              i.toLowerCase() === 'create table' ||
              i.toLowerCase() === 'statement'
            ) {
              data.tableInfo = element[i]
            }
          }
        });
      }
      if (type === 'base') {
        res.data.forEach((element) => {
          var tableBase = [];
          res.data.forEach((element) => {
            for (var key in element) {
              tableBase.push(`${key}  ${element[key]}`);
            }
          });
          data.tableInfo = tableBase.join('\n')
        });
      }

    } else {
      message.error(res.message)
    }
  })
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

:deep(#tree-container) {
  overflow: scroll;
  height: 580px;
  border-radius: 4px;
  border-left-width: 0px;
  border-right-width: 0px;
}

:deep(.ant-tree.ant-tree-show-line li span.ant-tree-switcher) {
  line-height: 20px;
}

:deep(.ant-tree li) {
  padding: 2px 0;
}</style>
