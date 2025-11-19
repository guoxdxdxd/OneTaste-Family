/**
 * HTTP请求封装
 * axios封装和拦截器
 */

import axios from 'axios'
import { getToken, clearAuth } from './auth'

// 创建 axios 实例
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    // 添加 token
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const res = response.data

    // 如果返回的状态码为200，说明接口请求成功
    if (res.code === 200) {
      return res
    }

    // 处理业务错误
    if (res.code === 401) {
      // 未授权，清除认证信息并跳转到登录页
      clearAuth()
      // 这里需要根据实际的路由方式来处理跳转
      // 如果是小程序，使用 wx.reLaunch 或 uni.reLaunch
      // 如果是 Web，使用 router.push
      if (typeof window !== 'undefined') {
        window.location.href = '/login'
      }
      return Promise.reject(new Error(res.message || '未授权'))
    }

    // 其他错误
    return Promise.reject(new Error(res.message || '请求失败'))
  },
  (error) => {
    // 处理 HTTP 错误
    let message = '网络错误，请稍后重试'
    
    if (error.response) {
      // 服务器返回了错误状态码
      const status = error.response.status
      switch (status) {
        case 400:
          message = '请求参数错误'
          break
        case 401:
          message = '未授权，请重新登录'
          clearAuth()
          if (typeof window !== 'undefined') {
            window.location.href = '/login'
          }
          break
        case 403:
          message = '无权限访问'
          break
        case 404:
          message = '请求的资源不存在'
          break
        case 500:
          message = '服务器错误'
          break
        default:
          message = error.response.data?.message || `请求失败 (${status})`
      }
    } else if (error.request) {
      // 请求已发出，但没有收到响应
      message = '网络连接失败，请检查网络'
    }

    return Promise.reject(new Error(message))
  }
)

export default request

