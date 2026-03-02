<template>
  <a-modal :open="props.open" title="用户详细信息" :width="600" @cancel="handleCancel">
    <template #footer>
      <a-button @click="handleCancel">关闭</a-button>
    </template>

    <div class="user-detail">
      <!-- 基本信息 -->
      <div class="section">
        <h3 class="section-title">基本信息</h3>
        <div class="info-grid">
          <div class="info-item">
            <span class="label">用户名：</span>
            <span class="value">{{ userInfo.username }}</span>
          </div>
          <div class="info-item">
            <span class="label">昵称：</span>
            <span class="value">{{ userInfo.nick_name }}</span>
          </div>
          <div class="info-item">
            <span class="label">邮箱：</span>
            <span class="value">{{ userInfo.email }}</span>
          </div>
          <div class="info-item">
            <span class="label">手机号：</span>
            <span class="value">{{ userInfo.mobile }}</span>
          </div>
          <div class="info-item">
            <span class="label">激活状态：</span>
            <span class="value">
              <a-tag :color="userInfo.is_active ? 'green' : 'volcano'">
                {{ userInfo.is_active ? '是' : '否' }}
              </a-tag>
            </span>
          </div>
          <div class="info-item">
            <span class="label">2FA认证：</span>
            <span class="value">
              <a-tag :color="userInfo.is_two_fa ? 'green' : 'volcano'">
                {{ userInfo.is_two_fa ? '是' : '否' }}
              </a-tag>
            </span>
          </div>
          <div class="info-item">
            <span class="label">管理员：</span>
            <span class="value">
              <a-tag :color="userInfo.is_superuser ? 'green' : 'volcano'">
                {{ userInfo.is_superuser ? '是' : '否' }}
              </a-tag>
            </span>
          </div>
          <div class="info-item">
            <span class="label">加入时间：</span>
            <span class="value">{{ userInfo.date_joined }}</span>
          </div>
        </div>
      </div>

      <!-- 组织信息 -->
      <div class="section">
        <h3 class="section-title">组织与角色</h3>
        <a-table
          :columns="tableColumns"
          :data-source="organizations"
          :pagination="false"
          size="small"
          :loading="loading"
        >
        </a-table>
        <div v-if="organizations.length === 0 && !loading" class="empty-state">
          <Empty description="该用户未加入任何组织" />
        </div>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { getUserOrganizationsApi } from '@/api/admin'
import { Empty } from 'ant-design-vue'
import { onMounted, reactive, ref, watch } from 'vue'

// 定义props和emits
const props = defineProps({
  open: Boolean,
  user: Object,
})
const emit = defineEmits(['update:open'])

// 状态
const loading = ref(false)
const userInfo = reactive({ ...props.user })
const organizations = ref([])

// 表格列
const tableColumns = [
  {
    title: '组织名称',
    dataIndex: 'organization_name',
    key: 'organization_name',
  },
  {
    title: '角色',
    dataIndex: 'role_name',
    key: 'role_name',
  },
]

// 取消按钮
const handleCancel = () => {
  emit('update:open', false)
}

// 获取用户组织信息
const fetchUserOrganizations = async () => {
  if (!props.user?.uid) return

  loading.value = true
  try {
    const res = await getUserOrganizationsApi(props.user.uid).catch(() => {})
    if (res) {
      organizations.value = res.data || []
    }
  } finally {
    loading.value = false
  }
}

// 监听用户变化
watch(
  () => props.user,
  (newUser) => {
    if (newUser) {
      Object.assign(userInfo, newUser)
      fetchUserOrganizations()
    }
  },
  { deep: true, immediate: true }
)

// 监听弹窗打开
watch(
  () => props.open,
  (isOpen) => {
    if (isOpen && props.user?.uid) {
      fetchUserOrganizations()
    }
  }
)

// 初始化
onMounted(() => {
  if (props.open && props.user?.uid) {
    fetchUserOrganizations()
  }
})
</script>

<style scoped>
.user-detail {
  padding: 16px 0;
}

.section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 12px;
  color: #333;
  border-bottom: 1px solid #f0f0f0;
  padding-bottom: 8px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.info-item {
  display: flex;
  align-items: center;
  padding: 8px 0;
}

.label {
  width: 80px;
  font-weight: 500;
  color: #666;
}

.value {
  flex: 1;
  color: #333;
}

.empty-state {
  margin-top: 24px;
  text-align: center;
}

:deep(.ant-table) {
  margin-top: 8px;
}
</style>
