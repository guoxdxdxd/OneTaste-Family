<template>
  <div class="page profile-page">
    <!-- 用户信息卡片 -->
    <section class="user-card">
      <div class="user-card__main">
        <div class="avatar avatar--xl">
          {{ userInitial }}
        </div>
        <div class="user-card__info">
          <h1 class="user-card__name">{{ userStore.nickname || '未设置昵称' }}</h1>
          <p class="user-card__phone">{{ userStore.phone || '未绑定手机' }}</p>
        </div>
      </div>
      <div class="user-card__meta">
        <div class="user-card__tag">
          <span class="tag tag--primary tag--pill">{{ membershipLabel }}</span>
        </div>
      </div>
    </section>

    <!-- 家庭状态 -->
    <section v-if="!familyStore.hasFamily" class="family-create card">
      <div class="family-create__header">
        <IconFamily class="family-create__icon" />
        <div>
          <h2 class="family-create__title">创建你的家庭</h2>
          <p class="family-create__desc">创建家庭后可邀请家人一起管理菜单</p>
        </div>
      </div>

      <form @submit.prevent="handleCreateFamily" class="family-create__form">
        <div class="form-group">
          <label class="form-label">家庭名称</label>
          <input
            v-model="createForm.name"
            type="text"
            class="input"
            :class="{ 'input--error': createErrors.name }"
            placeholder="例如：温馨小家"
            maxlength="20"
          />
          <span v-if="createErrors.name" class="form-error">{{ createErrors.name }}</span>
        </div>

        <div class="form-group">
          <label class="form-label">家庭描述 <span class="form-hint">（选填）</span></label>
          <textarea
            v-model="createForm.description"
            class="textarea"
            placeholder="给家庭写一段温馨的介绍"
            maxlength="100"
            rows="2"
          />
        </div>

        <button 
          type="submit" 
          class="btn btn--primary btn--full" 
          :disabled="familyStore.createLoading"
        >
          {{ familyStore.createLoading ? '创建中...' : '创建家庭' }}
        </button>
      </form>
    </section>

    <!-- 家庭信息 -->
    <section v-else class="family-info card">
      <div class="family-info__header">
        <div>
          <p class="family-info__label">我的家庭</p>
          <h2 class="family-info__name">{{ familyStore.familyName }}</h2>
        </div>
        <span class="tag tag--success tag--pill">{{ roleLabel }}</span>
      </div>

      <p v-if="familyStore.familyInfo?.description" class="family-info__desc">
        {{ familyStore.familyInfo.description }}
      </p>

      <div class="family-info__stats">
        <div class="stat stat--sm">
          <span class="stat__value">{{ familyStore.memberCount }}</span>
          <span class="stat__label">成员</span>
        </div>
        <div class="stat stat--sm">
          <span class="stat__value">{{ dishCount }}</span>
          <span class="stat__label">菜式</span>
        </div>
        <div class="stat stat--sm">
          <span class="stat__value">{{ maxDishes }}</span>
          <span class="stat__label">上限</span>
        </div>
      </div>
    </section>

    <!-- 功能入口 -->
    <section v-if="familyStore.hasFamily" class="feature-list">
      <!-- 成员管理 -->
      <div class="feature-card card card--flat" @click="toggleMembers">
        <div class="feature-card__main">
          <div class="feature-card__icon">
            <IconFamily />
          </div>
          <div class="feature-card__content">
            <h3>家庭成员</h3>
            <p>{{ familyStore.memberCount }} 位成员</p>
          </div>
        </div>
        <IconChevronRight 
          class="feature-card__arrow" 
          :class="{ 'feature-card__arrow--open': showMembers }"
        />
      </div>

      <!-- 成员列表（展开显示） -->
      <transition name="slide">
        <div v-if="showMembers" class="members-panel">
          <div v-if="familyStore.membersLoading" class="members-loading">
            <span class="loading-spinner"></span>
            <span>加载中...</span>
          </div>
          <template v-else>
            <div 
              v-for="member in familyStore.members" 
              :key="member.user_id"
              class="member-item"
            >
              <div class="avatar avatar--sm">{{ getMemberInitial(member) }}</div>
              <div class="member-item__info">
                <span class="member-item__name">{{ member.nickname || '未命名' }}</span>
                <span class="member-item__date">{{ formatJoinedAt(member.joined_at) }} 加入</span>
              </div>
              <span 
                class="tag tag--pill" 
                :class="member.role === 'owner' ? 'tag--success' : 'tag--default'"
              >
                {{ member.role === 'owner' ? '管理员' : '成员' }}
              </span>
            </div>
          </template>
        </div>
      </transition>

      <!-- 菜谱管理 -->
      <router-link to="/recipes" class="feature-card card card--flat">
        <div class="feature-card__main">
          <div class="feature-card__icon feature-card__icon--primary">
            <IconBook />
          </div>
          <div class="feature-card__content">
            <h3>菜谱管理</h3>
            <p>管理家庭菜式库</p>
          </div>
        </div>
        <IconChevronRight class="feature-card__arrow" />
      </router-link>

      <!-- 邀请成员（仅管理员可见） -->
      <div v-if="isOwner" class="feature-card card card--flat" @click="toggleInvite">
        <div class="feature-card__main">
          <div class="feature-card__icon feature-card__icon--warning">
            <IconPlus />
          </div>
          <div class="feature-card__content">
            <h3>邀请成员</h3>
            <p>分享二维码邀请家人</p>
          </div>
        </div>
        <IconChevronRight 
          class="feature-card__arrow" 
          :class="{ 'feature-card__arrow--open': showInvite }"
        />
      </div>

      <!-- 邀请面板 -->
      <transition name="slide">
        <div v-if="showInvite && isOwner" class="invite-panel card card--flat">
          <div class="invite-panel__qr">
            <img v-if="qrDataUrl" :src="qrDataUrl" alt="邀请二维码" />
            <div v-else class="invite-panel__qr-placeholder">
              <span class="loading-spinner" v-if="qrLoading"></span>
              <span v-else>{{ qrError || '生成二维码中...' }}</span>
            </div>
          </div>
          <div class="invite-panel__link">
            <label class="form-label">邀请链接</label>
            <div class="invite-panel__link-row">
              <input 
                :value="inviteLink" 
                class="input" 
                readonly 
                @click="$event.target.select()"
              />
              <button 
                type="button" 
                class="btn btn--primary btn--sm" 
                @click="copyInviteLink"
                :disabled="!inviteLink"
              >
                {{ copyText }}
              </button>
            </div>
          </div>
          <p class="invite-panel__hint">家人扫码后需要登录才能加入</p>
        </div>
      </transition>
    </section>

    <!-- 设置区域 -->
    <section class="settings-section">
      <div class="settings-card card card--flat" @click="handleRefresh">
        <div class="settings-card__main">
          <span>刷新数据</span>
        </div>
        <span v-if="refreshing" class="loading-spinner loading-spinner--sm"></span>
      </div>

      <div class="settings-card settings-card--danger card card--flat" @click="handleLogout">
        <div class="settings-card__main">
          <span>退出登录</span>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
/**
 * 个人中心页面
 * 
 * 功能：
 * - 用户信息展示
 * - 家庭管理（创建/查看）
 * - 成员管理
 * - 邀请功能
 * - 退出登录
 */
import { computed, reactive, ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import QRCode from 'qrcode'
import { useUserStore } from '@/stores/user'
import { useFamilyStore } from '@/stores/family'
import IconFamily from '@/components/icons/IconFamily.vue'
import IconBook from '@/components/icons/IconBook.vue'
import IconPlus from '@/components/icons/IconPlus.vue'
import IconChevronRight from '@/components/icons/IconChevronRight.vue'

const router = useRouter()
const userStore = useUserStore()
const familyStore = useFamilyStore()

// 展开状态
const showMembers = ref(false)
const showInvite = ref(false)
const refreshing = ref(false)

// 创建家庭表单
const createForm = reactive({
  name: '',
  description: ''
})
const createErrors = reactive({ name: '' })

// 邀请相关
const qrDataUrl = ref('')
const qrLoading = ref(false)
const qrError = ref('')
const copyText = ref('复制')

// 用户首字母
const userInitial = computed(() => {
  const name = userStore.nickname || userStore.phone || ''
  return name.charAt(0).toUpperCase()
})

// 会员标签
const membershipLabel = computed(() => {
  const labels = {
    premium: '黄金会员',
    vip: '尊享会员',
    pro: '暖厨 PRO',
    free: '免费版'
  }
  return labels[userStore.membershipType] || '免费版'
})

// 是否为管理员
const isOwner = computed(() => {
  if (!familyStore.hasFamily) return false
  const myId = String(userStore.userId)
  const ownerId = String(familyStore.familyInfo?.owner_id || familyStore.familyInfo?.ownerId)
  return myId === ownerId
})

// 角色标签
const roleLabel = computed(() => isOwner.value ? '管理员' : '成员')

// 菜式数量
const dishCount = computed(() => familyStore.familyInfo?.dish_count || 0)
const maxDishes = computed(() => familyStore.familyInfo?.max_dishes || 30)

// 邀请链接
const inviteBase = import.meta.env.VITE_INVITE_BASE_URL || 
  (typeof window !== 'undefined' ? window.location.origin : '')

const inviteLink = computed(() => {
  if (!familyStore.hasFamily || !userStore.userId) return ''
  const familyId = familyStore.familyInfo?.family_id
  if (!familyId) return ''
  const params = new URLSearchParams({
    family_id: familyId,
    family_name: familyStore.familyName,
    inviter_id: userStore.userId,
    inviter_nickname: userStore.nickname || '家人'
  })
  return `${inviteBase}/invite?${params.toString()}`
})

// 获取成员首字母
const getMemberInitial = (member) => {
  const name = member.nickname || ''
  return name.charAt(0).toUpperCase() || '?'
}

// 格式化加入时间
const formatJoinedAt = (value) => {
  if (!value) return '未知'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return `${date.getMonth() + 1}/${date.getDate()}`
}

// 切换成员列表
const toggleMembers = () => {
  showMembers.value = !showMembers.value
  if (showMembers.value && !familyStore.members.length) {
    familyStore.fetchMembers()
  }
}

// 切换邀请面板
const toggleInvite = () => {
  showInvite.value = !showInvite.value
  if (showInvite.value && !qrDataUrl.value) {
    generateQr()
  }
}

// 生成二维码
const generateQr = async () => {
  if (!inviteLink.value) return
  qrLoading.value = true
  qrError.value = ''
  try {
    qrDataUrl.value = await QRCode.toDataURL(inviteLink.value, {
      width: 200,
      margin: 1,
      errorCorrectionLevel: 'M'
    })
  } catch (error) {
    qrError.value = '生成失败'
    console.error('QR error:', error)
  } finally {
    qrLoading.value = false
  }
}

// 复制链接
const copyInviteLink = async () => {
  if (!inviteLink.value) return
  try {
    await navigator.clipboard.writeText(inviteLink.value)
    copyText.value = '已复制'
    setTimeout(() => { copyText.value = '复制' }, 2000)
  } catch (error) {
    copyText.value = '失败'
    setTimeout(() => { copyText.value = '复制' }, 2000)
  }
}

// 创建家庭
const handleCreateFamily = async () => {
  createErrors.name = ''
  const trimmed = createForm.name.trim()
  if (!trimmed) {
    createErrors.name = '请输入家庭名称'
    return
  }
  if (trimmed.length < 2) {
    createErrors.name = '至少 2 个字'
    return
  }

  try {
    await familyStore.createFamily({
      name: trimmed,
      description: createForm.description.trim()
    })
    await familyStore.fetchMembers(true)
  } catch (error) {
    createErrors.name = error.message || '创建失败'
  }
}

// 刷新数据
const handleRefresh = async () => {
  if (refreshing.value) return
  refreshing.value = true
  try {
    await familyStore.fetchFamilyInfo(true)
    if (familyStore.hasFamily) {
      await familyStore.fetchMembers(true)
    }
  } catch (error) {
    console.error('Refresh error:', error)
  } finally {
    refreshing.value = false
  }
}

// 退出登录
const handleLogout = () => {
  if (window.confirm('确定要退出登录吗？')) {
    userStore.logout()
    router.push('/login')
  }
}

// 初始化
onMounted(async () => {
  await familyStore.fetchFamilyInfo()
  if (familyStore.hasFamily) {
    await familyStore.fetchMembers()
  }
})

// 监听邀请链接变化重新生成二维码
watch(inviteLink, () => {
  if (showInvite.value && inviteLink.value) {
    generateQr()
  }
})
</script>

<style scoped>
.profile-page {
  padding-top: var(--space-4);
}

/* 用户卡片 */
.user-card {
  background: var(--gradient-primary);
  border-radius: var(--radius-2xl);
  padding: var(--space-5);
  margin-bottom: var(--space-5);
  color: white;
}

.user-card__main {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  margin-bottom: var(--space-4);
}

.user-card .avatar {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  font-weight: var(--font-weight-bold);
}

.user-card__name {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  margin: 0 0 var(--space-1);
}

.user-card__phone {
  font-size: var(--font-size-sm);
  opacity: 0.85;
  margin: 0;
}

.user-card__tag .tag {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: none;
}

/* 创建家庭 */
.family-create {
  margin-bottom: var(--space-5);
}

.family-create__header {
  display: flex;
  align-items: flex-start;
  gap: var(--space-3);
  margin-bottom: var(--space-5);
}

.family-create__icon {
  width: 40px;
  height: 40px;
  padding: var(--space-2);
  background: var(--color-primary-100);
  color: var(--color-primary);
  border-radius: var(--radius-lg);
}

.family-create__title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
  margin: 0 0 var(--space-1);
}

.family-create__desc {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0;
}

.family-create__form {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

/* 家庭信息 */
.family-info {
  margin-bottom: var(--space-4);
}

.family-info__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--space-3);
}

.family-info__label {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  margin: 0 0 var(--space-1);
}

.family-info__name {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-text-heading);
  margin: 0;
}

.family-info__desc {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0 0 var(--space-4);
}

.family-info__stats {
  display: flex;
  gap: var(--space-6);
  padding-top: var(--space-4);
  border-top: 1px solid var(--color-border-light);
}

/* 功能列表 */
.feature-list {
  margin-bottom: var(--space-5);
}

.feature-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4);
  margin-bottom: var(--space-2);
  cursor: pointer;
  text-decoration: none;
  transition: background-color var(--transition-fast);
}

.feature-card:hover {
  background: var(--color-bg-sunken);
}

.feature-card__main {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.feature-card__icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-sunken);
  border-radius: var(--radius-lg);
  color: var(--color-text-secondary);
}

.feature-card__icon--primary {
  background: var(--color-primary-100);
  color: var(--color-primary);
}

.feature-card__icon--warning {
  background: var(--color-warning-100);
  color: var(--color-warning-600);
}

.feature-card__icon svg {
  width: 20px;
  height: 20px;
}

.feature-card__content h3 {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-heading);
  margin: 0 0 var(--space-1);
}

.feature-card__content p {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0;
}

.feature-card__arrow {
  width: 18px;
  height: 18px;
  color: var(--color-text-tertiary);
  transition: transform var(--transition-fast);
}

.feature-card__arrow--open {
  transform: rotate(90deg);
}

/* 成员面板 */
.members-panel {
  background: var(--color-bg-sunken);
  border-radius: var(--radius-lg);
  padding: var(--space-3);
  margin-bottom: var(--space-2);
}

.members-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: var(--space-4);
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}

.member-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3);
  background: var(--color-bg-elevated);
  border-radius: var(--radius-md);
  margin-bottom: var(--space-2);
}

.member-item:last-child {
  margin-bottom: 0;
}

.member-item__info {
  flex: 1;
}

.member-item__name {
  display: block;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-heading);
}

.member-item__date {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
}

/* 邀请面板 */
.invite-panel {
  padding: var(--space-4);
  margin-bottom: var(--space-2);
}

.invite-panel__qr {
  width: 160px;
  height: 160px;
  margin: 0 auto var(--space-4);
  background: white;
  border-radius: var(--radius-lg);
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.invite-panel__qr img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.invite-panel__qr-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-2);
  color: var(--color-text-tertiary);
  font-size: var(--font-size-sm);
}

.invite-panel__link {
  margin-bottom: var(--space-3);
}

.invite-panel__link-row {
  display: flex;
  gap: var(--space-2);
}

.invite-panel__link-row .input {
  flex: 1;
  font-size: var(--font-size-xs);
}

.invite-panel__hint {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
  text-align: center;
  margin: 0;
}

/* 设置区域 */
.settings-section {
  margin-top: var(--space-6);
}

.settings-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4);
  margin-bottom: var(--space-2);
  cursor: pointer;
  transition: background-color var(--transition-fast);
}

.settings-card:hover {
  background: var(--color-bg-sunken);
}

.settings-card__main {
  font-size: var(--font-size-base);
  color: var(--color-text-primary);
}

.settings-card--danger .settings-card__main {
  color: var(--color-danger-500);
}

/* 动画 */
.slide-enter-active,
.slide-leave-active {
  transition: all var(--duration-normal) var(--ease-out);
  overflow: hidden;
}

.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  max-height: 0;
  padding-top: 0;
  padding-bottom: 0;
  margin-bottom: 0;
}

.slide-enter-to,
.slide-leave-from {
  max-height: 500px;
}
</style>
