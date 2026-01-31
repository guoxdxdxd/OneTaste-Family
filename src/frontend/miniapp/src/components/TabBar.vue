<template>
  <!-- 仅在主页面显示 TabBar，子页面（如创建/编辑）不显示 -->
  <nav v-if="showTabBar" class="tabbar">
    <div class="tabbar__inner">
      <button
        v-for="item in tabs"
        :key="item.path"
        type="button"
        class="tabbar__item"
        :class="{ 'tabbar__item--active': isActive(item.path) }"
        @click="navigate(item.path)"
      >
        <div class="tabbar__icon-wrapper">
          <component :is="item.icon" class="tabbar__icon" />
          <span v-if="item.badge" class="tabbar__badge">{{ item.badge }}</span>
        </div>
        <span class="tabbar__label">{{ item.label }}</span>
      </button>
    </div>
  </nav>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import IconMenu from './icons/IconMenu.vue'
import IconShopping from './icons/IconShopping.vue'
import IconUser from './icons/IconUser.vue'

const router = useRouter()
const route = useRoute()

// 需要隐藏 TabBar 的路由（子页面）
const hiddenRoutes = [
  '/menus/create',
  '/menus/daily',
  '/menus/weekly',
  '/menus/',  // 匹配 /menus/:id/edit
  '/recipes'  // 菜谱管理页面有 FAB，暂时保留 TabBar
]

// 是否显示 TabBar
const showTabBar = computed(() => {
  const path = route.path
  // 精确匹配主页面
  const mainPages = ['/menus', '/shopping', '/profile']
  if (mainPages.includes(path)) return true
  // 菜谱页面显示 TabBar
  if (path === '/recipes') return true
  // 其他子页面不显示
  return false
})

const tabs = computed(() => [
  { 
    path: '/menus', 
    label: '菜单',
    icon: IconMenu,
    badge: null
  },
  { 
    path: '/shopping', 
    label: '买菜',
    icon: IconShopping,
    badge: null
  },
  { 
    path: '/profile', 
    label: '我的',
    icon: IconUser,
    badge: null
  }
])

const isActive = (path) => {
  return route.path.startsWith(path)
}

const navigate = (path) => {
  if (route.path !== path) {
    router.push(path)
  }
}
</script>

<style scoped>
.tabbar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: var(--z-fixed);
  padding-bottom: var(--safe-area-bottom);
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-top: 1px solid var(--color-border-light);
}

.tabbar__inner {
  display: flex;
  justify-content: space-around;
  align-items: center;
  max-width: var(--container-max-width);
  margin: 0 auto;
  height: calc(var(--tabbar-height) - 1px);
  padding: 0 var(--space-2);
}

.tabbar__item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--space-1);
  padding: var(--space-2) var(--space-1);
  background: transparent;
  border: none;
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
  position: relative;
  -webkit-tap-highlight-color: transparent;
}

.tabbar__item:active {
  transform: scale(0.95);
}

.tabbar__icon-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 32px;
  border-radius: var(--radius-full);
  transition: all var(--duration-normal) var(--ease-spring);
}

.tabbar__item--active .tabbar__icon-wrapper {
  background: var(--color-primary-100);
}

.tabbar__icon {
  width: 24px;
  height: 24px;
  color: var(--color-text-tertiary);
  transition: color var(--duration-fast) var(--ease-out);
}

.tabbar__item--active .tabbar__icon {
  color: var(--color-primary);
}

.tabbar__label {
  font-size: 11px;
  font-weight: var(--font-weight-medium);
  color: var(--color-text-tertiary);
  transition: color var(--duration-fast) var(--ease-out);
  letter-spacing: 0.02em;
}

.tabbar__item--active .tabbar__label {
  color: var(--color-primary);
  font-weight: var(--font-weight-semibold);
}

.tabbar__badge {
  position: absolute;
  top: -2px;
  right: 4px;
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  background: var(--color-danger-500);
  color: white;
  font-size: 10px;
  font-weight: var(--font-weight-semibold);
  line-height: 16px;
  text-align: center;
  border-radius: var(--radius-full);
}
</style>
