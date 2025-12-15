<template>
  <div class="org-management">
    <div class="page-header">
      <div class="header-content">
        <div class="header-title">
          <h1>组织管理</h1>
          <p class="header-desc">管理组织架构和用户分配</p>
        </div>
        <a-button type="primary" @click="addRootNode">
          <template #icon>
            <PlusOutlined />
          </template>
          新增根组织
        </a-button>
      </div>
    </div>

    <div class="main-content">
      <div class="tree-panel">
        <div class="panel-header">
          <span class="panel-title">
            <ApartmentOutlined class="panel-icon" />
            组织架构
          </span>
          <a-input-search
            v-model:value="uiData.searchValue"
            placeholder="搜索组织"
            style="width: 180px"
            allow-clear
          />
        </div>
        <div class="tree-container">
          <a-empty v-if="uiData.treeData.length === 0" description="暂无组织数据">
            <a-button type="primary" size="small" @click="addRootNode"> 创建第一个组织 </a-button>
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
              <div
                class="tree-node"
                :class="{ 'is-selected': uiData.selectedKeys.includes(nodeKey) }"
              >
                <span class="tree-node-title" :title="title">
                  <FolderOutlined class="folder-icon" />
                  {{ title }}
                </span>
                <div class="tree-actions">
                  <a-tooltip title="新增子组织">
                    <a-button
                      type="text"
                      size="small"
                      class="action-btn"
                      @click.stop="addChildNode(dataRef)"
                    >
                      <PlusOutlined />
                    </a-button>
                  </a-tooltip>
                  <a-tooltip title="编辑">
                    <a-button
                      type="text"
                      size="small"
                      class="action-btn"
                      @click.stop="editCurrentNode(dataRef)"
                    >
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
                      <a-button type="text" size="small" class="action-btn danger" @click.stop>
                        <DeleteOutlined />
                      </a-button>
                    </a-tooltip>
                  </a-popconfirm>
                </div>
              </div>
            </template>
          </a-tree>
        </div>
      </div>

      <div class="users-panel">
        <template v-if="uiData.selectedNodeKey">
          <OrgUsers :node-key="uiData.selectedNodeKey" />
        </template>
        <template v-else>
          <div class="empty-state">
            <div class="empty-icon">
              <TeamOutlined />
            </div>
            <h3>请选择组织</h3>
            <p>从左侧选择一个组织以查看和管理其成员</p>
          </div>
        </template>
      </div>
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
      @refresh="fetchData"
    />
  </div>
</template>

<script setup>
import { deleteOrganizationsApi, getOrganizationsApi } from '@/api/admin'
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
import { computed, onMounted, reactive } from 'vue'
import AddChildOrg from './AddChildOrg.vue'
import AddRootOrg from './AddRootOrg.vue'
import EditOrgName from './EditOrgName.vue'
import OrgUsers from './OrgUsers.vue'

// 状态
const uiState = reactive({
  isAddRootNodeOpen: false,
  isAddChildNodeOpen: false,
  isEditNodeNameOpen: false,
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

// 获取列表数据
const fetchData = async () => {
  const res = await getOrganizationsApi().catch(() => {})
  uiData.treeData = res?.data || []
  if (uiData.treeData.length > 0) {
    const allKeys = getAllKeys(uiData.treeData)
    uiData.expandedKeys = allKeys
  }
}

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

// 初始化
onMounted(async () => {
  await fetchData()
})
</script>

<style scoped>
.org-management {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8ec 100%);
  padding: 12px;
}

.page-header {
  margin-bottom: 12px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fff;
  padding: 12px 32px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.header-title h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #1a1a2e;
}

.header-desc {
  margin: 4px 0 0;
  color: #8c8c8c;
  font-size: 14px;
}

.main-content {
  display: flex;
  gap: 12px;
  min-height: calc(100vh - 180px);
}

.tree-panel {
  width: 360px;
  flex-shrink: 0;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.users-panel {
  flex: 1;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  overflow: hidden;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
  background: #fafafa;
}

.panel-title {
  font-size: 15px;
  font-weight: 600;
  color: #1a1a2e;
  display: flex;
  align-items: center;
  gap: 8px;
}

.panel-icon {
  font-size: 18px;
  color: #1890ff;
}

.tree-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border-radius: 8px;
  transition: all 0.2s ease;
  margin: 2px 0;
}

.tree-node:hover {
  background: #f5f5f5;
}

.tree-node.is-selected {
  background: #e6f7ff;
}

.tree-node-title {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 14px;
  color: #333;
}

.folder-icon {
  color: #faad14;
  font-size: 16px;
}

.tree-actions {
  display: flex;
  align-items: center;
  gap: 2px;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.tree-node:hover .tree-actions {
  opacity: 1;
}

.action-btn {
  width: 28px;
  height: 28px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
  border-radius: 6px;
}

.action-btn:hover {
  color: #1890ff;
  background: #e6f7ff;
}

.action-btn.danger:hover {
  color: #ff4d4f;
  background: #fff1f0;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 400px;
  color: #8c8c8c;
}

.empty-icon {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #e6f7ff 0%, #bae7ff 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
  font-size: 36px;
  color: #1890ff;
}

.empty-state h3 {
  margin: 0 0 8px;
  font-size: 18px;
  font-weight: 500;
  color: #333;
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}

/* 覆盖 ant-design 树组件样式 */
:deep(.ant-tree) {
  background: transparent;
}

:deep(.ant-tree-treenode) {
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
  height: 40px;
  line-height: 40px;
}

:deep(.ant-empty) {
  margin: 40px 0;
}
</style>
