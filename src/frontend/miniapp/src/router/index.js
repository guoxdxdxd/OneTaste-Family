/**
 * 路由主文件
 * 定义所有路由规则
 */

import { createRouter, createWebHistory } from 'vue-router'
import { beforeEach, afterEach } from './guards'

// 路由配置
const routes = [
  {
    path: '/',
    component: () => import('@/pages/layouts/AppLayout.vue'),
    meta: {
      requiresAuth: true
    },
    children: [
      {
        path: '',
        redirect: '/menus'
      },
      {
        path: 'recipes',
        name: 'Recipes',
        component: () => import('@/pages/recipes/RecipeManagement.vue'),
        meta: {
          requiresAuth: true,
          title: '菜谱管理'
        }
      },
      {
        path: 'menus',
        name: 'Menus',
        component: () => import('@/pages/menu/MenuHub.vue'),
        meta: {
          requiresAuth: true,
          title: '菜单中心'
        }
      },
      {
        path: 'shopping',
        name: 'Shopping',
        component: () => import('@/pages/shopping/ShoppingHub.vue'),
        meta: {
          requiresAuth: true,
          title: '买菜清单'
        }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/pages/profile/ProfileHome.vue'),
        meta: {
          requiresAuth: true,
          title: '我的'
        }
      }
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/auth/Login.vue'),
    meta: {
      requiresAuth: false,
      title: '登录'
    }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/pages/auth/Register.vue'),
    meta: {
      requiresAuth: false,
      title: '注册'
    }
  },
  {
    path: '/invite',
    name: 'InviteLanding',
    component: () => import('@/pages/invite/InviteLanding.vue'),
    meta: {
      requiresAuth: false,
      title: '家庭邀请'
    }
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach(beforeEach)
router.afterEach(afterEach)

export default router
