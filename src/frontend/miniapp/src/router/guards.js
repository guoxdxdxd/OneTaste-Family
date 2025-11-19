/**
 * 路由守卫
 * 处理路由权限和登录验证
 */

import { useUserStore } from '@/stores/user'

/**
 * 检查是否已登录
 * @returns {boolean}
 */
export function isAuthenticated() {
  const userStore = useUserStore()
  return userStore.isLoggedIn
}

/**
 * 路由前置守卫
 * @param {object} to 目标路由
 * @param {object} from 来源路由
 * @param {Function} next 下一步函数
 */
export function beforeEach(to, from, next) {
  // 检查路由是否需要登录
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  if (requiresAuth && !isAuthenticated()) {
    // 需要登录但未登录，跳转到登录页
    next({
      path: '/login',
      query: { redirect: to.fullPath }
    })
  } else {
    next()
  }
}

/**
 * 路由后置守卫
 * @param {object} to 目标路由
 * @param {object} from 来源路由
 */
export function afterEach(to, from) {
  // 可以在这里添加页面访问统计等逻辑
  // console.log('Route changed:', to.path)
}

