/**
 * 用户API接口
 * 用户相关的API调用
 */

import request from './index'

/**
 * 用户注册
 * @param {object} data 注册数据
 * @param {string} data.phone 手机号
 * @param {string} data.password 密码
 * @param {string} data.verify_code 图形验证码内容
 * @param {string} data.captcha_key 图形验证码编码
 * @param {string} data.nickname 昵称
 * @returns {Promise}
 */
export function register(data) {
  return request.post('/auth/register', data)
}

/**
 * 用户登录
 * @param {object} data 登录数据
 * @param {string} data.phone 手机号
 * @param {string} data.password 密码
 * @returns {Promise}
 */
export function login(data) {
  return request.post('/auth/login', data)
}

/**
 * 获取用户信息
 * @returns {Promise}
 */
export function getUserInfo() {
  return request.get('/user/info')
}

/**
 * 获取图形验证码
 * @param {object} params 查询参数
 * @param {number} params.width 图片宽度
 * @param {number} params.height 图片高度
 * @returns {Promise}
 */
export function getCaptcha(params) {
  return request.get('/auth/captcha', { params })
}
