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
          <label for="verify_code">验证码</label>
          <div class="verify-row">
            <input
              id="verify_code"
              v-model="form.verify_code"
              type="text"
              :class="{ error: errors.verify_code }"
              placeholder="请输入验证码"
              maxlength="6"
              @blur="validateVerifyCode"
              @input="clearError('verify_code')"
            />
            <button type="button" class="ghost" :disabled="codeCountdown > 0 || loading" @click="sendVerifyCode">
              <span v-if="codeCountdown > 0">{{ codeCountdown }}秒</span>
              <span v-else>发送验证码</span>
            </button>
          </div>
          <span v-if="errors.verify_code" class="error-message">{{ errors.verify_code }}</span>
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
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
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
const codeCountdown = ref(0)

const validatePhone = () => {
  errors.phone = getPhoneError(form.phone)
  return !errors.phone
}

const validateVerifyCode = () => {
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
  return phoneValid && verifyCodeValid && passwordValid && nicknameValid
}

const sendVerifyCode = async () => {
  if (!validatePhone()) {
    return
  }

  try {
    codeCountdown.value = 60
    const timer = setInterval(() => {
      codeCountdown.value--
      if (codeCountdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
    console.log('验证码已发送')
  } catch (error) {
    errorMessage.value = error.message || '发送验证码失败，请稍后重试'
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
      password: form.password,
      nickname: form.nickname
    })

    router.push('/')
  } catch (error) {
    errorMessage.value = error.message || '注册失败，请稍后重试'
  } finally {
    loading.value = false
  }
}
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
  gap: 12px;
}

.verify-row button {
  flex-shrink: 0;
  padding: 0 16px;
  height: 46px;
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

  .verify-row {
    flex-direction: column;
  }

  .verify-row button {
    width: 100%;
  }
}
</style>
