<template>
  <div class="menu-page">
    <section class="menu-hero card">
      <div>
        <p class="eyebrow">一家一味 · 菜单中枢</p>
        <h1>{{ heroTitle }}</h1>
        <p class="subtitle">
          这里汇集每日三餐计划与灵感，规划好菜单就能自动衔接购物清单与家庭饮食节奏。
        </p>
      </div>
      <div class="today-plan">
        <p class="plan-label">今日灵感</p>
        <ul>
          <li v-for="tip in planTips" :key="tip">{{ tip }}</li>
        </ul>
      </div>
    </section>

    <section class="menu-actions">
      <article class="action-card">
        <header>
          <h3>快速创建菜单</h3>
          <p>选择日期与餐次，从家庭食谱挑选菜式，形成今日菜单。</p>
        </header>
        <div class="action-row">
          <div>
            <p class="label">需要做的事</p>
            <p class="value">早餐 · 午餐 · 晚餐</p>
          </div>
          <button type="button">即将开放</button>
        </div>
      </article>

      <article class="action-card">
        <header>
          <h3>菜单视图切换</h3>
          <p>预留每日/每周/月度视图切换区域，统一了解全家饮食节奏。</p>
        </header>
        <div class="view-pills">
          <span class="pill active">日视图</span>
          <span class="pill">周视图</span>
          <span class="pill">月视图</span>
        </div>
      </article>
    </section>

    <section class="meal-grid">
      <article
        v-for="meal in sampleMeals"
        :key="meal.label"
        class="meal-card"
      >
        <header>
          <span class="meal-label">{{ meal.label }}</span>
          <strong>{{ meal.time }}</strong>
        </header>
        <ul>
          <li v-for="dish in meal.dishes" :key="dish">
            {{ dish }}
          </li>
        </ul>
        <footer>
          <span>预计 {{ meal.duration }}</span>
          <button type="button">替换菜式</button>
        </footer>
      </article>
    </section>

    <section class="ai-tools card">
      <div class="ai-header">
        <div>
          <p class="eyebrow">AI 菜单助手</p>
          <h2>预留智能功能入口</h2>
          <p class="subtitle">
            规划 AI 生成/分析菜单的操作区，未来可直接在此触发付费能力。
          </p>
        </div>
        <div class="quota">
          <p>本周剩余调用</p>
          <strong>5 次</strong>
        </div>
      </div>
      <div class="ai-actions">
        <article
          v-for="action in aiActions"
          :key="action.title"
          class="ai-card"
        >
          <h3>{{ action.title }}</h3>
          <p>{{ action.desc }}</p>
          <button type="button">{{ action.cta }}</button>
        </article>
      </div>
    </section>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

const heroTitle = computed(() => {
  return userStore.nickname ? `${userStore.nickname} 的家庭菜单` : '温柔规划今日三餐'
})

const planTips = [
  '优先安排家人提到的想吃菜式',
  '兼顾口味与营养，预留轻断食一餐',
  '确认冰箱库存，避免重复采购'
]

const sampleMeals = [
  { label: '早餐', time: '07:30', dishes: ['牛奶燕麦粥', '煎蛋', '蜜桔一份'], duration: '20分钟' },
  { label: '午餐', time: '12:00', dishes: ['清蒸鲈鱼', '彩椒鸡丁', '番茄南瓜汤'], duration: '40分钟' },
  { label: '晚餐', time: '18:30', dishes: ['红烧牛肉', '蒜蓉生菜', '紫菜蛋花汤'], duration: '45分钟' }
]

const aiActions = [
  {
    title: 'AI 智能菜单生成',
    desc: '结合家庭食谱与身体状态，自动生成 7 日菜谱，附带推荐理由。',
    cta: '开放后提醒我'
  },
  {
    title: 'AI 菜单合理性分析',
    desc: '评估营养均衡与热量，指出需要补充或调整的餐次。',
    cta: '期待更新'
  }
]
</script>

<style scoped>
.menu-page {
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
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow-card);
  padding: 28px;
}

.menu-hero {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.eyebrow {
  font-size: 12px;
  letter-spacing: 0.4em;
  text-transform: uppercase;
  color: var(--color-text-secondary);
  margin-bottom: 8px;
}

.menu-hero h1 {
  margin: 0 0 8px;
  font-size: 30px;
}

.subtitle {
  margin: 0;
  color: var(--color-text-secondary);
  line-height: 1.6;
}

.today-plan {
  padding: 18px;
  border-radius: var(--radius-medium);
  background: var(--color-surface);
}

.plan-label {
  margin: 0 0 8px;
  font-size: 13px;
  text-transform: uppercase;
  letter-spacing: 0.3em;
}

.today-plan ul {
  margin: 0;
  padding-left: 18px;
  color: var(--color-text-primary);
  line-height: 1.5;
}

.menu-actions {
  display: grid;
  gap: 20px;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
}

.action-card {
  background: var(--color-card);
  border-radius: var(--radius-large);
  border: 1px solid var(--color-border);
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.action-card header h3 {
  margin: 0 0 6px;
}

.action-card header p {
  margin: 0;
  color: var(--color-text-secondary);
}

.action-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.action-row .label {
  margin: 0 0 4px;
  font-size: 12px;
  letter-spacing: 0.2em;
  color: var(--color-text-secondary);
}

.action-row .value {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.action-row button,
.meal-card button,
.ai-card button {
  border: none;
  background: linear-gradient(120deg, var(--color-accent), var(--color-accent-soft));
  color: #fff;
  border-radius: var(--radius-medium);
  padding: 10px 20px;
  cursor: pointer;
}

.view-pills {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.pill {
  border-radius: 999px;
  padding: 6px 14px;
  border: 1px solid var(--color-border);
  color: var(--color-text-secondary);
}

.pill.active {
  border-color: transparent;
  background: var(--color-surface);
  color: var(--color-text-primary);
}

.meal-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 16px;
}

.meal-card {
  background: var(--color-card);
  border-radius: var(--radius-medium);
  border: 1px solid var(--color-border);
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.meal-card header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.meal-label {
  font-size: 12px;
  letter-spacing: 0.3em;
  text-transform: uppercase;
  color: var(--color-text-secondary);
}

.meal-card ul {
  margin: 0;
  padding-left: 18px;
  color: var(--color-text-primary);
}

.meal-card footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
  color: var(--color-text-secondary);
}

.ai-tools {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.ai-header {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.quota {
  padding: 12px 18px;
  border-radius: var(--radius-medium);
  background: var(--color-surface);
  align-self: flex-start;
}

.quota p {
  margin: 0 0 6px;
  font-size: 13px;
  color: var(--color-text-secondary);
}

.quota strong {
  font-size: 24px;
}

.ai-actions {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
}

.ai-card {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-medium);
  padding: 18px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ai-card h3 {
  margin: 0;
  font-size: 18px;
}

.ai-card p {
  margin: 0;
  color: var(--color-text-secondary);
  line-height: 1.5;
}

@media (min-width: 720px) {
  .menu-hero {
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
  }

  .today-plan {
    width: 280px;
  }
}
</style>
