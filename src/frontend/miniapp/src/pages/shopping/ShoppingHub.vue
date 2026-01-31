<template>
  <div class="page shopping-page">
    <!-- 页面头部 -->
    <header class="shopping-header">
      <div class="shopping-header__info">
        <h1 class="shopping-header__title">买菜清单</h1>
        <p class="shopping-header__desc">根据菜单自动汇总食材</p>
      </div>
      <button class="icon-btn" @click="handleRefresh" :disabled="loading">
        <span v-if="loading" class="loading-spinner loading-spinner--sm"></span>
        <IconRefresh v-else />
      </button>
    </header>

    <!-- 快捷统计 -->
    <section class="shopping-stats">
      <div class="stat-card">
        <span class="stat-card__value">{{ totalItems }}</span>
        <span class="stat-card__label">待购项</span>
      </div>
      <div class="stat-card">
        <span class="stat-card__value">{{ purchasedItems }}</span>
        <span class="stat-card__label">已购买</span>
      </div>
      <div class="stat-card stat-card--accent">
        <span class="stat-card__value">{{ categoryCount }}</span>
        <span class="stat-card__label">分类</span>
      </div>
    </section>

    <!-- 生成清单入口 -->
    <section class="generate-card card card--highlight">
      <div class="generate-card__content">
        <h3>生成购物清单</h3>
        <p>选择日期范围，根据菜单自动计算食材用量</p>
      </div>
      <button class="btn btn--primary btn--sm" disabled>
        即将上线
      </button>
    </section>

    <!-- 清单内容 -->
    <section class="shopping-list">
      <div class="section-header">
        <h2 class="section-title">清单内容</h2>
        <span class="section-hint">点击可标记为已购买</span>
      </div>

      <!-- 空状态 -->
      <div v-if="!categories.length" class="empty-state">
        <div class="empty-state__icon">🛒</div>
        <h3 class="empty-state__title">暂无购物清单</h3>
        <p class="empty-state__description">创建菜单后可自动生成购物清单</p>
      </div>

      <!-- 分类列表 -->
      <div v-else class="category-list">
        <article 
          v-for="category in categories" 
          :key="category.name"
          class="category-section"
        >
          <div class="category-section__header">
            <span class="category-section__icon">{{ category.icon }}</span>
            <h3 class="category-section__title">{{ category.name }}</h3>
            <span class="category-section__count">{{ category.items.length }} 项</span>
          </div>

          <div class="item-list">
            <div 
              v-for="item in category.items" 
              :key="item.name"
              class="shopping-item"
              :class="{ 'shopping-item--done': item.purchased }"
              @click="toggleItem(item)"
            >
              <div class="shopping-item__checkbox">
                <IconCheck v-if="item.purchased" />
              </div>
              <div class="shopping-item__info">
                <span class="shopping-item__name">{{ item.name }}</span>
                <span class="shopping-item__quantity">{{ item.quantity }}</span>
              </div>
              <span v-if="item.storage" class="shopping-item__storage">
                {{ item.storage }}
              </span>
            </div>
          </div>
        </article>
      </div>
    </section>

    <!-- 底部操作 -->
    <section v-if="categories.length" class="shopping-actions">
      <button class="btn btn--ghost btn--full" @click="clearPurchased">
        清除已购买项
      </button>
    </section>
  </div>
</template>

<script setup>
/**
 * 购物清单页面
 * 
 * 功能：
 * - 展示购物清单
 * - 按分类显示食材
 * - 标记已购买
 */
import { ref, computed } from 'vue'
import IconCheck from '@/components/icons/IconCheck.vue'

// 添加刷新图标组件（简单实现）
const IconRefresh = {
  template: `
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2"/>
    </svg>
  `
}

const loading = ref(false)

// 模拟数据（实际应从接口获取）
const categories = ref([
  {
    name: '蔬菜',
    icon: '🥬',
    items: [
      { name: '西兰花', quantity: '2 颗', storage: '2天内', purchased: false },
      { name: '生菜', quantity: '1 颗', storage: '3天内', purchased: false },
      { name: '番茄', quantity: '4 个', storage: '5天内', purchased: true },
      { name: '青椒', quantity: '3 个', storage: '4天内', purchased: false }
    ]
  },
  {
    name: '肉类',
    icon: '🥩',
    items: [
      { name: '鸡胸肉', quantity: '500g', storage: '冷冻3天', purchased: false },
      { name: '五花肉', quantity: '400g', storage: '当天用', purchased: false }
    ]
  },
  {
    name: '调料',
    icon: '🧂',
    items: [
      { name: '生抽', quantity: '1 瓶', storage: null, purchased: true },
      { name: '蚝油', quantity: '1 瓶', storage: null, purchased: false }
    ]
  },
  {
    name: '其他',
    icon: '🥚',
    items: [
      { name: '鸡蛋', quantity: '10 个', storage: '7天内', purchased: false }
    ]
  }
])

// 统计数据
const totalItems = computed(() => {
  return categories.value.reduce((sum, cat) => {
    return sum + cat.items.filter(item => !item.purchased).length
  }, 0)
})

const purchasedItems = computed(() => {
  return categories.value.reduce((sum, cat) => {
    return sum + cat.items.filter(item => item.purchased).length
  }, 0)
})

const categoryCount = computed(() => categories.value.length)

// 切换购买状态
const toggleItem = (item) => {
  item.purchased = !item.purchased
}

// 刷新数据
const handleRefresh = async () => {
  loading.value = true
  // 模拟加载
  await new Promise(resolve => setTimeout(resolve, 1000))
  loading.value = false
}

// 清除已购买项
const clearPurchased = () => {
  if (!window.confirm('确定清除所有已购买的项目吗？')) return
  categories.value.forEach(cat => {
    cat.items = cat.items.filter(item => !item.purchased)
  })
  // 移除空分类
  categories.value = categories.value.filter(cat => cat.items.length > 0)
}
</script>

<style scoped>
.shopping-page {
  padding-top: var(--space-4);
}

/* 页面头部 */
.shopping-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--space-5);
}

.shopping-header__title {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-text-heading);
  margin: 0 0 var(--space-1);
}

.shopping-header__desc {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0;
}

.icon-btn svg {
  width: 20px;
  height: 20px;
}

/* 统计卡片 */
.shopping-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-3);
  margin-bottom: var(--space-5);
}

.stat-card {
  background: var(--color-bg-elevated);
  border-radius: var(--radius-xl);
  padding: var(--space-4);
  text-align: center;
  box-shadow: var(--shadow-card);
}

.stat-card--accent {
  background: var(--gradient-primary);
  color: white;
}

.stat-card__value {
  display: block;
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-text-heading);
  margin-bottom: var(--space-1);
}

.stat-card--accent .stat-card__value {
  color: white;
}

.stat-card__label {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
}

.stat-card--accent .stat-card__label {
  color: rgba(255, 255, 255, 0.85);
}

/* 生成清单卡片 */
.generate-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
  margin-bottom: var(--space-6);
}

.generate-card__content h3 {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
  margin: 0 0 var(--space-1);
}

.generate-card__content p {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0;
}

/* 区块标题 */
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-4);
}

.section-title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
  margin: 0;
}

.section-hint {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: var(--space-10) var(--space-6);
}

.empty-state__icon {
  font-size: 48px;
  margin-bottom: var(--space-4);
}

.empty-state__title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
  margin: 0 0 var(--space-2);
}

.empty-state__description {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0;
}

/* 分类列表 */
.category-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.category-section {
  background: var(--color-bg-elevated);
  border-radius: var(--radius-xl);
  padding: var(--space-4);
  box-shadow: var(--shadow-card);
}

.category-section__header {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  margin-bottom: var(--space-3);
  padding-bottom: var(--space-3);
  border-bottom: 1px solid var(--color-border-light);
}

.category-section__icon {
  font-size: 20px;
}

.category-section__title {
  flex: 1;
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
  margin: 0;
}

.category-section__count {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
  background: var(--color-bg-sunken);
  padding: 2px 8px;
  border-radius: var(--radius-full);
}

/* 商品列表 */
.item-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.shopping-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3);
  background: var(--color-bg-sunken);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.shopping-item:hover {
  background: var(--color-border-light);
}

.shopping-item--done {
  opacity: 0.6;
}

.shopping-item--done .shopping-item__name {
  text-decoration: line-through;
}

.shopping-item__checkbox {
  width: 22px;
  height: 22px;
  border: 2px solid var(--color-border-strong);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all var(--transition-fast);
}

.shopping-item--done .shopping-item__checkbox {
  background: var(--color-success-500);
  border-color: var(--color-success-500);
  color: white;
}

.shopping-item__checkbox svg {
  width: 14px;
  height: 14px;
}

.shopping-item__info {
  flex: 1;
  min-width: 0;
}

.shopping-item__name {
  display: block;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-heading);
  margin-bottom: 2px;
}

.shopping-item__quantity {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
}

.shopping-item__storage {
  font-size: var(--font-size-xs);
  color: var(--color-warning-600);
  background: var(--color-warning-100);
  padding: 2px 6px;
  border-radius: var(--radius-sm);
  white-space: nowrap;
}

/* 底部操作 */
.shopping-actions {
  margin-top: var(--space-6);
  padding-top: var(--space-4);
  border-top: 1px solid var(--color-border-light);
}
</style>
