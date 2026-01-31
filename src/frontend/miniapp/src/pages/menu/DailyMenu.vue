<template>
  <div class="page daily-menu-page">
    <!-- 页面头部 -->
    <header class="page-header">
      <button class="icon-btn" @click="$router.back()">
        <IconChevronLeft />
      </button>
      <div class="page-header__info">
        <h1 class="page-header__title">每日菜单</h1>
        <p class="page-header__subtitle">{{ formatDateTitle(selectedDate) }}</p>
      </div>
      <router-link :to="`/menus/create?date=${selectedDate}`" class="btn btn--primary btn--sm">
        添加
      </router-link>
    </header>

    <!-- 日历选择 -->
    <section class="calendar-section card">
      <Calendar
        v-model="selectedDate"
        :menu-dates="menuDates"
        @select="loadDailyMenu"
      />
    </section>

    <!-- 菜单内容 -->
    <section class="menu-content">
      <!-- 加载状态 -->
      <div v-if="loading" class="loading-state">
        <span class="loading-spinner"></span>
        <span>加载菜单中...</span>
      </div>

      <!-- 空状态 -->
      <div v-else-if="!menus.length" class="empty-card card">
        <div class="empty-card__icon">📋</div>
        <h3 class="empty-card__title">这一天还没有菜单</h3>
        <p class="empty-card__desc">点击右上角按钮创建菜单</p>
        <router-link 
          :to="`/menus/create?date=${selectedDate}`" 
          class="btn btn--primary"
        >
          创建菜单
        </router-link>
      </div>

      <!-- 三餐列表 -->
      <div v-else class="meal-list">
        <article 
          v-for="meal in mealSections" 
          :key="meal.type"
          class="meal-section"
        >
          <div class="meal-section__header">
            <span class="meal-section__icon">{{ meal.icon }}</span>
            <div class="meal-section__info">
              <h3 class="meal-section__title">{{ meal.label }}</h3>
              <span class="meal-section__time">{{ meal.time }}</span>
            </div>
          </div>

          <!-- 有菜单 -->
          <div v-if="meal.menu" class="meal-section__content">
            <MenuCard 
              :menu="meal.menu" 
              :editable="true"
              @edit="handleEdit"
            />
          </div>

          <!-- 无菜单 -->
          <div v-else class="meal-section__empty">
            <p>暂未安排</p>
            <router-link 
              :to="`/menus/create?date=${selectedDate}&meal=${meal.type}`"
              class="btn btn--ghost btn--sm"
            >
              添加
            </router-link>
          </div>
        </article>
      </div>
    </section>
  </div>
</template>

<script setup>
/**
 * 每日菜单页面
 * 
 * 功能：
 * - 日历选择日期
 * - 展示当天三餐菜单
 * - 快速创建/编辑菜单
 */
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Calendar from '@/components/Calendar.vue'
import MenuCard from '@/components/MenuCard.vue'
import IconChevronLeft from '@/components/icons/IconChevronLeft.vue'
import { getDailyMenu, getWeeklyMenu } from '@/api/menus'

const router = useRouter()

// 当前选中日期
const today = new Date()
const todayStr = formatDate(today)
const selectedDate = ref(todayStr)

// 数据状态
const menus = ref([])
const menuDates = ref([])
const loading = ref(false)

// 餐次配置
const mealConfig = [
  { type: 'breakfast', label: '早餐', icon: '🌅', time: '07:00 - 09:00' },
  { type: 'lunch', label: '午餐', icon: '☀️', time: '11:30 - 13:00' },
  { type: 'dinner', label: '晚餐', icon: '🌙', time: '18:00 - 20:00' }
]

// 三餐区块数据
const mealSections = computed(() => {
  return mealConfig.map(config => {
    const menu = menus.value.find(m => m.meal_type === config.type)
    return {
      ...config,
      menu
    }
  })
})

// 格式化日期字符串
function formatDate(date) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

// 格式化日期标题
function formatDateTitle(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const month = date.getMonth() + 1
  const day = date.getDate()
  const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  const weekday = weekdays[date.getDay()]
  
  // 判断是否是今天
  if (dateStr === todayStr) {
    return `今天 ${month}月${day}日 ${weekday}`
  }
  
  return `${month}月${day}日 ${weekday}`
}

// 加载每日菜单
async function loadDailyMenu(date) {
  if (!date) return

  loading.value = true
  try {
    const res = await getDailyMenu(date)
    if (res.code === 200 && res.data) {
      menus.value = res.data.menus || []
    }
  } catch (error) {
    console.error('加载每日菜单失败:', error)
    menus.value = []
  } finally {
    loading.value = false
  }
}

// 加载周菜单获取有菜单的日期
async function loadWeeklyMenuDates() {
  try {
    const res = await getWeeklyMenu(todayStr)
    if (res.code === 200 && res.data) {
      const allMenus = res.data.menus || []
      menuDates.value = [...new Set(allMenus.map(m => m.date))]
    }
  } catch (error) {
    console.error('加载周菜单日期失败:', error)
  }
}

// 编辑菜单
function handleEdit(menu) {
  router.push(`/menus/${menu.menu_id}/edit`)
}

// 初始化
onMounted(() => {
  loadDailyMenu(selectedDate.value)
  loadWeeklyMenuDates()
})
</script>

<style scoped>
.daily-menu-page {
  padding-top: var(--space-4);
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

/* 日历区块 */
.calendar-section {
  margin-bottom: var(--space-5);
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

/* 三餐列表 */
.meal-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.meal-section {
  background: var(--color-bg-elevated);
  border-radius: var(--radius-xl);
  padding: var(--space-4);
  box-shadow: var(--shadow-card);
}

.meal-section__header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-4);
  padding-bottom: var(--space-3);
  border-bottom: 1px solid var(--color-border-light);
}

.meal-section__icon {
  font-size: 28px;
}

.meal-section__info {
  flex: 1;
}

.meal-section__title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
  margin: 0;
}

.meal-section__time {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
}

.meal-section__content {
  /* MenuCard 样式由组件自带 */
}

.meal-section__empty {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4);
  background: var(--color-bg-sunken);
  border-radius: var(--radius-lg);
  border: 1px dashed var(--color-border-default);
}

.meal-section__empty p {
  font-size: var(--font-size-sm);
  color: var(--color-text-tertiary);
  margin: 0;
}
</style>
