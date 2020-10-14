<template>
  <div class="account-settings-info-view">
    <a-row :gutter="16">
      <a-col :md="24" :lg="16">
        <el-form
          :label-position="labelPosition"
          :model="ruleForm"
          :rules="rules"
          ref="ruleForm"
          label-width="80px"
          size="small"
        >
          <el-form-item label="用户名（修改请联系系统管理员）" prop="username">
            <el-input v-model="ruleForm.username" disabled placeholder="" style="width: 95%" />
          </el-form-item>

          <el-form-item label="昵称" prop="displayname">
            <el-input v-model="ruleForm.displayname" placeholder="" style="width: 95%" />
          </el-form-item>

          <el-form-item label="Email" prop="email">
            <el-input v-model="ruleForm.email" placeholder="" style="width: 95%" />
          </el-form-item>

          <el-form-item label="Mobile" prop="mobile">
            <el-input v-model="ruleForm.mobile" placeholder="" style="width: 95%" />
          </el-form-item>
          <el-form-item style="text-align: left">
            <el-button type="primary" :disabled="pushing" @click="submitForm('ruleForm')">提交</el-button>
          </el-form-item>
        </el-form>
      </a-col>
      <a-col :md="24" :lg="8" :style="{ minHeight: '180px' }">
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
import { UpdateUserInfo } from '@/api/user'

export default {
  components: {
    AvatarModal
  },
  data() {
    return {
      pushing: false,
      // cropper
      labelPosition: 'top',
      preview: {},
      ruleForm: {
        username: this.$store.getters.userInfo.username,
        displayname: this.$store.getters.userInfo.displayname,
        email: this.$store.getters.userInfo.email,
        mobile: this.$store.getters.userInfo.mobile
      },
      rules: {
        title: [
          { required: true, message: '请输入用户名', trigger: 'blur' },
          {
            min: 2,
            max: 32,
            message: '长度在 2 到 64 个字符',
            trigger: 'blur'
          }
        ]
      },
      option: {
        img: this.$store.getters.userInfo.avatar_file,
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
        fixedNumber: [1, 1]
      }
    }
  },
  methods: {
    setavatar(url) {
      this.option.img = url
    },
    submitForm(formName) {
      this.pushing = true
      this.$refs[formName].validate(valid => {
        if (valid) {
          const commitData = {
            ...this.ruleForm
          }
          UpdateUserInfo(commitData).then(response => {
            this.$message.info(response.message)
          })
        } else {
          return false
        }
      })
      this.pushing = false
    }
  }
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
