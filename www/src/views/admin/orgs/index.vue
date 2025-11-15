<template>
  <a-card title="组织管理">
    <template #extra>
      <a-button type="primary" @click="addRootNode()"> <PlusOutlined /> 新增根节点 </a-button>
    </template>
    <a-row :gutter="16">
      <a-col class="gutter-row" :span="state.visible ? 10 : 24">
        <a-card title="组织列表" class="box-card">
          <a-tree
            block-node
            :tree-data="treeData"
            show-icon
            :expanded-keys="expandedKeys"
            @update:expandedKeys="(val) => (expandedKeys = val)"
            @select="onSelect"
          >
            <!-- icon 插槽 -->

            <template #icon>
              <ApartmentOutlined />
            </template>

            <!-- title 插槽 -->
            <template #title="{ dataRef }">
              <span>{{ dataRef.title }}</span>
              <a-tooltip title="新增子节点">
                <a-button
                  type="text"
                  class="btn-type"
                  style="right: 80px"
                  @click.stop="addChildNode(dataRef)"
                >
                  <PlusOutlined />
                </a-button>
              </a-tooltip>
              <a-tooltip title="编辑当前节点">
                <a-button
                  type="text"
                  class="btn-type"
                  style="right: 50px"
                  @click.stop="editChildNode(dataRef)"
                >
                  <EditOutlined />
                </a-button>
              </a-tooltip>
              <a-tooltip title="删除当前节点">
                <a-button type="text" class="btn-type" @click.stop="DeleteConfirm(record)">
                  <DeleteOutlined />
                </a-button>
              </a-tooltip>
            </template>
          </a-tree>
        </a-card>
      </a-col>
      <!-- 新增弹窗 -->
      <RootNodeAddModal :open="state.isRootNodeAdd" @update:open="state.isRootNodeAdd = $event" />
      <!-- <a-col v-show="state.visible" class="gutter-row" :span="state.visible ? 14 : 0">
        <NodeUsers ref="NodeUsers" :nodeKey="nodeKey"></NodeUsers>
      </a-col> -->
      <!-- <RootNodeAddComponent ref="RootNodeAddComponent" @refresh="getOrganizations"></RootNodeAddComponent>
      <ChildNodeAddComponent ref="ChildNodeAddComponent" @refresh="getOrganizations"></ChildNodeAddComponent>
      <NodeEditComponent ref="NodeEditComponent" @refresh="getOrganizations"></NodeEditComponent> -->
    </a-row>
  </a-card>
</template>

<script setup>
import { getOrganizationsApi } from '@/api/admin'
import RootNodeAddModal from '@/views/admin/orgs/RootNode.vue'
import { onMounted, reactive, ref } from 'vue'

import {
  ApartmentOutlined,
  DeleteOutlined,
  EditOutlined,
  PlusOutlined,
} from '@ant-design/icons-vue'

const state = reactive({
  visible: false,
  isRootNodeAdd: false,
})
const expandedKeys = ref([])

const treeData = ref([])

const addRootNode = () => {
  state.isRootNodeAdd = true
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
  console.log('selectedKeys, keys: ', selectedKeys, keys)
}

onMounted(async () => {
  const res = await getOrganizationsApi().catch(() => {})
  treeData.value = res?.data || []
  console.log('treeData.value: ', treeData.value)
})
</script>
