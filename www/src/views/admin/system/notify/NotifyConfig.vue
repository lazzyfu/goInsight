<template>
  <a-card title="消息通知">
    <a-form
      ref="formRef"
      :model="formState"
      layout="vertical"
      :label-col="{ span: 24 }"
      :wrapper-col="{ span: 24 }"
      class="notify-form"
    >
      <a-form-item
        label="通知地址前缀（notice_url）"
        name="notice_url"
        :rules="[
          { required: true, message: '请输入通知地址前缀', trigger: ['blur', 'change'] },
          { validator: validateNoticeURL, trigger: 'blur' }
        ]"
      >
        <a-input
          v-model:value="formState.notice_url"
          placeholder="例如：https://goinsight.example.com"
          @blur="handleAutoSave('notice_url')"
        />
      </a-form-item>

      <a-collapse v-model:activeKey="uiState.activePanels" class="notify-collapse">
        <a-collapse-panel key="wechat">
          <template #header>
            <div class="notify-panel-title">
              企业微信
              <span class="notify-panel-status">{{
                formState.wechat.enable ? '已启用' : '未启用'
              }}</span>
            </div>
          </template>
          <template #extra>
            <div class="notify-panel-extra" @click.stop>
              <a-switch
                :checked="formState.wechat.enable"
                :loading="uiState.savingBySection.wechat"
                @change="(checked) => onToggleChannel('wechat', checked)"
              />
              <a-button
                size="small"
                :loading="uiState.testingChannel === 'wechat'"
                @click="handleTestSend('wechat')"
              >
                测试
              </a-button>
            </div>
          </template>

          <a-form-item
            label="Webhook"
            :name="['wechat', 'webhook']"
            :required="formState.wechat.enable"
            :rules="[{ validator: validateWechatWebhook, trigger: 'blur' }]"
          >
            <a-input
              v-model:value="formState.wechat.webhook"
              :placeholder="formState.wechat.enable ? '启用后此项必填' : '可先填写，启用后生效'"
              @blur="handleAutoSave('wechat')"
            />
          </a-form-item>
        </a-collapse-panel>

        <a-collapse-panel key="dingtalk">
          <template #header>
            <div class="notify-panel-title">
              钉钉
              <span class="notify-panel-status">{{
                formState.dingtalk.enable ? '已启用' : '未启用'
              }}</span>
            </div>
          </template>
          <template #extra>
            <div class="notify-panel-extra" @click.stop>
              <a-switch
                :checked="formState.dingtalk.enable"
                :loading="uiState.savingBySection.dingtalk"
                @change="(checked) => onToggleChannel('dingtalk', checked)"
              />
              <a-button
                size="small"
                :loading="uiState.testingChannel === 'dingtalk'"
                @click="handleTestSend('dingtalk')"
              >
                测试
              </a-button>
            </div>
          </template>

          <a-form-item
            label="Webhook"
            :name="['dingtalk', 'webhook']"
            :required="formState.dingtalk.enable"
            :rules="[{ validator: validateDingTalkWebhook, trigger: 'blur' }]"
          >
            <a-input
              v-model:value="formState.dingtalk.webhook"
              :placeholder="formState.dingtalk.enable ? '启用后此项必填' : '可先填写，启用后生效'"
              @blur="handleAutoSave('dingtalk')"
            />
          </a-form-item>
          <a-form-item
            label="关键字"
            :name="['dingtalk', 'keywords']"
            :required="formState.dingtalk.enable"
            :rules="[{ validator: validateDingTalkKeywords, trigger: 'blur' }]"
          >
            <a-input
              v-model:value="formState.dingtalk.keywords"
              :placeholder="
                formState.dingtalk.enable ? '启用后此项必填（机器人关键字）' : '例如：GoInsight'
              "
              @blur="handleAutoSave('dingtalk')"
            />
          </a-form-item>
        </a-collapse-panel>

        <a-collapse-panel key="mail">
          <template #header>
            <div class="notify-panel-title">
              邮件
              <span class="notify-panel-status">{{
                formState.mail.enable ? '已启用' : '未启用'
              }}</span>
            </div>
          </template>
          <template #extra>
            <div class="notify-panel-extra" @click.stop>
              <a-switch
                :checked="formState.mail.enable"
                :loading="uiState.savingBySection.mail"
                @change="(checked) => onToggleChannel('mail', checked)"
              />
              <a-button
                size="small"
                :loading="uiState.testingChannel === 'mail'"
                @click="handleTestSend('mail')"
              >
                测试
              </a-button>
            </div>
          </template>

          <a-form-item
            label="发件账号"
            :name="['mail', 'username']"
            :required="formState.mail.enable"
            :rules="[{ validator: validateMailUsername, trigger: 'blur' }]"
          >
            <a-input
              v-model:value="formState.mail.username"
              :placeholder="formState.mail.enable ? '启用后此项必填' : '例如：ops@example.com'"
              @blur="handleAutoSave('mail')"
            />
          </a-form-item>

          <a-form-item
            label="SMTP 主机"
            :name="['mail', 'host']"
            :required="formState.mail.enable"
            :rules="[{ validator: validateMailHost, trigger: 'blur' }]"
          >
            <a-input
              v-model:value="formState.mail.host"
              :placeholder="formState.mail.enable ? '启用后此项必填' : '例如：smtp.163.com'"
              @blur="handleAutoSave('mail')"
            />
          </a-form-item>

          <a-form-item
            label="SMTP 端口"
            :name="['mail', 'port']"
            :required="formState.mail.enable"
            :rules="[{ validator: validateMailPort, trigger: ['blur', 'change'] }]"
          >
            <a-input-number
              v-model:value="formState.mail.port"
              :min="1"
              :max="65535"
              :precision="0"
              style="width: 240px"
              placeholder="例如：465"
              @change="handleAutoSave('mail')"
            />
          </a-form-item>

          <a-form-item
            label="SMTP 密码"
            :name="['mail', 'password']"
            :required="formState.mail.enable"
            :extra="mailPasswordExtra"
            :rules="[{ validator: validateMailPassword, trigger: ['blur', 'change'] }]"
          >
            <a-input-password
              v-model:value="formState.mail.password"
              :placeholder="mailPasswordPlaceholder"
              @blur="handleAutoSave('mail')"
            />
          </a-form-item>
        </a-collapse-panel>
      </a-collapse>
    </a-form>
  </a-card>
</template>

<script setup>
import { getNotifyConfigApi, testNotifyConfigApi, updateNotifyConfigApi } from '@/api/admin'
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { computed, onMounted, reactive, ref } from 'vue'

const formRef = ref()

const normalizeText = (value) => (value || '').trim()

const createNotifyState = () => ({
  notice_url: '',
  wechat: {
    enable: false,
    webhook: '',
  },
  dingtalk: {
    enable: false,
    webhook: '',
    keywords: '',
  },
  mail: {
    enable: false,
    username: '',
    password: '',
    host: '',
    port: 465,
    has_password: false,
  },
})

const uiState = reactive({
  testingChannel: '',
  activePanels: ['wechat', 'dingtalk', 'mail'],
  savingBySection: {
    notice_url: false,
    wechat: false,
    dingtalk: false,
    mail: false,
  },
})

const formState = reactive(createNotifyState())
const savedState = reactive(createNotifyState())

const mailPasswordExtra = computed(() => {
  if (formState.mail.has_password) {
    return '已设置 SMTP 密码，留空表示保持不变。'
  }
  return '启用邮件通知时，SMTP 密码为必填项。'
})

const mailPasswordPlaceholder = computed(() => {
  if (formState.mail.has_password && !formState.mail.password) {
    return '******'
  }
  return formState.mail.enable ? '启用后此项必填' : '请输入 SMTP 密码'
})

const validateNoticeURL = async (_rule, value) => {
  const raw = normalizeText(value)
  if (!raw) {
    return Promise.resolve()
  }
  let parsed
  try {
    parsed = new URL(raw)
  } catch {
    return Promise.reject(new Error('请输入合法的 URL'))
  }
  if (!['http:', 'https:'].includes(parsed.protocol)) {
    return Promise.reject(new Error('仅支持 http/https 协议'))
  }
  return Promise.resolve()
}

const validateWechatWebhook = async (_rule, value) => {
  if (!formState.wechat.enable) {
    return Promise.resolve()
  }
  if (normalizeText(value)) {
    return Promise.resolve()
  }
  return Promise.reject(new Error('启用企业微信通知时，Webhook 不能为空'))
}

const validateDingTalkWebhook = async (_rule, value) => {
  if (!formState.dingtalk.enable) {
    return Promise.resolve()
  }
  if (normalizeText(value)) {
    return Promise.resolve()
  }
  return Promise.reject(new Error('启用钉钉通知时，Webhook 不能为空'))
}

const validateDingTalkKeywords = async (_rule, value) => {
  if (!formState.dingtalk.enable) {
    return Promise.resolve()
  }
  if (normalizeText(value)) {
    return Promise.resolve()
  }
  return Promise.reject(new Error('启用钉钉通知时，关键字不能为空'))
}

const validateMailUsername = async (_rule, value) => {
  if (!formState.mail.enable) {
    return Promise.resolve()
  }
  if (normalizeText(value)) {
    return Promise.resolve()
  }
  return Promise.reject(new Error('启用邮件通知时，发件账号不能为空'))
}

const validateMailHost = async (_rule, value) => {
  if (!formState.mail.enable) {
    return Promise.resolve()
  }
  if (normalizeText(value)) {
    return Promise.resolve()
  }
  return Promise.reject(new Error('启用邮件通知时，SMTP 主机不能为空'))
}

const validateMailPort = async (_rule, value) => {
  if (!formState.mail.enable) {
    return Promise.resolve()
  }
  const port = Number(value)
  if (!Number.isInteger(port) || port < 1 || port > 65535) {
    return Promise.reject(new Error('启用邮件通知时，SMTP 端口范围必须为 1-65535'))
  }
  return Promise.resolve()
}

const validateMailPassword = async (_rule, value) => {
  if (!formState.mail.enable) {
    return Promise.resolve()
  }
  const input = normalizeText(value)
  if (input) {
    return Promise.resolve()
  }
  if (formState.mail.has_password) {
    return Promise.resolve()
  }
  return Promise.reject(new Error('启用邮件通知时，SMTP 密码不能为空'))
}

const buildPayloadFromState = (state) => ({
  notice_url: normalizeText(state.notice_url),
  wechat: {
    enable: !!state.wechat.enable,
    webhook: normalizeText(state.wechat.webhook),
  },
  dingtalk: {
    enable: !!state.dingtalk.enable,
    webhook: normalizeText(state.dingtalk.webhook),
    keywords: normalizeText(state.dingtalk.keywords),
  },
  mail: {
    enable: !!state.mail.enable,
    username: normalizeText(state.mail.username),
    password: normalizeText(state.mail.password),
    host: normalizeText(state.mail.host),
    port: Number(state.mail.port || 0),
    has_password: !!state.mail.has_password,
  },
})

const applyState = (target, data = {}, options = {}) => {
  target.notice_url = data.notice_url || ''

  target.wechat.enable = !!data?.wechat?.enable
  target.wechat.webhook = data?.wechat?.webhook || ''

  target.dingtalk.enable = !!data?.dingtalk?.enable
  target.dingtalk.webhook = data?.dingtalk?.webhook || ''
  target.dingtalk.keywords = data?.dingtalk?.keywords || ''

  target.mail.enable = !!data?.mail?.enable
  target.mail.username = data?.mail?.username || ''
  target.mail.password = options.resetMailPassword ? '' : data?.mail?.password || ''
  target.mail.host = data?.mail?.host || ''
  target.mail.port = data?.mail?.port || 465
  target.mail.has_password = !!data?.mail?.has_password
}

const sectionValidateFields = {
  notice_url: ['notice_url'],
  wechat: ['notice_url', ['wechat', 'webhook']],
  dingtalk: ['notice_url', ['dingtalk', 'webhook'], ['dingtalk', 'keywords']],
  mail: ['notice_url', ['mail', 'username'], ['mail', 'host'], ['mail', 'port'], ['mail', 'password']],
}

const buildPayloadFromSnapshot = (section) => {
  const payload = buildPayloadFromState(savedState)
  if (section === 'notice_url') {
    payload.notice_url = normalizeText(formState.notice_url)
  }
  if (section === 'wechat') {
    payload.wechat = buildPayloadFromState(formState).wechat
  }
  if (section === 'dingtalk') {
    payload.dingtalk = buildPayloadFromState(formState).dingtalk
  }
  if (section === 'mail') {
    payload.mail = buildPayloadFromState(formState).mail
  }
  return payload
}

const isSectionDirty = (section) => {
  if (section === 'notice_url') {
    return normalizeText(formState.notice_url) !== normalizeText(savedState.notice_url)
  }
  if (section === 'wechat') {
    return (
      !!formState.wechat.enable !== !!savedState.wechat.enable ||
      normalizeText(formState.wechat.webhook) !== normalizeText(savedState.wechat.webhook)
    )
  }
  if (section === 'dingtalk') {
    return (
      !!formState.dingtalk.enable !== !!savedState.dingtalk.enable ||
      normalizeText(formState.dingtalk.webhook) !== normalizeText(savedState.dingtalk.webhook) ||
      normalizeText(formState.dingtalk.keywords) !== normalizeText(savedState.dingtalk.keywords)
    )
  }
  if (section === 'mail') {
    return (
      !!formState.mail.enable !== !!savedState.mail.enable ||
      normalizeText(formState.mail.username) !== normalizeText(savedState.mail.username) ||
      normalizeText(formState.mail.host) !== normalizeText(savedState.mail.host) ||
      Number(formState.mail.port || 0) !== Number(savedState.mail.port || 0) ||
      normalizeText(formState.mail.password) !== ''
    )
  }
  return false
}

const fetchConfig = async () => {
  const res = await getNotifyConfigApi().catch(() => {})
  if (res?.data) {
    applyState(savedState, res.data, { resetMailPassword: true })
    applyState(formState, res.data, { resetMailPassword: true })
  }
}

const saveSection = async (section, options = {}) => {
  if (uiState.savingBySection[section]) {
    return false
  }
  try {
    if (!options.skipValidate) {
      await formRef.value?.validateFields(sectionValidateFields[section] || [])
    }
  } catch {
    return false
  }
  uiState.savingBySection[section] = true
  const res = await updateNotifyConfigApi(buildPayloadFromSnapshot(section)).catch(() => {})
  uiState.savingBySection[section] = false
  if (!res) {
    return false
  }
  if (!options.silentSuccess) {
    message.success(options.successMessage || '保存成功')
  }
  await fetchConfig()
  return true
}

const handleAutoSave = useThrottleFn(async (section) => {
  if (!isSectionDirty(section)) {
    return
  }
  await saveSection(section, { silentSuccess: true })
}, 500)

const validateChannelBeforeEnable = async (channel) => {
  try {
    await formRef.value?.validateFields(sectionValidateFields[channel] || [])
    return true
  } catch {
    return false
  }
}

const onToggleChannel = async (channel, checked) => {
  if (checked) {
    formState[channel].enable = true
    const valid = await validateChannelBeforeEnable(channel)
    if (!valid) {
      formState[channel].enable = false
      message.warning('请先填写该渠道必填项后再启用')
      return
    }
    const success = await saveSection(channel, {
      skipValidate: true,
      successMessage: '启用成功',
    })
    if (!success) {
      formState[channel].enable = false
      message.warning('启用保存失败，请稍后重试')
    }
    return
  }

  formState[channel].enable = false
  const success = await saveSection(channel, {
    skipValidate: true,
    successMessage: '已关闭',
  })
  if (!success) {
    formState[channel].enable = savedState[channel].enable
  }
}

const handleTestSend = useThrottleFn(async (channel) => {
  if (isSectionDirty('notice_url')) {
    const noticeSaved = await saveSection('notice_url', { silentSuccess: true })
    if (!noticeSaved) {
      message.warning('通知地址未保存，请检查后重试')
      return
    }
  }
  if (isSectionDirty(channel)) {
    const saved = await saveSection(channel, { silentSuccess: true })
    if (!saved) {
      message.warning('请先检查并完善当前渠道配置')
      return
    }
  }
  uiState.testingChannel = channel
  const res = await testNotifyConfigApi(channel).catch(() => {})
  if (res) {
    message.success('测试发送成功')
  }
  uiState.testingChannel = ''
})

onMounted(() => {
  fetchConfig()
})
</script>

<style scoped>
.notify-form {
  max-width: 1080px;
}

.notify-collapse {
  margin-top: 8px;
}

.notify-panel-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.notify-panel-status {
  color: #8c8c8c;
  font-size: 12px;
  font-weight: 400;
}

.notify-panel-extra {
  display: flex;
  align-items: center;
  gap: 8px;
}

.notify-collapse :deep(.ant-collapse-item) {
  border-radius: 8px;
  border: 1px solid #f0f0f0;
  margin-bottom: 12px;
  overflow: hidden;
}

.notify-collapse :deep(.ant-collapse-header) {
  align-items: center !important;
  padding: 12px 16px !important;
}

.notify-collapse :deep(.ant-collapse-content-box) {
  padding: 8px 16px 12px !important;
}
</style>
