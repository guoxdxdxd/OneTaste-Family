<template>
  <div class="ingredient-selector">
    <!-- 搜索区域 -->
    <div class="selector-section">
      <div class="selector-section__header">
        <span class="selector-section__label">关键字搜索</span>
        <h3 class="selector-section__title">快速找到常用食材</h3>
      </div>
      <div class="search-input-wrapper">
        <input
          v-model="searchKeyword"
          type="text"
          :disabled="loading.search"
          class="input"
          placeholder="输入食材名称或拼音，例如：五花肉"
        />
      </div>
      <ul class="ingredient-list">
        <li v-if="loading.search" class="ingredient-list__loading">搜索中...</li>
        <li v-else-if="searchKeyword && !searchResults.length" class="ingredient-list__empty">
          暂无搜索结果
        </li>
        <li
          v-for="item in searchResults"
          :key="item.ingredient_id"
          class="ingredient-item"
          :class="{ 'ingredient-item--selected': isSelected(item.ingredient_id) }"
          @click="toggleIngredient(item)"
        >
          <div class="ingredient-item__info">
            <span class="ingredient-item__name">{{ item.name }}</span>
            <span class="ingredient-item__meta">{{ item.category || '未分类' }} · {{ item.default_unit || '份' }}</span>
          </div>
          <span class="ingredient-item__action">
            {{ isSelected(item.ingredient_id) ? '已选' : '选择' }}
          </span>
        </li>
      </ul>
    </div>

    <!-- 分类浏览区域 -->
    <div class="selector-section">
      <div class="selector-section__header">
        <span class="selector-section__label">按分类浏览</span>
        <h3 class="selector-section__title">灵感库</h3>
      </div>
      
      <!-- 分类标签选择 -->
      <div class="category-tabs">
        <button
          v-for="cat in categoryOptions"
          :key="cat.value"
          type="button"
          class="category-tab"
          :class="{ 'category-tab--active': categoryQuery.category === cat.value }"
          @click="selectCategory(cat.value)"
        >
          {{ cat.label }}
        </button>
      </div>

      <!-- 分类内搜索 -->
      <div class="search-input-wrapper search-input-wrapper--sm">
        <input
          v-model="categoryQuery.keyword"
          type="text"
          class="input"
          placeholder="分类内搜索"
        />
      </div>

      <ul class="ingredient-list">
        <li v-if="loading.category" class="ingredient-list__loading">载入中...</li>
        <li v-else-if="!categoryResults.length" class="ingredient-list__empty">
          暂无匹配食材
        </li>
        <li
          v-for="item in categoryResults"
          :key="`cat-${item.ingredient_id}`"
          class="ingredient-item"
          :class="{ 'ingredient-item--selected': isSelected(item.ingredient_id) }"
          @click="toggleIngredient(item)"
        >
          <div class="ingredient-item__info">
            <span class="ingredient-item__name">{{ item.name }}</span>
            <span class="ingredient-item__meta">{{ item.default_unit || '份' }}</span>
          </div>
          <span class="ingredient-item__action">
            {{ isSelected(item.ingredient_id) ? '已选' : '选择' }}
          </span>
        </li>
      </ul>

      <!-- 分页 -->
      <div class="pagination" v-if="categoryPaging.total > categoryQuery.pageSize">
        <button
          type="button"
          class="btn btn--ghost btn--sm"
          :disabled="categoryPaging.page === 1"
          @click="changeCategoryPage(-1)"
        >
          上一页
        </button>
        <span class="pagination__info">{{ categoryPaging.page }} / {{ totalCategoryPages }}</span>
        <button
          type="button"
          class="btn btn--ghost btn--sm"
          :disabled="categoryPaging.page >= totalCategoryPages"
          @click="changeCategoryPage(1)"
        >
          下一页
        </button>
      </div>
    </div>

    <!-- 已选食材区域 -->
    <div class="selector-section selector-section--selected">
      <div class="selector-section__header">
        <span class="selector-section__label">已选食材</span>
        <h3 class="selector-section__title">{{ selectedList.length }} 项</h3>
      </div>
      
      <p v-if="!selectedList.length" class="selected-empty">
        尚未选择食材，从上方列表点击选择
      </p>
      
      <ul v-else class="selected-list">
        <li v-for="(item, index) in selectedList" :key="item.ingredient_id" class="selected-item">
          <div class="selected-item__header">
            <span class="selected-item__name">{{ item.ingredient_name || item.name }}</span>
            <button type="button" class="btn btn--text btn--sm" @click="removeIngredient(index)">
              移除
            </button>
          </div>
          <div class="selected-item__inputs">
            <div class="form-group form-group--inline">
              <label class="form-label">数量</label>
              <input
                type="number"
                min="0"
                step="0.1"
                class="input input--sm"
                v-model.number="item.amount"
                @change="emitChange()"
              />
            </div>
            <div class="form-group form-group--inline">
              <label class="form-label">单位</label>
              <input
                type="text"
                maxlength="10"
                class="input input--sm"
                v-model="item.unit"
                @change="emitChange()"
              />
            </div>
            <div class="form-group form-group--inline form-group--full">
              <label class="form-label">备注</label>
              <input
                type="text"
                maxlength="50"
                class="input input--sm"
                v-model="item.notes"
                @change="emitChange()"
                placeholder="例如：切片/去籽"
              />
            </div>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
/**
 * 食材选择器组件
 * 
 * 功能：
 * - 关键字搜索食材
 * - 按分类浏览食材
 * - 选择/取消选择食材
 * - 编辑已选食材的数量、单位、备注
 */
import { computed, reactive, ref, watch } from 'vue'
import { searchIngredients, fetchIngredientsByCategory } from '@/api/ingredients'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue'])

// 深拷贝工具
const clone = (value) => {
  try {
    return JSON.parse(JSON.stringify(value || []))
  } catch (error) {
    return Array.isArray(value) ? [...value] : []
  }
}

// 状态
const selectedList = ref([])
const searchKeyword = ref('')
const searchResults = ref([])
const categoryResults = ref([])

const categoryQuery = reactive({
  category: 'meat',
  keyword: ''
})

const categoryPaging = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const loading = reactive({
  search: false,
  category: false
})

// 分类选项
const categoryOptions = [
  { label: '肉类', value: 'meat' },
  { label: '蔬菜', value: 'vegetable' },
  { label: '海鲜', value: 'seafood' },
  { label: '蛋白', value: 'protein' },
  { label: '调味', value: 'condiment' },
  { label: '谷物', value: 'grain' },
  { label: '其他', value: 'other' }
]

// 计算总页数
const totalCategoryPages = computed(() => {
  return Math.max(1, Math.ceil(categoryPaging.total / categoryPaging.pageSize))
})

// 检查食材是否已选择
const isSelected = (ingredientId) => {
  return selectedList.value.some(item => item.ingredient_id === ingredientId)
}

// 同步外部数据
watch(
  () => props.modelValue,
  (val) => {
    selectedList.value = clone(val)
  },
  { deep: true, immediate: true }
)

// 搜索关键字变化时触发搜索
watch(searchKeyword, (val) => {
  debouncedSearch(val)
})

// 分类或关键字变化时重新加载
watch(
  () => [categoryQuery.category, categoryQuery.keyword, categoryPaging.page],
  () => {
    loadCategory()
  }
)

// 触发数据变更
const emitChange = () => {
  emit('update:modelValue', clone(selectedList.value))
}

// 防抖搜索
const debouncedSearch = (() => {
  let timer
  return (keyword) => {
    clearTimeout(timer)
    timer = setTimeout(() => {
      if (!keyword) {
        searchResults.value = []
        return
      }
      fetchSearch(keyword)
    }, 400)
  }
})()

// 执行搜索
const fetchSearch = async (keyword) => {
  try {
    loading.search = true
    const res = await searchIngredients({ keyword, limit: 10 })
    searchResults.value = res?.data?.items || res?.data || []
  } catch (error) {
    searchResults.value = []
    console.error(error)
  } finally {
    loading.search = false
  }
}

// 加载分类数据
const loadCategory = async () => {
  try {
    loading.category = true
    const params = {
      category: categoryQuery.category,
      keyword: categoryQuery.keyword || undefined,
      page: categoryPaging.page,
      page_size: categoryPaging.pageSize
    }
    const res = await fetchIngredientsByCategory(params)
    const payload = res?.data || {}
    categoryResults.value = payload.items || []
    categoryPaging.total = payload.total || categoryResults.value.length
  } catch (error) {
    categoryResults.value = []
    categoryPaging.total = 0
    console.error(error)
  } finally {
    loading.category = false
  }
}

// 选择分类
const selectCategory = (category) => {
  categoryQuery.category = category
  categoryPaging.page = 1
}

// 切换分页
const changeCategoryPage = (delta) => {
  const next = categoryPaging.page + delta
  if (next < 1 || next > totalCategoryPages.value) return
  categoryPaging.page = next
}

// 切换食材选择（选择/取消选择）
const toggleIngredient = (item) => {
  if (!item?.ingredient_id) return
  
  const index = selectedList.value.findIndex(
    target => target.ingredient_id === item.ingredient_id
  )
  
  if (index > -1) {
    // 已存在，取消选择
    selectedList.value.splice(index, 1)
  } else {
    // 不存在，添加选择
    selectedList.value.push({
      ingredient_id: item.ingredient_id,
      ingredient_name: item.name,
      category: item.category,
      amount: 1,
      unit: item.default_unit || '份',
      notes: ''
    })
  }
  emitChange()
}

// 移除已选食材
const removeIngredient = (index) => {
  selectedList.value.splice(index, 1)
  emitChange()
}

// 初始化加载分类数据
loadCategory()
</script>

<style scoped>
.ingredient-selector {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

/* 区块样式 */
.selector-section {
  background: var(--color-bg-sunken);
  border-radius: var(--radius-lg);
  padding: var(--space-4);
}

.selector-section--selected {
  background: var(--color-primary-50);
  border: 1px solid var(--color-primary-100);
}

.selector-section__header {
  margin-bottom: var(--space-3);
}

.selector-section__label {
  display: block;
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: var(--space-1);
}

.selector-section__title {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
  margin: 0;
}

/* 搜索输入框 */
.search-input-wrapper {
  margin-bottom: var(--space-3);
}

.search-input-wrapper--sm {
  margin-top: var(--space-3);
}

/* 分类标签 */
.category-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-2);
}

.category-tab {
  height: 32px;
  padding: 0 var(--space-3);
  background: var(--color-bg-elevated);
  border: 2px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.category-tab:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.category-tab--active {
  background: var(--color-primary-100);
  border-color: var(--color-primary);
  color: var(--color-primary-700);
  font-weight: var(--font-weight-semibold);
}

/* 食材列表 */
.ingredient-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  max-height: 240px;
  overflow-y: auto;
}

.ingredient-list__loading,
.ingredient-list__empty {
  padding: var(--space-4);
  text-align: center;
  color: var(--color-text-tertiary);
  font-size: var(--font-size-sm);
}

.ingredient-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-3);
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.ingredient-item:hover {
  border-color: var(--color-primary-300);
}

.ingredient-item--selected {
  background: var(--color-primary-100);
  border-color: var(--color-primary);
}

.ingredient-item__info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.ingredient-item__name {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-heading);
}

.ingredient-item__meta {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
}

.ingredient-item__action {
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
  color: var(--color-primary);
  padding: var(--space-1) var(--space-2);
  background: var(--color-primary-100);
  border-radius: var(--radius-sm);
}

.ingredient-item--selected .ingredient-item__action {
  background: var(--color-primary);
  color: white;
}

/* 分页 */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-3);
  margin-top: var(--space-3);
  padding-top: var(--space-3);
  border-top: 1px solid var(--color-border-light);
}

.pagination__info {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

/* 已选食材区域 */
.selected-empty {
  text-align: center;
  padding: var(--space-4);
  color: var(--color-text-tertiary);
  font-size: var(--font-size-sm);
  margin: 0;
}

.selected-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.selected-item {
  background: var(--color-bg-elevated);
  border-radius: var(--radius-lg);
  padding: var(--space-3);
}

.selected-item__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-3);
}

.selected-item__name {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-heading);
}

.selected-item__inputs {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-2);
}

.form-group--inline {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.form-group--inline .form-label {
  font-size: 11px;
  color: var(--color-text-tertiary);
}

.form-group--full {
  grid-column: 1 / -1;
}

.input--sm {
  height: 36px;
  font-size: var(--font-size-sm);
}
</style>
