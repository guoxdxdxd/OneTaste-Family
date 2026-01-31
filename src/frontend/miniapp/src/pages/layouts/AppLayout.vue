<template>
  <div class="app-layout" :class="{ 'app-layout--with-tabbar': showTabBar }">
    <!-- 主内容区域 -->
    <main class="app-layout__main">
      <router-view v-slot="{ Component }">
        <transition name="page-fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
    
    <!-- 底部导航栏 -->
    <TabBar />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import TabBar from '@/components/TabBar.vue'

const route = useRoute()

// 显示 TabBar 的主页面
const showTabBar = computed(() => {
  const mainPages = ['/menus', '/shopping', '/profile', '/recipes']
  return mainPages.includes(route.path)
})
</script>

<style scoped>
.app-layout {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg-base);
}

.app-layout__main {
  flex: 1;
  display: flex;
  flex-direction: column;
}

/* 有底部导航栏时，为其留出空间 */
.app-layout--with-tabbar .app-layout__main {
  padding-bottom: calc(var(--tabbar-height) + var(--safe-area-bottom));
}

/* 页面切换动画 */
.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity var(--duration-normal) var(--ease-out);
}

.page-fade-enter-from,
.page-fade-leave-to {
  opacity: 0;
}
</style>
