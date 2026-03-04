<template>
  <div class="org-page compact-mode">
    <div class="page-hero">
      <div class="hero-content">
        <h2>组织管理中心</h2>
        <p>统一维护组织架构，快速查看成员归属关系，并在一个界面完成增删改查操作。</p>
        <div class="hero-stats">
          <div class="stat-item">
            <span>组织总数</span>
            <strong>{{ orgStats.total }}</strong>
          </div>
          <div class="stat-item">
            <span>叶子组织</span>
            <strong>{{ orgStats.leaf }}</strong>
          </div>
          <div class="stat-item">
            <span>最大层级</span>
            <strong>{{ orgStats.depth }}</strong>
          </div>
        </div>
      </div>
      <a-space class="hero-actions" wrap>
        <a-button type="primary" @click="addRootNode">
          <PlusOutlined />
          新增根组织
        </a-button>
        <a-button :disabled="!uiData.selectedNodeKey" @click="editSelectedNode">
          <EditOutlined />
          编辑当前组织
        </a-button>
      </a-space>
    </div>

    <div class="main-content">
      <SplitPanel :left-width="uiState.leftPanelWidth" :min-left-width="0" :collapsible="true" :collapse-threshold="0">
        <template #left-content>
          <section class="tree-panel">
            <div class="panel-header">
              <div>
                <div class="panel-title">
                  <ApartmentOutlined class="panel-icon" />
                  组织架构
                </div>
                <p class="panel-desc">点击节点后可在右侧管理成员绑定与角色关系。</p>
              </div>
              <a-input-search
                v-model:value="uiData.searchValue"
                placeholder="搜索组织名称"
                class="tree-search-input"
              />
            </div>

            <div class="panel-meta">
              <span>可见节点 {{ filteredNodeCount }}</span>
              <span v-if="uiData.selectedNodeKey" class="selected-node-field" :title="uiData.selectedNode">
                <span class="selected-node-label">当前选中：</span>
                <span class="selected-node-value">{{ uiData.selectedNode }}</span>
              </span>
            </div>

            <div class="tree-container">
              <a-empty v-if="uiData.treeData.length === 0" description="暂无组织数据">
                <a-button type="primary" @click="addRootNode">创建第一个组织</a-button>
              </a-empty>
              <a-tree
                v-else
                :tree-data="filteredTreeData"
                :selected-keys="uiData.selectedKeys"
                :expanded-keys="uiData.expandedKeys"
                :auto-expand-parent="uiData.autoExpandParent"
                block-node
                @select="handleTreeSelect"
                @expand="handleExpand"
              >
                <template #title="{ key: nodeKey, title, dataRef }">
                  <div class="tree-node" :class="{ 'is-selected': uiData.selectedKeys.includes(nodeKey) }">
                    <span class="tree-node-title">
                      <FolderOutlined class="folder-icon" />
                      <span class="tree-node-name-wrap">
                        <a-tooltip :title="title" placement="topLeft">
                          <span class="tree-node-name">{{ title }}</span>
                        </a-tooltip>
                      </span>
                    </span>
                    <div class="tree-actions">
                      <a-tooltip title="新增子组织">
                        <a-button type="text" @click.stop="addChildNode(dataRef)">
                          <PlusOutlined />
                        </a-button>
                      </a-tooltip>
                      <a-tooltip title="编辑">
                        <a-button type="text" @click.stop="editCurrentNode(dataRef)">
                          <EditOutlined />
                        </a-button>
                      </a-tooltip>
                      <a-popconfirm
                        title="确定要删除该组织吗？"
                        description="删除后将无法恢复"
                        @confirm="handleDelete(dataRef)"
                        ok-text="确定"
                        cancel-text="取消"
                      >
                        <a-tooltip title="删除">
                          <a-button type="text" danger @click.stop>
                            <DeleteOutlined />
                          </a-button>
                        </a-tooltip>
                      </a-popconfirm>
                    </div>
                  </div>
                </template>
              </a-tree>
            </div>
          </section>
        </template>

        <template #right-content>
          <section class="users-panel">
            <template v-if="uiData.selectedNodeKey">
              <OrgUsers :node-key="uiData.selectedNodeKey" :node-name="uiData.selectedNode" compact-mode />
            </template>
            <template v-else>
              <div class="empty-state">
                <div class="empty-icon">
                  <TeamOutlined />
                </div>
                <h3>先选择一个组织</h3>
                <p>从左侧组织树选择节点后，即可管理该组织的成员和角色。</p>
              </div>
            </template>
          </section>
        </template>
      </SplitPanel>
    </div>

    <AddRootOrg v-model:open="uiState.isAddRootNodeOpen" @refresh="fetchData" />

    <AddChildOrg
      v-model:open="uiState.isAddChildNodeOpen"
      :parent_node_key="uiData.selectedNodeKey"
      :parent_node_name="uiData.selectedNode"
      @refresh="fetchData"
    />

    <EditOrgName
      v-model:open="uiState.isEditNodeNameOpen"
      :node-key="uiData.selectedNodeKey"
      :node-name="uiData.selectedNode"
      @refresh="fetchData"
    />
  </div>
</template>

<script setup>
import { deleteOrganizationsApi, getOrganizationsApi } from '@/api/admin'
import SplitPanel from '@/components/panel/index.vue'
import {
  ApartmentOutlined,
  DeleteOutlined,
  EditOutlined,
  FolderOutlined,
  PlusOutlined,
  TeamOutlined,
} from '@ant-design/icons-vue'
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { computed, onMounted, reactive, watch } from 'vue'
import AddChildOrg from './AddChildOrg.vue'
import AddRootOrg from './AddRootOrg.vue'
import EditOrgName from './EditOrgName.vue'
import OrgUsers from './OrgUsers.vue'

// 状态
const uiState = reactive({
  isAddRootNodeOpen: false,
  isAddChildNodeOpen: false,
  isEditNodeNameOpen: false,
  leftPanelWidth: '420px',
})

// 数据
const uiData = reactive({
  treeData: [],
  selectedKeys: [],
  expandedKeys: [],
  autoExpandParent: true,
  selectedNodeKey: '',
  selectedNode: '',
  searchValue: '',
})

const getAllKeys = (data) => {
  let keys = []
  data.forEach((item) => {
    keys.push(item.key)
    if (item.children) {
      keys = keys.concat(getAllKeys(item.children))
    }
  })
  return keys
}

const findNodeByKey = (data, key) => {
  for (const item of data) {
    if (item.key === key) {
      return item
    }
    if (item.children?.length) {
      const childNode = findNodeByKey(item.children, key)
      if (childNode) {
        return childNode
      }
    }
  }
  return null
}

// 获取列表数据
const fetchData = async () => {
  const res = await getOrganizationsApi().catch(() => {})
  uiData.treeData = res?.data || []
  const allKeys = getAllKeys(uiData.treeData)
  uiData.expandedKeys = allKeys

  if (!uiData.selectedNodeKey) return

  if (!allKeys.includes(uiData.selectedNodeKey)) {
    uiData.selectedNodeKey = ''
    uiData.selectedNode = ''
    uiData.selectedKeys = []
    return
  }

  const selectedNode = findNodeByKey(uiData.treeData, uiData.selectedNodeKey)
  uiData.selectedNode = selectedNode?.title || ''
}

// 搜索过滤树数据
const filteredTreeData = computed(() => {
  if (!uiData.searchValue) return uiData.treeData
  return filterTree(uiData.treeData, uiData.searchValue.toLowerCase())
})

const filterTree = (data, keyword) => {
  return data
    .map((node) => {
      const children = node.children ? filterTree(node.children, keyword) : []
      if (node.title.toLowerCase().includes(keyword) || children.length > 0) {
        return { ...node, children }
      }
      return null
    })
    .filter(Boolean)
}

const filteredNodeCount = computed(() => getAllKeys(filteredTreeData.value).length)

const getAdaptiveLeftPanelWidth = () => {
  if (typeof window === 'undefined') return '420px'
  const viewportWidth = window.innerWidth || document.documentElement.clientWidth || 1440
  if (viewportWidth <= 1200) return '360px'
  if (viewportWidth <= 1440) return '400px'
  if (viewportWidth <= 1680) return '420px'
  return '460px'
}

const orgStats = computed(() => {
  const metrics = {
    total: 0,
    leaf: 0,
    depth: 0,
  }

  const walk = (nodes = [], level = 1) => {
    if (!nodes.length) return

    metrics.depth = Math.max(metrics.depth, level)

    nodes.forEach((node) => {
      metrics.total += 1
      if (!node.children?.length) {
        metrics.leaf += 1
      }
      walk(node.children || [], level + 1)
    })
  }

  walk(uiData.treeData)
  return metrics
})

// 组织树操作
const handleExpand = (keys) => {
  uiData.expandedKeys = keys
  uiData.autoExpandParent = false
}

const handleTreeSelect = (keys, { node }) => {
  // 始终保持节点选中状态，不允许取消选中
  const nodeKey = keys.length > 0 ? keys[0] : node.key
  uiData.selectedKeys = [nodeKey]
  uiData.selectedNodeKey = nodeKey
  uiData.selectedNode = node.title
}

const addRootNode = () => {
  uiState.isAddRootNodeOpen = true
}

const editSelectedNode = () => {
  if (!uiData.selectedNodeKey) return
  uiState.isEditNodeNameOpen = true
}

const addChildNode = (item) => {
  uiData.selectedNodeKey = item.key
  uiData.selectedNode = item.title
  uiData.selectedKeys = [item.key]
  uiState.isAddChildNodeOpen = true
}

const editCurrentNode = (item) => {
  uiData.selectedNodeKey = item.key
  uiData.selectedNode = item.title
  uiData.selectedKeys = [item.key]
  uiState.isEditNodeNameOpen = true
}

const handleDelete = useThrottleFn(async (item) => {
  const payload = {
    key: item.key,
    name: item.title,
  }
  const res = await deleteOrganizationsApi(payload).catch(() => {})
  if (res) {
    message.success('删除成功')
    uiData.selectedNodeKey = ''
    uiData.selectedNode = ''
    uiData.selectedKeys = []
    await fetchData()
  }
})

watch(
  () => uiData.searchValue,
  (value) => {
    if (!value) {
      uiData.expandedKeys = getAllKeys(uiData.treeData)
      uiData.autoExpandParent = false
      return
    }
    uiData.expandedKeys = getAllKeys(filteredTreeData.value)
    uiData.autoExpandParent = true
  },
)

// 初始化
onMounted(async () => {
  uiState.leftPanelWidth = getAdaptiveLeftPanelWidth()
  await fetchData()
})
</script>

<style scoped>
.org-page {
  --org-bg-soft: #f3f7ff;
  --org-bg-card: #ffffff;
  --org-border: #dce5f5;
  --org-text-main: #16213c;
  --org-text-sub: #5f6b8a;
  --org-accent: #1f6feb;
  --org-accent-deep: #154ea4;
  --org-shadow-lg: 0 18px 45px -28px rgba(25, 55, 115, 0.45);
  --org-shadow-sm: 0 10px 24px -22px rgba(17, 35, 78, 0.5);
  font-family: 'Avenir Next', 'PingFang SC', 'Hiragino Sans GB', 'Noto Sans SC', 'Microsoft YaHei', sans-serif;
  background: #ffffff;
  border: 1px solid #eceff5;
  border-radius: 20px;
  padding: 20px;
}

.page-hero {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 24px;
  border-radius: 16px;
  border: 1px solid var(--org-border);
  background:
    radial-gradient(circle at 90% 0%, rgba(22, 84, 194, 0.1), rgba(22, 84, 194, 0) 50%),
    var(--org-bg-card);
  box-shadow: var(--org-shadow-sm);
  padding: 22px 24px;
}

.hero-content h2 {
  margin: 10px 0 8px;
  color: var(--org-text-main);
  font-size: 28px;
  font-weight: 700;
  letter-spacing: 0.4px;
}

.hero-content p {
  margin: 0;
  color: var(--org-text-sub);
  font-size: 14px;
  line-height: 1.7;
  max-width: 680px;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 999px;
  color: var(--org-accent-deep);
  border: 1px solid rgba(31, 111, 235, 0.28);
  background: rgba(31, 111, 235, 0.08);
  font-size: 12px;
  font-weight: 600;
}

.hero-stats {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 16px;
}

.stat-item {
  min-width: 124px;
  padding: 10px 12px;
  border-radius: 12px;
  background: var(--org-bg-soft);
  border: 1px solid rgba(31, 111, 235, 0.15);
}

.stat-item span {
  display: block;
  color: var(--org-text-sub);
  font-size: 12px;
  margin-bottom: 6px;
}

.stat-item strong {
  color: var(--org-text-main);
  font-size: 24px;
  line-height: 1;
}

.main-content {
  margin-top: 14px;
  min-height: calc(100vh - 220px);
}

.main-content :deep(.split-wrapper) {
  --split-panel-gap: 10px;
  --split-handle-bg: #f5f7fa;
  --split-handle-border: var(--org-border);
  height: auto;
  min-height: calc(100vh - 220px);
  gap: 0;
}

.main-content :deep(.split-wrapper .left-content) {
  padding: 0;
}

.main-content :deep(.split-wrapper .right-content) {
  padding: 0;
  overflow: hidden;
}

.main-content :deep(.split-wrapper .separator),
.main-content :deep(.split-wrapper .collapsed-handle) {
  border-color: var(--org-border);
  background-color: #f5f7fa;
}

.main-content :deep(.split-wrapper .separator) {
  top: 0;
  bottom: 0;
  width: 10px;
  border-radius: 0;
}

.main-content :deep(.split-wrapper .separator i),
.main-content :deep(.split-wrapper .collapsed-handle i) {
  background-color: rgb(95 107 138 / 48%);
}

.tree-panel {
  width: 100%;
  height: 100%;
  background: var(--org-bg-card);
  border-radius: 0;
  border: 1px solid var(--org-border);
  box-shadow: var(--org-shadow-sm);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.users-panel {
  width: 100%;
  height: 100%;
  min-height: 0;
  background: var(--org-bg-card);
  border-radius: 0;
  border: 1px solid var(--org-border);
  box-shadow: var(--org-shadow-lg);
  overflow: hidden;
}

.panel-header {
  padding: 18px 20px 12px;
  border-bottom: 1px solid #ecf1fb;
  background: linear-gradient(180deg, rgba(245, 249, 255, 1) 0%, rgba(255, 255, 255, 1) 100%);
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--org-text-main);
  font-size: 16px;
  font-weight: 700;
}

.panel-desc {
  margin: 6px 0 0;
  color: var(--org-text-sub);
  font-size: 12px;
  line-height: 1.6;
}

.tree-search-input {
  width: 100%;
  margin-top: 14px;
}

.panel-icon {
  font-size: 16px;
  color: var(--org-accent);
}

.panel-meta {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  flex-wrap: wrap;
  padding: 10px 20px;
  font-size: 12px;
  color: var(--org-text-sub);
  border-bottom: 1px dashed #ebf1fc;
}

.selected-node-field {
  display: inline-flex;
  align-items: center;
  min-width: 0;
  max-width: min(60%, 320px);
}

.selected-node-label {
  flex-shrink: 0;
}

.selected-node-value {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tree-container {
  flex: 1;
  overflow-y: auto;
  padding: 12px 14px 16px;
}

.tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 10px;
  min-height: 39px;
  border-radius: 10px;
  transition: all 0.2s ease;
  margin: 3px 0;
  border: 1px solid transparent;
}

.tree-node:hover {
  background: #f5f9ff;
  border-color: #dde8fc;
}

.tree-node.is-selected {
  background: linear-gradient(90deg, rgba(31, 111, 235, 0.11), rgba(31, 111, 235, 0.03));
  border-color: rgba(31, 111, 235, 0.35);
}

.tree-node-title {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 14px;
  color: #22304f;
}

.tree-node-name {
  display: block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tree-node-name-wrap {
  display: block;
  flex: 1;
  min-width: 0;
  overflow: hidden;
}

.tree-node-name-wrap :deep(.ant-tooltip-open) {
  display: block;
  min-width: 0;
  max-width: 100%;
  overflow: hidden;
}

.folder-icon {
  color: #f5a623;
  font-size: 16px;
}

.tree-actions {
  display: flex;
  align-items: center;
  gap: 2px;
  flex-shrink: 0;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.tree-node:hover .tree-actions {
  opacity: 1;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 400px;
  color: var(--org-text-sub);
  text-align: center;
  padding: 0 16px;
}

.empty-icon {
  width: 92px;
  height: 92px;
  background: linear-gradient(135deg, rgba(31, 111, 235, 0.18), rgba(33, 184, 132, 0.18));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
  font-size: 40px;
  color: var(--org-accent);
}

.empty-state h3 {
  margin: 0 0 8px;
  font-size: 20px;
  font-weight: 700;
  color: var(--org-text-main);
}

.empty-state p {
  margin: 0 0 18px;
  font-size: 14px;
  line-height: 1.6;
}

/* 覆盖 ant-design 树组件样式 */
:deep(.ant-tree) {
  background: transparent;
}

:deep(.ant-tree-treenode) {
  display: flex;
  align-items: center;
  padding: 0;
  width: 100%;
}

:deep(.ant-tree-node-content-wrapper) {
  flex: 1;
  min-height: auto;
  padding: 0;
  line-height: 1;
}

:deep(.ant-tree-node-content-wrapper:hover) {
  background: transparent;
}

:deep(.ant-tree-node-selected) {
  background: transparent !important;
}

:deep(.ant-tree-switcher) {
  width: 24px;
  height: auto;
  min-height: 39px;
  line-height: 1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

:deep(.ant-empty) {
  margin: 40px 0;
}

.compact-mode {
  padding: 14px;
}

.compact-mode .page-hero {
  padding: 16px 18px;
}

.compact-mode .hero-content h2 {
  font-size: 22px;
  margin: 6px 0;
}

.compact-mode .hero-content p {
  font-size: 13px;
  line-height: 1.5;
}

.compact-mode .hero-stats {
  gap: 8px;
  margin-top: 12px;
}

.compact-mode .stat-item {
  padding: 8px 10px;
  min-width: 104px;
}

.compact-mode .stat-item strong {
  font-size: 18px;
}

.compact-mode .main-content {
  margin-top: 12px;
  min-height: calc(100vh - 190px);
}

.compact-mode .main-content :deep(.split-wrapper) {
  min-height: calc(100vh - 190px);
}

.compact-mode .panel-header {
  padding: 14px 14px 10px;
}

.compact-mode .panel-meta {
  padding: 8px 14px;
}

.compact-mode .tree-container {
  padding: 8px 10px 12px;
}

.compact-mode .tree-node {
  padding: 6px 8px;
  min-height: 35px;
}

.compact-mode :deep(.ant-tree-switcher) {
  min-height: 35px;
}

.compact-mode .empty-state {
  min-height: 320px;
}

.compact-mode .empty-icon {
  width: 74px;
  height: 74px;
  margin-bottom: 14px;
  font-size: 30px;
}

.compact-mode .empty-state h3 {
  font-size: 17px;
}

@media (max-width: 1280px) {
  .page-hero {
    flex-direction: column;
  }

  .hero-actions {
    width: 100%;
    justify-content: flex-start;
    flex-wrap: wrap;
  }
}

@media (max-width: 1080px) {
  .org-page {
    border-radius: 16px;
    padding: 14px;
  }

  .main-content {
    min-height: auto;
  }

  .main-content :deep(.split-wrapper) {
    display: block;
    min-height: 0;
  }

  .main-content :deep(.split-wrapper .scalable) {
    width: 100% !important;
    max-width: none;
  }

  .main-content :deep(.split-wrapper .separator),
  .main-content :deep(.split-wrapper .collapsed-handle) {
    display: none;
  }

  .main-content :deep(.split-wrapper .right-content) {
    padding-top: 12px;
  }

  .users-panel {
    min-height: 520px;
  }

  .tree-actions {
    opacity: 1;
  }
}
</style>
