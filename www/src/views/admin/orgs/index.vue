<template>
  <a-card title="组织管理">
    <template #extra>
      <a-button type="primary" @click="addRootNode()"> <PlusOutlined /> 新增根节点 </a-button>
    </template>
    <a-row :gutter="16">
      <a-col class="gutter-row" :span="state.visible ? 10 : 24">
        <div class="tree-container">
          <a-tree block-node :tree-data="treeData" show-icon @select="onSelect">
            <template #title="{ dataRef }">
              <div class="tree-node">
                <span class="tree-node-title">{{ dataRef.title }}</span>
                <div class="actions">
                  <a-tooltip title="新增子节点">
                    <a-button type="text" size="small" @click.stop="addChildNode(dataRef)">
                      <template #icon><PlusOutlined /></template>
                    </a-button>
                  </a-tooltip>
                  <a-tooltip title="编辑当前节点">
                    <a-button type="text" size="small" @click.stop="editChildNode(dataRef)">
                      <template #icon><EditOutlined /></template>
                    </a-button>
                  </a-tooltip>
                  <a-tooltip title="删除当前节点">
                    <a-popconfirm
                      title="确认删除吗？"
                      ok-text="是"
                      cancel-text="否"
                      @confirm="DeleteConfirm(dataRef)"
                    >
                      <a-button type="text" size="small" danger @click.stop>
                        <template #icon><DeleteOutlined /></template>
                      </a-button>
                    </a-popconfirm>
                  </a-tooltip>
                </div>
              </div>
            </template>
          </a-tree>
        </div>
      </a-col>
      <!-- 新增弹窗 -->
      <AddRootNodeModal
        :open="state.isAddRootNodeOpen"
        @update:open="state.isAddRootNodeOpen = $event"
      />
      <a-col v-show="state.visible" class="gutter-row" :span="state.visible ? 14 : 0">
        <NodeUsers
          :open="state.isNodeUsersOpen"
          @update:open="state.isNodeUsersOpen = $event"
          :nodeKey="selectedNodeKey"
        ></NodeUsers>
      </a-col>

      <AddChildNode
        :open="state.isChildNodeOpen"
        @update:open="state.isChildNodeOpen = $event"
        :parent_node_key="selectedNodeKey"
        :parent_node_name="selectedNode"
      ></AddChildNode>
      <EditNode
        :open="state.isEditNodeOpen"
        @update:open="state.isEditNodeOpen = $event"
        :nodeKey="selectedNodeKey"
      ></EditNode>
    </a-row>
  </a-card>
</template>

<script setup>
import { deleteOrganizationsApi, getOrganizationsApi } from '@/api/admin'
import AddChildNode from '@/views/admin/orgs/AddChildNode.vue'
import AddRootNodeModal from '@/views/admin/orgs/AddRootNodeModal.vue'
import EditNode from '@/views/admin/orgs/EditNode.vue'
import NodeUsers from '@/views/admin/orgs/NodeUsers.vue'
import { onMounted, reactive, ref } from 'vue'

import { DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons-vue'

const state = reactive({
  visible: false,
  isAddRootNodeOpen: false,
  isNodeUsersOpen: false,
  isNodeUsersOpen: false,
  isEditNodeOpen: false,
})

const treeData = ref([])
const selectedNodeKey = ref('')
const selectedNode = ref('')

const addRootNode = () => {
  state.isAddRootNodeOpen = true
}

const addChildNode = (item) => {
  state.isChildNodeOpen = true
  selectedNodeKey.value = item.key
  console.log('selectedNodeKey.value: ', selectedNodeKey.value)
  selectedNode.value = item.title
  console.log('selectedNode.value : ', selectedNode.value)
}

const editChildNode = (item) => {
  state.isEditNodeOpen = true
  console.log('editChildNode', item)
}

const DeleteConfirm = async (item) => {
  console.log('DeleteConfirm', item)
  const payload = {
    key: item.key,
    name: item.title,
  }
  const res = await deleteOrganizationsApi(payload).catch(() => {})
  if (res?.code === '0000') {
    message.info('操作成功')
    fetchData()
  }
}

const onSelect = (selectedKeys, keys) => {
  if (keys.selected) {
    state.visible = true
    selectedNodeKey.value = selectedKeys[0]
    state.isNodeUsersOpen = true
  } else {
    state.visible = false
  }
}

onMounted(async () => {
  const res = await getOrganizationsApi().catch(() => {})
  treeData.value = res?.data || []
  console.log('treeData.value: ', treeData.value)
})
</script>

<style scoped>
:deep(.tree-node) {
  display: flex;
  align-items: center;
  justify-content: space-between; /* 让按钮靠右 */
  width: 100%;
}

:deep(.actions) {
  display: flex;
  align-items: center;
}
:deep(.tree-container) {
  padding: 24px;
  background: #fff;
}

:deep(.tree-node) {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  gap: 8px;
}

.tree-node-title {
  flex: 1;
  line-height: 32px;
}

.actions {
  display: flex;
  align-items: center;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.tree-node:hover .actions {
  opacity: 1;
}

:deep(.ant-tree-node-content-wrapper) {
  flex: 1;
  min-width: 0;
}
</style>
