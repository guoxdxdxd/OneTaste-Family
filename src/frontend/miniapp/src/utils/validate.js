/**
 * 验证工具
 * 表单验证和数据验证函数
 */

/**
 * 验证手机号
 * @param {string} phone 手机号
 * @returns {boolean}
 */
export function validatePhone(phone) {
  const phoneReg = /^1[3-9]\d{9}$/
  return phoneReg.test(phone)
}

/**
 * 验证密码
 * @param {string} password 密码
 * @param {number} minLength 最小长度，默认6
 * @param {number} maxLength 最大长度，默认20
 * @returns {boolean}
 */
export function validatePassword(password, minLength = 6, maxLength = 20) {
  if (!password) return false
  if (password.length < minLength || password.length > maxLength) return false
  return true
}

/**
 * 验证图形验证码
 * @param {string} code 验证码
 * @param {number} length 验证码长度，默认4
 * @returns {boolean}
 */
export function validateVerifyCode(code, length = 4) {
  if (!code) return false
  const codeReg = new RegExp(`^[A-Za-z0-9]{${length}}$`)
  return codeReg.test(code)
}

/**
 * 验证昵称
 * @param {string} nickname 昵称
 * @param {number} minLength 最小长度，默认1
 * @param {number} maxLength 最大长度，默认20
 * @returns {boolean}
 */
export function validateNickname(nickname, minLength = 1, maxLength = 20) {
  if (!nickname) return false
  if (nickname.length < minLength || nickname.length > maxLength) return false
  return true
}

/**
 * 获取手机号验证错误信息
 * @param {string} phone 手机号
 * @returns {string|null}
 */
export function getPhoneError(phone) {
  if (!phone) {
    return '请输入手机号'
  }
  if (!validatePhone(phone)) {
    return '手机号格式不正确'
  }
  return null
}

/**
 * 获取密码验证错误信息
 * @param {string} password 密码
 * @returns {string|null}
 */
export function getPasswordError(password) {
  if (!password) {
    return '请输入密码'
  }
  if (password.length < 6) {
    return '密码长度不能少于6位'
  }
  if (password.length > 20) {
    return '密码长度不能超过20位'
  }
  return null
}

/**
 * 获取验证码验证错误信息
 * @param {string} code 验证码
 * @returns {string|null}
 */
export function getVerifyCodeError(code) {
  if (!code) {
    return '请输入图形验证码'
  }
  if (!validateVerifyCode(code)) {
    return '验证码格式不正确，请输入4位数字或字母'
  }
  return null
}

/**
 * 获取昵称验证错误信息
 * @param {string} nickname 昵称
 * @returns {string|null}
 */
export function getNicknameError(nickname) {
  if (!nickname) {
    return '请输入昵称'
  }
  if (nickname.length < 1) {
    return '昵称不能为空'
  }
  if (nickname.length > 20) {
    return '昵称长度不能超过20位'
  }
  return null
}
