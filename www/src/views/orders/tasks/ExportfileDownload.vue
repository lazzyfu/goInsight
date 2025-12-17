<template>
  <a-card style="width: 400px; margin: 20px auto; text-align: center;">
    <div v-if="uiData.type === 200">
      <p>下载文件：{{ uiData.file_name }}</p>
    </div>

    <div v-else>
      <a-result :status="uiData.type" :title="uiData.type" :sub-title="uiData.statusText" />
    </div>
  </a-card>
</template>

<script setup>
import { downloadExportFileApi } from '@/api/order'
import { onMounted, reactive } from 'vue'
import { useRoute } from 'vue-router'

// 路由拿到文件名
const route = useRoute()
const urlSuffix = route.params.filename || ''

// UI 状态
const uiData = reactive({
  file_name: urlSuffix,
  type: 200,       // HTTP 状态，用于 a-result
  statusText: '',  // 错误描述
})

// 获取 task_id
const getTaskId = (filename) => {
  // 去掉最后两个后缀：如 86b0ac87.xlsx.zip -> 86b0ac87
  return filename.replace(/\.[^/.]+$/, '').replace(/\.[^/.]+$/, '')
}

// 下载逻辑
const download = async () => {
  if (!urlSuffix) {
    uiData.type = 404
    uiData.statusText = '文件名不存在'
    return
  }

  const task_id = getTaskId(urlSuffix)
  console.log('task_id: ', task_id);

  try {
    const res = await downloadExportFileApi(task_id)

    // 构造 blob
    const blob = new Blob([res.data], { type: 'application/zip' })
    const url = URL.createObjectURL(blob)

    // 创建临时 a 标签触发下载
    const link = document.createElement('a')
    link.href = url
    link.download = urlSuffix
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)

    // 延迟释放 URL
    setTimeout(() => URL.revokeObjectURL(url), 1000)
  } catch (error) {
    console.log('error: ', error);
    const status = error?.status
    uiData.type = status || 500

    if (status === 403) {
      uiData.statusText = '无权限下载该文件'
    } else if (status === 404) {
      uiData.statusText = '文件不存在或已过期'
    } else {
      uiData.statusText = '下载失败，请稍后重试'
    }
  }
}

// 页面挂载时自动下载
onMounted(() => {
  download()
})
</script>

<style scoped>
a-card {
  padding: 20px;
}

p {
  font-size: 16px;
  color: #333;
}
</style>
