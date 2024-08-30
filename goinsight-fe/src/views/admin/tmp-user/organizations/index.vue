<template>
  <a-card>
    <a-row :gutter="16">
      <a-col class="gutter-row" :span="nodeUsersVisible ? 10 : 24">
        <a-card title="组织列表" class="box-card">
          <a-button slot="extra" type="primary" @click="addRootNode()"> <a-icon type="plus" /> 新增根节点 </a-button>
          <a-tree blockNode :tree-data="treeData" show-icon @select="onSelect" :expandedKeys.sync="expandedKeys">
            <a-icon slot="dir" type="apartment" />
            <template slot="custom" slot-scope="item">
              <span>{{ item.title }}</span>
              <a-tooltip>
                <template slot="title"> 新增子节点 </template>
                <a type="btn" class="btn-type" style="right: 80px" @click="addChildNode(item)">
                  <a-icon type="plus-circle"
                /></a>
              </a-tooltip>
              <a-tooltip>
                <template slot="title"> 编辑当前节点 </template>
                <a type="btn" class="btn-type" style="right: 50px" @click="editChildNode(item)"
                  ><a-icon type="edit"
                /></a>
              </a-tooltip>
              <a-tooltip>
                <template slot="title"> 删除当前节点 </template>
                <a type="btn" class="btn-type" @click="DeleteConfirm(item)"><a-icon type="delete" /></a>
              </a-tooltip>
            </template>
          </a-tree>
        </a-card>
      </a-col>
      <a-col v-show="nodeUsersVisible" class="gutter-row" :span="nodeUsersVisible ? 14 : 0">
        <NodeUsers ref="NodeUsers" :nodeKey="nodeKey"></NodeUsers>
      </a-col>
      <RootNodeAddComponent ref="RootNodeAddComponent" @refresh="getOrganizations"></RootNodeAddComponent>
      <ChildNodeAddComponent ref="ChildNodeAddComponent" @refresh="getOrganizations"></ChildNodeAddComponent>
      <NodeEditComponent ref="NodeEditComponent" @refresh="getOrganizations"></NodeEditComponent>
    </a-row>
  </a-card>
</template>

<script>
import { getOrganizationsApi, deleteOrganizationsApi } from '@/api/users'

import RootNodeAddComponent from './RootNodeAdd'
import ChildNodeAddComponent from './ChildNodeAdd'
import NodeEditComponent from './NodeEdit'
import NodeUsers from './NodeUsers'

export default {
  components: {
    RootNodeAddComponent,
    ChildNodeAddComponent,
    NodeEditComponent,
    NodeUsers,
  },
  data() {
    return {
      loading: false,
      treeData: [],
      expandedKeys: [],
      nodeKey: '',
      nodeUsersVisible: false,
    }
  },
  methods: {
    addRootNode() {
      this.$refs.RootNodeAddComponent.showModal()
    },
    addChildNode(item) {
      this.$refs.ChildNodeAddComponent.showModal(item)
    },
    editChildNode(item) {
      this.$refs.NodeEditComponent.showModal(item)
    },
    onSelect(selectedKeys, keys, event) {
      if (keys.selected) {
        this.nodeKey = selectedKeys[0]
        this.$nextTick(() => {
          this.$refs.NodeUsers.getOrganizationsUsers()
        })
        this.nodeUsersVisible = true
      } else {
        this.nodeUsersVisible = false
      }
    },
    getOrganizations() {
      this.loading = true
      getOrganizationsApi()
        .then((res) => {
          this.treeData = res.data
        }).catch((_error) => {})
        .finally(() => {
          this.loading = false
        })
    },
    DeleteConfirm(item) {
      const _this = this
      this.$confirm({
        title: '警告',
        content: '你确定删除？这将会删除当前节点所有的子节点！！！',
        okText: 'Yes',
        okType: 'danger',
        cancelText: 'No',
        onOk() {
          const data = {
            key: item.key,
            name: item.title,
          }
          deleteOrganizationsApi(data)
            .then((res) => {
              if (res.code === '0001') {
                _this.$message.warning(res.message)
              } else {
                _this.$message.info(res.message)
              }
            }).catch((_error) => {})
            .finally(() => {
              _this.getOrganizations()
            })
        },
        onCancel() {},
      })
    },
  },
  mounted() {
    this.getOrganizations()
  },
}
</script>

<style lang="less" scoped>
::v-deep .ant-card-body {
  padding: 8px;
}

::v-deep .ant-form {
  margin-bottom: 8px;
}

.ant-card-body .btn-type {
  float: right;
  position: absolute;
  right: 20px;
}
</style>
