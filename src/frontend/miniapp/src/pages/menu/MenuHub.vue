<template>
  <div class="page menu-hub">
    <!-- 顶部问候区域 -->
    <header class="menu-hub__header">
      <div class="menu-hub__greeting">
        <p class="menu-hub__date">{{ formattedDate }}</p>
        <h1 class="menu-hub__title">{{ greeting }}</h1>
      </div>
      <button class="icon-btn" @click="goToRecipes">
        <IconBook class="icon-btn__icon" />
      </button>
    </header>

    <!-- 今日菜单卡片 -->
    <section class="today-menu">
      <div class="section-header">
        <h2 class="section-title">今日菜单</h2>
        <router-link to="/menus/daily" class="section-link">
          查看详情
          <IconChevronRight class="section-link__icon" />
        </router-link>
      </div>

      <div class="meal-cards">
        <article 
          v-for="meal in todayMeals" 
          :key="meal.type"
          class="meal-card"
          :class="{ 'meal-card--empty': !meal.dishes.length }"
          @click="handleMealClick(meal)"
        >
          <div class="meal-card__header">
            <span class="meal-card__icon">{{ meal.icon }}</span>
            <div class="meal-card__info">
              <span class="meal-card__label">{{ meal.label }}</span>
              <span class="meal-card__time">{{ meal.time }}</span>
            </div>
          </div>
          <div class="meal-card__content">
            <template v-if="meal.dishes.length">
              <span 
                v-for="(dish, index) in meal.dishes.slice(0, 3)" 
                :key="index"
                class="meal-card__dish"
              >
                {{ dish }}
              </span>
              <span v-if="meal.dishes.length > 3" class="meal-card__more">
                +{{ meal.dishes.length - 3 }}
              </span>
            </template>
            <span v-else class="meal-card__empty-text">点击添加</span>
          </div>
        </article>
      </div>
    </section>

    <!-- 快捷操作 -->
    <section class="quick-actions">
      <router-link to="/menus/create" class="quick-action-card quick-action-card--primary">
        <div class="quick-action-card__icon">
          <IconPlus />
        </div>
        <div class="quick-action-card__content">
          <h3>创建菜单</h3>
          <p>规划今日或未来的三餐</p>
        </div>
      </router-link>

      <router-link to="/menus/weekly" class="quick-action-card">
        <div class="quick-action-card__icon">
          <IconCalendar />
        </div>
        <div class="quick-action-card__content">
          <h3>周菜单</h3>
          <p>查看本周计划</p>
        </div>
      </router-link>
    </section>

    <!-- 日历视图 -->
    <section class="calendar-section">
      <div class="section-header">
        <h2 class="section-title">选择日期</h2>
      </div>
      <Calendar 
        v-model="selectedDate" 
        :menu-dates="menuDates"
        @select="handleDateSelect"
      />
    </section>

    <!-- 选中日期的菜单 -->
    <section v-if="selectedDateMenus.length" class="selected-menu">
      <div class="section-header">
        <h2 class="section-title">{{ selectedDateLabel }} 的菜单</h2>
      </div>
      <div class="selected-menu__list">
        <article 
          v-for="menu in selectedDateMenus" 
          :key="menu.id"
          class="menu-item"
          @click="viewMenuDetail(menu)"
        >
          <div class="menu-item__info">
            <span class="menu-item__label">{{ getMealLabel(menu.meal_type) }}</span>
            <span class="menu-item__dishes">
              {{ menu.dishes?.map(d => d.name).join('、') || '暂无菜品' }}
            </span>
          </div>
          <IconChevronRight class="menu-item__arrow" />
        </article>
      </div>
    </section>

    <!-- 空状态提示 -->
    <section v-else-if="selectedDate" class="empty-state-card">
      <div class="empty-state-card__content">
        <p class="empty-state-card__text">{{ selectedDateLabel }} 还没有菜单</p>
        <router-link 
          :to="`/menus/create?date=${selectedDate}`" 
          class="btn btn--primary btn--sm"
        >
          创建菜单
        </router-link>
      </div>
    </section>

    <!-- AI 助手入口（预留） -->
    <section class="ai-entry card card--highlight">
      <div class="ai-entry__content">
        <div class="ai-entry__badge">即将上线</div>
        <h3 class="ai-entry__title">AI 智能助手</h3>
        <p class="ai-entry__desc">根据家庭成员的健康状况和口味偏好，智能推荐每日菜单</p>
      </div>
      <button class="btn btn--ghost btn--sm" disabled>敬请期待</button>
    </section>
  </div>
</template>

<script setup>
/**
 * 菜单中心页面
 * 
 * 功能：
 * - 展示今日菜单
 * - 日历视图选择日期查看菜单
 * - 快捷入口创建/查看菜单
 */
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Calendar from '@/components/Calendar.vue'
import IconBook from '@/components/icons/IconBook.vue'
import IconPlus from '@/components/icons/IconPlus.vue'
import IconCalendar from '@/components/icons/IconCalendar.vue'
import IconChevronRight from '@/components/icons/IconChevronRight.vue'
import { getDailyMenu, getWeeklyMenu } from '@/api/menus'

const router = useRouter()
const userStore = useUserStore()

// 当前选中日期
const today = new Date()
const todayStr = formatDateStr(today)
const selectedDate = ref(todayStr)

// 菜单数据
const menuDates = ref([])
const todayMenuData = ref([])
const selectedDateMenus = ref([])

// 格式化日期字符串 YYYY-MM-DD
function formatDateStr(date) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

// 格式化显示日期
const formattedDate = computed(() => {
  const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  const month = today.getMonth() + 1
  const day = today.getDate()
  const weekday = weekdays[today.getDay()]
  return `${month}月${day}日 ${weekday}`
})

// 问候语
const greeting = computed(() => {
  const hour = today.getHours()
  const name = userStore.nickname || '主人'
  if (hour < 6) return `${name}，夜深了`
  if (hour < 11) return `${name}，早上好`
  if (hour < 14) return `${name}，中午好`
  if (hour < 18) return `${name}，下午好`
  return `${name}，晚上好`
})

// 选中日期的标签
const selectedDateLabel = computed(() => {
  if (selectedDate.value === todayStr) return '今天'
  const date = new Date(selectedDate.value)
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${month}月${day}日`
})

// 今日三餐数据
const todayMeals = computed(() => {
  const meals = [
    { type: 'breakfast', label: '早餐', time: '07:00', icon: '🌅', dishes: [] },
    { type: 'lunch', label: '午餐', time: '12:00', icon: '☀️', dishes: [] },
    { type: 'dinner', label: '晚餐', time: '18:00', icon: '🌙', dishes: [] }
  ]

  todayMenuData.value.forEach(menu => {
    const meal = meals.find(m => m.type === menu.meal_type)
    if (meal && menu.dishes) {
      meal.dishes = menu.dishes.map(d => d.name)
    }
  })

  return meals
})

// 获取餐次标签
const getMealLabel = (type) => {
  const labels = {
    breakfast: '早餐',
    lunch: '午餐',
    dinner: '晚餐'
  }
  return labels[type] || type
}

// 加载今日菜单
const loadTodayMenus = async () => {
  try {
    const res = await getDailyMenu(todayStr)
    if (res.code === 200) {
      todayMenuData.value = res.data?.menus || []
    }
  } catch (error) {
    console.error('加载今日菜单失败:', error)
  }
}

// 加载周菜单（获取有菜单的日期）
const loadWeeklyMenus = async () => {
  try {
    const res = await getWeeklyMenu(todayStr)
    if (res.code === 200) {
      const menus = res.data?.menus || []
      menuDates.value = [...new Set(menus.map(m => m.date))]
    }
  } catch (error) {
    console.error('加载周菜单失败:', error)
  }
}

// 加载选中日期的菜单
const loadSelectedDateMenus = async (date) => {
  try {
    const res = await getDailyMenu(date)
    if (res.code === 200) {
      selectedDateMenus.value = res.data?.menus || []
    }
  } catch (error) {
    console.error('加载日期菜单失败:', error)
    selectedDateMenus.value = []
  }
}

// 处理日期选择
const handleDateSelect = (date) => {
  selectedDate.value = date
  loadSelectedDateMenus(date)
}

// 处理餐次点击
const handleMealClick = (meal) => {
  if (meal.dishes.length) {
    router.push('/menus/daily')
  } else {
    router.push(`/menus/create?date=${todayStr}&meal=${meal.type}`)
  }
}

// 查看菜单详情
const viewMenuDetail = (menu) => {
  router.push(`/menus/${menu.menu_id}/edit`)
}

// 跳转到菜谱管理
const goToRecipes = () => {
  router.push('/recipes')
}

// 页面初始化
onMounted(async () => {
  await Promise.all([
    loadTodayMenus(),
    loadWeeklyMenus()
  ])
  // 默认加载今日菜单
  selectedDateMenus.value = todayMenuData.value
})
</script>

<style scoped>
.menu-hub {
  padding-top: var(--space-4);
}

/* 顶部问候区域 */
.menu-hub__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--space-6);
}

.menu-hub__date {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0 0 var(--space-1);
}

.menu-hub__title {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-text-heading);
  margin: 0;
}

.icon-btn__icon {
  width: 24px;
  height: 24px;
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

.section-link {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.section-link:hover {
  color: var(--color-primary);
}

.section-link__icon {
  width: 16px;
  height: 16px;
}

/* 今日菜单区域 */
.today-menu {
  margin-bottom: var(--space-6);
}

.meal-cards {
  display: grid;
  /* 三列自适应，每列最小100px */
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: var(--space-3);
}

/* 窄屏幕（<380px）改为单列 */
@media (max-width: 379px) {
  .meal-cards {
    grid-template-columns: 1fr;
  }
}

/* 中等屏幕（380px-500px）改为两列，第三个占满一行 */
@media (min-width: 380px) and (max-width: 500px) {
  .meal-cards {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .meal-cards .meal-card:nth-child(3) {
    grid-column: 1 / -1;
  }
}

.meal-card {
  background: var(--color-bg-elevated);
  border-radius: var(--radius-xl);
  padding: var(--space-3);
  box-shadow: var(--shadow-card);
  cursor: pointer;
  transition: all var(--transition-normal);
  /* 防止内容溢出 */
  min-width: 0;
  overflow: hidden;
}

.meal-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-card-hover);
}

.meal-card--empty {
  background: var(--color-bg-sunken);
  box-shadow: none;
  border: 1px dashed var(--color-border-default);
}

.meal-card__header {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  margin-bottom: var(--space-2);
}

.meal-card__icon {
  font-size: 20px;
  flex-shrink: 0;
}

.meal-card__info {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.meal-card__label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
  white-space: nowrap;
}

.meal-card__time {
  font-size: 10px;
  color: var(--color-text-tertiary);
  white-space: nowrap;
}

.meal-card__content {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  min-height: 24px;
}

.meal-card__dish {
  font-size: 11px;
  color: var(--color-text-primary);
  background: var(--color-bg-sunken);
  padding: 2px 6px;
  border-radius: var(--radius-sm);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
}

.meal-card__more {
  font-size: 11px;
  color: var(--color-text-tertiary);
  white-space: nowrap;
}

.meal-card__empty-text {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
}

/* 窄屏单列时，卡片改为横向布局 */
@media (max-width: 379px) {
  .meal-card {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-3) var(--space-4);
  }
  
  .meal-card__header {
    margin-bottom: 0;
    flex-shrink: 0;
  }
  
  .meal-card__content {
    flex: 1;
    min-height: auto;
    justify-content: flex-end;
  }
}

/* 快捷操作 */
.quick-actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-3);
  margin-bottom: var(--space-6);
}

.quick-action-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  padding: var(--space-4);
  background: var(--color-bg-elevated);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-card);
  text-decoration: none;
  transition: all var(--transition-normal);
}

.quick-action-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-card-hover);
}

.quick-action-card--primary {
  background: var(--gradient-primary);
  color: var(--color-text-inverse);
}

.quick-action-card--primary .quick-action-card__icon {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

.quick-action-card--primary h3,
.quick-action-card--primary p {
  color: white;
}

.quick-action-card--primary p {
  opacity: 0.85;
}

.quick-action-card__icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-primary-100);
  color: var(--color-primary);
  border-radius: var(--radius-lg);
}

.quick-action-card__icon svg {
  width: 20px;
  height: 20px;
}

.quick-action-card__content h3 {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
  margin: 0 0 var(--space-1);
}

.quick-action-card__content p {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  margin: 0;
}

/* 日历区域 */
.calendar-section {
  margin-bottom: var(--space-6);
}

/* 选中日期菜单 */
.selected-menu {
  margin-bottom: var(--space-6);
}

.selected-menu__list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.menu-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4);
  background: var(--color-bg-elevated);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.menu-item:hover {
  background: var(--color-bg-sunken);
}

.menu-item__info {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.menu-item__label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
}

.menu-item__dishes {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.menu-item__arrow {
  width: 16px;
  height: 16px;
  color: var(--color-text-tertiary);
}

/* 空状态卡片 */
.empty-state-card {
  background: var(--color-bg-sunken);
  border-radius: var(--radius-xl);
  padding: var(--space-6);
  text-align: center;
  margin-bottom: var(--space-6);
}

.empty-state-card__text {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0 0 var(--space-4);
}

/* AI 入口 */
.ai-entry {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
}

.ai-entry__content {
  flex: 1;
}

.ai-entry__badge {
  display: inline-block;
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
  color: var(--color-primary);
  background: var(--color-primary-100);
  padding: 2px 8px;
  border-radius: var(--radius-full);
  margin-bottom: var(--space-2);
}

.ai-entry__title {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
  margin: 0 0 var(--space-1);
}

.ai-entry__desc {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0;
}
</style>
