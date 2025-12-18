<template>
  <a-modal :open="open" title="任务详情" width="60%" @cancel="handleCancel" :footer="null">
    <div class="task-detail-wrapper">
      <!-- 执行信息 -->
      <a-card size="small" class="info-card">
        <template #title>
          <div class="card-title">
            <ClockCircleOutlined class="title-icon" />
            <span>执行信息</span>
          </div>
        </template>
        <a-row :gutter="[16, 16]">
          <a-col :xs="24" :sm="12" :md="8">
            <div class="stat-item">
              <a-statistic title="执行耗时" :value="formData.result.execute_cost_time || ''">
                <template #prefix>
                  <ThunderboltOutlined style="color: #1890ff" />
                </template>
              </a-statistic>
            </div>
          </a-col>
          <a-col :xs="24" :sm="12" :md="8">
            <div class="stat-item">
              <a-statistic title="备份耗时" :value="formData.result.backup_cost_time || ''">
                <template #prefix>
                  <SaveOutlined style="color: #52c41a" />
                </template>
              </a-statistic>
            </div>
          </a-col>
          <a-col :xs="24" :sm="12" :md="8">
            <div class="stat-item">
              <a-statistic title="影响行数" :value="formData.result.affected_rows || 0" :value-style="{ color: '#fa8c16' }">
                <template #prefix>
                  <DatabaseOutlined style="color: #fa8c16" />
                </template>
              </a-statistic>
            </div>
          </a-col>
        </a-row>
      </a-card>

      <!-- 导出文件信息 -->
      <a-card v-if="formData?.sql_type === 'EXPORT'" size="small" class="info-card">
        <template #title>
          <div class="card-title">
            <FileTextOutlined class="title-icon" />
            <span>导出文件信息</span>
          </div>
        </template>
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="文件名">
            <a-typography-text copyable>
              {{ formData.result.file_name || '-' }}
            </a-typography-text>
          </a-descriptions-item>
          <a-descriptions-item label="文件大小">
            <a-tag color="blue">{{ formatFileSize(formData.result.file_size) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="导出行数">
            <a-tag color="green">{{ formData.result.export_rows || 0 }} 行</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="文件加密秘钥">
            <a @click="copyRecord(formData.result.encryption_key)">
              {{ formData.result.encryption_key || '-' }}
            </a>

          </a-descriptions-item>
          <a-descriptions-item label="文件下载路径">
            <a v-if="formData.result.download_url" :href="formData.result.download_url" target="_blank"
              class="download-link">
              <DownloadOutlined /> 点击下载
            </a>
            <span v-else>-</span>
          </a-descriptions-item>
        </a-descriptions>
      </a-card>

      <!-- 错误信息 -->
      <a-card v-if="formData?.result?.error" size="small" class="info-card error-card">
        <template #title>
          <div class="card-title error-title">
            <CloseCircleOutlined class="title-icon" />
            <span>错误信息</span>
          </div>
        </template>
        <div class="code-wrapper">
          <CodeMirror :initVal="formData.result.error" />
        </div>
      </a-card>

      <!-- 执行日志 -->
      <a-card v-if="formData?.result?.execute_log" size="small" class="info-card">
        <template #title>
          <div class="card-title">
            <FileSearchOutlined class="title-icon" />
            <span>执行日志</span>
          </div>
        </template>
        <div class="code-wrapper">
          <CodeMirror :initVal="formData.result.execute_log" :height="'350px'" />
        </div>
      </a-card>

      <!-- 回滚 SQL -->
      <a-card v-if="formData?.result?.rollback_sql" size="small" class="info-card">
        <template #title>
          <div class="card-title">
            <RollbackOutlined class="title-icon" />
            <span>回滚 SQL</span>
          </div>
        </template>
        <div class="code-wrapper">
          <CodeMirror :initVal="formData.result.rollback_sql" />
        </div>
      </a-card>

      <!-- Footer -->
      <div class="modal-footer">
        <a-button type="primary" @click="handleCancel">关闭</a-button>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import {
  ClockCircleOutlined,
  ThunderboltOutlined,
  SaveOutlined,
  DatabaseOutlined,
  FileTextOutlined,
  DownloadOutlined,
  CloseCircleOutlined,
  FileSearchOutlined,
  RollbackOutlined
} from '@ant-design/icons-vue'
import CodeMirror from '@/components/edit/Codemirror.vue'
import { message } from 'ant-design-vue'
import useClipboard from 'vue-clipboard3'

const { toClipboard } = useClipboard()

const props = defineProps({
  open: Boolean,
})

const emit = defineEmits(['update:open'])

// 任务详情数据
const formData = defineModel('modelValue', {
  type: Object,
  required: true
})

// 复制：支持传入字符串或对象 { sqltext }
const copyRecord = async (value) => {
  try {
    if (!value) {
      message.warning('没有可拷贝的内容')
      return
    }
    await toClipboard(value)
    message.success('已拷贝到剪贴板')
  } catch (e) {
    message.error(String(e))
  }
}

// 关闭
const handleCancel = () => {
  emit('update:open', false)
}

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}
</script>

<style scoped>
.task-detail-wrapper {
  max-height: 70vh;
  overflow-y: auto;
  padding-right: 4px;
}

/* 自定义滚动条 */
.task-detail-wrapper::-webkit-scrollbar {
  width: 6px;
}

.task-detail-wrapper::-webkit-scrollbar-track {
  background: #f0f0f0;
  border-radius: 3px;
}

.task-detail-wrapper::-webkit-scrollbar-thumb {
  background: #bfbfbf;
  border-radius: 3px;
}

.task-detail-wrapper::-webkit-scrollbar-thumb:hover {
  background: #999;
}

/* 卡片样式 */
.info-card {
  margin-bottom: 16px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
}

.info-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.info-card:last-of-type {
  margin-bottom: 0;
}

/* 卡片标题 */
.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #262626;
}

.title-icon {
  font-size: 16px;
  color: #1890ff;
}

/* 错误卡片图标颜色 */
.error-title .title-icon {
  color: #ff4d4f;
}

/* 统计项样式 */
.stat-item {
  padding: 12px;
  background: #fafafa;
  border-radius: 6px;
  transition: all 0.3s ease;
}

.stat-item:hover {
  background: #f0f0f0;
  transform: translateY(-2px);
}

:deep(.ant-statistic-title) {
  font-size: 13px;
  color: #8c8c8c;
  margin-bottom: 8px;
}

:deep(.ant-statistic-content) {
  font-size: 20px;
  font-weight: 600;
}

/* 下载链接样式 */
.download-link {
  color: #1890ff;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  transition: all 0.3s ease;
}

.download-link:hover {
  color: #40a9ff;
  text-decoration: underline;
}

/* 代码区域包装 */
.code-wrapper {
  border-radius: 4px;
  overflow: hidden;
  border: 1px solid #f0f0f0;
}

/* 描述列表优化 */
:deep(.ant-descriptions-item-label) {
  font-weight: 500;
  background: #fafafa;
  width: 140px;
}

:deep(.ant-descriptions-item-content) {
  background: #fff;
}

/* Footer */
.modal-footer {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
  display: flex;
  justify-content: flex-end;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .task-detail-wrapper {
    max-height: 60vh;
  }

  .card-title {
    font-size: 14px;
  }

  :deep(.ant-descriptions-item-label) {
    width: 100px;
  }
}
</style>
