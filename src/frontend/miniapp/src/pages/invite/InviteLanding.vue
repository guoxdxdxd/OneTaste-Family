<template>
  <div class="invite-page">
    <!-- 装饰背景 -->
    <div class="invite-page__bg">
      <div class="invite-page__bg-circle invite-page__bg-circle--1"></div>
      <div class="invite-page__bg-circle invite-page__bg-circle--2"></div>
    </div>

    <div class="invite-page__content">
      <!-- 品牌标识 -->
      <header class="invite-header">
        <div class="invite-header__logo">
          <span>🏠</span>
        </div>
        <p class="invite-header__brand">一家一味</p>
      </header>

      <!-- 邀请卡片 -->
      <div class="invite-card card">
        <!-- 邀请信息 -->
        <div class="invite-card__main">
          <h1 class="invite-card__title">{{ titleCopy }}</h1>
          <p class="invite-card__subtitle">{{ subtitleCopy }}</p>
        </div>

        <!-- 邀请详情 -->
        <div v-if="hasParams" class="invite-info">
          <div class="invite-info__item">
            <div class="avatar avatar--lg">
              {{ inviterInitial }}
            </div>
            <div class="invite-info__detail">
              <span class="invite-info__label">邀请人</span>
              <span class="invite-info__value">{{ inviteParams.inviter_nickname }}</span>
            </div>
          </div>
          <div class="invite-info__item">
            <div class="invite-info__icon">
              <IconFamily />
            </div>
            <div class="invite-info__detail">
              <span class="invite-info__label">家庭名称</span>
              <span class="invite-info__value">{{ inviteParams.family_name }}</span>
            </div>
          </div>
        </div>

        <!-- 错误状态 -->
        <div v-if="!hasParams" class="invite-error">
          <div class="invite-error__icon">⚠️</div>
          <p class="invite-error__text">邀请参数不完整，请联系邀请人重新扫码</p>
        </div>

        <!-- 未登录提示 -->
        <div v-else-if="!userStore.loggedIn" class="invite-login">
          <p class="invite-login__text">登录后才能确认是否加入该家庭</p>
          <button class="btn btn--primary btn--lg btn--full" @click="goLogin">
            去登录
          </button>
          <p class="invite-login__hint">
            还没有账号？<router-link to="/register">立即注册</router-link>
          </p>
        </div>

        <!-- 操作按钮 -->
        <div v-else class="invite-actions">
          <button 
            class="btn btn--ghost btn--lg" 
            @click="handleReject" 
            :disabled="accepting"
          >
            暂不加入
          </button>
          <button 
            class="btn btn--primary btn--lg" 
            @click="handleAccept" 
            :disabled="accepting"
          >
            <span v-if="accepting" class="loading-spinner loading-spinner--sm"></span>
            {{ accepting ? '处理中...' : '同意加入' }}
          </button>
        </div>

        <!-- 反馈信息 -->
        <transition name="fade">
          <div 
            v-if="feedback" 
            class="invite-feedback"
            :class="{ 'invite-feedback--success': success }"
          >
            <IconCheck v-if="success" class="invite-feedback__icon" />
            <span>{{ feedback }}</span>
          </div>
        </transition>
      </div>

      <!-- 底部说明 -->
      <footer class="invite-footer">
        <p>加入家庭后可以同步菜单、购物清单与家庭记录</p>
      </footer>
    </div>
  </div>
</template>

<script setup>
/**
 * 邀请落地页
 * 
 * 功能：
 * - 展示邀请信息
 * - 确认/拒绝加入家庭
 * - 引导未登录用户登录
 */
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useFamilyStore } from '@/stores/family'
import IconFamily from '@/components/icons/IconFamily.vue'
import IconCheck from '@/components/icons/IconCheck.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const familyStore = useFamilyStore()

// 邀请参数
const inviteParams = computed(() => ({
  family_id: route.query.family_id,
  family_name: route.query.family_name,
  inviter_id: route.query.inviter_id,
  inviter_nickname: route.query.inviter_nickname
}))

// 参数是否完整
const hasParams = computed(() => {
  return Object.values(inviteParams.value).every(value => !!value)
})

// 邀请人首字母
const inviterInitial = computed(() => {
  const name = inviteParams.value.inviter_nickname || ''
  return name.charAt(0).toUpperCase() || '?'
})

// 标题文案
const titleCopy = computed(() => {
  if (!hasParams.value) return '邀请信息缺失'
  return `${inviteParams.value.inviter_nickname} 邀请你加入`
})

// 副标题文案
const subtitleCopy = computed(() => {
  if (!hasParams.value) return '二维码参数可能已过期或缺失'
  return `「${inviteParams.value.family_name}」`
})

// 状态
const accepting = ref(false)
const feedback = ref('')
const success = ref(false)

// 去登录
const goLogin = () => {
  router.push({
    path: '/login',
    query: { redirect: route.fullPath }
  })
}

// 同意加入
const handleAccept = async () => {
  if (!hasParams.value) return
  
  accepting.value = true
  feedback.value = ''
  success.value = false
  
  try {
    await familyStore.acceptInvite({
      family_id: inviteParams.value.family_id,
      family_name: inviteParams.value.family_name,
      inviter_id: inviteParams.value.inviter_id,
      inviter_nickname: inviteParams.value.inviter_nickname
    })
    
    success.value = true
    feedback.value = '加入成功！正在跳转...'
    
    setTimeout(() => {
      router.push('/profile')
    }, 1500)
  } catch (error) {
    feedback.value = error.message || '加入失败，请稍后再试'
    success.value = false
  } finally {
    accepting.value = false
  }
}

// 拒绝加入
const handleReject = () => {
  feedback.value = '已忽略此次邀请'
  success.value = false
  
  setTimeout(() => {
    router.push('/profile')
  }, 1000)
}
</script>

<style scoped>
.invite-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg-base);
  position: relative;
  overflow: hidden;
}

/* 装饰背景 */
.invite-page__bg {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: hidden;
}

.invite-page__bg-circle {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0.5;
}

.invite-page__bg-circle--1 {
  width: 300px;
  height: 300px;
  background: var(--color-primary-200);
  top: -100px;
  left: -100px;
}

.invite-page__bg-circle--2 {
  width: 250px;
  height: 250px;
  background: var(--color-secondary-200);
  bottom: 10%;
  right: -80px;
}

/* 内容区域 */
.invite-page__content {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-6) var(--space-5);
  position: relative;
  z-index: 1;
}

/* 品牌标识 */
.invite-header {
  text-align: center;
  margin-bottom: var(--space-6);
  animation: slideInDown var(--duration-slow) var(--ease-out);
}

.invite-header__logo {
  width: 64px;
  height: 64px;
  margin: 0 auto var(--space-3);
  background: var(--gradient-primary);
  border-radius: var(--radius-2xl);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  box-shadow: var(--shadow-lg);
}

.invite-header__brand {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
  margin: 0;
  letter-spacing: 0.05em;
}

/* 邀请卡片 */
.invite-card {
  width: 100%;
  max-width: 400px;
  padding: var(--space-6);
  animation: slideInUp var(--duration-slow) var(--ease-out);
  animation-delay: 100ms;
  animation-fill-mode: both;
}

.invite-card__main {
  text-align: center;
  margin-bottom: var(--space-6);
}

.invite-card__title {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-text-heading);
  margin: 0 0 var(--space-2);
}

.invite-card__subtitle {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-primary);
  margin: 0;
}

/* 邀请信息 */
.invite-info {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  padding: var(--space-5);
  background: var(--color-bg-sunken);
  border-radius: var(--radius-xl);
  margin-bottom: var(--space-6);
}

.invite-info__item {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.invite-info__icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-primary-100);
  color: var(--color-primary);
  border-radius: var(--radius-lg);
}

.invite-info__icon svg {
  width: 24px;
  height: 24px;
}

.invite-info__detail {
  flex: 1;
}

.invite-info__label {
  display: block;
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
  margin-bottom: var(--space-1);
}

.invite-info__value {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
}

/* 错误状态 */
.invite-error {
  text-align: center;
  padding: var(--space-6);
  background: var(--color-danger-50);
  border-radius: var(--radius-xl);
  margin-bottom: var(--space-5);
}

.invite-error__icon {
  font-size: 32px;
  margin-bottom: var(--space-3);
}

.invite-error__text {
  font-size: var(--font-size-sm);
  color: var(--color-danger-600);
  margin: 0;
}

/* 登录提示 */
.invite-login {
  text-align: center;
}

.invite-login__text {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0 0 var(--space-5);
}

.invite-login__hint {
  font-size: var(--font-size-sm);
  color: var(--color-text-tertiary);
  margin: var(--space-4) 0 0;
}

.invite-login__hint a {
  color: var(--color-primary);
  font-weight: var(--font-weight-semibold);
}

/* 操作按钮 */
.invite-actions {
  display: flex;
  gap: var(--space-3);
}

.invite-actions .btn {
  flex: 1;
}

/* 反馈信息 */
.invite-feedback {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  margin-top: var(--space-5);
  padding: var(--space-3) var(--space-4);
  background: var(--color-gray-100);
  border-radius: var(--radius-lg);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.invite-feedback--success {
  background: var(--color-success-100);
  color: var(--color-success-600);
}

.invite-feedback__icon {
  width: 16px;
  height: 16px;
}

/* 底部说明 */
.invite-footer {
  margin-top: var(--space-6);
  text-align: center;
  animation: fadeIn var(--duration-slow) var(--ease-out);
  animation-delay: 300ms;
  animation-fill-mode: both;
}

.invite-footer p {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
  margin: 0;
}

/* 动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity var(--duration-normal) var(--ease-out);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
