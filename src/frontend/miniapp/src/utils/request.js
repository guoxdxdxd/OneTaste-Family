/**
 * HTTP请求封装
 * axios封装和拦截器
 */

import axios from 'axios'
import { getToken, clearAuth } from './auth'

function createRequestError(message, code, data) {
  const error = new Error(message || '请求失败')
  error.code = code
  error.data = data
  return error
}

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
    if (res?.code === 200) {
      return res
    }

    if (res?.code === 401) {
      clearAuth()
      if (typeof window !== 'undefined') {
        window.location.href = '/login'
      }
    }

    return Promise.reject(createRequestError(res?.message, res?.code, res))
  },
  (error) => {
    let message = '网络错误，请稍后重试'
    let code = error?.response?.status

    if (error.response) {
      const data = error.response.data || {}
      message = data.message || message

      if (code === 401) {
        clearAuth()
        if (typeof window !== 'undefined') {
          window.location.href = '/login'
        }
      }

      return Promise.reject(createRequestError(message, data.code || code, data))
    }

    return Promise.reject(createRequestError(message, code))
  }
)

export default request
