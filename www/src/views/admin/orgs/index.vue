<template>
  <a-card title="组织管理">
    <template #extra>
      <a-button type="primary" @click="addRootNode()"> <PlusOutlined /> 新增根节点 </a-button>
    </template>
    <a-row :gutter="16">
      <a-col class="gutter-row" :span="state.visible ? 10 : 24">
        <a-tree block-node :tree-data="treeData" show-icon @select="onSelect">
          <template #title="{ dataRef }">
            <div class="tree-node">
              <span>{{ dataRef.title }}</span>
              <div class="actions">
                <a-tooltip title="新增子节点">
                  <a-button type="text" @click.stop="addChildNode(dataRef)">
                    <PlusOutlined />
                  </a-button>
                </a-tooltip>
                <a-tooltip title="编辑当前节点">
                  <a-button type="text" @click.stop="editChildNode(dataRef)">
                    <EditOutlined />
                  </a-button>
                </a-tooltip>
                <a-tooltip title="删除当前节点">
                  <a-button type="text" @click.stop="DeleteConfirm(dataRef)">
                    <DeleteOutlined />
                  </a-button>
                </a-tooltip>
              </div>
            </div>
          </template>
        </a-tree>
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
          :nodeKey="nodeKey"
        ></NodeUsers>
      </a-col>
      <!-- <RootNodeAddComponent ref="RootNodeAddComponent" @refresh="getOrganizations"></RootNodeAddComponent>
      <ChildNodeAddComponent ref="ChildNodeAddComponent" @refresh="getOrganizations"></ChildNodeAddComponent>
      <NodeEditComponent ref="NodeEditComponent" @refresh="getOrganizations"></NodeEditComponent> -->
    </a-row>
  </a-card>
</template>

<script setup>
import { getOrganizationsApi } from '@/api/admin'
import AddRootNodeModal from '@/views/admin/orgs/AddRootNodeModal.vue'
import NodeUsers from '@/views/admin/orgs/NodeUsers.vue'
import { onMounted, reactive, ref } from 'vue'

import { DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons-vue'

const state = reactive({
  visible: false,
  isAddRootNodeOpen: false,
  isNodeUsersOpen: false,
})

const treeData = ref([])
const nodeKey = ref('')

const addRootNode = () => {
  state.isAddRootNodeOpen = true
}

const addChildNode = (item) => {
  console.log('addChildNode', item)
}

const editChildNode = (item) => {
  console.log('editChildNode', item)
}

const DeleteConfirm = (item) => {
  console.log('DeleteConfirm', item)
}

const onSelect = (selectedKeys, keys) => {
  if (keys.selected) {
    state.visible = true
    nodeKey.value = selectedKeys[0]
    console.log('nodeKey.value: ', nodeKey.value)
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
</style>
