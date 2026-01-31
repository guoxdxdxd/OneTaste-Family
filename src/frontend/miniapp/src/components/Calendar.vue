<template>
  <div class="calendar card card--flat">
    <!-- 日历头部 -->
    <div class="calendar__header">
      <button 
        type="button" 
        class="calendar__nav-btn" 
        @click="prevMonth"
        :disabled="!canGoPrev"
      >
        <IconChevronLeft />
      </button>
      <h3 class="calendar__title">{{ currentMonthLabel }}</h3>
      <button 
        type="button" 
        class="calendar__nav-btn" 
        @click="nextMonth"
        :disabled="!canGoNext"
      >
        <IconChevronRight />
      </button>
    </div>

    <!-- 星期标题 -->
    <div class="calendar__weekdays">
      <span 
        v-for="day in weekdays" 
        :key="day" 
        class="calendar__weekday"
        :class="{ 'calendar__weekday--weekend': day === '日' || day === '六' }"
      >
        {{ day }}
      </span>
    </div>

    <!-- 日期网格 -->
    <div class="calendar__grid">
      <button
        v-for="day in days"
        :key="day.key"
        type="button"
        class="calendar__day"
        :class="{
          'calendar__day--other': day.isOtherMonth,
          'calendar__day--today': day.isToday,
          'calendar__day--selected': day.isSelected,
          'calendar__day--has-menu': day.hasMenu
        }"
        @click="selectDate(day.date)"
      >
        <span class="calendar__day-number">{{ day.day }}</span>
        <span v-if="day.hasMenu && !day.isSelected" class="calendar__day-dot"></span>
      </button>
    </div>

    <!-- 图例说明 -->
    <div class="calendar__legend">
      <div class="calendar__legend-item">
        <span class="calendar__legend-dot calendar__legend-dot--menu"></span>
        <span>有菜单</span>
      </div>
      <div class="calendar__legend-item">
        <span class="calendar__legend-dot calendar__legend-dot--today"></span>
        <span>今天</span>
      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * 日历组件
 * 
 * Props:
 * - modelValue: 选中的日期 (YYYY-MM-DD)
 * - menuDates: 有菜单的日期列表
 * 
 * Events:
 * - update:modelValue: 选中日期变化
 * - select: 日期被选中
 */
import { ref, computed, watch } from 'vue'
import IconChevronLeft from './icons/IconChevronLeft.vue'
import IconChevronRight from './icons/IconChevronRight.vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: () => {
      // 直接内联格式化，因为 defineProps 会被提升，不能引用本地函数
      const today = new Date()
      const year = today.getFullYear()
      const month = String(today.getMonth() + 1).padStart(2, '0')
      const day = String(today.getDate()).padStart(2, '0')
      return `${year}-${month}-${day}`
    }
  },
  menuDates: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue', 'select'])

const weekdays = ['日', '一', '二', '三', '四', '五', '六']

// 当前显示的月份
const currentDate = ref(new Date(props.modelValue))
const selectedDate = ref(props.modelValue)

// 今天的日期
const today = new Date()
const todayStr = formatDate(today)

// 月份标签
const currentMonthLabel = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth() + 1
  return `${year}年${month}月`
})

// 是否可以切换上一月（限制到当前月之前的3个月）
const canGoPrev = computed(() => {
  const minDate = new Date(today.getFullYear(), today.getMonth() - 3, 1)
  return currentDate.value > minDate
})

// 是否可以切换下一月（限制到当前月之后的3个月）
const canGoNext = computed(() => {
  const maxDate = new Date(today.getFullYear(), today.getMonth() + 3, 1)
  return currentDate.value < maxDate
})

// 生成日期网格
const days = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth()
  
  // 获取当月第一天和最后一天
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)
  
  // 获取第一天是星期几（0=周日）
  const firstDayWeek = firstDay.getDay()
  
  // 获取上个月的最后几天
  const prevMonthLastDay = new Date(year, month, 0).getDate()
  
  const daysList = []
  
  // 添加上个月的日期
  for (let i = firstDayWeek - 1; i >= 0; i--) {
    const date = new Date(year, month - 1, prevMonthLastDay - i)
    daysList.push(createDayObject(date, true))
  }
  
  // 添加当月的日期
  for (let day = 1; day <= lastDay.getDate(); day++) {
    const date = new Date(year, month, day)
    daysList.push(createDayObject(date, false))
  }
  
  // 添加下个月的日期（填满6行）
  const remainingDays = 42 - daysList.length
  for (let day = 1; day <= remainingDays; day++) {
    const date = new Date(year, month + 1, day)
    daysList.push(createDayObject(date, true))
  }
  
  return daysList
})

// 创建日期对象
function createDayObject(date, isOtherMonth) {
  const dateStr = formatDate(date)
  return {
    key: `${isOtherMonth ? 'other' : 'current'}-${dateStr}`,
    day: date.getDate(),
    date: dateStr,
    isOtherMonth,
    isToday: dateStr === todayStr,
    isSelected: dateStr === selectedDate.value,
    hasMenu: props.menuDates.includes(dateStr)
  }
}

// 格式化日期
function formatDate(date) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

// 上一月
function prevMonth() {
  if (!canGoPrev.value) return
  currentDate.value = new Date(
    currentDate.value.getFullYear(), 
    currentDate.value.getMonth() - 1, 
    1
  )
}

// 下一月
function nextMonth() {
  if (!canGoNext.value) return
  currentDate.value = new Date(
    currentDate.value.getFullYear(), 
    currentDate.value.getMonth() + 1, 
    1
  )
}

// 选择日期
function selectDate(date) {
  selectedDate.value = date
  emit('update:modelValue', date)
  emit('select', date)
}

// 监听 props 变化
watch(() => props.modelValue, (newVal) => {
  if (newVal !== selectedDate.value) {
    selectedDate.value = newVal
    // 如果选中的日期不在当前显示的月份，切换到该月份
    const selectedMonth = new Date(newVal)
    if (
      selectedMonth.getFullYear() !== currentDate.value.getFullYear() ||
      selectedMonth.getMonth() !== currentDate.value.getMonth()
    ) {
      currentDate.value = new Date(selectedMonth.getFullYear(), selectedMonth.getMonth(), 1)
    }
  }
})
</script>

<style scoped>
.calendar {
  padding: var(--space-4);
}

/* 头部 */
.calendar__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-4);
}

.calendar__title {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
  margin: 0;
}

.calendar__nav-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: var(--radius-md);
  background: transparent;
  color: var(--color-text-secondary);
  transition: all var(--transition-fast);
}

.calendar__nav-btn:hover:not(:disabled) {
  background: var(--color-bg-sunken);
  color: var(--color-text-primary);
}

.calendar__nav-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.calendar__nav-btn svg {
  width: 18px;
  height: 18px;
}

/* 星期标题 */
.calendar__weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: var(--space-1);
  margin-bottom: var(--space-2);
}

.calendar__weekday {
  text-align: center;
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-tertiary);
  padding: var(--space-2) 0;
}

.calendar__weekday--weekend {
  color: var(--color-primary-400);
}

/* 日期网格 */
.calendar__grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: var(--space-1);
}

.calendar__day {
  aspect-ratio: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--transition-fast);
  position: relative;
  font-size: var(--font-size-sm);
  color: var(--color-text-primary);
}

.calendar__day:hover {
  background: var(--color-bg-sunken);
}

.calendar__day--other {
  color: var(--color-text-tertiary);
  opacity: 0.5;
}

.calendar__day--today {
  background: var(--color-primary-100);
  color: var(--color-primary-700);
  font-weight: var(--font-weight-semibold);
}

.calendar__day--today:hover {
  background: var(--color-primary-200);
}

.calendar__day--selected {
  background: var(--color-primary) !important;
  color: white !important;
  font-weight: var(--font-weight-semibold);
}

.calendar__day--selected:hover {
  background: var(--color-primary-dark) !important;
}

.calendar__day-number {
  line-height: 1;
}

.calendar__day-dot {
  position: absolute;
  bottom: 4px;
  width: 4px;
  height: 4px;
  border-radius: var(--radius-full);
  background: var(--color-primary);
}

.calendar__day--today .calendar__day-dot {
  background: var(--color-primary-700);
}

/* 图例 */
.calendar__legend {
  display: flex;
  justify-content: center;
  gap: var(--space-5);
  margin-top: var(--space-4);
  padding-top: var(--space-3);
  border-top: 1px solid var(--color-border-light);
}

.calendar__legend-item {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
}

.calendar__legend-dot {
  width: 8px;
  height: 8px;
  border-radius: var(--radius-full);
}

.calendar__legend-dot--menu {
  background: var(--color-primary);
}

.calendar__legend-dot--today {
  background: var(--color-primary-100);
  border: 1px solid var(--color-primary-300);
}
</style>
