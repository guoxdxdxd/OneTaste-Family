<template>
  <div class="auth-page">
    <!-- 装饰背景 -->
    <div class="auth-page__bg">
      <div class="auth-page__bg-circle auth-page__bg-circle--1"></div>
      <div class="auth-page__bg-circle auth-page__bg-circle--2"></div>
    </div>

    <div class="auth-page__content">
      <!-- 品牌区域 -->
      <header class="auth-header">
        <div class="auth-header__logo">
          <span class="auth-header__logo-icon">🍳</span>
        </div>
        <h1 class="auth-header__title">一家一味</h1>
        <p class="auth-header__subtitle">记录家人口味，规划温暖三餐</p>
      </header>

      <!-- 登录表单卡片 -->
      <div class="auth-card">
        <div class="auth-card__header">
          <h2 class="auth-card__title">欢迎回来</h2>
          <p class="auth-card__desc">登录后同步家庭数据与菜单计划</p>
        </div>

        <form @submit.prevent="handleLogin" class="auth-form">
          <!-- 手机号 -->
          <div class="form-group">
            <label class="form-label" for="phone">手机号</label>
            <input
              id="phone"
              v-model="form.phone"
              type="tel"
              class="input"
              :class="{ 'input--error': errors.phone }"
              placeholder="请输入手机号"
              maxlength="11"
              autocomplete="tel"
              @blur="validatePhone"
              @input="clearError('phone')"
            />
            <span v-if="errors.phone" class="form-error">{{ errors.phone }}</span>
          </div>

          <!-- 密码 -->
          <div class="form-group">
            <label class="form-label" for="password">密码</label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              class="input"
              :class="{ 'input--error': errors.password }"
              placeholder="请输入密码"
              autocomplete="current-password"
              @blur="validatePassword"
              @input="clearError('password')"
            />
            <span v-if="errors.password" class="form-error">{{ errors.password }}</span>
          </div>

          <!-- 错误提示 -->
          <div v-if="errorMessage" class="auth-alert auth-alert--error">
            <IconClose class="auth-alert__icon" />
            <span>{{ errorMessage }}</span>
          </div>

          <!-- 提交按钮 -->
          <button 
            type="submit" 
            class="btn btn--primary btn--lg btn--full" 
            :disabled="loading"
          >
            <span v-if="loading" class="loading-spinner loading-spinner--sm"></span>
            <span>{{ loading ? '登录中...' : '登录' }}</span>
          </button>
        </form>

        <!-- 底部链接 -->
        <footer class="auth-card__footer">
          <span>还没有账号？</span>
          <router-link to="/register" class="auth-link">立即注册</router-link>
        </footer>
      </div>

      <!-- 底部说明 -->
      <div class="auth-footer">
        <p>登录即表示同意 <a href="#">服务协议</a> 和 <a href="#">隐私政策</a></p>
      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * 登录页面
 * 
 * 功能：
 * - 手机号密码登录
 * - 表单验证
 * - 登录状态管理
 */
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getPhoneError, getPasswordError } from '@/utils/validate'
import IconClose from '@/components/icons/IconClose.vue'

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
    errorMessage.value = error.message || '登录失败，请检查账号密码'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg-base);
  position: relative;
  overflow: hidden;
}

/* 装饰背景 */
.auth-page__bg {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: hidden;
}

.auth-page__bg-circle {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0.6;
}

.auth-page__bg-circle--1 {
  width: 300px;
  height: 300px;
  background: var(--color-primary-200);
  top: -100px;
  right: -100px;
}

.auth-page__bg-circle--2 {
  width: 200px;
  height: 200px;
  background: var(--color-secondary-200);
  bottom: 10%;
  left: -60px;
}

/* 内容区域 */
.auth-page__content {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-6) var(--space-5);
  position: relative;
  z-index: 1;
}

/* 品牌区域 */
.auth-header {
  text-align: center;
  margin-bottom: var(--space-8);
  animation: slideInDown var(--duration-slow) var(--ease-out);
}

.auth-header__logo {
  width: 72px;
  height: 72px;
  margin: 0 auto var(--space-4);
  background: var(--gradient-primary);
  border-radius: var(--radius-2xl);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-lg);
}

.auth-header__logo-icon {
  font-size: 36px;
}

.auth-header__title {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-text-heading);
  margin: 0 0 var(--space-2);
  letter-spacing: 0.05em;
}

.auth-header__subtitle {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0;
}

/* 表单卡片 */
.auth-card {
  width: 100%;
  max-width: 380px;
  background: var(--color-bg-elevated);
  border-radius: var(--radius-2xl);
  padding: var(--space-6);
  box-shadow: var(--shadow-xl);
  animation: slideInUp var(--duration-slow) var(--ease-out);
  animation-delay: 100ms;
  animation-fill-mode: both;
}

.auth-card__header {
  text-align: center;
  margin-bottom: var(--space-6);
}

.auth-card__title {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
  margin: 0 0 var(--space-2);
}

.auth-card__desc {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0;
}

/* 表单 */
.auth-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.auth-form .form-group {
  animation: fadeIn var(--duration-normal) var(--ease-out);
  animation-fill-mode: both;
}

.auth-form .form-group:nth-child(1) {
  animation-delay: 200ms;
}

.auth-form .form-group:nth-child(2) {
  animation-delay: 300ms;
}

/* 错误提示框 */
.auth-alert {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-4);
  border-radius: var(--radius-lg);
  font-size: var(--font-size-sm);
  animation: scaleIn var(--duration-fast) var(--ease-spring);
}

.auth-alert--error {
  background: var(--color-danger-50);
  color: var(--color-danger-600);
  border: 1px solid var(--color-danger-100);
}

.auth-alert__icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

/* 卡片底部 */
.auth-card__footer {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  margin-top: var(--space-6);
  padding-top: var(--space-5);
  border-top: 1px solid var(--color-border-light);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.auth-link {
  color: var(--color-primary);
  font-weight: var(--font-weight-semibold);
}

.auth-link:hover {
  color: var(--color-primary-dark);
}

/* 页面底部 */
.auth-footer {
  margin-top: var(--space-8);
  text-align: center;
  animation: fadeIn var(--duration-slow) var(--ease-out);
  animation-delay: 400ms;
  animation-fill-mode: both;
}

.auth-footer p {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
}

.auth-footer a {
  color: var(--color-text-secondary);
}

.auth-footer a:hover {
  color: var(--color-primary);
}

/* 按钮加载状态 */
.btn .loading-spinner {
  margin-right: var(--space-2);
}
</style>
