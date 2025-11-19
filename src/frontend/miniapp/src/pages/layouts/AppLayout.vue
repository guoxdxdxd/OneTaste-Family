<template>
  <div class="app-shell">
    <main class="app-content">
      <router-view />
    </main>
    <nav class="tab-bar">
      <button
        v-for="tab in tabs"
        :key="tab.path"
        type="button"
        class="tab-item"
        :class="{ active: isActive(tab.path) }"
        @click="go(tab.path)"
      >
        <span class="tab-label">{{ tab.label }}</span>
        <small>{{ tab.caption }}</small>
      </button>
    </nav>
  </div>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const route = useRoute()

const tabs = [
  { path: '/menus', label: '菜单', caption: '计划' },
  { path: '/shopping', label: '买菜', caption: '清单' },
  { path: '/profile', label: '我的', caption: '家庭' }
]

const isActive = (path) => {
  return route.path.startsWith(path)
}

const go = (path) => {
  if (route.path !== path) {
    router.push(path)
  }
}
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
  background: var(--color-background);
  display: flex;
  flex-direction: column;
}

.app-content {
  flex: 1;
  padding-bottom: 80px;
}

.tab-bar {
  position: sticky;
  bottom: 0;
  left: 0;
  right: 0;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  border-top: 1px solid var(--color-border);
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(16px);
  padding: 4px;
  column-gap: 8px;
}

.tab-item {
  border: none;
  background: transparent;
  padding: 10px 8px;
  border-radius: var(--radius-large);
  display: flex;
  flex-direction: column;
  gap: 2px;
  font-size: 15px;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab-item small {
  font-size: 12px;
  letter-spacing: 0.3em;
}

.tab-item.active {
  background: var(--color-surface);
  color: var(--color-text-primary);
  box-shadow: var(--shadow-card);
}
</style>
