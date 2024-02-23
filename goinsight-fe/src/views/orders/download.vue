<template>
  <a-card title="下载导出文件">
    {{ this.type }}
  </a-card>
</template>

<script>
import { downloadExportFile } from '@/api/orders'

export default {
  data() {
    return {
      file_name: '',
      type: 200,
    }
  },
  methods: {
    download() {
      const urlSuffix = this.$route.params.filename
      // 86b0ac87-812a-4cee-899e-9d4703367b94.xlsx.zip
      var task_id = urlSuffix.split('.')[0]
      downloadExportFile(task_id)
        .then((response) => {
            console.log('response: ', response);
          let blob = new Blob([response], { type: 'application/zip' })
          let url = URL.createObjectURL(blob)
          let link = document.createElement('a')
          link.href = url
          link.download = this.urlSuffix
          link.click()
          URL.revokeObjectURL(url)
        })
        .catch((error) => {
          const errors = [403, 404]
          if (errors.includes(error.response.status)) {
            this.type = error.response.status
          }
        })
    },
  },
  mounted() {
    this.download()
  },
}
</script>

<style>
</style>