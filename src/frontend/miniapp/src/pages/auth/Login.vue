<template>
  <div class="auth-page">
    <div class="auth-card">
      <header class="card-header">
        <p class="eyebrow">一家一味 · 登录</p>
        <h1>欢迎回来</h1>
        <p class="description">同步家庭口味与计划，让日常三餐少一些纠结，多一份安定。</p>
      </header>

      <ul class="comfort-list">
        <li v-for="item in comfortNotes" :key="item">{{ item }}</li>
      </ul>

      <form @submit.prevent="handleLogin" class="auth-form">
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
            placeholder="请输入密码"
            @blur="validatePassword"
            @input="clearError('password')"
          />
          <span v-if="errors.password" class="error-message">{{ errors.password }}</span>
        </div>

        <div v-if="errorMessage" class="error-alert">
          {{ errorMessage }}
        </div>

        <button type="submit" :disabled="loading">
          <span v-if="loading">登录中...</span>
          <span v-else>进入家庭空间</span>
        </button>
      </form>

      <footer class="card-footer">
        <span>还没有账号？</span>
        <router-link to="/register">创建家庭账户</router-link>
      </footer>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getPhoneError, getPasswordError } from '@/utils/validate'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// 表单数据
const form = reactive({
  phone: '',
  password: ''
})

// 错误信息
const errors = reactive({
  phone: '',
  password: ''
})

const errorMessage = ref('')
const loading = ref(false)
const comfortNotes = [
  '记录家人的口味喜好',
  '同步菜单与采购提醒',
  '守护家庭的健康节奏'
]

// 验证手机号
const validatePhone = () => {
  errors.phone = getPhoneError(form.phone)
  return !errors.phone
}

// 验证密码
const validatePassword = () => {
  errors.password = getPasswordError(form.password)
  return !errors.password
}

// 清除错误
const clearError = (field) => {
  if (errors[field]) {
    errors[field] = ''
  }
  if (errorMessage.value) {
    errorMessage.value = ''
  }
}

// 表单验证
const validateForm = () => {
  const phoneValid = validatePhone()
  const passwordValid = validatePassword()
  return phoneValid && passwordValid
}

// 处理登录
const handleLogin = async () => {
  errorMessage.value = ''

  if (!validateForm()) {
    return
  }

  loading.value = true

  try {
    await userStore.login({
      phone: form.phone,
      password: form.password
    })

    const redirect = route.query.redirect || '/'
    router.push(redirect)
  } catch (error) {
    errorMessage.value = error.message || '登录失败，请稍后重试'
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
  max-width: 420px;
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

.comfort-list {
  list-style: none;
  margin: 0;
  padding: 14px 18px;
  border-radius: var(--radius-medium);
  background: var(--color-surface);
  color: var(--color-text-secondary);
  font-size: 14px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.comfort-list li::before {
  content: '•';
  color: var(--color-accent);
  margin-right: 8px;
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

@media (max-width: 480px) {
  .auth-card {
    padding: 24px;
  }
}
</style>
