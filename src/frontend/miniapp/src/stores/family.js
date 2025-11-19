/**
 * 家庭状态管理
 * 负责家庭信息与成员数据
 */

import { defineStore } from 'pinia'
import { createFamily, getFamilyInfo, getFamilyMembers, acceptInvite } from '@/api/family'

export const useFamilyStore = defineStore('family', {
  state: () => ({
    familyInfo: null,
    members: [],
    infoLoading: false,
    membersLoading: false,
    createLoading: false,
    lastLoadedAt: null,
    membersLoadedAt: null,
    error: ''
  }),

  getters: {
    hasFamily: (state) => !!state.familyInfo?.family_id,
    memberCount: (state) => state.members.length || state.familyInfo?.member_count || 0,
    role: (state) => state.familyInfo?.role || 'member',
    familyName: (state) => state.familyInfo?.name || ''
  },

  actions: {
    async fetchFamilyInfo(force = false) {
      if (this.infoLoading) return
      if (!force && this.familyInfo && Date.now() - this.lastLoadedAt < 60_000) {
        return this.familyInfo
      }
      this.infoLoading = true
      this.error = ''
      try {
        const res = await getFamilyInfo()
        if (res.data) {
          this.familyInfo = res.data
          this.lastLoadedAt = Date.now()
        }
        return res.data
      } catch (error) {
        this.error = error.message || '获取家庭信息失败'
        throw error
      } finally {
        this.infoLoading = false
      }
    },

    async fetchMembers(force = false) {
      if (this.membersLoading) return
      if (!this.hasFamily) {
        this.members = []
        return []
      }
      if (!force && this.members.length && Date.now() - this.membersLoadedAt < 60_000) {
        return this.members
      }
      this.membersLoading = true
      this.error = ''
      try {
        const res = await getFamilyMembers()
        if (res.data?.members) {
          this.members = res.data.members
          this.membersLoadedAt = Date.now()
        } else {
          this.members = []
        }
        return this.members
      } catch (error) {
        this.error = error.message || '获取成员列表失败'
        throw error
      } finally {
        this.membersLoading = false
      }
    },

    async createFamily(payload) {
      if (this.createLoading) return
      this.createLoading = true
      this.error = ''
      try {
        const res = await createFamily(payload)
        if (res.data) {
          this.familyInfo = res.data
          this.lastLoadedAt = Date.now()
        }
        return res
      } catch (error) {
        this.error = error.message || '创建家庭失败'
        throw error
      } finally {
        this.createLoading = false
      }
    },

    async acceptInvite(params) {
      this.error = ''
      try {
        const res = await acceptInvite({ ...params, action: 'accept' })
        if (res.data) {
          await this.fetchFamilyInfo(true)
          await this.fetchMembers(true)
        }
        return res
      } catch (error) {
        this.error = error.message || '加入家庭失败'
        throw error
      }
    },

    reset() {
      this.familyInfo = null
      this.members = []
      this.infoLoading = false
      this.membersLoading = false
      this.createLoading = false
      this.error = ''
      this.lastLoadedAt = null
      this.membersLoadedAt = null
    }
  }
})
