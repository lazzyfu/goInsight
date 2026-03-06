<template>
  <div class="notify-page">
    <div class="notify-page-header">
      <h2 class="notify-page-title">消息通知</h2>
      <p class="notify-page-subtitle">
        用于配置工单相关消息通知渠道。输入框离焦自动保存，启用开关会立即生效。
      </p>
    </div>

    <a-card class="notify-card" :bordered="false">
      <a-form
        ref="formRef"
        :model="formState"
        layout="vertical"
        :label-col="{ span: 24 }"
        :wrapper-col="{ span: 24 }"
        class="notify-form"
      >
        <section class="notify-section">
          <div class="notify-section-title">基础配置</div>
          <a-form-item
            label="通知地址前缀（notice_url）"
            name="notice_url"
            extra="用于拼接通知中的工单跳转地址，例如：https://goinsight.example.com/orders/{id}"
            :rules="[
              { required: true, message: '请输入通知地址前缀', trigger: ['blur', 'change'] },
              { validator: validateNoticeURL, trigger: 'blur' },
            ]"
          >
            <a-input
              v-model:value="formState.notice_url"
              placeholder="例如：https://goinsight.example.com"
              @blur="handleAutoSave('notice_url')"
            />
          </a-form-item>
        </section>

        <section class="notify-section">
          <div class="notify-section-title">通知渠道</div>
          <a-collapse v-model:activeKey="uiState.activePanels" class="notify-collapse" :bordered="false">
            <a-collapse-panel key="wechat">
              <template #header>
                <div class="notify-panel-title">
                  <span class="notify-panel-name">企业微信</span>
                  <a-tag :color="formState.wechat.enable ? 'success' : 'default'">
                    {{ formState.wechat.enable ? '已启用' : '未启用' }}
                  </a-tag>
                </div>
              </template>
              <template #extra>
                <div class="notify-panel-extra" @click.stop>
                  <a-button
                    size="small"
                    :loading="uiState.testingChannel === 'wechat'"
                    :disabled="!formState.wechat.enable || uiState.savingBySection.wechat"
                    @click="handleTestSend('wechat')"
                  >
                    测试
                  </a-button>
                  <a-switch
                    :checked="formState.wechat.enable"
                    :loading="uiState.savingBySection.wechat"
                    :disabled="uiState.savingBySection.wechat"
                    @change="(checked) => onToggleChannel('wechat', checked)"
                  />
                </div>
              </template>

              <div class="notify-channel-body">
                <p class="notify-channel-tip">用于工单审批与执行结果通知。</p>
                <a-form-item
                  label="Webhook"
                  :name="['wechat', 'webhook']"
                  :required="formState.wechat.enable"
                  :rules="[{ validator: validateWechatWebhook, trigger: 'blur' }]"
                >
                  <a-input
                    v-model:value="formState.wechat.webhook"
                    :placeholder="
                      formState.wechat.enable
                        ? '启用后此项必填，例如：https://qyapi.weixin.qq.com/...'
                        : '可先填写，启用后生效'
                    "
                    @blur="handleAutoSave('wechat')"
                  />
                </a-form-item>
              </div>
            </a-collapse-panel>

            <a-collapse-panel key="dingtalk">
              <template #header>
                <div class="notify-panel-title">
                  <span class="notify-panel-name">钉钉</span>
                  <a-tag :color="formState.dingtalk.enable ? 'success' : 'default'">
                    {{ formState.dingtalk.enable ? '已启用' : '未启用' }}
                  </a-tag>
                </div>
              </template>
              <template #extra>
                <div class="notify-panel-extra" @click.stop>
                  <a-button
                    size="small"
                    :loading="uiState.testingChannel === 'dingtalk'"
                    :disabled="!formState.dingtalk.enable || uiState.savingBySection.dingtalk"
                    @click="handleTestSend('dingtalk')"
                  >
                    测试
                  </a-button>
                  <a-switch
                    :checked="formState.dingtalk.enable"
                    :loading="uiState.savingBySection.dingtalk"
                    :disabled="uiState.savingBySection.dingtalk"
                    @change="(checked) => onToggleChannel('dingtalk', checked)"
                  />
                </div>
              </template>

              <div class="notify-channel-body">
                <p class="notify-channel-tip">适用于群机器人通知，启用前请配置关键字规则。</p>
                <a-form-item
                  label="Webhook"
                  :name="['dingtalk', 'webhook']"
                  :required="formState.dingtalk.enable"
                  :rules="[{ validator: validateDingTalkWebhook, trigger: 'blur' }]"
                >
                  <a-input
                    v-model:value="formState.dingtalk.webhook"
                    :placeholder="
                      formState.dingtalk.enable
                        ? '启用后此项必填，例如：https://oapi.dingtalk.com/...'
                        : '可先填写，启用后生效'
                    "
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
              </div>
            </a-collapse-panel>

            <a-collapse-panel key="mail">
              <template #header>
                <div class="notify-panel-title">
                  <span class="notify-panel-name">邮件</span>
                  <a-tag :color="formState.mail.enable ? 'success' : 'default'">
                    {{ formState.mail.enable ? '已启用' : '未启用' }}
                  </a-tag>
                </div>
              </template>
              <template #extra>
                <div class="notify-panel-extra" @click.stop>
                  <a-button
                    size="small"
                    :loading="uiState.testingChannel === 'mail'"
                    :disabled="!formState.mail.enable || uiState.savingBySection.mail"
                    @click="handleTestSend('mail')"
                  >
                    测试
                  </a-button>
                  <a-switch
                    :checked="formState.mail.enable"
                    :loading="uiState.savingBySection.mail"
                    :disabled="uiState.savingBySection.mail"
                    @change="(checked) => onToggleChannel('mail', checked)"
                  />
                </div>
              </template>

              <div class="notify-channel-body">
                <p class="notify-channel-tip">默认发送到平台用户邮箱，建议使用专用发件账号。</p>
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

                <a-row :gutter="[16, 0]">
                  <a-col :xs="24" :md="16">
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
                  </a-col>
                  <a-col :xs="24" :md="8">
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
                        class="notify-port-input"
                        placeholder="465"
                        @change="handleAutoSave('mail')"
                      />
                    </a-form-item>
                  </a-col>
                </a-row>

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
              </div>
            </a-collapse-panel>
          </a-collapse>
        </section>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { getNotifyConfigApi, testNotifyConfigApi, updateNotifyConfigApi } from '@/api/admin'
import { regEmail } from '@/utils/validate'
import { useThrottleFn } from '@vueuse/core'
import { message } from 'ant-design-vue'
import { computed, onMounted, reactive, ref } from 'vue'

const formRef = ref()

const normalizeText = (value) => (value || '').trim()

const isValidHTTPURL = (raw) => {
  try {
    const parsed = new URL(raw)
    return ['http:', 'https:'].includes(parsed.protocol) && !!normalizeText(parsed.host)
  } catch {
    return false
  }
}

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
  if (!isValidHTTPURL(raw)) {
    return Promise.reject(new Error('请输入合法的 URL（仅支持 http/https）'))
  }
  return Promise.resolve()
}

const validateWechatWebhook = async (_rule, value) => {
  if (!formState.wechat.enable) {
    return Promise.resolve()
  }
  const webhook = normalizeText(value)
  if (!webhook) {
    return Promise.reject(new Error('启用企业微信通知时，Webhook 不能为空'))
  }
  if (!isValidHTTPURL(webhook)) {
    return Promise.reject(new Error('企业微信 Webhook 地址格式不正确'))
  }
  return Promise.resolve()
}

const validateDingTalkWebhook = async (_rule, value) => {
  if (!formState.dingtalk.enable) {
    return Promise.resolve()
  }
  const webhook = normalizeText(value)
  if (!webhook) {
    return Promise.reject(new Error('启用钉钉通知时，Webhook 不能为空'))
  }
  if (!isValidHTTPURL(webhook)) {
    return Promise.reject(new Error('钉钉 Webhook 地址格式不正确'))
  }
  return Promise.resolve()
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
  const username = normalizeText(value)
  if (!username) {
    return Promise.reject(new Error('启用邮件通知时，发件账号不能为空'))
  }
  if (!regEmail.test(username)) {
    return Promise.reject(new Error('发件账号必须是合法邮箱'))
  }
  return Promise.resolve()
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
  if (!formState[channel].enable) {
    message.warning('请先启用该通知渠道')
    return
  }
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
.notify-page {
  max-width: 1120px;
}

.notify-page-header {
  margin-bottom: 12px;
}

.notify-page-title {
  margin: 0;
  color: var(--gi-color-text-primary);
  font-size: 20px;
  line-height: 32px;
  font-weight: 600;
}

.notify-page-subtitle {
  margin: 4px 0 0;
  color: var(--gi-color-text-secondary);
  font-size: 14px;
  line-height: 22px;
}

.notify-card {
  border-radius: 12px;
}

.notify-form {
  max-width: 1040px;
}

.notify-section + .notify-section {
  margin-top: 24px;
}

.notify-section-title {
  margin-bottom: 12px;
  color: var(--gi-color-text-primary);
  font-size: 16px;
  line-height: 24px;
  font-weight: 600;
}

.notify-collapse :deep(.ant-collapse-item) {
  border-radius: 8px;
  border: 1px solid var(--gi-color-border);
  margin-bottom: 12px;
  overflow: hidden;
}

.notify-collapse :deep(.ant-collapse) {
  border: 0;
  background: transparent;
}

.notify-collapse :deep(.ant-collapse-header) {
  align-items: center !important;
  padding: 12px 16px !important;
}

.notify-collapse :deep(.ant-collapse-content-box) {
  padding: 12px 16px 16px !important;
}

.notify-panel-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.notify-panel-name {
  color: var(--gi-color-text-primary);
  font-size: 14px;
  line-height: 22px;
  font-weight: 600;
}

.notify-panel-extra {
  display: flex;
  align-items: center;
  gap: 8px;
}

.notify-channel-body {
  padding-top: 4px;
}

.notify-channel-tip {
  margin: 0 0 12px;
  color: var(--gi-color-text-secondary);
  font-size: 12px;
  line-height: 20px;
}

.notify-port-input {
  width: 100%;
}

@media (max-width: 1023px) {
  .notify-page-title {
    font-size: 18px;
    line-height: 28px;
  }

  .notify-page-subtitle {
    font-size: 13px;
    line-height: 20px;
  }
}

@media (max-width: 767px) {
  .notify-section + .notify-section {
    margin-top: 16px;
  }

  .notify-collapse :deep(.ant-collapse-header) {
    padding: 12px !important;
  }

  .notify-collapse :deep(.ant-collapse-content-box) {
    padding: 12px !important;
  }
}
</style>
