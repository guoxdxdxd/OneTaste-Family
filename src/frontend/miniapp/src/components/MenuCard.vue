<template>
  <article 
    class="menu-card" 
    :class="{ 
      'menu-card--compact': compact,
      'menu-card--interactive': editable
    }"
    @click="editable && $emit('edit', menu)"
  >
    <!-- 卡片头部 -->
    <div class="menu-card__header">
      <div class="menu-card__meal">
        <span class="menu-card__icon">{{ mealIcon }}</span>
        <div class="menu-card__info">
          <span class="menu-card__label">{{ mealTypeLabel }}</span>
          <span v-if="!compact" class="menu-card__date">{{ formatDate(menu.date) }}</span>
        </div>
      </div>
      <button
        v-if="editable"
        type="button"
        class="btn btn--text btn--sm"
        @click.stop="$emit('edit', menu)"
      >
        编辑
      </button>
    </div>

    <!-- 菜式列表 -->
    <div v-if="hasDishes" class="menu-card__dishes">
      <div 
        v-for="dish in displayDishes" 
        :key="dish.dish_id"
        class="menu-card__dish"
      >
        <span class="menu-card__dish-name">{{ dish.name }}</span>
        <span v-if="dish.category" class="menu-card__dish-category">
          {{ getCategoryLabel(dish.category) }}
        </span>
      </div>
      <div v-if="remainingCount > 0" class="menu-card__more">
        +{{ remainingCount }} 道菜
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="menu-card__empty">
      <p>暂未安排菜式</p>
      <button
        v-if="editable"
        type="button"
        class="btn btn--primary btn--sm"
        @click.stop="$emit('add', menu)"
      >
        添加菜式
      </button>
    </div>

    <!-- AI 标记 -->
    <div v-if="menu.source === 'ai'" class="menu-card__footer">
      <span class="menu-card__ai-badge">
        <IconStar />
        AI 推荐
      </span>
    </div>
  </article>
</template>

<script setup>
/**
 * 菜单卡片组件
 * 
 * 用于展示单个餐次的菜单信息
 * 
 * Props:
 * - menu: 菜单数据对象
 * - editable: 是否可编辑
 * - compact: 紧凑模式
 * 
 * Events:
 * - edit: 点击编辑
 * - add: 点击添加菜式
 */
import { computed } from 'vue'

// 简单的星星图标
const IconStar = {
  template: `
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" stroke="none">
      <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
    </svg>
  `
}

const props = defineProps({
  menu: {
    type: Object,
    required: true
  },
  editable: {
    type: Boolean,
    default: false
  },
  compact: {
    type: Boolean,
    default: false
  }
})

defineEmits(['edit', 'add'])

// 餐次映射
const mealTypeConfig = {
  breakfast: { label: '早餐', icon: '🌅' },
  lunch: { label: '午餐', icon: '☀️' },
  dinner: { label: '晚餐', icon: '🌙' }
}

// 分类映射
const categoryLabels = {
  meat: '肉类',
  vegetable: '蔬菜',
  soup: '汤羹',
  staple: '主食',
  dessert: '甜品',
  other: '其他'
}

// 餐次标签
const mealTypeLabel = computed(() => {
  return mealTypeConfig[props.menu.meal_type]?.label || props.menu.meal_type
})

// 餐次图标
const mealIcon = computed(() => {
  return mealTypeConfig[props.menu.meal_type]?.icon || '🍽️'
})

// 是否有菜式
const hasDishes = computed(() => {
  return props.menu.dishes && props.menu.dishes.length > 0
})

// 显示的菜式（最多显示3个）
const displayDishes = computed(() => {
  if (!hasDishes.value) return []
  return props.compact 
    ? props.menu.dishes.slice(0, 2) 
    : props.menu.dishes.slice(0, 4)
})

// 剩余菜式数量
const remainingCount = computed(() => {
  if (!hasDishes.value) return 0
  const max = props.compact ? 2 : 4
  return Math.max(0, props.menu.dishes.length - max)
})

// 获取分类标签
const getCategoryLabel = (value) => categoryLabels[value] || ''

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${month}月${day}日`
}
</script>

<style scoped>
.menu-card {
  background: var(--color-bg-elevated);
  border-radius: var(--radius-lg);
  padding: var(--space-4);
  border: 1px solid var(--color-border-light);
  transition: all var(--transition-fast);
}

.menu-card--interactive {
  cursor: pointer;
}

.menu-card--interactive:hover {
  border-color: var(--color-primary-200);
  box-shadow: var(--shadow-md);
}

.menu-card--compact {
  padding: var(--space-3);
}

/* 头部 */
.menu-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-3);
}

.menu-card__meal {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.menu-card__icon {
  font-size: 20px;
}

.menu-card--compact .menu-card__icon {
  font-size: 16px;
}

.menu-card__label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
}

.menu-card__date {
  display: block;
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
  margin-top: 2px;
}

/* 菜式列表 */
.menu-card__dishes {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.menu-card__dish {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-2) var(--space-3);
  background: var(--color-bg-sunken);
  border-radius: var(--radius-md);
}

.menu-card--compact .menu-card__dish {
  padding: var(--space-1) var(--space-2);
}

.menu-card__dish-name {
  font-size: var(--font-size-sm);
  color: var(--color-text-primary);
}

.menu-card--compact .menu-card__dish-name {
  font-size: var(--font-size-xs);
}

.menu-card__dish-category {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
}

.menu-card__more {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
  text-align: center;
  padding: var(--space-1);
}

/* 空状态 */
.menu-card__empty {
  text-align: center;
  padding: var(--space-4);
  background: var(--color-bg-sunken);
  border-radius: var(--radius-lg);
  border: 1px dashed var(--color-border-default);
}

.menu-card__empty p {
  font-size: var(--font-size-sm);
  color: var(--color-text-tertiary);
  margin: 0 0 var(--space-3);
}

/* AI 标记 */
.menu-card__footer {
  margin-top: var(--space-3);
  padding-top: var(--space-3);
  border-top: 1px solid var(--color-border-light);
}

.menu-card__ai-badge {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  padding: var(--space-1) var(--space-2);
  background: var(--color-warning-100);
  color: var(--color-warning-600);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
}

.menu-card__ai-badge svg {
  width: 12px;
  height: 12px;
}
</style>
