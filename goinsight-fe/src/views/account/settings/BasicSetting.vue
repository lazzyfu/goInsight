<template>
  <div class="account-settings-info-view">
    <a-row :gutter="16" type="flex" justify="center">
      <a-col :order="isMobile ? 2 : 1" :md="24" :lg="16">
        <a-form layout="vertical" :form="form">
          <a-form-item label="用户名">
            <a-input disabled v-decorator="['username', { initialValue: this.$store.getters.userInfo.username }]" />
          </a-form-item>
          <a-form-item label="昵称">
            <a-input
              v-decorator="[
                'nick_name',
                {
                  rules: [{ required: true, min: 1, max: 64 }],
                  initialValue: this.$store.getters.userInfo.nick_name,
                  validateTrigger: 'blur',
                },
              ]"
            />
          </a-form-item>
          <a-form-item label="邮箱">
            <a-input
              v-decorator="[
                'email',
                {
                  rules: [
                    {
                      type: 'email',
                      message: 'The input is not valid E-mail!',
                    },
                    {
                      required: true,
                      message: 'Please input your E-mail!',
                    },
                  ],
                  initialValue: this.$store.getters.userInfo.email,
                },
              ]"
            />
          </a-form-item>
          <a-form-item label="手机号">
            <a-input
              v-decorator="[
                'mobile',
                {
                  rules: [{ required: true, min: 3, max: 32 }],
                  initialValue: this.$store.getters.userInfo.mobile,
                  validateTrigger: 'blur',
                },
              ]"
            />
          </a-form-item>
          <a-form-item label="角色">
           <a-input
              v-decorator="[
                'role',
                {
                  initialValue: this.$store.getters.userInfo.role,
                },
              ]"
              disabled
            />
          </a-form-item>
          <a-form-item label="组织">
            <a-input
              v-decorator="[
                'organization',
                {
                  initialValue: this.$store.getters.userInfo.organization,
                },
              ]"
              disabled
            />
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="onSubmit">更新</a-button>
          </a-form-item>
        </a-form>
      </a-col>
      <a-col :order="1" :md="24" :lg="8" :style="{ minHeight: '180px' }">
        <div class="ant-upload-preview" @click="$refs.modal.edit(1)">
          <a-icon type="cloud-upload-o" class="upload-icon" />
          <div class="mask">
            <a-icon type="plus" />
          </div>
          <img :src="option.img" />
        </div>
      </a-col>
    </a-row>
    <avatar-modal ref="modal" @ok="setavatar" />
  </div>
</template>

<script>
import AvatarModal from './AvatarModal'
import { baseMixin } from '@/store/app-mixin'
import { UpdateUserInfo } from '@/api/profile'

export default {
  mixins: [baseMixin],
  components: {
    AvatarModal,
  },
  data() {
    return {
      form: this.$form.createForm(this),
      // cropper
      preview: {},
      option: {
        img: this.$store.getters.avatar,
        info: true,
        size: 1,
        outputType: 'jpeg',
        canScale: false,
        autoCrop: true,
        // 只有自动截图开启 宽度高度才生效
        autoCropWidth: 180,
        autoCropHeight: 180,
        fixedBox: true,
        // 开启宽度和高度比例
        fixed: true,
        fixedNumber: [1, 1],
      },
    }
  },
  methods: {
    setavatar(url) {
      this.option.img = url
    },
    onSubmit(e) {
      e.preventDefault()
      this.form.validateFields((err, values) => {
        if (!err) {
          values['uid'] = this.$store.getters.userInfo.uid
          UpdateUserInfo(values).then((response) => {
            this.$message.info(response.message)
          })
        }
      })
    },
  },
}
</script>

<style lang="less" scoped>
.avatar-upload-wrapper {
  height: 200px;
  width: 100%;
}

.ant-upload-preview {
  position: relative;
  margin: 0 auto;
  width: 100%;
  max-width: 180px;
  border-radius: 50%;
  box-shadow: 0 0 4px #ccc;

  .upload-icon {
    position: absolute;
    top: 0;
    right: 10px;
    font-size: 1.4rem;
    padding: 0.5rem;
    background: rgba(222, 221, 221, 0.7);
    border-radius: 50%;
    border: 1px solid rgba(0, 0, 0, 0.2);
  }
  .mask {
    opacity: 0;
    position: absolute;
    background: rgba(0, 0, 0, 0.4);
    cursor: pointer;
    transition: opacity 0.4s;

    &:hover {
      opacity: 1;
    }

    i {
      font-size: 2rem;
      position: absolute;
      top: 50%;
      left: 50%;
      margin-left: -1rem;
      margin-top: -1rem;
      color: #d6d6d6;
    }
  }

  img,
  .mask {
    width: 100%;
    max-width: 180px;
    height: 100%;
    border-radius: 50%;
    overflow: hidden;
  }
}
</style>
@/api/profile