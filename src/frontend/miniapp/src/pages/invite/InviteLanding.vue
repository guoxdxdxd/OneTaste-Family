<template>
  <div class="invite-page">
    <section class="card">
      <p class="eyebrow">一家一味 · 家庭邀请</p>
      <h1>{{ titleCopy }}</h1>
      <p class="subtitle">
        {{ subtitleCopy }}
      </p>

      <dl class="invite-info">
        <div>
          <dt>邀请人</dt>
          <dd>{{ inviteParams.inviter_nickname || '未知' }}</dd>
        </div>
        <div>
          <dt>家庭名称</dt>
          <dd>{{ inviteParams.family_name || '未提供' }}</dd>
        </div>
      </dl>

      <p v-if="!hasParams" class="error-text">
        邀请参数不完整，请联系邀请人重新扫码。
      </p>

      <div v-else-if="!userStore.loggedIn" class="login-hint">
        <p>登录后才能确认是否加入该家庭。</p>
        <button type="button" @click="goLogin">去登录</button>
      </div>

      <div v-else class="actions">
        <button type="button" class="ghost" @click="handleReject" :disabled="accepting">
          暂不加入
        </button>
        <button type="button" @click="handleAccept" :disabled="accepting">
          {{ accepting ? '处理中...' : '同意加入' }}
        </button>
      </div>

      <p v-if="feedback" class="feedback" :class="{ success }">{{ feedback }}</p>
    </section>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useFamilyStore } from '@/stores/family'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const familyStore = useFamilyStore()

const inviteParams = computed(() => ({
  family_id: route.query.family_id,
  family_name: route.query.family_name,
  inviter_id: route.query.inviter_id,
  inviter_nickname: route.query.inviter_nickname
}))

const hasParams = computed(() => {
  return Object.values(inviteParams.value).every((value) => !!value)
})

const titleCopy = computed(() => {
  if (!hasParams.value) return '邀约信息缺失'
  return `${inviteParams.value.inviter_nickname || '家人'} 邀请你加入`
})

const subtitleCopy = computed(() => {
  if (!hasParams.value) return '二维码参数可能已过期或缺失。'
  return `加入「${inviteParams.value.family_name}」后即可同步菜单、购物清单与家庭记录。`
})

const accepting = ref(false)
const feedback = ref('')
const success = ref(false)

const goLogin = () => {
  router.push({
    path: '/login',
    query: { redirect: route.fullPath }
  })
}

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
    feedback.value = '加入成功，正在为你跳转到家庭页面。'
    setTimeout(() => {
      router.push('/profile')
    }, 1000)
  } catch (error) {
    feedback.value = error.message || '加入失败，请稍后再试'
    success.value = false
  } finally {
    accepting.value = false
  }
}

const handleReject = () => {
  feedback.value = '已忽略此次邀请。'
  success.value = false
  setTimeout(() => {
    router.push('/profile')
  }, 800)
}
</script>

<style scoped>
.invite-page {
  min-height: 100vh;
  padding: 40px 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card {
  width: 100%;
  max-width: 460px;
  background: var(--color-card);
  border-radius: var(--radius-large);
  padding: 32px;
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow-card);
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.eyebrow {
  font-size: 12px;
  letter-spacing: 0.4em;
  text-transform: uppercase;
  color: var(--color-text-secondary);
}

.subtitle {
  margin: 0;
  color: var(--color-text-secondary);
  line-height: 1.5;
}

.invite-info {
  margin: 0;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
}

.invite-info dt {
  font-size: 12px;
  letter-spacing: 0.2em;
  color: var(--color-text-secondary);
}

.invite-info dd {
  margin: 4px 0 0;
  font-size: 18px;
  font-weight: 600;
}

.login-hint {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.actions {
  display: flex;
  gap: 12px;
}

.actions button,
.login-hint button {
  flex: 1;
  border: none;
  border-radius: var(--radius-medium);
  padding: 12px 18px;
  cursor: pointer;
  background: linear-gradient(120deg, var(--color-accent), var(--color-accent-soft));
  color: #fff;
}

.ghost {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  color: var(--color-text-primary);
}

.error-text {
  color: #ff6b6b;
}

.feedback {
  margin: 0;
  color: var(--color-text-secondary);
}

.feedback.success {
  color: #2cb67d;
}
</style>
