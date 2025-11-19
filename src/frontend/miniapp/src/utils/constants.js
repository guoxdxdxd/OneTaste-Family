/**
 * 常量定义
 * 项目常量配置
 */

// API 基础URL
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

// Token 存储键名
export const TOKEN_KEY = 'token'

// 用户信息存储键名
export const USER_INFO_KEY = 'user_info'

// 餐次类型
export const MEAL_TYPES = {
  BREAKFAST: 'breakfast',
  LUNCH: 'lunch',
  DINNER: 'dinner'
}

// 餐次类型中文
export const MEAL_TYPE_NAMES = {
  [MEAL_TYPES.BREAKFAST]: '早餐',
  [MEAL_TYPES.LUNCH]: '午餐',
  [MEAL_TYPES.DINNER]: '晚餐'
}

// 会员类型
export const MEMBERSHIP_TYPES = {
  FREE: 'free',
  PREMIUM: 'premium'
}

// 会员类型中文
export const MEMBERSHIP_TYPE_NAMES = {
  [MEMBERSHIP_TYPES.FREE]: '免费版',
  [MEMBERSHIP_TYPES.PREMIUM]: '付费版'
}

// 分页默认值
export const PAGINATION = {
  DEFAULT_PAGE: 1,
  DEFAULT_PAGE_SIZE: 20
}

