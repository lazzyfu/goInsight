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
                        <a-button type="text" size="small" @click.stop>
                          <template #icon><DeleteOutlined /></template>
                        </a-button>
                      </a-popconfirm>
                    </a-tooltip>
                  </div>
                </div>
              </template>
            </a-tree>
          <!-- </div> -->
        </a-card>
      </a-col>
      <!-- 右侧详情区/用户列表 -->
      <a-col v-show="state.visible" class="gutter-row" :span="state.visible ? 14 : 0">
        <NodeUsers
          :open="state.isNodeUsersOpen"
          @update:open="state.isNodeUsersOpen = $event"
          :nodeKey="selectedNodeKey"
        ></NodeUsers>
      </a-col>
    </a-row>
    <!-- 增加root节点 -->
    <AddRootNode :open="state.isAddRootNodeOpen" @update:open="state.isAddRootNodeOpen = $event" />
    <!-- 增加子节点 -->
    <AddChildNode
      :open="state.isAddChildNodeOpen"
      @update:open="state.isAddChildNodeOpen = $event"
      :parent_node_key="selectedNodeKey"
      :parent_node_name="selectedNode"
    ></AddChildNode>
    <!-- 编辑节点名 -->
    <EditNodeName
      :open="state.isEditNodeNameOpen"
      @update:open="state.isEditNodeNameOpen = $event"
      :nodeKey="selectedNodeKey"
    ></EditNodeName>
  </a-card>
</template>

<script setup>
import { deleteOrganizationsApi, getOrganizationsApi } from '@/api/admin'
import { DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref } from 'vue'
import AddChildNode from './AddChildNode.vue'
import AddRootNode from './AddRootNode.vue'
import EditNodeName from './EditNodeName.vue'
import NodeUsers from './NodeUsers.vue'

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
  state.isAddChildNodeOpen = true
  selectedNodeKey.value = item.key
  selectedNode.value = item.title
}

const editChildNode = (item) => {
  state.isEditNodeNameOpen = true
}

const DeleteConfirm = async (item) => {
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

const handleTreeSelect = (selectedKeys, keys) => {
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
})
</script>

<style scoped>
/* 根容器 */
.tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-radius: 4px;
  transition: background-color 0.12s;
  gap: 6px;
}

.tree-node:hover {
  background-color: #f6f6f6;
}

/* 左侧结构：图标 + 文本 */
.node-left {
  display: flex;
  align-items: center;
  gap: 4px;
  min-width: 0;
}

/* 标题（重要：长文本避免撑开） */
.node-title {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 操作区 */
.node-actions {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* ⭐⭐⭐ 修复按钮和文字错位关键代码 ⭐⭐⭐
   AntD 的 button 内部 line-height ≠ icon 默认高度 → 造成垂直错位
*/
:deep(.node-actions .ant-btn) {
  padding: 0;
  width: 26px;
  height: 26px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 图标垂直居中 */
:deep(.node-actions .ant-btn .anticon) {
  font-size: 14px;
  line-height: 1;
  display: flex;
  align-items: center;
}



</style>
