/**
 * 家庭管理 API
 * 覆盖创建家庭、家庭信息、成员与邀请
 */

import request from './index'

/**
 * 创建家庭
 * @param {object} data
 * @param {string} data.name 家庭名称
 * @param {string} [data.description] 家庭描述
 */
export function createFamily(data) {
  return request.post('/family/create', data)
}

/**
 * 获取家庭信息
 */
export function getFamilyInfo() {
  return request.get('/family/info')
}

/**
 * 获取家庭成员列表
 */
export function getFamilyMembers() {
  return request.get('/family/members')
}

/**
 * 接受家庭邀请
 * @param {object} data
 * action 在当前阶段默认 accept
 */
export function acceptInvite(data) {
  return request.post('/family/member/invite', data)
}
