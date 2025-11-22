<template>
  <div class="recipe-page page page--dense">
    <section class="card card--surface hero">
      <div>
        <p class="eyebrow">家庭菜谱管理</p>
        <h1>{{ headerTitle }}</h1>
        <p class="subtitle">
          查看、创建与维护家庭菜式，配合基础食材库实现精准录入。
        </p>
      </div>
      <div class="hero-meta">
        <div>
          <p class="label">菜式上限</p>
          <strong>{{ dishUsage }}</strong>
          <small>免费版 30 道，升级会员扩容至 60 道</small>
        </div>
        <div>
          <p class="label">成员数</p>
          <strong>{{ familyStore.memberCount }}</strong>
          <small>同步管理可见权限</small>
        </div>
      </div>
      <div class="hero-actions">
        <button type="button" class="btn btn-primary" @click="openCreate" :disabled="editorLoading">
          {{ editorMode === 'create' && editorVisible ? '录入中...' : '创建菜式' }}
        </button>
        <button type="button" class="btn btn-ghost" @click="loadDishes(true)" :disabled="loading">
          {{ loading ? '刷新中...' : '刷新列表' }}
        </button>
      </div>
    </section>

    <section class="card card--surface filters">
      <div class="filter-group">
        <label>关键字</label>
        <input v-model.trim="filters.keyword" class="form-control" type="text" placeholder="输入菜式名称" />
      </div>
      <div class="filter-group">
        <label>分类</label>
        <button type="button" class="form-control category-select-btn" @click="openCategoryPicker('filter')">
          {{ getCategoryLabel(filters.category) || '全部' }}
        </button>
      </div>
      <div class="filter-actions">
        <button type="button" class="btn btn-primary btn--full" @click="loadDishes(true)">筛选</button>
      </div>
    </section>

    <section class="card card--surface dish-list">
      <div v-if="loading" class="empty">加载菜式列表中...</div>
      <div v-else-if="!dishList.length" class="empty">
        <p>当前还没有菜式，点击上方“创建菜式”开启你的家庭菜谱。</p>
      </div>
      <div v-else class="dish-cards">
        <article v-for="dish in dishList" :key="dish.dish_id" class="dish-card">
          <header>
            <div>
              <p class="label">{{ dish.category || '未分类' }}</p>
              <h3>{{ dish.name }}</h3>
            </div>
            <small>更新于 {{ formatTime(dish.updated_at) }}</small>
          </header>
          <p class="desc">{{ dish.description || '这道菜还没有简介。' }}</p>
          <footer>
            <button type="button" class="btn btn-ghost" @click="openEdit(dish.dish_id)">编辑</button>
            <button type="button" class="btn btn-danger" @click="confirmDelete(dish)">删除</button>
          </footer>
        </article>
      </div>
    </section>

    <section class="card card--surface helper">
      <header>
        <div>
          <p class="eyebrow">食材输入助手</p>
          <h2>调用基础食材接口</h2>
          <p>结合模糊搜索与分类浏览，快速找到标准食材 ID。</p>
        </div>
      </header>
      <IngredientSelector v-model="ingredientPreview" />
      <p class="helper-note">
        此区域为演示用途，所选食材不会提交，仅提供操作手感。正式建菜请使用“创建/编辑菜式”抽屉内的录入组件。
      </p>
    </section>

    <div v-if="editorVisible" class="drawer-wrapper">
      <div class="drawer">
        <header>
          <h2>{{ editorMode === 'create' ? '创建菜式' : '编辑菜式' }}</h2>
          <div class="drawer-header-actions">
            <button type="button" class="btn btn-ghost" @click="closeEditor">关闭</button>
            <button type="button" class="btn btn-primary" @click="submitEditor" :disabled="editorLoading">
              {{ editorLoading ? '保存中...' : '保存' }}
            </button>
          </div>
        </header>
        <form @submit.prevent="submitEditor" class="editor-form">
          <label>
            <span>菜式名称</span>
            <input
              v-model.trim="editorForm.name"
              class="form-control"
              type="text"
              maxlength="100"
              required
            />
          </label>
          <label>
            <span>分类</span>
            <button type="button" class="form-control category-select-btn" @click="openCategoryPicker('editor')">
              {{ getCategoryLabel(editorForm.category) || '未分类' }}
            </button>
          </label>
          <label>
            <span>描述</span>
            <textarea
              v-model.trim="editorForm.description"
              class="form-control"
              maxlength="2000"
              rows="3"
              placeholder="描述这道菜的口味或适用场景"
            />
          </label>
          <div>
            <span class="label">食材列表</span>
            <IngredientSelector v-model="editorForm.ingredients" />
          </div>
          <div class="steps-editor">
            <span class="label">烹饪步骤</span>
            <div class="steps-list">
              <article v-for="(step, index) in editorForm.steps" :key="index" class="step-item">
                <header>
                  <strong>步骤 {{ index + 1 }}</strong>
                  <button type="button" class="btn btn-ghost" @click="removeStep(index)" :disabled="editorForm.steps.length === 1">
                    移除
                  </button>
                </header>
                <textarea
                  v-model.trim="step.content"
                  class="form-control"
                  rows="2"
                  placeholder="描述本步骤（例如：锅中倒油小火煎香）"
                />
              </article>
            </div>
            <button type="button" class="btn btn-ghost" @click="addStep">新增步骤</button>
          </div>
          <p v-if="editorError" class="error-text">{{ editorError }}</p>
          <div class="drawer-actions">
            <button type="submit" class="btn btn-primary" :disabled="editorLoading">
              {{ editorLoading ? '提交中...' : editorMode === 'create' ? '保存菜式' : '更新菜式' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 底部滑入分类选择器 -->
    <div v-if="categoryPickerVisible" class="bottom-sheet-overlay" @click="closeCategoryPicker">
      <div class="bottom-sheet" @click.stop>
        <div class="bottom-sheet-header">
          <h3>选择分类</h3>
          <button type="button" class="btn btn-ghost" @click="closeCategoryPicker">取消</button>
        </div>
        <div class="bottom-sheet-content">
          <button
            type="button"
            class="category-option"
            :class="{ active: currentPickerValue === '' }"
            @click="selectCategory('')"
          >
            <span>{{ categoryPickerMode === 'filter' ? '全部' : '未分类' }}</span>
            <span v-if="currentPickerValue === ''" class="check-icon">✓</span>
          </button>
          <button
            v-for="cat in dishCategories"
            :key="cat.value"
            type="button"
            class="category-option"
            :class="{ active: currentPickerValue === cat.value }"
            @click="selectCategory(cat.value)"
          >
            <span>{{ cat.label }}</span>
            <span v-if="currentPickerValue === cat.value" class="check-icon">✓</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import IngredientSelector from '@/components/IngredientSelector.vue'
import { fetchDishes, createDish, updateDish, deleteDish, getDishDetail } from '@/api/dishes'
import { useFamilyStore } from '@/stores/family'

const familyStore = useFamilyStore()
const router = useRouter()

const filters = reactive({
  keyword: '',
  category: ''
})

const dishCategories = [
  { label: '肉类', value: 'meat' },
  { label: '蔬菜', value: 'vegetable' },
  { label: '汤羹', value: 'soup' },
  { label: '主食', value: 'staple' },
  { label: '甜品', value: 'dessert' },
  { label: '其他', value: 'other' }
]

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const dishList = ref([])
const loading = ref(false)
const ingredientPreview = ref([])
const editorVisible = ref(false)
const editorMode = ref('create')
const editorError = ref('')
const editorLoading = ref(false)
const categoryPickerVisible = ref(false)
const categoryPickerMode = ref('filter') // 'filter' or 'editor'
const currentPickerValue = ref('')
const editorForm = reactive({
  id: null,
  name: '',
  category: '',
  description: '',
  ingredients: [],
  steps: [{ order: 1, content: '' }]
})

const headerTitle = computed(() => {
  return familyStore.familyName ? `${familyStore.familyName} · 菜谱管理` : '家庭菜谱管理'
})

const dishUsage = computed(() => {
  const used = familyStore.familyInfo?.dish_count || 0
  const max = familyStore.familyInfo?.max_dishes || 30
  return `${used} / ${max}`
})

const addStep = () => {
  editorForm.steps.push({ order: editorForm.steps.length + 1, content: '' })
}

const removeStep = (index) => {
  if (editorForm.steps.length === 1) return
  editorForm.steps.splice(index, 1)
}

const resetEditor = () => {
  editorForm.id = null
  editorForm.name = ''
  editorForm.category = ''
  editorForm.description = ''
  editorForm.ingredients = []
  editorForm.steps = [{ order: 1, content: '' }]
  editorError.value = ''
}

const openCreate = () => {
  editorMode.value = 'create'
  resetEditor()
  editorVisible.value = true
}

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
    editorForm.ingredients = (detail.ingredients || []).map((item) => ({
      ingredient_id: item.ingredient_id,
      ingredient_name: item.ingredient_name,
      category: item.category,
      amount: item.amount,
      unit: item.unit || item.default_unit || '份',
      notes: item.notes || ''
    }))
    editorForm.steps = (detail.steps || []).map((step) => ({
      order: step.order,
      content: step.content
    }))
    if (!editorForm.steps.length) {
      editorForm.steps.push({ order: 1, content: '' })
    }
    editorVisible.value = true
  } catch (error) {
    editorError.value = error.message || '加载菜式详情失败'
  }
}

const closeEditor = () => {
  if (editorLoading.value) return
  editorVisible.value = false
}

const normalizeIngredients = () => {
  return editorForm.ingredients
    .filter((item) => item.ingredient_id)
    .map((item, index) => ({
      ingredient_id: item.ingredient_id,
      amount: Number(item.amount) || 0,
      unit: item.unit || '份',
      notes: item.notes || '',
      sort_order: index + 1
    }))
}

const normalizeSteps = () => {
  return editorForm.steps
    .filter((step) => step.content)
    .map((step, index) => ({
      order: index + 1,
      content: step.content
    }))
}

const submitEditor = async () => {
  editorError.value = ''
  if (!editorForm.name.trim()) {
    editorError.value = '请填写菜式名称'
    return
  }
  const ingredients = normalizeIngredients()
  if (!ingredients.length) {
    editorError.value = '请至少添加一个食材'
    return
  }
  const steps = normalizeSteps()
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
    console.error(error)
  } finally {
    loading.value = false
  }
}

const formatTime = (value) => {
  if (!value) return '刚刚'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '刚刚'
  return date.toLocaleDateString()
}

const getCategoryLabel = (value) => {
  if (!value) return ''
  const cat = dishCategories.find((c) => c.value === value)
  return cat ? cat.label : ''
}

const openCategoryPicker = (mode) => {
  categoryPickerMode.value = mode
  currentPickerValue.value = mode === 'filter' ? filters.category : editorForm.category
  categoryPickerVisible.value = true
}

const closeCategoryPicker = () => {
  categoryPickerVisible.value = false
}

const selectCategory = (value) => {
  if (categoryPickerMode.value === 'filter') {
    filters.category = value
  } else {
    editorForm.category = value
  }
  closeCategoryPicker()
}

onMounted(async () => {
  if (!familyStore.familyInfo) {
    await familyStore.fetchFamilyInfo()
  }
  await loadDishes()
})
</script>

<style scoped>
.hero {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.hero-meta {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 12px;
}

.hero-meta .label {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.hero-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.filters {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
  align-items: end;
}

.filter-group label {
  display: block;
  font-size: 13px;
  margin-bottom: 4px;
}

.filter-actions {
  display: flex;
  justify-content: flex-end;
}

.dish-list .dish-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 16px;
}

.dish-card {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-medium);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.dish-card header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dish-card footer {
  display: flex;
  gap: 10px;
}

.dish-card footer button {
  flex: 1;
}

.empty {
  text-align: center;
  color: var(--color-text-secondary);
}

.helper-note {
  margin-top: 12px;
  font-size: 13px;
  color: var(--color-text-secondary);
}

.drawer-wrapper {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  justify-content: flex-end;
  z-index: 20;
}

.drawer {
  width: min(480px, 100%);
  background: #fff;
  height: 100vh;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.drawer header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.editor-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
  overflow-y: auto;
}

.steps-editor .steps-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.step-item {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-medium);
  padding: 8px;
}

.drawer-actions {
  display: flex;
  justify-content: flex-end;
}

.drawer-header-actions {
  display: flex;
  gap: 8px;
}

.drawer-actions button {
  min-width: 120px;
}

.error-text {
  color: #f44336;
  font-size: 13px;
}

.category-select-btn {
  text-align: left;
  cursor: pointer;
  background: #fff;
  border: 1px solid var(--color-border);
  padding: 8px 12px;
  border-radius: var(--radius-medium);
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
}

.category-select-btn:active {
  background: #f5f5f5;
}

/* 底部滑入选择器 */
.bottom-sheet-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 100;
  display: flex;
  align-items: flex-end;
  animation: fadeIn 0.2s ease-out;
}

.bottom-sheet {
  width: 100%;
  max-height: 70vh;
  background: #fff;
  border-radius: 16px 16px 0 0;
  display: flex;
  flex-direction: column;
  animation: slideUp 0.3s ease-out;
  overflow: hidden;
}

.bottom-sheet-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--color-border);
}

.bottom-sheet-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.bottom-sheet-content {
  padding: 8px 0;
  overflow-y: auto;
  max-height: calc(70vh - 60px);
}

.category-option {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: #fff;
  border: none;
  text-align: left;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.category-option:active {
  background: #f5f5f5;
}

.category-option.active {
  color: var(--color-primary, #1976d2);
  font-weight: 500;
}

.check-icon {
  color: var(--color-primary, #1976d2);
  font-weight: bold;
  font-size: 18px;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    transform: translateY(100%);
  }
  to {
    transform: translateY(0);
  }
}
</style>
