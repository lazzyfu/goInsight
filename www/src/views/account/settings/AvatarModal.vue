<template>
  <a-modal
    title="修改头像"
    :open="props.open"
    width="800px"
    :mask-closable="false"
    :footer="null"
    @cancel="handleCancel"
  >
    <a-row :gutter="24">
      <a-col :xs="24" :md="14" :style="{ height: '350px' }">
        <vue-cropper
          ref="cropperRef"
          :img="options.img"
          :info="true"
          :autoCrop="options.autoCrop"
          :autoCropWidth="options.autoCropWidth"
          :autoCropHeight="options.autoCropHeight"
          :fixedBox="options.fixedBox"
          :outputType="options.outputType"
          @realTime="realTime"
        >
        </vue-cropper>
      </a-col>
      <a-col :xs="24" :md="10" class="preview-area">
        <div class="avatar-upload-preview">
          <div :style="previews.div" class="preview-wrapper">
            <img :src="previews.url" :style="previews.img" />
          </div>
        </div>
      </a-col>
    </a-row>
    <br />
    <a-row type="flex" justify="space-between" align="middle">
      <a-col>
        <a-upload name="file" :beforeUpload="beforeUpload" :showUploadList="false" accept="image/*">
          <a-button
            ><template #icon><UploadOutlined /></template>选择图片</a-button
          >
        </a-upload>
      </a-col>
      <a-col>
        <a-space>
          <a-button @click="changeScale(1)"
            ><template #icon><PlusOutlined /></template
          ></a-button>
          <a-button @click="changeScale(-1)"
            ><template #icon><MinusOutlined /></template
          ></a-button>
          <a-button @click="rotateLeft"
            ><template #icon><RotateLeftOutlined /></template
          ></a-button>
          <a-button @click="rotateRight"
            ><template #icon><RotateRightOutlined /></template
          ></a-button>
          <a-button type="primary" @click="onSave" :loading="loading">保存</a-button>
        </a-space>
      </a-col>
    </a-row>
  </a-modal>
</template>

<script setup>
import { ChangeAvatarApi } from '@/api/profile'
import {
    MinusOutlined,
    PlusOutlined,
    RotateLeftOutlined,
    RotateRightOutlined,
    UploadOutlined,
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { reactive, ref } from 'vue'
import { VueCropper } from 'vue-cropper'
import 'vue-cropper/dist/index.css'

const props = defineProps({ open: Boolean })
const emit = defineEmits(['update:open', 'ok'])

const cropperRef = ref(null)
const loading = ref(false)
const previews = ref({})
const options = reactive({
  img: '',
  autoCrop: true,
  autoCropWidth: 200,
  autoCropHeight: 200,
  fixedBox: true,
  outputType: 'png',
  fileName: '',
})

// ... changeScale, rotateLeft, beforeUpload 等方法保持上面的代码不变 ...
const changeScale = (num) => {
  num = num || 1
  cropperRef.value?.changeScale(num)
}
const rotateLeft = () => cropperRef.value?.rotateLeft()
const rotateRight = () => cropperRef.value?.rotateRight()
const realTime = (data) => {
  previews.value = data
}

const beforeUpload = (rawFile) => {
  if (rawFile.type.indexOf('image/') === -1) {
    message.error('请上传图片类型文件!')
    return false
  }
  if (rawFile.size / 1024 / 1024 > 5) {
    message.error('文件大小不能超过 5MB!')
    return false
  }
  options.fileName = rawFile.name
  const reader = new FileReader()
  reader.readAsDataURL(rawFile)
  reader.onload = () => {
    options.img = reader.result
  }
  return false
}

// --- 核心修复部分 ---
const onSave = () => {
  if (!options.img) return message.warning('请先选择图片')
  loading.value = true

  cropperRef.value.getCropBlob(async (blobData) => {
    if (!blobData) {
      loading.value = false
      return message.error('裁剪失败')
    }

    // 1. 生成本地 Blob URL (这能保证图片绝对能显示，不需要等后端返回 URL)
    const localBlobUrl = window.URL.createObjectURL(blobData)

    const formData = new FormData()
    formData.append('file', blobData, options.fileName || 'avatar.png')

    try {
      const res = await ChangeAvatarApi(formData)

      // 只要 code 是 0000，就认为上传成功
      if (res.code === '0000') {
        message.success('头像修改成功')

        // 2. 将本地的 Blob URL 传给父组件用于立即展示
        // 这样父组件使用的是 "blob:http://..." 地址，无需刷新，速度极快
        emit('ok', localBlobUrl)
        handleCancel()
      } else {
        message.error(res.message || '上传失败')
      }
    } catch {
      message.error('网络请求异常')
    } finally {
      loading.value = false
    }
  })
}

const handleCancel = () => {
  emit('update:open', false)
}
</script>

<style lang="less" scoped>
/* 样式保持不变，参考上一个回答 */
.preview-area {
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #fafafa;
}
.avatar-upload-preview {
  width: 200px;
  height: 200px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid #fff;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  display: flex;
  justify-content: center;
  align-items: center;
}
.preview-wrapper {
  overflow: hidden;
  border-radius: 50%;
}
</style>
