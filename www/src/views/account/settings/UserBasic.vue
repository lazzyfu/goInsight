<template>
  <div class="account-settings-info-view">
    <a-row :gutter="16" type="flex" justify="center">
      <a-col :order="1" :md="24" :lg="16">
        <a-form
          layout="vertical"
          :model="formState"
          :rules="rules"
          name="basic"
          autocomplete="off"
          @finish="onFinish"
        >
          <a-form-item label="用户名" name="username">
            <a-input disabled v-model:value="formState.username" placeholder="用户名不可修改" />
          </a-form-item>

          <a-form-item
            label="昵称"
            has-feedback
            name="nick_name"
            :rules="[{ required: true, min: 1, max: 32, message: '昵称不能为空且长度为1-32位' }]"
          >
            <a-input v-model:value="formState.nick_name" placeholder="请输入昵称" />
          </a-form-item>

          <a-form-item label="邮箱" has-feedback name="email">
            <a-input v-model:value="formState.email" placeholder="请输入邮箱" />
          </a-form-item>

          <a-form-item label="手机号" has-feedback name="mobile">
            <a-input v-model:value="formState.mobile" placeholder="请输入手机号" />
          </a-form-item>

          <a-form-item label="角色" help="请联系系统管理员配置" name="role">
            <a-input v-model:value="formState.role" disabled />
          </a-form-item>

          <a-form-item label="组织" help="请联系系统管理员配置" name="organization">
            <a-input v-model:value="formState.organization" disabled />
          </a-form-item>

          <a-form-item>
            <a-button type="primary" html-type="submit" :loading="loading">更新基本信息</a-button>
          </a-form-item>
        </a-form>
      </a-col>

      <a-col :order="1" :md="24" :lg="8" :style="{ minHeight: '180px' }">
        <div class="ant-upload-preview" @click="openModal()">
          <CloudUploadOutlined class="upload-icon" />
          <div class="mask">
            <PlusOutlined />
          </div>
          <img :src="option.img" alt="头像" />
        </div>
      </a-col>
    </a-row>

    <avatar-modal :open="open" @update:open="open = $event" @ok="setavatar" />
  </div>
</template>

<script setup>
import { GetUserProfileApi } from '@/api/login' // 保持原引用
import { UpdateUserInfoApi } from '@/api/profile' // 保持原引用
import { useUserStore } from '@/store/user'
import { regEmail, regPhone } from '@/utils/validate'
import { CloudUploadOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref, watch } from 'vue'
import AvatarModal from './AvatarModal.vue'

const userStore = useUserStore()
const open = ref(false)
const loading = ref(false) // 新增 loading 状态

// 表单数据
const formState = reactive({
  username: '',
  nick_name: '',
  email: '',
  mobile: '',
  role: '',
  organization: '',
})

// 头像显示数据
const option = ref({
  img: userStore.avatar || '', // 防止初始化为 undefined
})

// 校验规则：邮箱
const validateEmail = async (_rule, value) => {
  if (!value) {
    return Promise.reject('邮箱不可为空')
  } else if (!regEmail.test(value)) {
    return Promise.reject('请输入合法的邮箱')
  }
  return Promise.resolve()
}

// 校验规则：手机号
const validateMobile = async (_rule, value) => {
  if (!value) {
    return Promise.reject('手机号不可为空')
  } else if (!regPhone.test(value)) {
    return Promise.reject('请输入合法的手机号')
  }
  return Promise.resolve()
}

const rules = {
  email: [{ required: true, validator: validateEmail, trigger: 'blur' }], // 改为 blur 触发体验更好
  mobile: [{ required: true, validator: validateMobile, trigger: 'blur' }],
}

// 核心逻辑：刷新用户信息（替代 location.reload）
const reloadUserProfile = async () => {
  try {
    const res = await GetUserProfileApi()
    if (res.code === '0000') {
      const data = res.data

      // 1. 更新 Store (严格使用原代码的 Setters)
      userStore.setUid(data.uid)
      userStore.setUserName(data.username)
      userStore.setNickName(data.nick_name)
      userStore.setUserAvatar(data.avatar_file)
      userStore.setUserEmail(data.email)
      userStore.setUserMobile(data.mobile)
      userStore.setUserOrganization(data.organization)
      userStore.setUserRole(data.role)
      userStore.setUserDateJoined(data.date_joined)

      // 2. 更新当前表单数据 (实现无刷新回显)
      formState.username = data.username
      formState.nick_name = data.nick_name
      formState.email = data.email
      formState.mobile = data.mobile
      formState.role = data.role
      formState.organization = data.organization

      // 3. 更新头像显示
      option.value.img = data.avatar_file
    }
  } catch (error) {
    console.error('获取用户信息失败', error)
  }
}

// 提交表单
const onFinish = async (values) => {
  loading.value = true
  try {
    // 确保带上 uid
    const params = {
      ...values,
      uid: userStore.uid,
    }

    const res = await UpdateUserInfoApi(params)
    if (res.code === '0000') {
      message.success('更新成功')
      // 更新成功后，重新拉取最新数据并更新 Store 和界面
      await reloadUserProfile()
    } else {
      message.error(res.message || '更新失败')
    }
  } catch (error) {
    message.error('网络请求出错')
  } finally {
    loading.value = false
  }
}

// 头像修改回调
const setavatar = (imgUrl) => {
  if (!imgUrl) return

  // 🔴 关键修复：检查是否为 Blob URL
  if (imgUrl.startsWith('blob:')) {
    option.value.img = imgUrl
    userStore.setUserAvatar(imgUrl)
  } else {
    const timestamp = new Date().getTime()
    const newUrl = imgUrl.includes('?') ? `${imgUrl}&t=${timestamp}` : `${imgUrl}?t=${timestamp}`

    option.value.img = newUrl
    userStore.setUserAvatar(newUrl)
  }
}

const openModal = () => {
  open.value = true
}

// 生命周期：挂载时主动拉取一次最新数据，确保表单有值
onMounted(() => {
  reloadUserProfile()
})

watch(
  () => userStore.username,
  (val) => {
    if (val && !formState.username) {
      formState.username = userStore.username
      formState.nick_name = userStore.nickname
      formState.email = userStore.email
      formState.mobile = userStore.mobile
      formState.role = userStore.role
      formState.organization = userStore.organization
      option.value.img = userStore.avatar
    }
  },
  { immediate: true },
)
</script>

<style lang="less" scoped>
.account-settings-info-view {
  // 增加一点内边距看起来更舒服
  padding-top: 12px;
}

.avatar-upload-wrapper {
  height: 200px;
  width: 100%;
}

.ant-upload-preview {
  position: relative;
  margin: 0 auto;
  width: 100%;
  max-width: 180px;
  height: 180px; // 显式定高，防止坍塌
  border-radius: 50%;
  box-shadow: 0 0 4px #ccc;
  cursor: pointer; // 增加手型样式

  .upload-icon {
    position: absolute;
    top: 0;
    right: 10px;
    font-size: 1.4rem;
    padding: 0.5rem;
    background: rgba(222, 221, 221, 0.7);
    border-radius: 50%;
    border: 1px solid rgba(0, 0, 0, 0.2);
    z-index: 2; // 确保在图片之上
  }

  .mask {
    opacity: 0;
    position: absolute;
    background: rgba(0, 0, 0, 0.4);
    cursor: pointer;
    transition: opacity 0.4s;
    z-index: 3;

    &:hover {
      opacity: 1;
    }

    span {
      // 修改 i 为 span 或者 ant-design 图标组件
      font-size: 2rem;
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%); // 更稳健的居中
      color: #d6d6d6;
    }
  }

  img,
  .mask {
    width: 100%;
    height: 100%;
    border-radius: 50%;
    overflow: hidden;
    object-fit: cover; // 防止图片变形
  }
}
</style>
