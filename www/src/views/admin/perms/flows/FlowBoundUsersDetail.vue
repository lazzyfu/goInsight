<template>
  <a-modal
    :open="props.open"
    :title="`绑定用户列表 - 《${props.flowName}》`"
    :width="650"
    :footer="null"
    @cancel="handleCancel"
  >
    <a-spin :spinning="loading">
      <div class="user-list-container">
        <a-list :data-source="userList" size="small" :split="false">
          <template #header>
            <div class="list-header">
              <UserOutlined />
              共有 <strong>{{ userList.length }}</strong> 位用户绑定了此审批流
            </div>
          </template>
          <template #renderItem="{ item }">
            <a-list-item>
              <a-list-item-meta>
                <template #title>
                  <span class="user-name">
                    {{ item.nickname || item.username }}
                  </span>
                </template>
                <template #description>
                  <a-tag color="default"> 用户名: {{ item.username }} </a-tag>
                  <a-tag color="processing"> UID: {{ item.uid || item.id }} </a-tag>
                </template>
              </a-list-item-meta>
            </a-list-item>
          </template>
          <template #empty>
            <a-empty description="当前没有用户绑定此审批流" />
          </template>
        </a-list>
      </div>
    </a-spin>
  </a-modal>
</template>

<script setup>
import { getApprovalFlowUsersApi } from '@/api/admin'
import { UserOutlined } from '@ant-design/icons-vue'
import { ref, watch } from 'vue'

const emit = defineEmits(['update:open'])
const props = defineProps({
  open: Boolean,
  flowId: { type: [Number, String], default: null }, // 审批流ID
  flowName: { type: String, default: '' }, // 审批流名称
})

const loading = ref(false)
const userList = ref([])

const fetchData = async (approval_id) => {
  console.log('approval_id: ', approval_id)
  loading.value = true
  userList.value = []
  const res = await getApprovalFlowUsersApi(approval_id)
  userList.value = res.data
  loading.value = false
}

watch(
  () => props.open,
  (newVal) => {
    if (newVal && props.flowId) {
      fetchData(props.flowId)
    }
  },
)

const handleCancel = () => {
  emit('update:open', false)
  userList.value = []
}
</script>

<style scoped>
.user-list-container {
  max-height: 500px;
  overflow-y: auto;
  padding: 10px;
}
.list-header {
  padding: 8px 16px;
  background: #f0f5ff;
  border: 1px solid #bae7ff;
  border-radius: 4px;
  font-weight: 500;
  color: #1890ff;
  font-size: 14px;
}
.list-header strong {
  font-size: 16px;
  margin: 0 4px;
}
.user-name {
  font-weight: 600;
  color: #333;
}
</style>
