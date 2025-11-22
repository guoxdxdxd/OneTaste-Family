<template>
  <div class="ingredient-selector">
    <div class="selector-column">
      <div class="section-header">
        <div>
          <p class="eyebrow">关键字搜索</p>
          <h3>快速找到常用食材</h3>
        </div>
      </div>
      <div class="search-box">
        <input
          v-model="searchKeyword"
          type="text"
          :disabled="loading.search"
          class="form-control"
          placeholder="输入食材名称或拼音，例如：五花肉"
        />
      </div>
      <ul class="result-list">
        <li v-if="loading.search">搜索中...</li>
        <li v-else-if="searchKeyword && !searchResults.length" class="empty">暂无搜索结果</li>
        <li
          v-for="item in searchResults"
          :key="item.ingredient_id"
          class="result-item"
        >
          <div>
            <strong>{{ item.name }}</strong>
            <small>{{ item.category || '未分类' }} · 默认 {{ item.default_unit || '份' }}</small>
          </div>
          <button type="button" class="btn btn-primary btn--sm" @click="addIngredient(item)">
            添加
          </button>
        </li>
      </ul>
    </div>

    <div class="selector-column">
      <div class="section-header">
        <div>
          <p class="eyebrow">按分类浏览</p>
          <h3>灵感库</h3>
        </div>
        <div class="category-filters">
          <select v-model="categoryQuery.category" class="form-control">
            <option v-for="item in categoryOptions" :key="item.value" :value="item.value">
              {{ item.label }}
            </option>
          </select>
          <input
            v-model="categoryQuery.keyword"
            type="text"
            class="form-control"
            placeholder="分类内搜索"
          />
        </div>
      </div>
      <ul class="result-list">
        <li v-if="loading.category">载入分类...</li>
        <li v-else-if="!categoryResults.length" class="empty">分类中暂无匹配食材</li>
        <li
          v-for="item in categoryResults"
          :key="`cat-${item.ingredient_id}`"
          class="result-item"
        >
          <div>
            <strong>{{ item.name }}</strong>
            <small>{{ item.category || '未分类' }} · {{ item.default_unit || '份' }}</small>
          </div>
          <button type="button" class="btn btn-primary btn--sm" @click="addIngredient(item)">
            选择
          </button>
        </li>
      </ul>
      <div class="pagination" v-if="categoryPaging.total > categoryQuery.pageSize">
        <button
          type="button"
          class="btn btn-ghost btn--sm"
          :disabled="categoryPaging.page === 1"
          @click="changeCategoryPage(-1)"
        >
          上一页
        </button>
        <span>{{ categoryPaging.page }} / {{ totalCategoryPages }}</span>
        <button
          type="button"
          class="btn btn-ghost btn--sm"
          :disabled="categoryPaging.page >= totalCategoryPages"
          @click="changeCategoryPage(1)"
        >
          下一页
        </button>
      </div>
    </div>

    <div class="selector-column selected-column">
      <div class="section-header">
        <div>
          <p class="eyebrow">已选食材</p>
          <h3>{{ selectedList.length }} 项</h3>
        </div>
      </div>
      <p v-if="!selectedList.length" class="empty">尚未选择食材，先从左侧列表添加。</p>
      <ul v-else class="selected-list">
        <li v-for="(item, index) in selectedList" :key="item.ingredient_id" class="selected-item">
          <div
            class="swipe-content"
            :style="{ transform: `translateX(${swipeOffsets[item.ingredient_id] || 0}px)` }"
            @touchstart="handleSwipeStart($event, item.ingredient_id)"
            @touchmove="handleSwipeMove($event)"
            @touchend="handleSwipeEnd"
            @touchcancel="handleSwipeEnd"
            @mousedown="handleSwipeStart($event, item.ingredient_id)"
            @mousemove="handleSwipeMove($event)"
            @mouseup="handleSwipeEnd"
            @mouseleave="handleSwipeEnd"
            @click.stop="handleSwipeTap(item.ingredient_id, $event)"
          >
            <div class="selected-meta">
              <strong>{{ item.ingredient_name || item.name }}</strong>
              <small>{{ item.category || '未分类' }}</small>
            </div>
            <div class="selected-inputs">
              <label>
                <span>数量</span>
                <input
                  type="number"
                  min="0"
                  step="0.1"
                  class="form-control"
                  v-model.number="item.amount"
                  @change="emitChange()"
                />
              </label>
              <label>
                <span>单位</span>
                <input
                  type="text"
                  maxlength="10"
                  class="form-control"
                  v-model="item.unit"
                  @change="emitChange()"
                />
              </label>
              <label class="note-input">
                <span>备注</span>
                <input
                  type="text"
                  maxlength="50"
                  class="form-control"
                  v-model="item.notes"
                  @change="emitChange()"
                  placeholder="例如：切片/去籽"
                />
              </label>
            </div>
          </div>
          <button type="button" class="swipe-delete" @click="removeIngredient(index)">
            移除
          </button>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { searchIngredients, fetchIngredientsByCategory } from '@/api/ingredients'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue'])

const clone = (value) => {
  try {
    return JSON.parse(JSON.stringify(value || []))
  } catch (error) {
    return Array.isArray(value) ? [...value] : []
  }
}

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

const categoryOptions = [
  { label: '肉类', value: 'meat' },
  { label: '蔬菜', value: 'vegetable' },
  { label: '海鲜', value: 'seafood' },
  { label: '蛋白', value: 'protein' },
  { label: '调味', value: 'condiment' },
  { label: '谷物', value: 'grain' },
  { label: '其他', value: 'other' }
]

const totalCategoryPages = computed(() => {
  return Math.max(1, Math.ceil(categoryPaging.total / categoryPaging.pageSize))
})

watch(
  () => props.modelValue,
  (val) => {
    selectedList.value = clone(val)
  },
  { deep: true, immediate: true }
)

watch(searchKeyword, (val) => {
  debouncedSearch(val)
})

watch(
  () => [categoryQuery.category, categoryQuery.keyword, categoryPaging.page],
  () => {
    loadCategory()
  }
)

const emitChange = () => {
  emit('update:modelValue', clone(selectedList.value))
}

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

const changeCategoryPage = (delta) => {
  const next = categoryPaging.page + delta
  if (next < 1) return
  if (next > totalCategoryPages.value) return
  categoryPaging.page = next
}

const addIngredient = (item) => {
  if (!item?.ingredient_id) return
  const exists = selectedList.value.some((target) => target.ingredient_id === item.ingredient_id)
  if (exists) return
  selectedList.value.push({
    ingredient_id: item.ingredient_id,
    ingredient_name: item.name,
    category: item.category,
    amount: 1,
    unit: item.default_unit || '份',
    notes: ''
  })
  emitChange()
}

const ACTION_WIDTH = 88
const swipeOffsets = reactive({})
const openItemId = ref(null)
const swipeState = reactive({
  startX: 0,
  activeId: null,
  initialOffset: 0,
  isPointerDown: false
})

const clamp = (value, min, max) => {
  return Math.min(max, Math.max(min, value))
}

const getPointX = (event) => {
  if (event?.touches?.length) return event.touches[0].clientX
  if (event?.changedTouches?.length) return event.changedTouches[0].clientX
  return event?.clientX || 0
}

const isInteractiveTarget = (event) => {
  const target = event?.target
  if (!target) return false
  return Boolean(target.closest('input, textarea, select, button'))
}

const closeSwipe = (id) => {
  if (id === undefined || id === null) return
  swipeOffsets[id] = 0
  if (openItemId.value === id) {
    openItemId.value = null
  }
}

const handleSwipeStart = (event, id) => {
  if (id === undefined || id === null) return
  if (isInteractiveTarget(event)) return
  swipeState.startX = getPointX(event)
  swipeState.initialOffset = swipeOffsets[id] || 0
  swipeState.activeId = id
  swipeState.isPointerDown = true
  if (
    openItemId.value !== null &&
    openItemId.value !== undefined &&
    openItemId.value !== id
  ) {
    closeSwipe(openItemId.value)
  }
}

const handleSwipeMove = (event) => {
  if (!swipeState.isPointerDown || swipeState.activeId === null || swipeState.activeId === undefined)
    return
  if (event?.cancelable) {
    event.preventDefault()
  }
  const currentX = getPointX(event)
  const delta = currentX - swipeState.startX
  const offset = clamp(swipeState.initialOffset + delta, -ACTION_WIDTH, 0)
  swipeOffsets[swipeState.activeId] = offset
}

const handleSwipeEnd = () => {
  if (swipeState.activeId === null || swipeState.activeId === undefined) return
  const offset = swipeOffsets[swipeState.activeId] || 0
  if (offset < -ACTION_WIDTH / 2) {
    swipeOffsets[swipeState.activeId] = -ACTION_WIDTH
    openItemId.value = swipeState.activeId
  } else {
    closeSwipe(swipeState.activeId)
  }
  swipeState.activeId = null
  swipeState.isPointerDown = false
  swipeState.startX = 0
  swipeState.initialOffset = 0
}

const handleSwipeTap = (id, event) => {
  if (id === undefined || id === null) return
  if (isInteractiveTarget(event)) return
  if (openItemId.value !== null && openItemId.value !== undefined && openItemId.value === id) {
    closeSwipe(id)
  } else if (
    openItemId.value !== null &&
    openItemId.value !== undefined &&
    openItemId.value !== id
  ) {
    closeSwipe(openItemId.value)
  }
}

const removeIngredient = (index) => {
  const [removed] = selectedList.value.splice(index, 1)
  if (removed?.ingredient_id !== undefined) {
    delete swipeOffsets[removed.ingredient_id]
    if (openItemId.value === removed.ingredient_id) {
      openItemId.value = null
    }
  }
  emitChange()
}

watch(
  selectedList,
  (list) => {
    const ids = new Set(list.map((item) => String(item.ingredient_id)))
    Object.keys(swipeOffsets).forEach((id) => {
      if (!ids.has(id)) {
        delete swipeOffsets[id]
      }
    })
    const openKey = openItemId.value !== null && openItemId.value !== undefined ? String(openItemId.value) : null
    if (openKey && !ids.has(openKey)) {
      openItemId.value = null
    }
  },
  { deep: true }
)

loadCategory()
</script>

<style scoped>
.ingredient-selector {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  margin-top: 16px;
}

.selector-column {
  background: var(--color-surface);
  border-radius: var(--radius-large);
  padding: 16px;
  box-shadow: var(--shadow-card);
  display: flex;
  flex-direction: column;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.category-filters {
  display: flex;
  gap: 8px;
}

.category-filters select,
.category-filters input {
  flex: 1;
}

.result-list {
  margin-top: 12px;
  list-style: none;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.result-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-medium);
}

.result-item small {
  display: block;
  color: var(--color-text-secondary);
}

.empty {
  color: var(--color-text-secondary);
  font-size: 14px;
}

.selected-column {
  grid-column: span 1;
}

.selected-list {
  list-style: none;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.selected-item {
  position: relative;
  border-radius: var(--radius-medium);
  overflow: hidden;
}

.swipe-content {
  position: relative;
  width: 100%;
  border: 1px dashed var(--color-border);
  border-radius: var(--radius-medium);
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  background: var(--color-surface);
  transition: transform 0.2s ease;
  z-index: 1;
}

.swipe-delete {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 88px;
  border: none;
  background: var(--color-danger);
  color: #fff;
  font-weight: 600;
  letter-spacing: 2px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 0;
}

.selected-inputs {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 8px;
  align-items: end;
}

.note-input {
  grid-column: 1 / -1;
}

.selected-inputs label {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 13px;
}
.pagination {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 12px;
  justify-content: center;
}
</style>
