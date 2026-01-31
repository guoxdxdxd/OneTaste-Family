<template>
  <div class="page edit-menu-page">
    <!-- 页面头部 -->
    <header class="page-header">
      <button class="icon-btn" @click="handleCancel">
        <IconChevronLeft />
      </button>
      <div class="page-header__info">
        <h1 class="page-header__title">编辑菜单</h1>
        <p class="page-header__subtitle">修改菜单内容</p>
      </div>
      <button 
        class="btn btn--text btn--sm" 
        @click="handleDelete"
        :disabled="deleting"
      >
        {{ deleting ? '删除中...' : '删除' }}
      </button>
    </header>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <span class="loading-spinner"></span>
      <span>加载菜单信息...</span>
    </div>

    <template v-else>
      <!-- 菜单信息摘要 -->
      <section class="menu-summary card">
        <div class="menu-summary__item">
          <span class="menu-summary__label">日期</span>
          <span class="menu-summary__value">{{ formatDateDisplay(form.date) }}</span>
        </div>
        <div class="menu-summary__item">
          <span class="menu-summary__label">餐次</span>
          <span class="menu-summary__value">{{ getMealLabel(form.meal_type) }}</span>
        </div>
      </section>

      <!-- 日期选择 -->
      <section class="card">
        <div class="card-header-toggle" @click="showDatePicker = !showDatePicker">
          <h2 class="section-title">修改日期</h2>
          <IconChevronRight 
            class="toggle-icon" 
            :class="{ 'toggle-icon--open': showDatePicker }"
          />
        </div>
        <transition name="slide">
          <div v-if="showDatePicker" class="picker-content">
            <Calendar 
              v-model="form.date" 
              :menu-dates="[]" 
              @select="handleDateSelect" 
            />
          </div>
        </transition>
      </section>

      <!-- 餐次选择 -->
      <section class="card">
        <div class="card-header-toggle" @click="showMealPicker = !showMealPicker">
          <h2 class="section-title">修改餐次</h2>
          <IconChevronRight 
            class="toggle-icon" 
            :class="{ 'toggle-icon--open': showMealPicker }"
          />
        </div>
        <transition name="slide">
          <div v-if="showMealPicker" class="picker-content">
            <div class="meal-type-grid">
              <button
                v-for="type in mealTypes"
                :key="type.value"
                type="button"
                class="meal-type-card"
                :class="{ 'meal-type-card--active': form.meal_type === type.value }"
                @click="form.meal_type = type.value"
              >
                <span class="meal-type-card__icon">{{ type.icon }}</span>
                <span class="meal-type-card__label">{{ type.label }}</span>
              </button>
            </div>
          </div>
        </transition>
      </section>

      <!-- 菜式选择 -->
      <section class="card">
        <h2 class="section-title">选择菜式</h2>
        <p class="section-desc">已选择 {{ form.dish_ids.length }} 道菜</p>

        <!-- 加载状态 -->
        <div v-if="loadingDishes" class="loading-state loading-state--sm">
          <span class="loading-spinner"></span>
          <span>加载菜式中...</span>
        </div>

        <!-- 空状态 -->
        <div v-else-if="!dishList.length" class="empty-state">
          <div class="empty-state__icon">📖</div>
          <h3 class="empty-state__title">还没有菜式</h3>
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
      </section>
    </template>

    <!-- 底部操作栏 -->
    <div class="action-bar">
      <button 
        type="button" 
        class="btn btn--ghost" 
        @click="handleCancel"
      >
        取消
      </button>
      <button 
        type="button" 
        class="btn btn--primary" 
        :disabled="submitting || !canSubmit"
        @click="handleSubmit"
      >
        <span v-if="submitting" class="loading-spinner loading-spinner--sm"></span>
        {{ submitting ? '保存中...' : '保存修改' }}
      </button>
    </div>
  </div>
</template>

<script setup>
/**
 * 编辑菜单页面
 * 
 * 功能：
 * - 加载菜单详情
 * - 编辑日期、餐次、菜式
 * - 删除菜单
 */
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Calendar from '@/components/Calendar.vue'
import IconChevronLeft from '@/components/icons/IconChevronLeft.vue'
import IconChevronRight from '@/components/icons/IconChevronRight.vue'
import IconCheck from '@/components/icons/IconCheck.vue'
import { updateMenu, deleteMenu, getDailyMenu } from '@/api/menus'
import { fetchDishes } from '@/api/dishes'

const router = useRouter()
const route = useRoute()

const menuId = route.params.id

// 餐次选项
const mealTypes = [
  { value: 'breakfast', label: '早餐', icon: '🌅' },
  { value: 'lunch', label: '午餐', icon: '☀️' },
  { value: 'dinner', label: '晚餐', icon: '🌙' }
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

// 状态
const loading = ref(true)
const loadingDishes = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const showDatePicker = ref(false)
const showMealPicker = ref(false)

// 表单数据
const form = ref({
  date: '',
  meal_type: 'lunch',
  dish_ids: []
})

const dishList = ref([])

// 是否可以提交
const canSubmit = computed(() => {
  return form.value.date && form.value.meal_type && form.value.dish_ids.length > 0
})

// 获取餐次标签
const getMealLabel = (value) => {
  const type = mealTypes.find(t => t.value === value)
  return type ? type.label : value
}

// 获取分类标签
const getCategoryLabel = (value) => categoryLabels[value] || value

// 格式化日期显示
const formatDateDisplay = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const month = date.getMonth() + 1
  const day = date.getDate()
  const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  const weekday = weekdays[date.getDay()]
  return `${month}月${day}日 ${weekday}`
}

// 选择日期
const handleDateSelect = (date) => {
  form.value.date = date
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

// 加载菜单详情
const loadMenu = async () => {
  loading.value = true
  try {
    // 尝试通过获取每日菜单来找到当前菜单
    const today = new Date().toISOString().split('T')[0]
    const res = await getDailyMenu(today)
    if (res.code === 200 && res.data) {
      const menu = res.data.menus?.find(m => String(m.menu_id) === String(menuId))
      if (menu) {
        form.value.date = menu.date
        form.value.meal_type = menu.meal_type
        form.value.dish_ids = menu.dishes?.map(d => d.dish_id) || []
      }
    }
  } catch (error) {
    console.error('加载菜单失败:', error)
  } finally {
    loading.value = false
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

// 提交修改
const handleSubmit = async () => {
  if (!canSubmit.value) return

  submitting.value = true
  try {
    const res = await updateMenu(menuId, {
      date: form.value.date,
      meal_type: form.value.meal_type,
      dish_ids: form.value.dish_ids
    })

    if (res.code === 200) {
      router.back()
    } else {
      alert(res.message || '保存失败')
    }
  } catch (error) {
    console.error('保存失败:', error)
    alert('保存失败，请稍后重试')
  } finally {
    submitting.value = false
  }
}

// 删除菜单
const handleDelete = async () => {
  if (!window.confirm('确定要删除这个菜单吗？')) return

  deleting.value = true
  try {
    const res = await deleteMenu(menuId)
    if (res.code === 200) {
      router.back()
    } else {
      alert(res.message || '删除失败')
    }
  } catch (error) {
    console.error('删除失败:', error)
    alert('删除失败，请稍后重试')
  } finally {
    deleting.value = false
  }
}

// 取消
const handleCancel = () => {
  router.back()
}

// 初始化
onMounted(() => {
  loadMenu()
  loadDishes()
})
</script>

<style scoped>
.edit-menu-page {
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

.page-header__info {
  flex: 1;
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

/* 加载状态 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-10);
  color: var(--color-text-secondary);
}

.loading-state--sm {
  padding: var(--space-6);
}

/* 菜单摘要 */
.menu-summary {
  display: flex;
  gap: var(--space-6);
  margin-bottom: var(--space-4);
}

.menu-summary__item {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.menu-summary__label {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.menu-summary__value {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
}

/* 卡片区块 */
.card {
  margin-bottom: var(--space-4);
}

.card-header-toggle {
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  padding: var(--space-1) 0;
}

.toggle-icon {
  width: 18px;
  height: 18px;
  color: var(--color-text-tertiary);
  transition: transform var(--transition-fast);
}

.toggle-icon--open {
  transform: rotate(90deg);
}

.picker-content {
  margin-top: var(--space-4);
}

.section-title {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
  margin: 0;
}

.section-desc {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: var(--space-1) 0 var(--space-4);
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
  padding: var(--space-4) var(--space-3);
  background: var(--color-bg-sunken);
  border: 2px solid transparent;
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.meal-type-card:hover {
  background: var(--color-bg-elevated);
}

.meal-type-card--active {
  background: var(--color-primary-100);
  border-color: var(--color-primary);
}

.meal-type-card__icon {
  font-size: 24px;
}

.meal-type-card__label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
}

.meal-type-card--active .meal-type-card__label {
  color: var(--color-primary-700);
}

/* 菜式列表 */
.dish-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  max-height: 400px;
  overflow-y: auto;
}

.dish-select-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3);
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
  width: 20px;
  height: 20px;
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
  width: 12px;
  height: 12px;
}

.dish-select-item__content {
  flex: 1;
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.dish-select-item__name {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-heading);
  margin: 0;
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
}

.slide-enter-to,
.slide-leave-from {
  max-height: 500px;
}
</style>
