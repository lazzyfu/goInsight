<template>
  <div class="account-settings-info-view">
    <a-row :gutter="16" type="flex" justify="center">
      <a-col :order="1" :md="24" :lg="16">
        <a-form layout="vertical" :model="formState" :rules="rules" name="basic" autocomplete="off" @finish="onFinish">
          <a-form-item label="用户名" name="username" :rules="[{ required: true }]">
            <a-input disabled v-model:value="formState.username" />
          </a-form-item>
          <a-form-item label="昵称" has-feedback name="nick_name"
            :rules="[{ required: true, min: 1, max: 32, message: '昵称不能为空' }]">
            <a-input v-model:value="formState.nick_name" />
          </a-form-item>
          <a-form-item label="邮箱" has-feedback name="email">
            <a-input v-model:value="formState.email" />
          </a-form-item>
          <a-form-item label="手机号" has-feedback name="mobile">
            <a-input v-model:value="formState.mobile" />
          </a-form-item>
          <a-form-item label="角色" help="请联系系统管理员配置" name="role" :rules="[{ required: true }]">
            <a-input v-model:value="formState.role" disabled />
          </a-form-item>
          <a-form-item label="组织" help="请联系系统管理员配置" name="organization" :rules="[{ required: true }]">
            <a-input v-model:value="formState.organization" disabled />
          </a-form-item>
          <a-form-item>
            <a-button type="primary" html-type="submit">更新</a-button>
          </a-form-item>

        </a-form>
      </a-col>
      <a-col :order="1" :md="24" :lg="8" :style="{ minHeight: '180px' }">
        <div class="ant-upload-preview" @click="openModal()">
          <CloudUploadOutlined class="upload-icon" />
          <div class="mask">
            <PlusOutlined />
          </div>
          <img :src="option.img" />
        </div>
      </a-col>
    </a-row>
    <avatar-modal :open="open" @update:open="open = false" @ok="setavatar" />
  </div>
</template>


<script setup>
import { GetUserProfileApi } from "@/api/login";
import { UpdateUserInfoApi } from '@/api/profile';
import { useUserStore } from '@/store/user';
import { regEmail, regPhone } from '@/utils/validate';
import { CloudUploadOutlined, PlusOutlined } from '@ant-design/icons-vue';
import { message } from 'ant-design-vue';
import { reactive, ref } from 'vue';
import AvatarModal from './avatar-modal.vue';

const userStore = useUserStore()
const open = ref(false)

const formState = reactive({
  username: userStore.username,
  nick_name: userStore.nickname,
  email: userStore.email,
  mobile: userStore.mobile,
  role: userStore.role,
  organization: userStore.organization,
});

const validateEmail = async (_rule, value) => {
  if (!value) {
    return Promise.reject("邮箱不可为空");
  } else if (!regEmail.test(value)) {
    return Promise.reject("请输入合法的邮箱");
  }
  return Promise.resolve();
};

const validateMobile = async (_rule, value) => {
  if (!value) {
    return Promise.reject("手机号不可为空");
  } else if (!regPhone.test(value)) {
    return Promise.reject("请输入合法的手机号");
  }
  return Promise.resolve();
};

const rules = {
  email: [
    {
      required: true,
      validator: validateEmail,
      trigger: 'change',
    }
  ],
  mobile: [
    {
      required: true,
      validator: validateMobile,
      trigger: 'change',
    }
  ]
}

const reloadUserProfile = () => {
  GetUserProfileApi().then(res => {
    if (res.code === '0000') {
      userStore.setUid(res.data.uid)
      userStore.setUserName(res.data.username)
      userStore.setNickName(res.data.nick_name)
      userStore.setUserAvatar(res.data.avatar_file)
      userStore.setUserEmail(res.data.email)
      userStore.setUserMobile(res.data.mobile)
      userStore.setUserOrganization(res.data.organization)
      userStore.setUserRole(res.data.role)
      userStore.setUserDateJoined(res.data.date_joined)
      location.reload()
    }
  })
}

const onFinish = values => {
  values.uid = userStore.uid
  UpdateUserInfoApi(values).then(async res => {
    if (res.code === '0000') {
      message.info('更新成功');
      reloadUserProfile()
    } else {
      message.error(res.message);
    }
  });
};

const option = ref({
  img: userStore.avatar,
});

const setavatar = (img) => {
  option.value.img = img;
};

const openModal = () => {
  open.value = true;

};
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
