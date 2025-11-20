<template>
  <div class="profile-page">
    <section class="card hero">
      <div>
        <p class="eyebrow">我的家庭空间</p>
        <h1>{{ greeting }}</h1>
        <p class="subtitle">
          管理家庭资料、邀请成员与权限，一次配置即可串联菜单、食谱与购物清单。
        </p>
      </div>
      <div class="hero-tags">
        <span>账号：{{ userStore.phone || '未设置' }}</span>
        <span>会员：{{ membershipCopy }}</span>
      </div>
    </section>

    <section class="card summary">
      <article>
        <p class="label">个人昵称</p>
        <h3>{{ userStore.nickname || '待完善' }}</h3>
        <p class="desc">昵称会出现在二维码邀请与成员列表中。</p>
      </article>
      <article>
        <p class="label">家庭状态</p>
        <h3>{{ familyStore.hasFamily ? '已加入' : '未创建' }}</h3>
        <p class="desc">{{ familyStore.hasFamily ? familyStore.familyName : '请先创建或接受邀请' }}</p>
      </article>
      <article>
        <p class="label">角色</p>
        <h3>{{ roleCopy }}</h3>
        <p class="desc">{{ roleDesc }}</p>
      </article>
    </section>

    <section v-if="showCreate" class="card create-card">
      <header>
        <div>
          <h2>创建家庭</h2>
          <p>每位用户仅能创建或加入一个家庭，请谨慎填写名称。</p>
        </div>
        <button type="button" class="ghost" @click="refreshInfo" :disabled="familyStore.infoLoading">
          {{ familyStore.infoLoading ? '同步中...' : '刷新状态' }}
        </button>
      </header>

      <form @submit.prevent="handleCreate" class="create-form">
        <label>
          <span>家庭名称</span>
          <input
            v-model="createForm.name"
            type="text"
            :class="{ error: createErrors.name }"
            maxlength="20"
            placeholder="例如：张家的厨房"
          />
          <small v-if="createErrors.name" class="error-text">{{ createErrors.name }}</small>
        </label>
        <label>
          <span>家庭描述（可选）</span>
          <textarea
            v-model="createForm.description"
            maxlength="100"
            placeholder="给家人一段温柔介绍"
          />
        </label>
        <p v-if="feedback" class="feedback" :class="{ error: feedbackType === 'error', success: feedbackType === 'success' }">{{ feedback }}</p>
        <button type="submit" :disabled="familyStore.createLoading">
          {{ familyStore.createLoading ? '创建中...' : '创建家庭' }}
        </button>
      </form>
    </section>

    <section v-else class="family-grid">
      <article class="card family-info">
        <header>
          <div>
            <p class="label">家庭名称</p>
            <h2>{{ familyStore.familyName }}</h2>
          </div>
          <button type="button" class="ghost" @click="refreshInfo" :disabled="familyStore.infoLoading">
            {{ familyStore.infoLoading ? '同步中...' : '刷新信息' }}
          </button>
        </header>
        <p class="desc">{{ familyStore.familyInfo?.description || '还没有简介，快去补充一段温馨的话吧。' }}</p>
        <dl class="info-stats">
          <div>
            <dt>成员数量</dt>
            <dd>{{ familyStore.memberCount }} / {{ familyStore.familyInfo?.member_limit || 10 }}</dd>
          </div>
          <div>
            <dt>菜式数量</dt>
            <dd>{{ familyStore.familyInfo?.dish_count || 0 }} / {{ familyStore.familyInfo?.max_dishes || 30 }}</dd>
          </div>
          <div>
            <dt>身份</dt>
            <dd>{{ roleCopy }}</dd>
          </div>
        </dl>
      </article>

      <article class="card members-card">
        <header>
          <div>
            <h2>家庭成员</h2>
            <p class="desc">查看角色、加入时间，Owner 可在此扩展编辑操作。</p>
          </div>
          <button type="button" class="ghost" @click="refreshMembers" :disabled="familyStore.membersLoading">
            {{ familyStore.membersLoading ? '加载中...' : '刷新列表' }}
          </button>
        </header>
        <ul v-if="familyStore.members.length" class="member-list">
          <li v-for="member in familyStore.members" :key="member.user_id">
            <div class="member-meta">
              <strong>{{ member.nickname || '未命名' }}</strong>
              <span>{{ formatJoinedAt(member.joined_at) }}</span>
            </div>
            <span class="role-pill" :class="member.role">{{ formatRole(member.role) }}</span>
          </li>
        </ul>
        <p v-else class="desc">尚无成员加入，快邀请家人一起管理菜单。</p>
      </article>

      <article class="card invite-card" :class="{ disabled: !canInvite }">
        <header>
          <div>
            <h2>扫码邀请</h2>
            <p class="desc">
              {{ canInvite ? '分享二维码或链接，家人扫码后可选择同意/拒绝加入。' : '仅家庭管理员可生成邀请。' }}
            </p>
          </div>
        </header>
        <div class="invite-link">
          <label>邀请链接</label>
          <div class="link-row">
            <input :value="inviteLink" readonly />
            <button type="button" @click="copyInviteLink" :disabled="!canInvite || copying || !inviteLink">
              {{ copying ? '复制中...' : '复制' }}
            </button>
          </div>
          <small
            v-if="copyResult"
            class="feedback"
            :class="{ error: copyState === 'error', success: copyState === 'success' }"
          >
            {{ copyResult }}
          </small>
        </div>
        <div class="qr-area">
          <div class="qr-box" v-if="qrDataUrl && canInvite">
            <img :src="qrDataUrl" alt="邀请二维码" />
          </div>
          <div class="qr-placeholder" v-else>
            <p>{{ qrPlaceholder }}</p>
          </div>
          <ul>
            <li>二维码包含家庭名称与邀请人昵称。</li>
            <li>扫码后将跳转至邀请落地页，需登录确认。</li>
            <li>拒绝邀请不会触发接口，保持轻量。</li>
          </ul>
        </div>
      </article>
    </section>

    <section class="card followup">
      <h2>下一步规划</h2>
      <p>家庭设置完成后，可继续完善菜谱、AI 服务、购物清单等功能模块。</p>
    </section>
  </div>
</template>

<script setup>
import { computed, reactive, ref, onMounted, watch } from 'vue'
import QRCode from 'qrcode'
import { useUserStore } from '@/stores/user'
import { useFamilyStore } from '@/stores/family'

const userStore = useUserStore()
const familyStore = useFamilyStore()

const greeting = computed(() => {
  return userStore.nickname ? `${userStore.nickname}，欢迎回家` : '欢迎回到家庭空间'
})

const membershipCopy = computed(() => {
  const mapping = {
    premium: '黄金会员',
    vip: '尊享会员',
    pro: '暖厨 PRO',
    free: '体验版'
  }
  return mapping[userStore.membershipType] || '体验版'
})

const createForm = reactive({
  name: '',
  description: ''
})

const createErrors = reactive({
  name: ''
})

const feedback = ref('')
const feedbackType = ref('info')
const copyResult = ref('')
const copyState = ref('info')
const copying = ref(false)
const qrDataUrl = ref('')
const qrLoading = ref(false)
const qrError = ref('')

const showCreate = computed(() => !familyStore.hasFamily)

const normalizeId = (value) => {
  if (value === undefined || value === null) return null
  return String(value)
}

const currentUserId = computed(() => {
  if (userStore.userId == null) return null
  return normalizeId(userStore.userId)
})

const membershipRole = computed(() => {
  if (!familyStore.hasFamily) return 'member'
  const myId = currentUserId.value
  const possibleOwnerIds = [familyStore.familyInfo?.owner_id, familyStore.familyInfo?.ownerId]
    .map(normalizeId)
    .filter(Boolean)
  if (myId && possibleOwnerIds.some((id) => id === myId)) {
    return 'owner'
  }
  const me = familyStore.members.find((member) => {
    const ids = [member.user_id, member.userId, member.id]
      .map(normalizeId)
      .filter(Boolean)
    return myId && ids.some((id) => id === myId)
  })
  return me?.role || familyStore.familyInfo?.member_role || familyStore.familyInfo?.role || 'member'
})

const roleCopy = computed(() => {
  return membershipRole.value === 'owner' ? '家庭管理员' : '家庭成员'
})

const roleDesc = computed(() => {
  return membershipRole.value === 'owner'
    ? '可以管理家庭信息、邀请成员与编辑权限。'
    : '可查看家庭信息、接收菜单与清单同步。'
})

const canInvite = computed(() => membershipRole.value === 'owner')

const inviteBase =
  import.meta.env.VITE_INVITE_BASE_URL ||
  (typeof window !== 'undefined' ? window.location.origin : '')

const inviteLink = computed(() => {
  if (!familyStore.hasFamily || !userStore.userId || !inviteBase) return ''
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

const qrPlaceholder = computed(() => {
  if (!familyStore.hasFamily) return '请先创建或加入家庭'
  if (!canInvite.value) return '仅家庭管理员可生成邀请二维码'
  if (qrLoading.value) return '二维码生成中...'
  if (qrError.value) return qrError.value
  return '暂无二维码可用'
})

const generateInviteQr = async () => {
  if (!inviteLink.value || !canInvite.value) {
    qrDataUrl.value = ''
    qrError.value = ''
    qrLoading.value = false
    return
  }
  qrLoading.value = true
  qrError.value = ''
  try {
    qrDataUrl.value = await QRCode.toDataURL(inviteLink.value, {
      width: 240,
      margin: 1,
      errorCorrectionLevel: 'M'
    })
  } catch (error) {
    console.error('Generate QR error:', error)
    qrError.value = '二维码生成失败，请稍后再试'
    qrDataUrl.value = ''
  } finally {
    qrLoading.value = false
  }
}

const validateCreate = () => {
  createErrors.name = ''
  const trimmed = createForm.name.trim()
  if (!trimmed) {
    createErrors.name = '请输入家庭名称'
  } else if (trimmed.length < 2) {
    createErrors.name = '家庭名称至少 2 个字'
  }
  return !createErrors.name
}

const handleCreate = async () => {
  feedback.value = ''
  feedbackType.value = 'info'

  if (!validateCreate()) return

  try {
    await familyStore.createFamily({
      name: createForm.name.trim(),
      description: createForm.description.trim()
    })
    feedback.value = '家庭创建成功，已自动加入。'
    feedbackType.value = 'success'
    await familyStore.fetchMembers(true)
  } catch (error) {
    feedback.value = error.message || '创建失败，请稍后重试'
    feedbackType.value = 'error'
  }
}

const refreshInfo = async () => {
  feedback.value = ''
  try {
    await familyStore.fetchFamilyInfo(true)
    if (familyStore.hasFamily) {
      await familyStore.fetchMembers(true)
    }
    await generateInviteQr()
  } catch (error) {
    feedback.value = error.message || '同步失败'
    feedbackType.value = 'error'
  }
}

const refreshMembers = async () => {
  try {
    await familyStore.fetchMembers(true)
    await generateInviteQr()
  } catch (error) {
    feedback.value = error.message || '成员列表刷新失败'
    feedbackType.value = 'error'
  }
}

const copyInviteLink = async () => {
  if (!inviteLink.value) return
  copying.value = true
  copyResult.value = ''
  copyState.value = 'info'
  try {
    await navigator.clipboard.writeText(inviteLink.value)
    copyResult.value = '已复制到剪贴板'
    copyState.value = 'success'
  } catch (error) {
    copyResult.value = '复制失败，请手动复制链接'
    copyState.value = 'error'
  } finally {
    copying.value = false
  }
}

const formatRole = (role) => (role === 'owner' ? '管理员' : '成员')

const formatJoinedAt = (value) => {
  if (!value) return '时间未知'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleDateString()
}

onMounted(async () => {
  await refreshInfo()
})

watch(
  () => familyStore.hasFamily,
  (has) => {
    if (has) {
      familyStore.fetchMembers()
    }
  }
)

watch(
  [inviteLink, canInvite],
  async () => {
    await generateInviteQr()
  },
  { immediate: true }
)
</script>

<style scoped>
.profile-page {
  max-width: 960px;
  margin: 0 auto;
  padding: 32px 20px 120px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.card {
  background: var(--color-card);
  border-radius: var(--radius-large);
  padding: 28px;
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow-card);
}

.hero {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.eyebrow {
  font-size: 12px;
  letter-spacing: 0.4em;
  text-transform: uppercase;
  color: var(--color-text-secondary);
  margin-bottom: 4px;
}

.subtitle {
  margin: 0;
  color: var(--color-text-secondary);
  line-height: 1.5;
}

.hero-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  color: var(--color-text-secondary);
}

.hero-tags span {
  padding: 8px 14px;
  border-radius: var(--radius-medium);
  background: var(--color-surface);
}

.summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.summary article {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-medium);
  padding: 18px;
}

.label {
  margin: 0 0 8px;
  font-size: 12px;
  letter-spacing: 0.3em;
  color: var(--color-text-secondary);
}

.summary h3 {
  margin: 0 0 6px;
}

.desc {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: 14px;
  line-height: 1.5;
}

.create-card header,
.family-info header,
.members-card header,
.invite-card header {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: flex-start;
  margin-bottom: 20px;
}

.ghost {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-medium);
  padding: 8px 14px;
  background: var(--color-surface);
  cursor: pointer;
}

.create-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.create-form label {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.create-form input,
.create-form textarea {
  border-radius: var(--radius-medium);
  border: 1px solid var(--color-border);
  padding: 12px;
  background: var(--color-surface);
}

.create-form textarea {
  min-height: 80px;
  resize: vertical;
}

.create-form input.error {
  border-color: #ff6b6b;
}

.create-form button,
.members-card button:not(.ghost),
.invite-card button:not(.ghost) {
  border: none;
  border-radius: var(--radius-medium);
  background: linear-gradient(120deg, var(--color-accent), var(--color-accent-soft));
  color: #fff;
  padding: 12px 20px;
  cursor: pointer;
}

.create-form button:disabled,
.invite-card button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error-text {
  color: #ff6b6b;
  font-size: 13px;
}

.feedback {
  margin: 0;
  font-size: 14px;
  color: var(--color-text-secondary);
}

.feedback.error {
  color: #ff6b6b;
}

.feedback.success {
  color: #2cb67d;
}

.family-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 16px;
}

.info-stats {
  margin: 20px 0 0;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 12px;
}

.info-stats dt {
  font-size: 12px;
  letter-spacing: 0.2em;
  color: var(--color-text-secondary);
}

.info-stats dd {
  margin: 4px 0 0;
  font-size: 20px;
  font-weight: 600;
}

.member-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.member-list li {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-radius: var(--radius-medium);
  background: var(--color-surface);
}

.member-meta strong {
  display: block;
}

.member-meta span {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.role-pill {
  border-radius: 999px;
  padding: 6px 14px;
  font-size: 13px;
  border: 1px solid var(--color-border);
}

.role-pill.owner {
  background: rgba(44, 182, 125, 0.1);
  border: none;
  color: #2cb67d;
}

.invite-card.disabled {
  opacity: 0.7;
}

.invite-link label {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.link-row {
  display: flex;
  gap: 8px;
  margin-top: 6px;
}

.link-row input {
  flex: 1;
  border-radius: var(--radius-medium);
  border: 1px solid var(--color-border);
  padding: 10px;
  background: var(--color-surface);
}

.qr-area {
  display: flex;
  gap: 16px;
  margin-top: 16px;
  align-items: flex-start;
  flex-wrap: wrap;
}

.qr-box {
  width: 220px;
  height: 220px;
  max-width: 100%;
  max-height: 100%;
  border-radius: var(--radius-medium);
  overflow: hidden;
  border: 1px solid var(--color-border);
  background: #fff;
  flex: 0 0 220px;
}

.qr-box img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.qr-placeholder {
  width: 220px;
  height: 220px;
  max-width: 100%;
  max-height: 100%;
  border-radius: var(--radius-medium);
  border: 1px dashed var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-secondary);
  flex: 0 0 220px;
}

.qr-area ul {
  margin: 0;
  padding-left: 18px;
  color: var(--color-text-secondary);
  line-height: 1.5;
  flex: 1 1 200px;
}

.followup h2 {
  margin: 0 0 8px;
}

@media (min-width: 720px) {
  .hero {
    flex-direction: row;
    justify-content: space-between;
  }
}
</style>
