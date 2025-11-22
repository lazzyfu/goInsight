<template>
  <a-card title="组织管理">
    <a-row :gutter="16">
      <!-- 左侧组织节点 -->
      <a-col class="gutter-row" :span="state.visible ? 10 : 24">
        <a-card title="组织列表">
          <template #extra>
            <a-button type="primary" @click="addRootNode()"> <PlusOutlined /> 新增根节点 </a-button>
          </template>
          <!-- <div class="tree-container"> -->
          <a-tree block-node :tree-data="treeData" show-line @select="handleTreeSelect">
            <template #title="{ dataRef }">
              <div class="tree-node">
                <span class="tree-node-title">{{ dataRef.title }}</span>
                <div class="actions">
                  <a-dropdown>
                    <EllipsisOutlined />
                    <template #overlay>
                      <a-menu>
                        <a-menu-item key="1">
                          <a-tooltip title="新增子节点">
                            <a @click="addChildNode(dataRef)"> <PlusOutlined /> 新增 </a>
                          </a-tooltip>
                        </a-menu-item>
                        <a-menu-item key="2">
                          <a-tooltip title="编辑当前节点">
                            <a @click="editCurrentNode(dataRef)"> <EditOutlined /> 编辑 </a>
                          </a-tooltip>
                        </a-menu-item>
                        <a-menu-item key="3">
                          <a-tooltip title="删除当前节点">
                            <a-popconfirm
                              title="确认删除吗？"
                              ok-text="是"
                              cancel-text="否"
                              @confirm="handleDelete(dataRef)"
                            >
                              <a> <DeleteOutlined /> 删除 </a>
                            </a-popconfirm>
                          </a-tooltip>
                        </a-menu-item>
                      </a-menu>
                    </template>
                  </a-dropdown>
                </div>
              </div>
            </template>
          </a-tree>
          <!-- </div> -->
        </a-card>
      </a-col>
      <!-- 右侧详情区/用户列表 -->
      <a-col v-show="state.visible" class="gutter-row" :span="state.visible ? 14 : 0">
        <OrgUsers
          :open="state.isNodeUsersOpen"
          @update:open="state.isNodeUsersOpen = $event"
          :nodeKey="selectedNodeKey"
        ></OrgUsers>
      </a-col>
    </a-row>
    <!-- 增加root节点 -->
    <AddRootOrg
      :open="state.isAddRootNodeOpen"
      @update:open="state.isAddRootNodeOpen = $event"
      @refresh="fetchData"
    />
    <!-- 增加子节点 -->
    <AddChildOrg
      :open="state.isAddChildNodeOpen"
      @update:open="state.isAddChildNodeOpen = $event"
      :parent_node_key="selectedNodeKey"
      :parent_node_name="selectedNode"
      @refresh="fetchData"
    ></AddChildOrg>
    <!-- 编辑节点名 -->
    <EditOrgName
      :open="state.isEditNodeNameOpen"
      @update:open="state.isEditNodeNameOpen = $event"
      :nodeKey="selectedNodeKey"
      @refresh="fetchData"
    ></EditOrgName>
  </a-card>
</template>

<script setup>
import { deleteOrganizationsApi, getOrganizationsApi } from '@/api/admin'
import { DeleteOutlined, EditOutlined, EllipsisOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref } from 'vue'
import AddChildOrg from './AddChildOrg.vue'
import AddRootOrg from './AddRootOrg.vue'
import EditOrgName from './EditOrgName.vue'
import OrgUsers from './OrgUsers.vue'

const state = reactive({
  visible: false,
  isAddRootNodeOpen: false,
  isNodeUsersOpen: false,
  isAddChildNodeOpen: false,
  isEditNodeNameOpen: false,
})

const treeData = ref([])
const selectedNodeKey = ref('')
const selectedNode = ref('')

const addRootNode = () => {
  state.isAddRootNodeOpen = true
}

const addChildNode = (item) => {
  if (item.key !== selectedNodeKey.value) {
    message.warning('请先鼠标选中需要增加的节点')
    return
  }
  state.isAddChildNodeOpen = true
  selectedNodeKey.value = item.key
  selectedNode.value = item.title
}

const editCurrentNode = (item) => {
  if (item.key !== selectedNodeKey.value) {
    message.warning('请先鼠标选中需要增加的节点')
    return
  }
  state.isEditNodeNameOpen = true
}

const handleDelete = async (item) => {
  if (item.key !== selectedNodeKey.value) {
    message.warning('请先鼠标选中需要删除的节点')
    return
  }
  const payload = {
    key: item.key,
    name: item.title,
  }
  const res = await deleteOrganizationsApi(payload).catch(() => {})
  if (res?.code === '0000') {
    message.info('操作成功')
    await fetchData()
  }
}

const handleTreeSelect = (selectedKeys, keys) => {
  if (keys.selected) {
    state.visible = true
    selectedNodeKey.value = selectedKeys[0]
    state.isNodeUsersOpen = true
  } else {
    state.visible = false
  }
}

const fetchData = async () => {
  const res = await getOrganizationsApi().catch(() => {})
  treeData.value = res?.data || []
}

onMounted(async () => {
  await fetchData()
})
</script>

<style scoped>
.tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-right: 8px;
  border-radius: 4px;
}

.tree-node-title {
  flex-grow: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 28px;
  padding-right: 16px;
  color: #333;
}

.actions {
  opacity: 0;
  visibility: hidden;
  transition:
    opacity 0.15s ease-in-out,
    visibility 0.15s;
  display: flex;
  align-items: center;
  margin-left: 8px;
}

.action-icon {
  font-size: 16px;
  color: #999;
  padding: 4px;
  border-radius: 4px;
  cursor: pointer;
}
.action-icon:hover {
  color: #1890ff;
  background-color: #f0f2f5;
}

:deep(.ant-tree-node-content-wrapper:hover .actions) {
  opacity: 1;
  visibility: visible;
}

:deep(.ant-tree-node-content-wrapper.ant-tree-node-selected) {
  background-color: #e6f7ff !important;
  border-left: 3px solid #1890ff;
  padding-left: 5px !important;
}

:deep(.ant-tree-node-content-wrapper.ant-tree-node-selected:hover) {
  background-color: #e6f7ff !important;
}

:deep(.ant-tree-list-holder-inner) {
  padding-right: 0 !important;
}

:deep(.ant-tree-treenode) {
  padding: 4px 0;
}
</style>
