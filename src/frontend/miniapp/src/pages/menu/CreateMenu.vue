<template>
  <div class="page create-menu-page">
    <!-- 页面头部 -->
    <header class="page-header">
      <button class="icon-btn" @click="handleCancel">
        <IconChevronLeft />
      </button>
      <div class="page-header__info">
        <h1 class="page-header__title">创建菜单</h1>
        <p class="page-header__subtitle">规划三餐，轻松生活</p>
      </div>
    </header>

    <!-- 步骤指示器 -->
    <div class="steps-indicator">
      <div 
        v-for="(step, index) in steps" 
        :key="step.key"
        class="step-item"
        :class="{ 
          'step-item--active': currentStep === index,
          'step-item--done': currentStep > index
        }"
      >
        <span class="step-item__number">{{ index + 1 }}</span>
        <span class="step-item__label">{{ step.label }}</span>
      </div>
    </div>

    <!-- 步骤1: 选择日期 -->
    <section v-show="currentStep === 0" class="step-content">
      <div class="card">
        <h2 class="section-title">选择日期</h2>
        <p class="section-desc">选择需要创建菜单的日期</p>
        <Calendar 
          v-model="form.date" 
          :menu-dates="[]" 
          @select="handleDateSelect" 
        />
      </div>
    </section>

    <!-- 步骤2: 选择餐次 -->
    <section v-show="currentStep === 1" class="step-content">
      <div class="card">
        <h2 class="section-title">选择餐次</h2>
        <p class="section-desc">选择要规划的是早餐、午餐还是晚餐</p>
        <div class="meal-type-grid">
          <button
            v-for="type in mealTypes"
            :key="type.value"
            type="button"
            class="meal-type-card"
            :class="{ 'meal-type-card--active': form.meal_type === type.value }"
            @click="selectMealType(type.value)"
          >
            <span class="meal-type-card__icon">{{ type.icon }}</span>
            <span class="meal-type-card__label">{{ type.label }}</span>
            <span class="meal-type-card__time">{{ type.time }}</span>
          </button>
        </div>
      </div>
    </section>

    <!-- 步骤3: 选择菜式 -->
    <section v-show="currentStep === 2" class="step-content">
      <div class="card">
        <h2 class="section-title">选择菜式</h2>
        <p class="section-desc">从家庭菜谱中选择今日要做的菜</p>

        <!-- 加载状态 -->
        <div v-if="loadingDishes" class="loading-state">
          <span class="loading-spinner"></span>
          <span>加载菜式中...</span>
        </div>

        <!-- 空状态 -->
        <div v-else-if="!dishList.length" class="empty-state">
          <div class="empty-state__icon">📖</div>
          <h3 class="empty-state__title">还没有菜式</h3>
          <p class="empty-state__description">请先去创建一些菜式</p>
          <router-link to="/recipes" class="btn btn--primary btn--sm">
            去创建菜式
          </router-link>
        </div>

        <!-- 菜式列表 -->
        <div v-else class="dish-list">
          <div
            v-for="dish in dishList"
            :key="dish.dish_id"
            class="dish-select-item"
            :class="{ 'dish-select-item--selected': form.dish_ids.includes(dish.dish_id) }"
            @click="toggleDish(dish.dish_id)"
          >
            <div class="dish-select-item__checkbox">
              <IconCheck v-if="form.dish_ids.includes(dish.dish_id)" />
            </div>
            <div class="dish-select-item__content">
              <h4 class="dish-select-item__name">{{ dish.name }}</h4>
              <span v-if="dish.category" class="tag tag--default tag--pill">
                {{ getCategoryLabel(dish.category) }}
              </span>
            </div>
          </div>
        </div>

        <!-- 已选数量 -->
        <div v-if="form.dish_ids.length" class="selected-count">
          已选择 <strong>{{ form.dish_ids.length }}</strong> 道菜
        </div>
      </div>
    </section>

    <!-- 底部操作栏 -->
    <div class="action-bar">
      <button 
        v-if="currentStep > 0"
        type="button" 
        class="btn btn--ghost" 
        @click="prevStep"
      >
        上一步
      </button>
      <button 
        v-if="currentStep < 2"
        type="button" 
        class="btn btn--primary" 
        :disabled="!canNext"
        @click="nextStep"
      >
        下一步
      </button>
      <button 
        v-else
        type="button" 
        class="btn btn--primary" 
        :disabled="submitting || !canSubmit"
        @click="handleSubmit"
      >
        <span v-if="submitting" class="loading-spinner loading-spinner--sm"></span>
        {{ submitting ? '创建中...' : '确认创建' }}
      </button>
    </div>
  </div>
</template>

<script setup>
/**
 * 创建菜单页面
 * 
 * 功能：
 * - 分步骤创建菜单（日期 -> 餐次 -> 菜式）
 * - 支持多选菜式
 */
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Calendar from '@/components/Calendar.vue'
import IconChevronLeft from '@/components/icons/IconChevronLeft.vue'
import IconCheck from '@/components/icons/IconCheck.vue'
import { createMenu } from '@/api/menus'
import { fetchDishes } from '@/api/dishes'

const router = useRouter()
const route = useRoute()

// 步骤配置
const steps = [
  { key: 'date', label: '日期' },
  { key: 'meal', label: '餐次' },
  { key: 'dishes', label: '菜式' }
]

const currentStep = ref(0)

// 餐次选项
const mealTypes = [
  { value: 'breakfast', label: '早餐', icon: '🌅', time: '07:00 - 09:00' },
  { value: 'lunch', label: '午餐', icon: '☀️', time: '11:30 - 13:00' },
  { value: 'dinner', label: '晚餐', icon: '🌙', time: '18:00 - 20:00' }
]

// 分类标签
const categoryLabels = {
  meat: '肉类',
  vegetable: '蔬菜',
  soup: '汤羹',
  staple: '主食',
  dessert: '甜品',
  other: '其他'
}

// 表单数据
const form = ref({
  date: new Date().toISOString().split('T')[0],
  meal_type: 'lunch',
  dish_ids: []
})

const dishList = ref([])
const loadingDishes = ref(false)
const submitting = ref(false)

// 从路由获取初始参数
onMounted(() => {
  if (route.query.date) {
    form.value.date = route.query.date
  }
  if (route.query.meal) {
    form.value.meal_type = route.query.meal
  }
  loadDishes()
})

// 是否可以进入下一步
const canNext = computed(() => {
  if (currentStep.value === 0) return !!form.value.date
  if (currentStep.value === 1) return !!form.value.meal_type
  return true
})

// 是否可以提交
const canSubmit = computed(() => {
  return form.value.date && form.value.meal_type && form.value.dish_ids.length > 0
})

// 获取分类标签
const getCategoryLabel = (value) => categoryLabels[value] || value

// 选择日期（选择后自动进入下一步）
const handleDateSelect = (date) => {
  form.value.date = date
  // 选择日期后自动进入下一步
  if (currentStep.value === 0) {
    nextStep()
  }
}

// 选择餐次（选择后自动进入下一步）
const selectMealType = (type) => {
  form.value.meal_type = type
  // 选择餐次后自动进入下一步
  if (currentStep.value === 1) {
    nextStep()
  }
}

// 切换菜式选择
const toggleDish = (dishId) => {
  const index = form.value.dish_ids.indexOf(dishId)
  if (index > -1) {
    form.value.dish_ids.splice(index, 1)
  } else {
    form.value.dish_ids.push(dishId)
  }
}

// 加载菜式列表
const loadDishes = async () => {
  loadingDishes.value = true
  try {
    const res = await fetchDishes()
    if (res.code === 200 && res.data) {
      dishList.value = res.data.dishes || []
    }
  } catch (error) {
    console.error('加载菜式列表失败:', error)
  } finally {
    loadingDishes.value = false
  }
}

// 上一步
const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

// 下一步
const nextStep = () => {
  if (currentStep.value < 2 && canNext.value) {
    currentStep.value++
  }
}

// 提交创建
const handleSubmit = async () => {
  if (!canSubmit.value) return

  submitting.value = true
  try {
    const res = await createMenu({
      date: form.value.date,
      meal_type: form.value.meal_type,
      dish_ids: form.value.dish_ids
    })

    if (res.code === 200) {
      router.back()
    } else {
      alert(res.message || '创建菜单失败')
    }
  } catch (error) {
    console.error('创建菜单失败:', error)
    alert('创建菜单失败，请稍后重试')
  } finally {
    submitting.value = false
  }
}

// 取消
const handleCancel = () => {
  router.back()
}
</script>

<style scoped>
.create-menu-page {
  padding-top: var(--space-4);
  padding-bottom: 100px;
}

/* 页面头部 */
.page-header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-5);
}

.page-header__title {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-text-heading);
  margin: 0;
}

.page-header__subtitle {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0;
}

/* 步骤指示器 */
.steps-indicator {
  display: flex;
  justify-content: center;
  gap: var(--space-2);
  margin-bottom: var(--space-6);
}

.step-item {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  background: var(--color-bg-sunken);
  border-radius: var(--radius-full);
  transition: all var(--transition-fast);
}

.step-item--active {
  background: var(--color-primary);
  color: white;
}

.step-item--done {
  background: var(--color-success-100);
  color: var(--color-success-600);
}

.step-item__number {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.1);
  border-radius: var(--radius-full);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
}

.step-item--active .step-item__number {
  background: rgba(255, 255, 255, 0.2);
}

.step-item--done .step-item__number {
  background: var(--color-success-500);
  color: white;
}

.step-item__label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
}

/* 步骤内容 */
.step-content {
  animation: fadeIn var(--duration-normal) var(--ease-out);
}

.section-title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
  margin: 0 0 var(--space-1);
}

.section-desc {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0 0 var(--space-5);
}

/* 餐次选择 */
.meal-type-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-3);
}

.meal-type-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-5) var(--space-3);
  background: var(--color-bg-sunken);
  border: 2px solid transparent;
  border-radius: var(--radius-xl);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.meal-type-card:hover {
  background: var(--color-bg-elevated);
  box-shadow: var(--shadow-md);
}

.meal-type-card--active {
  background: var(--color-primary-100);
  border-color: var(--color-primary);
}

.meal-type-card__icon {
  font-size: 32px;
}

.meal-type-card__label {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
}

.meal-type-card__time {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
}

.meal-type-card--active .meal-type-card__label {
  color: var(--color-primary-700);
}

/* 菜式列表 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-10);
  color: var(--color-text-secondary);
}

.dish-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.dish-select-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-4);
  background: var(--color-bg-sunken);
  border: 2px solid transparent;
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.dish-select-item:hover {
  background: var(--color-bg-elevated);
}

.dish-select-item--selected {
  background: var(--color-primary-100);
  border-color: var(--color-primary);
}

.dish-select-item__checkbox {
  width: 22px;
  height: 22px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-elevated);
  border: 2px solid var(--color-border-strong);
  border-radius: var(--radius-sm);
  flex-shrink: 0;
  transition: all var(--transition-fast);
}

.dish-select-item--selected .dish-select-item__checkbox {
  background: var(--color-primary);
  border-color: var(--color-primary);
  color: white;
}

.dish-select-item__checkbox svg {
  width: 14px;
  height: 14px;
}

.dish-select-item__content {
  flex: 1;
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.dish-select-item__name {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-heading);
  margin: 0;
}

.selected-count {
  margin-top: var(--space-4);
  padding: var(--space-3) var(--space-4);
  background: var(--color-primary-100);
  border-radius: var(--radius-lg);
  font-size: var(--font-size-sm);
  color: var(--color-primary-700);
  text-align: center;
}

.selected-count strong {
  font-weight: var(--font-weight-bold);
}

/* 底部操作栏 */
.action-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  gap: var(--space-3);
  padding: var(--space-4) var(--space-5);
  padding-bottom: calc(var(--space-4) + var(--safe-area-bottom));
  background: var(--color-bg-elevated);
  border-top: 1px solid var(--color-border-light);
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.08);
  z-index: var(--z-fixed);
}

.action-bar .btn {
  flex: 1;
}
</style>
