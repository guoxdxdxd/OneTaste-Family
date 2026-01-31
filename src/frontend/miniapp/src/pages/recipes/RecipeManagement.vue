<template>
  <div class="page recipe-page">
    <!-- 页面头部 -->
    <header class="recipe-header">
      <div class="recipe-header__info">
        <h1 class="recipe-header__title">家庭菜谱</h1>
        <p class="recipe-header__subtitle">{{ familyStore.familyName || '我的菜谱' }}</p>
      </div>
      <div class="recipe-header__stats">
        <span class="recipe-header__count">{{ dishCount }} / {{ maxDishes }}</span>
        <span class="recipe-header__label">菜式</span>
      </div>
    </header>

    <!-- 搜索与筛选 -->
    <section class="filter-section">
      <div class="search-box">
        <IconSearch class="search-box__icon" />
        <input
          v-model.trim="filters.keyword"
          type="text"
          class="search-box__input"
          placeholder="搜索菜式名称"
          @keyup.enter="loadDishes(true)"
        />
      </div>
      <div class="filter-tags">
        <button 
          v-for="cat in categoryOptions" 
          :key="cat.value"
          type="button"
          class="filter-tag"
          :class="{ 'filter-tag--active': filters.category === cat.value }"
          @click="setCategory(cat.value)"
        >
          {{ cat.label }}
        </button>
      </div>
    </section>

    <!-- 菜式列表 -->
    <section class="dish-section">
      <!-- 加载状态 -->
      <div v-if="loading" class="loading-state">
        <span class="loading-spinner"></span>
        <span>加载中...</span>
      </div>

      <!-- 空状态 -->
      <div v-else-if="!dishList.length" class="empty-state">
        <div class="empty-state__icon">📖</div>
        <h3 class="empty-state__title">还没有菜式</h3>
        <p class="empty-state__description">点击下方按钮创建第一道菜式</p>
        <button class="btn btn--primary" @click="openCreate">
          <IconPlus class="btn__icon" />
          创建菜式
        </button>
      </div>

      <!-- 菜式卡片网格 -->
      <div v-else class="dish-grid">
        <article 
          v-for="dish in dishList" 
          :key="dish.dish_id"
          class="dish-card card card--interactive"
          @click="openEdit(dish.dish_id)"
        >
          <div class="dish-card__header">
            <span class="tag tag--primary tag--pill">{{ getCategoryLabel(dish.category) || '未分类' }}</span>
            <span class="dish-card__date">{{ formatTime(dish.updated_at) }}</span>
          </div>
          <h3 class="dish-card__title">{{ dish.name }}</h3>
          <p class="dish-card__desc">{{ dish.description || '暂无描述' }}</p>
          <div class="dish-card__footer">
            <button 
              class="btn btn--text btn--sm" 
              @click.stop="confirmDelete(dish)"
            >
              删除
            </button>
          </div>
        </article>
      </div>
    </section>

    <!-- 悬浮添加按钮 -->
    <button class="fab" @click="openCreate">
      <IconPlus class="fab__icon" />
    </button>

    <!-- 编辑抽屉 -->
    <transition name="drawer">
      <div v-if="editorVisible" class="drawer-overlay" @click.self="closeEditor">
        <div class="drawer">
          <div class="drawer__header">
            <h2 class="drawer__title">
              {{ editorMode === 'create' ? '创建菜式' : '编辑菜式' }}
            </h2>
            <button class="icon-btn" @click="closeEditor">
              <IconClose />
            </button>
          </div>

          <form @submit.prevent="submitEditor" class="drawer__body">
            <!-- 菜式名称 -->
            <div class="form-group">
              <label class="form-label form-label--required">菜式名称</label>
              <input
                v-model.trim="editorForm.name"
                type="text"
                class="input"
                placeholder="输入菜式名称"
                maxlength="100"
                required
              />
            </div>

            <!-- 分类选择 -->
            <div class="form-group">
              <label class="form-label">分类</label>
              <div class="category-picker">
                <button
                  v-for="cat in dishCategories"
                  :key="cat.value"
                  type="button"
                  class="category-picker__item"
                  :class="{ 'category-picker__item--active': editorForm.category === cat.value }"
                  @click="editorForm.category = cat.value"
                >
                  {{ cat.label }}
                </button>
              </div>
            </div>

            <!-- 描述 -->
            <div class="form-group">
              <label class="form-label">描述</label>
              <textarea
                v-model.trim="editorForm.description"
                class="textarea"
                placeholder="描述这道菜的口味或特点"
                maxlength="2000"
                rows="2"
              />
            </div>

            <!-- 食材 -->
            <div class="form-group">
              <label class="form-label form-label--required">食材</label>
              <IngredientSelector v-model="editorForm.ingredients" />
            </div>

            <!-- 烹饪步骤 -->
            <div class="form-group">
              <label class="form-label form-label--required">烹饪步骤</label>
              <div class="steps-editor">
                <div 
                  v-for="(step, index) in editorForm.steps" 
                  :key="index"
                  class="step-item"
                >
                  <div class="step-item__header">
                    <span class="step-item__number">{{ index + 1 }}</span>
                    <button 
                      v-if="editorForm.steps.length > 1"
                      type="button" 
                      class="btn btn--text btn--sm"
                      @click="removeStep(index)"
                    >
                      移除
                    </button>
                  </div>
                  <textarea
                    v-model.trim="step.content"
                    class="textarea"
                    placeholder="描述这一步骤..."
                    rows="2"
                  />
                </div>
                <button type="button" class="btn btn--ghost btn--sm btn--full" @click="addStep">
                  <IconPlus class="btn__icon" />
                  添加步骤
                </button>
              </div>
            </div>

            <!-- 错误提示 -->
            <div v-if="editorError" class="form-error-alert">
              {{ editorError }}
            </div>
          </form>

          <div class="drawer__footer">
            <button type="button" class="btn btn--ghost" @click="closeEditor">
              取消
            </button>
            <button 
              type="button" 
              class="btn btn--primary" 
              @click="submitEditor"
              :disabled="editorLoading"
            >
              {{ editorLoading ? '保存中...' : '保存' }}
            </button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
/**
 * 菜谱管理页面
 * 
 * 功能：
 * - 菜式列表展示与搜索
 * - 创建/编辑菜式
 * - 删除菜式
 */
import { computed, onMounted, reactive, ref } from 'vue'
import IngredientSelector from '@/components/IngredientSelector.vue'
import { fetchDishes, createDish, updateDish, deleteDish, getDishDetail } from '@/api/dishes'
import { useFamilyStore } from '@/stores/family'
import IconPlus from '@/components/icons/IconPlus.vue'
import IconClose from '@/components/icons/IconClose.vue'

// 搜索图标组件
const IconSearch = {
  template: `
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <circle cx="11" cy="11" r="8"/>
      <line x1="21" y1="21" x2="16.65" y2="16.65"/>
    </svg>
  `
}

const familyStore = useFamilyStore()

// 筛选
const filters = reactive({
  keyword: '',
  category: ''
})

// 分类选项
const dishCategories = [
  { label: '肉类', value: 'meat' },
  { label: '蔬菜', value: 'vegetable' },
  { label: '汤羹', value: 'soup' },
  { label: '主食', value: 'staple' },
  { label: '甜品', value: 'dessert' },
  { label: '其他', value: 'other' }
]

const categoryOptions = computed(() => [
  { label: '全部', value: '' },
  ...dishCategories
])

// 菜式数据
const dishList = ref([])
const loading = ref(false)
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 编辑器状态
const editorVisible = ref(false)
const editorMode = ref('create')
const editorLoading = ref(false)
const editorError = ref('')
const editorForm = reactive({
  id: null,
  name: '',
  category: '',
  description: '',
  ingredients: [],
  steps: [{ order: 1, content: '' }]
})

// 统计数据
const dishCount = computed(() => familyStore.familyInfo?.dish_count || dishList.value.length)
const maxDishes = computed(() => familyStore.familyInfo?.max_dishes || 30)

// 获取分类标签
const getCategoryLabel = (value) => {
  const cat = dishCategories.find(c => c.value === value)
  return cat ? cat.label : ''
}

// 格式化时间
const formatTime = (value) => {
  if (!value) return '刚刚'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '刚刚'
  const now = new Date()
  const diff = now - date
  if (diff < 86400000) return '今天'
  if (diff < 172800000) return '昨天'
  return `${date.getMonth() + 1}/${date.getDate()}`
}

// 设置分类筛选
const setCategory = (value) => {
  filters.category = value
  loadDishes(true)
}

// 加载菜式列表
const loadDishes = async (resetPage = false) => {
  if (resetPage) {
    pagination.page = 1
  }
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      category: filters.category || undefined,
      keyword: filters.keyword || undefined
    }
    const res = await fetchDishes(params)
    const payload = res?.data || {}
    dishList.value = payload.dishes || []
    pagination.total = payload.total || dishList.value.length
  } catch (error) {
    console.error('加载菜式失败:', error)
  } finally {
    loading.value = false
  }
}

// 重置编辑器
const resetEditor = () => {
  editorForm.id = null
  editorForm.name = ''
  editorForm.category = ''
  editorForm.description = ''
  editorForm.ingredients = []
  editorForm.steps = [{ order: 1, content: '' }]
  editorError.value = ''
}

// 打开创建
const openCreate = () => {
  editorMode.value = 'create'
  resetEditor()
  editorVisible.value = true
}

// 打开编辑
const openEdit = async (id) => {
  try {
    editorMode.value = 'update'
    editorError.value = ''
    const res = await getDishDetail(id)
    const detail = res?.data
    if (!detail) return
    
    editorForm.id = detail.dish_id
    editorForm.name = detail.name
    editorForm.category = detail.category || ''
    editorForm.description = detail.description || ''
    editorForm.ingredients = (detail.ingredients || []).map(item => ({
      ingredient_id: item.ingredient_id,
      ingredient_name: item.ingredient_name,
      category: item.category,
      amount: item.amount,
      unit: item.unit || item.default_unit || '份',
      notes: item.notes || ''
    }))
    editorForm.steps = (detail.steps || []).map(step => ({
      order: step.order,
      content: step.content
    }))
    if (!editorForm.steps.length) {
      editorForm.steps.push({ order: 1, content: '' })
    }
    editorVisible.value = true
  } catch (error) {
    console.error('加载菜式详情失败:', error)
  }
}

// 关闭编辑器
const closeEditor = () => {
  if (editorLoading.value) return
  editorVisible.value = false
}

// 添加步骤
const addStep = () => {
  editorForm.steps.push({ order: editorForm.steps.length + 1, content: '' })
}

// 移除步骤
const removeStep = (index) => {
  if (editorForm.steps.length === 1) return
  editorForm.steps.splice(index, 1)
}

// 提交编辑器
const submitEditor = async () => {
  editorError.value = ''
  
  if (!editorForm.name.trim()) {
    editorError.value = '请填写菜式名称'
    return
  }
  
  const ingredients = editorForm.ingredients
    .filter(item => item.ingredient_id)
    .map((item, index) => ({
      ingredient_id: item.ingredient_id,
      amount: Number(item.amount) || 0,
      unit: item.unit || '份',
      notes: item.notes || '',
      sort_order: index + 1
    }))
  
  if (!ingredients.length) {
    editorError.value = '请至少添加一个食材'
    return
  }
  
  const steps = editorForm.steps
    .filter(step => step.content)
    .map((step, index) => ({
      order: index + 1,
      content: step.content
    }))
  
  if (!steps.length) {
    editorError.value = '请至少填写一个步骤'
    return
  }
  
  const payload = {
    name: editorForm.name,
    category: editorForm.category,
    description: editorForm.description,
    ingredients,
    steps
  }
  
  editorLoading.value = true
  try {
    if (editorMode.value === 'create') {
      await createDish(payload)
    } else {
      await updateDish(editorForm.id, payload)
    }
    await loadDishes(true)
    await familyStore.fetchFamilyInfo(true)
    editorVisible.value = false
  } catch (error) {
    editorError.value = error.message || '保存失败'
  } finally {
    editorLoading.value = false
  }
}

// 确认删除
const confirmDelete = async (dish) => {
  if (!window.confirm(`确定要删除「${dish.name}」吗？`)) return
  try {
    await deleteDish(dish.dish_id)
    await loadDishes(true)
    await familyStore.fetchFamilyInfo(true)
  } catch (error) {
    window.alert(error.message || '删除失败')
  }
}

// 初始化
onMounted(async () => {
  if (!familyStore.familyInfo) {
    await familyStore.fetchFamilyInfo()
  }
  await loadDishes()
})
</script>

<style scoped>
.recipe-page {
  padding-top: var(--space-4);
  padding-bottom: calc(var(--tabbar-height) + var(--safe-area-bottom) + 80px);
}

/* 页面头部 */
.recipe-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--space-5);
}

.recipe-header__title {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-text-heading);
  margin: 0 0 var(--space-1);
}

.recipe-header__subtitle {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0;
}

.recipe-header__stats {
  text-align: right;
}

.recipe-header__count {
  display: block;
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-primary);
}

.recipe-header__label {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
}

/* 搜索与筛选 */
.filter-section {
  margin-bottom: var(--space-5);
}

.search-box {
  position: relative;
  margin-bottom: var(--space-3);
}

.search-box__icon {
  position: absolute;
  left: var(--space-4);
  top: 50%;
  transform: translateY(-50%);
  width: 18px;
  height: 18px;
  color: var(--color-text-tertiary);
}

.search-box__input {
  width: 100%;
  height: 44px;
  padding: 0 var(--space-4) 0 calc(var(--space-4) + 24px);
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  font-size: var(--font-size-sm);
}

.search-box__input:focus {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(240, 138, 93, 0.1);
}

.filter-tags {
  display: flex;
  gap: var(--space-2);
  overflow-x: auto;
  padding-bottom: var(--space-1);
  -webkit-overflow-scrolling: touch;
}

.filter-tags::-webkit-scrollbar {
  display: none;
}

.filter-tag {
  flex-shrink: 0;
  height: 32px;
  padding: 0 var(--space-4);
  background: var(--color-bg-sunken);
  border: none;
  border-radius: var(--radius-full);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  transition: all var(--transition-fast);
}

.filter-tag--active {
  background: var(--color-primary);
  color: white;
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
  margin: 0 0 var(--space-5);
}

/* 菜式网格 */
.dish-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: var(--space-3);
}

.dish-card {
  padding: var(--space-4);
}

.dish-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-3);
}

.dish-card__date {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
}

.dish-card__title {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
  margin: 0 0 var(--space-2);
}

.dish-card__desc {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.dish-card__footer {
  margin-top: var(--space-3);
  padding-top: var(--space-3);
  border-top: 1px solid var(--color-border-light);
  text-align: right;
}

.dish-card__footer .btn--text {
  color: var(--color-danger-500);
}

/* 悬浮按钮 */
.fab {
  position: fixed;
  right: var(--space-5);
  bottom: calc(var(--tabbar-height) + var(--safe-area-bottom) + var(--space-5));
}

/* 抽屉 */
.drawer-overlay {
  position: fixed;
  inset: 0;
  background: var(--color-bg-overlay);
  z-index: var(--z-modal-backdrop);
  display: flex;
  justify-content: flex-end;
}

.drawer {
  width: 100%;
  max-width: 480px;
  height: 100%;
  background: var(--color-bg-elevated);
  display: flex;
  flex-direction: column;
  animation: slideInRight var(--duration-normal) var(--ease-out);
}

.drawer__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-4) var(--space-5);
  border-bottom: 1px solid var(--color-border-light);
  flex-shrink: 0;
}

.drawer__title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
  margin: 0;
}

.drawer__body {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-5);
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.drawer__footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  padding: var(--space-4) var(--space-5);
  border-top: 1px solid var(--color-border-light);
  flex-shrink: 0;
}

/* 分类选择器 */
.category-picker {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-2);
}

.category-picker__item {
  height: 32px;
  padding: 0 var(--space-3);
  background: var(--color-bg-sunken);
  border: 1px solid transparent;
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  transition: all var(--transition-fast);
}

.category-picker__item--active {
  background: var(--color-primary-100);
  border-color: var(--color-primary);
  color: var(--color-primary);
}

/* 步骤编辑器 */
.steps-editor {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.step-item {
  background: var(--color-bg-sunken);
  border-radius: var(--radius-lg);
  padding: var(--space-3);
}

.step-item__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-2);
}

.step-item__number {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-primary);
  color: white;
  border-radius: var(--radius-full);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
}

/* 错误提示 */
.form-error-alert {
  padding: var(--space-3) var(--space-4);
  background: var(--color-danger-50);
  border: 1px solid var(--color-danger-100);
  border-radius: var(--radius-lg);
  color: var(--color-danger-600);
  font-size: var(--font-size-sm);
}

/* 动画 */
.drawer-enter-active,
.drawer-leave-active {
  transition: opacity var(--duration-normal) var(--ease-out);
}

.drawer-enter-active .drawer,
.drawer-leave-active .drawer {
  transition: transform var(--duration-normal) var(--ease-out);
}

.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;
}

.drawer-enter-from .drawer,
.drawer-leave-to .drawer {
  transform: translateX(100%);
}
</style>
