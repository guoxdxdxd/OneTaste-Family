/**
 * 用户状态管理
 * 管理用户信息和登录状态
 */

import { defineStore } from 'pinia'
import { setAuth, clearAuth, getToken, getUserInfo } from '@/utils/auth'
import { userInfoStorage } from '@/utils/storage'
import { register as registerApi, login as loginApi, getUserInfo as getUserInfoApi } from '@/api/user'
import { useFamilyStore } from './family'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: getToken(),
    userInfo: getUserInfo(),
    isLoggedIn: !!getToken()
  }),

  getters: {
    /**
     * 是否已登录
     */
    loggedIn: (state) => state.isLoggedIn,

    /**
     * 用户ID
     */
    userId: (state) => state.userInfo?.user_id || null,

    /**
     * 用户昵称
     */
    nickname: (state) => state.userInfo?.nickname || '',

    /**
     * 用户手机号
     */
    phone: (state) => state.userInfo?.phone || '',

    /**
     * 会员类型
     */
    membershipType: (state) => state.userInfo?.membership?.type || 'free'
  },

  actions: {
    /**
     * 用户注册
     * @param {object} data 注册数据
     * @returns {Promise}
     */
    async register(data) {
      try {
        const res = await registerApi(data)
        if (res.code === 200 && res.data) {
          // 保存 token 和用户信息
          this.setAuth(res.data.token, {
            user_id: res.data.user_id,
            phone: data.phone,
            nickname: data.nickname
          })
          return res
        }
        throw new Error(res.message || '注册失败')
      } catch (error) {
        console.error('Register error:', error)
        throw error
      }
    },

    /**
     * 用户登录
     * @param {object} data 登录数据
     * @returns {Promise}
     */
    async login(data) {
      try {
        const res = await loginApi(data)
        if (res.code === 200 && res.data) {
          // 保存 token
          this.setAuth(res.data.token, {
            user_id: res.data.user_id,
            phone: data.phone
          })
          // 获取完整用户信息
          await this.fetchUserInfo()
          return res
        }
        throw new Error(res.message || '登录失败')
      } catch (error) {
        console.error('Login error:', error)
        throw error
      }
    },

    /**
     * 获取用户信息
     * @returns {Promise}
     */
    async fetchUserInfo() {
      try {
        const res = await getUserInfoApi()
        if (res.code === 200 && res.data) {
          this.userInfo = res.data
          // 更新本地存储
          userInfoStorage.setUserInfo(res.data)
          return res.data
        }
        throw new Error(res.message || '获取用户信息失败')
      } catch (error) {
        console.error('Fetch user info error:', error)
        throw error
      }
    },

    /**
     * 设置认证信息
     * @param {string} token Token
     * @param {object} userInfo 用户信息
     */
    setAuth(token, userInfo) {
      this.token = token
      this.userInfo = userInfo
      this.isLoggedIn = true
      setAuth(token, userInfo)
    },

    /**
     * 退出登录
     */
    logout() {
      this.token = null
      this.userInfo = null
      this.isLoggedIn = false
      clearAuth()
      const familyStore = useFamilyStore()
      familyStore.reset()
    }
  }
})
