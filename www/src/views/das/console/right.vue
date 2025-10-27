<template>
  <a-tabs v-model="data.activeKey" type="editable-card" size="small" @edit="handleTabEdit" @change="handleTabChange">
    <a-tab-pane v-for="pane in data.panes" :key="pane.key" :tab="pane.title" :closable="pane.closable">
    </a-tab-pane>
  </a-tabs>
  <a-space size="small">
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
      <a-select style="width: 120px" v-model:value="data.characterSet" @change="saveTabToCache">
        <a-select-option v-for="item in characterSets" :key="item.key" :value="item.value">
          {{ item.key }}
        </a-select-option>
      </a-select>
    </span>

  </a-space>
  <div style="margin-top: 6px">
    <a-spin :spinning="data.tableLoading" tip="Loading...">
      <CodeMirror ref="codemirrorRef" />
      <a-card class="box-card">
        <span v-html="data.queryResultMessage"></span>
      </a-card>
    </a-spin>
  </div>
  <!-- 数据字典 -->
  <DbDict ref="dbDictRef" :open="isDbDictOpen" @update:open="isDbDictOpen = false" />
  <!-- 新增收藏SQL -->
  <FavoritesAdd :open="isFavoritesOpen" :formState="favoritesFormState" :btnType="favoritesBtnType"
    @update:open="isFavoritesOpen = $event" @submit="handleFavoritesSubmit" />
</template>

<script setup>
import { CreateFavoritesApi, ExecuteClickHouseQueryApi, ExecuteMySQLQueryApi, GetDBDictApi } from "@/api/das";
import CodeMirror from '@/components/edit/Codemirror.vue';
import FavoritesAdd from "@/views/das/favorite/modal.vue";
import { BookOutlined, CodeOutlined, PlayCircleOutlined, StarOutlined } from "@ant-design/icons-vue";
import { message } from "ant-design-vue";
import { inject, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue';
import DbDict from './dbdict.vue';

const characterSets = [
  { key: 'utf8', value: 'utf8' },
  { key: 'utf8mb4', value: 'utf8mb4' },
  { key: 'latin1', value: 'latin1' },
]
const isDbDictOpen = ref(false)
const dbDictRef = ref(null)

const isFavoritesOpen = ref(false)
const favoritesBtnType = ref('收藏')
const favoritesFormState = ref({
  sqltext: ''
});

const dasInstanceData = inject('dasInstanceData');
const codemirrorRef = ref(null);
const emit = defineEmits(["renderResultTable"])
const data = reactive({
  panes: [],
  activeKey: 1,
  characterSet: 'utf8',
  tableLoading: false,
  schema: '',
  instance_id: '',
  db_type: '',
  queryResultMessage: undefined,
})

const handleTabEdit = (targetKey, action) => {
  if (action === 'add') {
    addTab()
  } else {
    removeTab(targetKey)
  }
}

const addTab = () => {
  const tabKey = data.panes.length + 1
  data.panes.push({
    title: `New Console ${tabKey}`,
    key: tabKey,
  })
}

const removeTab = (targetKey) => {
  let activeKey = data.activeKey
  let lastIndex
  data.panes.forEach((pane, i) => {
    if (pane.key === targetKey) {
      lastIndex = i - 1
    }
  })
  const newPanes = data.panes.filter((pane) => pane.key !== targetKey)
  if (newPanes.length && activeKey === targetKey) {
    if (lastIndex >= 0) {
      activeKey = newPanes[lastIndex].key
    } else {
      activeKey = newPanes[0].key
    }
  }
  data.panes = newPanes
  data.activeKey = activeKey
  localStorage.removeItem('das#tab#' + targetKey)
  localStorage.setItem('das#panes', JSON.stringify(data.panes))
}

const handleTabChange = (value) => {
  saveTabToCache()
  data.activeKey = value
  loadTabFromCache()
}

const loadTab = () => {
  const panes = JSON.parse(localStorage.getItem('das#panes'))
  if (panes?.length > 0) {
    data.panes = panes
  } else {
    data.panes = [{ key: 1, title: 'Console 1', closable: false }]
  }
}

// 保存tab cache
const saveTabToCache = () => {
  var tabData = {
    "characterSet": data.characterSet,
    "userInput": codemirrorRef.value.getContent(),
  }
  localStorage.setItem('das#tab#' + data.activeKey, JSON.stringify(tabData))
}

// 加载tab cache
const loadTabFromCache = () => {
  var tabData = JSON.parse(localStorage.getItem('das#tab#' + data.activeKey))
  if (tabData != null) {
    codemirrorRef.value.setContent(tabData.userInput)
    data.characterSet = tabData.characterSet
  } else {
    data.characterSet = "utf8"
    codemirrorRef.value.setContent("")
  }
}

// 生成数据字典
const generatorDataDictionary = () => {
  if (!data.schema) {
    message.warning('请先选择左侧的DB库')
    return
  }
  const payLoad = {
    instance_id: data.instance_id,
    schema: data.schema,
    db_type: data.db_type
  }
  GetDBDictApi(payLoad).then((res) => {
    isDbDictOpen.value = true
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
  isFavoritesOpen.value = true
  favoritesFormState.value.sqltext = sqltext
}

// 收藏SQL提交
const handleFavoritesSubmit = (data) => {
  CreateFavoritesApi(data).then(res => {
    if (res.code == '0000') {
      message.success('收藏成功')
    } else {
      message.error(res.message)
    }
    isFavoritesOpen.value = false
  })
}

// 处理执行SQL返回的结果
const handleQueryResult = (res, resMsgs, emit) => {
  data.tableLoading = false;
  if (res.code === '0000') {
    message.success('执行成功');
    resMsgs.push('结果: 执行成功');
    resMsgs.push(`耗时: ${res.data.duration}${data.db_type === 'clickhouse' ? 'ms' : ''}`);
    resMsgs.push(`SQL: ${res.data.sqltext}`);
    emit('renderResultTable', res.data);
  } else {
    message.error('执行失败, ' + res.message);
    resMsgs.push('结果: 执行失败');
    resMsgs.push(`错误: ${res.message}`);
    resMsgs.push(`请求ID: ${res.request_id}`);
    emit('renderResultTable', null);
  }
  data.queryResultMessage = resMsgs.join('<br>');
};

// 执行SQL
const executeSqlQuery = () => {
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
  saveTabToCache()
  if (!data.schema) {
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

  const payLoad = {
    instance_id: data.instance_id,
    schema: data.schema,
    db_type: data.db_type,
    sqltext: sqltext
  }

  var resMsgs = []
  data.tableLoading = true

  const queryHandlers = {
    mysql: ExecuteMySQLQueryApi,
    tidb: ExecuteMySQLQueryApi,
    clickhouse: ExecuteClickHouseQueryApi,
  };

  const dbType = data.db_type.toLowerCase();
  if (queryHandlers[dbType]) {
    payLoad.params = {
      character_set_client: data.characterSet,
      character_set_connection: data.characterSet,
      character_set_results: data.characterSet,
    };
    queryHandlers[dbType](payLoad)
      .then(res => handleQueryResult(res, resMsgs, emit));
  }
}

// 格式化
const formatSqlContent = () => {
  codemirrorRef.value.formatContent();
  saveTabToCache()
}

// 监控用户切库
watch(dasInstanceData, (newVal) => {
  data.instance_id = newVal.value.instance_id
  data.schema = newVal.value.schema
  data.db_type = newVal.value.db_type
  // 设置自动补全
  codemirrorRef.value.setCompletion(newVal.value.tables)
})

onMounted(() => {
  loadTab()
  loadTabFromCache()
})

onBeforeUnmount(() => {
  saveTabToCache()
})
</script>

<style scoped>
:deep(.cm-gutters) {
  color: #9f9d9d;
  background-color: #f7f7f7;
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
