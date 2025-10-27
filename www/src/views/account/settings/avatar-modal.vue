<template>
  <a-modal title="修改头像" :open="open" width="40%" @cancel="handleCancel" :footer="null">
    <a-row>
      <a-col :xs="24" :md="12" :style="{ height: '350px' }">
        <vue-cropper ref="cropper" :img="options.img" :info="true" :autoCrop="options.autoCrop"
          :autoCropWidth="options.autoCropWidth" :autoCropHeight="options.autoCropHeight" :fixedBox="options.fixedBox"
          @realTime="realTime">
        </vue-cropper>
      </a-col>
      <a-col :xs="24" :md="12" :style="{ height: '350px' }">
        <div class="avatar-upload-preview">
          <img :src="previews.url" :style="previews.img" />
        </div>
      </a-col>
    </a-row>
    <br />
    <a-row :gutter="8" style="margin-top: 8px; flex-wrap: wrap;">
      <a-col :xs="8" :sm="4">
        <a-upload name="file" :beforeUpload="beforeUpload" :showUploadList="false">
          <a-button size="small" icon="upload">选择图片</a-button>
        </a-upload>
      </a-col>
      <a-col :xs="4" :sm="2" style="padding-right: 4px;">
        <a-button size="small" icon="plus" @click="changeScale(1)" style="width: 100%;" />
      </a-col>
      <a-col :xs="4" :sm="2" style="padding-right: 4px;">
        <a-button size="small" icon="minus" @click="changeScale(-1)" style="width: 100%;" />
      </a-col>
      <a-col :xs="4" :sm="2" style="padding-right: 4px;">
        <a-button size="small" icon="undo" @click="rotateLeft" style="width: 100%;" />
      </a-col>
      <a-col :xs="4" :sm="2" style="padding-left: 4px;">
        <a-button size="small" icon="redo" @click="rotateRight" style="width: 100%;" />
      </a-col>
      <a-col :xs="8" :sm="3" offset="4">
        <a-button size="small" type="primary" @click="finish('blob')" style="width: 100%;">保存</a-button>
      </a-col>
    </a-row>
  </a-modal>
</template>
<script setup>

import { ChangeAvatarApi } from "@/api/profile";
import { message } from 'ant-design-vue';
import { getCurrentInstance, reactive, ref } from "vue";

const { proxy } = getCurrentInstance();

const emit = defineEmits(["update:open"])
const props = defineProps({
  open: Boolean
})

const previews = ref({ url: '' })

const options = reactive({
  fileName: '',
  img: '',
  autoCrop: true,
  autoCropWidth: 200,
  autoCropHeight: 200,
  fixedBox: true,
})

// 修改图片大小 正数为变大 负数变小
const changeScale = (num) => {
  num = num || 1
  proxy.$refs.cropper.changeScale(num)
}
// 向左边旋转90度
const rotateLeft = () => {
  proxy.$refs.cropper.rotateLeft()
}
// 向右边旋转90度
const rotateRight = () => {
  proxy.$refs.cropper.rotateRight()
}
// 上传图片处理
const beforeUpload = (rawFile) => {
  if (rawFile.type.indexOf("image/") == -1) {
    message.error("请上传图片类型文件!");
    return false;
  }
  if (rawFile.size / 1024 / 1024 > 2) {
    message.error("文件大小不能超过2MB!");
    return false;
  }

  options.fileName = rawFile.name

  const reader = new FileReader()
  // 把Array Buffer转化为blob 如果是base64不需要
  // 转化为base64
  reader.readAsDataURL(rawFile)
  reader.onload = () => {
    options.img = reader.result
  }
  return false
}
// 实时预览事件
const realTime = (data) => {
  previews.value = data
}

// 上传图片（点击上传按钮）
const finish = (type) => {
  if (type === 'blob') {
    proxy.$refs.cropper.getCropBlob(async (data) => {
      const formData = new FormData()
      formData.append('file', data, options.fileName)
      ChangeAvatarApi(formData).then((res) => {
        message.success('修改成功，请刷新页面')
        proxy.$emit('ok', res.url)
        proxy.visible = false
      })
      handleCancel()
    })
  } else {
    proxy.$refs.cropper.getCropData((data) => {
      proxy.model = true
      proxy.modelSrc = data
    })
  }
}

const handleCancel = (e) => {
  emit('update:open')
}

</script>

<style lang="less" scoped>
.avatar-upload-preview {
  position: absolute;
  top: 50%;
  transform: translate(50%, -50%);
  width: 180px;
  height: 180px;
  border-radius: 50%;
  box-shadow: 0 0 4px #ccc;
  overflow: hidden;

  img {
    width: 100%;
    height: 100%;
  }
}
</style>
