<template>
  <div class="auth-page">
    <div class="auth-card">
      <header class="card-header">
        <p class="eyebrow">一家一味 · 注册</p>
        <h1>创建家庭账号</h1>
        <p class="description">三步完成注册，与你家人一起记录味道与体贴的提醒。</p>
      </header>

      <form @submit.prevent="handleRegister" class="auth-form">
        <div class="form-group">
          <label for="phone">手机号</label>
          <input
            id="phone"
            v-model="form.phone"
            type="tel"
            :class="{ error: errors.phone }"
            placeholder="请输入手机号"
            maxlength="11"
            @blur="validatePhone"
            @input="clearError('phone')"
          />
          <span v-if="errors.phone" class="error-message">{{ errors.phone }}</span>
        </div>

        <div class="form-group">
          <label for="password">密码</label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            :class="{ error: errors.password }"
            placeholder="请输入密码（6-20位）"
            @blur="validatePassword"
            @input="clearError('password')"
          />
          <span v-if="errors.password" class="error-message">{{ errors.password }}</span>
        </div>

        <div class="form-group">
          <label for="nickname">昵称</label>
          <input
            id="nickname"
            v-model="form.nickname"
            type="text"
            :class="{ error: errors.nickname }"
            placeholder="请输入昵称"
            maxlength="20"
            @blur="validateNickname"
            @input="clearError('nickname')"
          />
          <span v-if="errors.nickname" class="error-message">{{ errors.nickname }}</span>
        </div>

        <div class="form-group">
          <label for="verify_code">验证码</label>
          <div class="verify-row">
            <input
              id="verify_code"
              v-model="form.verify_code"
              type="text"
              :class="{ error: errors.verify_code }"
              placeholder="请输入图形验证码"
              maxlength="4"
              autocomplete="off"
              @blur="validateVerifyCode"
              @input="handleVerifyCodeInput"
            />
            
            <div class="captcha-panel">
              <div
                class="captcha-image"
                :style="captchaStyle"
                role="button"
                tabindex="0"
                @click="refreshCaptcha"
              >
                <img v-if="captcha.image" :src="captcha.image" alt="图形验证码" />
                <span v-else class="captcha-placeholder">
                  {{ captcha.loading ? '生成中…' : '点击获取' }}
                </span>
              </div>
              <button type="button" class="ghost small" :disabled="captcha.loading" @click="refreshCaptcha">
                {{ captcha.loading ? '生成中…' : '换一张' }}
              </button>
            </div>
          </div>
          <span v-if="errors.verify_code" class="error-message">{{ errors.verify_code }}</span>
        </div>

        <div v-if="errorMessage" class="error-alert">
          {{ errorMessage }}
        </div>

        <button type="submit" :disabled="loading">
          <span v-if="loading">注册中...</span>
          <span v-else>完成注册</span>
        </button>
      </form>

      <footer class="card-footer">
        <span>已有账号？</span>
        <router-link to="/login">立即登录</router-link>
      </footer>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getCaptcha } from '@/api/user'
import {
  getPhoneError,
  getPasswordError,
  getVerifyCodeError,
  getNicknameError
} from '@/utils/validate'

const router = useRouter()
const userStore = useUserStore()

const form = reactive({
  phone: '',
  verify_code: '',
  password: '',
  nickname: ''
})

const errors = reactive({
  phone: '',
  verify_code: '',
  password: '',
  nickname: ''
})

const errorMessage = ref('')
const loading = ref(false)

const captchaSize = ref(getResponsiveCaptchaSize())
const captcha = reactive({
  key: '',
  image: '',
  loading: false
})

const captchaStyle = computed(() => ({
  width: `${captchaSize.value.width}px`,
  height: `${captchaSize.value.height}px`
}))

function getResponsiveCaptchaSize() {
  if (typeof window === 'undefined') {
    return { width: 200, height: 64 }
  }
  return window.innerWidth <= 520
    ? { width: 160, height: 54 }
    : { width: 200, height: 64 }
}

const validatePhone = () => {
  errors.phone = getPhoneError(form.phone)
  return !errors.phone
}

const validateVerifyCode = () => {
  form.verify_code = form.verify_code.trim().toUpperCase()
  errors.verify_code = getVerifyCodeError(form.verify_code)
  return !errors.verify_code
}

const validatePassword = () => {
  errors.password = getPasswordError(form.password)
  return !errors.password
}

const validateNickname = () => {
  errors.nickname = getNicknameError(form.nickname)
  return !errors.nickname
}

const clearError = (field) => {
  if (errors[field]) {
    errors[field] = ''
  }
  if (errorMessage.value) {
    errorMessage.value = ''
  }
}

const validateForm = () => {
  const phoneValid = validatePhone()
  const verifyCodeValid = validateVerifyCode()
  const passwordValid = validatePassword()
  const nicknameValid = validateNickname()
  if (!captcha.key) {
    errorMessage.value = '验证码加载失败，请点击图片刷新'
    return false
  }
  return phoneValid && verifyCodeValid && passwordValid && nicknameValid
}

const formatCaptchaImage = (image) => {
  if (!image) return ''
  return image.startsWith('data:') ? image : `data:image/png;base64,${image}`
}

const fetchCaptcha = async (silent = false, resetError = true) => {
  if (captcha.loading) return
  captcha.loading = true
  try {
    const res = await getCaptcha({
      width: captchaSize.value.width,
      height: captchaSize.value.height
    })
    const data = res.data || {}
    captcha.key = data.captcha_key || ''
    captcha.image = formatCaptchaImage(data.image_base64)
    form.verify_code = ''
    if (resetError) {
      errors.verify_code = ''
    }
  } catch (error) {
    console.error('Fetch captcha error:', error)
    captcha.key = ''
    captcha.image = ''
    form.verify_code = ''
    if (resetError) {
      errors.verify_code = ''
    }
    if (!silent) {
      errorMessage.value = error.message || '获取验证码失败，请稍后重试'
    }
  } finally {
    captcha.loading = false
  }
}

const refreshCaptcha = () => {
  fetchCaptcha()
}

const handleVerifyCodeInput = (event) => {
  const nextValue = event.target.value.replace(/[^a-zA-Z0-9]/g, '').slice(0, 4).toUpperCase()
  form.verify_code = nextValue
  clearError('verify_code')
}

const handleResize = () => {
  const nextSize = getResponsiveCaptchaSize()
  if (
    nextSize.width !== captchaSize.value.width ||
    nextSize.height !== captchaSize.value.height
  ) {
    captchaSize.value = nextSize
    fetchCaptcha(true)
  }
}

const handleRegister = async () => {
  errorMessage.value = ''

  if (!validateForm()) {
    return
  }

  loading.value = true

  try {
    await userStore.register({
      phone: form.phone,
      verify_code: form.verify_code,
      captcha_key: captcha.key,
      password: form.password,
      nickname: form.nickname
    })

    router.push('/')
  } catch (error) {
    const message = error?.message || '注册失败，请稍后重试'
    if (message.includes('验证码')) {
      errors.verify_code = message
      form.verify_code = ''
      await fetchCaptcha(true, false)
      return
    } else {
      errorMessage.value = message
      await fetchCaptcha(true)
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCaptcha()
  if (typeof window !== 'undefined') {
    window.addEventListener('resize', handleResize)
  }
})

onBeforeUnmount(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('resize', handleResize)
  }
})
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  padding: 40px 20px 60px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.auth-card {
  width: 100%;
  max-width: 460px;
  background: var(--color-card);
  border-radius: var(--radius-large);
  padding: 32px;
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow-card);
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.card-header .eyebrow {
  font-size: 12px;
  letter-spacing: 0.4em;
  text-transform: uppercase;
  margin: 0 0 12px;
  color: var(--color-text-secondary);
}

.card-header h1 {
  margin: 0 0 8px;
  font-size: 30px;
  color: var(--color-text-primary);
}

.card-header .description {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: 15px;
  line-height: 1.6;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 14px;
  color: var(--color-text-primary);
}

.form-group input {
  height: 46px;
  border-radius: var(--radius-medium);
  border: 1px solid var(--color-border);
  padding: 0 14px;
  font-size: 15px;
  transition: border 0.2s ease;
}

.form-group input:focus {
  border-color: var(--color-accent);
  outline: none;
}

.form-group input.error {
  border-color: #e17055;
}

.verify-row {
  display: flex;
  flex-direction: column;
  gap: 12px;
  align-items: stretch;
}

.verify-row input {
  width: 100%;
}

.captcha-panel {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: flex-start;
  flex-shrink: 0;
  align-self: flex-start;
  width: min(220px, 45vw);
}

.captcha-image {
  border-radius: var(--radius-medium);
  overflow: hidden;
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  max-width: 100%;
}

.captcha-image::before {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  border: 1px solid var(--color-border);
  pointer-events: none;
}

.captcha-image img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  display: block;
}

.captcha-placeholder {
  font-size: 13px;
  color: var(--color-text-secondary);
  padding: 0 14px;
  text-align: center;
}

button.ghost.small {
  height: 34px;
  padding: 0 12px;
  font-size: 13px;
}

.error-message {
  font-size: 12px;
  color: #c44536;
}

.error-alert {
  padding: 12px 14px;
  border-radius: var(--radius-small);
  background: #ffe7e1;
  color: #a3412b;
  font-size: 14px;
  border: 1px solid #ffd3c7;
}

button {
  height: 46px;
  border: none;
  border-radius: var(--radius-medium);
  background: linear-gradient(120deg, var(--color-accent), var(--color-accent-soft));
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.2s ease;
}

button.ghost {
  background: transparent;
  border: 1px solid var(--color-border);
  color: var(--color-text-secondary);
}

button.ghost:hover:not(:disabled) {
  border-color: var(--color-accent);
  color: var(--color-accent);
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.card-footer {
  display: flex;
  gap: 6px;
  font-size: 14px;
  color: var(--color-text-secondary);
  justify-content: center;
}

.card-footer a {
  font-weight: 600;
}

@media (max-width: 520px) {
  .auth-card {
    padding: 24px;
  }

  .captcha-panel {
    width: min(200px, 100%);
  }

  .captcha-image {
    width: 100% !important;
  }
}
</style>
