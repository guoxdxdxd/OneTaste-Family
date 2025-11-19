/**
 * 本地存储工具
 * 封装 localStorage 和 sessionStorage
 */

const TOKEN_KEY = 'token'
const USER_INFO_KEY = 'user_info'

/**
 * localStorage 操作
 */
export const storage = {
  /**
   * 设置值
   * @param {string} key 键名
   * @param {any} value 值
   */
  set(key, value) {
    try {
      const strValue = JSON.stringify(value)
      localStorage.setItem(key, strValue)
    } catch (error) {
      console.error('Storage set error:', error)
    }
  },

  /**
   * 获取值
   * @param {string} key 键名
   * @param {any} defaultValue 默认值
   * @returns {any}
   */
  get(key, defaultValue = null) {
    try {
      const value = localStorage.getItem(key)
      return value ? JSON.parse(value) : defaultValue
    } catch (error) {
      console.error('Storage get error:', error)
      return defaultValue
    }
  },

  /**
   * 删除值
   * @param {string} key 键名
   */
  remove(key) {
    try {
      localStorage.removeItem(key)
    } catch (error) {
      console.error('Storage remove error:', error)
    }
  },

  /**
   * 清空所有
   */
  clear() {
    try {
      localStorage.clear()
    } catch (error) {
      console.error('Storage clear error:', error)
    }
  }
}

/**
 * Token 相关操作
 */
export const tokenStorage = {
  /**
   * 设置 token
   * @param {string} token
   */
  setToken(token) {
    storage.set(TOKEN_KEY, token)
  },

  /**
   * 获取 token
   * @returns {string|null}
   */
  getToken() {
    return storage.get(TOKEN_KEY)
  },

  /**
   * 删除 token
   */
  removeToken() {
    storage.remove(TOKEN_KEY)
  }
}

/**
 * 用户信息相关操作
 */
export const userInfoStorage = {
  /**
   * 设置用户信息
   * @param {object} userInfo
   */
  setUserInfo(userInfo) {
    storage.set(USER_INFO_KEY, userInfo)
  },

  /**
   * 获取用户信息
   * @returns {object|null}
   */
  getUserInfo() {
    return storage.get(USER_INFO_KEY)
  },

  /**
   * 删除用户信息
   */
  removeUserInfo() {
    storage.remove(USER_INFO_KEY)
  }
}

