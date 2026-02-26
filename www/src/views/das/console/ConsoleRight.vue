<template>
  <a-tabs
    class="console-tabs"
    v-model="uiData.activeKey"
    type="editable-card"
    size="small"
    @edit="handleTabEdit"
    @change="handleTabChange"
  >
    <a-tab-pane
      v-for="pane in uiData.panes"
      :key="pane.key"
      :tab="pane.title"
      :closable="pane.closable"
    >
    </a-tab-pane>
  </a-tabs>
  <a-space class="console-toolbar" size="small">
    <a-button type="primary" @click="executeSqlQuery()">
      <template #icon>
        <PlayCircleOutlined />
      </template>
      执行SQL
    </a-button>
    <a-button @click="formatSqlContent()">
      <template #icon>
        <CodeOutlined />
      </template>
      格式化
    </a-button>
    <a-button @click="addToFavorites()">
      <template #icon>
        <StarOutlined />
      </template>
      <a-tooltip>
        <template #title>先鼠标选中SQL，然后点击“收藏SQL”按钮</template>
        收藏SQL
      </a-tooltip>
    </a-button>
    <a-button @click="generatorDataDictionary()">
      <template #icon>
        <BookOutlined />
      </template>
      数据字典
    </a-button>
    <span>
      字符集
      <a-select style="width: 120px" v-model:value="uiData.characterSet" @change="saveTabToCache">
        <a-select-option v-for="item in characterSets" :key="item.key" :value="item.value">
          {{ item.key }}
        </a-select-option>
      </a-select>
    </span>
  </a-space>
  <div class="console-editor-wrap">
    <a-spin :spinning="currentTabLoading" tip="Loading...">
      <div class="console-editor-surface">
        <CodeMirror ref="codemirrorRef" :height="'300px'" />
      </div>
    </a-spin>
  </div>
  <!-- 数据字典 -->
  <ConsoleDbDict
    ref="dbDictRef"
    :open="uiState.isDbDictOpen"
    @update:open="uiState.isDbDictOpen = false"
  />
  <!-- 新增收藏SQL -->
  <DasFavoriteFormModal
    :open="uiState.isFavoritesOpen"
    v-model:modelValue="favoritesFormState"
    @update:open="uiState.isFavoritesOpen = $event"
    @submit="handleFavoritesSubmit"
  />
</template>

<script setup>
import {
  CreateFavoritesApi,
  ExecuteClickHouseQueryApi,
  ExecuteMySQLQueryApi,
  GetDBDictApi,
} from '@/api/das'
import CodeMirror from '@/components/edit/Codemirror.vue'
import DasFavoriteFormModal from '@/views/das/favorite/DasFavoriteFormModal.vue'
import { BookOutlined, CodeOutlined, PlayCircleOutlined, StarOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { computed, inject, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import ConsoleDbDict from './ConsoleDbDict.vue'

// 每个浏览器窗口(tab)独立：使用 sessionStorage 持久化 Console 状态
const storage = sessionStorage

// 每个 Console Tab(内部窗口)独立：运行时状态按 activeKey 隔离
const runtime = reactive({
  queryResultMessageByTab: {},
  resultByTab: {},
  tableLoadingByTab: {},
})

// 字符集选项
const characterSets = [
  { key: 'utf8', value: 'utf8' },
  { key: 'utf8mb4', value: 'utf8mb4' },
  { key: 'latin1', value: 'latin1' },
]

// 状态
const uiState = reactive({
  isDbDictOpen: false,
  isFavoritesOpen: false,
})

// 引用
const dbDictRef = ref(null)
const codemirrorRef = ref(null)

const favoritesFormState = ref({
  sqltext: '',
})

// 获取共享数据
const dasInstanceData = inject('dasInstanceData')

const emit = defineEmits(['renderResultTable', 'renderExecutionMessage'])
const uiData = reactive({
  panes: [],
  activeKey: 1,
  characterSet: 'utf8',
  schema: '',
  instance_id: '',
  db_type: '',
  queryResultMessage: undefined,
})

const currentTabLoading = computed(() => {
  return !!runtime.tableLoadingByTab[uiData.activeKey]
})

const emitActiveExecutionMessage = () => {
  emit('renderExecutionMessage', runtime.queryResultMessageByTab[uiData.activeKey])
}

const applyActiveTabRuntime = () => {
  const key = uiData.activeKey
  uiData.queryResultMessage = runtime.queryResultMessageByTab[key]
  const res = runtime.resultByTab[key]
  if (res) {
    emit('renderResultTable', res)
  } else {
    emit('renderResultTable', null)
  }
  emitActiveExecutionMessage()
}

// tab编辑
const handleTabEdit = (targetKey, action) => {
  if (action === 'add') {
    addTab()
  } else {
    removeTab(targetKey)
  }
}

// 新增tab
const addTab = () => {
  const tabKey = uiData.panes.length + 1
  uiData.panes.push({
    title: `New Console ${tabKey}`,
    key: tabKey,
  })
}

// 移除tab
const removeTab = (targetKey) => {
  let activeKey = uiData.activeKey
  let lastIndex
  uiData.panes.forEach((pane, i) => {
    if (pane.key === targetKey) {
      lastIndex = i - 1
    }
  })
  const newPanes = uiData.panes.filter((pane) => pane.key !== targetKey)
  if (newPanes.length && activeKey === targetKey) {
    if (lastIndex >= 0) {
      activeKey = newPanes[lastIndex].key
    } else {
      activeKey = newPanes[0].key
    }
  }
  uiData.panes = newPanes
  uiData.activeKey = activeKey
  delete runtime.queryResultMessageByTab[targetKey]
  delete runtime.resultByTab[targetKey]
  delete runtime.tableLoadingByTab[targetKey]
  storage.removeItem('das#tab#' + targetKey)
  storage.setItem('das#panes', JSON.stringify(uiData.panes))

  // 关闭 Tab 后同步恢复当前 Tab 的独立状态
  loadTabFromCache()
  applyActiveTabRuntime()
}

// 切换tab
const handleTabChange = (value) => {
  saveTabToCache()
  uiData.activeKey = value
  loadTabFromCache()
  applyActiveTabRuntime()
}

// 加载tab
const loadTab = () => {
  const panes = JSON.parse(storage.getItem('das#panes'))
  if (panes?.length > 0) {
    uiData.panes = panes
  } else {
    uiData.panes = [{ key: 1, title: 'Console 1', closable: false }]
  }
}

// 保存tab cache
const saveTabToCache = () => {
  var tabData = {
    characterSet: uiData.characterSet,
    userInput: codemirrorRef.value.getContent(),
  }
  storage.setItem('das#tab#' + uiData.activeKey, JSON.stringify(tabData))
}

// 加载tab cache
const loadTabFromCache = () => {
  var tabData = JSON.parse(storage.getItem('das#tab#' + uiData.activeKey))
  if (tabData != null) {
    codemirrorRef.value.setContent(tabData.userInput)
    uiData.characterSet = tabData.characterSet
  } else {
    uiData.characterSet = 'utf8'
    codemirrorRef.value.setContent('')
  }
}

// 生成数据字典
const generatorDataDictionary = () => {
  if (!uiData.schema) {
    message.warning('请先选择左侧的DB库')
    return
  }
  const payLoad = {
    instance_id: uiData.instance_id,
    schema: uiData.schema,
    db_type: uiData.db_type,
  }
  GetDBDictApi(payLoad).then((res) => {
    uiState.isDbDictOpen = true
    dbDictRef.value.render(res.data)
  })
}

// 收藏SQL
const addToFavorites = () => {
  // 获取选择的SQL
  const sqltext = codemirrorRef.value.getSelectedText()
  // 判断是否输入
  if (sqltext == 0) {
    message.warning('请鼠标选中要收藏的SQL')
    return
  }
  uiState.isFavoritesOpen = true
  favoritesFormState.value.sqltext = sqltext
}

// 收藏SQL提交
const handleFavoritesSubmit = (data) => {
  CreateFavoritesApi(data).then((res) => {
    if (res.code == '0000') {
      message.success('收藏成功')
    } else {
      message.error(res.message)
    }
    uiState.isFavoritesOpen = false
  })
}

// 执行SQL
const executeSqlQuery = async () => {
  /* post data:
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
  const tabKey = uiData.activeKey
  saveTabToCache()
  if (!uiData.schema) {
    message.warning('请先选择左侧的DB库')
    return
  }
  // 获取选择的SQL
  const sqltext = codemirrorRef.value.getSelectedText()
  // 判断是否输入
  if (sqltext.length == 0) {
    message.warning('请鼠标选中要执行的SQL')
    return
  }
  const dbType = uiData.db_type?.toLowerCase()
  const payLoad = {
    instance_id: uiData.instance_id,
    schema: uiData.schema,
    db_type: uiData.db_type,
    sqltext: sqltext,
  }
  if (dbType === 'mysql' || dbType === 'tidb') {
    payLoad.params = {
      character_set_client: uiData.characterSet,
      character_set_connection: uiData.characterSet,
      character_set_results: uiData.characterSet,
    }
  }

  let res = null
  runtime.tableLoadingByTab[tabKey] = true
  const resMsgs = []
  try {
    if (dbType === 'mysql' || dbType === 'tidb') {
      res = await ExecuteMySQLQueryApi(payLoad).catch(() => {})
    } else if (dbType === 'clickhouse') {
      res = await ExecuteClickHouseQueryApi(payLoad).catch(() => {})
    } else {
      message.error(`不支持的数据库类型: ${uiData.db_type}`)
      return
    }
    if (res) {
      message.success('执行成功')
      resMsgs.push('结果: 执行成功')
      resMsgs.push(`耗时: ${res.data?.duration ?? '-'}${dbType === 'clickhouse' ? 'ms' : ''}`)
      resMsgs.push(`SQL: ${res.data?.sqltext ?? sqltext}`)
      resMsgs.push(`请求ID: ${res.request_id}`)
      runtime.resultByTab[tabKey] = res.data
      if (uiData.activeKey === tabKey) {
        emit('renderResultTable', res.data)
      }
    }
  } finally {
    runtime.tableLoadingByTab[tabKey] = false
    const msg = resMsgs.join('\n')
    runtime.queryResultMessageByTab[tabKey] = msg
    if (uiData.activeKey === tabKey) {
      uiData.queryResultMessage = msg
      emit('renderExecutionMessage', msg)
    }
  }
}

// 格式化
const formatSqlContent = () => {
  codemirrorRef.value.formatContent()
  saveTabToCache()
}

// 监控用户切库
watch(dasInstanceData, (newVal) => {
  uiData.instance_id = newVal.value.instance_id
  uiData.schema = newVal.value.schema
  uiData.db_type = newVal.value.db_type
  // 设置自动补全
  codemirrorRef.value.setCompletion(newVal.value.tables)
})

onMounted(() => {
  loadTab()
  loadTabFromCache()
  applyActiveTabRuntime()
})

onBeforeUnmount(() => {
  saveTabToCache()
})
</script>

<style scoped>
:deep(.ant-tabs-nav::before) {
  border-bottom-color: var(--ant-colorSplit, #f0f0f0);
}

.console-tabs :deep(.ant-tabs-nav) {
  padding: 4px 8px;
  background: var(--ant-colorFillAlter, #fafafa);
  border: 1px solid var(--ant-colorSplit, #f0f0f0);
  border-radius: var(--ant-borderRadiusLG, 8px);
}

.console-tabs :deep(.ant-tabs-nav-wrap) {
  margin: 0;
}

.console-toolbar {
  padding: 8px;
  background: var(--ant-colorFillAlter, #fafafa);
  border: 1px solid var(--ant-colorSplit, #f0f0f0);
  border-radius: var(--ant-borderRadiusLG, 8px);
}

.console-editor-wrap {
  margin-top: 8px;
}

.console-editor-surface {
  border: 1px solid var(--ant-colorSplit, #f0f0f0);
  border-radius: var(--ant-borderRadiusLG, 8px);
  background: var(--ant-colorBgContainer, #ffffff);
  overflow: hidden;
}

:deep(.cm-gutters) {
  color: var(--ant-colorTextTertiary, rgba(0, 0, 0, 0.45));
  background-color: var(--ant-colorFillAlter, #fafafa);
}

:deep(.box-card) {
  font-size: 12px;
  margin-top: 4px;
  height: 100px;
  width: 100%;
  overflow: auto;
  zoom: 1;
  white-space: normal;
  word-break: break-all;
}
</style>
