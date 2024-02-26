<template>
  <a-card>
    {{ type }}
    {{ statusText }}
  </a-card>
</template>

<script>
import { downloadExportFileApi } from '@/api/orders'

export default {
  data() {
    return {
      file_name: '',
      type: 200,
      statusText: '',
    }
  },
  methods: {
    download() {
      var urlSuffix = this.$route.params.filename
      // 86b0ac87-812a-4cee-899e-9d4703367b94.xlsx.zip
      var task_id = urlSuffix.split('.')[0]
      downloadExportFileApi(task_id)
        .then((response) => {
          let blob = new Blob([response], { type: 'application/zip' })
          let url = URL.createObjectURL(blob)
          let link = document.createElement('a')
          link.href = url
          link.download = urlSuffix
          link.click()
          URL.revokeObjectURL(url)
        })
        .catch((error) => {
          console.log('error: ', error)
          const errors = [403, 404]
          if (errors.includes(error.response.status)) {
            this.type = error.response.status
            this.statusText = error.response.statusText
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