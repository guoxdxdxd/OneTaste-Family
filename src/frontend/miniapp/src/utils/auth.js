/**
 * 认证工具
 * Token管理和认证相关函数
 */

import { tokenStorage, userInfoStorage } from './storage'

/**
 * 设置认证信息
 * @param {string} token Token
 * @param {object} userInfo 用户信息
 */
export function setAuth(token, userInfo) {
  tokenStorage.setToken(token)
  if (userInfo) {
    userInfoStorage.setUserInfo(userInfo)
  }
}

/**
 * 清除认证信息
 */
export function clearAuth() {
  tokenStorage.removeToken()
  userInfoStorage.removeUserInfo()
}

/**
 * 获取 Token
 * @returns {string|null}
 */
export function getToken() {
  return tokenStorage.getToken()
}

/**
 * 获取用户信息
 * @returns {object|null}
 */
export function getUserInfo() {
  return userInfoStorage.getUserInfo()
}

/**
 * 检查是否已登录
 * @returns {boolean}
 */
export function isAuthenticated() {
  const token = getToken()
  return !!token
}

