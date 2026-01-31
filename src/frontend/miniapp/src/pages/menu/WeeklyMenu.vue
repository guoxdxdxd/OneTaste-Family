<template>
  <div class="page weekly-menu-page">
    <!-- 页面头部 -->
    <header class="page-header">
      <button class="icon-btn" @click="$router.back()">
        <IconChevronLeft />
      </button>
      <div class="page-header__info">
        <h1 class="page-header__title">周菜单</h1>
        <p class="page-header__subtitle">{{ weekTitle }}</p>
      </div>
    </header>

    <!-- 周导航 -->
    <section class="week-nav">
      <button class="week-nav__btn" @click="prevWeek">
        <IconChevronLeft />
        <span>上周</span>
      </button>
      <button class="week-nav__today" @click="thisWeek">
        本周
      </button>
      <button class="week-nav__btn" @click="nextWeek">
        <span>下周</span>
        <IconChevronRight />
      </button>
    </section>

    <!-- 周日期选择器 -->
    <section class="week-selector card">
      <div class="week-days">
        <button
          v-for="day in weekDays"
          :key="day.date"
          type="button"
          class="week-day"
          :class="{ 
            'week-day--active': selectedDay === day.date,
            'week-day--today': day.isToday,
            'week-day--has-menu': day.hasMenu
          }"
          @click="selectDay(day.date)"
        >
          <span class="week-day__weekday">{{ day.weekday }}</span>
          <span class="week-day__date">{{ day.day }}</span>
          <span v-if="day.hasMenu" class="week-day__dot"></span>
        </button>
      </div>
    </section>

    <!-- 菜单内容 -->
    <section class="menu-content">
      <!-- 加载状态 -->
      <div v-if="loading" class="loading-state">
        <span class="loading-spinner"></span>
        <span>加载菜单中...</span>
      </div>

      <!-- 空状态 -->
      <div v-else-if="!weeklyMenus.length" class="empty-card card">
        <div class="empty-card__icon">📅</div>
        <h3 class="empty-card__title">本周还没有菜单</h3>
        <p class="empty-card__desc">开始规划你的一周美食吧</p>
        <router-link to="/menus/create" class="btn btn--primary">
          创建菜单
        </router-link>
      </div>

      <!-- 周菜单列表 -->
      <div v-else class="day-list">
        <article 
          v-for="day in weekDays" 
          :key="day.date"
          class="day-section"
          :id="`day-${day.date}`"
        >
          <div class="day-section__header">
            <div class="day-section__date">
              <span class="day-section__weekday">{{ day.weekday }}</span>
              <span class="day-section__month-day">{{ day.month }}/{{ day.day }}</span>
            </div>
            <span v-if="day.isToday" class="tag tag--primary tag--pill">今天</span>
            <router-link 
              :to="`/menus/create?date=${day.date}`"
              class="btn btn--text btn--sm"
            >
              添加
            </router-link>
          </div>

          <!-- 有菜单 -->
          <div v-if="getDayMenus(day.date).length" class="day-section__menus">
            <MenuCard
              v-for="menu in getDayMenus(day.date)"
              :key="menu.menu_id"
              :menu="menu"
              :editable="true"
              :compact="true"
              @edit="handleEdit"
            />
          </div>

          <!-- 无菜单 -->
          <div v-else class="day-section__empty">
            <span>暂无安排</span>
          </div>
        </article>
      </div>
    </section>
  </div>
</template>

<script setup>
/**
 * 每周菜单页面
 * 
 * 功能：
 * - 周视图展示
 * - 切换上/下周
 * - 快速创建/编辑菜单
 */
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import MenuCard from '@/components/MenuCard.vue'
import IconChevronLeft from '@/components/icons/IconChevronLeft.vue'
import IconChevronRight from '@/components/icons/IconChevronRight.vue'
import { getWeeklyMenu } from '@/api/menus'

const router = useRouter()

// 当前周的起始日期（周日）
const startDate = ref(getWeekStart(new Date()))
const weeklyMenus = ref([])
const loading = ref(false)
const selectedDay = ref(null)

// 今天
const today = new Date()
const todayStr = formatDate(today)

// 周标题
const weekTitle = computed(() => {
  const endDate = new Date(startDate.value)
  endDate.setDate(endDate.getDate() + 6)
  const start = formatDateShort(startDate.value)
  const end = formatDateShort(endDate)
  return `${start} - ${end}`
})

// 周日期数据
const weekDays = computed(() => {
  const days = []
  const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  
  for (let i = 0; i < 7; i++) {
    const date = new Date(startDate.value)
    date.setDate(date.getDate() + i)
    const dateStr = formatDate(date)
    const hasMenu = weeklyMenus.value.some(menu => menu.date === dateStr)
    
    days.push({
      date: dateStr,
      weekday: weekdays[date.getDay()],
      day: date.getDate(),
      month: date.getMonth() + 1,
      hasMenu,
      isToday: dateStr === todayStr
    })
  }
  
  return days
})

// 获取周起始日期（周日）
function getWeekStart(date) {
  const d = new Date(date)
  const day = d.getDay()
  const diff = d.getDate() - day
  return new Date(d.setDate(diff))
}

// 格式化日期 YYYY-MM-DD
function formatDate(date) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

// 格式化短日期 M/D
function formatDateShort(date) {
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${month}/${day}`
}

// 上一周
function prevWeek() {
  const date = new Date(startDate.value)
  date.setDate(date.getDate() - 7)
  startDate.value = date
  loadWeeklyMenu()
}

// 下一周
function nextWeek() {
  const date = new Date(startDate.value)
  date.setDate(date.getDate() + 7)
  startDate.value = date
  loadWeeklyMenu()
}

// 本周
function thisWeek() {
  startDate.value = getWeekStart(new Date())
  loadWeeklyMenu()
}

// 选择日期
function selectDay(date) {
  selectedDay.value = date
  // 滚动到对应日期
  nextTick(() => {
    const element = document.getElementById(`day-${date}`)
    if (element) {
      element.scrollIntoView({ behavior: 'smooth', block: 'start' })
    }
  })
}

// 获取某天的菜单
function getDayMenus(date) {
  return weeklyMenus.value.filter(menu => menu.date === date)
}

// 加载周菜单
async function loadWeeklyMenu() {
  loading.value = true
  try {
    const dateStr = formatDate(startDate.value)
    const res = await getWeeklyMenu(dateStr)
    if (res.code === 200 && res.data) {
      weeklyMenus.value = res.data.menus || []
    }
  } catch (error) {
    console.error('加载每周菜单失败:', error)
    weeklyMenus.value = []
  } finally {
    loading.value = false
  }
}

// 编辑菜单
function handleEdit(menu) {
  router.push(`/menus/${menu.menu_id}/edit`)
}

// 初始化
onMounted(() => {
  loadWeeklyMenu()
  // 默认选中今天
  if (weekDays.value.some(d => d.isToday)) {
    selectedDay.value = todayStr
  }
})
</script>

<style scoped>
.weekly-menu-page {
  padding-top: var(--space-4);
}

/* 页面头部 */
.page-header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-4);
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

/* 周导航 */
.week-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-4);
}

.week-nav__btn {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  padding: var(--space-2) var(--space-3);
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  cursor: pointer;
  transition: color var(--transition-fast);
}

.week-nav__btn:hover {
  color: var(--color-primary);
}

.week-nav__btn svg {
  width: 16px;
  height: 16px;
}

.week-nav__today {
  padding: var(--space-2) var(--space-4);
  background: var(--color-primary-100);
  border: none;
  border-radius: var(--radius-full);
  color: var(--color-primary);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.week-nav__today:hover {
  background: var(--color-primary-200);
}

/* 周选择器 */
.week-selector {
  margin-bottom: var(--space-5);
  padding: var(--space-3);
}

.week-days {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: var(--space-2);
}

.week-day {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-1);
  padding: var(--space-3) var(--space-1);
  background: transparent;
  border: 2px solid transparent;
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--transition-fast);
  position: relative;
}

.week-day:hover {
  background: var(--color-bg-sunken);
}

.week-day--active {
  background: var(--color-primary-100);
  border-color: var(--color-primary);
}

.week-day--today {
  background: var(--color-primary-50);
}

.week-day--today .week-day__date {
  color: var(--color-primary);
  font-weight: var(--font-weight-bold);
}

.week-day__weekday {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
}

.week-day__date {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
}

.week-day__dot {
  position: absolute;
  bottom: 6px;
  width: 4px;
  height: 4px;
  background: var(--color-primary);
  border-radius: var(--radius-full);
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

/* 空状态卡片 */
.empty-card {
  text-align: center;
  padding: var(--space-8);
}

.empty-card__icon {
  font-size: 48px;
  margin-bottom: var(--space-4);
}

.empty-card__title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
  margin: 0 0 var(--space-2);
}

.empty-card__desc {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0 0 var(--space-5);
}

/* 日列表 */
.day-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.day-section {
  background: var(--color-bg-elevated);
  border-radius: var(--radius-xl);
  padding: var(--space-4);
  box-shadow: var(--shadow-sm);
}

.day-section__header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-3);
}

.day-section__date {
  flex: 1;
  display: flex;
  align-items: baseline;
  gap: var(--space-2);
}

.day-section__weekday {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
}

.day-section__month-day {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.day-section__menus {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.day-section__empty {
  padding: var(--space-4);
  background: var(--color-bg-sunken);
  border-radius: var(--radius-lg);
  text-align: center;
  color: var(--color-text-tertiary);
  font-size: var(--font-size-sm);
}
</style>
